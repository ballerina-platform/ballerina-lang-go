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

package tomlparser

import (
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type Toml struct {
	rootNode    map[string]any
	metadata    toml.MetaData
	diagnostics []Diagnostic
	content     string
}

type Diagnostic struct {
	Message  string
	Severity string
	Location *Location
}

type Location struct {
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

func Read(path string) (*Toml, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ReadString(string(content))
}

func ReadWithSchema(path string, schema Schema) (*Toml, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ReadStringWithSchema(string(content), schema)
}

func ReadStream(reader io.Reader) (*Toml, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return ReadString(string(content))
}

func ReadStreamWithSchema(reader io.Reader, schema Schema) (*Toml, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return ReadStringWithSchema(string(content), schema)
}

func ReadString(content string) (*Toml, error) {
	var data map[string]any
	metadata, err := toml.Decode(content, &data)

	t := &Toml{
		rootNode:    data,
		metadata:    metadata,
		diagnostics: make([]Diagnostic, 0),
		content:     content,
	}

	if err != nil {
		t.diagnostics = append(t.diagnostics, Diagnostic{
			Message:  err.Error(),
			Severity: "ERROR",
		})
	}

	return t, err
}

func ReadStringWithSchema(content string, schema Schema) (*Toml, error) {
	t, err := ReadString(content)
	if err != nil {
		return t, err
	}

	validator := NewValidator(schema)
	validationErr := validator.Validate(t)
	if validationErr != nil {
		return t, validationErr
	}

	return t, nil
}

func (t *Toml) Validate(schema Schema) error {
	validator := NewValidator(schema)
	return validator.Validate(t)
}

func (t *Toml) Get(dottedKey string) (any, bool) {
	keys := splitDottedKey(dottedKey)
	return t.getValueByPath(keys)
}

func (t *Toml) GetString(dottedKey string) (string, bool) {
	val, ok := t.Get(dottedKey)
	if !ok {
		return "", false
	}
	str, ok := val.(string)
	return str, ok
}

func (t *Toml) GetInt(dottedKey string) (int64, bool) {
	val, ok := t.Get(dottedKey)
	if !ok {
		return 0, false
	}

	switch v := val.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	default:
		return 0, false
	}
}

func (t *Toml) GetFloat(dottedKey string) (float64, bool) {
	val, ok := t.Get(dottedKey)
	if !ok {
		return 0, false
	}
	f, ok := val.(float64)
	return f, ok
}

func (t *Toml) GetBool(dottedKey string) (bool, bool) {
	val, ok := t.Get(dottedKey)
	if !ok {
		return false, false
	}
	b, ok := val.(bool)
	return b, ok
}

func (t *Toml) GetArray(dottedKey string) ([]any, bool) {
	val, ok := t.Get(dottedKey)
	if !ok {
		return nil, false
	}
	arr, ok := val.([]any)
	return arr, ok
}

func (t *Toml) GetTable(dottedKey string) (*Toml, bool) {
	val, ok := t.Get(dottedKey)
	if !ok {
		return nil, false
	}

	table, ok := val.(map[string]any)
	if !ok {
		return nil, false
	}

	return &Toml{
		rootNode:    table,
		metadata:    t.metadata,
		diagnostics: make([]Diagnostic, 0),
		content:     "",
	}, true
}

func (t *Toml) GetTables(dottedKey string) ([]*Toml, bool) {
	val, ok := t.Get(dottedKey)
	if !ok {
		return nil, false
	}

	arr, ok := val.([]any)
	if !ok {
		if tableArr, ok := val.([]map[string]any); ok {
			result := make([]*Toml, len(tableArr))
			for i, table := range tableArr {
				result[i] = &Toml{
					rootNode:    table,
					metadata:    t.metadata,
					diagnostics: make([]Diagnostic, 0),
					content:     "",
				}
			}
			return result, true
		}
		return nil, false
	}

	result := make([]*Toml, 0)
	for _, item := range arr {
		if table, ok := item.(map[string]any); ok {
			result = append(result, &Toml{
				rootNode:    table,
				metadata:    t.metadata,
				diagnostics: make([]Diagnostic, 0),
				content:     "",
			})
		}
	}

	if len(result) == 0 {
		return nil, false
	}

	return result, true
}

func (t *Toml) Diagnostics() []Diagnostic {
	return t.diagnostics
}

func (t *Toml) RootNode() map[string]any {
	return t.rootNode
}

func (t *Toml) ToMap() map[string]any {
	return t.rootNode
}

func (t *Toml) To(target any) {
	_, err := toml.Decode(t.content, target)
	if err != nil {
		t.diagnostics = append(t.diagnostics, Diagnostic{
			Message:  err.Error(),
			Severity: "ERROR",
		})
	}
}

func splitDottedKey(dottedKey string) []string {
	return strings.Split(dottedKey, ".")
}

func (t *Toml) getValueByPath(keys []string) (any, bool) {
	current := any(t.rootNode)

	for _, key := range keys {
		key = strings.Trim(key, "\"")

		currentMap, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}

		val, exists := currentMap[key]
		if !exists {
			return nil, false
		}

		current = val
	}

	return current, true
}
