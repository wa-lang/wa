// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseBlockStmt(endToken token.Token, orOtherEndTokens ...token.Token) *ast.BlockStmt {
	if p.trace {
		defer un(trace(p, "BlockStmt"))
	}

	lbrace := p.expect(token.COLON)
	p.openScope()
	list := p.parseStmtList()
	p.closeScope()

	// if 结束块有多种
	if p.tok != endToken {
		for _, tok := range orOtherEndTokens {
			if p.tok == tok {
				endToken = tok
				break
			}
		}
	}

	rbrace := p.pos
	if p.tok != endToken {
		p.errorExpected(p.pos, "'"+p.tok.String()+"'")
	}

	// 只有 完毕 可以吃掉
	// 或者/否则 不能吃掉
	if endToken == token.Zh_完毕 {
		p.next()
	}

	return &ast.BlockStmt{Lbrace: lbrace, List: list, Rbrace: rbrace}
}
