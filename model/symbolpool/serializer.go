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

	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

const (
	symMagic = "\x53\x59\x4d\x42"
	// This will perpetually remain 1 unless we create a spec for this
	symVersion = 1
)

const (
	symTagType     uint8 = 0
	symTagClass    uint8 = 1
	symTagValue    uint8 = 2
	symTagFunction uint8 = 3
)

type symbolWriter struct {
	cp  *constantPool
	tp  *semtypes.TypePool
	env semtypes.Env
}

func Marshal(exported model.ExportedSymbolSpace, env semtypes.Env) ([]byte, error) {
	sw := &symbolWriter{
		cp:  newConstantPool(),
		tp:  semtypes.NewTypePool(),
		env: env,
	}
	return sw.serialize(exported)
}

func (sw *symbolWriter) serialize(exported model.ExportedSymbolSpace) ([]byte, error) {
	body := &bytes.Buffer{}
	if err := sw.writeSymbolSpace(body, exported.Main); err != nil {
		return nil, err
	}
	if err := sw.writeSymbolSpace(body, exported.Annotation); err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if _, err := buf.Write([]byte(symMagic)); err != nil {
		return nil, fmt.Errorf("writing magic: %v", err)
	}
	if err := write(buf, int32(symVersion)); err != nil {
		return nil, err
	}

	tpBytes := semtypes.MarshalTypePool(sw.tp, sw.env)
	if err := write(buf, int64(len(tpBytes))); err != nil {
		return nil, err
	}
	if _, err := buf.Write(tpBytes); err != nil {
		return nil, fmt.Errorf("writing type pool: %v", err)
	}

	cpBytes, err := sw.cp.serialize()
	if err != nil {
		return nil, fmt.Errorf("writing constant pool: %v", err)
	}
	if _, err := buf.Write(cpBytes); err != nil {
		return nil, fmt.Errorf("writing constant pool: %v", err)
	}

	if _, err := buf.Write(body.Bytes()); err != nil {
		return nil, fmt.Errorf("writing body: %v", err)
	}

	return buf.Bytes(), nil
}

func (sw *symbolWriter) writePackageIdentifier(buf *bytes.Buffer, pkg model.PackageIdentifier) error {
	if err := sw.writeStringCP(buf, pkg.Organization); err != nil {
		return err
	}
	if err := sw.writeStringCP(buf, pkg.Package); err != nil {
		return err
	}
	return sw.writeStringCP(buf, pkg.Version)
}

func (sw *symbolWriter) writeSymbolSpace(buf *bytes.Buffer, space *model.SymbolSpace) error {
	if space == nil {
		return write(buf, int64(0))
	}

	if err := write(buf, int64(space.Len())); err != nil {
		return err
	}
	if err := sw.writePackageIdentifier(buf, space.Pkg); err != nil {
		return err
	}

	for _, sym := range space.Symbols() {
		if err := sw.writeSymbol(buf, sym); err != nil {
			return err
		}
	}
	return nil
}

func (sw *symbolWriter) writeSymbol(buf *bytes.Buffer, sym model.Symbol) error {
	switch s := sym.(type) {
	case *model.ClassSymbol:
		return sw.writeClassSymbol(buf, s)
	case *model.TypeSymbol:
		return sw.writeTypeSymbol(buf, s)
	case *model.ValueSymbol:
		return sw.writeValueSymbol(buf, s)
	case model.FunctionSymbol:
		return sw.writeFunctionSymbol(buf, s)
	default:
		return fmt.Errorf("unsupported symbol type: %T", sym)
	}
}

func (sw *symbolWriter) writeSymbolBase(buf *bytes.Buffer, sym model.Symbol) error {
	if err := sw.writeStringCP(buf, sym.Name()); err != nil {
		return err
	}
	if err := write(buf, sym.IsPublic()); err != nil {
		return err
	}
	return sw.writeType(buf, sym.Type())
}

func (sw *symbolWriter) writeTypeSymbol(buf *bytes.Buffer, sym *model.TypeSymbol) error {
	if err := write(buf, symTagType); err != nil {
		return err
	}
	return sw.writeSymbolBase(buf, sym)
}

func (sw *symbolWriter) writeClassSymbol(buf *bytes.Buffer, sym *model.ClassSymbol) error {
	if err := write(buf, symTagClass); err != nil {
		return err
	}
	return sw.writeSymbolBase(buf, sym)
}

func (sw *symbolWriter) writeValueSymbol(buf *bytes.Buffer, sym *model.ValueSymbol) error {
	if err := write(buf, symTagValue); err != nil {
		return err
	}
	if err := sw.writeSymbolBase(buf, sym); err != nil {
		return err
	}
	if err := write(buf, sym.Kind() == model.SymbolKindConstant); err != nil {
		return err
	}
	return write(buf, sym.Kind() == model.SymbolKindParemeter)
}

func (sw *symbolWriter) writeFunctionSymbol(buf *bytes.Buffer, sym model.FunctionSymbol) error {
	if err := write(buf, symTagFunction); err != nil {
		return err
	}
	if err := sw.writeSymbolBase(buf, sym); err != nil {
		return err
	}
	sig := sym.Signature()
	if err := write(buf, int64(len(sig.ParamTypes))); err != nil {
		return err
	}
	for _, pt := range sig.ParamTypes {
		if err := sw.writeType(buf, pt); err != nil {
			return err
		}
	}
	if err := sw.writeType(buf, sig.ReturnType); err != nil {
		return err
	}
	if err := write(buf, sig.RestParamType != nil); err != nil {
		return err
	}
	if sig.RestParamType != nil {
		if err := sw.writeType(buf, sig.RestParamType); err != nil {
			return err
		}
	}
	return nil
}

func (sw *symbolWriter) writeStringCP(buf *bytes.Buffer, s string) error {
	return write(buf, sw.cp.addString(s))
}

func (sw *symbolWriter) writeType(buf *bytes.Buffer, ty semtypes.SemType) error {
	if ty == nil {
		return write(buf, int32(-1))
	}
	return write(buf, int32(sw.tp.Put(ty)))
}
