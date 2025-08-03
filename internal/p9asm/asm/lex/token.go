// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package lex

import (
	"fmt"
	"strings"
	"text/scanner"
)

// 对应符号类别
// 一般由 text/scanner.Scanner 返回(一般是一个 rune), 外加自定义的扩展类型
type ScanToken rune

const (
	// 汇编中由一些2个符号构成的运算符, 用负数来定义为扩展类型
	LSH       ScanToken = -1000 - iota // << 左移
	RSH                                // >> 逻辑右移
	ARR                                // -> 数学右移, 用于 ARM 的第 3 类移动指令
	ROT                                // @> 循环右移, 用于 ARM 的第 4 类移动指令
	macroName                          // 宏的名字, 内部用
)

// 类别和文本构成一个词法符号
type Token struct {
	ScanToken ScanToken
	Text      string
}

// 解析字符串为 Token 列表
func LexString(str string) []Token {
	t := newTokenizer(nil, "command line", strings.NewReader(str), nil)
	var tokens []Token
	for {
		tok := t.Next()
		if tok == scanner.EOF {
			break
		}
		tokens = append(tokens, Token{tok, t.Text()})
	}
	return tokens
}

func (l Token) String() string {
	return l.Text
}

// IsRegisterShift reports whether the token is one of the ARM register shift operators.
func (r ScanToken) IsRegisterShift() bool {
	return ROT <= r && r <= LSH // Order looks backwards because these are negative.
}

func (t ScanToken) String() string {
	switch t {
	case scanner.EOF:
		return "EOF"
	case scanner.Ident:
		return "identifier"
	case scanner.Int:
		return "integer constant"
	case scanner.Float:
		return "float constant"
	case scanner.Char:
		return "rune constant"
	case scanner.String:
		return "string constant"
	case scanner.RawString:
		return "raw string constant"
	case scanner.Comment:
		return "comment"
	default:
		return fmt.Sprintf("%q", rune(t))
	}
}
