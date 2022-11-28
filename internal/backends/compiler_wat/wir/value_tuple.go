// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
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

func (t Tuple) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.Struct.EmitLoadFromAddr(addr, offset)
}

/**************************************
aTuple:
**************************************/
type aTuple struct {
	aStruct
	typ Tuple
}

func newValueTuple(name string, kind ValueKind, typ Tuple) *aTuple {
	var v aTuple
	v.typ = typ
	v.aStruct = *newValueStruct(name, kind, typ.Struct)
	return &v
}

func (v *aTuple) Type() ValueType { return v.typ }

func (v *aTuple) raw() []wat.Value        { return v.aStruct.raw() }
func (v *aTuple) EmitInit() []wat.Inst    { return v.aStruct.EmitInit() }
func (v *aTuple) EmitPush() []wat.Inst    { return v.aStruct.EmitPush() }
func (v *aTuple) EmitPop() []wat.Inst     { return v.aStruct.EmitPop() }
func (v *aTuple) EmitRelease() []wat.Inst { return v.aStruct.EmitRelease() }

func (v *aTuple) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	return v.aStruct.emitStoreToAddr(addr, offset)
}

func (v *aTuple) Extract(id int) Value {
	st := v.Type().(Tuple).Struct
	if id >= len(st.Members) {
		panic("id >= len(st.Members)")
	}

	return v.genSubValue(st.Members[id])
}
