// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import "fmt"

type RtImp interface {
	underlyingStruct(t Type) Struct
	initTank(t Type) *tank
}

var rtimp RtImp

func init() {
	imp := &rtImp{}
	imp.init()
	rtimp = imp
}

type rtImp struct {
	i32, i64, f32, f64, chunk Type
}

func (ri *rtImp) init() {
	ri.i32 = &I32{}
	ri.i64 = &I64{}
	ri.f32 = &F32{}
	ri.f64 = &F64{}
	ri.chunk = &chunk{Base: &Void{}}
}

func (ri *rtImp) underlyingStruct(t Type) (underlying Struct) {
	switch t := t.(type) {
	case *Ref:
		c := StructMember{Name: "c", Type: &chunk{Base: t.Base}, id: 0}
		d := StructMember{Name: "d", Type: &Ptr{Base: t.Base}, id: 1}
		underlying.member = []StructMember{c, d}

	case *String:
		c := StructMember{Name: "c", Type: &chunk{Base: &U8{}}, id: 0}
		d := StructMember{Name: "d", Type: &Ptr{Base: &U8{}}, id: 1}
		l := StructMember{Name: "l", Type: &Uint{}, id: 2}
		underlying.member = []StructMember{c, d, l}

	default:
		panic(fmt.Sprintf("Todo: %T", t))
	}

	return
}

func (ri *rtImp) initTank(t Type) *tank {
	tank := &tank{}
	tank.typ = t
	tank.register.id = -1

	switch t := t.(type) {
	case *Void:
		tank.register.typ = t

	case *Bool, *U8, *U16, *U32, *Uint, *I8, *I16, *I32, *Int, *Rune:
		tank.register.typ = ri.i32

	case *I64, *U64:
		tank.register.typ = ri.i64

	case *F32:
		tank.register.typ = ri.f32

	case *F64:
		tank.register.typ = ri.f64

	case *Ptr:
		tank.register.typ = ri.i32

	case *chunk:
		tank.register.typ = ri.chunk

	case *Tuple:
		for _, m := range t.fields {
			tm := ri.initTank(m)
			tank.member = append(tank.member, tm)
		}

	case *Struct:
		for _, m := range t.member {
			tm := ri.initTank(m.Type)
			tank.member = append(tank.member, tm)
		}

	case *String:
		tank = ri.initTank(t.Underlying())

	case *Ref:
		tank = ri.initTank(t.Underlying())

	default:
		panic(fmt.Sprintf("Todo: %T", t))
	}

	return tank
}
