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
	AADDI:  "ADDI",  // 某加常, 加常量
	ASLTI:  "SLTI",  // 某小常, 小于有符号常量时置位
	ASLTIU: "SLTIU", // 某低常, 小于无符号常量时置位
	AANDI:  "ANDI",  // 某与常, 常量与
	AORI:   "ORI",   // 某或常, 常量或
	AXORI:  "XORI",  // 某异常, 常量异或
	ASLLI:  "SLLI",  // 某左常, 逻辑左移常量位
	ASRLI:  "SRLI",  // 某右常, 逻辑右移常量位
	ASRAI:  "SRAI",  // 某佑常, 算术右移常量位
	ALUI:   "LUI",   // 常赋某, 加载常量
	AAUIPC: "AUIPC", // 计加常, 常量加程序计数器PC
	AADD:   "ADD",   // 某加某, 加法
	ASLT:   "SLT",   // 某小某, 有符号小于时置位
	ASLTU:  "SLTU",  // 某低某, 无符号小于时置位
	AAND:   "AND",   // 某与某, 与
	AOR:    "OR",    // 某或某, 或
	AXOR:   "XOR",   // 某异某, 异或
	ASLL:   "SLL",   // 某左某, 逻辑左移
	ASRL:   "SRL",   // 某右某, 逻辑右移
	ASUB:   "SUB",   // 某减某, 减法
	ASRA:   "SRA",   // 某佑某, 算术右移

	// 2.5: Control Transfer Instructions (RV32I)
	AJAL:  "JAL",  // 转常址, 跳转到常量偏移地址
	AJALR: "JALR", // 转某址, 跳转到寄存器偏移地址
	ABEQ:  "BEQ",  // 若等转, 相等则跳转
	ABNE:  "BNE",  // 若异转, 不等则跳转
	ABLT:  "BLT",  // 若小转, 有符号小于则跳转
	ABLTU: "BLTU", // 若低转, 无符号小于则跳转
	ABGE:  "BGE",  // 若大转, 有符号答疑则跳转
	ABGEU: "BGEU", // 若高转, 无符号答疑则跳转

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
	AECALL:  "ECALL",  // 调用某
	AEBREAK: "EBREAK", // 暂断点

	// 4.2: Integer Computational Instructions (RV64I)
	AADDIW: "ADDIW", // 某加常.普整
	ASLLIW: "SLLIW", // 某左常.普整
	ASRLIW: "SRLIW", // 某右常.普整
	ASRAIW: "SRAIW", // 某佑常.普整
	AADDW:  "ADDW",  // 某加某.普整
	ASLLW:  "SLLW",  // 某左某.普整
	ASRLW:  "SRLW",  // 某右某.普整
	ASUBW:  "SUBW",  // 某减某.普整
	ASRAW:  "SRAW",  // 某佑某.普整

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
	AMUL:    "MUL",    // 某乘某
	AMULH:   "MULH",   // 某乘某.高整
	AMULHU:  "MULHU",  // 某乘某.高正
	AMULHSU: "MULHSU", // 某乘某.高混
	AMULW:   "MULW",   // 某乘某.普整, RV64M

	// 13.2: Division Operations (RV32M/RV64M)
	ADIV:   "DIV",   // 某除某
	ADIVU:  "DIVU",  // 某模某.正整
	AREM:   "REM",   // 某模某
	AREMU:  "REMU",  // 某模某.正整
	ADIVW:  "DIVW",  // 某除某.普整, RV64M
	ADIVUW: "DIVUW", // 某除某.普正, RV64M
	AREMW:  "REMW",  // 某模某.普整, RV64M
	AREMUW: "REMUW", // 某模某.普正, RV64M

	// 20.5: Single-Precision Load and Store Instructions (F)
	AFLW: "FLW", // 读浮点
	AFSW: "FSW", // 写浮点

	// 20.6: Single-Precision Floating-Point Computational Instructions
	AFADD_S:   "FADD_S",   // 浮加浮
	AFSUB_S:   "FSUB_S",   // 浮减浮
	AFMUL_S:   "FMUL_S",   // 浮乘浮
	AFDIV_S:   "FDIV_S",   // 浮除浮
	AFSQRT_S:  "FSQRT_S",  // 浮平方根
	AFMIN_S:   "FMIN_S",   // 浮最小值
	AFMAX_S:   "FMAX_S",   // 浮最大值
	AFMADD_S:  "FMADD_S",  // 浮正乘加
	AFMSUB_S:  "FMSUB_S",  // 浮正乘减
	AFNMADD_S: "FNMADD_S", // 浮负乘加
	AFNMSUB_S: "FNMSUB_S", // 浮负乘减

	// 20.7: Single-Precision Floating-Point Conversion and Move Instructions
	AFCVT_W_S:  "FCVT_W_S",  // 浮转普整
	AFCVT_L_S:  "FCVT_L_S",  // 浮转长整(RV64)
	AFCVT_S_W:  "FCVT_S_W",  // 普整转浮
	AFCVT_S_L:  "FCVT_S_L",  // 长整转浮(RV64)
	AFCVT_WU_S: "FCVT_WU_S", // 浮转普正
	AFCVT_LU_S: "FCVT_LU_S", // 浮转长正(RV64)
	AFCVT_S_WU: "FCVT_S_WU", // 普正转浮
	AFCVT_S_LU: "FCVT_S_LU", // 长正转浮(RV64)
	AFSGNJ_S:   "FSGNJ_S",   // 浮符号正复制
	AFSGNJN_S:  "FSGNJN_S",  // 浮符号负复制
	AFSGNJX_S:  "FSGNJX_S",  // 浮符号异复制
	AFMV_X_W:   "FMV_X_W",   // 浮赋某.位模式
	AFMV_W_X:   "FMV_W_X",   // 某赋浮.位模式

	// 20.8: Single-Precision Floating-Point Compare Instructions
	AFEQ_S: "FEQ_S", // 浮等浮
	AFLT_S: "FLT_S", // 浮小浮
	AFLE_S: "FLE_S", // 浮弱浮

	// 20.9: Single-Precision Floating-Point Classify Instruction
	AFCLASS_S: "FCLASS_S", // 浮点分类

	// 21.3: Double-Precision Load and Store Instructions (D)
	AFLD: "FLD", // 读双精
	AFSD: "FSD", // 写双精

	// 21.4: Double-Precision Floating-Point Computational Instructions
	AFADD_D:   "FADD_D",   // 双精加
	AFSUB_D:   "FSUB_D",   // 双精减
	AFMUL_D:   "FMUL_D",   // 双精乘
	AFDIV_D:   "FDIV_D",   // 双精除
	AFMIN_D:   "FMIN_D",   // 双精平方根
	AFMAX_D:   "FMAX_D",   // 双精最小值
	AFSQRT_D:  "FSQRT_D",  // 双精最大值
	AFMADD_D:  "FMADD_D",  // 双精正乘加
	AFMSUB_D:  "FMSUB_D",  // 双精正乘减
	AFNMADD_D: "FNMADD_D", // 双精负乘加
	AFNMSUB_D: "FNMSUB_D", // 双精负乘减

	// 21.5: Double-Precision Floating-Point Conversion and Move Instructions
	AFCVT_W_D:  "FCVT_W_D",  // 双精转普整
	AFCVT_L_D:  "FCVT_L_D",  // 双精转长整
	AFCVT_D_W:  "FCVT_D_W",  // 普整转双精
	AFCVT_D_L:  "FCVT_D_L",  // 长整转双精
	AFCVT_WU_D: "FCVT_WU_D", // 双精转普正
	AFCVT_LU_D: "FCVT_LU_D", // 双精转长正
	AFCVT_D_WU: "FCVT_D_WU", // 普正转双精
	AFCVT_D_LU: "FCVT_D_LU", // 长正转双精
	AFCVT_S_D:  "FCVT_S_D",  // 双精转单精
	AFCVT_D_S:  "FCVT_D_S",  // 单精转双精
	AFSGNJ_D:   "FSGNJ_D",   // 双精符号正复制
	AFSGNJN_D:  "FSGNJN_D",  // 双精符号负复制
	AFSGNJX_D:  "FSGNJX_D",  // 双精符号异复制
	AFMV_X_D:   "FMV_X_D",   // 双精赋某.位模式
	AFMV_D_X:   "FMV_D_X",   // 某赋双精.位模式

	// 21.6: Double-Precision Floating-Point Compare Instructions
	AFEQ_D: "FEQ_D", // 双精等
	AFLT_D: "FLT_D", // 双精小
	AFLE_D: "FLE_D", // 双精弱

	// 21.7: Double-Precision Floating-Point Classify Instruction
	AFCLASS_D: "FCLASS_D", // 双精分类

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

var _ZhAnames = []string{
	//
	// Unprivileged ISA (version 20240411)
	//

	// 2.4: Integer Computational Instructions (RV32I)
	AADDI:  "某加常", // 某加常, 加常量
	ASLTI:  "某小常", // 某小常, 小于有符号常量时置位
	ASLTIU: "某低常", // 某低常, 小于无符号常量时置位
	AANDI:  "某与常", // 某与常, 常量与
	AORI:   "某或常", // 某或常, 常量或
	AXORI:  "某异常", // 某异常, 常量异或
	ASLLI:  "某左常", // 某左常, 逻辑左移常量位
	ASRLI:  "某右常", // 某右常, 逻辑右移常量位
	ASRAI:  "某佑常", // 某佑常, 算术右移常量位
	ALUI:   "常赋某", // 常赋某, 加载常量
	AAUIPC: "计加常", // 计加常, 常量加程序计数器PC
	AADD:   "某加某", // 某加某, 加法
	ASLT:   "某小某", // 某小某, 有符号小于时置位
	ASLTU:  "某低某", // 某低某, 无符号小于时置位
	AAND:   "某与某", // 某与某, 与
	AOR:    "某或某", // 某或某, 或
	AXOR:   "某异某", // 某异某, 异或
	ASLL:   "某左某", // 某左某, 逻辑左移
	ASRL:   "某右某", // 某右某, 逻辑右移
	ASUB:   "某减某", // 某减某, 减法
	ASRA:   "某佑某", // 某佑某, 算术右移

	// 2.5: Control Transfer Instructions (RV32I)
	AJAL:  "转常址", // 转常址, 跳转到常量偏移地址
	AJALR: "转某址", // 转某址, 跳转到寄存器偏移地址
	ABEQ:  "若等转", // 若等转, 相等则跳转
	ABNE:  "若异转", // 若异转, 不等则跳转
	ABLT:  "若小转", // 若小转, 有符号小于则跳转
	ABLTU: "若低转", // 若低转, 无符号小于则跳转
	ABGE:  "若大转", // 若大转, 有符号答疑则跳转
	ABGEU: "若高转", // 若高转, 无符号答疑则跳转

	// 2.6: Load and Store Instructions (RV32I)
	ALW:  "读普整", // 读普整, 读I32
	ALH:  "读短整", // 读短整, 读I16
	ALHU: "读短正", // 读短正, 读U16
	ALB:  "读微整", // 读微整, 读I8
	ALBU: "读微正", // 读微正, 读U8
	ASW:  "写普整", // 写普整, 写I32
	ASH:  "写短整", // 写短整, 写I16
	ASB:  "写微整", // 写微整, 写I8

	// 2.7: Memory Ordering Instructions (RV32I)
	AFENCE: "读写毕", // 读写毕, 读写屏障

	// 3.3.1: Environment Call and Breakpoint
	AECALL:  "调用某", // 调用某
	AEBREAK: "暂断点", // 暂断点

	// 4.2: Integer Computational Instructions (RV64I)
	AADDIW: "某加常.普整", // 某加常.普整
	ASLLIW: "某左常.普整", // 某左常.普整
	ASRLIW: "某右常.普整", // 某右常.普整
	ASRAIW: "某佑常.普整", // 某佑常.普整
	AADDW:  "某加某.普整", // 某加某.普整
	ASLLW:  "某左某.普整", // 某左某.普整
	ASRLW:  "某右某.普整", // 某右某.普整
	ASUBW:  "某减某.普整", // 某减某.普整
	ASRAW:  "某佑某.普整", // 某佑某.普整

	// 4.3: Load and Store Instructions (RV64I)
	ALWU: "读普正", // 读普正, 读U32
	ALD:  "读长整", // 读长整, 读I64
	ASD:  "写长整", // 写长整, 写I64

	// 7.1: CSR Instructions (Zicsr)
	ACSRRW:  "CSR写读",  // CSR写读
	ACSRRS:  "CSR置读",  // CSR置读
	ACSRRC:  "CSR清读",  // CSR清读
	ACSRRWI: "CSR常写读", // CSR常写读
	ACSRRSI: "CSR常置读", // CSR常置读
	ACSRRCI: "CSR常清读", // CSR常清读

	// 13.1: Multiplication Operations (RV32M/RV64M)
	AMUL:    "某乘某",    // 某乘某
	AMULH:   "某乘某.高整", // 某乘某.高整
	AMULHU:  "某乘某.高正", // 某乘某.高正
	AMULHSU: "某乘某.高混", // 某乘某.高混
	AMULW:   "某乘某.普整", // 某乘某.普整, RV64M

	// 13.2: Division Operations (RV32M/RV64M)
	ADIV:   "某除某",    // 某除某
	ADIVU:  "某模某.正整", // 某模某.正整
	AREM:   "某模某",    // 某模某
	AREMU:  "某模某.正整", // 某模某.正整
	ADIVW:  "某除某.普整", // 某除某.普整, RV64M
	ADIVUW: "某除某.普正", // 某除某.普正, RV64M
	AREMW:  "某模某.普整", // 某模某.普整, RV64M
	AREMUW: "某模某.普正", // 某模某.普正, RV64M

	// 20.5: Single-Precision Load and Store Instructions (F)
	AFLW: "读浮点", // 读浮点
	AFSW: "写浮点", // 写浮点

	// 20.6: Single-Precision Floating-Point Computational Instructions
	AFADD_S:   "浮加浮",  // 浮加浮
	AFSUB_S:   "浮减浮",  // 浮减浮
	AFMUL_S:   "浮乘浮",  // 浮乘浮
	AFDIV_S:   "浮除浮",  // 浮除浮
	AFSQRT_S:  "浮平方根", // 浮平方根
	AFMIN_S:   "浮最小值", // 浮最小值
	AFMAX_S:   "浮最大值", // 浮最大值
	AFMADD_S:  "浮正乘加", // 浮正乘加
	AFMSUB_S:  "浮正乘减", // 浮正乘减
	AFNMADD_S: "浮负乘加", // 浮负乘加
	AFNMSUB_S: "浮负乘减", // 浮负乘减

	// 20.7: Single-Precision Floating-Point Conversion and Move Instructions
	AFCVT_W_S:  "浮转普整",    // 浮转普整
	AFCVT_L_S:  "浮转长整",    // 浮转长整(RV64)
	AFCVT_S_W:  "普整转浮",    // 普整转浮
	AFCVT_S_L:  "长整转浮",    // 长整转浮(RV64)
	AFCVT_WU_S: "浮转普正",    // 浮转普正
	AFCVT_LU_S: "浮转长正",    // 浮转长正(RV64)
	AFCVT_S_WU: "普正转浮",    // 普正转浮
	AFCVT_S_LU: "长正转浮",    // 长正转浮(RV64)
	AFSGNJ_S:   "浮符号正复制",  // 浮符号正复制
	AFSGNJN_S:  "浮符号负复制",  // 浮符号负复制
	AFSGNJX_S:  "浮符号异复制",  // 浮符号异复制
	AFMV_X_W:   "浮赋某.位模式", // 浮赋某.位模式
	AFMV_W_X:   "某赋浮.位模式", // 某赋浮.位模式

	// 20.8: Single-Precision Floating-Point Compare Instructions
	AFEQ_S: "浮等浮", // 浮等浮
	AFLT_S: "浮小浮", // 浮小浮
	AFLE_S: "浮弱浮", // 浮弱浮

	// 20.9: Single-Precision Floating-Point Classify Instruction
	AFCLASS_S: "浮点分类", // 浮点分类

	// 21.3: Double-Precision Load and Store Instructions (D)
	AFLD: "读双精", // 读双精
	AFSD: "写双精", // 写双精

	// 21.4: Double-Precision Floating-Point Computational Instructions
	AFADD_D:   "双精加",   // 双精加
	AFSUB_D:   "双精减",   // 双精减
	AFMUL_D:   "双精乘",   // 双精乘
	AFDIV_D:   "双精除",   // 双精除
	AFMIN_D:   "双精平方根", // 双精平方根
	AFMAX_D:   "双精最小值", // 双精最小值
	AFSQRT_D:  "双精最大值", // 双精最大值
	AFMADD_D:  "双精正乘加", // 双精正乘加
	AFMSUB_D:  "双精正乘减", // 双精正乘减
	AFNMADD_D: "双精负乘加", // 双精负乘加
	AFNMSUB_D: "双精负乘减", // 双精负乘减

	// 21.5: Double-Precision Floating-Point Conversion and Move Instructions
	AFCVT_W_D:  "双精转普整",    // 双精转普整
	AFCVT_L_D:  "双精转长整",    // 双精转长整
	AFCVT_D_W:  "普整转双精",    // 普整转双精
	AFCVT_D_L:  "长整转双精",    // 长整转双精
	AFCVT_WU_D: "双精转普正",    // 双精转普正
	AFCVT_LU_D: "双精转长正",    // 双精转长正
	AFCVT_D_WU: "普正转双精",    // 普正转双精
	AFCVT_D_LU: "长正转双精",    // 长正转双精
	AFCVT_S_D:  "双精转单精",    // 双精转单精
	AFCVT_D_S:  "单精转双精",    // 单精转双精
	AFSGNJ_D:   "双精符号正复制",  // 双精符号正复制
	AFSGNJN_D:  "双精符号负复制",  // 双精符号负复制
	AFSGNJX_D:  "双精符号异复制",  // 双精符号异复制
	AFMV_X_D:   "双精赋某.位模式", // 双精赋某.位模式
	AFMV_D_X:   "某赋双精.位模式", // 某赋双精.位模式

	// 21.6: Double-Precision Floating-Point Compare Instructions
	AFEQ_D: "双精等", // 双精等
	AFLT_D: "双精小", // 双精小
	AFLE_D: "双精弱", // 双精弱

	// 21.7: Double-Precision Floating-Point Classify Instruction
	AFCLASS_D: "双精分类", // 双精分类
}
