// 版权 @2024 凹语言 作者。保留所有权利。

// wat记号类型
package token

import (
	"strconv"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm"
)

// 记号类型
type Token int

// 指令码
type Opcode = wasm.Opcode

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
	ANYFUNC // anyfunc
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

	LABEL // label

	keyword_end

	// TODO: 指令转关键字
	// if/else/loop/end/...
	INSTRUCTION // 指令码
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
	ANYFUNC: "anyfunc",
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

	LABEL: "label",

	INSTRUCTION: "INSTRUCTION",
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

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)

	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}

	// $name
	if strings.HasPrefix(ident, "$") {
		return IDENT
	}

	// i32.load
	return INSTRUCTION
}

// 指令
func LookupOpcode(s string) (Opcode, bool) {
	return wasm.LookupOpcode(s)
}
func OpcodeName(op Opcode) string {
	return wasm.OpcodeName(op)
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
	return tok == INSTRUCTION
}

// 标识符
func IsIdent(s string) bool {
	return strings.HasPrefix(s, "$")
}
