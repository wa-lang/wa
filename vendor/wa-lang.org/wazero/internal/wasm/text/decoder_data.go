package text

import (
	"encoding/hex"
	"fmt"

	"wa-lang.org/wazero/internal/leb128"
	"wa-lang.org/wazero/internal/wasm"
)

// 解析 data
// (data (i32.const 2048) "\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
func (p *moduleParser) parseData(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
	if tok != tokenLParen {
		return nil, unexpectedToken(tok, tokenBytes)
	}

	p.currentModuleField = &wasm.DataSegment{
		OffsetExpression: &wasm.ConstantExpression{},
		Init:             []byte{},
	}

	return p.parseData_offset, nil
}

func (p *moduleParser) parseData_offset(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
	if tok != tokenKeyword {
		return nil, unexpectedToken(tok, tokenBytes)
	}
	if string(tokenBytes) != wasm.OpcodeI32ConstName {
		return nil, unexpectedToken(tok, tokenBytes)
	}

	t := p.currentModuleField.(*wasm.DataSegment)
	t.OffsetExpression.Opcode = wasm.OpcodeI32Const

	return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
		if tok != tokenUN {
			return nil, unexpectedToken(tok, tokenBytes)
		}
		if i, overflow := decodeUint32(tokenBytes); overflow { // TODO: negative and hex
			return nil, fmt.Errorf("i32 outside range of uint32: %s", tokenBytes)
		} else { // See /RATIONALE.md we can't tell the signed interpretation of a constant, so default to signed.
			t.OffsetExpression.Data = leb128.EncodeInt32(int32(i))
		}
		return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
			if tok != tokenRParen {
				return nil, unexpectedToken(tok, tokenBytes)
			}
			return p.parseData_init, nil
		}, nil
	}, nil

}
func (p *moduleParser) parseData_init(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
	if tok != tokenString {
		return nil, unexpectedToken(tok, tokenBytes)
	}

	src := make([]byte, 0, len(tokenBytes))
	for _, c := range tokenBytes {
		if (c >= '0' && c <= '9') || (c >= 'a' && c < 'z') || (c >= 'A' && c <= 'Z') {
			src = append(src, c)
		}
	}
	dst := make([]byte, hex.DecodedLen(len(src)))

	n, err := hex.Decode(dst, src)
	if err != nil {
		return nil, err
	}

	t := p.currentModuleField.(*wasm.DataSegment)
	t.Init = dst[:n]

	return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
		if tok != tokenRParen {
			return nil, unexpectedToken(tok, tokenBytes)
		}

		p.currentModuleField = nil
		p.module.DataSection = append(p.module.DataSection, t)

		p.pos = positionInitial
		return p.parseModule, nil
	}, nil
}
