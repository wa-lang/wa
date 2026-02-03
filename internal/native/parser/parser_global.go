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

func (p *parser) parseGlobal(tok token.Token) *ast.Global {
	g := &ast.Global{Pos: p.pos, Tok: tok}

	g.Doc = p.parseDocComment(&p.prog.Comments, g.Pos)
	if g.Doc != nil {
		p.prog.Objects = p.prog.Objects[:len(p.prog.Objects)-1]
	}

	g.Tok = p.acceptToken(token.GLOBAL_zh)
	g.Name = p.parseIdent()

	p.acceptToken(token.COLON)
	switch p.tok {
	case token.BYTE_zh:
		g.TypeTok = token.BYTE_zh
		g.Type = token.I8
		g.Size = 1
		p.acceptToken(p.tok)
	case token.SHORT_zh:
		g.TypeTok = token.SHORT_zh
		g.Type = token.I16
		g.Size = 2
		p.acceptToken(p.tok)
	case token.LONG_zh:
		g.TypeTok = token.LONG_zh
		g.Type = token.I32
		g.Size = 4
		p.acceptToken(p.tok)
	case token.QUAD_zh:
		g.TypeTok = token.QUAD_zh
		g.Type = token.I64
		g.Size = 8
		p.acceptToken(p.tok)

	case token.ASCII_zh:
		g.TypeTok = token.ASCII_zh
		g.Type = token.I8
		p.acceptToken(p.tok)
	case token.SKIP_zh:
		g.TypeTok = token.SKIP_zh
		g.Type = token.I8
		p.acceptToken(p.tok)
	case token.FILE_zh:
		g.TypeTok = token.FILE_zh
		g.Type = token.I8
		p.acceptToken(p.tok)

	default:
		panic("unreachable: lit:" + p.lit)
	}

	// 全局变量必须显式初始化
	p.acceptToken(token.ASSIGN)

	// TODO: 初始化值和类型要匹配

	g.Init = &ast.InitValue{}

	if p.tok == token.IDENT {
		g.Init.Symbal = p.parseIdent()
		if g.Size == 0 {
			g.Size = 8
		}
	} else {
		g.Init.Lit = p.parseBasicLit()
		if g.Size == 0 {
			if g.Init.Lit.LitKind == token.STRING {
				g.Size = len(g.Init.Lit.ConstV.(string))
			}
		}
	}

	p.consumeSemicolonList()

	return g
}
