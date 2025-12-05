// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package loong64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

// 解析机器码指令
func Decode(x uint32) (as abi.As, arg *abi.AsArgument, err error) {
	as, arg, _, err = DecodeEx(x)
	return
}

// 解析机器码指令
func DecodeEx(x uint32) (as abi.As, arg *abi.AsArgument, argRaw *abi.AsRawArgument, err error) {
	return decodeInst(x)
}

func decodeInst(x uint32) (as abi.As, arg *abi.AsArgument, argRaw *abi.AsRawArgument, err error) {
	for _, op := range _AOpContextTable {
		if op.mask == 0 && op.value == 0 {
			continue
		}
		if x&op.mask == op.value {
			return op.decodeInst(x)
		}
	}
	err = fmt.Errorf("loong64.decodeInst(%x): not found", x)
	return
}

func (op _OpContextType) decodeInst(x uint32) (as abi.As, arg *abi.AsArgument, argRaw *abi.AsRawArgument, err error) {
	as = op.op
	arg = new(abi.AsArgument)
	argRaw = new(abi.AsRawArgument)

	rd := (x >> (5 * 0)) & 0b11111 // 4:0 位
	rj := (x >> (5 * 1)) & 0b11111 // 9:5 位
	rk := (x >> (5 * 2)) & 0b11111 // 14:10 位
	fa := (x >> (5 * 3)) & 0b11111 // 19:15 位

	switch op.fmt {
	case OpFormatType_NULL:
		return
	case OpFormatType_2R:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegI(rj)
		return
	case OpFormatType_2F:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = op.decodeRegF(rd)
		arg.Rs1 = op.decodeRegF(rj)
		return
	case OpFormatType_1F_1R:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = op.decodeRegF(rd)
		arg.Rs1 = op.decodeRegI(rj)
		return
	case OpFormatType_1R_1F:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegF(rj)
		return
	case OpFormatType_3R:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Rs2 = rk
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Rs2 = op.decodeRegI(rk)
		return
	case OpFormatType_3F:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Rs2 = rk
		arg.Rd = op.decodeRegF(rd)
		arg.Rs1 = op.decodeRegF(rj)
		arg.Rs2 = op.decodeRegF(rk)
		return
	case OpFormatType_1F_2R:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Rs2 = rk
		arg.Rd = op.decodeRegF(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Rs2 = op.decodeRegI(rk)
		return
	case OpFormatType_4F:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Rs2 = rk
		argRaw.Rs3 = fa
		arg.Rd = op.decodeRegF(rd)
		arg.Rs1 = op.decodeRegF(rj)
		arg.Rs2 = op.decodeRegF(rk)
		arg.Rs3 = op.decodeRegF(fa)
		return
	case OpFormatType_2R_ui5:
		imm := int32(uimm(x, 10, 5))
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Imm = imm
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Imm = imm
		return
	case OpFormatType_2R_ui6:
		imm := int32(uimm(x, 10, 6))
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Imm = imm
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Imm = imm
		return
	case OpFormatType_2R_si12:
		imm := simm(x, 10, 12)
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Imm = imm
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Imm = imm
		return
	case OpFormatType_2R_ui12:
		imm := int32(uimm(x, 10, 12))
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Imm = imm
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Imm = imm
		return
	case OpFormatType_2R_si14:
		imm := int32(uimm(x, 10, 14))
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Imm = imm
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Imm = imm
		return
	case OpFormatType_2R_si16:
		imm := int32(uimm(x, 10, 16))
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Imm = imm
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Imm = imm
		return
	case OpFormatType_1R_si20:
		imm := int32(uimm(x, 5, 20))
		argRaw.Rd = rd
		argRaw.Imm = imm
		arg.Rd = op.decodeRegF(rd)
		arg.Imm = imm
		return
	case OpFormatType_0_2R:
		argRaw.Rs1 = rj
		argRaw.Rs2 = rk
		arg.Rs1 = op.decodeRegI(rj)
		arg.Rs2 = op.decodeRegI(rk)
		return
	case OpFormatType_3R_sa2:
		imm := int32(uimm(x, 10, 2))
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Rs2 = rk
		argRaw.Imm = imm
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Rs2 = op.decodeRegI(rk)
		arg.Imm = imm
		return
	case OpFormatType_3R_sa3:
		imm := int32(uimm(x, 10, 3))
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Rs2 = rk
		argRaw.Imm = imm
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Rs2 = op.decodeRegI(rk)
		arg.Imm = imm
		return
	case OpFormatType_code:
		imm := int32(uimm(x, 0, 15))
		argRaw.Imm = imm
		arg.Imm = imm
		return
	case OpFormatType_code_1R_si12:
		code := rd
		imm := int32(uimm(x, 10, 12))
		argRaw.Rd = code
		argRaw.Rs1 = rj
		argRaw.Imm = imm
		arg.Rd = abi.RegType(code)
		argRaw.Rs1 = rj
		arg.Imm = imm
		return
	case OpFormatType_2R_msbw_lsbw:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Rs2 = rk
		argRaw.Rs3 = fa
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Rs2 = abi.RegType(rk)
		arg.Rs3 = abi.RegType(fa)
		return
	case OpFormatType_2R_msbd_lsbd:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Rs2 = rk
		argRaw.Rs3 = fa
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Rs2 = abi.RegType(rk)
		arg.Rs3 = abi.RegType(fa)
		return
	case OpFormatType_fcsr_1R:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = abi.RegType(rd)
		arg.Rs1 = op.decodeRegI(rj)
		return
	case OpFormatType_1R_fcsr:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = abi.RegType(rj)
		return
	case OpFormatType_cd_1R:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = abi.RegType(rd)
		arg.Rs1 = op.decodeRegI(rj)
		return
	case OpFormatType_cd_1F:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = abi.RegType(rd)
		arg.Rs1 = op.decodeRegF(rj)
		return
	case OpFormatType_cd_2F:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Rs2 = rk
		arg.Rd = abi.RegType(rd)
		arg.Rs1 = op.decodeRegF(rj)
		arg.Rs2 = op.decodeRegF(rk)
		return
	case OpFormatType_1R_cj:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = abi.RegType(rj)
		return
	case OpFormatType_1F_cj:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = op.decodeRegF(rd)
		arg.Rs1 = abi.RegType(rj)
		return
	case OpFormatType_1R_csr:
		imm := int32(uimm(x, 10, 14))
		argRaw.Rd = rd
		argRaw.Imm = imm
		arg.Rd = op.decodeRegI(rd)
		arg.Imm = imm
		return
	case OpFormatType_2R_csr:
		imm := int32(uimm(x, 10, 14))
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Imm = imm
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Imm = imm
		return
	case OpFormatType_2R_level:
		imm := int32(uimm(x, 10, 8))
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Imm = imm
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Imm = imm
		return
	case OpFormatType_level:
		imm := int32(uimm(x, 0, 15))
		argRaw.Imm = imm
		arg.Imm = imm
		return
	case OpFormatType_0_1R_seq:
		imm := int32(uimm(x, 10, 8))
		argRaw.Rs1 = rj
		argRaw.Imm = imm
		arg.Rs1 = op.decodeRegI(rj)
		arg.Imm = imm
		return
	case OpFormatType_op_2R:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Rs2 = rk
		arg.Rd = abi.RegType(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Rs2 = abi.RegType(rk)
		return
	case OpFormatType_3F_ca:
		imm := int32(uimm(x, 15, 3))
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Rs2 = rk
		argRaw.Imm = imm
		arg.Rd = op.decodeRegF(rd)
		arg.Rs1 = op.decodeRegF(rj)
		arg.Rs2 = op.decodeRegF(rk)
		arg.Imm = imm
		return
	case OpFormatType_hint_1R_si12:
		imm := int32(uimm(x, 10, 12))
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Imm = imm
		arg.Rd = abi.RegType(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Imm = imm
		return
	case OpFormatType_hint_2R:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Rs2 = rk
		arg.Rd = abi.RegType(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Rs2 = op.decodeRegI(rk)
		return
	case OpFormatType_hint:
		imm := int32(uimm(x, 0, 15))
		argRaw.Imm = imm
		arg.Imm = imm
		return
	case OpFormatType_cj_offset:
		imm := int32(rd | (uimm(x, 10, 16) << 16))
		argRaw.Rs1 = rj
		argRaw.Imm = imm
		arg.Rs1 = abi.RegType(rj)
		arg.Imm = imm
		return
	case OpFormatType_rj_offset:
		imm := int32(rd | (uimm(x, 10, 16) << 16))
		argRaw.Rs1 = rj
		argRaw.Imm = imm
		arg.Rs1 = op.decodeRegI(rj)
		arg.Imm = imm
		return
	case OpFormatType_rj_rd_offset:
		imm := int32(uimm(x, 10, 16))
		argRaw.Rs1 = rj
		argRaw.Rd = rd
		argRaw.Imm = imm
		arg.Rs1 = op.decodeRegI(rj)
		arg.Rd = op.decodeRegI(rd)
		arg.Imm = imm
		return
	case OpFormatType_rd_rj_offset:
		imm := int32(uimm(x, 10, 16))
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		argRaw.Imm = imm
		arg.Rd = op.decodeRegI(rd)
		arg.Rs1 = op.decodeRegI(rj)
		arg.Imm = imm
		return
	case OpFormatType_offset:
		imm := int32((uimm(x, 0, 10) << 10) | (uimm(x, 10, 16)))
		argRaw.Imm = imm
		arg.Imm = imm
		return
	default:
		panic("unreachable")
	}
}

// 解码寄存器
func (op _OpContextType) decodeRegI(r uint32) abi.RegType {
	return abi.RegType(r) + REG_R0
}

// 解码寄存器(浮点数)
func (op _OpContextType) decodeRegF(r uint32) abi.RegType {
	return abi.RegType(r) + REG_F0
}
