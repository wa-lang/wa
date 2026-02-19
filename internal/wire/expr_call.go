// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
)

/**************************************
Call: 函数调用指令，Callee 为 StaticCall、BuiltinCall、MethodCall、InterfaceCall、ClosureCall 之一，Call 实现了 Expr
**************************************/

type Call struct {
	aStmt
	Callee Callee
}

func (i *Call) Name() string   { return i.String() }
func (i *Call) Type() Type     { return i.Callee.Type() }
func (i *Call) Pos() int       { return i.Callee.Pos() }
func (i *Call) retained() bool { return rtimp.hasChunk(i.Callee.Type()) }
func (i *Call) String() string { return i.Callee.(fmt.Stringer).String() }

// 在 Block 中添加一条 InstCall 指令，callee 为 StaticCall、BuiltinCall、MethodCall、InterfaceCall、ClosureCall 之一
func (b *Block) NewCall(callee Callee) *Call {
	v := &Call{Callee: callee}
	v.Stringer = v
	v.pos = callee.Pos()
	return v
}

/**************************************
Callee: 调用接口
**************************************/

type Callee interface {
	Name() string
	Type() Type
	Pos() int
}

/**************************************
StaticCall: 包函数、非闭包匿名函数调用。满足 Calle 接口。
**************************************/

type StaticCall struct {
	CallCommon
}

func (p *StaticCall) Pos() int { return p.CallCommon.Pos }

/**************************************
BuiltinCall: 内置函数调用。内置函数调用为特殊的静态调用，满足 Expr 接口
**************************************/

type BuiltinCall struct {
	CallCommon
}

func (p *BuiltinCall) Pos() int { return p.CallCommon.Pos }

/**************************************
MethodCall: 对象方法调用，满足 Expr 接口
**************************************/

type MethodCall struct {
	Recv Expr // recv/接收器
	CallCommon
}

func (p *MethodCall) Pos() int       { return p.CallCommon.Pos }
func (p *MethodCall) String() string { return p.Recv.Name() + "." + p.CallCommon.String() }

/**************************************
InterfaceCall: 接口方法调用（既 Invoke），满足 Expr 接口
**************************************/

type InterfaceCall struct {
	Interface Expr // 被调用的接口
	CallCommon
}

func (p *InterfaceCall) Pos() int       { return p.CallCommon.Pos }
func (p *InterfaceCall) String() string { return p.Interface.Name() + "." + p.CallCommon.String() }

/**************************************
ClosureCall: 闭包调用，满足 Expr 接口
**************************************/

type ClosureCall struct {
	Closure Expr
	CallCommon
}

func (p *ClosureCall) Pos() int       { return p.CallCommon.Pos }
func (p *ClosureCall) String() string { return p.Closure.Name() + "." + p.CallCommon.String() }

/**************************************
CallCommon: 调用基本信息（函数名、函数签名、参数、调用位置）
**************************************/

type CallCommon struct {
	FnName string
	Sig    FnSig
	Args   []Expr
	Pos    int
}

func (v *CallCommon) Name() string { return v.String() }
func (v *CallCommon) Type() Type   { return v.Sig.Results }
func (v *CallCommon) String() string {
	s := v.FnName
	s += "("
	for i, p := range v.Args {
		if i > 0 {
			s += ", "
		}
		s += p.Name()
	}
	s += ")"
	return s
}
