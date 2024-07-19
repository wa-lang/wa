// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// func ::= (func id? typeuse (t:local)* (in:instr)*)
//
// local := (local id? t:valtype)

func (p *parser) parseModuleSection_func() *ast.Func {
	p.acceptToken(token.FUNC)

	fn := &ast.Func{
		Type: &ast.FuncType{},
		Body: &ast.FuncBody{},
	}

	if p.tok == token.IDENT {
		fn.Name = p.parseIdent()
	}

	// todo

	if p.tok == token.LPAREN {
		p.acceptToken(token.LPAREN)

	Loop0:
		for {
			p.consumeComments()
			switch p.tok {
			case token.EXPORT:
				p.acceptToken(token.EXPORT)

				p.consumeComments()
				fn.ExportName = p.parseStringLit()

				p.consumeComments()
				p.acceptToken(token.RPAREN)

			case token.PARAM:
				var field ast.Field
				p.next()
				if p.tok == token.IDENT {
					field.Name = p.parseIdent()
				}
				switch p.tok {
				case token.I32, token.I64, token.F32, token.F64:
					field.Type = p.tok
					p.next()
					p.consumeComments()
					p.acceptToken(token.RPAREN)
					fn.Type.Params = append(fn.Type.Params, field)
				default:
					p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
				}
			default:
				break Loop0
			}
		}
	}

	for {
		p.consumeComments()
		if p.tok == token.LPAREN {
			p.acceptToken(token.LPAREN)
		} else {
			break
		}

		switch {
		case p.tok.IsIsntruction():
			p.parseInstruction()
			p.acceptToken(token.RPAREN)
		default:
			p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
		}
	}

	return fn
}
