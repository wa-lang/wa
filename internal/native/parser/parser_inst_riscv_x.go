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
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_immAddr(&inst)
		return inst
	case riscv.ASLTIU:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_immAddr(&inst)
		return inst
	case riscv.AANDI:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_immAddr(&inst)
		return inst
	case riscv.AORI:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_immAddr(&inst)
		return inst
	case riscv.AXORI:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_immAddr(&inst)
		return inst
	case riscv.ASLLI:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit() // shamt
		return inst
	case riscv.ASRLI:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit() // shamt
		return inst
	case riscv.ASRAI:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit() // shamt
		return inst
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
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.ASLT:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.ASLTU:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AAND:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AOR:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AXOR:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.ASLL:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.ASRL:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.ASUB:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.ASRA:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst

	// 2.5: Control Transfer Instructions (RV32I)
	case riscv.AJAL:
		// jal x0, print_loop
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_immAddr(&inst)
		return inst
	case riscv.AJALR:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst
	case riscv.ABEQ:
		// beq a1, x0, finished
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_relAddr(&inst)
		return inst
	case riscv.ABNE:
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_relAddr(&inst)
		return inst
	case riscv.ABLT:
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_relAddr(&inst)
		return inst
	case riscv.ABLTU:
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_relAddr(&inst)
		return inst
	case riscv.ABGE:
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_relAddr(&inst)
		return inst
	case riscv.ABGEU:
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_relAddr(&inst)
		return inst

	// 2.6: Load and Store Instructions (RV32I)
	case riscv.ALW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst
	case riscv.ALH:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst
	case riscv.ALHU:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst
	case riscv.ALB:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst
	case riscv.ALBU:
		// lbu a1, 0(a0) # 取一个字节
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst
	case riscv.ASW:
		// sw t1, 0(t0)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst
	case riscv.ASH:
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst
	case riscv.ASB:
		// sb a1, 0(t0) # 写到UART寄存器
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst

	// 2.7: Memory Ordering Instructions (RV32I)
	case riscv.AFENCE:
		return inst

	// 3.3.1: Environment Call and Breakpoint
	case riscv.AECALL:
		return inst

	case riscv.AEBREAK:
		return inst

	// 4.2: Integer Computational Instructions (RV64I)
	case riscv.AADDIW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_riscv_immAddr(&inst)
		return inst
	case riscv.ASLLIW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit() // shamt
		return inst
	case riscv.ASRLIW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit() // shamt
		return inst
	case riscv.ASRAIW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit() // shamt
		return inst
	case riscv.AADDW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.ASLLW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.ASRLW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.ASUBW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.ASRAW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst

	// 4.3: Load and Store Instructions (RV64I)
	case riscv.ALWU:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst
	case riscv.ALD:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst
	case riscv.ASD:
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst

	// 7.1: CSR Instructions (Zicsr)
	case riscv.ACSRRW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		return inst
	case riscv.ACSRRS:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		return inst
	case riscv.ACSRRC:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		return inst
	case riscv.ACSRRWI:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.ACSRRSI:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.ACSRRCI:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst

	// 13.1: Multiplication Operations (RV32M/RV64M)
	case riscv.AMUL:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AMULH:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AMULHU:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AMULHSU:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AMULW: // RV64M
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst

	// 13.2: Division Operations (RV32M/RV64M)
	case riscv.ADIV:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.ADIVU:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AREM:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AREMU:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.ADIVW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.ADIVUW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AREMW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AREMUW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst

	// 20.5: Single-Precision Load and Store Instructions (F)
	case riscv.AFLW:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst
	case riscv.AFSW:
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst

	// 20.6: Single-Precision Floating-Point Computational Instructions
	case riscv.AFADD_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFSUB_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFMUL_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFDIV_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFMIN_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFMAX_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFSQRT_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFMADD_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs3 = p.parseRegister()
		return inst
	case riscv.AFMSUB_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs3 = p.parseRegister()
		return inst
	case riscv.AFNMADD_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs3 = p.parseRegister()
		return inst
	case riscv.AFNMSUB_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs3 = p.parseRegister()
		return inst

	// 20.7: Single-Precision Floating-Point Conversion and Move Instructions
	case riscv.AFCVT_W_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_L_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_S_W:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_S_L:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_WU_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_LU_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_S_WU:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_S_LU:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFSGNJ_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFSGNJN_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFSGNJX_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFMV_X_W:
		inst.Arg.Rd = p.parseRegister()
		return inst
	case riscv.AFMV_W_X:
		inst.Arg.Rd = p.parseRegister()
		return inst

	// 20.8: Single-Precision Floating-Point Compare Instructions
	case riscv.AFEQ_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFLT_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFLE_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst

	// 20.9: Single-Precision Floating-Point Classify Instruction
	case riscv.AFCLASS_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst

	// 21.3: Double-Precision Load and Store Instructions (D)
	case riscv.AFLD:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst
	case riscv.AFSD:
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1, inst.Arg.Imm = p.parseInst_riscv_baseOffset()
		return inst

	// 21.4: Double-Precision Floating-Point Computational Instructions
	case riscv.AFADD_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFSUB_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFMUL_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFDIV_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFMIN_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFMAX_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFSQRT_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFMADD_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs3 = p.parseRegister()
		return inst
	case riscv.AFMSUB_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs3 = p.parseRegister()
		return inst
	case riscv.AFNMADD_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs3 = p.parseRegister()
		return inst
	case riscv.AFNMSUB_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs3 = p.parseRegister()
		return inst

	// 21.5: Double-Precision Floating-Point Conversion and Move Instructions
	case riscv.AFCVT_W_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_L_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_D_W:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_D_L:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_WU_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_LU_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_D_WU:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_D_LU:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_S_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFCVT_D_S:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFSGNJ_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFSGNJN_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFSGNJX_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFMV_X_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case riscv.AFMV_D_X:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst

	// 21.6: Double-Precision Floating-Point Compare Instructions
	case riscv.AFEQ_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFLT_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case riscv.AFLE_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst

	// 21.7: Double-Precision Floating-Point Classify Instruction
	case riscv.AFCLASS_D:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst

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
