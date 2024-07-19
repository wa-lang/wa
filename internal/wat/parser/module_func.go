// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/token"
)

// func ::= (func id? typeuse (t:local)* (in:instr)*)
//
// local := (local id? t:valtype)

func (p *parser) parseModuleSection_func() {
	p.acceptToken(token.FUNC)

	// (func $main (export "_start")

	if p.tok == token.IDENT {
		p.acceptToken(token.IDENT)
	}

	if p.tok == token.LPAREN {
		p.acceptToken(token.LPAREN)

	Loop0:
		for {
			p.consumeComments()
			switch p.tok {
			case token.EXPORT:
				p.acceptToken(token.EXPORT)
				p.parseStringLit()
				p.acceptToken(token.RPAREN)

			case token.PARAM:
				p.next()
				if p.tok == token.IDENT {
					p.parseIdent()
				}
				switch p.tok {
				case token.I32, token.I64, token.F32, token.F64:
					p.next()
					p.acceptToken(token.RPAREN)
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
}
