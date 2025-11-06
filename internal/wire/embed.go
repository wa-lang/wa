// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
	"strings"
)

/**************************************
本文件包含了 wire 中的辅助对象
**************************************/

//-------------------------------------

/**************************************
anInstruction: 实现 Instruction 接口与 Pos、Scope 相关的方法
**************************************/
type anInstruction struct {
	fmt.Stringer
	pos   int
	scope Scope
}

func (v *anInstruction) ParentScope() Scope { return v.scope }
func (v *anInstruction) setScope(s Scope)   { v.scope = s }
func (v *anInstruction) Pos() int           { return v.pos }
func (v *anInstruction) Format(tab string, sb *strings.Builder) {
	sb.WriteString(tab)
	if va, ok := v.Stringer.(Value); ok {
		sb.WriteString(va.Name())
		sb.WriteString(" = ")
	}
	sb.WriteString(v.String())
}

type imv interface {
	setId(id int)
}

/**************************************
imv: InterMediateValue，中间值，实现 Value 接口与 Name、Kind、Type 相关的方法
**************************************/
type aImv struct {
	anInstruction
	id  int       // 在函数内虚拟寄存器数组中的下标，在 Function.EndBody() 前该值无意义
	typ ValueType // 值类型
}

func (v *aImv) Name() string    { return fmt.Sprintf("$t%d", v.id) }
func (v *aImv) Kind() ValueKind { return ValueKindLocal }
func (v *aImv) Type() ValueType { return v.typ }
func (v *aImv) setId(id int)    { v.id = id }
