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

package birserializer

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/model"
)

type CPEntryType uint8

const (
	CP_ENTRY_INTEGER CPEntryType = iota + 1
	CP_ENTRY_FLOAT
	CP_ENTRY_BOOLEAN
	CP_ENTRY_STRING
	CP_ENTRY_PACKAGE
	CP_ENTRY_BYTE
	CP_ENTRY_SHAPE
)

type CPEntry interface {
	EntryType() CPEntryType
}

type (
	IntegerCPEntry struct {
		Value int64
	}
	FloatCPEntry struct {
		Value float64
	}
	BooleanCPEntry struct {
		Value bool
	}
	StringCPEntry struct {
		Value string
	}
	ByteCPEntry struct {
		Value byte
	}
	PackageCPEntry struct {
		OrgNameCPIndex    int32
		PkgNameCPIndex    int32
		ModuleNameCPIndex int32
		VersionCPIndex    int32
	}
	ShapeCPEntry struct {
		Shape ast.BType
	}
)

func (e *IntegerCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_INTEGER
}

func (e *FloatCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_FLOAT
}

func (e *BooleanCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_BOOLEAN
}

func (e *StringCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_STRING
}

func (e *PackageCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_PACKAGE
}

func (e *ByteCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_BYTE
}

func (e *ShapeCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_SHAPE
}

type ConstantPool struct {
	entries  []CPEntry
	entryMap map[string]int
}

func (cp *ConstantPool) EntryKey(entry CPEntry) string {
	switch e := entry.(type) {
	case *IntegerCPEntry:
		return fmt.Sprintf("int:%d", e.Value)
	case *FloatCPEntry:
		return fmt.Sprintf("float:%g", e.Value)
	case *BooleanCPEntry:
		return fmt.Sprintf("bool:%v", e.Value)
	case *StringCPEntry:
		return fmt.Sprintf("str:%s", e.Value)
	case *PackageCPEntry:
		return fmt.Sprintf("pkg:%d:%d:%d:%d", e.OrgNameCPIndex, e.PkgNameCPIndex, e.ModuleNameCPIndex, e.VersionCPIndex)
	case *ByteCPEntry:
		return fmt.Sprintf("byte:%d", e.Value)
	case *ShapeCPEntry:
		panic("shape key generation not implemented")
	default:
		panic("unknown CPEntry type")
	}
}

func NewConstantPool() *ConstantPool {
	return &ConstantPool{
		entries:  make([]CPEntry, 0),
		entryMap: make(map[string]int),
	}
}

func (cp *ConstantPool) AddEntry(entry CPEntry) int32 {
	key := cp.EntryKey(entry)
	if index, exists := cp.entryMap[key]; exists {
		return int32(index)
	}

	index := len(cp.entries)
	cp.entries = append(cp.entries, entry)
	cp.entryMap[key] = index
	return int32(index)
}

func (cp *ConstantPool) AddIntegerCPEntry(value int64) int32 {
	return cp.AddEntry(&IntegerCPEntry{Value: value})
}

func (cp *ConstantPool) AddFloatCPEntry(value float64) int32 {
	return cp.AddEntry(&FloatCPEntry{Value: value})
}

func (cp *ConstantPool) AddBooleanCPEntry(value bool) int32 {
	return cp.AddEntry(&BooleanCPEntry{Value: value})
}

func (cp *ConstantPool) AddStringCPEntry(value string) int32 {
	return cp.AddEntry(&StringCPEntry{Value: value})
}

func (cp *ConstantPool) AddPackageCPEntry(pkg *model.PackageID) int32 {
	return cp.AddEntry(&PackageCPEntry{
		OrgNameCPIndex:    cp.AddStringCPEntry(pkg.OrgName.Value()),
		PkgNameCPIndex:    cp.AddStringCPEntry(pkg.PkgName.Value()),
		ModuleNameCPIndex: cp.AddStringCPEntry(pkg.Name.Value()),
		VersionCPIndex:    cp.AddStringCPEntry(pkg.Version.Value()),
	})
}

func (cp *ConstantPool) AddByteCPEntry(value byte) int32 {
	return cp.AddEntry(&ByteCPEntry{Value: value})
}

func (cp *ConstantPool) AddShapeCPEntry(shape ast.BType) int32 {
	panic("shape entry addition not implemented")
}

func (cp *ConstantPool) WriteCPEntry(buf *bytes.Buffer, entry CPEntry) error {
	entryType := entry.EntryType()
	if err := cp.writeInt8(buf, int8(entryType)); err != nil {
		return err
	}

	switch e := entry.(type) {
	case *IntegerCPEntry:
		return cp.writeInt64(buf, e.Value)
	case *FloatCPEntry:
		return cp.writeFloat64(buf, e.Value)
	case *BooleanCPEntry:
		var b byte
		if e.Value {
			b = 1
		}
		return cp.writeUInt8(buf, uint8(b))
	case *StringCPEntry:
		strBytes := []byte(e.Value)
		if err := cp.writeInt32(buf, int32(len(strBytes))); err != nil {
			return err
		}
		_, err := buf.Write(strBytes)
		return err
	case *ByteCPEntry:
		return cp.writeInt32(buf, int32(e.Value))
	case *PackageCPEntry:
		if err := cp.writeInt32(buf, int32(e.OrgNameCPIndex)); err != nil {
			return err
		}
		if err := cp.writeInt32(buf, int32(e.PkgNameCPIndex)); err != nil {
			return err
		}
		if err := cp.writeInt32(buf, int32(e.ModuleNameCPIndex)); err != nil {
			return err
		}
		return cp.writeInt32(buf, int32(e.VersionCPIndex))
	case *ShapeCPEntry:
		panic("shape serialization not implemented")
	default:
		return fmt.Errorf("unsupported constant pool entry type: %T", entry)
	}
}

func (cp *ConstantPool) Serialize() ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := cp.writeInt32(buf, int32(-1)); err != nil {
		return nil, err
	}

	for _, entry := range cp.entries {
		if err := cp.WriteCPEntry(buf, entry); err != nil {
			return nil, err
		}
	}

	bytes := buf.Bytes()
	entryCount := int32(len(cp.entries))
	binary.BigEndian.PutUint32(bytes[0:4], uint32(entryCount))

	return bytes, nil
}

func (cp *ConstantPool) writeInt8(buf *bytes.Buffer, val int8) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (cp *ConstantPool) writeUInt8(buf *bytes.Buffer, val uint8) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (cp *ConstantPool) writeInt32(buf *bytes.Buffer, val int32) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (cp *ConstantPool) writeInt64(buf *bytes.Buffer, val int64) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (cp *ConstantPool) writeUInt64(buf *bytes.Buffer, val uint64) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (cp *ConstantPool) writeFloat64(buf *bytes.Buffer, val float64) error {
	return binary.Write(buf, binary.BigEndian, val)
}
