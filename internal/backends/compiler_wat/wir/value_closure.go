package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

/**************************************
FnType:
**************************************/
type FnType struct {
	Name    string
	Params  []ValueType
	Results []ValueType
}

func (s *FnType) Equal(u *FnType) bool {
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
Closure:
**************************************/
type Closure struct {
	Struct
	Sig FnType
}

func NewClosure(sig FnType) Closure {
	var v Closure
	var fields []Field
	fields = append(fields, NewField("fn_index", U32{}))
	fields = append(fields, NewField("data", NewRef(VOID{})))
	v.Struct = NewStruct("closure$", fields)
	v.Sig = sig

	return v
}

func (t Closure) Name() string { return "closure$" }

func (t Closure) Equal(u ValueType) bool {
	if ut, ok := u.(Closure); ok {
		return t.Sig.Equal(&ut.Sig)
	}
	return false
}

func (t Closure) emitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}

	return t.Struct.emitLoadFromAddr(addr, offset)
}

/**************************************
aClosure:
**************************************/
type aClosure struct {
	aStruct
	typ Closure
}

func newValueClosure(name string, kind ValueKind, sig FnType) *aClosure {
	var v aClosure
	v.typ = NewClosure(sig)
	v.aStruct = *newValueStruct(name, kind, v.typ.Struct)
	return &v
}

func (v *aClosure) Type() ValueType { return v.typ }
