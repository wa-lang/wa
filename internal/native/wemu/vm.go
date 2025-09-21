// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wemu

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/riscv"
	"wa-lang.org/wa/internal/native/wemu/device"
	"wa-lang.org/wa/internal/native/wemu/device/dram"
	"wa-lang.org/wa/internal/native/wemu/device/power"
	"wa-lang.org/wa/internal/native/wemu/device/uart"
	"wa-lang.org/wa/internal/native/wemu/riscv64"
)

// 外设地址
const (
	KDeviceAddr_power = 0
	KDeviceAddr_uart  = 0
)

// 模拟器
type WEmu struct {
	CPU   riscv64.CPU        // 处理器
	Bus   *device.Bus        // 外设总线
	Power *power.Power       // 电源设备
	Dram  *dram.DRAM         // 内存设备
	Uart  *uart.UART         // 串口设备
	Prog  *abi.LinkedProgram // 程序
}

// 构建模拟器
func NewWEmu(prog *abi.LinkedProgram) *WEmu {
	p := &WEmu{
		Bus:   device.NewBus(),
		Power: power.NewPower("power", power.POWER_BASE),
		Dram:  dram.NewDRAM("memory", dram.DRAM_BASE, dram.DRAM_SIZE, false),
		Uart:  uart.NewUART("uart", uart.UART_BASE),
		Prog:  prog,
	}

	// 映射总线设备
	p.Bus.MapDevice(p.Power)
	p.Bus.MapDevice(p.Dram)
	p.Bus.MapDevice(p.Uart)
	return p
}

// 重新加载程序
func (p *WEmu) resetProgram() error {
	// 重新加载指令段
	if err := p.Dram.Fill(uint64(p.Prog.TextAddr), p.Prog.TextData); err != nil {
		return err
	}
	// 重新加载数据段
	if len(p.Prog.DataData) != 0 {
		if err := p.Dram.Fill(uint64(p.Prog.DataAddr), p.Prog.DataData); err != nil {
			return err
		}
	}

	// 重新设置PC和SP
	p.CPU.Reset(uint64(p.Prog.TextAddr), p.Dram.AddrEnd())
	return nil
}

// 运行程序
func (p *WEmu) Run() error {
	if err := p.resetProgram(); err != nil {
		return err
	}
	for {
		if p.Power.IsShutdown() {
			break
		}
		if err := p.CPU.StepRun(p.Bus); err != nil {
			return err
		}
	}
	if p.Power.Status() == power.ExitFail {
		return fmt.Errorf("power off")
	}
	return nil
}

// 调试模式执行
func (p *WEmu) DebugRun() error {
	if err := p.resetProgram(); err != nil {
		return err
	}

	var (
		stepcnt int
		pntflag bool
		traflag bool
	)

	stdin := bufio.NewReader(os.Stdin)

	fmt.Println("Debug (enter h for help)...")
	fmt.Println()

	for {
		fmt.Print("Enter command: ")
		line, _, _ := stdin.ReadLine()

		// 删除空白字符
		line = bytes.TrimSpace(line)

		// 跳过空白行
		if string(line) == "" {
			fmt.Println()
			continue
		}

		var cmd, x1, x2 = "", 0, 0
		n, _ := fmt.Fscanf(bytes.NewBuffer(line), "%s%x%x", &cmd, &x1, &x2)

		switch cmd {
		case "help", "h":
			fmt.Println(p.DebugHelp())
		case "go", "g":
			if p.Power.IsShutdown() {
				fmt.Println("halted, enter `clear` to reset VM")
				continue
			}
			stepcnt = 0
			for !p.Power.IsShutdown() {
				stepcnt++
				if traflag {
					fmt.Print(p.FormatInstruction(p.CPU.PC, 1))
				}

				// 单步执行(可能执行HALT关机指令)
				p.CPU.StepRun(p.Bus)
			}
			if pntflag {
				fmt.Printf("step count = %d\n", stepcnt)
			}

		case "step", "s":
			if p.Power.IsShutdown() {
				fmt.Println("halted, enter `clear` to reset VM")
				continue
			}

			if n >= 2 {
				stepcnt = x1
			} else {
				stepcnt = 1
			}

			var i int
			for i = 0; i < stepcnt && !p.Power.IsShutdown(); i++ {
				if traflag {
					fmt.Print(p.FormatInstruction(p.CPU.PC, 1))
				}

				// 单步执行(可能执行关机指令)
				p.CPU.StepRun(p.Bus)
			}
			if pntflag {
				fmt.Printf("step count = %d\n", i)
			}

		case "jump", "j":
			if n >= 2 {
				p.CPU.PC = uint64(x1)
				fmt.Printf("PC = 0x%08X\n", x1)
			} else {
				fmt.Println("invalid command")
			}

		case "xregs", "x":
			regStart := 0
			regNum := len(p.CPU.RegX)
			if n > 1 {
				regStart = x1
			}
			if n > 2 {
				regNum = x2
			}

			fmt.Printf("PC  = 0x%08X\n", p.CPU.PC)
			for i := regStart; i < len(p.CPU.RegX) && i < (regStart+regNum); i++ {
				reg, ok := riscv.LookupRegister(fmt.Sprintf("X%d", i))
				if !ok {
					break
				}
				fmt.Printf("X%-2d = 0x%08X # %s\n",
					i, p.CPU.RegX[i], riscv.RegAliasString(reg),
				)
			}

		case "fregs", "f":
			regStart := 0
			regNum := len(p.CPU.RegF)
			if n > 1 {
				regStart = x1
			}
			if n > 2 {
				regNum = x2
			}

			fmt.Printf("PC  = 0x%08X\n", p.CPU.PC)
			for i := regStart; i < len(p.CPU.RegX) && i < (regStart+regNum); i++ {
				reg, ok := riscv.LookupRegister(fmt.Sprintf("F%d", i))
				if !ok {
					break
				}
				fmt.Printf("F%-2d = %v (0x%08X) # %s\n",
					i, p.CPU.RegF[i], math.Float64bits(p.CPU.RegF[i]),
					riscv.RegAliasString(reg),
				)
			}

		case "iMem", "imem", "i":
			pc := p.CPU.PC
			cnt := 1
			if n >= 2 {
				pc = uint64(x1)
			}
			if n >= 3 {
				cnt = x2
			}

			fmt.Print(p.FormatInstruction(pc, cnt))

		case "dMem", "dmem", "d":
			pc := p.CPU.PC
			cnt := 1

			if n >= 2 {
				pc = uint64(x1)
			}
			if n >= 3 {
				cnt = x2
			}

			for i := 0; i < cnt; i++ {
				addr := pc + uint64(i)*4
				v, err := p.Dram.Read(addr, 4)
				if err != nil {
					fmt.Printf("mem[%08X] = ERR:%v\n", addr, err)
					break
				}
				fmt.Printf("mem[%08x] = %08x\n", addr, v)
				x1++
			}

		case "alter", "a":
			if n == 3 {
				fmt.Printf("mem[%x] = %x\n", x1, x2)
				if err := p.Dram.Write(uint64(x1), 4, uint64(x2)); err != nil {
					fmt.Println("ERR:", err)
				}
			} else {
				fmt.Println("invalid command")
			}

		case "trace", "t":
			traflag = !traflag
			if traflag {
				fmt.Println("trace instruction now on")
			} else {
				fmt.Println("trace instruction now off")
			}

		case "print", "p":
			pntflag = !pntflag
			if pntflag {
				fmt.Println("Printing instruction count now on")
			} else {
				fmt.Println("Printing instruction count now off")
			}

		case "clear", "c":
			if err := p.resetProgram(); err != nil {
				fmt.Println("ERR:", err)
			}
			stepcnt = 0

		case "quit", "q":
			return nil

		default:
			fmt.Println("unknown command:", cmd)
		}
	}
}

func (p *WEmu) DebugHelp() string {
	return `Commands are:
  h)elp           show help command list
  g)o             run instructions until power off
  s)tep  <n>      run n (default 1) instructions
  j)ump  <b>      jump to the b (default is current location)
  x)regs          print the contents of the int registers
  f)regs          print the contents of the float registers
  i)Mem  <b <n>>  print n iMem locations starting at b
  d)Mem  <b <n>>  print n dMem locations starting at b
  a)lter <b <v>>  change the memory value at b
  t)race          toggle instruction trace
  p)rint          toggle print of total instructions executed
  c)lear          reset VM
  q)uit           exit
`
}

// 格式化pc开始的n个指令
func (p *WEmu) FormatInstruction(pc uint64, n int) string {
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		addr := pc + uint64(i)*4
		inst, err := p.Bus.Read(addr, 4)
		if err != nil {
			fmt.Fprintf(&buf, "mem[%08X]: %v\n", addr, err)
			continue
		}
		as, arg, err := riscv.Decode(uint32(inst))
		if err != nil {
			fmt.Fprintf(&buf, "mem[%08X]: %v\n", addr, err)
			continue
		}
		fmt.Fprintf(&buf, "mem[%08X]: %s\n", addr,
			riscv.AsmSyntax(as, arg),
		)
	}

	return buf.String()
}
