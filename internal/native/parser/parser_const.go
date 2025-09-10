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

func (p *parser) parseConst() *ast.Const {
	pConst := &ast.Const{Pos: p.pos}

	pConst.Doc = p.parseDocComment(&p.prog.Comments, pConst.Pos)
	p.acceptTokenAorB(token.CONST, token.CONST_zh)
	pConst.Name = p.parseIdent()
	p.acceptToken(token.ASSIGN)
	pConst.Value = p.parseValue()
	pConst.Comment = p.parseTailComment(pConst.Pos)

	p.consumeSemicolonList()

	return pConst
}
