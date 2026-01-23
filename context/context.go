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

package context

import (
	"ballerina-lang-go/model"
	"strconv"
)

type CompilerContext struct {
	anonTypeCount map[packageKey]int
	packageContext *packageContext
}

func (this *CompilerContext) GetDefaultPackage() *model.PackageID {
	return this.packageContext.defaultPackage
}

type packageContext struct {
	defaultPackage *model.PackageID
}

func NewCompilerContext() *CompilerContext {
	return &CompilerContext{
		anonTypeCount: make(map[packageKey]int),
        packageContext: &packageContext{
            defaultPackage: model.DEFAULT,
        },
	}
}

type packageKey struct {
	orgName model.Name
	pkgName model.Name
	name    model.Name
	version model.Name
}

func packageKeyFromPackageID(packageID *model.PackageID) packageKey {
	if packageID == nil || packageID.IsUnnamed() {
		return packageKey{
			orgName: model.ANON_ORG,
			pkgName: model.DEFAULT_PACKAGE,
			version: model.DEFAULT_VERSION,
			name:    model.DEFAULT_PACKAGE,
		}
	}
	return packageKey{
		orgName: *packageID.OrgName,
		pkgName: *packageID.PkgName,
		version: *packageID.Version,
		name:    *packageID.Name,
	}
}

const (
	ANON_PREFIX       = "$anon"
	BUILTIN_ANON_TYPE = ANON_PREFIX + "Type$builtin$"
	ANON_TYPE         = ANON_PREFIX + "Type$"
)

func (this *CompilerContext) GetNextAnonymousTypeKey(packageID *model.PackageID) string {
	packageKey := packageKeyFromPackageID(packageID)
	nextValue := this.anonTypeCount[packageKey]
	this.anonTypeCount[packageKey] = nextValue + 1
	if packageID != nil && model.ANNOTATIONS_PKG != packageID {
		return BUILTIN_ANON_TYPE + "_" + strconv.Itoa(nextValue)
	}
	return ANON_TYPE + "_" + strconv.Itoa(nextValue)
}
