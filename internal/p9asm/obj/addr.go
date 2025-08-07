// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package obj

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"wa-lang.org/wa/internal/p9asm/objabi"
)

// Addr 对应汇编指令的一个参数. 有以下几种寻址形式:
//
// - `sym±offset(symkind)(reg)(index*scale)`
//
// 对应的内存地址为: `&sym(symkind) + offset + reg + index*scale`.
// 任何一个 `sym(symkind)`, `±offset`, `(reg)`, `(index*scale)` 和 `*scale` 都是可省略的.
// 如果 `(reg)` 和 `*scale` 都被省略, 那么解析得到的 `(index)` 将被视为 `(reg)`.
// 这就要求 `(index)` 必须写成 `(index*1)` 的形式, 避免被误认为是 `(reg)`.
// `sym` 前缀类似这个地址的名字, 用于注释或调试信息, 并不参与真正的地址计算.
//
// 对应: `Addr{Type:TYPE_MEM; sym:sym; Offset:ofset; Reg:REG_*; index:REG_*; Scale: (1, 2, 4, 8) }`
//
// - `$<mem>`
//
// `<mem>`对应内存的地址. 对应: `Addr{Type:TYPE_ADDR}`, 其他和`TYPE_MEM`相似.
//
// - `$<±integer_value>`
//
// 这是`$<mem>`形式的特殊情况, 只是其中只有 `±offset` 部分.
// 对应 `Addr{Type:TYPE_CONST; Offset:±integer_value}`
//
// - `*<mem>`
//
// 间接内存寻址. 仅在 X86 平台用于 `CALL/JMP *sym(SB)` 指令, `sym(SB)` 对应的内存是指针数据.
// 表示形式同上, 但是对应 `TYPE_INDIR` 类型.
//
// - `$<floating_point_literal>`
//
// 浮点数常量. 表示形式 `Addr{Type:TYPE_FCONST: Val:floating_point_literal}`
//
// - `$<string_literal>`, 最多8个字符
//
// 字符串面值, 用于 `DATA` 指令的原始数据.
// 表示形式 `Addr{Type:TYPE_SCONST; Val:string}`
//
// - `<register name>`
//
// 可以是任何的寄存器名字, 比如 整数/浮点/控制/段等寄存器类型.
// 如果要查找特殊种类的寄存器, 必须检查类型和寄存器的范围(比如有些是浮点数寄存器, 且有数量限制).
// 表示形式 `Addr{Type:TYPE_REG; Reg:REG_*}`
//
// - `x(PC)`
//
// 对应跳转指令的地址.
// 表示形式 `Addr{Type:TYPE_BRANCH; Val: Prog* reference OR ELSE offset = target pc (branch takes priority) }`
//
// - `$±x-±y`
//
// `TEXT` 伪指令最后的参数(TODO:第二个已经废弃), 用于表示函数帧的大小为`x`, 参数的大小为`y`.
// 在这里形式中 `x` 和 `y` 都是整数, 主要是为了避免中间的 `-` 对解析产生的干扰.
// 而 `±` 部分是可省的.
//
// 如果 `TEXT` 指令省略了第2个参数 `-±y` 部分, 类型依然是 `TYPE_TEXTSIZE`,
// 但是参数大小用 `ArgsSizeUnknown` 表示.
//
// 表示形式 `Addr{Type:TYPE_TEXTSIZE; Offset:x; Val:int32(y)}`
//
// - `reg<<shift` 逻辑左移, `reg>>shift` 逻辑右移, `reg->shift` 数学右移, `reg@>shift` 循环右移
//
// ARM平台的移位寄存器. 这种形式中, reg 部分必须是寄存器, shift 部分可以是寄存器或整数常量.
//
// 表示形式:
//
//	type = TYPE_SHIFT
//	offset = (reg&15) | shifttype<<5 | count
//	shifttype = 0, 1, 2, 3 for <<, >>, ->, @>
//	count = (reg&15)<<8 | 1<<4 for a register shift count, (n&31)<<7 for an integer constant.
//
// - `(reg1, reg2)`
//
// 目标寄存器对. 当用作指令的最后一个参数时, 这种形式明确表示两个寄存器都是目标寄存器.
// 表示形式 `Addr{Type:TYPE_REGREG; Reg:reg1; Offset:reg2}`
//
// - `[reg, reg, reg-reg]`
//
// ARM 平台的寄存器列表.
// 表示形式 `Addr{Type: TYPE_REGLIST; Offset: bit mask of registers in list; R0 is low bit.}`
//
// - `reg, reg`
//
// ARM 的寄存器对. 对应 `Addr{Type:TYPE_REGREG2}`
type Addr struct {
	Type AddrType
	Name AddrName
	Sym  *LSym

	Offset int64
	Reg    objabi.RBaseType
	Index  int16
	Scale  int16 // Sometimes holds a register.

	// argument value:
	//	for TYPE_SCONST, a string
	//	for TYPE_FCONST, a float64
	//	for TYPE_BRANCH, a *Prog (optional)
	//	for TYPE_TEXTSIZE, an int32 (optional)
	Val interface{}
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

func (dt AddrType) String() string {
	switch dt {
	default:
		return fmt.Sprintf("AddrType(%d)", int(dt))
	case TYPE_NONE:
		return "TYPE_NONE"
	case TYPE_BRANCH:
		return "TYPE_BRANCH"
	case TYPE_TEXTSIZE:
		return "TYPE_TEXTSIZE"
	case TYPE_MEM:
		return "TYPE_MEM"
	case TYPE_CONST:
		return "TYPE_CONST"
	case TYPE_FCONST:
		return "TYPE_FCONST"
	case TYPE_SCONST:
		return "TYPE_SCONST"
	case TYPE_REG:
		return "TYPE_REG"
	case TYPE_ADDR:
		return "TYPE_ADDR"
	case TYPE_SHIFT:
		return "TYPE_SHIFT"
	case TYPE_REGREG:
		return "TYPE_REGREG"
	case TYPE_REGREG2:
		return "TYPE_REGREG2"
	case TYPE_INDIR:
		return "TYPE_INDIR"
	case TYPE_REGLIST:
		return "TYPE_REGLIST"
	}
}

// 根据地址的类型生成不同格式的字符串
// 通常用于生成完整的操作数字符串, 可能包含寄存器/符号/偏移量等复合信息
func (a *Addr) Dconv(p *Prog) string {
	var str string

	switch a.Type {
	default:
		str = fmt.Sprintf("type=%d", a.Type)

	case TYPE_NONE:
		str = ""
		if a.Name != NAME_NONE || a.Reg != 0 || a.Sym != nil {
			str = fmt.Sprintf("%v(%v)(NONE)", a.Mconv(), a.Reg)
		}

	case TYPE_REG:
		// 寄存器
		// TODO(chai2010): This special case is for x86 instructions like
		//	PINSRQ	CX,$1,X6
		// where the $1 is included in the p->to Addr.
		// Move into a new field.
		if a.Offset != 0 {
			str = fmt.Sprintf("$%d,%v", a.Offset, a.Reg)
			break
		}

		str = a.Reg.String()
		if a.Name != NAME_NONE || a.Sym != nil {
			str = fmt.Sprintf("%v(%v)(REG)", a.Mconv(), a.Reg)
		}

	case TYPE_BRANCH:
		// 跳转的地址
		if a.Sym != nil {
			str = fmt.Sprintf("%s(SB)", a.Sym.Name)
		} else if p != nil && p.Pcond != nil {
			str = fmt.Sprint(p.Pcond.Pc)
		} else if a.Val != nil {
			str = fmt.Sprint(a.Val.(*Prog).Pc)
		} else {
			str = fmt.Sprintf("%d(PC)", a.Offset)
		}

	case TYPE_INDIR:
		// 间接寻址
		str = fmt.Sprintf("*%v", a.Mconv())

	case TYPE_MEM:
		// 内存地址
		str = a.Mconv()
		if a.Index != int16(objabi.REG_NONE) {
			str += fmt.Sprintf("(%v*%d)", objabi.RBaseType(a.Index), int(a.Scale))
		}

	case TYPE_CONST:
		// 常量
		if a.Reg != 0 {
			str = fmt.Sprintf("$%v(%v)", a.Mconv(), a.Reg)
		} else {
			str = fmt.Sprintf("$%v", a.Mconv())
		}

	case TYPE_TEXTSIZE:
		// TEXT pkg·foo(SB),$framesize-argsize
		// 参数的大小必须是固定的, 不再支持 printf 形式的可变参数
		str = fmt.Sprintf("$%d-%d", a.Offset, a.Val.(int32))

	case TYPE_FCONST:
		// 浮点数常量
		str = fmt.Sprintf("%.17g", a.Val.(float64))
		// Make sure 1 prints as 1.0
		if !strings.ContainsAny(str, ".e") {
			str += ".0"
		}
		str = fmt.Sprintf("$(%s)", str)

	case TYPE_SCONST:
		// 字符串常量
		str = fmt.Sprintf("$%q", a.Val.(string))

	case TYPE_ADDR:
		// 地址是必须满足特定的格式特征, 以便于和 TYPE_SCONST 区分
		str = fmt.Sprintf("$%v", a.Mconv())

	case TYPE_SHIFT:
		v := int(a.Offset)
		// `<<`, `>>`, `->`, `@>`
		op := string("<<>>->@>"[((v>>5)&3)<<1:])
		if v&(1<<4) != 0 {
			str = fmt.Sprintf("R%d%c%cR%d", v&15, op[0], op[1], (v>>8)&15)
		} else {
			str = fmt.Sprintf("R%d%c%c%d", v&15, op[0], op[1], (v>>7)&31)
		}
		if a.Reg != 0 {
			str += fmt.Sprintf("(%v)", a.Reg)
		}

	case TYPE_REGREG:
		str = fmt.Sprintf("(%v, %v)", a.Reg, objabi.RBaseType(a.Offset))

	case TYPE_REGREG2:
		str = fmt.Sprintf("%v, %v", a.Reg, objabi.RBaseType(a.Offset))

	case TYPE_REGLIST:
		// 通常出现在ARM, 最多有16个寄存器列表
		var sb strings.Builder
		for i := 0; i < 16; i++ {
			if a.Offset&(1<<uint(i)) != 0 {
				if sb.Len() == 0 {
					sb.WriteRune('[')
				} else {
					sb.WriteRune(',')
				}
				// 寄存器列表是 ARM 的用法, R10 对应 g 寄存器
				if i == 10 {
					sb.WriteRune('g')
				} else {
					sb.WriteRune('R')
					sb.WriteString(strconv.Itoa(i))
				}
			}
		}
		if sb.Len() == 0 {
			sb.WriteRune('[')
		}
		sb.WriteRune(']')
		str = sb.String()
	}

	return str
}

// 根据地址的名称类型生成内存引用的字符串表示
// 专注于符号/偏移量和基址寄存器的组合
func (a *Addr) Mconv() string {
	var str string

	switch a.Name {
	default:
		str = fmt.Sprintf("name=%d", a.Name)

	case NAME_NONE:
		switch {
		case a.Reg == objabi.REG_NONE:
			str = fmt.Sprint(a.Offset)
		case a.Offset == 0:
			str = fmt.Sprintf("(%v)", a.Reg)
		case a.Offset != 0:
			str = fmt.Sprintf("%d(%v)", a.Offset, a.Reg)
		}

	case NAME_EXTERN:
		if a.Offset != 0 {
			str = fmt.Sprintf("%s%+d(SB)", a.Sym.Name, a.Offset)
		} else {
			str = fmt.Sprintf("%s(SB)", a.Sym.Name)
		}

	case NAME_GOTREF:
		if a.Offset != 0 {
			str = fmt.Sprintf("%s%+d@GOT(SB)", a.Sym.Name, a.Offset)
		} else {
			str = fmt.Sprintf("%s@GOT(SB)", a.Sym.Name)
		}

	case NAME_STATIC:
		if a.Offset != 0 {
			str = fmt.Sprintf("%s<>%+d(SB)", a.Sym.Name, a.Offset)
		} else {
			str = fmt.Sprintf("%s<>(SB)", a.Sym.Name)
		}

	case NAME_AUTO:
		if a.Sym != nil {
			if a.Offset != 0 {
				str = fmt.Sprintf("%s%+d(SP)", a.Sym.Name, a.Offset)
			} else {
				str = fmt.Sprintf("%s(SP)", a.Sym.Name)
			}
		} else {
			if a.Offset != 0 {
				str = fmt.Sprintf("%+d(SP)", a.Offset)
			} else {
				str = "(SP)"
			}
		}

	case NAME_PARAM:
		if a.Sym != nil {
			if a.Offset != 0 {
				str = fmt.Sprintf("%s%+d(FP)", a.Sym.Name, a.Offset)
			} else {
				str = fmt.Sprintf("%s(FP)", a.Sym.Name)
			}
		} else {
			if a.Offset != 0 {
				str = fmt.Sprintf("%+d(FP)", a.Offset)
			} else {
				str = "(FP)"
			}
		}
	}
	return str
}

// 检查地址形式是否合法
// 特别是 TYPE_CONST 和 TYPE_ADDR 类型
func (a *Addr) checkAddr() error {
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
