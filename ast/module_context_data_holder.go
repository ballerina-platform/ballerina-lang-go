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

package ast

// ModuleContextDataHolder is a data holder struct for ModuleContext.
// Migrated from io.ballerina.projects.internal.ModuleContextDataHolder
type ModuleContextDataHolder struct {
	Exported               bool
	ProjectKind            ProjectKind
	SkipTests              bool
	SourceRoot             string
	IsObservabiltyIncluded bool
	DumpBir                bool
	Cloud                  string
	Descriptor             ModuleDescriptor
}

// ModuleDescriptor uniquely describes a Ballerina module in terms of its name and PackageDescriptor.
// Migrated from io.ballerina.projects.ModuleDescriptor
// Note: ModuleName, PackageDescriptor, PackageName, PackageOrg, and PackageVersion types need to be migrated
type ModuleDescriptor struct {
	moduleName              *ModuleName
	packageDesc             *PackageDescriptor
	moduleCompilationId     *PackageID
	moduleTestCompilationId *PackageID
}

// newModuleDescriptor creates a new ModuleDescriptor (private constructor from Java)
// Migrated from ModuleDescriptor.java:37:5
func newModuleDescriptor(moduleName *ModuleName, packageDesc *PackageDescriptor) *ModuleDescriptor {
	this := &ModuleDescriptor{}
	this.moduleName = moduleName
	this.packageDesc = packageDesc

	if packageDesc.name().value() == "." && packageDesc.org().anonymous() {
		this.moduleCompilationId = &DEFAULT
		this.moduleTestCompilationId = this.moduleCompilationId
	} else {
		panic("Not implemented")
	}
	return this
}

// From creates a ModuleDescriptor from a ModuleName and PackageDescriptor
// Migrated from ModuleDescriptor.java:54:5
func ModuleDescriptorFrom(moduleName *ModuleName, packageDescriptor *PackageDescriptor) *ModuleDescriptor {
	return newModuleDescriptor(moduleName, packageDescriptor)
}

// PackageName returns the package name
// Migrated from ModuleDescriptor.java:58:5
func (this *ModuleDescriptor) PackageName() *PackageName {
	return this.packageDesc.name()
}

// Org returns the package organization
// Migrated from ModuleDescriptor.java:62:5
func (this *ModuleDescriptor) Org() *PackageOrg {
	return this.packageDesc.org()
}

// Version returns the package version
// Migrated from ModuleDescriptor.java:66:5
func (this *ModuleDescriptor) Version() *PackageVersion {
	return this.packageDesc.version()
}

// Name returns the module name
// Migrated from ModuleDescriptor.java:70:5
func (this *ModuleDescriptor) Name() *ModuleName {
	return this.moduleName
}

// ModuleCompilationId returns the module compilation ID
// Migrated from ModuleDescriptor.java:74:5
func (this *ModuleDescriptor) ModuleCompilationId() *PackageID {
	return this.moduleCompilationId
}

// ModuleTestCompilationId returns the module test compilation ID
// Migrated from ModuleDescriptor.java:78:5
func (this *ModuleDescriptor) ModuleTestCompilationId() *PackageID {
	return this.moduleTestCompilationId
}

// ModuleName represents the name of a Package.
// Migrated from io.ballerina.projects.ModuleName
type ModuleName struct {
	packageName    *PackageName
	moduleNamePart string
}

// newModuleName creates a new ModuleName (private constructor from Java)
// Migrated from ModuleName.java:32:5
func newModuleName(packageName *PackageName, moduleNamePart string) *ModuleName {
	this := &ModuleName{}
	this.packageName = packageName
	this.moduleNamePart = moduleNamePart
	return this
}

// FromPackageName creates the default ModuleName from a PackageName
// Migrated from ModuleName.java:37:5
func FromPackageName(packageName *PackageName) *ModuleName {
	return newModuleName(packageName, "")
}

// FromPackageNameWithModuleNamePart creates a ModuleName from a PackageName and moduleNamePart
// Migrated from ModuleName.java:42:5
func FromPackageNameWithModuleNamePart(packageName *PackageName, moduleNamePart string) *ModuleName {
	if moduleNamePart != "" && len(moduleNamePart) == 0 {
		panic("moduleNamePart should be a non-empty string or null")
	}
	// TODO Check whether the moduleNamePart is a valid list of identifiers
	return newModuleName(packageName, moduleNamePart)
}

// PackageName returns the package name
// Migrated from ModuleName.java:50:5
func (this *ModuleName) PackageName() *PackageName {
	return this.packageName
}

// ModuleNamePart returns the module name part
// Migrated from ModuleName.java:54:5
func (this *ModuleName) ModuleNamePart() string {
	return this.moduleNamePart
}

// IsDefaultModuleName checks if this is the default module name
// Migrated from ModuleName.java:58:5
func (this *ModuleName) IsDefaultModuleName() bool {
	return this.moduleNamePart == ""
}

// PackageName represents the name of a Package.
// Migrated from io.ballerina.projects.PackageName
type PackageName struct {
	packageNameStr string
}

var LANG_LIB_PACKAGE_NAME_PREFIX = "lang"

// newPackageName creates a new PackageName (private constructor from Java)
// Migrated from PackageName.java:31:5
func newPackageName(packageNameStr string) *PackageName {
	this := &PackageName{}
	this.packageNameStr = packageNameStr
	return this
}

// FromString creates a PackageName from a string
// Migrated from PackageName.java:35:5
func FromString(packageNameStr string) *PackageName {
	// TODO Check whether the packageName is a valid Ballerina identifier
	return newPackageName(packageNameStr)
}

// value returns the package name string
// Migrated from PackageName.java:54:5
func (this *PackageName) value() string {
	return this.packageNameStr
}

// PackageOrg represents the organization of a Package.
// Migrated from io.ballerina.projects.PackageOrg
type PackageOrg struct {
	packageOrgStr string
}

const (
	PKG_ORG_BALLERINA_NAME  = "ballerina"
	PKG_ORG_BALLERINAI_NAME = "ballerinai"
	PKG_ORG_BALLERINAX_NAME = "ballerinax"
	PKG_ORG_ANON_NAME       = "$anon"
)

var (
	PKG_ORG_BALLERINA  = newPackageOrg(PKG_ORG_BALLERINA_NAME)
	PKG_ORG_BALLERINAI = newPackageOrg(PKG_ORG_BALLERINAI_NAME)
	PKG_ORG_BALLERINAX = newPackageOrg(PKG_ORG_BALLERINAX_NAME)
)

// newPackageOrg creates a new PackageOrg (private constructor from Java)
// Migrated from PackageOrg.java:38:5
func newPackageOrg(packageOrgStr string) *PackageOrg {
	this := &PackageOrg{}
	this.packageOrgStr = packageOrgStr
	return this
}

// FromPackageOrg creates a PackageOrg from a string
// Migrated from PackageOrg.java:42:5
func FromPackageOrg(packageOrgStr string) *PackageOrg {
	if PKG_ORG_BALLERINA_NAME == packageOrgStr {
		return PKG_ORG_BALLERINA
	}
	// TODO Check whether the packageOrg is a valid Ballerina identifier
	return newPackageOrg(packageOrgStr)
}

// value returns the package org string
// Migrated from PackageOrg.java:51:5
func (this *PackageOrg) value() string {
	return this.packageOrgStr
}

// anonymous checks if this is an anonymous organization
// Migrated from PackageOrg.java:79:5
func (this *PackageOrg) anonymous() bool {
	return PKG_ORG_ANON_NAME == this.packageOrgStr
}

// isBallerinaOrg checks if this is the ballerina organization
// Migrated from PackageOrg.java:83:5
func (this *PackageOrg) isBallerinaOrg() bool {
	return this == PKG_ORG_BALLERINA
}

// isBallerinaxOrg checks if this is the ballerinax organization
// Migrated from PackageOrg.java:87:5
func (this *PackageOrg) isBallerinaxOrg() bool {
	return this == PKG_ORG_BALLERINAX
}

// PackageVersion represents the version of a Package.
// Migrated from io.ballerina.projects.PackageVersion
// Note: SemanticVersion type needs to be migrated
type PackageVersion struct {
	version *SemanticVersion
}

type SemanticVersion string

// BUILTIN_PACKAGE_VERSION is the version of all built-in packages
// Note: BUILTIN_PACKAGE_VERSION_STR needs to be defined
// var BUILTIN_PACKAGE_VERSION = FromPackageVersionString(BUILTIN_PACKAGE_VERSION_STR)

// newPackageVersion creates a new PackageVersion (private constructor from Java)
// Migrated from PackageVersion.java:38:5
func newPackageVersion(version *SemanticVersion) *PackageVersion {
	panic("Not implemented")
}

// FromPackageVersionString creates a PackageVersion from a string
// Migrated from PackageVersion.java:42:5
func FromPackageVersionString(versionString string) *PackageVersion {
	panic("Not implemented")
}

// FromSemanticVersion creates a PackageVersion from a SemanticVersion
// Migrated from PackageVersion.java:47:5
func FromSemanticVersion(version *SemanticVersion) *PackageVersion {
	panic("Not implemented")
}

// value returns the semantic version
// Migrated from PackageVersion.java:52:5
func (this *PackageVersion) value() *SemanticVersion {
	panic("Not implemented")
}

// PackageDescriptor uniquely describes a Ballerina package in terms of its name, organization,
// version and the loaded repository.
// Migrated from io.ballerina.projects.PackageDescriptor
type PackageDescriptor struct {
	packageName    *PackageName
	packageOrg     *PackageOrg
	packageVersion *PackageVersion
	repository     string
	isDeprecated   bool
	deprecationMsg string
}

// newPackageDescriptor creates a new PackageDescriptor (private constructor from Java)
// Migrated from PackageDescriptor.java:41:5
func newPackageDescriptor(packageOrg *PackageOrg, packageName *PackageName, packageVersion *PackageVersion, repository string) *PackageDescriptor {
	panic("Not implemented")
}

// newPackageDescriptorWithDeprecation creates a new PackageDescriptor with deprecation info (private constructor from Java)
// Migrated from PackageDescriptor.java:54:5
func newPackageDescriptorWithDeprecation(packageOrg *PackageOrg, packageName *PackageName, packageVersion *PackageVersion, repository string, isDeprecated bool, deprecationMsg string) *PackageDescriptor {
	panic("Not implemented")
}

// FromPackageOrgAndName creates a PackageDescriptor from org and name
// Migrated from PackageDescriptor.java:67:5
func FromPackageOrgAndName(packageOrg *PackageOrg, packageName *PackageName) *PackageDescriptor {
	return newPackageDescriptor(packageOrg, packageName, nil, "")
}

// FromPackageOrgNameVersion creates a PackageDescriptor from org, name, and version
// Migrated from PackageDescriptor.java:71:5
func FromPackageOrgNameVersion(packageOrg *PackageOrg, packageName *PackageName, packageVersion *PackageVersion) *PackageDescriptor {
	return newPackageDescriptor(packageOrg, packageName, packageVersion, "")
}

// FromPackageOrgNameVersionRepository creates a PackageDescriptor from org, name, version, and repository
// Migrated from PackageDescriptor.java:76:5
func FromPackageOrgNameVersionRepository(packageOrg *PackageOrg, packageName *PackageName, packageVersion *PackageVersion, repository string) *PackageDescriptor {
	return newPackageDescriptor(packageOrg, packageName, packageVersion, repository)
}

// FromPackageOrgNameVersionWithDeprecation creates a PackageDescriptor with deprecation info
// Migrated from PackageDescriptor.java:81:5
func FromPackageOrgNameVersionWithDeprecation(packageOrg *PackageOrg, packageName *PackageName, packageVersion *PackageVersion, isDeprecated bool, deprecationMsg string) *PackageDescriptor {
	return newPackageDescriptorWithDeprecation(packageOrg, packageName, packageVersion, "", isDeprecated, deprecationMsg)
}

// name returns the package name
// Migrated from PackageDescriptor.java:86:5
func (this *PackageDescriptor) name() *PackageName {
	return this.packageName
}

// org returns the package organization
// Migrated from PackageDescriptor.java:90:5
func (this *PackageDescriptor) org() *PackageOrg {
	return this.packageOrg
}

// version returns the package version
// Migrated from PackageDescriptor.java:94:5
func (this *PackageDescriptor) version() *PackageVersion {
	return this.packageVersion
}

// Repository returns the repository (Optional)
// Migrated from PackageDescriptor.java:98:5
// Note: Using common.Optional when available
func (this *PackageDescriptor) Repository() string {
	return this.repository
}

// isLangLibPackage checks if this is a lang lib package
// Migrated from PackageDescriptor.java:102:5
func (this *PackageDescriptor) isLangLibPackage() bool {
	panic("Not implemented")
}

// isBuiltInPackage checks if this is a built-in package
// Migrated from PackageDescriptor.java:106:5
func (this *PackageDescriptor) isBuiltInPackage() bool {
	panic("Not implemented")
}

// getDeprecated returns whether this package is deprecated
// Migrated from PackageDescriptor.java:110:5
func (this *PackageDescriptor) getDeprecated() bool {
	return this.isDeprecated
}

// getDeprecationMsg returns the deprecation message
// Migrated from PackageDescriptor.java:114:5
func (this *PackageDescriptor) getDeprecationMsg() string {
	return this.deprecationMsg
}
