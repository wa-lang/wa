package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

func NewVar(name string, kind ValueKind, typ ValueType) Value {
	switch typ.(type) {
	case Int32:
		return NewVarI32(name, kind)

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
VarI32:
**************************************/
type VarI32 struct {
	aVar
}

func NewVarI32(name string, kind ValueKind) *VarI32 {
	return &VarI32{aVar: aVar{name: name, kind: kind, typ: Int32{}}}
}
func (v *VarI32) raw() []wat.Value { return []wat.Value{wat.NewVarI32(v.name)} }
func (v *VarI32) EmitInit() []wat.Inst {
	return []wat.Inst{wat.NewInstConst(wat.I32{}, "0"), v.rawGet(v.name)}
}
func (v *VarI32) EmitGet() []wat.Inst     { return []wat.Inst{v.rawGet(v.name)} }
func (v *VarI32) EmitSet() []wat.Inst     { return []wat.Inst{v.rawSet(v.name)} }
func (v *VarI32) EmitRelease() []wat.Inst { return nil }

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
