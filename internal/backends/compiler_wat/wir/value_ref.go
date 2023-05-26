// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
)

/**************************************
Ref:
**************************************/
type Ref struct {
	tCommon
	name string
	Base ValueType
}

func (m *Module) GenValueType_Ref(name string, base ValueType) *Ref {
	ref_t := Ref{name: name, Base: base}
	t, ok := m.findValueType(ref_t.Name())
	if ok {
		return t.(*Ref)
	}

	m.addValueType(&ref_t)
	return &ref_t
}

func (t *Ref) Name() string         { return t.name }
func (t *Ref) Size() int            { return t.Base.Size() }
func (t *Ref) align() int           { return t.Base.align() }
func (t *Ref) Kind() TypeKind       { return kRef }
func (t *Ref) onFree() int          { return t.Base.onFree() }
func (t *Ref) Raw() []wat.ValueType { return t.Base.Raw() }
func (t *Ref) Equal(u ValueType) bool {
	if ut, ok := u.(*SPtr); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}

func (t *Ref) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.Base.EmitLoadFromAddr(addr, offset)
}

/**************************************
aRef:
**************************************/
type aRef struct {
	underlying Value
	typ        *Ref
}

func newValue_Ref(name string, kind ValueKind, typ *Ref) *aRef {
	var v aRef
	v.typ = typ
	v.underlying = newValue(name, kind, typ.Base)
	return &v
}

func (v *aRef) Name() string    { return v.underlying.Name() }
func (v *aRef) Kind() ValueKind { return v.underlying.Kind() }
func (v *aRef) Type() ValueType { return v.typ }

func (v *aRef) raw() []wat.Value        { return v.underlying.raw() }
func (v *aRef) EmitInit() []wat.Inst    { return v.underlying.EmitInit() }
func (v *aRef) EmitPush() []wat.Inst    { return v.underlying.EmitPush() }
func (v *aRef) EmitPop() []wat.Inst     { return v.underlying.EmitPop() }
func (v *aRef) EmitRelease() []wat.Inst { return v.underlying.EmitRelease() }

func (v *aRef) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	return v.underlying.emitStoreToAddr(addr, offset)
}

func (v *aRef) emitStore(offset int) []wat.Inst {
	return v.underlying.emitStore(offset)
}

func (v *aRef) Bin() []byte {
	return v.underlying.Bin()
}
