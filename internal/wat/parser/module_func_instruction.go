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

	return ins
}

func (p *parser) parseInstruction_i32_const() *ast.Ins_I32Const {
	ins := &ast.Ins_I32Const{}
	p.acceptToken(token.INS_I32_CONST)
	ins.Value = int32(p.parseIntLit())
	return ins
}

func (p *parser) parseInstruction_call() *ast.Ins_Call {
	ins := &ast.Ins_Call{}
	p.acceptToken(token.INS_CALL)
	ins.Name = p.parseIdent()
	return ins
}

func (p *parser) parseInstruction_drop() *ast.Ins_Drop {
	ins := &ast.Ins_Drop{}
	p.acceptToken(token.INS_DROP)
	return ins
}
