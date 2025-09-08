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

	p.acceptToken(token.CONST)
	pConst.Name = p.parseIdent()
	p.acceptToken(token.ASSIGN)

	// 无类型的面值
	switch p.tok {
	case token.INT:
		pConst.Value.IntValue = p.parseInt64Lit()
	case token.FLOAT:
		pConst.Value.FloatValue = p.parseFloat64Lit()
	case token.STRING:
		pConst.Value.StrValue = p.parseStringLit()
	}

	p.consumeTokenList(token.SEMICOLON)

	return pConst
}
