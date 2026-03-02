// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package arm64

import (
	"wa-lang.org/wa/internal/native/abi"
)

const (
	// 通用寄存器
	_ abi.RegType = iota // 0 是无效的编号

	// 特殊寄存器 (硬件编码31)
	REG_SP
	REG_WZR // zero
	REG_XZR // zero

	// 低32位寄存器
	REG_W0
	REG_W1
	REG_W2
	REG_W3
	REG_W4
	REG_W5
	REG_W6
	REG_W7
	REG_W8
	REG_W9
	REG_W10
	REG_W11
	REG_W12
	REG_W13
	REG_W14
	REG_W15
	REG_W16
	REG_W17
	REG_W18
	REG_W19
	REG_W20
	REG_W21
	REG_W22
	REG_W23
	REG_W24
	REG_W25
	REG_W26
	REG_W27
	REG_W28
	REG_W29
	REG_W30

	// 通用寄存器
	REG_X0
	REG_X1
	REG_X2
	REG_X3
	REG_X4
	REG_X5
	REG_X6
	REG_X7
	REG_X8
	REG_X9
	REG_X10
	REG_X11
	REG_X12
	REG_X13
	REG_X14
	REG_X15
	REG_X16
	REG_X17
	REG_X18
	REG_X19
	REG_X20
	REG_X21
	REG_X22
	REG_X23
	REG_X24
	REG_X25
	REG_X26
	REG_X27
	REG_X28
	REG_X29
	REG_X30

	// 浮点数寄存器
	REG_V0
	REG_V1
	REG_V2
	REG_V3
	REG_V4
	REG_V5
	REG_V6
	REG_V7
	REG_V8
	REG_V9
	REG_V10
	REG_V11
	REG_V12
	REG_V13
	REG_V14
	REG_V15
	REG_V16
	REG_V17
	REG_V18
	REG_V19
	REG_V20
	REG_V21
	REG_V22
	REG_V23
	REG_V24
	REG_V25
	REG_V26
	REG_V27
	REG_V28
	REG_V29
	REG_V30
	REG_V31

	// 浮点数寄存器: float32
	REG_S0
	REG_S1
	REG_S2
	REG_S3
	REG_S4
	REG_S5
	REG_S6
	REG_S7
	REG_S8
	REG_S9
	REG_S10
	REG_S11
	REG_S12
	REG_S13
	REG_S14
	REG_S15
	REG_S16
	REG_S17
	REG_S18
	REG_S19
	REG_S20
	REG_S21
	REG_S22
	REG_S23
	REG_S24
	REG_S25
	REG_S26
	REG_S27
	REG_S28
	REG_S29
	REG_S30
	REG_S31

	// 浮点数寄存器: float64
	REG_D0
	REG_D1
	REG_D2
	REG_D3
	REG_D4
	REG_D5
	REG_D6
	REG_D7
	REG_D8
	REG_D9
	REG_D10
	REG_D11
	REG_D12
	REG_D13
	REG_D14
	REG_D15
	REG_D16
	REG_D17
	REG_D18
	REG_D19
	REG_D20
	REG_D21
	REG_D22
	REG_D23
	REG_D24
	REG_D25
	REG_D26
	REG_D27
	REG_D28
	REG_D29
	REG_D30
	REG_D31

	// 寄存器编号结束
	REG_END

	// 通用寄存器: 别名
	REG_IP0 = REG_X16
	REG_IP1 = REG_X17
	REG_FP  = REG_X29
	REG_LR  = REG_X30
)

// 凹语言用到的部分指令
// 优先支持最小指令集(LEGv8)
const (
	_ abi.As = iota

	AADD    // R[Rd] = R[Rn] + R[Rm]
	AADDI   // R[Rd] = R[Rn] + ALUImm
	AADDIS  // R[Rd], FLAGS = R[Rn] + ALUImm
	AADDS   // R[Rd], FLAGS = R[Rn] +R[Rm]
	AAND    // R[Rd] = R[Rn] & R[Rm]
	AANDI   // R[Rd] = R[Rn] & ALUImm
	AANDIS  // R[Rd], FLAGS = R[Rn] & ALUImm
	AANDS   // R[Rd], FLAGS=R[Rn] & R[Rm]
	AB      // PC = PC + BranchAddr
	AB_EQ   // PC = PC+ CondBranchAddr if FLAGS equal
	AB_NE   // PC = PC+ CondBranchAddr if FLAGS not equal
	AB_LT   // PC = PC+ CondBranchAddr if FLAGS less than
	AB_LE   // PC = PC+ CondBranchAddr if FLAGS less than or equal
	AB_GT   // PC = PC+ CondBranchAddr if FLAGS greater than
	AB_GE   // PC = PC+ CondBranchAddr if FLAGS greater than or equal
	AB_LO   // PC = PC+ CondBranchAddr if FLAGS lower
	AB_LS   // PC = PC+ CondBranchAddr if FLAGS lower or same
	AB_HI   // PC = PC+ CondBranchAddr if FLAGS higher
	AB_HS   // PC = PC+ CondBranchAddr if FLAGS higher or same
	ABL     // R[30] = PC+4; PC = PC+ BranchAddr
	ABR     // PC= R[Rt]
	ACBNZ   // PC = PC + CondBranchAddr if(R[Rt]!=0)
	ACBZ    // PC = PC + CondBranchAddr if(R[Rt]==0)
	AEOR    // R[Rd] = R[Rn] ^ R[Rm]
	AEORI   // R[Rd] = R[Rn] ^ ALUImm
	ALDUR   // R[Rt] = M[R[Rn] + DTAddr]
	ALDURB  // R[Rt]={56'bо, M[R[Rn] + DTAddr](7:0)}
	ALDURH  // R[Rt]={48'bо, M[R[Rn] + DTAddr](15:0)}
	ALDURSW // R[Rt]={32{M[R[Rn] + DTAddr][31]}, M[R[Rn] + DTAddr] (31:0)}
	ALDXR   // R[Rd] = M[R[Rn] + DTAddr]
	ALSL    // R[Rd] = R[Rn] << shamt
	ALSR    // R[Rd] = R[Rn]>>> shamt
	AMOVK   // R[Rd](Instruction[22:21]*16: Instruction[22:21]*16-15)=MOVImm
	AMOVZ   // R[Rd] = {MOVImm <<(Instruction[22:21]*16)}
	AORR    // R[Rd] = R[Rn] | R[Rm]
	AORRI   // R[Rd] = R[Rn] | ALUImm
	ASTUR   // M[R[Rn] + DTAddr] = R[Rt]
	ASTURB  // M[R[Rn]+DTAddr](7:0) = R[Rt](7:0)
	ASTURH  // M[R[Rn] + DTAddr](15:0) = R[Rt](15:0)
	ASTURW  // M[R[Rn] + DTAddr](31:0) = R[Rt](31:0)
	ASTXR   // M[R[Rn] + DTAddr] = R[Rt]; R[Rm] = (atomic)?0:1
	ASUB    // R[Rd] = R[Rn] - R[Rm]
	ASUBI   // R[Rd] = R[Rn] - ALUImm
	ASUBIS  // R[Rd], FLAGS = R[Rn] - ALUImm
	ASUBS   // R[Rd], FLAGS = R[Rn] - R[Rm]
	AFADDS  // S[Rd] = S[Rn] + S[Rm]
	AFADDD  // D[Rd] = D[Rn] + D[Rm]
	AFCMPS  // FLAGS = (S[Rn] vs S[Rm])
	AFCMPD  // FLAGS = (D[Rn] vs D[Rm])
	AFDIVS  // S[Rd] = S[Rn]/S[Rm]
	AFDIVD  // D[Rd] = D[Rn]/D[Rm]
	AFMULS  // S[Rd] = S[Rn] * S[Rm]
	AFMULD  // D[Rd] = D[Rn] * D[Rm]
	AFSUBS  // S[Rd] = S[Rn] – S[Rm]
	AFSUBD  // D[Rd] = D[Rn] - D[Rm]
	ALDURS  // S[Rt] = M[R[Rn] + DTAddr]
	ALDURD  // D[Rt] = M[R[Rn] + DTAddr]
	AMUL    // R[Rd] = (R[Rn] * R[Rm]) (63:0)
	ASDIV   // R[Rd] = R[Rn]/R[Rm]
	ASMULH  // R[Rd] = (R[Rn] * R[Rm]) (127:64)
	ASTURS  // M[R[Rn] + DTAddr] = S[Rt]
	ASTURD  // M[R[Rn] + DTAddr] = D[Rt]
	AUDIV   // R[Rd] = R[Rn] / R[Rm]
	AUMULH  // R[Rd] = (R[Rn] * R[Rm]) (127:64)

	// End marker
	ALAST
)
