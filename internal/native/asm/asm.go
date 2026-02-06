// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"encoding/binary"
	"fmt"
	"math"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
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

func (p *_Assembler) asmFile(filename string, source []byte, opt *abi.LinkOptions) (prog *abi.LinkedProgram, err error) {
	p.init(filename, source, opt)

	switch p.prog.CPU {
	case abi.LOONG64:
		return p.asmFile_loong64(filename, source, opt)
	case abi.RISCV32:
		return p.asmFile_riscv(filename, source, opt)
	case abi.RISCV64:
		return p.asmFile_riscv(filename, source, opt)
	default:
		return nil, fmt.Errorf("unknonw cpu: %v", p.prog.CPU)
	}
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
	case abi.LOONG64:
		p.defaultAlign = 4
	case abi.RISCV32, abi.RISCV64:
		p.defaultAlign = 4
	default:
		panic("unreachable")
	}

	p.dramNextAddr = align(opt.DRAMBase, 4)
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
	case abi.LOONG64:
		return 4
	case abi.RISCV32, abi.RISCV64:
		return 4
	default:
		panic("unreachable")
	}
}

func (p *_Assembler) asmGlobal(g *ast.Global) (err error) {
	// g.LinkInfo.Data 空间需要提前初始化
	if g.Init.Symbal != "" {
		v, ok := p.symbolAddress(g.Init.Symbal)
		if !ok {
			panic(fmt.Errorf("symbol %q not found", g.Init.Symbal))
		}
		if p.opt.CPU == abi.RISCV32 {
			binary.LittleEndian.PutUint32(g.LinkInfo.Data, uint32(v))
		} else {
			binary.LittleEndian.PutUint64(g.LinkInfo.Data, uint64(v))
		}
	}

	// 常量面值初始化
	switch g.Type {
	case token.I8:
		v := g.Init.Lit.ConstV.(int64)
		switch {
		case v >= 0 && v <= math.MaxUint8:
			g.LinkInfo.Data = append(g.LinkInfo.Data, uint8(v))
		case v >= math.MinInt8 && v <= math.MaxInt8:
			g.LinkInfo.Data = append(g.LinkInfo.Data, uint8(int8(v)))
		default:
			panic(fmt.Errorf("global %q init value overflow: %v", g.Name, g.Init.Lit))
		}
	case token.I16:
		v := g.Init.Lit.ConstV.(int64)
		switch {
		case v >= 0 && v <= math.MaxUint16:
			binary.LittleEndian.PutUint16(g.LinkInfo.Data, uint16(v))
		case v >= math.MinInt16 && v <= math.MaxInt16:
			binary.LittleEndian.PutUint16(g.LinkInfo.Data, uint16(int16(v)))
		default:
			panic(fmt.Errorf("global %q init value overflow: %v", g.Name, g.Init.Lit))
		}
	case token.I32:
		v := g.Init.Lit.ConstV.(int64)
		switch {
		case v >= 0 && v <= math.MaxUint32:
			binary.LittleEndian.PutUint32(g.LinkInfo.Data, uint32(v))
		case v >= math.MinInt32 && v <= math.MaxInt32:
			binary.LittleEndian.PutUint32(g.LinkInfo.Data, uint32(int32(v)))
		default:
			panic(fmt.Errorf("global %q init value overflow: %v", g.Name, g.Init.Lit))
		}
	case token.I64:
		switch v := g.Init.Lit.ConstV.(type) {
		case int64:
			if g.TypeTok == token.GAS_SKIP {
				g.LinkInfo.Data = make([]byte, v)
			} else {
				binary.LittleEndian.PutUint64(g.LinkInfo.Data, uint64(v))
			}
		case string:
			g.LinkInfo.Data = []byte(v)
		case []byte:
			g.LinkInfo.Data = v
		default:
			panic("unreachable")
		}
	case token.F32:
		v := g.Init.Lit.ConstV.(float64)
		binary.LittleEndian.PutUint32(g.LinkInfo.Data, math.Float32bits(float32(v)))
	case token.F64:
		v := g.Init.Lit.ConstV.(float64)
		binary.LittleEndian.PutUint64(g.LinkInfo.Data, math.Float64bits(float64(v)))
	default:
		assert(g.Init.Lit.LitKind == token.STRING)
		copy(g.LinkInfo.Data, []byte(g.Init.Lit.ConstV.(string)))
	}

	p.symbalMap[g.Name] = g.LinkInfo
	return nil
}

// 全局变量或函数的地址
func (p *_Assembler) symbolAddress(s string) (int64, bool) {
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
