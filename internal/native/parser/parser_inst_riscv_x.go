// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/riscv"
	"wa-lang.org/wa/internal/native/token"
)

func (p *parser) parseInst_riscv() (inst ast.Instruction) {
	assert(p.cpu == abi.RISCV64 || p.cpu == abi.RISCV32)

	inst.Pos = p.pos
	inst.As = p.parseAs()
	p.acceptToken(p.tok)
	inst.Arg = new(abi.AsArgument)

	switch inst.As {
	default:
		p.errorf(p.pos, "%v is not riscv instruction", p.tok)

	// 2.4: Integer Computational Instructions (RV32I)
	case riscv.AADDI:
		// addi a0, a0, %pcrel_lo(_start) # PC相对地址的低12bit
		// addi a0, a0, %lo(UART0) # 绝对地址的低12bit
		// addi a0, a0, 1
		// addi t1, t1, 0x555
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_immAddr(&inst)
		return inst

	case riscv.ASLTI:
		panic("TODO")
	case riscv.ASLTIU:
		panic("TODO")
	case riscv.AANDI:
		panic("TODO")
	case riscv.AORI:
		panic("TODO")
	case riscv.AXORI:
		panic("TODO")
	case riscv.ASLLI:
		panic("TODO")
	case riscv.ASRLI:
		panic("TODO")
	case riscv.ASRAI:
		panic("TODO")
	case riscv.ALUI:
		// lui t0, %hi(UART0) # UART0 高20位
		// lui t1, 0x5 # 高 20 位 (0x5 << 12 = 0x5000)
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_immAddr(&inst)
		return inst
	case riscv.AAUIPC:
		// auipc a0, %pcrel_hi(message) # 高20位 = 当前PC + 偏移
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_immAddr(&inst)
	case riscv.AADD:
		panic("TODO")
	case riscv.ASLT:
		panic("TODO")
	case riscv.ASLTU:
		panic("TODO")
	case riscv.AAND:
		panic("TODO")
	case riscv.AOR:
		panic("TODO")
	case riscv.AXOR:
		panic("TODO")
	case riscv.ASLL:
		panic("TODO")
	case riscv.ASRL:
		panic("TODO")
	case riscv.ASUB:
		panic("TODO")
	case riscv.ASRA:
		panic("TODO")

	// 2.5: Control Transfer Instructions (RV32I)
	case riscv.AJAL:
		// jal x0, print_loop
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_immAddr(&inst)
		return inst
	case riscv.AJALR:
		panic("TODO")
	case riscv.ABEQ:
		// beq a1, x0, finished
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_relAddr(&inst)
		return inst
	case riscv.ABNE:
		panic("TODO")
	case riscv.ABLT:
		panic("TODO")
	case riscv.ABLTU:
		panic("TODO")
	case riscv.ABGE:
		panic("TODO")
	case riscv.ABGEU:
		panic("TODO")

	// 2.6: Load and Store Instructions (RV32I)
	case riscv.ALW:
		panic("TODO")
	case riscv.ALH:
		panic("TODO")
	case riscv.ALHU:
		panic("TODO")
	case riscv.ALB:
		panic("TODO")
	case riscv.ALBU:
		// lbu a1, 0(a0) # 取一个字节
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		reg, off := p.parseInst_riscv_baseOffset()
		inst.Arg.Rs1 = reg
		inst.Arg.Imm = off
		return inst
	case riscv.ASW:
		// sw t1, 0(t0)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		reg, off := p.parseInst_riscv_baseOffset()
		inst.Arg.Rs2 = reg
		inst.Arg.Imm = off
		return inst
	case riscv.ASH:
		panic("TODO")
	case riscv.ASB:
		// sb a1, 0(t0) # 写到UART寄存器
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		reg, off := p.parseInst_riscv_baseOffset()
		inst.Arg.Rs2 = reg
		inst.Arg.Imm = off
		return inst

	// 2.7: Memory Ordering Instructions (RV32I)
	case riscv.AFENCE:
		panic("TODO")

	// 3.3.1: Environment Call and Breakpoint
	case riscv.AECALL:
		panic("TODO")
	case riscv.AEBREAK:
		panic("TODO")

	// 4.2: Integer Computational Instructions (RV64I)
	case riscv.AADDIW:
		panic("TODO")
	case riscv.ASLLIW:
		panic("TODO")
	case riscv.ASRLIW:
		panic("TODO")
	case riscv.ASRAIW:
		panic("TODO")
	case riscv.AADDW:
		panic("TODO")
	case riscv.ASLLW:
		panic("TODO")
	case riscv.ASRLW:
		panic("TODO")
	case riscv.ASUBW:
		panic("TODO")
	case riscv.ASRAW:
		panic("TODO")

	// 4.3: Load and Store Instructions (RV64I)
	case riscv.ALWU:
		panic("TODO")
	case riscv.ALD:
		panic("TODO")
	case riscv.ASD:
		panic("TODO")

	// 7.1: CSR Instructions (Zicsr)
	case riscv.ACSRRW:
		panic("TODO")
	case riscv.ACSRRS:
		panic("TODO")
	case riscv.ACSRRC:
		panic("TODO")
	case riscv.ACSRRWI:
		panic("TODO")
	case riscv.ACSRRSI:
		panic("TODO")
	case riscv.ACSRRCI:
		panic("TODO")

	// 13.1: Multiplication Operations (RV32M/RV64M)
	case riscv.AMUL:
		panic("TODO")
	case riscv.AMULH:
		panic("TODO")
	case riscv.AMULHU:
		panic("TODO")
	case riscv.AMULHSU:
		panic("TODO")
	case riscv.AMULW: // RV64M
		panic("TODO")

	// 13.2: Division Operations (RV32M/RV64M)
	case riscv.ADIV:
		panic("TODO")
	case riscv.ADIVU:
		panic("TODO")
	case riscv.AREM:
		panic("TODO")
	case riscv.AREMU:
		panic("TODO")
	case riscv.ADIVW:
		panic("TODO") // RV64M
	case riscv.ADIVUW:
		panic("TODO") // RV64M
	case riscv.AREMW:
		panic("TODO") // RV64M
	case riscv.AREMUW:
		panic("TODO") // RV64M

	// 20.5: Single-Precision Load and Store Instructions (F)
	case riscv.AFLW:
		panic("TODO")
	case riscv.AFSW:
		panic("TODO")

	// 20.6: Single-Precision Floating-Point Computational Instructions
	case riscv.AFADD_S:
		panic("TODO")
	case riscv.AFSUB_S:
		panic("TODO")
	case riscv.AFMUL_S:
		panic("TODO")
	case riscv.AFDIV_S:
		panic("TODO")
	case riscv.AFMIN_S:
		panic("TODO")
	case riscv.AFMAX_S:
		panic("TODO")
	case riscv.AFSQRT_S:
		panic("TODO")
	case riscv.AFMADD_S:
		panic("TODO")
	case riscv.AFMSUB_S:
		panic("TODO")
	case riscv.AFNMADD_S:
		panic("TODO")
	case riscv.AFNMSUB_S:
		panic("TODO")

	// 20.7: Single-Precision Floating-Point Conversion and Move Instructions
	case riscv.AFCVT_W_S:
		panic("TODO")
	case riscv.AFCVT_L_S:
		panic("TODO")
	case riscv.AFCVT_S_W:
		panic("TODO")
	case riscv.AFCVT_S_L:
		panic("TODO")
	case riscv.AFCVT_WU_S:
		panic("TODO")
	case riscv.AFCVT_LU_S:
		panic("TODO")
	case riscv.AFCVT_S_WU:
		panic("TODO")
	case riscv.AFCVT_S_LU:
		panic("TODO")
	case riscv.AFSGNJ_S:
		panic("TODO")
	case riscv.AFSGNJN_S:
		panic("TODO")
	case riscv.AFSGNJX_S:
		panic("TODO")
	case riscv.AFMV_X_W:
		panic("TODO")
	case riscv.AFMV_W_X:
		panic("TODO")

	// 20.8: Single-Precision Floating-Point Compare Instructions
	case riscv.AFEQ_S:
		panic("TODO")
	case riscv.AFLT_S:
		panic("TODO")
	case riscv.AFLE_S:
		panic("TODO")

	// 20.9: Single-Precision Floating-Point Classify Instruction
	case riscv.AFCLASS_S:
		panic("TODO")

	// 21.3: Double-Precision Load and Store Instructions (D)
	case riscv.AFLD:
		panic("TODO")
	case riscv.AFSD:
		panic("TODO")

	// 21.4: Double-Precision Floating-Point Computational Instructions
	case riscv.AFADD_D:
		panic("TODO")
	case riscv.AFSUB_D:
		panic("TODO")
	case riscv.AFMUL_D:
		panic("TODO")
	case riscv.AFDIV_D:
		panic("TODO")
	case riscv.AFMIN_D:
		panic("TODO")
	case riscv.AFMAX_D:
		panic("TODO")
	case riscv.AFSQRT_D:
		panic("TODO")
	case riscv.AFMADD_D:
		panic("TODO")
	case riscv.AFMSUB_D:
		panic("TODO")
	case riscv.AFNMADD_D:
		panic("TODO")
	case riscv.AFNMSUB_D:
		panic("TODO")

	// 21.5: Double-Precision Floating-Point Conversion and Move Instructions
	case riscv.AFCVT_W_D:
		panic("TODO")
	case riscv.AFCVT_L_D:
		panic("TODO")
	case riscv.AFCVT_D_W:
		panic("TODO")
	case riscv.AFCVT_D_L:
		panic("TODO")
	case riscv.AFCVT_WU_D:
		panic("TODO")
	case riscv.AFCVT_LU_D:
		panic("TODO")
	case riscv.AFCVT_D_WU:
		panic("TODO")
	case riscv.AFCVT_D_LU:
		panic("TODO")
	case riscv.AFCVT_S_D:
		panic("TODO")
	case riscv.AFCVT_D_S:
		panic("TODO")
	case riscv.AFSGNJ_D:
		panic("TODO")
	case riscv.AFSGNJN_D:
		panic("TODO")
	case riscv.AFSGNJX_D:
		panic("TODO")
	case riscv.AFMV_X_D:
		panic("TODO")
	case riscv.AFMV_D_X:
		panic("TODO")

	// 21.6: Double-Precision Floating-Point Compare Instructions
	case riscv.AFEQ_D:
		panic("TODO")
	case riscv.AFLT_D:
		panic("TODO")
	case riscv.AFLE_D:
		panic("TODO")

	// 21.7: Double-Precision Floating-Point Classify Instruction
	case riscv.AFCLASS_D:

	// 伪指令(A_开头以区分)
	// ISA (version 20191213)
	// 25: RISC-V Assembly Programmer's Handbook
	// 只保留可以1:1映射到原生指令的类型
	// 长地址跳转需要用户手动处理

	case riscv.A_NOP:
		panic("TODO")
	case riscv.A_MV:
		panic("TODO")
	case riscv.A_NOT:
		panic("TODO")
	case riscv.A_NEG:
		panic("TODO")
	case riscv.A_NEGW:
		panic("TODO")
	case riscv.A_SEXT_W:
		panic("TODO")
	case riscv.A_SEQZ:
		panic("TODO")
	case riscv.A_SNEZ:
		panic("TODO")
	case riscv.A_SLTZ:
		panic("TODO")
	case riscv.A_SGTZ:
		panic("TODO")
	case riscv.A_FMV_S:
		panic("TODO")
	case riscv.A_FABS_S:
		panic("TODO")
	case riscv.A_FNEG_S:
		panic("TODO")
	case riscv.A_FMV_D:
		panic("TODO")
	case riscv.A_FABS_D:
		panic("TODO")
	case riscv.A_FNEG_D:
		panic("TODO")
	case riscv.A_BEQZ:
		panic("TODO")
	case riscv.A_BNEZ:
		panic("TODO")
	case riscv.A_BLEZ:
		panic("TODO")
	case riscv.A_BGEZ:
		panic("TODO")
	case riscv.A_BLTZ:
		panic("TODO")
	case riscv.A_BGTZ:
		panic("TODO")
	case riscv.A_BGT:
		panic("TODO")
	case riscv.A_BLE:
		panic("TODO")
	case riscv.A_BGTU:
		panic("TODO")
	case riscv.A_BLEU:
		panic("TODO")
	case riscv.A_J:
		panic("TODO")
	case riscv.A_JR:
		panic("TODO")
	case riscv.A_RET:
		panic("TODO")
	case riscv.A_RDINSTRET:
		panic("TODO")
	case riscv.A_RDCYCLE:
		panic("TODO")
	case riscv.A_RDTIME:
		panic("TODO")
	case riscv.A_CSRR:
		panic("TODO")
	case riscv.A_CSRW:
		panic("TODO")
	case riscv.A_CSRS:
		panic("TODO")
	case riscv.A_CSRC:
		panic("TODO")
	case riscv.A_CSRWI:
		panic("TODO")
	case riscv.A_CSRSI:
		panic("TODO")
	case riscv.A_CSRCI:
		panic("TODO")
	case riscv.A_FRCSR:
		panic("TODO")
	case riscv.A_FSCSR:
		panic("TODO")
	case riscv.A_FRRM:
		panic("TODO")
	case riscv.A_FSRM:
		panic("TODO")
	case riscv.A_FRFLAGS:
		panic("TODO")
	case riscv.A_FSFLAGS:
		panic("TODO")
	}

	panic("unreachable")
}

// 基于寄存器的内地地址解析
// 只能出现在 I-type 和 S-type 这两类指令中
// (t0)
// 4(t0)
// -4(t0)
func (p *parser) parseInst_riscv_baseOffset() (reg abi.RegType, offset int32) {
	if p.tok == token.INT {
		offset = p.parseInt32Lit()
	}
	p.acceptToken(token.LPAREN)
	reg = p.parseRegister()
	p.acceptToken(token.RPAREN)
	return
}

// 解析相对地址
// 只能是label或相对PC的数值
func (p *parser) parseInst_riscv_relAddr(inst *ast.Instruction) {
	switch p.tok {
	case token.IDENT:
		inst.Arg.Symbol = p.parseIdent()
		return
	case token.INT:
		inst.Arg.Imm = p.parseInt32Lit()
		return
	default:
		p.errorf(p.pos, "expect label or int, got %v", p.tok)
	}
	panic("unreachable")
}

// 解析地址立即数
// 不涉及寄存器解析
// 12
// _start
// %lo(_start)
func (p *parser) parseInst_riscv_immAddr(inst *ast.Instruction) {
	if p.tok == token.INT {
		inst.Arg.Imm = p.parseInt32Lit()
		return
	}

	if p.tok != token.IDENT {
		p.errorf(p.pos, "export IDENT, got %v", p.tok)
	}

	pos := p.pos
	symbolOrDecor := p.parseIdent()

	// 没有重定位修饰函数
	if p.tok != token.LPAREN {
		inst.Arg.Symbol = symbolOrDecor
		return
	}

	// 判断重定位修饰函数
	switch symbolOrDecor {
	case "%hi":
		inst.Arg.SymbolDecor = abi.BuiltinFn_HI
	case "%lo":
		inst.Arg.SymbolDecor = abi.BuiltinFn_LO
	case "%pcrel_hi":
		inst.Arg.SymbolDecor = abi.BuiltinFn_PCREL_HI
	case "%pcrel_lo":
		inst.Arg.SymbolDecor = abi.BuiltinFn_PCREL_LO
	default:
		p.errorf(pos, "unknow symbol decorator %s", symbolOrDecor)
	}

	p.acceptToken(token.LPAREN)
	inst.Arg.Symbol = p.parseIdent()
	p.acceptToken(token.RPAREN)
}
