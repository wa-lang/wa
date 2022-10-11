package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

/**************************************
Struct:
**************************************/
type Struct struct {
	name    string
	Members []Field
	_size   int
	_align  int
}

type Field struct {
	name   string
	typ    ValueType
	_start int
}

func NewField(n string, t ValueType) Field { return Field{name: n, typ: t} }
func (i Field) Name() string               { return i.name }
func (i Field) Type() ValueType            { return i.typ }
func (i Field) Equal(u Field) bool         { return i.name == u.name && i.typ.Equal(u.typ) }

func makeAlign(i, a int) int {
	return (i + a - 1) / a * a
}

func NewStruct(name string, members []Field) Struct {
	var s Struct
	s.name = name

	for _, m := range members {
		ma := m.Type().align()
		m._start = makeAlign(s._size, ma)
		s.Members = append(s.Members, m)

		s._size = m._start + m.Type().size()
		if ma > s._align {
			s._align += ma
		}
	}
	s._size = makeAlign(s._size, s._align)

	return s
}

func (t Struct) Name() string { return t.name }
func (t Struct) size() int    { return t._size }
func (t Struct) align() int   { return t._align }

func (t Struct) onFree(m *Module) int {
	logger.Fatal("Todo")
	return 0
}

func (t Struct) Raw() []wat.ValueType {
	var r []wat.ValueType
	for _, f := range t.Members {
		r = append(r, f.Type().Raw()...)
	}
	return r
}

func (t Struct) Equal(u ValueType) bool {
	if u, ok := u.(Struct); ok {
		if len(t.Members) != len(u.Members) {
			return false
		}

		for i := range t.Members {
			if !t.Members[i].Equal(u.Members[i]) {
				return false
			}
		}

		return true
	}
	return false
}

func (t Struct) emitLoadFromAddr(addr Value, offset int) (insts []wat.Inst) {
	for _, m := range t.Members {
		ptr := newVarPointer(addr.Name(), addr.Kind(), m.Type())
		insts = append(insts, m.Type().emitLoadFromAddr(ptr, m._start+offset)...)
	}
	return
}

/**************************************
VarStruct:
**************************************/
type VarStruct struct {
	aVar
}

func newVarStruct(name string, kind ValueKind, typ ValueType) *VarStruct {
	return &VarStruct{aVar: aVar{name: name, kind: kind, typ: typ}}
}

func (v *VarStruct) raw() []wat.Value {
	var r []wat.Value
	st := v.Type().(Struct)
	for _, m := range st.Members {
		t := newValue(v.Name()+"."+m.Name(), v.kind, m.Type())
		r = append(r, t.raw()...)
	}
	return r
}

func (v *VarStruct) EmitInit() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(Struct)
	for _, m := range st.Members {
		t := newValue(v.Name()+"."+m.Name(), v.kind, m.Type())
		insts = append(insts, t.EmitInit()...)
	}
	return insts
}

func (v *VarStruct) EmitPush() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(Struct)
	for _, m := range st.Members {
		t := newValue(v.Name()+"."+m.Name(), v.kind, m.Type())
		insts = append(insts, t.EmitPush()...)
	}
	return insts
}

func (v *VarStruct) EmitPop() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(Struct)
	for i := range st.Members {
		m := st.Members[len(st.Members)-i-1]
		t := newValue(v.Name()+"."+m.Name(), v.kind, m.Type())
		insts = append(insts, t.EmitPop()...)
	}
	return insts
}

func (v *VarStruct) EmitRelease() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(Struct)
	for i := range st.Members {
		m := st.Members[len(st.Members)-i-1]
		t := newValue(v.Name()+"."+m.Name(), v.kind, m.Type())
		insts = append(insts, t.EmitRelease()...)
	}
	return insts
}

func (v *VarStruct) Extract(member_name string) Value {
	st := v.Type().(Struct)
	for _, m := range st.Members {
		if m.Name() == member_name {
			return newValue(v.Name()+"."+m.Name(), v.kind, m.Type())
		}
	}
	return nil
}

func (v *VarStruct) emitStoreToAddr(addr Value, offset int) (insts []wat.Inst) {
	st := v.Type().(Struct)
	for i := range st.Members {
		m := st.Members[len(st.Members)-i-1]
		t := newValue(v.Name()+"."+m.Name(), v.kind, m.Type())
		a := newVarPointer(addr.Name(), addr.Kind(), m.Type())
		insts = append(insts, t.emitStoreToAddr(a, m._start+offset)...)
	}
	return
}
