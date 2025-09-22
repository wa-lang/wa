// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"encoding/binary"
	"fmt"
	"math"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/parser"
	"wa-lang.org/wa/internal/native/riscv"
	"wa-lang.org/wa/internal/native/token"
)

// 将汇编语法树转为固定位置的机器码
func AssembleFile(filename string, source []byte, opt *abi.LinkOptions) (prog *abi.LinkedProgram, err error) {
	return new(_Assembler).asmFile(filename, source, opt)
}

type _Assembler struct {
	opt    *abi.LinkOptions
	path   string
	source []byte

	fset *token.FileSet
	file *ast.File
	prog *abi.LinkedProgram

	// 默认的对齐字节数
	defaultAlign int

	// 下个内存分配地址
	dramNextAddr int64
	dramEndAddr  int64

	// 符号表(不含const)
	symbalMap map[string]*abi.LinkedSymbol
}

func (p *_Assembler) init(filename string, source []byte, opt *abi.LinkOptions) {
	p.opt = opt
	p.path = filename
	p.source = source

	p.fset = token.NewFileSet()

	p.prog = &abi.LinkedProgram{
		CPU: opt.CPU,
	}

	switch p.opt.CPU {
	case abi.RISCV32, abi.RISCV64:
		// RISCV64 也是 4 字节对齐
		p.defaultAlign = 4
	default:
		panic("unreachable")
	}

	p.dramNextAddr = align(opt.DRAMBase, 8)
	p.dramEndAddr = opt.DRAMBase + opt.DRAMSize

	p.symbalMap = make(map[string]*abi.LinkedSymbol)
}

// 分配内存空间
func (p *_Assembler) alloc(memSize, addrAlign int64) (addr int64) {
	if addrAlign == 0 {
		addrAlign = int64(p.defaultAlign)
	}
	assert(addrAlign > 0)
	p.dramNextAddr = align(p.dramNextAddr, addrAlign)
	addr, p.dramNextAddr = p.dramNextAddr, p.dramNextAddr+memSize
	assert(p.dramNextAddr < p.dramEndAddr)
	return addr
}

func (p *_Assembler) asmFile(filename string, source []byte, opt *abi.LinkOptions) (prog *abi.LinkedProgram, err error) {
	p.init(filename, source, opt)

	// 解析汇编程序
	p.file, err = parser.ParseFile(opt.CPU, p.fset, filename, source)
	if err != nil {
		return nil, err
	}

	// 全局函数分配内存空间
	for _, fn := range p.file.Funcs {
		fn.Size = int(p.funcBodyLen(fn))
		fn.LinkInfo = &abi.LinkedSymbol{
			Name: fn.Name,
			Addr: p.alloc(int64(fn.Size), 0),
			Data: make([]byte, fn.Size),
		}
	}

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
		if err := p.asmFunc(fn); err != nil {
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
		p.prog.TextAddr = 0

		// 优先查找指定的入口函数
		if opt.EntryFunc != "" {
			for _, fn := range p.file.Funcs {
				if fn.Name == opt.EntryFunc {
					p.prog.TextAddr = fn.LinkInfo.Addr
				}
			}
		}

		// 然后查找默认的入口函数(中文)
		if p.prog.TextAddr == 0 {
			for _, fn := range p.file.Funcs {
				if fn.Name == abi.DefaultEntryFuncZh {
					p.prog.TextAddr = fn.LinkInfo.Addr
				}
			}
		}

		// 然后查找默认的入口函数(英文)
		if p.prog.TextAddr == 0 {
			for _, fn := range p.file.Funcs {
				if fn.Name == abi.DefaultEntryFunc {
					p.prog.TextAddr = fn.LinkInfo.Addr
				}
			}
		}

		// 查找失败
		if p.prog.TextAddr == 0 {
			return nil, fmt.Errorf("entry %q not found", opt.EntryFunc)
		}

		// data 段地址
		p.prog.DataAddr = opt.DRAMBase
		if len(p.file.Globals) > 0 {
			p.prog.DataAddr = p.file.Globals[0].LinkInfo.Addr
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

func (p *_Assembler) asmFunc(fn *ast.Func) (err error) {
	// 标签上下文信息
	type LabelContext struct {
		Name            string // 标签名称
		PC              int64  // 对应的PC
		Pcrel_hi_symbol string // 用于 %pcrel_lo(label) 查询
	}

	// 第一遍扫描Label, 生成对应的地址
	labelContextMap := make(map[string]*LabelContext)
	labelAddr := fn.LinkInfo.Addr
	for i, inst := range fn.Body.Insts {
		if inst.Label != "" {
			labelContextMap[inst.Label] = &LabelContext{
				Name: inst.Label,
				PC:   labelAddr,
			}
		}
		if inst.As == 0 {
			continue
		}

		// 记录 %pcrel_hi(symbol) 符号参数
		if inst.Arg.SymbolDecor == abi.BuiltinFn_PCREL_HI || inst.Arg.SymbolDecor == abi.BuiltinFn_PCREL_HI_zh {
			for k := i; k >= 0; k-- {
				if x := fn.Body.Insts[k]; x.Label != "" {
					if labelCtx, ok := labelContextMap[x.Label]; ok {
						if labelCtx.Pcrel_hi_symbol == "" {
							labelCtx.Pcrel_hi_symbol = inst.Arg.Symbol
						}
					}
					break
				}
			}
		}

		labelAddr += p.instLen(inst)
	}

	// 第二遍编码指令
	var pc = fn.LinkInfo.Addr
	for _, inst := range fn.Body.Insts {
		if inst.As == 0 {
			// 跳过空的指令, 比如标号
			continue
		}

		// 指令对 label 或全局的符号引用
		// 因为指令长度的关系, 指令并不会直接访问符号对应的绝对地址
		// 需要解决 %hi/%lo/%pcrel_hi/%pcrel_lo 等转化为最终可编码到指令的值
		if inst.Arg.Symbol != "" {
			PCREL_LO_PC := pc // 可能不是当前的 pc
			addr, ok := int64(0), bool(false)
			if inst.Arg.SymbolDecor == abi.BuiltinFn_PCREL_LO || inst.Arg.SymbolDecor == abi.BuiltinFn_PCREL_LO_zh {
				labelCtx := labelContextMap[inst.Arg.Symbol]
				if labelCtx == nil {
					panic(fmt.Errorf("label %q not found", inst.Arg.Symbol))
				}
				PCREL_LO_PC = labelCtx.PC
				addr, ok = p.symbolAddress(labelCtx.Pcrel_hi_symbol)
				if !ok {
					panic(fmt.Errorf("symbol %q not found", inst.Arg.Symbol))
				}
			} else {
				labelCtx := labelContextMap[inst.Arg.Symbol]
				if labelCtx != nil {
					addr, ok = labelCtx.PC, true
				} else {
					addr, ok = p.symbolAddress(inst.Arg.Symbol)
					if !ok {
						panic(fmt.Errorf("symbol %q not found", inst.Arg.Symbol))
					}
				}
			}

			switch inst.Arg.SymbolDecor {
			case abi.BuiltinFn_HI, abi.BuiltinFn_HI_zh: // 高20bit
				x := int32(addr)
				// 检查 lo(x) 的符号位（第 11 位）
				// (x & 0x800) == 0x800 等价于 (x << 20) >> 31 == -1
				if (x & 0x800) == 0x800 {
					// 如果 lo(x) 是负数，则 hi(x) 加 1
					inst.Arg.Imm = (x >> 12) + 1
				} else {
					inst.Arg.Imm = x >> 12
				}
			case abi.BuiltinFn_LO, abi.BuiltinFn_LO_zh: // 低12bit
				x := int32(addr)
				// 简单地取低 12 位
				inst.Arg.Imm = x & 0xFFF
			case abi.BuiltinFn_PCREL_HI, abi.BuiltinFn_PCREL_HI_zh:
				offset := int32(addr - pc)
				// 检查 lo(offset) 的符号位
				if (offset & 0x800) == 0x800 {
					inst.Arg.Imm = (offset >> 12) + 1
				} else {
					inst.Arg.Imm = offset >> 12
				}
			case abi.BuiltinFn_PCREL_LO, abi.BuiltinFn_PCREL_LO_zh:
				// https://sourceware.org/binutils/docs/as/RISC_002dV_002dModifiers.html
				// https://stackoverflow.com/questions/65879012/what-do-pcrel-hi-and-pcrel-lo-actually-do
				offset := int32(addr - PCREL_LO_PC)
				inst.Arg.Imm = offset & 0xFFF
			default:
				// 因为riscv指令只有32bit宽度, 默认是无法完全编码绝对地址的
				// 所以其他情况都也作为相对pc的地址
				inst.Arg.Imm = int32(addr - pc)
			}
		}

		// 编码使用的是符号被处理后对应的立即数
		x, err := riscv.Encode(p.opt.CPU, inst.As, inst.Arg)
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

// 计算函数指令需要的内存大小
func (p *_Assembler) funcBodyLen(fn *ast.Func) (n int64) {
	for _, inst := range fn.Body.Insts {
		n += p.instLen(inst)
	}
	return
}

func (p *_Assembler) instLen(inst *ast.Instruction) int64 {
	if inst.As == 0 {
		return 0
	}
	switch p.opt.CPU {
	case abi.RISCV32, abi.RISCV64:
		return 4
	default:
		panic("unreachable")
	}
}

func (p *_Assembler) asmGlobal(g *ast.Global) (err error) {
	for _, xInit := range g.Init {
		if xInit.Symbal != "" {
			v, ok := p.symbolAddress(xInit.Symbal)
			if !ok {
				panic(fmt.Errorf("symbol %q not found", xInit.Symbal))
			}
			if p.opt.CPU == abi.RISCV32 {
				binary.LittleEndian.PutUint32(g.LinkInfo.Data, uint32(v))
			} else {
				binary.LittleEndian.PutUint64(g.LinkInfo.Data, uint64(v))
			}
			continue
		}

		// 常量面值初始化
		switch xInit.Lit.TypeCast {
		case token.I32, token.I32_zh:
			v := xInit.Lit.ConstV.(int64)
			binary.LittleEndian.PutUint32(g.LinkInfo.Data, uint32(v))
		case token.I64, token.I64_zh:
			v := xInit.Lit.ConstV.(int64)
			binary.LittleEndian.PutUint64(g.LinkInfo.Data, uint64(v))
		case token.U32, token.U32_zh:
			v := xInit.Lit.ConstV.(int64)
			binary.LittleEndian.PutUint32(g.LinkInfo.Data, uint32(v))
		case token.U64, token.U64_zh:
			v := xInit.Lit.ConstV.(int64)
			binary.LittleEndian.PutUint64(g.LinkInfo.Data, uint64(v))
		case token.F32, token.F32_zh:
			v := xInit.Lit.ConstV.(float64)
			binary.LittleEndian.PutUint32(g.LinkInfo.Data, math.Float32bits(float32(v)))
		case token.F64, token.F64_zh:
			v := xInit.Lit.ConstV.(float64)
			binary.LittleEndian.PutUint64(g.LinkInfo.Data, math.Float64bits(float64(v)))
		case token.PTR, token.PTR_zh:
			v := xInit.Lit.ConstV.(int64)
			if p.opt.CPU == abi.RISCV32 {
				binary.LittleEndian.PutUint32(g.LinkInfo.Data, uint32(v))
			} else {
				binary.LittleEndian.PutUint64(g.LinkInfo.Data, uint64(v))
			}
		default:
			assert(xInit.Lit.LitKind == token.STRING)
			copy(g.LinkInfo.Data, []byte(xInit.Lit.ConstV.(string)))
		}
	}

	p.symbalMap[g.Name] = g.LinkInfo
	p.prog.DataData = append(p.prog.DataData, g.LinkInfo.Data...)
	return nil
}

// 全局变量或函数的地址
func (p *_Assembler) symbolAddress(s string) (int64, bool) {
	// 查找全局的常量(只会引用整数类型)
	for _, x := range p.file.Consts {
		if x.Name == s {
			v, ok := x.Value.ConstV.(int64)
			if !ok {
				panic(fmt.Sprintf("const %q is %v", s, x.Value.ConstV))
			}
			return v, true
		}
	}

	// 查找全局变量
	for _, x := range p.file.Globals {
		if x.Name == s {
			return x.LinkInfo.Addr, true
		}
	}

	// 查找全局函数
	for _, x := range p.file.Funcs {
		if x.Name == s {
			return x.LinkInfo.Addr, true
		}
	}

	// 查找失败
	return 0, false
}
