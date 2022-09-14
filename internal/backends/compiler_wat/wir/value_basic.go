// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

func NewVar(name string, kind ValueKind, typ ValueType) Value {
	switch typ := typ.(type) {
	case I32, U32, I64, U64, F32, F64, Pointer:
		return newVarBasic(name, kind, typ)

	/*
		case I32:
			return newVarI32(name, kind)

		case U32:
			return newVarU32(name, kind)

		case I64:
			return newVarI64(name, kind)

		case U64:
			return newVarU64(name, kind)

		case F32:
			return newVarF32(name, kind)

		case F64:
			return newVarF64(name, kind)

		case Pointer:
			return newVarPointer(name, kind, typ.Base) //*/

	case Block:
		return newVarBlock(name, kind, typ.Base)

	case Struct:
		return newVarStruct(name, kind, typ)

	case Ref:
		return NewVarRef(name, kind, typ.Base)

	default:
		logger.Fatalf("Todo: %T", typ)
	}

	return nil
}

/**************************************
aVar:
**************************************/
type aVar struct {
	name string
	kind ValueKind
	typ  ValueType
}

func (v *aVar) Name() string    { return v.name }
func (v *aVar) Kind() ValueKind { return v.kind }
func (v *aVar) Type() ValueType { return v.typ }
func (v *aVar) get(name string) wat.Inst {
	switch v.kind {
	case ValueKindLocal:
		return wat.NewInstGetLocal(name)

	case ValueKindGlobal:
		return wat.NewInstGetGlobal(name)

	default:
		logger.Fatal("Todo")
		return nil
	}
}
func (v *aVar) set(name string) wat.Inst {
	switch v.kind {
	case ValueKindLocal:
		return wat.NewInstSetLocal(name)

	case ValueKindGlobal:
		return wat.NewInstSetGlobal(name)

	default:
		logger.Fatal("Todo")
		return nil
	}
}

/**************************************
varBasic:
**************************************/
type varBasic struct {
	aVar
}

func newVarBasic(name string, kind ValueKind, typ ValueType) *varBasic {
	return &varBasic{aVar: aVar{name: name, kind: kind, typ: typ}}
}
func (v *varBasic) raw() []wat.Value { return []wat.Value{wat.NewVar(v.name, toWatType(v.Type()))} }
func (v *varBasic) EmitInit() []wat.Inst {
	return []wat.Inst{wat.NewInstConst(toWatType(v.Type()), "0"), v.set(v.name)}
}
func (v *varBasic) EmitGet() []wat.Inst     { return []wat.Inst{v.get(v.name)} }
func (v *varBasic) EmitSet() []wat.Inst     { return []wat.Inst{v.set(v.name)} }
func (v *varBasic) EmitRelease() []wat.Inst { return nil }
func (v *varBasic) emitLoad(addr Value) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitGet()
	insts = append(insts, wat.NewInstLoad(toWatType(v.Type()), 0, 1))
	return insts
}
func (v *varBasic) emitStore(addr Value) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitGet()
	insts = append(insts, v.EmitGet()...)
	insts = append(insts, wat.NewInstStore(toWatType(v.Type()), 0, 1))
	return insts
}

///**************************************
//varI32:
//**************************************/
//type varI32 struct {
//	aVar
//}
//
//func newVarI32(name string, kind ValueKind) *varI32 {
//	return &varI32{aVar: aVar{name: name, kind: kind, typ: I32{}}}
//}
//func (v *varI32) raw() []wat.Value { return []wat.Value{wat.NewVarI32(v.name)} }
//func (v *varI32) EmitInit() []wat.Inst {
//	return []wat.Inst{wat.NewInstConst(wat.I32{}, "0"), v.set(v.name)}
//}
//func (v *varI32) EmitGet() []wat.Inst     { return []wat.Inst{v.get(v.name)} }
//func (v *varI32) EmitSet() []wat.Inst     { return []wat.Inst{v.set(v.name)} }
//func (v *varI32) EmitRelease() []wat.Inst { return nil }
//func (v *varI32) emitLoad(addr varPointer) []wat.Inst {
//	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
//		logger.Fatal("Type not match")
//		return nil
//	}
//	insts := addr.EmitGet()
//	insts = append(insts, wat.NewInstLoad(wat.I32{}, 0, 1))
//	return insts
//}
//func (v *varI32) emitStore(addr varPointer) []wat.Inst {
//	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
//		logger.Fatal("Type not match")
//		return nil
//	}
//	insts := addr.EmitGet()
//	insts = append(insts, v.EmitGet()...)
//	insts = append(insts, wat.NewInstStore(wat.I32{}, 0, 1))
//	return insts
//}
//
///**************************************
//varU32:
//**************************************/
//type varU32 struct {
//	aVar
//}
//
//func newVarU32(name string, kind ValueKind) *varU32 {
//	return &varU32{aVar: aVar{name: name, kind: kind, typ: U32{}}}
//}
//func (v *varU32) raw() []wat.Value { return []wat.Value{wat.NewVarU32(v.name)} }
//func (v *varU32) EmitInit() []wat.Inst {
//	return []wat.Inst{wat.NewInstConst(wat.U32{}, "0"), v.set(v.name)}
//}
//func (v *varU32) EmitGet() []wat.Inst     { return []wat.Inst{v.get(v.name)} }
//func (v *varU32) EmitSet() []wat.Inst     { return []wat.Inst{v.set(v.name)} }
//func (v *varU32) EmitRelease() []wat.Inst { return nil }
//func (v *varU32) emitLoad(addr varPointer) []wat.Inst {
//	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
//		logger.Fatal("Type not match")
//		return nil
//	}
//	insts := addr.EmitGet()
//	insts = append(insts, wat.NewInstLoad(wat.U32{}, 0, 1))
//	return insts
//}
//func (v *varU32) emitStore(addr varPointer) []wat.Inst {
//	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
//		logger.Fatal("Type not match")
//		return nil
//	}
//	insts := addr.EmitGet()
//	insts = append(insts, v.EmitGet()...)
//	insts = append(insts, wat.NewInstStore(wat.U32{}, 0, 1))
//	return insts
//}
//
///**************************************
//varI64:
//**************************************/
//type varI64 struct {
//	aVar
//}
//
//func newVarI64(name string, kind ValueKind) *varI64 {
//	return &varI64{aVar: aVar{name: name, kind: kind, typ: I64{}}}
//}
//func (v *varI64) raw() []wat.Value { return []wat.Value{wat.NewVarI64(v.name)} }
//func (v *varI64) EmitInit() []wat.Inst {
//	return []wat.Inst{wat.NewInstConst(wat.I64{}, "0"), v.set(v.name)}
//}
//func (v *varI64) EmitGet() []wat.Inst     { return []wat.Inst{v.get(v.name)} }
//func (v *varI64) EmitSet() []wat.Inst     { return []wat.Inst{v.set(v.name)} }
//func (v *varI64) EmitRelease() []wat.Inst { return nil }
//func (v *varI64) emitLoad(addr varPointer) []wat.Inst {
//	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
//		logger.Fatal("Type not match")
//		return nil
//	}
//	insts := addr.EmitGet()
//	insts = append(insts, wat.NewInstLoad(wat.I64{}, 0, 1))
//	return insts
//}
//func (v *varI64) emitStore(addr varPointer) []wat.Inst {
//	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
//		logger.Fatal("Type not match")
//		return nil
//	}
//	insts := addr.EmitGet()
//	insts = append(insts, v.EmitGet()...)
//	insts = append(insts, wat.NewInstStore(wat.I64{}, 0, 1))
//	return insts
//}
//
///**************************************
//varU64:
//**************************************/
//type varU64 struct {
//	aVar
//}
//
//func newVarU64(name string, kind ValueKind) *varU64 {
//	return &varU64{aVar: aVar{name: name, kind: kind, typ: U64{}}}
//}
//func (v *varU64) raw() []wat.Value { return []wat.Value{wat.NewVarU64(v.name)} }
//func (v *varU64) EmitInit() []wat.Inst {
//	return []wat.Inst{wat.NewInstConst(wat.U64{}, "0"), v.set(v.name)}
//}
//func (v *varU64) EmitGet() []wat.Inst     { return []wat.Inst{v.get(v.name)} }
//func (v *varU64) EmitSet() []wat.Inst     { return []wat.Inst{v.set(v.name)} }
//func (v *varU64) EmitRelease() []wat.Inst { return nil }
//func (v *varU64) emitLoad(addr varPointer) []wat.Inst {
//	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
//		logger.Fatal("Type not match")
//		return nil
//	}
//	insts := addr.EmitGet()
//	insts = append(insts, wat.NewInstLoad(wat.U64{}, 0, 1))
//	return insts
//}
//func (v *varU64) emitStore(addr varPointer) []wat.Inst {
//	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
//		logger.Fatal("Type not match")
//		return nil
//	}
//	insts := addr.EmitGet()
//	insts = append(insts, v.EmitGet()...)
//	insts = append(insts, wat.NewInstStore(wat.U64{}, 0, 1))
//	return insts
//}
//
///**************************************
//varF32:
//**************************************/
//type varF32 struct {
//	aVar
//}
//
//func newVarF32(name string, kind ValueKind) *varF32 {
//	return &varF32{aVar: aVar{name: name, kind: kind, typ: F32{}}}
//}
//func (v *varF32) raw() []wat.Value { return []wat.Value{wat.NewVarF32(v.name)} }
//func (v *varF32) EmitInit() []wat.Inst {
//	return []wat.Inst{wat.NewInstConst(wat.F32{}, "0"), v.set(v.name)}
//}
//func (v *varF32) EmitGet() []wat.Inst     { return []wat.Inst{v.get(v.name)} }
//func (v *varF32) EmitSet() []wat.Inst     { return []wat.Inst{v.set(v.name)} }
//func (v *varF32) EmitRelease() []wat.Inst { return nil }
//func (v *varF32) emitLoad(addr varPointer) []wat.Inst {
//	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
//		logger.Fatal("Type not match")
//		return nil
//	}
//	insts := addr.EmitGet()
//	insts = append(insts, wat.NewInstLoad(wat.F32{}, 0, 1))
//	return insts
//}
//func (v *varF32) emitStore(addr varPointer) []wat.Inst {
//	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
//		logger.Fatal("Type not match")
//		return nil
//	}
//	insts := addr.EmitGet()
//	insts = append(insts, v.EmitGet()...)
//	insts = append(insts, wat.NewInstStore(wat.F32{}, 0, 1))
//	return insts
//}
//
///**************************************
//varF64:
//**************************************/
//type varF64 struct {
//	aVar
//}
//
//func newVarF64(name string, kind ValueKind) *varF64 {
//	return &varF64{aVar: aVar{name: name, kind: kind, typ: F64{}}}
//}
//func (v *varF64) raw() []wat.Value { return []wat.Value{wat.NewVarF64(v.name)} }
//func (v *varF64) EmitInit() []wat.Inst {
//	return []wat.Inst{wat.NewInstConst(wat.F64{}, "0"), v.set(v.name)}
//}
//func (v *varF64) EmitGet() []wat.Inst     { return []wat.Inst{v.get(v.name)} }
//func (v *varF64) EmitSet() []wat.Inst     { return []wat.Inst{v.set(v.name)} }
//func (v *varF64) EmitRelease() []wat.Inst { return nil }
//func (v *varF64) emitLoad(addr varPointer) []wat.Inst {
//	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
//		logger.Fatal("Type not match")
//		return nil
//	}
//	insts := addr.EmitGet()
//	insts = append(insts, wat.NewInstLoad(wat.F64{}, 0, 1))
//	return insts
//}
//func (v *varF64) emitStore(addr varPointer) []wat.Inst {
//	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
//		logger.Fatal("Type not match")
//		return nil
//	}
//	insts := addr.EmitGet()
//	insts = append(insts, v.EmitGet()...)
//	insts = append(insts, wat.NewInstStore(wat.F64{}, 0, 1))
//	return insts
//}
//
///**************************************
//varPointer:
//**************************************/
//type varPointer struct {
//	aVar
//}
//
//func newVarPointer(name string, kind ValueKind, base_type ValueType) *varPointer {
//	return &varPointer{aVar: aVar{name: name, kind: kind, typ: NewPointer(base_type)}}
//}
//func (v *varPointer) raw() []wat.Value { return []wat.Value{wat.NewVarI32(v.name)} }
//func (v *varPointer) EmitInit() []wat.Inst {
//	return []wat.Inst{wat.NewInstConst(wat.I32{}, "0"), v.set(v.name)}
//}
//func (v *varPointer) EmitGet() []wat.Inst     { return []wat.Inst{v.get(v.name)} }
//func (v *varPointer) EmitSet() []wat.Inst     { return []wat.Inst{v.set(v.name)} }
//func (v *varPointer) EmitRelease() []wat.Inst { return nil }
//func (v *varPointer) emitLoad(addr varPointer) []wat.Inst {
//	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
//		logger.Fatal("Type not match")
//		return nil
//	}
//	insts := addr.EmitGet()
//	insts = append(insts, wat.NewInstLoad(wat.I32{}, 0, 1))
//	return insts
//}
//func (v *varPointer) emitStore(addr varPointer) []wat.Inst {
//	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
//		logger.Fatal("Type not match")
//		return nil
//	}
//	insts := addr.EmitGet()
//	insts = append(insts, v.EmitGet()...)
//	insts = append(insts, wat.NewInstStore(wat.I32{}, 0, 1))
//	return insts
//}
