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
		n += i.Name()
		n += "$"
	}
	n += "$$"
	for _, r := range s.Results {
		n += r.Name()
		n += "$"
	}
	n += "$"
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
	tCommon
	underlying *Struct
	_fn_type   FnType
	_u32       ValueType
}

var _closure_id int

func (m *Module) GenValueType_Closure(sig FnSig) *Closure {
	var closure_t Closure
	closure_t._fn_type.FnSig = sig
	t, ok := m.findValueType(closure_t.Name())
	if ok {
		return t.(*Closure)
	}

	closure_t._u32 = m.U32
	closure_t._fn_type.Name = "closure$" + strconv.Itoa(_closure_id)
	_closure_id++
	m.AddFnType(&closure_t._fn_type)
	var fields []Field
	fields = append(fields, NewField("fn_index", m.U32))
	fields = append(fields, NewField("data", m.GenValueType_Ref(m.VOID)))
	closure_t.underlying = m.GenValueType_Struct(closure_t.Name()+".underlying", fields)

	m.addValueType(&closure_t)
	return &closure_t
}

func (t *Closure) Name() string         { return t._fn_type.String() }
func (t *Closure) Size() int            { return t.underlying.Size() }
func (t *Closure) align() int           { return t.underlying.align() }
func (t *Closure) onFree() int          { return t.underlying.onFree() }
func (t *Closure) Raw() []wat.ValueType { return t.underlying.Raw() }

func (t *Closure) Equal(u ValueType) bool {
	if ut, ok := u.(*Closure); ok {
		return t._fn_type.FnSig.Equal(&ut._fn_type.FnSig)
	}
	return false
}

func (t *Closure) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(*Ptr).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}

	return t.underlying.EmitLoadFromAddr(addr, offset)
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
	fn_index := currentModule.AddTableElem(fn_name)

	closure_t := m.GenValueType_Closure(sig)
	aClosure := newValue_Closure("0", ValueKindConst, closure_t)
	aClosure.aStruct.setFieldConstValue("fn_index", NewConst(strconv.Itoa(fn_index), m.U32))

	return aClosure
}

func EmitCallClosure(c Value, params []Value) (insts []wat.Inst) {
	for _, p := range params {
		insts = append(insts, p.EmitPush()...)
	}
	closure := c.(*aClosure)
	insts = append(insts, closure.Extract("fn_index").EmitPush()...)

	insts = append(insts, closure.Extract("data").EmitPush()...)
	insts = append(insts, currentModule.FindGlobalByName("$wa.RT.closure_data").EmitPop()...)

	insts = append(insts, wat.NewInstCallIndirect(closure.typ._fn_type.Name))
	return
}
