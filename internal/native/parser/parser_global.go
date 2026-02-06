// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"math"
	"os"

	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

// global $f64: f64 = 12.34567
// global $name = "wa native assembly language"

// global $f32: 20 = f32(12.5)

func (p *parser) parseGlobal(tok token.Token) *ast.Global {
	g := &ast.Global{
		Pos: p.pos,
		Tok: tok,

		Init: &ast.InitValue{},

		Section: ".text",
	}

	g.Doc = p.parseDocComment(&p.prog.Comments, g.Pos)
	if g.Doc != nil {
		p.prog.Objects = p.prog.Objects[:len(p.prog.Objects)-1]
	}

	g.Tok = p.acceptToken(token.GLOBAL_zh)
	g.Name = p.parseIdent()

	p.acceptToken(token.COLON)

	// 解析类型Token
	g.TypeTok = p.acceptToken(
		token.BYTE_zh,
		token.SHORT_zh,
		token.LONG_zh,
		token.QUAD_zh,
		token.ASCII_zh,
		token.SKIP_zh,
		token.FILE_zh,
	)

	// 全局变量必须显式初始化
	p.acceptToken(token.ASSIGN)

	// 解析初始化值
	g.Init.Pos = p.pos
	if p.tok == token.IDENT {
		g.Init.Symbal = p.parseIdent()
		g.Init.Comment = p.parseTailComment(g.Init.Pos)
		p.consumeSemicolonList()
	} else {
		g.Init.Lit = p.parseBasicLit()
		g.Init.Comment = p.parseTailComment(g.Init.Pos)
		p.consumeSemicolonList()
	}

	// 验证初始化值的合法性, 填充类型和Size
	switch g.TypeTok {
	case token.BYTE_zh:
		g.Type = token.I8
		g.Size = 1
		v := int(g.Init.Lit.ConstV.(int64))
		assert(v >= math.MinInt8 && v < math.MaxUint8)
	case token.SHORT_zh:
		g.Type = token.I16
		g.Size = 2
		v := int(g.Init.Lit.ConstV.(int64))
		assert(v >= math.MinInt16 && v < math.MaxUint16)
	case token.LONG_zh:
		switch v := g.Init.Lit.ConstV.(type) {
		case int64:
			g.Type = token.I32
			g.Size = 4
			assert(v >= math.MinInt32 && v < math.MaxUint32)
		case float64:
			g.Type = token.F32
			g.Size = 4
		default:
			panic("unreachable")
		}
	case token.QUAD_zh:
		switch g.Init.Lit.ConstV.(type) {
		case int64:
			g.Type = token.I64
			g.Size = 8
		case float64:
			g.Type = token.F64
			g.Size = 8
		default:
			panic("unreachable")
		}
	case token.ASCII_zh:
		g.Type = token.I64 // 地址类型
		g.Size = len(g.Init.Lit.ConstV.(string))
	case token.SKIP_zh:
		g.Type = token.I64 // 地址类型
		g.Size = int(g.Init.Lit.ConstV.(int64))
	case token.FILE_zh:
		g.Type = token.I64 // 地址类型
		filename := g.Init.Lit.ConstV.(string)
		if fi, err := os.Lstat(filename); err != nil {
			p.errorf(p.pos, "%v %v failed: %v", token.GAS_INCBIN, filename, err)
		} else {
			const maxSize = 2 << 20
			if fi.Size() > maxSize {
				p.errorf(p.pos, "%v %v file size large than 2MB: %v", token.GAS_INCBIN, filename, err)
			}
			g.Size = int(fi.Size())
		}
	default:
		panic("unreachable")
	}

	return g
}
