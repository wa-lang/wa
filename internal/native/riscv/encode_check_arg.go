// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

// 检查模板和参数
func (ctx *_OpContextType) checkArgMarks(xlen int, as abi.As, arg *abi.AsArgument, marks _ArgMarks) error {
	if marks&_ARG_RD_IS_X == 0 {
		if marks&_ARG_RD != 0 {
			if arg.Rd == 0 {
				return fmt.Errorf("%s: rd missing; arg: %v", AsString(as), arg)
			}
		} else {
			if arg.Rd != 0 {
				return fmt.Errorf("%s: rd was nozero; arg: %v", AsString(as), arg)
			}
		}
	}
	if marks&_ARG_RS1 != 0 {
		if arg.Rs1 == 0 {
			return fmt.Errorf("%s: rs1 missing; arg: %v", AsString(as), arg)
		}
	} else {
		if arg.Rs1 != 0 {
			return fmt.Errorf("%s: rs1 was nozero; arg: %v", AsString(as), arg)
		}
	}
	if marks&_ARG_RS2 != 0 {
		if arg.Rs2 == 0 {
			return fmt.Errorf("%s: rs2 missing; arg: %v", AsString(as), arg)
		}
	} else {
		if arg.Rs2 != 0 {
			return fmt.Errorf("%s: rs2 was nozero; arg: %v", AsString(as), arg)
		}
	}
	if marks&_ARG_RS3 != 0 {
		if arg.Rs3 == 0 {
			return fmt.Errorf("%s: rs3 missing; arg: %v", AsString(as), arg)
		}
	} else {
		if arg.Rs3 != 0 {
			return fmt.Errorf("%s: rs3 was nozero; arg: %v", AsString(as), arg)
		}
	}

	return ctx.checkArgImm(xlen, as, arg, marks)
}

func (ctx *_OpContextType) checkArgImm(xlen int, as abi.As, arg *abi.AsArgument, marks _ArgMarks) error {
	if marks&_ARG_IMM == 0 {
		if arg.Imm != 0 {
			return fmt.Errorf("%s: imm was nozero; arg: %v", AsString(as), arg)
		}
		return nil
	}

	if ctx.HasShamt {
		switch xlen {
		case 32:
			if err := immFitsRange(int64(arg.Imm), _ImmRanges_Shamt32); err != nil {
				return fmt.Errorf("%s: %w", AsString(as), err)
			}
		case 64:
			if err := immFitsRange(int64(arg.Imm), _ImmRanges_Shamt64); err != nil {
				return fmt.Errorf("%s: %w", AsString(as), err)
			}
		default:
			panic("unreachable")
		}
		return nil
	}

	var err error
	switch ctx.Opcode.FormatType() {
	case _I:
		err = immFitsRange(int64(arg.Imm), _ImmRanges_IType)
	case _S:
		err = immFitsRange(int64(arg.Imm), _ImmRanges_SType)
	case _B:
		err = immFitsRange(int64(arg.Imm), _ImmRanges_BType)
	case _U:
		err = immFitsRange(int64(arg.Imm), _ImmRanges_UType)
	case _J:
		err = immFitsRange(int64(arg.Imm), _ImmRanges_JType)
	default:
		panic("unreachable")
	}
	if err != nil {
		return fmt.Errorf("%s: %w", AsString(as), err)
	}

	return nil
}
