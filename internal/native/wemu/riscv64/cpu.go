// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/riscv"
	"wa-lang.org/wa/internal/native/wemu/device"
)

// 处理器
type CPU struct {
	Mode riscv.Mode   // 工作模式
	RegX [32]uint64   // 整数寄存器
	RegF [32]float64  // 浮点数寄存器
	CSR  [4096]uint64 // 状态寄存器
	PC   uint64       // PC指针
}

// 构建处理器
func (p *CPU) NewCPU() *CPU {
	return new(CPU)
}

// 重置CPU
func (p *CPU) Reset(pc, sp uint64) {
	*p = CPU{}
	p.RegX[riscv.REG_SP] = sp
	p.PC = pc
}

// 单步执行
func (p *CPU) StepRun(bus *device.Bus) error {
	inst, err := bus.Read(p.PC, 4)
	if err != nil {
		return err
	}
	as, arg, err := riscv.Decode(uint32(inst))
	if err != nil {
		return fmt.Errorf("fetch instruntion failed at 0x%08X: %v", p.PC, err)
	}

	if err := p.execInst(bus, as, arg); err != nil {
		return fmt.Errorf("exec instruntion(%08d:%s) failed at 0x%08X: %v", inst,
			riscv.AsmSyntaxEx(int64(p.PC), as, arg, riscv.RegAliasString),
			p.PC,
			err,
		)
	}

	return nil
}

func (p *CPU) execInst(bus *device.Bus, as abi.As, arg *abi.AsArgument) error {
	// 重置0寄存器
	p.RegX[0] = 0
	p.RegF[0] = 0

	// 调整PC
	p.PC += 4

	// 执行指令
	switch as {
	default:
		return fmt.Errorf("unsupport: %s", riscv.AsmSyntaxEx(int64(p.PC), as, arg, riscv.RegAliasString))
	case riscv.AAUIPC:
		// AUIPC: rd = PC + imm
		p.RegX[riscv.RegI(arg.Rd)] = uint64(p.PC) + uint64(arg.Imm)
	case riscv.AADDI:
		// ADDI: rd = rs1 + imm
		p.RegX[riscv.RegI(arg.Rd)] = p.RegX[riscv.RegI(arg.Rs1)] + uint64(arg.Imm)
	case riscv.ALBU:
		// LBU: rd = zero-extended(M[rs1+imm][7:0])
		addr := p.RegX[riscv.RegI(arg.Rs1)] + uint64(arg.Imm)
		if val, err := bus.Read(addr, 1); err == nil {
			p.RegX[riscv.RegI(arg.Rd)] = val // zero-extend
		} else {
			return err
		}
	case riscv.ABEQ:
		// BEQ: if (rs1 == rs2) PC += imm (branch)
		if p.RegX[riscv.RegI(arg.Rs1)] == p.RegX[riscv.RegI(arg.Rs2)] {
			p.PC = uint64(int64(p.PC) + int64(arg.Imm))
		}
	case riscv.ALUI:
		// LUI: rd = imm (高 20 位加载)
		p.RegX[riscv.RegI(arg.Rd)] = uint64(arg.Imm)

	case riscv.ASB:
		// SB: store byte (rs2[7:0]) to M[rs1+imm]
		addr := p.RegX[riscv.RegI(arg.Rs1)] + uint64(arg.Imm)
		if err := bus.Write(addr, 1, p.RegX[riscv.RegI(arg.Rs2)]); err != nil {
			return err
		}
	case riscv.AJAL:
		// JAL: rd = PC + 4; PC += imm
		p.RegX[riscv.RegI(arg.Rd)] = uint64(p.PC + 4)
		p.PC = uint64(int64(p.PC) + int64(arg.Imm))
	case riscv.ASW:
		// SW: store word (rs2) to M[rs1+imm]
		addr := p.RegX[riscv.RegI(arg.Rs1)] + uint64(arg.Imm)
		if err := bus.Write(addr, 4, p.RegX[riscv.RegI(arg.Rs2)]); err != nil {
			return err
		}
	}
	return nil
}
