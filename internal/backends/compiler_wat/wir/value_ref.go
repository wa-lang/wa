// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
Ref:
**************************************/
type Ref struct {
	tCommon
	Base        ValueType
	underlying  *Struct
	_base_block *Block
	_void       ValueType
}

func (m *Module) GenValueType_Ref(base ValueType) *Ref {
	ref_t := Ref{Base: base}
	t, ok := m.findValueType(ref_t.Name())
	if ok {
		return t.(*Ref)
	}

	ref_t._base_block = m.GenValueType_Block(base)
	ref_t._void = m.VOID
	base_ptr := m.GenValueType_Ptr(base)

	var found bool
	ref_t.underlying, found = m.GenValueType_Struct(ref_t.Name() + ".underlying")
	if found {
		logger.Fatalf("Type: %s already registered.", ref_t.Name()+".underlying")
	}

	ref_t.underlying.AppendField(m.NewStructField("block", ref_t._base_block))
	ref_t.underlying.AppendField(m.NewStructField("data", base_ptr))
	ref_t.underlying.Finish()

	m.addValueType(&ref_t)
	return &ref_t
}

func (t *Ref) Name() string         { return t.Base.Name() + ".$ref" }
func (t *Ref) Size() int            { return t.underlying.Size() }
func (t *Ref) align() int           { return t.underlying.align() }
func (t *Ref) onFree() int          { return t.underlying.onFree() }
func (t *Ref) Raw() []wat.ValueType { return t.underlying.Raw() }
func (t *Ref) Equal(u ValueType) bool {
	if ut, ok := u.(*Ref); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}

func (t *Ref) emitHeapAlloc() (insts []wat.Inst) {
	//insts = append(insts, wat.NewBlank())
	//insts = append(insts, wat.NewComment("Ref.emitHeapAlloc start"))

	insts = append(insts, t._base_block.emitHeapAlloc(NewConst("1", t.underlying._u32))...)
	insts = append(insts, wat.NewInstCall("$wa.RT.DupI32"))
	insts = append(insts, NewConst("16", t.underlying._u32).EmitPush()...)
	insts = append(insts, wat.NewInstAdd(wat.U32{}))

	//insts = append(insts, wat.NewComment("Ref.emitHeapAlloc end"))
	//insts = append(insts, wat.NewBlank())

	return
}

func (t *Ref) emitStackAlloc() (insts []wat.Inst) {
	//insts = append(insts, wat.NewBlank())
	//insts = append(insts, wat.NewComment("Ref.emitStackAlloc start"))

	logger.Fatal("Todo")

	insts = append(insts, NewConst("0", t.underlying._u32).EmitPush()...)
	insts = append(insts, NewConst(strconv.Itoa(t.Base.Size()), t.underlying._u32).EmitPush()...)
	insts = append(insts, wat.NewInstCall("$waStackAlloc"))

	//insts = append(insts, wat.NewComment("Ref.emitStackAlloc end"))
	//insts = append(insts, wat.NewBlank())
	return
}

func (t *Ref) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.underlying.EmitLoadFromAddr(addr, offset)
}

/**************************************
aRef:
**************************************/
type aRef struct {
	aStruct
	typ *Ref
}

func newValue_Ref(name string, kind ValueKind, typ *Ref) *aRef {
	var v aRef
	v.typ = typ
	v.aStruct = *newValue_Struct(name, kind, typ.underlying)
	return &v
}

func (v *aRef) Type() ValueType { return v.typ }

func (v *aRef) raw() []wat.Value        { return v.aStruct.raw() }
func (v *aRef) EmitInit() []wat.Inst    { return v.aStruct.EmitInit() }
func (v *aRef) EmitPush() []wat.Inst    { return v.aStruct.EmitPush() }
func (v *aRef) EmitPop() []wat.Inst     { return v.aStruct.EmitPop() }
func (v *aRef) EmitRelease() []wat.Inst { return v.aStruct.EmitRelease() }

func (v *aRef) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	return v.aStruct.emitStoreToAddr(addr, offset)
}

func (v *aRef) emitGetValue() []wat.Inst {
	return v.typ.Base.EmitLoadFromAddr(v.aStruct.Extract("data"), 0)
}

func (v *aRef) emitSetValue(d Value) []wat.Inst {
	if !d.Type().Equal(v.typ.Base) && !v.typ.Base.Equal(v.typ._void) {
		logger.Fatal("Type not match")
		return nil
	}
	return d.emitStoreToAddr(v.aStruct.Extract("data"), 0)
}
