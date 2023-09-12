// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
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
	case *tI8, *tU8, *tI16, *tU16, *tI32, *tU32, *tI64, *tU64, *tF32, *tF64, *tRune, *tBool:
		return newValue_Basic(name, kind, typ)

	case *Ptr:
		return newValue_Ptr(name, kind, typ)

	case *Block:
		return newValue_Block(name, kind, typ)

	case *Ref:
		return newValue_Ref(name, kind, typ)

	case *Array:
		return newValue_Array(name, kind, typ)

	case *Slice:
		return newValue_Slice(name, kind, typ)

	case *String:
		return newValue_String(name, kind, typ)

	case *Closure:
		return newValue_Closure(name, kind, typ)

	case *Interface:
		return newValue_Interface(name, kind, typ)

	case *Tuple:
		return newValue_Tuple(name, kind, typ)

	case *Struct:
		return newValue_Struct(name, kind, typ)

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

func newValue_Basic(name string, kind ValueKind, typ ValueType) *aBasic {
	return &aBasic{aValue: aValue{name: name, kind: kind, typ: typ}}
}

func (v *aBasic) raw() []wat.Value             { return []wat.Value{wat.NewVar(v.name, toWatType(v.Type()))} }
func (v *aBasic) EmitPush() []wat.Inst         { return []wat.Inst{v.push(v.name)} }
func (v *aBasic) EmitPushNoRetain() []wat.Inst { return []wat.Inst{v.push(v.name)} }
func (v *aBasic) EmitPop() []wat.Inst          { return []wat.Inst{v.pop(v.name)} }
func (v *aBasic) EmitRelease() []wat.Inst      { return nil }

func (v *aBasic) EmitInit() (insts []wat.Inst) {
	insts = append(insts, wat.NewInstConst(toWatType(v.Type()), "0"))
	insts = append(insts, v.pop(v.name))
	return
}

func (v *aBasic) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(v.Type()) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	insts := addr.EmitPush()
	insts = append(insts, v.EmitPush()...)
	switch v.Type().(type) {
	case *tU8, *tI8:
		insts = append(insts, wat.NewInstStore8(offset, 1))

	case *tU16, *tI16:
		insts = append(insts, wat.NewInstStore16(offset, 1))

	default:
		insts = append(insts, wat.NewInstStore(toWatType(v.Type()), offset, 1))
	}
	return insts
}

func (v *aBasic) emitStore(offset int) (insts []wat.Inst) {
	insts = append(insts, wat.NewInstCall("$wa.runtime.DupI32"))
	insts = append(insts, v.EmitPush()...)
	switch v.Type().(type) {
	case *tU8, *tI8:
		insts = append(insts, wat.NewInstStore8(offset, 1))

	case *tU16, *tI16:
		insts = append(insts, wat.NewInstStore16(offset, 1))

	default:
		insts = append(insts, wat.NewInstStore(toWatType(v.Type()), offset, 1))
	}

	return
}

func (v *aBasic) Bin() (b []byte) {
	if v.Kind() != ValueKindConst {
		panic("Value.bin(): const only!")
	}

	switch v.Type().(type) {
	case *tU8, *tI8:
		b = make([]byte, 1)
		i, _ := strconv.Atoi(v.Name())
		b[0] = byte(i & 0xFF)

	case *tU16, *tI16:
		b = make([]byte, 2)
		i, _ := strconv.Atoi(v.Name())
		b[0] = byte(i & 0xFF)
		b[1] = byte((i >> 8) & 0xFF)

	case *tU32, *tI32:
		b = make([]byte, 4)
		i, _ := strconv.Atoi(v.Name())
		b[0] = byte(i & 0xFF)
		b[1] = byte((i >> 8) & 0xFF)
		b[2] = byte((i >> 16) & 0xFF)
		b[3] = byte((i >> 24) & 0xFF)

	default:
		panic("todo")
	}

	return
}

func (v *aBasic) emitEq(r Value) (insts []wat.Inst, ok bool) {
	if !v.Type().Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}
	insts = append(insts, v.EmitPush()...)
	insts = append(insts, r.EmitPush()...)
	insts = append(insts, wat.NewInstEq(toWatType(v.Type())))

	ok = true

	return
}
