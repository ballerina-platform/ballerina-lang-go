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

type BasicTypeOpsPanicImpl struct {
}

var _ BasicTypeOps = &BasicTypeOpsPanicImpl{}

func NewBasicTypeOpsPanicImpl() BasicTypeOpsPanicImpl {
	this := BasicTypeOpsPanicImpl{}
	return this
}

func (this *BasicTypeOpsPanicImpl) Union(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from BasicTypeOpsPanicImpl.java:30:5
	panic("Binary operation should not be called")
}

func (this *BasicTypeOpsPanicImpl) Intersect(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from BasicTypeOpsPanicImpl.java:35:5
	panic("Binary operation should not be called")
}

func (this *BasicTypeOpsPanicImpl) Diff(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from BasicTypeOpsPanicImpl.java:40:5
	panic("Binary operation should not be called")
}

func (this *BasicTypeOpsPanicImpl) Complement(t SubtypeData) SubtypeData {
	// migrated from BasicTypeOpsPanicImpl.java:45:5
	panic("Unary operation should not be called")
}

func (this *BasicTypeOpsPanicImpl) IsEmpty(cx Context, t SubtypeData) bool {
	// migrated from BasicTypeOpsPanicImpl.java:50:5
	panic("Unary boolean operation should not be called")
}
