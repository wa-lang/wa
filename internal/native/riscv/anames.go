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

	// End marker
	"LAST",
}
