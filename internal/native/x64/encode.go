// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x64

import (
	"wa-lang.org/wa/internal/native/abi"
)

// 编码后指令的长度, 忽略 Symbol 的具体值
func EncodeLen(as abi.As, arg *abi.X64Argument) (size int, err error) {
	prog, err := BuildProg(as, arg)
	if err != nil {
		return 0, err
	}
	code := prog.Encode()
	return len(code), nil
}

// 指令编码, Symbol 对应的值需要提前解析到 Offset 属性中
func Encode(as abi.As, arg *abi.X64Argument) (code []byte, err error) {
	prog, err := BuildProg(as, arg)
	if err != nil {
		return nil, err
	}

	code = prog.Encode()
	return code, nil
}
