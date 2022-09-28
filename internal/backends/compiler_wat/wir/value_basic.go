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
func (v *aVar) push(name string) wat.Inst {
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
func (v *aVar) pop(name string) wat.Inst {
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
	return []wat.Inst{wat.NewInstConst(toWatType(v.Type()), "0"), v.pop(v.name)}
}
func (v *varBasic) EmitPush() []wat.Inst    { return []wat.Inst{v.push(v.name)} }
func (v *varBasic) EmitPop() []wat.Inst     { return []wat.Inst{v.pop(v.name)} }
func (v *varBasic) EmitRelease() []wat.Inst { return nil }
func (v *varBasic) emitLoad(addr Value) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(v.Type()), 0, 1))
	return insts
}
func (v *varBasic) emitStore(addr Value) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, v.EmitPush()...)
	insts = append(insts, wat.NewInstStore(toWatType(v.Type()), 0, 1))
	return insts
}
