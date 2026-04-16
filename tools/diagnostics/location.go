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

package diagnostics

import "fmt"

// Location represents the location in TextDocument.
// It is a combination of source file path, start and end line numbers, and start and end column numbers.
type Location struct {
	filePath    string
	startLine   int
	endLine     int
	startColumn int
	endColumn   int
}

// NewLocation creates a new Location with the given file path, line numbers, and column numbers.
func NewLocation(
	filePath string,
	startLine, endLine, startColumn, endColumn int,
) Location {
	return Location{
		filePath:    filePath,
		startLine:   startLine,
		endLine:     endLine,
		startColumn: startColumn,
		endColumn:   endColumn,
	}
}

// IsLocationEmpty returns true if all fields of the Location have zero values.
func IsLocationEmpty(loc Location) bool {
	return loc.filePath == "" && loc.startLine == 0 && loc.endLine == 0 &&
		loc.startColumn == 0 && loc.endColumn == 0
}

// FilePath returns the file path of the Location.
func (loc *Location) FilePath() string {
	return loc.filePath
}

// StartLine returns the start line of the Location.
func (loc *Location) StartLine() int {
	return loc.startLine
}

// StartColumn returns the start column of the Location.
func (loc *Location) StartColumn() int {
	return loc.startColumn
}

// EndLine returns the end line of the Location.
func (loc *Location) EndLine() int {
	return loc.endLine
}

// EndColumn returns the end column of the Location.
func (loc *Location) EndColumn() int {
	return loc.endColumn
}

// String returns a string representation of the Location in the format (startLine:startColumn,endLine:endColumn).
func (loc Location) String() string {
	return fmt.Sprintf("(%d:%d,%d:%d)", loc.startLine, loc.startColumn, loc.endLine, loc.endColumn)
}
