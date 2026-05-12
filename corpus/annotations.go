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

package corpus

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// outputAnn is an `// @output <value>` annotation.
type outputAnn struct {
	line  int    // 1-based source line
	value string // text after `@output ` with trailing whitespace trimmed
}

// errorAnn is an `// @error` annotation. Payload (if any) is ignored.
type errorAnn struct {
	line int
}

// panicAnn is an `// @panic <message>` annotation. The message after the
// keyword is ignored; assertions match on (file, line) of the annotation.
type panicAnn struct {
	line int
}

// fileAnns groups every annotation found in a single source file. Slices
// are kept in line order (which is the source order within that file).
type fileAnns struct {
	outputs []outputAnn
	errors  []errorAnn
	panics  []panicAnn
}

// annotations maps a file key (single-file: base name; project: slash-
// separated path relative to the project root) to its annotations. Files
// with no annotations are not present in the map.
type annotations map[string]*fileAnns

// annSourceFile pairs a path key (used for matching against diagnostic file
// names and stack frames) with the absolute filesystem path to read.
type annSourceFile struct {
	key     string
	absPath string
}

var (
	// `// @output value`, `//@output value`, `//  @output  value`.
	outputRe = regexp.MustCompile(`(?://)\s*@output\b[ \t]?(.*)$`)
	errorRe  = regexp.MustCompile(`(?://)\s*@error\b`)
	panicRe  = regexp.MustCompile(`(?://)\s*@panic\b`)
)

func parseAnnotations(sources []annSourceFile) (annotations, error) {
	anns := annotations{}
	for _, src := range sources {
		content, err := os.ReadFile(src.absPath)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", src.absPath, err)
		}
		if fa := parseAnnotationsInFile(string(content)); fa != nil {
			anns[src.key] = fa
		}
	}
	return anns, nil
}

func parseAnnotationsInFile(content string) *fileAnns {
	var fa fileAnns
	lines := strings.Split(content, "\n")
	for i, raw := range lines {
		lineNo := i + 1
		line := strings.TrimSuffix(raw, "\r")
		commentIdx := strings.Index(line, "//")
		if commentIdx < 0 {
			continue
		}
		comment := line[commentIdx:]
		if m := outputRe.FindStringSubmatch(comment); m != nil {
			fa.outputs = append(fa.outputs, outputAnn{
				line:  lineNo,
				value: strings.TrimRight(m[1], " \t"),
			})
			continue
		}
		if panicRe.MatchString(comment) {
			fa.panics = append(fa.panics, panicAnn{line: lineNo})
			continue
		}
		if errorRe.MatchString(comment) {
			fa.errors = append(fa.errors, errorAnn{line: lineNo})
		}
	}
	if fa.outputs == nil && fa.errors == nil && fa.panics == nil {
		return nil
	}
	return &fa
}

// collectSingleFileSources returns the one source file used by a single-file
// corpus test. The key is the base name (matches diagnostic registrations and
// stack frames).
func collectSingleFileSources(balFile string) []annSourceFile {
	return []annSourceFile{{
		key:     filepath.Base(balFile),
		absPath: balFile,
	}}
}

// collectProjectSources returns every `.bal` file under projectDir, keyed by
// its slash-separated path relative to projectDir.
func collectProjectSources(projectDir string) ([]annSourceFile, error) {
	var sources []annSourceFile
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".bal") {
			return nil
		}
		rel, err := filepath.Rel(projectDir, path)
		if err != nil {
			return err
		}
		sources = append(sources, annSourceFile{
			key:     filepath.ToSlash(rel),
			absPath: path,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return sources, nil
}
