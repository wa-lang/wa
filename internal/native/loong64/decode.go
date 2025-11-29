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

	rd := (x >> (5 * 0)) & 0b11111
	rj := (x >> (5 * 1)) & 0b11111
	rk := (x >> (5 * 2)) & 0b11111
	ra := (x >> (5 * 3)) & 0b11111

	_, _, _, _ = rd, rj, rk, ra

	for _, argTyp := range op.args {
		switch argTyp {
		case 0: // 空参数
			return

		// 1-5 bit
		case arg_fd:
			argRaw.Rd = rd
			arg.Rd, _ = op.decodeRegF(rd)
		case arg_fj:
			argRaw.Rs1 = rj
			arg.Rs1, _ = op.decodeRegF(rd)
		case arg_fk:
			argRaw.Rs2 = rk
			arg.Rs2, _ = op.decodeRegF(rd)
		case arg_fa:
			argRaw.Rs3 = ra
			arg.Rs3, _ = op.decodeRegF(rd)
		case arg_rd:
			argRaw.Rd = rd
			arg.Rd, _ = op.decodeRegI(rd)

		// 6-10 bit
		case arg_rj:
			argRaw.Rs1 = rj
			arg.Rs1, _ = op.decodeRegI(rd)
		case arg_rk:
			argRaw.Rs2 = rk
			arg.Rs2, _ = op.decodeRegI(rd)
		case arg_op_4_0:
			panic("TODO")
		case arg_fcsr_4_0:
			panic("TODO")
		case arg_fcsr_9_5:
			panic("TODO")

		// 11-15 bit
		case arg_csr_23_10:
			panic("TODO")
		case arg_cd:
			panic("TODO")
		case arg_cj:
			panic("TODO")
		case arg_ca:
			panic("TODO")
		case arg_sa2_16_15:
			panic("TODO")

		// 16-20 bit
		case arg_sa3_17_15:
			panic("TODO")
		case arg_code_4_0:
			panic("TODO")
		case arg_code_14_0:
			panic("TODO")
		case arg_ui5_14_10:
			panic("TODO")
		case arg_ui6_15_10:
			panic("TODO")

		// 21-25 bit
		case arg_ui12_21_10:
			panic("TODO")
		case arg_lsbw:
			panic("TODO") // 需要2个立即数
		case arg_msbw:
			panic("TODO") // 需要2个立即数
		case arg_lsbd:
			panic("TODO") // 需要2个立即数
		case arg_msbd:
			panic("TODO") // 需要2个立即数

		// 26-30 bit
		case arg_hint_4_0:
			panic("TODO")
		case arg_hint_14_0:
			panic("TODO")
		case arg_level_14_0:
			panic("TODO")
		case arg_level_17_10:
			panic("TODO")

		case arg_seq_17_10:
			panic("TODO")

		// 31-35 bit
		case arg_si12_21_10:
			panic("TODO")
		case arg_si14_23_10:
			panic("TODO")
		case arg_si16_25_10:
			panic("TODO")
		case arg_si20_24_5:
			panic("TODO")
		case arg_offset_20_0:
			panic("TODO")

		// 36~
		case arg_offset_25_0:
			panic("TODO")
		case arg_offset_15_0:
			panic("TODO")
		}

		panic("unreachable")
	}

	panic("TODO")
}

// 解码寄存器
func (op _OpContextType) decodeRegI(r uint32) (reg abi.RegType, err error) {
	if r <= 31 {
		return abi.RegType(r) + REG_R0, nil
	}
	return 0, fmt.Errorf("badreg(%d)", r)
}

// 解码寄存器(浮点数)
func (op _OpContextType) decodeRegF(r uint32) (reg abi.RegType, err error) {
	if r <= 31 {
		return abi.RegType(r) + REG_F0, nil
	}
	return 0, fmt.Errorf("badreg(%d)", r)
}
