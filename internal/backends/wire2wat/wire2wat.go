// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire2wat

import (
	"fmt"
	"strconv"
	"strings"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/wire"
)

func NewWire2Wat() *Wire2Wat {
	w := &Wire2Wat{}
	w.init()
	return w
}

func (w *Wire2Wat) BuildFunction(f *wire.Function) {
	wf := &wat.Function{}
	wf.InternalName = f.InternalName
	wf.ExternalName = f.ExternalName

	paramRegNum := 0
	for _, param := range f.Params {
		regs := param.Tank().Raw()
		for _, reg := range regs {
			paramRegNum++
			param := wat.NewVar(reg.String(), w.watType(reg.OType()))
			wf.Params = append(wf.Params, param)
		}
	}

	for _, result := range f.Results {
		raw := w.rawof(result.Type())
		for _, t := range raw {
			wf.Results = append(wf.Results, w.watType(t))
		}
	}

	for i := paramRegNum; i < len(f.CommonRegs); i++ {
		r := f.CommonRegs[i]
		wf.Locals = append(wf.Locals, wat.NewVar(r.String(), w.watType(r.OType())))
	}
	for _, r := range f.ChunkRegs {
		wf.Locals = append(wf.Locals, wat.NewVar(r.String(), w.watType(r.OType())))
	}

	for _, stmt := range f.Body.Stmts {
		wf.Insts = append(wf.Insts, w.emitStmt(stmt)...)
	}

	sb := strings.Builder{}
	wf.Format(&sb)
	println(sb.String())
}

func (w *Wire2Wat) emitStmt(stmt wire.Stmt) (ret []wat.Inst) {
	switch s := stmt.(type) {
	case *wire.Block:
		block := wat.NewInstBlock(s.Label)
		for _, s := range s.Stmts {
			block.Insts = append(block.Insts, w.emitStmt(s)...)
		}
		ret = append(ret, wat.NewComment(s.Label))
		ret = append(ret, block)

	case *wire.Alloc:
		ret = append(ret, wat.NewComment(s.String()))
		ret = append(ret, w.emitAlloc(s)...)

	case *wire.Imv:
		ret = append(ret, wat.NewComment(s.String()))
		ret = append(ret, w.emitImv(s)...)

	case *wire.Assign:
		ret = append(ret, wat.NewComment(s.String()))
		ret = append(ret, w.emitAssign(s)...)

	case *wire.Store:
		ret = append(ret, wat.NewComment(s.String()))
		ret = append(ret, w.emitStore(s)...)

	case *wire.If:
		ret = append(ret, wat.NewComment("if "+s.Cond.Name()))
		ret = append(ret, w.emitIf(s)...)

	case *wire.Loop:
		ret = append(ret, wat.NewComment("loop "+s.Label))
		ret = append(ret, w.emitLoop(s)...)

	case *wire.Discard:
		return

	case *wire.Release:
		ret = append(ret, wat.NewComment(stmt.String()))
		ret = append(ret, w.emitRelease_register(s.X)...)

	case *wire.Return:
		ret = append(ret, wat.NewComment(stmt.String()))
		ret = append(ret, w.emitReturn(s)...)

	case *wire.NilCheck:
		ret = append(ret, wat.NewComment(stmt.String()))
		ret = append(ret, w.emitNilCheck(s)...)

	default:
		panic(fmt.Sprintf("Todo: %T", stmt))
	}

	ret = append(ret, wat.NewBlank())
	return
}

func (w *Wire2Wat) emitAlloc(s *wire.Alloc) (ret []wat.Inst) {
	switch s.Kind() {
	case wire.AllocKindRegister:
		if s.Init != nil {
			ret = append(ret, w.emitExpr(s.Init)...)
			rhRtained, rhConst := w.exprRCState(s.Init)
			ret = append(ret, w.emitPop_tank(s.Tank(), true, !rhRtained && !rhConst)...)
		} else {
			ret = append(ret, w.emitZero_tank(s.Tank(), true)...)
		}

	case wire.AllocKindHeap:
		ss := fmt.Sprintf("%d", w.sizeof(s.DataType()))
		ret = append(ret, wat.NewInstConst(w.baseTypes[wire.TypeKindInt].watType, "1")) // item_count
		ret = append(ret, wat.NewInstConst(w.baseTypes[wire.TypeKindInt].watType, "0")) // todo: free_method
		ret = append(ret, wat.NewInstConst(w.baseTypes[wire.TypeKindInt].watType, ss))  // item_size
		ret = append(ret, wat.NewInstCall("runtime.Block.HeapAlloc"))
		ret = append(ret, w.emitPop_tank(s.Tank(), true, false)...)
	}

	return
}

func (w *Wire2Wat) emitImv(s *wire.Imv) (ret []wat.Inst) {
	ret = append(ret, w.emitExpr(s.Val)...)
	ret = append(ret, w.emitPop_tank(s.Tank(), false, false)...)
	return
}

func (w *Wire2Wat) exprRCState(expr wire.Expr) (retained bool, isConst bool) {
	retained = expr.Retained()
	_, isConst = expr.(*wire.Const)
	return
}

func (w *Wire2Wat) emitAssign(s *wire.Assign) (ret []wat.Inst) {
	for _, rh := range s.Rhs {
		ret = append(ret, w.emitExpr(rh)...)
	}

	for i := len(s.Lhs) - 1; i >= 0; i-- {
		lh := s.Lhs[i]

		if _, ok := lh.(*wire.Imv); ok {
			panic("Imv is not supported")
		}

		var rhRtained, rhConst bool
		if i < len(s.Rhs) {
			rhRtained, rhConst = w.exprRCState(s.Rhs[i])
		} else {
			rhRtained, rhConst = w.exprRCState(s.Rhs[0])
		}

		if lh == nil {
			ret = append(ret, w.emitDrop(s.RhsType[i], rhRtained && !rhConst)...)
		} else {
			ret = append(ret, w.emitPop_tank(lh.Tank(), true, !rhRtained && !rhConst)...)
		}
	}

	return
}

func (w *Wire2Wat) emitStore(s *wire.Store) (ret []wat.Inst) {
	var ptr wire.Register
	switch s.Loc.Type().(type) {
	case *wire.Ref:
		ptr = s.Loc.Tank().Member[1].Register

	case *wire.Ptr:
		ptr = s.Loc.Tank().Register

	default:
		panic(fmt.Sprintf("Todo: %T", s.Loc.Type()))
	}

	retained, isConst := w.exprRCState(s.Val)
	switch val := s.Val.(type) {
	case wire.Var:
		ret = w.emitStore_Tank(ptr, 0, val.Tank(), !retained && !isConst)

	case *wire.Const:
		ret = w.emitStore_const(ptr, 0, val)

	default:
		panic(fmt.Sprintf("Todo: %T", val))
	}

	return
}

func (w *Wire2Wat) emitStore_Tank(ptr wire.Register, offset int, tank *wire.Tank, needRetain bool) (ret []wat.Inst) {
	if len(tank.Member) == 0 {
		return w.emitStore_register(ptr, offset, tank.Register, needRetain)
	}

	t, ok := tank.Type().(*wire.Struct)
	if !ok {
		panic("tank.Type() is not a struct")
	}
	if !t.SizeInitialized {
		w.initStructSize(t)
	}
	for i := 0; i < t.Len(); i++ {
		member := t.At(i)
		ret = append(ret, w.emitStore_Tank(ptr, offset+member.Start, tank.Member[i], needRetain)...)
	}

	return
}

func (w *Wire2Wat) emitIf(s *wire.If) (ret []wat.Inst) {
	ret = append(ret, w.emitExpr(s.Cond)...)

	i := wat.NewInstIf(nil, nil, nil)
	for _, s := range s.True.Stmts {
		i.True = append(i.True, w.emitStmt(s)...)
	}
	for _, s := range s.False.Stmts {
		i.False = append(i.False, w.emitStmt(s)...)
	}

	ret = append(ret, i)
	return
}

func (w *Wire2Wat) emitLoop(s *wire.Loop) (ret []wat.Inst) {
	loopLabel := s.Label + ".loop"
	loop := wat.NewInstLoop(loopLabel)

	for _, preCond := range s.PreCond {
		loop.Insts = append(loop.Insts, w.emitStmt(preCond)...)
	}
	loop.Insts = append(loop.Insts, wat.NewComment(s.Cond.Name()))
	loop.Insts = append(loop.Insts, w.emitExpr(s.Cond)...)

	body := wat.NewInstBlock(s.Label + ".body")
	for _, stmt := range s.Body.Stmts {
		body.Insts = append(body.Insts, w.emitStmt(stmt)...)
	}

	_if := wat.NewInstIf(nil, nil, nil)
	_if.Name = s.Label + ".do"
	_if.True = append(_if.True, body)
	for _, post := range s.Post.Stmts {
		_if.True = append(_if.True, w.emitStmt(post)...)
	}
	_if.True = append(_if.True, wat.NewInstBr(loopLabel))

	loop.Insts = append(loop.Insts, _if)
	ret = append(ret, loop)
	return
}

func (w *Wire2Wat) emitStore_register(ptr wire.Register, offset int, value wire.Register, needRetain bool) (ret []wat.Inst) {
	hasChunk := value.OType().HasChunk()

	ret = append(ret, w.emitPush_register(ptr))   // p
	ret = append(ret, w.emitPush_register(value)) // p v

	if needRetain && hasChunk {
		ret = append(ret, w.emitRetain())
	}

	if hasChunk {
		ret = append(ret, w.emitLoad_BaseType(ptr, offset, value.OType())...) // p v o
		ret = append(ret, w.emitRelease_stack())                              // p v
	}

	switch value.OType().Kind() {
	case wire.TypeKindBool, wire.TypeKindU8, wire.TypeKindI8:
		ret = append(ret, wat.NewInstStore8(offset, 1))

	case wire.TypeKindI16, wire.TypeKindU16:
		ret = append(ret, wat.NewInstStore16(offset, 2))

	default:
		ret = append(ret, wat.NewInstStore(w.watType(value.OType()), offset, w.sizeof(value.OType())))
	}

	return ret
}

func (w *Wire2Wat) emitStore_const(ptr wire.Register, offset int, value *wire.Const) (ret []wat.Inst) {
	if value.Type().Kind() >= wire.BaseTypeNum {
		panic(fmt.Sprintf("Todo: %T", value.Type()))
	}

	ret = append(ret, w.emitPush_register(ptr))
	ret = append(ret, wat.NewInstConst(w.watType(value.Type()), value.Name()))

	switch value.Type().Kind() {
	case wire.TypeKindBool, wire.TypeKindU8, wire.TypeKindI8:
		ret = append(ret, wat.NewInstStore8(offset, 1))

	case wire.TypeKindI16, wire.TypeKindU16:
		ret = append(ret, wat.NewInstStore16(offset, 2))

	default:
		ret = append(ret, wat.NewInstStore(w.watType(value.Type()), offset, w.sizeof(value.Type())))
	}

	return ret
}

func (w *Wire2Wat) emitReturn(s *wire.Return) (ret []wat.Inst) {
	for _, expr := range s.Results {
		ret = append(ret, w.emitExpr(expr)...)
	}
	ret = append(ret, wat.NewInstReturn())
	return
}

func (w *Wire2Wat) emitPush_tank(tank *wire.Tank) (ret []wat.Inst) {
	regs := tank.Raw()
	for _, reg := range regs {
		ret = append(ret, w.emitPush_register(reg))
	}
	return
}

func (w *Wire2Wat) emitPush_register(reg wire.Register) wat.Inst {
	if reg.Kind() == wire.KGlobal {
		return wat.NewInstGetGlobal(reg.String())
	} else {
		return wat.NewInstGetLocal(reg.String())
	}
}

func (w *Wire2Wat) emitPop_tank(tank *wire.Tank, needRelease bool, needRetain bool) (ret []wat.Inst) {
	regs := tank.Raw()
	for i := len(regs) - 1; i >= 0; i-- {
		ret = append(ret, w.emitPop_register(regs[i], needRelease, needRetain)...)
	}
	return
}

func (w *Wire2Wat) emitPop_register(reg wire.Register, needRelease bool, needRetain bool) (ret []wat.Inst) {
	hasChunk := reg.OType().HasChunk()

	if needRetain && hasChunk {
		ret = append(ret, w.emitRetain())
	}

	if needRelease && hasChunk {
		ret = append(ret, w.emitRelease_register(reg)...)
	}

	if reg.Kind() == wire.KGlobal {
		ret = append(ret, wat.NewInstSetGlobal(reg.String()))
	} else {
		ret = append(ret, wat.NewInstSetLocal(reg.String()))
	}
	return
}

func (w *Wire2Wat) emitDrop(t wire.Type, needRelease bool) (ret []wat.Inst) {
	raw := w.rawof(t)
	for i := len(raw) - 1; i >= 0; i-- {
		ret = append(ret, w.emitDrop_BaseType(raw[i], needRelease))
	}
	return
}

func (w *Wire2Wat) emitDrop_BaseType(t wire.Type, needRelease bool) (ret wat.Inst) {
	if t.Kind() >= wire.BaseTypeNum {
		panic("t is not a base type")
	}
	if needRelease && t.HasChunk() {
		ret = wat.NewInstCall("runtime.Block.Release")
	} else {
		ret = wat.NewInstDrop()
	}
	return
}

func (w *Wire2Wat) emitZero_tank(tank *wire.Tank, needRelease bool) (ret []wat.Inst) {
	regs := tank.Raw()
	for i := len(regs) - 1; i >= 0; i-- {
		ret = append(ret, w.emitZero_register(regs[i], needRelease)...)
	}
	return
}

func (w *Wire2Wat) emitZero_register(reg wire.Register, needRelease bool) (ret []wat.Inst) {
	ret = append(ret, wat.NewInstConst(w.watType(reg.OType()), "0"))
	ret = append(ret, w.emitPop_register(reg, needRelease, false)...)
	return
}

func (w *Wire2Wat) emitRetain() wat.Inst {
	return wat.NewInstCall("runtime.Block.Retain")
}

func (w *Wire2Wat) emitRelease_register(reg wire.Register) (ret []wat.Inst) {
	ret = make([]wat.Inst, 2)
	ret[0] = w.emitPush_register(reg)
	ret[1] = w.emitRelease_stack()
	return
}

func (w *Wire2Wat) emitRelease_stack() wat.Inst {
	return wat.NewInstCall("runtime.Block.Release")
}

func (w *Wire2Wat) emitLoad_BaseType(ptr wire.Register, offset int, t wire.Type) (ret []wat.Inst) {
	if t.Kind() >= wire.BaseTypeNum {
		panic("t is not a base type")
	}

	ret = make([]wat.Inst, 2)
	ret[0] = w.emitPush_register(ptr)

	switch t.Kind() {
	case wire.TypeKindBool, wire.TypeKindU8:
		ret[1] = wat.NewInstLoad8u(offset, 1)
	case wire.TypeKindI8:
		ret[1] = wat.NewInstLoad8s(offset, 1)
	case wire.TypeKindI16:
		ret[1] = wat.NewInstLoad16s(offset, 2)
	case wire.TypeKindU16:
		ret[1] = wat.NewInstLoad16u(offset, 2)

	default:
		ret[1] = wat.NewInstLoad(w.watType(t), offset, w.sizeof(t))
	}
	return
}

func (w *Wire2Wat) emitExpr(expr wire.Expr) (ret []wat.Inst) {
	if c, ok := expr.(*wire.Const); ok {
		if expr.Type().Kind() < wire.BaseTypeNum {
			ret = append(ret, wat.NewInstConst(w.watType(expr.Type()), c.Name()))
		} else {
			panic(fmt.Sprintf("Todo: %T", expr))
		}
		return
	}

	switch e := expr.(type) {
	case wire.Var:
		ret = append(ret, w.emitPush_tank(e.Tank())...)

	case *wire.Call:
		var call_common *wire.CallCommon
		switch callee := e.Callee.(type) {
		case *wire.BuiltinCall:
			call_common = &callee.CallCommon
		case *wire.StaticCall:
			call_common = &callee.CallCommon
		default:
			panic(fmt.Sprintf("Todo: %T", callee))
		}

		for _, arg := range call_common.Args {
			ret = append(ret, w.emitExpr(arg)...)
		}
		ret = append(ret, wat.NewInstCall(call_common.String()))

	case *wire.Load:
		ret = append(ret, w.emitLoad(e)...)

	case *wire.MemberAddr:
		ret = append(ret, w.emitMemberAddr(e)...)

	case *wire.Biop:
		ret = append(ret, w.emitBiop(e)...)

	case *wire.Unop:
		ret = append(ret, w.emitUnop(e)...)

	default:
		panic(fmt.Sprintf("Todo: %T", e))
	}

	return
}

func (w *Wire2Wat) emitNilCheck(e *wire.NilCheck) (ret []wat.Inst) {
	var ptr wire.Register
	switch e.X.Type().(type) {
	case *wire.Ref:
		ptr = e.X.Tank().Member[1].Register
	case *wire.Ptr:
		ptr = e.X.Tank().Register
	default:
		panic("X is not a Ref or Ptr")
	}
	ret = append(ret, w.emitPush_register(ptr))
	ret = append(ret, wat.NewInstCall("runtime.Block.NilCheck"))

	return
}

func (w *Wire2Wat) emitLoad(e *wire.Load) (ret []wat.Inst) {
	loc, ok := e.Loc.(wire.Var)
	if !ok {
		panic("Loc is not a Var")
	}

	var ptr wire.Register
	switch e.Loc.Type().(type) {
	case *wire.Ref:
		ptr = loc.Tank().Member[1].Register

	case *wire.Ptr:
		ptr = loc.Tank().Register

	default:
		panic("Loc is not a Ref or Ptr")
	}

	return w.emitLoad_Type(ptr, 0, e.Type())
}

func (w *Wire2Wat) emitLoad_Type(ptr wire.Register, offset int, t wire.Type) (ret []wat.Inst) {
	if t.Kind() < wire.BaseTypeNum {
		return w.emitLoad_BaseType(ptr, offset, t)
	}

	switch t := t.(type) {
	case *wire.Tuple:
		panic("Tuple is not supported")

	case *wire.Struct:
		if !t.SizeInitialized {
			w.initStructSize(t)
		}
		for i := 0; i < t.Len(); i++ {
			member := t.At(i)
			ret = append(ret, w.emitLoad_Type(ptr, offset+member.Start, member.Type)...)
		}

	case *wire.Ref:
		us := t.Underlying()
		if us == nil {
			panic("String.Underlying() is nil")
		}
		return w.emitLoad_Type(ptr, offset, us)

	case *wire.String:
		us := t.Underlying()
		if us == nil {
			panic("String.Underlying() is nil")
		}
		return w.emitLoad_Type(ptr, offset, us)

	case *wire.Named:
		return w.emitLoad_Type(ptr, offset, t.Underlying())

	default:
		panic(fmt.Sprintf("Todo: %T", t))
	}
	return
}

func (w *Wire2Wat) emitMemberAddr(e *wire.MemberAddr) (ret []wat.Inst) {
	x := e.X
	offset := w.memberOffset(e)

	for sm, ok := x.(*wire.MemberAddr); ok; sm, ok = x.(*wire.MemberAddr) {
		x = sm.X
		offset += w.memberOffset(sm)
	}

	ret = append(ret, w.emitExpr(x)...)
	if offset > 0 {
		ret = append(ret, wat.NewInstConst(w.baseTypes[wire.TypeKindPtr].watType, strconv.Itoa(offset)))
		ret = append(ret, wat.NewInstAdd(w.baseTypes[wire.TypeKindPtr].watType))
	}
	return
}

func (w *Wire2Wat) memberOffset(e *wire.MemberAddr) int {
	un := wire.Unnamed(wire.Deref(e.X.Type()))
	if st, ok := un.(*wire.Struct); !ok {
		panic(fmt.Sprintf("Invalid X.Type():%s", un.Name()))
	} else {
		return st.At(e.Id).Start
	}
}

func (w *Wire2Wat) emitBiop(e *wire.Biop) (ret []wat.Inst) {
	switch e.Op {
	case wire.XOR:
		ret = append(ret, w.emitExpr(e.X)...)
		ret = append(ret, w.emitExpr(e.Y)...)
		ret = append(ret, wat.NewInstXor(w.watType(e.X.Type())))

	case wire.LAND:
		ret = append(ret, w.emitExpr(e.X)...)

		i := wat.NewInstIf(nil, nil, nil)
		i.Ret = []wat.ValueType{w.watType(&wire.Bool{})}
		i.True = w.emitExpr(e.Y)
		i.False = []wat.Inst{wat.NewInstConst(w.watType(&wire.Bool{}), "0")}

		ret = append(ret, i)

	case wire.LOR:
		ret = append(ret, w.emitExpr(e.X)...)

		i := wat.NewInstIf(nil, nil, nil)
		i.Ret = []wat.ValueType{w.watType(&wire.Bool{})}
		i.True = []wat.Inst{wat.NewInstConst(w.watType(&wire.Bool{}), "1")}
		i.False = w.emitExpr(e.Y)

		ret = append(ret, i)

	case wire.SHL:
		ret = append(ret, w.emitExpr(e.X)...)
		ret = append(ret, w.emitExpr(e.Y)...)
		if w.sizeof(e.X.Type()) <= 4 && w.sizeof(e.Y.Type()) == 8 {
			ret = append(ret, wat.NewInstConvert_i32_wrap_i64())
		} else if w.sizeof(e.X.Type()) == 8 && w.sizeof(e.Y.Type()) <= 4 {
			ret = append(ret, wat.NewInstConvert_i64_extend_i32_u())
		}
		ret = append(ret, wat.NewInstShl(w.watType(e.X.Type())))

	case wire.SHR:
		ret = append(ret, w.emitExpr(e.X)...)
		ret = append(ret, w.emitExpr(e.Y)...)
		if w.sizeof(e.X.Type()) <= 4 && w.sizeof(e.Y.Type()) == 8 {
			ret = append(ret, wat.NewInstConvert_i32_wrap_i64())
		} else if w.sizeof(e.X.Type()) == 8 && w.sizeof(e.Y.Type()) <= 4 {
			ret = append(ret, wat.NewInstConvert_i64_extend_i32_u())
		}
		ret = append(ret, wat.NewInstShr(w.watType(e.X.Type())))

	case wire.ADD:
		ret = append(ret, w.emitExpr(e.X)...)
		ret = append(ret, w.emitExpr(e.Y)...)

		switch e.X.Type().Kind() {
		case wire.TypeKindString, wire.TypeKindComplex64, wire.TypeKindComplex128:
			panic("Todo")

		default:
			ret = append(ret, wat.NewInstAdd(w.watType(e.X.Type())))
		}

	case wire.SUB:
		ret = append(ret, w.emitExpr(e.X)...)
		ret = append(ret, w.emitExpr(e.Y)...)

		switch e.X.Type().Kind() {
		case wire.TypeKindComplex64, wire.TypeKindComplex128:
			panic("Todo")

		default:
			ret = append(ret, wat.NewInstSub(w.watType(e.X.Type())))
		}

	case wire.MUL:
		ret = append(ret, w.emitExpr(e.X)...)
		ret = append(ret, w.emitExpr(e.Y)...)

		switch e.X.Type().Kind() {
		case wire.TypeKindComplex64, wire.TypeKindComplex128:
			panic("Todo")

		default:
			ret = append(ret, wat.NewInstMul(w.watType(e.X.Type())))
		}

	case wire.QUO:
		switch e.X.Type().Kind() {
		case wire.TypeKindComplex64, wire.TypeKindComplex128:
			panic("Todo")

		case wire.TypeKindI8:
			ret = append(ret, w.emitExpr(e.X)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))
			ret = append(ret, w.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))

			ret = append(ret, wat.NewInstDiv(wat.I32{}))

		case wire.TypeKindI16:
			ret = append(ret, w.emitExpr(e.X)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))
			ret = append(ret, w.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))

			ret = append(ret, wat.NewInstDiv(wat.I32{}))

		default:
			ret = append(ret, w.emitExpr(e.X)...)
			ret = append(ret, w.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstDiv(w.watType(e.X.Type())))
		}

	case wire.REM:
		switch e.X.Type().Kind() {
		case wire.TypeKindComplex64, wire.TypeKindComplex128:
			panic("Todo")

		case wire.TypeKindI8:
			ret = append(ret, w.emitExpr(e.X)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))
			ret = append(ret, w.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))

			ret = append(ret, wat.NewInstRem(wat.I32{}))

		case wire.TypeKindI16:
			ret = append(ret, w.emitExpr(e.X)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))
			ret = append(ret, w.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))

			ret = append(ret, wat.NewInstRem(wat.I32{}))

		default:
			ret = append(ret, w.emitExpr(e.X)...)
			ret = append(ret, w.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstRem(w.watType(e.X.Type())))
		}

	case wire.AND:
		ret = append(ret, w.emitExpr(e.X)...)
		ret = append(ret, w.emitExpr(e.Y)...)
		ret = append(ret, wat.NewInstAnd(w.watType(e.X.Type())))

	case wire.ANDNOT:
		t := w.watType(e.X.Type())
		ret = append(ret, w.emitExpr(e.X)...)
		ret = append(ret, w.emitExpr(e.Y)...)
		ret = append(ret, wat.NewInstConst(t, "-1"))
		ret = append(ret, wat.NewInstXor(t))
		ret = append(ret, wat.NewInstAnd(t))

	case wire.OR:
		ret = append(ret, w.emitExpr(e.X)...)
		ret = append(ret, w.emitExpr(e.Y)...)
		ret = append(ret, wat.NewInstOr(w.watType(e.X.Type())))

	case wire.EQL:
		ret = append(ret, w.emitEql(e.X, e.Y)...)

	case wire.NEQ:
		ret = append(ret, w.emitEql(e.X, e.Y)...)
		ret = append(ret, wat.NewInstEqz(wat.I32{}))

	case wire.GTR:
		ret = append(ret, w.emitExpr(e.X)...)
		ret = append(ret, w.emitExpr(e.Y)...)
		if e.X.Type().Kind() == wire.TypeKindString {
			ret = append(ret, wat.NewInstCall("$wa.runtime.string_GTR"))
		} else {
			ret = append(ret, wat.NewInstGt(w.watType(e.X.Type())))
		}

	case wire.LSS:
		ret = append(ret, w.emitExpr(e.X)...)
		ret = append(ret, w.emitExpr(e.Y)...)
		if e.X.Type().Kind() == wire.TypeKindString {
			ret = append(ret, wat.NewInstCall("$wa.runtime.string_LSS"))
		} else {
			ret = append(ret, wat.NewInstLt(w.watType(e.X.Type())))
		}

	case wire.GEQ:
		ret = append(ret, w.emitExpr(e.X)...)
		ret = append(ret, w.emitExpr(e.Y)...)
		if e.X.Type().Kind() == wire.TypeKindString {
			ret = append(ret, wat.NewInstCall("$wa.runtime.string_GEQ"))
		} else {
			ret = append(ret, wat.NewInstGe(w.watType(e.X.Type())))
		}

	case wire.LEQ:
		ret = append(ret, w.emitExpr(e.X)...)
		ret = append(ret, w.emitExpr(e.Y)...)
		if e.X.Type().Kind() == wire.TypeKindString {
			ret = append(ret, wat.NewInstCall("$wa.runtime.string_LEQ"))
		} else {
			ret = append(ret, wat.NewInstLe(w.watType(e.X.Type())))
		}

	case wire.LEG:
		if e.X.Type().Kind() == wire.TypeKindString {
			ret = append(ret, w.emitExpr(e.X)...)
			ret = append(ret, w.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstCall("$wa.runtime.string_Comp"))
		} else {
			t := w.watType(e.X.Type())
			ret = append(ret, w.emitExpr(e.X)...)
			ret = append(ret, w.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstLt(t))

			inst_lt := wat.NewInstIf(nil, nil, nil)
			inst_lt.Ret = append(inst_lt.Ret, wat.I32{})
			inst_lt.True = append(inst_lt.True, wat.NewInstConst(wat.I32{}, "-1"))
			inst_lt.False = append(inst_lt.False, w.emitExpr(e.X)...)
			inst_lt.False = append(inst_lt.False, w.emitExpr(e.Y)...)
			inst_lt.False = append(inst_lt.False, wat.NewInstGt(t))

			ret = append(ret, inst_lt)
		}

	default:
		panic(fmt.Sprintf("Todo: %s", e.Op.String()))
	}

	switch e.Type().Kind() {
	case wire.TypeKindU8:
		ret = append(ret, wat.NewInstConst(w.watType(e.X.Type()), "255"))
		ret = append(ret, wat.NewInstAnd(w.watType(e.X.Type())))
	case wire.TypeKindU16:
		ret = append(ret, wat.NewInstConst(w.watType(e.X.Type()), "65535"))
		ret = append(ret, wat.NewInstAnd(w.watType(e.X.Type())))
	case wire.TypeKindI8:
		ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
		ret = append(ret, wat.NewInstShl(wat.I32{}))
		ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
		ret = append(ret, wat.NewInstShr(wat.I32{}))
	case wire.TypeKindI16:
		ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
		ret = append(ret, wat.NewInstShl(wat.I32{}))
		ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
		ret = append(ret, wat.NewInstShr(wat.I32{}))
	}

	return
}

func (w *Wire2Wat) emitEql(x, y wire.Expr) (ret []wat.Inst) {
	if x.Type().Kind() < wire.BaseTypeNum {
		ret = append(ret, w.emitExpr(x)...)
		ret = append(ret, w.emitExpr(y)...)
		ret = append(ret, wat.NewInstEq(w.watType(x.Type())))
	} else {
		panic("Todo")
	}

	return
}

func (w *Wire2Wat) emitUnop(e *wire.Unop) (ret []wat.Inst) {
	switch e.Op {
	case wire.NOT:
		ret = append(ret, w.emitExpr(e.X)...)
		ret = append(ret, wat.NewInstEqz(wat.I32{}))

	case wire.NEG:
		switch e.X.Type().Kind() {
		case wire.TypeKindF32, wire.TypeKindF64:
			ret = append(ret, w.emitExpr(e.X)...)
			ret = append(ret, wat.NewInstNeg(w.watType(e.X.Type())))
		default:
			ret = append(ret, wat.NewInstConst(w.watType(e.X.Type()), "0"))
			ret = append(ret, w.emitExpr(e.X)...)
			ret = append(ret, wat.NewInstSub(w.watType(e.X.Type())))
		}

	case wire.XOR:
		ret = append(ret, wat.NewInstConst(w.watType(e.X.Type()), "-1"))
		ret = append(ret, w.emitExpr(e.X)...)
		ret = append(ret, wat.NewInstXor(w.watType(e.X.Type())))

	default:
		panic(fmt.Sprintf("Todo: %s", e.Op.String()))
	}

	switch e.Type().Kind() {
	case wire.TypeKindU8:
		ret = append(ret, wat.NewInstConst(w.watType(e.X.Type()), "255"))
		ret = append(ret, wat.NewInstAnd(w.watType(e.X.Type())))
	case wire.TypeKindU16:
		ret = append(ret, wat.NewInstConst(w.watType(e.X.Type()), "65535"))
		ret = append(ret, wat.NewInstAnd(w.watType(e.X.Type())))
	case wire.TypeKindI8:
		ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
		ret = append(ret, wat.NewInstShl(wat.I32{}))
		ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
		ret = append(ret, wat.NewInstShr(wat.I32{}))
	case wire.TypeKindI16:
		ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
		ret = append(ret, wat.NewInstShl(wat.I32{}))
		ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
		ret = append(ret, wat.NewInstShr(wat.I32{}))
	}

	return ret
}
