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
	tCommon
	underlying *Struct
	Fields     []ValueType
}

func (m *Module) GenValueType_Tuple(fields []ValueType) *Tuple {
	tuple_t := Tuple{Fields: fields}
	tuple_t.name = "$"
	for _, t := range fields {
		tuple_t.name += t.Named()
	}

	t, ok := m.findValueType(tuple_t.name)
	if ok {
		return t.(*Tuple)
	}

	tuple_t.underlying = m.genInternalStruct(tuple_t.name + ".underlying")
	for i, t := range fields {
		name := "m" + strconv.Itoa(i)
		tuple_t.underlying.AppendField(m.NewStructField(name, t))
	}
	tuple_t.underlying.Finish()

	m.addValueType(&tuple_t)
	return &tuple_t
}

func (t *Tuple) Size() int            { return t.underlying.Size() }
func (t *Tuple) align() int           { return t.underlying.align() }
func (t *Tuple) Kind() TypeKind       { return kTuple }
func (t *Tuple) OnFree() int          { return t.underlying.OnFree() }
func (t *Tuple) Raw() []wat.ValueType { return t.underlying.Raw() }
func (t *Tuple) Equal(u ValueType) bool {
	ut, ok := u.(*Tuple)
	if !ok {
		return false
	}

	if len(t.Fields) != len(ut.Fields) {
		return false
	}

	for i := range t.Fields {
		if !t.Fields[i].Equal(ut.Fields[i]) {
			return false
		}
	}

	return true
}

func (t *Tuple) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.underlying.EmitLoadFromAddr(addr, offset)
}

func (t *Tuple) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	return t.underlying.EmitLoadFromAddrNoRetain(addr, offset)
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
	v.aStruct = *newValue_Struct(name, kind, typ.underlying)
	return &v
}

func (v *aTuple) Type() ValueType { return v.typ }

func (v *aTuple) raw() []wat.Value        { return v.aStruct.raw() }
func (v *aTuple) EmitInit() []wat.Inst    { return v.aStruct.EmitInit() }
func (v *aTuple) EmitPush() []wat.Inst    { return v.aStruct.EmitPush() }
func (v *aTuple) EmitPop() []wat.Inst     { return v.aStruct.EmitPop() }
func (v *aTuple) EmitRelease() []wat.Inst { return v.aStruct.EmitRelease() }
