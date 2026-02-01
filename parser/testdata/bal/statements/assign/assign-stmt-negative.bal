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

function testIncompatibleTypeAssign () {
    boolean b;
    b = 12;
}

function testAssignCountMismatch1 () returns [int, string, int] {
    int a;
    string name;
    int b = 0;

    [a, name] = testMultiReturnValid();
    return [a, name, b];
}

function testAssignCountMismatch2 () returns [int, string, int] {
    int a;
    string name;
    int b;
    int c;

    [a, name, b, c] = testMultiReturnValid();
    return [a, name, b];
}

function testAssignTypeMismatch1 () returns [int, string, int] {
    int a;
    string name;
    int b;

    [a, name, b] = testMultiReturnInvalid();
    return [a, name, b];
}

function testMultiReturnInvalid () returns [string, string, int] {
    return [5, "john", 6];
}

function testAssignTypeMismatch2 () returns [int, string, int] {
    int a;
    int name;
    int b;

    [a, name, b] = testMultiReturnValid();
    return [a, name, b];
}

function testVarRepeatedReturn1 () returns [int, string, int] {
    var [a, name, a] = testMultiReturnValid();
    return [a, name, b];
}

function testVarRepeatedReturn2 () returns [int, string, int] {
    var [a, name, name] = testMultiReturnValid();
    return [a, name, b];
}

function testMultiReturnValid () returns [int, string, int] {
    return [5, "john", 6];
}

public client class Client {
    remote function foo() returns [int, int] {
        return [0, 0];
    }

    remote function foo1() returns record { string a; } {
        return { a: "a" };
    }

    remote function foo2() returns error {
        return error("the error reason");
    }
}

public function restActionResultAssignment() {
    Client c = new();
    int a = c->foo();
    map<string> sm = c->foo1();
    var { a: d } = c->foo2();
}


function assignErrorArrayToAnyTypeArrayViseVersa() {
    error[] ea = [];
    any[] j = ea;
    any[] anyArray = [];
    error[] errorArray = anyArray;
}

public const C_ERROR = "CError";
public const L_ERROR = "LError";

public type Detail record {
    string message;
    error cause?;
};

type CError distinct error<Detail>;
type LError distinct error<Detail>;
type CLError CError|LError;

function nonAssingableErrorTypeArrayAssign() {
    CLError?[] err = [];
    error?[] errs = err;
    ProcessErrors(errs);
    err = errs;
}

function ProcessErrors(CLError?[] errors) {}

function assignErrorArrayToUnionWithError() {
    error e1 = error("E1");
    error[] x = [e1];
    error|int[] y = x;
}

function assignErrorToUnionWithErrorArray() {
    error e1 = error("E1");
    int|error[] y = e1;
}

function assignFunctionParameterAnyToParameterUnionWithErrorAndAny() {
    function (any|error...) returns () func = function (any... y) {};
}

type Type A|Tuple|List;
type A "A";
type Tuple ["tuple", Type, Type...];
type List ["list", int?, Type];

function incompatibilityAssignInTupleTypes() {
    Type a = ["tuple", "A", "A", "A"];
    List b = a; // error: incompatible types: expected '["list",int?,Type]', found 'Type'
}

function assignTableCtrToIncompatibleType() {
    record{int a;} b = {a: 1};
    b.a = table[]; // error

    [int, string] a = [1, "a"];
    a[0] = table[]; // error
}

type Topt record {
    int a;
    record {
        int b?;
    }[1] c?;
};

function testOptionalFieldAssignment() {
    Topt t = {a: 2, c: [{b: 4}]};
    (t.c)[0].b = (); // error
}
