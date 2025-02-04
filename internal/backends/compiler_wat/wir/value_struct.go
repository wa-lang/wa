package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
StructField:
**************************************/
type StructField struct {
	name     string
	id       int
	typ      ValueType
	_start   int
	_typ_ptr *Ptr //type of *typ
}

func (m *Module) NewStructField(name string, typ ValueType) *StructField {
	f := StructField{name: name, typ: typ}
	f._typ_ptr = m.GenValueType_Ptr(typ)
	return &f
}

func (i *StructField) Name() string              { return i.name }
func (i *StructField) Type() ValueType           { return i.typ }
func (i *StructField) Equal(u *StructField) bool { return i.name == u.name && i.typ.Equal(u.typ) }

/**************************************
Struct:
**************************************/
type Struct struct {
	tCommon
	fields []*StructField
	_size  int
	_align int
	_u32   ValueType
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

func (m *Module) GenValueType_Struct(name string) (*Struct, bool) {
	t, ok := m.findValueType(name)
	if ok {
		return t.(*Struct), true
	}

	var struct_type Struct
	struct_type.name = name
	struct_type._u32 = m.U32
	m.addValueType(&struct_type)

	return &struct_type, false
}

func (m *Module) genInternalStruct(name string) *Struct {
	var struct_type Struct
	struct_type.name = name
	struct_type._u32 = m.U32

	return &struct_type
}

func (t *Struct) Size() int      { return t._size }
func (t *Struct) align() int     { return t._align }
func (t *Struct) Kind() TypeKind { return kStruct }

func (t *Struct) AppendField(f *StructField) {
	t.fields = append(t.fields, f)
}

func (t *Struct) Finish() {
	t._size = 0
	for i, field := range t.fields {
		field.id = i
		fa := field.Type().align()
		field._start = makeAlign(t._size, fa)

		t._size = field._start + field.Type().Size()
		if fa > t._align {
			t._align = fa
		}
	}
	t._size = makeAlign(t._size, t._align)
}

type fn_offset_pair struct {
	fn     int
	offset int
}

func (t *Struct) genRawFree() (ret []fn_offset_pair) {
	for _, member := range t.fields {
		member_type := member.Type()
		if istruct, ok := member_type.(iStruct); ok {
			rfs := istruct.genRawFree()
			for _, rf := range rfs {
				ret = append(ret, fn_offset_pair{fn: rf.fn, offset: rf.offset + member._start})
			}
		} else {
			mff := member_type.OnFree()
			if mff != 0 {
				ret = append(ret, fn_offset_pair{fn: mff, offset: member._start})
			}
		}
	}

	return
}

func (t *Struct) OnFree() int {
	var f Function
	f.InternalName = "$" + GenSymbolName(t.Named()) + ".$$OnFree"

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
		f.Insts = append(f.Insts, wat.NewInstCallIndirect("$OnFree"))
	}
	currentModule.AddFunc(&f)
	return currentModule.AddTableElem(f.InternalName)
}

func (t *Struct) Raw() []wat.ValueType {
	var r []wat.ValueType
	for _, f := range t.fields {
		r = append(r, f.Type().Raw()...)
	}
	return r
}

func (t *Struct) Equal(u ValueType) bool {
	ut, ok := u.(*Struct)
	if !ok {
		return false
	}

	if len(t.fields) != len(ut.fields) {
		return false
	}

	for i := range t.fields {
		if !t.fields[i].Equal(ut.fields[i]) {
			return false
		}
	}

	return true
}

func (t *Struct) EmitLoadFromAddr(addr Value, offset int) (insts []wat.Inst) {
	for _, m := range t.fields {
		ptr := newValue_Ptr(addr.Name(), addr.Kind(), m._typ_ptr)
		insts = append(insts, m.Type().EmitLoadFromAddr(ptr, m._start+offset)...)
	}
	return
}

func (t *Struct) EmitLoadFromAddrNoRetain(addr Value, offset int) (insts []wat.Inst) {
	for _, m := range t.fields {
		ptr := newValue_Ptr(addr.Name(), addr.Kind(), m._typ_ptr)
		insts = append(insts, m.Type().EmitLoadFromAddrNoRetain(ptr, m._start+offset)...)
	}
	return
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

func (v *aStruct) genSubValue(m *StructField) Value {
	if v.Kind() != ValueKindConst {
		return newValue(v.Name()+"."+strconv.Itoa(m.id), v.Kind(), m.Type())
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
	for _, m := range v.typ.fields {
		t := v.genSubValue(m)
		r = append(r, t.raw()...)
	}
	return r
}

func (v *aStruct) EmitInit() []wat.Inst {
	var insts []wat.Inst
	for _, m := range v.typ.fields {
		t := v.genSubValue(m)
		insts = append(insts, t.EmitInit()...)
	}
	return insts
}

func (v *aStruct) EmitPush() (insts []wat.Inst) {
	for _, m := range v.typ.fields {
		t := v.genSubValue(m)
		insts = append(insts, t.EmitPush()...)
	}
	return
}

func (v *aStruct) EmitPushNoRetain() (insts []wat.Inst) {
	for _, m := range v.typ.fields {
		t := v.genSubValue(m)
		insts = append(insts, t.EmitPushNoRetain()...)
	}
	return
}

func (v *aStruct) EmitPop() []wat.Inst {
	var insts []wat.Inst
	for i := range v.typ.fields {
		m := v.typ.fields[len(v.typ.fields)-i-1]
		t := v.genSubValue(m)
		insts = append(insts, t.EmitPop()...)
	}
	return insts
}

func (v *aStruct) EmitRelease() []wat.Inst {
	var insts []wat.Inst
	for i := range v.typ.fields {
		m := v.typ.fields[len(v.typ.fields)-i-1]
		t := v.genSubValue(m)
		insts = append(insts, t.EmitRelease()...)
	}
	return insts
}

func (v *aStruct) ExtractByName(member_name string) Value {
	for _, m := range v.typ.fields {
		if m.Name() == member_name {
			return v.genSubValue(m)
		}
	}
	logger.Fatal("field not found:", member_name)
	return nil
}

func (v *aStruct) ExtractByID(id int) Value {
	if id >= len(v.typ.fields) {
		logger.Fatal("id > len(fields)")
	}

	return v.genSubValue(v.typ.fields[id])
}

func (v *aStruct) emitStoreToAddr(addr Value, offset int) (insts []wat.Inst) {
	for _, m := range v.typ.fields {
		t := v.genSubValue(m)
		ptr := newValue_Ptr(addr.Name(), addr.Kind(), m._typ_ptr)
		insts = append(insts, t.emitStoreToAddr(ptr, m._start+offset)...)
	}
	return
}

func (v *aStruct) emitStore(offset int) (insts []wat.Inst) {
	for _, m := range v.typ.fields {
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
	for _, m := range v.typ.fields {
		d := b[m._start:]
		copy(d, v.genSubValue(m).Bin())
	}

	return
}

func (v *aStruct) emitEq(r Value) (insts []wat.Inst, ok bool) {
	if !v.typ.Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}

	d := r.(*aStruct)
	for i := range v.typ.fields {
		t1 := v.genSubValue(v.typ.fields[i])
		t2 := d.genSubValue(d.typ.fields[i])

		if ins, o := t1.emitEq(t2); !o {
			return nil, false
		} else {
			insts = append(insts, ins...)
		}

		if i > 0 {
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		}
	}

	ok = true

	return
}

func (v *aStruct) emitCompare(r Value) (insts []wat.Inst) {
	if !v.typ.Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}

	block := wat.NewInstBlock("")
	block.Ret = append(block.Ret, wat.I32{})

	d := r.(*aStruct)
	for i := range v.typ.fields {
		t1 := v.genSubValue(v.typ.fields[i])
		t2 := d.genSubValue(d.typ.fields[i])

		if i > 0 {
			block.Insts = append(block.Insts, wat.NewInstCall("runtime.DupI32"))
			block.Insts = append(block.Insts, wat.NewInstBrIf(0))
			block.Insts = append(block.Insts, wat.NewInstDrop())
		}

		block.Insts = append(block.Insts, t1.emitCompare(t2)...)
	}

	insts = append(insts, block)
	return
}
