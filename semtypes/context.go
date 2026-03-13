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
	Env() Env
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
	mappingMemo() map[string]*BddMemo
	functionMemo() map[string]*BddMemo
	listMemo() map[string]*BddMemo
	FunctionAtomType(atom Atom) *FunctionAtomicType
	ListAtomType(atom Atom) *ListAtomicType
	MappingAtomType(atom Atom) *MappingAtomicType
	comparableMemo(t1, t2 SemType) *comparableMemo
	setComparableMemo(t1, t2 SemType, memo *comparableMemo)
}

var _ Context = &contextImpl{}

type contextImpl struct {
	_env          Env
	_memoStack    []*BddMemo
	_listMemo     map[string]*BddMemo
	_mappingMemo  map[string]*BddMemo
	_functionMemo map[string]*BddMemo

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

func (c *contextImpl) pushToMemoStack(m *BddMemo) {
	c._memoStack = append(c._memoStack, m)
}

func (c *contextImpl) getMemoStackDepth() int {
	return len(c._memoStack)
}

func (c *contextImpl) getMemoStack(i int) *BddMemo {
	return c._memoStack[i]
}

func (c *contextImpl) popFromMemoStack() *BddMemo {
	lastIndex := len(c._memoStack) - 1
	memo := c._memoStack[lastIndex]
	c._memoStack = c._memoStack[:lastIndex]
	return memo
}

func (c *contextImpl) Env() Env {
	return c._env
}

func (c *contextImpl) jsonMemo() SemType {
	return c._jsonMemo
}

func (c *contextImpl) setJsonMemo(t SemType) {
	c._jsonMemo = t
}

func (c *contextImpl) anydataMemo() SemType {
	return c._anydataMemo
}

func (c *contextImpl) setAnydataMemo(t SemType) {
	c._anydataMemo = t
}

func (c *contextImpl) cloneableMemo() SemType {
	return c._cloneableMemo
}

func (c *contextImpl) setCloneableMemo(t SemType) {
	c._cloneableMemo = t
}

func (c *contextImpl) isolatedObjectMemo() SemType {
	return c._isolatedObjectMemo
}

func (c *contextImpl) setIsolatedObjectMemo(t SemType) {
	c._isolatedObjectMemo = t
}

func (c *contextImpl) serviceObjectMemo() SemType {
	return c._serviceObjectMemo
}

func (c *contextImpl) setServiceObjectMemo(t SemType) {
	c._serviceObjectMemo = t
}

func (c *contextImpl) mappingMemo() map[string]*BddMemo {
	return c._mappingMemo
}

func (c *contextImpl) functionMemo() map[string]*BddMemo {
	return c._functionMemo
}

func (c *contextImpl) listMemo() map[string]*BddMemo {
	return c._listMemo
}

func (c *contextImpl) FunctionAtomType(atom Atom) *FunctionAtomicType {
	return c._env.functionAtomType(atom)
}

func (c *contextImpl) ListAtomType(atom Atom) *ListAtomicType {
	return c._env.listAtomType(atom)
}

func (c *contextImpl) MappingAtomType(atom Atom) *MappingAtomicType {
	return c._env.mappingAtomType(atom)
}

func ContextFrom(env Env) Context {
	return &contextImpl{
		_env:            env,
		_listMemo:       make(map[string]*BddMemo),
		_mappingMemo:    make(map[string]*BddMemo),
		_functionMemo:   make(map[string]*BddMemo),
		_comparableMemo: make(map[comparableMemoKey]*comparableMemo),
	}
}

func (c *contextImpl) comparableMemo(t1, t2 SemType) *comparableMemo {
	return c._comparableMemo[comparableMemoKey{semType1: t1, semType2: t2}]
}

func (c *contextImpl) setComparableMemo(t1, t2 SemType, memo *comparableMemo) {
	c._comparableMemo[comparableMemoKey{semType1: t1, semType2: t2}] = memo
}
