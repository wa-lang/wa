// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"wa-lang.org/wa/internal/p9asm/asm/lex"
	"wa-lang.org/wa/internal/p9asm/obj"
)

// asmFuncData assembles a FUNCDATA pseudo-op.
// FUNCDATA $1, funcdata<>+4(SB)
func (p *Parser) asmFuncData(operands [][]lex.Token) {
	if len(operands) != 2 {
		p.errorf("expect two operands for FUNCDATA")
	}

	// Operand 0 must be an immediate constant.
	valueAddr := p.address(operands[0])
	p.validateImmediate("FUNCDATA", &valueAddr)

	// Operand 1 is a symbol name in the form foo(SB).
	nameAddr := p.address(operands[1])
	p.validateSymbol("FUNCDATA", &nameAddr, true)

	prog := &obj.Prog{
		Ctxt:   p.ctxt,
		As:     obj.AFUNCDATA,
		Lineno: p.lineNum,
		From:   valueAddr,
		To:     nameAddr,
	}
	p.append(prog, "", true)
}
