// Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
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

package bfs

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"testing"
)

func TestWriteFileOpenAndStat(t *testing.T) {
	fsys := NewMemFS()

	name := "dir/hello.txt"
	content := []byte("hello world")
	perm := fs.FileMode(0o644)

	if err := WriteFile(fsys, name, content, perm); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	f, err := fsys.Open(name)
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}
	if !bytes.Equal(data, content) {
		t.Fatalf("content mismatch: got %q want %q", string(data), string(content))
	}

	info, err := f.Stat()
	if err != nil {
		t.Fatalf("Stat error: %v", err)
	}
	if info.Size() != int64(len(content)) {
		t.Errorf("Size mismatch: got %d want %d", info.Size(), len(content))
	}
	if info.Mode() != perm {
		t.Errorf("Mode mismatch: got %v want %v", info.Mode(), perm)
	}
	if info.IsDir() {
		t.Errorf("file unexpectedly marked as directory")
	}
}

func TestCreateAndOpen_InvalidPath(t *testing.T) {
	fsys := NewMemFS()
	invalid := "/abs/path.txt"

	if _, err := Create(fsys, invalid); err == nil || !errors.Is(err, fs.ErrInvalid) {
		t.Fatalf("Create: expected ErrInvalid, got %v", err)
	}
	if _, err := fsys.Open(invalid); err == nil || !errors.Is(err, fs.ErrInvalid) {
		t.Fatalf("Open: expected ErrInvalid, got %v", err)
	}
	if _, err := OpenFile(fsys, invalid, 0, 0o644); err == nil || !errors.Is(err, fs.ErrInvalid) {
		t.Fatalf("OpenFile: expected ErrInvalid, got %v", err)
	}
}

func TestOpen_NotExist(t *testing.T) {
	fsys := NewMemFS()
	if _, err := fsys.Open("no/such/file.txt"); err == nil || !os.IsNotExist(err) {
		t.Fatalf("expected ErrNotExist, got %v", err)
	}
}

func TestRemove_FileAndDirectoryPrefix(t *testing.T) {
	fsys := NewMemFS()

	_ = WriteFile(fsys, "dir/a.txt", []byte("A"), 0o644)
	_ = WriteFile(fsys, "dir/sub/b.txt", []byte("B"), 0o644)
	_ = WriteFile(fsys, "c.txt", []byte("C"), 0o644)

	if err := Remove(fsys, "c.txt"); err != nil {
		t.Fatalf("Remove file error: %v", err)
	}
	if _, err := fsys.Open("c.txt"); !os.IsNotExist(err) {
		t.Fatalf("file should be removed, got err: %v", err)
	}

	if err := Remove(fsys, "dir"); err != nil {
		t.Fatalf("Remove dir prefix error: %v", err)
	}
	if _, err := fsys.Open("dir/a.txt"); !os.IsNotExist(err) {
		t.Fatalf("dir/a.txt should be removed, got err: %v", err)
	}
	if _, err := fsys.Open("dir/sub/b.txt"); !os.IsNotExist(err) {
		t.Fatalf("dir/sub/b.txt should be removed, got err: %v", err)
	}

	if err := Remove(fsys, "does/not/exist"); err == nil || !os.IsNotExist(err) {
		t.Fatalf("expected ErrNotExist when removing non-existent, got %v", err)
	}
}

func TestMove_FileAndDirectory(t *testing.T) {
	fsys := NewMemFS()

	_ = WriteFile(fsys, "a.txt", []byte("alpha"), 0o644)
	if err := Move(fsys, "a.txt", "b.txt"); err != nil {
		t.Fatalf("Move file error: %v", err)
	}
	if _, err := fsys.Open("a.txt"); !os.IsNotExist(err) {
		t.Fatalf("source should be gone after move, err: %v", err)
	}
	bf, err := fsys.Open("b.txt")
	if err != nil {
		t.Fatalf("destination open error: %v", err)
	}
	got, _ := io.ReadAll(bf)
	_ = bf.Close()
	if string(got) != "alpha" {
		t.Fatalf("moved file content mismatch: got %q", string(got))
	}

	_ = WriteFile(fsys, "src/x.txt", []byte("x"), 0o644)
	_ = WriteFile(fsys, "src/sub/y.txt", []byte("y"), 0o644)
	if err := Move(fsys, "src", "dst"); err != nil {
		t.Fatalf("Move dir error: %v", err)
	}
	if _, err := fsys.Open("src/x.txt"); !os.IsNotExist(err) {
		t.Fatalf("old prefixed file should not exist, err: %v", err)
	}
	if f, err := fsys.Open("dst/x.txt"); err != nil {
		t.Fatalf("dst/x.txt open error: %v", err)
	} else {
		f.Close()
	}
	if f, err := fsys.Open("dst/sub/y.txt"); err != nil {
		t.Fatalf("dst/sub/y.txt open error: %v", err)
	} else {
		f.Close()
	}

	if err := Move(fsys, "nope", "new"); err == nil || !os.IsNotExist(err) {
		t.Fatalf("expected ErrNotExist for move, got %v", err)
	}
}
