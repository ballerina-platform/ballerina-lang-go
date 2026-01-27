/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// prints an error message in the standard Ballerina CLI format.
func printError(err error, usage string, showHelp bool) {
	fmt.Fprintf(os.Stderr, "ballerina: %s\n", err.Error())
	if usage != "" {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "USAGE:")
		fmt.Fprintf(os.Stderr, "    %s\n", usage)
	}
	if showHelp {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "For more information try --help")
	}
}

// validates the source file argument for the 'run' command.
func validateSourceFile(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		err := fmt.Errorf("source file not provided")
		printError(err, cmd.Use, true)
		return err
	}

	return nil
}
