/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package projects

import (
	"fmt"
	"strconv"
	"strings"
)

// SemanticVersion represents a semantic version (major.minor.patch).
// Java source: io.ballerina.projects.SemanticVersion
type SemanticVersion struct {
	major      int
	minor      int
	patch      int
	preRelease string
	build      string
}

// NewSemanticVersion creates a new SemanticVersion with the given components.
func NewSemanticVersion(major, minor, patch int) SemanticVersion {
	return SemanticVersion{major: major, minor: minor, patch: patch}
}

// NewSemanticVersionWithPreRelease creates a SemanticVersion with pre-release info.
func NewSemanticVersionWithPreRelease(major, minor, patch int, preRelease string) SemanticVersion {
	return SemanticVersion{major: major, minor: minor, patch: patch, preRelease: preRelease}
}

// ParseSemanticVersion parses a version string into a SemanticVersion.
// Supports formats: "1.0.0", "1.0.0-beta", "1.0.0+build", "1.0.0-beta+build"
func ParseSemanticVersion(version string) (SemanticVersion, error) {
	if version == "" {
		return SemanticVersion{}, fmt.Errorf("version string is empty")
	}

	var preRelease, build string

	// Extract build metadata
	if idx := strings.Index(version, "+"); idx != -1 {
		build = version[idx+1:]
		version = version[:idx]
	}

	// Extract pre-release
	if idx := strings.Index(version, "-"); idx != -1 {
		preRelease = version[idx+1:]
		version = version[:idx]
	}

	// Parse major.minor.patch
	parts := strings.Split(version, ".")
	if len(parts) < 1 || len(parts) > 3 {
		return SemanticVersion{}, fmt.Errorf("invalid version format: %s", version)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return SemanticVersion{}, fmt.Errorf("invalid major version: %s", parts[0])
	}

	minor := 0
	if len(parts) >= 2 {
		minor, err = strconv.Atoi(parts[1])
		if err != nil {
			return SemanticVersion{}, fmt.Errorf("invalid minor version: %s", parts[1])
		}
	}

	patch := 0
	if len(parts) >= 3 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			return SemanticVersion{}, fmt.Errorf("invalid patch version: %s", parts[2])
		}
	}

	return SemanticVersion{
		major:      major,
		minor:      minor,
		patch:      patch,
		preRelease: preRelease,
		build:      build,
	}, nil
}

// Major returns the major version number.
func (v SemanticVersion) Major() int {
	return v.major
}

// Minor returns the minor version number.
func (v SemanticVersion) Minor() int {
	return v.minor
}

// Patch returns the patch version number.
func (v SemanticVersion) Patch() int {
	return v.patch
}

// PreRelease returns the pre-release identifier.
func (v SemanticVersion) PreRelease() string {
	return v.preRelease
}

// Build returns the build metadata.
func (v SemanticVersion) Build() string {
	return v.build
}

// IsPreRelease returns true if this is a pre-release version.
func (v SemanticVersion) IsPreRelease() bool {
	return v.preRelease != ""
}

// String returns the string representation of the version.
func (v SemanticVersion) String() string {
	result := fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	if v.preRelease != "" {
		result += "-" + v.preRelease
	}
	if v.build != "" {
		result += "+" + v.build
	}
	return result
}

// Equals checks if two SemanticVersion instances are equal.
func (v SemanticVersion) Equals(other SemanticVersion) bool {
	return v.major == other.major &&
		v.minor == other.minor &&
		v.patch == other.patch &&
		v.preRelease == other.preRelease
	// Note: build metadata is NOT considered in equality per SemVer spec
}

// Compare compares two versions. Returns -1 if v < other, 0 if equal, 1 if v > other.
func (v SemanticVersion) Compare(other SemanticVersion) int {
	if v.major != other.major {
		if v.major < other.major {
			return -1
		}
		return 1
	}
	if v.minor != other.minor {
		if v.minor < other.minor {
			return -1
		}
		return 1
	}
	if v.patch != other.patch {
		if v.patch < other.patch {
			return -1
		}
		return 1
	}
	// Pre-release versions have lower precedence
	if v.preRelease == "" && other.preRelease != "" {
		return 1
	}
	if v.preRelease != "" && other.preRelease == "" {
		return -1
	}
	if v.preRelease < other.preRelease {
		return -1
	}
	if v.preRelease > other.preRelease {
		return 1
	}
	return 0
}

// PackageVersion represents a Ballerina package version.
// Java source: io.ballerina.projects.PackageVersion
type PackageVersion struct {
	value SemanticVersion
}

// NewPackageVersion creates a new PackageVersion from a SemanticVersion.
func NewPackageVersion(value SemanticVersion) PackageVersion {
	return PackageVersion{value: value}
}

// NewPackageVersionFromString creates a PackageVersion by parsing a version string.
func NewPackageVersionFromString(version string) (PackageVersion, error) {
	sv, err := ParseSemanticVersion(version)
	if err != nil {
		return PackageVersion{}, err
	}
	return PackageVersion{value: sv}, nil
}

// Value returns the underlying SemanticVersion.
func (v PackageVersion) Value() SemanticVersion {
	return v.value
}

// String returns the string representation of the version.
func (v PackageVersion) String() string {
	return v.value.String()
}

// Major returns the major version number.
func (v PackageVersion) Major() int {
	return v.value.Major()
}

// Minor returns the minor version number.
func (v PackageVersion) Minor() int {
	return v.value.Minor()
}

// Patch returns the patch version number.
func (v PackageVersion) Patch() int {
	return v.value.Patch()
}

// Equals checks if two PackageVersion instances are equal.
func (v PackageVersion) Equals(other PackageVersion) bool {
	return v.value.Equals(other.value)
}

// Compare compares two versions. Returns -1 if v < other, 0 if equal, 1 if v > other.
func (v PackageVersion) Compare(other PackageVersion) int {
	return v.value.Compare(other.value)
}
