// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// global ::= (global id? globaltype expr)
//
// globaltype ::= valtype | (mut valtype)

// (global $__stack_ptr (mut i32) (i32.const 1024))     ;; index=0
// (global $__heap_max  i32       (i32.const 67108864)) ;; 64MB, 1024 page
func (p *parser) parseModuleSection_global() *ast.Global {
	p.acceptToken(token.GLOBAL)

	g := &ast.Global{}

	p.consumeComments()
	if p.tok == token.IDENT {
		g.Name = p.parseIdent()
	}

	p.consumeComments()
	if p.tok != token.LPAREN {
		g.Type = p.parseNumberType()

		p.consumeComments()
		p.acceptToken(token.LPAREN)

		p.consumeComments()
		p.acceptToken(token.INS_I32_CONST)
		g.Value = p.parseIntLit()

		p.consumeComments()
		p.acceptToken(token.RPAREN)

	} else {
		p.acceptToken(token.LPAREN)
		p.consumeComments()

		p.acceptToken(token.MUT)
		g.Mutable = true

		p.consumeComments()
		g.Type = p.parseNumberType()

		p.consumeComments()
		p.acceptToken(token.RPAREN)

		p.consumeComments()
		p.acceptToken(token.LPAREN)

		p.consumeComments()
		p.acceptToken(token.INS_I32_CONST)
		g.Value = p.parseIntLit()

		p.consumeComments()
		p.acceptToken(token.RPAREN)
	}

	return g
}
