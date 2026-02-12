// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	_ "embed"
	"encoding/csv"
	"strings"
)

// ----------------------------------------------------------------------------

//go:embed x86.csv
var x86_csv_data string

// CSV 文件说明:
// - 第1列: Mnemonic, 指令助记符, 如 "ADD, MOV, RET"
// - 第2列: Args, 操作数占位符, 如 "r/m64, r64", 定义了指令接受的操作数类型和位宽
// - 第3列: Encoding, 核心 Opcode 规则, 如 "01 /r, B8+rd io" 包含基础字节及 "ModR/M" 填充规则
// - 第4列: Flags, 指令属性标志, 描述指令对 EFLAGS 的影响或特定的 CPU 模式约束
// - 第5列: CPU Feature, 硬件扩展需求, 指明所属扩展集，如 "SSE2, AVX, VMX" 为空则通常是基础指令
// - 第6列: Description, 补充说明(可能为空)

// ----------------------------------------------------------------------------

// 指令操作码前缀决策树
type Prog struct {
	Path   string
	Action string
	Child  map[string]*Prog
	PC     int
	TailID int
}

func NewProg() *Prog {
	r := csv.NewReader(strings.NewReader(x86_csv_data))
	r.FieldsPerRecord = -1
	r.Comment = '#'

	table, err := r.ReadAll()
	if err != nil {
		panic(err)
	}
	if len(table) == 0 || len(table[0]) < 6 {
		panic("invalid x86.csv file")
	}

	p := new(Prog)
	for _, row := range table {
		p.add(row[0], row[1], row[2], row[3], row[4], row[5])
	}
	p.check()

	return p
}
