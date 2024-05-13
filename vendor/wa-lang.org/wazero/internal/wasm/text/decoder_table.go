// 版权 @2024 凹语言 作者。保留所有权利。

package text

import (
	"fmt"

	"wa-lang.org/wazero/internal/wasm"
)

// 解析 table
// (table 14 funcref)
func (p *moduleParser) parseTable(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
	if tok != tokenUN {
		return nil, unexpectedToken(tok, tokenBytes)
	}

	t := &wasm.Table{}
	p.currentModuleField = t

	if i, overflow := decodeUint32(tokenBytes); overflow { // TODO: negative and hex
		return nil, fmt.Errorf("i32 outside range of uint32: %s", tokenBytes)
	} else { // See /RATIONALE.md we can't tell the signed interpretation of a constant, so default to signed.
		t.Min = i
	}

	return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
		if tok != tokenKeyword {
			return nil, unexpectedToken(tok, tokenBytes)
		}
		if string(tokenBytes) != "funcref" {
			return nil, fmt.Errorf("donot support %s", tokenBytes)
		}

		t.Type = wasm.RefTypeFuncref

		return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
			if tok != tokenRParen {
				return nil, unexpectedToken(tok, tokenBytes)
			}

			p.currentModuleField = nil
			p.module.TableSection = append(p.module.TableSection, t)

			p.pos = positionInitial
			return p.parseModule, nil
		}, nil

	}, nil
}
