// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package loong64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

// 指令码有效
func AsValid(as abi.As) bool {
	return as > 0 && as < ALAST
}

// 寄存器有效
func RegValid(reg abi.RegType) bool {
	return reg > 0 && reg < REG_END
}

// 根据名字查找寄存器(忽略大小写, 忽略下划线和点的区别)
func LookupRegister(regName string) (r abi.RegType, ok bool) {
	if regName == "" {
		return
	}
	for i, s := range _Register {
		if strEqualFold(s, regName) {
			return abi.RegType(i), true
		}
	}
	for i, s := range _ZhRegister {
		if strEqualFold(s, regName) {
			return abi.RegType(i), true
		}
	}
	for i, s := range _RegisterAlias {
		if strEqualFold(s, regName) {
			return abi.RegType(i), true
		}
	}
	for i, s := range _ZhRegisterAlias {
		if strEqualFold(s, regName) {
			return abi.RegType(i), true
		}
	}
	return 0, false
}

// 寄存器转字符串格式
func RegString(r abi.RegType) string {
	switch {
	case REG_R0 <= r && r <= REG_R31:
		return fmt.Sprintf("R%d", r-REG_R0)
	case REG_F0 <= r && r <= REG_F31:
		return fmt.Sprintf("F%d", r-REG_F0)
	}
	return fmt.Sprintf("riscv.badreg(%d)", r)
}
func ZhRegString(r abi.RegType) string {
	switch {
	case REG_R0 <= r && r <= REG_R31:
		return _ZhRegister[r]
	case REG_F0 <= r && r <= REG_F31:
		return _ZhRegister[r]
	}
	return fmt.Sprintf("riscv.badreg(%d)", r)
}

// 寄存器别名
func RegAliasString(r abi.RegType) string {
	if r >= REG_R0 && r < REG_END {
		if s := _RegisterAlias[r]; s != "" {
			return s
		}
	}
	return RegString(r)
}
func ZhRegAliasString(r abi.RegType) string {
	if r >= REG_R0 && r < REG_END {
		if s := _ZhRegisterAlias[r]; s != "" {
			return s
		}
	}
	return RegString(r)
}

// 根据名字查找汇编指令(忽略大小写, 忽略下划线和点的区别)
func LookupAs(asName string) (as abi.As, ok bool) {
	if asName == "" {
		return
	}
	for i, s := range _Anames {
		if strEqualFold(s, asName) {
			return abi.As(i), true
		}
	}
	for i, s := range _ZhAnames {
		if strEqualFold(s, asName) {
			return abi.As(i), true
		}
	}
	return 0, false
}

// 汇编指令转字符串格式
func AsString(as abi.As, asName string) string {
	if asName != "" {
		return asName
	}
	if int(as) < len(_Anames) {
		return _Anames[as]
	}
	return fmt.Sprintf("riscv.badas(%d)", int(as))
}

// 汇编指令的参数格式
func AsArgs(as abi.As) [5]InstArg {
	if as > 0 && int(as) < len(_AOpContextTable) {
		ctx := _AOpContextTable[as]
		return ctx.args
	}
	return [5]InstArg{}
}
