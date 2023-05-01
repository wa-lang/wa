// 版权 @2023 凹语言 作者。保留所有权利。

package wat

// Token 结构
type Token struct {
	Pos   TokenPos  // 位置
	Kind  TokenKind // 类型
	Value string    // 面值
}

// 位置信息
type TokenPos struct {
	Offset int // offset, 从 0 开始
	Line   int // 行号, 从 1 开始
	Column int // 列号, 从 1 开始的字节数
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

// 将 wat 转化为 token 列表
func Tokenize(input []byte) []Token {
	return nil
}
