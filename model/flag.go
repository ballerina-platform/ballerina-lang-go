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

package model

// PR-TODO: standardize these with AST
type Flag uint

const (
	Flag_PUBLIC Flag = iota
	Flag_PRIVATE
	Flag_REMOTE
	Flag_TRANSACTIONAL
	Flag_NATIVE
	Flag_FINAL
	Flag_ATTACHED
	Flag_LAMBDA
	Flag_WORKER
	Flag_PARALLEL
	Flag_LISTENER
	Flag_READONLY
	Flag_FUNCTION_FINAL
	Flag_INTERFACE
	Flag_REQUIRED
	Flag_RECORD
	Flag_ANONYMOUS
	Flag_OPTIONAL
	Flag_TESTABLE
	Flag_CLIENT
	Flag_RESOURCE
	Flag_ISOLATED
	Flag_SERVICE
	Flag_CONSTANT
	Flag_TYPE_PARAM
	Flag_LANG_LIB
	Flag_FORKED
	Flag_DISTINCT
	Flag_CLASS
	Flag_CONFIGURABLE
	Flag_OBJECT_CTOR
	Flag_ENUM
	Flag_INCLUDED
	Flag_REQUIRED_PARAM
	Flag_DEFAULTABLE_PARAM
	Flag_REST_PARAM
	Flag_FIELD
	Flag_ANY_FUNCTION
	Flag_NEVER_ALLOWED
	Flag_ENUM_MEMBER
	Flag_QUERY_LAMBDA
)
