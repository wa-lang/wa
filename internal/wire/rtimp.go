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
	initTank(t Type, isImv bool) *tank
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

func (ri *rtImp) initTank(t Type, isImv bool) *tank {
	tank := &tank{}
	tank.typ = t
	tank.register.id = -1

	if t.Kind() < BaseTypeNum {
		if !isImv || t.Kind() != TypeKindChunk {
			tank.register.typ = ri.baseTypes[t.Kind()].regType
		} else {
			tank.register.typ = ri.baseTypes[TypeKindInt].regType
		}

		return tank
	}

	switch t := t.(type) {
	case *Tuple:
		for _, m := range t.members {
			tm := ri.initTank(m, isImv)
			tank.member = append(tank.member, tm)
		}

	case *Struct:
		for _, m := range t.members {
			tm := ri.initTank(m.Type, isImv)
			tank.member = append(tank.member, tm)
		}

	case *String:
		if t.underlying == nil {
			t.underlying = ri.underlyingStruct(t)
		}

		tank = ri.initTank(t.underlying, isImv)

	case *Ref:
		if t.underlying == nil {
			t.underlying = ri.underlyingStruct(t)
		}

		tank = ri.initTank(t.underlying, isImv)

	case *Named:
		tank = ri.initTank(t.underlying, isImv)

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
		tank := ri.initTank(result.Type(), false)
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
			ret = append(ret, ri.emitPop_tank(s.Tank(), s.init.retained())...)
		} else {
			ret = append(ret, ri.emitZero_tank(s.Tank())...)
		}

	case AllocKindHeap:
		ss := fmt.Sprintf("%d", ri.sizeof(s.DataType()))
		ret = append(ret, wat.NewInstConst(ri.baseTypes[TypeKindInt].watType, "1")) // item_count
		ret = append(ret, wat.NewInstConst(ri.baseTypes[TypeKindInt].watType, "0")) // todo: free_method
		ret = append(ret, wat.NewInstConst(ri.baseTypes[TypeKindInt].watType, ss))  // item_size
		ret = append(ret, wat.NewInstCall("runtime.Block.HeapAlloc"))
		ret = append(ret, ri.emitPop_tank(s.Tank(), true)...)
	}

	return
}

func (ri *rtImp) emitPop_tank(tank *tank, retained bool) (ret []wat.Inst) {
	regs := tank.raw()
	for i := len(regs) - 1; i >= 0; i-- {
		ret = append(ret, ri.emitPop_register(regs[i], retained)...)
	}
	return
}

func (ri *rtImp) emitPop_register(reg register, retained bool) (ret []wat.Inst) {
	if ri.hasChunk(reg.typ) {
		if !retained {
			ret = append(ret, ri.emitRetain()...)
		}
		ret = append(ret, ri.emitRelease_register(reg)...)
	}
	ret = append(ret, wat.NewInstSetLocal(reg.String()))
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
	ret = append(ret, ri.emitPop_register(reg, true)...)
	return
}

func (ri *rtImp) emitRetain() (ret []wat.Inst) {
	ret = make([]wat.Inst, 1)
	ret[0] = wat.NewInstCall("runtime.Block.Retain")
	return
}

func (ri *rtImp) emitRelease_register(reg register) (ret []wat.Inst) {
	ret = make([]wat.Inst, 2)
	ret[0] = wat.NewInstGetLocal(reg.String())
	ret[1] = wat.NewInstCall("runtime.Block.Release")
	return
}

func (ri *rtImp) emitExpr(expr Expr) (ret []wat.Inst) {

	panic("Todo")

	return
}
