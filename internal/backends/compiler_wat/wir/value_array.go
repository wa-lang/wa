// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

/**************************************
Array:
**************************************/
type Array struct {
	Base     ValueType
	Capacity int
	Struct
}

func NewArray(base ValueType, capacity int) Array {
	var v Array
	v.Base = base
	v.Capacity = capacity

	var m []Field
	for i := 0; i < capacity; i++ {
		m = append(m, NewField("m"+strconv.Itoa(i), base))
	}
	v.Struct = NewStruct(v.Name()+".underlying", m)

	return v
}
func (t Array) Name() string         { return t.Base.Name() + ".$$array" + strconv.Itoa(t.Capacity) }
func (t Array) size() int            { return t.Struct.size() }
func (t Array) align() int           { return t.Struct.align() }
func (t Array) Raw() []wat.ValueType { return t.Struct.Raw() }
func (t Array) onFree() int          { return t.Struct.onFree() }

func (t Array) Equal(u ValueType) bool {
	if ut, ok := u.(Array); ok {
		return t.Base.Equal(ut.Base) && t.Capacity == ut.Capacity
	}
	return false
}

func (t Array) emitLoadFromAddr(addr Value, offset int) (insts []wat.Inst) {
	return t.Struct.emitLoadFromAddr(addr, offset)
}

/**************************************
aArray:
**************************************/
type aArray struct {
	aValue
	underlying aStruct
}

func newValueArray(name string, kind ValueKind, base_type ValueType, capacity int) *aArray {
	var v aArray
	array_typ := NewArray(base_type, capacity)
	v.aValue = aValue{name: name, kind: kind, typ: array_typ}
	v.underlying = *newValueStruct(name, kind, array_typ.Struct)
	return &v
}

func (v *aArray) raw() []wat.Value                { return v.underlying.raw() }
func (v *aArray) EmitInit() (insts []wat.Inst)    { return v.underlying.EmitInit() }
func (v *aArray) EmitPush() (insts []wat.Inst)    { return v.underlying.EmitPush() }
func (v *aArray) EmitPop() (insts []wat.Inst)     { return v.underlying.EmitPop() }
func (v *aArray) EmitRelease() (insts []wat.Inst) { return v.underlying.EmitRelease() }

func (v *aArray) emitStoreToAddr(addr Value, offset int) (insts []wat.Inst) {
	if !addr.Type().(Pointer).Base.Equal(v.Type()) {
		logger.Fatal("Type not match")
		return nil
	}

	return v.underlying.emitStoreToAddr(addr, offset)
}
