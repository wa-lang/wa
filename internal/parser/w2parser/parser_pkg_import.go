// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"strconv"
	"strings"
	"unicode"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseGenDecl_import(keyword token.Token) *ast.GenDecl {
	if p.trace {
		defer un(trace(p, "GenDecl("+keyword.String()+")"))
	}

	doc := p.leadComment
	pos := p.expect(keyword)

	var lparen, rparen token.Pos
	var spec = p.parseImportSpec(nil, keyword, 0)

	return &ast.GenDecl{
		Doc:    doc,
		TokPos: pos,
		Tok:    keyword,
		Lparen: lparen,
		Specs:  []ast.Spec{spec},
		Rparen: rparen,
	}
}

func (p *parser) parseImportSpec(doc *ast.CommentGroup, _ token.Token, _ int) ast.Spec {
	if p.trace {
		defer un(trace(p, "ImportSpec"))
	}

	pos := p.pos
	var path string
	if p.tok == token.STRING {
		path = p.lit
		if !isValidImport(path) {
			p.error(pos, "invalid import path: "+path)
		}
		p.next()
	} else {
		p.expect(token.STRING) // use expect() error handling
	}

	var ident *ast.Ident
	var arrowPos token.Pos

	// parse => asname
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

	p.expectSemi() // call before accessing p.linecomment

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

func isValidImport(lit string) bool {
	const illegalChars = `!"#$%&'()*,:;<=>?[\]^{|}` + "`\uFFFD"
	s, _ := strconv.Unquote(lit) // wa-lang.org/wa/internal/scanner returns a legal string literal
	for _, r := range s {
		if !unicode.IsGraphic(r) || unicode.IsSpace(r) || strings.ContainsRune(illegalChars, r) {
			return false
		}
	}
	return s != ""
}
