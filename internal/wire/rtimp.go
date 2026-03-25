// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
	"strconv"
	"strings"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
)

type RtImp interface {
	underlyingStruct(t Type) *Struct
	initTank(t Type, kind RegisterKind) *tank
	hasChunk(t Type) bool
	alignof(t Type) int
	sizeof(t Type) int

	BuildFunction(f *Function)
}

var rtimp RtImp

func init() {
	imp := &rtImp{}
	imp.init()
	rtimp = imp
}

type baseType struct {
	regType  Type
	watType  wat.ValueType
	hasChunk bool
	size     int
}

type rtImp struct {
	baseTypes []baseType
}

func (ri *rtImp) init() {
	ri.baseTypes = make([]baseType, BaseTypeNum)

	ri.baseTypes[TypeKindUnknown] = baseType{regType: nil, watType: nil, hasChunk: false, size: 0}
	ri.baseTypes[TypeKindVoid] = baseType{regType: &Void{}, watType: nil, hasChunk: false, size: 0}
	ri.baseTypes[TypeKindBool] = baseType{regType: &I32{}, watType: wat.I32{}, hasChunk: false, size: 1}
	ri.baseTypes[TypeKindU8] = baseType{regType: &I32{}, watType: wat.U32{}, hasChunk: false, size: 1}
	ri.baseTypes[TypeKindU16] = baseType{regType: &I32{}, watType: wat.U32{}, hasChunk: false, size: 2}
	ri.baseTypes[TypeKindU32] = baseType{regType: &I32{}, watType: wat.U32{}, hasChunk: false, size: 4}
	ri.baseTypes[TypeKindU64] = baseType{regType: &I64{}, watType: wat.U64{}, hasChunk: false, size: 8}
	ri.baseTypes[TypeKindI8] = baseType{regType: &I32{}, watType: wat.I32{}, hasChunk: false, size: 1}
	ri.baseTypes[TypeKindI16] = baseType{regType: &I32{}, watType: wat.I32{}, hasChunk: false, size: 2}
	ri.baseTypes[TypeKindI32] = baseType{regType: &I32{}, watType: wat.I32{}, hasChunk: false, size: 4}
	ri.baseTypes[TypeKindI64] = baseType{regType: &I64{}, watType: wat.I64{}, hasChunk: false, size: 8}
	ri.baseTypes[TypeKindInt] = baseType{regType: &I32{}, watType: wat.I32{}, hasChunk: false, size: 4}
	ri.baseTypes[TypeKindUint] = baseType{regType: &I32{}, watType: wat.U32{}, hasChunk: false, size: 4}
	ri.baseTypes[TypeKindF32] = baseType{regType: &F32{}, watType: wat.F32{}, hasChunk: false, size: 4}
	ri.baseTypes[TypeKindF64] = baseType{regType: &F64{}, watType: wat.F64{}, hasChunk: false, size: 8}
	ri.baseTypes[TypeKindRune] = baseType{regType: &I32{}, watType: wat.U32{}, hasChunk: false, size: 4}
	ri.baseTypes[TypeKindPtr] = baseType{regType: &I32{}, watType: wat.U32{}, hasChunk: false, size: 4}
	ri.baseTypes[TypeKindChunk] = baseType{regType: &chunk{Base: &Void{}}, watType: wat.U32{}, hasChunk: true, size: 4}
}

func (ri *rtImp) hasChunk(t Type) bool {
	if t.Kind() < BaseTypeNum {
		return ri.baseTypes[t.Kind()].hasChunk
	}

	switch t := t.(type) {
	case *Tuple:
		for _, m := range t.members {
			if ri.hasChunk(m) {
				return true
			}
		}

	case *Struct:
		if !t._built {
			t.build()
		}
		return t._hasChunk

	case *String, *Ref:
		return true

	case *Named:
		return ri.hasChunk(t.Underlying())

	default:
		panic(fmt.Sprintf("Todo: %T", t))
	}

	return false
}

func (ri *rtImp) alignof(t Type) int {
	if t.Kind() < BaseTypeNum {
		return ri.baseTypes[t.Kind()].size
	}

	switch t := t.(type) {
	case *Tuple:
		panic("Tuple is not supported")

	case *Struct:
		if !t._built {
			t.build()
		}
		return t._align

	case *String:
		if t.underlying == nil {
			t.underlying = ri.underlyingStruct(t)
		}
		return ri.alignof(t.underlying)

	case *Ref:
		if t.underlying == nil {
			t.underlying = ri.underlyingStruct(t)
		}
		return ri.alignof(t.underlying)

	case *Named:
		return ri.alignof(t.Underlying())

	default:
		panic(fmt.Sprintf("Todo: %T", t))
	}
}

func (ri *rtImp) sizeof(t Type) int {
	if t.Kind() < BaseTypeNum {
		return ri.baseTypes[t.Kind()].size
	}

	switch t := t.(type) {
	case *Tuple:
		panic("Tuple is not supported")

	case *Struct:
		if !t._built {
			t.build()
		}
		return t._size

	case *String:
		if t.underlying == nil {
			t.underlying = ri.underlyingStruct(t)
		}
		return ri.sizeof(t.underlying)

	case *Ref:
		if t.underlying == nil {
			t.underlying = ri.underlyingStruct(t)
		}
		return ri.sizeof(t.underlying)

	case *Named:
		return ri.sizeof(t.Underlying())

	default:
		panic(fmt.Sprintf("Todo: %T", t))
	}
}

func (ri *rtImp) underlyingStruct(t Type) (underlying *Struct) {
	underlying = &Struct{}
	switch t := t.(type) {
	case *Ref:
		c := StructMember{Name: "c", Type: &chunk{Base: t.Base}}
		d := StructMember{Name: "d", Type: &Ptr{Base: t.Base}}
		underlying.members = []StructMember{c, d}

	case *String:
		c := StructMember{Name: "c", Type: &chunk{Base: &U8{}}}
		d := StructMember{Name: "d", Type: &Ptr{Base: &U8{}}}
		l := StructMember{Name: "l", Type: &Uint{}}
		underlying.members = []StructMember{c, d, l}

	default:
		panic(fmt.Sprintf("Todo: %T", t))
	}

	underlying.build()
	return
}

func (ri *rtImp) initTank(t Type, kind RegisterKind) *tank {
	return ri.initTank_Offset(t, 0, kind)
}

func (ri *rtImp) initTank_Offset(t Type, offset int, kind RegisterKind) *tank {
	tank := &tank{}
	tank.typ = t
	tank.register.id = -1

	if t.Kind() < BaseTypeNum {
		if kind != KImv || t.Kind() != TypeKindChunk {
			tank.register.typ = ri.baseTypes[t.Kind()].regType
		} else {
			tank.register.typ = ri.baseTypes[TypeKindInt].regType
		}

		tank.register._offset = offset
		tank.register._kind = kind
		tank.register._otype = t
		return tank
	}

	switch t := t.(type) {
	case *Tuple:
		for _, m := range t.members {
			tm := ri.initTank_Offset(m, offset, kind)
			offset += ri.sizeof(m)
			tank.member = append(tank.member, tm)
		}

	case *Struct:
		for _, m := range t.members {
			tm := ri.initTank_Offset(m.Type, offset+m._start, kind)
			tank.member = append(tank.member, tm)
		}

	case *String:
		if t.underlying == nil {
			t.underlying = ri.underlyingStruct(t)
		}

		tank = ri.initTank_Offset(t.underlying, offset, kind)

	case *Ref:
		if t.underlying == nil {
			t.underlying = ri.underlyingStruct(t)
		}

		tank = ri.initTank_Offset(t.underlying, offset, kind)

	case *Named:
		tank = ri.initTank_Offset(t.underlying, offset, kind)

	default:
		panic(fmt.Sprintf("Todo: %T", t))
	}

	return tank
}

func (ri *rtImp) watType(t Type) wat.ValueType {
	if t.Kind() >= BaseTypeNum {
		panic(fmt.Sprintf("t is not a base type: %s", t.Name()))
	}
	return ri.baseTypes[t.Kind()].watType
}

func (ri *rtImp) BuildFunction(f *Function) {
	wf := &wat.Function{}
	wf.InternalName = f.InternalName
	wf.ExternalName = f.ExternalName

	paramRegNum := 0
	for _, param := range f.params {
		regs := param.Tank().raw()
		for _, reg := range regs {
			paramRegNum++
			param := wat.NewVar(reg.String(), ri.watType(reg.typ))
			wf.Params = append(wf.Params, param)
		}
	}

	for _, result := range f.results {
		tank := ri.initTank(result.Type(), KLocal)
		raw := tank.raw()
		for _, reg := range raw {
			wf.Results = append(wf.Results, ri.watType(reg.typ))
		}
	}

	for i := paramRegNum; i < len(f.commonRegs); i++ {
		r := f.commonRegs[i]
		wf.Locals = append(wf.Locals, wat.NewVar(r.String(), ri.watType(r.typ)))
	}
	for _, r := range f.chunkRegs {
		wf.Locals = append(wf.Locals, wat.NewVar(r.String(), ri.watType(r.typ)))
	}

	for _, stmt := range f.Body.Stmts {
		wf.Insts = append(wf.Insts, ri.emitStmt(stmt)...)
	}

	sb := strings.Builder{}
	wf.Format(&sb)
	println(sb.String())
}

func (ri *rtImp) emitStmt(stmt Stmt) (ret []wat.Inst) {
	switch s := stmt.(type) {
	case *Block:
		block := wat.NewInstBlock(s.Label)
		for _, s := range s.Stmts {
			block.Insts = append(block.Insts, ri.emitStmt(s)...)
		}
		ret = append(ret, wat.NewComment(s.Label))
		ret = append(ret, block)

	case *Alloc:
		ret = append(ret, wat.NewComment(s.String()))
		ret = append(ret, ri.emitAlloc(s)...)

	case *Imv:
		ret = append(ret, wat.NewComment(s.String()))
		ret = append(ret, ri.emitImv(s)...)

	case *Assign:
		ret = append(ret, wat.NewComment(s.String()))
		ret = append(ret, ri.emitAssign(s)...)

	case *Store:
		ret = append(ret, wat.NewComment(s.String()))
		ret = append(ret, ri.emitStore(s)...)

	case *If:
		ret = append(ret, wat.NewComment("if "+s.Cond.Name()))
		ret = append(ret, ri.emitIf(s)...)

	case *Loop:
		ret = append(ret, wat.NewComment("loop "+s.Label))
		ret = append(ret, ri.emitLoop(s)...)

	case *Discard:
		return

	case *Release:
		ret = append(ret, wat.NewComment(stmt.String()))
		ret = append(ret, ri.emitRelease_register(s.X)...)

	case *Return:
		ret = append(ret, wat.NewComment(stmt.String()))
		ret = append(ret, ri.emitReturn(s)...)

	case *NilCheck:
		ret = append(ret, wat.NewComment(stmt.String()))
		ret = append(ret, ri.emitNilCheck(s)...)

	default:
		panic(fmt.Sprintf("Todo: %T", stmt))
	}

	ret = append(ret, wat.NewBlank())
	return
}

func (ri *rtImp) emitAlloc(s *Alloc) (ret []wat.Inst) {
	switch s.Kind() {
	case AllocKindRegister:
		if s.init != nil {
			ret = append(ret, ri.emitExpr(s.init)...)
			rhRtained, rhConst := ri.exprRCState(s.init)
			ret = append(ret, ri.emitPop_tank(s.Tank(), true, !rhRtained && !rhConst)...)
		} else {
			ret = append(ret, ri.emitZero_tank(s.Tank(), true)...)
		}

	case AllocKindHeap:
		ss := fmt.Sprintf("%d", ri.sizeof(s.DataType()))
		ret = append(ret, wat.NewInstConst(ri.baseTypes[TypeKindInt].watType, "1")) // item_count
		ret = append(ret, wat.NewInstConst(ri.baseTypes[TypeKindInt].watType, "0")) // todo: free_method
		ret = append(ret, wat.NewInstConst(ri.baseTypes[TypeKindInt].watType, ss))  // item_size
		ret = append(ret, wat.NewInstCall("runtime.Block.HeapAlloc"))
		ret = append(ret, ri.emitPop_tank(s.Tank(), true, false)...)
	}

	return
}

func (ri *rtImp) emitImv(s *Imv) (ret []wat.Inst) {
	ret = append(ret, ri.emitExpr(s.val)...)
	ret = append(ret, ri.emitPop_tank(s.Tank(), false, false)...)
	return
}

func (ri *rtImp) exprRCState(expr Expr) (retained bool, isConst bool) {
	retained = expr.retained()
	_, isConst = expr.(*Const)
	return
}

func (ri *rtImp) emitAssign(s *Assign) (ret []wat.Inst) {
	for _, rh := range s.Rhs {
		ret = append(ret, ri.emitExpr(rh)...)
	}

	for i := len(s.Lhs) - 1; i >= 0; i-- {
		lh := s.Lhs[i]

		if _, ok := lh.(*Imv); ok {
			panic("Imv is not supported")
		}

		var rhRtained, rhConst bool
		if i < len(s.Rhs) {
			rhRtained, rhConst = ri.exprRCState(s.Rhs[i])
		} else {
			rhRtained, rhConst = ri.exprRCState(s.Rhs[0])
		}

		if lh == nil {
			ret = append(ret, ri.emitDrop(s.rhsType[i], rhRtained && !rhConst)...)
		} else {
			ret = append(ret, ri.emitPop_tank(lh.Tank(), true, !rhRtained && !rhConst)...)
		}
	}

	return
}

func (ri *rtImp) emitStore(s *Store) (ret []wat.Inst) {
	var ptr register
	switch s.Loc.Type().(type) {
	case *Ref:
		ptr = s.Loc.Tank().member[1].register

	case *Ptr:
		ptr = s.Loc.Tank().register

	default:
		panic(fmt.Sprintf("Todo: %T", s.Loc.Type()))
	}

	retained, isConst := ri.exprRCState(s.Val)
	switch val := s.Val.(type) {
	case Var:
		regs := val.Tank().raw()
		for _, reg := range regs {
			ret = append(ret, ri.emitStore_register(ptr, reg._offset, reg, !retained && !isConst)...)
		}

	case *Const:
		ret = append(ret, ri.emitStore_const(ptr, 0, val)...)

	default:
		panic(fmt.Sprintf("Todo: %T", val))
	}

	return
}

func (ri *rtImp) emitIf(s *If) (ret []wat.Inst) {
	ret = append(ret, ri.emitExpr(s.Cond)...)

	i := wat.NewInstIf(nil, nil, nil)
	for _, s := range s.True.Stmts {
		i.True = append(i.True, ri.emitStmt(s)...)
	}
	for _, s := range s.False.Stmts {
		i.False = append(i.False, ri.emitStmt(s)...)
	}

	ret = append(ret, i)
	return
}

func (ri *rtImp) emitLoop(s *Loop) (ret []wat.Inst) {
	loopLabel := s.Label + ".loop"
	loop := wat.NewInstLoop(loopLabel)

	for _, preCond := range s.PreCond {
		loop.Insts = append(loop.Insts, ri.emitStmt(preCond)...)
	}
	loop.Insts = append(loop.Insts, wat.NewComment(s.Cond.Name()))
	loop.Insts = append(loop.Insts, ri.emitExpr(s.Cond)...)

	body := wat.NewInstBlock(s.Label + ".body")
	for _, stmt := range s.Body.Stmts {
		body.Insts = append(body.Insts, ri.emitStmt(stmt)...)
	}

	_if := wat.NewInstIf(nil, nil, nil)
	_if.Name = s.Label + ".do"
	_if.True = append(_if.True, body)
	for _, post := range s.Post.Stmts {
		_if.True = append(_if.True, ri.emitStmt(post)...)
	}
	_if.True = append(_if.True, wat.NewInstBr(loopLabel))

	loop.Insts = append(loop.Insts, _if)
	ret = append(ret, loop)
	return
}

func (ri *rtImp) emitStore_register(ptr register, offset int, value register, needRetain bool) (ret []wat.Inst) {
	hasChunk := ri.hasChunk(value._otype)

	ret = append(ret, ri.emitPush_register(ptr))   // p
	ret = append(ret, ri.emitPush_register(value)) // p v

	if needRetain && hasChunk {
		ret = append(ret, ri.emitRetain())
	}

	if hasChunk {
		ret = append(ret, ri.emitLoad_BaseType(ptr, offset, value.typ)...) // p v o
		ret = append(ret, ri.emitRelease_stack())                          // p v
	}

	switch value.typ.Kind() {
	case TypeKindBool, TypeKindU8, TypeKindI8:
		ret = append(ret, wat.NewInstStore8(offset, 1))

	case TypeKindI16, TypeKindU16:
		ret = append(ret, wat.NewInstStore16(offset, 2))

	default:
		ret = append(ret, wat.NewInstStore(ri.watType(value.typ), offset, ri.sizeof(value.typ)))
	}

	return ret
}

func (ri *rtImp) emitStore_const(ptr register, offset int, value *Const) (ret []wat.Inst) {
	if value.Type().Kind() >= BaseTypeNum {
		panic(fmt.Sprintf("Todo: %T", value.Type()))
	}

	ret = append(ret, ri.emitPush_register(ptr))
	ret = append(ret, wat.NewInstConst(ri.watType(value.Type()), value.Name()))

	switch value.Type().Kind() {
	case TypeKindBool, TypeKindU8, TypeKindI8:
		ret = append(ret, wat.NewInstStore8(offset, 1))

	case TypeKindI16, TypeKindU16:
		ret = append(ret, wat.NewInstStore16(offset, 2))

	default:
		ret = append(ret, wat.NewInstStore(ri.watType(value.Type()), offset, ri.sizeof(value.Type())))
	}

	return ret
}

func (ri *rtImp) emitReturn(s *Return) (ret []wat.Inst) {
	for _, expr := range s.Results {
		ret = append(ret, ri.emitExpr(expr)...)
	}
	ret = append(ret, wat.NewInstReturn())
	return
}

func (ri *rtImp) emitPush_tank(tank *tank) (ret []wat.Inst) {
	regs := tank.raw()
	for _, reg := range regs {
		ret = append(ret, ri.emitPush_register(reg))
	}
	return
}

func (ri *rtImp) emitPush_register(reg register) wat.Inst {
	if reg._kind == KGlobal {
		return wat.NewInstGetGlobal(reg.String())
	} else {
		return wat.NewInstGetLocal(reg.String())
	}
}

func (ri *rtImp) emitPop_tank(tank *tank, needRelease bool, needRetain bool) (ret []wat.Inst) {
	regs := tank.raw()
	for i := len(regs) - 1; i >= 0; i-- {
		ret = append(ret, ri.emitPop_register(regs[i], needRelease, needRetain)...)
	}
	return
}

func (ri *rtImp) emitPop_register(reg register, needRelease bool, needRetain bool) (ret []wat.Inst) {
	hasChunk := ri.hasChunk(reg._otype)

	if needRetain && hasChunk {
		ret = append(ret, ri.emitRetain())
	}

	if needRelease && hasChunk {
		ret = append(ret, ri.emitRelease_register(reg)...)
	}

	if reg._kind == KGlobal {
		ret = append(ret, wat.NewInstSetGlobal(reg.String()))
	} else {
		ret = append(ret, wat.NewInstSetLocal(reg.String()))
	}
	return
}

func (ri *rtImp) emitDrop(t Type, needRelease bool) (ret []wat.Inst) {
	raw := ri.initTank(t, KLocal).raw()
	for i := len(raw) - 1; i >= 0; i-- {
		ret = append(ret, ri.emitDrop_BaseType(raw[i].typ, needRelease))
	}
	return
}

func (ri *rtImp) emitDrop_BaseType(t Type, needRelease bool) (ret wat.Inst) {
	if t.Kind() >= BaseTypeNum {
		panic("t is not a base type")
	}
	if needRelease && ri.hasChunk(t) {
		ret = wat.NewInstCall("runtime.Block.Release")
	} else {
		ret = wat.NewInstDrop()
	}
	return
}

func (ri *rtImp) emitZero_tank(tank *tank, needRelease bool) (ret []wat.Inst) {
	regs := tank.raw()
	for i := len(regs) - 1; i >= 0; i-- {
		ret = append(ret, ri.emitZero_register(regs[i], needRelease)...)
	}
	return
}

func (ri *rtImp) emitZero_register(reg register, needRelease bool) (ret []wat.Inst) {
	ret = append(ret, wat.NewInstConst(ri.watType(reg.typ), "0"))
	ret = append(ret, ri.emitPop_register(reg, needRelease, false)...)
	return
}

func (ri *rtImp) emitRetain() wat.Inst {
	return wat.NewInstCall("runtime.Block.Retain")
}

func (ri *rtImp) emitRelease_register(reg register) (ret []wat.Inst) {
	ret = make([]wat.Inst, 2)
	ret[0] = ri.emitPush_register(reg)
	ret[1] = ri.emitRelease_stack()
	return
}

func (ri *rtImp) emitRelease_stack() wat.Inst {
	return wat.NewInstCall("runtime.Block.Release")
}

func (ri *rtImp) emitLoad_BaseType(ptr register, offset int, t Type) (ret []wat.Inst) {
	if t.Kind() >= BaseTypeNum {
		panic("t is not a base type")
	}

	ret = make([]wat.Inst, 2)
	ret[0] = ri.emitPush_register(ptr)

	switch t.Kind() {
	case TypeKindBool, TypeKindU8:
		ret[1] = wat.NewInstLoad8u(offset, 1)
	case TypeKindI8:
		ret[1] = wat.NewInstLoad8s(offset, 1)
	case TypeKindI16:
		ret[1] = wat.NewInstLoad16s(offset, 2)
	case TypeKindU16:
		ret[1] = wat.NewInstLoad16u(offset, 2)

	default:
		ret[1] = wat.NewInstLoad(ri.watType(t), offset, ri.sizeof(t))
	}
	return
}

func (ri *rtImp) emitExpr(expr Expr) (ret []wat.Inst) {
	if c, ok := expr.(*Const); ok {
		if expr.Type().Kind() < BaseTypeNum {
			ret = append(ret, wat.NewInstConst(ri.watType(expr.Type()), c.Name()))
		} else {
			panic(fmt.Sprintf("Todo: %T", expr))
		}
		return
	}

	switch e := expr.(type) {
	case Var:
		ret = append(ret, ri.emitPush_tank(e.Tank())...)

	case *Call:
		var call_common *CallCommon
		switch callee := e.Callee.(type) {
		case *BuiltinCall:
			call_common = &callee.CallCommon
		case *StaticCall:
			call_common = &callee.CallCommon
		default:
			panic(fmt.Sprintf("Todo: %T", callee))
		}

		for _, arg := range call_common.Args {
			ret = append(ret, ri.emitExpr(arg)...)
		}
		ret = append(ret, wat.NewInstCall(call_common.String()))

	case *Load:
		ret = append(ret, ri.emitLoad(e)...)

	case *MemberAddr:
		ret = append(ret, ri.emitMemberAddr(e)...)

	case *Biop:
		ret = append(ret, ri.emitBiop(e)...)

	case *Unop:
		ret = append(ret, ri.emitUnop(e)...)

	default:
		panic(fmt.Sprintf("Todo: %T", e))
	}

	return
}

func (ri *rtImp) emitNilCheck(e *NilCheck) (ret []wat.Inst) {
	var ptr register
	switch e.X.Type().(type) {
	case *Ref:
		ptr = e.X.Tank().member[1].register
	case *Ptr:
		ptr = e.X.Tank().register
	default:
		panic("X is not a Ref or Ptr")
	}
	ret = append(ret, ri.emitPush_register(ptr))
	ret = append(ret, wat.NewInstCall("runtime.Block.NilCheck"))

	return
}

func (ri *rtImp) emitLoad(e *Load) (ret []wat.Inst) {
	loc, ok := e.Loc.(Var)
	if !ok {
		panic("Loc is not a Var")
	}

	var ptr register
	switch e.Loc.Type().(type) {
	case *Ref:
		ptr = loc.Tank().member[1].register

	case *Ptr:
		ptr = loc.Tank().register

	default:
		panic("Loc is not a Ref or Ptr")
	}

	regs := ri.initTank(e.Type(), KLocal).raw()
	for _, reg := range regs {
		load := ri.emitLoad_BaseType(ptr, reg._offset, reg.typ)
		ret = append(ret, load...)
	}
	return
}

func (ri *rtImp) emitMemberAddr(e *MemberAddr) (ret []wat.Inst) {
	x := e.X
	offset := e.member._start

	for sm, ok := x.(*MemberAddr); ok; sm, ok = x.(*MemberAddr) {
		x = sm.X
		offset += sm.member._start
	}

	ret = append(ret, ri.emitExpr(x)...)
	if offset > 0 {
		ret = append(ret, wat.NewInstConst(ri.baseTypes[TypeKindPtr].watType, strconv.Itoa(offset)))
		ret = append(ret, wat.NewInstAdd(ri.baseTypes[TypeKindPtr].watType))
	}
	return
}

func (ri *rtImp) emitBiop(e *Biop) (ret []wat.Inst) {
	switch e.Op {
	case XOR:
		ret = append(ret, ri.emitExpr(e.X)...)
		ret = append(ret, ri.emitExpr(e.Y)...)
		ret = append(ret, wat.NewInstXor(ri.watType(e.X.Type())))

	case LAND:
		ret = append(ret, ri.emitExpr(e.X)...)

		i := wat.NewInstIf(nil, nil, nil)
		i.Ret = []wat.ValueType{ri.watType(&Bool{})}
		i.True = ri.emitExpr(e.Y)
		i.False = []wat.Inst{wat.NewInstConst(ri.watType(&Bool{}), "0")}

		ret = append(ret, i)

	case LOR:
		ret = append(ret, ri.emitExpr(e.X)...)

		i := wat.NewInstIf(nil, nil, nil)
		i.Ret = []wat.ValueType{ri.watType(&Bool{})}
		i.True = []wat.Inst{wat.NewInstConst(ri.watType(&Bool{}), "1")}
		i.False = ri.emitExpr(e.Y)

		ret = append(ret, i)

	case SHL:
		ret = append(ret, ri.emitExpr(e.X)...)
		ret = append(ret, ri.emitExpr(e.Y)...)
		if ri.sizeof(e.X.Type()) <= 4 && ri.sizeof(e.Y.Type()) == 8 {
			ret = append(ret, wat.NewInstConvert_i32_wrap_i64())
		} else if ri.sizeof(e.X.Type()) == 8 && ri.sizeof(e.Y.Type()) <= 4 {
			ret = append(ret, wat.NewInstConvert_i64_extend_i32_u())
		}
		ret = append(ret, wat.NewInstShl(ri.watType(e.X.Type())))

	case SHR:
		ret = append(ret, ri.emitExpr(e.X)...)
		ret = append(ret, ri.emitExpr(e.Y)...)
		if ri.sizeof(e.X.Type()) <= 4 && ri.sizeof(e.Y.Type()) == 8 {
			ret = append(ret, wat.NewInstConvert_i32_wrap_i64())
		} else if ri.sizeof(e.X.Type()) == 8 && ri.sizeof(e.Y.Type()) <= 4 {
			ret = append(ret, wat.NewInstConvert_i64_extend_i32_u())
		}
		ret = append(ret, wat.NewInstShr(ri.watType(e.X.Type())))

	case ADD:
		ret = append(ret, ri.emitExpr(e.X)...)
		ret = append(ret, ri.emitExpr(e.Y)...)

		switch e.X.Type().Kind() {
		case TypeKindString, TypeKindComplex64, TypeKindComplex128:
			panic("Todo")

		default:
			ret = append(ret, wat.NewInstAdd(ri.watType(e.X.Type())))
		}

	case SUB:
		ret = append(ret, ri.emitExpr(e.X)...)
		ret = append(ret, ri.emitExpr(e.Y)...)

		switch e.X.Type().Kind() {
		case TypeKindComplex64, TypeKindComplex128:
			panic("Todo")

		default:
			ret = append(ret, wat.NewInstSub(ri.watType(e.X.Type())))
		}

	case MUL:
		ret = append(ret, ri.emitExpr(e.X)...)
		ret = append(ret, ri.emitExpr(e.Y)...)

		switch e.X.Type().Kind() {
		case TypeKindComplex64, TypeKindComplex128:
			panic("Todo")

		default:
			ret = append(ret, wat.NewInstMul(ri.watType(e.X.Type())))
		}

	case QUO:
		switch e.X.Type().Kind() {
		case TypeKindComplex64, TypeKindComplex128:
			panic("Todo")

		case TypeKindI8:
			ret = append(ret, ri.emitExpr(e.X)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))
			ret = append(ret, ri.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))

			ret = append(ret, wat.NewInstDiv(wat.I32{}))

		case TypeKindI16:
			ret = append(ret, ri.emitExpr(e.X)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))
			ret = append(ret, ri.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))

			ret = append(ret, wat.NewInstDiv(wat.I32{}))

		default:
			ret = append(ret, ri.emitExpr(e.X)...)
			ret = append(ret, ri.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstDiv(ri.watType(e.X.Type())))
		}

	case REM:
		switch e.X.Type().Kind() {
		case TypeKindComplex64, TypeKindComplex128:
			panic("Todo")

		case TypeKindI8:
			ret = append(ret, ri.emitExpr(e.X)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))
			ret = append(ret, ri.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))

			ret = append(ret, wat.NewInstRem(wat.I32{}))

		case TypeKindI16:
			ret = append(ret, ri.emitExpr(e.X)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))
			ret = append(ret, ri.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShl(wat.I32{}))
			ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
			ret = append(ret, wat.NewInstShr(wat.I32{}))

			ret = append(ret, wat.NewInstRem(wat.I32{}))

		default:
			ret = append(ret, ri.emitExpr(e.X)...)
			ret = append(ret, ri.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstRem(ri.watType(e.X.Type())))
		}

	case AND:
		ret = append(ret, ri.emitExpr(e.X)...)
		ret = append(ret, ri.emitExpr(e.Y)...)
		ret = append(ret, wat.NewInstAnd(ri.watType(e.X.Type())))

	case ANDNOT:
		t := ri.watType(e.X.Type())
		ret = append(ret, ri.emitExpr(e.X)...)
		ret = append(ret, ri.emitExpr(e.Y)...)
		ret = append(ret, wat.NewInstConst(t, "-1"))
		ret = append(ret, wat.NewInstXor(t))
		ret = append(ret, wat.NewInstAnd(t))

	case OR:
		ret = append(ret, ri.emitExpr(e.X)...)
		ret = append(ret, ri.emitExpr(e.Y)...)
		ret = append(ret, wat.NewInstOr(ri.watType(e.X.Type())))

	case EQL:
		ret = append(ret, ri.emitEql(e.X, e.Y)...)

	case NEQ:
		ret = append(ret, ri.emitEql(e.X, e.Y)...)
		ret = append(ret, wat.NewInstEqz(wat.I32{}))

	case GTR:
		ret = append(ret, ri.emitExpr(e.X)...)
		ret = append(ret, ri.emitExpr(e.Y)...)
		if e.X.Type().Kind() == TypeKindString {
			ret = append(ret, wat.NewInstCall("$wa.runtime.string_GTR"))
		} else {
			ret = append(ret, wat.NewInstGt(ri.watType(e.X.Type())))
		}

	case LSS:
		ret = append(ret, ri.emitExpr(e.X)...)
		ret = append(ret, ri.emitExpr(e.Y)...)
		if e.X.Type().Kind() == TypeKindString {
			ret = append(ret, wat.NewInstCall("$wa.runtime.string_LSS"))
		} else {
			ret = append(ret, wat.NewInstLt(ri.watType(e.X.Type())))
		}

	case GEQ:
		ret = append(ret, ri.emitExpr(e.X)...)
		ret = append(ret, ri.emitExpr(e.Y)...)
		if e.X.Type().Kind() == TypeKindString {
			ret = append(ret, wat.NewInstCall("$wa.runtime.string_GEQ"))
		} else {
			ret = append(ret, wat.NewInstGe(ri.watType(e.X.Type())))
		}

	case LEQ:
		ret = append(ret, ri.emitExpr(e.X)...)
		ret = append(ret, ri.emitExpr(e.Y)...)
		if e.X.Type().Kind() == TypeKindString {
			ret = append(ret, wat.NewInstCall("$wa.runtime.string_LEQ"))
		} else {
			ret = append(ret, wat.NewInstLe(ri.watType(e.X.Type())))
		}

	case LEG:
		if e.X.Type().Kind() == TypeKindString {
			ret = append(ret, ri.emitExpr(e.X)...)
			ret = append(ret, ri.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstCall("$wa.runtime.string_Comp"))
		} else {
			t := ri.watType(e.X.Type())
			ret = append(ret, ri.emitExpr(e.X)...)
			ret = append(ret, ri.emitExpr(e.Y)...)
			ret = append(ret, wat.NewInstLt(t))

			inst_lt := wat.NewInstIf(nil, nil, nil)
			inst_lt.Ret = append(inst_lt.Ret, wat.I32{})
			inst_lt.True = append(inst_lt.True, wat.NewInstConst(wat.I32{}, "-1"))
			inst_lt.False = append(inst_lt.False, ri.emitExpr(e.X)...)
			inst_lt.False = append(inst_lt.False, ri.emitExpr(e.Y)...)
			inst_lt.False = append(inst_lt.False, wat.NewInstGt(t))

			ret = append(ret, inst_lt)
		}

	default:
		panic(fmt.Sprintf("Todo: %s", e.Op.String()))
	}

	switch e.Type().Kind() {
	case TypeKindU8:
		ret = append(ret, wat.NewInstConst(ri.watType(e.X.Type()), "255"))
		ret = append(ret, wat.NewInstAnd(ri.watType(e.X.Type())))
	case TypeKindU16:
		ret = append(ret, wat.NewInstConst(ri.watType(e.X.Type()), "65535"))
		ret = append(ret, wat.NewInstAnd(ri.watType(e.X.Type())))
	case TypeKindI8:
		ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
		ret = append(ret, wat.NewInstShl(wat.I32{}))
		ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
		ret = append(ret, wat.NewInstShr(wat.I32{}))
	case TypeKindI16:
		ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
		ret = append(ret, wat.NewInstShl(wat.I32{}))
		ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
		ret = append(ret, wat.NewInstShr(wat.I32{}))
	}

	return
}

func (ri *rtImp) emitEql(x, y Expr) (ret []wat.Inst) {
	if x.Type().Kind() < BaseTypeNum {
		ret = append(ret, ri.emitExpr(x)...)
		ret = append(ret, ri.emitExpr(y)...)
		ret = append(ret, wat.NewInstEq(ri.watType(x.Type())))
	} else {
		panic("Todo")
	}

	return
}

func (ri *rtImp) emitUnop(e *Unop) (ret []wat.Inst) {
	switch e.Op {
	case NOT:
		ret = append(ret, ri.emitExpr(e.X)...)
		ret = append(ret, wat.NewInstEqz(wat.I32{}))

	case NEG:
		switch e.X.Type().Kind() {
		case TypeKindF32, TypeKindF64:
			ret = append(ret, ri.emitExpr(e.X)...)
			ret = append(ret, wat.NewInstNeg(ri.watType(e.X.Type())))
		default:
			ret = append(ret, wat.NewInstConst(ri.watType(e.X.Type()), "0"))
			ret = append(ret, ri.emitExpr(e.X)...)
			ret = append(ret, wat.NewInstSub(ri.watType(e.X.Type())))
		}

	case XOR:
		ret = append(ret, wat.NewInstConst(ri.watType(e.X.Type()), "-1"))
		ret = append(ret, ri.emitExpr(e.X)...)
		ret = append(ret, wat.NewInstXor(ri.watType(e.X.Type())))

	default:
		panic(fmt.Sprintf("Todo: %s", e.Op.String()))
	}

	switch e.Type().Kind() {
	case TypeKindU8:
		ret = append(ret, wat.NewInstConst(ri.watType(e.X.Type()), "255"))
		ret = append(ret, wat.NewInstAnd(ri.watType(e.X.Type())))
	case TypeKindU16:
		ret = append(ret, wat.NewInstConst(ri.watType(e.X.Type()), "65535"))
		ret = append(ret, wat.NewInstAnd(ri.watType(e.X.Type())))
	case TypeKindI8:
		ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
		ret = append(ret, wat.NewInstShl(wat.I32{}))
		ret = append(ret, wat.NewInstConst(wat.I32{}, "24"))
		ret = append(ret, wat.NewInstShr(wat.I32{}))
	case TypeKindI16:
		ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
		ret = append(ret, wat.NewInstShl(wat.I32{}))
		ret = append(ret, wat.NewInstConst(wat.I32{}, "16"))
		ret = append(ret, wat.NewInstShr(wat.I32{}))
	}

	return ret
}
