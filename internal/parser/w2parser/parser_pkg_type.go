// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2parser

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

// 只支持别名
// 类型 X = A
func (p *parser) parseGenDecl_type(keyword token.Token) *ast.GenDecl {
	if p.trace {
		defer un(trace(p, "TypeSpec"))
	}

	spec := &ast.TypeSpec{}

	spec.Doc = p.leadComment
	tokPos := p.expect(keyword)

	spec.Name = p.parseIdent()
	spec.Assign = p.expect(token.ASSIGN)

	// 只能引用其他类型名字
	spec.Type = p.parseTypeName()
	p.expectSemi()

	spec.Comment = p.lineComment

	return &ast.GenDecl{
		Doc:    spec.Doc,
		TokPos: tokPos,
		Tok:    keyword,
		Specs:  []ast.Spec{spec},
	}
}
