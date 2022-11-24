package wir

import (
	"strconv"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
	"github.com/wa-lang/wa/internal/types"
)

/**************************************
FnSig:
**************************************/
type FnSig struct {
	Params  []ValueType
	Results []ValueType
}

func NewFnSigFromSignature(s *types.Signature) FnSig {
	var sig FnSig
	for i := 0; i < s.Params().Len(); i++ {
		typ := ToWType(s.Params().At(i).Type())
		sig.Params = append(sig.Params, typ)
	}
	for i := 0; i < s.Results().Len(); i++ {
		typ := ToWType(s.Results().At(i).Type())
		sig.Results = append(sig.Results, typ)
	}
	return sig
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

func (s *FnSig) String() string {
	var n string
	for _, i := range s.Params {
		n += i.Name()
	}
	n += "$$"
	for _, r := range s.Results {
		n += r.Name()
	}
	return n
}

/**************************************
FnType:
**************************************/
type FnType struct {
	Name string
	FnSig
}

/**************************************
Closure:
**************************************/
type Closure struct {
	Struct
	Sig FnSig
}

func NewClosure(sig FnSig) Closure {
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

func (t Closure) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}

	return t.Struct.EmitLoadFromAddr(addr, offset)
}

/**************************************
aClosure:
**************************************/
type aClosure struct {
	aStruct
	typ Closure
}

func newValueClosure(name string, kind ValueKind, sig FnSig) *aClosure {
	var v aClosure
	v.typ = NewClosure(sig)
	v.aStruct = *newValueStruct(name, kind, v.typ.Struct)
	return &v
}

func (v *aClosure) Type() ValueType { return v.typ }

func GenConstFnValue(fn_name string, fn_sig FnSig) Value {
	fn_index := currentModule.AddTableElem(fn_name)

	var v aClosure
	v.typ = NewClosure(fn_sig)
	v.typ.findFieldByName("fn_index").const_val = NewConst(strconv.Itoa(fn_index), U32{})
	v.aStruct = *newValueStruct("0", ValueKindConst, v.typ.Struct)

	return &v
}

func EmitCallClosure(c Value, params []Value) (insts []wat.Inst) {
	for _, p := range params {
		insts = append(insts, p.EmitPush()...)
	}
	closure := c.(*aClosure)
	insts = append(insts, closure.Extract("fn_index").EmitPush()...)

	insts = append(insts, closure.Extract("data").EmitPush()...)
	insts = append(insts, currentModule.FindGlobal("$wa.RT.closure_data").EmitPop()...)

	var fn_type FnType
	fn_type.FnSig = closure.typ.Sig
	fn_type.Name = fn_type.FnSig.String()
	currentModule.addFnType(&fn_type)

	insts = append(insts, wat.NewInstCallIndirect(fn_type.Name))
	return
}
