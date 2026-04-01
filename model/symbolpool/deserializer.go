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

package symbolpool

import (
	"bytes"
	"fmt"

	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/semtypes/typepool"
)

type symbolReader struct {
	r   *bytes.Reader
	cp  []string
	tp  *typepool.TypePool
	env *context.CompilerEnvironment
}

func Unmarshal(env *context.CompilerEnvironment, data []byte) (model.ExportedSymbolSpace, error) {
	sr := &symbolReader{
		r:   bytes.NewReader(data),
		env: env,
	}
	return sr.deserialize()
}

func (sr *symbolReader) deserialize() (result model.ExportedSymbolSpace, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = model.ExportedSymbolSpace{}
			err = fmt.Errorf("symbol deserializer failed: %v", r)
		}
	}()

	magic := make([]byte, 4)
	_, err = sr.r.Read(magic)
	if err != nil {
		panic(fmt.Sprintf("reading magic: %v", err))
	}
	if string(magic) != symMagic {
		panic(fmt.Sprintf("invalid symbol magic: %x", magic))
	}

	var version int32
	read(sr.r, &version)
	if version != symVersion {
		panic(fmt.Sprintf("unsupported symbol version: %d", version))
	}

	var tpSize int64
	read(sr.r, &tpSize)
	tpBytes := make([]byte, tpSize)
	_, err = sr.r.Read(tpBytes)
	if err != nil {
		panic(fmt.Sprintf("reading type pool: %v", err))
	}
	sr.tp = typepool.UnmarshalTypePool(tpBytes, sr.env.GetTypeEnv())

	sr.cp = deserializeConstantPool(sr.r)

	mainSpace := sr.readSymbolSpace()
	annotationSpace := sr.readSymbolSpace()

	return model.NewExportedSymbolSpace(mainSpace, annotationSpace), nil
}

func (sr *symbolReader) readPackageIdentifier() *model.PackageID {
	org := sr.readStringCP()
	pkg := sr.readStringCP()
	version := sr.readStringCP()
	nameComps := model.CreateNameComps(model.Name(pkg))
	versionName := model.Name(version)
	if versionName == "" {
		versionName = model.DEFAULT_VERSION
	}
	return sr.env.NewPackageID(model.Name(org), nameComps, versionName)
}

func (sr *symbolReader) readSymbolSpace() *model.SymbolSpace {
	var count int64
	read(sr.r, &count)
	if count == 0 {
		return nil
	}

	pkgID := sr.readPackageIdentifier()
	space := sr.env.NewSymbolSpace(*pkgID)
	for i := int64(0); i < count; i++ {
		sr.readSymbol(space)
	}

	return space
}

func (sr *symbolReader) readSymbol(space *model.SymbolSpace) {
	var tag uint8
	read(sr.r, &tag)

	switch tag {
	case symTagType:
		sr.readTypeSymbol(space)
	case symTagClass:
		sr.readClassSymbol(space)
	case symTagValue:
		sr.readValueSymbol(space)
	case symTagFunction:
		sr.readFunctionSymbol(space)
	default:
		panic(fmt.Sprintf("unknown symbol tag: %d", tag))
	}
}

func (sr *symbolReader) readSymbolBase() (name string, isPublic bool, ty semtypes.SemType) {
	name = sr.readStringCP()
	read(sr.r, &isPublic)
	ty = sr.readType()
	return
}

func (sr *symbolReader) readTypeSymbol(space *model.SymbolSpace) {
	name, isPublic, ty := sr.readSymbolBase()
	sym := model.NewTypeSymbol(name, isPublic)
	sym.SetType(ty)
	space.AddSymbol(name, &sym)
}

func (sr *symbolReader) readClassSymbol(space *model.SymbolSpace) {
	name, isPublic, ty := sr.readSymbolBase()
	sym := model.NewClassSymbol(name, isPublic)
	sym.SetType(ty)
	space.AddSymbol(name, &sym)
}

func (sr *symbolReader) readValueSymbol(space *model.SymbolSpace) {
	name, isPublic, ty := sr.readSymbolBase()
	var isConst, isParameter bool
	read(sr.r, &isConst)
	read(sr.r, &isParameter)
	sym := model.NewValueSymbol(name, isPublic, isConst, isParameter)
	sym.SetType(ty)
	space.AddSymbol(name, &sym)
}

func (sr *symbolReader) readFunctionSymbol(space *model.SymbolSpace) {
	name, isPublic, ty := sr.readSymbolBase()

	var paramCount int64
	read(sr.r, &paramCount)
	paramTypes := make([]semtypes.SemType, paramCount)
	for i := int64(0); i < paramCount; i++ {
		paramTypes[i] = sr.readType()
	}
	returnType := sr.readType()
	var hasRestParam bool
	read(sr.r, &hasRestParam)
	var restParamType semtypes.SemType
	if hasRestParam {
		restParamType = sr.readType()
	}

	sig := model.FunctionSignature{
		ParamTypes:    paramTypes,
		ReturnType:    returnType,
		RestParamType: restParamType,
	}
	sym := model.NewFunctionSymbol(name, sig, isPublic)
	sym.SetType(ty)
	space.AddSymbol(name, sym)
}

func (sr *symbolReader) readStringCP() string {
	var idx int32
	read(sr.r, &idx)
	return sr.cp[idx]
}

func (sr *symbolReader) readType() semtypes.SemType {
	var idx int32
	read(sr.r, &idx)
	if idx == -1 {
		return nil
	}
	return sr.tp.Get(typepool.Index(idx))
}
