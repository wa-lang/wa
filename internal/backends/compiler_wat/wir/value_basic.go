// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

func NewConst(lit string, t ValueType) Value {
	return newValue(lit, ValueKindConst, t)
}

func NewLocal(name string, typ ValueType) Value {
	return newValue(name, ValueKindLocal, typ)
}

func NewGlobal(name string, typ ValueType, as_pointer bool) Value {
	if as_pointer {
		return newValue(name, ValueKindGlobal_Pointer, typ)
	} else {
		return newValue(name, ValueKindGlobal_Value, typ)
	}
}

func newValue(name string, kind ValueKind, typ ValueType) Value {
	switch typ := typ.(type) {
	case I32, U32, I64, U64, F32, F64, RUNE:
		return newVarBasic(name, kind, typ)

	case Pointer:
		return newVarPointer(name, kind, typ.Base)

	case Block:
		return newVarBlock(name, kind, typ.Base)

	case Struct:
		return newVarStruct(name, kind, typ)

	case Ref:
		return newVarRef(name, kind, typ.Base)

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
func (v *aVar) push(name string) wat.Inst {
	switch v.kind {
	case ValueKindLocal:
		return wat.NewInstGetLocal(name)

	case ValueKindGlobal_Value, ValueKindGlobal_Pointer:
		return wat.NewInstGetGlobal(name)

	case ValueKindConst:
		return wat.NewInstConst(toWatType(v.Type()), name)

	default:
		logger.Fatal("Unreachable.")
		return nil
	}
}
func (v *aVar) pop(name string) wat.Inst {
	switch v.kind {
	case ValueKindLocal:
		return wat.NewInstSetLocal(name)

	case ValueKindGlobal_Value, ValueKindGlobal_Pointer:
		return wat.NewInstSetGlobal(name)

	case ValueKindConst:
		logger.Fatal("Can't pop to const.")
		return nil

	default:
		logger.Fatal("Unreachable.")
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

func (v *varBasic) raw() []wat.Value        { return []wat.Value{wat.NewVar(v.name, toWatType(v.Type()))} }
func (v *varBasic) EmitPush() []wat.Inst    { return []wat.Inst{v.push(v.name)} }
func (v *varBasic) EmitPop() []wat.Inst     { return []wat.Inst{v.pop(v.name)} }
func (v *varBasic) EmitRelease() []wat.Inst { return nil }

func (v *varBasic) EmitInit() (insts []wat.Inst) {
	insts = append(insts, wat.NewInstConst(toWatType(v.Type()), "0"))
	insts = append(insts, v.pop(v.name))
	return
}

func (v *varBasic) emitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(v.Type()), offset, 1))
	return insts
}

func (v *varBasic) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, v.EmitPush()...)
	insts = append(insts, wat.NewInstStore(toWatType(v.Type()), offset, 1))
	return insts
}

/**************************************
VarPointer:
**************************************/
type VarPointer struct {
	varBasic
}

func newVarPointer(name string, kind ValueKind, base_type ValueType) *VarPointer {
	var v VarPointer
	pointer_type := NewPointer(base_type)
	v.aVar = aVar{name: name, kind: kind, typ: pointer_type}
	return &v
}

func (v *VarPointer) emitGetValue() []wat.Inst {
	t := newValue("", v.kind, v.Type().(Pointer).Base)
	return t.emitLoadFromAddr(v, 0)
}

func (v *VarPointer) emitSetValue(d Value) []wat.Inst {
	if !d.Type().Equal(v.Type().(Pointer).Base) {
		logger.Fatal("Type not match")
		return nil
	}
	return d.emitStoreToAddr(v, 0)
}
