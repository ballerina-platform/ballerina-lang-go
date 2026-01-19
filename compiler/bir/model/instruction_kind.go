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

package model

type InstructionKind uint8

const (
	INSTRUCTION_KIND_GOTO                InstructionKind = 1
	INSTRUCTION_KIND_CALL                InstructionKind = 2
	INSTRUCTION_KIND_BRANCH              InstructionKind = 3
	INSTRUCTION_KIND_RETURN              InstructionKind = 4
	INSTRUCTION_KIND_ASYNC_CALL          InstructionKind = 5
	INSTRUCTION_KIND_WAIT                InstructionKind = 6
	INSTRUCTION_KIND_FP_CALL             InstructionKind = 7
	INSTRUCTION_KIND_WK_RECEIVE          InstructionKind = 8
	INSTRUCTION_KIND_WK_SEND             InstructionKind = 9
	INSTRUCTION_KIND_FLUSH               InstructionKind = 10
	INSTRUCTION_KIND_LOCK                InstructionKind = 11
	INSTRUCTION_KIND_FIELD_LOCK          InstructionKind = 12
	INSTRUCTION_KIND_UNLOCK              InstructionKind = 13
	INSTRUCTION_KIND_WAIT_ALL            InstructionKind = 14
	INSTRUCTION_KIND_WK_ALT_RECEIVE      InstructionKind = 15
	INSTRUCTION_KIND_WK_MULTIPLE_RECEIVE InstructionKind = 16

	INSTRUCTION_KIND_MOVE                 InstructionKind = 20
	INSTRUCTION_KIND_CONST_LOAD           InstructionKind = 21
	INSTRUCTION_KIND_NEW_STRUCTURE        InstructionKind = 22
	INSTRUCTION_KIND_MAP_STORE            InstructionKind = 23
	INSTRUCTION_KIND_MAP_LOAD             InstructionKind = 24
	INSTRUCTION_KIND_NEW_ARRAY            InstructionKind = 25
	INSTRUCTION_KIND_ARRAY_STORE          InstructionKind = 26
	INSTRUCTION_KIND_ARRAY_LOAD           InstructionKind = 27
	INSTRUCTION_KIND_NEW_ERROR            InstructionKind = 28
	INSTRUCTION_KIND_TYPE_CAST            InstructionKind = 29
	INSTRUCTION_KIND_IS_LIKE              InstructionKind = 30
	INSTRUCTION_KIND_TYPE_TEST            InstructionKind = 31
	INSTRUCTION_KIND_NEW_INSTANCE         InstructionKind = 32
	INSTRUCTION_KIND_OBJECT_STORE         InstructionKind = 33
	INSTRUCTION_KIND_OBJECT_LOAD          InstructionKind = 34
	INSTRUCTION_KIND_PANIC                InstructionKind = 35
	INSTRUCTION_KIND_FP_LOAD              InstructionKind = 36
	INSTRUCTION_KIND_STRING_LOAD          InstructionKind = 37
	INSTRUCTION_KIND_NEW_XML_ELEMENT      InstructionKind = 38
	INSTRUCTION_KIND_NEW_XML_TEXT         InstructionKind = 39
	INSTRUCTION_KIND_NEW_XML_COMMENT      InstructionKind = 40
	INSTRUCTION_KIND_NEW_XML_PI           InstructionKind = 41
	INSTRUCTION_KIND_NEW_XML_SEQUENCE     InstructionKind = 42
	INSTRUCTION_KIND_NEW_XML_QNAME        InstructionKind = 43
	INSTRUCTION_KIND_NEW_STRING_XML_QNAME InstructionKind = 44
	INSTRUCTION_KIND_XML_SEQ_STORE        InstructionKind = 45
	INSTRUCTION_KIND_XML_SEQ_LOAD         InstructionKind = 46
	INSTRUCTION_KIND_XML_LOAD             InstructionKind = 47
	INSTRUCTION_KIND_XML_LOAD_ALL         InstructionKind = 48
	INSTRUCTION_KIND_XML_ATTRIBUTE_LOAD   InstructionKind = 49
	INSTRUCTION_KIND_XML_ATTRIBUTE_STORE  InstructionKind = 50
	INSTRUCTION_KIND_NEW_TABLE            InstructionKind = 51
	INSTRUCTION_KIND_NEW_TYPEDESC         InstructionKind = 52
	INSTRUCTION_KIND_NEW_STREAM           InstructionKind = 53
	INSTRUCTION_KIND_TABLE_STORE          InstructionKind = 54
	INSTRUCTION_KIND_TABLE_LOAD           InstructionKind = 55

	INSTRUCTION_KIND_ADD             InstructionKind = 61
	INSTRUCTION_KIND_SUB             InstructionKind = 62
	INSTRUCTION_KIND_MUL             InstructionKind = 63
	INSTRUCTION_KIND_DIV             InstructionKind = 64
	INSTRUCTION_KIND_MOD             InstructionKind = 65
	INSTRUCTION_KIND_EQUAL           InstructionKind = 66
	INSTRUCTION_KIND_NOT_EQUAL       InstructionKind = 67
	INSTRUCTION_KIND_GREATER_THAN    InstructionKind = 68
	INSTRUCTION_KIND_GREATER_EQUAL   InstructionKind = 69
	INSTRUCTION_KIND_LESS_THAN       InstructionKind = 70
	INSTRUCTION_KIND_LESS_EQUAL      InstructionKind = 71
	INSTRUCTION_KIND_AND             InstructionKind = 72
	INSTRUCTION_KIND_OR              InstructionKind = 73
	INSTRUCTION_KIND_REF_EQUAL       InstructionKind = 74
	INSTRUCTION_KIND_REF_NOT_EQUAL   InstructionKind = 75
	INSTRUCTION_KIND_CLOSED_RANGE    InstructionKind = 76
	INSTRUCTION_KIND_HALF_OPEN_RANGE InstructionKind = 77
	INSTRUCTION_KIND_ANNOT_ACCESS    InstructionKind = 78

	INSTRUCTION_KIND_TYPEOF                       InstructionKind = 80
	INSTRUCTION_KIND_NOT                          InstructionKind = 81
	INSTRUCTION_KIND_NEGATE                       InstructionKind = 82
	INSTRUCTION_KIND_BITWISE_AND                  InstructionKind = 83
	INSTRUCTION_KIND_BITWISE_OR                   InstructionKind = 84
	INSTRUCTION_KIND_BITWISE_XOR                  InstructionKind = 85
	INSTRUCTION_KIND_BITWISE_LEFT_SHIFT           InstructionKind = 86
	INSTRUCTION_KIND_BITWISE_RIGHT_SHIFT          InstructionKind = 87
	INSTRUCTION_KIND_BITWISE_UNSIGNED_RIGHT_SHIFT InstructionKind = 88

	INSTRUCTION_KIND_NEW_REG_EXP                InstructionKind = 89
	INSTRUCTION_KIND_NEW_RE_DISJUNCTION         InstructionKind = 90
	INSTRUCTION_KIND_NEW_RE_SEQUENCE            InstructionKind = 91
	INSTRUCTION_KIND_NEW_RE_ASSERTION           InstructionKind = 92
	INSTRUCTION_KIND_NEW_RE_ATOM_QUANTIFIER     InstructionKind = 93
	INSTRUCTION_KIND_NEW_RE_LITERAL_CHAR_ESCAPE InstructionKind = 94
	INSTRUCTION_KIND_NEW_RE_CHAR_CLASS          InstructionKind = 95
	INSTRUCTION_KIND_NEW_RE_CHAR_SET            InstructionKind = 96
	INSTRUCTION_KIND_NEW_RE_CHAR_SET_RANGE      InstructionKind = 97
	INSTRUCTION_KIND_NEW_RE_CAPTURING_GROUP     InstructionKind = 98
	INSTRUCTION_KIND_NEW_RE_FLAG_EXPR           InstructionKind = 99
	INSTRUCTION_KIND_NEW_RE_FLAG_ON_OFF         InstructionKind = 100
	INSTRUCTION_KIND_NEW_RE_QUANTIFIER          InstructionKind = 101
	INSTRUCTION_KIND_RECORD_DEFAULT_FP_LOAD     InstructionKind = 102
	INSTRUCTION_KIND_PLATFORM                   InstructionKind = 128
)

func (i InstructionKind) GetValue() uint8 {
	return uint8(i)
}
