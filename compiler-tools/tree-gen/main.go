/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

type Attribute struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Occurrences string `json:"occurrences,omitempty"`
	IsOptional  bool   `json:"isOptional,omitempty"`
}

type NodeDef struct {
	Name       string      `json:"name"`
	Base       string      `json:"base"`
	Kind       string      `json:"kind"`
	IsAbstract bool        `json:"isAbstract,omitempty"`
	Attributes []Attribute `json:"attributes"`
}

type Config struct {
	Nodes []NodeDef `json:"nodes"`
}

func main() {
	var (
		configPath       = flag.String("config", "", "Path to JSON configuration file")
		nodeType         = flag.String("type", "", "Node type: 'st-node' or 'node'")
		templatePath     = flag.String("template", "", "Path to Go template file")
		outputPath       = flag.String("output", "", "Output file path")
		utilTemplatePath = flag.String("util-template", "", "Path to util template file (optional)")
		utilOutputPath   = flag.String("util-output", "", "Util output file path (optional)")
	)
	flag.Parse()

	if *configPath == "" || *nodeType == "" || *templatePath == "" || *outputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: all flags required (config, type, template, output)")
		flag.Usage()
		os.Exit(1)
	}

	if *nodeType != "st-node" && *nodeType != "node" {
		fmt.Fprintln(os.Stderr, "Error: type must be 'st-node' or 'node'")
		os.Exit(1)
	}

	// Read and parse config
	configData, err := os.ReadFile(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		os.Exit(1)
	}

	var nodes []NodeDef
	if err := json.Unmarshal(configData, &nodes); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config JSON: %v\n", err)
		os.Exit(1)
	}

	config := Config{Nodes: nodes}

	// Prepare template data
	data := map[string]interface{}{
		"Nodes":    config.Nodes,
		"NodeType": *nodeType,
	}

	// Generate main file
	if err := generateFile(*templatePath, *outputPath, data); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating main file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully generated: %s\n", *outputPath)

	// Generate util file if specified
	if *utilTemplatePath != "" && *utilOutputPath != "" {
		if err := generateFile(*utilTemplatePath, *utilOutputPath, data); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating util file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully generated: %s\n", *utilOutputPath)
	}
}

func generateFile(templatePath, outputPath string, data map[string]interface{}) error {
	// Load template
	tmplFuncs := template.FuncMap{
		"title": func(s string) string {
			if s == "" {
				return s
			}
			runes := []rune(s)
			return string(append([]rune{unicode.ToUpper(runes[0])}, runes[1:]...))
		},
	}
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(tmplFuncs).ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("loading template: %w", err)
	}

	// Generate code
	var output strings.Builder
	if err := tmpl.Execute(&output, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	// Create output directory if needed
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Write output file
	if err := os.WriteFile(outputPath, []byte(output.String()), 0644); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	return nil
}
