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
	fb := f.Body

	// 参数置换：
	for i, param := range f.params {
		if param.kind != Register || param.DataType().hasRef() && ob.varUsageRange(param).first != -1 { //Todo: 待优化，若参数中携带的引用未被重新赋值则无需置换  ob.varStored(param)
			np := *param
			np.Stringer = &np
			np.kind = Register

			param.name = "$" + param.name
			fb.emit(param)
			if param.DataType().hasRef() {
				fb.EmitSet(param, NewRetain(NewGet(&np, np.pos), np.pos), np.pos)
			} else {
				fb.EmitSet(param, NewGet(&np, np.pos), np.pos)
			}

			f.params[i] = &np
		}
	}

	// 返回置换
	rets := make(map[*Var]bool)
	for _, ret := range f.results {
		fb.emit(ret)
		rets[ret] = true
	}
	f.retRepalce(ob)

	nb := &Block{}
	nb.scope = fb
	nb.types = f.types
	nb.Label = _FN_START
	nb.init()

	blockImvRcProc(ob, false, nb)
	fb.emit(nb)

	f.varRangeProc(fb, rets)

	// Todo: defer

	// 插入真返回指令
	var ret_exprs []Expr
	for _, r := range f.results {
		switch r.kind {
		case Register:
			ret_exprs = append(ret_exprs, NewGet(r, f.EndPos))

		case Heap:
			rr := *r
			rr.Stringer = &rr
			rr.kind = Register
			rr.name = fmt.Sprintf("$%d", fb.imvCount)
			fb.imvCount++
			fb.emit(&rr)
			fb.EmitSet(&rr, NewGet(r, f.EndPos), f.EndPos)
			fb.emit(NewDrop(r, f.EndPos))
			ret_exprs = append(ret_exprs, NewGet(&rr, f.EndPos))

		default:
			panic(fmt.Sprintf("Todo: VarKind: %v", r.kind))
		}
	}
	fb.EmitReturn(ret_exprs, f.EndPos)

	{
		var sb strings.Builder
		sb.WriteString("\n<=======EndBody() 后=======>\n")
		f.Format("  ", &sb)
		println(sb.String())
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
				set := &Set{}
				set.Stringer = set
				for _, loc := range f.results {
					set.Loc = append(set.Loc, loc)
				}
				set.Val = ret_stmt.Results
				set.pos = ret_stmt.pos

				stmt.Stmts[i] = set
				stmt.Stmts = append(stmt.Stmts, br)
			} else {
				stmt.Stmts[i] = br
			}
		}

	case *If:
		f.retRepalce(stmt.True)
		f.retRepalce(stmt.False)

	case *Loop:
		f.retRepalce(stmt.Body)
		f.retRepalce(stmt.Post)

	case *Return:
		panic("Return can only be at the end of a block")

	case *Var, *Get, *Set, *Br, *Unop, *Biop, *Call:
		return

	default:
		panic(fmt.Sprintf("Todo: %s", stmt.String()))
	}

}

func (f *Function) varRangeProc(b *Block, reserve map[*Var]bool) {
	for _, stmt := range b.Stmts {
		switch s := stmt.(type) {
		case *Block:
			f.varRangeProc(s, reserve)

		case *Var:
			if _, ok := reserve[s]; ok {
				continue
			}

			r := b.varUsageRange(s)
			if r.last == -1 {
				panic(fmt.Sprintf("var:%s not used", s.Name()))
			}

			var t []Stmt
			t = append(t, b.Stmts[:r.last+1]...)
			t = append(t, NewDrop(s, b.Stmts[r.last].Pos()))
			t = append(t, b.Stmts[r.last+1:]...)
			b.Stmts = t

		case *If:
			f.varRangeProc(s.True, reserve)
			f.varRangeProc(s.False, reserve)

		case *Loop:
			f.varRangeProc(s.Body, reserve)
			f.varRangeProc(s.Post, reserve)

		case *Get, *Set, *Br, *Return, *Unop, *Biop, *Retain, *Drop:
			continue

		default:
			panic(fmt.Sprintf("Todo: %s", stmt.String()))
		}

	}
}

func blockImvRcProc(b *Block, inloop bool, d *Block) {
	for _, stmt := range b.Stmts {
		stmtImvRcProc(stmt, inloop, d)
	}
}

func stmtImvRcProc(s Stmt, inloop bool, d *Block) {
	var post []Stmt
	switch s := s.(type) {
	case *Block:
		if len(s.Stmts) == 0 {
			return
		}

		block := d.EmitBlock(s.Label, s.pos)
		for _, stmt := range s.Stmts {
			stmtImvRcProc(stmt, inloop, block)
		}

	case *Var:
		d.emit(s)

	case *Get:
		panic("Get should not be here")

	case *Set:
		for i := range s.Loc {
			s.Loc[i] = exprImvRcProc(s.Loc[i], inloop, true, d, &post)
		}

		for i := range s.Val {
			s.Val[i] = exprImvRcProc(s.Val[i], inloop, false, d, &post)
			if s.Val[i].Type().hasRef() && !s.Val[i].retained() {
				s.Val[i] = NewRetain(s.Val[i], s.Val[i].Pos())
			}
		}

		d.emit(s)

	case *Br:
		d.emit(s)

	case *Return:
		panic("Return should not be here")

	case *If:
		cond := exprImvRcProc(s.Cond, inloop, true, d, &post)
		n := d.EmitIf(cond, s.pos)
		blockImvRcProc(s.True, inloop, n.True)
		blockImvRcProc(s.False, inloop, n.False)

	case *Loop:
		cond := exprImvRcProc(s.Cond, true, true, d, &post)

		l := d.EmitLoop(cond, s.Label, s.pos)
		blockImvRcProc(s.Body, true, l.Body)
		blockImvRcProc(s.Post, true, l.Post)

	case *Retain:
		d.emit(s)

	case *Drop:
		d.emit(s)

	default:
		panic(fmt.Sprintf("Todo: %s", s.String()))
	}

	for _, i := range post {
		d.emit(i)
	}
}

func exprImvRcProc(e Expr, inloop bool, replace bool, d *Block, post *[]Stmt) (ret Expr) {
	ret = e
	if e == nil {
		return
	}
	if _, ok := e.(Stmt); !ok {
		return
	}

	switch e := e.(type) {
	case *Var:
		return

	case *Get:
		e.Loc = exprImvRcProc(e.Loc, inloop, true, d, post)

	case *Unop:
		e.X = exprImvRcProc(e.X, inloop, true, d, post)

	case *Biop:
		e.X = exprImvRcProc(e.X, inloop, true, d, post)
		e.Y = exprImvRcProc(e.Y, inloop, true, d, post)

	case *Call:
		var call_common *CallCommon
		switch call := e.Callee.(type) {
		case *BuiltinCall:
			call_common = &call.CallCommon

		case *StaticCall:
			call_common = &call.CallCommon

		case *MethodCall:
			call_common = &call.CallCommon
			call.Recv = exprImvRcProc(call.Recv, inloop, true, d, post)

		case *InterfaceCall:
			call_common = &call.CallCommon
			call.Interface = exprImvRcProc(call.Interface, inloop, true, d, post)

		default:
			panic("Todo")
		}

		for i := range call_common.Args {
			call_common.Args[i] = exprImvRcProc(call_common.Args[i], inloop, true, d, post)
		}

	case Stmt:
		panic(fmt.Sprintf("Todo: %s", e.String()))
	}

	if e.retained() && replace {
		dupref := NewDupRef(e, fmt.Sprintf("$%d", d.imvCount), e.Pos())
		d.imvCount++
		ret = dupref

		drop := NewDrop(dupref.Imv, e.Pos())
		*post = append(*post, drop)

		//imv := &Var{}
		//imv.Stringer = imv
		//imv.name = fmt.Sprintf("$%d", d.imvCount)
		//d.imvCount++
		//imv.dtype = e.Type()
		//
		//d.emit(imv)
		//d.EmitSet(imv, e, e.Pos())
		//ret = NewGet(imv, e.Pos())
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

//func exprReplaceLocation(e Expr, ov, nv *Var) Expr {
//	if v, ok := e.(*Var); ok {
//		if v == ov {
//			return ov
//		}
//	}
//
//	if s, ok := e.(Stmt); ok {
//		stmtReplaceLocation(s, ov, nv)
//	}
//
//	return e
//}
//
//func stmtReplaceLocation(stmt Stmt, ov, nv *Var) {
//	switch stmt := stmt.(type) {
//	case *Block:
//		for _, s := range stmt.Stmts {
//			stmtReplaceLocation(s, ov, nv)
//		}
//
//	case *Var:
//		return
//
//	case *Get:
//		if stmt.Loc == ov {
//			stmt.Loc = nv
//			return
//		}
//
//		if s, ok := stmt.Loc.(Stmt); ok {
//			stmtReplaceLocation(s, ov, nv)
//		}
//
//	case *Set:
//		for i := range stmt.Loc {
//			if stmt.Loc[i] == ov {
//				stmt.Loc[i] = nv
//				continue
//			}
//
//			if s, ok := stmt.Loc[i].(Stmt); ok {
//				stmtReplaceLocation(s, ov, nv)
//			}
//		}
//
//		for i, val := range stmt.Val {
//			stmt.Val[i] = exprReplaceLocation(val, ov, nv)
//		}
//
//	case *Br:
//		return
//
//	case *Return:
//		for i, ret := range stmt.Results {
//			stmt.Results[i] = exprReplaceLocation(ret, ov, nv)
//		}
//
//	case *Unop:
//		stmt.X = exprReplaceLocation(stmt.X, ov, nv)
//
//	case *Biop:
//		stmt.X = exprReplaceLocation(stmt.X, ov, nv)
//		stmt.Y = exprReplaceLocation(stmt.Y, ov, nv)
//
//	case *If:
//		stmt.Cond = exprReplaceLocation(stmt.Cond, ov, nv)
//		stmtReplaceLocation(stmt.True, ov, nv)
//		stmtReplaceLocation(stmt.False, ov, nv)
//
//	case *Loop:
//		stmt.Cond = exprReplaceLocation(stmt.Cond, ov, nv)
//		stmtReplaceLocation(stmt.Body, ov, nv)
//		stmtReplaceLocation(stmt.Post, ov, nv)
//
//	case *Call:
//		var call_common *CallCommon
//		switch call := stmt.Callee.(type) {
//		case *BuiltinCall:
//			call_common = &call.CallCommon
//
//		case *StaticCall:
//			call_common = &call.CallCommon
//
//		case *MethodCall:
//			call_common = &call.CallCommon
//			call.Recv = exprReplaceLocation(call.Recv, ov, nv)
//
//		case *InterfaceCall:
//			call_common = &call.CallCommon
//			call.Interface = exprReplaceLocation(call.Interface, ov, nv)
//
//		default:
//			panic("Todo")
//		}
//
//		for i, arg := range call_common.Args {
//			call_common.Args[i] = exprReplaceLocation(arg, ov, nv)
//		}
//
//	default:
//		panic(fmt.Sprintf("Todo: %s", stmt.String()))
//	}
//
//}
