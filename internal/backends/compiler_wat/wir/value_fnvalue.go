package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

/**************************************
FnSig:
**************************************/
type FnSig struct {
	Params  []ValueType
	Results []ValueType
}

func (s *FnSig) Equal(u *FnSig) bool {
	if len(s.Params) != len(u.Params) {
		return false
	}
	for i := range s.Params {
		if !s.Params[i].Equal(u.Params[i]) {
			return false
		}
	}

	if len(s.Results) != len(u.Results) {
		return false
	}
	for i := range s.Results {
		if !s.Results[i].Equal(u.Results[i]) {
			return false
		}
	}

	return true
}

/**************************************
FnValue:
**************************************/
type FnValue struct {
	Struct
	Sig FnSig
}

func NewFnValue(sig FnSig) FnValue {
	var v FnValue
	var fields []Field
	fields = append(fields, NewField("fn_index", U32{}))
	fields = append(fields, NewField("free_var", NewRef(VOID{})))
	v.Sig = sig

	return v
}

func (t FnValue) Name() string { return "fnptr$" }

func (t FnValue) Equal(u ValueType) bool {
	if ut, ok := u.(FnValue); ok {
		return t.Sig.Equal(&ut.Sig)
	}
	return false
}

func (t FnValue) emitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}

	return t.Struct.emitLoadFromAddr(addr, offset)
}

/**************************************
aFnValue:
**************************************/
type aFnValue struct {
	aStruct
	typ FnValue
}

func newValueFnValue(name string, kind ValueKind, sig FnSig) *aFnValue {
	var v aFnValue
	v.typ = NewFnValue(sig)
	v.aStruct = *newValueStruct(name, kind, v.typ.Struct)
	return &v
}

func (v *aFnValue) Type() ValueType { return v.typ }
