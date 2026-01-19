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

package elements

type Flag uint8

const (
	FLAG_PUBLIC Flag = iota
	FLAG_PRIVATE
	FLAG_REMOTE
	FLAG_TRANSACTIONAL
	FLAG_NATIVE
	FLAG_FINAL
	FLAG_ATTACHED
	FLAG_LAMBDA
	FLAG_WORKER
	FLAG_PARALLEL
	FLAG_LISTENER
	FLAG_READONLY
	FLAG_FUNCTION_FINAL
	FLAG_INTERFACE
	FLAG_REQUIRED
	FLAG_RECORD
	FLAG_ANONYMOUS
	FLAG_OPTIONAL
	FLAG_TESTABLE
	FLAG_CLIENT
	FLAG_RESOURCE
	FLAG_ISOLATED
	FLAG_SERVICE
	FLAG_CONSTANT
	FLAG_TYPE_PARAM
	FLAG_LANG_LIB
	FLAG_FORKED
	FLAG_DISTINCT
	FLAG_CLASS
	FLAG_CONFIGURABLE
	FLAG_OBJECT_CTOR
	FLAG_ENUM
	FLAG_INCLUDED
	FLAG_REQUIRED_PARAM
	FLAG_DEFAULTABLE_PARAM
	FLAG_REST_PARAM
	FLAG_FIELD
	FLAG_ANY_FUNCTION
	FLAG_NEVER_ALLOWED
	FLAG_ENUM_MEMBER
	FLAG_QUERY_LAMBDA
)
