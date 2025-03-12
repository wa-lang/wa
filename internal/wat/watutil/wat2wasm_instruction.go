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

func (p *wat2wasmWorker) buildInstruction(dst *wasm.Code, fn *ast.Func, i ast.Instruction) {
	switch i.Token() {
	case token.INS_UNREACHABLE:
		dst.Body = append(dst.Body, wasm.OpcodeUnreachable)
	case token.INS_NOP:
		dst.Body = append(dst.Body, wasm.OpcodeNop)
	case token.INS_BLOCK:
		ins := i.(ast.Ins_Block)
		p.enterLabelScope(ins.Label)
		defer p.leaveLabelScope()
		dst.Body = append(dst.Body, wasm.OpcodeBlock)

		// block type
		if n := len(ins.Results); n == 0 {
			dst.Body = append(dst.Body, 0x40)
		} else if n == 1 {
			dst.Body = append(dst.Body, p.buildValueType(ins.Results[0]))
		} else {
			idx := p.findBlockTypeIndex(ins.Results)
			dst.Body = append(dst.Body, p.encodeInt32(idx)...)
		}

		for _, x := range ins.List {
			p.buildInstruction(dst, fn, x)
		}
		dst.Body = append(dst.Body, wasm.OpcodeEnd)
	case token.INS_LOOP:
		ins := i.(ast.Ins_Loop)
		p.enterLabelScope(ins.Label)
		defer p.leaveLabelScope()
		dst.Body = append(dst.Body, wasm.OpcodeLoop)

		// block type
		if n := len(ins.Results); n == 0 {
			dst.Body = append(dst.Body, 0x40)
		} else if n == 1 {
			dst.Body = append(dst.Body, p.buildValueType(ins.Results[0]))
		} else {
			idx := p.findBlockTypeIndex(ins.Results)
			dst.Body = append(dst.Body, p.encodeInt32(idx)...)
		}

		for _, x := range ins.List {
			p.buildInstruction(dst, fn, x)
		}
		dst.Body = append(dst.Body, wasm.OpcodeEnd)

	case token.INS_IF:
		ins := i.(ast.Ins_If)
		p.enterLabelScope(ins.Label)
		defer p.leaveLabelScope()
		dst.Body = append(dst.Body, wasm.OpcodeIf)

		// block type
		if n := len(ins.Results); n == 0 {
			dst.Body = append(dst.Body, 0x40)
		} else if n == 1 {
			dst.Body = append(dst.Body, p.buildValueType(ins.Results[0]))
		} else {
			idx := p.findBlockTypeIndex(ins.Results)
			dst.Body = append(dst.Body, p.encodeInt32(idx)...)
		}

		for _, x := range ins.Body {
			p.buildInstruction(dst, fn, x)
		}
		if len(ins.Else) > 0 {
			dst.Body = append(dst.Body, wasm.OpcodeElse)
			for _, x := range ins.Else {
				p.buildInstruction(dst, fn, x)
			}
		}
		dst.Body = append(dst.Body, wasm.OpcodeEnd)

	case token.INS_ELSE:
		panic("unreachable")
	case token.INS_END:
		panic("unreachable")
	case token.INS_BR:
		ins := i.(ast.Ins_Br)
		x := p.findLabelIndex(ins.X)
		dst.Body = append(dst.Body, wasm.OpcodeBr)
		dst.Body = append(dst.Body, p.encodeUint32(x)...)

	case token.INS_BR_IF:
		ins := i.(ast.Ins_BrIf)
		x := p.findLabelIndex(ins.X)
		dst.Body = append(dst.Body, wasm.OpcodeBrIf)
		dst.Body = append(dst.Body, p.encodeUint32(x)...)

	case token.INS_BR_TABLE:
		ins := i.(ast.Ins_BrTable)
		dst.Body = append(dst.Body, wasm.OpcodeBrTable)
		dst.Body = append(dst.Body, p.encodeUint32(uint32(len(ins.XList)-1))...)
		for _, x := range ins.XList {
			x := p.findLabelIndex(x)
			dst.Body = append(dst.Body, p.encodeUint32(x)...)
		}

	case token.INS_RETURN:
		dst.Body = append(dst.Body, wasm.OpcodeReturn)
	case token.INS_CALL:
		ins := i.(ast.Ins_Call)
		x := p.findFuncIndex(ins.X)
		dst.Body = append(dst.Body, wasm.OpcodeCall)
		dst.Body = append(dst.Body, p.encodeUint32(x)...)

	case token.INS_CALL_INDIRECT:
		ins := i.(ast.Ins_CallIndirect)
		tableIdx := p.findTableIndex(ins.TableIdx)
		typeIdx := p.findTypeIndexByIdent(ins.TypeIdx)
		dst.Body = append(dst.Body, wasm.OpcodeCallIndirect)
		dst.Body = append(dst.Body, p.encodeUint32(typeIdx)...)
		dst.Body = append(dst.Body, p.encodeUint32(tableIdx)...)

	case token.INS_DROP:
		dst.Body = append(dst.Body, wasm.OpcodeDrop)
	case token.INS_SELECT:
		ins := i.(ast.Ins_Select)
		if ins.ResultTyp != 0 {
			dst.Body = append(dst.Body, wasm.OpcodeTypedSelect)
			switch ins.ResultTyp {
			case token.I32:
				dst.Body = append(dst.Body, wasm.ValueTypeI32)
			case token.I64:
				dst.Body = append(dst.Body, wasm.ValueTypeI64)
			case token.F32:
				dst.Body = append(dst.Body, wasm.ValueTypeF32)
			case token.F64:
				dst.Body = append(dst.Body, wasm.ValueTypeF64)
			default:
				panic("unreachable")
			}
		} else {
			dst.Body = append(dst.Body, wasm.OpcodeSelect)
		}

	case token.INS_LOCAL_GET:
		ins := i.(ast.Ins_LocalGet)
		x := p.findFuncLocalIndex(fn, ins.X)
		dst.Body = append(dst.Body, wasm.OpcodeLocalGet)
		dst.Body = append(dst.Body, p.encodeUint32(x)...)

	case token.INS_LOCAL_SET:
		ins := i.(ast.Ins_LocalSet)
		x := p.findFuncLocalIndex(fn, ins.X)
		dst.Body = append(dst.Body, wasm.OpcodeLocalSet)
		dst.Body = append(dst.Body, p.encodeUint32(x)...)

	case token.INS_LOCAL_TEE:
		ins := i.(ast.Ins_LocalTee)
		x := p.findFuncLocalIndex(fn, ins.X)
		dst.Body = append(dst.Body, wasm.OpcodeLocalTee)
		dst.Body = append(dst.Body, p.encodeUint32(x)...)

	case token.INS_GLOBAL_GET:
		ins := i.(ast.Ins_GlobalGet)
		x := p.findGlobalIndex(ins.X)
		dst.Body = append(dst.Body, wasm.OpcodeGlobalGet)
		dst.Body = append(dst.Body, p.encodeUint32(x)...)

	case token.INS_GLOBAL_SET:
		ins := i.(ast.Ins_GlobalSet)
		x := p.findGlobalIndex(ins.X)
		dst.Body = append(dst.Body, wasm.OpcodeGlobalSet)
		dst.Body = append(dst.Body, p.encodeUint32(x)...)

	case token.INS_TABLE_GET:
		ins := i.(ast.Ins_TableGet)
		x := p.findTableIndex(ins.TableIdx)
		dst.Body = append(dst.Body, wasm.OpcodeTableGet)
		dst.Body = append(dst.Body, p.encodeUint32(x)...)

	case token.INS_TABLE_SET:
		ins := i.(ast.Ins_TableSet)
		x := p.findTableIndex(ins.TableIdx)
		dst.Body = append(dst.Body, wasm.OpcodeTableSet)
		dst.Body = append(dst.Body, p.encodeUint32(x)...)

	case token.INS_I32_LOAD:
		ins := i.(ast.Ins_I32Load)
		dst.Body = append(dst.Body, wasm.OpcodeI32Load)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I64_LOAD:
		ins := i.(ast.Ins_I64Load)
		dst.Body = append(dst.Body, wasm.OpcodeI64Load)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_F32_LOAD:
		ins := i.(ast.Ins_F32Load)
		dst.Body = append(dst.Body, wasm.OpcodeF32Load)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_F64_LOAD:
		ins := i.(ast.Ins_F64Load)
		dst.Body = append(dst.Body, wasm.OpcodeF64Load)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I32_LOAD8_S:
		ins := i.(ast.Ins_I32Load8S)
		dst.Body = append(dst.Body, wasm.OpcodeI32Load8S)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I32_LOAD8_U:
		ins := i.(ast.Ins_I32Load8U)
		dst.Body = append(dst.Body, wasm.OpcodeI32Load8U)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I32_LOAD16_S:
		ins := i.(ast.Ins_I32Load16S)
		dst.Body = append(dst.Body, wasm.OpcodeI32Load16S)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I32_LOAD16_U:
		ins := i.(ast.Ins_I32Load16U)
		dst.Body = append(dst.Body, wasm.OpcodeI32Load16U)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I64_LOAD8_S:
		ins := i.(ast.Ins_I64Load8S)
		dst.Body = append(dst.Body, wasm.OpcodeI64Load8S)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I64_LOAD8_U:
		ins := i.(ast.Ins_I64Load8U)
		dst.Body = append(dst.Body, wasm.OpcodeI64Load8U)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I64_LOAD16_S:
		ins := i.(ast.Ins_I64Load16S)
		dst.Body = append(dst.Body, wasm.OpcodeI64Load16S)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I64_LOAD16_U:
		ins := i.(ast.Ins_I64Load16U)
		dst.Body = append(dst.Body, wasm.OpcodeI64Load16U)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I64_LOAD32_S:
		ins := i.(ast.Ins_I64Load32S)
		dst.Body = append(dst.Body, wasm.OpcodeI64Load32S)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I64_LOAD32_U:
		ins := i.(ast.Ins_I64Load32U)
		dst.Body = append(dst.Body, wasm.OpcodeI64Load32U)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I32_STORE:
		ins := i.(ast.Ins_I32Store)
		dst.Body = append(dst.Body, wasm.OpcodeI32Store)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I64_STORE:
		ins := i.(ast.Ins_I64Store)
		dst.Body = append(dst.Body, wasm.OpcodeI64Store)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_F32_STORE:
		ins := i.(ast.Ins_F32Store)
		dst.Body = append(dst.Body, wasm.OpcodeF32Store)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_F64_STORE:
		ins := i.(ast.Ins_F64Store)
		dst.Body = append(dst.Body, wasm.OpcodeF64Store)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I32_STORE8:
		ins := i.(ast.Ins_I32Store8)
		dst.Body = append(dst.Body, wasm.OpcodeI32Store8)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I32_STORE16:
		ins := i.(ast.Ins_I32Store16)
		dst.Body = append(dst.Body, wasm.OpcodeI32Store16)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I64_STORE8:
		ins := i.(ast.Ins_I64Store8)
		dst.Body = append(dst.Body, wasm.OpcodeI64Store8)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I64_STORE16:
		ins := i.(ast.Ins_I64Store16)
		dst.Body = append(dst.Body, wasm.OpcodeI64Store16)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_I64_STORE32:
		ins := i.(ast.Ins_I64Store32)
		dst.Body = append(dst.Body, wasm.OpcodeI64Store32)
		dst.Body = append(dst.Body, p.encodeAlign(ins.Align))
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.Offset))...)

	case token.INS_MEMORY_SIZE:
		dst.Body = append(dst.Body, wasm.OpcodeMemorySize, 0x00)
	case token.INS_MEMORY_GROW:
		dst.Body = append(dst.Body, wasm.OpcodeMemoryGrow, 0x00)
	case token.INS_MEMORY_INIT:
		ins := i.(ast.Ins_MemoryInit)
		dst.Body = append(dst.Body, 0xfc, wasm.OpcodeMiscMemoryInit)
		dst.Body = append(dst.Body, p.encodeUint32(uint32(ins.DataIdx))...)
		dst.Body = append(dst.Body, 0x00)
	case token.INS_MEMORY_COPY:
		dst.Body = append(dst.Body, 0xfc, wasm.OpcodeMiscMemoryCopy, 0x00, 0x00)
	case token.INS_MEMORY_FILL:
		dst.Body = append(dst.Body, 0xfc, wasm.OpcodeMiscMemoryFill, 0x00)
	case token.INS_I32_CONST:
		ins := i.(ast.Ins_I32Const)
		dst.Body = append(dst.Body, wasm.OpcodeI32Const)
		dst.Body = append(dst.Body, p.encodeInt32(ins.X)...)
	case token.INS_I64_CONST:
		ins := i.(ast.Ins_I64Const)
		dst.Body = append(dst.Body, wasm.OpcodeI64Const)
		dst.Body = append(dst.Body, p.encodeInt64(ins.X)...)
	case token.INS_F32_CONST:
		ins := i.(ast.Ins_F32Const)
		dst.Body = append(dst.Body, wasm.OpcodeF32Const)
		dst.Body = append(dst.Body, p.encodeFloat32(ins.X)...)
	case token.INS_F64_CONST:
		ins := i.(ast.Ins_F64Const)
		dst.Body = append(dst.Body, wasm.OpcodeF64Const)
		dst.Body = append(dst.Body, p.encodeFloat64(ins.X)...)

	case token.INS_I32_EQZ:
		dst.Body = append(dst.Body, wasm.OpcodeI32Eqz)
	case token.INS_I32_EQ:
		dst.Body = append(dst.Body, wasm.OpcodeI32Eq)
	case token.INS_I32_NE:
		dst.Body = append(dst.Body, wasm.OpcodeI32Ne)
	case token.INS_I32_LT_S:
		dst.Body = append(dst.Body, wasm.OpcodeI32LtS)
	case token.INS_I32_LT_U:
		dst.Body = append(dst.Body, wasm.OpcodeI32LtU)
	case token.INS_I32_GT_S:
		dst.Body = append(dst.Body, wasm.OpcodeI32GtS)
	case token.INS_I32_GT_U:
		dst.Body = append(dst.Body, wasm.OpcodeI32GtU)
	case token.INS_I32_LE_S:
		dst.Body = append(dst.Body, wasm.OpcodeI32LeS)
	case token.INS_I32_LE_U:
		dst.Body = append(dst.Body, wasm.OpcodeI32LeU)
	case token.INS_I32_GE_S:
		dst.Body = append(dst.Body, wasm.OpcodeI32GeS)
	case token.INS_I32_GE_U:
		dst.Body = append(dst.Body, wasm.OpcodeI32GeU)
	case token.INS_I64_EQZ:
		dst.Body = append(dst.Body, wasm.OpcodeI64Eqz)
	case token.INS_I64_EQ:
		dst.Body = append(dst.Body, wasm.OpcodeI64Eq)
	case token.INS_I64_NE:
		dst.Body = append(dst.Body, wasm.OpcodeI64Ne)
	case token.INS_I64_LT_S:
		dst.Body = append(dst.Body, wasm.OpcodeI64LtS)
	case token.INS_I64_LT_U:
		dst.Body = append(dst.Body, wasm.OpcodeI64LtU)
	case token.INS_I64_GT_S:
		dst.Body = append(dst.Body, wasm.OpcodeI64GtS)
	case token.INS_I64_GT_U:
		dst.Body = append(dst.Body, wasm.OpcodeI64GtU)
	case token.INS_I64_LE_S:
		dst.Body = append(dst.Body, wasm.OpcodeI64LeS)
	case token.INS_I64_LE_U:
		dst.Body = append(dst.Body, wasm.OpcodeI64LeU)
	case token.INS_I64_GE_S:
		dst.Body = append(dst.Body, wasm.OpcodeI64GeS)
	case token.INS_I64_GE_U:
		dst.Body = append(dst.Body, wasm.OpcodeI64GeU)
	case token.INS_F32_EQ:
		dst.Body = append(dst.Body, wasm.OpcodeF32Eq)
	case token.INS_F32_NE:
		dst.Body = append(dst.Body, wasm.OpcodeF32Ne)
	case token.INS_F32_LT:
		dst.Body = append(dst.Body, wasm.OpcodeF32Lt)
	case token.INS_F32_GT:
		dst.Body = append(dst.Body, wasm.OpcodeF32Gt)
	case token.INS_F32_LE:
		dst.Body = append(dst.Body, wasm.OpcodeF32Le)
	case token.INS_F32_GE:
		dst.Body = append(dst.Body, wasm.OpcodeF32Ge)
	case token.INS_F64_EQ:
		dst.Body = append(dst.Body, wasm.OpcodeF64Eq)
	case token.INS_F64_NE:
		dst.Body = append(dst.Body, wasm.OpcodeF64Ne)
	case token.INS_F64_LT:
		dst.Body = append(dst.Body, wasm.OpcodeF64Lt)
	case token.INS_F64_GT:
		dst.Body = append(dst.Body, wasm.OpcodeF64Gt)
	case token.INS_F64_LE:
		dst.Body = append(dst.Body, wasm.OpcodeF64Le)
	case token.INS_F64_GE:
		dst.Body = append(dst.Body, wasm.OpcodeF64Ge)
	case token.INS_I32_CLZ:
		dst.Body = append(dst.Body, wasm.OpcodeI32Clz)
	case token.INS_I32_CTZ:
		dst.Body = append(dst.Body, wasm.OpcodeI32Ctz)
	case token.INS_I32_POPCNT:
		dst.Body = append(dst.Body, wasm.OpcodeI32Popcnt)
	case token.INS_I32_ADD:
		dst.Body = append(dst.Body, wasm.OpcodeI32Add)
	case token.INS_I32_SUB:
		dst.Body = append(dst.Body, wasm.OpcodeI32Sub)
	case token.INS_I32_MUL:
		dst.Body = append(dst.Body, wasm.OpcodeI32Mul)
	case token.INS_I32_DIV_S:
		dst.Body = append(dst.Body, wasm.OpcodeI32DivS)
	case token.INS_I32_DIV_U:
		dst.Body = append(dst.Body, wasm.OpcodeI32DivU)
	case token.INS_I32_REM_S:
		dst.Body = append(dst.Body, wasm.OpcodeI32RemS)
	case token.INS_I32_REM_U:
		dst.Body = append(dst.Body, wasm.OpcodeI32RemU)
	case token.INS_I32_AND:
		dst.Body = append(dst.Body, wasm.OpcodeI32And)
	case token.INS_I32_OR:
		dst.Body = append(dst.Body, wasm.OpcodeI32Or)
	case token.INS_I32_XOR:
		dst.Body = append(dst.Body, wasm.OpcodeI32Xor)
	case token.INS_I32_SHL:
		dst.Body = append(dst.Body, wasm.OpcodeI32Shl)
	case token.INS_I32_SHR_S:
		dst.Body = append(dst.Body, wasm.OpcodeI32ShrS)
	case token.INS_I32_SHR_U:
		dst.Body = append(dst.Body, wasm.OpcodeI32ShrU)
	case token.INS_I32_ROTL:
		dst.Body = append(dst.Body, wasm.OpcodeI32Rotl)
	case token.INS_I32_ROTR:
		dst.Body = append(dst.Body, wasm.OpcodeI32Rotr)
	case token.INS_I64_CLZ:
		dst.Body = append(dst.Body, wasm.OpcodeI64Clz)
	case token.INS_I64_CTZ:
		dst.Body = append(dst.Body, wasm.OpcodeI64Ctz)
	case token.INS_I64_POPCNT:
		dst.Body = append(dst.Body, wasm.OpcodeI64Popcnt)
	case token.INS_I64_ADD:
		dst.Body = append(dst.Body, wasm.OpcodeI64Add)
	case token.INS_I64_SUB:
		dst.Body = append(dst.Body, wasm.OpcodeI64Sub)
	case token.INS_I64_MUL:
		dst.Body = append(dst.Body, wasm.OpcodeI64Mul)
	case token.INS_I64_DIV_S:
		dst.Body = append(dst.Body, wasm.OpcodeI64DivS)
	case token.INS_I64_DIV_U:
		dst.Body = append(dst.Body, wasm.OpcodeI64DivU)
	case token.INS_I64_REM_S:
		dst.Body = append(dst.Body, wasm.OpcodeI64RemS)
	case token.INS_I64_REM_U:
		dst.Body = append(dst.Body, wasm.OpcodeI64RemU)
	case token.INS_I64_AND:
		dst.Body = append(dst.Body, wasm.OpcodeI64And)
	case token.INS_I64_OR:
		dst.Body = append(dst.Body, wasm.OpcodeI64Or)
	case token.INS_I64_XOR:
		dst.Body = append(dst.Body, wasm.OpcodeI64Xor)
	case token.INS_I64_SHL:
		dst.Body = append(dst.Body, wasm.OpcodeI64Shl)
	case token.INS_I64_SHR_S:
		dst.Body = append(dst.Body, wasm.OpcodeI64ShrS)
	case token.INS_I64_SHR_U:
		dst.Body = append(dst.Body, wasm.OpcodeI64ShrU)
	case token.INS_I64_ROTL:
		dst.Body = append(dst.Body, wasm.OpcodeI64Rotl)
	case token.INS_I64_ROTR:
		dst.Body = append(dst.Body, wasm.OpcodeI64Rotr)
	case token.INS_F32_ABS:
		dst.Body = append(dst.Body, wasm.OpcodeF32Abs)
	case token.INS_F32_NEG:
		dst.Body = append(dst.Body, wasm.OpcodeF32Neg)
	case token.INS_F32_CEIL:
		dst.Body = append(dst.Body, wasm.OpcodeF32Ceil)
	case token.INS_F32_FLOOR:
		dst.Body = append(dst.Body, wasm.OpcodeF32Floor)
	case token.INS_F32_TRUNC:
		dst.Body = append(dst.Body, wasm.OpcodeF32Trunc)
	case token.INS_F32_NEAREST:
		dst.Body = append(dst.Body, wasm.OpcodeF32Nearest)
	case token.INS_F32_SQRT:
		dst.Body = append(dst.Body, wasm.OpcodeF32Sqrt)
	case token.INS_F32_ADD:
		dst.Body = append(dst.Body, wasm.OpcodeF32Add)
	case token.INS_F32_SUB:
		dst.Body = append(dst.Body, wasm.OpcodeF32Sub)
	case token.INS_F32_MUL:
		dst.Body = append(dst.Body, wasm.OpcodeF32Mul)
	case token.INS_F32_DIV:
		dst.Body = append(dst.Body, wasm.OpcodeF32Div)
	case token.INS_F32_MIN:
		dst.Body = append(dst.Body, wasm.OpcodeF32Min)
	case token.INS_F32_MAX:
		dst.Body = append(dst.Body, wasm.OpcodeF32Max)
	case token.INS_F32_COPYSIGN:
		dst.Body = append(dst.Body, wasm.OpcodeF32Copysign)
	case token.INS_F64_ABS:
		dst.Body = append(dst.Body, wasm.OpcodeF64Abs)
	case token.INS_F64_NEG:
		dst.Body = append(dst.Body, wasm.OpcodeF64Neg)
	case token.INS_F64_CEIL:
		dst.Body = append(dst.Body, wasm.OpcodeF64Ceil)
	case token.INS_F64_FLOOR:
		dst.Body = append(dst.Body, wasm.OpcodeF64Floor)
	case token.INS_F64_TRUNC:
		dst.Body = append(dst.Body, wasm.OpcodeF64Trunc)
	case token.INS_F64_NEAREST:
		dst.Body = append(dst.Body, wasm.OpcodeF64Nearest)
	case token.INS_F64_SQRT:
		dst.Body = append(dst.Body, wasm.OpcodeF64Sqrt)
	case token.INS_F64_ADD:
		dst.Body = append(dst.Body, wasm.OpcodeF64Add)
	case token.INS_F64_SUB:
		dst.Body = append(dst.Body, wasm.OpcodeF64Sub)
	case token.INS_F64_MUL:
		dst.Body = append(dst.Body, wasm.OpcodeF64Mul)
	case token.INS_F64_DIV:
		dst.Body = append(dst.Body, wasm.OpcodeF64Div)
	case token.INS_F64_MIN:
		dst.Body = append(dst.Body, wasm.OpcodeF64Min)
	case token.INS_F64_MAX:
		dst.Body = append(dst.Body, wasm.OpcodeF64Max)
	case token.INS_F64_COPYSIGN:
		dst.Body = append(dst.Body, wasm.OpcodeF64Copysign)
	case token.INS_I32_WRAP_I64:
		dst.Body = append(dst.Body, wasm.OpcodeI32WrapI64)
	case token.INS_I32_TRUNC_F32_S:
		dst.Body = append(dst.Body, wasm.OpcodeI32TruncF32S)
	case token.INS_I32_TRUNC_F32_U:
		dst.Body = append(dst.Body, wasm.OpcodeI32TruncF32U)
	case token.INS_I32_TRUNC_F64_S:
		dst.Body = append(dst.Body, wasm.OpcodeI32TruncF64S)
	case token.INS_I32_TRUNC_F64_U:
		dst.Body = append(dst.Body, wasm.OpcodeI32TruncF64U)
	case token.INS_I64_EXTEND_I32_S:
		dst.Body = append(dst.Body, wasm.OpcodeI64ExtendI32S)
	case token.INS_I64_EXTEND_I32_U:
		dst.Body = append(dst.Body, wasm.OpcodeI64ExtendI32U)
	case token.INS_I64_TRUNC_F32_S:
		dst.Body = append(dst.Body, wasm.OpcodeI64TruncF32S)
	case token.INS_I64_TRUNC_F32_U:
		dst.Body = append(dst.Body, wasm.OpcodeI64TruncF32U)
	case token.INS_I64_TRUNC_F64_S:
		dst.Body = append(dst.Body, wasm.OpcodeI64TruncF64S)
	case token.INS_I64_TRUNC_F64_U:
		dst.Body = append(dst.Body, wasm.OpcodeI64TruncF64U)
	case token.INS_F32_CONVERT_I32_S:
		dst.Body = append(dst.Body, wasm.OpcodeF32ConvertI32S)
	case token.INS_F32_CONVERT_I32_U:
		dst.Body = append(dst.Body, wasm.OpcodeF32ConvertI32U)
	case token.INS_F32_CONVERT_I64_S:
		dst.Body = append(dst.Body, wasm.OpcodeF32ConvertI64S)
	case token.INS_F32_CONVERT_I64_U:
		dst.Body = append(dst.Body, wasm.OpcodeF32ConvertI64U)
	case token.INS_F32_DEMOTE_F64:
		dst.Body = append(dst.Body, wasm.OpcodeF32DemoteF64)
	case token.INS_F64_CONVERT_I32_S:
		dst.Body = append(dst.Body, wasm.OpcodeF64ConvertI32S)
	case token.INS_F64_CONVERT_I32_U:
		dst.Body = append(dst.Body, wasm.OpcodeF64ConvertI32U)
	case token.INS_F64_CONVERT_I64_S:
		dst.Body = append(dst.Body, wasm.OpcodeF64ConvertI64S)
	case token.INS_F64_CONVERT_I64_U:
		dst.Body = append(dst.Body, wasm.OpcodeF64ConvertI64U)
	case token.INS_F64_PROMOTE_F32:
		dst.Body = append(dst.Body, wasm.OpcodeF64PromoteF32)
	case token.INS_I32_REINTERPRET_F32:
		dst.Body = append(dst.Body, wasm.OpcodeI32ReinterpretF32)
	case token.INS_I64_REINTERPRET_F64:
		dst.Body = append(dst.Body, wasm.OpcodeI64ReinterpretF64)
	case token.INS_F32_REINTERPRET_I32:
		dst.Body = append(dst.Body, wasm.OpcodeF32ReinterpretI32)
	case token.INS_F64_REINTERPRET_I64:
		dst.Body = append(dst.Body, wasm.OpcodeF64ReinterpretI64)
	default:
		panic("unreachable")
	}
}
