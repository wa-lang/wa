// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 重要: 请保持 RV32 和 RV64 版本完全一致, LAUInt 在另一个文件定义.

package loong64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/loong64"
	"wa-lang.org/wa/internal/native/wemu/device"
)

// 处理器
type CPU struct {
	RegX [32]LAUInt  // 整数寄存器
	RegF [32]float64 // 浮点数寄存器
	PC   LAUInt      // PC指针
}

// 构建处理器
func NewCPU() *CPU {
	return new(CPU)
}

func (p *CPU) GetPC() uint64  { return uint64(p.PC) }
func (p *CPU) SetPC(v uint64) { p.PC = LAUInt(v) }

func (p *CPU) XRegNum() int            { return len(p.RegX) }
func (p *CPU) GetXReg(i int) uint64    { return uint64(p.RegX[i]) }
func (p *CPU) SetXReg(i int, v uint64) { p.RegX[i] = LAUInt(v) }

func (p *CPU) FRegNum() int             { return len(p.RegF) }
func (p *CPU) GetFReg(i int) float64    { return p.RegF[i] }
func (p *CPU) SetFReg(i int, v float64) { p.RegF[i] = v }

// 重置CPU
func (p *CPU) Reset(pc, sp uint64) {
	*p = CPU{}
	p.RegX[loong64.RegI(loong64.REG_SP)] = LAUInt(sp)
	p.PC = LAUInt(pc)
}

// 单步执行
func (p *CPU) StepRun(bus *device.Bus) error {
	inst, err := bus.Read(uint64(p.PC), 4)
	if err != nil {
		return err
	}
	as, arg, argRaw, err := loong64.DecodeEx(uint32(inst))
	if err != nil {
		return fmt.Errorf("fetch instruntion failed at 0x%08X: %v", p.PC, err)
	}

	if err := p.execInst(bus, as, argRaw); err != nil {
		return fmt.Errorf("exec instruntion(%08d:%s) failed at 0x%08X: %v", inst,
			loong64.AsmSyntax(as, "", arg),
			p.PC,
			err,
		)
	}

	return nil
}

func (p *CPU) execInst(bus *device.Bus, as abi.As, arg *abi.AsRawArgument) error {
	// 重置0寄存器
	p.RegX[0] = 0
	p.RegF[0] = 0

	// 当前的PC
	curPC := p.PC
	_ = curPC

	// 调整PC(少数跳转指令会覆盖)
	p.PC += 4

	// 执行指令
	switch as {
	default:
		return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))

	case loong64.ALU12I_W:
		p.RegX[arg.Rd] = LAUInt(arg.Imm << 12)
	case loong64.AADDI_W:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] + LAUInt(arg.Imm)
	case loong64.ALD_BU:
		addr := p.RegX[arg.Rs1] + LAUInt(arg.Imm)
		value, err := bus.Read(uint64(addr), 1)
		if err != nil {
			return err
		}
		p.RegX[arg.Rd] = LAUInt(uint8(value))
	case loong64.AST_B:
		addr := p.RegX[arg.Rs1] + LAUInt(arg.Imm)
		value := p.RegX[arg.Rs2]
		if err := bus.Write(uint64(addr), 1, uint64(value)); err != nil {
			return err
		}
	case loong64.ABEQ:
		if p.RegX[arg.Rs1] == p.RegX[arg.Rs2] {
			p.PC = curPC + LAUInt(arg.Imm)
		}
	case loong64.AB:
		p.PC = curPC + LAUInt(arg.Imm)
	}

	return nil
}
