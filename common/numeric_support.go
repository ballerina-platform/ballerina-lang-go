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

package common

// HasHexIndicator checks if literal has hex indicator
// migrated from NumericLiteralSupport.java:77:5
func HasHexIndicator(literalValue string) bool {
	length := len(literalValue)
	// There should be at least 3 characters to form hex literal.
	if length < 3 {
		return false
	}
	// Check whether hex prefix is with positive and negative inputs.
	firstChar := literalValue[1]
	secondChar := literalValue[2]
	return firstChar == 'x' || firstChar == 'X' || secondChar == 'x' || secondChar == 'X'
}

// IsDecimalDiscriminated checks if numeric literal has decimal discriminator (d/D suffix)
// migrated from NumericLiteralSupport.java:110:5
func IsDecimalDiscriminated(literalValue string) bool {
	length := len(literalValue)
	// There should be at least 2 characters to form discriminated decimal literal.
	if length < 2 {
		return false
	}
	lastChar := literalValue[length-1]
	hasDecimalSuffix := (lastChar == 'd' || lastChar == 'D')
	if !hasDecimalSuffix {
		return false
	}
	// Check if it's not a hex literal
	return !HasHexIndicator(literalValue)
}
