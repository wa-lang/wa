// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package framestack

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/token"
)

type Framestack interface {
	HeadSize() int
	ArgRegNum() int

	AllocArg(typ token.Token) (reg abi.RegType, off int)
	AllocRet(typ token.Token) (reg abi.RegType, off int)
	AllocLocal(typ token.Token, cap int) (off int)
}

func NewFramestack(cpu abi.CPUType) Framestack {
	switch cpu {
	case abi.LOONG64:
		return NewLAFramestack(cpu)
	case abi.RISCV32:
		return NewRVFramestack(cpu)
	case abi.RISCV64:
		return NewRVFramestack(cpu)
	}
	return nil
}
