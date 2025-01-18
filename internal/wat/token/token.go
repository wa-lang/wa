// 版权 @2024 凹语言 作者。保留所有权利。

// wat记号类型
package token

import (
	"strconv"
	"strings"
)

// 记号类型
type Token int

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
	FUNCREF // funcref
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

	keyword_end

	// 指令是关键字, 同时有指令码
	// 指令码需要单独查表计算, 不能通过 instruction_beg 计算
	// https://www.w3.org/TR/2019/REC-wasm-core-1-20191205/

	instruction_beg

	INS_UNREACHABLE         // 0x00, unreachable
	INS_NOP                 // 0x01, nop
	INS_BLOCK               // 0x02, block
	INS_LOOP                // 0x03, loop
	INS_IF                  // 0x04, if
	INS_ELSE                // 0x05, else
	INS_END                 // 0x0b, end
	INS_BR                  // 0x0c, br
	INS_BR_IF               // 0x0d, br_if
	INS_BR_TABLE            // 0x0e, br_table
	INS_RETURN              // 0x0f, return
	INS_CALL                // 0x10, call
	INS_CALL_INDIRECT       // 0x11, call_indirect
	INS_DROP                // 0x1a, drop
	INS_SELECT              // 0x1b, select
	INS_LOCAL_GET           // 0x20, local.get
	INS_LOCAL_SET           // 0x21, local.set
	INS_LOCAL_TEE           // 0x22, local.tee
	INS_GLOBAL_GET          // 0x23, global.get
	INS_GLOBAL_SET          // 0x24, global.set
	INS_TABLE_GET           // 0x25, table.get
	INS_TABLE_SET           // 0x26, table.set
	INS_I32_LOAD            // 0x28, i32.load
	INS_I64_LOAD            // 0x29, i64.load
	INS_F32_LOAD            // 0x2a, f32.load
	INS_F64_LOAD            // 0x2b, f64.load
	INS_I32_LOAD8_S         // 0x2c, i32.load8_s
	INS_I32_LOAD8_U         // 0x2d, i32.load8_u
	INS_I32_LOAD16_S        // 0x2e, i32.load16_s
	INS_I32_LOAD16_U        // 0x2f, i32.load16_u
	INS_I64_LOAD8_S         // 0x30, i64.load8_s
	INS_I64_LOAD8_U         // 0x31, i64.load8_u
	INS_I64_LOAD16_S        // 0x31, i64.load16_s
	INS_I64_LOAD16_U        // 0x33, i64.load16_u
	INS_I64_LOAD32_S        // 0x34, i64.load32_s
	INS_I64_LOAD32_U        // 0x35, i64.load32_u
	INS_I32_STORE           // 0x36, i32.store
	INS_I64_STORE           // 0x37, i64.store
	INS_F32_STORE           // 0x38, f32.store
	INS_F64_STORE           // 0x39, f64.store
	INS_I32_STORE8          // 0x3a, i32.store8
	INS_I32_STORE16         // 0x3b, i32.store16
	INS_I64_STORE8          // 0x3c, i64.store8
	INS_I64_STORE16         // 0x3d, i64.store16
	INS_I64_STORE32         // 0x3e, i64.store32
	INS_MEMORY_SIZE         // 0x3f, memory.size
	INS_MEMORY_GROW         // 0x40, memory.grow
	INS_MEMORY_INIT         // 0xfc 0x08, memory.init
	INS_MEMORY_COPY         // 0xfc 0x0a, memory.copy
	INS_MEMORY_FILL         // 0xfc 0x0b, memory.fill
	INS_I32_CONST           // 0x41, i32.const
	INS_I64_CONST           // 0x42, i64.const
	INS_F32_CONST           // 0x43, f32.const
	INS_F64_CONST           // 0x44, f64.const
	INS_I32_EQZ             // 0x45, i32.eqz
	INS_I32_EQ              // 0x46, i32.eq
	INS_I32_NE              // 0x47, i32.ne
	INS_I32_LT_S            // 0x48, i32.lt_s
	INS_I32_LT_U            // 0x49, i32.lt_u
	INS_I32_GT_S            // 0x4a, i32.gt_s
	INS_I32_GT_U            // 0x4b, i32.gt_u
	INS_I32_LE_S            // 0x4c, i32.le_s
	INS_I32_LE_U            // 0x4d, i32.le_u
	INS_I32_GE_S            // 0x4e, i32.ge_s
	INS_I32_GE_U            // 0x4f, i32.ge_u
	INS_I64_EQZ             // 0x50, i64.eqz
	INS_I64_EQ              // 0x51, i64.eq
	INS_I64_NE              // 0x52, i64.ne
	INS_I64_LT_S            // 0x53, i64.lt_s
	INS_I64_LT_U            // 0x54, i64.lt_u
	INS_I64_GT_S            // 0x55, i64.gt_s
	INS_I64_GT_U            // 0x56, i64.gt_u
	INS_I64_LE_S            // 0x57, i64.le_s
	INS_I64_LE_U            // 0x58, i64.le_u
	INS_I64_GE_S            // 0x59, i64.ge_s
	INS_I64_GE_U            // 0x5a, i64.ge_u
	INS_F32_EQ              // 0x5b, f32.eq
	INS_F32_NE              // 0x5c, f32.ne
	INS_F32_LT              // 0x5d, f32.lt
	INS_F32_GT              // 0x5e, f32.gt
	INS_F32_LE              // 0x5f, f32.le
	INS_F32_GE              // 0x60, f32.ge
	INS_F64_EQ              // 0x61, f64.eq
	INS_F64_NE              // 0x62, f64.ne
	INS_F64_LT              // 0x63, f64.lt
	INS_F64_GT              // 0x64, f64.gt
	INS_F64_LE              // 0x65, f64.le
	INS_F64_GE              // 0x66, f64.ge
	INS_I32_CLZ             // 0x67, i32.clz
	INS_I32_CTZ             // 0x68, i32.ctz
	INS_I32_POPCNT          // 0x69, i32.popcnt
	INS_I32_ADD             // 0x6a, i32.add
	INS_I32_SUB             // 0x6b, i32.sub
	INS_I32_MUL             // 0x6c, i32.mul
	INS_I32_DIV_S           // 0x6d, i32.div_s
	INS_I32_DIV_U           // 0x6e, i32.div_u
	INS_I32_REM_S           // 0x6f, i32.rem_s
	INS_I32_REM_U           // 0x70, i32.rem_u
	INS_I32_AND             // 0x71, i32.and
	INS_I32_OR              // 0x72, i32.or
	INS_I32_XOR             // 0x73, i32.xor
	INS_I32_SHL             // 0x74, i32.shl
	INS_I32_SHR_S           // 0x75, i32.shr_s
	INS_I32_SHR_U           // 0x76, i32.shr_u
	INS_I32_ROTL            // 0x77, i32.rotl
	INS_I32_ROTR            // 0x78, i32.rotr
	INS_I64_CLZ             // 0x79, i64.clz
	INS_I64_CTZ             // 0x7a, i64.ctz
	INS_I64_POPCNT          // 0x7b, i64.popcnt
	INS_I64_ADD             // 0x7c, i64.add
	INS_I64_SUB             // 0x7d, i64.sub
	INS_I64_MUL             // 0x7e, i64.mul
	INS_I64_DIV_S           // 0x7f, i64.div_s
	INS_I64_DIV_U           // 0x80, i64.div_u
	INS_I64_REM_S           // 0x81, i64.rem_s
	INS_I64_REM_U           // 0x82, i64.rem_u
	INS_I64_AND             // 0x83, i64.and
	INS_I64_OR              // 0x84, i64.or
	INS_I64_XOR             // 0x85, i64.xor
	INS_I64_SHL             // 0x86, i64.shl
	INS_I64_SHR_S           // 0x87, i64.shr_s
	INS_I64_SHR_U           // 0x88, i64.shr_u
	INS_I64_ROTL            // 0x89, i64.rotl
	INS_I64_ROTR            // 0x8a, i64.rotr
	INS_F32_ABS             // 0x8b, f32.abs
	INS_F32_NEG             // 0x8c, f32.neg
	INS_F32_CEIL            // 0x8d, f32.ceil
	INS_F32_FLOOR           // 0x8e, f32.floor
	INS_F32_TRUNC           // 0x8f, f32.trunc
	INS_F32_NEAREST         // 0x90, f32.nearest
	INS_F32_SQRT            // 0x91, f32.sqrt
	INS_F32_ADD             // 0x92, f32.add
	INS_F32_SUB             // 0x93, f32.sub
	INS_F32_MUL             // 0x94, f32.mul
	INS_F32_DIV             // 0x95, f32.div
	INS_F32_MIN             // 0x96, f32.min
	INS_F32_MAX             // 0x97, f32.max
	INS_F32_COPYSIGN        // 0x98, f32.copysign
	INS_F64_ABS             // 0x99, f64.abs
	INS_F64_NEG             // 0x9a, f64.neg
	INS_F64_CEIL            // 0x9b, f64.ceil
	INS_F64_FLOOR           // 0x9c, f64.floor
	INS_F64_TRUNC           // 0x9d, f64.trunc
	INS_F64_NEAREST         // 0x9e, f64.nearest
	INS_F64_SQRT            // 0x9f, f64.sqrt
	INS_F64_ADD             // 0xa0, f64.add
	INS_F64_SUB             // 0xa1, f64.sub
	INS_F64_MUL             // 0xa2, f64.mul
	INS_F64_DIV             // 0xa3, f64.div
	INS_F64_MIN             // 0xa4, f64.min
	INS_F64_MAX             // 0xa5, f64.max
	INS_F64_COPYSIGN        // 0xa6, f64.copysign
	INS_I32_WRAP_I64        // 0xa7, i32.wrap_i64
	INS_I32_TRUNC_F32_S     // 0xa8, i32.trunc_f32_s
	INS_I32_TRUNC_F32_U     // 0xa9, i32.trunc_f32_u
	INS_I32_TRUNC_F64_S     // 0xaa, i32.trunc_f64_s
	INS_I32_TRUNC_F64_U     // 0xab, i32.trunc_f64_u
	INS_I64_EXTEND_I32_S    // 0xac, i64.extend_i32_s
	INS_I64_EXTEND_I32_U    // 0xad, i64.extend_i32_u
	INS_I64_TRUNC_F32_S     // 0xae, i64.trunc_f32_s
	INS_I64_TRUNC_F32_U     // 0xaf, i64.trunc_f32_u
	INS_I64_TRUNC_F64_S     // 0xb0, i64.trunc_f64_s
	INS_I64_TRUNC_F64_U     // 0xb1, i64.trunc_f64_u
	INS_F32_CONVERT_I32_S   // 0xb2, f32.convert_i32_s
	INS_F32_CONVERT_I32_U   // 0xb3, f32.convert_i32_u
	INS_F32_CONVERT_I64_S   // 0xb4, f32.convert_i64_s
	INS_F32_CONVERT_I64_U   // 0xb5, f32.convert_i64_u
	INS_F32_DEMOTE_F64      // 0xb6, f32.demote_f64
	INS_F64_CONVERT_I32_S   // 0xb7, f64.convert_i32_s
	INS_F64_CONVERT_I32_U   // 0xb8, f64.convert_i32_u
	INS_F64_CONVERT_I64_S   // 0xb9, f64.convert_i64_s
	INS_F64_CONVERT_I64_U   // 0xba, f64.convert_i64_u
	INS_F64_PROMOTE_F32     // 0xbb, f64.promote_f32
	INS_I32_REINTERPRET_F32 // 0xbc, i32.reinterpret_f32
	INS_I64_REINTERPRET_F64 // 0xbd, i64.reinterpret_f64
	INS_F32_REINTERPRET_I32 // 0xbe, f32.reinterpret_i32
	INS_F64_REINTERPRET_I64 // 0xbf, f64.reinterpret_i64

	instruction_end
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
	FUNCREF: "funcref",
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

	INS_UNREACHABLE:         "unreachable",
	INS_NOP:                 "nop",
	INS_BLOCK:               "block",
	INS_LOOP:                "loop",
	INS_IF:                  "if",
	INS_ELSE:                "else",
	INS_END:                 "end",
	INS_BR:                  "br",
	INS_BR_IF:               "br_if",
	INS_BR_TABLE:            "br_table",
	INS_RETURN:              "return",
	INS_CALL:                "call",
	INS_CALL_INDIRECT:       "call_indirect",
	INS_DROP:                "drop",
	INS_SELECT:              "select",
	INS_LOCAL_GET:           "local.get",
	INS_LOCAL_SET:           "local.set",
	INS_LOCAL_TEE:           "local.tee",
	INS_GLOBAL_GET:          "global.get",
	INS_GLOBAL_SET:          "global.set",
	INS_TABLE_GET:           "table.get",
	INS_TABLE_SET:           "table.set",
	INS_I32_LOAD:            "i32.load",
	INS_I64_LOAD:            "i64.load",
	INS_F32_LOAD:            "f32.load",
	INS_F64_LOAD:            "f64.load",
	INS_I32_LOAD8_S:         "i32.load8_s",
	INS_I32_LOAD8_U:         "i32.load8_u",
	INS_I32_LOAD16_S:        "i32.load16_s",
	INS_I32_LOAD16_U:        "i32.load16_u",
	INS_I64_LOAD8_S:         "i64.load8_s",
	INS_I64_LOAD8_U:         "i64.load8_u",
	INS_I64_LOAD16_S:        "i64.load16_s",
	INS_I64_LOAD16_U:        "i64.load16_u",
	INS_I64_LOAD32_S:        "i64.load32_s",
	INS_I64_LOAD32_U:        "i64.load32_u",
	INS_I32_STORE:           "i32.store",
	INS_I64_STORE:           "i64.store",
	INS_F32_STORE:           "f32.store",
	INS_F64_STORE:           "f64.store",
	INS_I32_STORE8:          "i32.store8",
	INS_I32_STORE16:         "i32.store16",
	INS_I64_STORE8:          "i64.store8",
	INS_I64_STORE16:         "i64.store16",
	INS_I64_STORE32:         "i64.store32",
	INS_MEMORY_SIZE:         "memory.size",
	INS_MEMORY_GROW:         "memory.grow",
	INS_MEMORY_INIT:         "memory.init",
	INS_MEMORY_COPY:         "memory.copy",
	INS_MEMORY_FILL:         "memory.fill",
	INS_I32_CONST:           "i32.const",
	INS_I64_CONST:           "i64.const",
	INS_F32_CONST:           "f32.const",
	INS_F64_CONST:           "f64.const",
	INS_I32_EQZ:             "i32.eqz",
	INS_I32_EQ:              "i32.eq",
	INS_I32_NE:              "i32.ne",
	INS_I32_LT_S:            "i32.lt_s",
	INS_I32_LT_U:            "i32.lt_u",
	INS_I32_GT_S:            "i32.gt_s",
	INS_I32_GT_U:            "i32.gt_u",
	INS_I32_LE_S:            "i32.le_s",
	INS_I32_LE_U:            "i32.le_u",
	INS_I32_GE_S:            "i32.ge_s",
	INS_I32_GE_U:            "i32.ge_u",
	INS_I64_EQZ:             "i64.eqz",
	INS_I64_EQ:              "i64.eq",
	INS_I64_NE:              "i64.ne",
	INS_I64_LT_S:            "i64.lt_s",
	INS_I64_LT_U:            "i64.lt_u",
	INS_I64_GT_S:            "i64.gt_s",
	INS_I64_GT_U:            "i64.gt_u",
	INS_I64_LE_S:            "i64.le_s",
	INS_I64_LE_U:            "i64.le_u",
	INS_I64_GE_S:            "i64.ge_s",
	INS_I64_GE_U:            "i64.ge_u",
	INS_F32_EQ:              "f32.eq",
	INS_F32_NE:              "f32.ne",
	INS_F32_LT:              "f32.lt",
	INS_F32_GT:              "f32.gt",
	INS_F32_LE:              "f32.le",
	INS_F32_GE:              "f32.ge",
	INS_F64_EQ:              "f64.eq",
	INS_F64_NE:              "f64.ne",
	INS_F64_LT:              "f64.lt",
	INS_F64_GT:              "f64.gt",
	INS_F64_LE:              "f64.le",
	INS_F64_GE:              "f64.ge",
	INS_I32_CLZ:             "i32.clz",
	INS_I32_CTZ:             "i32.ctz",
	INS_I32_POPCNT:          "i32.popcnt",
	INS_I32_ADD:             "i32.add",
	INS_I32_SUB:             "i32.sub",
	INS_I32_MUL:             "i32.mul",
	INS_I32_DIV_S:           "i32.div_s",
	INS_I32_DIV_U:           "i32.div_u",
	INS_I32_REM_S:           "i32.rem_s",
	INS_I32_REM_U:           "i32.rem_u",
	INS_I32_AND:             "i32.and",
	INS_I32_OR:              "i32.or",
	INS_I32_XOR:             "i32.xor",
	INS_I32_SHL:             "i32.shl",
	INS_I32_SHR_S:           "i32.shr_s",
	INS_I32_SHR_U:           "i32.shr_u",
	INS_I32_ROTL:            "i32.rotl",
	INS_I32_ROTR:            "i32.rotr",
	INS_I64_CLZ:             "i64.clz",
	INS_I64_CTZ:             "i64.ctz",
	INS_I64_POPCNT:          "i64.popcnt",
	INS_I64_ADD:             "i64.add",
	INS_I64_SUB:             "i64.sub",
	INS_I64_MUL:             "i64.mul",
	INS_I64_DIV_S:           "i64.div_s",
	INS_I64_DIV_U:           "i64.div_u",
	INS_I64_REM_S:           "i64.rem_s",
	INS_I64_REM_U:           "i64.rem_u",
	INS_I64_AND:             "i64.and",
	INS_I64_OR:              "i64.or",
	INS_I64_XOR:             "i64.xor",
	INS_I64_SHL:             "i64.shl",
	INS_I64_SHR_S:           "i64.shr_s",
	INS_I64_SHR_U:           "i64.shr_u",
	INS_I64_ROTL:            "i64.rotl",
	INS_I64_ROTR:            "i64.rotr",
	INS_F32_ABS:             "f32.abs",
	INS_F32_NEG:             "f32.neg",
	INS_F32_CEIL:            "f32.ceil",
	INS_F32_FLOOR:           "f32.floor",
	INS_F32_TRUNC:           "f32.trunc",
	INS_F32_NEAREST:         "f32.nearest",
	INS_F32_SQRT:            "f32.sqrt",
	INS_F32_ADD:             "f32.add",
	INS_F32_SUB:             "f32.sub",
	INS_F32_MUL:             "f32.mul",
	INS_F32_DIV:             "f32.div",
	INS_F32_MIN:             "f32.min",
	INS_F32_MAX:             "f32.max",
	INS_F32_COPYSIGN:        "f32.copysign",
	INS_F64_ABS:             "f64.abs",
	INS_F64_NEG:             "f64.neg",
	INS_F64_CEIL:            "f64.ceil",
	INS_F64_FLOOR:           "f64.floor",
	INS_F64_TRUNC:           "f64.trunc",
	INS_F64_NEAREST:         "f64.nearest",
	INS_F64_SQRT:            "f64.sqrt",
	INS_F64_ADD:             "f64.add",
	INS_F64_SUB:             "f64.sub",
	INS_F64_MUL:             "f64.mul",
	INS_F64_DIV:             "f64.div",
	INS_F64_MIN:             "f64.min",
	INS_F64_MAX:             "f64.max",
	INS_F64_COPYSIGN:        "f64.copysign",
	INS_I32_WRAP_I64:        "i32.wrap_i64",
	INS_I32_TRUNC_F32_S:     "i32.trunc_f32_s",
	INS_I32_TRUNC_F32_U:     "i32.trunc_f32_u",
	INS_I32_TRUNC_F64_S:     "i32.trunc_f64_s",
	INS_I32_TRUNC_F64_U:     "i32.trunc_f64_u",
	INS_I64_EXTEND_I32_S:    "i64.extend_i32_s",
	INS_I64_EXTEND_I32_U:    "i64.extend_i32_u",
	INS_I64_TRUNC_F32_S:     "i64.trunc_f32_s",
	INS_I64_TRUNC_F32_U:     "i64.trunc_f32_u",
	INS_I64_TRUNC_F64_S:     "i64.trunc_f64_s",
	INS_I64_TRUNC_F64_U:     "i64.trunc_f64_u",
	INS_F32_CONVERT_I32_S:   "f32.convert_i32_s",
	INS_F32_CONVERT_I32_U:   "f32.convert_i32_u",
	INS_F32_CONVERT_I64_S:   "f32.convert_i64_s",
	INS_F32_CONVERT_I64_U:   "f32.convert_i64_u",
	INS_F32_DEMOTE_F64:      "f32.demote_f64",
	INS_F64_CONVERT_I32_S:   "f64.convert_i32_s",
	INS_F64_CONVERT_I32_U:   "f64.convert_i32_u",
	INS_F64_CONVERT_I64_S:   "f64.convert_i64_s",
	INS_F64_CONVERT_I64_U:   "f64.convert_i64_u",
	INS_F64_PROMOTE_F32:     "f64.promote_f32",
	INS_I32_REINTERPRET_F32: "i32.reinterpret_f32",
	INS_I64_REINTERPRET_F64: "i64.reinterpret_f64",
	INS_F32_REINTERPRET_I32: "f32.reinterpret_i32",
	INS_F64_REINTERPRET_I64: "f64.reinterpret_i64",
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

var (
	keywords     map[string]Token
	instructions map[string]Token
)

func init() {
	keywords = make(map[string]Token)
	instructions = make(map[string]Token)

	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}

	for i := instruction_beg + 1; i < instruction_end; i++ {
		instructions[tokens[i]] = i
	}
}

func Lookup(ident string) Token {
	// 标识符 $name
	if strings.HasPrefix(ident, "$") {
		return IDENT
	}

	// 关键字
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}

	// 指令
	if tok, is_ins := instructions[ident]; is_ins {
		return tok
	}

	return ILLEGAL
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
	return instruction_beg < tok && tok < instruction_end
}
