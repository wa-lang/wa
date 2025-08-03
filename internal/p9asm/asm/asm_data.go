// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"text/scanner"

	"wa-lang.org/wa/internal/p9asm/asm/lex"
	"wa-lang.org/wa/internal/p9asm/obj"
)

// DATA 汇编伪指令
// DATA masks<>+0x00(SB)/4, $0x00000000
// masks 是全局变量对应的符号, '<>'表示为当前文件可见, '+0x00'表示偏移量,
// '(SB)'表示相对于全局的静态地址, '/4' 表示4字节, `0x00000000` 为初值
func (p *Parser) asmData(operands [][]lex.Token) {
	if len(operands) != 2 {
		p.errorf("expect two operands for DATA")
	}

	// Operand 0 has the general form foo<>+0x04(SB)/4.
	op := operands[0]
	n := len(op)
	if n < 3 || op[n-2].ScanToken != '/' || op[n-1].ScanToken != scanner.Int {
		p.errorf("expect /size for DATA argument")
	}

	var scale int8
	switch s := op[n-1].String(); s {
	case "1", "2", "4", "8":
		scale = int8(s[0] - '0')
	default:
		p.errorf("bad scale: %s", s)
		scale = 0
	}

	op = op[:n-2]
	nameAddr := p.address(op)
	p.validateSymbol("DATA", &nameAddr, true)
	name := nameAddr.Sym.Name

	// Operand 1 is an immediate constant or address.
	valueAddr := p.address(operands[1])
	switch valueAddr.Type {
	case obj.TYPE_CONST, obj.TYPE_FCONST, obj.TYPE_SCONST, obj.TYPE_ADDR:
		// OK
	default:
		p.errorf("DATA value must be an immediate constant or address")
	}

	// The addresses must not overlap. Easiest test: require monotonicity.
	if lastAddr, ok := p.dataAddr[name]; ok && nameAddr.Offset < lastAddr {
		p.errorf("overlapping DATA entry for %s", name)
	}
	p.dataAddr[name] = nameAddr.Offset + int64(scale)

	prog := &obj.Prog{
		Ctxt:   p.ctxt,
		As:     obj.ADATA,
		Lineno: p.histLineNum,
		From:   nameAddr,
		From3: &obj.Addr{
			Offset: int64(scale),
		},
		To: valueAddr,
	}

	p.append(prog, "", false)
}
