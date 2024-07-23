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

	if p.tok == token.IDENT {
		g.Name = p.parseIdent()
	}

	if p.tok != token.LPAREN {
		g.Type = p.parseNumberType()

		p.acceptToken(token.LPAREN)
		switch g.Type {
		case token.I32:
			p.acceptToken(token.INS_I32_CONST)
			g.I32Value = p.parseInt32Lit()
		case token.I64:
			p.acceptToken(token.INS_I64_CONST)
			g.I64Value = p.parseInt64Lit()
		case token.F32:
			p.acceptToken(token.INS_F32_CONST)
			g.F32Value = p.parseFloat32Lit()
		case token.F64:
			p.acceptToken(token.INS_I32_CONST)
			g.F64Value = p.parseFloat64Lit()
		default:
			p.errorf(p.pos, "bad token %v %v", p.tok, p.lit)
		}

		p.acceptToken(token.RPAREN)

	} else {
		p.acceptToken(token.LPAREN)
		p.acceptToken(token.MUT)
		g.Mutable = true
		g.Type = p.parseNumberType()
		p.acceptToken(token.RPAREN)

		p.acceptToken(token.LPAREN)
		switch g.Type {
		case token.I32:
			p.acceptToken(token.INS_I32_CONST)
			g.I32Value = p.parseInt32Lit()
		case token.I64:
			p.acceptToken(token.INS_I64_CONST)
			g.I64Value = p.parseInt64Lit()
		case token.F32:
			p.acceptToken(token.INS_F32_CONST)
			g.F32Value = p.parseFloat32Lit()
		case token.F64:
			p.acceptToken(token.INS_I32_CONST)
			g.F64Value = p.parseFloat64Lit()
		default:
			p.errorf(p.pos, "bad token %v %v", p.tok, p.lit)
		}

		p.acceptToken(token.RPAREN)
	}

	return g
}
