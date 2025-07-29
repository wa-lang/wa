// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package obj

import (
	"errors"
	"fmt"
)

// An Addr is an argument to an instruction.
// The general forms and their encodings are:
//
//	sym±offset(symkind)(reg)(index*scale)
//		Memory reference at address &sym(symkind) + offset + reg + index*scale.
//		Any of sym(symkind), ±offset, (reg), (index*scale), and *scale can be omitted.
//		If (reg) and *scale are both omitted, the resulting expression (index) is parsed as (reg).
//		To force a parsing as index*scale, write (index*1).
//		Encoding:
//			type = TYPE_MEM
//			name = symkind (NAME_AUTO, ...) or 0 (NAME_NONE)
//			sym = sym
//			offset = ±offset
//			reg = reg (REG_*)
//			index = index (REG_*)
//			scale = scale (1, 2, 4, 8)
//
//	$<mem>
//		Effective address of memory reference <mem>, defined above.
//		Encoding: same as memory reference, but type = TYPE_ADDR.
//
//	$<±integer value>
//		This is a special case of $<mem>, in which only ±offset is present.
//		It has a separate type for easy recognition.
//		Encoding:
//			type = TYPE_CONST
//			offset = ±integer value
//
//	*<mem>
//		Indirect reference through memory reference <mem>, defined above.
//		Only used on x86 for CALL/JMP *sym(SB), which calls/jumps to a function
//		pointer stored in the data word sym(SB), not a function named sym(SB).
//		Encoding: same as above, but type = TYPE_INDIR.
//
//	$*$<mem>
//		No longer used.
//		On machines with actual SB registers, $*$<mem> forced the
//		instruction encoding to use a full 32-bit constant, never a
//		reference relative to SB.
//
//	$<floating point literal>
//		Floating point constant value.
//		Encoding:
//			type = TYPE_FCONST
//			val = floating point value
//
//	$<string literal, up to 8 chars>
//		String literal value (raw bytes used for DATA instruction).
//		Encoding:
//			type = TYPE_SCONST
//			val = string
//
//	<register name>
//		Any register: integer, floating point, control, segment, and so on.
//		If looking for specific register kind, must check type and reg value range.
//		Encoding:
//			type = TYPE_REG
//			reg = reg (REG_*)
//
//	x(PC)
//		Encoding:
//			type = TYPE_BRANCH
//			val = Prog* reference OR ELSE offset = target pc (branch takes priority)
//
//	$±x-±y
//		Final argument to TEXT, specifying local frame size x and argument size y.
//		In this form, x and y are integer literals only, not arbitrary expressions.
//		This avoids parsing ambiguities due to the use of - as a separator.
//		The ± are optional.
//		If the final argument to TEXT omits the -±y, the encoding should still
//		use TYPE_TEXTSIZE (not TYPE_CONST), with u.argsize = ArgsSizeUnknown.
//		Encoding:
//			type = TYPE_TEXTSIZE
//			offset = x
//			val = int32(y)
//
//	reg<<shift, reg>>shift, reg->shift, reg@>shift
//		Shifted register value, for ARM.
//		In this form, reg must be a register and shift can be a register or an integer constant.
//		Encoding:
//			type = TYPE_SHIFT
//			offset = (reg&15) | shifttype<<5 | count
//			shifttype = 0, 1, 2, 3 for <<, >>, ->, @>
//			count = (reg&15)<<8 | 1<<4 for a register shift count, (n&31)<<7 for an integer constant.
//
//	(reg, reg)
//		A destination register pair. When used as the last argument of an instruction,
//		this form makes clear that both registers are destinations.
//		Encoding:
//			type = TYPE_REGREG
//			reg = first register
//			offset = second register
//
//	[reg, reg, reg-reg]
//		Register list for ARM.
//		Encoding:
//			type = TYPE_REGLIST
//			offset = bit mask of registers in list; R0 is low bit.
//
//	reg, reg
//		Register pair for ARM.
//		TYPE_REGREG2
type Addr struct {
	Type   AddrType
	Reg    RBaseType
	Index  int16
	Scale  int16 // Sometimes holds a register.
	Name   AddrName
	Class  int8
	Etype  uint8
	Offset int64
	Width  int64
	Sym    *LSym
	Watype *LSym

	// argument value:
	//	for TYPE_SCONST, a string
	//	for TYPE_FCONST, a float64
	//	for TYPE_BRANCH, a *Prog (optional)
	//	for TYPE_TEXTSIZE, an int32 (optional)
	Val interface{}

	Node interface{} // for use by compiler
}

// 符号(LSym)名字类别定义
type AddrName int8

const (
	NAME_NONE   AddrName = 0 + iota
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
type AddrType int16

const (
	TYPE_NONE     AddrType = iota // 操作数不存在
	TYPE_BRANCH                   // 跳转地址, 例如 JMP target
	TYPE_TEXTSIZE                 // 用于 .text 段中设置函数大小的标记(用于链接器调试或优化目的)
	TYPE_MEM                      // 内存访问, 例如 [R1+4] 在 MOV 等指令中用于读取或写入内存
	TYPE_CONST                    // 整型常量值, 例如 MOV $1, R0 中的 $1
	TYPE_FCONST                   // 浮点常量, 例如 $3.14
	TYPE_SCONST                   // 字符串常量, 例 .string "hello"
	TYPE_REG                      // 单一寄存器, 如 AX, R1, X0
	TYPE_ADDR                     // 地址引用, 比如 MOV $foo(SB), AX 中的 foo(SB). 这是静态地址访问, 常用于全局变量
	TYPE_SHIFT                    // 表示移位操作, 常见于 ARM 或其他 RISC 架构. 如 LSL R1, R2, #3 (逻辑左移)
	TYPE_REGREG                   // 两个寄存器, 例如 CMP R1, R2
	TYPE_REGREG2                  // 用于三寄存器形式的指令, 或者链接器特殊用途
	TYPE_INDIR                    // 间接寻址, 如 *(R1) 或 [R1]. 可与 TYPE_MEM 组合使用?
	TYPE_REGLIST                  // 寄存器列表, 常用于 PUSH {R1-R4} 或 POP 操作, 在 ARM 架构中常见
)

func (a *Addr) checkaddr() error {
	// Check expected encoding, especially TYPE_CONST vs TYPE_ADDR.
	switch a.Type {
	case TYPE_NONE:
		return nil

	case TYPE_BRANCH:
		if a.Reg != 0 || a.Index != 0 || a.Scale != 0 || a.Name != 0 {
			break
		}
		return nil

	case TYPE_TEXTSIZE:
		if a.Reg != 0 || a.Index != 0 || a.Scale != 0 || a.Name != 0 {
			break
		}
		return nil

	case TYPE_MEM:
		return nil

	case TYPE_CONST:
		if a.Name != 0 || a.Sym != nil || a.Reg != 0 {
			return errors.New("argument is TYPE_CONST, should be TYPE_ADDR")
		}

		if a.Reg != 0 || a.Scale != 0 || a.Name != 0 || a.Sym != nil || a.Val != nil {
			break
		}
		return nil

	case TYPE_FCONST, TYPE_SCONST:
		if a.Reg != 0 || a.Index != 0 || a.Scale != 0 || a.Name != 0 || a.Offset != 0 || a.Sym != nil {
			break
		}
		return nil

	case TYPE_REG:
		if a.Scale != 0 || a.Name != 0 || a.Sym != nil {
			break
		}
		return nil

	case TYPE_ADDR:
		if a.Val != nil {
			break
		}
		if a.Reg == 0 && a.Index == 0 && a.Scale == 0 && a.Name == 0 && a.Sym == nil {
			return errors.New("argument is TYPE_ADDR, should be TYPE_CONST")
		}
		return nil

	case TYPE_SHIFT:
		if a.Index != 0 || a.Scale != 0 || a.Name != 0 || a.Sym != nil || a.Val != nil {
			break
		}
		return nil

	case TYPE_REGREG:
		if a.Index != 0 || a.Scale != 0 || a.Name != 0 || a.Sym != nil || a.Val != nil {
			break
		}
		return nil

	case TYPE_REGREG2:
		return nil

	case TYPE_REGLIST:
		return nil

	// Expect sym and name to be set, nothing else.
	// Technically more is allowed, but this is only used for *name(SB).
	case TYPE_INDIR:
		if a.Reg != 0 || a.Index != 0 || a.Scale != 0 || a.Name == 0 || a.Offset != 0 || a.Sym == nil || a.Val != nil {
			break
		}
		return nil
	}

	return fmt.Errorf("invalid encoding for argument %v", a)
}

// 将表示地址的结构 *Addr 转换为字符串
func (a *Addr) String() string {
	var str string

	switch a.Name {
	default:
		str = fmt.Sprintf("name=%d", a.Name)

	case NAME_NONE:
		switch {
		case a.Reg == REG_NONE:
			str = fmt.Sprint(a.Offset)
		case a.Offset == 0:
			str = fmt.Sprintf("(%v)", a.Reg)
		case a.Offset != 0:
			str = fmt.Sprintf("%d(%v)", a.Offset, a.Reg)
		}

	case NAME_EXTERN:
		str = fmt.Sprintf("%s%s(SB)", a.Sym.Name, a.offString())

	case NAME_GOTREF:
		str = fmt.Sprintf("%s%s@GOT(SB)", a.Sym.Name, a.offString())

	case NAME_STATIC:
		str = fmt.Sprintf("%s<>%s(SB)", a.Sym.Name, a.offString())

	case NAME_AUTO:
		if a.Sym != nil {
			str = fmt.Sprintf("%s%s(SP)", a.Sym.Name, a.offString())
		} else {
			str = fmt.Sprintf("%s(SP)", a.offString())
		}

	case NAME_PARAM:
		if a.Sym != nil {
			str = fmt.Sprintf("%s%s(FP)", a.Sym.Name, a.offString())
		} else {
			str = fmt.Sprintf("%s(FP)", a.offString())
		}
	}
	return str
}

// 为偏移量生成对应的汇编格式字符串
func (a *Addr) offString() string {
	if a.Offset == 0 {
		return ""
	}
	return fmt.Sprintf("%+d", a.Offset)
}
