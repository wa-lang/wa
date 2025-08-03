// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"wa-lang.org/wa/internal/p9asm/asm/lex"
	"wa-lang.org/wa/internal/p9asm/obj"
)

// asmPCData assembles a PCDATA pseudo-op.
// PCDATA $2, $705
func (p *Parser) asmPCData(operands [][]lex.Token) {
	if len(operands) != 2 {
		p.errorf("expect two operands for PCDATA")
	}

	// Operand 0 must be an immediate constant.
	key := p.address(operands[0])
	p.validateImmediate("PCDATA", &key)

	// Operand 1 must be an immediate constant.
	value := p.address(operands[1])
	p.validateImmediate("PCDATA", &value)

	// log.Printf("PCDATA $%d, $%d", key.Offset, value.Offset)
	prog := &obj.Prog{
		Ctxt:   p.ctxt,
		As:     obj.APCDATA,
		Lineno: p.lineNum,
		From:   key,
		To:     value,
	}
	p.append(prog, "", true)
}
