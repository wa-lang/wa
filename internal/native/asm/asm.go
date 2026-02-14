// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/loong64"
	"wa-lang.org/wa/internal/native/riscv"
	"wa-lang.org/wa/internal/native/token"
)

// 将汇编语法树转为固定位置的机器码
func AssembleFile(filename string, source []byte, opt *abi.LinkOptions) (prog *abi.LinkedProgram, err error) {
	switch opt.CPU {
	case abi.LOONG64:
		return new(_Assembler).asmFile(filename, source, opt)
	case abi.RISCV32:
		return new(_Assembler).asmFile(filename, source, opt)
	case abi.RISCV64:
		return new(_Assembler).asmFile(filename, source, opt)
	case abi.X64Unix:
		return new(_Assembler).asmFile(filename, source, opt)
	default:
		return nil, fmt.Errorf("unknonw cpu: %v", opt.CPU)
	}
}

type _Assembler struct {
	opt    *abi.LinkOptions
	path   string
	source []byte

	fset *token.FileSet
	file *ast.File
	prog *abi.LinkedProgram

	// 全局符号
	objectMap map[string]ast.Object

	// 下一个指令的位置
	x64NextPcMap map[*ast.Instruction]int64

	// 下个内存分配地址
	dramNextAddr int64
	dramEndAddr  int64
}

func (p *_Assembler) init(filename string, source []byte, opt *abi.LinkOptions) {
	p.opt = opt
	p.path = filename
	p.source = source

	p.fset = token.NewFileSet()

	p.prog = &abi.LinkedProgram{
		CPU: opt.CPU,
	}

	p.objectMap = make(map[string]ast.Object)

	p.dramNextAddr, _ = align(opt.DRAMBase, 4)
	p.dramEndAddr = opt.DRAMBase + opt.DRAMSize
}

// 分配内存空间
func (p *_Assembler) alloc(memSize, addrAlign int64) (addr int64, padding int) {
	assert(addrAlign > 0)
	p.dramNextAddr, padding = align(p.dramNextAddr, addrAlign)
	addr, p.dramNextAddr = p.dramNextAddr, p.dramNextAddr+memSize
	assert(p.dramNextAddr < p.dramEndAddr)
	return addr, padding
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
	case abi.X64Unix, abi.X64Windows:
		return p.instLen_x64(inst)
	default:
		panic("unreachable")
	}
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

func (p *_Assembler) encodeInst(as abi.As, arg *abi.AsArgument) (uint32, error) {
	switch p.file.CPU {
	case abi.LOONG64:
		return loong64.EncodeLA64(as, arg)
	case abi.RISCV64:
		return riscv.EncodeRV64(as, arg)
	default:
		panic("unreachable")
	}
}
