// 版权 @2023 凹语言 作者。保留所有权利。

package wat

import (
	"fmt"
	"strconv"
)

// Token 结构
type Token struct {
	Pos     Pos       // 位置
	Kind    TokenKind // 类型
	Literal string    // 程序中原始的字符串
}

// 记号类型
type TokenKind int

const (
	// 非法/结尾/注释
	ILLEGAL TokenKind = iota
	EOF
	COMMENT

	// 面值类型
	literal_beg
	INSTRUCTION // 指令, 比如 global.get
	IDENT       // 表示符, 比如 $name
	INT         // 12345
	FLOAT       // 123.45
	CHAR        // 'a'
	STRING      // "abc"
	literal_end

	// 特殊符号
	operator_beg
	LPAREN // (
	RPAREN // )
	operator_end

	// 关键字
	keyword_beg

	I32 // i32
	I64 // i64
	F32 // f32
	F64 // f64

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
	keyword_end
)

var tokens = [...]string{
	// TODO
}

func (tok Token) String() string {
	return fmt.Sprintf("%v:%q", tok.Kind, tok.Literal)
}

func (tokType TokenKind) String() string {
	s := ""
	if 0 <= tokType && tokType < TokenKind(len(tokens)) {
		s = tokens[tokType]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tokType)) + ")"
	}
	return s
}

var keywords = map[string]TokenKind{
	// TODO
}

func watLookup(ident string) TokenKind {
	if tok, isKeyword := keywords[ident]; isKeyword {
		return tok
	}
	return IDENT
}

func (i Token) IntValue() int64 {
	x, err := strconv.ParseInt(i.Literal, 10, 64)
	if err != nil {
		panic(err)
	}
	return x
}
func (i Token) FloatValue() float64 {
	x, err := strconv.ParseFloat(i.Literal, 64)
	if err != nil {
		panic(err)
	}
	return x
}
func (i Token) StrValue() string {
	x, err := strconv.Unquote(i.Literal)
	if err != nil {
		panic(err)
	}
	return x
}
