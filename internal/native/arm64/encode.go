// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package arm64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

// 返回寄存器机器码编号
func RegI(r abi.RegType) uint32 {
	//return (*_OpContextType)(nil).regI(r)
	panic("TODO")
}

// 返回浮点数寄存器机器码编号
func RegF(r abi.RegType) uint32 {
	//return (*_OpContextType)(nil).regF(r)
	panic("TODO")
}

// 编码龙芯指令
func Encode(cpu abi.CPUType, as abi.As, arg *abi.AsArgument) (uint32, error) {
	switch cpu {
	case abi.ARM64:
		return EncodeARM64(as, arg)
	default:
		return 0, fmt.Errorf("unknonw cpu: %v", cpu)
	}
}

// 编码ARM64指令
func EncodeARM64(as abi.As, arg *abi.AsArgument) (uint32, error) {
	panic("TODO")
}
