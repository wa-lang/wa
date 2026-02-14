// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/x64"
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
		if n, err := x64.EncodeLen(inst.As, inst.ArgX64); err == nil {
			pc += int64(n)
		} else {
			panic(fmt.Errorf("encode %v failed: %v", inst.As, err))
		}

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

		// 修复符号的地址
		// pc 需要使用下个指令的位置
		p.x64FixSymbol(inst, p.x64NextPcMap[inst], label2pcMap)

		// 编码使用的是符号被处理后对应的立即数
		if code, err := x64.Encode(inst.As, inst.ArgX64); err == nil {
			fn.LinkInfo.Data = append(fn.LinkInfo.Data, code...)
		} else {
			panic(fmt.Errorf("encode %v failed: %v", inst.As, err))
		}

		// 更新下一个指令对应的 pc 位置
		pc += int64(len(fn.LinkInfo.Data))

		// 之前预估的指令长度必须和真实的长度一致
		assert(pc == p.x64NextPcMap[inst])
	}

	return nil
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
