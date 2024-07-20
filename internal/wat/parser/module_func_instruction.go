// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"strconv"

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
	case token.INS_I32_LOAD:
		p.next()
		// i32.load offset=0 align=1
		// todo

		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		p.parseIntLit()

		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		p.parseIntLit()

		return &ast.Ins_I32Load{}
	case token.INS_I32_CONST:
		return p.parseInstruction_i32_const()

	case token.INS_I64_CONST:
		p.next()
		x := p.parseIntLit()
		return &ast.Ins_I64Const{X: int64(x)}
	case token.INS_I64_STORE:
		p.next()
		return &ast.Ins_I64Store{}
	case token.INS_I64_LOAD:
		p.next()
		return &ast.Ins_I64Load{}

	case token.INS_BR:
		p.next()
		x := p.parseIdent()
		return &ast.Ins_Br{X: x}
	case token.INS_CALL:
		return p.parseInstruction_call()
	case token.INS_CALL_INDIRECT:
		// call_indirect (type $$onFree)
		p.next()
		p.acceptToken(token.LPAREN)
		p.acceptToken(token.TYPE)
		p.parseIdent()
		p.acceptToken(token.RPAREN)
		return &ast.Ins_CallIndirect{}
	case token.INS_UNREACHABLE:
		p.next()
		return &ast.Ins_Unreachable{}

	case token.INS_DROP:
		return p.parseInstruction_drop()
	case token.INS_RETURN:
		p.next()
		return &ast.Ins_Return{}
	case token.INS_IF:
		p.next()

		// if (result i32 i32 i32 i32)
		p.consumeComments()
		if p.tok == token.LPAREN {
			p.acceptToken(token.LPAREN)
			p.acceptToken(token.RESULT)
			p.parseNumberTypeList()
			p.acceptToken(token.RPAREN)
		}

		for {
			p.consumeComments()
			if !p.tok.IsIsntruction() {
				break
			}

			if p.tok == token.INS_ELSE {
				p.next()
				continue
			}

			if p.tok == token.INS_END {
				p.next()
				break
			}

			p.parseInstruction()
		}
		return &ast.Ins_If{}

	case token.INS_LOOP:
		p.next()
		p.parseIdent()
		for {
			p.consumeComments()
			if !p.tok.IsIsntruction() {
				break
			}
			if p.tok == token.INS_END {
				p.next()
				break
			}

			p.parseInstruction()
		}
		return &ast.Ins_If{}

	case token.INS_GLOBAL_GET:
		p.next()
		x := p.parseIdent()
		return &ast.Ins_GlobalGet{X: x}
	case token.INS_GLOBAL_SET:
		p.next()
		x := p.parseIdent()
		return &ast.Ins_GlobalSet{X: x}

	case token.INS_LOCAL_GET:
		p.next()
		var x string
		if p.tok == token.IDENT {
			x = p.parseIdent()
		} else {
			x = strconv.Itoa(p.parseIntLit())
		}
		return &ast.Ins_LocalGet{X: x}
	case token.INS_LOCAL_SET:
		p.next()
		var x string
		if p.tok == token.IDENT {
			x = p.parseIdent()
		} else {
			x = strconv.Itoa(p.parseIntLit())
		}
		return &ast.Ins_LocalSet{X: x}
	case token.INS_LOCAL_TEE:
		p.next()
		var x string
		if p.tok == token.IDENT {
			x = p.parseIdent()
		} else {
			x = strconv.Itoa(p.parseIntLit())
		}
		return &ast.Ins_LocalTee{X: x}

	case token.INS_I32_SUB:
		p.next()
		return &ast.Ins_I32Sub{}
	case token.INS_I32_ADD:
		p.next()
		return &ast.Ins_I32Add{}
	case token.INS_I32_MUL:
		p.next()
		return &ast.Ins_I32Mul{}
	case token.INS_I32_DIV_S:
		p.next()
		return &ast.Ins_I32DivS{}
	case token.INS_I32_DIV_U:
		p.next()
		return &ast.Ins_I32DivU{}
	case token.INS_I32_REM_S:
		p.next()
		return &ast.Ins_I32RemS{}
	case token.INS_I32_REM_U:
		p.next()
		return &ast.Ins_I32RemU{}

	case token.INS_I32_EQ:
		p.next()
		return &ast.Ins_I32Eq{}
	case token.INS_I32_EQZ:
		p.next()
		return &ast.Ins_I32Eqz{}
	}
}

func (p *parser) parseInstruction_i32_store() *ast.Ins_I32Store {
	// i32.store offset=0 align=1

	p.acceptToken(token.INS_I32_STORE)

	p.acceptToken(token.OFFSET)
	p.acceptToken(token.ASSIGN)
	p.parseIntLit()

	p.acceptToken(token.ALIGN)
	p.acceptToken(token.ASSIGN)
	p.parseIntLit()

	ins := &ast.Ins_I32Store{}

	return ins
}

func (p *parser) parseInstruction_i32_const() *ast.Ins_I32Const {
	ins := &ast.Ins_I32Const{}
	p.acceptToken(token.INS_I32_CONST)
	ins.X = int32(p.parseIntLit())
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
