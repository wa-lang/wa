// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
)

/**************************************
asLoc: 将 Expr 转换为 Location，asLoc 实现了 Location 接口
**************************************/

type asLoc struct {
	aStmt
	expr  Expr
	dtype Type
}

func (i *asLoc) Name() string   { return fmt.Sprintf("asloc(%s)", i.expr.Name()) }
func (i *asLoc) DataType() Type { return i.dtype }
func (i *asLoc) String() string { return fmt.Sprintf("asloc(%s)", i.expr.Name()) }

func AsLocation(expr Expr) Location {
	if l, ok := expr.(Location); ok {
		return l
	}

	l := &asLoc{expr: expr}
	l.Stringer = l
	l.pos = expr.Pos()

	switch t := expr.Type().(type) {
	case *Ref:
		l.dtype = t.Base
		return l
	case *Ptr:
		l.dtype = t.Base
		return l
	default:
		panic(fmt.Sprintf("Invalid loc.Type():%s", expr.Type().Name()))
	}
}

/**************************************
asAddr: 将 Location 转换为类型为指针或引用的 Expr，asAddr 实现了 Expr 接口
**************************************/

type asAddr struct {
	aStmt
	loc Location
	typ Type
}

func (i *asAddr) Name() string   { return i.String() }
func (i *asAddr) Type() Type     { return i.typ }
func (i *asAddr) retained() bool { return false }
func (i *asAddr) String() string { return fmt.Sprintf("asaddr(%s)", i.loc.Name()) }

func AsAddr(loc Location, typ Type, pos int) Expr {
	if e, ok := loc.(Expr); ok {
		return e
	}

	e := &asAddr{loc: loc, typ: typ}
	e.Stringer = e
	e.pos = pos

	return e
}
