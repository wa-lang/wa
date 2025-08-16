// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import "strconv"

// 记号类型
type Token int

const (
	ILLEGAL Token = iota // 非法
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
	MUL          // *
	QUO          // /
	REM          // %

	LSH // << 左移
	RSH // >> 逻辑右移
	ARR // -> 数学右移, 用于 ARM 的第 3 类移动指令
	ROT // @> 循环右移, 用于 ARM 的第 4 类移动指令

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN       // )
	RBRACK       // ]
	RBRACE       // }
	SEMICOLON    // ;
	COLON        // :
	operator_end // 运算符结束

	keyword_beg // 关键字开始
	GLOBAL      // 全局符号
	FUNC        // 函数
	keyword_end // 关键字结束
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

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	LSH: "<<",
	RSH: ">>",
	ARR: "->",
	ROT: "@>",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",

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
func Lookup(ident string, isRegister, isInstruction func(ident string) bool) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	if isRegister != nil && isRegister(ident) {
		return REG
	}
	if isInstruction != nil && isInstruction(ident) {
		return INST
	}
	return IDENT
}

// 面值类型
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// 运算符类型
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// 关键字
func (tok Token) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }

// 是否是导出的符号
func IsExported(name string) bool {
	ch := name[0]
	isInternal := ch == '_' || (ch >= 'a' && ch <= 'z')
	return !isInternal
}
