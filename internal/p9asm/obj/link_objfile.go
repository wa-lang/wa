// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Writing of Wa object files.
//
// Originally, Wa object files were Plan 9 object files, but no longer.
// Now they are more like standard object files, in that each symbol is defined
// by an associated memory image (bytes) and a list of relocations to apply
// during linking. We do not (yet?) use a standard file format, however.
// For now, the format is chosen to be as simple as possible to read and write.
// It may change for reasons of efficiency, or we may even switch to a
// standard file format if there are compelling benefits to doing so.
//
// The file format is:
//
//	- magic header: "\x00\x00wa01ld"
//	- byte 1 - version number
//	- sequence of strings giving dependencies (imported packages)
//	- empty string (marks end of sequence)
//	- sequence of defined symbols
//	- byte 0xff (marks end of sequence)
//	- magic footer: "\xff\xffwa01ld"
//
// All integers are stored in a zigzag varint format.
// Data blocks and strings are both stored as an integer
// followed by that many bytes.
//
// A symbol reference is a string name followed by a version.
// An empty name corresponds to a nil LSym* pointer.
//
// Each symbol is laid out as the following fields (taken from LSym*):
//
//	- byte 0xfe (sanity check for synchronization)
//	- type [int]
//	- name [string]
//	- version [int]
//	- flags [int]
//		1 dupok
//	- size [int]
//	- watype [symbol reference]
//	- p [data block]
//	- nr [int]
//	- r [nr relocations, sorted by off]
//
// If type == STEXT, there are a few more fields:
//
//	- args [int]
//	- locals [int]
//	- nosplit [int]
//	- nlocal [int]
//	- local [nlocal automatics]
//	- pcln [pcln table]
//
// Each relocation has the encoding:
//
//	- off [int]
//	- siz [int]
//	- type [int]
//	- add [int]
//	- xadd [int]
//	- sym [symbol reference]
//	- xsym [symbol reference]
//
// Each local has the encoding:
//
//	- asym [symbol reference]
//	- offset [int]
//	- type [int]
//	- watype [symbol reference]
//
// The pcln table has the encoding:
//
//	- pcsp [data block]
//	- pcfile [data block]
//	- pcline [data block]
//	- npcdata [int]
//	- pcdata [npcdata data blocks]
//	- nfuncdata [int]
//	- funcdata [nfuncdata symbol references]
//	- funcdatasym [nfuncdata ints]
//	- nfile [int]
//	- file [nfile symbol references]
//
// The file layout and meaning of type integers are architecture-independent.
//
// TODO(chai2010): The file format is good for a first pass but needs work.
//	- There are SymID in the object file that should really just be strings.
//	- The actual symbol memory images are interlaced with the symbol
//	  metadata. They should be separated, to reduce the I/O required to
//	  load just the metadata.
//	- The symbol references should be shortened, either with a symbol
//	  table or by using a simple backward index to an earlier mentioned symbol.

package obj

import (
	"fmt"
	"io"
	"log"
	"math"
	"strings"
)

// The Wa and C compilers, and the assembler, call writeobj to write
// out a Wa object file.  The linker does not call this; the linker
// does not write out object files.
func (ctxt *Link) Writeobjdirect(b io.Writer) error {
	// Build list of symbols, and assign instructions to lists.
	// Ignore ctxt->plist boundaries. There are no guarantees there,
	// and the C compilers and assemblers just use one big list.

	var text *LSym
	var data *LSym

	var plink *Prog
	var curtext *LSym
	var etext *LSym
	var edata *LSym
	for pl := ctxt.Plist; pl != nil; pl = pl.Link {
		for p := pl.Firstpc; p != nil; p = plink {
			if ctxt.Debugasm != 0 && ctxt.Debugvlog != 0 {
				fmt.Printf("obj: %v\n", p)
			}
			plink = p.Link
			p.Link = nil

			if p.As == AEND {
				continue
			}

			// 表示一个局部变量或函数参数的类型声明(类似调试信息)
			// 会生成一个 Auto 节点, 并加到当前函数的 Autom 链表中
			if p.As == ATYPE {
				// Assume each TYPE instruction describes
				// a different local variable or parameter,
				// so no dedup.
				// Using only the TYPE instructions means
				// that we discard location information about local variables
				// in C and assembly functions; that information is inferred
				// from ordinary references, because there are no TYPE
				// instructions there. Without the type information, gdb can't
				// use the locations, so we don't bother to save them.
				// If something else could use them, we could arrange to
				// preserve them.
				if curtext == nil {
					continue
				}
				a := new(Auto)
				a.Asym = p.From.Sym
				a.Aoffset = int32(p.From.Offset)
				a.Name = p.From.Name
				a.Link = curtext.Autom
				curtext.Autom = a
				continue
			}

			// 全局变量定义
			// 会创建一个 LSym, 加入全局数据链表
			if p.As == AGLOBL {
				s := p.From.Sym
				tmp6 := s.Seenglobl
				s.Seenglobl++
				if tmp6 != 0 {
					fmt.Printf("duplicate %v\n", p)
				}
				if s.Onlist != 0 {
					log.Fatalf("symbol %s listed multiple times", s.Name)
				}
				s.Onlist = 1
				if data == nil {
					data = s
				} else {
					edata.Next = s
				}
				s.Next = nil
				s.Size = p.To.Offset
				if s.Type == 0 || s.Type == SXREF {
					s.Type = SBSS
				}
				flag := int(p.From3.Offset)
				if flag&DUPOK != 0 {
					s.Dupok = 1
				}
				if flag&RODATA != 0 {
					s.Type = SRODATA
				} else if flag&NOPTR != 0 {
					s.Type = SNOPTRBSS
				}
				edata = s
				continue
			}

			// 全局数据赋值语句
			// 会调用 savedata() 保存数据内容到 LSym.P 中
			if p.As == ADATA {
				if err := ctxt.savedata(p.From.Sym, p, "<input>"); err != nil {
					return err
				}
				continue
			}

			// 函数入口
			// 创建并初始化一个 LSym, 标记为 STEXT, 建立指令链表头尾
			if p.As == ATEXT {
				s := p.From.Sym
				if s == nil {
					// func _() { }
					curtext = nil

					continue
				}

				if s.Text != nil {
					log.Fatalf("duplicate TEXT for %s", s.Name)
				}
				if s.Onlist != 0 {
					log.Fatalf("symbol %s listed multiple times", s.Name)
				}
				s.Onlist = 1
				if text == nil {
					text = s
				} else {
					etext.Next = s
				}
				etext = s
				flag := int(p.From3Offset())
				if flag&DUPOK != 0 {
					s.Dupok = 1
				}
				s.Next = nil
				s.Type = STEXT
				s.Text = p
				s.Etext = p
				curtext = s
				continue
			}

			if p.As == AFUNCDATA {
				// Rewrite reference to wa_args_stackmap(SB) to the Wa-provided declaration information.
				if curtext == nil { // func _() {}
					continue
				}
				if p.To.Sym.Name == "wa_args_stackmap" {
					if p.From.Type != TYPE_CONST || p.From.Offset != FUNCDATA_ArgsPointerMaps {
						return fmt.Errorf("FUNCDATA use of wa_args_stackmap(SB) without FUNCDATA_ArgsPointerMaps")
					}
					p.To.Sym = ctxt.Lookup(fmt.Sprintf("%s.args_stackmap", curtext.Name), curtext.Version)
				}
			}

			if curtext == nil {
				continue
			}
			s := curtext
			s.Etext.Link = p
			s.Etext = p
		}
	}

	// 确保所有函数都有 FUNCDATA wa_args_stackmap 节点
	// GC 需要知道函数的参数指针布局(用于栈扫描)
	// 如果发现某函数缺少这个 FUNCDATA 条目，就补一个伪节点：

	var found int
	for s := text; s != nil; s = s.Next {
		if !strings.HasPrefix(s.Name, "\"\".") {
			continue
		}
		found = 0
		for p := s.Text; p != nil; p = p.Link {
			if p.As == AFUNCDATA && p.From.Type == TYPE_CONST && p.From.Offset == FUNCDATA_ArgsPointerMaps {
				found = 1
				break
			}
		}

		if found == 0 {
			p := NewProg(ctxt)
			p.Link, s.Text.Link = s.Text.Link, p
			p.Lineno = s.Text.Lineno
			p.Mode = s.Text.Mode
			p.As = AFUNCDATA
			p.From.Type = TYPE_CONST
			p.From.Offset = FUNCDATA_ArgsPointerMaps
			p.To.Type = TYPE_MEM
			p.To.Name = NAME_EXTERN
			p.To.Sym = ctxt.Lookup(fmt.Sprintf("%s.args_stackmap", s.Name), s.Version)
		}
	}

	// 生成最终的机器码
	// 调用 Arch 的各个处理函数
	for s := text; s != nil; s = s.Next {
		// 初始化跳转伪指令的目标修复
		s.mkfwd()

		// 补丁处理, 修正重定位项(如 call/jump 地址)
		if err := s.linkpatch(ctxt); err != nil {
			return err
		}

		// 架构相关, 可能展开指令
		ctxt.Arch.Follow(ctxt, s)

		// 架构相关, 预处理代码块(如压缩指令、编码优化)
		ctxt.Arch.Preprocess(ctxt, s)

		// 真正将指令编译为机器码, 写入 LSym.P
		ctxt.Arch.Assemble(ctxt, s)

		// 生成行号调试信息, 构建 LSym.Pcln
		s.linkpcln(ctxt)
	}

	// 写 waobj 开始标志
	b.Write([]byte(MagicHeader))

	// 版本信息
	b.Write([]byte{1})

	// 写 import 列表
	for _, pkg := range ctxt.Imports {
		ctxt.wrstring(b, pkg)
	}

	// 以空字符串作为 import 列表结尾
	ctxt.wrstring(b, "")

	// 生成函数标识符
	for s := text; s != nil; s = s.Next {
		ctxt.writesym(b, s)
	}
	// 生成变量标识符
	for s := data; s != nil; s = s.Next {
		ctxt.writesym(b, s)
	}

	// 写 waobj 结束标志
	b.Write([]byte(MagicFooter))

	return nil
}

// 写具名的对象
func (ctxt *Link) writesym(b io.Writer, s *LSym) {
	// 打印文本格式的汇编信息
	if ctxt.Debugasm != 0 {
		// main 函数打印格式
		// main.main t=1 nosplit size=32 value=0 args=0xffffffff80000000 locals=0x0

		fmt.Fprintf(ctxt.Bso, "%s ", s.Name)
		if s.Version != 0 {
			fmt.Fprintf(ctxt.Bso, "v=%d ", s.Version)
		}
		if s.Type != 0 {
			fmt.Fprintf(ctxt.Bso, "t=%d ", s.Type)
		}

		// 打印 flags 标志(仅针对数据)
		if s.Dupok != 0 {
			fmt.Fprintf(ctxt.Bso, "dupok ")
		}

		// 打印内存大小和值
		fmt.Fprintf(ctxt.Bso, "size=%d value=%d", int64(s.Size), int64(s.Value))

		// 如果是函数类型
		// 打印参数数目和局部变量数目, 是否叶子函数等信息
		if s.Type == STEXT {
			fmt.Fprintf(ctxt.Bso, " args=%#x locals=%#x", uint64(s.Args), uint64(s.Locals))
		}
		fmt.Fprintf(ctxt.Bso, "\n")

		// 打印 pcln 表, 每行对应汇编指令
		//
		// 类似以下的格式:
		// 0x0000 00000 (textflag.h:25)    TEXT    main.main(SB), 4, $0
		// 0x0000 00000 (textflag.h:25)    NOP
		// 0x0000 00000 (textflag.h:25)    NOP
		// 0x0000 00000 (textflag.h:26)    MOVQ    $60, AX
		// 0x0007 00007 (textflag.h:27)    MOVQ    $42, DI

		for p := s.Text; p != nil; p = p.Link {
			fmt.Fprintf(ctxt.Bso, "\t%#04x %v\n", uint(int(p.Pc)), p)
		}

		// 符号对应的原始内容, 最终写入目标文件
		// 如果是函数, 则对应汇编指令的机器码
		// 如果是变量, 则是其中的值内容
		//
		// 每次处理 16 个字节, 类似以下的格式:
		// 0x0000 48 c7 c0 3c 00 00 00 48 c7 c7 2a 00 00 00 0f 05  H..<...H..*.....
		// 0x0010 c3

		for i := 0; i < len(s.P); i += 16 {
			fmt.Fprintf(ctxt.Bso, "\t%#04x", uint(i))

			var j int
			for j = i; j < i+16 && j < len(s.P); j++ {
				fmt.Fprintf(ctxt.Bso, " %02x", s.P[j])
			}

			// 不足 16 字节补充空格
			for ; j < i+16; j++ {
				fmt.Fprintf(ctxt.Bso, "   ")
			}

			// 16 字节数据对应的 ASCII 码形式打印
			// 不能识别的符号以 "." 展示

			fmt.Fprintf(ctxt.Bso, "  ")
			for j := i; j < i+16 && j < len(s.P); j++ {
				c := int(s.P[j])
				if ' ' <= c && c <= 0x7e {
					fmt.Fprintf(ctxt.Bso, "%c", c)
				} else {
					fmt.Fprintf(ctxt.Bso, ".")
				}
			}

			fmt.Fprintf(ctxt.Bso, "\n")
		}

		// 符号内的重定位项列表(地址引用, 符号引用等)
		for i := 0; i < len(s.R); i++ {
			r := &s.R[i]
			var name string
			if r.Sym != nil {
				name = r.Sym.Name
			}

			// 对应的是 `s.P[r.Off:][:r.Siz]` 中的位置, 最终需要根据引用的标识符名字计算出地址回填
			if ctxt.Arch.Thechar == '5' || ctxt.Arch.Thechar == '9' {
				fmt.Fprintf(ctxt.Bso, "\trel %d+%d t=%d %s+%x\n", int(r.Off), r.Siz, r.Type, name, uint64(int64(r.Add)))
			} else {
				fmt.Fprintf(ctxt.Bso, "\trel %d+%d t=%d %s+%d\n", int(r.Off), r.Siz, r.Type, name, int64(r.Add))
			}
		}
	}

	// 写标识符的开始标志
	b.Write([]byte{MagicSymbolStart})

	// 标识符类型
	ctxt.wrint(b, int64(s.Type))

	// 标识符名字
	ctxt.wrstring(b, s.Name)

	// 标识符版本
	ctxt.wrint(b, int64(s.Version))

	// 写 flags 信息
	flags := int64(s.Dupok)
	if s.Local {
		flags |= 2
	}
	ctxt.wrint(b, flags)

	// 对应的内存大小
	ctxt.wrint(b, s.Size)

	// 引用的类型符号
	ctxt.wrsym(b, s.Watype)

	// 二进制数据
	// 如果是函数, 则对应机器码数据
	// 如果是变量, 则对应内存的值
	ctxt.wrdata(b, s.P)

	// 重定位信息
	ctxt.wrint(b, int64(len(s.R)))
	for i := 0; i < len(s.R); i++ {
		r := &s.R[i]
		ctxt.wrint(b, int64(r.Off))
		ctxt.wrint(b, int64(r.Siz))
		ctxt.wrint(b, int64(r.Type))
		ctxt.wrint(b, r.Add)
		ctxt.wrint(b, 0) // Xadd, ignored
		ctxt.wrsym(b, r.Sym)
		ctxt.wrsym(b, nil) // Xsym, ignored
	}

	// 如果是函数
	if s.Type == STEXT {
		// 写参数/局部变量
		ctxt.wrint(b, int64(s.Args))
		ctxt.wrint(b, int64(s.Locals))

		// 自动变量(和locals的区别?)
		n := 0
		for a := s.Autom; a != nil; a = a.Link {
			n++
		}
		ctxt.wrint(b, int64(n))
		for a := s.Autom; a != nil; a = a.Link {
			// 标识符
			ctxt.wrsym(b, a.Asym)
			ctxt.wrint(b, int64(a.Aoffset))

			// 类别
			switch a.Name {
			case NAME_AUTO:
				ctxt.wrint(b, A_AUTO)
			case NAME_PARAM:
				ctxt.wrint(b, A_PARAM)
			default:
				log.Fatalf("%s: invalid local variable type %d", s.Name, a.Name)
			}

			// 对应类型的标识符
			ctxt.wrsym(b, a.Watype)
		}

		// 写 pcln 数据
		ctxt.wrdata(b, s.Pcln.Pcsp.P)
		ctxt.wrdata(b, s.Pcln.Pcfile.P)
		ctxt.wrdata(b, s.Pcln.Pcline.P)
		ctxt.wrint(b, int64(len(s.Pcln.Pcdata)))
		for i := 0; i < len(s.Pcln.Pcdata); i++ {
			ctxt.wrdata(b, s.Pcln.Pcdata[i].P)
		}
		ctxt.wrint(b, int64(len(s.Pcln.Funcdataoff)))
		for i := 0; i < len(s.Pcln.Funcdataoff); i++ {
			ctxt.wrsym(b, s.Pcln.Funcdata[i])
		}
		for i := 0; i < len(s.Pcln.Funcdataoff); i++ {
			ctxt.wrint(b, s.Pcln.Funcdataoff[i])
		}
		ctxt.wrint(b, int64(len(s.Pcln.File)))
		for i := 0; i < len(s.Pcln.File); i++ {
			// TODO(chai2010): 修改了格式, 需要验证
			ctxt.wrstring(b, s.Pcln.File[i])
		}
	}
}

func (ctxt *Link) savedata(s *LSym, p *Prog, pn string) error {
	off := int32(p.From.Offset)
	siz := int32(p.From3.Offset)
	if off < 0 || siz < 0 || off >= 1<<30 || siz >= 100 {
		return fmt.Errorf("%s: mangled input file", pn)
	}
	if ctxt.Enforce_data_order != 0 && off < int32(len(s.P)) {
		return fmt.Errorf("data out of order (already have %d)\n%v", len(s.P), p)
	}

	// 调整符号对应的机器码切片的容量
	if len(s.P) < int(off+siz) {
		s.P = append(s.P, make([]byte, int(off+siz)-len(s.P))...)
	}

	switch p.To.Type {
	default:
		return fmt.Errorf("bad data: %v", p)

	case TYPE_FCONST:
		switch siz {
		default:
			return fmt.Errorf("unexpected %d-byte floating point constant", siz)

		case 4:
			flt := math.Float32bits(float32(p.To.Val.(float64)))
			ctxt.Arch.ByteOrder.PutUint32(s.P[off:], flt)

		case 8:
			flt := math.Float64bits(p.To.Val.(float64))
			ctxt.Arch.ByteOrder.PutUint64(s.P[off:], flt)
		}

	case TYPE_SCONST:
		copy(s.P[off:off+siz], p.To.Val.(string))

	case TYPE_CONST, TYPE_ADDR:
		if p.To.Sym != nil || p.To.Type == TYPE_ADDR {
			s.R = append(s.R, Reloc{
				Off:  off,
				Siz:  uint8(siz),
				Sym:  p.To.Sym,
				Type: R_ADDR,
				Add:  p.To.Offset,
			})
			break
		}
		o := p.To.Offset
		switch siz {
		default:
			return fmt.Errorf("unexpected %d-byte integer constant", siz)
		case 1:
			s.P[off] = byte(o)
		case 2:
			ctxt.Arch.ByteOrder.PutUint16(s.P[off:], uint16(o))
		case 4:
			ctxt.Arch.ByteOrder.PutUint32(s.P[off:], uint32(o))
		case 8:
			ctxt.Arch.ByteOrder.PutUint64(s.P[off:], uint64(o))
		}
	}

	return nil
}

func (ctxt *Link) wrint(b io.Writer, sval int64) {
	var varintbuf [10]uint8
	var v uint64
	uv := (uint64(sval) << 1) ^ uint64(int64(sval>>63))
	p := varintbuf[:]
	for v = uv; v >= 0x80; v >>= 7 {
		p[0] = uint8(v | 0x80)
		p = p[1:]
	}
	p[0] = uint8(v)
	p = p[1:]
	b.Write(varintbuf[:len(varintbuf)-len(p)])
}

func (ctxt *Link) wrstring(b io.Writer, s string) {
	ctxt.wrint(b, int64(len(s)))
	b.Write([]byte(s))
}

func (ctxt *Link) wrdata(b io.Writer, v []byte) {
	ctxt.wrint(b, int64(len(v)))
	b.Write(v)
}

func (ctxt *Link) wrsym(b io.Writer, s *LSym) {
	if s == nil {
		ctxt.wrint(b, 0)
		ctxt.wrint(b, 0)
		return
	}

	ctxt.wrstring(b, s.Name)
	ctxt.wrint(b, int64(s.Version))
}
