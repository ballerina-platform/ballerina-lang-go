// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"ballerina-lang-go/bir"
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/directory"
	"ballerina-lang-go/runtime"

	"github.com/spf13/cobra"
)

var runOpts struct {
	dumpTokens    bool
	dumpST        bool
	dumpAST       bool
	dumpCFG       bool
	dumpBIR       bool
	traceRecovery bool
	logFile       string
	format        string // Output format (dot, etc.)
}

var runCmd = &cobra.Command{
	Use:   "run [<source-file.bal> | <package-dir> | .]",
	Short: "Compile and run the current package or a Ballerina source file",
	Long: `	Compile the current package and run it.

	The 'run' command compiles and executes the given Ballerina source file.

	A Ballerina program consists of one or more modules; one of these modules
	is distinguished as the root module, which is the default module of
	current package.

	Ballerina program execution consists of two consecutive phases.
	The initialization phase initializes all modules of a program one after
	another. If a module defines a function named 'init()', it will be
	invoked during this phase. If the root module of the program defines a
	public function named 'main()', then it will be invoked.

	If the initialization phase of program execution completes successfully,
	then execution proceeds to the listening phase. If there are no module
	listeners, then the listening phase immediately terminates successfully.
	Otherwise, the listening phase initializes the module listeners.

	A service declaration is the syntactic sugar for creating a service object
	and attaching it to the module listener specified in the service
	declaration.`,
	Args: validateSourceFile,
	RunE: runBallerina,
}

func init() {
	runCmd.Flags().BoolVar(&runOpts.dumpTokens, "dump-tokens", false, "Dump lexer tokens")
	runCmd.Flags().BoolVar(&runOpts.dumpST, "dump-st", false, "Dump syntax tree")
	runCmd.Flags().BoolVar(&runOpts.dumpAST, "dump-ast", false, "Dump abstract syntax tree")
	runCmd.Flags().BoolVar(&runOpts.dumpCFG, "dump-cfg", false, "Dump control flow graph")
	runCmd.Flags().BoolVar(&runOpts.dumpBIR, "dump-bir", false, "Dump Ballerina Intermediate Representation")
	runCmd.Flags().BoolVar(&runOpts.traceRecovery, "trace-recovery", false, "Enable error recovery tracing")
	runCmd.Flags().StringVar(&runOpts.logFile, "log-file", "", "Write debug output to specified file")
	runCmd.Flags().StringVar(&runOpts.format, "format", "", "Output format for dump operations (dot)")
}

func runBallerina(cmd *cobra.Command, args []string) error {
	// Default to current directory if no path provided (bal run == bal run .)
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// Build options from CLI flags. Constructed before debug setup so
	// buildOpts can be the single source of truth for all flag reads.
	buildOpts := projects.NewBuildOptionsBuilder().
		WithDumpAST(runOpts.dumpAST).
		WithDumpBIR(runOpts.dumpBIR).
		WithDumpCFG(runOpts.dumpCFG).
		WithDumpCFGFormat(projects.ParseCFGFormat(runOpts.format)).
		WithDumpTokens(runOpts.dumpTokens).
		WithDumpST(runOpts.dumpST).
		WithTraceRecovery(runOpts.traceRecovery).
		Build()

	var debugCtx *debugcommon.DebugContext
	var wg sync.WaitGroup
	flags := uint16(0)

	if buildOpts.DumpTokens() {
		flags |= debugcommon.DUMP_TOKENS
	}
	if buildOpts.DumpST() {
		flags |= debugcommon.DUMP_ST
	}
	if buildOpts.TraceRecovery() {
		flags |= debugcommon.DEBUG_ERROR_RECOVERY
	}

	if flags != 0 {
		debugcommon.Init(flags)
		debugCtx = &debugcommon.DebugCtx

		var logWriter *os.File
		var err error
		if runOpts.logFile != "" {
			logWriter, err = os.Create(runOpts.logFile)
			if err != nil {
				cmdErr := fmt.Errorf("error creating log file %s: %w", runOpts.logFile, err)
				printError(cmdErr, "", false)
				return cmdErr
			}
		} else {
			logWriter = os.Stderr
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if runOpts.logFile != "" {
				defer logWriter.Close()
			}
			for msg := range debugCtx.Channel {
				fmt.Fprintf(logWriter, "%s\n", msg)
			}
		}()

		// Ensure debug context cleanup on any exit path
		defer func() {
			if debugCtx != nil {
				close(debugCtx.Channel)
				wg.Wait()
			}
		}()
	}

	// Load project using ProjectLoader (auto-detects type)
	result, err := directory.LoadProject(path, directory.WithBuildOptions(buildOpts))
	if err != nil {
		printError(err, "run [<source-file.bal> | <package-dir> | .]", false)
		return err
	}

	// Check for loading errors
	diagResult := result.Diagnostics()
	if diagResult.HasErrors() {
		printDiagnostics(diagResult)
		return fmt.Errorf("project loading contains errors")
	}

	project := result.Project()
	pkg := project.CurrentPackage()

	// Get package compilation (triggers parsing, type checking, semantic analysis, CFG analysis)
	compilation := pkg.Compilation()

	// Check for compilation errors
	compilationDiags := compilation.DiagnosticResult()
	if compilationDiags.HasErrors() {
		printDiagnostics(compilationDiags)
		return fmt.Errorf("compilation failed with errors")
	}

	// Create backend and generate BIR
	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()

	if birPkg == nil {
		return fmt.Errorf("BIR generation failed: no BIR package produced")
	}

	// Dump BIR if requested
	if buildOpts.DumpBIR() {
		prettyPrinter := bir.PrettyPrinter{}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "==================BEGIN BIR==================")
		fmt.Println(strings.TrimSpace(prettyPrinter.Print(*birPkg)))
		fmt.Fprintln(os.Stderr, "===================END BIR===================")
	}

	rt := runtime.NewRuntime()
	if err := rt.Interpret(*birPkg); err != nil {
		return err
	}
	return nil
}
