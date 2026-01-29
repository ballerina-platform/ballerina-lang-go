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

type CommonOps interface {
	CommonBasicTypeOps
	Union(t1 SubtypeData, t2 SubtypeData) SubtypeData
	Intersect(t1 SubtypeData, t2 SubtypeData) SubtypeData
	Diff(t1 SubtypeData, t2 SubtypeData) SubtypeData
	Complement(t SubtypeData) SubtypeData
}

type CommonOpsBase struct {
	CommonOpsMethods
}

type CommonOpsMethods struct {
	Self CommonOps
}

func (m *CommonOpsMethods) Union(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from CommonOps.java:30:5
	return BddUnion(t1.(Bdd), t2.(Bdd))
}

func (m *CommonOpsMethods) Intersect(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from CommonOps.java:35:5
	return BddIntersect(t1.(Bdd), t2.(Bdd))
}

func (m *CommonOpsMethods) Diff(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from CommonOps.java:40:5
	return BddDiff(t1.(Bdd), t2.(Bdd))
}

func (m *CommonOpsMethods) Complement(t SubtypeData) SubtypeData {
	// migrated from CommonOps.java:45:5
	return BddComplement(t.(Bdd))
}
