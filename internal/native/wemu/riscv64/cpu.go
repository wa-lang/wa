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
	as, arg, argRaw, err := riscv.DecodeEx(uint32(inst))
	if err != nil {
		return fmt.Errorf("fetch instruntion failed at 0x%08X: %v", p.PC, err)
	}

	if err := p.execInst(bus, as, argRaw); err != nil {
		return fmt.Errorf("exec instruntion(%08d:%s) failed at 0x%08X: %v", inst,
			riscv.AsmSyntaxEx(as, arg, riscv.RegAliasString),
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

	// 调整PC
	p.PC += 4

	// 执行指令
	switch as {
	default:
		return fmt.Errorf("unsupport: %s", riscv.AsString(as))

	// RV32I Base Instruction Set

	case riscv.ALUI:
		p.RegX[arg.Rd] = uint64(arg.Imm << 12)
	case riscv.AAUIPC:
		p.RegX[arg.Rd] = uint64(p.PC) + uint64(arg.Imm) - 4
	case riscv.AJAL:
		t := p.PC
		p.PC = p.PC + uint64(arg.Imm) - 4
		p.RegX[arg.Rd] = t
	case riscv.AJALR:
		t := p.PC
		p.PC = p.RegX[arg.Rs1] + uint64(arg.Imm)
		p.RegX[arg.Rd] = t
	case riscv.ABEQ:
		if p.RegX[arg.Rs1] == p.RegX[arg.Rs2] {
			p.PC = p.PC + uint64(arg.Imm) - 4
		}
	case riscv.ABNE:
		if p.RegX[arg.Rs1] != p.RegX[arg.Rs2] {
			p.PC = p.PC + uint64(arg.Imm) - 4
		}
	case riscv.ABLT:
		if int64(p.RegX[arg.Rs1]) < int64(p.RegX[arg.Rs2]) {
			p.PC = p.PC + uint64(arg.Imm) - 4
		}
	case riscv.ABGE:
		if int64(p.RegX[arg.Rs1]) > int64(p.RegX[arg.Rs2]) {
			p.PC = p.PC + uint64(arg.Imm) - 4
		}
	case riscv.ABLTU:
		if p.RegX[arg.Rs1] < p.RegX[arg.Rs2] {
			p.PC = p.PC + uint64(arg.Imm) - 4
		}
	case riscv.ABGEU:
		if p.RegX[arg.Rs1] > p.RegX[arg.Rs2] {
			p.PC = p.PC + uint64(arg.Imm) - 4
		}
	case riscv.ALB:
		addr := p.RegX[arg.Rs1] + uint64(arg.Imm)
		value, err := bus.Read(addr, 1)
		if err != nil {
			return err
		}
		p.RegX[arg.Rd] = uint64(int8(value))
	case riscv.ALH:
		addr := p.RegX[arg.Rs1] + uint64(arg.Imm)
		value, err := bus.Read(addr, 2)
		if err != nil {
			return err
		}
		p.RegX[arg.Rd] = uint64(int16(value))
	case riscv.ALW:
		addr := p.RegX[arg.Rs1] + uint64(arg.Imm)
		value, err := bus.Read(addr, 4)
		if err != nil {
			return err
		}
		p.RegX[arg.Rd] = uint64(int32(value))
	case riscv.ALBU:
		addr := p.RegX[arg.Rs1] + uint64(arg.Imm)
		value, err := bus.Read(addr, 1)
		if err != nil {
			return err
		}
		p.RegX[arg.Rd] = value
	case riscv.ALHU:
		addr := p.RegX[arg.Rs1] + uint64(arg.Imm)
		value, err := bus.Read(addr, 2)
		if err != nil {
			return err
		}
		p.RegX[arg.Rd] = value
	case riscv.ASB:
		addr := p.RegX[arg.Rs1] + uint64(arg.Imm)
		value := p.RegX[arg.Rs2]
		if err := bus.Write(addr, 1, value); err != nil {
			return err
		}
	case riscv.ASH:
		addr := p.RegX[arg.Rs1] + uint64(arg.Imm)
		value := p.RegX[arg.Rs2]
		if err := bus.Write(addr, 2, value); err != nil {
			return err
		}
	case riscv.ASW:
		addr := p.RegX[arg.Rs1] + uint64(arg.Imm)
		value := p.RegX[arg.Rs2]
		if err := bus.Write(addr, 4, value); err != nil {
			return err
		}
	case riscv.AADDI:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] + uint64(arg.Imm)
	case riscv.ASLTI:
		if int64(p.RegX[arg.Rs1]) < int64(arg.Imm) {
			p.RegX[arg.Rd] = 1
		} else {
			p.RegX[arg.Rd] = 0
		}
	case riscv.ASLTIU:
		if p.RegX[arg.Rs1] < uint64(arg.Imm) {
			p.RegX[arg.Rd] = 1
		} else {
			p.RegX[arg.Rd] = 0
		}
	case riscv.AXORI:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] ^ uint64(arg.Imm)
	case riscv.AORI:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] | uint64(arg.Imm)
	case riscv.AANDI:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] & uint64(arg.Imm)
	case riscv.ASLLI:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] << arg.Imm
	case riscv.ASRLI:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] >> arg.Imm
	case riscv.ASRAI:
		p.RegX[arg.Rd] = uint64(int64(p.RegX[arg.Rs1]) >> arg.Imm)
	case riscv.AADD:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] + p.RegX[arg.Rs2]
	case riscv.ASUB:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] - p.RegX[arg.Rs2]
	case riscv.ASLL:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] << uint64(arg.Imm)
	case riscv.ASLT:
		if int64(p.RegX[arg.Rs1]) < int64(p.RegX[arg.Rs2]) {
			p.RegX[arg.Rd] = 1
		} else {
			p.RegX[arg.Rd] = 0
		}
	case riscv.ASLTU:
		if p.RegX[arg.Rs1] < p.RegX[arg.Rs2] {
			p.RegX[arg.Rd] = 1
		} else {
			p.RegX[arg.Rd] = 0
		}
	case riscv.AXOR:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] ^ p.RegX[arg.Rs2]
	case riscv.ASRL:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] >> uint64(arg.Imm)
	case riscv.ASRA:
		p.RegX[arg.Rd] = uint64(int64(p.RegX[arg.Rs1]) >> uint64(arg.Imm))
	case riscv.AOR:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] | p.RegX[arg.Rs2]
	case riscv.AAND:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] & p.RegX[arg.Rs2]
	case riscv.AFENCE:
		// NOP
	case riscv.AECALL:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as))
	case riscv.AEBREAK:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as))

	// RV64I Base Instruction Set (in addition to RV32I)

	case riscv.ALWU:
		addr := p.RegX[arg.Rs1] + uint64(arg.Imm)
		value, err := bus.Read(addr, 4)
		if err != nil {
			return err
		}
		p.RegX[arg.Rd] = value
	case riscv.ALD:
		addr := p.RegX[arg.Rs1] + uint64(arg.Imm)
		value, err := bus.Read(addr, 8)
		if err != nil {
			return err
		}
		p.RegX[arg.Rd] = value
	case riscv.ASD:
	case riscv.AADDIW:
		p.RegX[arg.Rd] = uint64(int32(p.RegX[arg.Rs1] + uint64(arg.Imm)))
	case riscv.ASLLIW:
		p.RegX[arg.Rd] = uint64(int32(p.RegX[arg.Rs1] << arg.Imm))
	case riscv.ASRLIW:
		p.RegX[arg.Rd] = uint64(int32(uint32(p.RegX[arg.Rs1]) >> arg.Imm))
	case riscv.ASRAIW:
		p.RegX[arg.Rd] = uint64(int32(p.RegX[arg.Rs1]) >> arg.Imm)

	case riscv.AADDW:
		p.RegX[arg.Rd] = uint64(int32(p.RegX[arg.Rs1]) + int32(p.RegX[arg.Rs2]))
	case riscv.ASUBW:
		p.RegX[arg.Rd] = uint64(int32(p.RegX[arg.Rs1]) - int32(p.RegX[arg.Rs2]))
	case riscv.ASLLW:
		p.RegX[arg.Rd] = uint64(int32(p.RegX[arg.Rs1] << arg.Imm))
	case riscv.ASRLW:
		p.RegX[arg.Rd] = uint64(int32(uint32(p.RegX[arg.Rs1]) << arg.Imm))

	case riscv.ASRAW:
		p.RegX[arg.Rd] = uint64(int32(p.RegX[arg.Rs1] >> arg.Imm))

	// RV32/RV64 Zicsr Standard Extension

	case riscv.ACSRRW:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as))
	case riscv.ACSRRS:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as))
	case riscv.ACSRRC:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as))
	case riscv.ACSRRWI:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as))
	case riscv.ACSRRSI:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as))
	case riscv.ACSRRCI:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as))

	// RV32M Standard Extension

	case riscv.AMUL:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] * p.RegX[arg.Rs2]
	case riscv.AMULH:
		return fmt.Errorf("%s: unsupport", riscv.AsString(as))
	case riscv.AMULHSU:
		return fmt.Errorf("%s: unsupport", riscv.AsString(as))
	case riscv.AMULHU:
		return fmt.Errorf("%s: unsupport", riscv.AsString(as))
	case riscv.ADIV:
		if p.RegX[arg.Rs2] != 0 {
			p.RegX[arg.Rd] = uint64(int64(p.RegX[arg.Rs1]) / int64(p.RegX[arg.Rs2]))
		} else {
			v := int64(-1)
			p.RegX[arg.Rd] = uint64(v)
		}
	case riscv.ADIVU:
		if p.RegX[arg.Rs2] != 0 {
			p.RegX[arg.Rd] = p.RegX[arg.Rs1] / p.RegX[arg.Rs2]
		} else {
			v := int64(-1)
			p.RegX[arg.Rd] = uint64(v)
		}

	case riscv.AREM:
		if p.RegX[arg.Rs2] != 0 {
			p.RegX[arg.Rd] = uint64(int32(p.RegX[arg.Rs1]) % int32(p.RegX[arg.Rs2]))
		} else {
			p.RegX[arg.Rd] = p.RegX[arg.Rs1]
		}
	case riscv.AREMU:
		if p.RegX[arg.Rs2] != 0 {
			p.RegX[arg.Rd] = p.RegX[arg.Rs1] % p.RegX[arg.Rs2]
		} else {
			p.RegX[arg.Rd] = p.RegX[arg.Rs1]
		}

	// RV64M Standard Extension (in addition to RV32M)

	case riscv.AMULW:
		p.RegX[arg.Rd] = uint64(int32(p.RegX[arg.Rs1]) * int32(p.RegX[arg.Rs2]))
	case riscv.ADIVW:
		if p.RegX[arg.Rs2] != 0 {
			p.RegX[arg.Rd] = uint64(int32(p.RegX[arg.Rs1]) / int32(p.RegX[arg.Rs2]))
		} else {
			v := int64(-1)
			p.RegX[arg.Rd] = uint64(v)
		}
	case riscv.ADIVUW:
		if p.RegX[arg.Rs2] != 0 {
			p.RegX[arg.Rd] = uint64(uint32(p.RegX[arg.Rs1]) / uint32(p.RegX[arg.Rs2]))
		} else {
			v := int64(-1)
			p.RegX[arg.Rd] = uint64(v)
		}
	case riscv.AREMW:
		if p.RegX[arg.Rs2] != 0 {
			p.RegX[arg.Rd] = uint64(int32(p.RegX[arg.Rs1]) % int32(p.RegX[arg.Rs2]))
		} else {
			p.RegX[arg.Rd] = p.RegX[arg.Rs1]
		}
	case riscv.AREMUW:
		if p.RegX[arg.Rs2] != 0 {
			p.RegX[arg.Rd] = uint64(uint32(p.RegX[arg.Rs1]) % uint32(p.RegX[arg.Rs2]))
		} else {
			p.RegX[arg.Rd] = uint64(uint32(p.RegX[arg.Rs1]))
		}
	}
	return nil
}
