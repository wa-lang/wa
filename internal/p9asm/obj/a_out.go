// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package obj

import (
	"fmt"
	"strconv"
)

// 各个平台通用的指令
// 平台特殊的指令从 A_ARCHSPECIFIC 开始定义
const (
	AXXX = 0 + iota
	ACALL
	ACHECKNIL
	ADATA
	ADUFFCOPY
	ADUFFZERO
	AEND
	AFUNCDATA
	AGLOBL
	AJMP
	ANOP
	APCDATA
	ARET
	ATEXT
	ATYPE
	AUNDEF
	AUSEFIELD
	AVARDEF
	AVARKILL
	A_ARCHSPECIFIC
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
const (
	REG_NONE     = 0 // 寄存器编号为空
	RBase386     = 1 * 1024
	RBaseAMD64   = 2 * 1024
	RBaseARM     = 3 * 1024
	RBaseARM64   = 8 * 1024  // range [8k, 13k)
	RBaseRISCV   = 15 * 1024 // range [15k, 16k)
	RBaseLOONG64 = 19 * 1024 // range [19K, 22k)
)

// 机器码和对应的名字
type opSet struct {
	lo    int
	names []string
}

// 机器码区间集合
var aSpace []opSet

// 注册不同平台的指令区间
func RegisterOpcode(lo int, Anames []string) {
	aSpace = append(aSpace, opSet{lo, Anames})
}

// 指令机器码转为字符串名字
func Aconv(a int) string {
	if a < A_ARCHSPECIFIC {
		return Anames[a]
	}
	for i := range aSpace {
		as := &aSpace[i]
		if as.lo <= a && a < as.lo+len(as.names) {
			return as.names[a-as.lo]
		}
	}
	return fmt.Sprintf("A???%d", a)
}

// 寄存器区间
type regSet struct {
	lo    int              // 开始
	hi    int              // 结束(开区间)
	Rconv func(int) string // 用于打印
}

// 用于注册不同架构下的寄存器
// 不同架构的寄存器编号处于不同空间不会冲突
var regSpace []regSet

// 注册不同平台的寄存器区间
func RegisterRegister(lo, hi int, Rconv func(int) string) {
	regSpace = append(regSpace, regSet{lo, hi, Rconv})
}

// 将寄存器编号转为字符串名字
func Rconv(reg int) string {
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

// executable header types

// 可执行文件头类型
type HExeType int

const (
	Hunknown HExeType = 0 + iota
	Hdarwin
	Helf
	Hlinux
	Hwindows
)

var headers = []struct {
	name string
	val  HExeType
}{
	{"darwin", Hdarwin},
	{"elf", Helf},
	{"linux", Hlinux},
	{"android", Hlinux}, // must be after "linux" entry or else headstr(Hlinux) == "android"
	{"windows", Hwindows},
	{"windowsgui", Hwindows},
}

func Headtype(name string) HExeType {
	for i := 0; i < len(headers); i++ {
		if name == headers[i].name {
			return headers[i].val
		}
	}
	return -1
}

func (v HExeType) String() string {
	for i := 0; i < len(headers); i++ {
		if v == headers[i].val {
			return headers[i].name
		}
	}
	return strconv.Itoa(int(v))
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

// For the linkers. Must match Wa definitions.
// TODO(chai2010): Share Wa definitions with linkers directly.

// chaishushan: 临时增加, 用于跳过构建错误
const stackGuardMultiplier = 1

const (
	STACKSYSTEM = 0
	StackSystem = STACKSYSTEM
	StackBig    = 4096
	StackGuard  = 640*stackGuardMultiplier + StackSystem
	StackSmall  = 128
	StackLimit  = StackGuard - StackSystem - StackSmall
)

const (
	StackPreempt = -1314 // 0xfff...fade
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

// Auto.name
const (
	A_AUTO = 1 + iota
	A_PARAM
)

// Reloc.type
// 最终的值是bit位组合
type RelocType int32

const (
	R_ADDR RelocType = 1 + iota
	R_ADDRPOWER
	R_ADDRARM64
	R_SIZE
	R_CALL
	R_CALLARM
	R_CALLARM64
	R_CALLIND
	R_CALLPOWER
	R_CONST
	R_PCREL
	// R_TLS (only used on arm currently, and not on android and darwin where tlsg is
	// a regular variable) resolves to data needed to access the thread-local g. It is
	// interpreted differently depending on toolchain flags to implement either the
	// "local exec" or "inital exec" model for tls access.
	// TODO(chai2010): change to use R_TLS_LE or R_TLS_IE as appropriate, not having
	// R_TLS do double duty.
	R_TLS
	// R_TLS_LE (only used on 386 and amd64 currently) resolves to the offset of the
	// thread-local g from the thread local base and is used to implement the "local
	// exec" model for tls access (r.Sym is not set by the compiler for this case but
	// is set to Tlsg in the linker when externally linking).
	R_TLS_LE
	// R_TLS_IE (only used on 386 and amd64 currently) resolves to the PC-relative
	// offset to a GOT slot containing the offset the thread-local g from the thread
	// local base and is used to implemented the "initial exec" model for tls access
	// (r.Sym is not set by the compiler for this case but is set to Tlsg in the
	// linker when externally linking).
	R_TLS_IE
	R_GOTOFF
	R_PLT0
	R_PLT1
	R_PLT2
	R_USEFIELD
	R_POWER_TOC
	R_GOTPCREL
)

// LSym.type
// 符号(LSym)的类型常量
type SymKind int16

const (
	Sxxx              SymKind = iota // 无效或未使用
	STEXT                            // 函数代码().text段)
	SELFRXSECT                       // ELF 自定义可执行段
	STYPE                            // 类型信息(reflect.Type等)
	SSTRING                          // 字符串常量
	SWASTRING                        // Write barrier 相关字符串
	SWAFUNC                          // Write barrier 相关函数
	SGCBITS                          // GC bitmap 数据段
	SRODATA                          // 只读数据段
	SFUNCTAB                         // 函数表(调试/调用信息)
	STYPELINK                        // 类型链接表
	SSYMTAB                          // 符号表(调试信息)
	SPCLNTAB                         // pcln 表, 用于定位源代码行
	SELFROSECT                       // ELF 自定义只读段
	SMACHOPLT                        // Mach-O 的 PLT 表(用于动态链接)
	SELFSECT                         // ELF 自定义数据段
	SMACHO                           // Mach-O 文件相关段
	SMACHOGOT                        // Mach-O GOT 表(全局偏移表)
	SWINDOWS                         // Windows 特定符号段
	SELFGOT                          // ELF GOT 表
	SNOPTRDATA                       // 无指针数据段(适用于 GC 优化)
	SINITARR                         // 初始化数组段(例如用于运行时init)
	SDATA                            // 普通数据段(含指针)
	SBSS                             // 未初始化全局变量段
	SNOPTRBSS                        // 未初始化且无指针变量段
	STLSBSS                          // TLS BSS段(线程局部存储)
	SXREF                            // 外部引用(尚未解析符号)
	SMACHOSYMSTR                     // Mach-O 符号字符串表
	SMACHOSYMTAB                     // Mach-O 符号表
	SMACHOINDIRECTPLT                // Mach-O 间接PLT表
	SMACHOINDIRECTGOT                // Mach-O 间接GOT表
	SFILE                            // 用于 DWARF 的源文件名
	SFILEPATH                        // DWARF 路径常量池
	SCONST                           // 编译器常量数据(例如小整数, 浮点数)
	SDYNIMPORT                       // 动态库导入符号
	SHOSTOBJ                         // C 编译目标文件中导入的符号

	SSUB       = 1 << 8   // 子符号标记(用于函数内嵌符号, 如 DWARF 子项)
	SMASK      = SSUB - 1 // 掩码，提取主类型时用
	SHIDDEN    = 1 << 9   // 符号隐藏(如不导出符号)
	SCONTAINER = 1 << 10  // 包含子符号(如容器符号)
)

// 符号(LSym)名字类别定义
type NameType int8

const (
	NAME_NONE   NameType = 0 + iota
	NAME_EXTERN          // 外部符号, 对应 ELF中的 STB_GLOBAL
	NAME_STATIC          // 文件级私有静态符号, 对应 ELF 中的 STB_LOCAL
	NAME_AUTO            // 函数内部的自动变量(栈上变量, 较少直接暴露), 但在 debug 或 backend 编译阶段会存在
	NAME_PARAM           // 函数的参数(也可能映射到栈上的位置), 和 NAME_AUTO 类似, 是调试信息构建时区分符号用途的关键字段

	// 对 name@GOT(SB) 全局偏移表(GOT)中某个符号地址
	// 主要用于连接时生成位置无关代码(PIC), 对应 ELF 中的 Global Offset Table
	// 比如在 AMD64 Linux 下如果引用了某个外部变量, 它在编译为动态库时可能会被变成 GOT 引用形式
	// NAME_GOTREF 名字是 GOT REF
	NAME_GOTREF
)

// 操作数类型常量
type TypeType int16

const (
	TYPE_NONE     TypeType = 0        // 操作数不存在
	TYPE_BRANCH   TypeType = 5 + iota // 跳转地址, 例如 JMP target
	TYPE_TEXTSIZE                     // 用于 .text 段中设置函数大小的标记(用于链接器调试或优化目的)
	TYPE_MEM                          // 内存访问, 例如 [R1+4] 在 MOV 等指令中用于读取或写入内存
	TYPE_CONST                        // 整型常量值, 例如 MOV $1, R0 中的 $1
	TYPE_FCONST                       // 浮点常量, 例如 $3.14
	TYPE_SCONST                       // 字符串常量, 例 .string "hello"
	TYPE_REG                          // 单一寄存器, 如 AX, R1, X0
	TYPE_ADDR                         // 地址引用, 比如 MOV $foo(SB), AX 中的 foo(SB). 这是静态地址访问, 常用于全局变量
	TYPE_SHIFT                        // 表示移位操作, 常见于 ARM 或其他 RISC 架构. 如 LSL R1, R2, #3 (逻辑左移)
	TYPE_REGREG                       // 两个寄存器, 例如 CMP R1, R2
	TYPE_REGREG2                      // 用于三寄存器形式的指令, 或者链接器特殊用途
	TYPE_INDIR                        // 间接寻址, 如 *(R1) 或 [R1]. 可与 TYPE_MEM 组合使用?
	TYPE_REGLIST                      // 寄存器列表, 常用于 PUSH {R1-R4} 或 POP 操作, 在 ARM 架构中常见
)
