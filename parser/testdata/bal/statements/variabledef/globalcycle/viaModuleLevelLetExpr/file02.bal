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



string letRef = modVar;

string[] modVarQueryLet1 = from var item in [queryRef, "hello", "world"] let string prefix = queryRef select prefix + item;
string[] modVarQueryLet2 = from var item in ["hello", "world"] let string prefix = queryRef2 select prefix + item;
string[] modVarQuery = from var item in [queryRef, "hello", "world"] select item;

string queryRef = modVarQuery[0] + modVarQueryLet1[0];
string queryRef2 = modVarQueryLet2[0];

RecordTypeWithDefaultLetExpr recD = {};

final int moduleCode = recD.code;
