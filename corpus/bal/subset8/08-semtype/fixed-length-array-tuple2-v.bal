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

type IS int|string;

type ISArray IS[];

// @type ALL1 = ISArray
// @type ALL2 = ISArray
// @type ALL1 = ALL2
type ALL1 []|[IS, IS...];

type ALL2 IS[0]|[IS, IS...];

type ISTuple [IS, IS...];

// @type ALL3 = ISTuple
// @type ALL4 = ISTuple
// @type ALL3 = ALL4
type ALL3 [IS]|[IS, IS, IS...];

type ALL4 IS[1]|[IS, IS, IS...];
