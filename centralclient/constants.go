// Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
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

package centralclient

const (
	BallerinaPlatform                 = "Ballerina-Platform"
	Identity                          = "identity"
	ResolvedRequestedURI              = "RESOLVED_REQUESTED_URI"
	SSL                               = "SSL"
	Authorization                     = "Authorization"
	ContentType                       = "Content-Type"
	AcceptEncoding                    = "Accept-Encoding"
	UserAgent                         = "User-Agent"
	Location                          = "Location"
	Accept                            = "Accept"
	ContentDisposition                = "Content-Disposition"
	ApplicationOctetStream            = "application/octet-stream"
	ApplicationJSON                   = "application/json"
	BallerinaCentralTelemetryDisabled = "Ballerina-Central-Telemetry-Disabled"
	Digest                            = "digest"
)

const (
	TestModeActive        = "TEST_MODE_ACTIVE"
	BallerinaStageCentral = "BALLERINA_STAGE_CENTRAL"
	BallerinaDevCentral   = "BALLERINA_DEV_CENTRAL"
)

const (
	ConnectorsPath      = "connectors"
	TriggersPath        = "triggers"
	Separator           = "/"
	ResolveDependencies = "resolve-dependencies"
	ResolveModules      = "resolve-modules"
	PackagePathPrefix   = Separator + "packages" + Separator
	ConnectorPathPrefix = Separator + "connectors" + Separator
	TriggerPathPrefix   = Separator + "triggers" + Separator
)

const SHA256 = "sha-256="

const (
	ErrCannotFindPackage  = "error: could not connect to remote repository to find package: "
	ErrCannotFindVersions = "error: could not connect to remote repository to find versions for: "
	ErrCannotPush         = "error: failed to push the package: "
	ErrCannotPullPackage  = "error: failed to pull the package: "
	ErrCannotGetConnector = "error: failed to find connector: "
	ErrCannotGetTriggers  = "error: failed to find triggers: "
	ErrCannotGetTrigger   = "error: failed to find the trigger: "
	ErrPackageResolution  = "error: while connecting to central: "
)

const (
	DefaultConnectTimeout = 60
	DefaultReadTimeout    = 60
	DefaultWriteTimeout   = 60
	DefaultCallTimeout    = 0
	MaxRetry              = 1
	ConnectionReset       = "Connection reset"
)

const (
	MediaTypeJSON        = "application/json; charset=utf-8"
	MediaTypeJSONContent = "application/json"
)
