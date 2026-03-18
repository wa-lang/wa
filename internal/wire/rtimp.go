// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
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
	ri.baseTypes[TypeKindU8] = baseType{regType: &I32{}, watType: wat.I32{}, hasChunk: false, size: 1}
	ri.baseTypes[TypeKindU16] = baseType{regType: &I32{}, watType: wat.I32{}, hasChunk: false, size: 2}
	ri.baseTypes[TypeKindU32] = baseType{regType: &I32{}, watType: wat.I32{}, hasChunk: false, size: 4}
	ri.baseTypes[TypeKindU64] = baseType{regType: &I64{}, watType: wat.I64{}, hasChunk: false, size: 8}
	ri.baseTypes[TypeKindI8] = baseType{regType: &I32{}, watType: wat.I32{}, hasChunk: false, size: 1}
	ri.baseTypes[TypeKindI16] = baseType{regType: &I32{}, watType: wat.I32{}, hasChunk: false, size: 2}
	ri.baseTypes[TypeKindI32] = baseType{regType: &I32{}, watType: wat.I32{}, hasChunk: false, size: 4}
	ri.baseTypes[TypeKindI64] = baseType{regType: &I64{}, watType: wat.I64{}, hasChunk: false, size: 8}
	ri.baseTypes[TypeKindInt] = baseType{regType: &I32{}, watType: wat.I32{}, hasChunk: false, size: 4}
	ri.baseTypes[TypeKindUint] = baseType{regType: &I32{}, watType: wat.I32{}, hasChunk: false, size: 4}
	ri.baseTypes[TypeKindF32] = baseType{regType: &F32{}, watType: wat.F32{}, hasChunk: false, size: 4}
	ri.baseTypes[TypeKindF64] = baseType{regType: &F64{}, watType: wat.F64{}, hasChunk: false, size: 8}
	ri.baseTypes[TypeKindRune] = baseType{regType: &I32{}, watType: wat.I32{}, hasChunk: false, size: 4}
	ri.baseTypes[TypeKindPtr] = baseType{regType: &I32{}, watType: wat.I32{}, hasChunk: false, size: 4}
	ri.baseTypes[TypeKindChunk] = baseType{regType: &chunk{Base: &Void{}}, watType: wat.I32{}, hasChunk: true, size: 4}
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
		return ri.sizeof(t.Underlying())

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
		ret = append(ret, block)

	case *Alloc:
		return ri.emitAlloc(s)

	case *Assign:
		return ri.emitAssign(s)

	case *Store:
		return ri.emitStore(s)

	case *Discard:
		return

	case *Release:
		return ri.emitRelease_register(s.X)

	case *Return:
		return ri.emitReturn(s)

	default:
		panic(fmt.Sprintf("Todo: %T", stmt))
	}

	return
}

func (ri *rtImp) emitAlloc(s *Alloc) (ret []wat.Inst) {
	switch s.Kind() {
	case AllocKindRegister:
		if s.init != nil {
			ret = append(ret, ri.emitExpr(s.init)...)
			rhRtained, rhConst := ri.exprRCState(s.init)
			ret = append(ret, ri.emitPop_tank(s.Tank(), !rhRtained && !rhConst)...)
		} else {
			ret = append(ret, ri.emitZero_tank(s.Tank())...)
		}

	case AllocKindHeap:
		ss := fmt.Sprintf("%d", ri.sizeof(s.DataType()))
		ret = append(ret, wat.NewInstConst(ri.baseTypes[TypeKindInt].watType, "1")) // item_count
		ret = append(ret, wat.NewInstConst(ri.baseTypes[TypeKindInt].watType, "0")) // todo: free_method
		ret = append(ret, wat.NewInstConst(ri.baseTypes[TypeKindInt].watType, ss))  // item_size
		ret = append(ret, wat.NewInstCall("runtime.Block.HeapAlloc"))
		ret = append(ret, ri.emitPop_tank(s.Tank(), false)...)
	}

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

		var rhRtained, rhConst bool
		if i < len(s.Rhs) {
			rhRtained, rhConst = ri.exprRCState(s.Rhs[i])
		} else {
			rhRtained, rhConst = ri.exprRCState(s.Rhs[0])
		}

		if lh == nil {
			ret = append(ret, ri.emitDrop(s.rhsType[i], rhRtained && !rhConst)...)
		} else {
			ret = append(ret, ri.emitPop_tank(lh.Tank(), !rhRtained && !rhConst)...)
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

	default:
		panic(fmt.Sprintf("Todo: %T", val))
	}

	return
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
	return wat.NewInstGetLocal(reg.String())
}

func (ri *rtImp) emitPop_tank(tank *tank, needRetain bool) (ret []wat.Inst) {
	regs := tank.raw()
	for i := len(regs) - 1; i >= 0; i-- {
		hasChunk := ri.hasChunk(regs[i].typ)
		ret = append(ret, ri.emitPop_register(regs[i], needRetain && hasChunk)...)
	}
	return
}

func (ri *rtImp) emitPop_register(reg register, needRetain bool) (ret []wat.Inst) {
	hasChunk := ri.hasChunk(reg.typ)

	if needRetain {
		if !hasChunk {
			panic("needRetain but no chunk")
		}
		ret = append(ret, ri.emitRetain())
	}

	if hasChunk {
		ret = append(ret, ri.emitRelease_register(reg)...)
	}

	ret = append(ret, wat.NewInstSetLocal(reg.String()))
	return
}

func (ri *rtImp) emitDrop(t Type, needRelease bool) (ret []wat.Inst) {
	raw := ri.initTank(t, KLocal).raw()
	for i := len(raw) - 1; i >= 0; i-- {
		rt := raw[i].typ
		hasChunk := ri.hasChunk(rt)
		ret = append(ret, ri.emitDrop_BaseType(rt, needRelease && hasChunk)...)
	}
	return
}

func (ri *rtImp) emitDrop_BaseType(t Type, needRelease bool) (ret []wat.Inst) {
	if t.Kind() >= BaseTypeNum {
		panic("t is not a base type")
	}
	if needRelease {
		if !ri.hasChunk(t) {
			panic("needRelease but no chunk")
		}
		ret = append(ret, wat.NewInstCall("runtime.Block.Release"))
	} else {
		drop := wat.NewInstDrop()
		ret = append(ret, drop)
	}
	return
}

func (ri *rtImp) emitZero_tank(tank *tank) (ret []wat.Inst) {
	regs := tank.raw()
	for i := len(regs) - 1; i >= 0; i-- {
		ret = append(ret, ri.emitZero_register(regs[i])...)
	}
	return
}

func (ri *rtImp) emitZero_register(reg register) (ret []wat.Inst) {
	ret = append(ret, wat.NewInstConst(ri.watType(reg.typ), "0"))
	ret = append(ret, ri.emitPop_register(reg, false)...)
	return
}

func (ri *rtImp) emitRetain() wat.Inst {
	return wat.NewInstCall("runtime.Block.Retain")
}

func (ri *rtImp) emitRelease_register(reg register) (ret []wat.Inst) {
	ret = make([]wat.Inst, 2)
	ret[0] = wat.NewInstGetLocal(reg.String())
	ret[1] = ri.emitRelease_stack()
	return
}

func (ri *rtImp) emitRelease_stack() wat.Inst {
	return wat.NewInstCall("runtime.Block.Release")
}

func (ri *rtImp) emitStore_register(ptr register, offset int, value register, needRetain bool) (ret []wat.Inst) {
	hasChunk := ri.hasChunk(value.typ)

	ret = append(ret, ri.emitPush_register(ptr))   // p
	ret = append(ret, ri.emitPush_register(value)) // p v

	if needRetain {
		if !hasChunk {
			panic("needRetain but no chunk")
		}
		ret = append(ret, ri.emitRetain())
	}

	if hasChunk {
		ret = append(ret, ri.emitLoad_BaseType(ptr, offset, value.typ)...) // p v o
		ret = append(ret, ri.emitRelease_stack())                          // p v
	}

	ret = append(ret, wat.NewInstStore(ri.watType(value.typ), offset, ri.sizeof(value.typ)))
	return ret
}

func (ri *rtImp) emitLoad_BaseType(ptr register, offset int, t Type) (ret []wat.Inst) {
	if t.Kind() >= BaseTypeNum {
		panic("t is not a base type")
	}

	ret = make([]wat.Inst, 2)
	ret[0] = wat.NewInstGetLocal(ptr.String())
	ret[1] = wat.NewInstLoad(ri.watType(t), offset, ri.sizeof(t))
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

	default:
		panic(fmt.Sprintf("Todo: %T", e))
	}

	return
}
