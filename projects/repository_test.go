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

	"ballerina-lang-go/projects/repository"
)

func TestLocalCacheRepository_GetVersions(t *testing.T) {
	repo := repository.NewLocalCacheRepository("testdata/cache")

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
			expected: []string{},
		},
		{
			name:     "non-existent org",
			org:      "nonexistent",
			pkg:      "http",
			expected: []string{},
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
			versions, err := repo.GetVersions(context.Background(), tt.org, tt.pkg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !slices.Equal(versions, tt.expected) {
				t.Errorf("got %v, want %v", versions, tt.expected)
			}
		})
	}
}

func TestLocalCacheRepository_GetLatestVersion(t *testing.T) {
	repo := repository.NewLocalCacheRepository("testdata/cache")

	tests := []struct {
		name     string
		org      string
		pkg      string
		expected string
	}{
		{"latest of multiple", "ballerina", "http", "2.10.1"},
		{"single version", "ballerina", "io", "1.6.0"},
		{"non-existent package", "ballerina", "nonexistent", ""},
		{"non-existent org", "nonexistent", "http", ""},
		{"testorg latest", "testorg", "testpkg", "2.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := repo.GetLatestVersion(context.Background(), tt.org, tt.pkg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if version != tt.expected {
				t.Errorf("got %q, want %q", version, tt.expected)
			}
		})
	}
}

func TestLocalCacheRepository_Exists(t *testing.T) {
	repo := repository.NewLocalCacheRepository("testdata/cache")

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

func TestLocalCacheRepository_ContextCancellation(t *testing.T) {
	repo := repository.NewLocalCacheRepository("testdata/cache")
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	t.Run("GetVersions cancelled", func(t *testing.T) {
		_, err := repo.GetVersions(ctx, "ballerina", "http")
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("GetLatestVersion cancelled", func(t *testing.T) {
		_, err := repo.GetLatestVersion(ctx, "ballerina", "http")
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

func TestLocalCacheRepository_Name(t *testing.T) {
	repo := repository.NewLocalCacheRepository("testdata/cache")
	if repo.Name() != "local-cache" {
		t.Errorf("expected 'local-cache', got %q", repo.Name())
	}
}

func TestLocalCacheRepository_BasePath(t *testing.T) {
	repo := repository.NewLocalCacheRepository("testdata/cache")
	if repo.BasePath() != "testdata/cache" {
		t.Errorf("expected 'testdata/cache', got %q", repo.BasePath())
	}
}

func TestLocalCacheRepository_ImplementsInterface(t *testing.T) {
	// Compile-time check - if this compiles, the interface is satisfied
	var _ repository.Repository = (*repository.LocalCacheRepository)(nil)
}
