// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

// 注意: 禁止依赖 riscv 等子包

import (
	"strconv"

	"wa-lang.org/wa/internal/native/abi"
)

// 记号类型
type Token int

const (
	NONE    Token = iota // 零值, 表示空
	ILLEGAL              // 非法
	EOF                  // 结尾
	COMMENT              // 注释

	literal_beg // 面值开始
	IDENT       // 标识符, 以 $ 或 % 开头, $ 表示全局变量
	REG         // 寄存器
	INST        // 机器指令
	INT         // 12345
	FLOAT       // 123.45
	CHAR        // 'a'
	STRING      // "abc"
	literal_end // 面值结束

	operator_beg // 运算符开始
	ADD          // +
	SUB          // -
	ASSIGN       // =
	ARROW        // =>
	COLON        // :
	COMMA        // ,
	SEMICOLON    // ;
	LPAREN       // (
	RPAREN       // )
	LBRACK       // [
	RBRACK       // ]
	LBRACE       // {
	RBRACE       // }
	operator_end // 运算符结束

	keyword_beg // 关键字开始

	// 英文版关键字

	BYTE   // byte
	SHORT  // short
	LONG   // long
	QUAD   // quad
	I32    // int32
	I64    // int64
	U32    // uint32
	U64    // uint64
	F32    // float32
	F64    // float64
	CONST  // 常量
	GLOBAL // 全局
	LOCAL  // 局部
	FUNC   // 函数

	// 中文版关键字

	I32_zh    // 普整
	I64_zh    // 长整
	U32_zh    // 普正
	U64_zh    // 长正
	F32_zh    // 单精
	F64_zh    // 双精
	CONST_zh  // 常量
	GLOBAL_zh // 全局
	LOCAL_zh  // 局部
	FUNC_zh   // 函数
	END_zh    // 完毕

	// gas 汇编语法关键字(子集)

	GAS_X64_INTEL_SYNTAX // .intel_syntax noprefix, x64 专有
	GAS_X64_NOPREFIX     // .intel_syntax noprefix, x64 专有

	GAS_EXTERN  // .extern _write
	GAS_SET     // .set .Wa.Runtime.write, _write
	GAS_ALIGN   // .align 8
	GAS_GLOBA   // .globl .Wa.Memory.addr
	GAS_BYTE    // .name: .byte 0
	GAS_SHORT   // .name: .short 0
	GAS_LONG    // .name: .long 0
	GAS_QUAD    // .name: .quad 0
	GAS_ASSCII  // .name: .ascii "abc\000"
	GAS_ASSCIZ  // .name: .asciz "abc"
	GAS_SECTION // .section .text

	keyword_end // 关键字结束
)

// 寄存器到 Token 空间的映射
const (
	// 寄存器编号空间
	// 每个平台不超过 100 个, 至少保证 10 个独立空间
	REG_LOONG_BEGIN Token = 1000 + 100*iota
	REG_RISCV_BEGIN

	REG_BEGIN     = REG_LOONG_BEGIN
	REG_LOONG_END = REG_LOONG_BEGIN + 100
	REG_RISCV_END = REG_RISCV_BEGIN + 100
	REG_END       = REG_RISCV_END
)

// 指令到 Token 空间的映射
const (
	// 指令编号空间
	// 每个平台不超过 2000 个, 至少保证 10 个独立空间
	A_LOONG_BEGIN Token = 2000 + 2000*iota
	A_RISCV_BEGIN
	A_END

	A_BEGIN     = A_LOONG_BEGIN
	A_LOONG_END = A_LOONG_BEGIN + 2000
	A_RISCV_END = A_RISCV_BEGIN + 2000
)

var tokens = [...]string{
	NONE:    "NONE",
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	REG:    "REG",
	INST:   "INST",
	INT:    "INT",
	FLOAT:  "FLOAT",
	CHAR:   "CHAR",
	STRING: "STRING",

	ADD:       "+",
	ASSIGN:    "=",
	ARROW:     "=>",
	COLON:     ":",
	COMMA:     ",",
	SEMICOLON: ";",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACK:    "[",
	RBRACK:    "]",
	LBRACE:    "{",
	RBRACE:    "}",

	I32:    "i32",
	I64:    "i64",
	U32:    "u32",
	U64:    "u64",
	F32:    "f32",
	F64:    "f64",
	CONST:  "const",
	GLOBAL: "global",
	LOCAL:  "local",
	FUNC:   "func",

	I32_zh:    "普整",
	I64_zh:    "长整",
	U32_zh:    "普正",
	U64_zh:    "长正",
	F32_zh:    "单精",
	F64_zh:    "双精",
	CONST_zh:  "常量",
	GLOBAL_zh: "全局",
	LOCAL_zh:  "局部",
	FUNC_zh:   "函数",
	END_zh:    "完毕",

	GAS_X64_INTEL_SYNTAX: ".intel_syntax",
	GAS_X64_NOPREFIX:     "noprefix",

	GAS_EXTERN:  ".extern",
	GAS_SET:     ".set",
	GAS_ALIGN:   ".align",
	GAS_GLOBA:   ".globl",
	GAS_BYTE:    ".byte",
	GAS_SHORT:   ".short",
	GAS_LONG:    ".long",
	GAS_QUAD:    ".quad",
	GAS_ASSCII:  ".ascii",
	GAS_ASSCIZ:  ".asciz",
	GAS_SECTION: ".section",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "native.Token(" + strconv.Itoa(int(tok)) + ")"
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

// 查询标识符的类型
// 寄存器和指令负责和 Token 名字空间的映射关系, 查询失败返回 0
func Lookup(ident string, lookupRegisterOrAs func(ident string) Token) Token {
	// 关键字类型
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}

	// 寄存器或汇编指令
	if lookupRegisterOrAs != nil {
		if tok := lookupRegisterOrAs(ident); tok != 0 {
			return tok
		}
	}

	// 函数名
	return IDENT
}

// 面值类型
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// 运算符类型
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// 关键字
func (tok Token) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }

// Gas 关键字
func (tok Token) IsGasKeyword() bool {
	if tok.IsKeyword() {
		return tok == GAS_X64_NOPREFIX || tokens[tok][0] == '.'
	}
	return false
}

// 寄存器
func (tok Token) IsRegister() bool { return REG_BEGIN <= tok && tok < REG_END }

// 指令
func (tok Token) IsAs() bool { return A_BEGIN <= tok && tok < A_END }

// 原始寄存器值
func (tok Token) RawReg() abi.RegType {
	if REG_LOONG_BEGIN <= tok && tok < REG_LOONG_END {
		return abi.RegType(tok - REG_LOONG_BEGIN)
	}
	if REG_RISCV_BEGIN <= tok && tok < REG_RISCV_END {
		return abi.RegType(tok - REG_RISCV_BEGIN)
	}
	return 0
}

// 原始指令值
func (tok Token) RawAs() abi.As {
	if A_LOONG_BEGIN <= tok && tok < A_LOONG_END {
		return abi.As(tok - A_LOONG_BEGIN)
	}
	if A_RISCV_BEGIN <= tok && tok < A_RISCV_END {
		return abi.As(tok - A_RISCV_BEGIN)
	}
	return 0
}

// 默认的数值类型
func (tok Token) DefaultNumberType() Token {
	switch tok {
	case CHAR, INT:
		return I32
	case FLOAT:
		return F32
	default:
		return tok
	}
}

// 是否数值类型
func (tok Token) IsNumberType() bool {
	switch tok {
	case I32, I32_zh:
		return true
	case I64, I64_zh:
		return true
	case U32, U32_zh:
		return true
	case U64, U64_zh:
		return true
	case F32, F32_zh:
		return true
	case F64, F64_zh:
		return true
	default:
		return false
	}
}

// 数值类型的内存大小(不含指针和字符串类型)
func (tok Token) NumberTypeSize() Token {
	switch tok {
	case I32, I32_zh:
		return 4
	case I64, I64_zh:
		return 8
	case U32, U32_zh:
		return 4
	case U64, U64_zh:
		return 8
	case F32, F32_zh:
		return 4
	case F64, F64_zh:
		return 8
	default:
		panic("unreachable")
	}
}

// 是否是导出的符号
func IsExported(name string) bool {
	ch := name[0]
	isInternal := ch == '_' || (ch >= 'a' && ch <= 'z')
	return !isInternal
}
