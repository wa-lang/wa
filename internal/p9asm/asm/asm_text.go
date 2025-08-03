// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"text/scanner"

	"wa-lang.org/wa/internal/p9asm/asm/lex"
	"wa-lang.org/wa/internal/p9asm/obj"
)

// TEXT 汇编伪指令
// TEXT pkg·foo(SB),$framesize-argsize
// argsize 缺少时为 0
func (p *Parser) asmText(operands [][]lex.Token) {
	if len(operands) != 2 {
		p.errorf("expect two operands for TEXT") // 函数不再支持 flags 标志
	}

	// Labels are function scoped. Patch existing labels and
	// create a new label space for this TEXT.
	for _, patch := range p.toPatch {
		targetProg := p.labels[patch.label]
		if targetProg == nil {
			p.errorf("undefined label %s", patch.label)
		} else {
			patch.prog.To = obj.Addr{
				Type: obj.TYPE_BRANCH,
				Val:  targetProg,
			}
		}
	}

	// 清空 toPatch 列表, 已经完成
	p.toPatch = p.toPatch[:0]

	p.labels = make(map[string]*obj.Prog)

	// Operand 0 is the symbol name in the form foo(SB).
	// That means symbol plus indirect on SB and no offset.
	nameAddr := p.address(operands[0])
	p.validateSymbol("TEXT", &nameAddr, false)
	name := nameAddr.Sym.Name
	next := 1

	// Next operand is the frame and arg size.
	// Bizarre syntax: $frameSize-argSize is two words, not subtraction.
	// Both frameSize and argSize must be simple integers; only frameSize
	// can be negative.
	// The "-argSize" may be missing; if so, set it to obj.ArgsSizeUnknown.
	// Parse left to right.
	op := operands[next]
	if len(op) < 2 || op[0].ScanToken != '$' {
		p.errorf("TEXT %s: frame size must be an immediate constant", name)
		return
	}
	op = op[1:]
	negative := false
	if op[0].ScanToken == '-' {
		negative = true
		op = op[1:]
	}
	if len(op) == 0 || op[0].ScanToken != scanner.Int {
		p.errorf("TEXT %s: frame size must be an immediate constant", name)
		return
	}
	frameSize := p.positiveAtoi(op[0].String())
	if negative {
		frameSize = -frameSize
	}
	op = op[1:]

	// argsize不能省略, 不再支持 printf 那种参数未知的调用
	// There is an argument size. It must be a minus sign followed by a non-negative integer literal.
	argSize := int32(0)
	if len(op) > 0 {
		if len(op) != 2 || op[0].ScanToken != '-' || op[1].ScanToken != scanner.Int {
			p.errorf("TEXT %s: argument size must be of form -integer", name)
		}
		argSize = int32(p.positiveAtoi(op[1].String()))
	}

	prog := &obj.Prog{
		Ctxt:   p.ctxt,
		As:     obj.ATEXT,
		Lineno: p.histLineNum,
		From:   nameAddr,
		From3: &obj.Addr{
			Type: obj.TYPE_CONST,
		},
		To: obj.Addr{
			Type:   obj.TYPE_TEXTSIZE,
			Offset: frameSize,
			Val:    int32(argSize),
		},
	}

	p.append(prog, "", true)
}
