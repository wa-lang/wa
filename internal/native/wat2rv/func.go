// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2rv

import (
	"fmt"

	watast "wa-lang.org/wa/internal/wat/ast"
	wattoken "wa-lang.org/wa/internal/wat/token"
)

func (p *wat2rvWorker) buildFunc(fn *watast.Func) error {
	p.buildFunc_args(fn)
	p.buildFunc_local(fn)
	p.buildFunc_body(fn)
	return nil
}

func (p *wat2rvWorker) buildFunc_args(fn *watast.Func) error {
	return nil
}

func (p *wat2rvWorker) buildFunc_local(fn *watast.Func) error {
	return nil
}

func (p *wat2rvWorker) buildFunc_body(fn *watast.Func) error {
	for _, ins := range fn.Body.Insts {
		if err := p.buildFunc_ins(fn, ins, 1); err != nil {
			return err
		}
	}
	return nil
}

func (p *wat2rvWorker) buildFunc_ins(fn *watast.Func, i watast.Instruction, level int) error {
	switch tok := i.Token(); tok {
	default:
		panic(fmt.Sprintf("unreachable: %v", tok))
	case wattoken.INS_UNREACHABLE:
	case wattoken.INS_LOCAL_GET:
	case wattoken.INS_I64_CONST:
	case wattoken.INS_I64_LE_U:
	case wattoken.INS_I64_SUB:
	case wattoken.INS_I64_ADD:
	case wattoken.INS_IF:
	case wattoken.INS_CALL:
	}
	switch tok := i.Token(); tok {
	default:
		panic(fmt.Sprintf("unreachable: %v", tok))

	case wattoken.INS_UNREACHABLE: // 0x00, unreachable
		// TODO
	case wattoken.INS_NOP: // 0x01, nop
	case wattoken.INS_BLOCK: // 0x02, block
	case wattoken.INS_LOOP: // 0x03, loop
	case wattoken.INS_IF: // 0x04, if
		// TODO
	case wattoken.INS_ELSE: // 0x05, else
	case wattoken.INS_END: // 0x0b, end
	case wattoken.INS_BR: // 0x0c, br
	case wattoken.INS_BR_IF: // 0x0d, br_if
	case wattoken.INS_BR_TABLE: // 0x0e, br_table
	case wattoken.INS_RETURN: // 0x0f, return
	case wattoken.INS_CALL: // 0x10, call
		// TODO
	case wattoken.INS_CALL_INDIRECT: // 0x11, call_indirect
	case wattoken.INS_DROP: // 0x1a, drop
	case wattoken.INS_SELECT: // 0x1b, select
	case wattoken.INS_LOCAL_GET: // 0x20, local.get
		// TODO
	case wattoken.INS_LOCAL_SET: // 0x21, local.set
	case wattoken.INS_LOCAL_TEE: // 0x22, local.tee
	case wattoken.INS_GLOBAL_GET: // 0x23, global.get
	case wattoken.INS_GLOBAL_SET: // 0x24, global.set
	case wattoken.INS_TABLE_GET: // 0x25, table.get
	case wattoken.INS_TABLE_SET: // 0x26, table.set
	case wattoken.INS_I32_LOAD: // 0x28, i32.load
	case wattoken.INS_I64_LOAD: // 0x29, i64.load
	case wattoken.INS_F32_LOAD: // 0x2a, f32.load
	case wattoken.INS_F64_LOAD: // 0x2b, f64.load
	case wattoken.INS_I32_LOAD8_S: // 0x2c, i32.load8_s
	case wattoken.INS_I32_LOAD8_U: // 0x2d, i32.load8_u
	case wattoken.INS_I32_LOAD16_S: // 0x2e, i32.load16_s
	case wattoken.INS_I32_LOAD16_U: // 0x2f, i32.load16_u
	case wattoken.INS_I64_LOAD8_S: // 0x30, i64.load8_s
	case wattoken.INS_I64_LOAD8_U: // 0x31, i64.load8_u
	case wattoken.INS_I64_LOAD16_S: // 0x31, i64.load16_s
	case wattoken.INS_I64_LOAD16_U: // 0x33, i64.load16_u
	case wattoken.INS_I64_LOAD32_S: // 0x34, i64.load32_s
	case wattoken.INS_I64_LOAD32_U: // 0x35, i64.load32_u
	case wattoken.INS_I32_STORE: // 0x36, i32.store
	case wattoken.INS_I64_STORE: // 0x37, i64.store
	case wattoken.INS_F32_STORE: // 0x38, f32.store
	case wattoken.INS_F64_STORE: // 0x39, f64.store
	case wattoken.INS_I32_STORE8: // 0x3a, i32.store8
	case wattoken.INS_I32_STORE16: // 0x3b, i32.store16
	case wattoken.INS_I64_STORE8: // 0x3c, i64.store8
	case wattoken.INS_I64_STORE16: // 0x3d, i64.store16
	case wattoken.INS_I64_STORE32: // 0x3e, i64.store32
	case wattoken.INS_MEMORY_SIZE: // 0x3f, memory.size
	case wattoken.INS_MEMORY_GROW: // 0x40, memory.grow
	case wattoken.INS_MEMORY_INIT: // 0xfc 0x08, memory.init
	case wattoken.INS_MEMORY_COPY: // 0xfc 0x0a, memory.copy
	case wattoken.INS_MEMORY_FILL: // 0xfc 0x0b, memory.fill
	case wattoken.INS_I32_CONST: // 0x41, i32.const
	case wattoken.INS_I64_CONST: // 0x42, i64.const
		// TODO
	case wattoken.INS_F32_CONST: // 0x43, f32.const
	case wattoken.INS_F64_CONST: // 0x44, f64.const
	case wattoken.INS_I32_EQZ: // 0x45, i32.eqz
	case wattoken.INS_I32_EQ: // 0x46, i32.eq
	case wattoken.INS_I32_NE: // 0x47, i32.ne
	case wattoken.INS_I32_LT_S: // 0x48, i32.lt_s
	case wattoken.INS_I32_LT_U: // 0x49, i32.lt_u
	case wattoken.INS_I32_GT_S: // 0x4a, i32.gt_s
	case wattoken.INS_I32_GT_U: // 0x4b, i32.gt_u
	case wattoken.INS_I32_LE_S: // 0x4c, i32.le_s
	case wattoken.INS_I32_LE_U: // 0x4d, i32.le_u
	case wattoken.INS_I32_GE_S: // 0x4e, i32.ge_s
	case wattoken.INS_I32_GE_U: // 0x4f, i32.ge_u
	case wattoken.INS_I64_EQZ: // 0x50, i64.eqz
	case wattoken.INS_I64_EQ: // 0x51, i64.eq
	case wattoken.INS_I64_NE: // 0x52, i64.ne
	case wattoken.INS_I64_LT_S: // 0x53, i64.lt_s
	case wattoken.INS_I64_LT_U: // 0x54, i64.lt_u
	case wattoken.INS_I64_GT_S: // 0x55, i64.gt_s
	case wattoken.INS_I64_GT_U: // 0x56, i64.gt_u
	case wattoken.INS_I64_LE_S: // 0x57, i64.le_s
	case wattoken.INS_I64_LE_U: // 0x58, i64.le_u
		// TODO
	case wattoken.INS_I64_GE_S: // 0x59, i64.ge_s
	case wattoken.INS_I64_GE_U: // 0x5a, i64.ge_u
	case wattoken.INS_F32_EQ: // 0x5b, f32.eq
	case wattoken.INS_F32_NE: // 0x5c, f32.ne
	case wattoken.INS_F32_LT: // 0x5d, f32.lt
	case wattoken.INS_F32_GT: // 0x5e, f32.gt
	case wattoken.INS_F32_LE: // 0x5f, f32.le
	case wattoken.INS_F32_GE: // 0x60, f32.ge
	case wattoken.INS_F64_EQ: // 0x61, f64.eq
	case wattoken.INS_F64_NE: // 0x62, f64.ne
	case wattoken.INS_F64_LT: // 0x63, f64.lt
	case wattoken.INS_F64_GT: // 0x64, f64.gt
	case wattoken.INS_F64_LE: // 0x65, f64.le
	case wattoken.INS_F64_GE: // 0x66, f64.ge
	case wattoken.INS_I32_CLZ: // 0x67, i32.clz
	case wattoken.INS_I32_CTZ: // 0x68, i32.ctz
	case wattoken.INS_I32_POPCNT: // 0x69, i32.popcnt
	case wattoken.INS_I32_ADD: // 0x6a, i32.add
	case wattoken.INS_I32_SUB: // 0x6b, i32.sub
	case wattoken.INS_I32_MUL: // 0x6c, i32.mul
	case wattoken.INS_I32_DIV_S: // 0x6d, i32.div_s
	case wattoken.INS_I32_DIV_U: // 0x6e, i32.div_u
	case wattoken.INS_I32_REM_S: // 0x6f, i32.rem_s
	case wattoken.INS_I32_REM_U: // 0x70, i32.rem_u
	case wattoken.INS_I32_AND: // 0x71, i32.and
	case wattoken.INS_I32_OR: // 0x72, i32.or
	case wattoken.INS_I32_XOR: // 0x73, i32.xor
	case wattoken.INS_I32_SHL: // 0x74, i32.shl
	case wattoken.INS_I32_SHR_S: // 0x75, i32.shr_s
	case wattoken.INS_I32_SHR_U: // 0x76, i32.shr_u
	case wattoken.INS_I32_ROTL: // 0x77, i32.rotl
	case wattoken.INS_I32_ROTR: // 0x78, i32.rotr
	case wattoken.INS_I64_CLZ: // 0x79, i64.clz
	case wattoken.INS_I64_CTZ: // 0x7a, i64.ctz
	case wattoken.INS_I64_POPCNT: // 0x7b, i64.popcnt
	case wattoken.INS_I64_ADD: // 0x7c, i64.add
		// TODO
	case wattoken.INS_I64_SUB: // 0x7d, i64.sub
		// TODO
	case wattoken.INS_I64_MUL: // 0x7e, i64.mul
	case wattoken.INS_I64_DIV_S: // 0x7f, i64.div_s
	case wattoken.INS_I64_DIV_U: // 0x80, i64.div_u
	case wattoken.INS_I64_REM_S: // 0x81, i64.rem_s
	case wattoken.INS_I64_REM_U: // 0x82, i64.rem_u
	case wattoken.INS_I64_AND: // 0x83, i64.and
	case wattoken.INS_I64_OR: // 0x84, i64.or
	case wattoken.INS_I64_XOR: // 0x85, i64.xor
	case wattoken.INS_I64_SHL: // 0x86, i64.shl
	case wattoken.INS_I64_SHR_S: // 0x87, i64.shr_s
	case wattoken.INS_I64_SHR_U: // 0x88, i64.shr_u
	case wattoken.INS_I64_ROTL: // 0x89, i64.rotl
	case wattoken.INS_I64_ROTR: // 0x8a, i64.rotr
	case wattoken.INS_F32_ABS: // 0x8b, f32.abs
	case wattoken.INS_F32_NEG: // 0x8c, f32.neg
	case wattoken.INS_F32_CEIL: // 0x8d, f32.ceil
	case wattoken.INS_F32_FLOOR: // 0x8e, f32.floor
	case wattoken.INS_F32_TRUNC: // 0x8f, f32.trunc
	case wattoken.INS_F32_NEAREST: // 0x90, f32.nearest
	case wattoken.INS_F32_SQRT: // 0x91, f32.sqrt
	case wattoken.INS_F32_ADD: // 0x92, f32.add
	case wattoken.INS_F32_SUB: // 0x93, f32.sub
	case wattoken.INS_F32_MUL: // 0x94, f32.mul
	case wattoken.INS_F32_DIV: // 0x95, f32.div
	case wattoken.INS_F32_MIN: // 0x96, f32.min
	case wattoken.INS_F32_MAX: // 0x97, f32.max
	case wattoken.INS_F32_COPYSIGN: // 0x98, f32.copysign
	case wattoken.INS_F64_ABS: // 0x99, f64.abs
	case wattoken.INS_F64_NEG: // 0x9a, f64.neg
	case wattoken.INS_F64_CEIL: // 0x9b, f64.ceil
	case wattoken.INS_F64_FLOOR: // 0x9c, f64.floor
	case wattoken.INS_F64_TRUNC: // 0x9d, f64.trunc
	case wattoken.INS_F64_NEAREST: // 0x9e, f64.nearest
	case wattoken.INS_F64_SQRT: // 0x9f, f64.sqrt
	case wattoken.INS_F64_ADD: // 0xa0, f64.add
	case wattoken.INS_F64_SUB: // 0xa1, f64.sub
	case wattoken.INS_F64_MUL: // 0xa2, f64.mul
	case wattoken.INS_F64_DIV: // 0xa3, f64.div
	case wattoken.INS_F64_MIN: // 0xa4, f64.min
	case wattoken.INS_F64_MAX: // 0xa5, f64.max
	case wattoken.INS_F64_COPYSIGN: // 0xa6, f64.copysign
	case wattoken.INS_I32_WRAP_I64: // 0xa7, i32.wrap_i64
	case wattoken.INS_I32_TRUNC_F32_S: // 0xa8, i32.trunc_f32_s
	case wattoken.INS_I32_TRUNC_F32_U: // 0xa9, i32.trunc_f32_u
	case wattoken.INS_I32_TRUNC_F64_S: // 0xaa, i32.trunc_f64_s
	case wattoken.INS_I32_TRUNC_F64_U: // 0xab, i32.trunc_f64_u
	case wattoken.INS_I64_EXTEND_I32_S: // 0xac, i64.extend_i32_s
	case wattoken.INS_I64_EXTEND_I32_U: // 0xad, i64.extend_i32_u
	case wattoken.INS_I64_TRUNC_F32_S: // 0xae, i64.trunc_f32_s
	case wattoken.INS_I64_TRUNC_F32_U: // 0xaf, i64.trunc_f32_u
	case wattoken.INS_I64_TRUNC_F64_S: // 0xb0, i64.trunc_f64_s
	case wattoken.INS_I64_TRUNC_F64_U: // 0xb1, i64.trunc_f64_u
	case wattoken.INS_F32_CONVERT_I32_S: // 0xb2, f32.convert_i32_s
	case wattoken.INS_F32_CONVERT_I32_U: // 0xb3, f32.convert_i32_u
	case wattoken.INS_F32_CONVERT_I64_S: // 0xb4, f32.convert_i64_s
	case wattoken.INS_F32_CONVERT_I64_U: // 0xb5, f32.convert_i64_u
	case wattoken.INS_F32_DEMOTE_F64: // 0xb6, f32.demote_f64
	case wattoken.INS_F64_CONVERT_I32_S: // 0xb7, f64.convert_i32_s
	case wattoken.INS_F64_CONVERT_I32_U: // 0xb8, f64.convert_i32_u
	case wattoken.INS_F64_CONVERT_I64_S: // 0xb9, f64.convert_i64_s
	case wattoken.INS_F64_CONVERT_I64_U: // 0xba, f64.convert_i64_u
	case wattoken.INS_F64_PROMOTE_F32: // 0xbb, f64.promote_f32
	case wattoken.INS_I32_REINTERPRET_F32: // 0xbc, i32.reinterpret_f32
	case wattoken.INS_I64_REINTERPRET_F64: // 0xbd, i64.reinterpret_f64
	case wattoken.INS_F32_REINTERPRET_I32: // 0xbe, f32.reinterpret_i32
	case wattoken.INS_F64_REINTERPRET_I64: // 0xbf, f64.reinterpret_i64
	}
	return nil
}
