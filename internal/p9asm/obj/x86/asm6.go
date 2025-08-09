// Inferno utils/6l/span.c
// http://code.google.com/p/inferno-os/source/browse/utils/6l/span.c
//
//	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
//	Portions Copyright © 1995-1997 C H Forsyth (forsyth@terzarima.net)
//	Portions Copyright © 1997-1999 Vita Nuova Limited
//	Portions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com)
//	Portions Copyright © 2004,2006 Bruce Ellis
//	Portions Copyright © 2005-2007 C H Forsyth (forsyth@terzarima.net)
//	Revisions Copyright © 2000-2007 Lucent Technologies Inc. and others
//	Portions Copyright © 2009 The Go Authors.  All rights reserved.
//	Portions Copyright © 2025 武汉凹语言科技有限公司.  All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package x86

import (
	"fmt"
	"log"
	"strings"

	"wa-lang.org/wa/internal/p9asm/obj"
	"wa-lang.org/wa/internal/p9asm/objabi"
)

// 使用帧指针
const Framepointer_enabled = true

// 指令布局

const (
	// 数据最大对齐边界
	MaxAlign = 32 // max data alignment

	// Loop alignment constants:
	// want to align loop entry to LoopAlign-byte boundary,
	// and willing to insert at most MaxLoopPad bytes of NOP to do so.
	// We define a loop entry as the target of a backward jump.
	//
	// gcc uses MaxLoopPad = 10 for its 'generic x86-64' config,
	// and it aligns all jump targets, not just backward jump targets.
	//
	// As of 6/1/2012, the effect of setting MaxLoopPad = 10 here
	// is very slight but negative, so the alignment is disabled by
	// setting MaxLoopPad = 0. The code is here for reference and
	// for future experiments.
	//

	// 编译器可以让循环入口对齐到一个固定的边界, 以提升 CPU 指令缓存命中率和解码效率
	LoopAlign = 16

	// 对齐时允许最多插入多少字节的 NOP 指令
	// 这里设为 0 表示 禁用 loop alignment
	// GCC 默认会插入最多 10 字节 NOP 来对齐循环入口(甚至对齐所有跳转目标, 不仅是循环入口)
	// Go 编译器在 2012 年测试过, 没有显著效果, 所以默认禁用
	MaxLoopPad = 0

	// 函数入口对齐到 16 字节。
	// 现代 x86 CPU 对 16 字节对齐的函数入口支持得更好, 性能更好
	FuncAlign = 16
)

// 指令机器码规范信息
type Optab struct {
	as     objabi.As // 指令种类(opcode 编号), 比如 AADD、AMOV 等
	ytab   []ytab    // 类型映射表: 描述操作数类型组合对应的编码
	prefix uint8     // 前缀字节, X86 有 0xF2, 0xF3, 0x66 等
	op     [23]uint8 // 输出的机器码模板
}

// 指令的操作数类别和对应的模板位置
type ytab struct {
	from    uint8 // 第一个操作数的类型
	from3   uint8 // 可选的第三个操作数类型(部分指令有三个操作数, 比如 LEA)
	to      uint8 // 第二个操作数的类型
	zcase   uint8 // 对应机器码模板的序号
	zoffset uint8 // 偏移, 用于在模板里取哪一段
}

// 描述 MOV 指令编码方式 的表项, 用来把不同操作数类型映射到不同机器码模板
type Movtab struct {
	as   objabi.As // 指令代号, 比如 AMOVB, AMOVW, AMOVL, 代表不同大小的 MOV 指令
	ft   uint8     // 第一个操作数(from)的类型(内部枚举, 比如寄存器/内存/立即数等)
	f3t  uint8     // 可选第三个操作数的类型. 多数情况下是 0, 只对少数指令有用, 比如三操作数指令
	tt   uint8     // 第二个操作数(to)的类型
	code uint8     // 编码方式编号，告诉后端怎么生成机器码
	op   [4]uint8  // 机器码模板, 最多 4 字节(因为 MOV 类指令通常很短)
}

// 操作数类型枚举表
// 用来告诉汇编器在指令匹配
// 比如 Optab、Movtab 时, 操作数属于哪一类
//
// 每个操作数(寄存器/立即数/内存/段寄存器等)在生成机器码时都要归类:
// - 是寄存器还是内存？
// - 是 8 位还是 32 位？
// - 是普通寄存器还是浮点寄存器？
// - 是立即数还是绝对地址？
const (
	Yxxx  = iota // 空类型
	Ynone        // 无操作数

	// 立即数
	Yi0    // $0
	Yi1    // $1
	Yi8    // $x, x fits in int8
	Yu8    // $x, x fits in uint8
	Yu7    // $x, x in 0..127 (fits in both int8 and uint8)
	Ys32   // 符号(符号常量)
	Yi32   // int32
	Yi64   // int64
	Yiauto // ?

	// 寄存器分类
	Yal   // AL
	Ycl   // CL
	Yax   // AX
	Ycx   // CX
	Yrb   // 通用 8 位寄存器(比如 AL, BL, CL, DL, SIL, DIL, BPL, SPL, R8B-R15B)
	Yrl   // 通用 32/64 位寄存器(比如 AX, BX, CX, DX, SI, DI, BP, SP, R8~R15)
	Yrl32 // Yrl on 32-bit system
	Yrf   // 浮点寄存器 ST(0) - ST(7)
	Yf0   // ST(0)(浮点堆栈栈顶)
	Yrx   // XMM 寄存器(XMM0~XMM15)

	// 内存寻址/寄存器
	Ymb       // 内存地址, 目标是 8 位数据
	Yml       // 内存地址, 目标是 32/64 位数据
	Ym        // 泛指内存地址
	Ybr       // 段寄存器(CS, SS, DS, ES, FS, GS)
	Ycs       // 段寄存器 CS
	Yss       // 段寄存器 SS
	Yds       // 段寄存器 DS
	Yes       // 段寄存器 ES
	Yfs       // 段寄存器 FS
	Ygs       // 段寄存器 GS
	Ygdtr     // GDTR
	Yidtr     // IDTR
	Yldtr     // LDTR
	Ymsw      // 机器状态字
	Ytask     // Task Register
	Ycr0      // CR0 控制寄存器
	Ycr1      //
	Ycr2      //
	Ycr3      //
	Ycr4      //
	Ycr5      //
	Ycr6      //
	Ycr7      //
	Ycr8      //
	Ydr0      // DR0 调试寄存器
	Ydr1      //
	Ydr2      //
	Ydr3      //
	Ydr4      //
	Ydr5      //
	Ydr6      //
	Ydr7      //
	Ytr0      // 测试寄存器(历史遗留)
	Ytr1      //
	Ytr2      //
	Ytr3      //
	Ytr4      //
	Ytr5      //
	Ytr6      //
	Ytr7      //
	Ymr       // MMX 寄存器
	Ymm       // YMM 寄存器(AVX)
	Yxr       // XMM 寄存器
	Yxm       // XMM 内存操作数（movaps/movdqa 等指令里)
	Ytls      // TLS 基址寄存器(FS/GS)
	Ytextsize // .text 指令的大小参数(比如 TEXT sym(SB), 4, $40-8)
	Yindir    // 间接寻址
	Ymax      // 最大值, 通常用于边界检查或数组大小
)

// 指令模板/编码模式(zcase)
const (
	Zxxx        = iota // 无效
	Zlit               // 输出固定的字节序列(字面值)
	Zlitm_r            // 字面值 + 寄存器
	Z_rp               // 寄存器到寄存器/内存（Reg / memory）
	Zbr                // 分支指令：短跳、长跳
	Zcall              // 绝对地址 call
	Zcallcon           // call 符号常量
	Zcallduff          // duff’s device 调用
	Zcallind           // 间接调用（通过内存地址）
	Zcallindreg        // 通过寄存器调用
	Zib_               // 立即数（8位）
	Zib_rp             // opcode + 立即数8 + reg/mem
	Zibo_m             // opcode + imm8 + 内存操作数
	Zibo_m_xm          // opcode + imm8 + 内存 + xmm 寄存器
	Zil_               // 立即数（32位）
	Zil_rp             // opcode + 立即数32 + reg/mem
	Ziq_rp             // 立即数（64位）到寄存器/内存
	Zilo_m             // opcode + 32位立即数到内存
	Ziqo_m             // opcode + 64位立即数到内存
	Zjmp               // 无条件跳转
	Zjmpcon            // 条件跳转
	Zloop              // loop 指令
	Zo_iw              // opcode 后跟 16 位立即数
	Zm_o               // 内存偏移到 opcode（比如 MOV 内存偏移量到 AL）
	Zm_r               // 内存到寄存器
	Zm2_r              // 两字节内存到寄存器（16位 load）
	Zm_r_xm            // 内存到 XMM 寄存器
	Zm_r_i_xm          // 内存+立即数到 XMM
	Zm_r_3d            // 内存到 3DNow! 寄存器
	Zm_r_xm_nr         // 内存到寄存器，不需要 REX 前缀
	Zr_m_xm_nr         // 寄存器到内存，不需要 REX 前缀
	Zibm_r             // mmx1,mmx2/mem64,imm8
	Zmb_r              // 内存(字节)到寄存器
	Zaut_r             // 自动变量到寄存器（栈变量寻址）
	Zo_m               // opcode 后跟内存地址
	Zo_m64             // opcode 后跟 64 位内存地址
	Zpseudo            // 伪指令, 不需要生成机器码
	Zr_m               // 寄存器到内存
	Zr_m_xm            // opcode + reg (xmm) + 内存
	Zrp_               // opcode + reg/mem
	Z_ib               // opcode 后跟 8 位立即数
	Z_il               // opcode 后跟 32 位立即数
	Zm_ibo             // 内存地址 + 8位立即数
	Zm_ilo             // 内存地址 + 32位立即数
	Zib_rr             // opcode + imm8 + reg,reg
	Zil_rr             // opcode + imm32 + reg,reg
	Zclr               // 清零指令（特殊编码）
	Zbyte              // 输出原始字节（db 0x90 等）
	Zmax               // 枚举最大值(通常用来做数组大小)
)

// 指令前缀/REX前缀/内部flag
// 例如: 指令需要 f2+0f 前缀, 对应 Pf2|Pm
const (
	Px  = 0    // 默认无前缀
	Px1 = 1    // symbolic; exact value doesn't matter, 符号型占位（只在内部标记时用，具体值不重要）
	P32 = 0x32 // 32-bit only
	Pe  = 0x66 // operand escape
	Pm  = 0x0f // 2byte opcode escape
	Pq  = 0xff // both escapes: 66 0f
	Pb  = 0xfe // byte operands
	Pf2 = 0xf2 // xmm escape 1: f2 0f
	Pf3 = 0xf3 // xmm escape 2: f3 0f
	Pq3 = 0x67 // xmm escape 3: 66 48 0f
	Pw  = 0x48 // Rex.w
	Pw8 = 0x90 // symbolic; exact value doesn't matter
	Py  = 0x80 // defaults to 64-bit mode
	Py1 = 0x81 // symbolic; exact value doesn't matter
	Py3 = 0x83 // symbolic; exact value doesn't matter

	// REX 前缀标志位
	Rxf = 1 << 9 // internal flag for Rxr on from
	Rxt = 1 << 8 // internal flag for Rxr on to
	Rxw = 1 << 3 // =1, 64-bit operand size
	Rxr = 1 << 2 // extend modrm reg
	Rxx = 1 << 1 // extend sib index
	Rxb = 1 << 0 // extend modrm r/m, sib base, or opcode reg

	Maxand = 10 // in -a output width of the byte codes
)

// 覆盖表, 大小是所有可能的操作数类型组合数
// Ymax 是枚举操作数类型(比如 Yi8、Yax、Yrl 等)的最大值
//
// 比如要判断以下指令:
// 第一个操作数是 Yi8, 第二个是 Yax, 这种组合合法吗
//
// ycover 用于快速查找: 给定 (from, to) 操作数类型，返回一个标记值
//
// 例如: MOV $5, AX 会查 ycover[Yi8*Ymax + Yax]
var ycover [Ymax * Ymax]uint8

// 寄存器属性
// 给每个寄存器编号对应一个内部属性或编号
// MAXREG 是寄存器枚举表里寄存器的总数
// 通常在编译阶段/汇编阶段, 根据寄存器常量快速查到其内部编号或类型
var reg [MAXREG]int

// 保存每个寄存器需要加的 REX 前缀位(R, X, B)
// 在 x86-64 中, R8~R15 寄存器需要 REX 前缀
var regrex [MAXREG + 1]int

// 无操作数指令的模式
// 如 NOP、RET, 都是没有操作数的指令
var ynone = []ytab{
	// 操作数 from/to 都是 Ynone
	// 指令编码模式是 Zlit(输出固定字节)
	// zoffset=1, 表示机器码模板的偏移是 1
	{Ynone, Ynone, Ynone, Zlit, 1},
}

// SAHF指令(Store AH into Flags) 对应的2种模式
var ysahf = []ytab{
	{Ynone, Ynone, Ynone, Zlit, 2},
	{Ynone, Ynone, Ynone, Zlit, 1},
}

// TEXT 伪指令
//
// 比如以下2个指令:
// TEXT ·foo(SB), $0-8
// TEXT ·foo(SB), $16-8
var ytext = []ytab{
	// 第一个操作数是 Ymb(内存基址/符号, 比如 ·foo(SB))
	// 第二个操作数可以是: Ynone (省略 offset); Yi32 (提供 offset, 立即数 32)
	// 第三个是 Ytextsize (特殊内部类型: text 的大小参数)
	// zcase=Zpseudo: 表示这是伪指令
	// zoffset=0/1: 区分是否带立即数
	{Ymb, Ynone, Ytextsize, Zpseudo, 0},
	{Ymb, Yi32, Ytextsize, Zpseudo, 1},
}

// NOP 伪指令
//
// 为什么要这样写写?
// - 汇编里 NOP 有很多灵活用法(比如用于 patch/对齐/替代)
// - 汇编器允许 NOP 带上不同来源/目标: 比如 NOP 8(SP) 和 NOP X0
// - 这些表项就告诉汇编器: 这些组合是合法的 NOP 操作, 虽然最终不输出实际机器码
var ynop = []ytab{
	// zoffset 多数是 0 (表示匹配到的编码模板偏移为 0)
	{Ynone, Ynone, Ynone, Zpseudo, 0},  // NOP 无参数形式
	{Ynone, Ynone, Yiauto, Zpseudo, 0}, // 目标是自动变量偏移地址
	{Ynone, Ynone, Yml, Zpseudo, 0},    // 目标是内存地址
	{Ynone, Ynone, Yrf, Zpseudo, 0},    // 目标是浮点寄存器 ST(x)
	{Ynone, Ynone, Yxr, Zpseudo, 0},    // 目标是 XMM 寄存器
	{Yiauto, Ynone, Ynone, Zpseudo, 0}, // 来源是自动变量偏移地址
	{Yml, Ynone, Ynone, Zpseudo, 0},    // 来源是内存地址
	{Yrf, Ynone, Ynone, Zpseudo, 0},    // 来源是浮点寄存器

	// 表示对于 NOP XMM 寄存器 这种情况, 用到模板列表的第 2 个模板
	{Yxr, Ynone, Ynone, Zpseudo, 1}, // 来源是 XMM 寄存器
}

// FUNCDATA 伪指令
//
// 例子: FUNCDATA $index, sym+offset(SB)
var yfuncdata = []ytab{
	// from 是 Yi32: 第一个参数必须是立即数(比如 $1)
	// from3 是 Ynone: 中间没有第 3 个操作数
	// to 是 Ym: 第二个参数必须是内存符号
	{Yi32, Ynone, Ym, Zpseudo, 0},
}

// PCDATA 伪指令
// 例子: PCDATA $index, $value
var ypcdata = []ytab{
	// from 是 Yi32: 第一个参数是立即数
	// to 是 Yi32: 第二个参数也是立即数
	{Yi32, Ynone, Yi32, Zpseudo, 0},
}

// 8位版本的 XOR, 汇编指令 XORB
var yxorb = []ytab{
	{Yi32, Ynone, Yal, Zib_, 1},   // XOR $imm8, AL
	{Yi32, Ynone, Ymb, Zibo_m, 2}, // XOR $imm8, mem8
	{Yrb, Ynone, Ymb, Zr_m, 1},    // XOR reg8, mem8
	{Ymb, Ynone, Yrb, Zm_r, 1},    // XOR mem8, reg8
}

// 32位版本的 XOR, 汇编指令 XORL。
var yxorl = []ytab{
	{Yi8, Ynone, Yml, Zibo_m, 2},  // XOR $imm8, mem32
	{Yi32, Ynone, Yax, Zil_, 1},   // XOR $imm32, EAX
	{Yi32, Ynone, Yml, Zilo_m, 2}, // XOR $imm32, mem32
	{Yrl, Ynone, Yml, Zr_m, 1},    // XOR reg32, mem32
	{Yml, Ynone, Yrl, Zm_r, 1},    // XOR mem32, reg32
}

// 32位版本的 ADD, 汇编指令 ADDL
var yaddl = []ytab{
	{Yi8, Ynone, Yml, Zibo_m, 2},  // ADD $imm8, mem32
	{Yi32, Ynone, Yax, Zil_, 1},   // ADD $imm32, EAX
	{Yi32, Ynone, Yml, Zilo_m, 2}, // ADD $imm32, mem32
	{Yrl, Ynone, Yml, Zr_m, 1},    // ADD reg32, mem32
	{Yml, Ynone, Yrl, Zm_r, 1},    // ADD mem32, reg32
}

// 8 位 inc
var yincb = []ytab{
	// 无 from
	// 目标是 Ymb: 内存地址（byte）
	// zcase=Zo_m: 表示输出 opcode + modrm，只对内存做 inc。
	{Ynone, Ynone, Ymb, Zo_m, 2},
}

// 16 位 inc
var yincw = []ytab{
	// Yml: 目标是内存地址(word, long, qword)
	{Ynone, Ynone, Yml, Zo_m, 2},
}

// 32 位 inc
var yincl = []ytab{
	// 对目标是寄存器（Yrl, 32 位寄存器）的情况
	// zcase=Z_rp: 这是 x86 的特殊情况
	// - 对寄存器有短编码(opcode + reg number)
	// - 比如 INCL AX -> opcode 是 0x40+reg
	{Ynone, Ynone, Yrl, Z_rp, 1},

	// 对目标是内存(Yml)
	// zcase=Zo_m: 普通模式(opcode + modrm)
	{Ynone, Ynone, Yml, Zo_m, 2},
}

// 64 位 inc
var yincq = []ytab{
	// 只支持对内存(Yml)
	{Ynone, Ynone, Yml, Zo_m, 2},
}

// 8 位比较
var ycmpb = []ytab{
	{Yal, Ynone, Yi32, Z_ib, 1},
	{Ymb, Ynone, Yi32, Zm_ibo, 2},
	{Ymb, Ynone, Yrb, Zm_r, 1},
	{Yrb, Ynone, Ymb, Zr_m, 1},
}

// 32 位比较指令 CMP
var ycmpl = []ytab{
	{Yml, Ynone, Yi8, Zm_ibo, 2},  // CMP mem32, imm8
	{Yax, Ynone, Yi32, Z_il, 1},   // CMP EAX, imm32
	{Yml, Ynone, Yi32, Zm_ilo, 2}, // CMP mem32, imm32
	{Yml, Ynone, Yrl, Zm_r, 1},    // CMP mem32, reg32
	{Yrl, Ynone, Yml, Zr_m, 1},    // CMP reg32, mem32
}

// SHRB/SHLB/SALB 等 8 位移位指令
var yshb = []ytab{
	{Yi1, Ynone, Ymb, Zo_m, 2},    // shift by 1: 特殊短 opcode
	{Yi32, Ynone, Ymb, Zibo_m, 2}, // shift by imm8（可选立即数）
	{Ycx, Ynone, Ymb, Zo_m, 2},    // shift by CL
}

// SHRL/SHLL/SARL 等 32 位移位指令
var yshl = []ytab{
	{Yi1, Ynone, Yml, Zo_m, 2},    // shift by 1: 专用短指令
	{Yi32, Ynone, Yml, Zibo_m, 2}, // shift by imm8
	{Ycl, Ynone, Yml, Zo_m, 2},    // shift by CL
	{Ycx, Ynone, Yml, Zo_m, 2},    // shift by CX
}

// 8 位 TEST
var ytestb = []ytab{
	{Yi32, Ynone, Yal, Zib_, 1},   // TEST imm8, AL
	{Yi32, Ynone, Ymb, Zibo_m, 2}, // TEST imm8, mem8
	{Yrb, Ynone, Ymb, Zr_m, 1},    // TEST reg8, mem8
	{Ymb, Ynone, Yrb, Zm_r, 1},    // TEST mem8, reg8
}

// 32 位 TEST
var ytestl = []ytab{
	{Yi32, Ynone, Yax, Zil_, 1},   // TEST imm32, EAX
	{Yi32, Ynone, Yml, Zilo_m, 2}, // TEST imm32, mem32
	{Yrl, Ynone, Yml, Zr_m, 1},    // TEST reg32, mem32
	{Yml, Ynone, Yrl, Zm_r, 1},    // TEST mem32, reg32
}

// 8 位 MOV
var ymovb = []ytab{
	{Yrb, Ynone, Ymb, Zr_m, 1},    // MOV reg8, mem8
	{Ymb, Ynone, Yrb, Zm_r, 1},    // MOV mem8, reg8
	{Yi32, Ynone, Yrb, Zib_rp, 1}, // MOV imm8, reg8
	{Yi32, Ynone, Ymb, Zibo_m, 2}, // MOV imm8, mem8

}

// 特殊用法
// 通常对应: 比如 BSF, BSR 需要只对目标是内存的形式, 或 MOVS, SCAS 等字符串操作指令只带一个内存地址
var ymbs = []ytab{
	{Ymb, Ynone, Ynone, Zm_o, 2}, // 对内存地址做“单目操作”
}

// BTL (Bit Test Long, 测试位指令)
var ybtl = []ytab{
	{Yi8, Ynone, Yml, Zibo_m, 2}, // BT imm8, mem32: 用 imm8 做 bit index
	{Yrl, Ynone, Yml, Zr_m, 1},   // BT reg32, mem32: 用寄存器做 bit index
}

// 16 位 MOV
var ymovw = []ytab{
	{Yrl, Ynone, Yml, Zr_m, 1},      // 32位寄存器 -> 内存 (word)
	{Yml, Ynone, Yrl, Zm_r, 1},      // 内存 -> 32位寄存器 (word)
	{Yi0, Ynone, Yrl, Zclr, 1},      // 特殊情况: MOVW $0, reg (清零)
	{Yi32, Ynone, Yrl, Zil_rp, 1},   // MOVW $imm, reg (立即数 -> 寄存器)
	{Yi32, Ynone, Yml, Zilo_m, 2},   // MOVW $imm, mem (立即数 -> 内存)
	{Yiauto, Ynone, Yrl, Zaut_r, 2}, // 从自动变量偏移取值 -> 寄存器(编译器生成的自动变量访问)
}

// 32 位 MOV
var ymovl = []ytab{
	{Yrl, Ynone, Yml, Zr_m, 1},      // 寄存器 -> 内存
	{Yml, Ynone, Yrl, Zm_r, 1},      // 内存 -> 寄存器
	{Yi0, Ynone, Yrl, Zclr, 1},      // MOVL $0, reg (清零)
	{Yi32, Ynone, Yrl, Zil_rp, 1},   // MOVL $imm, reg
	{Yi32, Ynone, Yml, Zilo_m, 2},   // MOVL $imm, mem
	{Yml, Ynone, Ymr, Zm_r_xm, 1},   // MMX MOVD
	{Ymr, Ynone, Yml, Zr_m_xm, 1},   // MMX MOVD
	{Yml, Ynone, Yxr, Zm_r_xm, 2},   // XMM MOVD (32 bit)
	{Yxr, Ynone, Yml, Zr_m_xm, 2},   // XMM MOVD (32 bit)
	{Yiauto, Ynone, Yrl, Zaut_r, 2}, // 从自动变量偏移取值 -> 寄存器
}

// RET 指令
var yret = []ytab{
	{Ynone, Ynone, Ynone, Zo_iw, 1}, // 普通的 RET (无参数)
	{Yi32, Ynone, Ynone, Zo_iw, 1},  // RET imm16: 带立即数参数的返回
}

// 64 位移动指令 MOVQ
var ymovq = []ytab{
	// valid in 32-bit mode
	{Ym, Ynone, Ymr, Zm_r_xm_nr, 1},  // 0x6f MMX MOVQ (shorter encoding)
	{Ymr, Ynone, Ym, Zr_m_xm_nr, 1},  // 0x7f MMX MOVQ
	{Yxr, Ynone, Ymr, Zm_r_xm_nr, 2}, // Pf2, 0xd6 MOVDQ2Q
	{Yxm, Ynone, Yxr, Zm_r_xm_nr, 2}, // Pf3, 0x7e MOVQ xmm1/m64 -> xmm2
	{Yxr, Ynone, Yxm, Zr_m_xm_nr, 2}, // Pe, 0xd6 MOVQ xmm1 -> xmm2/m64

	// valid only in 64-bit mode, usually with 64-bit prefix
	{Yrl, Ynone, Yml, Zr_m, 1},      // 0x89
	{Yml, Ynone, Yrl, Zm_r, 1},      // 0x8b
	{Yi0, Ynone, Yrl, Zclr, 1},      // 0x31
	{Ys32, Ynone, Yrl, Zilo_m, 2},   // 32 bit signed 0xc7,(0)
	{Yi64, Ynone, Yrl, Ziq_rp, 1},   // 0xb8 -- 32/64 bit immediate
	{Yi32, Ynone, Yml, Zilo_m, 2},   // 0xc7,(0)
	{Ymm, Ynone, Ymr, Zm_r_xm, 1},   // 0x6e MMX MOVD
	{Ymr, Ynone, Ymm, Zr_m_xm, 1},   // 0x7e MMX MOVD
	{Yml, Ynone, Yxr, Zm_r_xm, 2},   // Pe, 0x6e MOVD xmm load
	{Yxr, Ynone, Yml, Zr_m_xm, 2},   // Pe, 0x7e MOVD xmm store
	{Yiauto, Ynone, Yrl, Zaut_r, 1}, // 0 built-in LEAQ
}

// 内存(64 位或其他宽度) -> 通用寄存器(reg long)
// 比如: MOVQ m, AX
var ym_rl = []ytab{
	{Ym, Ynone, Yrl, Zm_r, 1},
}

// 寄存器 -> 内存
// 比如: MOVQ AX, m
var yrl_m = []ytab{
	{Yrl, Ynone, Ym, Zr_m, 1},
}

// 内存字节(Ymb) -> reg long (Yrl)
// 比如: MOVZBL m8, EAX
var ymb_rl = []ytab{
	{Ymb, Ynone, Yrl, Zmb_r, 1},
}

// 内存 long (Yml，通常指 32/64 位内存值) -> 寄存器
// 比如: MOVL m32, EAX
var yml_rl = []ytab{
	{Yml, Ynone, Yrl, Zm_r, 1},
}

// 寄存器 -> 内存 long
// 比如: MOVL EAX, m32
var yrl_ml = []ytab{
	{Yrl, Ynone, Yml, Zr_m, 1},
}

// 内存 byte 和 byte 寄存器之间的互相拷贝(load/store)
// 名字 yml_mb 可能有误导, 可能不是指 memory long -> memory byte
// 比如 MOVB m8, AL
var yml_mb = []ytab{
	// byte 寄存器 -> 内存 byte
	{Yrb, Ynone, Ymb, Zr_m, 1},
	// 内存 byte -> byte 寄存器
	{Ymb, Ynone, Yrb, Zm_r, 1},
}

// byte 寄存器 -> 内存 byte
// 比如 MOVB AL, m8
var yrb_mb = []ytab{
	{Yrb, Ynone, Ymb, Zr_m, 1},
}

// XCHG 交换两个操作数的值
// 可以是寄存器与寄存器, 寄存器与内存, 也可以和特定寄存器 AX 优化
var yxchg = []ytab{
	{Yax, Ynone, Yrl, Z_rp, 1}, // 	AX <-> 通用寄存器，特例：有更短的编码
	{Yrl, Ynone, Yax, Zrp_, 1}, // 通用寄存器 <-> AX
	{Yrl, Ynone, Yml, Zr_m, 1}, // reg <-> memory long
	{Yml, Ynone, Yrl, Zm_r, 1}, // memory long <-> reg
}

// DIVL 32 位无符号除法
var ydivl = []ytab{
	{Yml, Ynone, Ynone, Zm_o, 2}, // 除数是内存 long (Yml), 结果在 AX/DX 中, DIV m32
}

// DIVB 8 位无符号除法
var ydivb = []ytab{
	{Ymb, Ynone, Ynone, Zm_o, 2}, // 除数是内存 byte
}

// IMUL 有多个变种: 单操作数/双操作数/三操作数
var yimul = []ytab{
	{Yml, Ynone, Ynone, Zm_o, 2},  // 单操作数: IMUL m32, 结果进 AX/DX
	{Yi8, Ynone, Yrl, Zib_rr, 1},  // 双操作数: IMUL r32, r32, imm8 (目标 <- 目标 *imm8)
	{Yi32, Ynone, Yrl, Zil_rr, 1}, // 双操作数: IMUL r32, r32, imm32
	{Yml, Ynone, Yrl, Zm_r, 2},    // 双操作数: IMUL r32, m32 (目标 <- 目标 * m32)
}

// IMUL 三操作数形式
// 将 memory long（m32） * imm8，结果放入目标寄存器
// IMUL r32, m32, imm8
var yimul3 = []ytab{
	{Yi8, Yml, Yrl, Zibm_r, 2},
}

// 生成一个字节
// 比如 .byte 0xXX
var ybyte = []ytab{
	// Yi64: 字面值(立即数) 64 位
	// Zbyte: 这个模板会直接把立即数编码成一个单字节写进指令流(比如 .byte 指令)
	{Yi64, Ynone, Ynone, Zbyte, 1},
}

// IN 指令(从 IO 端口读到 AL/EAX)
var yin = []ytab{
	{Yi32, Ynone, Ynone, Zib_, 1},  // IN AL/EAX, imm8 (端口号是立即数)
	{Ynone, Ynone, Ynone, Zlit, 1}, // IN AL/EAX, DX (端口号在 DX 寄存器里)
}

// INT 指令(软件中断)
var yint = []ytab{
	{Yi32, Ynone, Ynone, Zib_, 1}, // INT imm8 (中断号是立即数)
}

// 32位PUSHL
var ypushl = []ytab{
	{Yrl, Ynone, Ynone, Zrp_, 1},  // PUSH reg (寄存器入栈)
	{Ym, Ynone, Ynone, Zm_o, 2},   // PUSH mem (内存内容入栈)
	{Yi8, Ynone, Ynone, Zib_, 1},  // PUSH imm8 (短立即数)
	{Yi32, Ynone, Ynone, Zil_, 1}, // PUSH imm32 (长立即数)
}

// POPL 指令
var ypopl = []ytab{
	{Ynone, Ynone, Yrl, Z_rp, 1}, // POP 到寄存器
	{Ynone, Ynone, Ym, Zo_m, 2},  // POP 到内存
}

// BSWAP 指令(字节顺序翻转寄存器)
var ybswap = []ytab{
	{Ynone, Ynone, Yrl, Z_rp, 2}, // BSWAP reg (目标是寄存器)
}

// 条件指令类(setcc)
// 例如 SETNE BYTE PTR [rax]
var yscond = []ytab{
	// Zo_m: 使用模板 Zo_m, 需要 modrm
	// zoffset=2 表示指令长度是两字节指令前缀 + modrm
	{Ynone, Ynone, Ymb, Zo_m, 2},
}

// 条件跳转指令(jcc)
//
// Ybr: 目标是分支目标（label）
// Zbr: 表示按条件跳转的模板, 通常是 0F 8X disp32 (或短跳用 7X disp8)
// 三种组合:
// - 无条件数: 常规条件跳转
// - Yi0: 带立即数0, 当作跳转条件
// - Yi1: 带立即数1, zoffset=1 表示使用不同 opcode (例如短跳)
var yjcond = []ytab{
	{Ynone, Ynone, Ybr, Zbr, 0},
	{Yi0, Ynone, Ybr, Zbr, 0},
	{Yi1, Ynone, Ybr, Zbr, 1},
}

// loop 指令
// 例如 loop label
var yloop = []ytab{
	{Ynone, Ynone, Ybr, Zloop, 1},
}

// 调用指令
var ycall = []ytab{
	{Ynone, Ynone, Yml, Zcallindreg, 0}, // 无额外操作数, 目标是内存/间接调用
	{Yrx, Ynone, Yrx, Zcallindreg, 2},   // 寄存器调用
	{Ynone, Ynone, Yindir, Zcallind, 2}, // 无额外操作数, 目标是间接地址调用
	{Ynone, Ynone, Ybr, Zcall, 0},       // 直接调用(标签)
	{Ynone, Ynone, Yi32, Zcallcon, 1},   // 调用符号常量(符号地址)
}

// Duff's device 调用
var yduff = []ytab{
	{Ynone, Ynone, Yi32, Zcallduff, 1},
}

// JMP 指令
var yjmp = []ytab{
	{Ynone, Ynone, Yml, Zo_m64, 2},   // 跳转到内存/间接地址
	{Ynone, Ynone, Ybr, Zjmp, 0},     // 跳转到标签(直接跳转)
	{Ynone, Ynone, Yi32, Zjmpcon, 1}, // 跳转到符号常量
}

// 浮点数移动指令
// d 后缀表示 double, 通常表示双向, 即既支持 load 又支持 store
var yfmvd = []ytab{
	{Ym, Ynone, Yf0, Zm_o, 2},  // 内存 -> ST(0)
	{Yf0, Ynone, Ym, Zo_m, 2},  // ST(0) -> 内存
	{Yrf, Ynone, Yf0, Zm_o, 2}, // ST(i) -> ST(0)
	{Yf0, Ynone, Yrf, Zo_m, 2}, // ST(0) -> ST(i)
}

// 浮点数移动指令
// p 后缀可能是 pop, 浮点指令里经常有 FSTP (存储并弹栈), 或 FLD 指令后面跟弹栈等。
var yfmvdp = []ytab{
	{Yf0, Ynone, Ym, Zo_m, 2},  // ST(0) -> 内存
	{Yf0, Ynone, Yrf, Zo_m, 2}, // ST(0) -> ST(i)
}

// 浮点数移动指令
// 后缀 f 可能表示 fast / forward / from memory
// 一般是支持基本的浮点 load/store，不涉及弹栈，也不涉及寄存器之间交换。
var yfmvf = []ytab{
	{Ym, Ynone, Yf0, Zm_o, 2}, // 内存 -> ST(0)
	{Yf0, Ynone, Ym, Zo_m, 2}, // ST(0) -> 内存
}

// 浮点数移动指令
// 后缀 x 表示 exchange
// 比如浮点指令 FXCH 交换栈顶和 ST(i)
var yfmvx = []ytab{
	{Ym, Ynone, Yf0, Zm_o, 2}, // 内存 -> ST(0)
}

// 浮点数移动指令
// 后缀 p 表示 pop, 对应 ST(0) -> memory 并弹栈
var yfmvp = []ytab{
	{Yf0, Ynone, Ym, Zo_m, 2}, // ST(0) -> 内存
}

// 浮点数移动指令
// 对应 x86 的 FCMOVcc 指令, 根据条件码决定是否从 ST(i) 移动到 ST(0)
var yfcmv = []ytab{
	{Yrf, Ynone, Yf0, Zm_o, 2}, // ST(i) -> ST(0)
}

// 浮点加法(FADD/FADDxx)
var yfadd = []ytab{
	{Ym, Ynone, Yf0, Zm_o, 2},  // 内存值加到 ST(0): FADD mem
	{Yrf, Ynone, Yf0, Zm_o, 2}, // ST(i) 加到 ST(0): FADD ST(i), ST(0)
	{Yf0, Ynone, Yrf, Zo_m, 2}, // ST(0) 加到 ST(i): FADD ST(0), ST(i)
}

// 浮点加并弹栈(FADDP)
var yfaddp = []ytab{
	{Yf0, Ynone, Yrf, Zo_m, 2}, // ST(0) 加到 ST(i), 并弹栈
}

// 浮点交换(FXCH)
var yfxch = []ytab{
	{Yf0, Ynone, Yrf, Zo_m, 2}, // 交换 ST(0) 和 ST(i)：FXCH ST(i)
	{Yrf, Ynone, Yf0, Zm_o, 2}, // ST(i) -> ST(0): 实际也是 FXCH, 只是写法不同
}

// 浮点比较并弹栈(FCOMPP)
var ycompp = []ytab{
	{Yf0, Ynone, Yrf, Zo_m, 2}, // 比较 ST(0) 和 ST(i), 并弹两次栈顶. botch is really f0,f1
}

// 存储浮点状态字(FSTSW)
var ystsw = []ytab{
	{Ynone, Ynone, Ym, Zo_m, 2},  // 将状态字写到内存: FSTSW mem
	{Ynone, Ynone, Yax, Zlit, 1}, // 将状态字写到 AX: FSTSW AX
}

// 控制字指令 STCW / LDCW (Store Control Word / Load Control Word)
var ystcw = []ytab{
	{Ynone, Ynone, Ym, Zo_m, 2}, // 将控制字存储到内存: FSTCW mem
	{Ym, Ynone, Ynone, Zm_o, 2}, // 从内存加载控制字: FLDCW mem
}

// 类似系统指令
// 比如保存/恢复 x87 环境
var ysvrs = []ytab{
	{Ynone, Ynone, Ym, Zo_m, 2}, // 保存状态: FNSAVE/FSTENV mem
	{Ym, Ynone, Ynone, Zm_o, 2}, // 恢复状态: FRSTOR/FLDENV mem
}

// 向量指令里 ymm 和 xmm
var ymm = []ytab{
	{Ymm, Ynone, Ymr, Zm_r_xm, 1}, // memory -> ymm 寄存器(可能是 256 位 AVX)
	{Yxm, Ynone, Yxr, Zm_r_xm, 2}, // memory -> xmm 寄存器(128 位 SSE)
}

// 只涉及 xmm
var yxm = []ytab{
	{Yxm, Ynone, Yxr, Zm_r_xm, 1}, // memory -> xmm 寄存器
}

// 数据转换, 数据搬运指令
// xcv 可能是 convert 缩写
var yxcvm1 = []ytab{
	{Yxm, Ynone, Yxr, Zm_r_xm, 2}, // memory -> xmm
	{Yxm, Ynone, Ymr, Zm_r_xm, 2}, // memory -> memory/寄存器
}

// 向量/浮点数据的转换指令, 第二种操作数模式
var yxcvm2 = []ytab{
	{Yxm, Ynone, Yxr, Zm_r_xm, 2}, // 从 XMM memory -> XMM register, 长度/特征=2
	{Ymm, Ynone, Yxr, Zm_r_xm, 2}, // 从 YMM memory -> XMM register
}

/*
var yxmq = []ytab{
	{Yxm, Ynone, Yxr, Zm_r_xm, 2},
}
*/

// yxr: “XMM register -> XMM register”
// 用来描述涉及 XMM 寄存器作为源和目标的指令(比如 ADDPS xmm1, xmm2)
// 典型 SSE/AVX 指令: 寄存器间运算, 如 ADDPS xmm1, xmm2
var yxr = []ytab{
	{Yxr, Ynone, Yxr, Zm_r_xm, 1},
}

// 对应把向量寄存器里的值写到内存
// 比如 MOVDQA xmm1, m128
// ml 表示 memory low 或 memory load/store
var yxr_ml = []ytab{
	{Yxr, Ynone, Yml, Zr_m_xm, 1},
}

// 对应传统通用寄存器指令(非向量), 如 MOV r/m32, r32
var ymr = []ytab{
	{Ymr, Ynone, Ymr, Zm_r, 1},
}

// 寄存器/内存 -> memory
// ml: memory low
var ymr_ml = []ytab{
	{Ymr, Ynone, Yml, Zr_m_xm, 1},
}

// 对应比较指令时从内存读到寄存器
// 如 CMPPS xmm1, m128
var yxcmp = []ytab{
	{Yxm, Ynone, Yxr, Zm_r_xm, 1},
}

// 对应 PCMPxSTRI/CMPPS 带 immediate mask
// 比如 CMPPS xmm1, m128, imm8
var yxcmpi = []ytab{
	{Yxm, Yxr, Yi8, Zm_r_i_xm, 2},
}

// 对应向量数据移动
var yxmov = []ytab{
	{Yxm, Ynone, Yxr, Zm_r_xm, 1},
	{Yxr, Ynone, Yxm, Zr_m_xm, 1},
}

// 浮点数 → 整数转换
// convert float -> long (int)
// 比如 CVTTSS2SI
var yxcvfl = []ytab{
	{Yxm, Ynone, Yrl, Zm_r_xm, 1},
}

// 整数 → 浮点数
// convert long -> float
// 比如 CVTSI2SS
var yxcvlf = []ytab{
	{Yml, Ynone, Yxr, Zm_r_xm, 1},
}

// 浮点数 → 64位整数
// convert float -> quad (int64)
// 比如 CVTTSD2SI
var yxcvfq = []ytab{
	{Yxm, Ynone, Yrl, Zm_r_xm, 2},
}

// 64位整数 → 浮点数
// convert quad (int64) -> float
// 比如 CVTSI2SD
var yxcvqf = []ytab{
	{Yml, Ynone, Yxr, Zm_r_xm, 2},
}

// 大多数向量指令
// packed single / shift
// 比如 VPADDQ ymm, ymm, m256
var yps = []ytab{
	{Ymm, Ynone, Ymr, Zm_r_xm, 1},
	{Yi8, Ynone, Ymr, Zibo_m_xm, 2},
	{Yxm, Ynone, Yxr, Zm_r_xm, 2},
	{Yi8, Ynone, Yxr, Zibo_m_xm, 3},
}

// 从向量寄存器提取到标量寄存器。
// XMM register → register long
// MOVD r32, xmm
// MOVQ r64, xmm
var yxrrl = []ytab{
	{Yxr, Ynone, Yrl, Zm_r, 1},
}

// memory float packed
// MOVAPS ymm, m256
// MOVUPS ymm, m256
var ymfp = []ytab{
	{Ymm, Ynone, Ymr, Zm_r_3d, 1},
}

// MOVD xmm, r32
// MOVSS xmm, m32
// MOVDQA xmm, m128
var ymrxr = []ytab{
	{Ymr, Ynone, Yxr, Zm_r, 1},
	{Yxm, Ynone, Yxr, Zm_r_xm, 1},
}

// shuffle：重新排列向量里的元素
// PSHUFD xmm, xmm/m128, imm8
// PSHUFW mm, mm/m64, imm8
var ymshuf = []ytab{
	{Yi8, Ymm, Ymr, Zibm_r, 2},
}

// 根据另一个向量按字节重新排列
// PSHUFB xmm, xmm/m128
var ymshufb = []ytab{
	{Yxm, Ynone, Yxr, Zm2_r, 2},
}

// 控制两个向量合并后的排列方式
// SHUFPS xmm, xmm/m128, imm8
// SHUFPD xmm, xmm/m128, imm8
var yxshuf = []ytab{
	{Yu8, Yxm, Yxr, Zibm_r, 2},
}

// 从 XMM 寄存器里提取一个 16 位 word 到通用寄存器, 立即数指定哪个 word
// PEXTRW r32, xmm, imm8
var yextrw = []ytab{
	{Yu8, Yxr, Yrl, Zibm_r, 2},
}

// 把 16 位的值(来自内存或寄存器)插入到 XMM 寄存器指定的位置
// PINSRW xmm, r/m16, imm8
var yinsrw = []ytab{
	{Yu8, Yml, Yxr, Zibm_r, 2},
}

// 把 32 位的值插入到 XMM
// PINSRD xmm, r/m32, imm8
// VPINSRD xmm, xmm, r/m32, imm8
var yinsr = []ytab{
	{Yu8, Ymm, Yxr, Zibm_r, 3},
}

// 对 128 位向量按字节逻辑左移或右移
// PSLLDQ xmm, imm8
// PSRLDQ xmm, imm8
var ypsdq = []ytab{
	{Yi8, Ynone, Yxr, Zibo_m, 2},
}

// 提取向量里每个字节最高位, 打包成整数
// PMOVMSKB r32, xmm
// PMOVMSKB r32, mm
var ymskb = []ytab{
	{Yxr, Ynone, Yrl, Zm_r_xm, 2},
	{Ymr, Ynone, Yrl, Zm_r_xm, 1},
}

// 对内存(或寄存器)做 CRC32 校验, 更新结果到寄存器
// CRC32 r32, r/m8
// CRC32 r32, r/m32
var ycrc32l = []ytab{
	{Yml, Ynone, Yrl, Zlitm_r, 0},
}

// 给缓存做预取提示
// PREFETCHT0 m8
// PREFETCHT1 m8
// PREFETCHNTA m8
var yprefetch = []ytab{
	{Ym, Ynone, Ynone, Zm_o, 2},
}

// AES 指令
// AESENC xmm1, xmm2/m128
// AESENCLAST xmm1, xmm2/m128
// AESDEC xmm1, xmm2/m128
// AESDECLAST xmm1, xmm2/m128
var yaes = []ytab{
	{Yxm, Ynone, Yxr, Zlitm_r, 2},
}

// AES 指令
// AESKEYGENASSIST xmm1, xmm2/m128, imm8
var yaes2 = []ytab{
	{Yu8, Yxm, Yxr, Zibm_r, 2},
}

/*
 * You are doasm, holding in your hand a Prog* with p->as set to, say, ACRC32,
 * and p->from and p->to as operands (Addr*).  The linker scans optab to find
 * the entry with the given p->as and then looks through the ytable for that
 * instruction (the second field in the optab struct) for a line whose first
 * two values match the Ytypes of the p->from and p->to operands.  The function
 * oclass in span.c computes the specific Ytype of an operand and then the set
 * of more general Ytypes that it satisfies is implied by the ycover table, set
 * up in instinit.  For example, oclass distinguishes the constants 0 and 1
 * from the more general 8-bit constants, but instinit says
 *
 *        ycover[Yi0*Ymax + Ys32] = 1;
 *        ycover[Yi1*Ymax + Ys32] = 1;
 *        ycover[Yi8*Ymax + Ys32] = 1;
 *
 * which means that Yi0, Yi1, and Yi8 all count as Ys32 (signed 32)
 * if that's what an instruction can handle.
 *
 * In parallel with the scan through the ytable for the appropriate line, there
 * is a z pointer that starts out pointing at the strange magic byte list in
 * the Optab struct.  With each step past a non-matching ytable line, z
 * advances by the 4th entry in the line.  When a matching line is found, that
 * z pointer has the extra data to use in laying down the instruction bytes.
 * The actual bytes laid down are a function of the 3rd entry in the line (that
 * is, the Ztype) and the z bytes.
 *
 * For example, let's look at AADDL.  The optab line says:
 *        { AADDL,        yaddl,  Px, 0x83,(00),0x05,0x81,(00),0x01,0x03 },
 *
 * and yaddl says
 *        uchar   yaddl[] =
 *        {
 *                Yi8,    Yml,    Zibo_m, 2,
 *                Yi32,   Yax,    Zil_,   1,
 *                Yi32,   Yml,    Zilo_m, 2,
 *                Yrl,    Yml,    Zr_m,   1,
 *                Yml,    Yrl,    Zm_r,   1,
 *                0
 *        };
 *
 * so there are 5 possible types of ADDL instruction that can be laid down, and
 * possible states used to lay them down (Ztype and z pointer, assuming z
 * points at {0x83,(00),0x05,0x81,(00),0x01,0x03}) are:
 *
 *        Yi8, Yml -> Zibo_m, z (0x83, 00)
 *        Yi32, Yax -> Zil_, z+2 (0x05)
 *        Yi32, Yml -> Zilo_m, z+2+1 (0x81, 0x00)
 *        Yrl, Yml -> Zr_m, z+2+1+2 (0x01)
 *        Yml, Yrl -> Zm_r, z+2+1+2+1 (0x03)
 *
 * The Pconstant in the optab line controls the prefix bytes to emit.  That's
 * relatively straightforward as this program goes.
 *
 * The switch on t[2] in doasm implements the various Z cases.  Zibo_m, for
 * example, is an opcode byte (z[0]) then an asmando (which is some kind of
 * encoded addressing mode for the Yml arg), and then a single immediate byte.
 * Zilo_m is the same but a long (32-bit) immediate.
 */

/*
 * 你现在在 doasm 函数中，手里拿着一个 Prog*，其中 p->as 设置为某个指令，比如 ACRC32，
 * 而 p->from 和 p->to 则是操作数（Addr* 类型）。
 *
 * 链接器会先扫描 optab 表，根据 p->as 查找对应的表项。
 * 每个 optab 表项有一个 ytable（第二个字段），它描述了该指令支持的不同操作数字段组合。
 *
 * 接着，doasm 会遍历这个 ytable：
 * - ytable 的每一行前两个字段是操作数的 Ytype 类型（由 span.c 中的 oclass 函数计算）
 * - 如果 p->from 和 p->to 的类型与 ytable 中某一行匹配，就选中该行
 *
 * oclass 会根据操作数具体值判断其类型，例如区分常量 0 和 1 与一般的 8 位常量 Yi8。
 * 而 instinit 初始化时，会设置 ycover 表，指定更宽泛的匹配关系，例如：
 *
 *    ycover[Yi0*Ymax + Ys32] = 1;
 *    ycover[Yi1*Ymax + Ys32] = 1;
 *    ycover[Yi8*Ymax + Ys32] = 1;
 *
 * 意思是：Yi0、Yi1 和 Yi8 都可以算作 Ys32（有符号 32 位数）。
 *
 * 与遍历 ytable 同时，还有一个指针 z，最开始指向 optab 中的“魔法字节”数组。
 * 每次跳过一行不匹配时，就根据该行的第四个字段跳过相应字节。
 * 当找到匹配行时，z 指针正好指向需要用来生成机器码的模板字节。
 * 实际生成的机器码取决于该行的第三个字段（Ztype）和 z 指针指向的模板字节。
 *
 * 举个例子：看 AADDL 指令。
 * optab 表中的一行：
 *    { AADDL, yaddl, Px, 0x83,(00),0x05,0x81,(00),0x01,0x03 },
 *
 * 对应的 yaddl 表：
 *    uchar yaddl[] = {
 *        Yi8,  Yml, Zibo_m, 2,
 *        Yi32, Yax, Zil_,   1,
 *        Yi32, Yml, Zilo_m, 2,
 *        Yrl,  Yml, Zr_m,   1,
 *        Yml,  Yrl, Zm_r,   1,
 *        0
 *    };
 *
 * 表示 ADDL 有 5 种操作数字段组合，每种有对应的模板：
 *
 *    Yi8,  Yml  -> Zibo_m, z          (z 指向 0x83, 00)
 *    Yi32, Yax  -> Zil_,   z+2        (z 指向 0x05)
 *    Yi32, Yml  -> Zilo_m, z+3        (z 指向 0x81, 00)
 *    Yrl,  Yml  -> Zr_m,   z+5        (z 指向 0x01)
 *    Yml,  Yrl  -> Zm_r,   z+6        (z 指向 0x03)
 *
 * optab 中的 Pconstant（这里是 Px）控制是否要生成前缀字节。
 *
 * 在 doasm 函数中根据第三个字段 Ztype 进行 switch：
 * 例如 Zibo_m 表示：
 *   - 先输出一个 opcode 字节（z[0]）
 *   - 然后输出 asmando（根据目标操作数编码的 modrm）
 *   - 最后输出一个立即数字节
 *
 * Zilo_m 类似，只是立即数是 32 位。
 *
 * 整个机制总结：
 * - optab 通过 as 匹配指令
 * - ytable 描述合法的操作数字段组合
 * - ycover 定义更宽松的类型匹配
 * - 匹配成功后，用 Ztype + 模板字节（z）拼出机器码
 */
var optab = []Optab{
	// as, ytab, andproto, opcode
	{AAAA, ynone, P32, [23]uint8{0x37}},
	{AAAD, ynone, P32, [23]uint8{0xd5, 0x0a}},
	{AAAM, ynone, P32, [23]uint8{0xd4, 0x0a}},
	{AAAS, ynone, P32, [23]uint8{0x3f}},
	{AADCB, yxorb, Pb, [23]uint8{0x14, 0x80, 02, 0x10, 0x10}},
	{AADCL, yxorl, Px, [23]uint8{0x83, 02, 0x15, 0x81, 02, 0x11, 0x13}},
	{AADCQ, yxorl, Pw, [23]uint8{0x83, 02, 0x15, 0x81, 02, 0x11, 0x13}},
	{AADCW, yxorl, Pe, [23]uint8{0x83, 02, 0x15, 0x81, 02, 0x11, 0x13}},
	{AADDB, yxorb, Pb, [23]uint8{0x04, 0x80, 00, 0x00, 0x02}},
	{AADDL, yaddl, Px, [23]uint8{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
	{AADDPD, yxm, Pq, [23]uint8{0x58}},
	{AADDPS, yxm, Pm, [23]uint8{0x58}},
	{AADDQ, yaddl, Pw, [23]uint8{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
	{AADDSD, yxm, Pf2, [23]uint8{0x58}},
	{AADDSS, yxm, Pf3, [23]uint8{0x58}},
	{AADDW, yaddl, Pe, [23]uint8{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
	{AADJSP, nil, 0, [23]uint8{}},
	{AANDB, yxorb, Pb, [23]uint8{0x24, 0x80, 04, 0x20, 0x22}},
	{AANDL, yxorl, Px, [23]uint8{0x83, 04, 0x25, 0x81, 04, 0x21, 0x23}},
	{AANDNPD, yxm, Pq, [23]uint8{0x55}},
	{AANDNPS, yxm, Pm, [23]uint8{0x55}},
	{AANDPD, yxm, Pq, [23]uint8{0x54}},
	{AANDPS, yxm, Pq, [23]uint8{0x54}},
	{AANDQ, yxorl, Pw, [23]uint8{0x83, 04, 0x25, 0x81, 04, 0x21, 0x23}},
	{AANDW, yxorl, Pe, [23]uint8{0x83, 04, 0x25, 0x81, 04, 0x21, 0x23}},
	{AARPL, yrl_ml, P32, [23]uint8{0x63}},
	{ABOUNDL, yrl_m, P32, [23]uint8{0x62}},
	{ABOUNDW, yrl_m, Pe, [23]uint8{0x62}},
	{ABSFL, yml_rl, Pm, [23]uint8{0xbc}},
	{ABSFQ, yml_rl, Pw, [23]uint8{0x0f, 0xbc}},
	{ABSFW, yml_rl, Pq, [23]uint8{0xbc}},
	{ABSRL, yml_rl, Pm, [23]uint8{0xbd}},
	{ABSRQ, yml_rl, Pw, [23]uint8{0x0f, 0xbd}},
	{ABSRW, yml_rl, Pq, [23]uint8{0xbd}},
	{ABSWAPL, ybswap, Px, [23]uint8{0x0f, 0xc8}},
	{ABSWAPQ, ybswap, Pw, [23]uint8{0x0f, 0xc8}},
	{ABTCL, ybtl, Pm, [23]uint8{0xba, 07, 0xbb}},
	{ABTCQ, ybtl, Pw, [23]uint8{0x0f, 0xba, 07, 0x0f, 0xbb}},
	{ABTCW, ybtl, Pq, [23]uint8{0xba, 07, 0xbb}},
	{ABTL, ybtl, Pm, [23]uint8{0xba, 04, 0xa3}},
	{ABTQ, ybtl, Pw, [23]uint8{0x0f, 0xba, 04, 0x0f, 0xa3}},
	{ABTRL, ybtl, Pm, [23]uint8{0xba, 06, 0xb3}},
	{ABTRQ, ybtl, Pw, [23]uint8{0x0f, 0xba, 06, 0x0f, 0xb3}},
	{ABTRW, ybtl, Pq, [23]uint8{0xba, 06, 0xb3}},
	{ABTSL, ybtl, Pm, [23]uint8{0xba, 05, 0xab}},
	{ABTSQ, ybtl, Pw, [23]uint8{0x0f, 0xba, 05, 0x0f, 0xab}},
	{ABTSW, ybtl, Pq, [23]uint8{0xba, 05, 0xab}},
	{ABTW, ybtl, Pq, [23]uint8{0xba, 04, 0xa3}},
	{ABYTE, ybyte, Px, [23]uint8{1}},
	{objabi.ACALL, ycall, Px, [23]uint8{0xff, 02, 0xff, 0x15, 0xe8}},
	{ACDQ, ynone, Px, [23]uint8{0x99}},
	{ACLC, ynone, Px, [23]uint8{0xf8}},
	{ACLD, ynone, Px, [23]uint8{0xfc}},
	{ACLI, ynone, Px, [23]uint8{0xfa}},
	{ACLTS, ynone, Pm, [23]uint8{0x06}},
	{ACMC, ynone, Px, [23]uint8{0xf5}},
	{ACMOVLCC, yml_rl, Pm, [23]uint8{0x43}},
	{ACMOVLCS, yml_rl, Pm, [23]uint8{0x42}},
	{ACMOVLEQ, yml_rl, Pm, [23]uint8{0x44}},
	{ACMOVLGE, yml_rl, Pm, [23]uint8{0x4d}},
	{ACMOVLGT, yml_rl, Pm, [23]uint8{0x4f}},
	{ACMOVLHI, yml_rl, Pm, [23]uint8{0x47}},
	{ACMOVLLE, yml_rl, Pm, [23]uint8{0x4e}},
	{ACMOVLLS, yml_rl, Pm, [23]uint8{0x46}},
	{ACMOVLLT, yml_rl, Pm, [23]uint8{0x4c}},
	{ACMOVLMI, yml_rl, Pm, [23]uint8{0x48}},
	{ACMOVLNE, yml_rl, Pm, [23]uint8{0x45}},
	{ACMOVLOC, yml_rl, Pm, [23]uint8{0x41}},
	{ACMOVLOS, yml_rl, Pm, [23]uint8{0x40}},
	{ACMOVLPC, yml_rl, Pm, [23]uint8{0x4b}},
	{ACMOVLPL, yml_rl, Pm, [23]uint8{0x49}},
	{ACMOVLPS, yml_rl, Pm, [23]uint8{0x4a}},
	{ACMOVQCC, yml_rl, Pw, [23]uint8{0x0f, 0x43}},
	{ACMOVQCS, yml_rl, Pw, [23]uint8{0x0f, 0x42}},
	{ACMOVQEQ, yml_rl, Pw, [23]uint8{0x0f, 0x44}},
	{ACMOVQGE, yml_rl, Pw, [23]uint8{0x0f, 0x4d}},
	{ACMOVQGT, yml_rl, Pw, [23]uint8{0x0f, 0x4f}},
	{ACMOVQHI, yml_rl, Pw, [23]uint8{0x0f, 0x47}},
	{ACMOVQLE, yml_rl, Pw, [23]uint8{0x0f, 0x4e}},
	{ACMOVQLS, yml_rl, Pw, [23]uint8{0x0f, 0x46}},
	{ACMOVQLT, yml_rl, Pw, [23]uint8{0x0f, 0x4c}},
	{ACMOVQMI, yml_rl, Pw, [23]uint8{0x0f, 0x48}},
	{ACMOVQNE, yml_rl, Pw, [23]uint8{0x0f, 0x45}},
	{ACMOVQOC, yml_rl, Pw, [23]uint8{0x0f, 0x41}},
	{ACMOVQOS, yml_rl, Pw, [23]uint8{0x0f, 0x40}},
	{ACMOVQPC, yml_rl, Pw, [23]uint8{0x0f, 0x4b}},
	{ACMOVQPL, yml_rl, Pw, [23]uint8{0x0f, 0x49}},
	{ACMOVQPS, yml_rl, Pw, [23]uint8{0x0f, 0x4a}},
	{ACMOVWCC, yml_rl, Pq, [23]uint8{0x43}},
	{ACMOVWCS, yml_rl, Pq, [23]uint8{0x42}},
	{ACMOVWEQ, yml_rl, Pq, [23]uint8{0x44}},
	{ACMOVWGE, yml_rl, Pq, [23]uint8{0x4d}},
	{ACMOVWGT, yml_rl, Pq, [23]uint8{0x4f}},
	{ACMOVWHI, yml_rl, Pq, [23]uint8{0x47}},
	{ACMOVWLE, yml_rl, Pq, [23]uint8{0x4e}},
	{ACMOVWLS, yml_rl, Pq, [23]uint8{0x46}},
	{ACMOVWLT, yml_rl, Pq, [23]uint8{0x4c}},
	{ACMOVWMI, yml_rl, Pq, [23]uint8{0x48}},
	{ACMOVWNE, yml_rl, Pq, [23]uint8{0x45}},
	{ACMOVWOC, yml_rl, Pq, [23]uint8{0x41}},
	{ACMOVWOS, yml_rl, Pq, [23]uint8{0x40}},
	{ACMOVWPC, yml_rl, Pq, [23]uint8{0x4b}},
	{ACMOVWPL, yml_rl, Pq, [23]uint8{0x49}},
	{ACMOVWPS, yml_rl, Pq, [23]uint8{0x4a}},
	{ACMPB, ycmpb, Pb, [23]uint8{0x3c, 0x80, 07, 0x38, 0x3a}},
	{ACMPL, ycmpl, Px, [23]uint8{0x83, 07, 0x3d, 0x81, 07, 0x39, 0x3b}},
	{ACMPPD, yxcmpi, Px, [23]uint8{Pe, 0xc2}},
	{ACMPPS, yxcmpi, Pm, [23]uint8{0xc2, 0}},
	{ACMPQ, ycmpl, Pw, [23]uint8{0x83, 07, 0x3d, 0x81, 07, 0x39, 0x3b}},
	{ACMPSB, ynone, Pb, [23]uint8{0xa6}},
	{ACMPSD, yxcmpi, Px, [23]uint8{Pf2, 0xc2}},
	{ACMPSL, ynone, Px, [23]uint8{0xa7}},
	{ACMPSQ, ynone, Pw, [23]uint8{0xa7}},
	{ACMPSS, yxcmpi, Px, [23]uint8{Pf3, 0xc2}},
	{ACMPSW, ynone, Pe, [23]uint8{0xa7}},
	{ACMPW, ycmpl, Pe, [23]uint8{0x83, 07, 0x3d, 0x81, 07, 0x39, 0x3b}},
	{ACOMISD, yxcmp, Pe, [23]uint8{0x2f}},
	{ACOMISS, yxcmp, Pm, [23]uint8{0x2f}},
	{ACPUID, ynone, Pm, [23]uint8{0xa2}},
	{ACVTPL2PD, yxcvm2, Px, [23]uint8{Pf3, 0xe6, Pe, 0x2a}},
	{ACVTPL2PS, yxcvm2, Pm, [23]uint8{0x5b, 0, 0x2a, 0}},
	{ACVTPD2PL, yxcvm1, Px, [23]uint8{Pf2, 0xe6, Pe, 0x2d}},
	{ACVTPD2PS, yxm, Pe, [23]uint8{0x5a}},
	{ACVTPS2PL, yxcvm1, Px, [23]uint8{Pe, 0x5b, Pm, 0x2d}},
	{ACVTPS2PD, yxm, Pm, [23]uint8{0x5a}},
	{API2FW, ymfp, Px, [23]uint8{0x0c}},
	{ACVTSD2SL, yxcvfl, Pf2, [23]uint8{0x2d}},
	{ACVTSD2SQ, yxcvfq, Pw, [23]uint8{Pf2, 0x2d}},
	{ACVTSD2SS, yxm, Pf2, [23]uint8{0x5a}},
	{ACVTSL2SD, yxcvlf, Pf2, [23]uint8{0x2a}},
	{ACVTSQ2SD, yxcvqf, Pw, [23]uint8{Pf2, 0x2a}},
	{ACVTSL2SS, yxcvlf, Pf3, [23]uint8{0x2a}},
	{ACVTSQ2SS, yxcvqf, Pw, [23]uint8{Pf3, 0x2a}},
	{ACVTSS2SD, yxm, Pf3, [23]uint8{0x5a}},
	{ACVTSS2SL, yxcvfl, Pf3, [23]uint8{0x2d}},
	{ACVTSS2SQ, yxcvfq, Pw, [23]uint8{Pf3, 0x2d}},
	{ACVTTPD2PL, yxcvm1, Px, [23]uint8{Pe, 0xe6, Pe, 0x2c}},
	{ACVTTPS2PL, yxcvm1, Px, [23]uint8{Pf3, 0x5b, Pm, 0x2c}},
	{ACVTTSD2SL, yxcvfl, Pf2, [23]uint8{0x2c}},
	{ACVTTSD2SQ, yxcvfq, Pw, [23]uint8{Pf2, 0x2c}},
	{ACVTTSS2SL, yxcvfl, Pf3, [23]uint8{0x2c}},
	{ACVTTSS2SQ, yxcvfq, Pw, [23]uint8{Pf3, 0x2c}},
	{ACWD, ynone, Pe, [23]uint8{0x99}},
	{ACQO, ynone, Pw, [23]uint8{0x99}},
	{ADAA, ynone, P32, [23]uint8{0x27}},
	{ADAS, ynone, P32, [23]uint8{0x2f}},
	{objabi.ADATA, nil, 0, [23]uint8{}},
	{ADECB, yincb, Pb, [23]uint8{0xfe, 01}},
	{ADECL, yincl, Px1, [23]uint8{0x48, 0xff, 01}},
	{ADECQ, yincq, Pw, [23]uint8{0xff, 01}},
	{ADECW, yincw, Pe, [23]uint8{0xff, 01}},
	{ADIVB, ydivb, Pb, [23]uint8{0xf6, 06}},
	{ADIVL, ydivl, Px, [23]uint8{0xf7, 06}},
	{ADIVPD, yxm, Pe, [23]uint8{0x5e}},
	{ADIVPS, yxm, Pm, [23]uint8{0x5e}},
	{ADIVQ, ydivl, Pw, [23]uint8{0xf7, 06}},
	{ADIVSD, yxm, Pf2, [23]uint8{0x5e}},
	{ADIVSS, yxm, Pf3, [23]uint8{0x5e}},
	{ADIVW, ydivl, Pe, [23]uint8{0xf7, 06}},
	{AEMMS, ynone, Pm, [23]uint8{0x77}},
	{AENTER, nil, 0, [23]uint8{}}, /* botch */
	{AFXRSTOR, ysvrs, Pm, [23]uint8{0xae, 01, 0xae, 01}},
	{AFXSAVE, ysvrs, Pm, [23]uint8{0xae, 00, 0xae, 00}},
	{AFXRSTOR64, ysvrs, Pw, [23]uint8{0x0f, 0xae, 01, 0x0f, 0xae, 01}},
	{AFXSAVE64, ysvrs, Pw, [23]uint8{0x0f, 0xae, 00, 0x0f, 0xae, 00}},
	{objabi.AGLOBL, nil, 0, [23]uint8{}},
	{AHLT, ynone, Px, [23]uint8{0xf4}},
	{AIDIVB, ydivb, Pb, [23]uint8{0xf6, 07}},
	{AIDIVL, ydivl, Px, [23]uint8{0xf7, 07}},
	{AIDIVQ, ydivl, Pw, [23]uint8{0xf7, 07}},
	{AIDIVW, ydivl, Pe, [23]uint8{0xf7, 07}},
	{AIMULB, ydivb, Pb, [23]uint8{0xf6, 05}},
	{AIMULL, yimul, Px, [23]uint8{0xf7, 05, 0x6b, 0x69, Pm, 0xaf}},
	{AIMULQ, yimul, Pw, [23]uint8{0xf7, 05, 0x6b, 0x69, Pm, 0xaf}},
	{AIMULW, yimul, Pe, [23]uint8{0xf7, 05, 0x6b, 0x69, Pm, 0xaf}},
	{AIMUL3Q, yimul3, Pw, [23]uint8{0x6b, 00}},
	{AINB, yin, Pb, [23]uint8{0xe4, 0xec}},
	{AINCB, yincb, Pb, [23]uint8{0xfe, 00}},
	{AINCL, yincl, Px1, [23]uint8{0x40, 0xff, 00}},
	{AINCQ, yincq, Pw, [23]uint8{0xff, 00}},
	{AINCW, yincw, Pe, [23]uint8{0xff, 00}},
	{AINL, yin, Px, [23]uint8{0xe5, 0xed}},
	{AINSB, ynone, Pb, [23]uint8{0x6c}},
	{AINSL, ynone, Px, [23]uint8{0x6d}},
	{AINSW, ynone, Pe, [23]uint8{0x6d}},
	{AINT, yint, Px, [23]uint8{0xcd}},
	{AINTO, ynone, P32, [23]uint8{0xce}},
	{AINW, yin, Pe, [23]uint8{0xe5, 0xed}},
	{AIRETL, ynone, Px, [23]uint8{0xcf}},
	{AIRETQ, ynone, Pw, [23]uint8{0xcf}},
	{AIRETW, ynone, Pe, [23]uint8{0xcf}},
	{AJCC, yjcond, Px, [23]uint8{0x73, 0x83, 00}},
	{AJCS, yjcond, Px, [23]uint8{0x72, 0x82}},
	{AJCXZL, yloop, Px, [23]uint8{0xe3}},
	{AJCXZW, yloop, Px, [23]uint8{0xe3}},
	{AJCXZQ, yloop, Px, [23]uint8{0xe3}},
	{AJEQ, yjcond, Px, [23]uint8{0x74, 0x84}},
	{AJGE, yjcond, Px, [23]uint8{0x7d, 0x8d}},
	{AJGT, yjcond, Px, [23]uint8{0x7f, 0x8f}},
	{AJHI, yjcond, Px, [23]uint8{0x77, 0x87}},
	{AJLE, yjcond, Px, [23]uint8{0x7e, 0x8e}},
	{AJLS, yjcond, Px, [23]uint8{0x76, 0x86}},
	{AJLT, yjcond, Px, [23]uint8{0x7c, 0x8c}},
	{AJMI, yjcond, Px, [23]uint8{0x78, 0x88}},
	{objabi.AJMP, yjmp, Px, [23]uint8{0xff, 04, 0xeb, 0xe9}},
	{AJNE, yjcond, Px, [23]uint8{0x75, 0x85}},
	{AJOC, yjcond, Px, [23]uint8{0x71, 0x81, 00}},
	{AJOS, yjcond, Px, [23]uint8{0x70, 0x80, 00}},
	{AJPC, yjcond, Px, [23]uint8{0x7b, 0x8b}},
	{AJPL, yjcond, Px, [23]uint8{0x79, 0x89}},
	{AJPS, yjcond, Px, [23]uint8{0x7a, 0x8a}},
	{ALAHF, ynone, Px, [23]uint8{0x9f}},
	{ALARL, yml_rl, Pm, [23]uint8{0x02}},
	{ALARW, yml_rl, Pq, [23]uint8{0x02}},
	{ALDMXCSR, ysvrs, Pm, [23]uint8{0xae, 02, 0xae, 02}},
	{ALEAL, ym_rl, Px, [23]uint8{0x8d}},
	{ALEAQ, ym_rl, Pw, [23]uint8{0x8d}},
	{ALEAVEL, ynone, P32, [23]uint8{0xc9}},
	{ALEAVEQ, ynone, Py, [23]uint8{0xc9}},
	{ALEAVEW, ynone, Pe, [23]uint8{0xc9}},
	{ALEAW, ym_rl, Pe, [23]uint8{0x8d}},
	{ALOCK, ynone, Px, [23]uint8{0xf0}},
	{ALODSB, ynone, Pb, [23]uint8{0xac}},
	{ALODSL, ynone, Px, [23]uint8{0xad}},
	{ALODSQ, ynone, Pw, [23]uint8{0xad}},
	{ALODSW, ynone, Pe, [23]uint8{0xad}},
	{ALONG, ybyte, Px, [23]uint8{4}},
	{ALOOP, yloop, Px, [23]uint8{0xe2}},
	{ALOOPEQ, yloop, Px, [23]uint8{0xe1}},
	{ALOOPNE, yloop, Px, [23]uint8{0xe0}},
	{ALSLL, yml_rl, Pm, [23]uint8{0x03}},
	{ALSLW, yml_rl, Pq, [23]uint8{0x03}},
	{AMASKMOVOU, yxr, Pe, [23]uint8{0xf7}},
	{AMASKMOVQ, ymr, Pm, [23]uint8{0xf7}},
	{AMAXPD, yxm, Pe, [23]uint8{0x5f}},
	{AMAXPS, yxm, Pm, [23]uint8{0x5f}},
	{AMAXSD, yxm, Pf2, [23]uint8{0x5f}},
	{AMAXSS, yxm, Pf3, [23]uint8{0x5f}},
	{AMINPD, yxm, Pe, [23]uint8{0x5d}},
	{AMINPS, yxm, Pm, [23]uint8{0x5d}},
	{AMINSD, yxm, Pf2, [23]uint8{0x5d}},
	{AMINSS, yxm, Pf3, [23]uint8{0x5d}},
	{AMOVAPD, yxmov, Pe, [23]uint8{0x28, 0x29}},
	{AMOVAPS, yxmov, Pm, [23]uint8{0x28, 0x29}},
	{AMOVB, ymovb, Pb, [23]uint8{0x88, 0x8a, 0xb0, 0xc6, 00}},
	{AMOVBLSX, ymb_rl, Pm, [23]uint8{0xbe}},
	{AMOVBLZX, ymb_rl, Pm, [23]uint8{0xb6}},
	{AMOVBQSX, ymb_rl, Pw, [23]uint8{0x0f, 0xbe}},
	{AMOVBQZX, ymb_rl, Pm, [23]uint8{0xb6}},
	{AMOVBWSX, ymb_rl, Pq, [23]uint8{0xbe}},
	{AMOVBWZX, ymb_rl, Pq, [23]uint8{0xb6}},
	{AMOVO, yxmov, Pe, [23]uint8{0x6f, 0x7f}},
	{AMOVOU, yxmov, Pf3, [23]uint8{0x6f, 0x7f}},
	{AMOVHLPS, yxr, Pm, [23]uint8{0x12}},
	{AMOVHPD, yxmov, Pe, [23]uint8{0x16, 0x17}},
	{AMOVHPS, yxmov, Pm, [23]uint8{0x16, 0x17}},
	{AMOVL, ymovl, Px, [23]uint8{0x89, 0x8b, 0x31, 0xb8, 0xc7, 00, 0x6e, 0x7e, Pe, 0x6e, Pe, 0x7e, 0}},
	{AMOVLHPS, yxr, Pm, [23]uint8{0x16}},
	{AMOVLPD, yxmov, Pe, [23]uint8{0x12, 0x13}},
	{AMOVLPS, yxmov, Pm, [23]uint8{0x12, 0x13}},
	{AMOVLQSX, yml_rl, Pw, [23]uint8{0x63}},
	{AMOVLQZX, yml_rl, Px, [23]uint8{0x8b}},
	{AMOVMSKPD, yxrrl, Pq, [23]uint8{0x50}},
	{AMOVMSKPS, yxrrl, Pm, [23]uint8{0x50}},
	{AMOVNTO, yxr_ml, Pe, [23]uint8{0xe7}},
	{AMOVNTPD, yxr_ml, Pe, [23]uint8{0x2b}},
	{AMOVNTPS, yxr_ml, Pm, [23]uint8{0x2b}},
	{AMOVNTQ, ymr_ml, Pm, [23]uint8{0xe7}},
	{AMOVQ, ymovq, Pw8, [23]uint8{0x6f, 0x7f, Pf2, 0xd6, Pf3, 0x7e, Pe, 0xd6, 0x89, 0x8b, 0x31, 0xc7, 00, 0xb8, 0xc7, 00, 0x6e, 0x7e, Pe, 0x6e, Pe, 0x7e, 0}},
	{AMOVQOZX, ymrxr, Pf3, [23]uint8{0xd6, 0x7e}},
	{AMOVSB, ynone, Pb, [23]uint8{0xa4}},
	{AMOVSD, yxmov, Pf2, [23]uint8{0x10, 0x11}},
	{AMOVSL, ynone, Px, [23]uint8{0xa5}},
	{AMOVSQ, ynone, Pw, [23]uint8{0xa5}},
	{AMOVSS, yxmov, Pf3, [23]uint8{0x10, 0x11}},
	{AMOVSW, ynone, Pe, [23]uint8{0xa5}},
	{AMOVUPD, yxmov, Pe, [23]uint8{0x10, 0x11}},
	{AMOVUPS, yxmov, Pm, [23]uint8{0x10, 0x11}},
	{AMOVW, ymovw, Pe, [23]uint8{0x89, 0x8b, 0x31, 0xb8, 0xc7, 00, 0}},
	{AMOVWLSX, yml_rl, Pm, [23]uint8{0xbf}},
	{AMOVWLZX, yml_rl, Pm, [23]uint8{0xb7}},
	{AMOVWQSX, yml_rl, Pw, [23]uint8{0x0f, 0xbf}},
	{AMOVWQZX, yml_rl, Pw, [23]uint8{0x0f, 0xb7}},
	{AMULB, ydivb, Pb, [23]uint8{0xf6, 04}},
	{AMULL, ydivl, Px, [23]uint8{0xf7, 04}},
	{AMULPD, yxm, Pe, [23]uint8{0x59}},
	{AMULPS, yxm, Ym, [23]uint8{0x59}},
	{AMULQ, ydivl, Pw, [23]uint8{0xf7, 04}},
	{AMULSD, yxm, Pf2, [23]uint8{0x59}},
	{AMULSS, yxm, Pf3, [23]uint8{0x59}},
	{AMULW, ydivl, Pe, [23]uint8{0xf7, 04}},
	{ANEGB, yscond, Pb, [23]uint8{0xf6, 03}},
	{ANEGL, yscond, Px, [23]uint8{0xf7, 03}},
	{ANEGQ, yscond, Pw, [23]uint8{0xf7, 03}},
	{ANEGW, yscond, Pe, [23]uint8{0xf7, 03}},
	{objabi.ANOP, ynop, Px, [23]uint8{0, 0}},
	{ANOTB, yscond, Pb, [23]uint8{0xf6, 02}},
	{ANOTL, yscond, Px, [23]uint8{0xf7, 02}}, // TODO(chai2010): yscond is wrong here.
	{ANOTQ, yscond, Pw, [23]uint8{0xf7, 02}},
	{ANOTW, yscond, Pe, [23]uint8{0xf7, 02}},
	{AORB, yxorb, Pb, [23]uint8{0x0c, 0x80, 01, 0x08, 0x0a}},
	{AORL, yxorl, Px, [23]uint8{0x83, 01, 0x0d, 0x81, 01, 0x09, 0x0b}},
	{AORPD, yxm, Pq, [23]uint8{0x56}},
	{AORPS, yxm, Pm, [23]uint8{0x56}},
	{AORQ, yxorl, Pw, [23]uint8{0x83, 01, 0x0d, 0x81, 01, 0x09, 0x0b}},
	{AORW, yxorl, Pe, [23]uint8{0x83, 01, 0x0d, 0x81, 01, 0x09, 0x0b}},
	{AOUTB, yin, Pb, [23]uint8{0xe6, 0xee}},
	{AOUTL, yin, Px, [23]uint8{0xe7, 0xef}},
	{AOUTSB, ynone, Pb, [23]uint8{0x6e}},
	{AOUTSL, ynone, Px, [23]uint8{0x6f}},
	{AOUTSW, ynone, Pe, [23]uint8{0x6f}},
	{AOUTW, yin, Pe, [23]uint8{0xe7, 0xef}},
	{APACKSSLW, ymm, Py1, [23]uint8{0x6b, Pe, 0x6b}},
	{APACKSSWB, ymm, Py1, [23]uint8{0x63, Pe, 0x63}},
	{APACKUSWB, ymm, Py1, [23]uint8{0x67, Pe, 0x67}},
	{APADDB, ymm, Py1, [23]uint8{0xfc, Pe, 0xfc}},
	{APADDL, ymm, Py1, [23]uint8{0xfe, Pe, 0xfe}},
	{APADDQ, yxm, Pe, [23]uint8{0xd4}},
	{APADDSB, ymm, Py1, [23]uint8{0xec, Pe, 0xec}},
	{APADDSW, ymm, Py1, [23]uint8{0xed, Pe, 0xed}},
	{APADDUSB, ymm, Py1, [23]uint8{0xdc, Pe, 0xdc}},
	{APADDUSW, ymm, Py1, [23]uint8{0xdd, Pe, 0xdd}},
	{APADDW, ymm, Py1, [23]uint8{0xfd, Pe, 0xfd}},
	{APAND, ymm, Py1, [23]uint8{0xdb, Pe, 0xdb}},
	{APANDN, ymm, Py1, [23]uint8{0xdf, Pe, 0xdf}},
	{APAUSE, ynone, Px, [23]uint8{0xf3, 0x90}},
	{APAVGB, ymm, Py1, [23]uint8{0xe0, Pe, 0xe0}},
	{APAVGW, ymm, Py1, [23]uint8{0xe3, Pe, 0xe3}},
	{APCMPEQB, ymm, Py1, [23]uint8{0x74, Pe, 0x74}},
	{APCMPEQL, ymm, Py1, [23]uint8{0x76, Pe, 0x76}},
	{APCMPEQW, ymm, Py1, [23]uint8{0x75, Pe, 0x75}},
	{APCMPGTB, ymm, Py1, [23]uint8{0x64, Pe, 0x64}},
	{APCMPGTL, ymm, Py1, [23]uint8{0x66, Pe, 0x66}},
	{APCMPGTW, ymm, Py1, [23]uint8{0x65, Pe, 0x65}},
	{APEXTRW, yextrw, Pq, [23]uint8{0xc5, 00}},
	{APF2IL, ymfp, Px, [23]uint8{0x1d}},
	{APF2IW, ymfp, Px, [23]uint8{0x1c}},
	{API2FL, ymfp, Px, [23]uint8{0x0d}},
	{APFACC, ymfp, Px, [23]uint8{0xae}},
	{APFADD, ymfp, Px, [23]uint8{0x9e}},
	{APFCMPEQ, ymfp, Px, [23]uint8{0xb0}},
	{APFCMPGE, ymfp, Px, [23]uint8{0x90}},
	{APFCMPGT, ymfp, Px, [23]uint8{0xa0}},
	{APFMAX, ymfp, Px, [23]uint8{0xa4}},
	{APFMIN, ymfp, Px, [23]uint8{0x94}},
	{APFMUL, ymfp, Px, [23]uint8{0xb4}},
	{APFNACC, ymfp, Px, [23]uint8{0x8a}},
	{APFPNACC, ymfp, Px, [23]uint8{0x8e}},
	{APFRCP, ymfp, Px, [23]uint8{0x96}},
	{APFRCPIT1, ymfp, Px, [23]uint8{0xa6}},
	{APFRCPI2T, ymfp, Px, [23]uint8{0xb6}},
	{APFRSQIT1, ymfp, Px, [23]uint8{0xa7}},
	{APFRSQRT, ymfp, Px, [23]uint8{0x97}},
	{APFSUB, ymfp, Px, [23]uint8{0x9a}},
	{APFSUBR, ymfp, Px, [23]uint8{0xaa}},
	{APINSRW, yinsrw, Pq, [23]uint8{0xc4, 00}},
	{APINSRD, yinsr, Pq, [23]uint8{0x3a, 0x22, 00}},
	{APINSRQ, yinsr, Pq3, [23]uint8{0x3a, 0x22, 00}},
	{APMADDWL, ymm, Py1, [23]uint8{0xf5, Pe, 0xf5}},
	{APMAXSW, yxm, Pe, [23]uint8{0xee}},
	{APMAXUB, yxm, Pe, [23]uint8{0xde}},
	{APMINSW, yxm, Pe, [23]uint8{0xea}},
	{APMINUB, yxm, Pe, [23]uint8{0xda}},
	{APMOVMSKB, ymskb, Px, [23]uint8{Pe, 0xd7, 0xd7}},
	{APMULHRW, ymfp, Px, [23]uint8{0xb7}},
	{APMULHUW, ymm, Py1, [23]uint8{0xe4, Pe, 0xe4}},
	{APMULHW, ymm, Py1, [23]uint8{0xe5, Pe, 0xe5}},
	{APMULLW, ymm, Py1, [23]uint8{0xd5, Pe, 0xd5}},
	{APMULULQ, ymm, Py1, [23]uint8{0xf4, Pe, 0xf4}},
	{APOPAL, ynone, P32, [23]uint8{0x61}},
	{APOPAW, ynone, Pe, [23]uint8{0x61}},
	{APOPFL, ynone, P32, [23]uint8{0x9d}},
	{APOPFQ, ynone, Py, [23]uint8{0x9d}},
	{APOPFW, ynone, Pe, [23]uint8{0x9d}},
	{APOPL, ypopl, P32, [23]uint8{0x58, 0x8f, 00}},
	{APOPQ, ypopl, Py, [23]uint8{0x58, 0x8f, 00}},
	{APOPW, ypopl, Pe, [23]uint8{0x58, 0x8f, 00}},
	{APOR, ymm, Py1, [23]uint8{0xeb, Pe, 0xeb}},
	{APSADBW, yxm, Pq, [23]uint8{0xf6}},
	{APSHUFHW, yxshuf, Pf3, [23]uint8{0x70, 00}},
	{APSHUFL, yxshuf, Pq, [23]uint8{0x70, 00}},
	{APSHUFLW, yxshuf, Pf2, [23]uint8{0x70, 00}},
	{APSHUFW, ymshuf, Pm, [23]uint8{0x70, 00}},
	{APSHUFB, ymshufb, Pq, [23]uint8{0x38, 0x00}},
	{APSLLO, ypsdq, Pq, [23]uint8{0x73, 07}},
	{APSLLL, yps, Py3, [23]uint8{0xf2, 0x72, 06, Pe, 0xf2, Pe, 0x72, 06}},
	{APSLLQ, yps, Py3, [23]uint8{0xf3, 0x73, 06, Pe, 0xf3, Pe, 0x73, 06}},
	{APSLLW, yps, Py3, [23]uint8{0xf1, 0x71, 06, Pe, 0xf1, Pe, 0x71, 06}},
	{APSRAL, yps, Py3, [23]uint8{0xe2, 0x72, 04, Pe, 0xe2, Pe, 0x72, 04}},
	{APSRAW, yps, Py3, [23]uint8{0xe1, 0x71, 04, Pe, 0xe1, Pe, 0x71, 04}},
	{APSRLO, ypsdq, Pq, [23]uint8{0x73, 03}},
	{APSRLL, yps, Py3, [23]uint8{0xd2, 0x72, 02, Pe, 0xd2, Pe, 0x72, 02}},
	{APSRLQ, yps, Py3, [23]uint8{0xd3, 0x73, 02, Pe, 0xd3, Pe, 0x73, 02}},
	{APSRLW, yps, Py3, [23]uint8{0xd1, 0x71, 02, Pe, 0xe1, Pe, 0x71, 02}},
	{APSUBB, yxm, Pe, [23]uint8{0xf8}},
	{APSUBL, yxm, Pe, [23]uint8{0xfa}},
	{APSUBQ, yxm, Pe, [23]uint8{0xfb}},
	{APSUBSB, yxm, Pe, [23]uint8{0xe8}},
	{APSUBSW, yxm, Pe, [23]uint8{0xe9}},
	{APSUBUSB, yxm, Pe, [23]uint8{0xd8}},
	{APSUBUSW, yxm, Pe, [23]uint8{0xd9}},
	{APSUBW, yxm, Pe, [23]uint8{0xf9}},
	{APSWAPL, ymfp, Px, [23]uint8{0xbb}},
	{APUNPCKHBW, ymm, Py1, [23]uint8{0x68, Pe, 0x68}},
	{APUNPCKHLQ, ymm, Py1, [23]uint8{0x6a, Pe, 0x6a}},
	{APUNPCKHQDQ, yxm, Pe, [23]uint8{0x6d}},
	{APUNPCKHWL, ymm, Py1, [23]uint8{0x69, Pe, 0x69}},
	{APUNPCKLBW, ymm, Py1, [23]uint8{0x60, Pe, 0x60}},
	{APUNPCKLLQ, ymm, Py1, [23]uint8{0x62, Pe, 0x62}},
	{APUNPCKLQDQ, yxm, Pe, [23]uint8{0x6c}},
	{APUNPCKLWL, ymm, Py1, [23]uint8{0x61, Pe, 0x61}},
	{APUSHAL, ynone, P32, [23]uint8{0x60}},
	{APUSHAW, ynone, Pe, [23]uint8{0x60}},
	{APUSHFL, ynone, P32, [23]uint8{0x9c}},
	{APUSHFQ, ynone, Py, [23]uint8{0x9c}},
	{APUSHFW, ynone, Pe, [23]uint8{0x9c}},
	{APUSHL, ypushl, P32, [23]uint8{0x50, 0xff, 06, 0x6a, 0x68}},
	{APUSHQ, ypushl, Py, [23]uint8{0x50, 0xff, 06, 0x6a, 0x68}},
	{APUSHW, ypushl, Pe, [23]uint8{0x50, 0xff, 06, 0x6a, 0x68}},
	{APXOR, ymm, Py1, [23]uint8{0xef, Pe, 0xef}},
	{AQUAD, ybyte, Px, [23]uint8{8}},
	{ARCLB, yshb, Pb, [23]uint8{0xd0, 02, 0xc0, 02, 0xd2, 02}},
	{ARCLL, yshl, Px, [23]uint8{0xd1, 02, 0xc1, 02, 0xd3, 02, 0xd3, 02}},
	{ARCLQ, yshl, Pw, [23]uint8{0xd1, 02, 0xc1, 02, 0xd3, 02, 0xd3, 02}},
	{ARCLW, yshl, Pe, [23]uint8{0xd1, 02, 0xc1, 02, 0xd3, 02, 0xd3, 02}},
	{ARCPPS, yxm, Pm, [23]uint8{0x53}},
	{ARCPSS, yxm, Pf3, [23]uint8{0x53}},
	{ARCRB, yshb, Pb, [23]uint8{0xd0, 03, 0xc0, 03, 0xd2, 03}},
	{ARCRL, yshl, Px, [23]uint8{0xd1, 03, 0xc1, 03, 0xd3, 03, 0xd3, 03}},
	{ARCRQ, yshl, Pw, [23]uint8{0xd1, 03, 0xc1, 03, 0xd3, 03, 0xd3, 03}},
	{ARCRW, yshl, Pe, [23]uint8{0xd1, 03, 0xc1, 03, 0xd3, 03, 0xd3, 03}},
	{AREP, ynone, Px, [23]uint8{0xf3}},
	{AREPN, ynone, Px, [23]uint8{0xf2}},
	{objabi.ARET, ynone, Px, [23]uint8{0xc3}},
	{ARETFW, yret, Pe, [23]uint8{0xcb, 0xca}},
	{ARETFL, yret, Px, [23]uint8{0xcb, 0xca}},
	{ARETFQ, yret, Pw, [23]uint8{0xcb, 0xca}},
	{AROLB, yshb, Pb, [23]uint8{0xd0, 00, 0xc0, 00, 0xd2, 00}},
	{AROLL, yshl, Px, [23]uint8{0xd1, 00, 0xc1, 00, 0xd3, 00, 0xd3, 00}},
	{AROLQ, yshl, Pw, [23]uint8{0xd1, 00, 0xc1, 00, 0xd3, 00, 0xd3, 00}},
	{AROLW, yshl, Pe, [23]uint8{0xd1, 00, 0xc1, 00, 0xd3, 00, 0xd3, 00}},
	{ARORB, yshb, Pb, [23]uint8{0xd0, 01, 0xc0, 01, 0xd2, 01}},
	{ARORL, yshl, Px, [23]uint8{0xd1, 01, 0xc1, 01, 0xd3, 01, 0xd3, 01}},
	{ARORQ, yshl, Pw, [23]uint8{0xd1, 01, 0xc1, 01, 0xd3, 01, 0xd3, 01}},
	{ARORW, yshl, Pe, [23]uint8{0xd1, 01, 0xc1, 01, 0xd3, 01, 0xd3, 01}},
	{ARSQRTPS, yxm, Pm, [23]uint8{0x52}},
	{ARSQRTSS, yxm, Pf3, [23]uint8{0x52}},
	{ASAHF, ynone, Px1, [23]uint8{0x9e, 00, 0x86, 0xe0, 0x50, 0x9d}}, /* XCHGB AH,AL; PUSH AX; POPFL */
	{ASALB, yshb, Pb, [23]uint8{0xd0, 04, 0xc0, 04, 0xd2, 04}},
	{ASALL, yshl, Px, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASALQ, yshl, Pw, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASALW, yshl, Pe, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASARB, yshb, Pb, [23]uint8{0xd0, 07, 0xc0, 07, 0xd2, 07}},
	{ASARL, yshl, Px, [23]uint8{0xd1, 07, 0xc1, 07, 0xd3, 07, 0xd3, 07}},
	{ASARQ, yshl, Pw, [23]uint8{0xd1, 07, 0xc1, 07, 0xd3, 07, 0xd3, 07}},
	{ASARW, yshl, Pe, [23]uint8{0xd1, 07, 0xc1, 07, 0xd3, 07, 0xd3, 07}},
	{ASBBB, yxorb, Pb, [23]uint8{0x1c, 0x80, 03, 0x18, 0x1a}},
	{ASBBL, yxorl, Px, [23]uint8{0x83, 03, 0x1d, 0x81, 03, 0x19, 0x1b}},
	{ASBBQ, yxorl, Pw, [23]uint8{0x83, 03, 0x1d, 0x81, 03, 0x19, 0x1b}},
	{ASBBW, yxorl, Pe, [23]uint8{0x83, 03, 0x1d, 0x81, 03, 0x19, 0x1b}},
	{ASCASB, ynone, Pb, [23]uint8{0xae}},
	{ASCASL, ynone, Px, [23]uint8{0xaf}},
	{ASCASQ, ynone, Pw, [23]uint8{0xaf}},
	{ASCASW, ynone, Pe, [23]uint8{0xaf}},
	{ASETCC, yscond, Pb, [23]uint8{0x0f, 0x93, 00}},
	{ASETCS, yscond, Pb, [23]uint8{0x0f, 0x92, 00}},
	{ASETEQ, yscond, Pb, [23]uint8{0x0f, 0x94, 00}},
	{ASETGE, yscond, Pb, [23]uint8{0x0f, 0x9d, 00}},
	{ASETGT, yscond, Pb, [23]uint8{0x0f, 0x9f, 00}},
	{ASETHI, yscond, Pb, [23]uint8{0x0f, 0x97, 00}},
	{ASETLE, yscond, Pb, [23]uint8{0x0f, 0x9e, 00}},
	{ASETLS, yscond, Pb, [23]uint8{0x0f, 0x96, 00}},
	{ASETLT, yscond, Pb, [23]uint8{0x0f, 0x9c, 00}},
	{ASETMI, yscond, Pb, [23]uint8{0x0f, 0x98, 00}},
	{ASETNE, yscond, Pb, [23]uint8{0x0f, 0x95, 00}},
	{ASETOC, yscond, Pb, [23]uint8{0x0f, 0x91, 00}},
	{ASETOS, yscond, Pb, [23]uint8{0x0f, 0x90, 00}},
	{ASETPC, yscond, Pb, [23]uint8{0x0f, 0x9b, 00}},
	{ASETPL, yscond, Pb, [23]uint8{0x0f, 0x99, 00}},
	{ASETPS, yscond, Pb, [23]uint8{0x0f, 0x9a, 00}},
	{ASHLB, yshb, Pb, [23]uint8{0xd0, 04, 0xc0, 04, 0xd2, 04}},
	{ASHLL, yshl, Px, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASHLQ, yshl, Pw, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASHLW, yshl, Pe, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASHRB, yshb, Pb, [23]uint8{0xd0, 05, 0xc0, 05, 0xd2, 05}},
	{ASHRL, yshl, Px, [23]uint8{0xd1, 05, 0xc1, 05, 0xd3, 05, 0xd3, 05}},
	{ASHRQ, yshl, Pw, [23]uint8{0xd1, 05, 0xc1, 05, 0xd3, 05, 0xd3, 05}},
	{ASHRW, yshl, Pe, [23]uint8{0xd1, 05, 0xc1, 05, 0xd3, 05, 0xd3, 05}},
	{ASHUFPD, yxshuf, Pq, [23]uint8{0xc6, 00}},
	{ASHUFPS, yxshuf, Pm, [23]uint8{0xc6, 00}},
	{ASQRTPD, yxm, Pe, [23]uint8{0x51}},
	{ASQRTPS, yxm, Pm, [23]uint8{0x51}},
	{ASQRTSD, yxm, Pf2, [23]uint8{0x51}},
	{ASQRTSS, yxm, Pf3, [23]uint8{0x51}},
	{ASTC, ynone, Px, [23]uint8{0xf9}},
	{ASTD, ynone, Px, [23]uint8{0xfd}},
	{ASTI, ynone, Px, [23]uint8{0xfb}},
	{ASTMXCSR, ysvrs, Pm, [23]uint8{0xae, 03, 0xae, 03}},
	{ASTOSB, ynone, Pb, [23]uint8{0xaa}},
	{ASTOSL, ynone, Px, [23]uint8{0xab}},
	{ASTOSQ, ynone, Pw, [23]uint8{0xab}},
	{ASTOSW, ynone, Pe, [23]uint8{0xab}},
	{ASUBB, yxorb, Pb, [23]uint8{0x2c, 0x80, 05, 0x28, 0x2a}},
	{ASUBL, yaddl, Px, [23]uint8{0x83, 05, 0x2d, 0x81, 05, 0x29, 0x2b}},
	{ASUBPD, yxm, Pe, [23]uint8{0x5c}},
	{ASUBPS, yxm, Pm, [23]uint8{0x5c}},
	{ASUBQ, yaddl, Pw, [23]uint8{0x83, 05, 0x2d, 0x81, 05, 0x29, 0x2b}},
	{ASUBSD, yxm, Pf2, [23]uint8{0x5c}},
	{ASUBSS, yxm, Pf3, [23]uint8{0x5c}},
	{ASUBW, yaddl, Pe, [23]uint8{0x83, 05, 0x2d, 0x81, 05, 0x29, 0x2b}},
	{ASWAPGS, ynone, Pm, [23]uint8{0x01, 0xf8}},
	{ASYSCALL, ynone, Px, [23]uint8{0x0f, 0x05}}, /* fast syscall */
	{ATESTB, ytestb, Pb, [23]uint8{0xa8, 0xf6, 00, 0x84, 0x84}},
	{ATESTL, ytestl, Px, [23]uint8{0xa9, 0xf7, 00, 0x85, 0x85}},
	{ATESTQ, ytestl, Pw, [23]uint8{0xa9, 0xf7, 00, 0x85, 0x85}},
	{ATESTW, ytestl, Pe, [23]uint8{0xa9, 0xf7, 00, 0x85, 0x85}},
	{objabi.ATEXT, ytext, Px, [23]uint8{}},
	{AUCOMISD, yxcmp, Pe, [23]uint8{0x2e}},
	{AUCOMISS, yxcmp, Pm, [23]uint8{0x2e}},
	{AUNPCKHPD, yxm, Pe, [23]uint8{0x15}},
	{AUNPCKHPS, yxm, Pm, [23]uint8{0x15}},
	{AUNPCKLPD, yxm, Pe, [23]uint8{0x14}},
	{AUNPCKLPS, yxm, Pm, [23]uint8{0x14}},
	{AVERR, ydivl, Pm, [23]uint8{0x00, 04}},
	{AVERW, ydivl, Pm, [23]uint8{0x00, 05}},
	{AWAIT, ynone, Px, [23]uint8{0x9b}},
	{AWORD, ybyte, Px, [23]uint8{2}},
	{AXCHGB, yml_mb, Pb, [23]uint8{0x86, 0x86}},
	{AXCHGL, yxchg, Px, [23]uint8{0x90, 0x90, 0x87, 0x87}},
	{AXCHGQ, yxchg, Pw, [23]uint8{0x90, 0x90, 0x87, 0x87}},
	{AXCHGW, yxchg, Pe, [23]uint8{0x90, 0x90, 0x87, 0x87}},
	{AXLAT, ynone, Px, [23]uint8{0xd7}},
	{AXORB, yxorb, Pb, [23]uint8{0x34, 0x80, 06, 0x30, 0x32}},
	{AXORL, yxorl, Px, [23]uint8{0x83, 06, 0x35, 0x81, 06, 0x31, 0x33}},
	{AXORPD, yxm, Pe, [23]uint8{0x57}},
	{AXORPS, yxm, Pm, [23]uint8{0x57}},
	{AXORQ, yxorl, Pw, [23]uint8{0x83, 06, 0x35, 0x81, 06, 0x31, 0x33}},
	{AXORW, yxorl, Pe, [23]uint8{0x83, 06, 0x35, 0x81, 06, 0x31, 0x33}},
	{AFMOVB, yfmvx, Px, [23]uint8{0xdf, 04}},
	{AFMOVBP, yfmvp, Px, [23]uint8{0xdf, 06}},
	{AFMOVD, yfmvd, Px, [23]uint8{0xdd, 00, 0xdd, 02, 0xd9, 00, 0xdd, 02}},
	{AFMOVDP, yfmvdp, Px, [23]uint8{0xdd, 03, 0xdd, 03}},
	{AFMOVF, yfmvf, Px, [23]uint8{0xd9, 00, 0xd9, 02}},
	{AFMOVFP, yfmvp, Px, [23]uint8{0xd9, 03}},
	{AFMOVL, yfmvf, Px, [23]uint8{0xdb, 00, 0xdb, 02}},
	{AFMOVLP, yfmvp, Px, [23]uint8{0xdb, 03}},
	{AFMOVV, yfmvx, Px, [23]uint8{0xdf, 05}},
	{AFMOVVP, yfmvp, Px, [23]uint8{0xdf, 07}},
	{AFMOVW, yfmvf, Px, [23]uint8{0xdf, 00, 0xdf, 02}},
	{AFMOVWP, yfmvp, Px, [23]uint8{0xdf, 03}},
	{AFMOVX, yfmvx, Px, [23]uint8{0xdb, 05}},
	{AFMOVXP, yfmvp, Px, [23]uint8{0xdb, 07}},
	{AFCMOVCC, yfcmv, Px, [23]uint8{0xdb, 00}},
	{AFCMOVCS, yfcmv, Px, [23]uint8{0xda, 00}},
	{AFCMOVEQ, yfcmv, Px, [23]uint8{0xda, 01}},
	{AFCMOVHI, yfcmv, Px, [23]uint8{0xdb, 02}},
	{AFCMOVLS, yfcmv, Px, [23]uint8{0xda, 02}},
	{AFCMOVNE, yfcmv, Px, [23]uint8{0xdb, 01}},
	{AFCMOVNU, yfcmv, Px, [23]uint8{0xdb, 03}},
	{AFCMOVUN, yfcmv, Px, [23]uint8{0xda, 03}},
	{AFCOMB, nil, 0, [23]uint8{}},
	{AFCOMBP, nil, 0, [23]uint8{}},
	{AFCOMD, yfadd, Px, [23]uint8{0xdc, 02, 0xd8, 02, 0xdc, 02}},  /* botch */
	{AFCOMDP, yfadd, Px, [23]uint8{0xdc, 03, 0xd8, 03, 0xdc, 03}}, /* botch */
	{AFCOMDPP, ycompp, Px, [23]uint8{0xde, 03}},
	{AFCOMF, yfmvx, Px, [23]uint8{0xd8, 02}},
	{AFCOMFP, yfmvx, Px, [23]uint8{0xd8, 03}},
	{AFCOMI, yfmvx, Px, [23]uint8{0xdb, 06}},
	{AFCOMIP, yfmvx, Px, [23]uint8{0xdf, 06}},
	{AFCOML, yfmvx, Px, [23]uint8{0xda, 02}},
	{AFCOMLP, yfmvx, Px, [23]uint8{0xda, 03}},
	{AFCOMW, yfmvx, Px, [23]uint8{0xde, 02}},
	{AFCOMWP, yfmvx, Px, [23]uint8{0xde, 03}},
	{AFUCOM, ycompp, Px, [23]uint8{0xdd, 04}},
	{AFUCOMI, ycompp, Px, [23]uint8{0xdb, 05}},
	{AFUCOMIP, ycompp, Px, [23]uint8{0xdf, 05}},
	{AFUCOMP, ycompp, Px, [23]uint8{0xdd, 05}},
	{AFUCOMPP, ycompp, Px, [23]uint8{0xda, 13}},
	{AFADDDP, yfaddp, Px, [23]uint8{0xde, 00}},
	{AFADDW, yfmvx, Px, [23]uint8{0xde, 00}},
	{AFADDL, yfmvx, Px, [23]uint8{0xda, 00}},
	{AFADDF, yfmvx, Px, [23]uint8{0xd8, 00}},
	{AFADDD, yfadd, Px, [23]uint8{0xdc, 00, 0xd8, 00, 0xdc, 00}},
	{AFMULDP, yfaddp, Px, [23]uint8{0xde, 01}},
	{AFMULW, yfmvx, Px, [23]uint8{0xde, 01}},
	{AFMULL, yfmvx, Px, [23]uint8{0xda, 01}},
	{AFMULF, yfmvx, Px, [23]uint8{0xd8, 01}},
	{AFMULD, yfadd, Px, [23]uint8{0xdc, 01, 0xd8, 01, 0xdc, 01}},
	{AFSUBDP, yfaddp, Px, [23]uint8{0xde, 05}},
	{AFSUBW, yfmvx, Px, [23]uint8{0xde, 04}},
	{AFSUBL, yfmvx, Px, [23]uint8{0xda, 04}},
	{AFSUBF, yfmvx, Px, [23]uint8{0xd8, 04}},
	{AFSUBD, yfadd, Px, [23]uint8{0xdc, 04, 0xd8, 04, 0xdc, 05}},
	{AFSUBRDP, yfaddp, Px, [23]uint8{0xde, 04}},
	{AFSUBRW, yfmvx, Px, [23]uint8{0xde, 05}},
	{AFSUBRL, yfmvx, Px, [23]uint8{0xda, 05}},
	{AFSUBRF, yfmvx, Px, [23]uint8{0xd8, 05}},
	{AFSUBRD, yfadd, Px, [23]uint8{0xdc, 05, 0xd8, 05, 0xdc, 04}},
	{AFDIVDP, yfaddp, Px, [23]uint8{0xde, 07}},
	{AFDIVW, yfmvx, Px, [23]uint8{0xde, 06}},
	{AFDIVL, yfmvx, Px, [23]uint8{0xda, 06}},
	{AFDIVF, yfmvx, Px, [23]uint8{0xd8, 06}},
	{AFDIVD, yfadd, Px, [23]uint8{0xdc, 06, 0xd8, 06, 0xdc, 07}},
	{AFDIVRDP, yfaddp, Px, [23]uint8{0xde, 06}},
	{AFDIVRW, yfmvx, Px, [23]uint8{0xde, 07}},
	{AFDIVRL, yfmvx, Px, [23]uint8{0xda, 07}},
	{AFDIVRF, yfmvx, Px, [23]uint8{0xd8, 07}},
	{AFDIVRD, yfadd, Px, [23]uint8{0xdc, 07, 0xd8, 07, 0xdc, 06}},
	{AFXCHD, yfxch, Px, [23]uint8{0xd9, 01, 0xd9, 01}},
	{AFFREE, nil, 0, [23]uint8{}},
	{AFLDCW, ystcw, Px, [23]uint8{0xd9, 05, 0xd9, 05}},
	{AFLDENV, ystcw, Px, [23]uint8{0xd9, 04, 0xd9, 04}},
	{AFRSTOR, ysvrs, Px, [23]uint8{0xdd, 04, 0xdd, 04}},
	{AFSAVE, ysvrs, Px, [23]uint8{0xdd, 06, 0xdd, 06}},
	{AFSTCW, ystcw, Px, [23]uint8{0xd9, 07, 0xd9, 07}},
	{AFSTENV, ystcw, Px, [23]uint8{0xd9, 06, 0xd9, 06}},
	{AFSTSW, ystsw, Px, [23]uint8{0xdd, 07, 0xdf, 0xe0}},
	{AF2XM1, ynone, Px, [23]uint8{0xd9, 0xf0}},
	{AFABS, ynone, Px, [23]uint8{0xd9, 0xe1}},
	{AFCHS, ynone, Px, [23]uint8{0xd9, 0xe0}},
	{AFCLEX, ynone, Px, [23]uint8{0xdb, 0xe2}},
	{AFCOS, ynone, Px, [23]uint8{0xd9, 0xff}},
	{AFDECSTP, ynone, Px, [23]uint8{0xd9, 0xf6}},
	{AFINCSTP, ynone, Px, [23]uint8{0xd9, 0xf7}},
	{AFINIT, ynone, Px, [23]uint8{0xdb, 0xe3}},
	{AFLD1, ynone, Px, [23]uint8{0xd9, 0xe8}},
	{AFLDL2E, ynone, Px, [23]uint8{0xd9, 0xea}},
	{AFLDL2T, ynone, Px, [23]uint8{0xd9, 0xe9}},
	{AFLDLG2, ynone, Px, [23]uint8{0xd9, 0xec}},
	{AFLDLN2, ynone, Px, [23]uint8{0xd9, 0xed}},
	{AFLDPI, ynone, Px, [23]uint8{0xd9, 0xeb}},
	{AFLDZ, ynone, Px, [23]uint8{0xd9, 0xee}},
	{AFNOP, ynone, Px, [23]uint8{0xd9, 0xd0}},
	{AFPATAN, ynone, Px, [23]uint8{0xd9, 0xf3}},
	{AFPREM, ynone, Px, [23]uint8{0xd9, 0xf8}},
	{AFPREM1, ynone, Px, [23]uint8{0xd9, 0xf5}},
	{AFPTAN, ynone, Px, [23]uint8{0xd9, 0xf2}},
	{AFRNDINT, ynone, Px, [23]uint8{0xd9, 0xfc}},
	{AFSCALE, ynone, Px, [23]uint8{0xd9, 0xfd}},
	{AFSIN, ynone, Px, [23]uint8{0xd9, 0xfe}},
	{AFSINCOS, ynone, Px, [23]uint8{0xd9, 0xfb}},
	{AFSQRT, ynone, Px, [23]uint8{0xd9, 0xfa}},
	{AFTST, ynone, Px, [23]uint8{0xd9, 0xe4}},
	{AFXAM, ynone, Px, [23]uint8{0xd9, 0xe5}},
	{AFXTRACT, ynone, Px, [23]uint8{0xd9, 0xf4}},
	{AFYL2X, ynone, Px, [23]uint8{0xd9, 0xf1}},
	{AFYL2XP1, ynone, Px, [23]uint8{0xd9, 0xf9}},
	{ACMPXCHGB, yrb_mb, Pb, [23]uint8{0x0f, 0xb0}},
	{ACMPXCHGL, yrl_ml, Px, [23]uint8{0x0f, 0xb1}},
	{ACMPXCHGW, yrl_ml, Pe, [23]uint8{0x0f, 0xb1}},
	{ACMPXCHGQ, yrl_ml, Pw, [23]uint8{0x0f, 0xb1}},
	{ACMPXCHG8B, yscond, Pm, [23]uint8{0xc7, 01}},
	{AINVD, ynone, Pm, [23]uint8{0x08}},
	{AINVLPG, ymbs, Pm, [23]uint8{0x01, 07}},
	{ALFENCE, ynone, Pm, [23]uint8{0xae, 0xe8}},
	{AMFENCE, ynone, Pm, [23]uint8{0xae, 0xf0}},
	{AMOVNTIL, yrl_ml, Pm, [23]uint8{0xc3}},
	{AMOVNTIQ, yrl_ml, Pw, [23]uint8{0x0f, 0xc3}},
	{ARDMSR, ynone, Pm, [23]uint8{0x32}},
	{ARDPMC, ynone, Pm, [23]uint8{0x33}},
	{ARDTSC, ynone, Pm, [23]uint8{0x31}},
	{ARSM, ynone, Pm, [23]uint8{0xaa}},
	{ASFENCE, ynone, Pm, [23]uint8{0xae, 0xf8}},
	{ASYSRET, ynone, Pm, [23]uint8{0x07}},
	{AWBINVD, ynone, Pm, [23]uint8{0x09}},
	{AWRMSR, ynone, Pm, [23]uint8{0x30}},
	{AXADDB, yrb_mb, Pb, [23]uint8{0x0f, 0xc0}},
	{AXADDL, yrl_ml, Px, [23]uint8{0x0f, 0xc1}},
	{AXADDQ, yrl_ml, Pw, [23]uint8{0x0f, 0xc1}},
	{AXADDW, yrl_ml, Pe, [23]uint8{0x0f, 0xc1}},
	{ACRC32B, ycrc32l, Px, [23]uint8{0xf2, 0x0f, 0x38, 0xf0, 0}},
	{ACRC32Q, ycrc32l, Pw, [23]uint8{0xf2, 0x0f, 0x38, 0xf1, 0}},
	{APREFETCHT0, yprefetch, Pm, [23]uint8{0x18, 01}},
	{APREFETCHT1, yprefetch, Pm, [23]uint8{0x18, 02}},
	{APREFETCHT2, yprefetch, Pm, [23]uint8{0x18, 03}},
	{APREFETCHNTA, yprefetch, Pm, [23]uint8{0x18, 00}},
	{AMOVQL, yrl_ml, Px, [23]uint8{0x89}},
	{objabi.AUNDEF, ynone, Px, [23]uint8{0x0f, 0x0b}},
	{AAESENC, yaes, Pq, [23]uint8{0x38, 0xdc, 0}},
	{AAESENCLAST, yaes, Pq, [23]uint8{0x38, 0xdd, 0}},
	{AAESDEC, yaes, Pq, [23]uint8{0x38, 0xde, 0}},
	{AAESDECLAST, yaes, Pq, [23]uint8{0x38, 0xdf, 0}},
	{AAESIMC, yaes, Pq, [23]uint8{0x38, 0xdb, 0}},
	{AAESKEYGENASSIST, yaes2, Pq, [23]uint8{0x3a, 0xdf, 0}},
	{APSHUFD, yxshuf, Pq, [23]uint8{0x70, 0}},
	{APCLMULQDQ, yxshuf, Pq, [23]uint8{0x3a, 0x44, 0}},
	{objabi.AUSEFIELD, ynop, Px, [23]uint8{0, 0}},
	{objabi.ATYPE, nil, 0, [23]uint8{}},
	{objabi.AFUNCDATA, yfuncdata, Px, [23]uint8{0, 0}},
	{objabi.APCDATA, ypcdata, Px, [23]uint8{0, 0}},
	{objabi.ACHECKNIL, nil, 0, [23]uint8{}},
	{objabi.AVARDEF, nil, 0, [23]uint8{}},
	{objabi.AVARKILL, nil, 0, [23]uint8{}},
	{objabi.ADUFFCOPY, yduff, Px, [23]uint8{0xe8}},
	{objabi.ADUFFZERO, yduff, Px, [23]uint8{0xe8}},
	{objabi.AEND, nil, 0, [23]uint8{}},
	{0, nil, 0, [23]uint8{}},
}

// 快速索引表
// 因为原始的指令比较稀疏(涉及多个指令集), 取低bit可以更紧凑
var opindex [(ALAST + 1) & objabi.AMask]*Optab

// isextern reports whether s describes an external symbol that must avoid pc-relative addressing.
// This happens on systems like Solaris that call .so functions instead of system calls.
// It does not seem to be necessary for any other systems. This is probably working
// around a Solaris-specific bug that should be fixed differently, but we don't know
// what that bug is. And this does fix it.
func isextern(s *obj.LSym) bool {
	// All the Solaris dynamic imports from libc.so begin with "libc_".
	return strings.HasPrefix(s.Name, "libc_")
}

// single-instruction no-ops of various lengths.
// constructed by hand and disassembled with gdb to verify.
// see http://www.agner.org/optimize/optimizing_assembly.pdf for discussion.
var nop = [][16]uint8{
	{0x90},
	{0x66, 0x90},
	{0x0F, 0x1F, 0x00},
	{0x0F, 0x1F, 0x40, 0x00},
	{0x0F, 0x1F, 0x44, 0x00, 0x00},
	{0x66, 0x0F, 0x1F, 0x44, 0x00, 0x00},
	{0x0F, 0x1F, 0x80, 0x00, 0x00, 0x00, 0x00},
	{0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},
	{0x66, 0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},
}

// Native Client rejects the repeated 0x66 prefix.
// {0x66, 0x66, 0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},
func fillnop(p []byte, n int) {
	var m int

	for n > 0 {
		m = n
		if m > len(nop) {
			m = len(nop)
		}
		copy(p[:m], nop[m-1][:m])
		p = p[m:]
		n -= m
	}
}

func spadjop(ctxt *obj.Link, p *obj.Prog, l, q objabi.As) objabi.As {
	if p.Mode != 64 || ctxt.Arch.Ptrsize == 4 {
		return l
	}
	return q
}

func span6(ctxt *obj.Link, s *obj.LSym) {
	ctxt.Cursym = s

	if s.P != nil {
		return
	}

	if ycover[0] == 0 {
		instinit()
	}

	var v int32
	for p := ctxt.Cursym.Text; p != nil; p = p.Link {
		if p.To.Type == obj.TYPE_BRANCH {
			if p.Pcond == nil {
				p.Pcond = p
			}
		}
		if p.As == AADJSP {
			p.To.Type = obj.TYPE_REG
			p.To.Reg = REG_SP
			v = int32(-p.From.Offset)
			p.From.Offset = int64(v)
			p.As = spadjop(ctxt, p, AADDL, AADDQ)
			if v < 0 {
				p.As = spadjop(ctxt, p, ASUBL, ASUBQ)
				v = -v
				p.From.Offset = int64(v)
			}

			if v == 0 {
				p.As = objabi.ANOP
			}
		}
	}

	var q *obj.Prog
	for p := s.Text; p != nil; p = p.Link {
		p.Back = 2 // use short branches first time through
		q = p.Pcond
		if q != nil && (q.Back&2 != 0) {
			p.Back |= 1 // backward jump
			q.Back |= 4 // loop head
		}

		if p.As == AADJSP {
			p.To.Type = obj.TYPE_REG
			p.To.Reg = REG_SP
			v = int32(-p.From.Offset)
			p.From.Offset = int64(v)
			p.As = spadjop(ctxt, p, AADDL, AADDQ)
			if v < 0 {
				p.As = spadjop(ctxt, p, ASUBL, ASUBQ)
				v = -v
				p.From.Offset = int64(v)
			}

			if v == 0 {
				p.As = objabi.ANOP
			}
		}
	}

	n := 0
	var bp []byte
	var c int32
	var i int
	var loop int32
	var m int
	var p *obj.Prog
	for {
		loop = 0
		for i = 0; i < len(s.R); i++ {
			s.R[i] = obj.Reloc{}
		}
		s.R = s.R[:0]
		s.P = s.P[:0]
		c = 0
		for p = s.Text; p != nil; p = p.Link {
			if (p.Back&4 != 0) && c&(LoopAlign-1) != 0 {
				// pad with NOPs
				v = -c & (LoopAlign - 1)

				if v <= MaxLoopPad {
					obj_Symgrow(s, int64(c)+int64(v))
					fillnop(s.P[c:], int(v))
					c += v
				}
			}

			p.Pc = int64(c)

			// process forward jumps to p
			for q = p.Rel; q != nil; q = q.Forwd {
				v = int32(p.Pc - (q.Pc + int64(q.Mark)))
				if q.Back&2 != 0 { // short
					if v > 127 {
						loop++
						q.Back ^= 2
					}

					if q.As == AJCXZL {
						s.P[q.Pc+2] = byte(v)
					} else {
						s.P[q.Pc+1] = byte(v)
					}
				} else {
					bp = s.P[q.Pc+int64(q.Mark)-4:]
					bp[0] = byte(v)
					bp = bp[1:]
					bp[0] = byte(v >> 8)
					bp = bp[1:]
					bp[0] = byte(v >> 16)
					bp = bp[1:]
					bp[0] = byte(v >> 24)
				}
			}

			p.Rel = nil

			p.Pc = int64(c)
			asmins(ctxt, p)
			m = -cap(ctxt.Andptr) + cap(ctxt.And[:])
			if int(p.Isize) != m {
				p.Isize = uint8(m)
				loop++
			}

			obj_Symgrow(s, p.Pc+int64(m))
			copy(s.P[p.Pc:][:m], ctxt.And[:m])
			p.Mark = uint16(m)
			c += int32(m)
		}

		n++
		if n > 20 {
			ctxt.Diag("span must be looping")
			log.Fatalf("loop")
		}
		if loop == 0 {
			break
		}
	}

	c += -c & (FuncAlign - 1)
	s.Size = int64(c)

	if false { /* debug['a'] > 1 */
		fmt.Printf("span1 %s %d (%d tries)\n %.6x", s.Name, s.Size, n, 0)
		var i int
		for i = 0; i < len(s.P); i++ {
			fmt.Printf(" %.2x", s.P[i])
			if i%16 == 15 {
				fmt.Printf("\n  %.6x", uint(i+1))
			}
		}

		if i%16 != 0 {
			fmt.Printf("\n")
		}

		for i := 0; i < len(s.R); i++ {
			r := &s.R[i]
			fmt.Printf(" rel %#.4x/%d %s%+d\n", uint32(r.Off), r.Siz, r.Sym.Name, r.Add)
		}
	}
}

func instinit() {
	var c objabi.As

	for i := 1; optab[i].as != 0; i++ {
		c = optab[i].as
		if opindex[c&objabi.AMask] != nil {
			log.Fatalf("phase error in optab: %d (%v)", i, c)
		}
		opindex[c&objabi.AMask] = &optab[i]
	}

	for i := 0; i < Ymax; i++ {
		ycover[i*Ymax+i] = 1
	}

	ycover[Yi0*Ymax+Yi8] = 1
	ycover[Yi1*Ymax+Yi8] = 1
	ycover[Yu7*Ymax+Yi8] = 1

	ycover[Yi0*Ymax+Yu7] = 1
	ycover[Yi1*Ymax+Yu7] = 1

	ycover[Yi0*Ymax+Yu8] = 1
	ycover[Yi1*Ymax+Yu8] = 1
	ycover[Yu7*Ymax+Yu8] = 1

	ycover[Yi0*Ymax+Ys32] = 1
	ycover[Yi1*Ymax+Ys32] = 1
	ycover[Yu7*Ymax+Ys32] = 1
	ycover[Yu8*Ymax+Ys32] = 1
	ycover[Yi8*Ymax+Ys32] = 1

	ycover[Yi0*Ymax+Yi32] = 1
	ycover[Yi1*Ymax+Yi32] = 1
	ycover[Yu7*Ymax+Yi32] = 1
	ycover[Yu8*Ymax+Yi32] = 1
	ycover[Yi8*Ymax+Yi32] = 1
	ycover[Ys32*Ymax+Yi32] = 1

	ycover[Yi0*Ymax+Yi64] = 1
	ycover[Yi1*Ymax+Yi64] = 1
	ycover[Yu7*Ymax+Yi64] = 1
	ycover[Yu8*Ymax+Yi64] = 1
	ycover[Yi8*Ymax+Yi64] = 1
	ycover[Ys32*Ymax+Yi64] = 1
	ycover[Yi32*Ymax+Yi64] = 1

	ycover[Yal*Ymax+Yrb] = 1
	ycover[Ycl*Ymax+Yrb] = 1
	ycover[Yax*Ymax+Yrb] = 1
	ycover[Ycx*Ymax+Yrb] = 1
	ycover[Yrx*Ymax+Yrb] = 1
	ycover[Yrl*Ymax+Yrb] = 1 // but not Yrl32

	ycover[Ycl*Ymax+Ycx] = 1

	ycover[Yax*Ymax+Yrx] = 1
	ycover[Ycx*Ymax+Yrx] = 1

	ycover[Yax*Ymax+Yrl] = 1
	ycover[Ycx*Ymax+Yrl] = 1
	ycover[Yrx*Ymax+Yrl] = 1
	ycover[Yrl32*Ymax+Yrl] = 1

	ycover[Yf0*Ymax+Yrf] = 1

	ycover[Yal*Ymax+Ymb] = 1
	ycover[Ycl*Ymax+Ymb] = 1
	ycover[Yax*Ymax+Ymb] = 1
	ycover[Ycx*Ymax+Ymb] = 1
	ycover[Yrx*Ymax+Ymb] = 1
	ycover[Yrb*Ymax+Ymb] = 1
	ycover[Yrl*Ymax+Ymb] = 1 // but not Yrl32
	ycover[Ym*Ymax+Ymb] = 1

	ycover[Yax*Ymax+Yml] = 1
	ycover[Ycx*Ymax+Yml] = 1
	ycover[Yrx*Ymax+Yml] = 1
	ycover[Yrl*Ymax+Yml] = 1
	ycover[Yrl32*Ymax+Yml] = 1
	ycover[Ym*Ymax+Yml] = 1

	ycover[Yax*Ymax+Ymm] = 1
	ycover[Ycx*Ymax+Ymm] = 1
	ycover[Yrx*Ymax+Ymm] = 1
	ycover[Yrl*Ymax+Ymm] = 1
	ycover[Yrl32*Ymax+Ymm] = 1
	ycover[Ym*Ymax+Ymm] = 1
	ycover[Ymr*Ymax+Ymm] = 1

	ycover[Ym*Ymax+Yxm] = 1
	ycover[Yxr*Ymax+Yxm] = 1

	for i := 0; i < int(MAXREG); i++ {
		reg[i] = -1
		if i >= int(REG_AL) && i <= int(REG_R15B) {
			reg[i] = (i - int(REG_AL)) & 7
			if i >= int(REG_SPB) && i <= int(REG_DIB) {
				regrex[i] = 0x40
			}
			if i >= int(REG_R8B) && i <= int(REG_R15B) {
				regrex[i] = Rxr | Rxx | Rxb
			}
		}

		if i >= int(REG_AH) && i <= int(REG_BH) {
			reg[i] = 4 + ((i - int(REG_AH)) & 7)
		}
		if i >= int(REG_AX) && i <= int(REG_R15) {
			reg[i] = (i - int(REG_AX)) & 7
			if i >= int(REG_R8) {
				regrex[i] = Rxr | Rxx | Rxb
			}
		}

		if i >= int(REG_F0) && i <= int(REG_F0)+7 {
			reg[i] = (i - int(REG_F0)) & 7
		}
		if i >= int(REG_M0) && i <= int(REG_M0)+7 {
			reg[i] = (i - int(REG_M0)) & 7
		}
		if i >= int(REG_X0) && i <= int(REG_X0)+15 {
			reg[i] = (i - int(REG_X0)) & 7
			if i >= int(REG_X0)+8 {
				regrex[i] = Rxr | Rxx | Rxb
			}
		}

		if i >= int(REG_CR)+8 && i <= int(REG_CR)+15 {
			regrex[i] = Rxr
		}
	}
}

func prefixof(ctxt *obj.Link, p *obj.Prog, a *obj.Addr) int {
	if a.Reg < REG_CS && a.Index < int32(REG_CS) { // fast path
		return 0
	}
	if a.Type == obj.TYPE_MEM && a.Name == obj.NAME_NONE {
		switch a.Reg {
		case REG_CS:
			return 0x2e

		case REG_DS:
			return 0x3e

		case REG_ES:
			return 0x26

		case REG_FS:
			return 0x64

		case REG_GS:
			return 0x65

		case REG_TLS:
			// NOTE: Systems listed here should be only systems that
			// support direct TLS references like 8(TLS) implemented as
			// direct references from FS or GS. Systems that require
			// the initial-exec model, where you load the TLS base into
			// a register and then index from that register, do not reach
			// this code and should not be listed.
			if p.Mode == 32 {
				switch ctxt.Headtype {
				default:
					log.Fatalf("unknown TLS base register for %s", ctxt.Headtype)

				case objabi.Hdarwin:
					return 0x65 // GS
				}
			}

			switch ctxt.Headtype {
			default:
				log.Fatalf("unknown TLS base register for %s", ctxt.Headtype)

			case objabi.Hlinux:
				if ctxt.Flag_shared != 0 {
					log.Fatalf("unknown TLS base register for linux with -shared")
				} else {
					return 0x64 // FS
				}

			case objabi.Hdarwin:
				return 0x65 // GS
			}
		}
	}

	if p.Mode == 32 {
		return 0
	}

	switch objabi.RBaseType(a.Index) {
	case REG_CS:
		return 0x2e

	case REG_DS:
		return 0x3e

	case REG_ES:
		return 0x26

	case REG_TLS:
		if ctxt.Flag_shared != 0 {
			// When building for inclusion into a shared library, an instruction of the form
			//     MOV 0(CX)(TLS*1), AX
			// becomes
			//     mov %fs:(%rcx), %rax
			// which assumes that the correct TLS offset has been loaded into %rcx (today
			// there is only one TLS variable -- g -- so this is OK). When not building for
			// a shared library the instruction does not require a prefix.
			if a.Offset != 0 {
				log.Fatalf("cannot handle non-0 offsets to TLS")
			}
			return 0x64
		}

	case REG_FS:
		return 0x64

	case REG_GS:
		return 0x65
	}

	return 0
}

func oclass(ctxt *obj.Link, p *obj.Prog, a *obj.Addr) int {
	switch a.Type {
	case obj.TYPE_NONE:
		return Ynone

	case obj.TYPE_BRANCH:
		return Ybr

	case obj.TYPE_INDIR:
		if a.Name != obj.NAME_NONE && a.Reg == REG_NONE && a.Index == REG_NONE && a.Scale == 0 {
			return Yindir
		}
		return Yxxx

	case obj.TYPE_MEM:
		return Ym

	case obj.TYPE_ADDR:
		switch a.Name {
		case obj.NAME_EXTERN,
			obj.NAME_GOTREF,
			obj.NAME_STATIC:
			if a.Sym != nil && isextern(a.Sym) || p.Mode == 32 {
				return Yi32
			}
			return Yiauto // use pc-relative addressing

		case obj.NAME_AUTO,
			obj.NAME_PARAM:
			return Yiauto
		}

		// TODO(chai2010): DUFFZERO/DUFFCOPY encoding forgot to set a->index
		// and got Yi32 in an earlier version of this code.
		// Keep doing that until we fix yduff etc.
		if a.Sym != nil && strings.HasPrefix(a.Sym.Name, "runtime.duff") {
			return Yi32
		}

		if a.Sym != nil || a.Name != obj.NAME_NONE {
			ctxt.Diag("unexpected addr: %v", a.Dconv(p))
		}
		fallthrough

		// fall through

	case obj.TYPE_CONST:
		if a.Sym != nil {
			ctxt.Diag("TYPE_CONST with symbol: %v", a.Dconv(p))
		}

		v := a.Offset
		if p.Mode == 32 {
			v = int64(int32(v))
		}
		if v == 0 {
			return Yi0
		}
		if v == 1 {
			return Yi1
		}
		if v >= 0 && v <= 127 {
			return Yu7
		}
		if v >= 0 && v <= 255 {
			return Yu8
		}
		if v >= -128 && v <= 127 {
			return Yi8
		}
		if p.Mode == 32 {
			return Yi32
		}
		l := int32(v)
		if int64(l) == v {
			return Ys32 /* can sign extend */
		}
		if v>>32 == 0 {
			return Yi32 /* unsigned */
		}
		return Yi64

	case obj.TYPE_TEXTSIZE:
		return Ytextsize
	}

	if a.Type != obj.TYPE_REG {
		ctxt.Diag("unexpected addr1: type=%d %v", a.Type, a.Dconv(p))
		return Yxxx
	}

	switch a.Reg {
	case REG_AL:
		return Yal

	case REG_AX:
		return Yax

		/*
			case REG_SPB:
		*/
	case REG_BPB,
		REG_SIB,
		REG_DIB,
		REG_R8B,
		REG_R9B,
		REG_R10B,
		REG_R11B,
		REG_R12B,
		REG_R13B,
		REG_R14B,
		REG_R15B:
		if ctxt.Asmode != 64 {
			return Yxxx
		}
		fallthrough

	case REG_DL,
		REG_BL,
		REG_AH,
		REG_CH,
		REG_DH,
		REG_BH:
		return Yrb

	case REG_CL:
		return Ycl

	case REG_CX:
		return Ycx

	case REG_DX, REG_BX:
		return Yrx

	case REG_R8, /* not really Yrl */
		REG_R9,
		REG_R10,
		REG_R11,
		REG_R12,
		REG_R13,
		REG_R14,
		REG_R15:
		if ctxt.Asmode != 64 {
			return Yxxx
		}
		fallthrough

	case REG_SP, REG_BP, REG_SI, REG_DI:
		if p.Mode == 32 {
			return Yrl32
		}
		return Yrl

	case REG_F0 + 0:
		return Yf0

	case REG_F0 + 1,
		REG_F0 + 2,
		REG_F0 + 3,
		REG_F0 + 4,
		REG_F0 + 5,
		REG_F0 + 6,
		REG_F0 + 7:
		return Yrf

	case REG_M0 + 0,
		REG_M0 + 1,
		REG_M0 + 2,
		REG_M0 + 3,
		REG_M0 + 4,
		REG_M0 + 5,
		REG_M0 + 6,
		REG_M0 + 7:
		return Ymr

	case REG_X0 + 0,
		REG_X0 + 1,
		REG_X0 + 2,
		REG_X0 + 3,
		REG_X0 + 4,
		REG_X0 + 5,
		REG_X0 + 6,
		REG_X0 + 7,
		REG_X0 + 8,
		REG_X0 + 9,
		REG_X0 + 10,
		REG_X0 + 11,
		REG_X0 + 12,
		REG_X0 + 13,
		REG_X0 + 14,
		REG_X0 + 15:
		return Yxr

	case REG_CS:
		return Ycs
	case REG_SS:
		return Yss
	case REG_DS:
		return Yds
	case REG_ES:
		return Yes
	case REG_FS:
		return Yfs
	case REG_GS:
		return Ygs
	case REG_TLS:
		return Ytls

	case REG_GDTR:
		return Ygdtr
	case REG_IDTR:
		return Yidtr
	case REG_LDTR:
		return Yldtr
	case REG_MSW:
		return Ymsw
	case REG_TASK:
		return Ytask

	case REG_CR + 0:
		return Ycr0
	case REG_CR + 1:
		return Ycr1
	case REG_CR + 2:
		return Ycr2
	case REG_CR + 3:
		return Ycr3
	case REG_CR + 4:
		return Ycr4
	case REG_CR + 5:
		return Ycr5
	case REG_CR + 6:
		return Ycr6
	case REG_CR + 7:
		return Ycr7
	case REG_CR + 8:
		return Ycr8

	case REG_DR + 0:
		return Ydr0
	case REG_DR + 1:
		return Ydr1
	case REG_DR + 2:
		return Ydr2
	case REG_DR + 3:
		return Ydr3
	case REG_DR + 4:
		return Ydr4
	case REG_DR + 5:
		return Ydr5
	case REG_DR + 6:
		return Ydr6
	case REG_DR + 7:
		return Ydr7

	case REG_TR + 0:
		return Ytr0
	case REG_TR + 1:
		return Ytr1
	case REG_TR + 2:
		return Ytr2
	case REG_TR + 3:
		return Ytr3
	case REG_TR + 4:
		return Ytr4
	case REG_TR + 5:
		return Ytr5
	case REG_TR + 6:
		return Ytr6
	case REG_TR + 7:
		return Ytr7
	}

	return Yxxx
}

func asmidx(ctxt *obj.Link, scale int, index int, base int) {
	var i int

	switch objabi.RBaseType(index) {
	default:
		goto bad

	case REG_NONE:
		i = 4 << 3
		goto bas

	case REG_R8,
		REG_R9,
		REG_R10,
		REG_R11,
		REG_R12,
		REG_R13,
		REG_R14,
		REG_R15:
		if ctxt.Asmode != 64 {
			goto bad
		}
		fallthrough

	case REG_AX,
		REG_CX,
		REG_DX,
		REG_BX,
		REG_BP,
		REG_SI,
		REG_DI:
		i = reg[index] << 3
	}

	switch scale {
	default:
		goto bad

	case 1:
		break

	case 2:
		i |= 1 << 6

	case 4:
		i |= 2 << 6

	case 8:
		i |= 3 << 6
	}

bas:
	switch objabi.RBaseType(base) {
	default:
		goto bad

	case REG_NONE: /* must be mod=00 */
		i |= 5

	case REG_R8,
		REG_R9,
		REG_R10,
		REG_R11,
		REG_R12,
		REG_R13,
		REG_R14,
		REG_R15:
		if ctxt.Asmode != 64 {
			goto bad
		}
		fallthrough

	case REG_AX,
		REG_CX,
		REG_DX,
		REG_BX,
		REG_SP,
		REG_BP,
		REG_SI,
		REG_DI:
		i |= reg[base]
	}

	ctxt.Andptr[0] = byte(i)
	ctxt.Andptr = ctxt.Andptr[1:]
	return

bad:
	ctxt.Diag("asmidx: bad address %d/%d/%d", scale, index, base)
	ctxt.Andptr[0] = 0
	ctxt.Andptr = ctxt.Andptr[1:]
	return
}

func put4(ctxt *obj.Link, v int32) {
	ctxt.Andptr[0] = byte(v)
	ctxt.Andptr[1] = byte(v >> 8)
	ctxt.Andptr[2] = byte(v >> 16)
	ctxt.Andptr[3] = byte(v >> 24)
	ctxt.Andptr = ctxt.Andptr[4:]
}

func relput4(ctxt *obj.Link, p *obj.Prog, a *obj.Addr) {
	var rel obj.Reloc

	v := vaddr(ctxt, p, a, &rel)
	if rel.Siz != 0 {
		if rel.Siz != 4 {
			ctxt.Diag("bad reloc")
		}
		r := obj_Addrel(ctxt.Cursym)
		*r = rel
		r.Off = int32(p.Pc + int64(-cap(ctxt.Andptr)+cap(ctxt.And[:])))
	}

	put4(ctxt, int32(v))
}

func put8(ctxt *obj.Link, v int64) {
	ctxt.Andptr[0] = byte(v)
	ctxt.Andptr[1] = byte(v >> 8)
	ctxt.Andptr[2] = byte(v >> 16)
	ctxt.Andptr[3] = byte(v >> 24)
	ctxt.Andptr[4] = byte(v >> 32)
	ctxt.Andptr[5] = byte(v >> 40)
	ctxt.Andptr[6] = byte(v >> 48)
	ctxt.Andptr[7] = byte(v >> 56)
	ctxt.Andptr = ctxt.Andptr[8:]
}

/*
static void
relput8(Prog *p, Addr *a)

	{
		vlong v;
		Reloc rel, *r;

		v = vaddr(ctxt, p, a, &rel);
		if(rel.siz != 0) {
			r = addrel(ctxt->cursym);
			*r = rel;
			r->siz = 8;
			r->off = p->pc + ctxt->andptr - ctxt->and;
		}
		put8(ctxt, v);
	}
*/
func vaddr(ctxt *obj.Link, p *obj.Prog, a *obj.Addr, r *obj.Reloc) int64 {
	if r != nil {
		*r = obj.Reloc{}
	}

	switch a.Name {
	case obj.NAME_STATIC,
		obj.NAME_GOTREF,
		obj.NAME_EXTERN:
		s := a.Sym
		if r == nil {
			ctxt.Diag("need reloc for %v", a.Dconv(p))
			log.Fatalf("reloc")
		}

		if a.Name == obj.NAME_GOTREF {
			r.Siz = 4
			r.Type = obj.R_GOTPCREL
		} else if isextern(s) || p.Mode != 64 {
			r.Siz = 4
			r.Type = obj.R_ADDR
		} else {
			r.Siz = 4
			r.Type = obj.R_PCREL
		}

		r.Off = -1 // caller must fill in
		r.Sym = s
		r.Add = a.Offset

		return 0
	}

	if (a.Type == obj.TYPE_MEM || a.Type == obj.TYPE_ADDR) && a.Reg == REG_TLS {
		if r == nil {
			ctxt.Diag("need reloc for %v", a.Dconv(p))
			log.Fatalf("reloc")
		}

		r.Type = obj.R_TLS_LE
		r.Siz = 4
		r.Off = -1 // caller must fill in
		r.Add = a.Offset
		return 0
	}

	return a.Offset
}

func asmandsz(ctxt *obj.Link, p *obj.Prog, a *obj.Addr, r int, rex int, m64 int) {
	var base int
	var rel obj.Reloc

	rex &= 0x40 | Rxr
	v := int32(a.Offset)
	rel.Siz = 0

	switch a.Type {
	case obj.TYPE_ADDR:
		if a.Name == obj.NAME_NONE {
			ctxt.Diag("unexpected TYPE_ADDR with NAME_NONE")
		}
		if a.Index == int32(REG_TLS) {
			ctxt.Diag("unexpected TYPE_ADDR with index==REG_TLS")
		}
		goto bad

	case obj.TYPE_REG:
		if a.Reg < REG_AL || REG_X0+15 < a.Reg {
			goto bad
		}
		if v != 0 {
			goto bad
		}
		ctxt.Andptr[0] = byte(3<<6 | reg[a.Reg]<<0 | r<<3)
		ctxt.Andptr = ctxt.Andptr[1:]
		ctxt.Rexflag |= regrex[a.Reg]&(0x40|Rxb) | rex
		return
	}

	if a.Type != obj.TYPE_MEM {
		goto bad
	}

	if a.Index != REG_NONE && a.Index != int32(REG_TLS) {
		base := int(a.Reg)
		switch a.Name {
		case obj.NAME_EXTERN,
			obj.NAME_GOTREF,
			obj.NAME_STATIC:
			if !isextern(a.Sym) && p.Mode == 64 {
				goto bad
			}
			base = REG_NONE
			v = int32(vaddr(ctxt, p, a, &rel))

		case obj.NAME_AUTO,
			obj.NAME_PARAM:
			base = int(REG_SP)
		}

		ctxt.Rexflag |= regrex[int(a.Index)]&Rxx | regrex[base]&Rxb | rex
		if base == REG_NONE {
			ctxt.Andptr[0] = byte(0<<6 | 4<<0 | r<<3)
			ctxt.Andptr = ctxt.Andptr[1:]
			asmidx(ctxt, int(a.Scale), int(a.Index), base)
			goto putrelv
		}

		if v == 0 && rel.Siz == 0 && base != int(REG_BP) && base != int(REG_R13) {
			ctxt.Andptr[0] = byte(0<<6 | 4<<0 | r<<3)
			ctxt.Andptr = ctxt.Andptr[1:]
			asmidx(ctxt, int(a.Scale), int(a.Index), base)
			return
		}

		if v >= -128 && v < 128 && rel.Siz == 0 {
			ctxt.Andptr[0] = byte(1<<6 | 4<<0 | r<<3)
			ctxt.Andptr = ctxt.Andptr[1:]
			asmidx(ctxt, int(a.Scale), int(a.Index), base)
			ctxt.Andptr[0] = byte(v)
			ctxt.Andptr = ctxt.Andptr[1:]
			return
		}

		ctxt.Andptr[0] = byte(2<<6 | 4<<0 | r<<3)
		ctxt.Andptr = ctxt.Andptr[1:]
		asmidx(ctxt, int(a.Scale), int(a.Index), base)
		goto putrelv
	}

	base = int(a.Reg)
	switch a.Name {
	case obj.NAME_STATIC,
		obj.NAME_GOTREF,
		obj.NAME_EXTERN:
		if a.Sym == nil {
			ctxt.Diag("bad addr: %v", p)
		}
		base = REG_NONE
		v = int32(vaddr(ctxt, p, a, &rel))

	case obj.NAME_AUTO,
		obj.NAME_PARAM:
		base = int(REG_SP)
	}

	if base == int(REG_TLS) {
		v = int32(vaddr(ctxt, p, a, &rel))
	}

	ctxt.Rexflag |= regrex[base]&Rxb | rex
	if base == REG_NONE || (int(REG_CS) <= base && base <= int(REG_GS)) || base == int(REG_TLS) {
		if (a.Sym == nil || !isextern(a.Sym)) && base == REG_NONE && (a.Name == obj.NAME_STATIC || a.Name == obj.NAME_EXTERN || a.Name == obj.NAME_GOTREF) || p.Mode != 64 {
			if a.Name == obj.NAME_GOTREF && (a.Offset != 0 || a.Index != 0 || a.Scale != 0) {
				ctxt.Diag("%v has offset against gotref", p)
			}
			ctxt.Andptr[0] = byte(0<<6 | 5<<0 | r<<3)
			ctxt.Andptr = ctxt.Andptr[1:]
			goto putrelv
		}

		/* temporary */
		ctxt.Andptr[0] = byte(0<<6 | 4<<0 | r<<3)
		ctxt.Andptr = ctxt.Andptr[1:] /* sib present */
		ctxt.Andptr[0] = 0<<6 | 4<<3 | 5<<0
		ctxt.Andptr = ctxt.Andptr[1:] /* DS:d32 */
		goto putrelv
	}

	if base == int(REG_SP) || base == int(REG_R12) {
		if v == 0 {
			ctxt.Andptr[0] = byte(0<<6 | reg[base]<<0 | r<<3)
			ctxt.Andptr = ctxt.Andptr[1:]
			asmidx(ctxt, int(a.Scale), REG_NONE, base)
			return
		}

		if v >= -128 && v < 128 {
			ctxt.Andptr[0] = byte(1<<6 | reg[base]<<0 | r<<3)
			ctxt.Andptr = ctxt.Andptr[1:]
			asmidx(ctxt, int(a.Scale), REG_NONE, base)
			ctxt.Andptr[0] = byte(v)
			ctxt.Andptr = ctxt.Andptr[1:]
			return
		}

		ctxt.Andptr[0] = byte(2<<6 | reg[base]<<0 | r<<3)
		ctxt.Andptr = ctxt.Andptr[1:]
		asmidx(ctxt, int(a.Scale), REG_NONE, base)
		goto putrelv
	}

	if int(REG_AX) <= base && base <= int(REG_R15) {
		if a.Index == int32(REG_TLS) && ctxt.Flag_shared == 0 {
			rel = obj.Reloc{}
			rel.Type = obj.R_TLS_LE
			rel.Siz = 4
			rel.Sym = nil
			rel.Add = int64(v)
			v = 0
		}

		if v == 0 && rel.Siz == 0 && base != int(REG_BP) && base != int(REG_R13) {
			ctxt.Andptr[0] = byte(0<<6 | reg[base]<<0 | r<<3)
			ctxt.Andptr = ctxt.Andptr[1:]
			return
		}

		if v >= -128 && v < 128 && rel.Siz == 0 {
			ctxt.Andptr[0] = byte(1<<6 | reg[base]<<0 | r<<3)
			ctxt.Andptr[1] = byte(v)
			ctxt.Andptr = ctxt.Andptr[2:]
			return
		}

		ctxt.Andptr[0] = byte(2<<6 | reg[base]<<0 | r<<3)
		ctxt.Andptr = ctxt.Andptr[1:]
		goto putrelv
	}

	goto bad

putrelv:
	if rel.Siz != 0 {
		if rel.Siz != 4 {
			ctxt.Diag("bad rel")
			goto bad
		}

		r := obj_Addrel(ctxt.Cursym)
		*r = rel
		r.Off = int32(ctxt.Curp.Pc + int64(-cap(ctxt.Andptr)+cap(ctxt.And[:])))
	}

	put4(ctxt, v)
	return

bad:
	ctxt.Diag("asmand: bad address %v", a.Dconv(p))
	return
}

func asmand(ctxt *obj.Link, p *obj.Prog, a *obj.Addr, ra *obj.Addr) {
	asmandsz(ctxt, p, a, reg[ra.Reg], regrex[ra.Reg], 0)
}

func asmando(ctxt *obj.Link, p *obj.Prog, a *obj.Addr, o int) {
	asmandsz(ctxt, p, a, o, 0, 0)
}

func bytereg(a *obj.Addr, t *uint8) {
	if a.Type == obj.TYPE_REG && a.Index == REG_NONE && (REG_AX <= a.Reg && a.Reg <= REG_R15) {
		a.Reg += REG_AL - REG_AX
		*t = 0
	}
}

func unbytereg(a *obj.Addr, t *uint8) {
	if a.Type == obj.TYPE_REG && a.Index == REG_NONE && (REG_AL <= a.Reg && a.Reg <= REG_R15B) {
		a.Reg += REG_AX - REG_AL
		*t = 0
	}
}

const (
	E = 0xff
)

var ymovtab = []Movtab{
	/* push */
	{APUSHL, Ycs, Ynone, Ynone, 0, [4]uint8{0x0e, E, 0, 0}},
	{APUSHL, Yss, Ynone, Ynone, 0, [4]uint8{0x16, E, 0, 0}},
	{APUSHL, Yds, Ynone, Ynone, 0, [4]uint8{0x1e, E, 0, 0}},
	{APUSHL, Yes, Ynone, Ynone, 0, [4]uint8{0x06, E, 0, 0}},
	{APUSHL, Yfs, Ynone, Ynone, 0, [4]uint8{0x0f, 0xa0, E, 0}},
	{APUSHL, Ygs, Ynone, Ynone, 0, [4]uint8{0x0f, 0xa8, E, 0}},
	{APUSHQ, Yfs, Ynone, Ynone, 0, [4]uint8{0x0f, 0xa0, E, 0}},
	{APUSHQ, Ygs, Ynone, Ynone, 0, [4]uint8{0x0f, 0xa8, E, 0}},
	{APUSHW, Ycs, Ynone, Ynone, 0, [4]uint8{Pe, 0x0e, E, 0}},
	{APUSHW, Yss, Ynone, Ynone, 0, [4]uint8{Pe, 0x16, E, 0}},
	{APUSHW, Yds, Ynone, Ynone, 0, [4]uint8{Pe, 0x1e, E, 0}},
	{APUSHW, Yes, Ynone, Ynone, 0, [4]uint8{Pe, 0x06, E, 0}},
	{APUSHW, Yfs, Ynone, Ynone, 0, [4]uint8{Pe, 0x0f, 0xa0, E}},
	{APUSHW, Ygs, Ynone, Ynone, 0, [4]uint8{Pe, 0x0f, 0xa8, E}},

	/* pop */
	{APOPL, Ynone, Ynone, Yds, 0, [4]uint8{0x1f, E, 0, 0}},
	{APOPL, Ynone, Ynone, Yes, 0, [4]uint8{0x07, E, 0, 0}},
	{APOPL, Ynone, Ynone, Yss, 0, [4]uint8{0x17, E, 0, 0}},
	{APOPL, Ynone, Ynone, Yfs, 0, [4]uint8{0x0f, 0xa1, E, 0}},
	{APOPL, Ynone, Ynone, Ygs, 0, [4]uint8{0x0f, 0xa9, E, 0}},
	{APOPQ, Ynone, Ynone, Yfs, 0, [4]uint8{0x0f, 0xa1, E, 0}},
	{APOPQ, Ynone, Ynone, Ygs, 0, [4]uint8{0x0f, 0xa9, E, 0}},
	{APOPW, Ynone, Ynone, Yds, 0, [4]uint8{Pe, 0x1f, E, 0}},
	{APOPW, Ynone, Ynone, Yes, 0, [4]uint8{Pe, 0x07, E, 0}},
	{APOPW, Ynone, Ynone, Yss, 0, [4]uint8{Pe, 0x17, E, 0}},
	{APOPW, Ynone, Ynone, Yfs, 0, [4]uint8{Pe, 0x0f, 0xa1, E}},
	{APOPW, Ynone, Ynone, Ygs, 0, [4]uint8{Pe, 0x0f, 0xa9, E}},

	/* mov seg */
	{AMOVW, Yes, Ynone, Yml, 1, [4]uint8{0x8c, 0, 0, 0}},
	{AMOVW, Ycs, Ynone, Yml, 1, [4]uint8{0x8c, 1, 0, 0}},
	{AMOVW, Yss, Ynone, Yml, 1, [4]uint8{0x8c, 2, 0, 0}},
	{AMOVW, Yds, Ynone, Yml, 1, [4]uint8{0x8c, 3, 0, 0}},
	{AMOVW, Yfs, Ynone, Yml, 1, [4]uint8{0x8c, 4, 0, 0}},
	{AMOVW, Ygs, Ynone, Yml, 1, [4]uint8{0x8c, 5, 0, 0}},
	{AMOVW, Yml, Ynone, Yes, 2, [4]uint8{0x8e, 0, 0, 0}},
	{AMOVW, Yml, Ynone, Ycs, 2, [4]uint8{0x8e, 1, 0, 0}},
	{AMOVW, Yml, Ynone, Yss, 2, [4]uint8{0x8e, 2, 0, 0}},
	{AMOVW, Yml, Ynone, Yds, 2, [4]uint8{0x8e, 3, 0, 0}},
	{AMOVW, Yml, Ynone, Yfs, 2, [4]uint8{0x8e, 4, 0, 0}},
	{AMOVW, Yml, Ynone, Ygs, 2, [4]uint8{0x8e, 5, 0, 0}},

	/* mov cr */
	{AMOVL, Ycr0, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 0, 0}},
	{AMOVL, Ycr2, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 2, 0}},
	{AMOVL, Ycr3, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 3, 0}},
	{AMOVL, Ycr4, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 4, 0}},
	{AMOVL, Ycr8, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 8, 0}},
	{AMOVQ, Ycr0, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 0, 0}},
	{AMOVQ, Ycr2, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 2, 0}},
	{AMOVQ, Ycr3, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 3, 0}},
	{AMOVQ, Ycr4, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 4, 0}},
	{AMOVQ, Ycr8, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 8, 0}},
	{AMOVL, Yml, Ynone, Ycr0, 4, [4]uint8{0x0f, 0x22, 0, 0}},
	{AMOVL, Yml, Ynone, Ycr2, 4, [4]uint8{0x0f, 0x22, 2, 0}},
	{AMOVL, Yml, Ynone, Ycr3, 4, [4]uint8{0x0f, 0x22, 3, 0}},
	{AMOVL, Yml, Ynone, Ycr4, 4, [4]uint8{0x0f, 0x22, 4, 0}},
	{AMOVL, Yml, Ynone, Ycr8, 4, [4]uint8{0x0f, 0x22, 8, 0}},
	{AMOVQ, Yml, Ynone, Ycr0, 4, [4]uint8{0x0f, 0x22, 0, 0}},
	{AMOVQ, Yml, Ynone, Ycr2, 4, [4]uint8{0x0f, 0x22, 2, 0}},
	{AMOVQ, Yml, Ynone, Ycr3, 4, [4]uint8{0x0f, 0x22, 3, 0}},
	{AMOVQ, Yml, Ynone, Ycr4, 4, [4]uint8{0x0f, 0x22, 4, 0}},
	{AMOVQ, Yml, Ynone, Ycr8, 4, [4]uint8{0x0f, 0x22, 8, 0}},

	/* mov dr */
	{AMOVL, Ydr0, Ynone, Yml, 3, [4]uint8{0x0f, 0x21, 0, 0}},
	{AMOVL, Ydr6, Ynone, Yml, 3, [4]uint8{0x0f, 0x21, 6, 0}},
	{AMOVL, Ydr7, Ynone, Yml, 3, [4]uint8{0x0f, 0x21, 7, 0}},
	{AMOVQ, Ydr0, Ynone, Yml, 3, [4]uint8{0x0f, 0x21, 0, 0}},
	{AMOVQ, Ydr6, Ynone, Yml, 3, [4]uint8{0x0f, 0x21, 6, 0}},
	{AMOVQ, Ydr7, Ynone, Yml, 3, [4]uint8{0x0f, 0x21, 7, 0}},
	{AMOVL, Yml, Ynone, Ydr0, 4, [4]uint8{0x0f, 0x23, 0, 0}},
	{AMOVL, Yml, Ynone, Ydr6, 4, [4]uint8{0x0f, 0x23, 6, 0}},
	{AMOVL, Yml, Ynone, Ydr7, 4, [4]uint8{0x0f, 0x23, 7, 0}},
	{AMOVQ, Yml, Ynone, Ydr0, 4, [4]uint8{0x0f, 0x23, 0, 0}},
	{AMOVQ, Yml, Ynone, Ydr6, 4, [4]uint8{0x0f, 0x23, 6, 0}},
	{AMOVQ, Yml, Ynone, Ydr7, 4, [4]uint8{0x0f, 0x23, 7, 0}},

	/* mov tr */
	{AMOVL, Ytr6, Ynone, Yml, 3, [4]uint8{0x0f, 0x24, 6, 0}},
	{AMOVL, Ytr7, Ynone, Yml, 3, [4]uint8{0x0f, 0x24, 7, 0}},
	{AMOVL, Yml, Ynone, Ytr6, 4, [4]uint8{0x0f, 0x26, 6, E}},
	{AMOVL, Yml, Ynone, Ytr7, 4, [4]uint8{0x0f, 0x26, 7, E}},

	/* lgdt, sgdt, lidt, sidt */
	{AMOVL, Ym, Ynone, Ygdtr, 4, [4]uint8{0x0f, 0x01, 2, 0}},
	{AMOVL, Ygdtr, Ynone, Ym, 3, [4]uint8{0x0f, 0x01, 0, 0}},
	{AMOVL, Ym, Ynone, Yidtr, 4, [4]uint8{0x0f, 0x01, 3, 0}},
	{AMOVL, Yidtr, Ynone, Ym, 3, [4]uint8{0x0f, 0x01, 1, 0}},
	{AMOVQ, Ym, Ynone, Ygdtr, 4, [4]uint8{0x0f, 0x01, 2, 0}},
	{AMOVQ, Ygdtr, Ynone, Ym, 3, [4]uint8{0x0f, 0x01, 0, 0}},
	{AMOVQ, Ym, Ynone, Yidtr, 4, [4]uint8{0x0f, 0x01, 3, 0}},
	{AMOVQ, Yidtr, Ynone, Ym, 3, [4]uint8{0x0f, 0x01, 1, 0}},

	/* lldt, sldt */
	{AMOVW, Yml, Ynone, Yldtr, 4, [4]uint8{0x0f, 0x00, 2, 0}},
	{AMOVW, Yldtr, Ynone, Yml, 3, [4]uint8{0x0f, 0x00, 0, 0}},

	/* lmsw, smsw */
	{AMOVW, Yml, Ynone, Ymsw, 4, [4]uint8{0x0f, 0x01, 6, 0}},
	{AMOVW, Ymsw, Ynone, Yml, 3, [4]uint8{0x0f, 0x01, 4, 0}},

	/* ltr, str */
	{AMOVW, Yml, Ynone, Ytask, 4, [4]uint8{0x0f, 0x00, 3, 0}},
	{AMOVW, Ytask, Ynone, Yml, 3, [4]uint8{0x0f, 0x00, 1, 0}},

	/* load full pointer - unsupported
	{AMOVL, Yml, Ycol, 5, [4]uint8{0, 0, 0, 0}},
	{AMOVW, Yml, Ycol, 5, [4]uint8{Pe, 0, 0, 0}},
	*/

	/* double shift */
	{ASHLL, Yi8, Yrl, Yml, 6, [4]uint8{0xa4, 0xa5, 0, 0}},
	{ASHLL, Ycl, Yrl, Yml, 6, [4]uint8{0xa4, 0xa5, 0, 0}},
	{ASHLL, Ycx, Yrl, Yml, 6, [4]uint8{0xa4, 0xa5, 0, 0}},
	{ASHRL, Yi8, Yrl, Yml, 6, [4]uint8{0xac, 0xad, 0, 0}},
	{ASHRL, Ycl, Yrl, Yml, 6, [4]uint8{0xac, 0xad, 0, 0}},
	{ASHRL, Ycx, Yrl, Yml, 6, [4]uint8{0xac, 0xad, 0, 0}},
	{ASHLQ, Yi8, Yrl, Yml, 6, [4]uint8{Pw, 0xa4, 0xa5, 0}},
	{ASHLQ, Ycl, Yrl, Yml, 6, [4]uint8{Pw, 0xa4, 0xa5, 0}},
	{ASHLQ, Ycx, Yrl, Yml, 6, [4]uint8{Pw, 0xa4, 0xa5, 0}},
	{ASHRQ, Yi8, Yrl, Yml, 6, [4]uint8{Pw, 0xac, 0xad, 0}},
	{ASHRQ, Ycl, Yrl, Yml, 6, [4]uint8{Pw, 0xac, 0xad, 0}},
	{ASHRQ, Ycx, Yrl, Yml, 6, [4]uint8{Pw, 0xac, 0xad, 0}},
	{ASHLW, Yi8, Yrl, Yml, 6, [4]uint8{Pe, 0xa4, 0xa5, 0}},
	{ASHLW, Ycl, Yrl, Yml, 6, [4]uint8{Pe, 0xa4, 0xa5, 0}},
	{ASHLW, Ycx, Yrl, Yml, 6, [4]uint8{Pe, 0xa4, 0xa5, 0}},
	{ASHRW, Yi8, Yrl, Yml, 6, [4]uint8{Pe, 0xac, 0xad, 0}},
	{ASHRW, Ycl, Yrl, Yml, 6, [4]uint8{Pe, 0xac, 0xad, 0}},
	{ASHRW, Ycx, Yrl, Yml, 6, [4]uint8{Pe, 0xac, 0xad, 0}},

	/* load TLS base */
	{AMOVL, Ytls, Ynone, Yrl, 7, [4]uint8{0, 0, 0, 0}},
	{AMOVQ, Ytls, Ynone, Yrl, 7, [4]uint8{0, 0, 0, 0}},
	{0, 0, 0, 0, 0, [4]uint8{}},
}

func isax(a *obj.Addr) bool {
	switch a.Reg {
	case REG_AX, REG_AL, REG_AH:
		return true
	}

	if a.Index == int32(REG_AX) {
		return true
	}
	return false
}

func subreg(p *obj.Prog, from, to objabi.RBaseType) {
	if false { /* debug['Q'] */
		fmt.Printf("\n%v\ts/%v/%v/\n", p, Rconv(from), Rconv(to))
	}

	if p.From.Reg == from {
		p.From.Reg = to
		p.Ft = 0
	}

	if p.To.Reg == from {
		p.To.Reg = to
		p.Tt = 0
	}

	if objabi.RBaseType(p.From.Index) == from {
		p.From.Index = int32(to)
		p.Ft = 0
	}

	if objabi.RBaseType(p.To.Index) == from {
		p.To.Index = int32(to)
		p.Tt = 0
	}

	if false { /* debug['Q'] */
		fmt.Printf("%v\n", p)
	}
}

func mediaop(ctxt *obj.Link, o *Optab, op int, osize int, z int) int {
	switch op {
	case Pm, Pe, Pf2, Pf3:
		if osize != 1 {
			if op != Pm {
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
			}
			ctxt.Andptr[0] = Pm
			ctxt.Andptr = ctxt.Andptr[1:]
			z++
			op = int(o.op[z])
			break
		}
		fallthrough

	default:
		if -cap(ctxt.Andptr) == -cap(ctxt.And) || ctxt.And[-cap(ctxt.Andptr)+cap(ctxt.And[:])-1] != Pm {
			ctxt.Andptr[0] = Pm
			ctxt.Andptr = ctxt.Andptr[1:]
		}
	}

	ctxt.Andptr[0] = byte(op)
	ctxt.Andptr = ctxt.Andptr[1:]
	return z
}

var bpduff1 = []byte{
	0x48, 0x89, 0x6c, 0x24, 0xf0, // MOVQ BP, -16(SP)
	0x48, 0x8d, 0x6c, 0x24, 0xf0, // LEAQ -16(SP), BP
}

var bpduff2 = []byte{
	0x48, 0x8b, 0x6d, 0x00, // MOVQ 0(BP), BP
}

func doasm(ctxt *obj.Link, p *obj.Prog) {
	ctxt.Curp = p // TODO

	o := opindex[p.As&objabi.AMask]

	if o == nil {
		ctxt.Diag("asmins: missing op %v", p)
		return
	}

	pre := prefixof(ctxt, p, &p.From)
	if pre != 0 {
		ctxt.Andptr[0] = byte(pre)
		ctxt.Andptr = ctxt.Andptr[1:]
	}
	pre = prefixof(ctxt, p, &p.To)
	if pre != 0 {
		ctxt.Andptr[0] = byte(pre)
		ctxt.Andptr = ctxt.Andptr[1:]
	}

	// TODO(chai2010): This special case is for SHRQ $3, AX:DX,
	// which encodes as SHRQ $32(DX*0), AX.
	// Similarly SHRQ CX, AX:DX is really SHRQ CX(DX*0), AX.
	// Change encoding generated by assemblers and compilers and remove.
	if (p.From.Type == obj.TYPE_CONST || p.From.Type == obj.TYPE_REG) && p.From.Index != REG_NONE && p.From.Scale == 0 {
		p.From3 = new(obj.Addr)
		p.From3.Type = obj.TYPE_REG
		p.From3.Reg = objabi.RBaseType(p.From.Index)
		p.From.Index = 0
	}

	// TODO(chai2010): This special case is for PINSRQ etc, CMPSD etc.
	// Change encoding generated by assemblers and compilers (if any) and remove.
	switch p.As {
	case AIMUL3Q, APEXTRW, APINSRW, APINSRD, APINSRQ, APSHUFHW, APSHUFL, APSHUFW, ASHUFPD, ASHUFPS, AAESKEYGENASSIST, APSHUFD, APCLMULQDQ:
		if p.From3Type() == obj.TYPE_NONE {
			p.From3 = new(obj.Addr)
			*p.From3 = p.From
			p.From = obj.Addr{}
			p.From.Type = obj.TYPE_CONST
			p.From.Offset = p.To.Offset
			p.To.Offset = 0
		}
	case ACMPSD, ACMPSS, ACMPPS, ACMPPD:
		if p.From3Type() == obj.TYPE_NONE {
			p.From3 = new(obj.Addr)
			*p.From3 = p.To
			p.To = obj.Addr{}
			p.To.Type = obj.TYPE_CONST
			p.To.Offset = p.From3.Offset
			p.From3.Offset = 0
		}
	}

	if p.Ft == 0 {
		p.Ft = uint8(oclass(ctxt, p, &p.From))
	}
	if p.Tt == 0 {
		p.Tt = uint8(oclass(ctxt, p, &p.To))
	}

	ft := int(p.Ft) * Ymax
	f3t := Ynone * Ymax
	if p.From3 != nil {
		f3t = oclass(ctxt, p, p.From3) * Ymax
	}
	tt := int(p.Tt) * Ymax

	xo := 0
	if o.op[0] == 0x0f {
		xo = 1
	}

	z := 0
	var a *obj.Addr
	var l int
	var op int
	var q *obj.Prog
	var r *obj.Reloc
	var rel obj.Reloc
	var v int64
	for i := range o.ytab {
		yt := &o.ytab[i]
		if ycover[ft+int(yt.from)] != 0 && ycover[f3t+int(yt.from3)] != 0 && ycover[tt+int(yt.to)] != 0 {
			switch o.prefix {
			case Px1: /* first option valid only in 32-bit mode */
				if ctxt.Mode == 64 && z == 0 {
					z += int(yt.zoffset) + xo
					continue
				}
			case Pq: /* 16 bit escape and opcode escape */
				ctxt.Andptr[0] = Pe
				ctxt.Andptr = ctxt.Andptr[1:]

				ctxt.Andptr[0] = Pm
				ctxt.Andptr = ctxt.Andptr[1:]

			case Pq3: /* 16 bit escape, Rex.w, and opcode escape */
				ctxt.Andptr[0] = Pe
				ctxt.Andptr = ctxt.Andptr[1:]

				ctxt.Andptr[0] = Pw
				ctxt.Andptr = ctxt.Andptr[1:]
				ctxt.Andptr[0] = Pm
				ctxt.Andptr = ctxt.Andptr[1:]

			case Pf2, /* xmm opcode escape */
				Pf3:
				ctxt.Andptr[0] = byte(o.prefix)
				ctxt.Andptr = ctxt.Andptr[1:]

				ctxt.Andptr[0] = Pm
				ctxt.Andptr = ctxt.Andptr[1:]

			case Pm: /* opcode escape */
				ctxt.Andptr[0] = Pm
				ctxt.Andptr = ctxt.Andptr[1:]

			case Pe: /* 16 bit escape */
				ctxt.Andptr[0] = Pe
				ctxt.Andptr = ctxt.Andptr[1:]

			case Pw: /* 64-bit escape */
				if p.Mode != 64 {
					ctxt.Diag("asmins: illegal 64: %v", p)
				}
				ctxt.Rexflag |= Pw

			case Pw8: /* 64-bit escape if z >= 8 */
				if z >= 8 {
					if p.Mode != 64 {
						ctxt.Diag("asmins: illegal 64: %v", p)
					}
					ctxt.Rexflag |= Pw
				}

			case Pb: /* botch */
				if p.Mode != 64 && (isbadbyte(&p.From) || isbadbyte(&p.To)) {
					goto bad
				}
				// NOTE(rsc): This is probably safe to do always,
				// but when enabled it chooses different encodings
				// than the old cmd/internal/obj/i386 code did,
				// which breaks our "same bits out" checks.
				// In particular, CMPB AX, $0 encodes as 80 f8 00
				// in the original obj/i386, and it would encode
				// (using a valid, shorter form) as 3c 00 if we enabled
				// the call to bytereg here.
				if p.Mode == 64 {
					bytereg(&p.From, &p.Ft)
					bytereg(&p.To, &p.Tt)
				}

			case P32: /* 32 bit but illegal if 64-bit mode */
				if p.Mode == 64 {
					ctxt.Diag("asmins: illegal in 64-bit mode: %v", p)
				}

			case Py: /* 64-bit only, no prefix */
				if p.Mode != 64 {
					ctxt.Diag("asmins: illegal in %d-bit mode: %v", p.Mode, p)
				}

			case Py1: /* 64-bit only if z < 1, no prefix */
				if z < 1 && p.Mode != 64 {
					ctxt.Diag("asmins: illegal in %d-bit mode: %v", p.Mode, p)
				}

			case Py3: /* 64-bit only if z < 3, no prefix */
				if z < 3 && p.Mode != 64 {
					ctxt.Diag("asmins: illegal in %d-bit mode: %v", p.Mode, p)
				}
			}

			if z >= len(o.op) {
				log.Fatalf("asmins bad table %v", p)
			}
			op = int(o.op[z])
			if op == 0x0f {
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				z++
				op = int(o.op[z])
			}

			switch yt.zcase {
			default:
				ctxt.Diag("asmins: unknown z %d %v", yt.zcase, p)
				return

			case Zpseudo:
				break

			case Zlit:
				for ; ; z++ {
					op = int(o.op[z])
					if op == 0 {
						break
					}
					ctxt.Andptr[0] = byte(op)
					ctxt.Andptr = ctxt.Andptr[1:]
				}

			case Zlitm_r:
				for ; ; z++ {
					op = int(o.op[z])
					if op == 0 {
						break
					}
					ctxt.Andptr[0] = byte(op)
					ctxt.Andptr = ctxt.Andptr[1:]
				}
				asmand(ctxt, p, &p.From, &p.To)

			case Zmb_r:
				bytereg(&p.From, &p.Ft)
				fallthrough

				/* fall through */
			case Zm_r:
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]

				asmand(ctxt, p, &p.From, &p.To)

			case Zm2_r:
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				ctxt.Andptr[0] = byte(o.op[z+1])
				ctxt.Andptr = ctxt.Andptr[1:]
				asmand(ctxt, p, &p.From, &p.To)

			case Zm_r_xm:
				mediaop(ctxt, o, op, int(yt.zoffset), z)
				asmand(ctxt, p, &p.From, &p.To)

			case Zm_r_xm_nr:
				ctxt.Rexflag = 0
				mediaop(ctxt, o, op, int(yt.zoffset), z)
				asmand(ctxt, p, &p.From, &p.To)

			case Zm_r_i_xm:
				mediaop(ctxt, o, op, int(yt.zoffset), z)
				asmand(ctxt, p, &p.From, p.From3)
				ctxt.Andptr[0] = byte(p.To.Offset)
				ctxt.Andptr = ctxt.Andptr[1:]

			case Zm_r_3d:
				ctxt.Andptr[0] = 0x0f
				ctxt.Andptr = ctxt.Andptr[1:]
				ctxt.Andptr[0] = 0x0f
				ctxt.Andptr = ctxt.Andptr[1:]
				asmand(ctxt, p, &p.From, &p.To)
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]

			case Zibm_r:
				for {
					tmp1 := z
					z++
					op = int(o.op[tmp1])
					if op == 0 {
						break
					}
					ctxt.Andptr[0] = byte(op)
					ctxt.Andptr = ctxt.Andptr[1:]
				}
				asmand(ctxt, p, p.From3, &p.To)
				ctxt.Andptr[0] = byte(p.From.Offset)
				ctxt.Andptr = ctxt.Andptr[1:]

			case Zaut_r:
				ctxt.Andptr[0] = 0x8d
				ctxt.Andptr = ctxt.Andptr[1:] /* leal */
				if p.From.Type != obj.TYPE_ADDR {
					ctxt.Diag("asmins: Zaut sb type ADDR")
				}
				p.From.Type = obj.TYPE_MEM
				asmand(ctxt, p, &p.From, &p.To)
				p.From.Type = obj.TYPE_ADDR

			case Zm_o:
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				asmando(ctxt, p, &p.From, int(o.op[z+1]))

			case Zr_m:
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				asmand(ctxt, p, &p.To, &p.From)

			case Zr_m_xm:
				mediaop(ctxt, o, op, int(yt.zoffset), z)
				asmand(ctxt, p, &p.To, &p.From)

			case Zr_m_xm_nr:
				ctxt.Rexflag = 0
				mediaop(ctxt, o, op, int(yt.zoffset), z)
				asmand(ctxt, p, &p.To, &p.From)

			case Zo_m:
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				asmando(ctxt, p, &p.To, int(o.op[z+1]))

			case Zcallindreg:
				r = obj_Addrel(ctxt.Cursym)
				r.Off = int32(p.Pc)
				r.Type = obj.R_CALLIND
				r.Siz = 0
				fallthrough

			case Zo_m64:
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				asmandsz(ctxt, p, &p.To, int(o.op[z+1]), 0, 1)

			case Zm_ibo:
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				asmando(ctxt, p, &p.From, int(o.op[z+1]))
				ctxt.Andptr[0] = byte(vaddr(ctxt, p, &p.To, nil))
				ctxt.Andptr = ctxt.Andptr[1:]

			case Zibo_m:
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				asmando(ctxt, p, &p.To, int(o.op[z+1]))
				ctxt.Andptr[0] = byte(vaddr(ctxt, p, &p.From, nil))
				ctxt.Andptr = ctxt.Andptr[1:]

			case Zibo_m_xm:
				z = mediaop(ctxt, o, op, int(yt.zoffset), z)
				asmando(ctxt, p, &p.To, int(o.op[z+1]))
				ctxt.Andptr[0] = byte(vaddr(ctxt, p, &p.From, nil))
				ctxt.Andptr = ctxt.Andptr[1:]

			case Z_ib, Zib_:
				if yt.zcase == Zib_ {
					a = &p.From
				} else {
					a = &p.To
				}
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				ctxt.Andptr[0] = byte(vaddr(ctxt, p, a, nil))
				ctxt.Andptr = ctxt.Andptr[1:]

			case Zib_rp:
				ctxt.Rexflag |= regrex[p.To.Reg] & (Rxb | 0x40)
				ctxt.Andptr[0] = byte(op + reg[p.To.Reg])
				ctxt.Andptr = ctxt.Andptr[1:]
				ctxt.Andptr[0] = byte(vaddr(ctxt, p, &p.From, nil))
				ctxt.Andptr = ctxt.Andptr[1:]

			case Zil_rp:
				ctxt.Rexflag |= regrex[p.To.Reg] & Rxb
				ctxt.Andptr[0] = byte(op + reg[p.To.Reg])
				ctxt.Andptr = ctxt.Andptr[1:]
				if o.prefix == Pe {
					v = vaddr(ctxt, p, &p.From, nil)
					ctxt.Andptr[0] = byte(v)
					ctxt.Andptr = ctxt.Andptr[1:]
					ctxt.Andptr[0] = byte(v >> 8)
					ctxt.Andptr = ctxt.Andptr[1:]
				} else {
					relput4(ctxt, p, &p.From)
				}

			case Zo_iw:
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				if p.From.Type != obj.TYPE_NONE {
					v = vaddr(ctxt, p, &p.From, nil)
					ctxt.Andptr[0] = byte(v)
					ctxt.Andptr = ctxt.Andptr[1:]
					ctxt.Andptr[0] = byte(v >> 8)
					ctxt.Andptr = ctxt.Andptr[1:]
				}

			case Ziq_rp:
				v = vaddr(ctxt, p, &p.From, &rel)
				l = int(v >> 32)
				if l == 0 && rel.Siz != 8 {
					//p->mark |= 0100;
					//print("zero: %llux %v\n", v, p);
					ctxt.Rexflag &^= (0x40 | Rxw)

					ctxt.Rexflag |= regrex[p.To.Reg] & Rxb
					ctxt.Andptr[0] = byte(0xb8 + reg[p.To.Reg])
					ctxt.Andptr = ctxt.Andptr[1:]
					if rel.Type != 0 {
						r = obj_Addrel(ctxt.Cursym)
						*r = rel
						r.Off = int32(p.Pc + int64(-cap(ctxt.Andptr)+cap(ctxt.And[:])))
					}

					put4(ctxt, int32(v))
				} else if l == -1 && uint64(v)&(uint64(1)<<31) != 0 { /* sign extend */

					//p->mark |= 0100;
					//print("sign: %llux %v\n", v, p);
					ctxt.Andptr[0] = 0xc7
					ctxt.Andptr = ctxt.Andptr[1:]

					asmando(ctxt, p, &p.To, 0)
					put4(ctxt, int32(v)) /* need all 8 */
				} else {
					//print("all: %llux %v\n", v, p);
					ctxt.Rexflag |= regrex[p.To.Reg] & Rxb

					ctxt.Andptr[0] = byte(op + reg[p.To.Reg])
					ctxt.Andptr = ctxt.Andptr[1:]
					if rel.Type != 0 {
						r = obj_Addrel(ctxt.Cursym)
						*r = rel
						r.Off = int32(p.Pc + int64(-cap(ctxt.Andptr)+cap(ctxt.And[:])))
					}

					put8(ctxt, v)
				}

			case Zib_rr:
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				asmand(ctxt, p, &p.To, &p.To)
				ctxt.Andptr[0] = byte(vaddr(ctxt, p, &p.From, nil))
				ctxt.Andptr = ctxt.Andptr[1:]

			case Z_il, Zil_:
				if yt.zcase == Zil_ {
					a = &p.From
				} else {
					a = &p.To
				}
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				if o.prefix == Pe {
					v = vaddr(ctxt, p, a, nil)
					ctxt.Andptr[0] = byte(v)
					ctxt.Andptr = ctxt.Andptr[1:]
					ctxt.Andptr[0] = byte(v >> 8)
					ctxt.Andptr = ctxt.Andptr[1:]
				} else {
					relput4(ctxt, p, a)
				}

			case Zm_ilo, Zilo_m:
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				if yt.zcase == Zilo_m {
					a = &p.From
					asmando(ctxt, p, &p.To, int(o.op[z+1]))
				} else {
					a = &p.To
					asmando(ctxt, p, &p.From, int(o.op[z+1]))
				}

				if o.prefix == Pe {
					v = vaddr(ctxt, p, a, nil)
					ctxt.Andptr[0] = byte(v)
					ctxt.Andptr = ctxt.Andptr[1:]
					ctxt.Andptr[0] = byte(v >> 8)
					ctxt.Andptr = ctxt.Andptr[1:]
				} else {
					relput4(ctxt, p, a)
				}

			case Zil_rr:
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				asmand(ctxt, p, &p.To, &p.To)
				if o.prefix == Pe {
					v = vaddr(ctxt, p, &p.From, nil)
					ctxt.Andptr[0] = byte(v)
					ctxt.Andptr = ctxt.Andptr[1:]
					ctxt.Andptr[0] = byte(v >> 8)
					ctxt.Andptr = ctxt.Andptr[1:]
				} else {
					relput4(ctxt, p, &p.From)
				}

			case Z_rp:
				ctxt.Rexflag |= regrex[p.To.Reg] & (Rxb | 0x40)
				ctxt.Andptr[0] = byte(op + reg[p.To.Reg])
				ctxt.Andptr = ctxt.Andptr[1:]

			case Zrp_:
				ctxt.Rexflag |= regrex[p.From.Reg] & (Rxb | 0x40)
				ctxt.Andptr[0] = byte(op + reg[p.From.Reg])
				ctxt.Andptr = ctxt.Andptr[1:]

			case Zclr:
				ctxt.Rexflag &^= Pw
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				asmand(ctxt, p, &p.To, &p.To)

			case Zcallcon, Zjmpcon:
				if yt.zcase == Zcallcon {
					ctxt.Andptr[0] = byte(op)
					ctxt.Andptr = ctxt.Andptr[1:]
				} else {
					ctxt.Andptr[0] = byte(o.op[z+1])
					ctxt.Andptr = ctxt.Andptr[1:]
				}
				r = obj_Addrel(ctxt.Cursym)
				r.Off = int32(p.Pc + int64(-cap(ctxt.Andptr)+cap(ctxt.And[:])))
				r.Type = obj.R_PCREL
				r.Siz = 4
				r.Add = p.To.Offset
				put4(ctxt, 0)

			case Zcallind:
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				ctxt.Andptr[0] = byte(o.op[z+1])
				ctxt.Andptr = ctxt.Andptr[1:]
				r = obj_Addrel(ctxt.Cursym)
				r.Off = int32(p.Pc + int64(-cap(ctxt.Andptr)+cap(ctxt.And[:])))
				r.Type = obj.R_ADDR
				r.Siz = 4
				r.Add = p.To.Offset
				r.Sym = p.To.Sym
				put4(ctxt, 0)

			case Zcall, Zcallduff:
				if p.To.Sym == nil {
					ctxt.Diag("call without target")
					log.Fatalf("bad code")
				}

				if yt.zcase == Zcallduff && ctxt.Flag_dynlink {
					ctxt.Diag("directly calling duff when dynamically linking Wa")
				}

				if Framepointer_enabled && yt.zcase == Zcallduff && p.Mode == 64 {
					// Maintain BP around call, since duffcopy/duffzero can't do it
					// (the call jumps into the middle of the function).
					// This makes it possible to see call sites for duffcopy/duffzero in
					// BP-based profiling tools like Linux perf (which is the
					// whole point of Framepointer_enabled).
					// MOVQ BP, -16(SP)
					// LEAQ -16(SP), BP
					copy(ctxt.Andptr, bpduff1)
					ctxt.Andptr = ctxt.Andptr[len(bpduff1):]
				}
				ctxt.Andptr[0] = byte(op)
				ctxt.Andptr = ctxt.Andptr[1:]
				r = obj_Addrel(ctxt.Cursym)
				r.Off = int32(p.Pc + int64(-cap(ctxt.Andptr)+cap(ctxt.And[:])))
				r.Sym = p.To.Sym
				r.Add = p.To.Offset
				r.Type = obj.R_CALL
				r.Siz = 4
				put4(ctxt, 0)

				if Framepointer_enabled && yt.zcase == Zcallduff && p.Mode == 64 {
					// Pop BP pushed above.
					// MOVQ 0(BP), BP
					copy(ctxt.Andptr, bpduff2)
					ctxt.Andptr = ctxt.Andptr[len(bpduff2):]
				}

			// TODO: jump across functions needs reloc
			case Zbr, Zjmp, Zloop:
				if p.To.Sym != nil {
					if yt.zcase != Zjmp {
						ctxt.Diag("branch to ATEXT")
						log.Fatalf("bad code")
					}

					ctxt.Andptr[0] = byte(o.op[z+1])
					ctxt.Andptr = ctxt.Andptr[1:]
					r = obj_Addrel(ctxt.Cursym)
					r.Off = int32(p.Pc + int64(-cap(ctxt.Andptr)+cap(ctxt.And[:])))
					r.Sym = p.To.Sym
					r.Type = obj.R_PCREL
					r.Siz = 4
					put4(ctxt, 0)
					break
				}

				// Assumes q is in this function.
				// TODO: Check in input, preserve in brchain.

				// Fill in backward jump now.
				q = p.Pcond

				if q == nil {
					ctxt.Diag("jmp/branch/loop without target")
					log.Fatalf("bad code")
				}

				if p.Back&1 != 0 {
					v = q.Pc - (p.Pc + 2)
					if v >= -128 {
						if p.As == AJCXZL {
							ctxt.Andptr[0] = 0x67
							ctxt.Andptr = ctxt.Andptr[1:]
						}
						ctxt.Andptr[0] = byte(op)
						ctxt.Andptr = ctxt.Andptr[1:]
						ctxt.Andptr[0] = byte(v)
						ctxt.Andptr = ctxt.Andptr[1:]
					} else if yt.zcase == Zloop {
						ctxt.Diag("loop too far: %v", p)
					} else {
						v -= 5 - 2
						if yt.zcase == Zbr {
							ctxt.Andptr[0] = 0x0f
							ctxt.Andptr = ctxt.Andptr[1:]
							v--
						}

						ctxt.Andptr[0] = byte(o.op[z+1])
						ctxt.Andptr = ctxt.Andptr[1:]
						ctxt.Andptr[0] = byte(v)
						ctxt.Andptr = ctxt.Andptr[1:]
						ctxt.Andptr[0] = byte(v >> 8)
						ctxt.Andptr = ctxt.Andptr[1:]
						ctxt.Andptr[0] = byte(v >> 16)
						ctxt.Andptr = ctxt.Andptr[1:]
						ctxt.Andptr[0] = byte(v >> 24)
						ctxt.Andptr = ctxt.Andptr[1:]
					}

					break
				}

				// Annotate target; will fill in later.
				p.Forwd = q.Rel

				q.Rel = p
				if p.Back&2 != 0 { // short
					if p.As == AJCXZL {
						ctxt.Andptr[0] = 0x67
						ctxt.Andptr = ctxt.Andptr[1:]
					}
					ctxt.Andptr[0] = byte(op)
					ctxt.Andptr = ctxt.Andptr[1:]
					ctxt.Andptr[0] = 0
					ctxt.Andptr = ctxt.Andptr[1:]
				} else if yt.zcase == Zloop {
					ctxt.Diag("loop too far: %v", p)
				} else {
					if yt.zcase == Zbr {
						ctxt.Andptr[0] = 0x0f
						ctxt.Andptr = ctxt.Andptr[1:]
					}
					ctxt.Andptr[0] = byte(o.op[z+1])
					ctxt.Andptr = ctxt.Andptr[1:]
					ctxt.Andptr[0] = 0
					ctxt.Andptr = ctxt.Andptr[1:]
					ctxt.Andptr[0] = 0
					ctxt.Andptr = ctxt.Andptr[1:]
					ctxt.Andptr[0] = 0
					ctxt.Andptr = ctxt.Andptr[1:]
					ctxt.Andptr[0] = 0
					ctxt.Andptr = ctxt.Andptr[1:]
				}

				break

			/*
				v = q->pc - p->pc - 2;
				if((v >= -128 && v <= 127) || p->pc == -1 || q->pc == -1) {
					*ctxt->andptr++ = op;
					*ctxt->andptr++ = v;
				} else {
					v -= 5-2;
					if(yt.zcase == Zbr) {
						*ctxt->andptr++ = 0x0f;
						v--;
					}
					*ctxt->andptr++ = o->op[z+1];
					*ctxt->andptr++ = v;
					*ctxt->andptr++ = v>>8;
					*ctxt->andptr++ = v>>16;
					*ctxt->andptr++ = v>>24;
				}
			*/

			case Zbyte:
				v = vaddr(ctxt, p, &p.From, &rel)
				if rel.Siz != 0 {
					rel.Siz = uint8(op)
					r = obj_Addrel(ctxt.Cursym)
					*r = rel
					r.Off = int32(p.Pc + int64(-cap(ctxt.Andptr)+cap(ctxt.And[:])))
				}

				ctxt.Andptr[0] = byte(v)
				ctxt.Andptr = ctxt.Andptr[1:]
				if op > 1 {
					ctxt.Andptr[0] = byte(v >> 8)
					ctxt.Andptr = ctxt.Andptr[1:]
					if op > 2 {
						ctxt.Andptr[0] = byte(v >> 16)
						ctxt.Andptr = ctxt.Andptr[1:]
						ctxt.Andptr[0] = byte(v >> 24)
						ctxt.Andptr = ctxt.Andptr[1:]
						if op > 4 {
							ctxt.Andptr[0] = byte(v >> 32)
							ctxt.Andptr = ctxt.Andptr[1:]
							ctxt.Andptr[0] = byte(v >> 40)
							ctxt.Andptr = ctxt.Andptr[1:]
							ctxt.Andptr[0] = byte(v >> 48)
							ctxt.Andptr = ctxt.Andptr[1:]
							ctxt.Andptr[0] = byte(v >> 56)
							ctxt.Andptr = ctxt.Andptr[1:]
						}
					}
				}
			}

			return
		}
		z += int(yt.zoffset) + xo
	}
	for mo := ymovtab; mo[0].as != 0; mo = mo[1:] {
		var pp obj.Prog
		var t []byte
		if p.As == mo[0].as {
			if ycover[ft+int(mo[0].ft)] != 0 && ycover[f3t+int(mo[0].f3t)] != 0 && ycover[tt+int(mo[0].tt)] != 0 {
				t = mo[0].op[:]
				switch mo[0].code {
				default:
					ctxt.Diag("asmins: unknown mov %d %v", mo[0].code, p)

				case 0: /* lit */
					for z = 0; t[z] != E; z++ {
						ctxt.Andptr[0] = t[z]
						ctxt.Andptr = ctxt.Andptr[1:]
					}

				case 1: /* r,m */
					ctxt.Andptr[0] = t[0]
					ctxt.Andptr = ctxt.Andptr[1:]

					asmando(ctxt, p, &p.To, int(t[1]))

				case 2: /* m,r */
					ctxt.Andptr[0] = t[0]
					ctxt.Andptr = ctxt.Andptr[1:]

					asmando(ctxt, p, &p.From, int(t[1]))

				case 3: /* r,m - 2op */
					ctxt.Andptr[0] = t[0]
					ctxt.Andptr = ctxt.Andptr[1:]

					ctxt.Andptr[0] = t[1]
					ctxt.Andptr = ctxt.Andptr[1:]
					asmando(ctxt, p, &p.To, int(t[2]))
					ctxt.Rexflag |= regrex[p.From.Reg] & (Rxr | 0x40)

				case 4: /* m,r - 2op */
					ctxt.Andptr[0] = t[0]
					ctxt.Andptr = ctxt.Andptr[1:]

					ctxt.Andptr[0] = t[1]
					ctxt.Andptr = ctxt.Andptr[1:]
					asmando(ctxt, p, &p.From, int(t[2]))
					ctxt.Rexflag |= regrex[p.To.Reg] & (Rxr | 0x40)

				case 5: /* load full pointer, trash heap */
					if t[0] != 0 {
						ctxt.Andptr[0] = t[0]
						ctxt.Andptr = ctxt.Andptr[1:]
					}
					switch objabi.RBaseType(p.To.Index) {
					default:
						goto bad

					case REG_DS:
						ctxt.Andptr[0] = 0xc5
						ctxt.Andptr = ctxt.Andptr[1:]

					case REG_SS:
						ctxt.Andptr[0] = 0x0f
						ctxt.Andptr = ctxt.Andptr[1:]
						ctxt.Andptr[0] = 0xb2
						ctxt.Andptr = ctxt.Andptr[1:]

					case REG_ES:
						ctxt.Andptr[0] = 0xc4
						ctxt.Andptr = ctxt.Andptr[1:]

					case REG_FS:
						ctxt.Andptr[0] = 0x0f
						ctxt.Andptr = ctxt.Andptr[1:]
						ctxt.Andptr[0] = 0xb4
						ctxt.Andptr = ctxt.Andptr[1:]

					case REG_GS:
						ctxt.Andptr[0] = 0x0f
						ctxt.Andptr = ctxt.Andptr[1:]
						ctxt.Andptr[0] = 0xb5
						ctxt.Andptr = ctxt.Andptr[1:]
					}

					asmand(ctxt, p, &p.From, &p.To)

				case 6: /* double shift */
					if t[0] == Pw {
						if p.Mode != 64 {
							ctxt.Diag("asmins: illegal 64: %v", p)
						}
						ctxt.Rexflag |= Pw
						t = t[1:]
					} else if t[0] == Pe {
						ctxt.Andptr[0] = Pe
						ctxt.Andptr = ctxt.Andptr[1:]
						t = t[1:]
					}

					switch p.From.Type {
					default:
						goto bad

					case obj.TYPE_CONST:
						ctxt.Andptr[0] = 0x0f
						ctxt.Andptr = ctxt.Andptr[1:]
						ctxt.Andptr[0] = t[0]
						ctxt.Andptr = ctxt.Andptr[1:]
						asmandsz(ctxt, p, &p.To, reg[p.From3.Reg], regrex[p.From3.Reg], 0)
						ctxt.Andptr[0] = byte(p.From.Offset)
						ctxt.Andptr = ctxt.Andptr[1:]

					case obj.TYPE_REG:
						switch p.From.Reg {
						default:
							goto bad

						case REG_CL, REG_CX:
							ctxt.Andptr[0] = 0x0f
							ctxt.Andptr = ctxt.Andptr[1:]
							ctxt.Andptr[0] = t[1]
							ctxt.Andptr = ctxt.Andptr[1:]
							asmandsz(ctxt, p, &p.To, reg[p.From3.Reg], regrex[p.From3.Reg], 0)
						}
					}

				// NOTE: The systems listed here are the ones that use the "TLS initial exec" model,
				// where you load the TLS base register into a register and then index off that
				// register to access the actual TLS variables. Systems that allow direct TLS access
				// are handled in prefixof above and should not be listed here.
				case 7: /* mov tls, r */
					if p.Mode == 64 && p.As != AMOVQ || p.Mode == 32 && p.As != AMOVL {
						ctxt.Diag("invalid load of TLS: %v", p)
					}

					if p.Mode == 32 {
						// NOTE: The systems listed here are the ones that use the "TLS initial exec" model,
						// where you load the TLS base register into a register and then index off that
						// register to access the actual TLS variables. Systems that allow direct TLS access
						// are handled in prefixof above and should not be listed here.
						switch ctxt.Headtype {
						default:
							log.Fatalf("unknown TLS base location for %s", ctxt.Headtype)

						case objabi.Hlinux:
							// ELF TLS base is 0(GS).
							pp.From = p.From

							pp.From.Type = obj.TYPE_MEM
							pp.From.Reg = REG_GS
							pp.From.Offset = 0
							pp.From.Index = REG_NONE
							pp.From.Scale = 0
							ctxt.Andptr[0] = 0x65
							ctxt.Andptr = ctxt.Andptr[1:] // GS
							ctxt.Andptr[0] = 0x8B
							ctxt.Andptr = ctxt.Andptr[1:]
							asmand(ctxt, p, &pp.From, &p.To)

						case objabi.Hwindows:
							// Windows TLS base is always 0x14(FS).
							pp.From = p.From

							pp.From.Type = obj.TYPE_MEM
							pp.From.Reg = REG_FS
							pp.From.Offset = 0x14
							pp.From.Index = REG_NONE
							pp.From.Scale = 0
							ctxt.Andptr[0] = 0x64
							ctxt.Andptr = ctxt.Andptr[1:] // FS
							ctxt.Andptr[0] = 0x8B
							ctxt.Andptr = ctxt.Andptr[1:]
							asmand(ctxt, p, &pp.From, &p.To)
						}
						break
					}

					switch ctxt.Headtype {
					default:
						log.Fatalf("unknown TLS base location for %s", ctxt.Headtype)

					case objabi.Hlinux:
						if ctxt.Flag_shared == 0 {
							log.Fatalf("unknown TLS base location for linux without -shared")
						}
						// Note that this is not generating the same insn as the other cases.
						//     MOV TLS, R_to
						// becomes
						//     movq g@gottpoff(%rip), R_to
						// which is encoded as
						//     movq 0(%rip), R_to
						// and a R_TLS_IE reloc. This all assumes the only tls variable we access
						// is g, which we can't check here, but will when we assemble the second
						// instruction.
						ctxt.Rexflag = Pw | (regrex[p.To.Reg] & Rxr)

						ctxt.Andptr[0] = 0x8B
						ctxt.Andptr = ctxt.Andptr[1:]
						ctxt.Andptr[0] = byte(0x05 | (reg[p.To.Reg] << 3))
						ctxt.Andptr = ctxt.Andptr[1:]
						r = obj_Addrel(ctxt.Cursym)
						r.Off = int32(p.Pc + int64(-cap(ctxt.Andptr)+cap(ctxt.And[:])))
						r.Type = obj.R_TLS_IE
						r.Siz = 4
						r.Add = -4
						put4(ctxt, 0)

					case objabi.Hwindows:
						// Windows TLS base is always 0x28(GS).
						pp.From = p.From

						pp.From.Type = obj.TYPE_MEM
						pp.From.Name = obj.NAME_NONE
						pp.From.Reg = REG_GS
						pp.From.Offset = 0x28
						pp.From.Index = REG_NONE
						pp.From.Scale = 0
						ctxt.Rexflag |= Pw
						ctxt.Andptr[0] = 0x65
						ctxt.Andptr = ctxt.Andptr[1:] // GS
						ctxt.Andptr[0] = 0x8B
						ctxt.Andptr = ctxt.Andptr[1:]
						asmand(ctxt, p, &pp.From, &p.To)
					}
				}
				return
			}
		}
	}
	goto bad

bad:
	if p.Mode != 64 {
		/*
		 * here, the assembly has failed.
		 * if its a byte instruction that has
		 * unaddressable registers, try to
		 * exchange registers and reissue the
		 * instruction with the operands renamed.
		 */
		pp := *p

		unbytereg(&pp.From, &pp.Ft)
		unbytereg(&pp.To, &pp.Tt)

		z := p.From.Reg
		if p.From.Type == obj.TYPE_REG && z >= REG_BP && z <= REG_DI {
			// TODO(chai2010): Use this code for x86-64 too. It has bug fixes not present in the amd64 code base.
			// For now, different to keep bit-for-bit compatibility.
			if p.Mode == 32 {
				breg := byteswapreg(ctxt, &p.To)
				if breg != REG_AX {
					ctxt.Andptr[0] = 0x87
					ctxt.Andptr = ctxt.Andptr[1:] /* xchg lhs,bx */
					asmando(ctxt, p, &p.From, reg[breg])
					subreg(&pp, z, breg)
					doasm(ctxt, &pp)
					ctxt.Andptr[0] = 0x87
					ctxt.Andptr = ctxt.Andptr[1:] /* xchg lhs,bx */
					asmando(ctxt, p, &p.From, reg[breg])
				} else {
					ctxt.Andptr[0] = byte(0x90 + reg[z])
					ctxt.Andptr = ctxt.Andptr[1:] /* xchg lsh,ax */
					subreg(&pp, z, REG_AX)
					doasm(ctxt, &pp)
					ctxt.Andptr[0] = byte(0x90 + reg[z])
					ctxt.Andptr = ctxt.Andptr[1:] /* xchg lsh,ax */
				}
				return
			}

			if isax(&p.To) || p.To.Type == obj.TYPE_NONE {
				// We certainly don't want to exchange
				// with AX if the op is MUL or DIV.
				ctxt.Andptr[0] = 0x87
				ctxt.Andptr = ctxt.Andptr[1:] /* xchg lhs,bx */
				asmando(ctxt, p, &p.From, reg[REG_BX])
				subreg(&pp, z, REG_BX)
				doasm(ctxt, &pp)
				ctxt.Andptr[0] = 0x87
				ctxt.Andptr = ctxt.Andptr[1:] /* xchg lhs,bx */
				asmando(ctxt, p, &p.From, reg[REG_BX])
			} else {
				ctxt.Andptr[0] = byte(0x90 + reg[z])
				ctxt.Andptr = ctxt.Andptr[1:] /* xchg lsh,ax */
				subreg(&pp, z, REG_AX)
				doasm(ctxt, &pp)
				ctxt.Andptr[0] = byte(0x90 + reg[z])
				ctxt.Andptr = ctxt.Andptr[1:] /* xchg lsh,ax */
			}
			return
		}

		z = p.To.Reg
		if p.To.Type == obj.TYPE_REG && z >= REG_BP && z <= REG_DI {
			// TODO(chai2010): Use this code for x86-64 too. It has bug fixes not present in the amd64 code base.
			// For now, different to keep bit-for-bit compatibility.
			if p.Mode == 32 {
				breg := byteswapreg(ctxt, &p.From)
				if breg != REG_AX {
					ctxt.Andptr[0] = 0x87
					ctxt.Andptr = ctxt.Andptr[1:] /* xchg rhs,bx */
					asmando(ctxt, p, &p.To, reg[breg])
					subreg(&pp, z, breg)
					doasm(ctxt, &pp)
					ctxt.Andptr[0] = 0x87
					ctxt.Andptr = ctxt.Andptr[1:] /* xchg rhs,bx */
					asmando(ctxt, p, &p.To, reg[breg])
				} else {
					ctxt.Andptr[0] = byte(0x90 + reg[z])
					ctxt.Andptr = ctxt.Andptr[1:] /* xchg rsh,ax */
					subreg(&pp, z, REG_AX)
					doasm(ctxt, &pp)
					ctxt.Andptr[0] = byte(0x90 + reg[z])
					ctxt.Andptr = ctxt.Andptr[1:] /* xchg rsh,ax */
				}
				return
			}

			if isax(&p.From) {
				ctxt.Andptr[0] = 0x87
				ctxt.Andptr = ctxt.Andptr[1:] /* xchg rhs,bx */
				asmando(ctxt, p, &p.To, reg[REG_BX])
				subreg(&pp, z, REG_BX)
				doasm(ctxt, &pp)
				ctxt.Andptr[0] = 0x87
				ctxt.Andptr = ctxt.Andptr[1:] /* xchg rhs,bx */
				asmando(ctxt, p, &p.To, reg[REG_BX])
			} else {
				ctxt.Andptr[0] = byte(0x90 + reg[z])
				ctxt.Andptr = ctxt.Andptr[1:] /* xchg rsh,ax */
				subreg(&pp, z, REG_AX)
				doasm(ctxt, &pp)
				ctxt.Andptr[0] = byte(0x90 + reg[z])
				ctxt.Andptr = ctxt.Andptr[1:] /* xchg rsh,ax */
			}
			return
		}
	}

	ctxt.Diag("doasm: notfound ft=%d tt=%d %v %d %d", p.Ft, p.Tt, p, oclass(ctxt, p, &p.From), oclass(ctxt, p, &p.To))
	return
}

// byteswapreg returns a byte-addressable register (AX, BX, CX, DX)
// which is not referenced in a.
// If a is empty, it returns BX to account for MULB-like instructions
// that might use DX and AX.
func byteswapreg(ctxt *obj.Link, a *obj.Addr) objabi.RBaseType {
	cand := 1
	canc := cand
	canb := canc
	cana := canb

	if a.Type == obj.TYPE_NONE {
		cand = 0
		cana = cand
	}

	if a.Type == obj.TYPE_REG || ((a.Type == obj.TYPE_MEM || a.Type == obj.TYPE_ADDR) && a.Name == obj.NAME_NONE) {
		switch a.Reg {
		case REG_NONE:
			cand = 0
			cana = cand

		case REG_AX, REG_AL, REG_AH:
			cana = 0

		case REG_BX, REG_BL, REG_BH:
			canb = 0

		case REG_CX, REG_CL, REG_CH:
			canc = 0

		case REG_DX, REG_DL, REG_DH:
			cand = 0
		}
	}

	if a.Type == obj.TYPE_MEM || a.Type == obj.TYPE_ADDR {
		switch objabi.RBaseType(a.Index) {
		case REG_AX:
			cana = 0

		case REG_BX:
			canb = 0

		case REG_CX:
			canc = 0

		case REG_DX:
			cand = 0
		}
	}

	if cana != 0 {
		return REG_AX
	}
	if canb != 0 {
		return REG_BX
	}
	if canc != 0 {
		return REG_CX
	}
	if cand != 0 {
		return REG_DX
	}

	ctxt.Diag("impossible byte register")
	log.Fatalf("bad code")
	return 0
}

func isbadbyte(a *obj.Addr) bool {
	return a.Type == obj.TYPE_REG && (REG_BP <= a.Reg && a.Reg <= REG_DI || REG_BPB <= a.Reg && a.Reg <= REG_DIB)
}

var naclret = []uint8{
	0x5e, // POPL SI
	// 0x8b, 0x7d, 0x00, // MOVL (BP), DI - catch return to invalid address, for debugging
	0x83,
	0xe6,
	0xe0, // ANDL $~31, SI
	0x4c,
	0x01,
	0xfe, // ADDQ R15, SI
	0xff,
	0xe6, // JMP SI
}

var naclret8 = []uint8{
	0x5d, // POPL BP
	// 0x8b, 0x7d, 0x00, // MOVL (BP), DI - catch return to invalid address, for debugging
	0x83,
	0xe5,
	0xe0, // ANDL $~31, BP
	0xff,
	0xe5, // JMP BP
}

var naclspfix = []uint8{0x4c, 0x01, 0xfc} // ADDQ R15, SP

var naclbpfix = []uint8{0x4c, 0x01, 0xfd} // ADDQ R15, BP

var naclmovs = []uint8{
	0x89,
	0xf6, // MOVL SI, SI
	0x49,
	0x8d,
	0x34,
	0x37, // LEAQ (R15)(SI*1), SI
	0x89,
	0xff, // MOVL DI, DI
	0x49,
	0x8d,
	0x3c,
	0x3f, // LEAQ (R15)(DI*1), DI
}

var naclstos = []uint8{
	0x89,
	0xff, // MOVL DI, DI
	0x49,
	0x8d,
	0x3c,
	0x3f, // LEAQ (R15)(DI*1), DI
}

func nacltrunc(ctxt *obj.Link, reg objabi.RBaseType) {
	if reg >= REG_R8 {
		ctxt.Andptr[0] = 0x45
		ctxt.Andptr = ctxt.Andptr[1:]
	}
	reg = (reg - REG_AX) & 7
	ctxt.Andptr[0] = 0x89
	ctxt.Andptr = ctxt.Andptr[1:]
	ctxt.Andptr[0] = byte(3<<6 | reg<<3 | reg)
	ctxt.Andptr = ctxt.Andptr[1:]
}

func asmins(ctxt *obj.Link, p *obj.Prog) {
	ctxt.Andptr = ctxt.And[:]
	ctxt.Asmode = int(p.Mode)

	if p.As == objabi.AUSEFIELD {
		r := obj_Addrel(ctxt.Cursym)
		r.Off = 0
		r.Siz = 0
		r.Sym = p.From.Sym
		r.Type = obj.R_USEFIELD
		return
	}

	ctxt.Rexflag = 0
	and0 := ctxt.Andptr
	ctxt.Asmode = int(p.Mode)
	doasm(ctxt, p)
	if ctxt.Rexflag != 0 {
		/*
		 * as befits the whole approach of the architecture,
		 * the rex prefix must appear before the first opcode byte
		 * (and thus after any 66/67/f2/f3/26/2e/3e prefix bytes, but
		 * before the 0f opcode escape!), or it might be ignored.
		 * note that the handbook often misleadingly shows 66/f2/f3 in `opcode'.
		 */
		if p.Mode != 64 {
			ctxt.Diag("asmins: illegal in mode %d: %v (%d %d)", p.Mode, p, p.Ft, p.Tt)
		}
		n := -cap(ctxt.Andptr) + cap(and0)
		var c int
		var np int
		for np = 0; np < n; np++ {
			c = int(and0[np])
			if c != 0xf2 && c != 0xf3 && (c < 0x64 || c > 0x67) && c != 0x2e && c != 0x3e && c != 0x26 {
				break
			}
		}

		copy(and0[np+1:], and0[np:n])
		and0[np] = byte(0x40 | ctxt.Rexflag)
		ctxt.Andptr = ctxt.Andptr[1:]
	}

	n := -cap(ctxt.Andptr) + cap(ctxt.And[:])
	var r *obj.Reloc
	for i := len(ctxt.Cursym.R) - 1; i >= 0; i-- {
		r = &ctxt.Cursym.R[i:][0]
		if int64(r.Off) < p.Pc {
			break
		}
		if ctxt.Rexflag != 0 {
			r.Off++
		}
		if r.Type == obj.R_PCREL {
			// PC-relative addressing is relative to the end of the instruction,
			// but the relocations applied by the linker are relative to the end
			// of the relocation. Because immediate instruction
			// arguments can follow the PC-relative memory reference in the
			// instruction encoding, the two may not coincide. In this case,
			// adjust addend so that linker can keep relocating relative to the
			// end of the relocation.
			r.Add -= p.Pc + int64(n) - (int64(r.Off) + int64(r.Siz))
		}
	}
}
