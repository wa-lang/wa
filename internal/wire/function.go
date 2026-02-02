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

	params  []*Var // 参数列表
	results []*Var // 返回值列表
	Body    *Block // 函数体，为 nil 表明该函数为外部导入

	scope Scope  // 匿名函数的父域为 Block，全局（非匿名）函数的父域为 Module
	types *Types // 该函数所属 Module 的类型库，切勿手动修改

	StartPos, EndPos int

	numImv int
}

// Scope 接口相关
func (f *Function) ScopeKind() ScopeKind { return ScopeKindFunc }
func (f *Function) ParentScope() Scope   { return f.scope }
func (f *Function) Lookup(obj interface{}, level VarKind) *Var {
	for _, p := range f.params {
		if p.object == obj {
			if level > p.kind {
				p.kind = level
			}
			return p
		}
	}
	for _, r := range f.results {
		if r.object == obj {
			if level > r.kind {
				r.kind = level
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
func (f *Function) AddParam(name string, typ Type, pos int, obj interface{}) *Var {
	v := &Var{}
	v.Stringer = v
	v.name = name
	v.dtype = typ
	v.rtype = f.types.GenRef(typ)

	v.pos = pos
	v.object = obj

	f.params = append(f.params, v)
	return v
}

// 添加返回值
func (f *Function) AddResult(name string, typ Type, pos int, obj interface{}) *Var {
	v := &Var{}
	v.Stringer = v
	v.name = name
	v.dtype = typ
	v.rtype = f.types.GenRef(typ)

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

	ob := f.Body
	ob.Label = _FN_START

	f.StartBody()

	var pre, post []Stmt

	// 逃逸参数置换：
	for i, param := range f.params {
		if param.kind != Register || param.DataType().hasRef() {
			np := *param
			np.Stringer = &np
			np.kind = Register

			param.name = "$" + param.name
			f.Body.emit(param)
			if param.DataType().hasRef() {
				f.Body.EmitStore(param, NewRetain(NewLoad(&np, np.pos), np.pos), np.pos)
			} else {
				f.Body.EmitStore(param, NewLoad(&np, np.pos), np.pos)
			}
			post = append(post, NewRelease(param, param.pos))

			f.params[i] = &np
		}
	}

	// 返回置换
	for _, ret := range f.results {
		f.Body.emit(ret)
	}
	f.retRepalce(ob)

	nb := &Block{}
	nb.scope = f.Body
	nb.types = f.types
	nb.Label = _FN_START
	nb.init()

	f.blockRcProc(ob, false, nb, &pre, &post)
	f.Body.Stmts = append(f.Body.Stmts, pre...)
	f.Body.emit(nb)
	f.Body.Stmts = append(f.Body.Stmts, post...)

	// Todo: defer

	// 插入真返回指令
	var ret_exprs []Expr
	for _, r := range f.results {
		ret_exprs = append(ret_exprs, NewLoad(r, f.EndPos))
	}
	f.Body.EmitReturn(ret_exprs, f.EndPos)

	{
		var sb strings.Builder
		sb.WriteString("\n<=======EndBody() 后=======>\n")
		f.Format("  ", &sb)
		println(sb.String())
	}
}

func exprReplaceLocation(e Expr, ov, nv *Var) Expr {
	if v, ok := e.(*Var); ok {
		if v == ov {
			return ov
		}
	}

	if s, ok := e.(Stmt); ok {
		stmtReplaceLocation(s, ov, nv)
	}

	return e
}

func stmtReplaceLocation(stmt Stmt, ov, nv *Var) {
	switch stmt := stmt.(type) {
	case *Block:
		for _, s := range stmt.Stmts {
			stmtReplaceLocation(s, ov, nv)
		}

	case *Var:
		return

	case *Load:
		if stmt.Loc == ov {
			stmt.Loc = nv
			return
		}

		if s, ok := stmt.Loc.(Stmt); ok {
			stmtReplaceLocation(s, ov, nv)
		}

	case *Store:
		for i := range stmt.Loc {
			if stmt.Loc[i] == ov {
				stmt.Loc[i] = nv
				continue
			}

			if s, ok := stmt.Loc[i].(Stmt); ok {
				stmtReplaceLocation(s, ov, nv)
			}
		}

		for i, val := range stmt.Val {
			stmt.Val[i] = exprReplaceLocation(val, ov, nv)
		}

	case *Br:
		return

	case *Return:
		for i, ret := range stmt.Results {
			stmt.Results[i] = exprReplaceLocation(ret, ov, nv)
		}

	case *Unop:
		stmt.X = exprReplaceLocation(stmt.X, ov, nv)

	case *Biop:
		stmt.X = exprReplaceLocation(stmt.X, ov, nv)
		stmt.Y = exprReplaceLocation(stmt.Y, ov, nv)

	case *If:
		stmt.Cond = exprReplaceLocation(stmt.Cond, ov, nv)
		stmtReplaceLocation(stmt.True, ov, nv)
		stmtReplaceLocation(stmt.False, ov, nv)

	case *Loop:
		stmt.Cond = exprReplaceLocation(stmt.Cond, ov, nv)
		stmtReplaceLocation(stmt.Body, ov, nv)
		stmtReplaceLocation(stmt.Post, ov, nv)

	case *Call:
		var call_common *CallCommon
		switch call := stmt.Callee.(type) {
		case *BuiltinCall:
			call_common = &call.CallCommon

		case *StaticCall:
			call_common = &call.CallCommon

		case *MethodCall:
			call_common = &call.CallCommon
			call.Recv = exprReplaceLocation(call.Recv, ov, nv)

		case *InterfaceCall:
			call_common = &call.CallCommon
			call.Interface = exprReplaceLocation(call.Interface, ov, nv)
		}

		for i, arg := range call_common.Args {
			call_common.Args[i] = exprReplaceLocation(arg, ov, nv)
		}

	default:
		panic(fmt.Sprintf("Todo: %s", stmt.String()))
	}

}

func (f *Function) retRepalce(stmt Stmt) {
	if stmt == nil {
		return
	}

	switch stmt := stmt.(type) {
	case *Block:
		if len(stmt.Stmts) == 0 {
			return
		}

		var i int
		for i = 0; i < len(stmt.Stmts)-1; i++ {
			f.retRepalce(stmt.Stmts[i])
		}

		if ret_stmt, ok := stmt.Stmts[i].(*Return); ok {
			br := &Br{}
			br.Stringer = br
			br.Label = _FN_START
			br.pos = ret_stmt.pos
			if len(f.results) > 0 && len(ret_stmt.Results) > 0 {
				store := &Store{}
				store.Stringer = store
				for _, loc := range f.results {
					store.Loc = append(store.Loc, loc)
				}
				store.Val = ret_stmt.Results
				store.pos = ret_stmt.pos

				stmt.Stmts[i] = store
				stmt.Stmts = append(stmt.Stmts, br)
			} else {
				stmt.Stmts[i] = br
			}
		}

	case *Var:
		return

	case *Load:
		if s, ok := stmt.Loc.(Stmt); ok {
			f.retRepalce(s)
		}

	case *Store:
		for _, loc := range stmt.Loc {
			if s, ok := loc.(Stmt); ok {
				f.retRepalce(s)
			}
		}

		for _, val := range stmt.Val {
			if s, ok := val.(Stmt); ok {
				f.retRepalce(s)
			}
		}

	case *Br:
		return

	case *Return:
		panic("Return can only be at the end of a block")

	case *Unop:
		if s, ok := stmt.X.(Stmt); ok {
			f.retRepalce(s)
		}

	case *Biop:
		if s, ok := stmt.X.(Stmt); ok {
			f.retRepalce(s)
		}
		if s, ok := stmt.Y.(Stmt); ok {
			f.retRepalce(s)
		}

	case *If:
		if s, ok := stmt.Cond.(Stmt); ok {
			f.retRepalce(s)
		}
		f.retRepalce(stmt.True)
		f.retRepalce(stmt.False)

	case *Loop:
		if s, ok := stmt.Cond.(Stmt); ok {
			f.retRepalce(s)
		}
		f.retRepalce(stmt.Body)
		f.retRepalce(stmt.Post)

	case *Call:
		var call_common *CallCommon
		switch call := stmt.Callee.(type) {
		case *BuiltinCall:
			call_common = &call.CallCommon

		case *StaticCall:
			call_common = &call.CallCommon

		case *MethodCall:
			call_common = &call.CallCommon
			if s, ok := call.Recv.(Stmt); ok {
				f.retRepalce(s)
			}

		case *InterfaceCall:
			call_common = &call.CallCommon
			if s, ok := call.Interface.(Stmt); ok {
				f.retRepalce(s)
			}
		}

		for _, arg := range call_common.Args {
			if s, ok := arg.(Stmt); ok {
				f.retRepalce(s)
			}
		}

	default:
		panic(fmt.Sprintf("Todo: %s", stmt.String()))
	}

}

func (f *Function) blockRcProc(s *Block, inloop bool, d *Block, pre *[]Stmt, post *[]Stmt) {
	for _, stmt := range s.Stmts {
		f.stmtRcProc(stmt, inloop, d, pre, post)
	}
}

func (f *Function) stmtRcProc(s Stmt, inloop bool, b *Block, pre *[]Stmt, post *[]Stmt) {
	switch s := s.(type) {
	case *Block:
		if len(s.Stmts) == 0 {
			return
		}

		block := b.EmitBlock(s.Label, s.pos)
		for _, stmt := range s.Stmts {
			f.stmtRcProc(stmt, inloop, block, pre, post)
		}

	case *Var:
		*pre = append(*pre, s)
		switch s.kind {
		case Register:
			if s.DataType().hasRef() {
				release := NewRelease(s, 0)
				*post = append(*post, release)
			}

		case Heap:
			release := NewRelease(s, 0)
			*post = append(*post, release)

		default:
			panic(fmt.Sprintf("Invalid VarKind: %v", s.kind))
		}

	case *Load:
		panic("Load should not be here")

	case *Store:
		for i := range s.Val {
			s.Val[i] = f.exprRcProc(s.Val[i], inloop, false, b, pre, post)
			if s.Val[i].Type().hasRef() && !s.Val[i].retained() {
				s.Val[i] = NewRetain(s.Val[i], s.Val[i].Pos())
			}
		}

		for i := range s.Loc {
			s.Loc[i] = f.exprRcProc(s.Loc[i], inloop, true, b, pre, post)
		}

		b.emit(s)

	case *Br:
		b.emit(s)

	case *Return:
		panic("Return should not be here")

	case *If:
		cond := f.exprRcProc(s.Cond, inloop, true, b, pre, post)
		n := b.EmitIf(cond, s.pos)
		f.blockRcProc(s.True, inloop, n.True, pre, post)
		f.blockRcProc(s.False, inloop, n.False, pre, post)

	case *Loop:
		cond := f.exprRcProc(s.Cond, true, true, b, pre, post)

		l := b.EmitLoop(cond, s.Label, s.pos)
		f.blockRcProc(s.Body, true, l.Body, pre, post)
		f.blockRcProc(s.Post, true, l.Post, pre, post)

	case *Retain:
		panic("Retain should not be here")

	case *Release:
		panic("Release should not be here")

	default:
		panic(fmt.Sprintf("Todo: %s", s.String()))
	}
}

func (f *Function) exprRcProc(e Expr, inloop bool, replace bool, b *Block, pre *[]Stmt, post *[]Stmt) (ret Expr) {
	if e == nil {
		return
	}
	if _, ok := e.(Stmt); !ok {
		return
	}

	ret = e
	switch e := e.(type) {
	case *Var:
		return

	case *Load:
		e.Loc = f.exprRcProc(e.Loc, inloop, true, b, pre, post)

	case *Unop:
		e.X = f.exprRcProc(e.X, inloop, true, b, pre, post)

	case *Biop:
		e.X = f.exprRcProc(e.X, inloop, true, b, pre, post)
		e.Y = f.exprRcProc(e.Y, inloop, true, b, pre, post)

	case *Call:
		var call_common *CallCommon
		switch call := e.Callee.(type) {
		case *BuiltinCall:
			call_common = &call.CallCommon

		case *StaticCall:
			call_common = &call.CallCommon

		case *MethodCall:
			call_common = &call.CallCommon
			call.Recv = f.exprRcProc(call.Recv, inloop, true, b, pre, post)

		case *InterfaceCall:
			call_common = &call.CallCommon
			call.Interface = f.exprRcProc(call.Interface, inloop, true, b, pre, post)
		}

		for i := range call_common.Args {
			call_common.Args[i] = f.exprRcProc(call_common.Args[i], inloop, true, b, pre, post)
		}

	case Stmt:
		panic(fmt.Sprintf("Todo: %s", e.String()))
	}

	if e.retained() && replace {
		imv := &Var{}
		imv.Stringer = imv
		imv.name = fmt.Sprintf("$$imv%d", f.numImv)
		f.numImv++
		imv.dtype = e.Type()
		*pre = append(*pre, imv)

		b.EmitStore(imv, e, e.Pos())
		ret = NewLoad(imv, e.Pos())

		release := NewRelease(imv, 0)
		*post = append(*post, release)
	}

	return
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
