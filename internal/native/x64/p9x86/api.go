// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package p9x86

// 初始化查询表格
func Init() {
	instinit()
}

// 编码指令
func AsmInst(p *Prog) []byte {
	var ab AsmBuf
	ab.Asmins(p)
	return ab.Bytes()
}

// 指令码
type As int16

// Prog 表示一个机器指令
type Prog struct {
	As       As     // 指令码
	From     Addr   // 源对象
	To       Addr   // 目的对象
	RestArgs []Addr // 其他的对象
	Ft       uint8  // for x86 back end: type index of Prog.From
	Tt       uint8  // for x86 back end: type index of Prog.To
}

// Addr 表示一个指令操作对象
type Addr struct {
	Type   AddrType
	Reg    int16
	Index  int16
	Scale  int16 // Sometimes holds a register.
	Offset int64
}

type AddrType uint8

const (
	TYPE_NONE   AddrType = iota
	TYPE_REG             // AX
	TYPE_MEM             // [rip+name]
	TYPE_CONST           // 123
	TYPE_BRANCH          // jmp
	TYPE_ADDR            // lea
)

func (p *Prog) GetFrom3() *Addr {
	if p.RestArgs == nil {
		return nil
	}
	return &p.RestArgs[0]
}
