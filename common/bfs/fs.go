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
	"io/fs"
)

type MutableFS interface {
	fs.FS

	Create(name string) (fs.File, error)
	MkdirAll(path string, perm fs.FileMode) error
	OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error)
	Remove(name string) error
	Move(oldpath, newpath string) error
}

type WritableFS interface {
	fs.FS

	WriteFile(name string, data []byte, perm fs.FileMode) error
}

func Create(fsys fs.FS, name string) (fs.File, error) {
	mfs, ok := fsys.(MutableFS)
	if !ok {
		return nil, &fs.PathError{Op: "create", Path: name, Err: fs.ErrInvalid}
	}
	return mfs.Create(name)
}

func MkdirAll(fsys fs.FS, path string, perm fs.FileMode) error {
	mfs, ok := fsys.(MutableFS)
	if !ok {
		return &fs.PathError{Op: "mkdirall", Path: path, Err: fs.ErrInvalid}
	}
	return mfs.MkdirAll(path, perm)
}

func OpenFile(fsys fs.FS, name string, flag int, perm fs.FileMode) (fs.File, error) {
	mfs, ok := fsys.(MutableFS)
	if !ok {
		return nil, &fs.PathError{Op: "openfile", Path: name, Err: fs.ErrInvalid}
	}
	return mfs.OpenFile(name, flag, perm)
}

func Remove(fsys fs.FS, name string) error {
	mfs, ok := fsys.(MutableFS)
	if !ok {
		return &fs.PathError{Op: "remove", Path: name, Err: fs.ErrInvalid}
	}
	return mfs.Remove(name)
}

func Move(fsys fs.FS, oldpath, newpath string) error {
	mfs, ok := fsys.(MutableFS)
	if !ok {
		return &fs.PathError{Op: "move", Path: oldpath + "->" + newpath, Err: fs.ErrInvalid}
	}
	return mfs.Move(oldpath, newpath)
}

func WriteFile(fsys fs.FS, name string, data []byte, perm fs.FileMode) error {
	wfs, ok := fsys.(WritableFS)
	if !ok {
		return &fs.PathError{Op: "writefile", Path: name, Err: fs.ErrInvalid}
	}
	return wfs.WriteFile(name, data, perm)
}
