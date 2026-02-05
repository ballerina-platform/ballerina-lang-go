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

	var debugCtx *debugcommon.DebugContext
	var wg sync.WaitGroup
	flags := uint16(0)

	if runOpts.dumpTokens {
		flags |= debugcommon.DUMP_TOKENS
	}
	if runOpts.dumpST {
		flags |= debugcommon.DUMP_ST
	}
	if runOpts.traceRecovery {
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
	// Java: io.ballerina.projects.directory.ProjectLoader.loadProject
	result, err := directory.LoadProject(path)
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

	// Get default module for compilation
	project := result.Project()
	pkg := project.CurrentPackage()
	defaultModule := pkg.GetDefaultModule()

	// Compile the source
	fmt.Fprintln(os.Stderr, "Compiling source")
	// Print package identifier based on project type
	if project.Kind() == projects.ProjectKindSingleFile {
		// Single file project - show filename
		docId := defaultModule.DocumentIDs()[0]
		doc := defaultModule.Document(docId)
		fmt.Fprintf(os.Stderr, "\t%s\n", doc.Name())
	} else {
		// Build project - show org/packageName:version
		fmt.Fprintf(os.Stderr, "\t%s/%s:%s\n",
			pkg.PackageOrg().Value(),
			pkg.PackageName().Value(),
			pkg.PackageVersion().String())
	}

	cx := context.NewCompilerContext(semtypes.CreateTypeEnv())

	syntaxTree, err := parser.GetSyntaxTree(cx, debugCtx, fileName)
	if err != nil {
		printError(fmt.Errorf("compilation failed: %w", err), "", false)
		return fmt.Errorf("compilation failed: %w", err)
	}

	compilationUnit := ast.GetCompilationUnit(cx, syntaxTree)
	if runOpts.dumpAST {
		prettyPrinter := ast.PrettyPrinter{}
		fmt.Println(prettyPrinter.Print(compilationUnit))
	}
	pkg := ast.ToPackage(compilationUnit)
	// Resolve symbols (imports) before type resolution
	importedSymbols := semantics.ResolveImports(cx, pkg, semantics.GetImplicitImports(cx))
	semantics.ResolveSymbols(cx, pkg, importedSymbols)
	// Add type resolution step
	typeResolver := semantics.NewTypeResolver(cx, importedSymbols)
	typeResolver.ResolveTypes(cx, pkg)
	// Run control flow analysis after type resolution
	/// We need this before semantic analysis since we need to do conditional type narrowing before semantic analysis
	cfg := semantics.CreateControlFlowGraph(cx, pkg)
	if runOpts.dumpCFG {
		// Print the CFG with separators
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "==================BEGIN CFG==================")

		if runOpts.format == "dot" {
			// Use DOT exporter
			dotExporter := semantics.NewCFGDotExporter(cx)
			fmt.Println(strings.TrimSpace(dotExporter.Export(cfg)))
		} else {
			// Use default S-expression printer
			prettyPrinter := semantics.NewCFGPrettyPrinter(cx)
			fmt.Println(strings.TrimSpace(prettyPrinter.Print(cfg)))
		}

		fmt.Fprintln(os.Stderr, "===================END CFG===================")
	}
	// Run semantic analysis after type resolution
	semanticAnalyzer := semantics.NewSemanticAnalyzer(cx)
	semanticAnalyzer.Analyze(pkg)
	// Run CFG analyses (reachability and explicit return) concurrently
	semantics.AnalyzeCFG(cx, pkg, cfg)
	birPkg := bir.GenBir(cx, pkg)
	if runOpts.dumpBIR {
		prettyPrinter := bir.PrettyPrinter{}
		// Print the BIR with separators
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "==================BEGIN BIR==================")
		fmt.Println(strings.TrimSpace(prettyPrinter.Print(*birPkg)))
		fmt.Fprintln(os.Stderr, "===================END BIR===================")
	}

	// Run the executable
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Running executable")
	fmt.Fprintln(os.Stderr)

	rt := runtime.NewRuntime()
	if err := rt.Interpret(*birPkg); err != nil {
		return err
	}
	return nil
}

