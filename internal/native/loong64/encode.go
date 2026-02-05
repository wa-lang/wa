// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package loong64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

// 返回寄存器机器码编号
func RegI(r abi.RegType) uint32 {
	return (*_OpContextType)(nil).regI(r)
}

// 返回浮点数寄存器机器码编号
func RegF(r abi.RegType) uint32 {
	return (*_OpContextType)(nil).regF(r)
}

// 编码龙芯指令
func Encode(cpu abi.CPUType, as abi.As, arg *abi.AsArgument) (uint32, error) {
	switch cpu {
	case abi.RISCV64:
		return EncodeLA64(as, arg)
	default:
		return 0, fmt.Errorf("unknonw cpu: %v", cpu)
	}
}

// 编码龙芯64指令
func EncodeLA64(as abi.As, arg *abi.AsArgument) (uint32, error) {
	if as <= 0 || as >= ALAST {
		return 0, fmt.Errorf("loong64.EncodeLA64: bad as(%v), arg(%v)", as, arg)
	}

	ctx := _AOpContextTable[as]
	assert(ctx.mask != 0)

	return ctx.encodeRaw(as, arg)
}

func (ctx *_OpContextType) encodeRaw(as abi.As, arg *abi.AsArgument) (x uint32, err error) {
	assert(ctx != nil)
	assert(ctx.op == as)
	assert(arg != nil)

	x = ctx.mask & ctx.value

	switch ctx.fmt {
	case OpFormatType_NULL:
		return
	case OpFormatType_2R:
		rd := ctx.regI(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		x |= (rj << 5) | rd
		return
	case OpFormatType_2F:
		fd := ctx.regF(arg.Rd)
		fj := ctx.regF(arg.Rs1)
		x |= (fj << 5) | fd
		return
	case OpFormatType_1F_1R:
		fd := ctx.regF(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		x |= (rj << 5) | fd
		return
	case OpFormatType_1R_1F:
		rd := ctx.regI(arg.Rd)
		fj := ctx.regF(arg.Rs1)
		x |= (fj << 5) | rd
		return
	case OpFormatType_3R:
		rd := ctx.regI(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		rk := ctx.regI(arg.Rs2)
		x |= (rk << 10) | (rj << 5) | rd
		return
	case OpFormatType_3F:
		fd := ctx.regF(arg.Rd)
		fj := ctx.regF(arg.Rs1)
		fk := ctx.regF(arg.Rs2)
		x |= (fk << 10) | (fj << 5) | fd
		return
	case OpFormatType_1F_2R:
		fd := ctx.regF(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		rk := ctx.regI(arg.Rs2)
		x |= (rk << 10) | (rj << 5) | fd
		return
	case OpFormatType_4F:
		fd := ctx.regF(arg.Rd)
		fj := ctx.regF(arg.Rs1)
		fk := ctx.regF(arg.Rs2)
		fa := ctx.regF(arg.Rs3)
		x |= (fa << 15) | (fk << 10) | (fj << 5) | fd
		return
	case OpFormatType_2R_ui5:
		assert(arg.Imm >= 0 && arg.Imm < (1<<5))
		rd := ctx.regI(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		ui5 := arg.Imm & 0x1F
		x |= (uint32(ui5) << 10) | (rj << 5) | rd
		return
	case OpFormatType_2R_ui6:
		assert(arg.Imm >= 0 && arg.Imm < (1<<6))
		rd := ctx.regI(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		ui5 := arg.Imm & 0x3F
		x |= (uint32(ui5) << 10) | (rj << 5) | rd
		return
	case OpFormatType_2R_si12:
		// 编码时候带符号的立即数正数部分范围可以放宽到无符号
		assert(arg.Imm >= -(1<<11) && arg.Imm < (1<<12))
		rd := ctx.regI(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		si12 := uint32(arg.Imm) & 0xFFF
		x |= (si12 << 10) | (rj << 5) | rd
		return
	case OpFormatType_1F_1R_si12:
		// 编码时候带符号的立即数正数部分范围可以放宽到无符号
		assert(arg.Imm >= -(1<<11) && arg.Imm < (1<<12))
		fd := ctx.regF(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		si12 := uint32(arg.Imm) & 0xFFF
		x |= (si12 << 10) | (rj << 5) | fd
		return
	case OpFormatType_2R_ui12:
		assert(arg.Imm >= 0 && arg.Imm < (1<<12))
		rd := ctx.regI(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		ui12 := arg.Imm & 0xFFF
		x |= (uint32(ui12) << 10) | (rj << 5) | rd
		return
	case OpFormatType_2R_si14:
		// 编码时候带符号的立即数正数部分范围可以放宽到无符号
		assert(arg.Imm >= -(1<<13) && arg.Imm < (1<<14))
		rd := ctx.regI(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		si14 := arg.Imm & 0x3FFF
		x |= (uint32(si14) << 10) | (rj << 5) | rd
		return
	case OpFormatType_1R_si20:
		// 编码时候带符号的立即数正数部分范围可以放宽到无符号
		assert(arg.Imm >= -(1<<19) && arg.Imm < (1<<20))
		rd := ctx.regI(arg.Rd)
		si20 := arg.Imm & 0xFFFFF
		x |= (uint32(si20) << 5) | rd
		return
	case OpFormatType_0_2R:
		rj := ctx.regI(arg.Rs1)
		rk := ctx.regI(arg.Rs2)
		x |= (rk << 10) | (rj << 5)
		return
	case OpFormatType_3R_sa2:
		rd := ctx.regI(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		rk := ctx.regI(arg.Rs2)
		sa2 := arg.Imm & 0xF
		x |= (rk << 14) | (uint32(sa2) << 10) | (rj << 5) | rd
		return
	case OpFormatType_3R_sa3:
		rd := ctx.regI(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		rk := ctx.regI(arg.Rs2)
		sa2 := arg.Imm & 0x1F
		x |= (rk << 14) | (uint32(sa2) << 10) | (rj << 5) | rd
		return
	case OpFormatType_code:
		code := arg.Imm & 0x7FFF
		x |= uint32(code)
		return
	case OpFormatType_code_1R_si12:
		// 编码时候带符号的立即数正数部分范围可以放宽到无符号
		assert(arg.Imm >= -(1<<11) && arg.Imm < (1<<12))
		code := uint32(arg.Rd) & 0b_1_1111
		rj := ctx.regI(arg.Rs1)
		si12 := arg.Imm & 0xFFF
		x |= (uint32(si12) << 10) | (rj << 5) | code
		return
	case OpFormatType_2R_msbw_lsbw:
		rd := ctx.regI(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		msbw := uint32(arg.Rs2) & 0b_0_1_1111
		lsbw := uint32(arg.Rs3) & 0b_0_1_1111
		x |= (msbw << 16) | (lsbw << 10) | (rj << 5) | rd
		return
	case OpFormatType_2R_msbd_lsbd:
		rd := ctx.regI(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		msbd := uint32(arg.Rs2) & 0b_1_1_1111
		lsbd := uint32(arg.Rs3) & 0b_1_1_1111
		x |= (msbd << 16) | (lsbd << 10) | (rj << 5) | rd
		return
	case OpFormatType_fcsr_1R:
		fcsr := ctx.regFCSR(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		x |= (rj << 5) | fcsr
		return
	case OpFormatType_1R_fcsr:
		rd := ctx.regI(arg.Rd)
		fcsr := ctx.regFCSR(arg.Rs1)
		x |= (fcsr << 5) | rd
		return
	case OpFormatType_cd_1R:
		cd := ctx.regFCC(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		x |= (rj << 5) | cd
		return
	case OpFormatType_cd_1F:
		cd := ctx.regFCC(arg.Rd)
		fj := ctx.regF(arg.Rs1)
		x |= (fj << 5) | cd
		return
	case OpFormatType_cd_2F:
		cd := ctx.regFCC(arg.Rd)
		fj := ctx.regF(arg.Rs1)
		fk := ctx.regF(arg.Rs1)
		x |= (fk << 10) | (fj << 5) | cd
		return
	case OpFormatType_1R_cj:
		rd := ctx.regI(arg.Rd)
		cj := ctx.regFCC(arg.Rs1)
		x |= (cj << 5) | rd
		return
	case OpFormatType_1F_cj:
		fd := ctx.regF(arg.Rd)
		cj := ctx.regFCC(arg.Rs1)
		x |= (cj << 5) | fd
		return
	case OpFormatType_1R_csr:
		rd := ctx.regI(arg.Rd)
		csr := uint32(arg.Imm) & 0x3FFF
		x |= (csr << 10) | rd
		return
	case OpFormatType_2R_csr:
		rd := ctx.regI(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		csr := uint32(arg.Imm) & 0x3FFF
		x |= (csr << 10) | (rj << 5) | rd
		return
	case OpFormatType_2R_level:
		rd := ctx.regI(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		level := uint32(arg.Imm) & 0xFFFF
		x |= (level << 10) | (rj << 5) | rd
		return
	case OpFormatType_level:
		level := uint32(arg.Imm) & 0x7FFF
		x |= level
		return
	case OpFormatType_0_1R_seq:
		rj := ctx.regI(arg.Rs1)
		seq := uint32(arg.Imm) & 0xFFFF
		x |= (seq << 10) | (rj << 5)
		return
	case OpFormatType_op_2R:
		op := uint32(arg.Rd) & 0b_1_1111
		rj := ctx.regI(arg.Rs1)
		rk := ctx.regI(arg.Rs2)
		x |= (rk << 10) | (rj << 5) | op
		return
	case OpFormatType_3F_ca:
		fd := ctx.regF(arg.Rd)
		fj := ctx.regF(arg.Rs1)
		fk := ctx.regF(arg.Rs2)
		ca := uint32(arg.Imm) & 0b_111
		x |= (ca << 15) | (fk << 10) | (fj << 5) | fd
		return
	case OpFormatType_hint_1R_si12:
		// 编码时候带符号的立即数正数部分范围可以放宽到无符号
		assert(arg.Imm >= -(1<<11) && arg.Imm < (1<<12))
		hint := uint32(arg.Rd) & 0b_1_1111
		rj := ctx.regI(arg.Rs1)
		si12 := uint32(arg.Imm) & 0xFFF
		x |= (si12 << 10) | (rj << 5) | hint
		return
	case OpFormatType_hint_2R:
		hint := uint32(arg.Rd) & 0b_1_1111
		rj := ctx.regI(arg.Rs1)
		rk := ctx.regI(arg.Rs2)
		x |= (rk << 10) | (rj << 5) | hint
		return
	case OpFormatType_hint:
		hint := uint32(arg.Imm) & 0x7FFF
		x |= hint
		return
	case OpFormatType_cj_offset:
		assert(arg.Imm&0b11 == 0)
		imm := uint32(arg.Imm >> 2)
		off16_20 := (imm >> 16) & 0b_1_1111
		cj := ctx.regFCC(arg.Rs1)
		off0_15 := imm & 0xFFFF
		x |= (off0_15 << 10) | (cj << 5) | off16_20
		return
	case OpFormatType_rj_offset:
		assert(arg.Imm&0b11 == 0)
		imm := uint32(arg.Imm >> 2)
		off16_20 := (imm >> 16) & 0b_1_1111
		rj := ctx.regI(arg.Rs1)
		off0_15 := imm & 0xFFFF
		x |= (off0_15 << 10) | (rj << 5) | off16_20
		return
	case OpFormatType_rj_rd_offset:
		assert(arg.Imm&0b11 == 0)
		imm := uint32(arg.Imm >> 2)
		rj := ctx.regI(arg.Rs1)
		rd := ctx.regI(arg.Rd)
		offset := imm & 0xFFFF
		x |= (offset << 10) | (rj << 5) | rd
		return
	case OpFormatType_rd_rj_offset:
		assert(arg.Imm&0b11 == 0)
		imm := uint32(arg.Imm >> 2)
		rd := ctx.regI(arg.Rd)
		rj := ctx.regI(arg.Rs1)
		offset := imm & 0xFFFF
		x |= (offset << 10) | (rj << 5) | rd
		return
	case OpFormatType_offset:
		assert(arg.Imm&0b11 == 0)
		imm := uint32(arg.Imm >> 2)
		off16_25 := (imm >> 16) & 0b_11_1111_1111
		off0_15 := imm & 0xFFFF
		x |= (off0_15 << 10) | off16_25
		return
	default:
		panic("unreachable")
	}
}
