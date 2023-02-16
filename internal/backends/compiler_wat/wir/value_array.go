// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
Array:
**************************************/
type Array struct {
	Base ValueType
	*Struct
	Capacity int
}

func (m *Module) GenValueType_Array(base ValueType, capacity int) *Array {
	arr_t := Array{Base: base, Capacity: capacity}
	t, ok := m.findValueType(arr_t.Name())
	if ok {
		return t.(*Array)
	}

	var members []Field
	for i := 0; i < capacity; i++ {
		members = append(members, NewField("m"+strconv.Itoa(i), base))
	}
	arr_t.Struct = m.GenValueType_Struct(arr_t.Name()+".underlying", members)
	m.regValueType(&arr_t)
	return &arr_t
}

func (t *Array) Name() string         { return t.Base.Name() + ".$array" + strconv.Itoa(t.Capacity) }
func (t *Array) size() int            { return t.Struct.size() }
func (t *Array) align() int           { return t.Struct.align() }
func (t *Array) Raw() []wat.ValueType { return t.Struct.Raw() }
func (t *Array) onFree() int          { return t.Struct.onFree() }

func (t *Array) Equal(u ValueType) bool {
	if ut, ok := u.(*Array); ok {
		return t.Base.Equal(ut.Base) && t.Capacity == ut.Capacity
	}
	return false
}

func (t *Array) EmitLoadFromAddr(addr Value, offset int) (insts []wat.Inst) {
	return t.Struct.EmitLoadFromAddr(addr, offset)
}

/**************************************
aArray:
**************************************/
type aArray struct {
	aStruct
	typ *Array
}

func newValue_Array(name string, kind ValueKind, typ *Array) *aArray {
	var v aArray
	v.typ = typ
	v.aStruct = *newValue_Struct(name, kind, typ.Struct)
	return &v
}

func (v *aArray) Type() ValueType { return v.typ }

func (v *aArray) raw() []wat.Value                { return v.aStruct.raw() }
func (v *aArray) EmitInit() (insts []wat.Inst)    { return v.aStruct.EmitInit() }
func (v *aArray) EmitPush() (insts []wat.Inst)    { return v.aStruct.EmitPush() }
func (v *aArray) EmitPop() (insts []wat.Inst)     { return v.aStruct.EmitPop() }
func (v *aArray) EmitRelease() (insts []wat.Inst) { return v.aStruct.EmitRelease() }

func (v *aArray) emitStoreToAddr(addr Value, offset int) (insts []wat.Inst) {
	if !addr.Type().(*Ptr).Base.Equal(v.Type()) {
		logger.Fatal("Type not match")
		return nil
	}

	return v.aStruct.emitStoreToAddr(addr, offset)
}
