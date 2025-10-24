// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *parser) parseFile() *ast.File {
	if p.trace {
		defer un(trace(p, "File"))
	}

	// Don't bother parsing the rest if we had errors scanning the first token.
	// Likely not a Go source file at all.
	if p.errors.Len() != 0 {
		return nil
	}

	// package clause
	doc := p.leadComment
	pos := p.pos
	ident := &ast.Ident{}

	// Don't bother parsing the rest if we had errors parsing the package clause.
	// Likely not a Go source file at all.
	if p.errors.Len() != 0 {
		return nil
	}

	p.openScope()
	p.pkgScope = p.topScope
	var decls []ast.Decl
	if p.mode&PackageClauseOnly == 0 {
		// import decls
		for p.tok == token.Zh_引入 {
			decls = append(decls, p.parseGenDecl_import(p.tok))
		}

		if p.mode&ImportsOnly == 0 {
			// rest of package body
			for p.tok != token.EOF {
				switch p.tok {
				case token.Zh_类型:
					decls = append(decls, p.parseGenDecl_type(p.tok))
				case token.Zh_常量:
					decls = append(decls, p.parseGenDecl_const(p.tok))
				case token.Zh_全局:
					decls = append(decls, p.parseGenDecl_global(p.tok))
				case token.Zh_结构:
					decls = append(decls, p.parseGenDecl_struct(p.tok))
				case token.Zh_接口:
					decls = append(decls, p.parseGenDecl_interface(p.tok))
				case token.Zh_函数:
					decls = append(decls, p.parseFuncDecl(p.tok))

				default:
					pos := p.pos
					p.errorExpected(pos, "declaration:"+p.lit+p.tok.String())
					p.advance(declStart)
					decls = append(decls, &ast.BadDecl{From: pos, To: p.pos})
				}
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
		W2Mode:     true,
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
