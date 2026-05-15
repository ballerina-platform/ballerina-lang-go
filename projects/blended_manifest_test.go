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

// Package projects internal tests for blendedManifest.
// TestDepResFix_TransitiveLocalPin_HonoredViaBFS will land here after the
// local-repo branch rebase, where custom-repo dispatch in the resolver is
// fully wired.

package projects

import "testing"

// TestBlendedManifest_Dependency exercises blendedManifest.dependency in
// isolation.
func TestBlendedManifest_Dependency(t *testing.T) {
	// Build a root manifest with a known set of [[dependency]] entries.
	v100, err := NewPackageVersionFromString("1.0.0")
	if err != nil {
		t.Fatalf("NewPackageVersionFromString: %v", err)
	}
	v200, err := NewPackageVersionFromString("2.0.0")
	if err != nil {
		t.Fatalf("NewPackageVersionFromString: %v", err)
	}

	rootDesc := NewPackageDescriptor(
		NewPackageOrg("consumerorg"),
		NewPackageName("consumerpkg"),
		v100,
	)

	// Deps: localpkg with user-specified repository "local", plainpkg with no repository.
	deps := []Dependency{
		NewDependencyWithRepository(
			NewPackageOrg("myorg"),
			NewPackageName("localpkg"),
			v100,
			"local",
		),
		NewDependency(
			NewPackageOrg("myorg"),
			NewPackageName("plainpkg"),
			v200,
		),
	}

	rootManifest := NewPackageManifestFromParams(PackageManifestParams{
		PackageDesc:  rootDesc,
		Dependencies: deps,
	})

	// Empty manifest for the "no deps" subtests.
	emptyDesc := NewPackageDescriptor(
		NewPackageOrg("emptyorg"),
		NewPackageName("emptypkg"),
		v100,
	)
	emptyManifest := NewPackageManifestFromParams(PackageManifestParams{
		PackageDesc: emptyDesc,
	})

	tests := []struct {
		name     string
		manifest PackageManifest
		org      string
		pkgName  string
		wantOk   bool
		wantRepo string
	}{
		{
			name:     "present with repository=local",
			manifest: rootManifest,
			org:      "myorg",
			pkgName:  "localpkg",
			wantOk:   true,
			wantRepo: "local",
		},
		{
			name:     "present with no repository field",
			manifest: rootManifest,
			org:      "myorg",
			pkgName:  "plainpkg",
			wantOk:   true,
			wantRepo: "",
		},
		{
			name:     "absent from root manifest",
			manifest: rootManifest,
			org:      "myorg",
			pkgName:  "unknownpkg",
			wantOk:   false,
			wantRepo: "",
		},
		{
			name:     "org match but name mismatch",
			manifest: rootManifest,
			org:      "myorg",
			pkgName:  "localpkgX",
			wantOk:   false,
			wantRepo: "",
		},
		{
			name:     "empty root manifest",
			manifest: emptyManifest,
			org:      "myorg",
			pkgName:  "anypkg",
			wantOk:   false,
			wantRepo: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := newBlendedManifest(withPackageManifest(tt.manifest))

			got, ok := bm.dependency(tt.org, tt.pkgName)
			if ok != tt.wantOk {
				t.Errorf("dependency(%q, %q) ok = %v, want %v", tt.org, tt.pkgName, ok, tt.wantOk)
			}
			if ok && got.Repository() != tt.wantRepo {
				t.Errorf("dependency(%q, %q) Repository() = %q, want %q",
					tt.org, tt.pkgName, got.Repository(), tt.wantRepo)
			}
		})
	}
}

// TestBlendedManifest_NilReceiver verifies the nil-guard on dependency():
// a nil *blendedManifest must return (zero, false) without panicking.
func TestBlendedManifest_NilReceiver(t *testing.T) {
	var bm *blendedManifest
	got, ok := bm.dependency("anyorg", "anypkg")
	if ok {
		t.Errorf("nil receiver: dependency() ok = true, want false")
	}
	if got != (blendedDependency{}) {
		t.Errorf("nil receiver: dependency() returned non-zero value: %+v", got)
	}
}
