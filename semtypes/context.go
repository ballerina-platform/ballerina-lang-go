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

// Context is an opaque pointer to the thread local context for type system. As such it should never be used
// concurrently in multiple type check operations at the same time. That being said it is highly advised to
// reuse the Context given it acts as a cache for many type operations.
type Context = *context

type context struct {
	_env          Env
	_memoStack    []*bddMemo
	_listMemo     map[string]*bddMemo
	_mappingMemo  map[string]*bddMemo
	_functionMemo map[string]*bddMemo

	_conjunctions []conjunction

	_jsonMemo           SemType
	_anydataMemo        SemType
	_cloneableMemo      SemType
	_isolatedObjectMemo SemType
	_serviceObjectMemo  SemType
	_isolatedTopMemo    SemType
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

func (c *context) pushToMemoStack(m *bddMemo) {
	c._memoStack = append(c._memoStack, m)
}

func (c *context) getMemoStackDepth() int {
	return len(c._memoStack)
}

func (c *context) getMemoStack(i int) *bddMemo {
	return c._memoStack[i]
}

func (c *context) popFromMemoStack() *bddMemo {
	lastIndex := len(c._memoStack) - 1
	memo := c._memoStack[lastIndex]
	c._memoStack = c._memoStack[:lastIndex]
	return memo
}

func (c *context) Env() Env {
	return c._env
}

func (c *context) jsonMemo() SemType {
	return c._jsonMemo
}

func (c *context) setJsonMemo(t SemType) {
	c._jsonMemo = t
}

func (c *context) anydataMemo() SemType {
	return c._anydataMemo
}

func (c *context) setAnydataMemo(t SemType) {
	c._anydataMemo = t
}

func (c *context) cloneableMemo() SemType {
	return c._cloneableMemo
}

func (c *context) setCloneableMemo(t SemType) {
	c._cloneableMemo = t
}

func (c *context) isolatedObjectMemo() SemType {
	return c._isolatedObjectMemo
}

func (c *context) setIsolatedObjectMemo(t SemType) {
	c._isolatedObjectMemo = t
}

func (c *context) serviceObjectMemo() SemType {
	return c._serviceObjectMemo
}

func (c *context) setServiceObjectMemo(t SemType) {
	c._serviceObjectMemo = t
}

func (c *context) mappingMemo() map[string]*bddMemo {
	return c._mappingMemo
}

func (c *context) functionMemo() map[string]*bddMemo {
	return c._functionMemo
}

func (c *context) listMemo() map[string]*bddMemo {
	return c._listMemo
}

func (c *context) FunctionAtomType(atom atom) *functionAtomicType {
	return c._env.functionAtomType(atom)
}

func (c *context) ListAtomType(atom atom) *ListAtomicType {
	return c._env.listAtomType(atom)
}

func (c *context) MappingAtomType(atom atom) *MappingAtomicType {
	return c._env.mappingAtomType(atom)
}

func (c *context) pushConjunction(atom atom, next conjunctionHandle) conjunctionHandle {
	idx := conjunctionHandle(len(c._conjunctions) + 1)
	c._conjunctions = append(c._conjunctions, conjunction{atom: atom, Next: next})
	return idx
}

func (c *context) conjunctionAtom(h conjunctionHandle) atom {
	return c._conjunctions[h-1].atom
}

func (c *context) conjunctionNext(h conjunctionHandle) conjunctionHandle {
	return c._conjunctions[h-1].Next
}

func (c *context) conjunctionStackDepth() int32 {
	return int32(len(c._conjunctions))
}

func (c *context) resetConjunctionStack(depth int32) {
	c._conjunctions = c._conjunctions[:depth]
}

func ContextFrom(env Env) Context {
	return &context{
		_env:            env,
		_listMemo:       make(map[string]*bddMemo),
		_mappingMemo:    make(map[string]*bddMemo),
		_functionMemo:   make(map[string]*bddMemo),
		_comparableMemo: make(map[comparableMemoKey]*comparableMemo),
		_conjunctions:   make([]conjunction, 0, 64),
	}
}

func (c *context) comparableMemo(t1, t2 SemType) *comparableMemo {
	return c._comparableMemo[comparableMemoKey{semType1: t1, semType2: t2}]
}

func (c *context) setComparableMemo(t1, t2 SemType, memo *comparableMemo) {
	c._comparableMemo[comparableMemoKey{semType1: t1, semType2: t2}] = memo
}
