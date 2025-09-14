// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"bytes"
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

	p.dramNextAddr = align(opt.DRAMBase, 8)
	p.dramEndAddr = opt.DRAMBase + opt.DRAMSize

	p.symbalMap = make(map[string]*abi.LinkedSymbol)
}

// 分配内存空间
func (p *_Assembler) alloc(memSize, addrAlign int64) (addr int64) {
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

	// 编译全局变量
	p.prog.DataAddr = p.dramNextAddr
	for _, g := range p.file.Globals {
		if err := p.asmGlobal(g); err != nil {
			return nil, err
		}
	}

	// 编译函数
	p.prog.TextAddr = p.dramNextAddr
	for _, fn := range p.file.Funcs {
		if err := p.asmFunc(fn); err != nil {
			return nil, err
		}
	}

	// TODO: 回填未知的符号
	// TODO: 合并全部数据

	return p.prog, nil
}

func (p *_Assembler) asmGlobal(g *ast.Global) (err error) {
	g.LinkInfo = &abi.LinkedSymbol{
		Name: g.Name,
		Addr: p.alloc(int64(g.Size), 8),
		Size: int64(g.Size),
		Data: make([]byte, g.Size),
	}

	for _, xInit := range g.Init {
		if xInit.Symbal != "" {
			// 符号在后续同续处理
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
			copy(g.LinkInfo.Data, []byte(xInit.Lit.LitString))
		}
	}

	p.symbalMap[g.Name] = g.LinkInfo
	p.prog.DataData = append(p.prog.DataData, g.LinkInfo.Data...)
	return nil
}

func (p *_Assembler) asmFunc(fn *ast.Func) (err error) {
	fn.LinkInfo = &abi.LinkedSymbol{
		Addr: p.alloc(0, 8),
		Name: fn.Name,
	}

	// 第一遍扫描Label, 生成对应的地址
	labelAddrMap := make(map[string]int64)
	labelAddr := fn.LinkInfo.Addr
	for _, inst := range fn.Body.Insts {
		if inst.Label != "" {
			labelAddrMap[inst.Label] = labelAddr
		}
		if inst.As == 0 {
			continue
		}
		labelAddr += p.instLen(inst)
	}

	// 第二遍编码指令
	var bufText bytes.Buffer
	for _, inst := range fn.Body.Insts {
		if inst.As == 0 {
			// 跳过空的指令, 比如标号
			continue
		}

		// 优先处理参数中引用的局部标号
		// 其他全局的符号后续再处理(后续的处理将不包含局部符号引用)
		if inst.Arg.Symbol != "" {
			if x, ok := labelAddrMap[inst.Arg.Symbol]; ok {
				_ = x
			}
		}

		// TODO: 编码指令需要延后, 需要生成全部的全局符号地址
		x, err := riscv.Encode(p.opt.CPU, inst.As, inst.Arg)
		if err != nil {
			return fmt.Errorf("%v: %w", inst, err)
		}
		err = binary.Write(&bufText, binary.LittleEndian, x)
		if err != nil {
			return fmt.Errorf("%v: %w", inst, err)
		}
	}

	p.prog.TextData = bufText.Bytes()
	return nil
}

func (p *_Assembler) instLen(*ast.Instruction) int64 {
	switch p.opt.CPU {
	case abi.RISCV32, abi.RISCV64:
		return 4
	default:
		panic("unreachable")
	}
}
