// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

// const $EXIT_DEVICE = 0x100000

func (p *parser) parseConst() *ast.Const {
	pConst := &ast.Const{Pos: p.pos}

	p.acceptToken(token.CONST)

	if p.tok == token.IDENT {
		pConst.Name = p.parseIdent()
	}

	p.acceptToken(token.ASSIGN)

	// 无类型的面值
	switch p.tok {
	case token.INT:
	case token.FLOAT:
	case token.STRING:
	case token.I32, token.I64, token.F32, token.F64:
	}

	// 带类型转义

	p.next()
	// TODO

	return pConst
}
