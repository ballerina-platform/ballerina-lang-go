/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
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

package text

import "strings"

// TextDocumentChange represents textual changes on a single TextDocument.
type TextDocumentChange interface {
	GetTextEditCount() int
	GetTextEdit(index int) TextEdit
	String() string
}

type textDocumentChangeImpl struct {
	textEdits []TextEdit
}

func TextDocumentChangeFromTextEdits(textEdits []TextEdit) TextDocumentChange {
	// Create a copy of the slice to ensure immutability
	editsCopy := make([]TextEdit, len(textEdits))
	copy(editsCopy, textEdits)

	return &textDocumentChangeImpl{
		textEdits: editsCopy,
	}
}

func (tdc textDocumentChangeImpl) GetTextEditCount() int {
	return len(tdc.textEdits)
}

func (tdc textDocumentChangeImpl) GetTextEdit(index int) TextEdit {
	return tdc.textEdits[index]
}

func (tdc textDocumentChangeImpl) String() string {
	var editStrings []string
	for _, textEdit := range tdc.textEdits {
		editStrings = append(editStrings, textEdit.String())
	}

	return strings.Join(editStrings, ",")
}
