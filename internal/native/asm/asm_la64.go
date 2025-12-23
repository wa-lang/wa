// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"encoding/binary"
	"fmt"
	"sort"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/loong64"
	"wa-lang.org/wa/internal/native/parser"
	"wa-lang.org/wa/internal/native/pcrel"
	"wa-lang.org/wa/internal/native/token"
)

func (p *_Assembler) asmFile_loong64(filename string, source []byte, opt *abi.LinkOptions) (prog *abi.LinkedProgram, err error) {
	// 最大的页大小
	const maxPageSize = 64 << 10

	// 解析汇编程序
	p.file, err = parser.ParseFile(opt.CPU, p.fset, filename, source)
	if err != nil {
		return nil, err
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
		fn.Size = int(p.funcBodyLen(fn))
		fn.LinkInfo = &abi.LinkedSymbol{
			Name: fn.Name,
			Addr: p.alloc(int64(fn.Size), 0),
			Data: make([]byte, fn.Size),
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
		if err := p.asmFunc_loong64(fn); err != nil {
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

func (p *_Assembler) asmFunc_loong64(fn *ast.Func) (err error) {
	if err := p.asmFuncBody_local_loong64(fn); err != nil {
		return err
	}
	if err := p.asmFuncBody_inst_loong64(fn); err != nil {
		return err
	}
	return nil
}

func (p *_Assembler) asmFuncBody_local_loong64(fn *ast.Func) error {
	const (
		iArgRegEnd = loong64.REG_A7
		fArgRegEnd = loong64.REG_F7

		iRetRegEnd = loong64.REG_A1
		fRetRegEnd = loong64.REG_F1

		// 局部变量在 FP 低地址, 最终要取反为负数
		// 开头有 ra 和 fp 两个空间要跳过
		FP_localOffsetBase = (8 + 8)
	)
	var (
		iArgReg = loong64.REG_A0
		fArgReg = loong64.REG_F0

		iRetReg = loong64.REG_A0
		fRetReg = loong64.REG_F0

		argRetOffset = 0 // 参数和返回值偏移, 包含返回值 从 0 开始, 正整数
		localOffset  = 0 // 局部变量偏移
	)

	// 输入参数分配
	for _, arg := range fn.Type.Args {
		switch arg.Type {
		case token.I32, token.I32_zh:
			if iArgReg <= iArgRegEnd {
				arg.Reg = iArgReg
				iArgReg++
			} else {
				arg.Off = argRetOffset
				argRetOffset += 4
			}
		case token.I64, token.I64_zh:
			if iArgReg <= iArgRegEnd {
				arg.Reg = iArgReg
				iArgReg++
			} else {
				if argRetOffset%8 != 0 {
					argRetOffset += 4
				}
				arg.Off = argRetOffset
				argRetOffset += 8
			}
		case token.F32, token.F32_zh:
			if fArgReg <= fArgRegEnd {
				arg.Reg = fArgReg
				fArgReg++
			} else {
				arg.Off = argRetOffset
				argRetOffset += 4
			}
		case token.F64, token.F64_zh:
			if fArgReg <= fArgRegEnd {
				arg.Reg = fArgReg
				fArgReg++
			} else {
				if argRetOffset%8 != 0 {
					argRetOffset += 4
				}
				arg.Off = argRetOffset
				argRetOffset += 8
			}
		}
	}

	// 返回值寄存器分配
	// 栈传递部分和参数的栈转递部分相对关系?
	for _, ret := range fn.Type.Return {
		switch ret.Type {
		case token.I32, token.I32_zh:
			if iArgReg <= iRetRegEnd {
				ret.Reg = iRetReg
				iArgReg++
			} else {
				ret.Off = argRetOffset
				argRetOffset += 4
			}
		case token.I64, token.I64_zh:
			if iArgReg <= iRetRegEnd {
				ret.Reg = iRetReg
				iArgReg++
			} else {
				if argRetOffset%8 != 0 {
					argRetOffset += 4
				}
				ret.Off = argRetOffset
				argRetOffset += 8
			}
		case token.F32, token.F32_zh:
			if fArgReg <= fRetRegEnd {
				ret.Reg = fRetReg
				fArgReg++
			} else {
				ret.Off = argRetOffset
				argRetOffset += 4
			}
		case token.F64, token.F64_zh:
			if fArgReg <= fRetRegEnd {
				ret.Reg = fRetReg
				fArgReg++
			} else {
				if argRetOffset%8 != 0 {
					argRetOffset += 4
				}
				ret.Off = argRetOffset
				argRetOffset += 8
			}
		}
	}

	// 为了减少栈内存碎片, 局部变量会重新排序
	locals := append([]*ast.Local{}, fn.Body.Locals...)
	sort.Slice(locals, func(i, j int) bool {
		return localSize(locals[i]) < localSize(locals[j])
	})

	// 局部变量分配
	for _, local := range locals {
		switch local.Type {
		case token.I32, token.I32_zh:
			local.Off = -4 - (FP_localOffsetBase + localOffset)
			localOffset += 4
		case token.I64, token.I64_zh:
			if localOffset%8 != 0 {
				localOffset += 4
			}
			local.Off = -8 - (FP_localOffsetBase + localOffset)
			localOffset += 8
		case token.F32, token.F32_zh:
			local.Off = -4 - (FP_localOffsetBase + localOffset)
			localOffset += 4
		case token.F64, token.F64_zh:
			if localOffset%8 != 0 {
				localOffset += 4
			}
			local.Off = -8 - (FP_localOffsetBase + localOffset)
			localOffset += 8
		}
	}

	if localOffset%8 != 0 {
		localOffset += 4
	}

	// 函数栈上的参数和返回值大小
	fn.LinkInfo.ArgsSize = argRetOffset

	// 函数栈帧大小
	fn.LinkInfo.FrameSize = FP_localOffsetBase + localOffset

	return nil
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
