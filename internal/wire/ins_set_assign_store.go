// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
	"strings"
)

/**************************************
Set: Set 指令，实现赋值 `=` 操作，将值保存至变量、或地址
  - Set 支持多赋值
  - Lh 应为 Location；Lh 为 nil 是合法的，这等价于向匿名变量 _ 赋值
  - Rh 可能为元组（Tuple），此时 Rhs 的长度应为 1，Lhs 的长度应为元组长度
**************************************/
type Set struct {
	aStmt
	Lhs     []Location
	Rhs     []Expr
	rhsType []Type
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

		var loc_name string
		if v, ok := lh.(Var); ok && v.Kind() == AllocKindRegister {
			loc_name = lh.Name()
		} else {
			loc_name = "*(" + lh.Name() + ")"
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

func NewSet(lhs Location, rhs Expr, pos int) *Set {
	return NewSetN([]Location{lhs}, []Expr{rhs}, pos)
}

func NewSetN(lhs []Location, rhs []Expr, pos int) *Set {
	v := &Set{}
	v.Stringer = v
	v.Lhs = lhs
	v.Rhs = rhs

	v.rhsType = make([]Type, len(lhs))
	if len(lhs) != len(rhs) {
		if len(rhs) != 1 {
			panic("rhs is not a tuple")
		}
		tuple := rhs[0].Type().(*Tuple)
		if len(tuple.members) != len(lhs) {
			panic("tuple members length mismatch")
		}
		copy(v.rhsType, tuple.members)
	} else {
		for i := range rhs {
			v.rhsType[i] = rhs[i].Type()
		}
	}
	v.pos = pos
	return v
}

// 在 Block 中添加一条 Set 指令
func (b *Block) EmitSet(lhs Location, rhs Expr, pos int) *Set {
	v := NewSet(lhs, rhs, pos)
	b.emit(v)
	return v
}

// Block.EmitSet 的多重赋值版
func (b *Block) EmitSetN(lhs []Location, rhs []Expr, pos int) *Set {
	v := NewSetN(lhs, rhs, pos)
	b.emit(v)
	return v
}

/**************************************
Assign: Assign 指令，向 Register 变量赋值
  - Assgin 支持多赋值
  - Lh 必须为 Register 变量 或 nil；Lh 为 nil 时等价于向匿名变量 _ 赋值
  - Rh 可能为元组（Tuple），此时 Rhs 的长度应为 1，Lhs 的长度应为元组长度
  - 该指令仅供内部使用，上层高级语法不应直接使用
**************************************/
type Assign struct {
	aStmt
	Lhs     []Var
	Rhs     []Expr
	rhsType []Type
}

func (i *Assign) String() string {
	var sb strings.Builder

	var rh Expr
	for j, lh := range i.Lhs {
		if j > 0 {
			sb.WriteString(", ")
		}

		if j < len(i.Rhs) {
			rh = i.Rhs[j]
		}

		if lh == nil {
			sb.WriteRune('_')
			if _, ok := rh.(*Const); !ok && rh.retained() && rtimp.hasChunk(i.rhsType[j]) {
				sb.WriteString("↓")
			}
		} else {
			sb.WriteString(lh.Name())
			if _, ok := rh.(*Const); !ok && !rh.retained() && rtimp.hasChunk(i.rhsType[j]) {
				sb.WriteString("↑")
			}
		}
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

func newAssign(lh Var, rh Expr, pos int) *Assign {
	return newAssignN([]Var{lh}, []Expr{rh}, pos)
}

func newAssignN(lhs []Var, rhs []Expr, pos int) *Assign {
	v := &Assign{}
	v.Stringer = v
	v.Lhs = lhs
	v.Rhs = rhs
	v.pos = pos

	v.rhsType = make([]Type, len(lhs))
	if len(lhs) != len(rhs) {
		if len(rhs) != 1 {
			panic("rhs is not a tuple")
		}
		tuple := rhs[0].Type().(*Tuple)
		if len(tuple.members) != len(lhs) {
			panic("tuple members length mismatch")
		}
		copy(v.rhsType, tuple.members)
	} else {
		for i := range rhs {
			v.rhsType[i] = rhs[i].Type()
		}
	}
	return v
}

/**************************************
Store: Store 指令，将 Val 存储到 Loc 指定的位置
  - Loc 应为 Var
  - Val 应为 Var 或 Const
  - 该指令仅供内部使用，上层高级语法不应直接使用
**************************************/
type Store struct {
	aStmt
	Loc Var
	Val Expr
}

func (i *Store) String() string {
	if _, ok := i.Val.(*Const); !ok && !i.Val.retained() && rtimp.hasChunk(i.Val.Type()) {
		return fmt.Sprintf("store(%s↑, %s)", i.Loc.Name(), i.Val.Name())
	} else {
		return fmt.Sprintf("store(%s, %s)", i.Loc.Name(), i.Val.Name())
	}
}

func newStore(loc Var, val Expr, pos int) *Store {
	if val == nil {
		panic("val is nil")
	}

	if _, ok := val.(Var); !ok {
		if _, ok := val.(*Const); !ok {
			panic(fmt.Sprintf("val is not a Var or Const: %s", val.Name()))
		}
	}

	v := &Store{}
	v.Stringer = v
	v.Loc = loc
	v.Val = val
	v.pos = pos
	return v
}
