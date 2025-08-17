// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import (
	"fmt"
)

// 寄存器名字列表
var Register = make([]string, REG_END)

func init() {
	// 初始化寄存器名字表格
	for i := RegType(0); i < REG_END; i++ {
		Register[i] = RegString(i)
	}
}

// 根据名字查找寄存器
func LookupRegister(regName string) (r RegType, ok bool) {
	for i, s := range Register {
		if s == regName {
			return RegType(i), true
		}
	}
	return 0, false
}

// 寄存器转字符串格式
func RegString(r RegType) string {
	switch {
	case r == REG_SP:
		return "SP"
	case REG_X0 <= r && r <= REG_X31:
		return fmt.Sprintf("X%d", r-REG_X0)
	case REG_F0 <= r && r <= REG_F31:
		return fmt.Sprintf("F%d", r-REG_F0)
	}

	return fmt.Sprintf("riscv.badreg(%d)", r)

}

// 根据名字查找汇编指令
func LookupAs(asName string) (as As, ok bool) {
	for i, s := range Anames {
		if s == asName {
			return As(i), true
		}
	}
	return 0, false
}

// 汇编指令转字符串格式
func AsString(as As) string {
	return Anames[as]
}
