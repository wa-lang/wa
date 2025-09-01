// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

// 注意: 禁止依赖 riscv 等子包

import (
	"strconv"
	"strings"

	"wa-lang.org/wa/internal/native/abi"
)

// 记号类型
type Token int

const (
	ILLEGAL Token = iota // 非法
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
	ASSIGN       // =
	ARROW        // =>
	COLON        // :
	COMMA        // ,
	SEMICOLON    // ;
	LPAREN       // (
	RPAREN       // )
	LBRACE       // {
	RBRACE       // }
	operator_end // 运算符结束

	keyword_beg // 关键字开始
	I32         // int32
	I64         // int64
	F32         // float32
	F64         // float64
	PTR         // 地址类型
	CONST       // 常量
	GLOBAL      // 全局符号
	FUNC        // 函数
	keyword_end // 关键字结束
)

// 寄存器到 Token 空间的映射
const (
	// 寄存器编号空间
	// 每个平台不超过 100 个, 至少保证 10 个独立空间
	REG_RISCV_BEGIN Token = 1000 + 100*iota

	REG_BEGIN     = REG_RISCV_BEGIN
	REG_RISCV_END = REG_RISCV_BEGIN + 100
	REG_END       = REG_RISCV_END
)

// 指令到 Token 空间的映射
const (
	// 指令编号空间
	// 每个平台不超过 2000 个, 至少保证 10 个独立空间
	A_RISCV_BEGIN Token = 2000 + 2000*iota

	A_BEGIN     = A_RISCV_BEGIN
	A_RISCV_END = A_RISCV_BEGIN + 2000
	A_END       = A_RISCV_END
)

var tokens = [...]string{
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
	LBRACE:    "{",
	RBRACE:    "}",

	I32:    "i32",
	I64:    "i64",
	F32:    "f32",
	F64:    "f64",
	PTR:    "ptr",
	CONST:  "const",
	GLOBAL: "global",
	FUNC:   "func",
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
func Lookup(ident string,
	lookupRegister func(ident string) Token,
	lookupAs func(ident string) Token,
) Token {
	// 因为标识符采用特殊的前缀, 处理的流程可以再简化

	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	if lookupRegister != nil {
		if reg := lookupRegister(ident); reg != 0 {
			return reg
		}
	}
	if lookupAs != nil {
		if as := lookupAs(ident); as != 0 {
			return as
		}
	}

	// 标识符以 $ 或 % 开头
	if strings.HasPrefix(ident, "$") || strings.HasPrefix(ident, "%") {
		return IDENT
	}

	// 失败
	return ILLEGAL
}

// 面值类型
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// 运算符类型
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// 关键字
func (tok Token) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }

// 寄存器
func (tok Token) IsRegister() bool { return REG_BEGIN <= tok && tok < REG_END }

// 指令
func (tok Token) IsAs() bool { return A_BEGIN <= tok && tok < A_END }

// 原始寄存器值
func (tok Token) RawReg() abi.RegType {
	if REG_RISCV_BEGIN <= tok && tok < REG_RISCV_END {
		return abi.RegType(tok - REG_RISCV_BEGIN)
	}
	return 0
}

// 原始指令值
func (tok Token) RawAs() abi.As {
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
