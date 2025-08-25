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
		case 64:
		default:
			panic("unreachable")
		}
	}

	if ctx.ImmMin < ctx.ImmMax {
		if imm := int64(arg.Imm); imm < ctx.ImmMin || imm > ctx.ImmMax {
			return fmt.Errorf("%s: imm must be in [%d, %d]", AsString(as), ctx.ImmMin, ctx.ImmMax)
		}
	}
	return nil
}
