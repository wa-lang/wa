// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

// 检查模板和参数
func (ctx *OpContextType) checkArgMarks(xlen int, as abi.As, arg *abi.AsArgument, marks ArgMarks) error {
	if marks&ARG_RD != 0 {
		if arg.Rd == 0 {
			return fmt.Errorf("%s: rd missing; arg: %v", AsString(as), arg)
		}
	} else {
		if arg.Rd != 0 {
			return fmt.Errorf("%s: rd was nozero; arg: %v", AsString(as), arg)
		}
	}
	if marks&ARG_RS1 != 0 {
		if arg.Rs1 == 0 {
			return fmt.Errorf("%s: rs1 missing; arg: %v", AsString(as), arg)
		}
	} else {
		if arg.Rs1 != 0 {
			return fmt.Errorf("%s: rs1 was nozero; arg: %v", AsString(as), arg)
		}
	}
	if marks&ARG_RS2 != 0 {
		if arg.Rs2 == 0 {
			return fmt.Errorf("%s: rs2 missing; arg: %v", AsString(as), arg)
		}
	} else {
		if arg.Rs2 != 0 {
			return fmt.Errorf("%s: rs2 was nozero; arg: %v", AsString(as), arg)
		}
	}
	if marks&ARG_RS3 != 0 {
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

func (ctx *OpContextType) checkArgImm(xlen int, as abi.As, arg *abi.AsArgument, marks ArgMarks) error {
	if marks&ARG_IMM == 0 {
		if arg.Imm != 0 {
			return fmt.Errorf("%s: imm was nozero; arg: %v", AsString(as), arg)
		}
		return nil
	}

	if ctx.HasShamt {
		switch xlen {
		case 32:
			if err := immFitsRange(int64(arg.Imm), ImmRanges_Shamt32); err != nil {
				return fmt.Errorf("%s: %w", AsString(as), err)
			}
		case 64:
			if err := immFitsRange(int64(arg.Imm), ImmRanges_Shamt64); err != nil {
				return fmt.Errorf("%s: %w", AsString(as), err)
			}
		default:
			panic("unreachable")
		}
		return nil
	}

	var err error
	switch ctx.Opcode.FormatType() {
	case I:
		err = immFitsRange(int64(arg.Imm), ImmRanges_IType)
	case S:
		err = immFitsRange(int64(arg.Imm), ImmRanges_SType)
	case B:
		err = immFitsRange(int64(arg.Imm), ImmRanges_BType)
	case U:
		err = immFitsRange(int64(arg.Imm), ImmRanges_UType)
	case J:
		err = immFitsRange(int64(arg.Imm), ImmRanges_JType)
	default:
		panic("unreachable")
	}
	if err != nil {
		return fmt.Errorf("%s: %w", AsString(as), err)
	}

	return nil
}
