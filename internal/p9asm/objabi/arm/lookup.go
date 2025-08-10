// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package arm

import (
	"fmt"

	"wa-lang.org/wa/internal/p9asm/objabi"
)

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
	if r == 0 {
		return "NONE"
	}
	if REG_R0 <= r && r <= REG_R15 {
		return fmt.Sprintf("R%d", r-REG_R0)
	}
	if REG_F0 <= r && r <= REG_F15 {
		return fmt.Sprintf("F%d", r-REG_F0)
	}

	switch r {
	case REG_FPSR:
		return "FPSR"

	case REG_FPCR:
		return "FPCR"

	case REG_CPSR:
		return "CPSR"

	case REG_SPSR:
		return "SPSR"

	case REG_MB_SY:
		return "MB_SY"
	case REG_MB_ST:
		return "MB_ST"
	case REG_MB_ISH:
		return "MB_ISH"
	case REG_MB_ISHST:
		return "MB_ISHST"
	case REG_MB_NSH:
		return "MB_NSH"
	case REG_MB_NSHST:
		return "MB_NSHST"
	case REG_MB_OSH:
		return "MB_OSH"
	case REG_MB_OSHST:
		return "MB_OSHST"
	}

	return fmt.Sprintf("arm.badreg(%d)", r)
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
