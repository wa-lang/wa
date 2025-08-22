// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import "fmt"

// 解析机器码指令
func Decode(x uint32) (as As, arg *AsArgument, err error) {
	// 根据 opcode 查指令类型
	opcode := OpcodeType(x) & OpBase_Mask

	// 根据指令类型解析参数和功能码
	switch opcode.FormatType() {
	case R:
		return opcode.decodeR(x)
	case R4:
		return opcode.decodeR4(x)
	case I:
		return opcode.decodeI(x)
	case S:
		return opcode.decodeS(x)
	case B:
		return opcode.decodeB(x)
	case U:
		return opcode.decodeU(x)
	case J:
		return opcode.decodeJ(x)
	default:
		err = fmt.Errorf("unknown opcode %07b", opcode)
		return
	}
}

func (op OpcodeType) decodeR(x uint32) (as As, arg *AsArgument, err error) {
	arg = new(AsArgument)

	rd := (x >> 7) & 0b_1_1111
	rs1 := (x >> 15) & 0b_1_1111
	rs2 := (x >> 20) & 0b_1_1111

	funct3 := (x >> 12) & 0b_111
	funct7 := (x >> 25) & 0b_111_1111

	if op == OpBase_OP_FP {
		if arg.Rd, err = op.decodeRegF(rd); err != nil {
			return 0, nil, err
		}
		if arg.Rs1, err = op.decodeRegF(rs1); err != nil {
			return 0, nil, err
		}
		if arg.Rs2, err = op.decodeRegF(rs2); err != nil {
			return 0, nil, err
		}
	} else {
		if arg.Rd, err = op.decodeRegI(rd); err != nil {
			return 0, nil, err
		}
		if arg.Rs1, err = op.decodeRegI(rs1); err != nil {
			return 0, nil, err
		}
		if arg.Rs2, err = op.decodeRegI(rs2); err != nil {
			return 0, nil, err
		}
	}

	// 查询表格
	for i, ctx := range AOpContextTable {
		if ctx.Opcode == op {
			if ctx.Funct3 == funct3 && ctx.Funct7 == funct7 {
				as = As(i)
				break
			}
		}
	}
	if as == 0 {
		err = fmt.Errorf("decodeR: opcode=%07b, funct3=%03b, funct7=%07b", op, funct3, funct7)
		return
	}

	// OK
	return
}

func (op OpcodeType) decodeR4(x uint32) (as As, arg *AsArgument, err error) {
	arg = new(AsArgument)

	rd := (x >> 7) & 0b_1_1111
	rs1 := (x >> 15) & 0b_1_1111
	rs2 := (x >> 20) & 0b_1_1111
	rs3 := (x >> 27) & 0b_1_1111

	funct3 := (x >> 12) & 0b_111
	funct2 := (x >> 25) & 0b_11

	if arg.Rd, err = op.decodeRegF(rd); err != nil {
		return 0, nil, err
	}
	if arg.Rs1, err = op.decodeRegF(rs1); err != nil {
		return 0, nil, err
	}
	if arg.Rs2, err = op.decodeRegF(rs2); err != nil {
		return 0, nil, err
	}
	if arg.Rs3, err = op.decodeRegF(rs3); err != nil {
		return 0, nil, err
	}

	// 查询表格
	for i, ctx := range AOpContextTable {
		if ctx.Opcode == op {
			if ctx.Funct3 == funct3 && ctx.Funct7 == funct2 {
				as = As(i)
				break
			}
		}
	}
	if as == 0 {
		err = fmt.Errorf("decodeR4: opcode=%07b, funct3=%03b, funct7=%02b", op, funct3, funct2)
		return
	}

	// OK
	return
}

func (op OpcodeType) decodeI(x uint32) (as As, arg *AsArgument, err error) {
	arg = new(AsArgument)

	rd := (x >> 7) & 0b_1_1111
	rs1 := (x >> 15) & 0b_1_1111
	imm := int32(x) >> 20

	funct3 := (x >> 12) & 0b_111

	if op == OpBase_LOAD_FP {
		if arg.Rd, err = op.decodeRegF(rd); err != nil {
			return 0, nil, err
		}
		if arg.Rs1, err = op.decodeRegF(rs1); err != nil {
			return 0, nil, err
		}
	} else {
		if arg.Rd, err = op.decodeRegI(rd); err != nil {
			return 0, nil, err
		}
		if arg.Rs1, err = op.decodeRegI(rs1); err != nil {
			return 0, nil, err
		}
	}
	arg.Imm = imm

	// 查询表格
	for i, ctx := range AOpContextTable {
		if ctx.Opcode == op {
			if ctx.Funct3 == funct3 {
				as = As(i)
				break
			}
		}
	}
	if as == 0 {
		err = fmt.Errorf("decodeI: opcode=%07b, funct3=%03b", op, funct3)
		return
	}

	// OK
	return
}

func (op OpcodeType) decodeS(x uint32) (as As, arg *AsArgument, err error) {
	arg = new(AsArgument)

	rs1 := (x >> 15) & 0b_1_1111
	rs2 := (x >> 20) & 0b_1_1111
	imm := (int32(x)>>25)<<5 | int32(x>>7)&0b_1_1111

	funct3 := (x >> 12) & 0b_111

	if op == OpBase_STORE_FP {
		if arg.Rs1, err = op.decodeRegF(rs1); err != nil {
			return 0, nil, err
		}
		if arg.Rs2, err = op.decodeRegF(rs2); err != nil {
			return 0, nil, err
		}
	} else {
		if arg.Rs1, err = op.decodeRegI(rs1); err != nil {
			return 0, nil, err
		}
		if arg.Rs2, err = op.decodeRegI(rs2); err != nil {
			return 0, nil, err
		}
	}
	arg.Imm = imm

	// 查询表格
	for i, ctx := range AOpContextTable {
		if ctx.Opcode == op {
			if ctx.Funct3 == funct3 {
				as = As(i)
				break
			}
		}
	}
	if as == 0 {
		err = fmt.Errorf("decodeS: opcode=%07b, funct3=%03b", op, funct3)
		return
	}

	// OK
	return
}

func (op OpcodeType) decodeB(x uint32) (as As, arg *AsArgument, err error) {
	arg = new(AsArgument)

	rs1 := (x >> 15) & 0b_1_1111
	rs2 := (x >> 20) & 0b_1_1111

	imm12 := x & (1 << 31)
	imm5_10 := ((x >> 25) & 0b_11_1111) << 5
	imm1_4 := ((x >> 8) & 0b_1111) << 1
	imm11 := ((x >> 7) & 0b_1) << 11
	imm := int32(imm12 | imm11 | imm5_10 | imm1_4)

	funct3 := (x >> 12) & 0b_111

	if op == OpBase_STORE_FP {
		if arg.Rs1, err = op.decodeRegF(rs1); err != nil {
			return 0, nil, err
		}
		if arg.Rs2, err = op.decodeRegF(rs2); err != nil {
			return 0, nil, err
		}
	} else {
		if arg.Rs1, err = op.decodeRegI(rs1); err != nil {
			return 0, nil, err
		}
		if arg.Rs2, err = op.decodeRegI(rs2); err != nil {
			return 0, nil, err
		}
	}
	arg.Imm = imm

	// 查询表格
	for i, ctx := range AOpContextTable {
		if ctx.Opcode == op {
			if ctx.Funct3 == funct3 {
				as = As(i)
				break
			}
		}
	}
	if as == 0 {
		err = fmt.Errorf("decodeB: opcode=%07b, funct3=%03b", op, funct3)
		return
	}

	// OK
	return
}

func (op OpcodeType) decodeU(x uint32) (as As, arg *AsArgument, err error) {
	arg = new(AsArgument)

	rd := (x >> 7) & 0b_1_1111

	imm := int32(x >> 12) // U 模式汇编指令不包含低 12bit 部分

	if arg.Rd, err = op.decodeRegI(rd); err != nil {
		return 0, nil, err
	}

	arg.Imm = imm

	// 查询表格
	for i, ctx := range AOpContextTable {
		if ctx.Opcode == op {
			as = As(i)
			break
		}
	}
	if as == 0 {
		err = fmt.Errorf("decodeU: opcode=%07b", op)
		return
	}

	// OK
	return
}

func (op OpcodeType) decodeJ(x uint32) (as As, arg *AsArgument, err error) {
	arg = new(AsArgument)

	rd := (x >> 7) & 0b_1_1111

	imm20 := (x >> 31) << 20
	imm1_10 := (x << 1) >> 22 << 1
	imm11 := (x << 11) >> 31 << 11
	imm12_19 := (x << 12) >> 24 << 12
	imm := imm20 | imm12_19 | imm11 | imm1_10
	if imm>>uint32(21-1) == 1 {
		imm |= 0x7ff << 21
	}

	if arg.Rd, err = op.decodeRegI(rd); err != nil {
		return 0, nil, err
	}

	arg.Imm = int32(imm)

	// 查询表格
	for i, ctx := range AOpContextTable {
		if ctx.Opcode == op {
			as = As(i)
			break
		}
	}
	if as == 0 {
		err = fmt.Errorf("decodeJ: opcode=%07b", op)
		return
	}

	// OK
	return
}

// 解码寄存器
func (op OpcodeType) decodeRegI(r uint32) (reg RegType, err error) {
	if r <= 31 {
		return RegType(r) + REG_X0, nil
	}
	return 0, fmt.Errorf("badreg(%d)", r)
}

// 解码寄存器(浮点数)
func (op OpcodeType) decodeRegF(r uint32) (reg RegType, err error) {
	if r <= 31 {
		return RegType(r) + REG_F0, nil
	}
	return 0, fmt.Errorf("badreg(%d)", r)
}
