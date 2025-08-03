// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"wa-lang.org/wa/internal/p9asm/asm/arch"
	"wa-lang.org/wa/internal/p9asm/asm/lex"
	"wa-lang.org/wa/internal/p9asm/obj"
)

func (p *Parser) instruction(op obj.As, word, cond string, operands [][]lex.Token) {
	p.addr = p.addr[0:0]
	isJump := p.arch.IsJump(word)
	for _, op := range operands {
		addr := p.address(op)
		if !isJump && addr.Reg < 0 { // Jumps refer to PC, a pseudo.
			p.errorf("illegal use of pseudo-register in %s", word)
		}
		p.addr = append(p.addr, addr)
	}
	if isJump {
		p.asmJump(op, cond, p.addr)
		return
	}
	p.asmInstruction(op, cond, p.addr)
}

// asmJump assembles a jump instruction.
// JMP	R1
// JMP	exit
// JMP	3(PC)
func (p *Parser) asmJump(op obj.As, cond string, a []obj.Addr) {
	var target *obj.Addr
	prog := &obj.Prog{
		Ctxt:   p.ctxt,
		Lineno: p.histLineNum,
		As:     op,
	}
	switch len(a) {
	case 1:
		target = &a[0]
	case 2:
		// Special 2-operand jumps.
		target = &a[1]
		prog.From = a[0]
	case 3:
		if p.arch.LinkArch.Thechar == '9' {
			// Special 3-operand jumps.
			// First two must be constants; a[1] is a register number.
			target = &a[2]
			prog.From = obj.Addr{
				Type:   obj.TYPE_CONST,
				Offset: p.getConstant(prog, op, &a[0]),
			}
			reg := obj.RBaseType(p.getConstant(prog, op, &a[1]))
			reg, ok := p.arch.RegisterNumber("R", reg)
			if !ok {
				p.errorf("bad register number %d", reg)
			}
			prog.Reg = reg
			break
		}
		fallthrough
	default:
		p.errorf("wrong number of arguments to %v instruction", op)
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
			prog.To = obj.Addr{
				Type: obj.TYPE_BRANCH,
				Val:  targetProg,
			}
		}
	case target.Type == obj.TYPE_MEM && target.Name == obj.NAME_NONE:
		// JMP 4(R0)
		prog.To = *target
		// On the ppc64, 9a encodes BR (CTR) as BR CTR. We do the same.
		if p.arch.LinkArch.Thechar == '9' && target.Offset == 0 {
			prog.To.Type = obj.TYPE_REG
		}
	case target.Type == obj.TYPE_CONST:
		// JMP $4
		prog.To = a[0]
	default:
		p.errorf("cannot assemble jump %+v", target)
	}

	p.append(prog, cond, true)
}

// asmInstruction assembles an instruction.
// MOVW R9, (R10)
func (p *Parser) asmInstruction(op obj.As, cond string, a []obj.Addr) {
	// fmt.Printf("%s %+v\n", obj.Aconv(op), a)
	prog := &obj.Prog{
		Ctxt:   p.ctxt,
		Lineno: p.histLineNum,
		As:     op,
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
		if p.arch.LinkArch.Thechar == '5' {
			if arch.IsARMCMP(op) {
				prog.From = a[0]
				prog.Reg = p.getRegister(prog, op, &a[1])
				break
			}
			// Strange special cases.
			if arch.IsARMSTREX(op) {
				/*
					STREX x, (y)
						from=(y) reg=x to=x
					STREX (x), y
						from=(x) reg=y to=y
				*/
				if a[0].Type == obj.TYPE_REG && a[1].Type != obj.TYPE_REG {
					prog.From = a[1]
					prog.Reg = a[0].Reg
					prog.To = a[0]
					break
				} else if a[0].Type != obj.TYPE_REG && a[1].Type == obj.TYPE_REG {
					prog.From = a[0]
					prog.Reg = a[1].Reg
					prog.To = a[1]
					break
				}
				p.errorf("unrecognized addressing for %s", op)
			}
		} else if p.arch.LinkArch.Thechar == '7' && arch.IsARM64CMP(op) {
			prog.From = a[0]
			prog.Reg = p.getRegister(prog, op, &a[1])
			break
		}
		prog.From = a[0]
		prog.To = a[1]
	case 3:
		switch p.arch.LinkArch.Thechar {
		case '5':
			// Special cases.
			if arch.IsARMSTREX(op) {
				/*
					STREX x, (y), z
						from=(y) reg=x to=z
				*/
				prog.From = a[1]
				prog.Reg = p.getRegister(prog, op, &a[0])
				prog.To = a[2]
				break
			}
			// Otherwise the 2nd operand (a[1]) must be a register.
			prog.From = a[0]
			prog.Reg = p.getRegister(prog, op, &a[1])
			prog.To = a[2]
		case '7':
			// ARM64 instructions with one input and two outputs.
			if arch.IsARM64STLXR(op) {
				prog.From = a[0]
				prog.To = a[1]
				if a[2].Type != obj.TYPE_REG {
					p.errorf("invalid addressing modes for third operand to %s instruction, must be register", op)
				}
				prog.RegTo2 = a[2].Reg
				break
			}
			prog.From = a[0]
			prog.Reg = p.getRegister(prog, op, &a[1])
			prog.To = a[2]
		case '6', '8':
			prog.From = a[0]
			prog.From3 = new(obj.Addr)
			*prog.From3 = a[1]
			prog.To = a[2]
		case '9':
			// Arithmetic. Choices are:
			// reg reg reg
			// imm reg reg
			// reg imm reg
			// If the immediate is the middle argument, use From3.
			switch a[1].Type {
			case obj.TYPE_REG:
				prog.From = a[0]
				prog.Reg = p.getRegister(prog, op, &a[1])
				prog.To = a[2]
			case obj.TYPE_CONST:
				prog.From = a[0]
				prog.From3 = new(obj.Addr)
				*prog.From3 = a[1]
				prog.To = a[2]
			default:
				p.errorf("invalid addressing modes for %s instruction", op)
			}
		default:
			p.errorf("TODO: implement three-operand instructions for this architecture")
		}
	case 4:
		if p.arch.LinkArch.Thechar == '5' && arch.IsARMMULA(op) {
			// All must be registers.
			p.getRegister(prog, op, &a[0])
			r1 := p.getRegister(prog, op, &a[1])
			p.getRegister(prog, op, &a[2])
			r3 := p.getRegister(prog, op, &a[3])
			prog.From = a[0]
			prog.To = a[2]
			prog.To.Type = obj.TYPE_REGREG2
			prog.To.Offset = int64(r3)
			prog.Reg = r1
			break
		}
		if p.arch.LinkArch.Thechar == '7' {
			prog.From = a[0]
			prog.Reg = p.getRegister(prog, op, &a[1])
			prog.From3 = new(obj.Addr)
			*prog.From3 = a[2]
			prog.To = a[3]
			break
		}
		p.errorf("can't handle %s instruction with 4 operands", op)
	case 5:
		p.errorf("can't handle %s instruction with 5 operands", op)
	case 6:
		if p.arch.LinkArch.Thechar == '5' && arch.IsARMMRC(op) {
			// Strange special case: MCR, MRC.
			prog.To.Type = obj.TYPE_CONST
			x0 := p.getConstant(prog, op, &a[0])
			x1 := p.getConstant(prog, op, &a[1])
			x2 := int64(p.getRegister(prog, op, &a[2]))
			x3 := int64(p.getRegister(prog, op, &a[3]))
			x4 := int64(p.getRegister(prog, op, &a[4]))
			x5 := p.getConstant(prog, op, &a[5])
			// Cond is handled specially for this instruction.
			offset, MRC, ok := arch.ARMMRCOffset(op, cond, x0, x1, x2, x3, x4, x5)
			if !ok {
				p.errorf("unrecognized condition code .%q", cond)
			}
			prog.To.Offset = offset
			cond = ""
			prog.As = MRC // Both instructions are coded as MRC.
			break
		}
		fallthrough
	default:
		p.errorf("can't handle %s instruction with %d operands", op, len(a))
	}

	p.append(prog, cond, true)
}

// getRegister checks that addr represents a register and returns its value.
func (p *Parser) getRegister(prog *obj.Prog, op obj.As, addr *obj.Addr) obj.RBaseType {
	if addr.Type != obj.TYPE_REG || addr.Offset != 0 || addr.Name != 0 || addr.Index != 0 {
		p.errorf("%s: expected register; found %s", op, addr.Dconv(prog))
	}
	return addr.Reg
}

// getConstant checks that addr represents a plain constant and returns its value.
func (p *Parser) getConstant(prog *obj.Prog, op obj.As, addr *obj.Addr) int64 {
	if addr.Type != obj.TYPE_MEM || addr.Name != 0 || addr.Reg != 0 || addr.Index != 0 {
		p.errorf("%s: expected integer constant; found %s", op, addr.Dconv(prog))
	}
	return addr.Offset
}
