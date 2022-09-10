package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

func NewVar(name string, kind ValueKind, typ ValueType) Value {
	switch typ.(type) {
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

	default:
		logger.Fatalf("Todo: %T", typ)
	}

	return nil
}

type aVar struct {
	name string
	kind ValueKind
	typ  ValueType
}

func (v *aVar) Name() string    { return v.name }
func (v *aVar) Kind() ValueKind { return v.kind }
func (v *aVar) Type() ValueType { return v.typ }
func (v *aVar) rawGet(name string) wat.Inst {
	switch v.kind {
	case ValueKindLocal:
		return wat.NewInstGetLocal(name)

	default:
		logger.Fatal("Todo")
		return nil
	}
}
func (v *aVar) rawSet(name string) wat.Inst {
	switch v.kind {
	case ValueKindLocal:
		return wat.NewInstSetLocal(name)

	default:
		logger.Fatal("Todo")
		return nil
	}
}

/**************************************
varI32:
**************************************/
type varI32 struct {
	aVar
}

func newVarI32(name string, kind ValueKind) *varI32 {
	return &varI32{aVar: aVar{name: name, kind: kind, typ: I32{}}}
}
func (v *varI32) raw() []wat.Value { return []wat.Value{wat.NewVarI32(v.name)} }
func (v *varI32) EmitInit() []wat.Inst {
	return []wat.Inst{wat.NewInstConst(wat.I32{}, "0"), v.rawSet(v.name)}
}
func (v *varI32) EmitGet() []wat.Inst     { return []wat.Inst{v.rawGet(v.name)} }
func (v *varI32) EmitSet() []wat.Inst     { return []wat.Inst{v.rawSet(v.name)} }
func (v *varI32) EmitRelease() []wat.Inst { return nil }

/**************************************
varU32:
**************************************/
type varU32 struct {
	aVar
}

func newVarU32(name string, kind ValueKind) *varU32 {
	return &varU32{aVar: aVar{name: name, kind: kind, typ: U32{}}}
}
func (v *varU32) raw() []wat.Value { return []wat.Value{wat.NewVarU32(v.name)} }
func (v *varU32) EmitInit() []wat.Inst {
	return []wat.Inst{wat.NewInstConst(wat.U32{}, "0"), v.rawSet(v.name)}
}
func (v *varU32) EmitGet() []wat.Inst     { return []wat.Inst{v.rawGet(v.name)} }
func (v *varU32) EmitSet() []wat.Inst     { return []wat.Inst{v.rawSet(v.name)} }
func (v *varU32) EmitRelease() []wat.Inst { return nil }

/**************************************
varI64:
**************************************/
type varI64 struct {
	aVar
}

func newVarI64(name string, kind ValueKind) *varI64 {
	return &varI64{aVar: aVar{name: name, kind: kind, typ: I64{}}}
}
func (v *varI64) raw() []wat.Value { return []wat.Value{wat.NewVarI64(v.name)} }
func (v *varI64) EmitInit() []wat.Inst {
	return []wat.Inst{wat.NewInstConst(wat.I64{}, "0"), v.rawSet(v.name)}
}
func (v *varI64) EmitGet() []wat.Inst     { return []wat.Inst{v.rawGet(v.name)} }
func (v *varI64) EmitSet() []wat.Inst     { return []wat.Inst{v.rawSet(v.name)} }
func (v *varI64) EmitRelease() []wat.Inst { return nil }

/**************************************
varU64:
**************************************/
type varU64 struct {
	aVar
}

func newVarU64(name string, kind ValueKind) *varU64 {
	return &varU64{aVar: aVar{name: name, kind: kind, typ: U64{}}}
}
func (v *varU64) raw() []wat.Value { return []wat.Value{wat.NewVarU64(v.name)} }
func (v *varU64) EmitInit() []wat.Inst {
	return []wat.Inst{wat.NewInstConst(wat.U64{}, "0"), v.rawSet(v.name)}
}
func (v *varU64) EmitGet() []wat.Inst     { return []wat.Inst{v.rawGet(v.name)} }
func (v *varU64) EmitSet() []wat.Inst     { return []wat.Inst{v.rawSet(v.name)} }
func (v *varU64) EmitRelease() []wat.Inst { return nil }

/**************************************
varF32:
**************************************/
type varF32 struct {
	aVar
}

func newVarF32(name string, kind ValueKind) *varF32 {
	return &varF32{aVar: aVar{name: name, kind: kind, typ: F32{}}}
}
func (v *varF32) raw() []wat.Value { return []wat.Value{wat.NewVarF32(v.name)} }
func (v *varF32) EmitInit() []wat.Inst {
	return []wat.Inst{wat.NewInstConst(wat.F32{}, "0"), v.rawSet(v.name)}
}
func (v *varF32) EmitGet() []wat.Inst     { return []wat.Inst{v.rawGet(v.name)} }
func (v *varF32) EmitSet() []wat.Inst     { return []wat.Inst{v.rawSet(v.name)} }
func (v *varF32) EmitRelease() []wat.Inst { return nil }

/**************************************
varF64:
**************************************/
type varF64 struct {
	aVar
}

func newVarF64(name string, kind ValueKind) *varF64 {
	return &varF64{aVar: aVar{name: name, kind: kind, typ: F64{}}}
}
func (v *varF64) raw() []wat.Value { return []wat.Value{wat.NewVarF64(v.name)} }
func (v *varF64) EmitInit() []wat.Inst {
	return []wat.Inst{wat.NewInstConst(wat.F64{}, "0"), v.rawSet(v.name)}
}
func (v *varF64) EmitGet() []wat.Inst     { return []wat.Inst{v.rawGet(v.name)} }
func (v *varF64) EmitSet() []wat.Inst     { return []wat.Inst{v.rawSet(v.name)} }
func (v *varF64) EmitRelease() []wat.Inst { return nil }

/**************************************
VarPointer:
**************************************/
//type VarPointer struct {
//	aVar
//}
//
//func (v *VarPointer) raw() []Value { return []Value{NewVar(v.name, v.kind, Int32{})} }
//func (v *VarPointer) EmitInit() []wat.Inst {
//	return []wat.Inst{wat.NewInstConst(wat.I32{}, "0"), emitSetVar(v.kind, v.name)}
//}
//func (v *VarPointer) EmitGet() []wat.Inst     { return []wat.Inst{emitGetVar(v.kind, v.name)} }
//func (v *VarPointer) EmitSet() []wat.Inst     { return []wat.Inst{emitSetVar(v.kind, v.name)} }
//func (v *VarPointer) EmitRelease() []wat.Inst { return nil }

/**************************************
VarRef:
**************************************/
type VarRef struct {
	aVar

	data  wat.VarI32
	block wat.VarI32
}

func NewVarRef(name string, kind ValueKind, base_type ValueType) *VarRef {
	var v VarRef
	v.aVar = aVar{name: name, kind: kind, typ: NewRef(base_type)}
	v.data = *wat.NewVarI32(name + ".data")
	v.block = *wat.NewVarI32(name + ".block")
	return &v
}

func (v *VarRef) raw() []wat.Value {
	var wat_raw []wat.Value
	wat_raw = append(wat_raw, &v.data)
	wat_raw = append(wat_raw, &v.block)
	return wat_raw
}

func (v *VarRef) EmitInit() []wat.Inst {
	var insts []wat.Inst
	insts = append(insts, wat.NewInstConst(wat.I32{}, "0"))
	insts = append(insts, v.rawSet(v.data.Name()))
	insts = append(insts, wat.NewInstConst(wat.I32{}, "0"))
	insts = append(insts, v.rawSet(v.block.Name()))
	return insts
}

func (v *VarRef) EmitGet() []wat.Inst {
	var insts []wat.Inst
	insts = append(insts, v.rawGet(v.data.Name()))
	insts = append(insts, v.rawGet(v.block.Name()))
	return insts
}

func (v *VarRef) EmitSet() []wat.Inst {
	insts := v.EmitRelease()
	insts = append(insts, v.rawSet(v.block.Name()))
	insts = append(insts, v.rawSet(v.data.Name()))
	return insts
}

func (v *VarRef) EmitRelease() []wat.Inst {
	return nil
}
