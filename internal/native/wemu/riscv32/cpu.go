// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 重要: 请保持 RV32 和 RV64 版本完全一致, RVUInt 在另一个文件定义.

package riscv32

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/riscv"
	"wa-lang.org/wa/internal/native/wemu/device"
)

// 处理器
type CPU struct {
	Mode riscv.Mode   // 工作模式
	RegX [32]RVUInt   // 整数寄存器
	RegF [32]float64  // 浮点数寄存器
	CSR  [4096]RVUInt // 状态寄存器
	PC   RVUInt       // PC指针
}

// 构建处理器
func NewCPU() *CPU {
	return new(CPU)
}

func (p *CPU) GetPC() uint64  { return uint64(p.PC) }
func (p *CPU) SetPC(v uint64) { p.PC = RVUInt(v) }

func (p *CPU) XRegNum() int            { return len(p.RegX) }
func (p *CPU) GetXReg(i int) uint64    { return uint64(p.RegX[i]) }
func (p *CPU) SetXReg(i int, v uint64) { p.RegX[i] = RVUInt(v) }

func (p *CPU) FRegNum() int             { return len(p.RegF) }
func (p *CPU) GetFReg(i int) float64    { return p.RegF[i] }
func (p *CPU) SetFReg(i int, v float64) { p.RegF[i] = v }

// 重置CPU
func (p *CPU) Reset(pc, sp uint64) {
	*p = CPU{}
	p.RegX[riscv.RegI(riscv.REG_SP)] = RVUInt(sp)
	p.PC = RVUInt(pc)
}

// 单步执行
func (p *CPU) StepRun(bus *device.Bus) error {
	inst, err := bus.Read(uint64(p.PC), 4)
	if err != nil {
		return err
	}
	as, arg, argRaw, err := riscv.DecodeEx(uint32(inst))
	if err != nil {
		return fmt.Errorf("fetch instruntion failed at 0x%08X: %v", p.PC, err)
	}

	if err := p.execInst(bus, as, argRaw); err != nil {
		return fmt.Errorf("exec instruntion(%08d:%s) failed at 0x%08X: %v", inst,
			riscv.AsmSyntax(as, "", arg),
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

	// 调整PC(少数跳转指令会覆盖)
	p.PC += 4

	// 执行指令
	switch as {
	default:
		return fmt.Errorf("unsupport: %s", riscv.AsString(as, ""))

	// RV32I Base Instruction Set

	case riscv.ALUI:
		p.RegX[arg.Rd] = RVUInt(arg.Imm << 12)
	case riscv.AAUIPC:
		p.RegX[arg.Rd] = curPC + RVUInt(arg.Imm)
	case riscv.AJAL:
		// rd 寄存器先保存下一个指令对应的 PC
		p.RegX[arg.Rd] = p.PC
		// 然后根据当前指令对应的 PC 计算出跳转地址覆盖当前的 PC
		p.PC = curPC + RVUInt(arg.Imm)
	case riscv.AJALR:
		// rd 寄存器先保存下一个指令对应的 PC
		p.RegX[arg.Rd] = p.PC
		// 然后根据计算出跳转地址覆盖当前的 PC
		p.PC = p.RegX[arg.Rs1] + RVUInt(arg.Imm)
	case riscv.ABEQ:
		if p.RegX[arg.Rs1] == p.RegX[arg.Rs2] {
			p.PC = curPC + RVUInt(arg.Imm)
		}
	case riscv.ABNE:
		if p.RegX[arg.Rs1] != p.RegX[arg.Rs2] {
			p.PC = curPC + RVUInt(arg.Imm)
		}
	case riscv.ABLT:
		if int64(p.RegX[arg.Rs1]) < int64(p.RegX[arg.Rs2]) {
			p.PC = curPC + RVUInt(arg.Imm)
		}
	case riscv.ABGE:
		if int64(p.RegX[arg.Rs1]) > int64(p.RegX[arg.Rs2]) {
			p.PC = curPC + RVUInt(arg.Imm)
		}
	case riscv.ABLTU:
		if p.RegX[arg.Rs1] < p.RegX[arg.Rs2] {
			p.PC = curPC + RVUInt(arg.Imm)
		}
	case riscv.ABGEU:
		if p.RegX[arg.Rs1] > p.RegX[arg.Rs2] {
			p.PC = curPC + RVUInt(arg.Imm)
		}
	case riscv.ALB:
		addr := p.RegX[arg.Rs1] + RVUInt(arg.Imm)
		value, err := bus.Read(uint64(addr), 1)
		if err != nil {
			return err
		}
		p.RegX[arg.Rd] = RVUInt(int8(value))
	case riscv.ALH:
		addr := p.RegX[arg.Rs1] + RVUInt(arg.Imm)
		value, err := bus.Read(uint64(addr), 2)
		if err != nil {
			return err
		}
		p.RegX[arg.Rd] = RVUInt(int16(value))
	case riscv.ALW:
		addr := p.RegX[arg.Rs1] + RVUInt(arg.Imm)
		value, err := bus.Read(uint64(addr), 4)
		if err != nil {
			return err
		}
		p.RegX[arg.Rd] = RVUInt(int32(value))
	case riscv.ALBU:
		addr := p.RegX[arg.Rs1] + RVUInt(arg.Imm)
		value, err := bus.Read(uint64(addr), 1)
		if err != nil {
			return err
		}
		p.RegX[arg.Rd] = RVUInt(value)
	case riscv.ALHU:
		addr := p.RegX[arg.Rs1] + RVUInt(arg.Imm)
		value, err := bus.Read(uint64(addr), 2)
		if err != nil {
			return err
		}
		p.RegX[arg.Rd] = RVUInt(value)
	case riscv.ASB:
		addr := p.RegX[arg.Rs1] + RVUInt(arg.Imm)
		value := p.RegX[arg.Rs2]
		if err := bus.Write(uint64(addr), 1, uint64(value)); err != nil {
			return err
		}
	case riscv.ASH:
		addr := p.RegX[arg.Rs1] + RVUInt(arg.Imm)
		value := p.RegX[arg.Rs2]
		if err := bus.Write(uint64(addr), 2, uint64(value)); err != nil {
			return err
		}
	case riscv.ASW:
		addr := p.RegX[arg.Rs1] + RVUInt(arg.Imm)
		value := p.RegX[arg.Rs2]
		if err := bus.Write(uint64(addr), 4, uint64(value)); err != nil {
			return err
		}
	case riscv.AADDI:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] + RVUInt(arg.Imm)
	case riscv.ASLTI:
		if int64(p.RegX[arg.Rs1]) < int64(arg.Imm) {
			p.RegX[arg.Rd] = 1
		} else {
			p.RegX[arg.Rd] = 0
		}
	case riscv.ASLTIU:
		if p.RegX[arg.Rs1] < RVUInt(arg.Imm) {
			p.RegX[arg.Rd] = 1
		} else {
			p.RegX[arg.Rd] = 0
		}
	case riscv.AXORI:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] ^ RVUInt(arg.Imm)
	case riscv.AORI:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] | RVUInt(arg.Imm)
	case riscv.AANDI:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] & RVUInt(arg.Imm)
	case riscv.ASLLI:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] << arg.Imm
	case riscv.ASRLI:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] >> arg.Imm
	case riscv.ASRAI:
		p.RegX[arg.Rd] = RVUInt(int64(p.RegX[arg.Rs1]) >> arg.Imm)
	case riscv.AADD:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] + p.RegX[arg.Rs2]
	case riscv.ASUB:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] - p.RegX[arg.Rs2]
	case riscv.ASLL:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] << RVUInt(arg.Imm)
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
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] >> RVUInt(arg.Imm)
	case riscv.ASRA:
		p.RegX[arg.Rd] = RVUInt(int64(p.RegX[arg.Rs1]) >> RVUInt(arg.Imm))
	case riscv.AOR:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] | p.RegX[arg.Rs2]
	case riscv.AAND:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] & p.RegX[arg.Rs2]
	case riscv.AFENCE:
		// NOP
	case riscv.AECALL:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as, ""))
	case riscv.AEBREAK:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as, ""))

	// RV64I Base Instruction Set (in addition to RV32I)

	case riscv.ALWU:
		addr := p.RegX[arg.Rs1] + RVUInt(arg.Imm)
		value, err := bus.Read(uint64(addr), 4)
		if err != nil {
			return err
		}
		p.RegX[arg.Rd] = RVUInt(value)
	case riscv.ALD:
		addr := p.RegX[arg.Rs1] + RVUInt(arg.Imm)
		value, err := bus.Read(uint64(addr), 8)
		if err != nil {
			return err
		}
		p.RegX[arg.Rd] = RVUInt(value)
	case riscv.ASD:
	case riscv.AADDIW:
		p.RegX[arg.Rd] = RVUInt(int32(p.RegX[arg.Rs1] + RVUInt(arg.Imm)))
	case riscv.ASLLIW:
		p.RegX[arg.Rd] = RVUInt(int32(p.RegX[arg.Rs1] << arg.Imm))
	case riscv.ASRLIW:
		p.RegX[arg.Rd] = RVUInt(int32(uint32(p.RegX[arg.Rs1]) >> arg.Imm))
	case riscv.ASRAIW:
		p.RegX[arg.Rd] = RVUInt(int32(p.RegX[arg.Rs1]) >> arg.Imm)

	case riscv.AADDW:
		p.RegX[arg.Rd] = RVUInt(int32(p.RegX[arg.Rs1]) + int32(p.RegX[arg.Rs2]))
	case riscv.ASUBW:
		p.RegX[arg.Rd] = RVUInt(int32(p.RegX[arg.Rs1]) - int32(p.RegX[arg.Rs2]))
	case riscv.ASLLW:
		p.RegX[arg.Rd] = RVUInt(int32(p.RegX[arg.Rs1] << arg.Imm))
	case riscv.ASRLW:
		p.RegX[arg.Rd] = RVUInt(int32(uint32(p.RegX[arg.Rs1]) << arg.Imm))

	case riscv.ASRAW:
		p.RegX[arg.Rd] = RVUInt(int32(p.RegX[arg.Rs1] >> arg.Imm))

	// RV32/RV64 Zicsr Standard Extension

	case riscv.ACSRRW:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as, ""))
	case riscv.ACSRRS:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as, ""))
	case riscv.ACSRRC:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as, ""))
	case riscv.ACSRRWI:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as, ""))
	case riscv.ACSRRSI:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as, ""))
	case riscv.ACSRRCI:
		// 涉及 CSR, 暂不支持
		return fmt.Errorf("%s: unsupport", riscv.AsString(as, ""))

	// RV32M Standard Extension

	case riscv.AMUL:
		p.RegX[arg.Rd] = p.RegX[arg.Rs1] * p.RegX[arg.Rs2]
	case riscv.AMULH:
		return fmt.Errorf("%s: unsupport", riscv.AsString(as, ""))
	case riscv.AMULHSU:
		return fmt.Errorf("%s: unsupport", riscv.AsString(as, ""))
	case riscv.AMULHU:
		return fmt.Errorf("%s: unsupport", riscv.AsString(as, ""))
	case riscv.ADIV:
		if p.RegX[arg.Rs2] != 0 {
			p.RegX[arg.Rd] = RVUInt(int64(p.RegX[arg.Rs1]) / int64(p.RegX[arg.Rs2]))
		} else {
			v := int64(-1)
			p.RegX[arg.Rd] = RVUInt(v)
		}
	case riscv.ADIVU:
		if p.RegX[arg.Rs2] != 0 {
			p.RegX[arg.Rd] = p.RegX[arg.Rs1] / p.RegX[arg.Rs2]
		} else {
			v := int64(-1)
			p.RegX[arg.Rd] = RVUInt(v)
		}

	case riscv.AREM:
		if p.RegX[arg.Rs2] != 0 {
			p.RegX[arg.Rd] = RVUInt(int32(p.RegX[arg.Rs1]) % int32(p.RegX[arg.Rs2]))
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
		p.RegX[arg.Rd] = RVUInt(int32(p.RegX[arg.Rs1]) * int32(p.RegX[arg.Rs2]))
	case riscv.ADIVW:
		if p.RegX[arg.Rs2] != 0 {
			p.RegX[arg.Rd] = RVUInt(int32(p.RegX[arg.Rs1]) / int32(p.RegX[arg.Rs2]))
		} else {
			v := int64(-1)
			p.RegX[arg.Rd] = RVUInt(v)
		}
	case riscv.ADIVUW:
		if p.RegX[arg.Rs2] != 0 {
			p.RegX[arg.Rd] = RVUInt(uint32(p.RegX[arg.Rs1]) / uint32(p.RegX[arg.Rs2]))
		} else {
			v := int64(-1)
			p.RegX[arg.Rd] = RVUInt(v)
		}
	case riscv.AREMW:
		if p.RegX[arg.Rs2] != 0 {
			p.RegX[arg.Rd] = RVUInt(int32(p.RegX[arg.Rs1]) % int32(p.RegX[arg.Rs2]))
		} else {
			p.RegX[arg.Rd] = p.RegX[arg.Rs1]
		}
	case riscv.AREMUW:
		if p.RegX[arg.Rs2] != 0 {
			p.RegX[arg.Rd] = RVUInt(uint32(p.RegX[arg.Rs1]) % uint32(p.RegX[arg.Rs2]))
		} else {
			p.RegX[arg.Rd] = RVUInt(uint32(p.RegX[arg.Rs1]))
		}
	}
	return nil
}

func (p *CPU) InstString(x uint32) (string, error) {
	as, arg, err := riscv.Decode(x)
	if err != nil {
		return "", err
	}
	return riscv.AsmSyntax(as, "", arg), nil
}

func (p *CPU) LookupRegister(regName string) (r abi.RegType, ok bool) {
	return riscv.LookupRegister(regName)
}

func (p *CPU) RegAliasString(r abi.RegType) string {
	return riscv.RegAliasString(r)
}
