// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x64

import (
	"wa-lang.org/wa/internal/native/x64/x86asm"
)

// 底层的指令别名
type Inst x86asm.Inst

// 解析机器码指令
// 这是底层的 p9x86 定义的格式, 可用于格式化打印
// 到 abi.As/arg *abi.AsArgument 类型需要手动转化
func Decode(p []byte, mode int) (inst *Inst, err error) {
	if x, err := x86asm.Decode(p, 64); err != nil {
		return nil, err
	} else {
		return (*Inst)(&x), nil
	}
}

func (p *Inst) String() string {
	return (*x86asm.Inst)(p).String()
}
