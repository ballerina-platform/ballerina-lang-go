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

package projects_test

import (
	"context"
	"slices"
	"testing"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/repository"
)

func TestRepository_GetPackageVersions(t *testing.T) {
	repo := repository.NewRepository("testdata/repo/bala")

	tests := []struct {
		name     string
		org      string
		pkg      string
		expected []string
	}{
		{
			name:     "multiple versions sorted",
			org:      "ballerina",
			pkg:      "http",
			expected: []string{"2.9.0", "2.10.0", "2.10.1"},
		},
		{
			name:     "single version",
			org:      "ballerina",
			pkg:      "io",
			expected: []string{"1.6.0"},
		},
		{
			name:     "non-existent package",
			org:      "ballerina",
			pkg:      "nonexistent",
			expected: nil,
		},
		{
			name:     "non-existent org",
			org:      "nonexistent",
			pkg:      "http",
			expected: nil,
		},
		{
			name:     "testorg with multiple versions",
			org:      "testorg",
			pkg:      "testpkg",
			expected: []string{"1.0.0", "1.0.1", "2.0.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			versions, err := repo.GetPackageVersions(context.Background(), tt.org, tt.pkg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			// Convert to strings for comparison
			var got []string
			for _, v := range versions {
				got = append(got, v.String())
			}
			if !slices.Equal(got, tt.expected) {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRepository_GetLatestVersion(t *testing.T) {
	repo := repository.NewRepository("testdata/repo/bala")

	tests := []struct {
		name      string
		org       string
		pkg       string
		expected  string
		wantFound bool
	}{
		{"latest of multiple", "ballerina", "http", "2.10.1", true},
		{"single version", "ballerina", "io", "1.6.0", true},
		{"non-existent package", "ballerina", "nonexistent", "", false},
		{"non-existent org", "nonexistent", "http", "", false},
		{"testorg latest", "testorg", "testpkg", "2.0.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, found, err := repo.GetLatestVersion(context.Background(), tt.org, tt.pkg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if found != tt.wantFound {
				t.Errorf("found = %v, want %v", found, tt.wantFound)
			}
			var got string
			if found {
				got = version.String()
			}
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestRepository_Exists(t *testing.T) {
	repo := repository.NewRepository("testdata/repo/bala")

	tests := []struct {
		name     string
		org      string
		pkg      string
		version  string
		expected bool
	}{
		{"exists", "ballerina", "http", "2.10.0", true},
		{"exists latest", "ballerina", "http", "2.10.1", true},
		{"version not exists", "ballerina", "http", "1.0.0", false},
		{"package not exists", "ballerina", "nonexistent", "1.0.0", false},
		{"org not exists", "nonexistent", "http", "1.0.0", false},
		{"testorg exists", "testorg", "testpkg", "1.0.1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := repo.Exists(context.Background(), tt.org, tt.pkg, tt.version)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if exists != tt.expected {
				t.Errorf("got %v, want %v", exists, tt.expected)
			}
		})
	}
}

func TestRepository_ContextCancellation(t *testing.T) {
	repo := repository.NewRepository("testdata/repo/bala")
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	t.Run("GetPackageVersions cancelled", func(t *testing.T) {
		_, err := repo.GetPackageVersions(ctx, "ballerina", "http")
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("GetLatestVersion cancelled", func(t *testing.T) {
		_, _, err := repo.GetLatestVersion(ctx, "ballerina", "http")
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("Exists cancelled", func(t *testing.T) {
		_, err := repo.Exists(ctx, "ballerina", "http", "2.10.0")
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})
}

func TestRepository_Name(t *testing.T) {
	repo := repository.NewRepository("testdata/repo/bala")
	if repo.Name() != "filesystem" {
		t.Errorf("expected 'filesystem', got %q", repo.Name())
	}
}

func TestRepository_Root(t *testing.T) {
	repo := repository.NewRepository("testdata/repo/bala")
	if repo.Root() != "testdata/repo/bala" {
		t.Errorf("expected 'testdata/repo/bala', got %q", repo.Root())
	}
}

func TestRepository_ImplementsWritableRepository(t *testing.T) {
	// Compile-time check - if this compiles, the interface is satisfied
	var _ repository.WritableRepository = (*repository.Repository)(nil)
}

func TestRepository_ImplementsProjectsRepository(t *testing.T) {
	// Compile-time check - if this compiles, Repository implements projects.Repository
	var _ projects.Repository = (*repository.Repository)(nil)
}
