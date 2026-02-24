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

	StartPos, EndPos int

	freeCommonRegs []register
	freeChunkRegs  []register
	chunkRegs      []register
	regCount       int
}

// Scope 接口相关
func (f *Function) ScopeKind() ScopeKind { return ScopeKindFunc }
func (f *Function) ParentScope() Scope   { return f.scope }
func (f *Function) Lookup(obj interface{}, level VarKind) *Alloc {
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
func (f *Function) AddParam(name string, typ Type, pos int, obj interface{}) *Alloc {
	v := f.Body.NewAlloc(name, typ, pos, obj, nil)
	//v := &Var{}
	//v.Stringer = v
	//v.name = name
	//v.dtype = typ
	//v.rtype = f.types.GenRef(typ)
	//
	//v.pos = pos
	//v.object = obj

	f.params = append(f.params, v)
	return v
}

// 添加返回值
func (f *Function) AddResult(name string, typ Type, pos int, obj interface{}) *Alloc {
	v := f.Body.NewAlloc(name, typ, pos, obj, nil)
	//v := &Var{}
	//v.Stringer = v
	//v.name = name
	//v.dtype = typ
	//v.rtype = f.types.GenRef(typ)
	//
	//v.pos = pos
	//v.object = obj

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
		hasChunk := rtimp.hasChunk(param.DataType())
		if param.kind != Register || hasChunk && ob.varUsageRange(param).first != -1 { //Todo: 待优化，若参数中携带的引用未被重新赋值则无需置换  ob.varStored(param)
			np := *param
			np.Stringer = &np
			np.kind = Register

			param.name = "$" + param.name
			fb.emit(param)
			if hasChunk {
				fb.EmitSet(param, NewRetain(NewGet(&np, np.pos), np.pos), np.pos)
			} else {
				fb.EmitSet(param, NewGet(&np, np.pos), np.pos)
			}

			f.params[i] = &np
		}
	}

	// 返回置换
	for _, ret := range f.results {
		fb.emit(ret)
	}
	f.retRepalce(ob)

	nb := &Block{}
	nb.scope = fb
	nb.types = f.types
	nb.Label = _FN_START
	nb.init()

	rc_block(ob, false, nb)
	f.autoDrop(nb)

	fb.emit(nb)

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
			ret_exprs = append(ret_exprs, NewGet(&rr, f.EndPos))

		default:
			panic(fmt.Sprintf("Todo: VarKind: %v", r.kind))
		}
	}
	// Todo: 插入所有 chunk 的 release
	fb.EmitReturn(ret_exprs, f.EndPos)

	// 虚拟寄存器
	f.allocVR_block(f.Body, false)

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
			br := NewBr(_FN_START, ret_stmt.pos)
			if len(f.results) > 0 && len(ret_stmt.Results) > 0 {
				lhs := make([]Expr, len(f.results))
				for i, lh := range f.results {
					lhs[i] = lh
				}

				stmt.Stmts[i] = NewSetN(lhs, ret_stmt.Results, ret_stmt.pos)
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

	case *Alloc, *Get, *Set, *Br, *Unop, *Biop, *Call:
		return

	default:
		panic(fmt.Sprintf("Todo: %s", stmt.String()))
	}

}

func rc_block(b *Block, inloop bool, d *Block) {
	for _, stmt := range b.Stmts {
		rc_stmt(stmt, inloop, d)
	}
}

func rc_stmt(s Stmt, inloop bool, d *Block) {
	var pre []Stmt
	switch s := s.(type) {
	case *Block:
		if len(s.Stmts) == 0 {
			return
		}

		block := d.EmitBlock(s.Label, s.pos)
		for _, stmt := range s.Stmts {
			rc_stmt(stmt, inloop, block)
		}

	case *Alloc:
		s.init = rc_expr(s.init, inloop, false, false, d, &pre)
		d.Stmts = append(d.Stmts, pre...)
		if s.init != nil && rtimp.hasChunk(s.init.Type()) && !s.init.retained() {
			s.init = NewRetain(s.init, s.init.Pos())
		}
		d.emit(s)

	case *Get:
		panic("Get should not be here")

	case *Set:
		for i := range s.Lhs {
			s.Lhs[i] = rc_expr(s.Lhs[i], inloop, true, false, d, &pre)
		}

		for i := range s.Rhs {
			s.Rhs[i] = rc_expr(s.Rhs[i], inloop, false, false, d, &pre)
			if rtimp.hasChunk(s.Rhs[i].Type()) && !s.Rhs[i].retained() {
				s.Rhs[i] = NewRetain(s.Rhs[i], s.Rhs[i].Pos())
			}
		}

		d.Stmts = append(d.Stmts, pre...)
		d.emit(s)

	case *Br:
		d.emit(s)

	case *Return:
		panic("Return should not be here")

	case *If:
		cond := rc_expr(s.Cond, inloop, true, false, d, &pre)
		d.Stmts = append(d.Stmts, pre...)
		n := d.EmitIf(cond, s.pos)
		rc_block(s.True, inloop, n.True)
		rc_block(s.False, inloop, n.False)

	case *Loop:
		cond := rc_expr(s.Cond, true, true, true, d, &pre)

		l := d.EmitLoop(cond, s.Label, s.pos)
		l.PreCond = append(s.PreCond, pre...)

		rc_block(s.Body, true, l.Body)
		rc_block(s.Post, true, l.Post)

	case *Retain:
		d.emit(s)

	case *Drop:
		d.emit(s)

	default:
		panic(fmt.Sprintf("Todo: %s", s.String()))
	}
}

func rc_expr(e Expr, inloop bool, replace bool, isLoopCond bool, d *Block, pre *[]Stmt) (ret Expr) {
	ret = e
	if e == nil {
		return
	}
	if _, ok := e.(Stmt); !ok {
		return
	}

	switch e := e.(type) {
	case *Alloc:
		return

	case *Get:
		e.Loc = rc_expr(e.Loc, inloop, true, isLoopCond, d, pre)

	case *Unop:
		e.X = rc_expr(e.X, inloop, true, isLoopCond, d, pre)

	case *Biop:
		e.X = rc_expr(e.X, inloop, true, isLoopCond, d, pre)
		e.Y = rc_expr(e.Y, inloop, true, isLoopCond, d, pre)

	case *Call:
		var call_common *CallCommon
		switch call := e.Callee.(type) {
		case *BuiltinCall:
			call_common = &call.CallCommon

		case *StaticCall:
			call_common = &call.CallCommon

		case *MethodCall:
			call_common = &call.CallCommon
			call.Recv = rc_expr(call.Recv, inloop, true, isLoopCond, d, pre)

		case *InterfaceCall:
			call_common = &call.CallCommon
			call.Interface = rc_expr(call.Interface, inloop, true, isLoopCond, d, pre)

		default:
			panic("Todo")
		}

		for i := range call_common.Args {
			call_common.Args[i] = rc_expr(call_common.Args[i], inloop, true, isLoopCond, d, pre)
		}

	case Stmt:
		panic(fmt.Sprintf("Todo: %s", e.String()))
	}

	if e.retained() && replace {
		if isLoopCond {
			imv := d.NewAlloc(fmt.Sprintf("$%d", d.imvCount), e.Type(), e.Pos(), nil, nil)
			imv.noinit = true
			d.imvCount++
			d.emit(imv)

			*pre = append(*pre, NewSet(imv, e, e.Pos()))
			ret = NewGet(imv, e.Pos())
		} else {
			imv := d.NewAlloc(fmt.Sprintf("$%d", d.imvCount), e.Type(), e.Pos(), nil, e)
			d.imvCount++
			d.emit(imv)

			ret = NewGet(imv, e.Pos())
		}
	}

	return
}

func (f *Function) autoDrop(b *Block) {
	for _, stmt := range b.Stmts {
		switch s := stmt.(type) {
		case *Block:
			f.autoDrop(s)

		case *Alloc:
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
			f.autoDrop(s.True)
			f.autoDrop(s.False)

		case *Loop:
			f.autoDrop(s.Body)
			f.autoDrop(s.Post)

		case *Get, *Set, *Br, *Return, *Unop, *Biop, *Retain, *Drop:
			continue

		default:
			panic(fmt.Sprintf("Todo: %s", stmt.String()))
		}

	}
}

func (f *Function) allocVR_block(b *Block, inloop bool) {
	for _, s := range b.Stmts {
		f.allocVR_stmt(s, inloop)
	}
}

func (f *Function) allocVR_stmt(s Stmt, inloop bool) {
	switch s := s.(type) {
	case *Block:
		f.allocVR_block(s, inloop)

	case *Alloc:
		f.allocVR_expr(s.init, inloop)
		f.allocVR_var(s, inloop)

	case *Get:
		panic("Get should not be here")

	case *Set:
		for _, lh := range s.Lhs {
			f.allocVR_expr(lh, inloop)
		}
		for _, rh := range s.Rhs {
			f.allocVR_expr(rh, inloop)
		}

	case *If:
		f.allocVR_expr(s.Cond, inloop)
		f.allocVR_block(s.True, inloop)
		f.allocVR_block(s.False, inloop)

	case *Loop:
		f.allocVR_expr(s.Cond, true)
		f.allocVR_block(s.Body, true)
		f.allocVR_block(s.Post, true)

	case *Drop:
		f.dropVR_tank(s.X.tank)

	case *Retain:
		panic("Retain should not be here")

	case *Br, *Return:

	default:
		panic(fmt.Sprintf("Todo: %s", s.String()))
	}
}

func (f *Function) allocVR_expr(e Expr, inloop bool) {
	if e == nil {
		return
	}
	if _, ok := e.(Stmt); !ok {
		return
	}

	switch e := e.(type) {
	case *Alloc:
		return

	case *Get:
		f.allocVR_expr(e.Loc, inloop)

	case *Unop:
		f.allocVR_expr(e.X, inloop)

	case *Biop:
		f.allocVR_expr(e.X, inloop)
		f.allocVR_expr(e.Y, inloop)

	//case *Combo:
	//	for _, stmt := range e.Stmts {
	//		f.allocVR_stmt(stmt, inloop)
	//	}
	//	f.allocVR_expr(e.Result, inloop)

	case *Call:
		var call_common *CallCommon
		switch call := e.Callee.(type) {
		case *BuiltinCall:
			call_common = &call.CallCommon

		case *StaticCall:
			call_common = &call.CallCommon

		case *MethodCall:
			call_common = &call.CallCommon
			f.allocVR_expr(call.Recv, inloop)

		case *InterfaceCall:
			call_common = &call.CallCommon
			f.allocVR_expr(call.Interface, inloop)

		default:
			panic("Todo")
		}

		for _, arg := range call_common.Args {
			f.allocVR_expr(arg, inloop)
		}

	//case *DupRef:
	//	f.allocVR_expr(e.X, inloop)
	//	f.allocVR_var(e.Imv, inloop)

	case *Retain:
		f.allocVR_expr(e.X, inloop)

	default:
		panic(fmt.Sprintf("Todo: %s", e.(Stmt).String()))
	}
}

func (f *Function) allocVR_var(v *Alloc, inloop bool) {
	v.tank = rtimp.initTank(v.Type())
	f.allocVR_tank(v.tank, inloop)
}

func (f *Function) allocVR_tank(t *tank, inloop bool) {
	if len(t.member) == 0 {
		t.register.id = f.allocRegister(t.register.typ)
		return
	}

	for _, m := range t.member {
		f.allocVR_tank(m, inloop)
	}
}

func (f *Function) allocRegister(typ Type) (id int) {
	if typ.Kind() == TypeKindChunk {
		l := len(f.freeChunkRegs) - 1
		if l >= 0 {
			id = f.freeChunkRegs[l].id
			f.freeChunkRegs = f.freeChunkRegs[:l]
		} else {
			id = f.regCount
			f.regCount++
			f.chunkRegs = append(f.chunkRegs, register{id: id, typ: typ})
		}
	} else {
		free := false
		for i := len(f.freeCommonRegs) - 1; i >= 0; i-- {
			freeReg := f.freeCommonRegs[i]
			if freeReg.typ.Equal(typ) {
				id = freeReg.id
				f.freeCommonRegs = append(f.freeCommonRegs[:i], f.freeCommonRegs[i+1:]...)
				free = true
				break
			}
		}

		if !free {
			id = f.regCount
			f.regCount++
		}
	}
	return
}

func (f *Function) dropVR_tank(t *tank) {
	if len(t.member) > 0 {
		for _, m := range t.member {
			f.dropVR_tank(m)
		}
	} else {
		f.dropRegister(t.register)
	}
}

func (f *Function) dropRegister(r register) {
	if r.id == -1 {
		panic("dropRegister: id == -1")
	}

	if r.typ.Kind() == TypeKindChunk {
		f.freeChunkRegs = append(f.freeChunkRegs, r)
	} else {
		f.freeCommonRegs = append(f.freeCommonRegs, r)
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
