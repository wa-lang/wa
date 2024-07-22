package text

import "fmt"

// 解析elem
// (elem (i32.const 1) $$u8.$$block.$$onFree)
func (p *moduleParser) parseElem(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
	fmt.Println("moduleParser.parseElem")
	s := string(tokenBytes)
	_ = s
	return nil, unexpectedToken(tok, tokenBytes)
}
