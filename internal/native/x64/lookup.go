// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

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
	return 0, false
}

// 寄存器转字符串格式
func RegString(r abi.RegType) string {
	if int(r) < len(_Register) {
		if s := _Register[int(r)]; s != "" {
			return s
		}
	}
	return fmt.Sprintf("x64.badreg(%d)", r)
}

// 寄存器转字符串格式(32位格式)
func Reg32String(r abi.RegType) string {
	if int(r) < len(_Register32) {
		if s := _Register32[int(r)]; s != "" {
			return s
		}
	}
	return fmt.Sprintf("x64.badreg(%d)", r)
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
	return fmt.Sprintf("riscv.badas(%d)", int(as))
}

func AsOpFormatType(as abi.As) OpFormatType {
	if int(as) < len(x64ModeTable) {
		return x64ModeTable[as]
	}
	return 0
}
