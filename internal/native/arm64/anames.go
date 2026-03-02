// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package arm64

// 寄存器名字列表
var _Register = []string{
	// 特殊寄存器 (硬件编码31)
	REG_SP:  "sp",
	REG_WZR: "wzr", // zero
	REG_XZR: "xzr", // zero

	// 低32位寄存器
	REG_W0:  "w0",
	REG_W1:  "w1",
	REG_W2:  "w2",
	REG_W3:  "w3",
	REG_W4:  "w4",
	REG_W5:  "w5",
	REG_W6:  "w6",
	REG_W7:  "w7",
	REG_W8:  "w8",
	REG_W9:  "w9",
	REG_W10: "w10",
	REG_W11: "w11",
	REG_W12: "w12",
	REG_W13: "w13",
	REG_W14: "w14",
	REG_W15: "w15",
	REG_W16: "w16",
	REG_W17: "w17",
	REG_W18: "w18",
	REG_W19: "w19",
	REG_W20: "w20",
	REG_W21: "w21",
	REG_W22: "w22",
	REG_W23: "w23",
	REG_W24: "w24",
	REG_W25: "w25",
	REG_W26: "w26",
	REG_W27: "w27",
	REG_W28: "w28",
	REG_W29: "w29",
	REG_W30: "w30",

	// 通用寄存器
	REG_X0:  "x0",
	REG_X1:  "x1",
	REG_X2:  "x2",
	REG_X3:  "x3",
	REG_X4:  "x4",
	REG_X5:  "x5",
	REG_X6:  "x6",
	REG_X7:  "x7",
	REG_X8:  "x8",
	REG_X9:  "x9",
	REG_X10: "x10",
	REG_X11: "x11",
	REG_X12: "x12",
	REG_X13: "x13",
	REG_X14: "x14",
	REG_X15: "x15",
	REG_X16: "x16",
	REG_X17: "x17",
	REG_X18: "x18",
	REG_X19: "x19",
	REG_X20: "x20",
	REG_X21: "x21",
	REG_X22: "x22",
	REG_X23: "x23",
	REG_X24: "x24",
	REG_X25: "x25",
	REG_X26: "x26",
	REG_X27: "x27",
	REG_X28: "x28",
	REG_X29: "x29",
	REG_X30: "x30",

	// 浮点数寄存器
	REG_V0:  "v0",
	REG_V1:  "v1",
	REG_V2:  "v2",
	REG_V3:  "v3",
	REG_V4:  "v4",
	REG_V5:  "v5",
	REG_V6:  "v6",
	REG_V7:  "v7",
	REG_V8:  "v8",
	REG_V9:  "v9",
	REG_V10: "v10",
	REG_V11: "v11",
	REG_V12: "v12",
	REG_V13: "v13",
	REG_V14: "v14",
	REG_V15: "v15",
	REG_V16: "v16",
	REG_V17: "v17",
	REG_V18: "v18",
	REG_V19: "v19",
	REG_V20: "v20",
	REG_V21: "v21",
	REG_V22: "v22",
	REG_V23: "v23",
	REG_V24: "v24",
	REG_V25: "v25",
	REG_V26: "v26",
	REG_V27: "v27",
	REG_V28: "v28",
	REG_V29: "v29",
	REG_V30: "v30",
	REG_V31: "v31",

	// 浮点数寄存器: float32
	REG_S0:  "s0",
	REG_S1:  "s1",
	REG_S2:  "s2",
	REG_S3:  "s3",
	REG_S4:  "s4",
	REG_S5:  "s5",
	REG_S6:  "s6",
	REG_S7:  "s7",
	REG_S8:  "s8",
	REG_S9:  "s9",
	REG_S10: "s10",
	REG_S11: "s11",
	REG_S12: "s12",
	REG_S13: "s13",
	REG_S14: "s14",
	REG_S15: "s15",
	REG_S16: "s16",
	REG_S17: "s17",
	REG_S18: "s18",
	REG_S19: "s19",
	REG_S20: "s20",
	REG_S21: "s21",
	REG_S22: "s22",
	REG_S23: "s23",
	REG_S24: "s24",
	REG_S25: "s25",
	REG_S26: "s26",
	REG_S27: "s27",
	REG_S28: "s28",
	REG_S29: "s29",
	REG_S30: "s30",
	REG_S31: "s31",

	// 浮点数寄存器: float64
	REG_D0:  "d0",
	REG_D1:  "d1",
	REG_D2:  "d2",
	REG_D3:  "d3",
	REG_D4:  "d4",
	REG_D5:  "d5",
	REG_D6:  "d6",
	REG_D7:  "d7",
	REG_D8:  "d8",
	REG_D9:  "d9",
	REG_D10: "d10",
	REG_D11: "d11",
	REG_D12: "d12",
	REG_D13: "d13",
	REG_D14: "d14",
	REG_D15: "d15",
	REG_D16: "d16",
	REG_D17: "d17",
	REG_D18: "d18",
	REG_D19: "d19",
	REG_D20: "d20",
	REG_D21: "d21",
	REG_D22: "d22",
	REG_D23: "d23",
	REG_D24: "d24",
	REG_D25: "d25",
	REG_D26: "d26",
	REG_D27: "d27",
	REG_D28: "d28",
	REG_D29: "d29",
	REG_D30: "d30",
	REG_D31: "d31",
}

// 寄存器名字列表
var _Register32 = []string{
	// 通用寄存器
	REG_X0:  "w0",
	REG_X1:  "w1",
	REG_X2:  "w2",
	REG_X3:  "w3",
	REG_X4:  "w4",
	REG_X5:  "w5",
	REG_X6:  "w6",
	REG_X7:  "w7",
	REG_X8:  "w8",
	REG_X9:  "w9",
	REG_X10: "w10",
	REG_X11: "w11",
	REG_X12: "w12",
	REG_X13: "w13",
	REG_X14: "w14",
	REG_X15: "w15",
	REG_X16: "w16",
	REG_X17: "w17",
	REG_X18: "w18",
	REG_X19: "w19",
	REG_X20: "w20",
	REG_X21: "w21",
	REG_X22: "w22",
	REG_X23: "w23",
	REG_X24: "w24",
	REG_X25: "w25",
	REG_X26: "w26",
	REG_X27: "w27",
	REG_X28: "w28",
	REG_X29: "w29",
	REG_X30: "w30",

	// 浮点数寄存器
	REG_V0:  "s0",
	REG_V1:  "s1",
	REG_V2:  "s2",
	REG_V3:  "s3",
	REG_V4:  "s4",
	REG_V5:  "s5",
	REG_V6:  "s6",
	REG_V7:  "s7",
	REG_V8:  "s8",
	REG_V9:  "s9",
	REG_V10: "s10",
	REG_V11: "s11",
	REG_V12: "s12",
	REG_V13: "s13",
	REG_V14: "s14",
	REG_V15: "s15",
	REG_V16: "s16",
	REG_V17: "s17",
	REG_V18: "s18",
	REG_V19: "s19",
	REG_V20: "s20",
	REG_V21: "s21",
	REG_V22: "s22",
	REG_V23: "s23",
	REG_V24: "s24",
	REG_V25: "s25",
	REG_V26: "s26",
	REG_V27: "s27",
	REG_V28: "s28",
	REG_V29: "s29",
	REG_V30: "s30",
	REG_V31: "s31",
}

// 寄存器别名
var _RegisterAlias = []string{
	REG_SP:  "sp",
	REG_WZR: "wzr",
	REG_XZR: "xzr",

	REG_IP0: "ip0",
	REG_IP1: "ip1",
	REG_FP:  "fp",
	REG_LR:  "lr",
}

// 指令的名字
// 保持和指令定义相同的顺序
var _Anames = []string{
	AADD:    "add",    // R[Rd] = R[Rn] + R[Rm]
	AADDI:   "addi",   // R[Rd] = R[Rn] + ALUImm
	AADDIS:  "addis",  // R[Rd], FLAGS = R[Rn] + ALUImm
	AADDS:   "adds",   // R[Rd], FLAGS = R[Rn] +R[Rm]
	AAND:    "and",    // R[Rd] = R[Rn] & R[Rm]
	AANDI:   "andi",   // R[Rd] = R[Rn] & ALUImm
	AANDIS:  "andis",  // R[Rd], FLAGS = R[Rn] & ALUImm
	AANDS:   "ands",   // R[Rd], FLAGS=R[Rn] & R[Rm]
	AB:      "b",      // PC = PC + BranchAddr
	AB_EQ:   "b.eq",   // PC = PC+ CondBranchAddr if FLAGS equal
	AB_NE:   "b.ne",   // PC = PC+ CondBranchAddr if FLAGS not equal
	AB_LT:   "b.lt",   // PC = PC+ CondBranchAddr if FLAGS less than
	AB_LE:   "b.le",   // PC = PC+ CondBranchAddr if FLAGS less than or equal
	AB_GT:   "b.gt",   // PC = PC+ CondBranchAddr if FLAGS greater than
	AB_GE:   "b.ge",   // PC = PC+ CondBranchAddr if FLAGS greater than or equal
	AB_LO:   "b.lo",   // PC = PC+ CondBranchAddr if FLAGS lower
	AB_LS:   "b.ls",   // PC = PC+ CondBranchAddr if FLAGS lower or same
	AB_HI:   "b.hi",   // PC = PC+ CondBranchAddr if FLAGS higher
	AB_HS:   "b.hs",   // PC = PC+ CondBranchAddr if FLAGS higher or same
	ABL:     "bl",     // R[30] = PC+4; PC = PC+ BranchAddr
	ABR:     "br",     // PC= R[Rt]
	ACBNZ:   "cbnz",   // PC = PC + CondBranchAddr if(R[Rt]!=0)
	ACBZ:    "cbz",    // PC = PC + CondBranchAddr if(R[Rt]==0)
	AEOR:    "eor",    // R[Rd] = R[Rn] ^ R[Rm]
	AEORI:   "eori",   // R[Rd] = R[Rn] ^ ALUImm
	ALDUR:   "ldur",   // R[Rt] = M[R[Rn] + DTAddr]
	ALDURB:  "ldurb",  // R[Rt]={56'bо, M[R[Rn] + DTAddr](7:0)}
	ALDURH:  "ldurh",  // R[Rt]={48'bо, M[R[Rn] + DTAddr](15:0)}
	ALDURSW: "ldursw", // R[Rt]={32{M[R[Rn] + DTAddr][31]}, M[R[Rn] + DTAddr] (31:0)}
	ALDXR:   "ldxr",   // R[Rd] = M[R[Rn] + DTAddr]
	ALSL:    "lsl",    // R[Rd] = R[Rn] << shamt
	ALSR:    "lsr",    // R[Rd] = R[Rn]>>> shamt
	AMOVK:   "movk",   // R[Rd](Instruction[22:21]*16: Instruction[22:21]*16-15)=MOVImm
	AMOVZ:   "movz",   // R[Rd] = {MOVImm <<(Instruction[22:21]*16)}
	AORR:    "orr",    // R[Rd] = R[Rn] | R[Rm]
	AORRI:   "orri",   // R[Rd] = R[Rn] | ALUImm
	ASTUR:   "stur",   // M[R[Rn] + DTAddr] = R[Rt]
	ASTURB:  "sturb",  // M[R[Rn]+DTAddr](7:0) = R[Rt](7:0)
	ASTURH:  "sturh",  // M[R[Rn] + DTAddr](15:0) = R[Rt](15:0)
	ASTURW:  "sturw",  // M[R[Rn] + DTAddr](31:0) = R[Rt](31:0)
	ASTXR:   "stxr",   // M[R[Rn] + DTAddr] = R[Rt]; R[Rm] = (atomic)?0:1
	ASUB:    "sub",    // R[Rd] = R[Rn] - R[Rm]
	ASUBI:   "subi",   // R[Rd] = R[Rn] - ALUImm
	ASUBIS:  "subis",  // R[Rd], FLAGS = R[Rn] - ALUImm
	ASUBS:   "subs",   // R[Rd], FLAGS = R[Rn] - R[Rm]
	AFADDS:  "adds",   // S[Rd] = S[Rn] + S[Rm]
	AFADDD:  "addd",   // D[Rd] = D[Rn] + D[Rm]
	AFCMPS:  "fcmps",  // FLAGS = (S[Rn] vs S[Rm])
	AFCMPD:  "fcmpd",  // FLAGS = (D[Rn] vs D[Rm])
	AFDIVS:  "fdivs",  // S[Rd] = S[Rn]/S[Rm]
	AFDIVD:  "fdivd",  // D[Rd] = D[Rn]/D[Rm]
	AFMULS:  "fmuls",  // S[Rd] = S[Rn] * S[Rm]
	AFMULD:  "fmuld",  // D[Rd] = D[Rn] * D[Rm]
	AFSUBS:  "fsubs",  // S[Rd] = S[Rn] – S[Rm]
	AFSUBD:  "fsubd",  // D[Rd] = D[Rn] - D[Rm]
	ALDURS:  "ldurs",  // S[Rt] = M[R[Rn] + DTAddr]
	ALDURD:  "ldurd",  // D[Rt] = M[R[Rn] + DTAddr]
	AMUL:    "mul",    // R[Rd] = (R[Rn] * R[Rm]) (63:0)
	ASDIV:   "sdiv",   // R[Rd] = R[Rn]/R[Rm]
	ASMULH:  "smulh",  // R[Rd] = (R[Rn] * R[Rm]) (127:64)
	ASTURS:  "sturs",  // M[R[Rn] + DTAddr] = S[Rt]
	ASTURD:  "sturd",  // M[R[Rn] + DTAddr] = D[Rt]
	AUDIV:   "udiv",   // R[Rd] = R[Rn] / R[Rm]
	AUMULH:  "umulh",  // R[Rd] = (R[Rn] * R[Rm]) (127:64)
}
