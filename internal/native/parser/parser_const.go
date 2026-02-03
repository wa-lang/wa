// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

// 和 #define 类似, 只是简单的替换, 数值没有精确的类型
// const $A = 0x100000
// const $B = 12.5
// const $C = "abc"

func (p *parser) parseConst(tok token.Token) *ast.Const {
	pConst := &ast.Const{Pos: p.pos, Tok: tok}

	pConst.Doc = p.parseDocComment(&p.prog.Comments, pConst.Pos)
	if pConst.Doc != nil {
		p.prog.Objects = p.prog.Objects[:len(p.prog.Objects)-1]
	}

	pConst.Tok = p.acceptToken(token.CONST_zh)
	pConst.Name = p.parseIdent()
	p.acceptToken(token.ASSIGN)
	pConst.Value = p.parseBasicLit()
	if pConst.Value.LitKind == token.STRING {
		p.errorf(pConst.Value.Pos, "const(%s) donot support string type", pConst.Value.LitString)
	}
	pConst.Comment = p.parseTailComment(pConst.Pos)

	p.consumeSemicolonList()

	return pConst
}
