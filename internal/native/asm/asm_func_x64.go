// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/x64"
	"wa-lang.org/wa/internal/native/x64/p9x86"
)

func (p *_Assembler) asmFuncBody_inst_x64_linux(fn *ast.Func) (err error) {
	// 函数内部指令的PC列表
	p.x64NextPcMap = make(map[*ast.Instruction]int64)

	// label 的地址列表
	label2pcMap := make(map[string]int64)

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

		// 记录下个指针位置
		p.x64NextPcMap[inst] = pc
	}

	// 第二遍遍历编码指令
	pc = fn.LinkInfo.Addr
	for _, inst := range fn.Body.Insts {
		if inst.As == 0 {
			// 跳过空的指令, 比如标号
			continue
		}

		// 编码使用的是符号被处理后对应的立即数
		code := p.encodeInst_x64(inst, pc, label2pcMap)
		fn.LinkInfo.Data = append(fn.LinkInfo.Data, code...)

		// 更新下一个指令对应的 pc 位置
		pc += int64(len(code))

		// 之前预估的指令长度必须和真实的长度一致
		assert(pc == p.x64NextPcMap[inst])
	}

	return nil
}

func (p *_Assembler) encodeInst_x64(inst *ast.Instruction, pc int64, label2pcMap map[string]int64) []byte {
	p.x64FixSymbol(inst, pc, label2pcMap)
	prog := p.x64Inst2P9Prog(inst)
	code := p9x86.AsmInst(prog)
	return code
}

func (p *_Assembler) x64FixSymbol(inst *ast.Instruction, pc int64, label2pcMap map[string]int64) {
	// 遍历操作数, 查找是否有符号引用
	ops := []*abi.X64Operand{inst.ArgX64.Dst, inst.ArgX64.Src}
	for _, op := range ops {
		if op == nil || op.Symbol == "" {
			continue
		}

		// 获取真实的地址
		targetPC, ok := label2pcMap[op.Symbol]
		if !ok {
			// 是外部符号
			if v, ok := p.symbolAddress(op.Symbol); ok {
				targetPC = v
			} else {
				panic("unknown symbol:" + op.Symbol)
			}
		}

		switch {
		case op.Kind == abi.X64Operand_Mem && op.Reg == x64.REG_RIP:
			instrLen := p.x64NextPcMap[inst]
			op.Offset = targetPC - (pc + instrLen)
		case inst.As == x64.ACALL || inst.As == x64.AJMP:
			instrLen := p.x64NextPcMap[inst]
			op.Imm = targetPC - (pc + instrLen)
		default:
			panic("unreachable")
		}
	}
}

func (p *_Assembler) x64Inst2P9Prog(inst *ast.Instruction) *p9x86.Prog {
	prog := &p9x86.Prog{}
	prog.To = p.x64Operand2P9Addr(inst.ArgX64.Dst)
	prog.From = p.x64Operand2P9Addr(inst.ArgX64.Src)

	switch inst.As {
	case x64.AMOV:
		// TODO: 根据寄存器类型选择 MOV 指令
		prog.As = p9x86.AMOVQ
	case x64.AADD:
		prog.As = p9x86.AADDQ
	case x64.ASUB:
		prog.As = p9x86.ASUBQ
	case x64.APUSH:
		prog.As = p9x86.APUSHQ
	case x64.APOP:
		prog.As = p9x86.APOPQ
	case x64.ARET:
		prog.As = p9x86.ARET
	case x64.ASYSCALL:
		prog.As = p9x86.ASYSCALL
	case x64.ACALL:
		assert(prog.From.Type == p9x86.TYPE_NONE)
		prog.As = p9x86.ACALL
	case x64.AJMP:
		prog.As = p9x86.AJMP

	default:
		panic(fmt.Sprintf("TODO: %v", inst.As))
	}

	return prog
}

func (p *_Assembler) x64Operand2P9Addr(op *abi.X64Operand) p9x86.Addr {
	if op == nil {
		return p9x86.Addr{Type: p9x86.TYPE_NONE}
	}

	addr := p9x86.Addr{}

	switch op.Kind {
	case abi.X64Operand_Reg:
		// 映射为寄存器类型
		addr.Type = p9x86.TYPE_REG
		addr.Reg = int16(op.Reg)

	case abi.X64Operand_Imm:
		// 映射为常量（立即数）
		// 例如：MOV $123, RAX 中的 $123
		addr.Type = p9x86.TYPE_CONST
		addr.Offset = op.Imm

	case abi.X64Operand_Mem:
		// 映射为内存引用
		// 这里的逻辑处理 [Reg + Offset]
		addr.Type = p9x86.TYPE_MEM
		addr.Reg = int16(op.Reg) // 基址寄存器 (Base)
		addr.Offset = op.Offset  // 位移 (Displacement)

		// 如果是 RIP 相对寻址 [rip + 0x10]，
		// 你的解析器应确保 op.Reg 是 REG_RIP 的编号
	}

	return addr
}

func (p *_Assembler) instLen_x64(inst *ast.Instruction) int64 {
	switch inst.As {
	case x64.AMOV:
		panic("TODO")
	default:
		panic("TODO")
	}
}
