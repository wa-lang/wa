// 版权 @2024 凹语言 作者。保留所有权利。

package text

import (
	"fmt"

	"wa-lang.org/wa/internal/3rdparty/wazero/api"
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/leb128"
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/u64"
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm"
)

// 解析global
// (global (mut i32) (i32.const 42))
// (global $pi f32 (f32.const 3.14159))
// (global f32 (f32.const 3.14159))
func (p *moduleParser) parseGlobal(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
	g := &wasm.Global{
		Type: &wasm.GlobalType{},
		Init: &wasm.ConstantExpression{},
	}
	p.currentModuleField = g

	switch tok {
	case tokenID:
		// (global $pi f32 (f32.const 3.14159))
		// (global $pi (f32) (f32.const 3.14159))
		// (global $pi (mut f32) (f32.const 3.14159))
		// ........^

		g.Name = string(tokenBytes)

		return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
			switch tok {
			case tokenKeyword:
				// (global $pi f32 (f32.const 3.14159))
				// ............^...^

				switch string(tokenBytes) {
				default:
					return nil, unexpectedToken(tok, tokenBytes)

				case "i32":
					g.Type.ValType = wasm.ValueTypeI32
					return p.parseGlobal_init, nil
				case "i64":
					g.Type.ValType = wasm.ValueTypeI64
					return p.parseGlobal_init, nil
				case "f32":
					g.Type.ValType = wasm.ValueTypeF32
					return p.parseGlobal_init, nil
				case "f64":
					g.Type.ValType = wasm.ValueTypeF64
					return p.parseGlobal_init, nil
				}

			case tokenLParen:
				// (global $pi (f32) (f32.const 3.14159))
				// (global $pi (mut f32) (f32.const 3.14159))
				// ............^

				return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
					if tok != tokenKeyword {
						return nil, unexpectedToken(tok, tokenBytes)
					}

					switch string(tokenBytes) {
					default:
						return nil, unexpectedToken(tok, tokenBytes)

					case "mut":
						return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
							if tok != tokenKeyword {
								return nil, unexpectedToken(tok, tokenBytes)
							}

							switch string(tokenBytes) {
							default:
								return nil, unexpectedToken(tok, tokenBytes)

							case "i32":
								g.Type.ValType = wasm.ValueTypeI32
							case "i64":
								g.Type.ValType = wasm.ValueTypeI64
							case "f32":
								g.Type.ValType = wasm.ValueTypeF32
							case "f64":
								g.Type.ValType = wasm.ValueTypeF64
							}

							return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
								if tok != tokenRParen {
									return nil, unexpectedToken(tok, tokenBytes)
								}

								return p.parseGlobal_init, nil
							}, nil
						}, nil

					case "i32":
						g.Type.ValType = wasm.ValueTypeI32
					case "i64":
						g.Type.ValType = wasm.ValueTypeI64
					case "f32":
						g.Type.ValType = wasm.ValueTypeF32
					case "f64":
						g.Type.ValType = wasm.ValueTypeF64
					}

					return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
						if tok != tokenRParen {
							return nil, unexpectedToken(tok, tokenBytes)
						}

						return p.parseGlobal_init, nil
					}, nil

				}, nil

			default:
				return nil, unexpectedToken(tok, tokenBytes)
			}
		}, nil

	case tokenKeyword:
		// (global f32 (f32.const 3.14159))
		// ........^
		switch string(tokenBytes) {
		default:
			return nil, unexpectedToken(tok, tokenBytes)

		case "i32":
			g.Type.ValType = wasm.ValueTypeI32
			return p.parseGlobal_init, nil
		case "i64":
			g.Type.ValType = wasm.ValueTypeI64
			return p.parseGlobal_init, nil
		case "f32":
			g.Type.ValType = wasm.ValueTypeF32
			return p.parseGlobal_init, nil
		case "f64":
			g.Type.ValType = wasm.ValueTypeF64
			return p.parseGlobal_init, nil
		}

	case tokenLParen:
		// (global (mut i32) (i32.const 42))
		// (global (i32) (i32.const 42))
		// ........^

		return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
			if tok != tokenKeyword {
				return nil, unexpectedToken(tok, tokenBytes)
			}

			switch string(tokenBytes) {
			default:
				return nil, unexpectedToken(tok, tokenBytes)

			case "mut":
				return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
					if tok != tokenKeyword {
						return nil, unexpectedToken(tok, tokenBytes)
					}

					switch string(tokenBytes) {
					default:
						return nil, unexpectedToken(tok, tokenBytes)

					case "i32":
						g.Type.ValType = wasm.ValueTypeI32
					case "i64":
						g.Type.ValType = wasm.ValueTypeI64
					case "f32":
						g.Type.ValType = wasm.ValueTypeF32
					case "f64":
						g.Type.ValType = wasm.ValueTypeF64
					}

					return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
						if tok != tokenRParen {
							return nil, unexpectedToken(tok, tokenBytes)
						}

						return p.parseGlobal_init, nil
					}, nil
				}, nil

			case "i32":
				g.Type.ValType = wasm.ValueTypeI32
			case "i64":
				g.Type.ValType = wasm.ValueTypeI64
			case "f32":
				g.Type.ValType = wasm.ValueTypeF32
			case "f64":
				g.Type.ValType = wasm.ValueTypeF64
			}

			return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
				if tok != tokenRParen {
					return nil, unexpectedToken(tok, tokenBytes)
				}

				return p.parseGlobal_init, nil
			}, nil

		}, nil

	default:
		return nil, unexpectedToken(tok, tokenBytes)
	}
}

// (global $pi (mut f32) (f32.const 3.14159))
// ......................^
func (p *moduleParser) parseGlobal_init(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
	if tok != tokenLParen {
		return nil, unexpectedToken(tok, tokenBytes)
	}

	return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
		if tok != tokenKeyword {
			return nil, unexpectedToken(tok, tokenBytes)
		}

		g := p.currentModuleField.(*wasm.Global)

		switch string(tokenBytes) {
		default:
			return nil, unexpectedToken(tok, tokenBytes)
		case wasm.OpcodeI32ConstName:
			g.Init.Opcode = wasm.OpcodeI32Const
		case wasm.OpcodeI64ConstName:
			g.Init.Opcode = wasm.OpcodeI64Const
		case wasm.OpcodeF32ConstName:
			g.Init.Opcode = wasm.OpcodeF32Const
		case wasm.OpcodeF64ConstName:
			g.Init.Opcode = wasm.OpcodeF64Const
		}

		return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
			if tok != tokenUN {
				return nil, unexpectedToken(tok, tokenBytes)
			}

			switch g.Init.Opcode {
			case wasm.OpcodeI32Const:
				if i, overflow := decodeUint32(tokenBytes); overflow { // TODO: negative and hex
					return nil, fmt.Errorf("i32 outside range of uint32: %s", tokenBytes)
				} else { // See /RATIONALE.md we can't tell the signed interpretation of a constant, so default to signed.
					g.Init.Data = leb128.EncodeInt32(int32(i))
				}
			case wasm.OpcodeI64Const:
				if i, overflow := decodeUint64(tokenBytes); overflow { // TODO: negative and hex
					return nil, fmt.Errorf("i64 outside range of uint64: %s", tokenBytes)
				} else { // See /RATIONALE.md we can't tell the signed interpretation of a constant, so default to signed.
					g.Init.Data = leb128.EncodeInt64(int64(i))
				}

			case wasm.OpcodeF32Const:
				if i, overflow := decodeUint32(tokenBytes); overflow { // TODO: negative hex nan inf and actual float!
					return nil, fmt.Errorf("f32 outside range of uint32: %s", tokenBytes)
				} else {
					g.Init.Data = u64.LeBytes(api.EncodeF32(float32(i)))
				}
			case wasm.OpcodeF64Const:
				if i, overflow := decodeUint64(tokenBytes); overflow { // TODO: negative hex nan inf and actual float!
					return nil, fmt.Errorf("f64 outside range of uint64: %s", tokenBytes)
				} else {
					g.Init.Data = u64.LeBytes(api.EncodeF64(float64(i)))
				}
			}

			return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
				if tok != tokenRParen {
					return nil, unexpectedToken(tok, tokenBytes)
				}
				return func(tok tokenType, tokenBytes []byte, line, col uint32) (tokenParser, error) {
					if tok != tokenRParen {
						return nil, unexpectedToken(tok, tokenBytes)
					}

					p.currentModuleField = nil
					p.module.GlobalSection = append(p.module.GlobalSection, g)

					p.pos = positionInitial
					return p.parseModule, nil
				}, nil
			}, nil
		}, nil
	}, nil
}
