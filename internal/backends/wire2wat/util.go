// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire2wat

import (
	"fmt"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/wire"
)

type baseType struct {
	watType wat.ValueType
	size    int
}

type Wire2Wat struct {
	baseTypes []baseType
}

func (w *Wire2Wat) init() {
	w.baseTypes = make([]baseType, wire.BaseTypeNum)

	w.baseTypes[wire.TypeKindUnknown] = baseType{watType: nil, size: 0}
	w.baseTypes[wire.TypeKindVoid] = baseType{watType: nil, size: 0}
	w.baseTypes[wire.TypeKindBool] = baseType{watType: wat.I32{}, size: 1}
	w.baseTypes[wire.TypeKindU8] = baseType{watType: wat.U32{}, size: 1}
	w.baseTypes[wire.TypeKindU16] = baseType{watType: wat.U32{}, size: 2}
	w.baseTypes[wire.TypeKindU32] = baseType{watType: wat.U32{}, size: 4}
	w.baseTypes[wire.TypeKindU64] = baseType{watType: wat.U64{}, size: 8}
	w.baseTypes[wire.TypeKindI8] = baseType{watType: wat.I32{}, size: 1}
	w.baseTypes[wire.TypeKindI16] = baseType{watType: wat.I32{}, size: 2}
	w.baseTypes[wire.TypeKindI32] = baseType{watType: wat.I32{}, size: 4}
	w.baseTypes[wire.TypeKindI64] = baseType{watType: wat.I64{}, size: 8}
	w.baseTypes[wire.TypeKindInt] = baseType{watType: wat.I32{}, size: 4}
	w.baseTypes[wire.TypeKindUint] = baseType{watType: wat.U32{}, size: 4}
	w.baseTypes[wire.TypeKindF32] = baseType{watType: wat.F32{}, size: 4}
	w.baseTypes[wire.TypeKindF64] = baseType{watType: wat.F64{}, size: 8}
	w.baseTypes[wire.TypeKindRune] = baseType{watType: wat.U32{}, size: 4}
	w.baseTypes[wire.TypeKindPtr] = baseType{watType: wat.U32{}, size: 4}
	w.baseTypes[wire.TypeKindChunk] = baseType{watType: wat.U32{}, size: 4}
}

func (w *Wire2Wat) watType(t wire.Type) wat.ValueType {
	if t.Kind() >= wire.BaseTypeNum {
		panic(fmt.Sprintf("t is not a base type: %s", t.Name()))
	}
	return w.baseTypes[t.Kind()].watType
}

func (w *Wire2Wat) rawof(t wire.Type) (raw []wire.Type) {
	if t == nil {
		panic("t is nil")
	}

	if t.Kind() < wire.BaseTypeNum {
		return []wire.Type{t}
	}

	switch t := t.(type) {
	case *wire.Tuple:
		for i := 0; i < t.Len(); i++ {
			raw = append(raw, w.rawof(t.At(i))...)
		}

	case *wire.Struct:
		for i := 0; i < t.Len(); i++ {
			raw = append(raw, w.rawof(t.At(i).Type)...)
		}

	case *wire.String:
		us := t.Underlying()
		if us == nil {
			panic("String.Underlying() is nil")
		}
		raw = w.rawof(us)

	case *wire.Ref:
		us := t.Underlying()
		if us == nil {
			panic("String.Underlying() is nil")
		}
		raw = w.rawof(us)

	case *wire.Named:
		raw = w.rawof(t.Underlying())

	default:
		panic(fmt.Sprintf("Todo: %T", t))
	}

	return
}

func (w *Wire2Wat) sizeof(t wire.Type) int {
	if t.Kind() < wire.BaseTypeNum {
		return w.baseTypes[t.Kind()].size
	}

	switch t := t.(type) {
	case *wire.Tuple:
		panic("Tuple is not supported")

	case *wire.Struct:
		if !t.SizeInitialized {
			w.initStructSize(t)
		}
		return t.Size

	case *wire.String:
		us := t.Underlying()
		if us == nil {
			panic("String.Underlying() is nil")
		}
		if !us.SizeInitialized {
			w.initStructSize(us)
		}
		return us.Size

	case *wire.Ref:
		us := t.Underlying()
		if us == nil {
			panic("String.Underlying() is nil")
		}
		if !us.SizeInitialized {
			w.initStructSize(us)
		}
		return us.Size

	case *wire.Named:
		return w.sizeof(t.Underlying())

	default:
		panic(fmt.Sprintf("Todo: %T", t))
	}
}

func (w *Wire2Wat) alignof(t wire.Type) int {
	if t.Kind() < wire.BaseTypeNum {
		return w.baseTypes[t.Kind()].size
	}

	switch t := t.(type) {
	case *wire.Tuple:
		panic("Tuple is not supported")

	case *wire.Struct:
		if !t.SizeInitialized {
			w.initStructSize(t)
		}
		return t.Align

	case *wire.String:
		us := t.Underlying()
		if us == nil {
			panic("String.Underlying() is nil")
		}
		if !us.SizeInitialized {
			w.initStructSize(us)
		}
		return us.Align

	case *wire.Ref:
		us := t.Underlying()
		if us == nil {
			panic("String.Underlying() is nil")
		}
		if !us.SizeInitialized {
			w.initStructSize(us)
		}
		return us.Align

	case *wire.Named:
		return w.alignof(t.Underlying())

	default:
		panic(fmt.Sprintf("Todo: %T", t))
	}
}

func (w *Wire2Wat) initStructSize(t *wire.Struct) {
	t.Size = 0
	t.Align = 0
	for i := 0; i < t.Len(); i++ {
		member := t.At(i)
		ma := w.alignof(member.Type)
		ms := w.sizeof(member.Type)
		member.Start = makeAlign(t.Size, ma)

		t.Size = member.Start + ms
		if ma > t.Align {
			t.Align = ma
		}
	}

	t.Size = makeAlign(t.Size, t.Align)
}

func makeAlign(i, align int) int {
	if align == 1 || align == 0 {
		return i
	}
	return (i + align - 1) / align * align
}
