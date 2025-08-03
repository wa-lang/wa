// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"strconv"
	"text/scanner"

	"wa-lang.org/wa/internal/p9asm/asm/arch"
	"wa-lang.org/wa/internal/p9asm/asm/lex"
	"wa-lang.org/wa/internal/p9asm/obj"
)

// atStartOfRegister reports whether the parser is at the start of a register definition.
func (p *Parser) atStartOfRegister(name string) bool {
	// Simple register: R10.
	_, present := p.arch.Register[name]
	if present {
		return true
	}
	// Parenthesized register: R(10).
	return p.arch.RegisterPrefix[name] && p.peek() == '('
}

// atRegisterShift reports whether we are at the start of an ARM shifted register.
// We have consumed the register or R prefix.
func (p *Parser) atRegisterShift() bool {
	// ARM only.
	if p.arch.LinkArch.Thechar != '5' {
		return false
	}
	// R1<<...
	if p.peek().IsRegisterShift() {
		return true
	}
	// R(1)<<...   Ugly check. TODO: Rethink how we handle ARM register shifts to be
	// less special.
	if p.peek() != '(' || len(p.input)-p.inputPos < 4 {
		return false
	}
	return p.at('(', scanner.Int, ')') && p.input[p.inputPos+3].ScanToken.IsRegisterShift()
}

// registerReference parses a register given either the name, R10, or a parenthesized form, SPR(10).
func (p *Parser) registerReference(name string) (obj.RBaseType, bool) {
	r, present := p.arch.Register[name]
	if present {
		return r, true
	}
	if !p.arch.RegisterPrefix[name] {
		p.errorf("expected register; found %s", name)
		return 0, false
	}
	p.get('(')
	tok := p.get(scanner.Int)
	num, err := strconv.ParseInt(tok.String(), 10, 16)
	p.get(')')
	if err != nil {
		p.errorf("parsing register list: %s", err)
		return 0, false
	}
	r, ok := p.arch.RegisterNumber(name, obj.RBaseType(num))
	if !ok {
		p.errorf("illegal register %s(%d)", name, r)
		return 0, false
	}
	return r, true
}

// register parses a full register reference where there is no symbol present (as in 4(R0) or R(10) but not sym(SB))
// including forms involving multiple registers such as R1:R2.
func (p *Parser) register(name string, prefix rune) (r1, r2 obj.RBaseType, scale int8, ok bool) {
	// R1 or R(1) R1:R2 R1,R2 R1+R2, or R1*scale.
	r1, ok = p.registerReference(name)
	if !ok {
		return
	}
	if prefix != 0 && prefix != '*' { // *AX is OK.
		p.errorf("prefix %c not allowed for register: %c%s", prefix, prefix, name)
	}
	c := p.peek()
	if c == ':' || c == ',' || c == '+' {
		// 2nd register; syntax (R1+R2) etc. No two architectures agree.
		// Check the architectures match the syntax.
		char := p.arch.LinkArch.Thechar
		switch p.next().ScanToken {
		case ',':
			if char != '5' && char != '7' {
				p.errorf("(register,register) not supported on this architecture")
				return
			}
		case '+':
			if char != '9' {
				p.errorf("(register+register) not supported on this architecture")
				return
			}
		}
		name := p.next().String()
		r2, ok = p.registerReference(name)
		if !ok {
			return
		}
	}
	if p.peek() == '*' {
		// Scale
		p.next()

		switch s := p.next().String(); s {
		case "1", "2", "4", "8":
			scale = int8(s[0] - '0')
		default:
			p.errorf("bad scale: %s", s)
			scale = 0
		}
	}
	return r1, r2, scale, true
}

// registerShift parses an ARM shifted register reference and returns the encoded representation.
// There is known to be a register (current token) and a shift operator (peeked token).
func (p *Parser) registerShift(name string, prefix rune) int64 {
	if prefix != 0 {
		p.errorf("prefix %c not allowed for shifted register: $%s", prefix, name)
	}
	// R1 op R2 or r1 op constant.
	// op is:
	//	"<<" == 0
	//	">>" == 1
	//	"->" == 2
	//	"@>" == 3
	r1, ok := p.registerReference(name)
	if !ok {
		return 0
	}
	var op int16
	switch p.next().ScanToken {
	case lex.LSH:
		op = 0
	case lex.RSH:
		op = 1
	case lex.ARR:
		op = 2
	case lex.ROT:
		op = 3
	}
	tok := p.next()
	str := tok.String()
	var count int16
	switch tok.ScanToken {
	case scanner.Ident:
		r2, ok := p.registerReference(str)
		if !ok {
			p.errorf("rhs of shift must be register or integer: %s", str)
		}
		count = (int16(r2)&15)<<8 | 1<<4
	case scanner.Int, '(':
		p.back()
		x := int64(p.expr())
		if x >= 32 {
			p.errorf("register shift count too large: %s", str)
		}
		count = int16((x & 31) << 7)
	default:
		p.errorf("unexpected %s in register shift", tok.String())
	}
	return int64((int16(r1) & 15) | op<<5 | count)
}

// setPseudoRegister sets the NAME field of addr for a pseudo-register reference such as (SB).
func (p *Parser) setPseudoRegister(addr *obj.Addr, reg string, isStatic bool, prefix rune) {
	if addr.Reg != 0 {
		p.errorf("internal error: reg %s already set in pseudo", reg)
	}
	switch reg {
	case "FP":
		addr.Name = obj.NAME_PARAM
	case "PC":
		if prefix != 0 {
			p.errorf("illegal addressing mode for PC")
		}
		addr.Type = obj.TYPE_BRANCH // We set the type and leave NAME untouched. See asmJump.
	case "SB":
		addr.Name = obj.NAME_EXTERN
		if isStatic {
			addr.Name = obj.NAME_STATIC
		}
	case "SP":
		addr.Name = obj.NAME_AUTO // The pseudo-stack.
	default:
		p.errorf("expected pseudo-register; found %s", reg)
	}
	if prefix == '$' {
		addr.Type = obj.TYPE_ADDR
	}
}

// registerList parses an ARM register list expression, a list of registers in [].
// There may be comma-separated ranges or individual registers, as in
// [R1,R3-R5]. Only R0 through R15 may appear.
// The opening bracket has been consumed.
func (p *Parser) registerList(a *obj.Addr) {
	// One range per loop.
	var bits uint16
	for {
		tok := p.next()
		if tok.ScanToken == ']' {
			break
		}
		lo := p.registerNumber(tok.String())
		hi := lo
		if p.peek() == '-' {
			p.next()
			hi = p.registerNumber(p.next().String())
		}
		if hi < lo {
			lo, hi = hi, lo
		}
		for lo <= hi {
			if bits&(1<<lo) != 0 {
				p.errorf("register R%d already in list", lo)
			}
			bits |= 1 << lo
			lo++
		}
		if p.peek() != ']' {
			p.get(',')
		}
	}
	a.Type = obj.TYPE_REGLIST
	a.Offset = int64(bits)
}

// register number is ARM-specific. It returns the number of the specified register.
func (p *Parser) registerNumber(name string) uint16 {
	if p.arch.LinkArch.Thechar == '5' && name == "g" {
		return 10
	}
	if name[0] != 'R' {
		p.errorf("expected g or R0 through R15; found %s", name)
	}
	r, ok := p.registerReference(name)
	if !ok {
		return 0
	}
	return uint16(r - p.arch.Register["R0"])
}

// registerIndirect parses the general form of a register indirection.
// It is can be (R1), (R2*scale), or (R1)(R2*scale) where R1 may be a simple
// register or register pair R:R or (R, R) or (R+R).
// Or it might be a pseudo-indirection like (FP).
// We are sitting on the opening parenthesis.
func (p *Parser) registerIndirect(a *obj.Addr, prefix rune) {
	p.get('(')
	tok := p.next()
	name := tok.String()
	r1, r2, scale, ok := p.register(name, 0)
	if !ok {
		p.errorf("indirect through non-register %s", tok)
	}
	p.get(')')
	a.Type = obj.TYPE_MEM
	if r1 < 0 {
		// Pseudo-register reference.
		if r2 != 0 {
			p.errorf("cannot use pseudo-register in pair")
			return
		}
		// For SB, SP, and FP, there must be a name here. 0(FP) is not legal.
		if name != "PC" && a.Name == obj.NAME_NONE {
			p.errorf("cannot reference %s without a symbol", name)
		}
		p.setPseudoRegister(a, name, false, prefix)
		return
	}
	a.Reg = r1
	if r2 != 0 {
		// TODO: Consistency in the encoding would be nice here.
		if p.arch.LinkArch.Thechar == '5' || p.arch.LinkArch.Thechar == '7' {
			// Special form
			// ARM: destination register pair (R1, R2).
			// ARM64: register pair (R1, R2) for LDP/STP.
			if prefix != 0 || scale != 0 {
				p.errorf("illegal address mode for register pair")
				return
			}
			a.Type = obj.TYPE_REGREG
			a.Offset = int64(r2)
			// Nothing may follow
			return
		}
	}
	if r2 != 0 {
		p.errorf("indirect through register pair")
	}
	if prefix == '$' {
		a.Type = obj.TYPE_ADDR
	}
	if r1 == arch.RPC && prefix != 0 {
		p.errorf("illegal addressing mode for PC")
	}
	if scale == 0 && p.peek() == '(' {
		// General form (R)(R*scale).
		p.next()
		tok := p.next()
		r1, r2, scale, ok = p.register(tok.String(), 0)
		if !ok {
			p.errorf("indirect through non-register %s", tok)
		}
		if r2 != 0 {
			p.errorf("unimplemented two-register form")
		}
		a.Index = int16(r1)
		a.Scale = int16(scale)
		p.get(')')
	} else if scale != 0 {
		// First (R) was missing, all we have is (R*scale).
		a.Reg = 0
		a.Index = int16(r1)
		a.Scale = int16(scale)
	}
}
