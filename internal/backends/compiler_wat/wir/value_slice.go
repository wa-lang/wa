// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
)

/**************************************
Slice:
**************************************/
type Slice struct {
	Base ValueType
	Struct
}

func NewSlice(base ValueType) Slice {
	var v Slice
	v.Base = base
	var m []Field
	m = append(m, NewField("block", NewBlock(base)))
	m = append(m, NewField("data", NewPointer(base)))
	m = append(m, NewField("len", I32{}))
	m = append(m, NewField("cap", I32{}))
	v.Struct = NewStruct(base.Name()+".$$slice", m)
	return v
}
func (t Slice) Name() string         { return t.Base.Name() + ".$$slice" }
func (t Slice) size() int            { return 16 }
func (t Slice) align() int           { return 4 }
func (t Slice) onFree(m *Module) int { return t.Struct.onFree(m) }
func (t Slice) Raw() []wat.ValueType { return t.Struct.Raw() }
func (t Slice) Equal(u ValueType) bool {
	if ut, ok := u.(Slice); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}

func (t Slice) emitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.Struct.emitLoadFromAddr(addr, offset)
}

/**************************************
aSlice:
**************************************/
type aSlice struct {
	aValue
	underlying aStruct
}

func newValueSlice(name string, kind ValueKind, base_type ValueType) *aSlice {
	var v aSlice
	slice_type := NewSlice(base_type)
	v.aValue = aValue{name: name, kind: kind, typ: slice_type}
	v.underlying = *newValueStruct(name, kind, slice_type.Struct)
	return &v
}

func (v *aSlice) raw() []wat.Value        { return v.underlying.raw() }
func (v *aSlice) EmitInit() []wat.Inst    { return v.underlying.EmitInit() }
func (v *aSlice) EmitPush() []wat.Inst    { return v.underlying.EmitPush() }
func (v *aSlice) EmitPop() []wat.Inst     { return v.underlying.EmitPop() }
func (v *aSlice) EmitRelease() []wat.Inst { return v.underlying.EmitRelease() }

func (v *aSlice) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	return v.underlying.emitStoreToAddr(addr, offset)
}
