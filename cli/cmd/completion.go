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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

type shellType string

const (
	shellBash       shellType = "bash"
	shellZsh        shellType = "zsh"
	shellFish       shellType = "fish"
	shellPowerShell shellType = "powershell"
)

const (
	psMarkerBegin = "# BEGIN bal completion"
	psMarkerEnd   = "# END bal completion"
)

var shellFlag string

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Manage shell completions for bal",
	Long: `Manage shell completions for the bal CLI.

  bal completion install       Auto-detect shell and install completions
  bal completion uninstall     Remove installed completions
  bal completion bash          Print bash completion script to stdout
  bal completion zsh           Print zsh completion script to stdout
  bal completion fish          Print fish completion script to stdout
  bal completion powershell    Print powershell completion script to stdout`,
	ValidArgsFunction: cobra.NoFileCompletions,
}

var completionInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install shell completions for bal",
	Long:  `Auto-detects your shell and installs the completion script to the appropriate location.`,
	Args:  cobra.NoArgs,
	RunE:  runInstall,
}

var completionUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove shell completions for bal",
	Long:  `Removes the previously installed completion script.`,
	Args:  cobra.NoArgs,
	RunE:  runUninstall,
}

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generate the autocompletion script for bash",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().GenBashCompletionV2(cmd.OutOrStdout(), true)
	},
}

var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generate the autocompletion script for zsh",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
	},
}

var completionFishCmd = &cobra.Command{
	Use:   "fish",
	Short: "Generate the autocompletion script for fish",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
	},
}

var completionPowerShellCmd = &cobra.Command{
	Use:   "powershell",
	Short: "Generate the autocompletion script for powershell",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
	},
}

func init() {
	completionCmd.PersistentFlags().StringVar(&shellFlag, "shell", "",
		"Override shell detection (bash|zsh|fish|powershell)")

	completionCmd.AddCommand(completionInstallCmd)
	completionCmd.AddCommand(completionUninstallCmd)
	completionCmd.AddCommand(completionBashCmd)
	completionCmd.AddCommand(completionZshCmd)
	completionCmd.AddCommand(completionFishCmd)
	completionCmd.AddCommand(completionPowerShellCmd)
}

// detectShell determines the user's shell from --shell flag or $SHELL env var.
func detectShell(flagOverride string) (shellType, error) {
	if flagOverride != "" {
		switch shellType(flagOverride) {
		case shellBash, shellZsh, shellFish, shellPowerShell:
			return shellType(flagOverride), nil
		default:
			return "", fmt.Errorf("unsupported shell %q; valid options: bash, zsh, fish, powershell", flagOverride)
		}
	}

	shellEnv := os.Getenv("SHELL")
	if shellEnv != "" {
		base := filepath.Base(shellEnv)
		switch {
		case strings.Contains(base, "bash"):
			return shellBash, nil
		case strings.Contains(base, "zsh"):
			return shellZsh, nil
		case strings.Contains(base, "fish"):
			return shellFish, nil
		}
	}

	if os.Getenv("PSModulePath") != "" {
		return shellPowerShell, nil
	}

	return "", fmt.Errorf("unable to detect shell; use --shell (bash|zsh|fish|powershell)")
}

// completionFilePath returns the target path for the completion script.
// All paths are user-level (no sudo required).
func completionFilePath(shell shellType) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}

	switch shell {
	case shellBash:
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, ".bash_completion.d", "bal"), nil
		}
		return filepath.Join(home, ".local", "share", "bash-completion", "completions", "bal"), nil

	case shellZsh:
		omz := filepath.Join(home, ".oh-my-zsh")
		if info, err := os.Stat(omz); err == nil && info.IsDir() {
			return filepath.Join(omz, "completions", "_bal"), nil
		}
		return filepath.Join(home, ".zsh", "completions", "_bal"), nil

	case shellFish:
		return filepath.Join(home, ".config", "fish", "completions", "bal.fish"), nil

	case shellPowerShell:
		if runtime.GOOS == "windows" {
			return filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1"), nil
		}
		return filepath.Join(home, ".config", "powershell", "Microsoft.PowerShell_profile.ps1"), nil

	default:
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}
}

// generateCompletionScript produces the completion script bytes for the given shell.
func generateCompletionScript(root *cobra.Command, shell shellType) ([]byte, error) {
	var buf bytes.Buffer
	var err error

	switch shell {
	case shellBash:
		err = root.GenBashCompletionV2(&buf, true)
	case shellZsh:
		err = root.GenZshCompletion(&buf)
	case shellFish:
		err = root.GenFishCompletion(&buf, true)
	case shellPowerShell:
		err = root.GenPowerShellCompletionWithDesc(&buf)
	default:
		return nil, fmt.Errorf("unsupported shell: %s", shell)
	}

	if err != nil {
		return nil, fmt.Errorf("generating %s completion: %w", shell, err)
	}
	return buf.Bytes(), nil
}

func runInstall(cmd *cobra.Command, args []string) error {
	shell, err := detectShell(shellFlag)
	if err != nil {
		return err
	}

	targetPath, err := completionFilePath(shell)
	if err != nil {
		return err
	}

	script, err := generateCompletionScript(cmd.Root(), shell)
	if err != nil {
		return err
	}

	_, existed := os.Stat(targetPath)
	isUpdate := existed == nil

	if shell == shellPowerShell {
		if err := installPowerShell(targetPath, script); err != nil {
			return err
		}
	} else {
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("creating directory %s: %w", filepath.Dir(targetPath), err)
		}
		if err := os.WriteFile(targetPath, script, 0644); err != nil {
			return fmt.Errorf("writing completion to %s: %w", targetPath, err)
		}
	}

	verb := "Installed"
	if isUpdate {
		verb = "Updated"
	}
	fmt.Printf("%s bal completions for %s to %s\n", verb, shell, targetPath)

	printActivationHint(shell, targetPath)
	return nil
}

func runUninstall(cmd *cobra.Command, args []string) error {
	shell, err := detectShell(shellFlag)
	if err != nil {
		return err
	}

	targetPath, err := completionFilePath(shell)
	if err != nil {
		return err
	}

	if shell == shellPowerShell {
		if err := uninstallPowerShell(targetPath); err != nil {
			return err
		}
		fmt.Printf("Removed bal completions for %s from %s\n", shell, targetPath)
		return nil
	}

	err = os.Remove(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("No completion file found for %s at %s\n", shell, targetPath)
			return nil
		}
		return fmt.Errorf("removing %s: %w", targetPath, err)
	}

	fmt.Printf("Removed bal completions for %s from %s\n", shell, targetPath)
	return nil
}

// installPowerShell writes the completion script into the PowerShell profile
// using markers to make repeated installs idempotent.
func installPowerShell(profilePath string, script []byte) error {
	existing, _ := os.ReadFile(profilePath)
	block := fmt.Sprintf("%s\n%s\n%s", psMarkerBegin, strings.TrimSpace(string(script)), psMarkerEnd)

	content := string(existing)
	if beginIdx := strings.Index(content, psMarkerBegin); beginIdx >= 0 {
		if endIdx := strings.Index(content, psMarkerEnd); endIdx >= 0 {
			content = content[:beginIdx] + block + content[endIdx+len(psMarkerEnd):]
		}
	} else {
		if len(content) > 0 && !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		content += block + "\n"
	}

	if err := os.MkdirAll(filepath.Dir(profilePath), 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", filepath.Dir(profilePath), err)
	}
	return os.WriteFile(profilePath, []byte(content), 0644)
}

// uninstallPowerShell removes the marked completion block from the PowerShell profile.
func uninstallPowerShell(profilePath string) error {
	data, err := os.ReadFile(profilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("No completion found for powershell at %s\n", profilePath)
			return nil
		}
		return err
	}

	content := string(data)
	beginIdx := strings.Index(content, psMarkerBegin)
	if beginIdx < 0 {
		fmt.Printf("No completion found for powershell at %s\n", profilePath)
		return nil
	}

	endIdx := strings.Index(content, psMarkerEnd)
	if endIdx < 0 {
		return fmt.Errorf("found %s but no matching %s in %s", psMarkerBegin, psMarkerEnd, profilePath)
	}

	// Remove the block and any trailing newline
	after := endIdx + len(psMarkerEnd)
	if after < len(content) && content[after] == '\n' {
		after++
	}
	content = content[:beginIdx] + content[after:]

	if strings.TrimSpace(content) == "" {
		return os.Remove(profilePath)
	}
	return os.WriteFile(profilePath, []byte(content), 0644)
}

func printActivationHint(shell shellType, path string) {
	switch shell {
	case shellBash:
		if runtime.GOOS == "darwin" {
			dir := filepath.Dir(path)
			fmt.Println()
			fmt.Println("To activate, add to your ~/.bashrc or ~/.bash_profile:")
			fmt.Printf("  for f in %s/*; do [ -f \"$f\" ] && source \"$f\"; done\n", dir)
		} else {
			fmt.Println()
			fmt.Println("Completions will load automatically if the bash-completion package is installed.")
			fmt.Println("Otherwise, add to your ~/.bashrc:")
			fmt.Printf("  source %s\n", path)
		}

	case shellZsh:
		dir := filepath.Dir(path)
		if !strings.Contains(os.Getenv("FPATH"), dir) {
			fmt.Println()
			fmt.Println("To activate, add to your ~/.zshrc:")
			fmt.Printf("  fpath=(%s $fpath)\n", dir)
			fmt.Println("  autoload -U compinit && compinit")
		} else {
			fmt.Println()
			fmt.Println("Run `exec zsh` or open a new terminal to activate.")
		}

	case shellFish:
		fmt.Println()
		fmt.Println("Restart your shell or open a new terminal to activate.")

	case shellPowerShell:
		fmt.Println()
		fmt.Println("Restart PowerShell to activate.")
	}
}
