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
const (
	R_ADDR = 1 + iota
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
const (
	Sxxx = iota
	STEXT
	SELFRXSECT
	STYPE
	SSTRING
	SWASTRING
	SWAFUNC
	SGCBITS
	SRODATA
	SFUNCTAB
	STYPELINK
	SSYMTAB
	SPCLNTAB
	SELFROSECT
	SMACHOPLT
	SELFSECT
	SMACHO
	SMACHOGOT
	SWINDOWS
	SELFGOT
	SNOPTRDATA
	SINITARR
	SDATA
	SBSS
	SNOPTRBSS
	STLSBSS
	SXREF
	SMACHOSYMSTR
	SMACHOSYMTAB
	SMACHOINDIRECTPLT
	SMACHOINDIRECTGOT
	SFILE
	SFILEPATH
	SCONST
	SDYNIMPORT
	SHOSTOBJ
	SSUB       = 1 << 8
	SMASK      = SSUB - 1
	SHIDDEN    = 1 << 9
	SCONTAINER = 1 << 10 // has a sub-symbol
)

const (
	NAME_NONE = 0 + iota
	NAME_EXTERN
	NAME_STATIC
	NAME_AUTO
	NAME_PARAM
	// A reference to name@GOT(SB) is a reference to the entry in the global offset
	// table for 'name'.
	NAME_GOTREF
)

const (
	TYPE_NONE   = 0
	TYPE_BRANCH = 5 + iota
	TYPE_TEXTSIZE
	TYPE_MEM
	TYPE_CONST
	TYPE_FCONST
	TYPE_SCONST
	TYPE_REG
	TYPE_ADDR
	TYPE_SHIFT
	TYPE_REGREG
	TYPE_REGREG2
	TYPE_INDIR
	TYPE_REGLIST
)
