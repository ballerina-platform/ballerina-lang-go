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



# Documentation for Test annotation
#
# + a - `field a` documentation
# + a - `field a` documentation
# + b - `field b` documentation
# + c - `field c` documentation
type Tst record {
    string a = "";
    string b = "";
    string cd = "";
};

annotation Tst Test;

# Documentation for testConst constant
#
# + testConst - abc description
final string testConst = "TestConstantDocumentation";

# Documentation for Test struct
#
# + a - struct `field a` documentation
# + a - struct `field a` documentation
# + b - struct `field b` documentation
# + c - struct `field c` documentation
class Test {
    public int a = 0;
    public int b = 0;
    public int cdd = 0;
}

# Documentation for File object
#
# + path - file path. Example: ``C:\users\OddThinking\Documents\My Source\Widget\foo.src``
public class File {

    public string path = "";

    # Gets a access parameter value (`true` or `false`) for a given key. Please note that #foo will always be bigger than #bar.
    # Example:
    # ``SymbolEnv pkgEnv = symbolEnter.packageEnvs.get(pkgNode.symbol);``
    #
    # + accessMode - read or write mode
    # + accessMode - read or write mode
    # + successfuls - boolean `true` or `false`
    public function open(string accessMode) returns boolean {
        boolean successful = false;
        return successful;
    }
}

class Person {
    public string firstName = "";
    public string lastName = "";
    public int age = 0;
    public string city = "";
}

class Employee {
    public string name = "";
    public int age = 0;
    public string address = "";
    public any ageAny = ();
}

# Test Connector
#
# + url - url for endpoint
# + url - url for endpoint
# + urls - urls for endpoint
class TestConnector {
  public string url = "";
}

# PizzaService HTTP Service
#
# + conn - HTTP connection.
service "PizzaService" on new Listener() {

    # Check orderPizza resource.
    #
    # + req - In request.
    # + req - In request.
    # + reqest - In request.
    # + conn - Listener.

    resource function get pizza(Listener conn, map<json> req) {


    }
}

# Documentation for testConst constant
#
# + testConstd - abc description
final string testConsts = "TestConstantDocumentation";

# Docs for this record.
# + abc - the field documented at object level. Known type `Bar`, unknown type `Baz`.
public class Foo {

    # The field re-documented at field level.
    # Known type `Bar`, unknown type `Baz`.
    public string abc = "str";

    // The undocumented field.
    public int def = 1;

    # The field documented at field level.
    # Known type `Bar`, unknown type `Baz`.
    string ghi = "str2";
}

type Bar record {
};

public class Listener {
    public function attach(service object {} s, string[]|string? name = ()) returns error? {
        return;
    }

    public function detach(service object {} s) returns error? {
        return;
    }

    public function 'start() returns error? {
        return;
    }

    public function gracefulStop() returns error? {
        return;
    }

    public function immediateStop() returns error? {
        return;
    }
}
