// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
)

/**************************************
Get: Get 指令，获取变量 Loc 的值，Get 实现了 Expr 接口
  - Loc 应为 Location
**************************************/

type Get struct {
	aStmt
	Loc Location
}

func (i *Get) Name() string   { return i.String() }
func (i *Get) Type() Type     { return i.Loc.DataType() }
func (i *Get) retained() bool { return false }
func (i *Get) String() string {
	if v, ok := i.Loc.(Var); ok && v.Kind() == AllocKindRegister {
		return fmt.Sprintf("get(%s)", i.Loc.Name())
	} else {
		return fmt.Sprintf("get*(%s)", i.Loc.Name())
	}
}

// 生成一条 Get 指令
func NewGet(loc Location, pos int) *Get {
	if loc == nil {
		panic("loc is nil")
	}

	v := &Get{Loc: loc}
	v.Stringer = v
	v.pos = pos
	return v
}

/**************************************
Load: Load 指令，获取内存 Loc 处的值， Load 实现了 Expr 接口
  - Loc 类型 应为 Ref、Ptr
  - 该指令仅供内部使用，上层高级语法不应直接使用
**************************************/

type Load struct {
	aStmt
	Loc Expr
}

func (i *Load) Name() string { return i.String() }
func (i *Load) Type() Type {
	switch t := i.Loc.Type().(type) {
	case *Ref:
		return t.Base

	case *Ptr:
		return t.Base

	default:
		panic(fmt.Sprintf("Invalid Loc.Type():%s", i.Loc.Type().Name()))
	}
}
func (i *Load) retained() bool { return false }
func (i *Load) String() string {
	return fmt.Sprintf("load(%s)", i.Loc.Name())
}

// 生成一条 Load 指令
func newLoad(loc Expr, pos int) *Load {
	if loc == nil {
		panic("loc is nil")
	}
	if loc.Type().Kind() != TypeKindPtr && loc.Type().Kind() != TypeKindRef {
		panic(fmt.Sprintf("loc.Type() is not Ptr or Ref: %s", loc.Type().Name()))
	}

	v := &Load{}
	v.Stringer = v
	v.Loc = loc
	v.pos = pos
	return v
}
