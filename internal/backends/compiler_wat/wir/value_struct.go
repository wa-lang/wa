package wir

import (
	"strconv"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
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

type iStruct interface {
	ValueType
	genRawFree() (ret []fn_offset_pair)
}

type Field struct {
	name      string
	typ       ValueType
	const_val Value
	_start    int
}

func NewField(n string, t ValueType) Field { return Field{name: n, typ: t} }
func (i Field) Name() string               { return i.name }
func (i Field) Type() ValueType            { return i.typ }
func (i Field) Equal(u Field) bool         { return i.name == u.name && i.typ.Equal(u.typ) }

func makeAlign(i, a int) int {
	if a == 1 || a == 0 {
		return i
	}
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

type fn_offset_pair struct {
	fn     int
	offset int
}

func (t Struct) genRawFree() (ret []fn_offset_pair) {
	for _, member := range t.Members {
		member_type := member.Type()
		if istruct, ok := member_type.(iStruct); ok {
			rfs := istruct.genRawFree()
			for _, rf := range rfs {
				ret = append(ret, fn_offset_pair{fn: rf.fn, offset: rf.offset + member._start})
			}
		} else {
			mff := member_type.onFree()
			if mff != 0 {
				ret = append(ret, fn_offset_pair{fn: mff, offset: member._start})
			}
		}
	}

	return
}

func (t Struct) onFree() int {
	var f Function
	f.Name = "$" + t.Name() + ".$$onFree"

	if i := currentModule.findTableElem(f.Name); i != 0 {
		return i
	}

	f.Result = VOID{}
	ptr := NewLocal("$ptr", I32{})
	f.Params = append(f.Params, ptr)

	rfs := t.genRawFree()
	if len(rfs) == 0 {
		return 0
	}
	for _, rf := range rfs {
		f.Insts = append(f.Insts, ptr.EmitPush()...)
		if rf.offset != 0 {
			f.Insts = append(f.Insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(rf.offset)))
			f.Insts = append(f.Insts, wat.NewInstAdd(wat.I32{}))
		}

		f.Insts = append(f.Insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(rf.fn)))
		f.Insts = append(f.Insts, wat.NewInstCallIndirect("$onFree"))
	}
	return currentModule.addTableFunc(&f)

	/*	has_free := false
		for _, member := range t.Members {
			member_free_func := member.Type().onFree(module)
			if member_free_func == 0 {
				continue
			}

			f.Insts = append(f.Insts, ptr.EmitPush()...)
			f.Insts = append(f.Insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(member._start)))
			f.Insts = append(f.Insts, wat.NewInstAdd(wat.I32{}))

			f.Insts = append(f.Insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(member_free_func)))
			f.Insts = append(f.Insts, wat.NewInstCallIndirect("$onFree"))
			has_free = true
		}

		if has_free {
			return module.addTableFunc(&f)
		}

		return 0  //*/
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
		ptr := newValuePointer(addr.Name(), addr.Kind(), m.Type())
		insts = append(insts, m.Type().emitLoadFromAddr(ptr, m._start+offset)...)
	}
	return
}

func (t Struct) findFieldByName(field_name string) *Field {
	for i := range t.Members {
		if t.Members[i].Name() == field_name {
			return &t.Members[i]
		}
	}
	return nil
}

/**************************************
aStruct:
**************************************/
type aStruct struct {
	aValue
}

func newValueStruct(name string, kind ValueKind, typ ValueType) *aStruct {
	return &aStruct{aValue: aValue{name: name, kind: kind, typ: typ}}
}

func (v *aStruct) genSubValue(m Field) Value {
	if v.Kind() != ValueKindConst {
		return newValue(v.Name()+"."+m.Name(), v.Kind(), m.Type())
	} else {
		if m.const_val != nil {
			return m.const_val
		} else {
			return newValue(v.Name(), v.Kind(), m.Type())
		}
	}
}

func (v *aStruct) raw() []wat.Value {
	var r []wat.Value
	st := v.Type().(Struct)
	for _, m := range st.Members {
		t := v.genSubValue(m)
		r = append(r, t.raw()...)
	}
	return r
}

func (v *aStruct) EmitInit() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(Struct)
	for _, m := range st.Members {
		t := v.genSubValue(m)
		insts = append(insts, t.EmitInit()...)
	}
	return insts
}

func (v *aStruct) EmitPush() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(Struct)
	for _, m := range st.Members {
		t := v.genSubValue(m)
		insts = append(insts, t.EmitPush()...)
	}
	return insts
}

func (v *aStruct) EmitPop() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(Struct)
	for i := range st.Members {
		m := st.Members[len(st.Members)-i-1]
		t := v.genSubValue(m)
		insts = append(insts, t.EmitPop()...)
	}
	return insts
}

func (v *aStruct) EmitRelease() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(Struct)
	for i := range st.Members {
		m := st.Members[len(st.Members)-i-1]
		t := v.genSubValue(m)
		insts = append(insts, t.EmitRelease()...)
	}
	return insts
}

func (v *aStruct) Extract(member_name string) Value {
	st := v.Type().(Struct)
	for _, m := range st.Members {
		if m.Name() == member_name {
			return v.genSubValue(m)
		}
	}
	return nil
}

func (v *aStruct) emitStoreToAddr(addr Value, offset int) (insts []wat.Inst) {
	st := v.Type().(Struct)
	for i := range st.Members {
		m := st.Members[len(st.Members)-i-1]
		t := v.genSubValue(m)
		a := newValuePointer(addr.Name(), addr.Kind(), m.Type())
		insts = append(insts, t.emitStoreToAddr(a, m._start+offset)...)
	}
	return
}
