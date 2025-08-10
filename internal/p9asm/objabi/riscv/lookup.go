// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

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
	switch {
	case r == 0:
		return "NONE"
	case r == REG_SP:
		return "SP"
	case REG_X0 <= r && r <= REG_X31:
		return fmt.Sprintf("X%d", r-REG_X0)
	case REG_F0 <= r && r <= REG_F31:
		return fmt.Sprintf("F%d", r-REG_F0)
	case REG_V0 <= r && r <= REG_V31:
		return fmt.Sprintf("V%d", r-REG_V0)
	}

	return fmt.Sprintf("riscv.badreg(%d)", r)

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
