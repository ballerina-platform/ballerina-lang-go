/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
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

package lexer

const (
	// Keywords
	kwTrue  = "true"
	kwFalse = "false"
	kwInf   = "inf"
	kwNan   = "nan"

	// Single char constants — rune values
	charSemicolon      = ';'
	charDot            = '.'
	charComma          = ','
	charOpenBrace      = '{'
	charCloseBrace     = '}'
	charOpenBracket    = '['
	charCloseBracket   = ']'
	charDoubleQuote    = '"'
	charSingleQuote    = '\''
	charHash           = '#'
	charEqual          = '='
	charPlus           = '+'
	charMinus          = '-'
	charNewline        = '\n'
	charCarriageReturn = '\r'
	charTab            = '\t'
	charSpace          = ' '
	charFormFeed       = '\f'
	charBackslash      = '\\'
	charEOF            = rune(-1)
)
