// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strings"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
Complex64:
**************************************/
type Complex64 struct {
	tCommon
	underlying *Struct

	_f32 ValueType
}

func (m *Module) GenValueType_complex64(name string) *Complex64 {
	var t Complex64
	if len(name) > 0 {
		t.name = name
	} else {
		t.name = "complex64"
	}
	ot, ok := m.findValueType(t.name)
	if ok {
		return ot.(*Complex64)
	}

	t._f32 = m.F32
	t.underlying = m.genInternalStruct(t.name + ".underlying")
	t.underlying.AppendField(m.NewStructField("real", m.F32))
	t.underlying.AppendField(m.NewStructField("imag", m.F32))
	t.underlying.Finish()

	m.addValueType(&t)
	return &t
}

func (t *Complex64) Size() int              { return t.underlying.Size() }
func (t *Complex64) align() int             { return t.underlying.align() }
func (t *Complex64) Kind() TypeKind         { return kComplex64 }
func (t *Complex64) OnFree() int            { return t.underlying.OnFree() }
func (t *Complex64) Raw() []wat.ValueType   { return t.underlying.Raw() }
func (t *Complex64) Equal(u ValueType) bool { _, ok := u.(*Complex64); return ok }

func (t *Complex64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.underlying.EmitLoadFromAddr(addr, offset)
}

func (t *Complex64) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	return t.underlying.EmitLoadFromAddrNoRetain(addr, offset)
}

func (t *Complex64) emitAdd() (insts []wat.Inst) {
	return []wat.Inst{wat.NewInstCall("$wa.runtime.complex64_Add")}
}

func (t *Complex64) emitSub() (insts []wat.Inst) {
	return []wat.Inst{wat.NewInstCall("$wa.runtime.complex64_Sub")}
}

func (t *Complex64) emitMul() (insts []wat.Inst) {
	return []wat.Inst{wat.NewInstCall("$wa.runtime.complex64_Mul")}
}

func (t *Complex64) emitDiv() (insts []wat.Inst) {
	return []wat.Inst{wat.NewInstCall("$wa.runtime.complex64_Div")}
}

/**************************************
aComplex64:
**************************************/
type aComplex64 struct {
	aStruct
	typ *Complex64
}

func newValue_Complex64(name string, kind ValueKind, typ *Complex64) *aComplex64 {
	var v aComplex64
	v.typ = typ
	v.aStruct = *newValue_Struct(name, kind, typ.underlying)
	if kind == ValueKindConst {
		s := strings.Split(name, " ")
		if len(s) != 2 {
			logger.Fatal("Invalid const complex64.")
		}
		v.aStruct.setFieldConstValue("real", NewConst(s[0], typ._f32))
		v.aStruct.setFieldConstValue("imag", NewConst(s[1], typ._f32))
	}

	return &v
}

func (v *aComplex64) Type() ValueType { return v.typ }

func (v *aComplex64) emitEq(r Value) (insts []wat.Inst, ok bool) {
	if !v.Type().Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}
	return v.aStruct.emitEq(&r.(*aComplex64).aStruct)
}

func (v *aComplex64) emitCompare(r Value) (insts []wat.Inst) {
	if !v.Type().Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}

	return v.aStruct.emitCompare(&r.(*aComplex64).aStruct)
}

func ComplexExtractReal(v Value) Value {
	switch v := v.(type) {
	case *aComplex64:
		return v.aStruct.ExtractByName("real")

	case *aComplex128:
		return v.aStruct.ExtractByName("real")

	default:
		logger.Fatal("Value is not a complex.")
		return nil
	}
}

func ComplexExtractImag(v Value) Value {
	switch v := v.(type) {
	case *aComplex64:
		return v.aStruct.ExtractByName("imag")

	case *aComplex128:
		return v.aStruct.ExtractByName("imag")

	default:
		logger.Fatal("Value is not a complex.")
		return nil
	}
}
