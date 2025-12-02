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
	ra := (x >> (5 * 3)) & 0b11111 // 19:15 位

	_, _, _, _ = rd, rj, rk, ra

	for _, argTyp := range op.args {
		if argTyp == 0 {
			break
		}
		switch argTyp {
		case Arg_fd:
			argRaw.Rd = rd
			arg.Rd, _ = op.decodeRegF(rd)
		case Arg_fj:
			argRaw.Rs1 = rj
			arg.Rs1, _ = op.decodeRegF(rj)
		case Arg_fk:
			argRaw.Rs2 = rk
			arg.Rs2, _ = op.decodeRegF(rk)
		case Arg_fa:
			argRaw.Rs3 = ra
			arg.Rs3, _ = op.decodeRegF(ra)

		case Arg_rd:
			argRaw.Rd = rd
			arg.Rd, _ = op.decodeRegI(rd)
		case Arg_rj:
			argRaw.Rs1 = rj
			arg.Rs1, _ = op.decodeRegI(rj)
		case Arg_rk:
			argRaw.Rs2 = rk
			arg.Rs2, _ = op.decodeRegI(rk)

		case Arg_op_4_0:
			argRaw.Imm = int32(rd)
			arg.Imm = int32(rd)
		case Arg_fcsr_4_0:
			argRaw.Imm = int32(rd)
			arg.Imm = int32(rd)
		case Arg_fcsr_9_5:
			argRaw.Imm = int32(rj)
			arg.Imm = int32(rj)
		case Arg_csr_23_10:
			v := (x >> 10) & ((1 << 14) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)

		case Arg_cd:
			v := FCC0 + Fcc(x&((1<<3)-1))
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)
		case Arg_cj:
			v := FCC0 + Fcc((x>>5)&((1<<3)-1))
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)
		case Arg_ca:
			v := FCC0 + Fcc((x>>15)&((1<<3)-1))
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)

		case Arg_sa2_16_15:
			v := (x >> 15) & ((1 << 2) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)
		case Arg_sa3_17_15:
			v := (x >> 15) & ((1 << 3) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)

		case Arg_code_4_0:
			v := x & ((1 << 5) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)
		case Arg_code_14_0:
			v := x & ((1 << 15) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)

		case Arg_ui5_14_10:
			v := (x >> 10) & ((1 << 5) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)
		case Arg_ui6_15_10:
			v := (x >> 10) & ((1 << 6) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)
		case Arg_ui12_21_10:
			v := ((x >> 10) & ((1 << 12) - 1) & 0xfff)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)

		case Arg_lsbw:
			v := (x >> 10) & ((1 << 5) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)
		case Arg_msbw:
			v := (x >> 16) & ((1 << 5) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)
		case Arg_lsbd:
			v := (x >> 10) & ((1 << 6) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)
		case Arg_msbd:
			v := (x >> 16) & ((1 << 6) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)

		case Arg_hint_4_0:
			v := x & ((1 << 5) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)
		case Arg_hint_14_0:
			v := x & ((1 << 15) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)

		case Arg_level_14_0:
			v := x & ((1 << 15) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)
		case Arg_level_17_10:
			v := (x >> 10) & ((1 << 8) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)

		case Arg_seq_17_10:
			v := (x >> 10) & ((1 << 8) - 1)
			argRaw.Imm = int32(v)
			arg.Imm = int32(v)

		case Arg_si12_21_10:
			if (x & 0x200000) == 0x200000 {
				v := ((x >> 10) & ((1 << 12) - 1)) | 0xf000
				argRaw.Imm = int32(v)
				arg.Imm = int32(v)
			} else {
				v := (x >> 10) & ((1 << 12) - 1)
				argRaw.Imm = int32(v)
				arg.Imm = int32(v)
			}

		case Arg_si14_23_10:
			if (x & 0x800000) == 0x800000 {
				v := (((x >> 10) & ((1 << 14) - 1)) << 2) | 0xffff0000
				argRaw.Imm = int32(v)
				arg.Imm = int32(v)
			} else {
				v := (((x >> 10) & ((1 << 14) - 1)) << 2)
				argRaw.Imm = int32(v)
				arg.Imm = int32(v)
			}
		case Arg_si16_25_10:
			if (x & 0x2000000) == 0x2000000 {
				v := ((x >> 10) & ((1 << 16) - 1)) | 0xffff0000
				argRaw.Imm = int32(v)
				arg.Imm = int32(v)
			} else {
				v := ((x >> 10) & ((1 << 16) - 1))
				argRaw.Imm = int32(v)
				arg.Imm = int32(v)
			}
		case Arg_si20_24_5:
			if (x & 0x1000000) == 0x1000000 {
				v := ((x >> 5) & ((1 << 20) - 1)) | 0xfff00000
				argRaw.Imm = int32(v)
				arg.Imm = int32(v)
			} else {
				v := ((x >> 5) & ((1 << 20) - 1))
				argRaw.Imm = int32(v)
				arg.Imm = int32(v)
			}
		case Arg_offset_20_0:
			if (x & 0x10) == 0x10 {
				v := ((((x << 16) | ((x >> 10) & ((1 << 16) - 1))) & ((1 << 21) - 1)) << 2) | 0xff800000
				argRaw.Imm = int32(v)
				arg.Imm = int32(v)
			} else {
				v := (((x << 16) | ((x >> 10) & ((1 << 16) - 1))) & ((1 << 21) - 1)) << 2
				argRaw.Imm = int32(v)
				arg.Imm = int32(v)
			}

		case Arg_offset_25_0:
			if (x & 0x200) == 0x200 {
				v := ((((x << 16) | ((x >> 10) & ((1 << 16) - 1))) & ((1 << 26) - 1)) << 2) | 0xf0000000
				argRaw.Imm = int32(v)
				arg.Imm = int32(v)
			} else {
				v := ((((x << 16) | ((x >> 10) & ((1 << 16) - 1))) & ((1 << 26) - 1)) << 2)
				argRaw.Imm = int32(v)
				arg.Imm = int32(v)
			}
		case Arg_offset_15_0:
			if (x & 0x2000000) == 0x2000000 {
				v := (((x >> 10) & ((1 << 16) - 1)) << 2) | 0xfffc0000
				argRaw.Imm = int32(v)
				arg.Imm = int32(v)
			} else {
				v := (((x >> 10) & ((1 << 16) - 1)) << 2)
				argRaw.Imm = int32(v)
				arg.Imm = int32(v)
			}

		default:
			// 遇到无法识别的参数类型，返回错误
			return 0, nil, nil, fmt.Errorf("unknown argument type: %d", argTyp)
		}
	}

	// 成功解码所有参数后，返回结果
	return as, arg, argRaw, nil
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
