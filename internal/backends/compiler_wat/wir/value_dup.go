// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
Dup:
**************************************/
type Dup struct {
	tCommon
	name string
	Base ValueType
}

func (m *Module) GenValueType_Dup(name string, base ValueType) *Dup {
	dup_t := Dup{name: name, Base: base}
	t, ok := m.findValueType(dup_t.Name())
	if ok {
		return t.(*Dup)
	}

	m.addValueType(&dup_t)
	return &dup_t
}

func (t *Dup) Name() string         { return t.name }
func (t *Dup) Size() int            { return t.Base.Size() }
func (t *Dup) align() int           { return t.Base.align() }
func (t *Dup) Kind() TypeKind       { return kDup }
func (t *Dup) onFree() int          { return t.Base.onFree() }
func (t *Dup) Raw() []wat.ValueType { return t.Base.Raw() }
func (t *Dup) Equal(u ValueType) bool {
	if ut, ok := u.(*Dup); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}

func (t *Dup) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.Base.EmitLoadFromAddr(addr, offset)
}

/**************************************
aDup:
**************************************/
type aDup struct {
	underlying Value
	typ        *Dup
}

func newValue_Dup(name string, kind ValueKind, typ *Dup) *aDup {
	var v aDup
	v.typ = typ
	v.underlying = newValue(name, kind, typ.Base)
	return &v
}

func (v *aDup) Name() string    { return v.underlying.Name() }
func (v *aDup) Kind() ValueKind { return v.underlying.Kind() }
func (v *aDup) Type() ValueType { return v.typ }

func (v *aDup) raw() []wat.Value        { return v.underlying.raw() }
func (v *aDup) EmitInit() []wat.Inst    { return v.underlying.EmitInit() }
func (v *aDup) EmitPush() []wat.Inst    { return v.underlying.EmitPush() }
func (v *aDup) EmitPop() []wat.Inst     { return v.underlying.EmitPop() }
func (v *aDup) EmitRelease() []wat.Inst { return v.underlying.EmitRelease() }

func (v *aDup) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	return v.underlying.emitStoreToAddr(addr, offset)
}

func (v *aDup) emitStore(offset int) []wat.Inst {
	return v.underlying.emitStore(offset)
}

func (v *aDup) Bin() []byte {
	return v.underlying.Bin()
}

func (v *aDup) emitEq(r Value) (insts []wat.Inst, ok bool) {
	if !v.Type().Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}
	return v.underlying.emitEq(r.(*aDup).underlying)
}
