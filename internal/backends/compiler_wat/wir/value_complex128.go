// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strings"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
Complex128:
**************************************/
type Complex128 struct {
	tCommon
	underlying *Struct
	_f64       ValueType
}

func (m *Module) GenValueType_complex128(name string) *Complex128 {
	var t Complex128
	if len(name) > 0 {
		t.name = name
	} else {
		t.name = "complex128"
	}
	ot, ok := m.findValueType(t.name)
	if ok {
		return ot.(*Complex128)
	}

	t._f64 = m.F64
	t.underlying = m.genInternalStruct(t.name + ".underlying")
	t.underlying.AppendField(m.NewStructField("real", m.F64))
	t.underlying.AppendField(m.NewStructField("imag", m.F64))
	t.underlying.Finish()

	m.addValueType(&t)
	return &t
}

func (t *Complex128) Size() int              { return t.underlying.Size() }
func (t *Complex128) align() int             { return t.underlying.align() }
func (t *Complex128) Kind() TypeKind         { return kComplex128 }
func (t *Complex128) OnFree() int            { return t.underlying.OnFree() }
func (t *Complex128) Raw() []wat.ValueType   { return t.underlying.Raw() }
func (t *Complex128) Equal(u ValueType) bool { _, ok := u.(*Complex128); return ok }

func (t *Complex128) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.underlying.EmitLoadFromAddr(addr, offset)
}

func (t *Complex128) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	return t.underlying.EmitLoadFromAddrNoRetain(addr, offset)
}

func (t *Complex128) emitAdd() (insts []wat.Inst) {
	return []wat.Inst{wat.NewInstCall("$wa.runtime.complex128_Add")}
}

func (t *Complex128) emitSub() (insts []wat.Inst) {
	return []wat.Inst{wat.NewInstCall("$wa.runtime.complex128_Sub")}
}

func (t *Complex128) emitMul() (insts []wat.Inst) {
	return []wat.Inst{wat.NewInstCall("$wa.runtime.complex128_Mul")}
}

func (t *Complex128) emitDiv() (insts []wat.Inst) {
	return []wat.Inst{wat.NewInstCall("$wa.runtime.complex128_Div")}
}

/**************************************
aComplex128:
**************************************/
type aComplex128 struct {
	aStruct
	typ *Complex128
}

func newValue_Complex128(name string, kind ValueKind, typ *Complex128) *aComplex128 {
	var v aComplex128
	v.typ = typ
	v.aStruct = *newValue_Struct(name, kind, typ.underlying)
	if kind == ValueKindConst {
		s := strings.Split(name, " ")
		if len(s) != 2 {
			logger.Fatal("Invalid const complex128.")
		}
		v.aStruct.setFieldConstValue("real", NewConst(s[0], typ._f64))
		v.aStruct.setFieldConstValue("imag", NewConst(s[1], typ._f64))
	}

	return &v
}

func (v *aComplex128) Type() ValueType { return v.typ }

func (v *aComplex128) emitEq(r Value) (insts []wat.Inst, ok bool) {
	if !v.Type().Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}
	return v.aStruct.emitEq(&r.(*aComplex128).aStruct)
}

func (v *aComplex128) emitCompare(r Value) (insts []wat.Inst) {
	if !v.Type().Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}

	return v.aStruct.emitCompare(&r.(*aComplex128).aStruct)
}
