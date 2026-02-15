// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package abi

// 默认的入口函数
const (
	DefaultEntryFunc   = "_start"
	DefaultEntryFuncZh = "_启动"
)

// CPU类型
// 对应不同的汇编指令
type CPUType int

const (
	CPU_Nil CPUType = iota
	LOONG64
	RISCV64
	RISCV32
	X64Unix    // UNIX System V ABI
	X64Windows // Windows ABI
	CPU_Max
)

var _CPUType_strings = []string{
	LOONG64:    "loong64",
	RISCV64:    "riscv64",
	RISCV32:    "riscv32",
	X64Unix:    "X64-Unix", // "x64"
	X64Windows: "x64-Windows",
}

var _CPUType_x64Unix_strings = []string{
	"x64", "amd64",
}

// 操作系统
type OSType int

const (
	OS_Nil OSType = iota
	LINUX
	WINDOWS
)

var _OSType_strings = []string{
	LINUX:   "linux",
	WINDOWS: "windows",
}

// 链接参数
type LinkOptions struct {
	OS        OSType  // OS类型
	CPU       CPUType // CPU类型
	DRAMBase  int64   // 代码段开始地址
	DRAMSize  int64   // 数据段开始地址
	EntryFunc string  // 入口函数
}

// 链接后的符号
type LinkedSymbol struct {
	Name         string // 名字
	Addr         int64  // 内存地址
	Size         int64  // 内存大小
	Data         []byte // 内存数据
	AlignPadding int    // 开头对齐需要的填充数
}

// 程链接后的程序
type LinkedProgram struct {
	CPU CPUType // CPU类型

	Entry    int64  // 入口地址
	TextAddr int64  // 程序段地址
	TextData []byte // 程序段数据
	DataAddr int64  // 数据段地址
	DataData []byte // 数据段数据

	X64ImageBase int64 // 镜像基础地址
}

// 寄存器类型
type RegType int16

// 指令类型
type As int16

// 内置的重定位函数
type BuiltinFn int16

const (
	BuiltinFn_Nil BuiltinFn = iota

	// =============== 开始 ===============

	// 龙芯内置宏(基础)
	BuiltinFn_ABS_LO12   // %abs_lo12(symbol)
	BuiltinFn_ABS_HI20   // %abs_hi20(symbol)
	BuiltinFn_ABS64_LO20 // %abs_lo20(symbol)
	BuiltinFn_ABS64_HI12 // %abs_hi12(symbol)
	BuiltinFn_PC_LO12    // %pc_lo12(symbol)
	BuiltinFn_PC_HI20    // %pc_hi20(symbol)
	BuiltinFn_PC64_LO20  // %pc_lo20(symbol)
	BuiltinFn_PC64_HI12  // %pc_hi12(symbol)

	// RISCV 定义的宏(扩展)
	BuiltinFn_HI       // %hi(symbol)
	BuiltinFn_LO       // %lo(symbol)
	BuiltinFn_PCREL_HI // %pcrel_hi(symbol)
	BuiltinFn_PCREL_LO // %pcrel_lo(label|symbol) # 注意, RISCV 中该参数是 label, symbol 是凹汇编器扩展

	// 通用宏
	BuiltinFn_SIZEOF // %sizeof(symbol) # 获取全局变量的内存大小

	// =============== 以下是对应的中文定义 ===============

	// 龙芯内置宏(基础)
	BuiltinFn_ABS_LO12_zh   // %绝对.低12(符号)
	BuiltinFn_ABS_HI20_zh   // %绝对.高20(符号)
	BuiltinFn_ABS64_LO20_zh // %绝对64.低20(符号)
	BuiltinFn_ABS64_HI12_zh // %绝对64.高12(符号)
	BuiltinFn_PC_LO12_zh    // %相对.低12(符号)
	BuiltinFn_PC_HI20_zh    // %相对.高20(符号)
	BuiltinFn_PC64_LO20_zh  // %相对64.低20(符号)
	BuiltinFn_PC64_HI12_zh  // %相对64.高12(符号)

	// RISCV 定义的宏(扩展)
	BuiltinFn_HI_zh       // %高位(符号)
	BuiltinFn_LO_zh       // %低位(符号)
	BuiltinFn_PCREL_HI_zh // %相对高位(符号)
	BuiltinFn_PCREL_LO_zh // %相对低位(标签) # 注意, RISCV 中该参数是标签

	BuiltinFn_SIZEOF_zh // %内存字节数

	// =============== 结束 ===============

	BuiltinFn_Max // 标记结束
)

var _BuiltinFn_strings = []string{
	BuiltinFn_ABS_LO12:      "%abs_lo12",
	BuiltinFn_ABS_HI20:      "%abs_hi20",
	BuiltinFn_ABS64_LO20:    "%abs_lo20",
	BuiltinFn_ABS64_HI12:    "%abs_hi12",
	BuiltinFn_PC_LO12:       "%pc_lo12",
	BuiltinFn_PC_HI20:       "%pc_hi20",
	BuiltinFn_PC64_LO20:     "%pc_lo20",
	BuiltinFn_PC64_HI12:     "%pc_hi12",
	BuiltinFn_HI:            "%hi",
	BuiltinFn_LO:            "%lo",
	BuiltinFn_PCREL_HI:      "%pcrel_hi",
	BuiltinFn_PCREL_LO:      "%pcrel_lo",
	BuiltinFn_SIZEOF:        "%sizeof",
	BuiltinFn_ABS_LO12_zh:   "%绝对.低12",
	BuiltinFn_ABS_HI20_zh:   "%绝对.高20",
	BuiltinFn_ABS64_LO20_zh: "%绝对64.低20",
	BuiltinFn_ABS64_HI12_zh: "%绝对64.高12",
	BuiltinFn_PC_LO12_zh:    "%相对.低12",
	BuiltinFn_PC_HI20_zh:    "%相对.高20",
	BuiltinFn_PC64_LO20_zh:  "%相对64.低20",
	BuiltinFn_PC64_HI12_zh:  "%相对64.高12",
	BuiltinFn_HI_zh:         "%高位",
	BuiltinFn_LO_zh:         "%低位",
	BuiltinFn_PCREL_HI_zh:   "%相对高位",
	BuiltinFn_PCREL_LO_zh:   "%相对低位",
	BuiltinFn_SIZEOF_zh:     "%内存字节数",
}

// 指令参数
type AsArgument struct {
	Rd  RegType // 目标寄存器
	Rs1 RegType // 原寄存器1
	Rs2 RegType // 原寄存器2
	Rs3 RegType // 原寄存器3
	Imm int32   // 立即数

	// 参数的名字, 用于格式化
	RdName  string
	Rs1Name string
	Rs2Name string
	Rs3Name string

	Symbol      string    // 可能是 Label/全局符号, 用于重定位和输出文本(不支持Imm2)
	SymbolDecor BuiltinFn // 符号的修饰函数, 可能要重新计算(只用于跳转/加立即数部分)
}

// 指令原生参数
type AsRawArgument struct {
	Rd  uint32 // 目标寄存器
	Rs1 uint32 // 原寄存器1
	Rs2 uint32 // 原寄存器2
	Rs3 uint32 // 原寄存器3
	Imm int32  // 立即数
}
