// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/link/elf"
	"wa-lang.org/wa/internal/native/parser"
	"wa-lang.org/wa/internal/native/token"
)

func (p *_Assembler) asmFile(filename string, source []byte, opt *abi.LinkOptions) (prog *abi.LinkedProgram, err error) {
	p.init(filename, source, opt)

	// 最大的页大小
	const maxPageSize = 64 << 10

	// get source
	if source == nil {
		source, err = os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
	}

	// 解析汇编程序
	p.file, err = parser.ParseFile(opt.CPU, p.fset, filename, source)
	if err != nil {
		return nil, err
	}

	// 指令段的地址必须页对齐
	p.prog.TextAddr = p.dramNextAddr
	assert(p.prog.TextAddr%maxPageSize == 0)

	// 给 ELF 头和 程序头预留出空间
	// 如果有调试用的段数据, 放到后面
	//
	// gcc 会生成4个程序头, 和 .note.gnu.property 段
	// 有效的 text 数据可能从 0x144 附近开始
	const fileHeaderSize = elf.ELF64HDRSIZE + elf.ELF64PHDRSIZE*2
	p.prog.TextData = make([]byte, fileHeaderSize)
	assert(len(p.prog.TextData) == fileHeaderSize)
	p.dramNextAddr += int64(len(p.prog.TextData))

	// 全局函数分配内存空间
	for _, fn := range p.file.Funcs {
		if _, ok := p.objectMap[fn.Name]; ok {
			panic(fmt.Sprintf("object %s exists", fn.Name))
		}
		p.objectMap[fn.Name] = fn

		fn.BodySize = int(p.funcBodyLen(fn))
		fn.LinkInfo = &abi.LinkedSymbol{Name: fn.Name}
		fn.LinkInfo.Addr, fn.LinkInfo.AlignPadding = p.alloc(int64(fn.BodySize), 4)
		fn.LinkInfo.Data = make([]byte, fn.BodySize)
	}

	// 数据段从下个页面开始
	// 但是文件偏移量不变, 确保文件被映射到内存后依然是相似的布局
	p.dramNextAddr += maxPageSize
	p.prog.DataAddr = p.dramNextAddr

	// 全局变量分配内存空间
	for _, g := range p.file.Globals {
		if _, ok := p.objectMap[g.Name]; ok {
			panic(fmt.Sprintf("object %s exists", g.Name))
		}
		p.objectMap[g.Name] = g

		g.LinkInfo = &abi.LinkedSymbol{
			Name: g.Name,
		}

		switch g.TypeTok {
		case token.BYTE_zh:
			assert(g.Type == token.I8)
			assert(g.Size == 1)
			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 1)
			g.LinkInfo.Data = make([]byte, g.Size)
		case token.SHORT_zh:
			assert(g.Type == token.I16)
			assert(g.Size == 2)
			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 2)
			g.LinkInfo.Data = make([]byte, g.Size)
		case token.LONG_zh:
			assert(g.Type == token.I32)
			assert(g.Size == 4)
			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 4)
			g.LinkInfo.Data = make([]byte, g.Size)
		case token.QUAD_zh:
			assert(g.Type == token.I64)
			assert(g.Size == 8)
			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 8)
			g.LinkInfo.Data = make([]byte, g.Size)
		case token.FLOAT_zh:
			assert(g.Type == token.F32)
			assert(g.Size == 4)
			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 4)
			g.LinkInfo.Data = make([]byte, g.Size)
		case token.DOUBLE_zh:
			assert(g.Type == token.F64)
			assert(g.Size == 8)
			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 8)
			g.LinkInfo.Data = make([]byte, g.Size)
		case token.ADDR_zh:
			if p.file.CPU == abi.RISCV32 {
				assert(g.Type == token.I32)
				assert(g.Size == 4)
				g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 4)
				g.LinkInfo.Data = make([]byte, g.Size)
			} else {
				assert(g.Type == token.I64)
				assert(g.Size == 8)
				g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 8)
				g.LinkInfo.Data = make([]byte, g.Size)
			}
		case token.ASCII_zh:
			assert(g.Type == token.Bin)
			s := g.Init.Lit.ConstV.(string)
			assert(g.Size == len(s))
			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 1)
			g.LinkInfo.Data = make([]byte, g.Size)
		case token.GAS_BYTE:
			assert(g.Type == token.I8)
			assert(g.Size == 1)
			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 1)
			g.LinkInfo.Data = make([]byte, g.Size)
		case token.GAS_SHORT:
			assert(g.Type == token.I16)
			assert(g.Size == 2)
			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 2)
			g.LinkInfo.Data = make([]byte, g.Size)
		case token.GAS_LONG:
			assert(g.Type == token.I32)
			assert(g.Size == 4)
			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 4)
			g.LinkInfo.Data = make([]byte, g.Size)
		case token.GAS_QUAD:
			assert(g.Type == token.I64)
			assert(g.Size == 8)
			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 8)
			g.LinkInfo.Data = make([]byte, g.Size)
		case token.GAS_FLOAT:
			assert(g.Type == token.F32)
			assert(g.Size == 4)
			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 4)
			g.LinkInfo.Data = make([]byte, g.Size)
		case token.GAS_DOUBLE:
			assert(g.Type == token.F64)
			assert(g.Size == 8)
			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 8)
			g.LinkInfo.Data = make([]byte, g.Size)
		case token.GAS_ASCII:
			assert(g.Type == token.Bin)
			s := g.Init.Lit.ConstV.(string)
			assert(g.Size == len(s))
			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 1)
			g.LinkInfo.Data = make([]byte, g.Size)
		case token.GAS_SKIP:
			assert(g.Type == token.Bin)
			n := g.Init.Lit.ConstV.(int64)
			assert(n >= 0)
			assert(g.Size == int(n))
			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 1)
			g.LinkInfo.Data = make([]byte, g.Size)
		case token.GAS_INCBIN:
			assert(g.Type == token.Bin)

			// 检查文件的大小是否超出限制
			filename := g.Init.Lit.ConstV.(string)
			if fi, err := os.Lstat(filename); err != nil {
				panic(fmt.Sprintf("file %s not found", filename))
			} else {
				const maxSize = 2 << 20
				if fi.Size() > maxSize {
					panic(fmt.Sprintf("%v %v file size large than 2MB", token.GAS_INCBIN, filename))
				}
				g.Size = int(fi.Size())
			}

			// 读取文件数据
			data, err := os.ReadFile(filename)
			if err != nil {
				panic(fmt.Sprintf("read file %s failed: %v", filename, err))
			}
			if len(data) != g.Size {
				panic(fmt.Sprintf("read file %s failed: %v", filename, err))
			}

			g.LinkInfo.Addr, g.LinkInfo.AlignPadding = p.alloc(int64(g.Size), 1)
			g.LinkInfo.Data = data

		default:
			panic("unreachable")
		}
	}

	// 编译函数
	switch p.file.CPU {
	case abi.LOONG64:
		for _, fn := range p.file.Funcs {
			if err := p.asmFuncBody_inst_loong64(fn); err != nil {
				return nil, err
			}
		}
	case abi.RISCV32, abi.RISCV64:
		for _, fn := range p.file.Funcs {
			if err := p.asmFuncBody_inst_riscv(fn); err != nil {
				return nil, err
			}
		}

	case abi.X64Unix:
		for _, fn := range p.file.Funcs {
			if err := p.asmFuncBody_inst_x64_linux(fn); err != nil {
				return nil, err
			}
		}

	default:
		panic(fmt.Sprintf("unnsupport cpu: %v", p.file.CPU))
	}

	// 编译全局变量
	for _, g := range p.file.Globals {
		if err := p.asmGlobal(g); err != nil {
			return nil, err
		}
	}

	// 收集全部信息
	{
		// text 段数据(头部空间保留)
		assert(p.prog.Entry == 0)
		for _, fn := range p.file.Funcs {
			assert(fn.LinkInfo.AlignPadding == 0)
			p.prog.TextData = append(p.prog.TextData, fn.LinkInfo.Data...)

			// 设置入口函数
			if fn.Name == abi.DefaultEntryFunc || fn.Name == abi.DefaultEntryFuncZh {
				p.prog.Entry = fn.LinkInfo.Addr
			}
		}

		assert(p.prog.Entry > 0)
		assert((p.prog.Entry - p.prog.TextAddr) >= fileHeaderSize)

		// data 段数据
		assert(len(p.prog.DataData) == 0)
		for _, g := range p.file.Globals {
			for i := 0; i < g.LinkInfo.AlignPadding; i++ {
				p.prog.DataData = append(p.prog.DataData, 0)
			}
			p.prog.DataData = append(p.prog.DataData, g.LinkInfo.Data...)
		}
	}

	return p.prog, nil
}
