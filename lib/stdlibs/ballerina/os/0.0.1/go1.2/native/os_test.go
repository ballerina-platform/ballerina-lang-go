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

package native

import (
	"errors"
	"testing"

	"ballerina-lang-go/platform/pal"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util/testharness"
	"ballerina-lang-go/values"
)

// fakeProcess is a controllable pal.ProcessHandle for exercising the exec path
// without spawning a real subprocess (keeps the test cross-platform).
type fakeProcess struct {
	code    int
	stdout  []byte
	stderr  []byte
	waitErr error
	outErr  error
	killed  bool
}

func (f *fakeProcess) WaitForExit() (int, error)   { return f.code, f.waitErr }
func (f *fakeProcess) ReadStdout() ([]byte, error) { return f.stdout, f.outErr }
func (f *fakeProcess) ReadStderr() ([]byte, error) { return f.stderr, f.outErr }
func (f *fakeProcess) Kill()                       { f.killed = true }

type osTestEnv struct {
	rt        *runtime.Runtime
	tc        semtypes.Context
	byteArrTy semtypes.SemType
	strArrTy  semtypes.SemType
	mapAtomic *semtypes.MappingAtomicType
}

func newOSTestEnv(execFn func(string, []string, map[string]string) (pal.ProcessHandle, error)) osTestEnv {
	env := semtypes.CreateTypeEnv()
	tc := semtypes.ContextFrom(env)
	platform := testharness.NewTestPal().Platform()
	platform.OS.Exec = execFn
	bld := semtypes.NewListDefinition()
	sld := semtypes.NewListDefinition()
	return osTestEnv{
		rt:        runtime.NewRuntime(platform, env),
		tc:        tc,
		byteArrTy: bld.DefineListTypeWrappedWithEnvSemType(env, semtypes.BYTE),
		strArrTy:  sld.DefineListTypeWrappedWithEnvSemType(env, semtypes.STRING),
		mapAtomic: semtypes.ToMappingAtomicType(tc, semtypes.MAPPING),
	}
}

func (e osTestEnv) command(value string, argStrs ...string) *values.Map {
	args := make([]values.BalValue, len(argStrs))
	for i, s := range argStrs {
		args[i] = s
	}
	argList := values.NewList(e.strArrTy, semtypes.ToListAtomicType(e.tc, e.strArrTy), false, nil, 0, args)
	return values.NewMap(semtypes.MAPPING, e.mapAtomic, false, []values.MapEntry{
		{Key: "value", Value: value},
		{Key: "arguments", Value: argList},
	})
}

func TestExecCommandSuccess(t *testing.T) {
	var gotCmd string
	var gotArgs []string
	var gotEnv map[string]string
	handle := &fakeProcess{code: 0, stdout: []byte("hi"), stderr: []byte("err")}
	e := newOSTestEnv(func(cmd string, args []string, env map[string]string) (pal.ProcessHandle, error) {
		gotCmd, gotArgs, gotEnv = cmd, args, env
		return handle, nil
	})

	// exec with arguments and an environment override.
	envMap := values.NewMap(semtypes.MAPPING, e.mapAtomic, false, []values.MapEntry{{Key: "FOO", Value: "bar"}})
	result, err := execCommand(e.rt, []values.BalValue{e.command("echo", "hi"), envMap})
	if err != nil {
		t.Fatalf("execCommand: %v", err)
	}
	proc, ok := result.(*values.Object)
	if !ok {
		t.Fatalf("execCommand returned %T, want *values.Object", result)
	}
	if gotCmd != "echo" || len(gotArgs) != 1 || gotArgs[0] != "hi" || gotEnv["FOO"] != "bar" {
		t.Errorf("exec received cmd=%q args=%v env=%v", gotCmd, gotArgs, gotEnv)
	}

	// waitForExit returns the exit code.
	code, err := processWaitForExit([]values.BalValue{proc})
	if err != nil {
		t.Fatalf("waitForExit: %v", err)
	}
	if code != int64(0) {
		t.Errorf("waitForExit = %v, want 0", code)
	}

	// output(stdout) and output(stderr) return the respective byte streams.
	stdout, err := processOutput(e.tc, e.byteArrTy, []values.BalValue{proc, int64(1)})
	if err != nil {
		t.Fatalf("output stdout: %v", err)
	}
	if l, ok := stdout.(*values.List); !ok || l.Len() != 2 {
		t.Errorf("stdout length: got %v", stdout)
	}
	stderr, err := processOutput(e.tc, e.byteArrTy, []values.BalValue{proc, int64(2)})
	if err != nil {
		t.Fatalf("output stderr: %v", err)
	}
	if l, ok := stderr.(*values.List); !ok || l.Len() != 3 {
		t.Errorf("stderr length: got %v", stderr)
	}

	// exit kills the process.
	if _, err := processExit([]values.BalValue{proc}); err != nil {
		t.Fatalf("exit: %v", err)
	}
	if !handle.killed {
		t.Error("exit did not kill the process")
	}
}

func TestExecCommandErrors(t *testing.T) {
	// exec failure surfaces an os error value.
	failEnv := newOSTestEnv(func(string, []string, map[string]string) (pal.ProcessHandle, error) {
		return nil, errors.New("no such command")
	})
	res, err := execCommand(failEnv.rt, []values.BalValue{failEnv.command("missing")})
	if err != nil {
		t.Fatalf("execCommand: %v", err)
	}
	if _, ok := res.(*values.Error); !ok {
		t.Errorf("exec failure: got %T, want *values.Error", res)
	}

	// waitForExit and output failures surface os error values.
	handle := &fakeProcess{waitErr: errors.New("wait boom"), outErr: errors.New("read boom")}
	okEnv := newOSTestEnv(func(string, []string, map[string]string) (pal.ProcessHandle, error) {
		return handle, nil
	})
	proc, _ := execCommand(okEnv.rt, []values.BalValue{okEnv.command("x")})
	procObj := proc.(*values.Object)

	if res, _ := processWaitForExit([]values.BalValue{procObj}); !isErr(res) {
		t.Errorf("waitForExit error: got %T", res)
	}
	if res, _ := processOutput(okEnv.tc, okEnv.byteArrTy, []values.BalValue{procObj, int64(1)}); !isErr(res) {
		t.Errorf("output error: got %T", res)
	}
}

func isErr(v values.BalValue) bool {
	_, ok := v.(*values.Error)
	return ok
}
