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

package tomlparser

import (
	"os"
	"strings"
	"testing"
)

const sampleToml = `
title = "Server Configuration"

[owner]
name = "WSO2"

[database]
server = "192.168.1.1"
ports = [ 8001, 8001, 8002 ]
connection_max = 5000
enabled = true

[servers]
  [servers.alpha]
  ip = "10.0.0.1"
  dc = "eqdc10"

  [servers.beta]
  ip = "10.0.0.2"
  dc = "eqdc10"

[[routes]]
name = "health"
path = "/health"

[[routes]]
name = "metrics"
path = "/metrics"
method = "GET"
`

func TestReadString(t *testing.T) {
	toml, err := ReadString(sampleToml)
	if err != nil {
		t.Fatalf("Failed to parse TOML: %v", err)
	}

	if toml == nil {
		t.Fatal("Expected non-nil TOML object")
	}
}

func TestGet(t *testing.T) {
	toml, err := ReadString(sampleToml)
	if err != nil {
		t.Fatalf("Failed to parse TOML: %v", err)
	}

	tests := []struct {
		name     string
		key      string
		expected any
	}{
		{"root string", "title", "Server Configuration"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := toml.Get(tt.key)
			if !ok {
				t.Fatalf("Get(%q) returned ok=false; want true", tt.key)
			}
			if got != tt.expected {
				t.Errorf("Get(%q) = %v, want %v", tt.key, got, tt.expected)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	toml, err := ReadString(sampleToml)
	if err != nil {
		t.Fatalf("Failed to parse TOML: %v", err)
	}

	tests := []struct {
		key      string
		expected string
	}{
		{"owner.name", "WSO2"},
		{"servers.alpha.ip", "10.0.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, ok := toml.GetString(tt.key)
			if !ok {
				t.Fatalf("GetString(%q) returned ok=false; want true", tt.key)
			}
			if got != tt.expected {
				t.Errorf("GetString(%q) = %q, want %q", tt.key, got, tt.expected)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	toml, err := ReadString(sampleToml)
	if err != nil {
		t.Fatalf("Failed to parse TOML: %v", err)
	}

	tests := []struct {
		key      string
		expected int64
	}{
		{"database.connection_max", 5000},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, ok := toml.GetInt(tt.key)
			if !ok {
				t.Fatalf("GetInt(%q) returned ok=false; want true", tt.key)
			}
			if got != tt.expected {
				t.Errorf("GetInt(%q) = %d, want %d", tt.key, got, tt.expected)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	toml, err := ReadString(sampleToml)
	if err != nil {
		t.Fatalf("Failed to parse TOML: %v", err)
	}

	tests := []struct {
		key      string
		expected bool
	}{
		{"database.enabled", true},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, ok := toml.GetBool(tt.key)
			if !ok {
				t.Fatalf("GetBool(%q) returned ok=false; want true", tt.key)
			}
			if got != tt.expected {
				t.Errorf("GetBool(%q) = %v, want %v", tt.key, got, tt.expected)
			}
		})
	}
}

func TestGetArray(t *testing.T) {
	toml, err := ReadString(sampleToml)
	if err != nil {
		t.Fatalf("Failed to parse TOML: %v", err)
	}

	t.Run("database.ports length", func(t *testing.T) {
		ports, ok := toml.GetArray("database.ports")
		if !ok {
			t.Fatalf("GetArray(%q) returned ok=false; want true", "database.ports")
		}
		if len(ports) != 3 {
			t.Errorf("len(GetArray(%q)) = %d, want %d", "database.ports", len(ports), 3)
		}
	})
}

func TestGetTable(t *testing.T) {
	toml, err := ReadString(sampleToml)
	if err != nil {
		t.Fatalf("Failed to parse TOML: %v", err)
	}

	t.Run("servers.alpha ip", func(t *testing.T) {
		serverTable, ok := toml.GetTable("servers.alpha")
		if !ok {
			t.Fatalf("GetTable(%q) returned ok=false; want true", "servers.alpha")
		}
		ip, ok := serverTable.GetString("ip")
		if !ok {
			t.Fatalf("GetString(%q) on table returned ok=false; want true", "ip")
		}
		if ip != "10.0.0.1" {
			t.Errorf("servers.alpha.ip = %q, want %q", ip, "10.0.0.1")
		}
	})
}

func TestGetTables(t *testing.T) {
	toml, err := ReadString(sampleToml)
	if err != nil {
		t.Fatalf("Failed to parse TOML: %v", err)
	}

	routes, ok := toml.GetTables("routes")
	if !ok {
		t.Fatalf("GetTables(%q) returned ok=false; want true", "routes")
	}
	if len(routes) != 2 {
		t.Fatalf("len(GetTables(%q)) = %d, want %d", "routes", len(routes), 2)
	}

	cases := []struct {
		idx      int
		key      string
		wantStr  string
		wantInt  int64
		isString bool
	}{
		{0, "name", "health", 0, true},
		{1, "path", "/metrics", 0, true},
	}

	for _, c := range cases {
		t.Run("route-idx-"+c.key, func(t *testing.T) {
			if c.isString {
				got, ok := routes[c.idx].GetString(c.key)
				if !ok {
					t.Fatalf("GetString(%q) on routes[%d] returned ok=false; want true", c.key, c.idx)
				}
				if got != c.wantStr {
					t.Errorf("routes[%d].%s = %q, want %q", c.idx, c.key, got, c.wantStr)
				}
				return
			}
			got, ok := routes[c.idx].GetInt(c.key)
			if !ok {
				t.Fatalf("GetInt(%q) on routes[%d] returned ok=false; want true", c.key, c.idx)
			}
			if got != c.wantInt {
				t.Errorf("routes[%d].%s = %d, want %d", c.idx, c.key, got, c.wantInt)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	toml, err := ReadString(sampleToml)
	if err != nil {
		t.Fatalf("Failed to parse TOML: %v", err)
	}

	m := toml.ToMap()
	if m == nil {
		t.Fatal("ToMap() returned nil; want non-nil map")
	}
	if got := m["title"]; got != "Server Configuration" {
		t.Errorf("ToMap()[%q] = %v, want %v", "title", got, "Server Configuration")
	}
}

func TestTo(t *testing.T) {
	toml, err := ReadString(sampleToml)
	if err != nil {
		t.Fatalf("Failed to parse TOML: %v", err)
	}

	t.Run("successful unmarshal", func(t *testing.T) {
		type Config struct {
			Title string
			Owner struct {
				Name string
			}
			Database struct {
				Server        string
				Ports         []int
				ConnectionMax int `toml:"connection_max"`
				Enabled       bool
			}
		}

		var config Config
		toml.To(&config)

		if len(toml.Diagnostics()) > 0 {
			t.Errorf("To() added unexpected diagnostics: %v", toml.Diagnostics())
		}

		if config.Title != "Server Configuration" {
			t.Errorf("config.Title = %q, want %q", config.Title, "Server Configuration")
		}
		if config.Owner.Name != "WSO2" {
			t.Errorf("config.Owner.Name = %q, want %q", config.Owner.Name, "WSO2")
		}
		if config.Database.ConnectionMax != 5000 {
			t.Errorf("config.Database.ConnectionMax = %d, want %d", config.Database.ConnectionMax, 5000)
		}
	})

	t.Run("unmarshal type mismatch", func(t *testing.T) {
		freshToml, _ := ReadString(sampleToml)

		type BadConfig struct {
			Title int
		}

		var badConfig BadConfig
		freshToml.To(&badConfig)

		diagnostics := freshToml.Diagnostics()
		if len(diagnostics) == 0 {
			t.Error("To() should add diagnostics for type mismatch, but got none")
		}
	})
}

func TestReadFile(t *testing.T) {
	toml, err := Read(os.DirFS("."), "testdata/sample.toml")
	if err != nil {
		t.Fatalf("Failed to read TOML file: %v", err)
	}

	title, ok := toml.GetString("title")
	if !ok {
		t.Fatalf("GetString(%q) returned ok=false; want true", "title")
	}
	if title != "Server Configuration" {
		t.Errorf("GetString(%q) = %q, want %q", "title", title, "Server Configuration")
	}
}

func TestReadStream(t *testing.T) {
	reader := strings.NewReader(sampleToml)
	toml, err := ReadStream(reader)
	if err != nil {
		t.Fatalf("Failed to read TOML from stream: %v", err)
	}

	title, ok := toml.GetString("title")
	if !ok {
		t.Fatalf("GetString(%q) returned ok=false; want true", "title")
	}
	if title != "Server Configuration" {
		t.Errorf("GetString(%q) = %q, want %q", "title", title, "Server Configuration")
	}
}

func TestDiagnostics(t *testing.T) {
	invalidToml := `
	invalid toml syntax here
	missing = sign
	`

	toml, err := ReadString(invalidToml)
	if err == nil {
		t.Error("Expected error for invalid TOML")
	}

	if toml == nil {
		t.Fatal("Expected non-nil TOML object even with errors")
	}

	diagnostics := toml.Diagnostics()
	if len(diagnostics) == 0 {
		t.Error("Expected diagnostics for invalid TOML")
	}
}

func TestNonExistentKey(t *testing.T) {
	toml, err := ReadString(sampleToml)
	if err != nil {
		t.Fatalf("Failed to parse TOML: %v", err)
	}

	_, ok := toml.Get("nonexistent.key")
	if ok {
		t.Error("Expected false for non-existent key")
	}
}

// Consolidated essential location/diagnostics tests
func TestErrorLocation_DuplicateKeys(t *testing.T) {
	invalidToml := `
title = "Test"
title = "Duplicate"
`

	tomlDoc, err := ReadString(invalidToml)
	if err == nil {
		t.Fatal("Expected error for duplicate keys, got nil")
	}

	diags := tomlDoc.Diagnostics()
	if len(diags) == 0 {
		t.Fatal("Expected diagnostics, got none")
	}

	diag := diags[0]
	if diag.Location == nil {
		t.Fatal("Expected location information, got nil")
	}

	// The error should be on line 3 (duplicate title)
	if diag.Location.StartLine != 3 {
		t.Errorf("Expected StartLine 3, got %d", diag.Location.StartLine)
	}
	if diag.Location.StartColumn <= 0 {
		t.Errorf("Expected positive StartColumn, got %d", diag.Location.StartColumn)
	}
}

// readComprehensiveBallerinaToml parses testdata/comprehensive-ballerina.toml,
// which covers every construct manifest_builder reads from a real Ballerina.toml.
func readComprehensiveBallerinaToml(t *testing.T) *Toml {
	t.Helper()
	doc, err := Read(os.DirFS("testdata"), "comprehensive-ballerina.toml")
	if err != nil {
		t.Fatalf("failed to parse comprehensive-ballerina.toml: %v", err)
	}
	return doc
}

func TestBallerinaTomlParsesWithoutErrors(t *testing.T) {
	readComprehensiveBallerinaToml(t)
}

func TestBallerinaToml_PackageScalarFields(t *testing.T) {
	doc := readComprehensiveBallerinaToml(t)

	cases := []struct{ key, want string }{
		{"package.org", "testorg"},
		{"package.name", "my_service"},
		{"package.version", "1.2.3"},
		{"package.visibility", "private"},
		{"package.icon", "icon.png"},
		{"package.readme", "README.md"},
		{"package.distribution", "2201.15.0"},
		{"package.repository", "https://github.com/testorg/my_service"},
	}
	for _, c := range cases {
		t.Run(c.key, func(t *testing.T) {
			got, ok := doc.GetString(c.key)
			if !ok {
				t.Fatalf("GetString(%q) not found", c.key)
			}
			if got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestBallerinaToml_PackageBoolField(t *testing.T) {
	doc := readComprehensiveBallerinaToml(t)
	got, ok := doc.GetBool("package.template")
	if !ok {
		t.Fatal("package.template not found")
	}
	if got {
		t.Error("package.template should be false")
	}
}

func TestBallerinaToml_PackageStringArrays(t *testing.T) {
	doc := readComprehensiveBallerinaToml(t)

	t.Run("license", func(t *testing.T) {
		arr, ok := doc.GetArray("package.license")
		if !ok {
			t.Fatal("package.license not found")
		}
		if len(arr) != 2 || arr[0] != "Apache-2.0" || arr[1] != "MIT" {
			t.Errorf("license = %v, want [Apache-2.0 MIT]", arr)
		}
	})
	t.Run("authors", func(t *testing.T) {
		arr, ok := doc.GetArray("package.authors")
		if !ok {
			t.Fatal("package.authors not found")
		}
		if len(arr) != 2 {
			t.Errorf("len(authors) = %d, want 2", len(arr))
		}
	})
	t.Run("keywords", func(t *testing.T) {
		arr, ok := doc.GetArray("package.keywords")
		if !ok {
			t.Fatal("package.keywords not found")
		}
		if len(arr) != 4 {
			t.Errorf("len(keywords) = %d, want 4", len(arr))
		}
	})
}

func TestBallerinaToml_PackageModules(t *testing.T) {
	doc := readComprehensiveBallerinaToml(t)

	pkg, ok := doc.GetTable("package")
	if !ok {
		t.Fatal("package table not found")
	}
	modules, ok := pkg.GetTables("modules")
	if !ok {
		t.Fatal("package.modules not found")
	}
	if len(modules) != 2 {
		t.Fatalf("len(modules) = %d, want 2", len(modules))
	}

	t.Run("modules[0]", func(t *testing.T) {
		name, ok := modules[0].GetString("name")
		if !ok {
			t.Fatal("modules[0].name not found")
		}
		if name != "my_service.auth" {
			t.Errorf("got %q, want my_service.auth", name)
		}
		export, ok := modules[0].GetBool("export")
		if !ok {
			t.Fatal("modules[0].export not found")
		}
		if !export {
			t.Error("modules[0].export should be true")
		}
	})
	t.Run("modules[1]", func(t *testing.T) {
		name, _ := modules[1].GetString("name")
		if name != "my_service.db" {
			t.Errorf("got %q, want my_service.db", name)
		}
		export, _ := modules[1].GetBool("export")
		if export {
			t.Error("modules[1].export should be false")
		}
	})
}

func TestBallerinaToml_BuildOptions(t *testing.T) {
	doc := readComprehensiveBallerinaToml(t)

	boolCases := []struct {
		key  string
		want bool
	}{
		{"build-options.observabilityIncluded", true},
		{"build-options.offline", false},
		{"build-options.skipTests", false},
		{"build-options.testReport", true},
		{"build-options.codeCoverage", true},
	}
	for _, c := range boolCases {
		t.Run(c.key, func(t *testing.T) {
			got, ok := doc.GetBool(c.key)
			if !ok {
				t.Fatalf("GetBool(%q) not found", c.key)
			}
			if got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}

	t.Run("cloud", func(t *testing.T) {
		got, ok := doc.GetString("build-options.cloud")
		if !ok {
			t.Fatal("build-options.cloud not found")
		}
		if got != "k8s" {
			t.Errorf("got %q, want k8s", got)
		}
	})
}

func TestBallerinaToml_Dependencies(t *testing.T) {
	doc := readComprehensiveBallerinaToml(t)

	deps, ok := doc.GetTables("dependency")
	if !ok {
		t.Fatal("dependency not found")
	}
	if len(deps) != 2 {
		t.Fatalf("len(dependency) = %d, want 2", len(deps))
	}

	t.Run("dependency[0]", func(t *testing.T) {
		cases := []struct{ key, want string }{
			{"org", "testorg"},
			{"name", "mysql"},
			{"version", "1.5.0"},
			{"repository", "https://repo.central.ballerina.io"},
		}
		for _, c := range cases {
			got, ok := deps[0].GetString(c.key)
			if !ok {
				t.Errorf("dependency[0].%s not found", c.key)
				continue
			}
			if got != c.want {
				t.Errorf("dependency[0].%s = %q, want %q", c.key, got, c.want)
			}
		}
	})

	t.Run("dependency[1]", func(t *testing.T) {
		org, _ := deps[1].GetString("org")
		name, _ := deps[1].GetString("name")
		version, _ := deps[1].GetString("version")
		if org != "testorg" || name != "http" || version != "2.10.2" {
			t.Errorf("got {%s/%s@%s}, want {testorg/http@2.10.2}", org, name, version)
		}
		_, hasRepo := deps[1].GetString("repository")
		if hasRepo {
			t.Error("dependency[1] should not have a repository field")
		}
	})
}

func TestErrorLocation_SyntaxError(t *testing.T) {
	invalidToml := `
[section
key = "value"
`

	tomlDoc, err := ReadString(invalidToml)
	if err == nil {
		t.Fatal("Expected error for invalid syntax, got nil")
	}

	diags := tomlDoc.Diagnostics()
	if len(diags) == 0 {
		t.Fatal("Expected diagnostics, got none")
	}

	diag := diags[0]
	if diag.Location == nil {
		t.Fatal("Expected location information, got nil")
	}
	// The error is on line 2: the newline after "[section" is where ']' was expected.
	// (Previously the error landed on line 3 because skipLeadingTrivia was eating
	// the newline as trivia, making "key" the first visible token — that was wrong.)
	if diag.Location.StartLine != 2 {
		t.Errorf("Expected StartLine 2, got %d", diag.Location.StartLine)
	}
}

// TestErrorRecoveryStopsAtNewline verifies that a syntax error on one line
// does not swallow subsequent valid lines. Recovery must stop at the newline
// boundary so that the next valid key-value pair is still parsed.
// testdata/error-recovery.toml contains:
//
//	bad_key        ← missing '=', triggers skipToRecovery
//	key = "value"  ← must survive recovery
func TestErrorRecoveryStopsAtNewline(t *testing.T) {
	tomlDoc, err := Read(os.DirFS("testdata"), "error-recovery.toml")
	if err == nil {
		t.Fatal("expected a parse error for the bad_key line, got nil")
	}
	if tomlDoc == nil {
		t.Fatal("expected partial Toml even on parse error")
	}
	_, ok := tomlDoc.Get("key")
	if !ok {
		t.Error("key = \"value\" was swallowed by overshoot recovery; " +
			"skipToRecovery must stop at TokenNewline")
	}
}

