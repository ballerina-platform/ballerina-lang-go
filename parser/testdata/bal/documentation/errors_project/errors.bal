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

# Account error data.
public type AccountErrorData record {
    int accountID;
};

# All account errors.
public type AccountError distinct error<AccountErrorData>;

# Invlaid account id error.
public type InvalidAccountIdError distinct AccountError;

# Acc not found error
public type AccountNotFoundError distinct AccountError;

# Represents the total Cache error type.
public type TotalCacheError distinct CacheError;

# Represents the Cache error type with details. This will be returned if an error occurred while doing any of the cache
# operations.
public type CacheError distinct error;

# Represents the operation canceled(typically by the caller) error.
public type CancelledError distinct error;

# Represents unknown error.(e.g. Status value received is unknown)
public type UnKnownError distinct error;

# Represents Cache related errors.
public type Error CacheError;

# Represents gRPC related errors.
public type GrpcError CancelledError | UnKnownError | CacheError;

# Represents link to GrpcError.
public type LinkToGrpcError GrpcError;

# Represents union of builtin error and string.
public type YErrorType error | string;

# Represents link to YErrorType.
public type LinktoYError YErrorType;
