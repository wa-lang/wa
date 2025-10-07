// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"fmt"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseDeferStmt(keyword token.Token) ast.Stmt {
	if p.trace {
		defer un(trace(p, "DeferStmt"))
	}

	pos := p.expect(keyword)
	call := p.parseDeferStmt_CallExpr(keyword.String())
	p.expectSemi()
	if call == nil {
		return &ast.BadStmt{From: pos, To: pos + token.Pos(len(keyword.String()))} // len("defer")
	}

	return &ast.DeferStmt{
		TokPos: pos,
		Tok:    keyword,
		Call:   call,
	}
}

func (p *parser) parseDeferStmt_CallExpr(callType string) *ast.CallExpr {
	x := p.parseRhsOrType() // could be a conversion: (some type)(x)
	if call, isCall := x.(*ast.CallExpr); isCall {
		return call
	}
	if _, isBad := x.(*ast.BadExpr); !isBad {
		// only report error if it's a new one
		p.error(p.safePos(x.End()), fmt.Sprintf("function must be invoked in %s statement", callType))
	}
	return nil
}
