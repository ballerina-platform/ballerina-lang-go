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

var _ Context = &contextImpl{}

type contextImpl struct {
	_env Env
	_memoStack []*BddMemo
	_listMemo map[Bdd]*BddMemo
	_mappingMemo map[Bdd]*BddMemo
	_functionMemo map[Bdd]*BddMemo

	_jsonMemo SemType
	_anydataMemo SemType
	_cloneableMemo SemType
	_isolatedObjectMemo SemType
	_serviceObjectMemo SemType
}

func (this *contextImpl) pushToMemoStack(m *BddMemo) {
	this._memoStack = append(this._memoStack, m)
}

func (this *contextImpl) getMemoStackDepth() int {
	return len(this._memoStack)
}

func (this *contextImpl) getMemoStack(i int) *BddMemo {
	return this._memoStack[i]
}

func (this *contextImpl) popFromMemoStack() *BddMemo {
	lastIndex := len(this._memoStack) - 1
	memo := this._memoStack[lastIndex]
	this._memoStack = this._memoStack[:lastIndex]
	return memo
}

func (this *contextImpl) env() Env {
	return this._env
}

func (this *contextImpl) jsonMemo() SemType {
	return this._jsonMemo
}

func (this *contextImpl) setJsonMemo(t SemType) {
	this._jsonMemo = t
}

func (this *contextImpl) anydataMemo() SemType {
	return this._anydataMemo
}

func (this *contextImpl) setAnydataMemo(t SemType) {
	this._anydataMemo = t
}

func (this *contextImpl) cloneableMemo() SemType {
	return this._cloneableMemo
}

func (this *contextImpl) setCloneableMemo(t SemType) {
	this._cloneableMemo = t
}

func (this *contextImpl) isolatedObjectMemo() SemType {
	return this._isolatedObjectMemo
}

func (this *contextImpl) setIsolatedObjectMemo(t SemType) {
	this._isolatedObjectMemo = t
}

func (this *contextImpl) serviceObjectMemo() SemType {
	return this._serviceObjectMemo
}

func (this *contextImpl) setServiceObjectMemo(t SemType) {
	this._serviceObjectMemo = t
}

func (this *contextImpl) mappingMemo() map[Bdd]*BddMemo {
	return this._mappingMemo
}

func (this *contextImpl) functionMemo() map[Bdd]*BddMemo {
	return this._functionMemo
}

func (this *contextImpl) listMemo() map[Bdd]*BddMemo {
	return this._listMemo
}

func (this *contextImpl) functionAtomType(atom Atom) *FunctionAtomicType {
	return this._env.functionAtomType(atom)
}

func (this *contextImpl) listAtomType(atom Atom) *ListAtomicType {
	return this._env.listAtomType(atom)
}

func (this *contextImpl) mappingAtomType(atom Atom) *MappingAtomicType {
	return this._env.mappingAtomType(atom)
}

func ContextFrom(env Env) Context {
	return &contextImpl{
		_env:          env,
		_listMemo:     make(map[Bdd]*BddMemo),
		_mappingMemo:  make(map[Bdd]*BddMemo),
		_functionMemo: make(map[Bdd]*BddMemo),
	}
}
