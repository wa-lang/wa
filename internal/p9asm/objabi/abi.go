// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package objabi

import "math"

// 可执行文件头类型
type HeadType int

const (
	Hunknown HeadType = iota
	Helf
	Hdarwin
	Hlinux
	Hwindows
)

// CPU类型
// 对应不同的汇编指令
type CPUType int

const (
	X386 CPUType = iota + 1
	AMD64
	ARM
	ARM64
	Loong64
	RISCV
)

// 函数不再包含标志位, 默认都是固定栈

const (
	DUPOK  = 1 // 可以出现多个重名符号, 取第一个
	RODATA = 2 // 只读数据段
	NOPTR  = 4 // 不包含指针的数据
)

// 函数的参数/局部变量/帧大小信息
const (
	PCDATA_StackMapIndex       = 0
	FUNCDATA_ArgsPointerMaps   = 0
	FUNCDATA_LocalsPointerMaps = 1
	FUNCDATA_DeadValueMaps     = 2
)

// 指令机器码
type As int16

// 各个平台通用的指令
// 平台特殊的指令从 A_ARCHSPECIFIC 开始定义
// TODO(chai2010): 精简伪指令
const (
	AXXX           As = 0            // 无效或未初始化的指令
	ACALL          As = ABase + iota // 调用函数(ABase之前是给普通Token预留的空间)
	ACHECKNIL                        // 空指针检查, 用于 runtime 插入的 nil-check
	ADATA                            // 静态数据段的数据定义
	ADUFFCOPY                        // Duff's device 复制优化入口，用于快速 memcopy
	ADUFFZERO                        // Duff's device 清零优化入口
	AEND                             // 汇编文件结尾标志(需要吗?)
	AFUNCDATA                        // 函数级的元信息注入, 常见的是 gcmap、defer info 等
	AGLOBL                           // 全局变量定义(类似于 .globl)
	AJMP                             // 无条件跳转指令
	ANOP                             // 空操作指令, 用于填充, 对齐或占位
	APCDATA                          // 异常栈, 调试元信息(比如 PC 到 stack map 的映射)
	ARET                             // 函数返回指令
	ATEXT                            // 函数定义入口标记, 指定函数名和属性
	ATYPE                            // 类型信息
	AUNDEF                           // 未定义的操作, 或执行到这里就崩溃(像 trap)
	AUSEFIELD                        // 用于反射优化, 标记 struct field 被使用
	AVARDEF                          // 标记局部变量的生命周期开始(调试, GC 用)
	AVARKILL                         // 标记局部变量生命周期结束
	A_ARCHSPECIFIC                   // 架构专属操作码的起点
)

// 每个平台有独立的指令空间
// 比如 ABaseAMD64 + A_ARCHSPECIFIC 开始的是 AMD64 特有的指令
const (
	ABase    = 100 // 前的空间保留给普通的 Token
	ABase386 = (1 + iota) << 11
	ABaseARM
	ABaseAMD64
	ABaseARM64
	ABaseLoong64
	ABaseRISCV
	ABaseMax

	AllowedOpCodes = 1 << 11            // The number of opcodes available for any given architecture.
	AMask          = AllowedOpCodes - 1 // AND with this to use the opcode as an array index.
)

// 每个平台寄存也有独立的空间
// 比如 RBaseAMD64 开始的是 AMD64 平台的特有的寄存器
// 每个平台寄存器范围不超过 1k

type RBaseType int32

const (
	REG_NONE RBaseType = 0                          // 寄存器编号为空
	RBase    RBaseType = math.MaxUint16 + iota*1024 // int16 部分保留给汇编指令
	RBase386
	RBaseAMD64
	RBaseARM
	RBaseARM64
	RBaseRISCV
	RBaseLOONG64
	RBaseMax
)
