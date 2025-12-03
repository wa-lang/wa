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
		arg.Rd = op.decodeRegI(rd)
		return
	case OpFormatType_2F:
		argRaw.Rd = rd
		arg.Rd = op.decodeRegF(rd)
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
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = op.decodeRegF(rd)
		arg.Rs1 = op.decodeRegF(rj)
		arg.Imm = simm(x, 10, 5)
		return
	case OpFormatType_2R_ui6:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = op.decodeRegF(rd)
		arg.Rs1 = op.decodeRegF(rj)
		arg.Imm = simm(x, 10, 6)
		return
	case OpFormatType_2R_si12:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = op.decodeRegF(rd)
		arg.Rs1 = op.decodeRegF(rj)
		arg.Imm = simm(x, 10, 12)
		return
	case OpFormatType_2R_ui12:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = op.decodeRegF(rd)
		arg.Rs1 = op.decodeRegF(rj)
		arg.Imm = simm(x, 10, 12)
		return
	case OpFormatType_2R_si14:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = op.decodeRegF(rd)
		arg.Rs1 = op.decodeRegF(rj)
		arg.Imm = simm(x, 10, 14)
		return
	case OpFormatType_2R_si16:
		argRaw.Rd = rd
		argRaw.Rs1 = rj
		arg.Rd = op.decodeRegF(rd)
		arg.Rs1 = op.decodeRegF(rj)
		arg.Imm = simm(x, 10, 16)
		return
	case OpFormatType_1R_si20:
		argRaw.Rd = rd
		arg.Rd = op.decodeRegF(rd)
		arg.Imm = simm(x, 10, 20)
		return
	case OpFormatType_0_2R:
		panic("TODO")
	case OpFormatType_3R_sa2:
		panic("TODO")
	case OpFormatType_3R_sa3:
		panic("TODO")
	case OpFormatType_code:
		panic("TODO")
	case OpFormatType_code_1R_si12:
		panic("TODO")
	case OpFormatType_2R_msbw_lsbw:
		panic("TODO")
	case OpFormatType_2R_msbd_lsbd:
		panic("TODO")
	case OpFormatType_fcsr_1R:
		panic("TODO")
	case OpFormatType_1R_fcsr:
		panic("TODO")
	case OpFormatType_cd_1R:
		panic("TODO")
	case OpFormatType_cd_1F:
		panic("TODO")
	case OpFormatType_cd_2R:
		panic("TODO")
	case OpFormatType_cd_2F:
		panic("TODO")
	case OpFormatType_1R_cj:
		panic("TODO")
	case OpFormatType_1F_cj:
		panic("TODO")
	case OpFormatType_1R_csr:
		panic("TODO")
	case OpFormatType_2R_csr:
		panic("TODO")
	case OpFormatType_2R_level:
		panic("TODO")
	case OpFormatType_level:
		panic("TODO")
	case OpFormatType_0_1R_seq:
		panic("TODO")
	case OpFormatType_op_2R:
		panic("TODO")
	case OpFormatType_3R_ca:
		panic("TODO")
	case OpFormatType_hint_1R_si12:
		panic("TODO")
	case OpFormatType_hint_2R:
		panic("TODO")
	case OpFormatType_hint:
		panic("TODO")
	case OpFormatType_cj_offset:
		panic("TODO")
	case OpFormatType_rj_offset:
		panic("TODO")
	case OpFormatType_rj_rd_offset:
		panic("TODO")
	case OpFormatType_rd_rj_offset:
		panic("TODO")
	case OpFormatType_offset:
		panic("TODO")
	default:
		panic("unreachable")
	}

	// 成功解码所有参数后，返回结果
	return as, arg, argRaw, nil
}

// 解码寄存器
func (op _OpContextType) decodeRegI(r uint32) abi.RegType {
	return abi.RegType(r) + REG_R0
}

// 解码寄存器(浮点数)
func (op _OpContextType) decodeRegF(r uint32) abi.RegType {
	return abi.RegType(r) + REG_F0
}
