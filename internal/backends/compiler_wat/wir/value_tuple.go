// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
)

/**************************************
Tuple:
**************************************/
type Tuple struct {
	Struct
}

func NewTuple(fields []ValueType) Tuple {
	var m []Field
	for i, t := range fields {
		fname := "m" + strconv.Itoa(i)
		m = append(m, NewField(fname, t))
	}
	var v Tuple
	v.Struct = NewStruct("tuple", m)
	return v
}
func (t Tuple) Name() string         { return t.Struct.Name() }
func (t Tuple) size() int            { return t.Struct.size() }
func (t Tuple) align() int           { return t.Struct.align() }
func (t Tuple) onFree() int          { return t.Struct.onFree() }
func (t Tuple) Raw() []wat.ValueType { return t.Struct.Raw() }
func (t Tuple) Equal(u ValueType) bool {
	ut, ok := u.(Tuple)
	if !ok {
		return false
	}
	return t.Struct.Equal(ut.Struct)
}

func (t Tuple) emitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.Struct.emitLoadFromAddr(addr, offset)
}

/**************************************
aTuple:
**************************************/
type aTuple struct {
	aValue
	underlying aStruct
}

func newValueTuple(name string, kind ValueKind, typ Tuple) *aTuple {
	var v aTuple
	v.aValue = aValue{name: name, kind: kind, typ: typ}
	v.underlying = *newValueStruct(name, kind, typ.Struct)
	return &v
}

func (v *aTuple) raw() []wat.Value        { return v.underlying.raw() }
func (v *aTuple) EmitInit() []wat.Inst    { return v.underlying.EmitInit() }
func (v *aTuple) EmitPush() []wat.Inst    { return v.underlying.EmitPush() }
func (v *aTuple) EmitPop() []wat.Inst     { return v.underlying.EmitPop() }
func (v *aTuple) EmitRelease() []wat.Inst { return v.underlying.EmitRelease() }

func (v *aTuple) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	return v.underlying.emitStoreToAddr(addr, offset)
}

func (v *aTuple) Extract(id int) Value {
	st := v.Type().(Tuple).Struct
	if id >= len(st.Members) {
		panic("id >= len(st.Members)")
	}

	return v.underlying.genSubValue(st.Members[id])
}
