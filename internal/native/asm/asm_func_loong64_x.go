// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"encoding/binary"
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/pcrel"
)

func (p *_Assembler) asmFuncBody_inst_loong64(fn *ast.Func) (err error) {
	// label 的地址列表
	label2pcMap := make(map[string]int64)

	// 绝对地址拆分基于 symbol 名字
	abs_hi2loMap := make(map[string]int32)

	// 龙芯语法的 symbol 相对 PC 寻址映射
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
					panic(fmt.Errorf("%s: symbol %q not found", fn.Name, inst.Arg.Symbol))
				}
				inst.Arg.Imm = lo

			case abi.BuiltinFn_PC_HI20, abi.BuiltinFn_PC_HI20_zh:
				// TODO: 这2个宏不支持独立使用, 需要改进!!

				// 直接以符号名字作为 key 记录, 和 RISCV 处理不同
				hi, lo := pcrel.MakeLa64PCRel(addr, pc)
				pc_hi2loMap[inst.Arg.Symbol] = lo
				inst.Arg.Imm = hi

			case abi.BuiltinFn_PC_LO12, abi.BuiltinFn_PC_LO12_zh:
				// TODO: 这2个宏不支持独立使用, 需要改进!!

				lo, ok := pc_hi2loMap[inst.Arg.Symbol]
				if !ok {
					panic(fmt.Errorf("%s: symbol %q not found", fn.Name, inst.Arg.Symbol))
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

				// 龙芯平台, 符号作为地址大于等于 uint32
				// 必须处理为 pc 相对地址才有可能正常表示
				inst.Arg.Imm = int32(addr - pc)
			}
		}

		// 编码使用的是符号被处理后对应的立即数
		x, err := p.encodeInst(inst.As, inst.Arg)
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
