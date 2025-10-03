// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseStmtList_zh() (list []ast.Stmt) {
	if p.trace {
		defer un(trace(p, "StatementList"))
	}

	for p.tok != token.Zh_岔道 && p.tok != token.Zh_主道 && p.tok != token.Zh_函终 && p.tok != token.Zh_算终 && p.tok != token.EOF {
		list = append(list, p.parseStmt_zh())
	}

	return
}

func (p *parser) parseStmt_zh() (s ast.Stmt) {
	if p.trace {
		defer un(trace(p, "Statement"))
	}

	switch p.tok {
	case token.CONST, token.TYPE, token.VAR:
		// TODO(chai2010): var declaration not allowed in func body
		s = &ast.DeclStmt{Decl: p.parseDecl_zh(stmtStart)}
	case
		// tokens that may start an expression
		token.IDENT, token.INT, token.FLOAT, token.IMAG, token.CHAR, token.STRING, token.FUNC, token.LPAREN, // operands
		token.LBRACK, token.STRUCT, token.MAP, token.INTERFACE, // composite types
		token.ADD, token.SUB, token.MUL, token.AND, token.XOR, token.NOT: // unary operators
		s, _ = p.parseSimpleStmt(labelOk)
		// because of the required look-ahead, labeled statements are
		// parsed by parseSimpleStmt - don't expect a semicolon after
		// them
		if _, isLabeledStmt := s.(*ast.LabeledStmt); !isLabeledStmt {
			p.expectSemi()
		}
	case token.DEFER, token.Zh_押后:
		s = p.parseDeferStmt(p.tok)
	case token.RETURN, token.Zh_返回:
		s = p.parseReturnStmt(p.tok)
	case token.BREAK, token.CONTINUE:
		s = p.parseBranchStmt(p.tok)
	case token.LBRACE:
		s = p.parseBlockStmt()
		p.expectSemi()
	case token.IF, token.Zh_若始:
		s = p.parseIfStmt(p.tok)
	case token.SWITCH, token.Zh_岔始:
		s = p.parseSwitchStmt(p.tok)
	case token.FOR, token.Zh_当始:
		s = p.parseForStmt(p.tok)
	case token.SEMICOLON:
		// Is it ever possible to have an implicit semicolon
		// producing an empty statement in a valid program?
		// (handle correctly anyway)
		s = &ast.EmptyStmt{Semicolon: p.pos, Implicit: p.lit == "\n"}
		p.next()
	case token.RBRACE:
		// a semicolon may be omitted before a closing "}"
		s = &ast.EmptyStmt{Semicolon: p.pos, Implicit: true}
	default:
		// no statement found
		pos := p.pos
		p.errorExpected(pos, "statement")
		p.advance(stmtStart)
		s = &ast.BadStmt{From: pos, To: p.pos}
	}

	return
}
