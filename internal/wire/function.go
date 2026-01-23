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

			f.replaceLocation(body, param, loc)
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

func (f *Function) replaceLocation(stmt Stmt, ol, nl Location) {
	switch stmt := stmt.(type) {
	case *Block:
		for _, s := range stmt.Stmts {
			f.replaceLocation(s, ol, nl)
		}

	case *Alloc:
		return

	case *Load:
		if stmt.Loc == ol {
			stmt.Loc = nl
			return
		}

		if s, ok := stmt.Loc.(Stmt); ok {
			f.replaceLocation(s, ol, nl)
		}

	case *Store:
		for i := range stmt.Loc {
			if stmt.Loc[i] == ol {
				stmt.Loc[i] = nl
				continue
			}

			if s, ok := stmt.Loc[i].(Stmt); ok {
				f.replaceLocation(s, ol, nl)
			}
		}

		for i, val := range stmt.Val {
			if loc, ok := val.(Location); ok {
				if loc == ol {
					stmt.Val[i] = nl
					continue
				}
			}

			if s, ok := val.(Stmt); ok {
				f.replaceLocation(s, ol, nl)
			}
		}

	case *Br:
		return

	case *Return:
		for i, ret := range stmt.Results {
			if loc, ok := ret.(Location); ok {
				if loc == ol {
					stmt.Results[i] = nl
					continue
				}
			}

			if s, ok := ret.(Stmt); ok {
				f.replaceLocation(s, ol, nl)
			}
		}

	case *Unop:
		if loc, ok := stmt.X.(Location); ok {
			if loc == ol {
				stmt.X = nl
				return
			}
		}

		if s, ok := stmt.X.(Stmt); ok {
			f.replaceLocation(s, ol, nl)
		}

	case *Biop:
		if loc, ok := stmt.X.(Location); ok {
			if loc == ol {
				stmt.X = nl
			}
		}
		if s, ok := stmt.X.(Stmt); ok {
			f.replaceLocation(s, ol, nl)
		}

		if loc, ok := stmt.Y.(Location); ok {
			if loc == ol {
				stmt.Y = nl
			}
		}
		if s, ok := stmt.Y.(Stmt); ok {
			f.replaceLocation(s, ol, nl)
		}

	case *If:
		if loc, ok := stmt.Cond.(Location); ok {
			if loc == ol {
				stmt.Cond = nl
			}
		}
		if s, ok := stmt.Cond.(Stmt); ok {
			f.replaceLocation(s, ol, nl)
		}
		f.replaceLocation(stmt.True, ol, nl)
		f.replaceLocation(stmt.False, ol, nl)

	case *Loop:
		if loc, ok := stmt.Cond.(Location); ok {
			if loc == ol {
				stmt.Cond = nl
			}
		}
		if s, ok := stmt.Cond.(Stmt); ok {
			f.replaceLocation(s, ol, nl)
		}
		f.replaceLocation(stmt.Body, ol, nl)
		f.replaceLocation(stmt.Post, ol, nl)

	case *AsLocation:
		if loc, ok := stmt.addr.(Location); ok {
			if loc == ol {
				stmt.addr = nl
				return
			}
		}
		if s, ok := stmt.addr.(Stmt); ok {
			f.replaceLocation(s, ol, nl)
		}

	case *Call:
		var call_common *CallCommon
		switch call := stmt.Callee.(type) {
		case *BuiltinCall:
			call_common = &call.CallCommon

		case *StaticCall:
			call_common = &call.CallCommon

		case *MethodCall:
			call_common = &call.CallCommon
			if loc, ok := call.Recv.(Location); ok {
				if loc == ol {
					call.Recv = nl
				}
			}
			if s, ok := call.Recv.(Stmt); ok {
				f.replaceLocation(s, ol, nl)
			}

		case *InterfaceCall:
			call_common = &call.CallCommon
			if loc, ok := call.Interface.(Location); ok {
				if loc == ol {
					call.Interface = nl
				}
			}
			if s, ok := call.Interface.(Stmt); ok {
				f.replaceLocation(s, ol, nl)
			}
		}

		for i, arg := range call_common.Args {
			if loc, ok := arg.(Location); ok {
				if loc == ol {
					call_common.Args[i] = nl
					continue
				}
			}

			if s, ok := arg.(Stmt); ok {
				f.replaceLocation(s, ol, nl)
			}
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
