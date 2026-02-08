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

public const APPLICATION_ERROR_REASON = "{ballerina/sql}ApplicationError";

# Represents the properties which are related to Non SQL errors
#
# + message - Error message
# + cause - cause
public type ApplicationErrorData record {|
    string message;
    error cause?;
|};

public type ApplicationError error<ApplicationErrorData>;

public type OrderCreationError distinct ApplicationError;
public type OrderProcessingError distinct ApplicationError;
public type OrderCreationError2 distinct OrderCreationError;

public type NewPostDefinedError distinct PostDefinedError;
public type PostDefinedError error<ErrorData>;

public type ErrorData record {|
    string code;
|};

public type Data1 record {|
    int num;
|};

public type Data2 record {|
    byte num;
|};

public type Data3 record {|
    int:Signed16 num;
|};

public type ErrorIntersection1 distinct error<Data1> & error<Data2>;
public type ErrorIntersection2 distinct error<Data3> & error<Data2>;

public type StatusCodeError distinct error;
public type DefaultStatusCodeError distinct StatusCodeError & error<record { int code; }>;
