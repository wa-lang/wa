// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

// 指令的名字
// 保持和指令定义相同的顺序
var Anames = []string{
	// 2.4: Integer Computational Instructions (RV32I)
	"ADDI",
	"SLTI",
	"SLTIU",
	"ANDI",
	"ORI",
	"XORI",
	"SLLI",
	"SRLI",
	"SRAI",
	"LUI",
	"AUIPC",
	"ADD",
	"SLT",
	"SLTU",
	"AND",
	"OR",
	"XOR",
	"SLL",
	"SRL",
	"SUB",
	"SRA",

	// 2.5: Control Transfer Instructions (RV32I)
	"JAL",
	"JALR",
	"BEQ",
	"BNE",
	"BLT",
	"BLTU",
	"BGE",
	"BGEU",

	// 2.6: Load and Store Instructions (RV32I)
	"LW",
	"LH",
	"LHU",
	"LB",
	"LBU",
	"SW",
	"SH",
	"SB",

	// 2.8: Environment Call and Breakpoint (RV32I)
	"ECALL",
	"EBREAK",

	// 5.2: Integer Computational Instructions (RV64I)
	"ADDIW",
	"SLLIW",
	"SRLIW",
	"SRAIW",
	"ADDW",
	"SLLW",
	"SRLW",
	"SUBW",
	"SRAW",

	// 5.3: Load and Store Instructions (RV64I)
	"LWU",
	"LD",
	"SD",

	// 7.1: Multiplication Operations (RV32M/RV64M)
	"MUL",
	"MULH",
	"MULHU",
	"MULHSU",
	"MULW", // RV64M
	"DIV",
	"DIVU",
	"REM",
	"REMU",
	"DIVW",  // RV64M
	"DIVUW", // RV64M
	"REMW",  // RV64M
	"REMUW", // RV64M

	// 9.1: CSR Instructions (Zicsr)
	"CSRRW",
	"CSRRS",
	"CSRRC",
	"CSRRWI",
	"CSRRSI",
	"CSRRCI",

	// 11.2: Floating-Point Control and Status Register (F)
	"FRCSR",
	"FSCSR",
	"FRRM",
	"FSRM",
	"FRFLAGS",
	"FSFLAGS",

	// 11.5: Single-Precision Load and Store Instructions (F)
	"FLW",
	"FSW",

	// 11.6: Single-Precision Floating-Point Computational Instructions
	"FADDS",
	"FSUBS",
	"FMULS",
	"FDIVS",
	"FSQRTS",
	"FMINS",
	"FMAXS",
	"FMADDS",
	"FMSUBS",
	"FNMADDS",
	"FNMSUBS",

	// 11.7: Single-Precision Floating-Point Conversion and Move Instructions
	"FCVTWS",
	"FCVTLS",
	"FCVTSW",
	"FCVTSL",
	"FCVTWUS",
	"FCVTLUS",
	"FCVTSWU",
	"FCVTSLU",
	"FSGNJS",
	"FSGNJNS",
	"FSGNJXS",
	"FMVXS",
	"FMVSX",
	"FMVXW",
	"FMVWX",

	// 11.8: Single-Precision Floating-Point Compare Instructions
	"FEQS",
	"FLTS",
	"FLES",

	// 11.9: Single-Precision Floating-Point Classify Instruction
	"FCLASSS",

	// 12.3: Double-Precision Load and Store Instructions (D)
	"FLD",
	"FSD",

	// 12.4: Double-Precision Floating-Point Computational Instructions
	"FADDD",
	"FSUBD",
	"FMULD",
	"FDIVD",
	"FMIND",
	"FMAXD",
	"FSQRTD",
	"FMADDD",
	"FMSUBD",
	"FNMADDD",
	"FNMSUBD",

	// 12.5: Double-Precision Floating-Point Conversion and Move Instructions
	"FCVTWD",
	"FCVTLD",
	"FCVTDW",
	"FCVTDL",
	"FCVTWUD",
	"FCVTLUD",
	"FCVTDWU",
	"FCVTDLU",
	"FCVTSD",
	"FCVTDS",
	"FSGNJD",
	"FSGNJND",
	"FSGNJXD",
	"FMVXD",
	"FMVDX",

	// 12.6: Double-Precision Floating-Point Compare Instructions
	"FEQD",
	"FLTD",
	"FLED",

	// 12.7: Double-Precision Floating-Point Classify Instruction
	"FCLASSD",

	// End marker
	"LAST",
}
