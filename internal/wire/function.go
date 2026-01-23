// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
	"strings"
)

/**************************************
本文件包含了 function 对象的功能
**************************************/

//-------------------------------------

/**************************************
Function: 函数。不可直接声明该类型的对象，必须通过 Module.NewFunction() 创建
**************************************/
type Function struct {
	InternalName string // 函数的内部名称(含包路径)，是其身份标识，应进行名字修饰
	ExternalName string // 函数的导出名称，非导出函数应为 nil

	params  []*Alloc // 参数列表
	results []*Alloc // 返回值列表
	Body    *Block   // 函数体，为 nil 表明该函数为外部导入

	scope Scope  // 匿名函数的父域为 Block，全局（非匿名）函数的父域为 Module
	types *Types // 该函数所属 Module 的类型库，切勿手动修改

}

// Scope 接口相关
func (f *Function) ScopeKind() ScopeKind { return ScopeKindFunc }
func (f *Function) ParentScope() Scope   { return f.scope }
func (f *Function) Lookup(obj interface{}, level LocationKind) Location {
	for _, p := range f.params {
		if p.Object() == obj {
			if level > p.location {
				p.location = level
			}
			return p
		}
	}
	for _, r := range f.results {
		if r.Object() == obj {
			if level > r.location {
				r.location = level
			}
			return r
		}
	}
	return f.scope.Lookup(obj, level)
}
func (f *Function) Format(tab string, sb *strings.Builder) {
	sb.WriteString(tab)

	sb.WriteString("func ")
	sb.WriteString(f.InternalName)

	sb.WriteRune('(')
	for i, v := range f.params {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.Name())
		sb.WriteRune(' ')
		sb.WriteString(v.DataType().Name())
	}
	sb.WriteRune(')')

	sb.WriteString(" => (")
	for i, v := range f.results {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.Name())
		sb.WriteRune(' ')
		sb.WriteString(v.DataType().Name())
	}
	sb.WriteRune(')')

	sb.WriteString("//ExternalName:")
	sb.WriteString(f.ExternalName)
	sb.WriteRune('\n')

	if f.Body != nil {
		f.Body.Format(tab, sb)
	}
}

// 添加函数参数
func (f *Function) AddParam(name string, typ Type, pos int, obj interface{}, refMode bool) Location {
	v := &Alloc{}
	v.Stringer = v
	v.name = name
	v.dataType = typ
	if refMode {
		v.refType = f.types.GenRef(typ)
	} else {
		v.refType = f.types.genPtr(typ)
	}
	v.pos = pos
	v.object = obj

	f.params = append(f.params, v)
	return v
}

// 添加返回值
func (f *Function) AddResult(name string, typ Type, pos int, obj interface{}, refMode bool) Location {
	v := &Alloc{}
	v.Stringer = v
	v.name = name
	v.dataType = typ
	if refMode {
		v.refType = f.types.GenRef(typ)
	} else {
		v.refType = f.types.genPtr(typ)
	}
	v.pos = pos
	v.object = obj

	f.results = append(f.results, v)
	return v
}

// 开始函数体
func (f *Function) StartBody() {
	f.Body = &Block{}
	f.Body.scope = f
	f.Body.types = f.types
	f.Body.init()
}

const (
	_FN_START = "$fn_start"
)

//
func (f *Function) EndBody() {
	if f.Body == nil {
		panic("StartBody first")
	}

	{
		var sb strings.Builder
		sb.WriteString("<=======EndBody() 前=======>\n")
		f.Format("  ", &sb)
		println(sb.String())
	}

	body := f.Body
	body.Label = _FN_START

	f.StartBody()

	// 逃逸参数置换：
	for _, param := range f.params {
		if param.LocationKind() != LocationKindRegister {
			_, refMode := param.refType.(*Ref)
			loc := f.Body.AddLocal("$"+param.name, param.dataType, param.pos, param.object, refMode)
			loc.location = param.location

			param.location = LocationKindRegister
			f.Body.EmitStore(loc, f.Body.NewLoad(param, param.pos), param.pos)

			stmtReplaceLocation(body, param, loc)
		}
	}

	f.Body.emit(body)

	{
		var sb strings.Builder
		sb.WriteString("\n<=======EndBody() 后=======>\n")
		f.Format("  ", &sb)
		println(sb.String())
	}

	//setImvId(0, f.Body)
}

func exprReplaceLocation(e Expr, ol, nl Location) Expr {
	if loc, ok := e.(Location); ok {
		if loc == ol {
			return nl
		}
	}

	if s, ok := e.(Stmt); ok {
		stmtReplaceLocation(s, ol, nl)
	}

	return e
}

func stmtReplaceLocation(stmt Stmt, ol, nl Location) {
	switch stmt := stmt.(type) {
	case *Block:
		for _, s := range stmt.Stmts {
			stmtReplaceLocation(s, ol, nl)
		}

	case *Alloc:
		return

	case *Load:
		if stmt.Loc == ol {
			stmt.Loc = nl
			return
		}

		if s, ok := stmt.Loc.(Stmt); ok {
			stmtReplaceLocation(s, ol, nl)
		}

	case *Store:
		for i := range stmt.Loc {
			if stmt.Loc[i] == ol {
				stmt.Loc[i] = nl
				continue
			}

			if s, ok := stmt.Loc[i].(Stmt); ok {
				stmtReplaceLocation(s, ol, nl)
			}
		}

		for i, val := range stmt.Val {
			stmt.Val[i] = exprReplaceLocation(val, ol, nl)
		}

	case *Br:
		return

	case *Return:
		for i, ret := range stmt.Results {
			stmt.Results[i] = exprReplaceLocation(ret, ol, nl)
		}

	case *Unop:
		stmt.X = exprReplaceLocation(stmt.X, ol, nl)

	case *Biop:
		stmt.X = exprReplaceLocation(stmt.X, ol, nl)
		stmt.Y = exprReplaceLocation(stmt.Y, ol, nl)

	case *If:
		stmt.Cond = exprReplaceLocation(stmt.Cond, ol, nl)
		stmtReplaceLocation(stmt.True, ol, nl)
		stmtReplaceLocation(stmt.False, ol, nl)

	case *Loop:
		stmt.Cond = exprReplaceLocation(stmt.Cond, ol, nl)
		stmtReplaceLocation(stmt.Body, ol, nl)
		stmtReplaceLocation(stmt.Post, ol, nl)

	case *AsLocation:
		stmt.addr = exprReplaceLocation(stmt.addr, ol, nl)

	case *Call:
		var call_common *CallCommon
		switch call := stmt.Callee.(type) {
		case *BuiltinCall:
			call_common = &call.CallCommon

		case *StaticCall:
			call_common = &call.CallCommon

		case *MethodCall:
			call_common = &call.CallCommon
			call.Recv = exprReplaceLocation(call.Recv, ol, nl)

		case *InterfaceCall:
			call_common = &call.CallCommon
			call.Interface = exprReplaceLocation(call.Interface, ol, nl)
		}

		for i, arg := range call_common.Args {
			call_common.Args[i] = exprReplaceLocation(arg, ol, nl)
		}

	default:
		panic(fmt.Sprintf("Todo: %s", stmt.String()))
	}

}

//func setImvId(num int, b *Block) int {
//	for _, i := range b.Instrs {
//		if v, ok := i.(imv); ok {
//			if av, ok := v.(Value); ok {
//				avt := av.Type()
//				if avt != nil && !avt.Equal(&Void{}) {
//					v.setId(num)
//					num++
//				}
//			}
//		}
//
//		switch i := i.(type) {
//		case *Block:
//			num = setImvId(num, i)
//
//		case *InstIf:
//			num = setImvId(num, i.True)
//			num = setImvId(num, i.False)
//		}
//	}
//	return num
//}
