// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

import (
	"bytes"
	"fmt"

	"wa-lang.org/wa/internal/p9asm/asm/arch"
	"wa-lang.org/wa/internal/p9asm/asm/lex"
	"wa-lang.org/wa/internal/p9asm/obj"
)

// TODO: configure the architecture

var testOut *bytes.Buffer // Gathers output when testing.

// append adds the Prog to the end of the program-thus-far.
// If doLabel is set, it also defines the labels collect for this Prog.
func (p *Parser) append(prog *obj.Prog, cond string, doLabel bool) {
	if cond != "" {
		switch p.arch.LinkArch.Thechar {
		case '5':
			if !arch.ARMConditionCodes(prog, cond) {
				p.errorf("unrecognized condition code .%q", cond)
			}

		case '7':
			if !arch.ARM64Suffix(prog, cond) {
				p.errorf("unrecognized suffix .%q", cond)
			}

		default:
			p.errorf("unrecognized suffix .%q", cond)
		}
	}
	if p.firstProg == nil {
		p.firstProg = prog
	} else {
		p.lastProg.Link = prog
	}
	p.lastProg = prog
	if doLabel {
		p.pc++
		for _, label := range p.pendingLabels {
			if p.labels[label] != nil {
				p.errorf("label %q multiply defined", label)
			}
			p.labels[label] = prog
		}
		p.pendingLabels = p.pendingLabels[0:0]
	}
	prog.Pc = int64(p.pc)
	if p.Flags.Debug {
		fmt.Println(p.lineNum, prog)
	}
	if testOut != nil {
		fmt.Fprintln(testOut, p.lineNum, prog)
	}
}

// validateSymbol checks that addr represents a valid name for a pseudo-op.
func (p *Parser) validateSymbol(pseudo string, addr *obj.Addr, offsetOk bool) {
	if addr.Name != obj.NAME_EXTERN && addr.Name != obj.NAME_STATIC || addr.Scale != 0 || addr.Reg != 0 {
		p.errorf("%s symbol %q must be a symbol(SB)", pseudo, addr.Sym.Name)
	}
	if !offsetOk && addr.Offset != 0 {
		p.errorf("%s symbol %q must not be offset from SB", pseudo, addr.Sym.Name)
	}
}

// evalInteger evaluates an integer constant for a pseudo-op.
func (p *Parser) evalInteger(pseudo string, operands []lex.Token) int64 {
	addr := p.address(operands)

	if addr.Type != obj.TYPE_MEM || addr.Name != 0 || addr.Reg != 0 || addr.Index != 0 {
		var emptyProg obj.Prog
		p.errorf("%s: expected integer constant; found %s", pseudo, addr.Dconv(&emptyProg))
	}
	return addr.Offset
}

// validateImmediate checks that addr represents an immediate constant.
func (p *Parser) validateImmediate(pseudo string, addr *obj.Addr) {
	if addr.Type != obj.TYPE_CONST || addr.Name != 0 || addr.Reg != 0 || addr.Index != 0 {
		var emptyProg obj.Prog
		p.errorf("%s: expected immediate constant; found %s", pseudo, addr.Dconv(&emptyProg))
	}
}
