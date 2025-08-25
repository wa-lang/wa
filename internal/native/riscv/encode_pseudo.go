// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import "wa-lang.org/wa/internal/native/abi"

func (ctx *OpContextType) encodePseudo(xlen int, as abi.As, arg *abi.AsArgument) (uint32, error) {
	// TODO: 展开伪代码 后 继续调用 ctx.encode() 处理
	panic("TODO")
}

// 其他涉及地址的指令
// 超出范围的跳转需要拆分为2个指令
// if offset < -(1<<20) || (1<<20) <= offset {}
