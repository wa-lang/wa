// 版权 @2023 凹语言 作者。保留所有权利。

package wat

// 将 wat 转化为 token 列表
func Tokenize(filename, input string) (tokens, comments []Token) {
	l := newLexer(filename, input)
	tokens = l.Tokens()
	comments = l.Comments()
	return
}

// 解析模块文件
func ParseFile(filename, input string) (*Module, error) {
	panic("TODO")
}
