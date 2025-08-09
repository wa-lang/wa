// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

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
	// TODO
	// 参考lex包的移位运算符
	operator_end

	// 面值类型
	literal_beg
	IDENT  // 标识符, 比如 $name
	INT    // 12345
	FLOAT  // 123.45
	CHAR   // 'a'
	STRING // "abc"

	// 把指令当作一种面值
	AS // 汇编指令
	literal_end

	// 关键字
	keyword_beg

	// TODO

	TEXT
	GLOBL

	keyword_end

	// 指令是关键字, 同时有指令码
	// 前面需要空出一个空间, 确保后续的部分可以和指令机器码自然对应
	// 多普通的指令空间需要在这里划分好

	// 但是只能同时选择一个平台, 因为不同的平台指令名字可能相同
	// 只针对解析的限制

	instruction_beg

	// TODO

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

	// TODO
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
