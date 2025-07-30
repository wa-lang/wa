// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package obj

import (
	"fmt"
)

// 函数的参数/局部变量/帧大小信息
const (
	PCDATA_StackMapIndex       = 0
	FUNCDATA_ArgsPointerMaps   = 0
	FUNCDATA_LocalsPointerMaps = 1
	FUNCDATA_DeadValueMaps     = 2
	ArgsSizeUnknown            = -0x80000000 // 帧大小未知, 比如 C 的 printf 函数
)

// 指令机器码
type As int16

// 各个平台通用的指令
// 平台特殊的指令从 A_ARCHSPECIFIC 开始定义
// TODO(chai2010): 精简伪指令
const (
	AXXX           As = iota // 无效或未初始化的指令
	ACALL                    // 调用函数
	ACHECKNIL                // 空指针检查, 用于 runtime 插入的 nil-check
	ADATA                    // 静态数据段的数据定义
	ADUFFCOPY                // Duff's device 复制优化入口，用于快速 memcopy
	ADUFFZERO                // Duff's device 清零优化入口
	AEND                     // 汇编文件结尾标志(需要吗?)
	AFUNCDATA                // 函数级的元信息注入, 常见的是 gcmap、defer info 等
	AGLOBL                   // 全局变量定义(类似于 .globl)
	AJMP                     // 无条件跳转指令
	ANOP                     // 空操作指令, 用于填充, 对齐或占位
	APCDATA                  // 异常栈, 调试元信息(比如 PC 到 stack map 的映射)
	ARET                     // 函数返回指令
	ATEXT                    // 函数定义入口标记, 指定函数名和属性
	ATYPE                    // 类型信息
	AUNDEF                   // 未定义的操作, 或执行到这里就崩溃(像 trap)
	AUSEFIELD                // 用于反射优化, 标记 struct field 被使用
	AVARDEF                  // 标记局部变量的生命周期开始(调试, GC 用)
	AVARKILL                 // 标记局部变量生命周期结束
	A_ARCHSPECIFIC           // 架构专属操作码的起点
)

// 指令的名字
var Anames = []string{
	AXXX:      "XXX",
	ACALL:     "CALL",
	ACHECKNIL: "CHECKNIL",
	ADATA:     "DATA",
	ADUFFCOPY: "DUFFCOPY",
	ADUFFZERO: "DUFFZERO",
	AEND:      "END",
	AFUNCDATA: "FUNCDATA",
	AGLOBL:    "GLOBL",
	AJMP:      "JMP",
	ANOP:      "NOP",
	APCDATA:   "PCDATA",
	ARET:      "RET",
	ATEXT:     "TEXT",
	ATYPE:     "TYPE",
	AUNDEF:    "UNDEF",
	AUSEFIELD: "USEFIELD",
	AVARDEF:   "VARDEF",
	AVARKILL:  "VARKILL",
}

// 每个平台有独立的指令空间
// 比如 ABaseAMD64 + A_ARCHSPECIFIC 开始的是 AMD64 特有的指令
const (
	ABase386 = (1 + iota) << 11
	ABaseARM
	ABaseAMD64
	ABaseARM64
	ABaseLoong64
	ABaseRISCV

	AllowedOpCodes = 1 << 11            // The number of opcodes available for any given architecture.
	AMask          = AllowedOpCodes - 1 // AND with this to use the opcode as an array index.
)

// 每个平台寄存也有独立的空间
// 比如 RBaseAMD64 开始的是 AMD64 平台的特有的寄存器
// 每个平台寄存器范围不超过 1k

type RBaseType int16

const (
	REG_NONE RBaseType = iota * 1024 // 寄存器编号为空
	RBase386
	RBaseAMD64
	RBaseARM
	RBaseARM64
	RBaseRISCV
	RBaseLOONG64
)

// 机器码和对应的名字
type opSet struct {
	lo    As
	names []string
}

// 机器码区间集合
var aSpace []opSet

// 注册不同平台的指令区间
func RegisterOpcode(lo As, anames []string) {
	aSpace = append(aSpace, opSet{lo, anames})
}

// 指令机器码转为字符串名字
func (a As) String() string {
	if a < A_ARCHSPECIFIC {
		return Anames[a]
	}
	for i := range aSpace {
		as := &aSpace[i]
		if as.lo <= a && a < as.lo+As(len(as.names)) {
			return as.names[a-as.lo]
		}
	}
	return fmt.Sprintf("A???%d", a)
}

// 寄存器区间
type regSet struct {
	lo    RBaseType              // 开始
	hi    RBaseType              // 结束(开区间)
	Rconv func(RBaseType) string // 用于打印
}

// 用于注册不同架构下的寄存器
// 不同架构的寄存器编号处于不同空间不会冲突
var regSpace []regSet

// 注册不同平台的寄存器区间
func RegisterRegister(lo, hi RBaseType, Rconv func(RBaseType) string) {
	regSpace = append(regSpace, regSet{lo, hi, Rconv})
}

// 将寄存器编号转为字符串名字
func (reg RBaseType) String() string {
	if reg == REG_NONE {
		return "NONE"
	}
	for i := range regSpace {
		rs := &regSpace[i]
		if rs.lo <= reg && reg < rs.hi {
			return rs.Rconv(reg)
		}
	}
	return fmt.Sprintf("R???%d", reg)
}

// 可执行文件头类型
type HeadType int

const (
	Hunknown HeadType = iota
	Helf
	Hdarwin
	Hlinux
	Hwindows
)

func (h *HeadType) Set(name string) error {
	switch name {
	case "darwin", "ios":
		*h = Hdarwin
	case "elf":
		*h = Helf
	case "linux", "android":
		*h = Hlinux
	case "windows":
		*h = Hwindows
	default:
		return fmt.Errorf("invalid headtype: %q", name)
	}
	return nil
}

func (v HeadType) String() string {
	switch v {
	case Hdarwin:
		return "darwin"
	case Helf:
		return "elf"
	case Hlinux:
		return "linux"
	case Hwindows:
		return "windows"
	}
	return fmt.Sprintf("HeadType(%d)", int(v))
}

// ARM scond byte
const (
	C_SCOND     = (1 << 4) - 1
	C_SBIT      = 1 << 4
	C_PBIT      = 1 << 5
	C_WBIT      = 1 << 6
	C_FBIT      = 1 << 7
	C_UBIT      = 1 << 7
	C_SCOND_XOR = 14
)

// ARM条件指令
var armCondCode = []string{
	".EQ",
	".NE",
	".CS",
	".CC",
	".MI",
	".PL",
	".VS",
	".VC",
	".HI",
	".LS",
	".GE",
	".LT",
	".GT",
	".LE",
	"",
	".NV",
}

// CConv formats ARM condition codes.
func CConv(s uint8) string {
	if s == 0 {
		return ""
	}
	sc := armCondCode[(s&C_SCOND)^C_SCOND_XOR]
	if s&C_SBIT != 0 {
		sc += ".S"
	}
	if s&C_PBIT != 0 {
		sc += ".P"
	}
	if s&C_WBIT != 0 {
		sc += ".W"
	}
	if s&C_UBIT != 0 { // ambiguous with FBIT
		sc += ".U"
	}
	return sc
}

const (
	StackBig   = 4096 // 大帧判断阈值
	StackSmall = 128  // 最小预留空间
	StackGuard = 640  // 栈溢出保护区

	StackLimit = StackGuard - StackSmall // 剩余触发点, 可触发自动报告或任务失败处理
)

// Wa obj 文件魔数
// Writeobjdirect 函数使用
const (
	MagicHeader = "\x00\x00wa01ld"
	MagicFooter = "\xff\xffwa01ld"

	MagicSymbolStart = 0xfe
	MagicFooterStart = 0xff
)

// symbol version, incremented each time a file is loaded.
// version==1 is reserved for savehist.
const (
	HistVersion = 1
)

// 函数不再包含标志位, 默认都是固定栈

const (
	DUPOK  = 1 // 可以出现多个重名符号, 取第一个
	RODATA = 2 // 只读数据段
	NOPTR  = 4 // 不包含指针的数据
)
