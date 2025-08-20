// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

// 指令的名字
// 保持和指令定义相同的顺序
var Anames = []string{
	// 各平台通用的伪指令
	AXXX:   "XXX",
	AGLOBL: "GLOBL",
	ADATA:  "DATA",
	ATEXT:  "TEXT",
	ACALL:  "CALL",
	ARET:   "RET",
	AJMP:   "JMP",
	ANOP:   "NOP",

	//
	// Unprivileged ISA (version 20240411)
	//

	// 2.4: Integer Computational Instructions (RV32I)
	AADDI:  "ADDI",
	ASLTI:  "SLTI",
	ASLTIU: "SLTIU",
	AANDI:  "ANDI",
	AORI:   "ORI",
	AXORI:  "XORI",
	ASLLI:  "SLLI",
	ASRLI:  "SRLI",
	ASRAI:  "SRAI",
	ALUI:   "LUI",
	AAUIPC: "AUIPC",
	AADD:   "ADD",
	ASLT:   "SLT",
	ASLTU:  "SLTU",
	AAND:   "AND",
	AOR:    "OR",
	AXOR:   "XOR",
	ASLL:   "SLL",
	ASRL:   "SRL",
	ASUB:   "SUB",
	ASRA:   "SRA",

	// 2.5: Control Transfer Instructions (RV32I)
	AJAL:  "JAL",
	AJALR: "JALR",
	ABEQ:  "BEQ",
	ABNE:  "BNE",
	ABLT:  "BLT",
	ABLTU: "BLTU",
	ABGE:  "BGE",
	ABGEU: "BGEU",

	// 2.6: Load and Store Instructions (RV32I)
	ALW:  "LW",
	ALH:  "LH",
	ALHU: "LHU",
	ALB:  "LB",
	ALBU: "LBU",
	ASW:  "SW",
	ASH:  "SH",
	ASB:  "SB",

	// 2.7: Memory Ordering Instructions (RV32I)
	AFENCE: "FENCE",

	// 4.2: Integer Computational Instructions (RV64I)
	AADDIW: "ADDIW",
	ASLLIW: "SLLIW",
	ASRLIW: "SRLIW",
	ASRAIW: "SRAIW",
	AADDW:  "ADDW",
	ASLLW:  "SLLW",
	ASRLW:  "SRLW",
	ASUBW:  "SUBW",
	ASRAW:  "SRAW",

	// 4.3: Load and Store Instructions (RV64I)
	ALWU: "LWU",
	ALD:  "LD",
	ASD:  "SD",

	// 7.1: CSR Instructions (Zicsr)
	ACSRRW:  "CSRRW",
	ACSRRS:  "CSRRS",
	ACSRRC:  "CSRRC",
	ACSRRWI: "CSRRWI",
	ACSRRSI: "CSRRSI",
	ACSRRCI: "CSRRCI",

	// 13.1: Multiplication Operations (RV32M/RV64M)
	AMUL:    "MUL",
	AMULH:   "MULH",
	AMULHU:  "MULHU",
	AMULHSU: "MULHSU",
	AMULW:   "MULW", // RV64M

	// 13.2: Division Operations (RV32M/RV64M)
	ADIV:   "DIV",
	ADIVU:  "DIVU",
	AREM:   "REM",
	AREMU:  "REMU",
	ADIVW:  "DIVW",  // RV64M
	ADIVUW: "DIVUW", // RV64M
	AREMW:  "REMW",  // RV64M
	AREMUW: "REMUW", // RV64M

	// 20.5: Single-Precision Load and Store Instructions (F)
	AFLW: "FLW",
	AFSW: "FSW",

	// 20.6: Single-Precision Floating-Point Computational Instructions
	AFADDS:   "FADDS",
	AFSUBS:   "FSUBS",
	AFMULS:   "FMULS",
	AFDIVS:   "FDIVS",
	AFSQRTS:  "FSQRTS",
	AFMINS:   "FMINS",
	AFMAXS:   "FMAXS",
	AFMADDS:  "FMADDS",
	AFMSUBS:  "FMSUBS",
	AFNMADDS: "FNMADDS",
	AFNMSUBS: "FNMSUBS",

	// 20.7: Single-Precision Floating-Point Conversion and Move Instructions
	AFCVTWS:  "FCVTWS",
	AFCVTLS:  "FCVTLS",
	AFCVTSW:  "FCVTSW",
	AFCVTSL:  "FCVTSL",
	AFCVTWUS: "FCVTWUS",
	AFCVTLUS: "FCVTLUS",
	AFCVTSWU: "FCVTSWU",
	AFCVTSLU: "FCVTSLU",
	AFSGNJS:  "FSGNJS",
	AFSGNJNS: "FSGNJNS",
	AFSGNJXS: "FSGNJXS",
	AFMVXS:   "FMVXS",
	AFMVSX:   "FMVSX",
	AFMVXW:   "FMVXW",
	AFMVWX:   "FMVWX",

	// 20.8: Single-Precision Floating-Point Compare Instructions
	AFEQS: "FEQS",
	AFLTS: "FLTS",
	AFLES: "FLES",

	// 20.9: Single-Precision Floating-Point Classify Instruction
	AFCLASSS: "FCLASSS",

	// 21.3: Double-Precision Load and Store Instructions (D)
	AFLD: "FLD",
	AFSD: "FSD",

	// 21.4: Double-Precision Floating-Point Computational Instructions
	AFADDD:   "FADDD",
	AFSUBD:   "FSUBD",
	AFMULD:   "FMULD",
	AFDIVD:   "FDIVD",
	AFMIND:   "FMIND",
	AFMAXD:   "FMAXD",
	AFSQRTD:  "FSQRTD",
	AFMADDD:  "FMADDD",
	AFMSUBD:  "FMSUBD",
	AFNMADDD: "FNMADDD",
	AFNMSUBD: "FNMSUBD",

	// 21.5: Double-Precision Floating-Point Conversion and Move Instructions
	AFCVTWD:  "FCVTWD",
	AFCVTLD:  "FCVTLD",
	AFCVTDW:  "FCVTDW",
	AFCVTDL:  "FCVTDL",
	AFCVTWUD: "FCVTWUD",
	AFCVTLUD: "FCVTLUD",
	AFCVTDWU: "FCVTDWU",
	AFCVTDLU: "FCVTDLU",
	AFCVTSD:  "FCVTSD",
	AFCVTDS:  "FCVTDS",
	AFSGNJD:  "FSGNJD",
	AFSGNJND: "FSGNJND",
	AFSGNJXD: "FSGNJXD",
	AFMVXD:   "FMVXD",
	AFMVDX:   "FMVDX",

	// 21.6: Double-Precision Floating-Point Compare Instructions
	AFEQD: "FEQD",
	AFLTD: "FLTD",
	AFLED: "FLED",

	// 21.7: Double-Precision Floating-Point Classify Instruction
	AFCLASSD: "FCLASSD",

	//
	// Privileged ISA (version 20240411)
	//

	// 3.3.1: Environment Call and Breakpoint
	AECALL:  "ECALL",
	ASCALL:  "SCALL",
	AEBREAK: "EBREAK",
	ASBREAK: "SBREAK",

	// 3.3.2: Trap-Return Instructions
	AMRET: "MRET",
	ASRET: "SRET",
	ADRET: "DRET",

	// 3.3.3: Wait for Interrupt
	AWFI: "WFI",

	// End marker
	ALAST: "LAST",
}
