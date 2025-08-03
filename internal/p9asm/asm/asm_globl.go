// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"wa-lang.org/wa/internal/p9asm/asm/lex"
	"wa-lang.org/wa/internal/p9asm/obj"
)

// GLOBL 汇编伪指令
// GLOBL shifts<>(SB),8,$256
// GLOBL shifts<>(SB),$256
func (p *Parser) asmGlobl(operands [][]lex.Token) {
	if len(operands) != 2 && len(operands) != 3 {
		p.errorf("expect two or three operands for GLOBL")
	}

	// Operand 0 has the general form foo<>+0x04(SB).
	nameAddr := p.address(operands[0])
	p.validateSymbol("GLOBL", &nameAddr, false)
	next := 1

	// Next operand is the optional flag, a literal integer.
	var flag = int64(0)
	if len(operands) == 3 {
		flag = p.evalInteger("GLOBL", operands[1])
		next++
	}

	// Final operand is an immediate constant.
	addr := p.address(operands[next])
	p.validateImmediate("GLOBL", &addr)

	// log.Printf("GLOBL %s %d, $%d", name, flag, size)
	prog := &obj.Prog{
		Ctxt:   p.ctxt,
		As:     obj.AGLOBL,
		Lineno: p.lineNum,
		From:   nameAddr,
		From3: &obj.Addr{
			Offset: flag,
		},
		To: addr,
	}
	p.append(prog, "", false)
}
