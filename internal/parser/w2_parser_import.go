// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseImportSpec_zh(doc *ast.CommentGroup, _ token.Token, _ int) ast.Spec {
	if p.trace {
		defer un(trace(p, "ImportSpec"))
	}

	pos := p.pos
	var path string
	if p.tok == token.STRING {
		path = p.lit
		if !isValidImport(path) {
			p.error(pos, "无效的引入路径: "+path)
		}
		p.next()
	} else {
		p.expect(token.STRING) // use expect() error handling
	}

	// parse => asname
	var ident *ast.Ident
	var arrowPos token.Pos
	if p.tok == token.ARROW {
		arrowPos = p.pos
		p.next() // skip =>

		switch p.tok {
		case token.PERIOD:
			ident = &ast.Ident{NamePos: p.pos, Name: "."}
			p.next()
		case token.IDENT:
			ident = p.parseIdent()
		default:
			p.expect(token.IDENT)
		}
	}

	p.expectSemi_zh() // call before accessing p.linecomment

	// collect imports
	spec := &ast.ImportSpec{
		Doc:      doc,
		Name:     ident,
		Path:     &ast.BasicLit{ValuePos: pos, Kind: token.STRING, Value: path},
		ArrowPos: arrowPos,
		Comment:  p.lineComment,
	}
	p.imports = append(p.imports, spec)

	return spec
}
