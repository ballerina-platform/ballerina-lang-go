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

package exec

import (
	"ballerina-lang-go/bir"
	"ballerina-lang-go/runtime/internal/modules"
	"ballerina-lang-go/tools/diagnostics"
	"ballerina-lang-go/values"
)

type Frame struct {
	locals      []values.BalValue    // variable index → value (indexed by BIROperand.Index)
	functionKey string               // function key (package name + function name)
	location    diagnostics.Location // source location of the currently executing instruction/terminator
}

func getOperandValue(op *bir.BIROperand, currentFrame *Frame, reg *modules.Registry) values.BalValue {
	if gv, ok := op.VariableDcl.(*bir.BIRGlobalVariableDcl); ok {
		module := reg.GetModule(gv.PkgId)
		return module.Globals[*op.SymRef]
	}
	return currentFrame.locals[op.Index]
}

func setOperandValue(op *bir.BIROperand, currentFrame *Frame, reg *modules.Registry, value values.BalValue) {
	if gv, ok := op.VariableDcl.(*bir.BIRGlobalVariableDcl); ok {
		module := reg.GetModule(gv.PkgId)
		module.Globals[*op.SymRef] = value
	} else {
		currentFrame.locals[op.Index] = value
	}
}

// SetLocation updates the current source location associated with this frame.
func (f *Frame) SetLocation(loc diagnostics.Location) {
	f.location = loc
}

// Location returns the current source location associated with this frame.
func (f *Frame) Location() diagnostics.Location {
	return f.location
}
