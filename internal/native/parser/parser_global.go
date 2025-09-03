// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

// global $f64: f64 = 12.34567
// global $name = "wa native assembly language"

// global $f32: 20 = f32(12.5)

// global $info: 1024 = {
//     5: "abc",    # 从第5字节开始 `abc\0`
//     9: i32(123), # 从第9字节开始
// }

func (p *parser) parseGlobal() *ast.Global {
	g := &ast.Global{Pos: p.pos}

	p.acceptToken(token.GLOBAL)
	g.Name = p.parseIdent()

	p.acceptToken(token.COLON)
	switch p.tok {
	case token.I32:
		g.Size = 4
		p.acceptToken(token.I32)
	case token.I64:
		g.Size = 8
		p.acceptToken(token.I64)
	case token.F32:
		g.Size = 4
		p.acceptToken(token.F32)
	case token.F64:
		g.Size = 8
		p.acceptToken(token.F64)
	case token.INT:
		g.Size = p.parseIntLit()
	default:
		p.errorf(p.pos, "expect %v, got %v", "I32/I64/F32/F64/INT", p.tok)
	}

	p.acceptToken(token.ASSIGN)

	if p.tok == token.LBRACE {
		g.Init = p.parseGlobal_initGroup()
	} else {
		g.Init = []ast.InitValue{p.parseGlobal_initValue(0)}
	}

	return g
}

func (p *parser) parseGlobal_initGroup() []ast.InitValue {
	p.acceptToken(token.LBRACE)
	defer p.acceptToken(token.RBRACE)

	var initGroup []ast.InitValue
	for p.tok == token.INT {
		initGroup = append(initGroup, p.parseGlobal_initGroup_elem())
	}

	return initGroup
}

func (p *parser) parseGlobal_initGroup_elem() ast.InitValue {
	offset := p.parseIntLit()
	p.acceptToken(token.COLON)
	return p.parseGlobal_initValue(offset)
}

func (p *parser) parseGlobal_initValue(offset int) ast.InitValue {
	switch p.tok {
	case token.INT:
		return ast.InitValue{
			Offset: offset,
			Type:   token.I32,
			LitValue: &ast.Value{
				LitKind:  token.INT,
				IntValue: int64(p.parseInt32Lit()),
			},
		}

	case token.I32:
		p.acceptToken(token.I32)
		p.acceptToken(token.LPAREN)
		defer p.acceptToken(token.RPAREN)

		return ast.InitValue{
			Offset: offset,
			Type:   token.I32,
			LitValue: &ast.Value{
				LitKind:  token.INT,
				IntValue: int64(p.parseInt32Lit()),
			},
		}

	case token.I64:
		p.acceptToken(token.I32)
		p.acceptToken(token.LPAREN)
		defer p.acceptToken(token.RPAREN)

		return ast.InitValue{
			Offset: offset,
			Type:   token.I32,
			LitValue: &ast.Value{
				LitKind:  token.INT,
				IntValue: int64(p.parseInt64Lit()),
			},
		}

	case token.FLOAT:
		return ast.InitValue{
			Offset: offset,
			Type:   token.F32,
			LitValue: &ast.Value{
				LitKind:    token.FLOAT,
				FloatValue: float64(p.parseFloat32Lit()),
			},
		}

	case token.F32:
		p.acceptToken(token.F32)
		p.acceptToken(token.LPAREN)
		defer p.acceptToken(token.RPAREN)

		return ast.InitValue{
			Offset: offset,
			Type:   token.F32,
			LitValue: &ast.Value{
				LitKind:    token.FLOAT,
				FloatValue: float64(p.parseFloat32Lit()),
			},
		}

	case token.F64:
		p.acceptToken(token.F64)
		p.acceptToken(token.LPAREN)
		defer p.acceptToken(token.RPAREN)

		return ast.InitValue{
			Offset: offset,
			Type:   token.F32,
			LitValue: &ast.Value{
				LitKind:    token.FLOAT,
				FloatValue: float64(p.parseFloat64Lit()),
			},
		}

	case token.STRING:
		return ast.InitValue{
			Offset: offset,
			Type:   token.STRING,
			LitValue: &ast.Value{
				LitKind:  token.STRING,
				StrValue: p.parseStringLit(),
			},
		}

	default:
		p.errorf(p.pos, "expect global init value, got %v", p.tok)
	}

	panic("unreachable")
}
