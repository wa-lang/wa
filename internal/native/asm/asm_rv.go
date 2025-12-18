// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"encoding/binary"
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/parser"
	"wa-lang.org/wa/internal/native/pcrel"
	"wa-lang.org/wa/internal/native/riscv"
)

func (p *_Assembler) asmFile_riscv(filename string, source []byte, opt *abi.LinkOptions) (prog *abi.LinkedProgram, err error) {
	// 解析汇编程序
	p.file, err = parser.ParseFile(opt.CPU, p.fset, filename, source)
	if err != nil {
		return nil, err
	}

	// 全局函数分配内存空间
	textAddrStart := p.dramNextAddr
	for _, fn := range p.file.Funcs {
		fn.Size = int(p.funcBodyLen(fn))
		fn.LinkInfo = &abi.LinkedSymbol{
			Name: fn.Name,
			Addr: p.alloc(int64(fn.Size), 0),
			Data: make([]byte, fn.Size),
		}
	}

	// 龙芯数据段从下个页面开始
	dataAddrStart := p.dramNextAddr

	// 全局变量分配内存空间
	for _, g := range p.file.Globals {
		assert(g.Size > 0)
		g.LinkInfo = &abi.LinkedSymbol{
			Name: g.Name,
			Addr: p.alloc(int64(g.Size), 0),
			Data: make([]byte, g.Size),
		}
	}

	// 编译函数
	for _, fn := range p.file.Funcs {
		if err := p.asmFunc_riscv(fn); err != nil {
			return nil, err
		}
	}

	// 编译全局变量
	for _, g := range p.file.Globals {
		if err := p.asmGlobal(g); err != nil {
			return nil, err
		}
	}

	// 收集全部信息
	{
		p.prog.TextAddr = textAddrStart
		p.prog.DataAddr = dataAddrStart

		// 优先查找指定的入口函数
		if opt.EntryFunc != "" {
			for _, fn := range p.file.Funcs {
				if fn.Name == opt.EntryFunc {
					p.prog.Entry = fn.LinkInfo.Addr
				}
			}
		}

		// 然后查找默认的入口函数(中文)
		if p.prog.Entry == 0 {
			for _, fn := range p.file.Funcs {
				if fn.Name == abi.DefaultEntryFuncZh {
					p.prog.Entry = fn.LinkInfo.Addr
				}
			}
		}

		// 然后查找默认的入口函数(英文)
		if p.prog.Entry == 0 {
			for _, fn := range p.file.Funcs {
				if fn.Name == abi.DefaultEntryFunc {
					p.prog.Entry = fn.LinkInfo.Addr
				}
			}
		}
		if p.prog.Entry == 0 {
			p.prog.Entry = textAddrStart
		}

		// text 段数据
		p.prog.TextData = nil
		for _, fn := range p.file.Funcs {
			p.prog.TextData = append(p.prog.TextData, fn.LinkInfo.Data...)
		}

		// data 段数据
		p.prog.DataData = nil
		for _, g := range p.file.Globals {
			p.prog.DataData = append(p.prog.DataData, g.LinkInfo.Data...)
		}
	}

	return p.prog, nil
}

func (p *_Assembler) asmFunc_riscv(fn *ast.Func) (err error) {
	// label 的地址列表
	label2pcMap := make(map[string]int64)

	// 绝对地址拆分基于 symbol 名字
	hi2loMap := make(map[string]int32)

	// 相对地址拆分基于 label 名字
	// 一个 label 只能有一个 symbol 相对 PC 寻址
	label2loMap := make(map[string]int32)

	// 第一遍收集全部 label, 因为可能向前跳转没有出现的 label
	pc := fn.LinkInfo.Addr
	for _, inst := range fn.Body.Insts {
		if inst.Label != "" {
			if _, ok := label2pcMap[inst.Label]; ok {
				panic(fmt.Errorf("label %q exists", inst.Label))
			}
			label2pcMap[inst.Label] = pc
		}
		if inst.As == 0 {
			// 跳过空的指令, 比如标号
			continue
		}

		// 更新下一个指令对应的 pc 位置
		pc += p.instLen(inst)
	}

	// 第二遍遍历编码指令
	pc = fn.LinkInfo.Addr
	for i, inst := range fn.Body.Insts {
		if inst.As == 0 {
			// 跳过空的指令, 比如标号
			continue
		}

		// 指令对 label 或全局的符号引用
		// 因为指令长度的关系, 指令并不会直接访问符号对应的绝对地址
		// 需要解决 %hi/%lo/%pcrel_hi/%pcrel_lo 等转化为最终可编码到指令的值
		if inst.Arg.Symbol != "" {
			addr, ok := label2pcMap[inst.Arg.Symbol]
			if !ok {
				addr, ok = p.symbolAddress(inst.Arg.Symbol)
				if !ok {
					panic(fmt.Errorf("symbol %q not found", inst.Arg.Symbol))
				}
			}

			if inst.Arg.SymbolDecor != abi.BuiltinFn_Nil {
				assert(inst.Arg.SymbolDecor.IsValid(p.prog.CPU))
			}

			switch inst.Arg.SymbolDecor {
			case abi.BuiltinFn_HI, abi.BuiltinFn_HI_zh: // 高20bit
				hi, lo := pcrel.MakeAbs(uint32(addr))
				hi2loMap[inst.Arg.Symbol] = lo
				inst.Arg.Imm = hi

			case abi.BuiltinFn_LO, abi.BuiltinFn_LO_zh: // 低12bit
				lo, ok := hi2loMap[inst.Arg.Symbol]
				if !ok {
					panic(fmt.Errorf("symbol %q not found", inst.Arg.Symbol))
				}
				inst.Arg.Imm = lo

			case abi.BuiltinFn_PCREL_HI, abi.BuiltinFn_PCREL_HI_zh:
				var labelName string
				for k := i; k >= 0; k-- {
					if x := fn.Body.Insts[k]; x.Label != "" {
						labelName = x.Label
						break
					}
				}
				if labelName == "" {
					panic(fmt.Errorf("symbol %q not found found", inst.Arg.Symbol))
				}

				hi, lo := pcrel.MakePCRel(addr, pc)
				label2loMap[labelName] = lo
				inst.Arg.Imm = hi

			case abi.BuiltinFn_PCREL_LO, abi.BuiltinFn_PCREL_LO_zh:
				// https://sourceware.org/binutils/docs/as/RISC_002dV_002dModifiers.html
				// https://stackoverflow.com/questions/65879012/what-do-pcrel-hi-and-pcrel-lo-actually-do

				lo, ok := label2loMap[inst.Arg.Symbol]
				if !ok {
					panic(fmt.Errorf("symbol %q not found", inst.Arg.Symbol))
				}
				inst.Arg.Imm = lo

			case abi.BuiltinFn_SIZEOF, abi.BuiltinFn_SIZEOF_zh:
				var g *ast.Global
				for _, x := range p.file.Globals {
					if x.Name == inst.Arg.Symbol {
						g = x
						break
					}
				}
				if g == nil {
					panic(fmt.Errorf("global %q not found", inst.Arg.Symbol))
				}
				inst.Arg.Imm = int32(g.Size)

			default:
				assert(inst.Arg.SymbolDecor == abi.BuiltinFn_Nil)

				// 如果是label, 则为 pc 相对地址
				if _, isLabel := label2pcMap[inst.Arg.Symbol]; isLabel {
					inst.Arg.Imm = int32(addr - pc)
				} else {
					inst.Arg.Imm = int32(addr)
				}
			}
		}

		// 编码使用的是符号被处理后对应的立即数
		var x uint32
		switch p.opt.CPU {
		case abi.RISCV32:
			x, err = riscv.EncodeRV32(inst.As, inst.Arg)
		case abi.RISCV64:
			x, err = riscv.EncodeRV64(inst.As, inst.Arg)
		default:
			panic("unreachable")
		}
		if err != nil {
			return fmt.Errorf("%v: %w", inst, err)
		}

		// 保存指令编码后的机器码
		binary.LittleEndian.PutUint32(
			fn.LinkInfo.Data[int(pc-fn.LinkInfo.Addr):],
			x,
		)

		// 更新下一个指令对应的 pc 位置
		pc += p.instLen(inst)
	}

	return nil
}
