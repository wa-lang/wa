// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *parser) parseInstruction() ast.Instruction {
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
	case token.INS_MEMORY_INIT:
		return p.parseIns_MemoryInit()
	case token.INS_MEMORY_COPY:
		return p.parseIns_MemoryCopy()
	case token.INS_MEMORY_FILL:
		return p.parseIns_MemoryFill()
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
	case token.INS_F64_PROMOTE_F32:
		return p.parseIns_F64PromoteF32()
	case token.INS_I32_REINTERPRET_F32:
		return p.parseIns_I32ReintepretF32()
	case token.INS_I64_REINTERPRET_F64:
		return p.parseIns_I64ReintepretF64()
	case token.INS_F32_REINTERPRET_I32:
		return p.parseIns_F32ReintepretI32()
	case token.INS_F64_REINTERPRET_I64:
		return p.parseIns_F64ReintepretI64()
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

	if p.tok == token.IDENT {
		i.Label = p.parseIdent()
	}
	if p.tok == token.LPAREN {
		p.acceptToken(token.LPAREN)
		p.acceptToken(token.RESULT)
		i.Results = p.parseNumberTypeList()
		p.acceptToken(token.RPAREN)
	}

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

	if p.tok == token.IDENT {
		i.Label = p.parseIdent()
	}
	if p.tok == token.LPAREN {
		p.acceptToken(token.LPAREN)
		p.acceptToken(token.RESULT)
		i.Results = p.parseNumberTypeList()
		p.acceptToken(token.RPAREN)
	}

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

	// if $label (result i32 i32 i32 i32)
	if p.tok == token.IDENT {
		i.Label = p.parseIdent()
	}
	if p.tok == token.LPAREN {
		p.acceptToken(token.LPAREN)
		p.acceptToken(token.RESULT)
		i.Results = p.parseNumberTypeList()
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
	i.X = p.parseIdent() // 必须根据名字调用
	return
}

func (p *parser) parseIns_CallIndirect() (i ast.Ins_CallIndirect) {
	// call_indirect $idx? (type $$OnFree)
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_CALL_INDIRECT)
	if p.tok != token.LPAREN {
		i.TableIdx = p.parseIdentOrIndex()
	} else {
		i.TableIdx = "0"
	}
	p.acceptToken(token.LPAREN)
	p.acceptToken(token.TYPE)
	i.TypeIdx = p.parseIdentOrIndex()
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
	// wasm 2.0 支持带类型的T
	if p.tok == token.LPAREN {
		p.acceptToken(token.LPAREN)
		p.acceptToken(token.RESULT)
		tokTyp := p.tok
		p.acceptToken(token.I32, token.I64, token.F32, token.F64)
		i.ResultTyp = tokTyp
		p.acceptToken(token.RPAREN)
	}
	return
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
	i.TableIdx = p.parseIdentOrIndex()
	return
}
func (p *parser) parseIns_TableSet() (i ast.Ins_TableSet) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_TABLE_SET)
	i.TableIdx = p.parseIdentOrIndex()
	return
}

func (p *parser) parseIns_I32Load() (i ast.Ins_I32Load) {
	// i32.load offset=0 align=1
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_LOAD)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 4
	}
	return
}
func (p *parser) parseIns_I64Load() (i ast.Ins_I64Load) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_LOAD)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 8
	}
	return
}
func (p *parser) parseIns_F32Load() (i ast.Ins_F32Load) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_LOAD)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 4
	}
	return
}
func (p *parser) parseIns_F64Load() (i ast.Ins_F64Load) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_LOAD)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 8
	}
	return
}

func (p *parser) parseIns_I32Load8S() (i ast.Ins_I32Load8S) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_LOAD8_S)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 1
	}
	return
}
func (p *parser) parseIns_I32Load8U() (i ast.Ins_I32Load8U) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_LOAD8_U)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 1
	}
	return
}
func (p *parser) parseIns_I32Load16S() (i ast.Ins_I32Load16S) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_LOAD16_S)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 2
	}
	return
}
func (p *parser) parseIns_I32Load16U() (i ast.Ins_I32Load16U) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_LOAD16_U)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 2
	}
	return
}
func (p *parser) parseIns_I64Load8S() (i ast.Ins_I64Load8S) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_LOAD8_S)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 1
	}
	return
}
func (p *parser) parseIns_I64Load8U() (i ast.Ins_I64Load8U) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_LOAD8_U)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 1
	}
	return
}
func (p *parser) parseIns_I64Load16S() (i ast.Ins_I64Load16S) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_LOAD16_S)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 2
	}
	return
}
func (p *parser) parseIns_I64Load16U() (i ast.Ins_I64Load16U) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_LOAD16_U)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 2
	}
	return
}
func (p *parser) parseIns_I64Load32S() (i ast.Ins_I64Load32S) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_LOAD32_S)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 4
	}
	return
}
func (p *parser) parseIns_I64Load32U() (i ast.Ins_I64Load32U) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_LOAD32_U)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 4
	}
	return
}
func (p *parser) parseIns_I32Store() (i ast.Ins_I32Store) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_STORE)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 4
	}
	return
}
func (p *parser) parseIns_I64Store() (i ast.Ins_I64Store) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_STORE)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 8
	}
	return
}
func (p *parser) parseIns_F32Store() (i ast.Ins_F32Store) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_STORE)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 4
	}
	return
}
func (p *parser) parseIns_F64Store() (i ast.Ins_F64Store) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_STORE)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 8
	}
	return
}
func (p *parser) parseIns_I32Store8() (i ast.Ins_I32Store8) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_STORE8)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 1
	}
	return
}
func (p *parser) parseIns_I32Store16() (i ast.Ins_I32Store16) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_STORE16)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 2
	}
	return
}
func (p *parser) parseIns_I64Store8() (i ast.Ins_I64Store8) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_STORE8)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 1
	}
	return
}
func (p *parser) parseIns_I64Store16() (i ast.Ins_I64Store16) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_STORE16)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 2
	}
	return
}
func (p *parser) parseIns_I64Store32() (i ast.Ins_I64Store32) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_STORE32)
	if p.tok == token.OFFSET {
		p.acceptToken(token.OFFSET)
		p.acceptToken(token.ASSIGN)
		i.Offset = uint(p.parseIntLit())
	} else {
		i.Offset = 0
	}
	if p.tok == token.ALIGN {
		p.acceptToken(token.ALIGN)
		p.acceptToken(token.ASSIGN)
		i.Align = uint(p.parseIntLit())
	} else {
		i.Align = 4
	}
	return
}
func (p *parser) parseIns_MemorySize() (i ast.Ins_MemorySize) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_MEMORY_SIZE)
	return
}
func (p *parser) parseIns_MemoryGrow() (i ast.Ins_MemoryGrow) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_MEMORY_GROW)
	return
}
func (p *parser) parseIns_MemoryInit() (i ast.Ins_MemoryInit) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_MEMORY_INIT)
	i.DataIdx = p.parseInt32Lit() // TODO: 支持标识符
	return
}
func (p *parser) parseIns_MemoryCopy() (i ast.Ins_MemoryCopy) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_MEMORY_COPY)
	return
}
func (p *parser) parseIns_MemoryFill() (i ast.Ins_MemoryFill) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_MEMORY_FILL)
	return
}
func (p *parser) parseIns_I32Const() (i ast.Ins_I32Const) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_CONST)
	i.X = p.parseInt32Lit()
	return
}
func (p *parser) parseIns_I64Const() (i ast.Ins_I64Const) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_CONST)
	i.X = p.parseInt64Lit()
	return
}
func (p *parser) parseIns_F32Const() (i ast.Ins_F32Const) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_CONST)
	i.X = p.parseFloat32Lit()
	return
}
func (p *parser) parseIns_F64Const() (i ast.Ins_F64Const) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_CONST)
	i.X = p.parseFloat64Lit()
	return
}
func (p *parser) parseIns_I32Eqz() (i ast.Ins_I32Eqz) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_EQZ)
	return
}
func (p *parser) parseIns_I32Eq() (i ast.Ins_I32Eq) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_EQ)
	return
}
func (p *parser) parseIns_I32Ne() (i ast.Ins_I32Ne) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_NE)
	return
}
func (p *parser) parseIns_I32LtS() (i ast.Ins_I32LtS) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_LT_S)
	return
}
func (p *parser) parseIns_I32LtU() (i ast.Ins_I32LtU) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_LT_U)
	return
}
func (p *parser) parseIns_I32GtS() (i ast.Ins_I32GtS) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_GT_S)
	return
}
func (p *parser) parseIns_I32GtU() (i ast.Ins_I32GtU) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_GT_U)
	return
}
func (p *parser) parseIns_I32LeS() (i ast.Ins_I32LeS) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_LE_S)
	return
}
func (p *parser) parseIns_I32LeU() (i ast.Ins_I32LeU) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_LE_U)
	return
}
func (p *parser) parseIns_I32GeS() (i ast.Ins_I32GeS) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_GE_S)
	return
}
func (p *parser) parseIns_I32GeU() (i ast.Ins_I32GeU) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_GE_U)
	return
}
func (p *parser) parseIns_I64Eqz() (i ast.Ins_I64Eqz) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_EQZ)
	return
}
func (p *parser) parseIns_I64Eq() (i ast.Ins_I64Eq) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_EQ)
	return
}
func (p *parser) parseIns_I64Ne() (i ast.Ins_I64Ne) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_NE)
	return
}
func (p *parser) parseIns_I64LtS() (i ast.Ins_I64LtS) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_LT_S)
	return
}
func (p *parser) parseIns_I64LtU() (i ast.Ins_I64LtU) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_LT_U)
	return
}
func (p *parser) parseIns_I64GtS() (i ast.Ins_I64GtS) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_GT_S)
	return
}
func (p *parser) parseIns_I64GtU() (i ast.Ins_I64GtU) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_GT_U)
	return
}
func (p *parser) parseIns_I64LeS() (i ast.Ins_I64LeS) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_LE_S)
	return
}
func (p *parser) parseIns_I64LeU() (i ast.Ins_I64LeU) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_LE_U)
	return
}
func (p *parser) parseIns_I64GeS() (i ast.Ins_I64GeS) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_GE_S)
	return
}
func (p *parser) parseIns_I64GeU() (i ast.Ins_I64GeU) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_GE_U)
	return
}
func (p *parser) parseIns_F32Eq() (i ast.Ins_F32Eq) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_EQ)
	return
}
func (p *parser) parseIns_F32Ne() (i ast.Ins_F32Ne) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_NE)
	return
}
func (p *parser) parseIns_F32Lt() (i ast.Ins_F32Lt) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_LT)
	return
}
func (p *parser) parseIns_F32Gt() (i ast.Ins_F32Gt) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_GT)
	return
}
func (p *parser) parseIns_F32Le() (i ast.Ins_F32Le) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_LE)
	return
}
func (p *parser) parseIns_F32Ge() (i ast.Ins_F32Ge) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_GE)
	return
}
func (p *parser) parseIns_F64Eq() (i ast.Ins_F64Eq) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_EQ)
	return
}
func (p *parser) parseIns_F64Ne() (i ast.Ins_F64Ne) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_NE)
	return
}
func (p *parser) parseIns_F64Lt() (i ast.Ins_F64Lt) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_LT)
	return
}
func (p *parser) parseIns_F64Gt() (i ast.Ins_F64Gt) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_GT)
	return
}
func (p *parser) parseIns_F64Le() (i ast.Ins_F64Le) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_LE)
	return
}
func (p *parser) parseIns_F64Ge() (i ast.Ins_F64Ge) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_GE)
	return
}
func (p *parser) parseIns_I32Clz() (i ast.Ins_I32Clz) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_CLZ)
	return
}
func (p *parser) parseIns_I32Ctz() (i ast.Ins_I32Ctz) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_CTZ)
	return
}
func (p *parser) parseIns_I32Popcnt() (i ast.Ins_I32Popcnt) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_POPCNT)
	return
}
func (p *parser) parseIns_I32Add() (i ast.Ins_I32Add) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_ADD)
	return
}
func (p *parser) parseIns_I32Sub() (i ast.Ins_I32Sub) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_SUB)
	return
}
func (p *parser) parseIns_I32Mul() (i ast.Ins_I32Mul) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_MUL)
	return
}
func (p *parser) parseIns_I32DivS() (i ast.Ins_I32DivS) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_DIV_S)
	return
}
func (p *parser) parseIns_I32DivU() (i ast.Ins_I32DivU) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_DIV_U)
	return
}
func (p *parser) parseIns_I32RemS() (i ast.Ins_I32RemS) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_REM_S)
	return
}
func (p *parser) parseIns_I32RemU() (i ast.Ins_I32RemU) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_REM_U)
	return
}
func (p *parser) parseIns_I32And() (i ast.Ins_I32And) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_AND)
	return
}
func (p *parser) parseIns_I32Or() (i ast.Ins_I32Or) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_OR)
	return
}
func (p *parser) parseIns_I32Xor() (i ast.Ins_I32Xor) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_XOR)
	return
}
func (p *parser) parseIns_I32Shl() (i ast.Ins_I32Shl) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_SHL)
	return
}
func (p *parser) parseIns_I32ShrS() (i ast.Ins_I32ShrS) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_SHR_S)
	return
}
func (p *parser) parseIns_I32ShrU() (i ast.Ins_I32ShrU) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_SHR_U)
	return
}
func (p *parser) parseIns_I32Rotl() (i ast.Ins_I32Rotl) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_ROTL)
	return
}
func (p *parser) parseIns_I32Rotr() (i ast.Ins_I32Rotr) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_ROTR)
	return
}
func (p *parser) parseIns_I64Clz() (i ast.Ins_I64Clz) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_CLZ)
	return
}
func (p *parser) parseIns_I64Ctz() (i ast.Ins_I64Ctz) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_CTZ)
	return
}
func (p *parser) parseIns_I64Popcnt() (i ast.Ins_I64Popcnt) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_POPCNT)
	return
}
func (p *parser) parseIns_I64Add() (i ast.Ins_I64Add) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_ADD)
	return
}
func (p *parser) parseIns_I64Sub() (i ast.Ins_I64Sub) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_SUB)
	return
}
func (p *parser) parseIns_I64Mul() (i ast.Ins_I64Mul) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_MUL)
	return
}
func (p *parser) parseIns_I64DivS() (i ast.Ins_I64DivS) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_DIV_S)
	return
}
func (p *parser) parseIns_I64DivU() (i ast.Ins_I64DivU) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_DIV_U)
	return
}
func (p *parser) parseIns_I64RemS() (i ast.Ins_I64RemS) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_REM_S)
	return
}
func (p *parser) parseIns_I64RemU() (i ast.Ins_I64RemU) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_REM_U)
	return
}
func (p *parser) parseIns_I64And() (i ast.Ins_I64And) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_AND)
	return
}
func (p *parser) parseIns_I64Or() (i ast.Ins_I64Or) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_OR)
	return
}
func (p *parser) parseIns_I64Xor() (i ast.Ins_I64Xor) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_XOR)
	return
}
func (p *parser) parseIns_I64Shl() (i ast.Ins_I64Shl) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_SHL)
	return
}
func (p *parser) parseIns_I64ShrS() (i ast.Ins_I64ShrS) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_SHR_S)
	return
}
func (p *parser) parseIns_I64ShrU() (i ast.Ins_I64ShrU) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_SHR_U)
	return
}
func (p *parser) parseIns_I64Rotl() (i ast.Ins_I64Rotl) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_ROTL)
	return
}
func (p *parser) parseIns_I64Rotr() (i ast.Ins_I64Rotr) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_ROTR)
	return
}
func (p *parser) parseIns_F32Abs() (i ast.Ins_F32Abs) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_ABS)
	return
}
func (p *parser) parseIns_F32Neg() (i ast.Ins_F32Neg) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_NEG)
	return
}
func (p *parser) parseIns_F32Ceil() (i ast.Ins_F32Ceil) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_CEIL)
	return
}
func (p *parser) parseIns_F32Floor() (i ast.Ins_F32Floor) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_FLOOR)
	return
}
func (p *parser) parseIns_F32Trunc() (i ast.Ins_F32Trunc) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_TRUNC)
	return
}
func (p *parser) parseIns_F32Nearest() (i ast.Ins_F32Nearest) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_NEAREST)
	return
}
func (p *parser) parseIns_F32Sqrt() (i ast.Ins_F32Sqrt) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_SQRT)
	return
}
func (p *parser) parseIns_F32Add() (i ast.Ins_F32Add) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_ADD)
	return
}
func (p *parser) parseIns_F32Sub() (i ast.Ins_F32Sub) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_SUB)
	return
}
func (p *parser) parseIns_F32Mul() (i ast.Ins_F32Mul) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_MUL)
	return
}
func (p *parser) parseIns_F32Div() (i ast.Ins_F32Div) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_DIV)
	return
}
func (p *parser) parseIns_F32Min() (i ast.Ins_F32Min) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_MIN)
	return
}
func (p *parser) parseIns_F32Max() (i ast.Ins_F32Max) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_MAX)
	return
}
func (p *parser) parseIns_F32Copysign() (i ast.Ins_F32Copysign) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_COPYSIGN)
	return
}
func (p *parser) parseIns_F64Abs() (i ast.Ins_F64Abs) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_ABS)
	return
}
func (p *parser) parseIns_F64Neg() (i ast.Ins_F64Neg) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_NEG)
	return
}
func (p *parser) parseIns_F64Ceil() (i ast.Ins_F64Ceil) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_CEIL)
	return
}
func (p *parser) parseIns_F64Floor() (i ast.Ins_F64Floor) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_FLOOR)
	return
}
func (p *parser) parseIns_F64Trunc() (i ast.Ins_F64Trunc) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_TRUNC)
	return
}
func (p *parser) parseIns_F64Nearest() (i ast.Ins_F64Nearest) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_NEAREST)
	return
}
func (p *parser) parseIns_F64Sqrt() (i ast.Ins_F64Sqrt) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_SQRT)
	return
}
func (p *parser) parseIns_F64Add() (i ast.Ins_F64Add) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_ADD)
	return
}
func (p *parser) parseIns_F64Sub() (i ast.Ins_F64Sub) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_SUB)
	return
}
func (p *parser) parseIns_F64Mul() (i ast.Ins_F64Mul) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_MUL)
	return
}
func (p *parser) parseIns_F64Div() (i ast.Ins_F64Div) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_DIV)
	return
}
func (p *parser) parseIns_F64Min() (i ast.Ins_F64Min) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_MIN)
	return
}
func (p *parser) parseIns_F64Max() (i ast.Ins_F64Max) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_MAX)
	return
}
func (p *parser) parseIns_F64Copysign() (i ast.Ins_F64Copysign) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_COPYSIGN)
	return
}
func (p *parser) parseIns_I32WrapI64() (i ast.Ins_I32WrapI64) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_WRAP_I64)
	return
}
func (p *parser) parseIns_I32TruncF32S() (i ast.Ins_I32TruncF32S) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_TRUNC_F32_S)
	return
}
func (p *parser) parseIns_I32TruncF32U() (i ast.Ins_I32TruncF32U) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_TRUNC_F32_U)
	return
}
func (p *parser) parseIns_I32TruncF64S() (i ast.Ins_I32TruncF64S) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_TRUNC_F64_S)
	return
}
func (p *parser) parseIns_I32TruncF64U() (i ast.Ins_I32TruncF64U) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_TRUNC_F64_U)
	return
}
func (p *parser) parseIns_I64ExtendI32S() (i ast.Ins_I64ExtendI32S) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_EXTEND_I32_S)
	return
}
func (p *parser) parseIns_I64ExtendI32U() (i ast.Ins_I64ExtendI32U) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_EXTEND_I32_U)
	return
}
func (p *parser) parseIns_I64TruncF32S() (i ast.Ins_I64TruncF32S) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_TRUNC_F32_S)
	return
}
func (p *parser) parseIns_I64TruncF32U() (i ast.Ins_I64TruncF32U) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_TRUNC_F32_U)
	return
}
func (p *parser) parseIns_I64TruncF64S() (i ast.Ins_I64TruncF64S) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_TRUNC_F64_S)
	return
}
func (p *parser) parseIns_I64TruncF64U() (i ast.Ins_I64TruncF64U) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_TRUNC_F64_U)
	return
}
func (p *parser) parseIns_F32ConvertI32S() (i ast.Ins_F32ConvertI32S) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_CONVERT_I32_S)
	return
}
func (p *parser) parseIns_F32ConvertI32U() (i ast.Ins_F32ConvertI32U) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_CONVERT_I32_U)
	return
}
func (p *parser) parseIns_F32ConvertI64S() (i ast.Ins_F32ConvertI64S) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_CONVERT_I64_S)
	return
}
func (p *parser) parseIns_F32ConvertI64U() (i ast.Ins_F32ConvertI64U) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_CONVERT_I64_U)
	return
}
func (p *parser) parseIns_F32DemoteF64() (i ast.Ins_F32DemoteF64) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_DEMOTE_F64)
	return
}
func (p *parser) parseIns_F64ConvertI32S() (i ast.Ins_F64ConvertI32S) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_CONVERT_I32_S)
	return
}
func (p *parser) parseIns_F64ConvertI32U() (i ast.Ins_F64ConvertI32U) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_CONVERT_I32_U)
	return
}
func (p *parser) parseIns_F64ConvertI64S() (i ast.Ins_F64ConvertI64S) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_CONVERT_I64_S)
	return
}
func (p *parser) parseIns_F64ConvertI64U() (i ast.Ins_F64ConvertI64U) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_CONVERT_I64_U)
	return
}
func (p *parser) parseIns_F64PromoteF32() (i ast.Ins_F64PromoteF32) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_PROMOTE_F32)
	return
}
func (p *parser) parseIns_I32ReintepretF32() (i ast.Ins_I32ReintepretF32) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I32_REINTERPRET_F32)
	return
}
func (p *parser) parseIns_I64ReintepretF64() (i ast.Ins_I64ReintepretF64) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_I64_REINTERPRET_F64)
	return
}
func (p *parser) parseIns_F32ReintepretI32() (i ast.Ins_F32ReintepretI32) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F32_REINTERPRET_I32)
	return
}
func (p *parser) parseIns_F64ReintepretI64() (i ast.Ins_F64ReintepretI64) {
	i.OpToken = ast.OpToken(p.tok)
	p.acceptToken(token.INS_F64_REINTERPRET_I64)
	return
}
