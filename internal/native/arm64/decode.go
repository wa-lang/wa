// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package arm64

import "wa-lang.org/wa/internal/native/abi"

// 解析机器码指令
func Decode(x uint32) (as abi.As, arg *abi.AsArgument, err error) {
	as, arg, _, err = DecodeEx(x)
	return
}

// 解析机器码指令
func DecodeEx(x uint32) (as abi.As, arg *abi.AsArgument, argRaw *abi.AsRawArgument, err error) {
	//return decodeInst(x)
	panic("TODO")
}
