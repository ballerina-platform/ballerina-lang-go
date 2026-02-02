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


@Description{
    value:"Self annotating an annotation",
    paramValue:{
        value:"some parameter value"
    },
    queryParamValue:[{
        name:"first query param name",
        value:"first query param value"
    }],
    queryParamValue2:[{}],
    code:[7,8,9],
    args: {}
}
public annotation <service, resource, function, streamlet, struct, annotation, enum, parameter, transformer, endpoint> Description Desc;

struct Desc {
    string value = "Description of the service/function";
    int[] code;
    Param paramValue;
    QueryParam[] queryParamValue;
    QueryParam[] queryParamValue2;
    string[] paramValue2;
    Args args;
}

public struct Param {
    string value = "Description of the input param";
}

public struct QueryParam {
    string name = "default name";
    string value = "default value";
}

public struct Doc {
    Desc des;
}

public struct Args {
    string value = "default value for 'Args' annotation in doc package";
}