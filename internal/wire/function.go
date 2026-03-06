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
			param.init = NewGet(&np, np.pos)

			rc_stmt(param, false, fb)

			f.params[i] = &np
		}
	}

	// 返回置换
	for _, ret := range f.results {
		fb.emit(ret)
	}
	ob.Label = _FN_START
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
			rr := newImv(fb.newTempVarName(), NewGet(r, f.EndPos), f.EndPos)
			fb.emit(rr)
			ret_exprs = append(ret_exprs, rr)

			//rr := *r
			//rr.Stringer = &rr
			//rr.kind = Register
			//rr.name = fb.newTempVarName()
			//fb.emit(&rr)
			//fb.EmitSet(&rr, NewGet(r, f.EndPos), f.EndPos)
			//ret_exprs = append(ret_exprs, NewGet(&rr, f.EndPos))

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
		sb.WriteString("\n<=======allocVR_block() 后=======>\n")
		f.Format("  ", &sb)
		println(sb.String())
	}

	// Get指令替换
	getReplace_stmt(f.Body)

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

		for i, s := range stmt.Stmts {
			if ret_stmt, ok := s.(*Return); ok {
				if len(f.results) > 0 && len(ret_stmt.Results) > 0 {
					lhs := make([]Location, len(f.results))
					for i, lh := range f.results {
						lhs[i] = lh
					}

					stmt.Stmts[i] = NewSetN(lhs, ret_stmt.Results, ret_stmt.pos)
					stmt.Stmts = stmt.Stmts[:i+1]
				} else {
					stmt.Stmts = stmt.Stmts[:i]
				}
				if stmt.Label != _FN_START {
					stmt.Stmts = append(stmt.Stmts, NewBr(_FN_START, ret_stmt.pos))
				}
				break
			}

			f.retRepalce(s)
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
		if s.Kind() != Register && s.init != nil {
			init := s.init
			s.init = nil
			d.emit(s)
			rc_stmt(NewSet(s, init, s.pos), inloop, d)
		} else {
			d.emit(s)
		}

	case *Get:
		panic("Get should not be here")

	case *Set:
		for i := range s.Lhs {
			s.Lhs[i] = rc_location(s.Lhs[i], inloop, false, d, &pre)
		}

		for i := range s.Rhs {
			s.Rhs[i] = rc_expr(s.Rhs[i], inloop, false, false, d, &pre)
		}

		//rhNeedRetain := make([]bool, len(s.Lhs))
		//if len(s.Lhs) == len(s.Rhs) {
		//	for i := range s.Rhs {
		//		rc_expr(s.Rhs[i], inloop, false, false, d, &pre)
		//		if rtimp.hasChunk(s.Rhs[i].Type()) && !s.Rhs[i].retained() {
		//			rhNeedRetain[i] = true
		//		}
		//	}
		//} else {
		//	// rh 是元组
		//	notRetained := rtimp.hasChunk(s.Rhs[0].Type()) && !s.Rhs[0].retained()
		//	for i := range rhNeedRetain {
		//		rhNeedRetain[i] = notRetained
		//	}
		//}

		d.Stmts = append(d.Stmts, pre...)

		allLhsAssignable := true
		lhAssignable := make([]bool, len(s.Lhs))
		lhs := make([]Var, len(s.Lhs))
		for i, lh := range s.Lhs {
			if lh == nil {
				lhs[i] = nil
				lhAssignable[i] = true
				continue
			}

			if v, ok := lh.(Var); ok {
				if v.Kind() == Register { // v: T; v = expr; v 未逃逸
					lhs[i] = v
					lhAssignable[i] = true
				} else { // v: T; v = expr;  v 逃逸
					lhs[i] = v
					lhAssignable[i] = false
					allLhsAssignable = false
				}
				continue
			}

			allLhsAssignable = false
			lhAssignable[i] = false

			lhe := loc2expr(lh)
			if get, ok := lhe.(*Get); ok {
				if v, ok := get.Loc.(Var); ok && v.Kind() == Register { // p: *T; *p = expr; p 未逃逸
					lhs[i] = v
					continue
				}
			}

			// 其余情况，左部是指针或引用类型的表达式：
			loc := newImv(d.newTempVarName(), lhe, s.pos)
			d.emit(loc)
			lhs[i] = loc
		}

		if allLhsAssignable {
			assign := newAssignN(lhs, s.Rhs, s.pos)
			//copy(assign.needRetain, rhNeedRetain)
			d.emit(assign)
		} else {
			rhs := make([]Expr, len(s.Lhs))

			if len(s.Lhs) > len(s.Rhs) {
				// 元组展开
				if s.Rhs[0].Type().Kind() != TypeKindTuple {
					panic("RH is not a tuple")
				}

				tuple := newImv(d.newTempVarName(), s.Rhs[0], s.pos)
				d.emit(tuple)

				for i := range s.Lhs {
					rhs[i] = newExtract(tuple, i, s.pos)
				}
			} else {
				copy(rhs, s.Rhs)
			}

			for i := range lhs {
				loc := lhs[i]
				rh := getReplace_expr(rhs[i])

				if lhAssignable[i] {
					assign := newAssign(loc, rh, s.pos)
					//assign.needRetain[0] = rhNeedRetain[i]
					d.emit(assign)
					continue
				}

				if c, ok := rh.(*Const); ok {
					store := newStore(loc, c, s.pos)
					//store.needRetain = false
					d.emit(store)
					continue
				}

				if v, ok := rh.(Var); ok {
					if v.Kind() != Register {
						panic(fmt.Sprintf("rh: %s is not a Register", v.Name()))
					}
					store := newStore(loc, v, s.pos)
					//store.needRetain = rhNeedRetain[i]
					d.emit(store)
					continue
				}

				imv := newImv(d.newTempVarName(), rh, s.pos)
				d.emit(imv)
				store := newStore(loc, imv, s.pos)
				//store.needRetain = rhNeedRetain[i]
				d.emit(store)
			}
		}

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

	//case *Retain:
	//	d.emit(s)

	case *Drop:
		d.emit(s)

	default:
		panic(fmt.Sprintf("Todo: %s", s.String()))
	}
}

func rc_expr(expr Expr, inloop bool, replace bool, isLoopCond bool, d *Block, pre *[]Stmt) (ret Expr) {
	ret = expr
	if ret == nil {
		return
	}
	if _, ok := ret.(*Const); ok {
		return
	}

	switch e := ret.(type) {
	case *Alloc:
		return

	case *Get:
		e.Loc = rc_location(e.Loc, inloop, isLoopCond, d, pre)

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

	case *Member:
		return

	case *MemberAddr:
		e.X = rc_expr(e.X, inloop, true, isLoopCond, d, pre)

	case *asAddr:
		loc := rc_location(e.loc, inloop, isLoopCond, d, pre)
		ret = loc2expr(loc)
		if !ret.Type().Equal(e.Type()) {
			panic("asAddr() type mismatch")
		}

	case Stmt:
		panic(fmt.Sprintf("Todo: %s", e.String()))
	}

	if ret.retained() && replace {
		if isLoopCond {
			tmp := d.NewAlloc(d.newTempVarName(), ret.Type(), ret.Pos(), nil, nil)
			tmp.noinit = true
			d.emit(tmp)

			*pre = append(*pre, newAssign(tmp, ret, ret.Pos()))
			ret = NewGet(tmp, ret.Pos())
		} else {
			tmp := d.NewAlloc(d.newTempVarName(), ret.Type(), ret.Pos(), nil, ret)
			d.emit(tmp)

			ret = NewGet(tmp, ret.Pos())
		}
	}

	return
}

func rc_location(loc Location, inloop bool, isLoopCond bool, d *Block, pre *[]Stmt) (ret Location) {
	ret = loc
	if ret == nil {
		return
	}

	switch e := ret.(type) {
	case *Alloc, *Imv:
		return

	case *asLoc:
		e.expr = rc_expr(e.expr, inloop, true, isLoopCond, d, pre)
		e.expr = getReplace_expr(e.expr)

	case *MemberLocation:
		e.X = rc_location(e.X, inloop, isLoopCond, d, pre)

	default:
		panic(fmt.Sprintf("Todo: %T", e))
	}

	return
}

func loc2expr(loc Location) Expr {
	switch loc := loc.(type) {
	case Expr:
		return loc

	case *MemberLocation:
		x := loc2expr(loc.X)
		if v, ok := x.(Var); ok && v.Kind() == Register && unname(v.Type()).Kind() == TypeKindStruct {
			return newMember(v, loc.Id, loc.pos)
		} else {
			return newMemberAddr(x, loc.Id, loc.pos, loc.types)
		}

	case *asLoc:
		return loc.expr

	default:
		panic(fmt.Sprintf("Todo: %T", loc))
	}
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
			t = append(t, newDrop(s, b.Stmts[r.last].Pos()))
			t = append(t, b.Stmts[r.last+1:]...)
			b.Stmts = t

		case *Imv:
			r := b.varUsageRange(s)
			if r.last == -1 {
				panic(fmt.Sprintf("var:%s not used", s.Name()))
			}

			var t []Stmt
			t = append(t, b.Stmts[:r.last+1]...)
			t = append(t, newDrop(s, b.Stmts[r.last].Pos()))
			t = append(t, b.Stmts[r.last+1:]...)
			b.Stmts = t

		case *If:
			f.autoDrop(s.True)
			f.autoDrop(s.False)

		case *Loop:
			f.autoDrop(s.Body)
			f.autoDrop(s.Post)

		case *Get, *Set, *Assign, *Store, *Br, *Return, *Unop, *Biop, *Drop:
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
		f.allocVR_var(s, inloop)

	case *Imv:
		f.allocVR_imv(s, inloop)

	case *If:
		f.allocVR_block(s.True, inloop)
		f.allocVR_block(s.False, inloop)

	case *Loop:
		f.allocVR_block(s.Body, true)
		f.allocVR_block(s.Post, true)

	case *Drop:
		f.dropVR_tank(s.X.Tank())

	case *Br, *Return, *Get, *Set, *Assign, *Store:

	default:
		panic(fmt.Sprintf("Todo: %s", s.String()))
	}
}

func (f *Function) allocVR_var(v *Alloc, inloop bool) {
	v.tank = rtimp.initTank(v.Type(), false)
	f.allocVR_tank(v.tank, inloop)
}

func (f *Function) allocVR_imv(v *Imv, inloop bool) {
	v.tank = rtimp.initTank(v.Type(), true)
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

func getReplace_stmt(stmt Stmt) {
	switch stmt := stmt.(type) {
	case *Block:
		for _, s := range stmt.Stmts {
			getReplace_stmt(s)
		}

	case *Alloc:
		stmt.init = getReplace_expr(stmt.init)

	case *Imv:
		stmt.val = getReplace_expr(stmt.val)

	case *Set:
		panic("Set should not be here")

	case *Assign:
		for i := range stmt.Rhs {
			stmt.Rhs[i] = getReplace_expr(stmt.Rhs[i])
		}

	case *Store:
		// Store 指令的 Loc 和 Val 都是 Var，无需替换
		return

	case *Return:
		for i := range stmt.Results {
			stmt.Results[i] = getReplace_expr(stmt.Results[i])
		}

	case *If:
		stmt.Cond = getReplace_expr(stmt.Cond)
		getReplace_stmt(stmt.True)
		getReplace_stmt(stmt.False)

	case *Loop:
		for _, pre := range stmt.PreCond {
			getReplace_stmt(pre)
		}
		stmt.Cond = getReplace_expr(stmt.Cond)
		getReplace_stmt(stmt.Body)
		getReplace_stmt(stmt.Post)

	case *Br, *Drop:
		return

	default:
		panic(fmt.Sprintf("Todo: %s", stmt.String()))

	}
}

func getReplace_expr(e Expr) (ret Expr) {
	ret = e
	if e == nil {
		return
	}
	if _, ok := e.(*Const); ok {
		return
	}

	switch e := e.(type) {
	case *Get:
		loc := loc2expr(e.Loc)
		if v, ok := loc.(Var); !ok {
			return newLoad(getReplace_expr(loc), e.Pos())
		} else {
			if v.Kind() == Register {
				return v
			} else {
				return newLoad(v, e.Pos())
			}
		}

	case *Load:
		e.Loc = getReplace_expr(e.Loc)

	case *Extract:
		return

	case *Unop:
		e.X = getReplace_expr(e.X)

	case *Biop:
		e.X = getReplace_expr(e.X)
		e.Y = getReplace_expr(e.Y)

	case *Call:
		var call_common *CallCommon
		switch call := e.Callee.(type) {
		case *BuiltinCall:
			call_common = &call.CallCommon

		case *StaticCall:
			call_common = &call.CallCommon

		case *MethodCall:
			call_common = &call.CallCommon
			call.Recv = getReplace_expr(call.Recv)

		case *InterfaceCall:
			call_common = &call.CallCommon
			call.Interface = getReplace_expr(call.Interface)

		default:
			panic("Todo")
		}

		for i := range call_common.Args {
			call_common.Args[i] = getReplace_expr(call_common.Args[i])
		}

	case *Member:
		return

	case *MemberAddr:
		e.X = getReplace_expr(e.X)

	case *Alloc, *Imv:
		return

	default:
		panic(fmt.Sprintf("Todo: %t", e))

	}

	return
}
