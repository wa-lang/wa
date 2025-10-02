// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package token defines constants representing the lexical tokens of the Go
// programming language and basic operations on tokens (printing, predicates).
package token

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

// Token is the set of lexical tokens of the Go programming language.
type Token int

// The list of tokens.
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	literal_beg
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	IMAG   // 123.45i
	CHAR   // 'a'
	STRING // "abc"
	literal_end

	operator_beg
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND // &&
	LOR  // ||
	INC  // ++
	DEC  // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ       // !=
	LEQ       // <=
	GEQ       // >=
	DEFINE    // :=
	ELLIPSIS  // ...
	SPACESHIP // <=>

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :
	ARROW     // =>
	operator_end

	keyword_beg
	// Keywords
	BREAK
	CASE
	CONST
	CONTINUE

	DEFAULT
	DEFER
	ELSE
	FOR

	FUNC
	GLOBAL
	IF
	IMPORT

	INTERFACE
	MAP
	PACKAGE
	RANGE
	RETURN

	STRUCT
	SWITCH
	TYPE
	VAR

	keyword_end

	w2_keyword_beg
	// Keywords

	Zh_未定义 // undefined

	Zh_真 // true
	Zh_假 // false

	Zh_引入 // import

	Zh_常量 // const
	Zh_定义 // const
	Zh_全局 // flobal

	Zh_算始 // func
	Zh_算终 // func {}
	Zh_函始 // func
	Zh_函终 // func {}

	Zh_若始 // if
	Zh_若另 // else if
	Zh_若否 // else
	Zh_若终 // {}

	Zh_岔始 // switch
	Zh_岔终 // {}
	Zh_岔道 // case
	Zh_主道 // default

	Zh_当始 // for
	Zh_当终 // {}
	Zh_跳出 // break
	Zh_继续 // continue

	Zh_类始 // struct
	Zh_类终 // {}

	Zh_返回 // return
	Zh_押后 // defer
	Zh_字典 // map

	Zh_类型 // type
	Zh_接口 // interface

	w2_keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	IMAG:   "IMAG",
	CHAR:   "CHAR",
	STRING: "STRING",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND:     "&",
	OR:      "|",
	XOR:     "^",
	SHL:     "<<",
	SHR:     ">>",
	AND_NOT: "&^",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	AND_ASSIGN:     "&=",
	OR_ASSIGN:      "|=",
	XOR_ASSIGN:     "^=",
	SHL_ASSIGN:     "<<=",
	SHR_ASSIGN:     ">>=",
	AND_NOT_ASSIGN: "&^=",

	LAND: "&&",
	LOR:  "||",
	INC:  "++",
	DEC:  "--",

	EQL:    "==",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "!",

	NEQ:       "!=",
	LEQ:       "<=",
	GEQ:       ">=",
	DEFINE:    ":=",
	ELLIPSIS:  "...",
	SPACESHIP: "<=>",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",
	ARROW:     "=>",

	BREAK:    "break",
	CASE:     "case",
	CONST:    "const",
	CONTINUE: "continue",

	DEFAULT: "default",
	DEFER:   "defer",
	ELSE:    "else",
	FOR:     "for",

	FUNC:   "func",
	GLOBAL: "global",
	IF:     "if",
	IMPORT: "import",

	INTERFACE: "interface",
	MAP:       "map",
	PACKAGE:   "package",
	RANGE:     "range",
	RETURN:    "return",

	STRUCT: "struct",
	SWITCH: "switch",
	TYPE:   "type",
	VAR:    "var",

	Zh_未定义: "未定义",

	Zh_真: "真",
	Zh_假: "假",

	Zh_引入: "引入",

	Zh_常量: "常量",
	Zh_定义: "定义",
	Zh_全局: "全局",

	Zh_算始: "算始",
	Zh_算终: "算终",
	Zh_函始: "函始",
	Zh_函终: "函终",

	Zh_若始: "若始",
	Zh_若另: "若另",
	Zh_若否: "若否",
	Zh_若终: "若终",

	Zh_岔始: "岔始",
	Zh_岔终: "岔终",
	Zh_岔道: "岔道",
	Zh_主道: "主道",

	Zh_当始: "当始",
	Zh_当终: "当终",
	Zh_跳出: "跳出",
	Zh_继续: "继续",

	Zh_类始: "类始",
	Zh_类终: "类终",

	Zh_返回: "返回",
	Zh_押后: "押后",
	Zh_字典: "字典",

	Zh_类型: "类型",
	Zh_接口: "接口",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
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

func (tok Token) WaGoString() string {
	return tok.String()
}

func (tok Token) WaZhString() string {
	return tok.String()
}

// A set of constants for precedence-based expression parsing.
// Non-operators have lowest precedence, followed by operators
// starting with precedence 1 up to unary operators. The highest
// precedence serves as "catch-all" precedence for selector,
// indexing, and other operator and delimiter tokens.
const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

// Precedence returns the operator precedence of the binary
// operator op. If op is not a binary operator, the result
// is LowestPrecedence.
func (op Token) Precedence() int {
	switch op {
	case LOR:
		return 1
	case LAND:
		return 2
	case EQL, NEQ, LSS, LEQ, GTR, GEQ, SPACESHIP:
		return 3
	case ADD, SUB, OR, XOR:
		return 4
	case MUL, QUO, REM, SHL, SHR, AND, AND_NOT:
		return 5
	}
	return LowestPrec
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)

	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}

	for i := w2_keyword_beg + 1; i < w2_keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

func LookupEx(ident string, useW2Mode bool) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		if useW2Mode {
			if tok.IsW2Keyword() {
				return tok
			}
		} else {
			if tok.IsKeyword() {
				return tok
			}
		}
	}
	return IDENT
}

// Predicates

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
func (tok Token) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }

func (tok Token) IsW2Keyword() bool { return w2_keyword_beg < tok && tok < w2_keyword_end }

// IsExported reports whether name is exported.
func IsExported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	isInternal := ch == '_' || (ch >= 'a' && ch <= 'z') || ch == utf8.RuneError
	return !isInternal
}

// IsKeyword reports whether name is a Go keyword, such as "func" or "return".
func IsKeyword(name string) bool {
	// TODO: opt: use a perfect hash function instead of a global map.
	_, ok := keywords[name]
	return ok
}

// IsIdentifier reports whether name is a Go identifier, that is, a non-empty
// string made up of letters, digits, and underscores, where the first character
// is not a digit. Keywords are not identifiers.
func IsIdentifier(name string) bool {
	for i, c := range name {
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}
	return name != "" && !IsKeyword(name)
}
