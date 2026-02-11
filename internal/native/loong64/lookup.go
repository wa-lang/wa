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
		return fmt.Sprintf("$r%d", r-REG_R0)
	case REG_F0 <= r && r <= REG_F31:
		return fmt.Sprintf("$f%d", r-REG_F0)
	}
	if r >= 0 && r < REG_END {
		if int(r) < len(_Register) && _Register[r] != "" {
			return _Register[r]
		}
	}
	return fmt.Sprintf("loong64.badreg(%d)", r)
}
func ZhRegString(r abi.RegType) string {
	switch {
	case REG_R0 <= r && r <= REG_R31:
		return _ZhRegister[r]
	case REG_F0 <= r && r <= REG_F31:
		return _ZhRegister[r]
	}
	if r >= 0 && r < REG_END {
		if int(r) < len(_ZhRegister) && _ZhRegister[r] != "" {
			return _ZhRegister[r]
		}
	}
	return fmt.Sprintf("loong64.badreg(%d)", r)
}

// 寄存器别名
func RegAliasString(r abi.RegType) string {
	if r >= REG_R0 && r < REG_END {
		if int(r) < len(_RegisterAlias) {
			if s := _RegisterAlias[r]; s != "" {
				return s
			}
		}
	}
	return RegString(r)
}
func ZhRegAliasString(r abi.RegType) string {
	if r >= REG_R0 && r < REG_END {
		if int(r) < len(_ZhRegisterAlias) {
			if s := _ZhRegisterAlias[r]; s != "" {
				return s
			}
		}
	}
	return ZhRegString(r)
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
		if s := _Anames[as]; s != "" {
			return s
		}
	}
	return fmt.Sprintf("loong64.badas(%d)", int(as))
}

// 汇编指令转字符串格式
func ZhAsString(as abi.As, asName string) string {
	if asName != "" {
		return asName
	}
	if int(as) < len(_ZhAnames) {
		if s := _ZhAnames[as]; s != "" {
			return s
		}
	}
	return fmt.Sprintf("loong64.badas(%d)", int(as))
}

// 指令的编码格式
func AsFormatType(as abi.As) OpFormatType {
	if as > 0 && int(as) < len(_AOpContextTable) {
		return _AOpContextTable[as].fmt
	}
	return OpFormatType_NULL
}
