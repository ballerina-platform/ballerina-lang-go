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

import "fmt"

// TextEdit represents a text edit on a TextDocument.
type TextEdit interface {
	Range() TextRange
	Text() string
	String() string
}

type textEditImpl struct {
	textRange TextRange
	text      string
}

func TextEditFromTextRangeAndText(textRange TextRange, text string) TextEdit {
	return &textEditImpl{
		textRange: textRange,
		text:      text,
	}
}

func (te textEditImpl) Range() TextRange {
	return te.textRange
}

func (te textEditImpl) Text() string {
	return te.text
}

func (te textEditImpl) String() string {
	return fmt.Sprintf("%s%s", te.textRange.String(), te.text)
}
