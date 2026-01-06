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

type Context interface {
	pushToMemoStack(m *BddMemo)
	getMemoStackDepth() int
	getMemoStack(i int) *BddMemo
	popFromMemoStack() *BddMemo
	env() Env
	jsonMemo() SemType
	setJsonMemo(t SemType)
	anydataMemo() SemType
	setAnydataMemo(t SemType)
	cloneableMemo() SemType
	setCloneableMemo(t SemType)
	isolatedObjectMemo() SemType
	setIsolatedObjectMemo(t SemType)
	serviceObjectMemo() SemType
	setServiceObjectMemo(t SemType)
	mappingMemo() map[Bdd]*BddMemo
	functionMemo() map[Bdd]*BddMemo
	listMemo() map[Bdd]*BddMemo
	functionAtomType(atom Atom) *FunctionAtomicType
	listAtomType(atom Atom) *ListAtomicType
	mappingAtomType(atom Atom) *MappingAtomicType
}

func ContextFrom(env Env) Context {
	panic("not implemented")
}
