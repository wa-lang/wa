// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import "wa-lang.org/wa/internal/wat/token"

func (p *parser) parseInstruction() {
	switch p.tok {
	default:
		p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)

	case token.INS_I32_STORE:
		p.parseInstruction_i32_store()
	case token.INS_I32_CONST:
		p.parseInstruction_i32_const()

	case token.INS_CALL:
		p.parseInstruction_call()
	case token.INS_DROP:
		p.parseInstruction_drop()
	}
}

func (p *parser) parseInstruction_i32_store() {
	p.acceptToken(token.INS_I32_STORE)

	p.acceptToken(token.LPAREN)
	p.parseInstruction_i32_const()
	p.acceptToken(token.RPAREN)

	p.acceptToken(token.LPAREN)
	p.parseInstruction_i32_const()
	p.acceptToken(token.RPAREN)
}

func (p *parser) parseInstruction_i32_const() {
	p.acceptToken(token.INS_I32_CONST)
	p.parseIntLit()
}

func (p *parser) parseInstruction_call() {
	p.acceptToken(token.INS_CALL)

	p.parseIdent()

	for {
		if p.tok == token.RPAREN {
			return
		}

		if p.tok == token.COMMENT {
			p.parseComment()
			continue
		}

		for {
			if p.tok == token.COMMENT {
				p.parseComment()
				continue
			}

			if p.tok == token.LPAREN {
				p.acceptToken(token.LPAREN)
				if p.tok.IsIsntruction() {
					p.parseInstruction()
					p.acceptToken(token.RPAREN)
					continue
				}

				if p.tok == token.COMMENT {
					p.parseComment()
					continue
				}
			}

			break
		}
	}
}

func (p *parser) parseInstruction_drop() {
	p.acceptToken(token.INS_DROP)

}
