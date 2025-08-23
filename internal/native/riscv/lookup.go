// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

// 根据名字查找寄存器
func LookupRegister(regName string) (r abi.RegType, ok bool) {
	for i, s := range Register {
		if s == regName {
			return abi.RegType(i), true
		}
	}
	for i, s := range RegisterAlias {
		if s == regName {
			return abi.RegType(i), true
		}
	}
	return 0, false
}

// 寄存器转字符串格式
func RegString(r abi.RegType) string {
	switch {
	case REG_X0 <= r && r <= REG_X31:
		return fmt.Sprintf("X%d", r-REG_X0)
	case REG_F0 <= r && r <= REG_F31:
		return fmt.Sprintf("F%d", r-REG_F0)
	}
	return fmt.Sprintf("riscv.badreg(%d)", r)
}

// 寄存器别名
func RegAliasString(r abi.RegType) string {
	if r >= REG_X0 && r < REG_END {
		if s := RegisterAlias[r]; s != "" {
			return s
		}
	}
	return RegString(r)
}

// 根据名字查找汇编指令
func LookupAs(asName string) (as abi.As, ok bool) {
	for i, s := range Anames {
		if s == asName {
			return abi.As(i), true
		}
	}
	return 0, false
}

// 汇编指令转字符串格式
func AsString(as abi.As) string {
	if int(as) < len(Anames) {
		return Anames[as]
	}
	return fmt.Sprintf("riscv.badas(%d)", int(as))
}
