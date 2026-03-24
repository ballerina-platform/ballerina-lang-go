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

package values

import "ballerina-lang-go/semtypes"

type Object struct {
	Type       semtypes.SemType
	fields     map[string]BalValue
	methodKeys map[string]string
}

func NewObject(typ semtypes.SemType, fieldValues map[string]BalValue, methodKeys map[string]string) *Object {
	if fieldValues == nil {
		fieldValues = make(map[string]BalValue)
	}
	if methodKeys == nil {
		methodKeys = make(map[string]string)
	}
	return &Object{
		Type:       typ,
		fields:     fieldValues,
		methodKeys: methodKeys,
	}
}

func (o *Object) Put(field string, value BalValue) {
	o.fields[field] = value
}

func (o *Object) Get(field string) (BalValue, bool) {
	value, ok := o.fields[field]
	return value, ok
}

func (o *Object) MethodLookupKey(name string) (string, bool) {
	key, ok := o.methodKeys[name]
	return key, ok
}
