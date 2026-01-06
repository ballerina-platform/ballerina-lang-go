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

package semtypes

type Atom interface {
	Index() int
	Kind() Kind
}

type atomIdentifier struct {
	Index int
	Kind  Kind
}

type Kind uint

const (
	Kind_LIST_ATOM Kind = iota
	Kind_FUNCTION_ATOM
	Kind_MAPPING_ATOM
	Kind_CELL_ATOM
	Kind_XML_ATOM
	Kind_DISTINCT_ATOM
)

func GetIdentifier(this Atom) atomIdentifier {
	// migrated from Atom.java:43:5

    return atomIdentifier{
        Index: this.Index(),
        Kind:  this.Kind(),
    }
}
