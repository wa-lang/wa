// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

// 顶层的声明入口
func (p *parser) parseDecl(sync map[token.Token]bool) ast.Decl {
	if p.trace {
		defer un(trace(p, "Declaration"))
	}

	switch p.tok {
	case token.Zh_常量:
		return p.parseGenDecl_const(p.tok)
	case token.Zh_全局:
		return p.parseGenDecl_global(p.tok)
	case token.Zh_结构:
		return p.parseStructDecl(p.tok)
	case token.Zh_接口:
		return p.parseInterfaceDecl(p.tok)
	case token.Zh_函数:
		return p.parseFuncDecl(p.tok)

	default:
		pos := p.pos
		p.errorExpected(pos, "declaration:"+p.lit+p.tok.String())
		p.advance(sync)
		return &ast.BadDecl{From: pos, To: p.pos}
	}
}
