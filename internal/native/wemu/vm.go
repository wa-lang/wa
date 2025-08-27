package riscv64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/link"
	"wa-lang.org/wa/internal/native/wemu/device"
	"wa-lang.org/wa/internal/native/wemu/device/dram"
	"wa-lang.org/wa/internal/native/wemu/device/power"
	"wa-lang.org/wa/internal/native/wemu/device/uart"
)

// 外设地址
const (
	KDeviceAddr_power = 0
	KDeviceAddr_uart  = 0
)

// 处理器
type CoreRunner interface {
	Reset(pc, sp uint64)
	StepRun(bus *device.Bus) error
}

// 模拟器
type WEmu struct {
	CPU   CoreRunner
	Bus   *device.Bus  // 外设总线
	Power *power.Power // 电源设备
	Dram  *dram.DRAM   // 内存设备
	Uart  *uart.UART   // 串口设备
}

// 构建模拟器
func NewWEmu(cpu CoreRunner, memSize int) *WEmu {
	p := &WEmu{
		CPU:   cpu,
		Bus:   device.NewBus(),
		Power: power.NewPower("power"),
		Dram:  dram.NewDRAM("memory", uint64(memSize), false),
	}

	// 映射总线设备
	p.Bus.MapDevice(p.Power, power.POWER_BASE)
	p.Bus.MapDevice(p.Dram, dram.DRAM_BASE)
	p.Bus.MapDevice(p.Uart, uart.UART_BASE)

	return p
}

// 加载程序
func (p *WEmu) LoadProgram(prog *link.Program, pc, sp uint64) error {
	p.CPU.Reset(pc, sp)
	return nil
}

// 运行程序
func (p *WEmu) Run() error {
	for {
		switch p.Power.Status() {
		case power.ExitOK:
			return nil
		case power.ExitFail:
			return fmt.Errorf("power off")
		}
		if err := p.CPU.StepRun(p.Bus); err != nil {
			return err
		}
	}
}
