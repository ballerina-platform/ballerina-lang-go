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
	"path/filepath"
	"strings"

	"ballerina-lang-go/bir"
	debugcommon "ballerina-lang-go/common"
	_ "ballerina-lang-go/lib/rt"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/tools/diagnostics"

	"github.com/spf13/cobra"
)

var runOpts struct {
	dumpTokens    bool
	dumpST        bool
	dumpAST       bool
	dumpCFG       bool
	dumpBIR       bool
	traceRecovery bool
	stats         bool
	statsOneline  bool
	logFile       string
	format        string // Output format (dot, etc.)
}

var runCmd = &cobra.Command{
	Use:   "run [<source-file.bal> | <package-dir> | .]",
	Short: "Build and run the current package or a Ballerina source file",
	Long: `	Build the current package and run it.

	The 'run' command builds and executes the given Ballerina package or
	a source file.

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
	declaration.

	Note: Running individual '.bal' files of a package is not allowed.`,
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
	runCmd.Flags().BoolVar(&runOpts.stats, "stats", false, "Print per-stage compilation timing statistics")
	runCmd.Flags().BoolVar(&runOpts.statsOneline, "stats-oneline", false, "Print per-stage compilation timing totals only")
	runCmd.Flags().StringVar(&runOpts.logFile, "log-file", "", "Write debug output to specified file")
	runCmd.Flags().StringVar(&runOpts.format, "format", "", "Output format for dump operations (dot)")
	profiler.RegisterFlags(runCmd)
}

func runBallerina(cmd *cobra.Command, args []string) error {
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
		WithStats(runOpts.stats || runOpts.statsOneline).
		Build()

	if err := profiler.Start(); err != nil {
		profErr := fmt.Errorf("failed to start profiler: %w", err)
		printError(profErr, "", false)
		return profErr
	}
	defer func() { _ = profiler.Stop() }()

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
		var logWriter *os.File
		var err error
		if runOpts.logFile != "" {
			logWriter, err = os.Create(runOpts.logFile)
			if err != nil {
				cmdErr := fmt.Errorf("error creating log file %s: %w", runOpts.logFile, err)
				printError(cmdErr, "", false)
				return cmdErr
			}
			defer func() { _ = logWriter.Close() }()
			debugcommon.InitDebug(flags, logWriter)
		} else {
			debugcommon.InitDebug(flags, os.Stderr)
		}
	}

	// Default to current directory if no path provided (bal run == bal run .)
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	info, err := os.Stat(path)
	if err != nil {
		printRunError(err)
	}

	baseDir := path
	if !info.IsDir() {
		baseDir = filepath.Dir(path)
		path = filepath.Base(path)
	} else {
		path = "."
	}

	fsys := os.DirFS(baseDir)

	ballerinaHomePath, err := getBallerinaHomePath()
	if err != nil {
		printRunError(err)
		return err
	}
	ballerinaHomeFs := os.DirFS(ballerinaHomePath)

	result, err := projects.Load(fsys, path, projects.ProjectLoadConfig{
		BallerinaHomeFs: ballerinaHomeFs,
		BuildOptions:    &buildOpts,
	})
	if err != nil {
		printRunError(err)
		return err
	}

	// Check for loading errors
	diagResult := result.Diagnostics()
	if diagResult.HasErrors() {
		// Given we don't have sources at this point it is okay to pass an empty diagnostic env
		printDiagnostics(fsys, os.Stderr, diagResult, !isTerminal(), diagnostics.NewDiagnosticEnv())
		return fmt.Errorf("project loading contains errors")
	}

	project := result.Project()
	pkg := project.CurrentPackage()

	// Get package compilation (triggers parsing, type checking, semantic analysis, CFG analysis)
	compilation := pkg.Compilation()

	// Check for compilation errors
	compilationDiags := compilation.DiagnosticResult()
	if compilationDiags.HasErrors() {
		printDiagnostics(fsys, os.Stderr, compilationDiags, !isTerminal(), compilation.DiagnosticEnv())
		return fmt.Errorf("compilation failed with errors")
	}

	// Create backend and generate BIR
	backend := projects.NewBallerinaBackend(compilation)
	birPkgs := backend.BIRPackages()

	if len(birPkgs) == 0 {
		err := fmt.Errorf("BIR generation failed: no BIR package produced")
		printError(err, "", false)
		return err
	}

	if runOpts.statsOneline {
		fmt.Fprint(os.Stderr, compilation.StatsReportOneline())
	} else if buildOpts.Stats() {
		fmt.Fprint(os.Stderr, compilation.StatsReport())
	}

	// Dump BIR if requested
	if buildOpts.DumpBIR() {
		prettyPrinter := bir.PrettyPrinter{}
		for _, birPkg := range birPkgs {
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, "==================BEGIN BIR==================")
			fmt.Fprintln(os.Stderr, strings.TrimSpace(prettyPrinter.Print(*birPkg)))
			fmt.Fprintln(os.Stderr, "===================END BIR===================")
		}
	}

	rt := runtime.NewRuntime()
	for _, birPkg := range birPkgs {
		if err := rt.Interpret(*birPkg); err != nil {
			printRuntimeError(err)
			return err
		}
	}
	return nil
}

func printRunError(err error) {
	printError(err, "run [<source-file.bal> | <package-dir> | .]", false)
}

func getBallerinaHomePath() (string, error) {
	if balHome := os.Getenv(projects.BallerinaHomeEnvVar); balHome != "" {
		return balHome, nil
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userHome, projects.UserHomeDirName), nil
}
