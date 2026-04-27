// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
)

type RtImp interface {
	underlyingStruct(t Type) *Struct
	initTank(t Type, kind RegisterKind) *Tank
}

var rtimp RtImp

func init() {
	imp := &rtImp{}
	imp.init()
	rtimp = imp
}

type baseType struct {
	regType Type
}

type rtImp struct {
	baseTypes []baseType
}

func (ri *rtImp) init() {
	ri.baseTypes = make([]baseType, BaseTypeNum)

	ri.baseTypes[TypeKindUnknown] = baseType{regType: nil}
	ri.baseTypes[TypeKindVoid] = baseType{regType: &Void{}}
	ri.baseTypes[TypeKindBool] = baseType{regType: &I32{}}
	ri.baseTypes[TypeKindU8] = baseType{regType: &I32{}}
	ri.baseTypes[TypeKindU16] = baseType{regType: &I32{}}
	ri.baseTypes[TypeKindU32] = baseType{regType: &I32{}}
	ri.baseTypes[TypeKindU64] = baseType{regType: &I64{}}
	ri.baseTypes[TypeKindI8] = baseType{regType: &I32{}}
	ri.baseTypes[TypeKindI16] = baseType{regType: &I32{}}
	ri.baseTypes[TypeKindI32] = baseType{regType: &I32{}}
	ri.baseTypes[TypeKindI64] = baseType{regType: &I64{}}
	ri.baseTypes[TypeKindInt] = baseType{regType: &I32{}}
	ri.baseTypes[TypeKindUint] = baseType{regType: &I32{}}
	ri.baseTypes[TypeKindF32] = baseType{regType: &F32{}}
	ri.baseTypes[TypeKindF64] = baseType{regType: &F64{}}
	ri.baseTypes[TypeKindRune] = baseType{regType: &I32{}}
	ri.baseTypes[TypeKindPtr] = baseType{regType: &I32{}}
	ri.baseTypes[TypeKindChunk] = baseType{regType: &chunk{Base: &Void{}}}
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

	underlying.initChunk()
	return
}

func (ri *rtImp) initTank(t Type, kind RegisterKind) *Tank {
	tank := &Tank{}
	tank.typ = t
	tank.Register.id = -1

	if t.Kind() < BaseTypeNum {
		if kind != KImv || t.Kind() != TypeKindChunk {
			tank.Register.typ = ri.baseTypes[t.Kind()].regType
		} else {
			tank.Register.typ = ri.baseTypes[TypeKindInt].regType
		}

		tank.Register._kind = kind
		tank.Register._otype = t
		return tank
	}

	switch t := t.(type) {
	case *Tuple:
		for _, m := range t.members {
			tm := ri.initTank(m, kind)
			tank.Member = append(tank.Member, tm)
		}

	case *Struct:
		for _, m := range t.members {
			tm := ri.initTank(m.Type, kind)
			tank.Member = append(tank.Member, tm)
		}

	case *String:
		if t.underlying == nil {
			t.underlying = ri.underlyingStruct(t)
		}

		tank = ri.initTank(t.underlying, kind)

	case *Ref:
		if t.underlying == nil {
			t.underlying = ri.underlyingStruct(t)
		}

		tank = ri.initTank(t.underlying, kind)

	case *Named:
		tank = ri.initTank(t.underlying, kind)

	default:
		panic(fmt.Sprintf("Todo: %T", t))
	}

	return tank
}
