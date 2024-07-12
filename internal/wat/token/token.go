// 版权 @2024 凹语言 作者。保留所有权利。

// wat记号类型
package token

import (
	"strconv"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm"
)

// 记号类型
type Token int

// 指令码
type Opcode = wasm.Opcode

const (
	// 非法/结尾/注释
	ILLEGAL Token = iota
	EOF
	COMMENT

	// 特殊符号
	operator_beg
	LPAREN // (
	RPAREN // )
	ASSIGN // =, i32.load offset=8 align=4
	operator_end

	// 面值类型
	literal_beg
	IDENT  // 标识符, 比如 $name
	INT    // 12345
	FLOAT  // 123.45
	CHAR   // 'a'
	STRING // "abc"
	literal_end

	// 关键字
	keyword_beg

	I32 // i32
	I64 // i64
	F32 // f32
	F64 // f64

	MUT     // mut
	ANYFUNC // anyfunc
	OFFSET  // offset
	ALIGN   // align

	MODULE // module
	IMPORT // import
	EXPORT // export

	MEMORY // memory
	TABLE  // table
	GLOBAL // global
	LOCAL  // local
	DATA   // data
	ELEM   // elem
	TYPE   // type

	FUNC   // func
	PARAM  // param
	RESULT // result

	START // start

	// 指令是关键字, 同时有指令码
	// 指令码需要单独查表计算, 不能通过 instruction_beg 计算
	ins_beg

	// https://www.w3.org/TR/2019/REC-wasm-core-1-20191205/

	INS_UNREACHABLE   // 0x00, unreachable
	INS_NOP           // 0x01, nop
	INS_BLOCK         // 0x02, block
	INS_LOOP          // 0x03, loop
	INS_IF            // 0x04, if
	INS_ELSE          // 0x05, else
	INS_END           // 0x0b, end
	INS_BR            // 0x0c, br
	INS_BR_IF         // 0x0d, br_if
	INS_BR_TABLE      // 0x0e, br_table
	INS_RETURN        // 0x0f, return
	INS_CALL          // 0x10, call
	INS_CALL_INDIRECT // 0x11, call_indirect
	INS_DROP          // 0x1a, drop
	INS_SELECT        // 0x1b, select
	INS_TYPED_SELECT  // 0x1c, typed_select
	INS_LOCAL_GET     // 0x20, local.get
	INS_LOCAL_SET     // 0x21, local.set
	INS_LOCAL_TEE     // 0x22, local.tee
	INS_GLOBAL_GET    // 0x23, global.get
	INS_GLOBAL_SET    // 0x24, global.set
	INS_TABLE_GET     // 0x25, table.get
	INS_TABLE_SET     // 0x26, table.set
	INS_I32_LOAD      // 0x28, i32.load
	INS_I64_LOAD      // 0x29, i32.load
	INS_F32_LOAD      // 0x2a, f32.load
	INS_F64_LOAD      // 0x2b, f64.load
	INS_I32_LOAD_8S   // 0x2c, i32.load_8s
	INS_I32_LOAD_8U   // 0x2d, i32.load_8u
	INS_I32_LOAD_16S  // 0x2e, i32.load_16s
	INS_I32_LOAD_16U  // 0x2f, i32.load_16u
	INS_I64_LOAD_8S   // 0x30, i64.load_8s
	INS_I64_LOAD_8U   // 0x31, i64.load_8u
	INS_I64_LOAD_16S  // 0x31, i64.load_16s
	INS_I64_LOAD_16U  // 0x33, i64.load_16u
	INS_I64_LOAD_32S  // 0x34, i64.load_32s
	INS_I64_LOAD_32U  // 0x35, i64.load_32u
	INS_I32_STORE     // 0x36, i32.store
	INS_I64_STORE     // 0x37, i64.store
	INS_F32_STORE     // 0x38, f32.store
	INS_F64_STORE     // 0x39, f64.store
	INS_I32_STORE8    // 0x3a, i32.store8
	INS_I32_STORE16   // 0x3b, i32.store16
	INS_I64_STORE8    // 0x3c, i64.store8
	INS_I64_STORE16   // 0x3d, i64.store16
	INS_I64_STORE32   // 0x3e, i64.store32
	INS_MEMORY_SIZE   // 0x3f, memory.size
	INS_MEMORY_GROW   // 0x40, memory.grow
	INS_I32_CONST     // 0x41, i32.const
	INS_I64_CONST     // 0x42, i64.const
	INS_F32_CONST     // 0x43, f32.const
	INS_F64_CONST     // 0x44, f64.const
	INS_I32_EQZ       // 0x45, i32.eqz
	INS_I32_EQ        // 0x46, i32.eq
	INS_I32_NE        // 0x47, i32.ne
	INS_I32_LT_S      // 0x48, i32.lt_s
	INS_I32_LT_U      // 0x49, i32.lt_u
	INS_I32_GT_S      // 0x4a, i32.gt_s
	INS_I32_GT_U      // 0x4b, i32.gt_u
	INS_I32_LE_S      // 0x4c, i32.le_s
	INS_I32_LE_U      // 0x4d, i32.le_u
	INS_I32_GE_S      // 0x4e, i32.ge_s
	INS_I32_GE_U      // 0x4f, i32.ge_u
	INS_I64_EQZ       // 0x50, i64.eqz
	INS_I64_EQ        // 0x51, i64.eq
	INS_I64_NE        // 0x52, i64.ne
	INS_I64_LT_S      // 0x53, i64.lt_s
	INS_I64_LT_U      // 0x54, i64.lt_u
	INS_I64_GT_S      // 0x55, i64.gt_s
	INS_I64_GT_U      // 0x56, i64.gt_u
	INS_I64_LE_S      // 0x57, i64.le_s
	INS_I64_LE_U      // 0x58, i64.le_u
	INS_I64_GE_S      // 0x59, i64.ge_s
	INS_I64_GE_U      // 0x5a, i64.ge_u
	INS_F32_EQ        // 0x5b, f32.eq
	INS_F32_NE        // 0x5c, f32.ne
	INS_F32_LT        // 0x5d, f32.lt
	INS_F32_GT        // 0x5e, f32.gt
	INS_F32_LE        // 0x5f, f32.le
	INS_F32_GE        // 0x60, f32.ge
	INS_F64_EQ        // 0x61, f64.eq
	INS_F64_NE        // 0x62, f64.ne
	INS_F64_LT        // 0x63, f64.lt
	INS_F64_GT        // 0x64, f64.gt
	INS_F64_LE        // 0x65, f64.le
	INS_F64_GE        // 0x66, f64.ge

	LABEL // label

	ins_end
	keyword_end

	// TODO: 指令转关键字
	// if/else/loop/end/...
	INSTRUCTION // 指令码
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	CHAR:   "CHAR",
	STRING: "STRING",

	LPAREN: "(",
	RPAREN: ")",
	ASSIGN: "=",

	I32: "i32",
	I64: "i64",
	F32: "f32",
	F64: "f64",

	MUT:     "mut",
	ANYFUNC: "anyfunc",
	OFFSET:  "offset",
	ALIGN:   "align",

	MODULE: "module",
	IMPORT: "import",
	EXPORT: "export",

	MEMORY: "memory",
	TABLE:  "table",
	GLOBAL: "global",
	LOCAL:  "local",
	DATA:   "data",
	ELEM:   "elem",
	TYPE:   "type",

	FUNC:   "func",
	PARAM:  "param",
	RESULT: "result",

	START: "start",

	LABEL: "label",

	INS_UNREACHABLE:   "unreachable",
	INS_NOP:           "nop",
	INS_BLOCK:         "block",
	INS_LOOP:          "loop",
	INS_IF:            "if",
	INS_ELSE:          "else",
	INS_END:           "end",
	INS_BR:            "br",
	INS_BR_IF:         "br_if",
	INS_BR_TABLE:      "br_table",
	INS_RETURN:        "return",
	INS_CALL:          "call",
	INS_CALL_INDIRECT: "call_indirect",
	INS_DROP:          "drop",
	INS_SELECT:        "select",
	INS_TYPED_SELECT:  "typed_select",
	INS_LOCAL_GET:     "local.get",
	INS_LOCAL_SET:     "local.set",
	INS_LOCAL_TEE:     "local.tee",
	INS_GLOBAL_GET:    "global.get",
	INS_GLOBAL_SET:    "global.set",
	INS_TABLE_GET:     "table.get",
	INS_TABLE_SET:     "table.set",
	INS_I32_LOAD:      "i32.load",
	INS_I64_LOAD:      "i32.load",
	INS_F32_LOAD:      "f32.load",
	INS_F64_LOAD:      "f64.load",
	INS_I32_LOAD_8S:   "i32.load_8s",
	INS_I32_LOAD_8U:   "i32.load_8u",
	INS_I32_LOAD_16S:  "i32.load_16s",
	INS_I32_LOAD_16U:  "i32.load_16u",
	INS_I64_LOAD_8S:   "i64.load_8s",
	INS_I64_LOAD_8U:   "i64.load_8u",
	INS_I64_LOAD_16S:  "i64.load_16s",
	INS_I64_LOAD_16U:  "i64.load_16u",
	INS_I64_LOAD_32S:  "i64.load_32s",
	INS_I64_LOAD_32U:  "i64.load_32u",
	INS_I32_STORE:     "i32.store",
	INS_I64_STORE:     "i64.store",
	INS_F32_STORE:     "f32.store",
	INS_F64_STORE:     "f64.store",
	INS_I32_STORE8:    "i32.store8",
	INS_I32_STORE16:   "i32.store16",
	INS_I64_STORE8:    "i64.store8",
	INS_I64_STORE16:   "i64.store16",
	INS_I64_STORE32:   "i64.store32",
	INS_MEMORY_SIZE:   "memory.size",
	INS_MEMORY_GROW:   "memory.grow",
	INS_I32_CONST:     "i32.const",
	INS_I64_CONST:     "i64.const",
	INS_F32_CONST:     "f32.const",
	INS_F64_CONST:     "f64.const",
	INS_I32_EQZ:       "i32.eqz",
	INS_I32_EQ:        "i32.eq",
	INS_I32_NE:        "i32.ne",
	INS_I32_LT_S:      "i32.lt_s",
	INS_I32_LT_U:      "i32.lt_u",
	INS_I32_GT_S:      "i32.gt_s",
	INS_I32_GT_U:      "i32.gt_u",
	INS_I32_LE_S:      "i32.le_s",
	INS_I32_LE_U:      "i32.le_u",
	INS_I32_GE_S:      "i32.ge_s",
	INS_I32_GE_U:      "i32.ge_u",
	INS_I64_EQZ:       "i64.eqz",
	INS_I64_EQ:        "i64.eq",
	INS_I64_NE:        "i64.ne",
	INS_I64_LT_S:      "i64.lt_s",
	INS_I64_LT_U:      "i64.lt_u",
	INS_I64_GT_S:      "i64.gt_s",
	INS_I64_GT_U:      "i64.gt_u",
	INS_I64_LE_S:      "i64.le_s",
	INS_I64_LE_U:      "i64.le_u",
	INS_I64_GE_S:      "i64.ge_s",
	INS_I64_GE_U:      "i64.ge_u",
	INS_F32_EQ:        "f32.eq",
	INS_F32_NE:        "f32.ne",
	INS_F32_LT:        "f32.lt",
	INS_F32_GT:        "f32.gt",
	INS_F32_LE:        "f32.le",
	INS_F32_GE:        "f32.ge",
	INS_F64_EQ:        "f64.eq",
	INS_F64_NE:        "f64.ne",
	INS_F64_LT:        "f64.lt",
	INS_F64_GT:        "f64.gt",
	INS_F64_LE:        "f64.le",
	INS_F64_GE:        "f64.ge",

	INSTRUCTION: "INSTRUCTION",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)

	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}

	// $name
	if strings.HasPrefix(ident, "$") {
		return IDENT
	}

	// i32.load
	return INSTRUCTION
}

// 指令
func LookupOpcode(s string) (Opcode, bool) {
	return wasm.LookupOpcode(s)
}
func OpcodeName(op Opcode) string {
	return wasm.OpcodeName(op)
}

// 常量面值
func (tok Token) IsLiteral() bool {
	return literal_beg < tok && tok < literal_end
}

// 关键字
func (tok Token) IsKeyword() bool {
	return keyword_beg < tok && tok < keyword_end
}

// 特殊符号
func (tok Token) IsOperator() bool {
	return operator_beg < tok && tok < operator_end
}

func (tok Token) IsIsntruction() bool {
	if tok > ins_beg && tok < ins_end {
		return true
	}
	return tok == INSTRUCTION
}

// 标识符
func IsIdent(s string) bool {
	return strings.HasPrefix(s, "$")
}
