// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
	"strings"
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

/**************************************
NilCheckWrapper: 检查 X 是否为 nil，为 nil 则 panic，不为 nil 则返回 X。NilCheck 实现了 Expr 接口
  - X 类型应为 Ptr 或 Ref
**************************************/
type NilCheckWrapper struct {
	aStmt
	X Expr
}

func (i *NilCheckWrapper) Name() string   { return i.String() }
func (i *NilCheckWrapper) Type() Type     { return i.X.Type() }
func (i *NilCheckWrapper) retained() bool { return false }
func (i *NilCheckWrapper) String() string { return fmt.Sprintf("nilcheckwrapper(%s)", i.X.Name()) }

func NewNilCheckWrapper(x Expr) Expr {
	if x.Type().Kind() != TypeKindPtr && x.Type().Kind() != TypeKindRef {
		panic(fmt.Sprintf("Invalid X.Type():%s", x.Type().Name()))
	}

	e := &NilCheckWrapper{X: x}
	e.Stringer = e
	e.pos = x.Pos()
	return e
}

/**************************************
NilCheck: 检查 X 是否为 nil
  - X 类型应为 Ptr 或 Ref 类型的 Var
**************************************/
type NilCheck struct {
	aStmt
	X Var
}

func (i *NilCheck) String() string { return fmt.Sprintf("nilcheck(%s)", i.X.Name()) }

func newNilCheck(x Var) *NilCheck {
	if x.Type().Kind() != TypeKindPtr && x.Type().Kind() != TypeKindRef {
		panic(fmt.Sprintf("Invalid X.Type():%s", x.Type().Name()))
	}
	s := &NilCheck{X: x}
	s.Stringer = s
	s.pos = x.Pos()
	return s
}

/**************************************
Combo: 组合指令，将多个指令组合成一个指令，实现了 Expr 、Var 接口，返回 Result
  - Result 应为 Var
**************************************/
type Combo struct {
	aStmt
	Stmts  []Stmt
	Result Var
}

func (i *Combo) Name() string    { return i.String() }
func (i *Combo) Type() Type      { return i.Result.Type() }
func (i *Combo) retained() bool  { return i.Result.retained() }
func (i *Combo) DataType() Type  { return i.Result.DataType() }
func (i *Combo) Kind() AllocKind { return i.Result.Kind() }
func (i *Combo) Tank() *tank     { return i.Result.Tank() }

func (i *Combo) String() string {
	var sb strings.Builder
	sb.WriteRune('{')
	for _, stmt := range i.Stmts {
		sb.WriteString(stmt.String())
		sb.WriteString("; ")
	}
	sb.WriteString(i.Result.Name())
	sb.WriteRune('}')
	return sb.String()
}

func NewCombo(stmts []Stmt, result Var, pos int) *Combo {
	v := &Combo{Stmts: stmts, Result: result}
	v.Stringer = v
	v.pos = pos
	return v
}
