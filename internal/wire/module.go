// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
	"strings"
)

/**************************************
本文件包含了 Module 对象的功能
**************************************/

//-------------------------------------

/**************************************
Module: 定义了一个 wire 模块，对应于 ast.Program
Module 内含的所有成员极其子成员的 ValueType，必须经由本 Module 的 Types 创建
**************************************/

type Module struct {
	Types   Types
	Globals map[interface{}]*Var // 全局变量，对凹语言前端，键（key）为 types.Object
	Funcs   []*Function
}

// Scope 接口相关
func (m *Module) ScopeKind() ScopeKind { return ScopeKindModule }
func (m *Module) ParentScope() Scope   { return nil }
func (m *Module) Lookup(obj interface{}, level VarKind) *Var {
	v, ok := m.Globals[obj]
	if !ok {
		panic(fmt.Sprintf("no Value for: %v", obj))
	}
	return v
}
func (m *Module) Format(tab string, sb *strings.Builder) {
	// Todo: 输出全局变量

	for _, f := range m.Funcs {
		f.Format("", sb)
	}
}

// 初始化 Module
func (m *Module) Init() {
	m.Types.Init()
	m.Globals = make(map[interface{}]*Var)
}

// 创建一个 Function。该函数仅创建值，并不会将其合并至 Module 的相应位置
func (m *Module) NewFunction() *Function {
	f := Function{
		scope: m,
		types: &m.Types,
	}
	return &f
}

// 创建一个 常量值
func (m *Module) NewConst(name string, typ Type, pos int) Expr {
	c := Const{
		name: name,
		typ:  typ,
		pos:  pos,
	}
	return &c
}
