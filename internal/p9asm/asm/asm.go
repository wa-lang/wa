// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

import (
	"bytes"
	"text/scanner"

	"wa-lang.org/wa/internal/p9asm/arch"
	"wa-lang.org/wa/internal/p9asm/lex"
	"wa-lang.org/wa/internal/p9asm/obj"
)

// TODO: configure the architecture

var testOut *bytes.Buffer // Gathers output when testing.

// append adds the Prog to the end of the program-thus-far.
// If doLabel is set, it also defines the labels collect for this Prog.
func (p *Parser) append(prog *obj.Prog, cond string, doLabel bool) {
	if cond != "" {
		switch p.arch.CPU {
		case arch.ARM64:
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
	return p.getConstantPseudo(pseudo, &addr)
}

// validateImmediate checks that addr represents an immediate constant.
func (p *Parser) validateImmediate(pseudo string, addr *obj.Addr) {
	if addr.Type != obj.TYPE_CONST || addr.Name != 0 || addr.Reg != 0 || addr.Index != 0 {
		p.errorf("%s: expected immediate constant; found %s", pseudo, obj.Dconv(&emptyProg, addr))
	}
}

// asmText assembles a TEXT pseudo-op.
// TEXT runtime·sigtramp(SB),4,$0-0
func (p *Parser) asmText(word string, operands [][]lex.Token) {
	if len(operands) != 2 && len(operands) != 3 {
		p.errorf("expect two or three operands for TEXT")
	}

	// Labels are function scoped. Patch existing labels and
	// create a new label space for this TEXT.
	p.patch()
	p.labels = make(map[string]*obj.Prog)

	// Operand 0 is the symbol name in the form foo(SB).
	// That means symbol plus indirect on SB and no offset.
	nameAddr := p.address(operands[0])
	p.validateSymbol("TEXT", &nameAddr, false)
	name := nameAddr.Sym.Name
	next := 1

	// Next operand is the optional text flag, a literal integer.
	var flag = int64(0)
	if len(operands) == 3 {
		flag = p.evalInteger("TEXT", operands[1])
		next++
	}

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
	argSize := int64(obj.ArgsSizeUnknown)
	if len(op) > 0 {
		// There is an argument size. It must be a minus sign followed by a non-negative integer literal.
		if len(op) != 2 || op[0].ScanToken != '-' || op[1].ScanToken != scanner.Int {
			p.errorf("TEXT %s: argument size must be of form -integer", name)
		}
		argSize = p.positiveAtoi(op[1].String())
	}
	prog := &obj.Prog{
		Ctxt: p.ctxt,
		As:   obj.ATEXT,
		From: nameAddr,
		From3: &obj.Addr{
			Type:   obj.TYPE_CONST,
			Offset: flag,
		},
		To: obj.Addr{
			Type:   obj.TYPE_TEXTSIZE,
			Offset: frameSize,
			// Argsize set below.
		},
	}
	prog.To.Val = int32(argSize)

	p.append(prog, "", true)
}

// asmData assembles a DATA pseudo-op.
// DATA masks<>+0x00(SB)/4, $0x00000000
func (p *Parser) asmData(word string, operands [][]lex.Token) {
	if len(operands) != 2 {
		p.errorf("expect two operands for DATA")
	}

	// Operand 0 has the general form foo<>+0x04(SB)/4.
	op := operands[0]
	n := len(op)
	if n < 3 || op[n-2].ScanToken != '/' || op[n-1].ScanToken != scanner.Int {
		p.errorf("expect /size for DATA argument")
	}
	scale := p.parseScale(op[n-1].String())
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
		Ctxt: p.ctxt,
		As:   obj.ADATA,
		From: nameAddr,
		From3: &obj.Addr{
			Offset: int64(scale),
		},
		To: valueAddr,
	}

	p.append(prog, "", false)
}

// asmGlobl assembles a GLOBL pseudo-op.
// GLOBL shifts<>(SB),8,$256
// GLOBL shifts<>(SB),$256
func (p *Parser) asmGlobl(word string, operands [][]lex.Token) {
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
		Ctxt: p.ctxt,
		As:   obj.AGLOBL,
		From: nameAddr,
		From3: &obj.Addr{
			Offset: flag,
		},
		To: addr,
	}
	p.append(prog, "", false)
}

// asmPCData assembles a PCDATA pseudo-op.
// PCDATA $2, $705
func (p *Parser) asmPCData(word string, operands [][]lex.Token) {
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
		Ctxt: p.ctxt,
		As:   obj.APCDATA,
		From: key,
		To:   value,
	}
	p.append(prog, "", true)
}

// asmFuncData assembles a FUNCDATA pseudo-op.
// FUNCDATA $1, funcdata<>+4(SB)
func (p *Parser) asmFuncData(word string, operands [][]lex.Token) {
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
		Ctxt: p.ctxt,
		As:   obj.AFUNCDATA,
		From: valueAddr,
		To:   nameAddr,
	}
	p.append(prog, "", true)
}

// asmJump assembles a jump instruction.
// JMP	R1
// JMP	exit
// JMP	3(PC)
func (p *Parser) asmJump(op int, cond string, a []obj.Addr) {
	var target *obj.Addr
	prog := &obj.Prog{
		Ctxt: p.ctxt,
		As:   int16(op),
	}
	switch len(a) {
	case 1:
		target = &a[0]
	case 2:
		// Special 2-operand jumps.
		target = &a[1]
		prog.From = a[0]
	default:
		p.errorf("wrong number of arguments to %s instruction", obj.Aconv(op))
		return
	}
	switch {
	case target.Type == obj.TYPE_BRANCH:
		// JMP 4(PC)
		prog.To = obj.Addr{
			Type:   obj.TYPE_BRANCH,
			Offset: p.pc + 1 + target.Offset, // +1 because p.pc is incremented in append, below.
		}
	case target.Type == obj.TYPE_REG:
		// JMP R1
		prog.To = *target
	case target.Type == obj.TYPE_MEM && (target.Name == obj.NAME_EXTERN || target.Name == obj.NAME_STATIC):
		// JMP main·morestack(SB)
		prog.To = *target
	case target.Type == obj.TYPE_INDIR && (target.Name == obj.NAME_EXTERN || target.Name == obj.NAME_STATIC):
		// JMP *main·morestack(SB)
		prog.To = *target
		prog.To.Type = obj.TYPE_INDIR
	case target.Type == obj.TYPE_MEM && target.Reg == 0 && target.Offset == 0:
		// JMP exit
		if target.Sym == nil {
			// Parse error left name unset.
			return
		}
		targetProg := p.labels[target.Sym.Name]
		if targetProg == nil {
			p.toPatch = append(p.toPatch, Patch{prog, target.Sym.Name})
		} else {
			p.branch(prog, targetProg)
		}
	case target.Type == obj.TYPE_MEM && target.Name == obj.NAME_NONE:
		// JMP 4(R0)
		prog.To = *target
	case target.Type == obj.TYPE_CONST:
		// JMP $4
		prog.To = a[0]
	default:
		p.errorf("cannot assemble jump %+v", target)
	}

	p.append(prog, cond, true)
}

func (p *Parser) patch() {
	for _, patch := range p.toPatch {
		targetProg := p.labels[patch.label]
		if targetProg == nil {
			p.errorf("undefined label %s", patch.label)
		} else {
			p.branch(patch.prog, targetProg)
		}
	}
	p.toPatch = p.toPatch[:0]
}

func (p *Parser) branch(jmp, target *obj.Prog) {
	jmp.To = obj.Addr{
		Type:  obj.TYPE_BRANCH,
		Index: 0,
	}
	jmp.To.Val = target
}

// asmInstruction assembles an instruction.
// MOVW R9, (R10)
func (p *Parser) asmInstruction(op int, cond string, a []obj.Addr) {
	// fmt.Printf("%s %+v\n", obj.Aconv(op), a)
	prog := &obj.Prog{
		Ctxt: p.ctxt,
		As:   int16(op),
	}
	switch len(a) {
	case 0:
		// Nothing to do.
	case 1:
		if p.arch.LinkArch.UnaryDst[op] {
			// prog.From is no address.
			prog.To = a[0]
		} else {
			prog.From = a[0]
			// prog.To is no address.
		}
	case 2:
		if p.arch.CPU == arch.ARM64 && arch.IsARM64CMP(op) {
			prog.From = a[0]
			prog.Reg = p.getRegister(prog, op, &a[1])
			break
		}
		prog.From = a[0]
		prog.To = a[1]
	case 3:
		switch p.arch.CPU {
		case arch.ARM64:
			// ARM64 instructions with one input and two outputs.
			if arch.IsARM64STLXR(op) {
				prog.From = a[0]
				prog.To = a[1]
				if a[2].Type != obj.TYPE_REG {
					p.errorf("invalid addressing modes for third operand to %s instruction, must be register", obj.Aconv(op))
				}
				prog.RegTo2 = a[2].Reg
				break
			}
			prog.From = a[0]
			prog.Reg = p.getRegister(prog, op, &a[1])
			prog.To = a[2]
		case arch.AMD64:
			prog.From = a[0]
			prog.From3 = newAddr(a[1])
			prog.To = a[2]
		default:
			p.errorf("TODO: implement three-operand instructions for this architecture")
		}
	case 4:
		if p.arch.CPU == arch.ARM64 {
			prog.From = a[0]
			prog.Reg = p.getRegister(prog, op, &a[1])
			prog.From3 = newAddr(a[2])
			prog.To = a[3]
			break
		}
		p.errorf("can't handle %s instruction with 4 operands", obj.Aconv(op))
	case 5:
		p.errorf("can't handle %s instruction with 5 operands", obj.Aconv(op))
	case 6:
		fallthrough
	default:
		p.errorf("can't handle %s instruction with %d operands", obj.Aconv(op), len(a))
	}

	p.append(prog, cond, true)
}

// newAddr returns a new(Addr) initialized to x.
func newAddr(x obj.Addr) *obj.Addr {
	p := new(obj.Addr)
	*p = x
	return p
}

var emptyProg obj.Prog

// getConstantPseudo checks that addr represents a plain constant and returns its value.
func (p *Parser) getConstantPseudo(pseudo string, addr *obj.Addr) int64 {
	if addr.Type != obj.TYPE_MEM || addr.Name != 0 || addr.Reg != 0 || addr.Index != 0 {
		p.errorf("%s: expected integer constant; found %s", pseudo, obj.Dconv(&emptyProg, addr))
	}
	return addr.Offset
}

// getConstant checks that addr represents a plain constant and returns its value.
func (p *Parser) getConstant(prog *obj.Prog, op int, addr *obj.Addr) int64 {
	if addr.Type != obj.TYPE_MEM || addr.Name != 0 || addr.Reg != 0 || addr.Index != 0 {
		p.errorf("%s: expected integer constant; found %s", obj.Aconv(op), obj.Dconv(prog, addr))
	}
	return addr.Offset
}

// getImmediate checks that addr represents an immediate constant and returns its value.
func (p *Parser) getImmediate(prog *obj.Prog, op int, addr *obj.Addr) int64 {
	if addr.Type != obj.TYPE_CONST || addr.Name != 0 || addr.Reg != 0 || addr.Index != 0 {
		p.errorf("%s: expected immediate constant; found %s", obj.Aconv(op), obj.Dconv(prog, addr))
	}
	return addr.Offset
}

// getRegister checks that addr represents a register and returns its value.
func (p *Parser) getRegister(prog *obj.Prog, op int, addr *obj.Addr) int16 {
	if addr.Type != obj.TYPE_REG || addr.Offset != 0 || addr.Name != 0 || addr.Index != 0 {
		p.errorf("%s: expected register; found %s", obj.Aconv(op), obj.Dconv(prog, addr))
	}
	return addr.Reg
}
