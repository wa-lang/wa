// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"encoding/binary"
	"fmt"
	"os"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/loong64"
	"wa-lang.org/wa/internal/native/parser"
	"wa-lang.org/wa/internal/native/pcrel"
	"wa-lang.org/wa/internal/token"
	"wa-lang.org/wa/internal/xlang"
)

func (p *_Assembler) asmFile_loong64(filename string, source []byte, opt *abi.LinkOptions) (prog *abi.LinkedProgram, err error) {
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
	xtype := xlang.DetectLang(filename, source)
	switch xtype {
	case token.LangType_Nasm_zh:
		p.file, err = parser.ParseFile(opt.CPU, p.fset, filename, source)
		if err != nil {
			return nil, err
		}
	case token.LangType_Nasm_gas:
		p.file, err = parser.ParseFile(opt.CPU, p.fset, filename, source)
		if err != nil {
			return nil, err
		}
	default:
		err = fmt.Errorf("unsupport file format: %s: %v", filename, xtype)
		return
	}

	// 指令段的地址必须页对齐
	p.prog.TextAddr = p.dramNextAddr
	assert(p.prog.TextAddr%maxPageSize == 0)

	// 给 ELF 头和 程序头预留出空间
	const fileHeaderSize = 64 + 56*2
	p.prog.TextData = make([]byte, 0x10c)
	assert(len(p.prog.TextData) > fileHeaderSize)
	p.dramNextAddr += int64(len(p.prog.TextData))
	p.prog.Entry = p.dramNextAddr

	// 全局函数分配内存空间
	for _, fn := range p.file.Funcs {
		fn.BodySize = int(p.funcBodyLen(fn))
		fn.LinkInfo = &abi.LinkedSymbol{
			Name: fn.Name,
			Addr: p.alloc(int64(fn.BodySize), 0),
			Data: make([]byte, fn.BodySize),
		}
	}

	// 数据段从下个页面开始
	// 但是文件偏移量不变, 确保文件被映射到内存后依然是相似的布局
	p.dramNextAddr += maxPageSize
	p.prog.DataAddr = p.dramNextAddr

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
		if err := p.asmFuncBody_inst_loong64(fn); err != nil {
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
		// text 段数据(头部空间保留)
		assert(len(p.prog.TextData) == int(p.prog.Entry)-int(p.prog.TextAddr))
		for _, fn := range p.file.Funcs {
			p.prog.TextData = append(p.prog.TextData, fn.LinkInfo.Data...)
		}

		// data 段数据
		assert(len(p.prog.DataData) == 0)
		for _, g := range p.file.Globals {
			p.prog.DataData = append(p.prog.DataData, g.LinkInfo.Data...)
		}
	}

	return p.prog, nil
}

func (p *_Assembler) asmFuncBody_inst_loong64(fn *ast.Func) (err error) {
	// label 的地址列表
	label2pcMap := make(map[string]int64)

	// 绝对地址拆分基于 symbol 名字
	abs_hi2loMap := make(map[string]int32)

	// 相对地址拆分基于 label 名字
	// 一个 label 只能有一个 symbol 相对 PC 寻址
	pc_hi2loMap := make(map[string]int32)

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
	for _, inst := range fn.Body.Insts {
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
			case abi.BuiltinFn_ABS_HI20, abi.BuiltinFn_ABS_HI20_zh: // 高20bit
				hi, lo := pcrel.MakeAbs(uint32(addr)) // TODO: 验证拆分算法是否 OK
				abs_hi2loMap[inst.Arg.Symbol] = lo
				inst.Arg.Imm = hi

			case abi.BuiltinFn_ABS_LO12, abi.BuiltinFn_ABS_LO12_zh: // 低12bit
				lo, ok := abs_hi2loMap[inst.Arg.Symbol]
				if !ok {
					panic(fmt.Errorf("symbol %q not found", inst.Arg.Symbol))
				}
				inst.Arg.Imm = lo

			case abi.BuiltinFn_PC_HI20, abi.BuiltinFn_PC_HI20_zh:
				// 直接以符号名字作为 key 记录, 和 RISCV 处理不同
				hi, lo := pcrel.MakeLa64PCRel(addr, pc)
				pc_hi2loMap[inst.Arg.Symbol] = lo
				inst.Arg.Imm = hi

			case abi.BuiltinFn_PC_LO12, abi.BuiltinFn_PC_LO12_zh:
				lo, ok := pc_hi2loMap[inst.Arg.Symbol]
				if !ok {
					panic(fmt.Errorf("symbol %q not found", inst.Arg.Symbol))
				}
				inst.Arg.Imm = lo

			case abi.BuiltinFn_ABS64_LO20, abi.BuiltinFn_ABS64_LO20_zh:
				panic("TODO: %abs64_lo20")
			case abi.BuiltinFn_ABS64_HI12, abi.BuiltinFn_ABS64_HI12_zh:
				panic("TODO: %abs64_hi12")
			case abi.BuiltinFn_PC64_LO20, abi.BuiltinFn_PC64_LO20_zh:
				panic("TODO: %pc64_lo20")
			case abi.BuiltinFn_PC64_HI12, abi.BuiltinFn_PC64_HI12_zh:
				panic("TODO: %pc64_hi12")

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
		x, err := loong64.EncodeLA64(inst.As, inst.Arg)
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
