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

func (i *anInstruction) Parent() *Block     { return i.parent }
func (i *anInstruction) setParent(p *Block) { i.parent = p }
func (i *anInstruction) Pos() int           { return i.pos }

/**************************************
aRegister: 实现 Value 接口与 Type、Name 相关的方法
**************************************/
type aRegister struct {
	anInstruction
	id  int       // 在函数内虚拟寄存器数组中的下标，函数内唯一
	typ ValueType // 值类型
}

func (r *aRegister) Name() string    { return fmt.Sprintf("$t%d", r.id) }
func (r *aRegister) Type() ValueType { return r.typ }
