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
	*Struct
	_fields []ValueType
}

func (m *Module) GenValueType_Tuple(fields []ValueType) *Tuple {
	tuple_t := Tuple{_fields: fields}
	t, ok := m.findValueType(tuple_t.Name())
	if ok {
		return t.(*Tuple)
	}

	var members []Field
	for i, t := range fields {
		fname := "m" + strconv.Itoa(i)
		members = append(members, NewField(fname, t))
	}
	tuple_t.Struct = m.GenValueType_Struct(tuple_t.Name()+".underlying", members)
	m.regValueType(&tuple_t)
	return &tuple_t
}

func (t *Tuple) Name() string {
	s := "$"
	for _, t := range t._fields {
		s += t.Name()
	}
	return s
}

func (t *Tuple) size() int            { return t.Struct.size() }
func (t *Tuple) align() int           { return t.Struct.align() }
func (t *Tuple) onFree() int          { return t.Struct.onFree() }
func (t *Tuple) Raw() []wat.ValueType { return t.Struct.Raw() }
func (t *Tuple) Equal(u ValueType) bool {
	ut, ok := u.(*Tuple)
	if !ok {
		return false
	}
	return t.Struct.Equal(ut.Struct)
}

func (t *Tuple) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.Struct.EmitLoadFromAddr(addr, offset)
}

/**************************************
aTuple:
**************************************/
type aTuple struct {
	aStruct
	typ *Tuple
}

func newValue_Tuple(name string, kind ValueKind, typ *Tuple) *aTuple {
	var v aTuple
	v.typ = typ
	v.aStruct = *newValue_Struct(name, kind, typ.Struct)
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
	st := v.typ.Struct
	if id >= len(st.Members) {
		panic("id >= len(st.Members)")
	}

	return v.aStruct.genSubValue(st.Members[id])
}
