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

func (p *wat2wasmWorker) lookupTokenOpcode(tok token.Token) wasm.Opcode {
	return tokOpcodeMap[tok]
}

func (p *wat2wasmWorker) buildInstruction_enterBlock() {}
func (p *wat2wasmWorker) buildInstruction_leaveBlock() {}

func (p *wat2wasmWorker) buildInstruction(dst []byte, i ast.Instruction) []byte {
	tok := i.Token()
	opcode := p.lookupTokenOpcode(tok)

	switch tok {
	case token.INS_UNREACHABLE:
		return append(dst, opcode)
	case token.INS_NOP:
		return append(dst, opcode)
	case token.INS_BLOCK:
		blk := i.(ast.Ins_Block)
		_ = blk
		return append(dst, opcode)
	case token.INS_LOOP:
		return append(dst, opcode)
	case token.INS_IF:
		return append(dst, opcode)
	case token.INS_ELSE:
		return append(dst, opcode)
	case token.INS_END:
		panic("unreachable")
	case token.INS_BR:
		return append(dst, opcode)
	case token.INS_BR_IF:
		return append(dst, opcode)
	case token.INS_BR_TABLE:
		return append(dst, opcode)
	case token.INS_RETURN:
		return append(dst, opcode)
	case token.INS_CALL:
		return append(dst, opcode)
	case token.INS_CALL_INDIRECT:
		return append(dst, opcode)
	case token.INS_DROP:
		return append(dst, opcode)
	case token.INS_SELECT:
		return append(dst, opcode)
	case token.INS_TYPED_SELECT:
		return append(dst, opcode)
	case token.INS_LOCAL_GET:
		return append(dst, opcode)
	case token.INS_LOCAL_SET:
		return append(dst, opcode)
	case token.INS_LOCAL_TEE:
		return append(dst, opcode)
	case token.INS_GLOBAL_GET:
		return append(dst, opcode)
	case token.INS_GLOBAL_SET:
		return append(dst, opcode)
	case token.INS_TABLE_GET:
		return append(dst, opcode)
	case token.INS_TABLE_SET:
		return append(dst, opcode)
	case token.INS_I32_LOAD:
		return append(dst, opcode)
	case token.INS_I64_LOAD:
		return append(dst, opcode)
	case token.INS_F32_LOAD:
		return append(dst, opcode)
	case token.INS_F64_LOAD:
		return append(dst, opcode)
	case token.INS_I32_LOAD8_S:
		return append(dst, opcode)
	case token.INS_I32_LOAD8_U:
		return append(dst, opcode)
	case token.INS_I32_LOAD16_S:
		return append(dst, opcode)
	case token.INS_I32_LOAD16_U:
		return append(dst, opcode)
	case token.INS_I64_LOAD8_S:
		return append(dst, opcode)
	case token.INS_I64_LOAD8_U:
		return append(dst, opcode)
	case token.INS_I64_LOAD16_S:
		return append(dst, opcode)
	case token.INS_I64_LOAD16_U:
		return append(dst, opcode)
	case token.INS_I64_LOAD32_S:
		return append(dst, opcode)
	case token.INS_I64_LOAD32_U:
		return append(dst, opcode)
	case token.INS_I32_STORE:
		return append(dst, opcode)
	case token.INS_I64_STORE:
		return append(dst, opcode)
	case token.INS_F32_STORE:
		return append(dst, opcode)
	case token.INS_F64_STORE:
		return append(dst, opcode)
	case token.INS_I32_STORE8:
		return append(dst, opcode)
	case token.INS_I32_STORE16:
		return append(dst, opcode)
	case token.INS_I64_STORE8:
		return append(dst, opcode)
	case token.INS_I64_STORE16:
		return append(dst, opcode)
	case token.INS_I64_STORE32:
		return append(dst, opcode)
	case token.INS_MEMORY_SIZE:
		return append(dst, opcode)
	case token.INS_MEMORY_GROW:
		return append(dst, opcode)
	case token.INS_I32_CONST:
		return append(dst, opcode)
	case token.INS_I64_CONST:
		return append(dst, opcode)
	case token.INS_F32_CONST:
		return append(dst, opcode)
	case token.INS_F64_CONST:
		return append(dst, opcode)
	case token.INS_I32_EQZ:
		return append(dst, opcode)
	case token.INS_I32_EQ:
		return append(dst, opcode)
	case token.INS_I32_NE:
		return append(dst, opcode)
	case token.INS_I32_LT_S:
		return append(dst, opcode)
	case token.INS_I32_LT_U:
		return append(dst, opcode)
	case token.INS_I32_GT_S:
		return append(dst, opcode)
	case token.INS_I32_GT_U:
		return append(dst, opcode)
	case token.INS_I32_LE_S:
		return append(dst, opcode)
	case token.INS_I32_LE_U:
		return append(dst, opcode)
	case token.INS_I32_GE_S:
		return append(dst, opcode)
	case token.INS_I32_GE_U:
		return append(dst, opcode)
	case token.INS_I64_EQZ:
		return append(dst, opcode)
	case token.INS_I64_EQ:
		return append(dst, opcode)
	case token.INS_I64_NE:
		return append(dst, opcode)
	case token.INS_I64_LT_S:
		return append(dst, opcode)
	case token.INS_I64_LT_U:
		return append(dst, opcode)
	case token.INS_I64_GT_S:
		return append(dst, opcode)
	case token.INS_I64_GT_U:
		return append(dst, opcode)
	case token.INS_I64_LE_S:
		return append(dst, opcode)
	case token.INS_I64_LE_U:
		return append(dst, opcode)
	case token.INS_I64_GE_S:
		return append(dst, opcode)
	case token.INS_I64_GE_U:
		return append(dst, opcode)
	case token.INS_F32_EQ:
		return append(dst, opcode)
	case token.INS_F32_NE:
		return append(dst, opcode)
	case token.INS_F32_LT:
		return append(dst, opcode)
	case token.INS_F32_GT:
		return append(dst, opcode)
	case token.INS_F32_LE:
		return append(dst, opcode)
	case token.INS_F32_GE:
		return append(dst, opcode)
	case token.INS_F64_EQ:
		return append(dst, opcode)
	case token.INS_F64_NE:
		return append(dst, opcode)
	case token.INS_F64_LT:
		return append(dst, opcode)
	case token.INS_F64_GT:
		return append(dst, opcode)
	case token.INS_F64_LE:
		return append(dst, opcode)
	case token.INS_F64_GE:
		return append(dst, opcode)
	case token.INS_I32_CLZ:
		return append(dst, opcode)
	case token.INS_I32_CTZ:
		return append(dst, opcode)
	case token.INS_I32_POPCNT:
		return append(dst, opcode)
	case token.INS_I32_ADD:
		return append(dst, opcode)
	case token.INS_I32_SUB:
		return append(dst, opcode)
	case token.INS_I32_MUL:
		return append(dst, opcode)
	case token.INS_I32_DIV_S:
		return append(dst, opcode)
	case token.INS_I32_DIV_U:
		return append(dst, opcode)
	case token.INS_I32_REM_S:
		return append(dst, opcode)
	case token.INS_I32_REM_U:
		return append(dst, opcode)
	case token.INS_I32_AND:
		return append(dst, opcode)
	case token.INS_I32_OR:
		return append(dst, opcode)
	case token.INS_I32_XOR:
		return append(dst, opcode)
	case token.INS_I32_SHL:
		return append(dst, opcode)
	case token.INS_I32_SHR_S:
		return append(dst, opcode)
	case token.INS_I32_SHR_U:
		return append(dst, opcode)
	case token.INS_I32_ROTL:
		return append(dst, opcode)
	case token.INS_I32_ROTR:
		return append(dst, opcode)
	case token.INS_I64_CLZ:
		return append(dst, opcode)
	case token.INS_I64_CTZ:
		return append(dst, opcode)
	case token.INS_I64_POPCNT:
		return append(dst, opcode)
	case token.INS_I64_ADD:
		return append(dst, opcode)
	case token.INS_I64_SUB:
		return append(dst, opcode)
	case token.INS_I64_MUL:
		return append(dst, opcode)
	case token.INS_I64_DIV_S:
		return append(dst, opcode)
	case token.INS_I64_DIV_U:
		return append(dst, opcode)
	case token.INS_I64_REM_S:
		return append(dst, opcode)
	case token.INS_I64_REM_U:
		return append(dst, opcode)
	case token.INS_I64_AND:
		return append(dst, opcode)
	case token.INS_I64_OR:
		return append(dst, opcode)
	case token.INS_I64_XOR:
		return append(dst, opcode)
	case token.INS_I64_SHL:
		return append(dst, opcode)
	case token.INS_I64_SHR_S:
		return append(dst, opcode)
	case token.INS_I64_SHR_U:
		return append(dst, opcode)
	case token.INS_I64_ROTL:
		return append(dst, opcode)
	case token.INS_I64_ROTR:
		return append(dst, opcode)
	case token.INS_F32_ABS:
		return append(dst, opcode)
	case token.INS_F32_NEG:
		return append(dst, opcode)
	case token.INS_F32_CEIL:
		return append(dst, opcode)
	case token.INS_F32_FLOOR:
		return append(dst, opcode)
	case token.INS_F32_TRUNC:
		return append(dst, opcode)
	case token.INS_F32_NEAREST:
		return append(dst, opcode)
	case token.INS_F32_SQRT:
		return append(dst, opcode)
	case token.INS_F32_ADD:
		return append(dst, opcode)
	case token.INS_F32_SUB:
		return append(dst, opcode)
	case token.INS_F32_MUL:
		return append(dst, opcode)
	case token.INS_F32_DIV:
		return append(dst, opcode)
	case token.INS_F32_MIN:
		return append(dst, opcode)
	case token.INS_F32_MAX:
		return append(dst, opcode)
	case token.INS_F32_COPYSIGN:
		return append(dst, opcode)
	case token.INS_F64_ABS:
		return append(dst, opcode)
	case token.INS_F64_NEG:
		return append(dst, opcode)
	case token.INS_F64_CEIL:
		return append(dst, opcode)
	case token.INS_F64_FLOOR:
		return append(dst, opcode)
	case token.INS_F64_TRUNC:
		return append(dst, opcode)
	case token.INS_F64_NEAREST:
		return append(dst, opcode)
	case token.INS_F64_SQRT:
		return append(dst, opcode)
	case token.INS_F64_ADD:
		return append(dst, opcode)
	case token.INS_F64_SUB:
		return append(dst, opcode)
	case token.INS_F64_MUL:
		return append(dst, opcode)
	case token.INS_F64_DIV:
		return append(dst, opcode)
	case token.INS_F64_MIN:
		return append(dst, opcode)
	case token.INS_F64_MAX:
		return append(dst, opcode)
	case token.INS_F64_COPYSIGN:
		return append(dst, opcode)
	case token.INS_I32_WRAP_I64:
		return append(dst, opcode)
	case token.INS_I32_TRUNC_F32_S:
		return append(dst, opcode)
	case token.INS_I32_TRUNC_F32_U:
		return append(dst, opcode)
	case token.INS_I32_TRUNC_F64_S:
		return append(dst, opcode)
	case token.INS_I32_TRUNC_F64_U:
		return append(dst, opcode)
	case token.INS_I64_EXTEND_I32_S:
		return append(dst, opcode)
	case token.INS_I64_EXTEND_I32_U:
		return append(dst, opcode)
	case token.INS_I64_TRUNC_F32_S:
		return append(dst, opcode)
	case token.INS_I64_TRUNC_F32_U:
		return append(dst, opcode)
	case token.INS_I64_TRUNC_F64_S:
		return append(dst, opcode)
	case token.INS_I64_TRUNC_F64_U:
		return append(dst, opcode)
	case token.INS_F32_CONVERT_I32_S:
		return append(dst, opcode)
	case token.INS_F32_CONVERT_I32_U:
		return append(dst, opcode)
	case token.INS_F32_CONVERT_I64_S:
		return append(dst, opcode)
	case token.INS_F32_CONVERT_I64_U:
		return append(dst, opcode)
	case token.INS_F32_DEMOTE_F64:
		return append(dst, opcode)
	case token.INS_F64_CONVERT_I32_S:
		return append(dst, opcode)
	case token.INS_F64_CONVERT_I32_U:
		return append(dst, opcode)
	case token.INS_F64_CONVERT_I64_S:
		return append(dst, opcode)
	case token.INS_F64_CONVERT_I64_U:
		return append(dst, opcode)
	case token.INS_F64_DEMOTE_F32:
		return append(dst, opcode)
	case token.INS_I32_REINTERPRET_F32:
		return append(dst, opcode)
	case token.INS_I64_REINTERPRET_F64:
		return append(dst, opcode)
	case token.INS_I32_REINTERPRET_I32:
		return append(dst, opcode)
	case token.INS_I64_REINTERPRET_I64:
		return append(dst, opcode)
	}
	panic("unreachable")
}
