// Derived from Inferno utils/6l/obj.c and utils/6l/span.c
// http://code.google.com/p/inferno-os/source/browse/utils/6l/obj.c
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

package obj

import "fmt"

// 对应的是将被写入目标文件的符号, 比如函数/变量/常量/调试信息等
type LSym struct {
	Name      string  // 符号名字
	Type      SymKind // 符号类型
	Version   int16   // 符号版本
	Dupok     uint8   // 是否可重复定义(比如内联函数)
	Seenglobl uint8   // 是否已处理全局定义
	Onlist    uint8   // 是否已加入输出列表
	Local     bool    // 是否为本地符号(模块内可见)
	Args      int32   // 参数大小
	Locals    int32   // 局部变量大小
	Value     int64   // 符号地址或常量值
	Size      int64   // 大小, 通常对于数据符号
	Next      *LSym   // 链表
	Watype    *LSym   // 类型信息对应的符号
	Autom     *Auto   // 自动变量(用于函数)
	Text      *Prog   // 如果是函数, 对应第一个汇编指令
	Etext     *Prog   // 最后一个汇编指令
	Pcln      *Pcln   // PC对应的行号位图信息
	P         []byte  // 指令的机器码, 或者是数据
	R         []Reloc // 重定位表
}

// 描述一个函数的局部变量或临时变量
type Auto struct {
	Asym    *LSym    // 变量名
	Link    *Auto    // 下一个 Auto, 形成链表
	Aoffset int32    // 该变量在栈帧中的偏移量
	Name    AddrName // 变量的类别(参数/局部变量/临时变量)
	Watype  *LSym    // 变量的类型信息
}

type Pcln struct {
	Pcsp        Pcdata
	Pcfile      Pcdata
	Pcline      Pcdata
	Pcdata      []Pcdata
	Funcdata    []*LSym
	Funcdataoff []int64
	File        []*LSym
	Lastfile    *LSym
	Lastindex   int
}

type Reloc struct {
	Off  int32
	Siz  uint8
	Type RelocType
	Add  int64
	Sym  *LSym
}

// Auto.name
const (
	A_AUTO = 1 + iota
	A_PARAM
)

// Reloc.type
// 最终的值是bit位组合
type RelocType int32

const (
	R_ADDR      RelocType = 1 + iota // 表示一个绝对地址, 需要填入的是目标符号的实际地址
	R_ADDRPOWER                      // 特定平台上绝对地址重定位方式的差异
	R_ADDRARM64                      // 特定平台上绝对地址重定位方式的差异。
	R_SIZE
	R_CALL // 常规函数调用
	R_CALLARM
	R_CALLARM64
	R_CALLIND // 间接调用(通过函数指针)
	R_CALLPOWER
	R_CONST // 一个常数值重定位, 不会被链接器修改
	R_PCREL // 表示 PC 相对寻址, 典型用于跳转指令, 比如 call 或 jmp
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
	R_GOTOFF // 全局变量地址与 GOT 表基地址的偏移
	R_PLT0   // 用于动态链接库调用时的跳板(plt stub)入口构建
	R_PLT1
	R_PLT2
	R_USEFIELD
	R_GOTPCREL // 通过 GOT 表(Global Offset Table)进行 PC 相对寻址
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

func (sym *LSym) linkpatch(ctxt *Link) error {
	var c int32
	var name string
	var q *Prog

	ctxt.Cursym = sym

	for p := sym.Text; p != nil; p = p.Link {
		if err := p.From.checkaddr(); err != nil {
			return fmt.Errorf("%w in %v", err, p)
		}
		if p.From3 != nil {
			if err := p.From3.checkaddr(); err != nil {
				return fmt.Errorf("%w in %v", err, p)
			}
		}
		if err := p.To.checkaddr(); err != nil {
			return fmt.Errorf("%w in %v", err, p)
		}

		if ctxt.Arch.Progedit != nil {
			ctxt.Arch.Progedit(ctxt, p)
		}
		if p.To.Type != TYPE_BRANCH {
			continue
		}
		if p.To.Val != nil {
			// TODO: Remove To.Val.(*Prog) in favor of p->pcond.
			p.Pcond = p.To.Val.(*Prog)
			continue
		}

		if p.To.Sym != nil {
			continue
		}
		c = int32(p.To.Offset)
		for q = sym.Text; q != nil; {
			if int64(c) == q.Pc {
				break
			}
			if q.Forwd != nil && int64(c) >= q.Forwd.Pc {
				q = q.Forwd
			} else {
				q = q.Link
			}
		}

		if q == nil {
			name = "<nil>"
			if p.To.Sym != nil {
				name = p.To.Sym.Name
			}
			p.To.Type = TYPE_NONE
			return fmt.Errorf("branch out of range (%#x)\n%v [%s]", uint32(c), p, name)
		}

		p.To.Val = q
		p.Pcond = q
	}

	for p := sym.Text; p != nil; p = p.Link {
		p.Mark = 0 /* initialization for follow */
		if p.Pcond != nil {
			p.Pcond = p.Pcond.brloop()
			if p.Pcond != nil {
				if p.To.Type == TYPE_BRANCH {
					p.To.Offset = p.Pcond.Pc
				}
			}
		}
	}

	return nil
}

// 加速指令链的遍历
// 通过为 sym.Text 链表上的某些节点设置 Forwd 快捷指针来构建 跳跃表 式的结构
func (sym *LSym) mkfwd() {
	// LOG 决定了快捷指针的等级, 类似 skip list 的层级
	const LOG = 5

	var dwn [LOG]int32 // 表示第 i 层的跳跃间隔
	var cnt [LOG]int32 // 当前第 i 层剩余的步数
	var lst [LOG]*Prog // 每一层上上一个被连接的节点

	// 初始化
	for i := 0; i < LOG; i++ {
		if i == 0 {
			cnt[i] = 1
		} else {
			cnt[i] = LOG * cnt[i-1] // 几何级数增长
		}
		dwn[i] = 1 // 初始都触发一次
		lst[i] = nil
	}

	// 构造 Forwd 快捷指针
	// 每次处理一个指令节点 p, 尝试设置若干 Forwd 指针
	// i 每次减一形成一种循环滑动窗口(0~4)
	// 每个 dwn[i] 是该层倒计数器, 等于 0 时触发设置 Forwd
	// 这个结构类似 跳表(skip list), 通过不同间隔建立部分链接, 提高访问效率

	i := 0
	for p := sym.Text; p != nil && p.Link != nil; p = p.Link {
		i--
		if i < 0 {
			i = LOG - 1
		}
		p.Forwd = nil
		dwn[i]--
		if dwn[i] <= 0 {
			dwn[i] = cnt[i] // 重置该层倒计数器
			if lst[i] != nil {
				lst[i].Forwd = p // 上一次节点指向当前
			}
			lst[i] = p // 更新该层最新节点
		}
	}
}
