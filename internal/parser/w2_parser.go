// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseFile_zh() *ast.File {
	if p.trace {
		defer un(trace(p, "文件"))
	}

	if p.errors.Len() != 0 {
		return nil
	}

	doc := p.leadComment
	pos := p.pos
	ident := &ast.Ident{}

	p.openScope()
	p.pkgScope = p.topScope
	var decls []ast.Decl
	if p.mode&PackageClauseOnly == 0 {
		for p.tok == token.Zh_引入 {
			decls = append(decls, p.parseGenDecl_zh(p.tok, p.parseImportSpec_zh))
		}

		if p.mode&ImportsOnly == 0 {
			for p.tok != token.EOF {
				decls = append(decls, p.parseDecl_zh(declStart))
			}
		}
	}
	p.closeScope()
	assert(p.topScope == nil, "unbalanced scopes")
	assert(p.labelScope == nil, "unbalanced label scopes")

	// resolve global identifiers within the same file
	i := 0
	for _, ident := range p.unresolved {
		// i <= index for current ident
		assert(ident.Obj == unresolved, "object already resolved")
		ident.Obj = p.pkgScope.Lookup(ident.Name) // also removes unresolved sentinel
		if ident.Obj == nil {
			p.unresolved[i] = ident
			i++
		}
	}

	return &ast.File{
		Doc:        doc,
		Package:    pos,
		Name:       ident,
		Decls:      decls,
		Scope:      p.pkgScope,
		Imports:    p.imports,
		Unresolved: p.unresolved[0:i],
		Comments:   p.comments,
	}
}

func (p *parser) expectSemi_zh() {
	// semicolon is optional before a closing ')' or '}'
	if p.tok != token.RPAREN && p.tok != token.RBRACE {
		switch p.tok {
		case token.COMMA:
			// permit a ',' instead of a ';' but complain
			p.errorExpected(p.pos, "';'")
			fallthrough
		case token.SEMICOLON:
			p.next()
		default:
			p.errorExpected(p.pos, "';'")
			p.advance(stmtStart)
		}
	}
}

func (p *parser) errorExpected_zh(pos token.Pos, msg string) {
	msg = "期望 " + msg
	if pos == p.pos {
		// the error happened at the current position;
		// make the error message more specific
		switch {
		case p.tok == token.SEMICOLON && p.lit == "\n":
			msg += ", 实获 换行符"
		case p.tok.IsLiteral():
			// print 123 rather than 'INT', etc.
			msg += ", 实获 " + p.lit
		default:
			msg += ", 实获 '" + p.tok.String() + "'"
		}
	}
	p.error(pos, msg)
}
