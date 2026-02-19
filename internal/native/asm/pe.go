// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

type _PEAssembler struct{}

func (p *_PEAssembler) asmFile(filename string, source []byte, opt *abi.LinkOptions) (prog *abi.LinkedProgram, err error) {
	return nil, fmt.Errorf("TODO: support windows/PE")
}
