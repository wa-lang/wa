// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wair

import "fmt"

/**************************************
本文件包含了 wair 中的辅助对象
**************************************/

//-------------------------------------

/**************************************
anInstruction: 实现 Instruction 接口与 Pos、parent 相关的方法
**************************************/
type anInstruction struct {
	pos    int
	parent *Block
}

func (v *anInstruction) Parent() *Block     { return v.parent }
func (v *anInstruction) setParent(p *Block) { v.parent = p }
func (v *anInstruction) Pos() int           { return v.pos }

/**************************************
imv: InterMediateValue，中间值，实现 Value 接口与 Type、Name 相关的方法
**************************************/
type imv struct {
	anInstruction
	id  int       // 在函数内虚拟寄存器数组中的下标，函数内唯一
	typ ValueType // 值类型
}

func (v *imv) Name() string    { return fmt.Sprintf("$t%d", v.id) }
func (v *imv) Type() ValueType { return v.typ }
