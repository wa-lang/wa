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

	g.Doc = p.parseDocComment(&p.prog.Comments, g.Pos)
	if g.Doc != nil {
		p.prog.Objects = p.prog.Objects[:len(p.prog.Objects)-1]
	}

	g.Tok = p.acceptTokenAorB(token.GLOBAL, token.GLOBAL_zh)
	g.Name = p.parseIdent()

	if p.tok == token.COLON {
		p.acceptToken(token.COLON)
		switch p.tok {
		case token.I32, token.I32_zh:
			g.Type = p.tok
			g.Size = 4
			p.acceptToken(p.tok)
		case token.I64, token.I64_zh:
			g.Type = p.tok
			g.Size = 8
			p.acceptToken(p.tok)
		case token.U32, token.U32_zh:
			g.Type = p.tok
			g.Size = 4
			p.acceptToken(p.tok)
		case token.U64, token.U64_zh:
			g.Type = p.tok
			g.Size = 8
			p.acceptToken(p.tok)
		case token.F32:
			g.Type = p.tok
			g.Size = 4
			p.acceptToken(token.F32)
		case token.F64:
			g.Type = p.tok
			g.Size = 8
			p.acceptToken(token.F64)
		case token.PTR:
			// ptr 大小依赖于平台
			g.Type = token.PTR
			g.Size = 0
		case token.INT:
			// 没有固定类型, 只有内存大小
			g.Type = token.NONE
			g.Size = p.parseIntLit()
		default:
			// 不需要显式指定类型或内存大小的情况
			// global x = INT/FLOAT/STRING
		}
	}

	// 全局变量必须显式初始化
	p.acceptToken(token.ASSIGN)

	if p.tok == token.LBRACE {
		p.parseGlobal_initGroup(g)
	} else {
		initV := &ast.InitValue{Pos: p.pos}
		initV.Doc = p.parseDocComment(&g.Comments, initV.Pos)
		if initV.Doc != nil {
			g.Objects = g.Objects[:len(g.Objects)-1]
		}

		if p.tok == token.IDENT {
			initV.Symbal = p.parseIdent()
		} else {
			initV.Lit = p.parseBasicLit()
		}
		initV.Comment = p.parseTailComment(initV.Pos)
		g.Init = []*ast.InitValue{initV}
	}

	p.consumeSemicolonList()

	return g
}

func (p *parser) parseGlobal_initGroup(g *ast.Global) {
	p.acceptToken(token.LBRACE)
	defer p.acceptToken(token.RBRACE)

	// 结构体初始化
	// 必须显式以整数字面值指定要初始化的偏移地址
Loop:
	for {
		switch p.tok {
		case token.RBRACE:
			break Loop
		case token.COMMENT:
			commentObj := p.parseCommentGroup(false)
			g.Comments = append(g.Comments, commentObj)
			g.Objects = append(g.Objects, commentObj)
		case token.INT:
			initV := &ast.InitValue{Pos: p.pos}
			{
				initV.Doc = p.parseDocComment(&g.Comments, initV.Pos)
				if initV.Doc != nil {
					g.Objects = g.Objects[:len(g.Objects)-1]
				}

				initV.Offset = p.parseIntLit()
				p.acceptToken(token.COLON)
				if p.tok == token.IDENT {
					initV.Symbal = p.parseIdent()
				} else {
					initV.Lit = p.parseBasicLit()
				}
				initV.Comment = p.parseTailComment(initV.Pos)
			}
			g.Init = append(g.Init, initV)
			g.Objects = append(g.Objects, initV)

		default:
			p.errorf(p.pos, "unknown token %v", p.tok)
		}
	}
}
