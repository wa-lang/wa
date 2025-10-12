// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package token defines constants representing the lexical tokens of the Go
// programming language and basic operations on tokens (printing, predicates).
package token

import (
	"strconv"
	"strings"
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

	wz_keyword_beg
	// Keywords
	Zh_引入 // import

	Zh_常量 // const
	Zh_全局 // global
	Zh_类型 // type
	Zh_函数 // func

	Zh_设定 // var

	Zh_结构 // struct
	Zh_字典 // map
	Zh_接口 // interface

	Zh_如果 // if
	Zh_或者 // else if [新加]
	Zh_否则 // else

	Zh_找辙 // switch
	Zh_有辙 // case
	Zh_没辙 // default

	Zh_循环 // for
	Zh_迭代 // range
	Zh_继续 // continue
	Zh_跳出 // break

	Zh_押后 // defer
	Zh_返回 // return

	Zh_区块 // { [新加]
	Zh_完毕 // } [新加]
	wz_keyword_end
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

	Zh_引入: K_引入,

	Zh_常量: K_常量,
	Zh_全局: K_全局,
	Zh_类型: K_类型,
	Zh_函数: K_函数,

	Zh_设定: K_设定,

	Zh_结构: K_结构,
	Zh_字典: K_字典,
	Zh_接口: K_接口,

	Zh_如果: K_如果,
	Zh_或者: K_或者,
	Zh_否则: K_否则,

	Zh_找辙: K_找辙,
	Zh_有辙: K_有辙,
	Zh_没辙: K_没辙,

	Zh_循环: K_循环,
	Zh_迭代: K_迭代,
	Zh_继续: K_继续,
	Zh_跳出: K_跳出,

	Zh_押后: K_押后,
	Zh_返回: K_返回,

	Zh_区块: K_区块,
	Zh_完毕: K_完毕,
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

	for i := wz_keyword_beg + 1; i < wz_keyword_end; i++ {
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

func LookupEx(ident string, useWzMode bool) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		if useWzMode {
			if tok.IsWzKeyword() {
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

func (tok Token) IsWzKeyword() bool { return wz_keyword_beg < tok && tok < wz_keyword_end }

// 是否为中文语法的注释
func (tok Token) IsWzComment(lit string) bool {
	return tok == COMMENT && strings.HasPrefix(lit, K_注)
}

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
