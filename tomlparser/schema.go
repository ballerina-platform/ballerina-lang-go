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
	"encoding/json"
	"fmt"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Schema interface {
	Validate(data any) error
	FromPath(path string) (Schema, error)
	FromString(content string) (Schema, error)
}

type schemaImpl struct {
	compiled *jsonschema.Schema
}

func NewSchemaFromPath(path string) (Schema, error) {
	compiler := jsonschema.NewCompiler()

	schema, err := compiler.Compile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema from path %s: %w", path, err)
	}

	return &schemaImpl{
		compiled: schema,
	}, nil
}

func NewSchemaFromString(content string) (Schema, error) {
	compiler := jsonschema.NewCompiler()

	var schemaDoc any
	if err := json.Unmarshal([]byte(content), &schemaDoc); err != nil {
		return nil, fmt.Errorf("failed to parse schema JSON: %w", err)
	}

	if err := compiler.AddResource("schema.json", schemaDoc); err != nil {
		return nil, fmt.Errorf("failed to add schema resource: %w", err)
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}

	return &schemaImpl{
		compiled: schema,
	}, nil
}

func NewSchemaFromFile(file *os.File) (Schema, error) {
	compiler := jsonschema.NewCompiler()

	var schemaDoc any
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&schemaDoc); err != nil {
		return nil, fmt.Errorf("failed to decode schema JSON: %w", err)
	}

	if err := compiler.AddResource("schema.json", schemaDoc); err != nil {
		return nil, fmt.Errorf("failed to add schema resource: %w", err)
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}

	return &schemaImpl{
		compiled: schema,
	}, nil
}

func (s *schemaImpl) Validate(data any) error {
	if err := s.compiled.Validate(data); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}
	return nil
}

func (s *schemaImpl) FromPath(path string) (Schema, error) {
	return NewSchemaFromPath(path)
}

func (s *schemaImpl) FromString(content string) (Schema, error) {
	return NewSchemaFromString(content)
}
