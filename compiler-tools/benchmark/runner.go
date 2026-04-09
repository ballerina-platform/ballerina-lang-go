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
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	baseInterpreter = "bal-base"
	headInterpreter = "bal-base"
)

type benchmark struct {
	config
	workRoot string
}

func (b *benchmark) run() error {
	workRoot, err := os.MkdirTemp("", "bal-bench-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(workRoot) }()
	b.workRoot = workRoot

	baseWorktree, err := b.checkoutWorktree(b.baseRef)
	if err != nil {
		return err
	}
	defer b.removeWorktree(baseWorktree)

	headWorktree, err := b.checkoutWorktree(b.headRef)
	if err != nil {
		return err
	}
	defer b.removeWorktree(headWorktree)

	fmt.Printf("Building interpreter for base ref %q...\n", b.baseRef)
	if err := b.buildInterpreter(baseWorktree, b.baseRef, baseInterpreter); err != nil {
		return err
	}
	fmt.Printf("Building interpreter for head ref %q...\n", b.headRef)
	if err := b.buildInterpreter(headWorktree, b.headRef, headInterpreter); err != nil {
		return err
	}

	target, err := resolveTarget(b.target)
	if err != nil {
		return fmt.Errorf("failed to resolve benchmark target: %w", err)
	}

	for _, path := range target.paths {
		cmds := b.benchmarkCmdPair(baseWorktree, headWorktree, path, target.mode)
		if err := b.runHyperfine(cmds); err != nil {
			return err
		}
	}

	return nil
}

func (b *benchmark) checkoutWorktree(ref string) (string, error) {
	path := filepath.Join(b.workRoot, "worktree"+sanitizeRef(ref))
	b.removeWorktree(path)

	if err := runCmd(".", "git", "worktree", "add", "--detach", path, ref); err != nil {
		return "", fmt.Errorf("failed to checkout worktree for ref %q: %w", ref, err)
	}
	return path, nil
}

func (b *benchmark) removeWorktree(path string) {
	_ = runCmdSilent(".", "git", "worktree", "remove", "--force", path)
}

func (b *benchmark) buildInterpreter(worktreePath, ref, output string) error {
	if err := runCmd(worktreePath, "go", "build", "-o", output, "./cli/cmd"); err != nil {
		return fmt.Errorf("failed to build interpreter for ref %q: %w", ref, err)
	}
	return nil
}

func (b *benchmark) hyperfineFlags() []string {
	var args []string
	if b.warmup > 0 {
		args = append(args, "--warmup", strconv.Itoa(b.warmup))
	}
	if b.runs > 0 {
		args = append(args, "--runs", strconv.Itoa(b.runs))
	}
	return args
}

func (b *benchmark) runHyperfine(cmds []string) error {
	args := b.hyperfineFlags()
	args = append(args, cmds...)
	return runCmd(".", "hyperfine", args...)
}

func (b *benchmark) benchmarkCmdArgs(ref, interpreter, target string, mode targetMode) []string {
	command := fmt.Sprintf("%s run %s", shellQuote(interpreter), shellQuote(target))
	if mode == multipleFilesMode {
		ref = fmt.Sprintf("%s (%s)", ref, filepath.Base(target))
	}
	return []string{"--command-name", ref, command}
}

func (b *benchmark) benchmarkCmdPair(baseWorktree, headWorktree, target string, mode targetMode) []string {
	baseCmd := b.benchmarkCmdArgs(b.baseRef, filepath.Join(baseWorktree, baseInterpreter), target, mode)
	headCmd := b.benchmarkCmdArgs(b.headRef, filepath.Join(headWorktree, headInterpreter), target, mode)
	return append(baseCmd, headCmd...)
}

func sanitizeRef(ref string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '/', '\\', ':':
			return '-'
		default:
			return r
		}
	}, ref)
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'"'"'`) + "'"
}

func runCmd(dir, name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func runCmdSilent(dir, name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Dir = dir
	c.Stdout = nil
	c.Stderr = nil
	return c.Run()
}
