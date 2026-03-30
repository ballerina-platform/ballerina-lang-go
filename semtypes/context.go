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

type Context = *context

type context struct {
	_env          Env
	_memoStack    []*bddMemo
	_listMemo     map[string]*bddMemo
	_mappingMemo  map[string]*bddMemo
	_functionMemo map[string]*bddMemo

	_jsonMemo           SemType
	_anydataMemo        SemType
	_cloneableMemo      SemType
	_isolatedObjectMemo SemType
	_serviceObjectMemo  SemType
	_comparableMemo     map[comparableMemoKey]*comparableMemo
}

type comparableMemo struct {
	semType1   SemType
	semType2   SemType
	comparable bool
}

type comparableMemoKey struct {
	semType1 SemType
	semType2 SemType
}

func (this *context) pushToMemoStack(m *bddMemo) {
	this._memoStack = append(this._memoStack, m)
}

func (this *context) getMemoStackDepth() int {
	return len(this._memoStack)
}

func (this *context) getMemoStack(i int) *bddMemo {
	return this._memoStack[i]
}

func (this *context) popFromMemoStack() *bddMemo {
	lastIndex := len(this._memoStack) - 1
	memo := this._memoStack[lastIndex]
	this._memoStack = this._memoStack[:lastIndex]
	return memo
}

func (this *context) Env() Env {
	return this._env
}

func (this *context) jsonMemo() SemType {
	return this._jsonMemo
}

func (this *context) setJsonMemo(t SemType) {
	this._jsonMemo = t
}

func (this *context) anydataMemo() SemType {
	return this._anydataMemo
}

func (this *context) setAnydataMemo(t SemType) {
	this._anydataMemo = t
}

func (this *context) cloneableMemo() SemType {
	return this._cloneableMemo
}

func (this *context) setCloneableMemo(t SemType) {
	this._cloneableMemo = t
}

func (this *context) isolatedObjectMemo() SemType {
	return this._isolatedObjectMemo
}

func (this *context) setIsolatedObjectMemo(t SemType) {
	this._isolatedObjectMemo = t
}

func (this *context) serviceObjectMemo() SemType {
	return this._serviceObjectMemo
}

func (this *context) setServiceObjectMemo(t SemType) {
	this._serviceObjectMemo = t
}

func (this *context) mappingMemo() map[string]*bddMemo {
	return this._mappingMemo
}

func (this *context) functionMemo() map[string]*bddMemo {
	return this._functionMemo
}

func (this *context) listMemo() map[string]*bddMemo {
	return this._listMemo
}

func (this *context) FunctionAtomType(atom Atom) *functionAtomicType {
	return this._env.functionAtomType(atom)
}

func (this *context) ListAtomType(atom Atom) *ListAtomicType {
	return this._env.listAtomType(atom)
}

func (this *context) MappingAtomType(atom Atom) *MappingAtomicType {
	return this._env.mappingAtomType(atom)
}

func ContextFrom(env Env) Context {
	return &context{
		_env:            env,
		_listMemo:       make(map[string]*bddMemo),
		_mappingMemo:    make(map[string]*bddMemo),
		_functionMemo:   make(map[string]*bddMemo),
		_comparableMemo: make(map[comparableMemoKey]*comparableMemo),
	}
}

func (this *context) comparableMemo(t1, t2 SemType) *comparableMemo {
	return this._comparableMemo[comparableMemoKey{semType1: t1, semType2: t2}]
}

func (this *context) setComparableMemo(t1, t2 SemType, memo *comparableMemo) {
	this._comparableMemo[comparableMemoKey{semType1: t1, semType2: t2}] = memo
}
