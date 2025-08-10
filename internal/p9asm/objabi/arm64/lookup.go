// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package arm64

import (
	"fmt"

	"wa-lang.org/wa/internal/p9asm/objabi"
)

var strcond = [16]string{
	"EQ",
	"NE",
	"HS",
	"LO",
	"MI",
	"PL",
	"VS",
	"VC",
	"HI",
	"LS",
	"GE",
	"LT",
	"GT",
	"LE",
	"AL",
	"NV",
}

// 寄存器名字列表
var Register = make([]string, RegMax-RBase)

func init() {
	// 初始化寄存器名字表格
	for i := RBase; i < RegMax; i++ {
		Register[i-RBase] = RegString(i)
	}
}

// 根据名字查找寄存器, 失败返回 objabi.REG_NONE
func LookupRegister(regName string) objabi.RBaseType {
	for i, s := range Register {
		if s == regName {
			return RBase + objabi.RBaseType(i)
		}
	}
	return objabi.REG_NONE
}

// 寄存器转字符串格式
func RegString(r objabi.RBaseType) string {
	switch {
	case REG_R0 <= r && r <= REG_R30:
		return fmt.Sprintf("R%d", r-REG_R0)
	case r == REG_R31:
		return "ZR"
	case REG_F0 <= r && r <= REG_F31:
		return fmt.Sprintf("F%d", r-REG_F0)
	case REG_V0 <= r && r <= REG_V31:
		return fmt.Sprintf("V%d", r-REG_F0)
	case COND_EQ <= r && r <= COND_NV:
		return strcond[r-COND_EQ]
	case r == REGSP:
		return "RSP"
	case r == REG_DAIF:
		return "DAIF"
	case r == REG_NZCV:
		return "NZCV"
	case r == REG_FPSR:
		return "FPSR"
	case r == REG_FPCR:
		return "FPCR"
	case r == REG_SPSR_EL1:
		return "SPSR_EL1"
	case r == REG_ELR_EL1:
		return "ELR_EL1"
	case r == REG_SPSR_EL2:
		return "SPSR_EL2"
	case r == REG_ELR_EL2:
		return "ELR_EL2"
	case r == REG_CurrentEL:
		return "CurrentEL"
	case r == REG_SP_EL0:
		return "SP_EL0"
	case r == REG_SPSel:
		return "SPSel"
	case r == REG_DAIFSet:
		return "DAIFSet"
	case r == REG_DAIFClr:
		return "DAIFClr"
	}
	return fmt.Sprintf("arm64.badreg(%d)", r)
}

// 根据名字查找汇编指令, 失败返回 objabi.AXXX
func LookupAs(asName string) objabi.As {
	for i, s := range Anames {
		if s == asName {
			return ABase + objabi.As(i)
		}
	}
	return objabi.AXXX
}

// 汇编指令转字符串格式
func AsString(as objabi.As) string {
	if ABase <= as && as < AsMax {
		return Anames[as-ABase]
	}
	panic("unreachable")
}
