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

type Member struct {
	Name       string
	ValueTy    SemType
	Kind       MemberKind
	Visibility Visibility
	Immutable  bool
}

func NewMember(name string, valueTy SemType, kind MemberKind, visibility Visibility, immutable bool) *Member {
	return &Member{Name: name, ValueTy: valueTy, Kind: kind, Visibility: visibility, Immutable: immutable}
}

type memberTag interface {
	field() Field
}

type MemberKind uint8

const (
	MemberKindField MemberKind = iota
	MemberKindMethod
)

func (k *MemberKind) field() Field {
	switch *k {
	case MemberKindField:
		return Field{Name: "kind", Ty: StringConst("field"), Ro: true, Opt: false}
	case MemberKindMethod:
		return Field{Name: "kind", Ty: StringConst("method"), Ro: true, Opt: false}
	default:
		panic("invalid member kind")
	}
}

type Visibility uint8

const (
	VisibilityPublic Visibility = iota
	VisibilityPrivate
)

var (
	visibilityPublicTag  = StringConst("public")
	visibilityPrivateTag = StringConst("private")
	visibilityAll        = Field{Name: "visibility", Ty: Union(visibilityPublicTag, visibilityPrivateTag), Ro: true, Opt: false}
)

func (v *Visibility) field() Field {
	switch *v {
	case VisibilityPublic:
		return Field{Name: "visibility", Ty: visibilityPublicTag, Ro: true, Opt: false}
	case VisibilityPrivate:
		return Field{Name: "visibility", Ty: visibilityPrivateTag, Ro: true, Opt: false}
	default:
		panic("invalid visibility")
	}
}
