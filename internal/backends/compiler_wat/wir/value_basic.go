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
		return newValueBasic(name, kind, typ)

	case Pointer:
		return newValuePointer(name, kind, typ.Base)

	case Block:
		return newValueBlock(name, kind, typ.Base)

	case Struct:
		return newValueStruct(name, kind, typ)

	case Ref:
		return newValueRef(name, kind, typ.Base)

	default:
		logger.Fatalf("Todo: %T", typ)
	}

	return nil
}

/**************************************
aValue:
**************************************/
type aValue struct {
	name string
	kind ValueKind
	typ  ValueType
}

func (v *aValue) Name() string    { return v.name }
func (v *aValue) Kind() ValueKind { return v.kind }
func (v *aValue) Type() ValueType { return v.typ }
func (v *aValue) push(name string) wat.Inst {
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
func (v *aValue) pop(name string) wat.Inst {
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
aBasic:
**************************************/
type aBasic struct {
	aValue
}

func newValueBasic(name string, kind ValueKind, typ ValueType) *aBasic {
	return &aBasic{aValue: aValue{name: name, kind: kind, typ: typ}}
}

func (v *aBasic) raw() []wat.Value        { return []wat.Value{wat.NewVar(v.name, toWatType(v.Type()))} }
func (v *aBasic) EmitPush() []wat.Inst    { return []wat.Inst{v.push(v.name)} }
func (v *aBasic) EmitPop() []wat.Inst     { return []wat.Inst{v.pop(v.name)} }
func (v *aBasic) EmitRelease() []wat.Inst { return nil }

func (v *aBasic) EmitInit() (insts []wat.Inst) {
	insts = append(insts, wat.NewInstConst(toWatType(v.Type()), "0"))
	insts = append(insts, v.pop(v.name))
	return
}

func (v *aBasic) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, v.EmitPush()...)
	insts = append(insts, wat.NewInstStore(toWatType(v.Type()), offset, 1))
	return insts
}
