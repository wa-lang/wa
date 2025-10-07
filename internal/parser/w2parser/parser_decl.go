// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseGenDecl(
	keyword token.Token,
	fn func(doc *ast.CommentGroup, keyword token.Token, iota int) ast.Spec,
) *ast.GenDecl {
	if p.trace {
		defer un(trace(p, "GenDecl("+keyword.String()+")"))
	}

	doc := p.leadComment
	pos := p.expect(keyword)

	var lparen, rparen token.Pos
	var list []ast.Spec
	if p.tok == token.COLON {
		// XXX: ... 完毕
		lparen = p.pos
		p.next()
		for iota := 0; p.tok != token.Zh_完毕 && p.tok != token.EOF; iota++ {
			list = append(list, fn(p.leadComment, keyword, iota))
		}
		rparen = p.expect(token.Zh_完毕)
		p.expectSemi()
	} else {
		list = append(list, fn(nil, keyword, 0))
	}

	return &ast.GenDecl{
		Doc:    doc,
		TokPos: pos,
		Tok:    keyword,
		Lparen: lparen,
		Specs:  list,
		Rparen: rparen,
	}
}

func (p *parser) parseDecl(sync map[token.Token]bool) ast.Decl {
	if p.trace {
		defer un(trace(p, "Declaration"))
	}

	var fn func(doc *ast.CommentGroup, keyword token.Token, iota int) ast.Spec
	switch p.tok {
	case token.Zh_常量, token.Zh_全局:
		fn = p.parseValueSpec

	case token.Zh_类型:
		fn = p.parseTypeSpec

	case token.Zh_函数:
		return p.parseFuncDecl(p.tok)

	default:
		pos := p.pos
		p.errorExpected(pos, "declaration:"+p.lit+p.tok.String())
		p.advance(sync)
		return &ast.BadDecl{From: pos, To: p.pos}
	}

	return p.parseGenDecl(p.tok, fn)
}
