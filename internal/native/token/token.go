// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

// 注意: 禁止依赖 loong64/riscv 等子包

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
	IDENT       // 标识符
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
	COLON        // :
	COMMA        // ,
	SEMICOLON    // ;
	LPAREN       // (
	RPAREN       // )
	LBRACK       // [
	RBRACK       // ]
	operator_end // 运算符结束

	// 中文版关键字

	zh_keyword_beg // 关键字开始

	EXTERN_zh   // 声明
	CONST_zh    // 常量
	GLOBAL_zh   // 全局
	READONLY_zh // 只读
	FUNC_zh     // 函数
	END_zh      // 完毕
	BYTE_zh     // 字节
	SHORT_zh    // 单字
	LONG_zh     // 双字
	QUAD_zh     // 四字
	FLOAT_zh    // 单精
	DOUBLE_zh   // 双精
	ADDR_zh     // 地址
	ASCII_zh    // 字串

	zh_keyword_end // 关键字结束

	// gas 汇编语法关键字(子集)

	gas_keyword_beg // 关键字开始

	GAS_X64_INTEL_SYNTAX // .intel_syntax noprefix, x64 专有
	GAS_X64_NOPREFIX     // .intel_syntax noprefix, x64 专有

	GAS_EXTERN  // .extern _write
	GAS_ALIGN   // .align 8
	GAS_GLOBL   // .globl .Wa.Memory.addr # .global 是别名, 等价
	GAS_BYTE    // .name: .byte 0
	GAS_SHORT   // .name: .short 0
	GAS_LONG    // .name: .long 0
	GAS_QUAD    // .name: .quad 0
	GAS_FLOAT   // .name: .quad 0.0
	GAS_DOUBLE  // .name: .quad 0.0
	GAS_ASCII   // .name: .ascii "abc\000"
	GAS_SKIP    // .name: .skip 100
	GAS_INCBIN  // .name: .incbin "lena.jpg"
	GAS_SECTION // .section .text

	gas_keyword_end // 关键字结束
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
	COLON:     ":",
	COMMA:     ",",
	SEMICOLON: ";",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACK:    "[",
	RBRACK:    "]",

	EXTERN_zh:   "声明",
	CONST_zh:    "常量",
	GLOBAL_zh:   "全局",
	READONLY_zh: "只读",
	FUNC_zh:     "函数",
	END_zh:      "完毕",
	BYTE_zh:     "字节",
	SHORT_zh:    "单字",
	LONG_zh:     "双字",
	QUAD_zh:     "四字",
	FLOAT_zh:    "单精",
	DOUBLE_zh:   "双精",
	ADDR_zh:     "地址",
	ASCII_zh:    "字串",

	GAS_X64_INTEL_SYNTAX: ".intel_syntax",
	GAS_X64_NOPREFIX:     "noprefix",

	GAS_EXTERN:  ".extern",
	GAS_ALIGN:   ".align",
	GAS_GLOBL:   ".globl",
	GAS_BYTE:    ".byte",
	GAS_SHORT:   ".short",
	GAS_LONG:    ".long",
	GAS_QUAD:    ".quad",
	GAS_FLOAT:   ".float",
	GAS_DOUBLE:  ".double",
	GAS_ASCII:   ".ascii",
	GAS_SKIP:    ".skip",
	GAS_INCBIN:  ".incbin",
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

var gas_keywords map[string]Token
var zh_keywords map[string]Token

func init() {
	gas_keywords = make(map[string]Token)
	zh_keywords = make(map[string]Token)

	for i := gas_keyword_beg + 1; i < gas_keyword_end; i++ {
		gas_keywords[tokens[i]] = i
	}
	for i := zh_keyword_beg + 1; i < zh_keyword_end; i++ {
		zh_keywords[tokens[i]] = i
	}
}

// 查询标识符的类型
// 寄存器和指令负责和 Token 名字空间的映射关系, 查询失败返回 0
func Lookup(ident string, lookupRegisterOrAs func(ident string) Token) Token {
	// 特殊处理 .global 别名
	const gas_global = ".global"
	if ident == gas_global {
		return GAS_GLOBL
	}

	// 关键字类型
	if tok, is_keyword := gas_keywords[ident]; is_keyword {
		return tok
	}
	if tok, is_keyword := zh_keywords[ident]; is_keyword {
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
func (tok Token) IsKeyword() bool { return tok.IsGasKeyword() || tok.IsZhKeyword() }

// Gas 关键字
func (tok Token) IsGasKeyword() bool {
	return gas_keyword_beg < tok && tok < gas_keyword_end
}

// 中文关键字
func (tok Token) IsZhKeyword() bool {
	return zh_keyword_beg < tok && tok < zh_keyword_end
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

// 是否是导出的符号
func IsExported(name string) bool {
	ch := name[0]
	isInternal := ch == '_' || (ch >= 'a' && ch <= 'z')
	return !isInternal
}
