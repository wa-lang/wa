// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"math"
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

func NewGlobal(name string, typ ValueType) Value {
	return newValue(name, ValueKindGlobal, typ)
}

func newValue(name string, kind ValueKind, typ ValueType) Value {
	switch typ := typ.(type) {
	case *I8, *U8, *I16, *U16, *I32, *U32, *I64, *U64, *F32, *F64, *Rune, *Bool:
		return newValue_Basic(name, kind, typ)

	case *Complex64:
		return newValue_Complex64(name, kind, typ)

	case *Complex128:
		return newValue_Complex128(name, kind, typ)

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

	case *Map:
		return newValue_Map(name, kind, typ)

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

	case ValueKindGlobal:
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

	case ValueKindGlobal:
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
	case *U8, *I8, *Bool:
		insts = append(insts, wat.NewInstStore8(offset, 1))

	case *U16, *I16:
		insts = append(insts, wat.NewInstStore16(offset, 2))

	default:
		insts = append(insts, wat.NewInstStore(toWatType(v.Type()), offset, v.Type().align()))
	}
	return insts
}

func (v *aBasic) emitStore(offset int) (insts []wat.Inst) {
	insts = append(insts, wat.NewInstCall("runtime.DupI32"))
	insts = append(insts, v.EmitPush()...)
	switch v.Type().(type) {
	case *U8, *I8, *Bool:
		insts = append(insts, wat.NewInstStore8(offset, 1))

	case *U16, *I16:
		insts = append(insts, wat.NewInstStore16(offset, 2))

	default:
		insts = append(insts, wat.NewInstStore(toWatType(v.Type()), offset, v.Type().align()))
	}

	return
}

func (v *aBasic) Bin() (b []byte) {
	if v.Kind() != ValueKindConst {
		panic("Value.bin(): const only!")
	}

	switch v.Type().(type) {
	case *U8, *Bool:
		b = make([]byte, 1)
		i, _ := strconv.ParseUint(v.Name(), 0, 8)
		b[0] = byte(i)

	case *I8:
		b = make([]byte, 1)
		i, _ := strconv.ParseInt(v.Name(), 0, 8)
		si := uint8(int8(i))
		b[0] = byte(si)

	case *U16:
		b = make([]byte, 2)
		i, _ := strconv.ParseUint(v.Name(), 0, 16)
		si := uint16(i)
		b[0] = byte(si & 0xFF)
		b[1] = byte((si >> 8) & 0xFF)

	case *I16:
		b = make([]byte, 2)
		i, _ := strconv.ParseInt(v.Name(), 0, 16)
		si := uint16(int16(i))
		b[0] = byte(si & 0xFF)
		b[1] = byte((si >> 8) & 0xFF)

	case *U32, *Rune:
		b = make([]byte, 4)
		i, _ := strconv.ParseUint(v.Name(), 0, 32)
		si := uint32(i)
		b[0] = byte(si & 0xFF)
		b[1] = byte((si >> 8) & 0xFF)
		b[2] = byte((si >> 16) & 0xFF)
		b[3] = byte((si >> 24) & 0xFF)

	case *I32:
		b = make([]byte, 4)
		i, _ := strconv.ParseInt(v.Name(), 0, 32)
		si := uint32(int32(i))
		b[0] = byte(si & 0xFF)
		b[1] = byte((si >> 8) & 0xFF)
		b[2] = byte((si >> 16) & 0xFF)
		b[3] = byte((si >> 24) & 0xFF)

	case *U64:
		b = make([]byte, 8)
		i, _ := strconv.ParseUint(v.Name(), 0, 64)
		si := uint64(i)
		b[0] = byte(si & 0xFF)
		b[1] = byte((si >> 8) & 0xFF)
		b[2] = byte((si >> 16) & 0xFF)
		b[3] = byte((si >> 24) & 0xFF)
		b[4] = byte((si >> 32) & 0xFF)
		b[5] = byte((si >> 40) & 0xFF)
		b[6] = byte((si >> 48) & 0xFF)
		b[7] = byte((si >> 56) & 0xFF)

	case *I64:
		b = make([]byte, 8)
		i, _ := strconv.ParseInt(v.Name(), 0, 6)
		si := uint64(int64(i))
		b[0] = byte(si & 0xFF)
		b[1] = byte((si >> 8) & 0xFF)
		b[2] = byte((si >> 16) & 0xFF)
		b[3] = byte((si >> 24) & 0xFF)
		b[4] = byte((si >> 32) & 0xFF)
		b[5] = byte((si >> 40) & 0xFF)
		b[6] = byte((si >> 48) & 0xFF)
		b[7] = byte((si >> 56) & 0xFF)

	case *F32:
		b = make([]byte, 4)
		f, _ := strconv.ParseFloat(v.Name(), 32)
		si := math.Float32bits(float32(f))
		b[0] = byte(si & 0xFF)
		b[1] = byte((si >> 8) & 0xFF)
		b[2] = byte((si >> 16) & 0xFF)
		b[3] = byte((si >> 24) & 0xFF)

	case *F64:
		b = make([]byte, 8)
		f, _ := strconv.ParseFloat(v.Name(), 64)
		si := math.Float64bits(f)
		b[0] = byte(si & 0xFF)
		b[1] = byte((si >> 8) & 0xFF)
		b[2] = byte((si >> 16) & 0xFF)
		b[3] = byte((si >> 24) & 0xFF)
		b[4] = byte((si >> 32) & 0xFF)
		b[5] = byte((si >> 40) & 0xFF)
		b[6] = byte((si >> 48) & 0xFF)
		b[7] = byte((si >> 56) & 0xFF)

	default:
		logger.Fatalf("todo: %T", v.Type())
	}

	return
}

func (v *aBasic) emitEq(r Value) (insts []wat.Inst, ok bool) {
	if !v.Type().Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}
	insts = append(insts, v.EmitPushNoRetain()...)
	insts = append(insts, r.EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstEq(toWatType(v.Type())))

	ok = true

	return
}

func (v *aBasic) emitCompare(r Value) (insts []wat.Inst) {
	if !v.Type().Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}

	insts = append(insts, v.EmitPushNoRetain()...)
	insts = append(insts, r.EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstLt(toWatType(v.Type())))

	instLe := wat.NewInstIf(nil, nil, []wat.ValueType{wat.I32{}})

	instLe.True = append(instLe.True, wat.NewInstConst(wat.I32{}, "-1"))

	instLe.False = append(instLe.False, v.EmitPushNoRetain()...)
	instLe.False = append(instLe.False, r.EmitPushNoRetain()...)
	instLe.False = append(instLe.False, wat.NewInstGt(toWatType(v.Type())))

	insts = append(insts, instLe)
	return
}
