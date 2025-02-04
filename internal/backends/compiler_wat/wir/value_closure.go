package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
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

func (s *FnSig) String() string {
	n := "$"
	for _, i := range s.Params {
		n += i.Named()
		n += "$"
	}
	n += "$$"
	for _, r := range s.Results {
		n += r.Named()
		n += "$"
	}
	n += "$"
	return n
}

/**************************************
Closure:
**************************************/
type Closure struct {
	tCommon
	underlying  *Struct
	_fnSig      FnSig
	_fnTypeName string
	_u32        ValueType
}

func (m *Module) GenValueType_Closure(sig FnSig) *Closure {
	var closure_t Closure
	closure_t._fnSig = sig
	closure_t.name = sig.String()

	t, ok := m.findValueType(closure_t.name)
	if ok {
		return t.(*Closure)
	}

	closure_t._u32 = m.U32
	closure_t._fnTypeName = m.AddFnSig(&sig)

	closure_t.underlying = m.genInternalStruct(closure_t.name + ".underlying")
	closure_t.underlying.AppendField(m.NewStructField("fn_index", m.U32))
	closure_t.underlying.AppendField(m.NewStructField("d", m.GenValueType_Ref(m.VOID)))
	closure_t.underlying.Finish()

	m.addValueType(&closure_t)
	return &closure_t
}

func (t *Closure) Size() int            { return t.underlying.Size() }
func (t *Closure) align() int           { return t.underlying.align() }
func (t *Closure) Kind() TypeKind       { return kStruct }
func (t *Closure) OnFree() int          { return t.underlying.OnFree() }
func (t *Closure) Raw() []wat.ValueType { return t.underlying.Raw() }

func (t *Closure) Equal(u ValueType) bool {
	if ut, ok := u.(*Closure); ok {
		return t._fnSig.Equal(&ut._fnSig)
	}
	return false
}

func (t *Closure) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}

	return t.underlying.EmitLoadFromAddr(addr, offset)
}

func (t *Closure) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}

	return t.underlying.EmitLoadFromAddrNoRetain(addr, offset)
}

/**************************************
aClosure:
**************************************/
type aClosure struct {
	aStruct
	typ *Closure
}

func newValue_Closure(name string, kind ValueKind, typ *Closure) *aClosure {
	var v aClosure
	v.typ = typ
	v.aStruct = *newValue_Struct(name, kind, typ.underlying)
	return &v
}

func (v *aClosure) Type() ValueType { return v.typ }

func (m *Module) GenConstFnValue(fn_name string, sig FnSig) Value {
	fn_index := m.AddTableElem(fn_name)

	closure_t := m.GenValueType_Closure(sig)
	aClosure := newValue_Closure("0", ValueKindConst, closure_t)
	aClosure.aStruct.setFieldConstValue("fn_index", NewConst(strconv.Itoa(fn_index), m.U32))

	return aClosure
}

func EmitCallClosure(c Value, params []Value) (insts []wat.Inst) {
	for _, p := range params {
		insts = append(insts, p.EmitPushNoRetain()...)
	}
	closure := c.(*aClosure)
	insts = append(insts, closure.ExtractByName("fn_index").EmitPush()...)

	insts = append(insts, closure.ExtractByName("d").(*aRef).ExtractByName("d").EmitPushNoRetain()...)
	insts = append(insts, currentModule.FindGlobalByName("$wa.runtime.closure_data").EmitPop()...)

	insts = append(insts, wat.NewInstCallIndirect(closure.typ._fnTypeName))
	return
}

func (v *aClosure) emitEq(r Value) ([]wat.Inst, bool) {
	if r.Kind() != ValueKindConst {
		return nil, false
	}

	r_c, ok := r.(*aClosure)
	if !ok {
		logger.Fatal("r is not a Closure")
	}

	insts, ok := v.ExtractByName("fn_index").emitEq(r_c.ExtractByName("fn_index"))
	if !ok {
		logger.Fatal("fn_index is not comparable")
	}
	return insts, ok
}

func (v *aClosure) emitCompare(r Value) (insts []wat.Inst) {
	if !v.Type().Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}
	d := r.(*aClosure)

	return v.aStruct.emitCompare(&d.aStruct)
}
