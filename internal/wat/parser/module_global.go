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
//
// (global $bufio.ErrAdvanceTooFar.0 (export "bufio.ErrAdvanceTooFar.0") i32 (i32.const 0))
func (p *parser) parseModuleSection_global() *ast.Global {
	p.acceptToken(token.GLOBAL)

	g := &ast.Global{}

	gInit := false
	gPos := p.pos

	if p.tok == token.IDENT {
		g.Name = p.parseIdent()
	}

	for {
		if p.tok == token.RPAREN {
			break
		}

		// i32/i64/f32/f64
		if p.tok != token.LPAREN {
			if g.Type != 0 {
				p.errorf(p.pos, "bad token %v %v", p.tok, p.lit)
			}
			g.Type = p.parseNumberType()
			continue
		}

		p.acceptToken(token.LPAREN)
		switch p.tok {
		case token.EXPORT:
			p.acceptToken(token.EXPORT)
			g.ExportName = p.parseStringLit()
			p.acceptToken(token.RPAREN)

		case token.MUT: // (mut i32)
			p.acceptToken(token.MUT)
			g.Mutable = true
			if g.Type != 0 {
				p.errorf(p.pos, "bad token %v %v", p.tok, p.lit)
			}
			g.Type = p.parseNumberType()
			p.acceptToken(token.RPAREN)

		case token.INS_I32_CONST:
			if gInit {
				p.errorf(p.pos, "init twice: %v %v", p.tok, p.lit)
			}
			p.acceptToken(token.INS_I32_CONST)
			g.I32Value = p.parseInt32Lit()
			p.acceptToken(token.RPAREN)
			gInit = true
		case token.INS_I64_CONST:
			if gInit {
				p.errorf(p.pos, "init twice: %v %v", p.tok, p.lit)
			}
			p.acceptToken(token.INS_I64_CONST)
			g.I64Value = p.parseInt64Lit()
			p.acceptToken(token.RPAREN)
			gInit = true

		case token.INS_F32_CONST:
			if gInit {
				p.errorf(p.pos, "init twice: %v %v", p.tok, p.lit)
			}
			p.acceptToken(token.INS_F32_CONST)
			g.F32Value = p.parseFloat32Lit()
			p.acceptToken(token.RPAREN)
			gInit = true

		case token.INS_F64_CONST:
			if gInit {
				p.errorf(p.pos, "init twice: %v %v", p.tok, p.lit)
			}
			p.acceptToken(token.INS_I32_CONST)
			g.F64Value = p.parseFloat64Lit()
			p.acceptToken(token.RPAREN)
			gInit = true

		default:
			p.errorf(p.pos, "bad token %v %v", p.tok, p.lit)
		}
	}

	if g.Type == 0 {
		p.errorf(gPos, "missing type")
	}
	if !gInit {
		p.errorf(gPos, "missing init expr")
	}

	return g
}
