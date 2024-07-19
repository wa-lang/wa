// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *parser) parseInstruction() ast.Instruction {
	switch p.tok {
	default:
		p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
		panic("unreachable")

	case token.INS_I32_STORE:
		return p.parseInstruction_i32_store()
	case token.INS_I32_CONST:
		return p.parseInstruction_i32_const()

	case token.INS_CALL:
		return p.parseInstruction_call()
	case token.INS_DROP:
		return p.parseInstruction_drop()
	}
}

func (p *parser) parseInstruction_i32_store() *ast.Ins_I32Store {
	p.acceptToken(token.INS_I32_STORE)

	ins := &ast.Ins_I32Store{}

	p.acceptToken(token.LPAREN)
	addr := p.parseInstruction_i32_const()
	p.acceptToken(token.RPAREN)

	p.acceptToken(token.LPAREN)
	val := p.parseInstruction_i32_const()
	p.acceptToken(token.RPAREN)

	ins.Offset = uint32(addr.Value)
	ins.Value = val.Value

	return ins
}

func (p *parser) parseInstruction_i32_const() *ast.Ins_I32Const {
	p.acceptToken(token.INS_I32_CONST)
	ins := &ast.Ins_I32Const{}
	ins.Value = int32(p.parseIntLit())
	return ins
}

func (p *parser) parseInstruction_call() *ast.Ins_Call {
	p.acceptToken(token.INS_CALL)

	ins := &ast.Ins_Call{}

	ins.Name = p.parseIdent()

	for {
		if p.tok == token.RPAREN {
			return ins
		}

		p.consumeComments()

		for {
			p.consumeComments()

			if p.tok == token.LPAREN {
				p.acceptToken(token.LPAREN)
				if p.tok.IsIsntruction() {
					ins.Args = append(ins.Args, p.parseInstruction())
					p.acceptToken(token.RPAREN)
					continue
				}
			}

			break
		}
	}
}

func (p *parser) parseInstruction_drop() *ast.Ins_Drop {
	p.acceptToken(token.INS_DROP)
	ins := &ast.Ins_Drop{}
	return ins
}
