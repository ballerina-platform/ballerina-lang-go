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

package elements

import (
	"strings"

	"ballerina-lang-go/compiler/util"
)

type PackageID struct {
	OrgName        util.Name
	PkgName        util.Name
	Name           util.Name
	Version        util.Name
	NameComps      []util.Name
	IsUnnamed      bool
	SkipTests      bool
	IsTestPkg      bool
	SourceFileName util.Name
	SourceRoot     string
}

func NewPackageIDWithOrgNameCompsVersion(orgName util.Name, nameComps []util.Name, version util.Name) PackageID {
	nameParts := make([]string, len(nameComps))
	for i, comp := range nameComps {
		nameParts[i] = comp.GetValue()
	}
	name := util.NewName(strings.Join(nameParts, "."))
	return PackageID{
		OrgName:   orgName,
		NameComps: nameComps,
		Name:      name,
		PkgName:   name,
		Version:   version,
	}
}

func NewPackageIDWithOrgPkgNameVersionSourceFileName(orgName util.Name, pkgName util.Name, name util.Name, version util.Name, sourceFileName util.Name) PackageID {
	return PackageID{
		OrgName:        orgName,
		Name:           name,
		PkgName:        pkgName,
		Version:        version,
		NameComps:      createNameComps(name),
		SourceFileName: sourceFileName,
	}
}

func NewPackageIDWithOrgNameVersion(orgName util.Name, name util.Name, version util.Name) PackageID {
	return PackageID{
		OrgName:   orgName,
		Name:      name,
		PkgName:   name,
		Version:   version,
		NameComps: createNameComps(name),
	}
}

func NewPackageIDWithOrgNameVersionSourceFileName(orgName util.Name, name util.Name, version util.Name, sourceFileName util.Name) PackageID {
	pkg := NewPackageIDWithOrgNameVersion(orgName, name, version)
	pkg.SourceFileName = sourceFileName
	return pkg
}

func NewPackageID(orgName util.Name, pkgName util.Name, name util.Name, version util.Name, sourceFileName util.Name, sourceRoot string, isTestPkg bool, skipTest bool) PackageID {
	return PackageID{
		OrgName:        orgName,
		Name:           name,
		PkgName:        pkgName,
		Version:        version,
		NameComps:      createNameComps(name),
		SourceFileName: sourceFileName,
		SourceRoot:     sourceRoot,
		IsTestPkg:      isTestPkg,
		SkipTests:      skipTest,
	}
}

func NewPackageIDWithOrgSourceFileNameVersion(orgName util.Name, sourceFileName string, version util.Name) PackageID {
	return PackageID{
		OrgName:        orgName,
		Name:           util.NewName("."),
		PkgName:        util.NewName("."),
		Version:        version,
		NameComps:      []util.Name{util.NewName(".")},
		IsUnnamed:      true,
		SourceFileName: util.NewName(sourceFileName),
	}
}

func NewPackageIDWithSourceFileName(sourceFileName string) PackageID {
	return PackageID{
		OrgName:        util.NewName(""),
		Name:           util.NewName("."),
		PkgName:        util.NewName("."),
		NameComps:      []util.Name{util.NewName(".")},
		IsUnnamed:      true,
		SourceFileName: util.NewName(sourceFileName),
		Version:        util.NewName(""),
	}
}

func createNameComps(name util.Name) []util.Name {
	if name.GetValue() == "." {
		return []util.Name{util.NewName(".")}
	}
	parts := strings.Split(name.GetValue(), ".")
	result := make([]util.Name, len(parts))
	for i, part := range parts {
		result[i] = util.NewName(part)
	}
	return result
}

func (p PackageID) GetPkgName() util.Name {
	return p.PkgName
}

func (p PackageID) GetName() util.Name {
	return p.Name
}

func (p PackageID) GetNameComp(index int) util.Name {
	return p.NameComps[index]
}

func (p PackageID) GetNameComps() []util.Name {
	return p.NameComps
}

func (p PackageID) GetPackageVersion() util.Name {
	return p.Version
}

func (p PackageID) Equals(o interface{}) bool {
	if o == nil {
		return false
	}
	other, ok := o.(PackageID)
	if !ok {
		return false
	}
	samePkg := false
	if p.IsUnnamed == other.IsUnnamed {
		samePkg = (!p.IsUnnamed) || (p.SourceFileName.Equals(other.SourceFileName))
	}
	return samePkg && p.OrgName.Equals(other.OrgName) && p.PkgName.Equals(other.PkgName) && p.Name.Equals(other.Name) && p.Version.Equals(other.Version)
}

func (p PackageID) HashCode() int {
	result := 0
	if p.OrgName.GetValue() != "" {
		result = len(p.OrgName.GetValue())
	}
	result = 31*result + len(p.Name.GetValue())
	result = 31*result + len(p.Version.GetValue())
	return result
}

func (p PackageID) String() string {
	if p.Name.GetValue() == "." {
		return p.Name.GetValue()
	}
	org := ""
	if p.OrgName.GetValue() != "" && p.OrgName.GetValue() != "" {
		org = p.OrgName.GetValue() + "/"
	}
	if p.Version.GetValue() == "" {
		return org + p.Name.GetValue()
	}
	return org + p.Name.GetValue() + ":" + p.Version.GetValue()
}

func (p PackageID) GetOrgName() util.Name {
	return p.OrgName
}

func IsLangLibPackageID(packageID PackageID) bool {
	if packageID.GetOrgName().GetValue() != "ballerina" {
		return false
	}
	return len(packageID.NameComps) > 1 && packageID.NameComps[0].GetValue() == "lang" || packageID.Name.GetValue() == "jballerina.java"
}
