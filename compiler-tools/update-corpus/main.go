/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	numJobs := flag.Int("j", 1, "Number of parallel jobs")
	jobsFlag := flag.Int("jobs", 1, "Number of parallel jobs (alternative to -j)")
	flag.Parse()

	// Determine which flag was used by checking command line args
	finalNumJobs := *numJobs
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "jobs" {
			finalNumJobs = *jobsFlag
		}
	})

	if finalNumJobs < 1 {
		fmt.Fprintf(os.Stderr, "Error: number of jobs must be at least 1\n")
		os.Exit(1)
	}

	corpusBalDir := "./corpus/bal"
	corpusTokensDir := "./corpus/tokens"
	corpusParserDir := "./corpus/parser"

	// Find ballerina-lang-go binary
	ballerinaLangGo := "ballerina-lang-go"
	if path, err := exec.LookPath(ballerinaLangGo); err == nil {
		ballerinaLangGo = path
	} else {
		// Try relative path from current directory
		if _, err := os.Stat("./ballerina-lang-go"); err == nil {
			var absPath string
			absPath, err = filepath.Abs("./ballerina-lang-go")
			if err == nil {
				ballerinaLangGo = absPath
			} else {
				ballerinaLangGo = "./ballerina-lang-go"
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error: ballerina-lang-go binary not found in PATH or current directory\n")
			os.Exit(1)
		}
	}

	// Find all .bal files
	var balFiles []string
	err := filepath.Walk(corpusBalDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".bal") {
			balFiles = append(balFiles, path)
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking corpus/bal directory: %v\n", err)
		os.Exit(1)
	}

	if len(balFiles) == 0 {
		fmt.Fprintf(os.Stderr, "No .bal files found in %s\n", corpusBalDir)
		os.Exit(1)
	}

	// Create job channel and worker pool
	jobChan := make(chan string, len(balFiles))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < finalNumJobs; i++ {
		wg.Add(1)
		go worker(jobChan, &wg, ballerinaLangGo, corpusBalDir, corpusTokensDir, corpusParserDir)
	}

	// Send jobs
	for _, file := range balFiles {
		jobChan <- file
	}
	close(jobChan)

	// Wait for all workers to finish
	wg.Wait()
}

func worker(jobChan <-chan string, wg *sync.WaitGroup, ballerinaLangGo, corpusBalDir, corpusTokensDir, corpusParserDir string) {
	defer wg.Done()

	for balFile := range jobChan {
		if err := processFile(ballerinaLangGo, balFile, corpusBalDir, corpusTokensDir, corpusParserDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", balFile, err)
		}
	}
}

func processFile(ballerinaLangGo, balFile, corpusBalDir, corpusTokensDir, corpusParserDir string) error {
	// Print progress
	fmt.Printf("Processing: %s\n", balFile)

	// Calculate relative path (used for both tokens and parser output)
	relPath, err := filepath.Rel(corpusBalDir, balFile)
	if err != nil {
		return fmt.Errorf("calculating relative path: %w", err)
	}

	// Process tokens: Run ballerina-lang-go with -dump-tokens
	tokenCmd := exec.Command(ballerinaLangGo, balFile, "-dump-tokens")
	tokenStderr, err := tokenCmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("creating stderr pipe for tokens: %w", err)
	}

	if err := tokenCmd.Start(); err != nil {
		return fmt.Errorf("starting token command: %w", err)
	}

	// Read all stderr output (tokens)
	tokenOutput, err := io.ReadAll(tokenStderr)
	if err != nil {
		tokenCmd.Wait()
		return fmt.Errorf("reading token stderr: %w", err)
	}

	// Wait for command to complete
	tokenRelPath := strings.TrimSuffix(relPath, ".bal") + ".token"
	tokenOutputPath := filepath.Join(corpusTokensDir, tokenRelPath)
	tokenOutputDir := filepath.Dir(tokenOutputPath)
	if err := os.MkdirAll(tokenOutputDir, 0755); err != nil {
		return fmt.Errorf("creating token output directory: %w", err)
	}

	if err := tokenCmd.Wait(); err != nil {
		// Command crashed or failed - write error output to file for test comparison
		errorMsg := fmt.Sprintf("ERROR: ballerina-lang-go -dump-tokens failed for %s\nExit code: %v\nOutput:\n%s", balFile, err, string(tokenOutput))
		if writeErr := os.WriteFile(tokenOutputPath, []byte(errorMsg), 0644); writeErr != nil {
			return fmt.Errorf("writing token error file: %w", writeErr)
		}
		fmt.Fprintf(os.Stderr, "Warning: %s produced error output (written to %s)\n", balFile, tokenOutputPath)
	} else {
		// Write token file
		if err := os.WriteFile(tokenOutputPath, tokenOutput, 0644); err != nil {
			return fmt.Errorf("writing token file: %w", err)
		}
	}

	// Process parser: Run ballerina-lang-go with -dump-st
	parserCmd := exec.Command(ballerinaLangGo, balFile, "-dump-st")
	parserStderr, err := parserCmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("creating stderr pipe for parser: %w", err)
	}

	if err := parserCmd.Start(); err != nil {
		return fmt.Errorf("starting parser command: %w", err)
	}

	// Read all stderr output (parser JSON)
	parserOutput, err := io.ReadAll(parserStderr)
	if err != nil {
		parserCmd.Wait()
		return fmt.Errorf("reading parser stderr: %w", err)
	}

	// Wait for command to complete
	parserRelPath := strings.TrimSuffix(relPath, ".bal") + ".json"
	parserOutputPath := filepath.Join(corpusParserDir, parserRelPath)
	parserOutputDir := filepath.Dir(parserOutputPath)
	if err := os.MkdirAll(parserOutputDir, 0755); err != nil {
		return fmt.Errorf("creating parser output directory: %w", err)
	}

	if err := parserCmd.Wait(); err != nil {
		// Command crashed or failed - write error output to file for test comparison
		errorMsg := fmt.Sprintf("ERROR: ballerina-lang-go -dump-st failed for %s\nExit code: %v\nOutput:\n%s", balFile, err, string(parserOutput))
		if writeErr := os.WriteFile(parserOutputPath, []byte(errorMsg), 0644); writeErr != nil {
			return fmt.Errorf("writing parser error file: %w", writeErr)
		}
		fmt.Fprintf(os.Stderr, "Warning: %s produced error output (written to %s)\n", balFile, parserOutputPath)
	} else {
		// Write parser JSON file
		if err := os.WriteFile(parserOutputPath, parserOutput, 0644); err != nil {
			return fmt.Errorf("writing parser JSON file: %w", err)
		}
	}

	return nil
}
