// 版权 @2024 凹语言 作者。保留所有权利。

package watutil

import (
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// https://webassembly.github.io/spec/core/binary/instructions.html

var tokOpcodeMap = make(map[token.Token]wasm.Opcode)

func init() {
	for tok := token.INS_UNREACHABLE; tok.IsIsntruction(); tok++ {
		if opcode, ok := wasm.LookupOpcode(tok.String()); ok {
			tokOpcodeMap[tok] = opcode
		}
	}
}

func (p *wat2wasmWorker) appendInstruction(dst []byte, fn *ast.Func, i ast.Instruction) []byte {
	switch i.Token() {
	case token.INS_UNREACHABLE:
		return append(dst, wasm.OpcodeUnreachable)
	case token.INS_NOP:
		return append(dst, wasm.OpcodeNop)
	case token.INS_BLOCK:
		ins := i.(ast.Ins_Block)
		p.enterLabelScope(ins.Label)
		defer p.leaveLabelScope()
		dst = append(dst, wasm.OpcodeBlock)
		for _, x := range ins.List {
			dst = append(dst, p.appendInstruction(dst, fn, x)...)
		}
		dst = append(dst, wasm.OpcodeEnd)
		return dst
	case token.INS_LOOP:
		ins := i.(ast.Ins_Loop)
		p.enterLabelScope(ins.Label)
		defer p.leaveLabelScope()
		dst = append(dst, wasm.OpcodeLoop)
		for _, x := range ins.List {
			dst = append(dst, p.appendInstruction(dst, fn, x)...)
		}
		dst = append(dst, wasm.OpcodeEnd)
		return dst
	case token.INS_IF:
		ins := i.(ast.Ins_If)
		p.enterLabelScope(ins.Label)
		defer p.leaveLabelScope()
		dst = append(dst, wasm.OpcodeIf)
		for _, x := range ins.Body {
			dst = append(dst, p.appendInstruction(dst, fn, x)...)
		}
		if len(ins.Else) > 0 {
			dst = append(dst, wasm.OpcodeElse)
			for _, x := range ins.Else {
				dst = append(dst, p.appendInstruction(dst, fn, x)...)
			}
		}
		dst = append(dst, wasm.OpcodeEnd)
		return dst
	case token.INS_ELSE:
		panic("unreachable")
	case token.INS_END:
		panic("unreachable")
	case token.INS_BR:
		ins := i.(ast.Ins_Br)
		x := p.finfLabelIndex(ins.X)
		dst = append(dst, wasm.OpcodeBr)
		dst = append(dst, p.encodeUint32(x)...)
		return dst
	case token.INS_BR_IF:
		ins := i.(ast.Ins_BrIf)
		x := p.finfLabelIndex(ins.X)
		dst = append(dst, wasm.OpcodeBrIf)
		dst = append(dst, p.encodeUint32(x)...)
		return dst
	case token.INS_BR_TABLE:
		ins := i.(ast.Ins_BrTable)
		dst = append(dst, wasm.OpcodeBrTable)
		for _, x := range ins.XList {
			x := p.finfLabelIndex(x)
			dst = append(dst, p.encodeUint32(x)...)
		}
		// todo: 是否有结尾标志
		return dst
	case token.INS_RETURN:
		return append(dst, wasm.OpcodeReturn)
	case token.INS_CALL:
		ins := i.(ast.Ins_Call)
		x := p.findFuncIndex(ins.X)
		dst = append(dst, wasm.OpcodeCall)
		dst = append(dst, p.encodeUint32(x)...)
		return dst
	case token.INS_CALL_INDIRECT:
		ins := i.(ast.Ins_CallIndirect)
		x := p.findTableIndex(ins.X)
		dst = append(dst, wasm.OpcodeCallIndirect)
		dst = append(dst, p.encodeUint32(x)...)
		return dst
	case token.INS_DROP:
		return append(dst, wasm.OpcodeDrop)
	case token.INS_SELECT:
		return append(dst, wasm.OpcodeSelect)
	case token.INS_TYPED_SELECT:
		ins := i.(ast.Ins_TypedSelect)
		x := p.findTableIndex(ins.Typ)
		dst = append(dst, wasm.OpcodeTypedSelect)
		dst = append(dst, p.encodeUint32(x)...)
		return dst
	case token.INS_LOCAL_GET:
		ins := i.(ast.Ins_LocalGet)
		x := p.findFuncLocalIndex(fn, ins.X)
		dst = append(dst, wasm.OpcodeLocalGet)
		dst = append(dst, p.encodeUint32(x)...)
		return dst
	case token.INS_LOCAL_SET:
		ins := i.(ast.Ins_LocalSet)
		x := p.findFuncLocalIndex(fn, ins.X)
		dst = append(dst, wasm.OpcodeLocalSet)
		dst = append(dst, p.encodeUint32(x)...)
		return dst
	case token.INS_LOCAL_TEE:
		ins := i.(ast.Ins_LocalTee)
		x := p.findFuncLocalIndex(fn, ins.X)
		dst = append(dst, wasm.OpcodeLocalTee)
		dst = append(dst, p.encodeUint32(x)...)
		return dst
	case token.INS_GLOBAL_GET:
		ins := i.(ast.Ins_GlobalGet)
		x := p.findGlobalIndex(ins.X)
		dst = append(dst, wasm.OpcodeGlobalGet)
		dst = append(dst, p.encodeUint32(x)...)
		return dst
	case token.INS_GLOBAL_SET:
		ins := i.(ast.Ins_GlobalSet)
		x := p.findGlobalIndex(ins.X)
		dst = append(dst, wasm.OpcodeGlobalSet)
		dst = append(dst, p.encodeUint32(x)...)
		return dst
	case token.INS_TABLE_GET:
		ins := i.(ast.Ins_TableGet)
		x := p.findTableIndex(ins.X)
		dst = append(dst, wasm.OpcodeTableGet)
		dst = append(dst, p.encodeUint32(x)...)
		return dst
	case token.INS_TABLE_SET:
		ins := i.(ast.Ins_TableSet)
		x := p.findTableIndex(ins.X)
		dst = append(dst, wasm.OpcodeTableSet)
		dst = append(dst, p.encodeUint32(x)...)
		return dst
	case token.INS_I32_LOAD:
		ins := i.(ast.Ins_I32Load)
		dst = append(dst, wasm.OpcodeI32Load)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I64_LOAD:
		ins := i.(ast.Ins_I64Load)
		dst = append(dst, wasm.OpcodeI64Load)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_F32_LOAD:
		ins := i.(ast.Ins_F32Load)
		dst = append(dst, wasm.OpcodeF32Load)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_F64_LOAD:
		ins := i.(ast.Ins_F64Load)
		dst = append(dst, wasm.OpcodeF64Load)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I32_LOAD8_S:
		ins := i.(ast.Ins_I32Load8S)
		dst = append(dst, wasm.OpcodeI32Load8S)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I32_LOAD8_U:
		ins := i.(ast.Ins_I32Load8U)
		dst = append(dst, wasm.OpcodeI32Load8U)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I32_LOAD16_S:
		ins := i.(ast.Ins_I32Load16S)
		dst = append(dst, wasm.OpcodeI32Load16S)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I32_LOAD16_U:
		ins := i.(ast.Ins_I32Load16U)
		dst = append(dst, wasm.OpcodeI32Load16U)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I64_LOAD8_S:
		ins := i.(ast.Ins_I64Load8S)
		dst = append(dst, wasm.OpcodeI64Load8S)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I64_LOAD8_U:
		ins := i.(ast.Ins_I64Load8U)
		dst = append(dst, wasm.OpcodeI64Load8U)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I64_LOAD16_S:
		ins := i.(ast.Ins_I64Load16S)
		dst = append(dst, wasm.OpcodeI64Load16S)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I64_LOAD16_U:
		ins := i.(ast.Ins_I64Load16U)
		dst = append(dst, wasm.OpcodeI64Load16U)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I64_LOAD32_S:
		ins := i.(ast.Ins_I64Load32S)
		dst = append(dst, wasm.OpcodeI64Load32S)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I64_LOAD32_U:
		ins := i.(ast.Ins_I64Load32U)
		dst = append(dst, wasm.OpcodeI64Load32U)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I32_STORE:
		ins := i.(ast.Ins_I32Store)
		dst = append(dst, wasm.OpcodeI32Store)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I64_STORE:
		ins := i.(ast.Ins_I64Store)
		dst = append(dst, wasm.OpcodeI64Store)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_F32_STORE:
		ins := i.(ast.Ins_F32Store)
		dst = append(dst, wasm.OpcodeF32Store)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_F64_STORE:
		ins := i.(ast.Ins_F64Store)
		dst = append(dst, wasm.OpcodeF64Store)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I32_STORE8:
		ins := i.(ast.Ins_I32Store8)
		dst = append(dst, wasm.OpcodeI32Store8)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I32_STORE16:
		ins := i.(ast.Ins_I32Store16)
		dst = append(dst, wasm.OpcodeI32Store16)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I64_STORE8:
		ins := i.(ast.Ins_I64Store8)
		dst = append(dst, wasm.OpcodeI64Store8)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I64_STORE16:
		ins := i.(ast.Ins_I64Store16)
		dst = append(dst, wasm.OpcodeI64Store16)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_I64_STORE32:
		ins := i.(ast.Ins_I64Store32)
		dst = append(dst, wasm.OpcodeI64Store32)
		dst = append(dst, p.encodeUint32(uint32(ins.Align))...)
		dst = append(dst, p.encodeUint32(uint32(ins.Offset))...)
		return dst
	case token.INS_MEMORY_SIZE:
		return append(dst, wasm.OpcodeMemorySize)
	case token.INS_MEMORY_GROW:
		return append(dst, wasm.OpcodeMemoryGrow)
	case token.INS_I32_CONST:
		ins := i.(ast.Ins_I32Const)
		dst = append(dst, wasm.OpcodeI32Const)
		dst = append(dst, p.encodeInt32(ins.X)...)
		return dst
	case token.INS_I64_CONST:
		ins := i.(ast.Ins_I64Const)
		dst = append(dst, wasm.OpcodeI64Const)
		dst = append(dst, p.encodeInt64(ins.X)...)
		return dst
	case token.INS_F32_CONST:
		ins := i.(ast.Ins_F32Const)
		dst = append(dst, wasm.OpcodeF32Const)
		dst = append(dst, p.encodeFloat32(ins.X)...)
		return dst
	case token.INS_F64_CONST:
		ins := i.(ast.Ins_F64Const)
		dst = append(dst, wasm.OpcodeF64Const)
		dst = append(dst, p.encodeFloat64(ins.X)...)
		return dst
	case token.INS_I32_EQZ:
		return append(dst, wasm.OpcodeI32Eqz)
	case token.INS_I32_EQ:
		return append(dst, wasm.OpcodeI32Eq)
	case token.INS_I32_NE:
		return append(dst, wasm.OpcodeI32Ne)
	case token.INS_I32_LT_S:
		return append(dst, wasm.OpcodeI32LtS)
	case token.INS_I32_LT_U:
		return append(dst, wasm.OpcodeI32LtU)
	case token.INS_I32_GT_S:
		return append(dst, wasm.OpcodeI32GtS)
	case token.INS_I32_GT_U:
		return append(dst, wasm.OpcodeI32GtU)
	case token.INS_I32_LE_S:
		return append(dst, wasm.OpcodeI32LeS)
	case token.INS_I32_LE_U:
		return append(dst, wasm.OpcodeI32LeU)
	case token.INS_I32_GE_S:
		return append(dst, wasm.OpcodeI32GeS)
	case token.INS_I32_GE_U:
		return append(dst, wasm.OpcodeI32GeU)
	case token.INS_I64_EQZ:
		return append(dst, wasm.OpcodeI64Eqz)
	case token.INS_I64_EQ:
		return append(dst, wasm.OpcodeI64Eq)
	case token.INS_I64_NE:
		return append(dst, wasm.OpcodeI64Ne)
	case token.INS_I64_LT_S:
		return append(dst, wasm.OpcodeI64LtS)
	case token.INS_I64_LT_U:
		return append(dst, wasm.OpcodeI64LtU)
	case token.INS_I64_GT_S:
		return append(dst, wasm.OpcodeI64GtS)
	case token.INS_I64_GT_U:
		return append(dst, wasm.OpcodeI64GtU)
	case token.INS_I64_LE_S:
		return append(dst, wasm.OpcodeI64LeS)
	case token.INS_I64_LE_U:
		return append(dst, wasm.OpcodeI64LeU)
	case token.INS_I64_GE_S:
		return append(dst, wasm.OpcodeI64GeS)
	case token.INS_I64_GE_U:
		return append(dst, wasm.OpcodeI64GeU)
	case token.INS_F32_EQ:
		return append(dst, wasm.OpcodeF32Eq)
	case token.INS_F32_NE:
		return append(dst, wasm.OpcodeF32Ne)
	case token.INS_F32_LT:
		return append(dst, wasm.OpcodeF32Lt)
	case token.INS_F32_GT:
		return append(dst, wasm.OpcodeF32Gt)
	case token.INS_F32_LE:
		return append(dst, wasm.OpcodeF32Le)
	case token.INS_F32_GE:
		return append(dst, wasm.OpcodeF32Ge)
	case token.INS_F64_EQ:
		return append(dst, wasm.OpcodeF64Eq)
	case token.INS_F64_NE:
		return append(dst, wasm.OpcodeF64Ne)
	case token.INS_F64_LT:
		return append(dst, wasm.OpcodeF64Lt)
	case token.INS_F64_GT:
		return append(dst, wasm.OpcodeF64Gt)
	case token.INS_F64_LE:
		return append(dst, wasm.OpcodeF64Le)
	case token.INS_F64_GE:
		return append(dst, wasm.OpcodeF64Ge)
	case token.INS_I32_CLZ:
		return append(dst, wasm.OpcodeI32Clz)
	case token.INS_I32_CTZ:
		return append(dst, wasm.OpcodeI32Ctz)
	case token.INS_I32_POPCNT:
		return append(dst, wasm.OpcodeI32Popcnt)
	case token.INS_I32_ADD:
		return append(dst, wasm.OpcodeI32Add)
	case token.INS_I32_SUB:
		return append(dst, wasm.OpcodeI32Sub)
	case token.INS_I32_MUL:
		return append(dst, wasm.OpcodeI32Mul)
	case token.INS_I32_DIV_S:
		return append(dst, wasm.OpcodeI32DivS)
	case token.INS_I32_DIV_U:
		return append(dst, wasm.OpcodeI32DivU)
	case token.INS_I32_REM_S:
		return append(dst, wasm.OpcodeI32RemS)
	case token.INS_I32_REM_U:
		return append(dst, wasm.OpcodeI32RemU)
	case token.INS_I32_AND:
		return append(dst, wasm.OpcodeI32And)
	case token.INS_I32_OR:
		return append(dst, wasm.OpcodeI32Or)
	case token.INS_I32_XOR:
		return append(dst, wasm.OpcodeI32Xor)
	case token.INS_I32_SHL:
		return append(dst, wasm.OpcodeI32Shl)
	case token.INS_I32_SHR_S:
		return append(dst, wasm.OpcodeI32ShrS)
	case token.INS_I32_SHR_U:
		return append(dst, wasm.OpcodeI32ShrU)
	case token.INS_I32_ROTL:
		return append(dst, wasm.OpcodeI32Rotl)
	case token.INS_I32_ROTR:
		return append(dst, wasm.OpcodeI32Rotr)
	case token.INS_I64_CLZ:
		return append(dst, wasm.OpcodeI64Clz)
	case token.INS_I64_CTZ:
		return append(dst, wasm.OpcodeI64Ctz)
	case token.INS_I64_POPCNT:
		return append(dst, wasm.OpcodeI64Popcnt)
	case token.INS_I64_ADD:
		return append(dst, wasm.OpcodeI64Add)
	case token.INS_I64_SUB:
		return append(dst, wasm.OpcodeI64Sub)
	case token.INS_I64_MUL:
		return append(dst, wasm.OpcodeI64Mul)
	case token.INS_I64_DIV_S:
		return append(dst, wasm.OpcodeI64DivS)
	case token.INS_I64_DIV_U:
		return append(dst, wasm.OpcodeI64DivU)
	case token.INS_I64_REM_S:
		return append(dst, wasm.OpcodeI64RemS)
	case token.INS_I64_REM_U:
		return append(dst, wasm.OpcodeI64RemU)
	case token.INS_I64_AND:
		return append(dst, wasm.OpcodeI64And)
	case token.INS_I64_OR:
		return append(dst, wasm.OpcodeI64Or)
	case token.INS_I64_XOR:
		return append(dst, wasm.OpcodeI64Xor)
	case token.INS_I64_SHL:
		return append(dst, wasm.OpcodeI64Shl)
	case token.INS_I64_SHR_S:
		return append(dst, wasm.OpcodeI64ShrS)
	case token.INS_I64_SHR_U:
		return append(dst, wasm.OpcodeI64ShrU)
	case token.INS_I64_ROTL:
		return append(dst, wasm.OpcodeI64Rotl)
	case token.INS_I64_ROTR:
		return append(dst, wasm.OpcodeI64Rotr)
	case token.INS_F32_ABS:
		return append(dst, wasm.OpcodeF32Abs)
	case token.INS_F32_NEG:
		return append(dst, wasm.OpcodeF32Neg)
	case token.INS_F32_CEIL:
		return append(dst, wasm.OpcodeF32Ceil)
	case token.INS_F32_FLOOR:
		return append(dst, wasm.OpcodeF32Floor)
	case token.INS_F32_TRUNC:
		return append(dst, wasm.OpcodeF32Trunc)
	case token.INS_F32_NEAREST:
		return append(dst, wasm.OpcodeF32Nearest)
	case token.INS_F32_SQRT:
		return append(dst, wasm.OpcodeF32Sqrt)
	case token.INS_F32_ADD:
		return append(dst, wasm.OpcodeF32Add)
	case token.INS_F32_SUB:
		return append(dst, wasm.OpcodeF32Sub)
	case token.INS_F32_MUL:
		return append(dst, wasm.OpcodeF32Mul)
	case token.INS_F32_DIV:
		return append(dst, wasm.OpcodeF32Div)
	case token.INS_F32_MIN:
		return append(dst, wasm.OpcodeF32Min)
	case token.INS_F32_MAX:
		return append(dst, wasm.OpcodeF32Max)
	case token.INS_F32_COPYSIGN:
		return append(dst, wasm.OpcodeF32Copysign)
	case token.INS_F64_ABS:
		return append(dst, wasm.OpcodeF64Abs)
	case token.INS_F64_NEG:
		return append(dst, wasm.OpcodeF64Neg)
	case token.INS_F64_CEIL:
		return append(dst, wasm.OpcodeF64Ceil)
	case token.INS_F64_FLOOR:
		return append(dst, wasm.OpcodeF64Floor)
	case token.INS_F64_TRUNC:
		return append(dst, wasm.OpcodeF64Trunc)
	case token.INS_F64_NEAREST:
		return append(dst, wasm.OpcodeF64Nearest)
	case token.INS_F64_SQRT:
		return append(dst, wasm.OpcodeF64Sqrt)
	case token.INS_F64_ADD:
		return append(dst, wasm.OpcodeF64Add)
	case token.INS_F64_SUB:
		return append(dst, wasm.OpcodeF64Sub)
	case token.INS_F64_MUL:
		return append(dst, wasm.OpcodeF64Mul)
	case token.INS_F64_DIV:
		return append(dst, wasm.OpcodeF64Div)
	case token.INS_F64_MIN:
		return append(dst, wasm.OpcodeF64Min)
	case token.INS_F64_MAX:
		return append(dst, wasm.OpcodeF64Max)
	case token.INS_F64_COPYSIGN:
		return append(dst, wasm.OpcodeF64Copysign)
	case token.INS_I32_WRAP_I64:
		return append(dst, wasm.OpcodeI32WrapI64)
	case token.INS_I32_TRUNC_F32_S:
		return append(dst, wasm.OpcodeI32TruncF32S)
	case token.INS_I32_TRUNC_F32_U:
		return append(dst, wasm.OpcodeI32TruncF32U)
	case token.INS_I32_TRUNC_F64_S:
		return append(dst, wasm.OpcodeI32TruncF64S)
	case token.INS_I32_TRUNC_F64_U:
		return append(dst, wasm.OpcodeI32TruncF64U)
	case token.INS_I64_EXTEND_I32_S:
		return append(dst, wasm.OpcodeI64Extend32S)
	case token.INS_I64_EXTEND_I32_U:
		return append(dst, wasm.OpcodeI64ExtendI32U)
	case token.INS_I64_TRUNC_F32_S:
		return append(dst, wasm.OpcodeI64TruncF32S)
	case token.INS_I64_TRUNC_F32_U:
		return append(dst, wasm.OpcodeI64TruncF32U)
	case token.INS_I64_TRUNC_F64_S:
		return append(dst, wasm.OpcodeI64TruncF64S)
	case token.INS_I64_TRUNC_F64_U:
		return append(dst, wasm.OpcodeI64TruncF64U)
	case token.INS_F32_CONVERT_I32_S:
		return append(dst, wasm.OpcodeF32ConvertI32S)
	case token.INS_F32_CONVERT_I32_U:
		return append(dst, wasm.OpcodeF32ConvertI32U)
	case token.INS_F32_CONVERT_I64_S:
		return append(dst, wasm.OpcodeF32ConvertI64S)
	case token.INS_F32_CONVERT_I64_U:
		return append(dst, wasm.OpcodeF32ConvertI64U)
	case token.INS_F32_DEMOTE_F64:
		return append(dst, wasm.OpcodeF32DemoteF64)
	case token.INS_F64_CONVERT_I32_S:
		return append(dst, wasm.OpcodeF64ConvertI32S)
	case token.INS_F64_CONVERT_I32_U:
		return append(dst, wasm.OpcodeF64ConvertI32U)
	case token.INS_F64_CONVERT_I64_S:
		return append(dst, wasm.OpcodeF64ConvertI64S)
	case token.INS_F64_CONVERT_I64_U:
		return append(dst, wasm.OpcodeF64ConvertI64U)
	case token.INS_F64_PROMOTE_F32:
		return append(dst, wasm.OpcodeF64PromoteF32)
	case token.INS_I32_REINTERPRET_F32:
		return append(dst, wasm.OpcodeI32ReinterpretF32)
	case token.INS_I64_REINTERPRET_F64:
		return append(dst, wasm.OpcodeI64ReinterpretF64)
	case token.INS_F32_REINTERPRET_I32:
		return append(dst, wasm.OpcodeF32ReinterpretI32)
	case token.INS_F64_REINTERPRET_I64:
		return append(dst, wasm.OpcodeF64ReinterpretI64)
	}
	panic("unreachable")
}
