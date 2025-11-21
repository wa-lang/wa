// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package loong64

const (
	REG_R0 = iota // must be a multiple of 32
	REG_R1
	REG_R2
	REG_R3
	REG_R4
	REG_R5
	REG_R6
	REG_R7
	REG_R8
	REG_R9
	REG_R10
	REG_R11
	REG_R12
	REG_R13
	REG_R14
	REG_R15
	REG_R16
	REG_R17
	REG_R18
	REG_R19
	REG_R20
	REG_R21
	REG_R22
	REG_R23
	REG_R24
	REG_R25
	REG_R26
	REG_R27
	REG_R28
	REG_R29
	REG_R30
	REG_R31

	REG_F0 // must be a multiple of 32
	REG_F1
	REG_F2
	REG_F3
	REG_F4
	REG_F5
	REG_F6
	REG_F7
	REG_F8
	REG_F9
	REG_F10
	REG_F11
	REG_F12
	REG_F13
	REG_F14
	REG_F15
	REG_F16
	REG_F17
	REG_F18
	REG_F19
	REG_F20
	REG_F21
	REG_F22
	REG_F23
	REG_F24
	REG_F25
	REG_F26
	REG_F27
	REG_F28
	REG_F29
	REG_F30
	REG_F31

	REGZERO = REG_R0 // set to zero
	REGLINK = REG_R1
	REGSP   = REG_R3
	REGRET  = REG_R20 // not use
	REGARG  = -1      // -1 disables passing the first argument in register
	REGRT1  = REG_R20 // reserved for runtime, duffzero and duffcopy
	REGRT2  = REG_R21 // reserved for runtime, duffcopy
	REGCTXT = REG_R29 // context for closures
	REGG    = REG_R22 // G in loong64
	REGTMP  = REG_R30 // used by the assembler
	FREGRET = REG_F0  // not use
)

const (
	// ----------------------------------------------------
	// # I. 核心整数运算 (Core Arithmetic and Logic)
	// ----------------------------------------------------
	AADD   = iota + 1 // 32-bit Add
	AADDV             // 64-bit Add (AADD.V)
	AADDU             // 32-bit Add Unsigned
	AADDVU            // 64-bit Add Unsigned (AADD.VU)
	ASUB              // 32-bit Subtract
	ASUBV             // 64-bit Subtract (ASUB.V)
	ASUBU             // 32-bit Subtract Unsigned
	ASUBVU            // 64-bit Subtract Unsigned (ASUB.VU)

	// 逻辑运算 (Logic)
	AAND  // Logical AND (32-bit, LA64下可用于64位)
	AANDV // 64 位逻辑 AND (可选，如果 AAND 默认非 64 位)
	AOR   // Logical OR
	AXOR  // Logical XOR
	ANOR  // Logical NOR

	// ----------------------------------------------------
	// # II. 乘除法 (Multiply and Divide)
	// ----------------------------------------------------
	AMUL    // 32-bit Multiply
	AMULV   // 64-bit Multiply (AMUL.V)
	AMULU   // 32-bit Multiply Unsigned
	AMULHV  // 64-bit Multiply High (Signed)
	AMULHVU // 64-bit Multiply High Unsigned

	ADIV   // 32-bit Divide
	ADIVV  // 64-bit Divide (ADIV.V)
	ADIVU  // 32-bit Divide Unsigned
	ADIVVU // 64-bit Divide Unsigned
	AREM   // 32-bit Remainder
	AREMV  // 64-bit Remainder
	AREMU  // 32-bit Remainder Unsigned
	AREMVU // 64-bit Remainder Unsigned

	// ----------------------------------------------------
	// # III. 移位 (Shift)
	// ----------------------------------------------------
	ASLL   // 32-bit Shift Left Logical
	ASLLV  // 64-bit Shift Left Logical
	ASRA   // 32-bit Shift Right Arithmetic
	ASRAV  // 64-bit Shift Right Arithmetic
	ASRL   // 32-bit Shift Right Logical
	ASRLV  // 64-bit Shift Right Logical
	AROTR  // 32-bit Rotate Right
	AROTRV // 64-bit Rotate Right

	// ----------------------------------------------------
	// # IV. 数据传输/内存访问 (Load/Store & Move)
	// ----------------------------------------------------
	ALUI   // Load Upper Immediate (用于构建大立即数)
	AMOVB  // Move/Load/Store Byte
	AMOVBU // Move/Load/Store Byte Unsigned
	AMOVH  // Move/Load/Store Halfword
	AMOVHU // Move/Load/Store Halfword Unsigned
	AMOVW  // Move/Load/Store Word (32-bit)
	AMOVV  // Move/Load/Store Value (64-bit)

	// 原子和独占加载/存储 (Atomic & Exclusive)
	ALL  // 32-bit Load Linked
	ALLV // 64-bit Load Linked
	ASC  // 32-bit Store Conditional
	ASCV // 64-bit Store Conditional

	// ----------------------------------------------------
	// # V. 控制流 (Branches and Jumps)
	// ----------------------------------------------------
	// 条件分支 (Condition Branches)
	ABEQ  // Branch if Equal
	ABNE  // Branch if Not Equal
	ABLGT // Branch if Less Than (signed)
	ABLTU // Branch if Less Than Unsigned
	ABGE  // Branch if Greater than or Equal (signed)
	ABGEU // Branch if Greater than or Equal Unsigned

	// 与零比较的分支 (Branches vs Zero)
	ABGEZ // Branch if Greater than or Equal to Zero
	ABLEZ // Branch if Less than or Equal to Zero
	ABGTZ // Branch if Greater than Zero
	ABLTZ // Branch if Less than Zero

	// 跳转 (Jumps)
	AJIRL // Jump and Link Register (函数调用/返回)

	// ----------------------------------------------------
	// # VI. 比较、杂项与伪指令 (Misc)
	// ----------------------------------------------------
	// 比较 (Set/Mask/Trap)
	ASGT     // Set if Greater Than
	ASGTU    // Set if Greater Than Unsigned
	ACMPV    // 64 位比较（占位或别名）
	AMASKEQZ // Mask if Equal to Zero
	AMASKNEZ // Mask if Not Equal to Zero
	ATEQ     // Trap if Equal
	ATNE     // Trap if Not Equal

	// 系统与同步
	ASYSCALL // System Call (程序与OS交互)
	ADBAR    // Data Barrier (内存同步)
	ABREAK   // Software breakpoint/Exception

	// 伪指令/优化相关
	ANOOP // No Operation (硬件 NOP)
	AWORD // 伪指令：定义一个 32 位数据

	// PC 相对寻址
	ALU12IW    // Load Upper Immediate 12-bit
	APCALAU12I // PC-relative Load Upper Immediate
	APCADDU12I // PC-relative Add Upper Immediate

	// 符号扩展与位操作 (Sign/Bit Manipulation)
	AEXTWB // Extend Word from Byte
	AEXTWH // Extend Word from Halfword
	ACLW   // Count Leading Zeros Word
	ACLZW  // Count Leading Ones Word
	ACTZW  // Count Trailing Zeros Word

	ALAST // 指令集结束标记
)
