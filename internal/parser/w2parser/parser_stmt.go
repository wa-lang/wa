// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

// Parsing modes for parseSimpleStmt.
const (
	basic = iota
	labelOk
	rangeOk
)

func (p *parser) parseStmtList() (list []ast.Stmt) {
	if p.trace {
		defer un(trace(p, "StatementList"))
	}

	for p.tok != token.Zh_有辙 && p.tok != token.Zh_没辙 && p.tok != token.Zh_或者 && p.tok != token.Zh_否则 && p.tok != token.Zh_完毕 && p.tok != token.EOF {
		list = append(list, p.parseStmt())
	}

	return
}

func (p *parser) parseStmt() (s ast.Stmt) {
	if p.trace {
		defer un(trace(p, "Statement"))
	}

	switch p.tok {
	case token.Zh_常量:
		s = &ast.DeclStmt{Decl: p.parseGenDecl_const(p.tok)}
	case
		// tokens that may start an expression
		token.IDENT, token.INT, token.FLOAT, token.IMAG, token.CHAR, token.STRING, token.Zh_函数, token.LPAREN, // operands
		token.LBRACK, token.Zh_结构, token.Zh_字典, token.Zh_接口, // composite types
		token.ADD, token.SUB, token.MUL, token.AND, token.XOR, token.NOT: // unary operators
		s, _ = p.parseSimpleStmt(0, labelOk)
		// because of the required look-ahead, labeled statements are
		// parsed by parseSimpleStmt - don't expect a semicolon after
		// them
		if _, isLabeledStmt := s.(*ast.LabeledStmt); !isLabeledStmt {
			p.expectSemi()
		}
	case token.Zh_押后:
		s = p.parseDeferStmt(p.tok)
	case token.Zh_返回:
		s = p.parseReturnStmt(p.tok)
	case token.Zh_跳出, token.Zh_继续:
		s = p.parseBranchStmt(p.tok)
	case token.Zh_区块:
		p.expect(token.Zh_区块)
		s = p.parseBlockStmt(token.Zh_完毕)
		p.expectSemi()
	case token.Zh_如果:
		s = p.parseIfStmt(p.tok)
	case token.Zh_找辙:
		s = p.parseSwitchStmt(p.tok)
	case token.Zh_循环:
		s = p.parseForStmt(p.tok)
	case token.SEMICOLON:
		// Is it ever possible to have an implicit semicolon
		// producing an empty statement in a valid program?
		// (handle correctly anyway)
		s = &ast.EmptyStmt{Semicolon: p.pos, Implicit: p.lit == "\n"}
		p.next()
	case token.Zh_完毕:
		// a semicolon may be omitted before a closing "}"
		s = &ast.EmptyStmt{Semicolon: p.pos, Implicit: true}
	default:
		// no statement found
		pos := p.pos
		p.errorExpected(pos, "statement"+p.tok.String())
		p.advance(stmtStart)
		s = &ast.BadStmt{From: pos, To: p.pos}
	}

	return
}

// parseSimpleStmt returns true as 2nd result if it parsed the assignment
// of a range clause (with mode == rangeOk). The returned statement is an
// assignment with a right-hand side that is a single unary expression of
// the form "range x". No guarantees are given for the left-hand side.
func (p *parser) parseSimpleStmt(keyword token.Token, mode int) (ast.Stmt, bool) {
	if p.trace {
		defer un(trace(p, "SimpleStmt"))
	}

	var x = p.parseLhsList()
	var colonPos token.Pos

	switch p.tok {
	case token.COLON:
		// ‘:’, 对应 if/switch/for 的区块开始, 不能吃掉
		if keyword != token.ILLEGAL {
			// break
		}

		colonPos = p.pos
		p.next()

		// Wa 只有 for 和 switch 有 Label
		if p.tok == token.Zh_循环 || p.tok == token.Zh_找辙 {
			break // 继续后面的 Label 解析
		}

		// x: int
		// x: int = 123
		// END: int = 123 // 这里的 END 是 变量
		// a, b: int
		// a, b: int = 123, 456

		// 解析 idents 列表
		var idents = make([]*ast.Ident, 0, len(x))
		for _, xi := range x {
			if ident, ok := xi.(*ast.Ident); ok {
				idents = append(idents, ident)
			} else {
				p.errorExpected(xi.Pos(), "identifier on left side of :type")
				break
			}
		}

		typ := p.tryType()
		if typ == nil {
			p.error(colonPos, "missing variable type")
		}

		var values []ast.Expr
		// always permit optional initialization for more tolerant parsing
		if p.tok == token.ASSIGN {
			p.next()
			values = p.parseRhsList()
		}

		// Go spec: The scope of a constant or variable identifier declared inside
		// a function begins at the end of the ConstSpec or VarSpec and ends at
		// the end of the innermost containing block.
		// (Global identifiers are resolved in a separate phase after parsing.)
		spec := &ast.ValueSpec{
			Names:    idents,
			ColonPos: colonPos,
			Type:     typ,
			Values:   values,
			Comment:  p.lineComment,
		}

		p.declare(spec, nil, p.topScope, ast.Var, idents...)

		declStmt := &ast.DeclStmt{
			Decl: &ast.GenDecl{
				TokPos: x[0].Pos(),
				Tok:    token.VAR,
				Specs:  []ast.Spec{spec},
			},
		}

		return declStmt, false

	case
		token.DEFINE, token.ASSIGN, token.ADD_ASSIGN,
		token.SUB_ASSIGN, token.MUL_ASSIGN, token.QUO_ASSIGN,
		token.REM_ASSIGN, token.AND_ASSIGN, token.OR_ASSIGN,
		token.XOR_ASSIGN, token.SHL_ASSIGN, token.SHR_ASSIGN, token.AND_NOT_ASSIGN:
		// assignment statement, possibly part of a range clause
		pos, tok := p.pos, p.tok
		p.next()
		var y []ast.Expr
		isRange := false
		if mode == rangeOk && p.tok == token.Zh_迭代 && (tok == token.DEFINE || tok == token.ASSIGN) {
			pos := p.pos
			p.next()
			y = []ast.Expr{&ast.UnaryExpr{OpPos: pos, Op: token.RANGE, X: p.parseRhs()}}
			isRange = true
		} else {
			y = p.parseRhsList()
		}
		as := &ast.AssignStmt{Lhs: x, TokPos: pos, Tok: tok, Rhs: y}
		if tok == token.DEFINE {
			p.shortVarDecl(as, x)
		}
		return as, isRange
	}

	if len(x) > 1 {
		p.errorExpected(x[0].Pos(), "1 expression")
		// continue with first expression
	}

	if colonPos != token.NoPos {
		// labeled statement
		colon := colonPos

		// 只有 for 和 switch 有 Label
		if p.tok != token.Zh_循环 && p.tok != token.Zh_找辙 {
			p.errorExpected(p.pos, "for or switch")
		}

		if label, isIdent := x[0].(*ast.Ident); mode == labelOk && isIdent {
			// Go spec: The scope of a label is the body of the function
			// in which it is declared and excludes the body of any nested
			// function.
			stmt := &ast.LabeledStmt{Label: label, Colon: colon, Stmt: p.parseStmt()}
			p.declare(stmt, nil, p.labelScope, ast.Lbl, label)
			return stmt, false
		}
		// The label declaration typically starts at x[0].Pos(), but the label
		// declaration may be erroneous due to a token after that position (and
		// before the ':'). If SpuriousErrors is not set, the (only) error
		// reported for the line is the illegal label error instead of the token
		// before the ':' that caused the problem. Thus, use the (latest) colon
		// position for error reporting.
		p.error(colon, "illegal label declaration")
		return &ast.BadStmt{From: x[0].Pos(), To: colon + 1}, false
	} else {
		switch p.tok {
		case token.INC, token.DEC:
			// increment or decrement
			s := &ast.IncDecStmt{X: x[0], TokPos: p.pos, Tok: p.tok}
			p.next()
			return s, false
		}
	}

	// expression
	return &ast.ExprStmt{X: x[0]}, false
}
