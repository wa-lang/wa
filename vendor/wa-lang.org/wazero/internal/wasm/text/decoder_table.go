// 版权 @2024 凹语言 作者。保留所有权利。

package text

import "fmt"

// 解析 table
// (table 14 funcref)
func (p *moduleParser) parseTable(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
	fmt.Println("moduleParser.parseTable")
	s := string(tokenBytes)
	_ = s
	return nil, unexpectedToken(tok, tokenBytes)
}
