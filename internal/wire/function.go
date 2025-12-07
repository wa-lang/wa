// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import "strings"

/**************************************
本文件包含了 function 对象的功能
**************************************/

//-------------------------------------

/**************************************
Function: 函数。不可直接声明该类型的对象，必须通过 Module.NewFunction() 创建
**************************************/
type Function struct {
	InternalName string  // 函数的内部名称(含包路径)，是其身份标识，应进行名字修饰
	ExternalName string  // 函数的导出名称，非导出函数应为 nil
	Params       []Value // 参数列表
	Results      []Type  // 返回值列表。非具名返回值的名字为空
	Body         *Block  // 函数体，为 nil 表明该函数为外部导入

	scope Scope  // 匿名函数的父域为 Block，全局（非匿名）函数的父域为 Module
	types *Types // 该函数所属 Module 的类型库，切勿手动修改

}

// Scope 接口相关
func (f *Function) ScopeKind() ScopeKind { return ScopeKindFunc }
func (f *Function) ParentScope() Scope   { return f.scope }
func (f *Function) Lookup(obj interface{}, escaping bool) Location {
	return f.scope.Lookup(obj, escaping)
}
func (f *Function) Format(tab string, sb *strings.Builder) {
	sb.WriteString(tab)

	sb.WriteString("func ")
	sb.WriteString(f.InternalName)

	sb.WriteRune('(')
	for i, v := range f.Params {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.Name())
		sb.WriteRune(' ')
		sb.WriteString(v.Type().Name())
	}
	sb.WriteRune(')')

	sb.WriteString(" => (")
	for i, v := range f.Results {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.Name())
	}
	sb.WriteRune(')')

	sb.WriteString("//ExternalName:")
	sb.WriteString(f.ExternalName)
	sb.WriteRune('\n')

	if f.Body != nil {
		f.Body.Format(tab, sb)
	}
}

// 开始函数体
func (f *Function) StartBody() {
	f.Body = &Block{}
	f.Body.scope = f
	f.Body.types = f.types
	f.Body.Init()
}

//
func (f *Function) EndBody() {
	if f.Body == nil {
		panic("StartBody first")
	}

	setImvId(0, f.Body)
}

func setImvId(num int, b *Block) int {
	for _, i := range b.Instrs {
		if block, ok := i.(*Block); ok {
			num = setImvId(num, block)
		} else if v, ok := i.(imv); ok {
			v.setId(num)
			num++
		}
	}
	return num
}
