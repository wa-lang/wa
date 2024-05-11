// 版权 @2024 凹语言 作者。保留所有权利。

package text

import (
	"fmt"

	"wa-lang.org/wazero/internal/wasm"
)

// 解析global
// (global $__stack_ptr (mut i32) (i32.const 1024))     ;; index=0
// (global $$wa.runtime.closure_data (mut i32) (i32.const 0))
func (p *moduleParser) parseGlobal(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
	fmt.Println("moduleParser.parseGlobal")

	switch tok {
	case tokenID: // Ex. $__stack_ptr
		p.currentModuleField = &wasm.Global{
			Name: string(tokenBytes),
			Type: &wasm.GlobalType{
				ValType: 1,
				Mutable: true,
			},
			Init: &wasm.ConstantExpression{
				Opcode: wasm.OpcodeI32Const,
			},
		}
		return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
			switch tok {
			case tokenLParen:
				return p.parseGlobalDesc, nil
			default:
				return nil, unexpectedToken(tok, tokenBytes)
			}
		}, nil

	default:
		return nil, unexpectedToken(tok, tokenBytes)
	}
}

func (p *moduleParser) parseGlobalDesc(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
	if tok != tokenKeyword {
		return nil, unexpectedToken(tok, tokenBytes)
	}

	g := p.currentModuleField.(*wasm.Global)
	switch string(tokenBytes) {
	case "mut":
		g.Type.Mutable = true
		// TODO

	case "i32.const":
		// TODO
	case "i64.const":
		// TODO
	case "f32.const":
		// TODO
	case "f64.const":
		// TODO
	}

	return nil, unexpectedToken(tok, tokenBytes)
}

// 解析 table
// (table 14 funcref)
func (p *moduleParser) parseTable(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
	fmt.Println("moduleParser.parseTable")
	s := string(tokenBytes)
	_ = s
	return nil, unexpectedToken(tok, tokenBytes)
}

// 解析 data
// (data (i32.const 2048) "\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
func (p *moduleParser) parseData(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
	fmt.Println("moduleParser.parseData")
	s := string(tokenBytes)
	_ = s
	return nil, unexpectedToken(tok, tokenBytes)
}

// 解析elem
// (elem (i32.const 1) $$u8.$$block.$$onFree)
func (p *moduleParser) parseElem(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
	fmt.Println("moduleParser.parseElem")
	s := string(tokenBytes)
	_ = s
	return nil, unexpectedToken(tok, tokenBytes)
}
