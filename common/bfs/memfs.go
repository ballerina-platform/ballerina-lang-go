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
	"fmt"
	"io/fs"
	"strings"
	"time"
)

type memFS struct {
	files map[string]*memFile
}

type memFile struct {
	name  string
	data  []byte
	mode  fs.FileMode
	isDir bool
}

// openMemFile embeds memFile and bytes.Reader to implement fs.File interface
type openMemFile struct {
	*memFile
	*bytes.Reader
}

func (mfs *memFS) Create(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "create", Path: name, Err: fs.ErrInvalid}
	}

	file := &memFile{
		name: name,
	}
	mfs.files[name] = file

	return &openMemFile{
		memFile: file,
		Reader:  bytes.NewReader(file.data),
	}, nil
}

func (mfs *memFS) MkdirAll(path string, perm fs.FileMode) error {
	return nil
}

func (mfs *memFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}

	file, ok := mfs.files[name]
	if !ok {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	return &openMemFile{
		memFile: file,
		Reader:  bytes.NewReader(file.data),
	}, nil
}

func (mfs *memFS) OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "openfile", Path: name, Err: fs.ErrInvalid}
	}

	file, ok := mfs.files[name]
	if !ok {
		file = &memFile{
			name: name,
			mode: perm,
		}
		mfs.files[name] = file
	}

	return &openMemFile{
		memFile: file,
		Reader:  bytes.NewReader(file.data),
	}, nil
}

func (mfs *memFS) Remove(name string) error {
	if _, exists := mfs.files[name]; exists {
		delete(mfs.files, name)
		return nil
	}

	dirPrefix := name
	if !strings.HasSuffix(dirPrefix, "/") {
		dirPrefix = fmt.Sprintf("%s/", dirPrefix)
	}

	removed := false
	for fname := range mfs.files {
		if strings.HasPrefix(fname, dirPrefix) {
			delete(mfs.files, fname)
			removed = true
		}
	}

	if !removed {
		return &fs.PathError{
			Op:   "remove",
			Path: name,
			Err:  fs.ErrNotExist,
		}
	}

	return nil
}

func (mfs *memFS) Move(oldpath, newpath string) error {
	if file, exists := mfs.files[oldpath]; exists {
		delete(mfs.files, oldpath)
		file.name = newpath
		mfs.files[newpath] = file
		return nil
	}

	oldDirPrefix := oldpath
	if !strings.HasSuffix(oldDirPrefix, "/") {
		oldDirPrefix = fmt.Sprintf("%s/", oldDirPrefix)
	}

	newDirPrefix := newpath
	if !strings.HasSuffix(newDirPrefix, "/") {
		newDirPrefix = fmt.Sprintf("%s/", newDirPrefix)
	}

	type moveEntry struct {
		oldName string
		newName string
		file    *memFile
	}

	var toMove []moveEntry

	for name, file := range mfs.files {
		if after, ok := strings.CutPrefix(name, oldDirPrefix); ok {
			newName := newDirPrefix + after
			toMove = append(toMove, moveEntry{
				oldName: name,
				newName: newName,
				file:    file,
			})
		}
	}

	if len(toMove) == 0 {
		return &fs.PathError{
			Op:   "move",
			Path: oldpath,
			Err:  fs.ErrNotExist,
		}
	}

	for _, entry := range toMove {
		delete(mfs.files, entry.oldName)
		entry.file.name = entry.newName
		mfs.files[entry.newName] = entry.file
	}

	return nil
}

func (mfs *memFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "writefile", Path: name, Err: fs.ErrInvalid}
	}

	mfs.files[name] = &memFile{
		name: name,
		data: data,
		mode: perm,
	}

	return nil
}

func (o *openMemFile) Close() error {
	return nil
}

func (o *openMemFile) IsDir() bool {
	return o.isDir
}

func (o *openMemFile) ModTime() time.Time {
	return time.Time{}
}

func (o *openMemFile) Mode() fs.FileMode {
	return o.mode
}

func (o *openMemFile) Name() string {
	return o.name
}

func (o *openMemFile) Read(p []byte) (int, error) {
	return o.Reader.Read(p)
}

func (o *openMemFile) Size() int64 {
	return int64(len(o.data))
}

func (o *openMemFile) Stat() (fs.FileInfo, error) {
	return o, nil
}

func (o *openMemFile) Sys() any {
	return nil
}

func NewMemFS() fs.FS {
	return &memFS{
		files: make(map[string]*memFile),
	}
}
