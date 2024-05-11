package text

import "fmt"

// 解析 data
// (data (i32.const 2048) "\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
func (p *moduleParser) parseData(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
	fmt.Println("moduleParser.parseData")
	s := string(tokenBytes)
	_ = s
	return nil, unexpectedToken(tok, tokenBytes)
}
