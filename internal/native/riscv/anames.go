// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

// 寄存器名字列表
var Register = []string{
	// 通用寄存器
	REG_X0:  "X0",
	REG_X1:  "X1",
	REG_X2:  "X2",
	REG_X3:  "X3",
	REG_X4:  "X4",
	REG_X5:  "X5",
	REG_X6:  "X6",
	REG_X7:  "X7",
	REG_X8:  "X8",
	REG_X9:  "X9",
	REG_X10: "X10",
	REG_X11: "X11",
	REG_X12: "X12",
	REG_X13: "X13",
	REG_X14: "X14",
	REG_X15: "X15",
	REG_X16: "X16",
	REG_X17: "X17",
	REG_X18: "X18",
	REG_X19: "X19",
	REG_X20: "X20",
	REG_X21: "X21",
	REG_X22: "X22",
	REG_X23: "X23",
	REG_X24: "X24",
	REG_X25: "X25",
	REG_X26: "X26",
	REG_X27: "X27",
	REG_X28: "X28",
	REG_X29: "X29",
	REG_X30: "X30",
	REG_X31: "X31",

	// 浮点数寄存器(F/D扩展)
	REG_F0:  "F0",
	REG_F1:  "F1",
	REG_F2:  "F2",
	REG_F3:  "F3",
	REG_F4:  "F4",
	REG_F5:  "F5",
	REG_F6:  "F6",
	REG_F7:  "F7",
	REG_F8:  "F8",
	REG_F9:  "F9",
	REG_F10: "F10",
	REG_F11: "F11",
	REG_F12: "F12",
	REG_F13: "F13",
	REG_F14: "F14",
	REG_F15: "F15",
	REG_F16: "F16",
	REG_F17: "F17",
	REG_F18: "F18",
	REG_F19: "F19",
	REG_F20: "F20",
	REG_F21: "F21",
	REG_F22: "F22",
	REG_F23: "F23",
	REG_F24: "F24",
	REG_F25: "F25",
	REG_F26: "F26",
	REG_F27: "F27",
	REG_F28: "F28",
	REG_F29: "F29",
	REG_F30: "F30",
	REG_F31: "F31",
}

// 寄存器别名
var RegisterAlias = []string{
	REG_ZERO: "ZERO", // 零寄存器
	REG_RA:   "RA",   // 返回地址
	REG_SP:   "SP",   // 栈指针
	REG_GP:   "GP",   // 全局基指针
	REG_TP:   "TP",   // 线程指针
	REG_T0:   "T0",   // 临时变量
	REG_T1:   "T1",   // 临时变量
	REG_T2:   "T2",   // 临时变量
	REG_S0:   "S0",   // Saved register, 帧指针
	REG_S1:   "S1",   // Saved register
	REG_A0:   "A0",   // 函数参数/返回值
	REG_A1:   "A1",   // 函数参数/返回值
	REG_A2:   "A2",   // 函数参数
	REG_A3:   "A3",   // 函数参数
	REG_A4:   "A4",   // 函数参数
	REG_A5:   "A5",   // 函数参数
	REG_A6:   "A6",   // 函数参数
	REG_A7:   "A7",   // 函数参数
	REG_S2:   "S2",   // Saved register
	REG_S3:   "S3",   // Saved register
	REG_S4:   "S4",   // Saved register
	REG_S5:   "S5",   // Saved register
	REG_S6:   "S6",   // Saved register
	REG_S7:   "S7",   // Saved register
	REG_S8:   "S8",   // Saved register
	REG_S9:   "S9",   // Saved register
	REG_S10:  "S10",  // Saved register
	REG_S11:  "S10",  // Saved register
	REG_T3:   "T3",   // 临时变量
	REG_T4:   "T4",   // 临时变量
	REG_T5:   "T5",   // 临时变量
	REG_T6:   "T6",   // 临时变量

	REG_FT0:  "FT0",  // 临时变量
	REG_FT1:  "FT1",  // 临时变量
	REG_FT2:  "FT2",  // 临时变量
	REG_FT3:  "FT3",  // 临时变量
	REG_FT4:  "FT4",  // 临时变量
	REG_FT5:  "FT5",  // 临时变量
	REG_FT6:  "FT6",  // 临时变量
	REG_FT7:  "FT7",  // 临时变量
	REG_FS0:  "FS0",  // Saved register
	REG_FS1:  "FS1",  // Saved register
	REG_FA0:  "FA0",  // 函数参数/返回值
	REG_FA1:  "FA1",  // 函数参数/返回值
	REG_FA2:  "FA2",  // 函数参数
	REG_FA3:  "FA3",  // 函数参数
	REG_FA4:  "FA4",  // 函数参数
	REG_FA5:  "FA5",  // 函数参数
	REG_FA6:  "FA6",  // 函数参数
	REG_FA7:  "FA7",  // 函数参数
	REG_FS2:  "FS2",  // Saved register
	REG_FS3:  "FS3",  // Saved register
	REG_FS4:  "FS4",  // Saved register
	REG_FS5:  "FS5",  // Saved register
	REG_FS6:  "FS6",  // Saved register
	REG_FS7:  "FS7",  // Saved register
	REG_FS8:  "FS8",  // Saved register
	REG_FS9:  "FS9",  // Saved register
	REG_FS10: "FS10", // Saved register
	REG_FS11: "FS11", // Saved register
	REG_FT8:  "FT8",  // 临时变量
	REG_FT9:  "FT9",  // 临时变量
	REG_FT10: "FT10", // 临时变量
	REG_FT11: "FT11", // 临时变量
}

// 指令的名字
// 保持和指令定义相同的顺序
var Anames = []string{
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

	// 3.3.1: Environment Call and Breakpoint
	AECALL:  "ECALL",
	AEBREAK: "EBREAK",

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

	// 伪指令
	// ISA (version 20191213)
	// 25: RISC-V Assembly Programmer's Handbook
	// 只保留可以1:1映射到原生指令的类型
	// 长地址跳转需要用户手动处理

	ANOP:       "NOP",
	AMV:        "MV",
	ANOT:       "NOT",
	ANEG:       "NEG",
	ANEGW:      "NEGW",
	ASEXT_W:    "SEXT_W",
	ASEQZ:      "SEQZ",
	ASNEZ:      "SNEZ",
	ASLTZ:      "SLTZ",
	ASGTZ:      "SGTZ",
	AFMV_S:     "FMV_S",
	AFABS_S:    "FABS_S",
	AFNEG_S:    "FNEG_S",
	AFMV_D:     "FMV_D",
	AFABS_D:    "FABS_D",
	AFNEG_D:    "FNEG_D",
	ABEQZ:      "BEQZ",
	ABNEZ:      "BNEZ",
	ABLEZ:      "BLEZ",
	ABGEZ:      "BGEZ",
	ABLTZ:      "BLTZ",
	ABGTZ:      "BGTZ",
	ABGT:       "BGT",
	ABLE:       "BLE",
	ABGTU:      "BGTU",
	ABLEU:      "BLEU",
	AJ:         "J",
	AJR:        "JR",
	ARET:       "RET",
	ARDINSTRET: "RDINSTRET",
	ARDCYCLE:   "RDCYCLE",
	ARDTIME:    "RDTIME",
	ACSRR:      "CSRR",
	ACSRW:      "CSRW",
	ACSRS:      "CSRS",
	ACSRC:      "CSRC",
	ACSRWI:     "CSRWI",
	ACSRSI:     "CSRSI",
	ACSRCI:     "CSRCI",
	AFRCSR:     "FRCSR",
	AFSCSR:     "FSCSR",
	AFRRM:      "FRRM",
	AFSRM:      "FSRM",
	AFRFLAGS:   "FRFLAGS",
	AFSFLAGS:   "FSFLAGS",

	// End marker
	ALAST: "LAST",
}
