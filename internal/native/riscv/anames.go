// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

// 机器模式
var _ModeNames = []string{
	UserMode:       "User",
	SupervisorMode: "Supervisor",
	MachineMode:    "Machine",
}

// 寄存器名字列表
var _Register = []string{
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

// 寄存器名字列表(中文)
var _ZhRegister = []string{
	// 通用寄存器
	REG_X0:  "甲格",
	REG_X1:  "乙格",
	REG_X2:  "丙格",
	REG_X3:  "丁格",
	REG_X4:  "戊格",
	REG_X5:  "己格",
	REG_X6:  "庚格",
	REG_X7:  "辛格",
	REG_X8:  "壬格",
	REG_X9:  "癸格",
	REG_X10: "子格",
	REG_X11: "丑格",
	REG_X12: "寅格",
	REG_X13: "卯格",
	REG_X14: "辰格",
	REG_X15: "巳格",
	REG_X16: "午格",
	REG_X17: "未格",
	REG_X18: "申格",
	REG_X19: "酉格",
	REG_X20: "戌格",
	REG_X21: "亥格",
	REG_X22: "乾格",
	REG_X23: "坤格",
	REG_X24: "震格",
	REG_X25: "巽格",
	REG_X26: "坎格",
	REG_X27: "离格",
	REG_X28: "艮格",
	REG_X29: "兑格",
	REG_X30: "阴格",
	REG_X31: "阳格",

	// 浮点数寄存器(F/D扩展)
	REG_F0:  "甲皿",
	REG_F1:  "乙皿",
	REG_F2:  "丙皿",
	REG_F3:  "丁皿",
	REG_F4:  "戊皿",
	REG_F5:  "己皿",
	REG_F6:  "庚皿",
	REG_F7:  "辛皿",
	REG_F8:  "壬皿",
	REG_F9:  "癸皿",
	REG_F10: "子皿",
	REG_F11: "丑皿",
	REG_F12: "寅皿",
	REG_F13: "卯皿",
	REG_F14: "辰皿",
	REG_F15: "巳皿",
	REG_F16: "午皿",
	REG_F17: "未皿",
	REG_F18: "申皿",
	REG_F19: "酉皿",
	REG_F20: "戌皿",
	REG_F21: "亥皿",
	REG_F22: "乾皿",
	REG_F23: "坤皿",
	REG_F24: "震皿",
	REG_F25: "巽皿",
	REG_F26: "坎皿",
	REG_F27: "离皿",
	REG_F28: "艮皿",
	REG_F29: "兑皿",
	REG_F30: "阴皿",
	REG_F31: "阳皿",
}

// 寄存器别名
var _RegisterAlias = []string{
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
var _ZhRegisterAlias = []string{
	REG_ZERO: "零格",  // 零寄存器
	REG_RA:   "回格",  // 返回地址
	REG_SP:   "栈格",  // 栈指针
	REG_GP:   "基格",  // 全局基指针
	REG_TP:   "线格",  // 线程指针
	REG_T0:   "暂甲格", // 临时变量
	REG_T1:   "暂乙格", // 临时变量
	REG_T2:   "暂丙格", // 临时变量
	REG_S0:   "守甲格", // Saved register, 帧指针
	REG_S1:   "守乙格", // Saved register
	REG_A0:   "参甲格", // 函数参数/返回值
	REG_A1:   "参乙格", // 函数参数/返回值
	REG_A2:   "参丙格", // 函数参数
	REG_A3:   "参丁格", // 函数参数
	REG_A4:   "参戊格", // 函数参数
	REG_A5:   "参己格", // 函数参数
	REG_A6:   "参庚格", // 函数参数
	REG_A7:   "参辛格", // 函数参数
	REG_S2:   "守丙格", // Saved register
	REG_S3:   "守丁格", // Saved register
	REG_S4:   "守戊格", // Saved register
	REG_S5:   "守己格", // Saved register
	REG_S6:   "守庚格", // Saved register
	REG_S7:   "守辛格", // Saved register
	REG_S8:   "守壬格", // Saved register
	REG_S9:   "守癸格", // Saved register
	REG_S10:  "守子格", // Saved register
	REG_S11:  "守丑格", // Saved register
	REG_T3:   "暂丁格", // 临时变量
	REG_T4:   "暂戊格", // 临时变量
	REG_T5:   "暂己格", // 临时变量
	REG_T6:   "暂庚格", // 临时变量

	REG_FT0:  "暂甲皿", // 临时变量
	REG_FT1:  "暂乙皿", // 临时变量
	REG_FT2:  "暂丙皿", // 临时变量
	REG_FT3:  "暂丁皿", // 临时变量
	REG_FT4:  "暂戊皿", // 临时变量
	REG_FT5:  "暂己皿", // 临时变量
	REG_FT6:  "暂庚皿", // 临时变量
	REG_FT7:  "暂辛皿", // 临时变量
	REG_FS0:  "守甲皿", // Saved register
	REG_FS1:  "守乙皿", // Saved register
	REG_FA0:  "参甲皿", // 函数参数/返回值
	REG_FA1:  "参乙皿", // 函数参数/返回值
	REG_FA2:  "参丙皿", // 函数参数
	REG_FA3:  "参丁皿", // 函数参数
	REG_FA4:  "参戊皿", // 函数参数
	REG_FA5:  "参己皿", // 函数参数
	REG_FA6:  "参庚皿", // 函数参数
	REG_FA7:  "参辛皿", // 函数参数
	REG_FS2:  "守丙皿", // Saved register
	REG_FS3:  "守丁皿", // Saved register
	REG_FS4:  "守戊皿", // Saved register
	REG_FS5:  "守己皿", // Saved register
	REG_FS6:  "守庚皿", // Saved register
	REG_FS7:  "守辛皿", // Saved register
	REG_FS8:  "守壬皿", // Saved register
	REG_FS9:  "守癸皿", // Saved register
	REG_FS10: "守子皿", // Saved register
	REG_FS11: "守丑皿", // Saved register
	REG_FT8:  "暂壬皿", // 临时变量
	REG_FT9:  "暂癸皿", // 临时变量
	REG_FT10: "暂子皿", // 临时变量
	REG_FT11: "暂丑皿", // 临时变量
}

// 指令的名字
// 保持和指令定义相同的顺序
var _Anames = []string{
	//
	// Unprivileged ISA (version 20240411)
	//

	// 2.4: Integer Computational Instructions (RV32I)
	AADDI:  "ADDI",  // 常加, 加常量
	ASLTI:  "SLTI",  // 小常, 小于有符号常量时置位
	ASLTIU: "SLTIU", // 低常, 小于无符号常量时置位
	AANDI:  "ANDI",  // 常与, 常量与
	AORI:   "ORI",   // 常或, 常量或
	AXORI:  "XORI",  // 常异, 常量异或
	ASLLI:  "SLLI",  // 常左, 逻辑左移常量位
	ASRLI:  "SRLI",  // 常右, 逻辑右移常量位
	ASRAI:  "SRAI",  // 常佑, 算术右移常量位
	ALUI:   "LUI",   // 常载, 加载常量
	AAUIPC: "AUIPC", // 加计, 常量加程序计数器PC
	AADD:   "ADD",   // 相加, 加法
	ASLT:   "SLT",   // 小于, 有符号小于时置位
	ASLTU:  "SLTU",  // 低于, 无符号小于时置位
	AAND:   "AND",   // 相与, 与
	AOR:    "OR",    // 相或, 或
	AXOR:   "XOR",   // 异或, 异或
	ASLL:   "SLL",   // 左移, 逻辑左移
	ASRL:   "SRL",   // 右移, 逻辑右移
	ASUB:   "SUB",   // 相减, 减法
	ASRA:   "SRA",   // 佑移, 算术右移

	// 2.5: Control Transfer Instructions (RV32I)
	AJAL:  "JAL",  // 常转, 跳转到常量偏移地址
	AJALR: "JALR", // 寄转, 跳转到寄存器偏移地址
	ABEQ:  "BEQ",  // 等转, 相等则跳转
	ABNE:  "BNE",  // 异转, 不等则跳转
	ABLT:  "BLT",  // 小转, 有符号小于则跳转
	ABLTU: "BLTU", // 低转, 无符号小于则跳转
	ABGE:  "BGE",  // 大转, 有符号答疑则跳转
	ABGEU: "BGEU", // 高转, 无符号答疑则跳转

	// 2.6: Load and Store Instructions (RV32I)
	ALW:  "LW",  // 读普整, 读I32
	ALH:  "LH",  // 读短整, 读I16
	ALHU: "LHU", // 读短正, 读U16
	ALB:  "LB",  // 读微整, 读I8
	ALBU: "LBU", // 读微正, 读U8
	ASW:  "SW",  // 写普整, 写I32
	ASH:  "SH",  // 写短整, 写I16
	ASB:  "SB",  // 写微整, 写I8

	// 2.7: Memory Ordering Instructions (RV32I)
	AFENCE: "FENCE", // 读写毕, 读写屏障

	// 3.3.1: Environment Call and Breakpoint
	AECALL:  "ECALL",  // 环境调用
	AEBREAK: "EBREAK", // 环境断点

	// 4.2: Integer Computational Instructions (RV64I)
	AADDIW: "ADDIW", // 常加.普整
	ASLLIW: "SLLIW", // 常左.普整
	ASRLIW: "SRLIW", // 常右.普整
	ASRAIW: "SRAIW", // 常右.普整
	AADDW:  "ADDW",  // 相加.普整
	ASLLW:  "SLLW",  // 左移.普整
	ASRLW:  "SRLW",  // 右移.普整
	ASUBW:  "SUBW",  // 相减.普整
	ASRAW:  "SRAW",  // 佑移.普整

	// 4.3: Load and Store Instructions (RV64I)
	ALWU: "LWU", // 读普正, 读U32
	ALD:  "LD",  // 读长整, 读I64
	ASD:  "SD",  // 写长整, 写I64

	// 7.1: CSR Instructions (Zicsr)
	ACSRRW:  "CSRRW",  // CSR写读
	ACSRRS:  "CSRRS",  // CSR置读
	ACSRRC:  "CSRRC",  // CSR清读
	ACSRRWI: "CSRRWI", // CSR常写读
	ACSRRSI: "CSRRSI", // CSR常置读
	ACSRRCI: "CSRRCI", // CSR常清读

	// 13.1: Multiplication Operations (RV32M/RV64M)
	AMUL:    "MUL",    // 相乘
	AMULH:   "MULH",   // 相乘.高整
	AMULHU:  "MULHU",  // 相乘.高正
	AMULHSU: "MULHSU", // 相乘.高混
	AMULW:   "MULW",   // 乘法.普整, RV64M

	// 13.2: Division Operations (RV32M/RV64M)
	ADIV:   "DIV",   // 相除
	ADIVU:  "DIVU",  // 相除.正整
	AREM:   "REM",   // 取余
	AREMU:  "REMU",  // 取余.正整
	ADIVW:  "DIVW",  // 相除.普整, RV64M
	ADIVUW: "DIVUW", // 相除.普正, RV64M
	AREMW:  "REMW",  // 取余.普整, RV64M
	AREMUW: "REMUW", // 取余.普正, RV64M

	// 20.5: Single-Precision Load and Store Instructions (F)
	AFLW: "FLW", // 读浮点.单精
	AFSW: "FSW", // 写浮点.单精

	// 20.6: Single-Precision Floating-Point Computational Instructions
	AFADD_S:   "FADD_S",   // 浮点加.单精
	AFSUB_S:   "FSUB_S",   // 浮点减.单精
	AFMUL_S:   "FMUL_S",   // 浮点乘.单精
	AFDIV_S:   "FDIV_S",   // 浮点除.单精
	AFSQRT_S:  "FSQRT_S",  // 浮点平方根.单精度
	AFMIN_S:   "FMIN_S",   // 浮点最小值.单精度
	AFMAX_S:   "FMAX_S",   // 浮点最大值.单精度
	AFMADD_S:  "FMADD_S",  // 浮点乘加.单精度
	AFMSUB_S:  "FMSUB_S",  // 浮点乘减.单精度
	AFNMADD_S: "FNMADD_S", // 浮点负乘加.单精度
	AFNMSUB_S: "FNMSUB_S", // 浮点负乘减.单精度

	// 20.7: Single-Precision Floating-Point Conversion and Move Instructions
	AFCVT_W_S:  "FCVT_W_S",  // 浮点转换为字（有符号）
	AFCVT_L_S:  "FCVT_L_S",  // 浮点转换为长字（有符号，RV64）
	AFCVT_S_W:  "FCVT_S_W",  // 字转换为浮点（有符号）
	AFCVT_S_L:  "FCVT_S_L",  // 长字转换为浮点（有符号，RV64）
	AFCVT_WU_S: "FCVT_WU_S", // 浮点转换为字（无符号）
	AFCVT_LU_S: "FCVT_LU_S", // 浮点转换为长字（无符号，RV64）
	AFCVT_S_WU: "FCVT_S_WU", // 无符号字转换为浮点
	AFCVT_S_LU: "FCVT_S_LU", // 无符号长字转换为浮点（RV64）
	AFSGNJ_S:   "FSGNJ_S",   // 浮点符号复制
	AFSGNJN_S:  "FSGNJN_S",  // 浮点符号取反复制
	AFSGNJX_S:  "FSGNJX_S",  // 浮点符号异或复制
	AFMV_X_W:   "FMV_X_W",   // 浮点位模式移动到整数寄存器
	AFMV_W_X:   "FMV_W_X",   // 整数位模式移动到浮点寄存器

	// 20.8: Single-Precision Floating-Point Compare Instructions
	AFEQ_S: "FEQ_S", // 浮点相等比较.单精
	AFLT_S: "FLT_S", // 浮点小于比较.单精
	AFLE_S: "FLE_S", // 浮点小于或等于比较.单精

	// 20.9: Single-Precision Floating-Point Classify Instruction
	AFCLASS_S: "FCLASS_S", // 浮点数分类.单精

	// 21.3: Double-Precision Load and Store Instructions (D)
	AFLD: "FLD", // 读浮点.双精
	AFSD: "FSD", // 写浮点.双精

	// 21.4: Double-Precision Floating-Point Computational Instructions
	AFADD_D:   "FADD_D",
	AFSUB_D:   "FSUB_D",
	AFMUL_D:   "FMUL_D",
	AFDIV_D:   "FDIV_D",
	AFMIN_D:   "FMIN_D",
	AFMAX_D:   "FMAX_D",
	AFSQRT_D:  "FSQRT_D",
	AFMADD_D:  "FMADD_D",
	AFMSUB_D:  "FMSUB_D",
	AFNMADD_D: "FNMADD_D",
	AFNMSUB_D: "FNMSUB_D",

	// 21.5: Double-Precision Floating-Point Conversion and Move Instructions
	AFCVT_W_D:  "FCVT_W_D",
	AFCVT_L_D:  "FCVT_L_D",
	AFCVT_D_W:  "FCVT_D_W",
	AFCVT_D_L:  "FCVT_D_L",
	AFCVT_WU_D: "FCVT_WU_D",
	AFCVT_LU_D: "FCVT_LU_D",
	AFCVT_D_WU: "FCVT_D_WU",
	AFCVT_D_LU: "FCVT_D_LU",
	AFCVT_S_D:  "FCVT_S_D",
	AFCVT_D_S:  "FCVT_D_S",
	AFSGNJ_D:   "FSGNJ_D",
	AFSGNJN_D:  "FSGNJN_D",
	AFSGNJX_D:  "FSGNJX_D",
	AFMV_X_D:   "FMV_X_D",
	AFMV_D_X:   "FMV_D_X",

	// 21.6: Double-Precision Floating-Point Compare Instructions
	AFEQ_D: "FEQ_D",
	AFLT_D: "FLT_D",
	AFLE_D: "FLE_D",

	// 21.7: Double-Precision Floating-Point Classify Instruction
	AFCLASS_D: "FCLASS_D",

	// 伪指令(A_开头以区分)
	// ISA (version 20191213)
	// 25: RISC-V Assembly Programmer's Handbook
	// 只保留可以1:1映射到原生指令的类型
	// 长地址跳转需要用户手动处理

	A_NOP:       "NOP",
	A_MV:        "MV",
	A_NOT:       "NOT",
	A_NEG:       "NEG",
	A_NEGW:      "NEGW",
	A_SEXT_W:    "SEXT_W",
	A_SEQZ:      "SEQZ",
	A_SNEZ:      "SNEZ",
	A_SLTZ:      "SLTZ",
	A_SGTZ:      "SGTZ",
	A_FMV_S:     "FMV_S",
	A_FABS_S:    "FABS_S",
	A_FNEG_S:    "FNEG_S",
	A_FMV_D:     "FMV_D",
	A_FABS_D:    "FABS_D",
	A_FNEG_D:    "FNEG_D",
	A_BEQZ:      "BEQZ",
	A_BNEZ:      "BNEZ",
	A_BLEZ:      "BLEZ",
	A_BGEZ:      "BGEZ",
	A_BLTZ:      "BLTZ",
	A_BGTZ:      "BGTZ",
	A_BGT:       "BGT",
	A_BLE:       "BLE",
	A_BGTU:      "BGTU",
	A_BLEU:      "BLEU",
	A_J:         "J",
	A_JR:        "JR",
	A_RET:       "RET",
	A_RDINSTRET: "RDINSTRET",
	A_RDCYCLE:   "RDCYCLE",
	A_RDTIME:    "RDTIME",
	A_CSRR:      "CSRR",
	A_CSRW:      "CSRW",
	A_CSRS:      "CSRS",
	A_CSRC:      "CSRC",
	A_CSRWI:     "CSRWI",
	A_CSRSI:     "CSRSI",
	A_CSRCI:     "CSRCI",
	A_FRCSR:     "FRCSR",
	A_FSCSR:     "FSCSR",
	A_FRRM:      "FRRM",
	A_FSRM:      "FSRM",
	A_FRFLAGS:   "FRFLAGS",
	A_FSFLAGS:   "FSFLAGS",

	// End marker
	ALAST: "LAST",
}
