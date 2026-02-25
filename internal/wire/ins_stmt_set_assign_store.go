// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
	"strings"
)

/**************************************
Set: Set 指令，将 Val 存储到 Loc 指定的位置，Set 支持多赋值，该指令应触发 RC-1 动作
  - Set 支持多赋值
  - Lh 应为 Var、或类型为 Ref、Ptr 的 Expr
  - Rh 可能为元组（Tuple），此时 Rhs 的长度应为 1，Lhs 的长度应为元组长度
  - 向 nil 的 Lh 赋值是合法的，这等价于向匿名变量 _ 赋值，此时若 Rh 已 retain，应触发 release
**************************************/

type Set struct {
	aStmt
	Lhs []Expr
	Rhs []Expr
}

func (i *Set) String() string {
	var sb strings.Builder

	sb.WriteString("set([")
	for i, lh := range i.Lhs {
		if i > 0 {
			sb.WriteString(", ")
		}

		if lh == nil {
			sb.WriteRune('_')
			continue
		}

		loc_name := "*(" + lh.Name() + ")"
		if v, ok := lh.(Var); ok {
			if v.Kind() == Register {
				loc_name = lh.Name()
			}
		}
		sb.WriteString(loc_name)
	}
	sb.WriteString("] , [")

	for i, rh := range i.Rhs {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(rh.Name())
	}
	sb.WriteString("])")

	return sb.String()
}

func NewSet(lhs Expr, rhs Expr, pos int) *Set {
	return NewSetN([]Expr{lhs}, []Expr{rhs}, pos)
}

func NewSetN(lhs []Expr, rhs []Expr, pos int) *Set {
	v := &Set{}
	v.Stringer = v
	v.Lhs = lhs
	v.Rhs = rhs
	v.pos = pos
	return v
}

// 在 Block 中添加一条 Set 指令
func (b *Block) EmitSet(lhs Expr, rhs Expr, pos int) *Set {
	v := NewSet(lhs, rhs, pos)
	b.emit(v)
	return v
}

// Block.EmitSet 的多重赋值版
func (b *Block) EmitSetN(lhs []Expr, rhs []Expr, pos int) *Set {
	v := NewSetN(lhs, rhs, pos)
	b.emit(v)
	return v
}

/**************************************
Assign: Assign 指令，将 Rhs 赋值给 Lhs
  - Assgin 支持多赋值
  - Lh 必须为 Register 型变量
  - Rh 可能为元组（Tuple），此时 Rhs 的长度应为 1，Lhs 的长度应为元组长度
  - 向 nil 的 Lh 赋值是合法的，这等价于向匿名变量 _ 赋值，此时若 Rh 已 retain，应触发 release
**************************************/

type Assign struct {
	aStmt
	Lhs []Var
	Rhs []Expr
}

func (i *Assign) String() string {
	var sb strings.Builder

	for i, lh := range i.Lhs {
		if i > 0 {
			sb.WriteString(", ")
		}

		if lh == nil {
			sb.WriteRune('_')
			continue
		}

		sb.WriteString(lh.Name())
	}

	sb.WriteString(" = ")
	for i, rh := range i.Rhs {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(rh.Name())
	}

	return sb.String()
}

func NewAssign(lh Var, rh Expr, pos int) *Assign {
	return NewAssignN([]Var{lh}, []Expr{rh}, pos)
}

func NewAssignN(lhs []Var, rhs []Expr, pos int) *Assign {
	v := &Assign{}
	v.Stringer = v
	v.Lhs = lhs
	v.Rhs = rhs
	v.pos = pos
	return v
}

// 在 Block 中添加一条 Assign 指令
func (b *Block) EmitAssign(lhs Var, rhs Expr, pos int) *Assign {
	v := NewAssign(lhs, rhs, pos)
	b.emit(v)
	return v
}

// EmitAssign 的多重赋值版
func (b *Block) EmitAssignN(lhs []Var, rhs []Expr, pos int) *Assign {
	v := NewAssignN(lhs, rhs, pos)
	b.emit(v)
	return v
}

/**************************************
Store: Store 指令，将 Val 存储到 Loc 指定的位置
  - Loc 类型应为 Ref、Ptr
**************************************/

type Store struct {
	aStmt
	Loc Expr
	Val Expr
}

func (i *Store) String() string {
	return fmt.Sprintf("store(%s, %s)", i.Loc.Name(), i.Val.Name())
}

func NewStore(loc Expr, val Expr, pos int) *Store {
	v := &Store{}
	v.Stringer = v
	v.Loc = loc
	v.Val = val
	v.pos = pos
	return v
}
