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

type IntInfo record {|
    int value;
|};

type DecimalInfo record {|
    decimal value;
|};

annotation IntInfo intInfo on type;
annotation DecimalInfo decimalInfo on type;

@intInfo {value: -(-9223372036854775807 - 1)} // @error
type UnaryOverflow int;

@intInfo {value: 9223372036854775807 + 1} // @error
type AddOverflow int;

@intInfo {value: -9223372036854775807 - 2} // @error
type SubOverflow int;

@intInfo {value: 9223372036854775807 * 2} // @error
type MulOverflow int;

@intInfo {value: 10 / 0} // @error
type DivisionByZero int;

@intInfo {value: (-9223372036854775807 - 1) / -1} // @error
type DivisionOverflow int;

@intInfo {value: 10 % 0} // @error
type RemainderByZero int;

@decimalInfo {value: 9.999999999999999999999999999999999E6144d * 2d} // @error
type DecimalOverflow int;

@intInfo {value: <int>(1.0 / 0.0)} // @error
type InfiniteFloatConversion int;

@intInfo {value: <int>1E100d} // @error
type LargeDecimalConversion int;
