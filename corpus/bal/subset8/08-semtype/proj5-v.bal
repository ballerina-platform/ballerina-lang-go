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

const THREE = 3;

type Float float;

type Int int;

type ISF int|string|float;

type IF int|float;

type SF string|float;

type T02 0|2;

type T12 1|2;

// @test T[THREE] = F
// @test T[Int] = ISF
// @test T[T02] = IF
// @test T[T12] = SF
type T [int, string, float...];
