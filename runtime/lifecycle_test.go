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

package runtime_test

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "ballerina-lang-go/lib/rt"
	"ballerina-lang-go/platform/pal"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/runtime"
)

const lifecycleTestSource = `
import ballerina/io;

class ListenerOne {
    public function attach(service object {} svc, () attachPoint = ()) returns error? {
        var _ = svc;
        var _ = attachPoint;
    }

    public function detach(service object {} svc) returns error? {
        var _ = svc;
    }

    public function 'start() returns error? {
        io:println("start:one");
    }

    public function gracefulStop() returns error? {
        io:println("graceful:one");
    }

    public function immediateStop() returns error? {
        io:println("immediate:one");
    }
}

class ListenerTwo {
    public function attach(service object {} svc, () attachPoint = ()) returns error? {
        var _ = svc;
        var _ = attachPoint;
    }

    public function detach(service object {} svc) returns error? {
        var _ = svc;
    }

    public function 'start() returns error? {
        io:println("start:two");
    }

    public function gracefulStop() returns error? {
        io:println("graceful:two");
    }

    public function immediateStop() returns error? {
        io:println("immediate:two");
    }
}

listener ListenerOne l1 = new ();
listener ListenerTwo l2 = new ();

service on l1 {
}

service on l2 {
}
`

func TestLifecycleGracefulStopSignal(t *testing.T) {
	pal := newLifecycleTestPal()
	rt := newLifecycleTestRuntime(t, lifecycleTestSource, pal)

	rt.Listen()
	pal.Send(palSignalGracefulStop)
	code := readExitStatus(t, rt)

	if code != 130 {
		t.Fatalf("expected graceful stop exit code 130, got %d", code)
	}
	if got, want := pal.Stdout(), "start:one\nstart:two\ngraceful:one\ngraceful:two\n"; got != want {
		t.Fatalf("unexpected stdout: got %q, want %q", got, want)
	}
}

func TestLifecycleImmediateStopSignal(t *testing.T) {
	pal := newLifecycleTestPal()
	rt := newLifecycleTestRuntime(t, lifecycleTestSource, pal)

	rt.Listen()
	pal.Send(palSignalImmediateStop)
	code := readExitStatus(t, rt)

	if code != 131 {
		t.Fatalf("expected immediate stop exit code 131, got %d", code)
	}
	if got, want := pal.Stdout(), "start:one\nstart:two\nimmediate:one\nimmediate:two\n"; got != want {
		t.Fatalf("unexpected stdout: got %q, want %q", got, want)
	}
}

func newLifecycleTestRuntime(t *testing.T, source string, platform *lifecycleTestPal) *runtime.Runtime {
	t.Helper()

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "Ballerina.toml"), `[package]
org = "testorg"
name = "lifecycletest"
version = "0.1.0"
`)
	writeFile(t, filepath.Join(dir, "main.bal"), source)

	ballerinaEnvFs, err := ballerinaEnvFS()
	if err != nil {
		t.Fatal(err)
	}
	result, err := projects.Load(os.DirFS(dir), ".", projects.ProjectLoadConfig{BallerinaEnvFs: ballerinaEnvFs})
	if err != nil {
		t.Fatal(err)
	}
	compilation := result.Project().CurrentPackage().Compilation()
	if result.Diagnostics().HasErrors() || compilation.DiagnosticResult().HasErrors() {
		t.Fatalf("lifecycle test project has diagnostics: load=%v compile=%v", result.Diagnostics().Errors(), compilation.DiagnosticResult().Errors())
	}
	pkgs := projects.NewBallerinaBackend(compilation).BIRPackages()
	if len(pkgs) == 0 {
		t.Fatal("compilation succeeded but produced no BIR packages")
	}

	rt := runtime.NewRuntime(platform.Platform(), result.Project().Environment().TypeEnv())
	for _, pkg := range pkgs {
		if err := rt.Init(*pkg); err != nil {
			t.Fatal(err)
		}
	}
	return rt
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func ballerinaEnvFS() (fs.FS, error) {
	if v := os.Getenv(projects.BallerinaEnvVar); v != "" {
		return os.DirFS(v), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return os.DirFS(filepath.Join(home, projects.UserHomeDirName)), nil
}

func readExitStatus(t *testing.T, rt *runtime.Runtime) uint8 {
	t.Helper()
	select {
	case code := <-rt.ExitStatus:
		return code
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for runtime exit status")
	}
	return 0
}

type lifecycleTestPal struct {
	stdout  bytes.Buffer
	stderr  bytes.Buffer
	signals chan pal.Signal
}

const (
	palSignalGracefulStop  = pal.GracefulStop
	palSignalImmediateStop = pal.ImmediateStop
)

func newLifecycleTestPal() *lifecycleTestPal {
	return &lifecycleTestPal{signals: make(chan pal.Signal, 4)}
}

func (p *lifecycleTestPal) Platform() pal.Platform {
	return pal.Platform{
		IO: pal.IO{
			Stdout: p.stdout.Write,
			Stderr: p.stderr.Write,
		},
		FS: pal.FS{
			ReadFile: func(path string) ([]byte, error) {
				return nil, &fs.PathError{Op: "open", Path: path, Err: fs.ErrNotExist}
			},
		},
		HTTP: pal.HTTP{
			NewClient: func(_ pal.ClientConfig) pal.HTTPClient { return nil },
		},
		Signals: pal.SignalSource{Signals: p.signals},
	}
}

func (p *lifecycleTestPal) Send(signal pal.Signal) {
	p.signals <- signal
}

func (p *lifecycleTestPal) Stdout() string {
	return p.stdout.String()
}

func (p *lifecycleTestPal) String() string {
	return fmt.Sprintf("stdout=%q stderr=%q", p.stdout.String(), p.stderr.String())
}
