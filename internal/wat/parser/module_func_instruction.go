// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"strconv"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *parser) parseInstruction() ast.Instruction {

	switch p.tok {
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

	switch p.tok {
	case token.INS_UNREACHABLE:
		return p.parseIns_Unreachable()
	case token.INS_NOP:
		return p.parseIns_Nop()
	case token.INS_BLOCK:
		return p.parseIns_Block()
	case token.INS_LOOP:
		return p.parseIns_Loop()
	case token.INS_IF:
		return p.parseIns_If()
	case token.INS_ELSE:
		return p.parseIns_Else()
	case token.INS_END:
		return p.parseIns_End()
	case token.INS_BR:
		return p.parseIns_Br()
	case token.INS_BR_IF:
		return p.parseIns_BrIf()
	case token.INS_BR_TABLE:
		return p.parseIns_BrTable()
	case token.INS_RETURN:
		return p.parseIns_Return()
	case token.INS_CALL:
		return p.parseIns_Call()
	case token.INS_CALL_INDIRECT:
		return p.parseIns_CallIndirect()
	case token.INS_DROP:
		return p.parseIns_Drop()
	case token.INS_SELECT:
		return p.parseIns_Select()
	case token.INS_TYPED_SELECT:
		return p.parseIns_TypedSelect()
	case token.INS_LOCAL_GET:
		return p.parseIns_LocalGet()
	case token.INS_LOCAL_SET:
		return p.parseIns_LocalSet()
	case token.INS_LOCAL_TEE:
		return p.parseIns_LocalTee()
	case token.INS_GLOBAL_GET:
		return p.parseIns_GlobalGet()
	case token.INS_GLOBAL_SET:
		return p.parseIns_GlobalSet()
	case token.INS_TABLE_GET:
		return p.parseIns_TableGet()
	case token.INS_TABLE_SET:
		return p.parseIns_TableSet()
	case token.INS_I32_LOAD:
		return p.parseIns_I32Load()
	case token.INS_I64_LOAD:
		return p.parseIns_I64Load()
	case token.INS_F32_LOAD:
		return p.parseIns_F32Load()
	case token.INS_F64_LOAD:
		return p.parseIns_F64Load()
	case token.INS_I32_LOAD8_S:
		return p.parseIns_I32Load8S()
	case token.INS_I32_LOAD8_U:
		return p.parseIns_I32Load8U()
	case token.INS_I32_LOAD16_S:
		return p.parseIns_I32Load16S()
	case token.INS_I32_LOAD16_U:
		return p.parseIns_I32Load16U()
	case token.INS_I64_LOAD8_S:
		return p.parseIns_I64Load8S()
	case token.INS_I64_LOAD8_U:
		return p.parseIns_I64Load8U()
	case token.INS_I64_LOAD16_S:
		return p.parseIns_I64Load16S()
	case token.INS_I64_LOAD16_U:
		return p.parseIns_I64Load16U()
	case token.INS_I64_LOAD32_S:
		return p.parseIns_I64Load32S()
	case token.INS_I64_LOAD32_U:
		return p.parseIns_I64Load32U()
	case token.INS_I32_STORE:
		return p.parseIns_I32Store()
	case token.INS_I64_STORE:
		return p.parseIns_I64Store()
	case token.INS_F32_STORE:
		return p.parseIns_F32Store()
	case token.INS_F64_STORE:
		return p.parseIns_F64Store()
	case token.INS_I32_STORE8:
		return p.parseIns_I32Store8()
	case token.INS_I32_STORE16:
		return p.parseIns_I32Store16()
	case token.INS_I64_STORE8:
		return p.parseIns_I64Store8()
	case token.INS_I64_STORE16:
		return p.parseIns_I64Store16()
	case token.INS_I64_STORE32:
		return p.parseIns_I64Store32()
	case token.INS_MEMORY_SIZE:
		return p.parseIns_MemorySize()
	case token.INS_MEMORY_GROW:
		return p.parseIns_MemoryGrow()
	case token.INS_I32_CONST:
		return p.parseIns_I32Const()
	case token.INS_I64_CONST:
		return p.parseIns_I64Const()
	case token.INS_F32_CONST:
		return p.parseIns_F32Const()
	case token.INS_F64_CONST:
		return p.parseIns_F64Const()
	case token.INS_I32_EQZ:
		return p.parseIns_I32Eqz()
	case token.INS_I32_EQ:
		return p.parseIns_I32Eq()
	case token.INS_I32_NE:
		return p.parseIns_I32Ne()
	case token.INS_I32_LT_S:
		return p.parseIns_I32LtS()
	case token.INS_I32_LT_U:
		return p.parseIns_I32LtU()
	case token.INS_I32_GT_S:
		return p.parseIns_I32GtS()
	case token.INS_I32_GT_U:
		return p.parseIns_I32GtU()
	case token.INS_I32_LE_S:
		return p.parseIns_I32LeS()
	case token.INS_I32_LE_U:
		return p.parseIns_I32LeU()
	case token.INS_I32_GE_S:
		return p.parseIns_I32GeS()
	case token.INS_I32_GE_U:
		return p.parseIns_I32GeU()
	case token.INS_I64_EQZ:
		return p.parseIns_I64Eqz()
	case token.INS_I64_EQ:
		return p.parseIns_I64Eq()
	case token.INS_I64_NE:
		return p.parseIns_I64Ne()
	case token.INS_I64_LT_S:
		return p.parseIns_I64LtS()
	case token.INS_I64_LT_U:
		return p.parseIns_I64LtU()
	case token.INS_I64_GT_S:
		return p.parseIns_I64GtS()
	case token.INS_I64_GT_U:
		return p.parseIns_I64GtU()
	case token.INS_I64_LE_S:
		return p.parseIns_I64LeS()
	case token.INS_I64_LE_U:
		return p.parseIns_I64LeU()
	case token.INS_I64_GE_S:
		return p.parseIns_I64GeS()
	case token.INS_I64_GE_U:
		return p.parseIns_I64GeU()
	case token.INS_F32_EQ:
		return p.parseIns_F32Eq()
	case token.INS_F32_NE:
		return p.parseIns_F32Ne()
	case token.INS_F32_LT:
		return p.parseIns_F32Lt()
	case token.INS_F32_GT:
		return p.parseIns_F32Gt()
	case token.INS_F32_LE:
		return p.parseIns_F32Le()
	case token.INS_F32_GE:
		return p.parseIns_F32Ge()
	case token.INS_F64_EQ:
		return p.parseIns_F64Eq()
	case token.INS_F64_NE:
		return p.parseIns_F64Ne()
	case token.INS_F64_LT:
		return p.parseIns_F64Lt()
	case token.INS_F64_GT:
		return p.parseIns_F64Gt()
	case token.INS_F64_LE:
		return p.parseIns_F64Le()
	case token.INS_F64_GE:
		return p.parseIns_F64Ge()
	case token.INS_I32_CLZ:
		return p.parseIns_I32Clz()
	case token.INS_I32_CTZ:
		return p.parseIns_I32Ctz()
	case token.INS_I32_POPCNT:
		return p.parseIns_I32Popcnt()
	case token.INS_I32_ADD:
		return p.parseIns_I32Add()
	case token.INS_I32_SUB:
		return p.parseIns_I32Sub()
	case token.INS_I32_MUL:
		return p.parseIns_I32Mul()
	case token.INS_I32_DIV_S:
		return p.parseIns_I32DivS()
	case token.INS_I32_DIV_U:
		return p.parseIns_I32DivU()
	case token.INS_I32_REM_S:
		return p.parseIns_I32RemS()
	case token.INS_I32_REM_U:
		return p.parseIns_I32RemU()
	case token.INS_I32_AND:
		return p.parseIns_I32And()
	case token.INS_I32_OR:
		return p.parseIns_I32Or()
	case token.INS_I32_XOR:
		return p.parseIns_I32Xor()
	case token.INS_I32_SHL:
		return p.parseIns_I32Shl()
	case token.INS_I32_SHR_S:
		return p.parseIns_I32ShrS()
	case token.INS_I32_SHR_U:
		return p.parseIns_I32ShrU()
	case token.INS_I32_ROTL:
		return p.parseIns_I32Rotl()
	case token.INS_I32_ROTR:
		return p.parseIns_I32Rotr()
	case token.INS_I64_CLZ:
		return p.parseIns_I64Clz()
	case token.INS_I64_CTZ:
		return p.parseIns_I64Ctz()
	case token.INS_I64_POPCNT:
		return p.parseIns_I64Popcnt()
	case token.INS_I64_ADD:
		return p.parseIns_I64Add()
	case token.INS_I64_SUB:
		return p.parseIns_I64Sub()
	case token.INS_I64_MUL:
		return p.parseIns_I64Mul()
	case token.INS_I64_DIV_S:
		return p.parseIns_I64DivS()
	case token.INS_I64_DIV_U:
		return p.parseIns_I64DivU()
	case token.INS_I64_REM_S:
		return p.parseIns_I64RemS()
	case token.INS_I64_REM_U:
		return p.parseIns_I64RemU()
	case token.INS_I64_AND:
		return p.parseIns_I64And()
	case token.INS_I64_OR:
		return p.parseIns_I64Or()
	case token.INS_I64_XOR:
		return p.parseIns_I64Xor()
	case token.INS_I64_SHL:
		return p.parseIns_I64Shl()
	case token.INS_I64_SHR_S:
		return p.parseIns_I64ShrS()
	case token.INS_I64_SHR_U:
		return p.parseIns_I64ShrU()
	case token.INS_I64_ROTL:
		return p.parseIns_I64Rotl()
	case token.INS_I64_ROTR:
		return p.parseIns_I64Rotr()
	case token.INS_F32_ABS:
		return p.parseIns_F32Abs()
	case token.INS_F32_NEG:
		return p.parseIns_F32Neg()
	case token.INS_F32_CEIL:
		return p.parseIns_F32Ceil()
	case token.INS_F32_FLOOR:
		return p.parseIns_F32Floor()
	case token.INS_F32_TRUNC:
		return p.parseIns_F32Trunc()
	case token.INS_F32_NEAREST:
		return p.parseIns_F32Nearest()
	case token.INS_F32_SQRT:
		return p.parseIns_F32Sqrt()
	case token.INS_F32_ADD:
		return p.parseIns_F32Add()
	case token.INS_F32_SUB:
		return p.parseIns_F32Sub()
	case token.INS_F32_MUL:
		return p.parseIns_F32Mul()
	case token.INS_F32_DIV:
		return p.parseIns_F32Div()
	case token.INS_F32_MIN:
		return p.parseIns_F32Min()
	case token.INS_F32_MAX:
		return p.parseIns_F32Max()
	case token.INS_F32_COPYSIGN:
		return p.parseIns_F32Copysign()
	case token.INS_F64_ABS:
		return p.parseIns_F64Abs()
	case token.INS_F64_NEG:
		return p.parseIns_F64Neg()
	case token.INS_F64_CEIL:
		return p.parseIns_F64Ceil()
	case token.INS_F64_FLOOR:
		return p.parseIns_F64Floor()
	case token.INS_F64_TRUNC:
		return p.parseIns_F64Trunc()
	case token.INS_F64_NEAREST:
		return p.parseIns_F64Nearest()
	case token.INS_F64_SQRT:
		return p.parseIns_F64Sqrt()
	case token.INS_F64_ADD:
		return p.parseIns_F64Add()
	case token.INS_F64_SUB:
		return p.parseIns_F64Sub()
	case token.INS_F64_MUL:
		return p.parseIns_F64Mul()
	case token.INS_F64_DIV:
		return p.parseIns_F64Div()
	case token.INS_F64_MIN:
		return p.parseIns_F64Min()
	case token.INS_F64_MAX:
		return p.parseIns_F64Max()
	case token.INS_F64_COPYSIGN:
		return p.parseIns_F64Copysign()
	case token.INS_I32_WRAP_I64:
		return p.parseIns_I32WrapI64()
	case token.INS_I32_TRUNC_F32_S:
		return p.parseIns_I32TruncF32S()
	case token.INS_I32_TRUNC_F32_U:
		return p.parseIns_I32TruncF32U()
	case token.INS_I32_TRUNC_F64_S:
		return p.parseIns_I32TruncF64S()
	case token.INS_I32_TRUNC_F64_U:
		return p.parseIns_I32TruncF64U()
	case token.INS_I64_EXTEND_I32_S:
		return p.parseIns_I64ExtendI32S()
	case token.INS_I64_EXTEND_I32_U:
		return p.parseIns_I64ExtendI32U()
	case token.INS_I64_TRUNC_F32_S:
		return p.parseIns_I64TruncF32S()
	case token.INS_I64_TRUNC_F32_U:
		return p.parseIns_I64TruncF32U()
	case token.INS_I64_TRUNC_F64_S:
		return p.parseIns_I64TruncF64S()
	case token.INS_I64_TRUNC_F64_U:
		return p.parseIns_I64TruncF64U()
	case token.INS_F32_CONVERT_I32_S:
		return p.parseIns_F32ConvertI32S()
	case token.INS_F32_CONVERT_I32_U:
		return p.parseIns_F32ConvertI32U()
	case token.INS_F32_CONVERT_I64_S:
		return p.parseIns_F32ConvertI64S()
	case token.INS_F32_CONVERT_I64_U:
		return p.parseIns_F32ConvertI64U()
	case token.INS_F32_DEMOTE_F64:
		return p.parseIns_F32DemoteF64()
	case token.INS_F64_CONVERT_I32_S:
		return p.parseIns_F64ConvertI32S()
	case token.INS_F64_CONVERT_I32_U:
		return p.parseIns_F64ConvertI32U()
	case token.INS_F64_CONVERT_I64_S:
		return p.parseIns_F64ConvertI64S()
	case token.INS_F64_CONVERT_I64_U:
		return p.parseIns_F64ConvertI64U()
	case token.INS_F64_DEMOTE_F32:
		return p.parseIns_F64DemoteF32()
	case token.INS_I32_REINTERPRET_F32:
		return p.parseIns_I32ReintepretF32()
	case token.INS_I64_REINTERPRET_F64:
		return p.parseIns_I64ReintepretF64()
	case token.INS_I32_REINTERPRET_I32:
		return p.parseIns_I32ReintepretI32()
	case token.INS_I64_REINTERPRET_I64:
		return p.parseIns_I64ReintepretI64()
	}

	p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
	panic("unreachable")
}

func (p *parser) parseIns_Unreachable() (i ast.Ins_Unreachable) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_UNREACHABLE)
	return
}
func (p *parser) parseIns_Nop() (i ast.Ins_Nop) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_DROP)
	return
}
func (p *parser) parseIns_Block() (i ast.Ins_Block) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_BLOCK)

	p.consumeComments()
	i.X = p.parseIdentOrIndex()

	for {
		p.consumeComments()
		if p.tok == token.INS_END {
			break
		}
		i.List = append(i.List, p.parseInstruction())
	}

	p.consumeComments()
	p.acceptToken(token.INS_END)
	return
}
func (p *parser) parseIns_Loop() (i ast.Ins_Loop) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_LOOP)

	p.consumeComments()
	i.X = p.parseIdentOrIndex()

	for {
		p.consumeComments()
		if p.tok == token.INS_END {
			break
		}
		i.List = append(i.List, p.parseInstruction())
	}

	p.consumeComments()
	p.acceptToken(token.INS_END)
	return
}

func (p *parser) parseIns_If() (i ast.Ins_If) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_IF)

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
		if p.tok == token.INS_END || p.tok == token.INS_ELSE {
			break
		}
		i.Body = append(i.Body, p.parseInstruction())
	}
	if p.tok == token.INS_ELSE {
		p.acceptToken(token.INS_ELSE)

		for {
			p.consumeComments()
			if p.tok == token.INS_END {
				break
			}
			i.Else = append(i.Else, p.parseInstruction())
		}
	}

	p.consumeComments()
	p.acceptToken(token.INS_END)
	return
}

func (p *parser) parseIns_Else() ast.Ins_Else {
	panic("unreachable")
}
func (p *parser) parseIns_End() ast.Ins_End {
	panic("unreachable")
}

func (p *parser) parseIns_Br() (i ast.Ins_Br) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_BR)
	i.X = p.parseIdentOrIndex()
	return
}
func (p *parser) parseIns_BrIf() (i ast.Ins_BrIf) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_BR_IF)

	i.X = p.parseIdentOrIndex()

	for {
		p.consumeComments()
		if p.tok == token.INS_END || p.tok == token.INS_ELSE {
			break
		}
		i.Body = append(i.Body, p.parseInstruction())
	}
	if p.tok == token.INS_ELSE {
		p.acceptToken(token.INS_ELSE)

		for {
			p.consumeComments()
			if p.tok == token.INS_END {
				break
			}
			i.Else = append(i.Else, p.parseInstruction())
		}
	}

	p.consumeComments()
	p.acceptToken(token.INS_END)
	return
}

func (p *parser) parseIns_BrTable() (i ast.Ins_BrTable) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_BR_TABLE)
	i.XList = p.parseIdentOrIndexList()
	return
}
func (p *parser) parseIns_Return() (i ast.Ins_Return) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_RETURN)
	return
}
func (p *parser) parseIns_Call() (i ast.Ins_Call) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_CALL)
	i.X = p.parseIdentOrIndex()
	return
}

func (p *parser) parseIns_CallIndirect() (i ast.Ins_CallIndirect) {
	// call_indirect (type $$onFree)
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_CALL_INDIRECT)
	p.acceptToken(token.LPAREN)
	p.acceptToken(token.TYPE)
	i.X = p.parseIdentOrIndex()
	p.acceptToken(token.RPAREN)
	return
}
func (p *parser) parseIns_Drop() (i ast.Ins_Drop) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_DROP)
	return
}
func (p *parser) parseIns_Select() (i ast.Ins_Select) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_SELECT)
	i.X = p.parseIdentOrIndex()
	i.Y = p.parseIdentOrIndex()
	return
}
func (p *parser) parseIns_TypedSelect() ast.Ins_TypedSelect {
	panic("TODO")
}

func (p *parser) parseIns_LocalGet() (i ast.Ins_LocalGet) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_LOCAL_GET)
	i.X = p.parseIdentOrIndex()
	return
}
func (p *parser) parseIns_LocalSet() (i ast.Ins_LocalSet) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_LOCAL_SET)
	i.X = p.parseIdentOrIndex()
	return
}
func (p *parser) parseIns_LocalTee() (i ast.Ins_LocalTee) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_LOCAL_TEE)
	i.X = p.parseIdentOrIndex()
	return
}

func (p *parser) parseIns_GlobalGet() (i ast.Ins_GlobalGet) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_GLOBAL_GET)
	i.X = p.parseIdentOrIndex()
	return
}
func (p *parser) parseIns_GlobalSet() (i ast.Ins_GlobalSet) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_GLOBAL_SET)
	i.X = p.parseIdentOrIndex()
	return
}

func (p *parser) parseIns_TableGet() (i ast.Ins_TableGet) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_TABLE_GET)
	i.X = p.parseIdentOrIndex()
	return
}
func (p *parser) parseIns_TableSet() (i ast.Ins_TableSet) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_TABLE_SET)
	i.X = p.parseIdentOrIndex()
	return
}

func (p *parser) parseIns_I32Load() (i ast.Ins_I32Load) {
	panic("TODO")
}
func (p *parser) parseIns_I64Load() (i ast.Ins_I64Load) {
	panic("TODO")
}
func (p *parser) parseIns_F32Load() (i ast.Ins_F32Load) {
	panic("TODO")
}
func (p *parser) parseIns_F64Load() (i ast.Ins_F64Load) {
	panic("TODO")
}

func (p *parser) parseIns_I32Load8S() ast.Ins_I32Load8S               { panic("TODO") }
func (p *parser) parseIns_I32Load8U() ast.Ins_I32Load8U               { panic("TODO") }
func (p *parser) parseIns_I32Load16S() ast.Ins_I32Load16S             { panic("TODO") }
func (p *parser) parseIns_I32Load16U() ast.Ins_I32Load16U             { panic("TODO") }
func (p *parser) parseIns_I64Load8S() ast.Ins_I64Load8S               { panic("TODO") }
func (p *parser) parseIns_I64Load8U() ast.Ins_I64Load8U               { panic("TODO") }
func (p *parser) parseIns_I64Load16S() ast.Ins_I64Load16S             { panic("TODO") }
func (p *parser) parseIns_I64Load16U() ast.Ins_I64Load16U             { panic("TODO") }
func (p *parser) parseIns_I64Load32S() ast.Ins_I64Load32S             { panic("TODO") }
func (p *parser) parseIns_I64Load32U() ast.Ins_I64Load32U             { panic("TODO") }
func (p *parser) parseIns_I32Store() ast.Ins_I32Store                 { panic("TODO") }
func (p *parser) parseIns_I64Store() ast.Ins_I64Store                 { panic("TODO") }
func (p *parser) parseIns_F32Store() ast.Ins_F32Store                 { panic("TODO") }
func (p *parser) parseIns_F64Store() ast.Ins_F64Store                 { panic("TODO") }
func (p *parser) parseIns_I32Store8() ast.Ins_I32Store8               { panic("TODO") }
func (p *parser) parseIns_I32Store16() ast.Ins_I32Store16             { panic("TODO") }
func (p *parser) parseIns_I64Store8() ast.Ins_I64Store8               { panic("TODO") }
func (p *parser) parseIns_I64Store16() ast.Ins_I64Store16             { panic("TODO") }
func (p *parser) parseIns_I64Store32() ast.Ins_I64Store32             { panic("TODO") }
func (p *parser) parseIns_MemorySize() ast.Ins_MemorySize             { panic("TODO") }
func (p *parser) parseIns_MemoryGrow() ast.Ins_MemoryGrow             { panic("TODO") }
func (p *parser) parseIns_I32Const() ast.Ins_I32Const                 { panic("TODO") }
func (p *parser) parseIns_I64Const() ast.Ins_I64Const                 { panic("TODO") }
func (p *parser) parseIns_F32Const() ast.Ins_F32Const                 { panic("TODO") }
func (p *parser) parseIns_F64Const() ast.Ins_F64Const                 { panic("TODO") }
func (p *parser) parseIns_I32Eqz() ast.Ins_I32Eqz                     { panic("TODO") }
func (p *parser) parseIns_I32Eq() ast.Ins_I32Eq                       { panic("TODO") }
func (p *parser) parseIns_I32Ne() ast.Ins_I32Ne                       { panic("TODO") }
func (p *parser) parseIns_I32LtS() ast.Ins_I32LtS                     { panic("TODO") }
func (p *parser) parseIns_I32LtU() ast.Ins_I32LtU                     { panic("TODO") }
func (p *parser) parseIns_I32GtS() ast.Ins_I32GtS                     { panic("TODO") }
func (p *parser) parseIns_I32GtU() ast.Ins_I32GtU                     { panic("TODO") }
func (p *parser) parseIns_I32LeS() ast.Ins_I32LeS                     { panic("TODO") }
func (p *parser) parseIns_I32LeU() ast.Ins_I32LeU                     { panic("TODO") }
func (p *parser) parseIns_I32GeS() ast.Ins_I32GeS                     { panic("TODO") }
func (p *parser) parseIns_I32GeU() ast.Ins_I32GeU                     { panic("TODO") }
func (p *parser) parseIns_I64Eqz() ast.Ins_I64Eqz                     { panic("TODO") }
func (p *parser) parseIns_I64Eq() ast.Ins_I64Eq                       { panic("TODO") }
func (p *parser) parseIns_I64Ne() ast.Ins_I64Ne                       { panic("TODO") }
func (p *parser) parseIns_I64LtS() ast.Ins_I64LtS                     { panic("TODO") }
func (p *parser) parseIns_I64LtU() ast.Ins_I64LtU                     { panic("TODO") }
func (p *parser) parseIns_I64GtS() ast.Ins_I64GtS                     { panic("TODO") }
func (p *parser) parseIns_I64GtU() ast.Ins_I64GtU                     { panic("TODO") }
func (p *parser) parseIns_I64LeS() ast.Ins_I64LeS                     { panic("TODO") }
func (p *parser) parseIns_I64LeU() ast.Ins_I64LeU                     { panic("TODO") }
func (p *parser) parseIns_I64GeS() ast.Ins_I64GeS                     { panic("TODO") }
func (p *parser) parseIns_I64GeU() ast.Ins_I64GeU                     { panic("TODO") }
func (p *parser) parseIns_F32Eq() ast.Ins_F32Eq                       { panic("TODO") }
func (p *parser) parseIns_F32Ne() ast.Ins_F32Ne                       { panic("TODO") }
func (p *parser) parseIns_F32Lt() ast.Ins_F32Lt                       { panic("TODO") }
func (p *parser) parseIns_F32Gt() ast.Ins_F32Gt                       { panic("TODO") }
func (p *parser) parseIns_F32Le() ast.Ins_F32Le                       { panic("TODO") }
func (p *parser) parseIns_F32Ge() ast.Ins_F32Ge                       { panic("TODO") }
func (p *parser) parseIns_F64Eq() ast.Ins_F64Eq                       { panic("TODO") }
func (p *parser) parseIns_F64Ne() ast.Ins_F64Ne                       { panic("TODO") }
func (p *parser) parseIns_F64Lt() ast.Ins_F64Lt                       { panic("TODO") }
func (p *parser) parseIns_F64Gt() ast.Ins_F64Gt                       { panic("TODO") }
func (p *parser) parseIns_F64Le() ast.Ins_F64Le                       { panic("TODO") }
func (p *parser) parseIns_F64Ge() ast.Ins_F64Ge                       { panic("TODO") }
func (p *parser) parseIns_I32Clz() ast.Ins_I32Clz                     { panic("TODO") }
func (p *parser) parseIns_I32Ctz() ast.Ins_I32Ctz                     { panic("TODO") }
func (p *parser) parseIns_I32Popcnt() ast.Ins_I32Popcnt               { panic("TODO") }
func (p *parser) parseIns_I32Add() ast.Ins_I32Add                     { panic("TODO") }
func (p *parser) parseIns_I32Sub() ast.Ins_I32Sub                     { panic("TODO") }
func (p *parser) parseIns_I32Mul() ast.Ins_I32Mul                     { panic("TODO") }
func (p *parser) parseIns_I32DivS() ast.Ins_I32DivS                   { panic("TODO") }
func (p *parser) parseIns_I32DivU() ast.Ins_I32DivU                   { panic("TODO") }
func (p *parser) parseIns_I32RemS() ast.Ins_I32RemS                   { panic("TODO") }
func (p *parser) parseIns_I32RemU() ast.Ins_I32RemU                   { panic("TODO") }
func (p *parser) parseIns_I32And() ast.Ins_I32And                     { panic("TODO") }
func (p *parser) parseIns_I32Or() ast.Ins_I32Or                       { panic("TODO") }
func (p *parser) parseIns_I32Xor() ast.Ins_I32Xor                     { panic("TODO") }
func (p *parser) parseIns_I32Shl() ast.Ins_I32Shl                     { panic("TODO") }
func (p *parser) parseIns_I32ShrS() ast.Ins_I32ShrS                   { panic("TODO") }
func (p *parser) parseIns_I32ShrU() ast.Ins_I32ShrU                   { panic("TODO") }
func (p *parser) parseIns_I32Rotl() ast.Ins_I32Rotl                   { panic("TODO") }
func (p *parser) parseIns_I32Rotr() ast.Ins_I32Rotr                   { panic("TODO") }
func (p *parser) parseIns_I64Clz() ast.Ins_I64Clz                     { panic("TODO") }
func (p *parser) parseIns_I64Ctz() ast.Ins_I64Ctz                     { panic("TODO") }
func (p *parser) parseIns_I64Popcnt() ast.Ins_I64Popcnt               { panic("TODO") }
func (p *parser) parseIns_I64Add() ast.Ins_I64Add                     { panic("TODO") }
func (p *parser) parseIns_I64Sub() ast.Ins_I64Sub                     { panic("TODO") }
func (p *parser) parseIns_I64Mul() ast.Ins_I64Mul                     { panic("TODO") }
func (p *parser) parseIns_I64DivS() ast.Ins_I64DivS                   { panic("TODO") }
func (p *parser) parseIns_I64DivU() ast.Ins_I64DivU                   { panic("TODO") }
func (p *parser) parseIns_I64RemS() ast.Ins_I64RemS                   { panic("TODO") }
func (p *parser) parseIns_I64RemU() ast.Ins_I64RemU                   { panic("TODO") }
func (p *parser) parseIns_I64And() ast.Ins_I64And                     { panic("TODO") }
func (p *parser) parseIns_I64Or() ast.Ins_I64Or                       { panic("TODO") }
func (p *parser) parseIns_I64Xor() ast.Ins_I64Xor                     { panic("TODO") }
func (p *parser) parseIns_I64Shl() ast.Ins_I64Shl                     { panic("TODO") }
func (p *parser) parseIns_I64ShrS() ast.Ins_I64ShrS                   { panic("TODO") }
func (p *parser) parseIns_I64ShrU() ast.Ins_I64ShrU                   { panic("TODO") }
func (p *parser) parseIns_I64Rotl() ast.Ins_I64Rotl                   { panic("TODO") }
func (p *parser) parseIns_I64Rotr() ast.Ins_I64Rotr                   { panic("TODO") }
func (p *parser) parseIns_F32Abs() ast.Ins_F32Abs                     { panic("TODO") }
func (p *parser) parseIns_F32Neg() ast.Ins_F32Neg                     { panic("TODO") }
func (p *parser) parseIns_F32Ceil() ast.Ins_F32Ceil                   { panic("TODO") }
func (p *parser) parseIns_F32Floor() ast.Ins_F32Floor                 { panic("TODO") }
func (p *parser) parseIns_F32Trunc() ast.Ins_F32Trunc                 { panic("TODO") }
func (p *parser) parseIns_F32Nearest() ast.Ins_F32Nearest             { panic("TODO") }
func (p *parser) parseIns_F32Sqrt() ast.Ins_F32Sqrt                   { panic("TODO") }
func (p *parser) parseIns_F32Add() ast.Ins_F32Add                     { panic("TODO") }
func (p *parser) parseIns_F32Sub() ast.Ins_F32Sub                     { panic("TODO") }
func (p *parser) parseIns_F32Mul() ast.Ins_F32Mul                     { panic("TODO") }
func (p *parser) parseIns_F32Div() ast.Ins_F32Div                     { panic("TODO") }
func (p *parser) parseIns_F32Min() ast.Ins_F32Min                     { panic("TODO") }
func (p *parser) parseIns_F32Max() ast.Ins_F32Max                     { panic("TODO") }
func (p *parser) parseIns_F32Copysign() ast.Ins_F32Copysign           { panic("TODO") }
func (p *parser) parseIns_F64Abs() ast.Ins_F64Abs                     { panic("TODO") }
func (p *parser) parseIns_F64Neg() ast.Ins_F64Neg                     { panic("TODO") }
func (p *parser) parseIns_F64Ceil() ast.Ins_F64Ceil                   { panic("TODO") }
func (p *parser) parseIns_F64Floor() ast.Ins_F64Floor                 { panic("TODO") }
func (p *parser) parseIns_F64Trunc() ast.Ins_F64Trunc                 { panic("TODO") }
func (p *parser) parseIns_F64Nearest() ast.Ins_F64Nearest             { panic("TODO") }
func (p *parser) parseIns_F64Sqrt() ast.Ins_F64Sqrt                   { panic("TODO") }
func (p *parser) parseIns_F64Add() ast.Ins_F64Add                     { panic("TODO") }
func (p *parser) parseIns_F64Sub() ast.Ins_F64Sub                     { panic("TODO") }
func (p *parser) parseIns_F64Mul() ast.Ins_F64Mul                     { panic("TODO") }
func (p *parser) parseIns_F64Div() ast.Ins_F64Div                     { panic("TODO") }
func (p *parser) parseIns_F64Min() ast.Ins_F64Min                     { panic("TODO") }
func (p *parser) parseIns_F64Max() ast.Ins_F64Max                     { panic("TODO") }
func (p *parser) parseIns_F64Copysign() ast.Ins_F64Copysign           { panic("TODO") }
func (p *parser) parseIns_I32WrapI64() ast.Ins_I32WrapI64             { panic("TODO") }
func (p *parser) parseIns_I32TruncF32S() ast.Ins_I32TruncF32S         { panic("TODO") }
func (p *parser) parseIns_I32TruncF32U() ast.Ins_I32TruncF32U         { panic("TODO") }
func (p *parser) parseIns_I32TruncF64S() ast.Ins_I32TruncF64S         { panic("TODO") }
func (p *parser) parseIns_I32TruncF64U() ast.Ins_I32TruncF64U         { panic("TODO") }
func (p *parser) parseIns_I64ExtendI32S() ast.Ins_I64ExtendI32S       { panic("TODO") }
func (p *parser) parseIns_I64ExtendI32U() ast.Ins_I64ExtendI32U       { panic("TODO") }
func (p *parser) parseIns_I64TruncF32S() ast.Ins_I64TruncF32S         { panic("TODO") }
func (p *parser) parseIns_I64TruncF32U() ast.Ins_I64TruncF32U         { panic("TODO") }
func (p *parser) parseIns_I64TruncF64S() ast.Ins_I64TruncF64S         { panic("TODO") }
func (p *parser) parseIns_I64TruncF64U() ast.Ins_I64TruncF64U         { panic("TODO") }
func (p *parser) parseIns_F32ConvertI32S() ast.Ins_F32ConvertI32S     { panic("TODO") }
func (p *parser) parseIns_F32ConvertI32U() ast.Ins_F32ConvertI32U     { panic("TODO") }
func (p *parser) parseIns_F32ConvertI64S() ast.Ins_F32ConvertI64S     { panic("TODO") }
func (p *parser) parseIns_F32ConvertI64U() ast.Ins_F32ConvertI64U     { panic("TODO") }
func (p *parser) parseIns_F32DemoteF64() ast.Ins_F32DemoteF64         { panic("TODO") }
func (p *parser) parseIns_F64ConvertI32S() ast.Ins_F64ConvertI32S     { panic("TODO") }
func (p *parser) parseIns_F64ConvertI32U() ast.Ins_F64ConvertI32U     { panic("TODO") }
func (p *parser) parseIns_F64ConvertI64S() ast.Ins_F64ConvertI64S     { panic("TODO") }
func (p *parser) parseIns_F64ConvertI64U() ast.Ins_F64ConvertI64U     { panic("TODO") }
func (p *parser) parseIns_F64DemoteF32() ast.Ins_F64DemoteF32         { panic("TODO") }
func (p *parser) parseIns_I32ReintepretF32() ast.Ins_I32ReintepretF32 { panic("TODO") }
func (p *parser) parseIns_I64ReintepretF64() ast.Ins_I64ReintepretF64 { panic("TODO") }
func (p *parser) parseIns_I32ReintepretI32() ast.Ins_I32ReintepretI32 { panic("TODO") }
func (p *parser) parseIns_I64ReintepretI64() ast.Ins_I64ReintepretI64 { panic("TODO") }

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
	ins.X = p.parseIdent()
	return ins
}

func (p *parser) parseInstruction_drop() *ast.Ins_Drop {
	ins := &ast.Ins_Drop{}
	p.acceptToken(token.INS_DROP)
	return ins
}
