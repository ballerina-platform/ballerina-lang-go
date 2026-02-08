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

# Gets a access parameter value (`true` or `false`) for a given key. Please note that #foo will always be bigger than #bar.
# Example:
# ``SymbolEnv pkgEnv = symbolEnter.packageEnvs.get(pkgNode.symbol);``
# + accessMode - read or write mode
# + return - success or not
public function open (string accessMode) returns (boolean) {
    return true;
}

# Represents a Person type in ballerina.
# + name - name of the person.
public class Person {

    private string name;

    # get the users name.
    # + val - integer value
    function getName(int val) {

    }

    # Indicate whether this is a male or female.
    # + return - True if male
    public function isMale() returns boolean {
        return true;
    }
}

# Documentation for Test annotation
# + a - 'field a' documentation
# + b - 'field b' documentation
# + c - 'field c' documentation
public type Tst record {
    string a = "";
    string b = "";
    string c = "";
};

# Documentation for Test annotation
public annotation Tst Test on function;
