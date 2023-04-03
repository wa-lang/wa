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

	watoken "wa-lang.org/wa/internal/token"
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

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // :=
	ELLIPSIS // ...

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

	TK_左方 // 【
	TK_右方 // 】
	TK_冒号 // ：
	TK_句号 // 。
	TK_顿号 // 、
	TK_叹号 // ！
	TK_点号 // ·
	TK_左书 // 《
	TK_右书 // 》
	TK_逗号 // ，

	// 类型
	TK_数
	TK_浮点数
	TK_布尔

	operator_end

	keyword_beg

	TK_引于
	TK_归于
	TK_若
	TK_则
	TK_否则
	TK_设
	TK_之
	TK_从
	TK_到
	TK_直
	TK_自
	TK_至
	TK_有
	TK_当
	TK_为

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

	// reserved keywords
	CLASS
	ENUM
	keyword_end
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

	NEQ:      "!=",
	LEQ:      "<=",
	GEQ:      ">=",
	DEFINE:   ":=",
	ELLIPSIS: "...",

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

	TK_左方: "【",
	TK_右方: "】",
	TK_冒号: "：",
	TK_句号: "。",
	TK_顿号: "、",
	TK_叹号: "！",
	TK_点号: "·",
	TK_左书: "《",
	TK_右书: "》",
	TK_逗号: "，",

	TK_引于: "引于",
	TK_归于: "归于",
	TK_若:  "若",
	TK_则:  "则",
	TK_否则: "否则",
	TK_设:  "设",
	TK_之:  "之",
	TK_从:  "从",
	TK_到:  "到",
	TK_直:  "直",
	TK_自:  "自",
	TK_至:  "至",
	TK_有:  "有",
	TK_当:  "当",
	TK_为:  "为",

	BREAK:    "break",
	CASE:     "case",
	CONST:    "const",
	CONTINUE: "continue",

	DEFAULT: "default",
	DEFER:   "defer",
	ELSE:    "else",
	FOR:     "for",

	FUNC:   "func",
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

	CLASS: "class",
	ENUM:  "enum",
}

// WaGo 关键字补丁
var tokens_wago = map[Token]string{
	FUNC: "func",
}

// 中文关键字(最终通过海选选择前几名, 凹开发者最终决定)
var tokens_zh = map[Token]string{
	IMPORT: "导入",

	CONST: "常量",
	VAR:   "变量",

	TYPE:      "类型",
	MAP:       "字典",
	STRUCT:    "结构",
	CLASS:     "类", // 有 2 个字吗
	ENUM:      "枚举",
	INTERFACE: "接口",

	FUNC:   "函数",
	DEFER:  "善后",
	RETURN: "返回",

	IF:       "如果",
	ELSE:     "否则",
	FOR:      "循环",
	RANGE:    "区间",
	CONTINUE: "继续",
	BREAK:    "跳出",

	SWITCH:  "找辙",
	CASE:    "有辙",
	DEFAULT: "没辙",

	PACKAGE: "包", // 废弃
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
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		if tok > keyword_beg && tok < keyword_end {
			s = tokens_zh[tok]
		} else {
			s = tokens[tok]
		}
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

// 返回 WaGo 关键字名字
func (tok Token) WaGoKeykword() string {
	if tok.IsKeyword() {
		if x, ok := tokens_wago[tok]; ok {
			return x
		}
		return tokens[tok]
	}
	return ""
}

// 返回中文关键字名字
func (tok Token) ZhKeykword() string {
	if tok.IsKeyword() {
		return tokens_zh[tok]
	}
	return ""
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
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 3
	case ADD, SUB, OR, XOR:
		return 4
	case MUL, QUO, REM, SHL, SHR, AND, AND_NOT:
		return 5
	}
	return LowestPrec
}

var (
	keywords    map[string]Token
	keywords_zh map[string]Token
)

func init() {
	keywords = make(map[string]Token)
	keywords_zh = make(map[string]Token)

	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
	for k, name := range tokens_zh {
		keywords_zh[name] = k
	}
}

var keywords_closing = map[string]Token{
	"若": TK_若,
	"设": TK_设,
}

var keywords_breaker = map[string]Token{
	"之": TK_之,
}

func IsClosingKeyword(ident string) bool {
	_, ok := keywords_closing[ident]
	return ok
}

func IsBreakerKeyword(ident string) bool {
	_, ok := keywords_breaker[ident]
	return ok
}

func IsSplittingKeyword(ident string) bool {
	return IsKeyword(ident)
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	// wa 支持中文关键字
	if tok, is_keyword := keywords_zh[ident]; is_keyword {
		return tok
	}
	return IDENT
}

// 解析 WaGo 关键字
// WaGo 不支持中文关键字
func LookupWaGo(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
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

var waTokMap = map[Token]watoken.Token{
	ILLEGAL: watoken.ILLEGAL,
	EOF:     watoken.EOF,
	COMMENT: watoken.COMMENT,

	IDENT:  watoken.IDENT,
	INT:    watoken.INT,
	FLOAT:  watoken.FLOAT,
	IMAG:   watoken.IMAG,
	CHAR:   watoken.CHAR,
	STRING: watoken.STRING,

	ADD: watoken.ADD,
	SUB: watoken.SUB,
	MUL: watoken.MUL,
	QUO: watoken.QUO,
	REM: watoken.REM,

	AND:     watoken.AND,
	OR:      watoken.OR,
	XOR:     watoken.XOR,
	SHL:     watoken.SHL,
	SHR:     watoken.SHR,
	AND_NOT: watoken.AND_NOT,

	ADD_ASSIGN: watoken.ADD_ASSIGN,
	SUB_ASSIGN: watoken.SUB_ASSIGN,
	MUL_ASSIGN: watoken.MUL_ASSIGN,
	QUO_ASSIGN: watoken.QUO_ASSIGN,
	REM_ASSIGN: watoken.REM_ASSIGN,

	AND_ASSIGN:     watoken.AND_ASSIGN,
	OR_ASSIGN:      watoken.OR_ASSIGN,
	XOR_ASSIGN:     watoken.XOR_ASSIGN,
	SHL_ASSIGN:     watoken.SHL_ASSIGN,
	SHR_ASSIGN:     watoken.SHR_ASSIGN,
	AND_NOT_ASSIGN: watoken.AND_NOT_ASSIGN,

	LAND: watoken.LAND,
	LOR:  watoken.LOR,
	INC:  watoken.INC,
	DEC:  watoken.DEC,

	EQL:    watoken.EQL,
	LSS:    watoken.LSS,
	GTR:    watoken.GTR,
	ASSIGN: watoken.ASSIGN,
	NOT:    watoken.NOT,

	NEQ:      watoken.NEQ,
	LEQ:      watoken.LEQ,
	GEQ:      watoken.GEQ,
	DEFINE:   watoken.DEFINE,
	ELLIPSIS: watoken.ELLIPSIS,

	LPAREN: watoken.LPAREN,
	LBRACK: watoken.LBRACK,
	LBRACE: watoken.LBRACE,
	COMMA:  watoken.COMMA,
	PERIOD: watoken.PERIOD,

	RPAREN:    watoken.RPAREN,
	RBRACK:    watoken.RBRACK,
	RBRACE:    watoken.RBRACE,
	SEMICOLON: watoken.SEMICOLON,
	COLON:     watoken.COLON,
	ARROW:     watoken.ARROW,

	TK_左方: watoken.LBRACK,
	TK_右方: watoken.RBRACK,
	TK_冒号: watoken.COLON,
	TK_句号: watoken.PERIOD,
	TK_顿号: watoken.COMMA,
	TK_点号: watoken.PERIOD,

	TK_引于: watoken.IMPORT,
	TK_归于: watoken.RETURN,
	TK_若:  watoken.IF,
	TK_否则: watoken.ELSE,
	TK_设:  watoken.VAR,
	TK_从:  watoken.FOR,

	// TK_叹号 :

	BREAK:    watoken.BREAK,
	CASE:     watoken.CASE,
	CONST:    watoken.CONST,
	CONTINUE: watoken.CONTINUE,

	DEFAULT: watoken.DEFAULT,
	DEFER:   watoken.DEFER,
	ELSE:    watoken.ELSE,
	FOR:     watoken.FOR,

	FUNC:   watoken.FUNC,
	IF:     watoken.IF,
	IMPORT: watoken.IMPORT,

	INTERFACE: watoken.INTERFACE,
	MAP:       watoken.MAP,
	PACKAGE:   watoken.PACKAGE,
	RANGE:     watoken.RANGE,
	RETURN:    watoken.RETURN,

	STRUCT: watoken.STRUCT,
	SWITCH: watoken.SWITCH,
	TYPE:   watoken.TYPE,
	VAR:    watoken.VAR,
}

func ToWaTok(tok Token) watoken.Token {
	if t, ok := waTokMap[tok]; ok {
		return t
	}
	return watoken.ILLEGAL
}
