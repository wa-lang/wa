package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
Field:
**************************************/
type Field struct {
	name     string
	typ      ValueType
	_start   int
	_typ_ptr *Ptr //type of *typ
}

func NewField(n string, t ValueType) Field { return Field{name: n, typ: t} }
func (i *Field) Name() string              { return i.name }
func (i *Field) Type() ValueType           { return i.typ }
func (i *Field) Equal(u Field) bool        { return i.name == u.name && i.typ.Equal(u.typ) }

/**************************************
Struct:
**************************************/
type Struct struct {
	tCommon
	name    string
	Members []Field
	_size   int
	_align  int
	_u32    ValueType
}

type iStruct interface {
	ValueType
	genRawFree() (ret []fn_offset_pair)
}

func makeAlign(i, align int) int {
	if align == 1 || align == 0 {
		return i
	}
	return (i + align - 1) / align * align
}

func (m *Module) GenValueType_Struct(name string, fields []Field) *Struct {
	t, ok := m.findValueType(name)
	if ok {
		return t.(*Struct)
	}

	var struct_type Struct
	struct_type.name = name
	struct_type._u32 = m.U32
	m.addValueType(&struct_type)

	for _, f := range fields {
		ma := f.Type().align()
		f._start = makeAlign(struct_type._size, ma)
		f._typ_ptr = m.GenValueType_Ptr(f.typ)
		struct_type.Members = append(struct_type.Members, f)

		struct_type._size = f._start + f.Type().Size()
		if ma > struct_type._align {
			struct_type._align = ma
		}
	}
	struct_type._size = makeAlign(struct_type._size, struct_type._align)
	return &struct_type
}

func (t *Struct) Name() string { return t.name }
func (t *Struct) Size() int    { return t._size }
func (t *Struct) align() int   { return t._align }

type fn_offset_pair struct {
	fn     int
	offset int
}

func (t *Struct) genRawFree() (ret []fn_offset_pair) {
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

func (t *Struct) onFree() int {
	var f Function
	f.InternalName = "$" + GenSymbolName(t.Name()) + ".$$onFree"

	if i := currentModule.findTableElem(f.InternalName); i != 0 {
		return i
	}

	ptr := NewLocal("$ptr", t._u32)
	f.Params = append(f.Params, ptr)

	rfs := t.genRawFree()
	if len(rfs) == 0 {
		return 0
	}
	for _, rf := range rfs {
		f.Insts = append(f.Insts, ptr.EmitPush()...)
		if rf.offset != 0 {
			f.Insts = append(f.Insts, wat.NewInstConst(wat.U32{}, strconv.Itoa(rf.offset)))
			f.Insts = append(f.Insts, wat.NewInstAdd(wat.U32{}))
		}

		f.Insts = append(f.Insts, wat.NewInstConst(wat.U32{}, strconv.Itoa(rf.fn)))
		f.Insts = append(f.Insts, wat.NewInstCallIndirect("$onFree"))
	}
	currentModule.AddFunc(&f)
	return currentModule.AddTableElem(f.InternalName)
}

func (t *Struct) Raw() []wat.ValueType {
	var r []wat.ValueType
	for _, f := range t.Members {
		r = append(r, f.Type().Raw()...)
	}
	return r
}

func (t *Struct) Equal(u ValueType) bool {
	ut, ok := u.(*Struct)
	if !ok {
		return false
	}

	if len(t.Members) != len(ut.Members) {
		return false
	}

	for i := range t.Members {
		if !t.Members[i].Equal(ut.Members[i]) {
			return false
		}
	}

	return true
}

func (t *Struct) EmitLoadFromAddr(addr Value, offset int) (insts []wat.Inst) {
	for _, m := range t.Members {
		ptr := newValue_Ptr(addr.Name(), addr.Kind(), m._typ_ptr)
		insts = append(insts, m.Type().EmitLoadFromAddr(ptr, m._start+offset)...)
	}
	return
}

func (t *Struct) findFieldByName(field_name string) *Field {
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
	typ              *Struct
	field_const_vals map[string]Value
}

func newValue_Struct(name string, kind ValueKind, typ *Struct) *aStruct {
	var v aStruct
	v.typ = typ
	v.aValue.name = name
	v.aValue.kind = kind
	v.aValue.typ = typ
	return &v
}

func (v *aStruct) genSubValue(m Field) Value {
	if v.Kind() != ValueKindConst {
		return newValue(v.Name()+"."+m.Name(), v.Kind(), m.Type())
	} else {
		fv, ok := v.field_const_vals[m.Name()]
		if ok {
			return fv
		} else {
			return newValue(v.Name(), v.Kind(), m.Type())
		}
	}
}

func (v *aStruct) setFieldConstValue(field string, sv Value) {
	if v.Kind() != ValueKindConst {
		logger.Fatal("Can't set field-val of none-const value")
	}
	if v.field_const_vals == nil {
		v.field_const_vals = make(map[string]Value)
	}
	v.field_const_vals[field] = sv
}

func (v *aStruct) raw() []wat.Value {
	var r []wat.Value
	st := v.Type().(*Struct)
	for _, m := range st.Members {
		t := v.genSubValue(m)
		r = append(r, t.raw()...)
	}
	return r
}

func (v *aStruct) EmitInit() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(*Struct)
	for _, m := range st.Members {
		t := v.genSubValue(m)
		insts = append(insts, t.EmitInit()...)
	}
	return insts
}

func (v *aStruct) EmitPush() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(*Struct)
	for _, m := range st.Members {
		t := v.genSubValue(m)
		insts = append(insts, t.EmitPush()...)
	}
	return insts
}

func (v *aStruct) EmitPop() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(*Struct)
	for i := range st.Members {
		m := st.Members[len(st.Members)-i-1]
		t := v.genSubValue(m)
		insts = append(insts, t.EmitPop()...)
	}
	return insts
}

func (v *aStruct) EmitRelease() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(*Struct)
	for i := range st.Members {
		m := st.Members[len(st.Members)-i-1]
		t := v.genSubValue(m)
		insts = append(insts, t.EmitRelease()...)
	}
	return insts
}

func (v *aStruct) Extract(member_name string) Value {
	st := v.Type().(*Struct)
	for _, m := range st.Members {
		if m.Name() == member_name {
			return v.genSubValue(m)
		}
	}
	return nil
}

func (v *aStruct) emitStoreToAddr(addr Value, offset int) (insts []wat.Inst) {
	st := v.Type().(*Struct)
	for _, m := range st.Members {
		t := v.genSubValue(m)
		ptr := newValue_Ptr(addr.Name(), addr.Kind(), m._typ_ptr)
		insts = append(insts, t.emitStoreToAddr(ptr, m._start+offset)...)
	}
	return
}

func (v *aStruct) emitStore(offset int) (insts []wat.Inst) {
	st := v.Type().(*Struct)
	for _, m := range st.Members {
		t := v.genSubValue(m)
		insts = append(insts, t.emitStore(m._start+offset)...)
	}
	return
}

func (v *aStruct) Bin() (b []byte) {
	if v.Kind() != ValueKindConst {
		panic("Value.bin(): const only!")
	}

	b = make([]byte, v.typ.Size())
	for _, m := range v.typ.Members {
		d := b[m._start:]
		copy(d, v.genSubValue(m).Bin())
	}

	return
}
