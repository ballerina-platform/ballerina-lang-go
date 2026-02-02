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

int 'int\ literal\ \$\ global = 23;
string 'string\ literal\ \$\$\ global = "literal";

public function main(string... args) {
    float 'float\ literal\ \$\$\$\ local = 34.43;
    decimal 'decimal\ literal\ \$\$\$\$\ local = 21.1;
    int nonLiteralInt = 5;
    float nonLiteralFloat = 5.1;
}
