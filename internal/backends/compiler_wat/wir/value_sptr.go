// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
SPtr:
**************************************/
type SPtr struct {
	tCommon
	Base        ValueType
	underlying  *Struct
	_base_block *Block
	_void       ValueType
}

func (m *Module) GenValueType_SPtr(base ValueType) *SPtr {
	sptr_t := SPtr{Base: base}
	t, ok := m.findValueType(sptr_t.Name())
	if ok {
		return t.(*SPtr)
	}

	sptr_t._base_block = m.GenValueType_Block(base)
	sptr_t._void = m.VOID
	base_ptr := m.GenValueType_Ptr(base)

	sptr_t.underlying = m.genInternalStruct(sptr_t.Name() + ".underlying")
	sptr_t.underlying.AppendField(m.NewStructField("block", sptr_t._base_block))
	sptr_t.underlying.AppendField(m.NewStructField("data", base_ptr))
	sptr_t.underlying.Finish()

	m.addValueType(&sptr_t)
	return &sptr_t
}

func (t *SPtr) Name() string         { return t.Base.Name() + ".$sptr" }
func (t *SPtr) Size() int            { return t.underlying.Size() }
func (t *SPtr) align() int           { return t.underlying.align() }
func (t *SPtr) Kind() TypeKind       { return kSPtr }
func (t *SPtr) onFree() int          { return t.underlying.onFree() }
func (t *SPtr) Raw() []wat.ValueType { return t.underlying.Raw() }
func (t *SPtr) Equal(u ValueType) bool {
	if ut, ok := u.(*SPtr); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}

func (t *SPtr) emitHeapAlloc() (insts []wat.Inst) {
	//insts = append(insts, wat.NewBlank())
	//insts = append(insts, wat.NewComment("Ref.emitHeapAlloc start"))

	insts = append(insts, t._base_block.emitHeapAlloc(NewConst("1", t.underlying._u32))...)
	insts = append(insts, wat.NewInstCall("$wa.runtime.DupI32"))
	insts = append(insts, NewConst("16", t.underlying._u32).EmitPush()...)
	insts = append(insts, wat.NewInstAdd(wat.U32{}))

	//insts = append(insts, wat.NewComment("Ref.emitHeapAlloc end"))
	//insts = append(insts, wat.NewBlank())

	return
}

func (t *SPtr) emitStackAlloc() (insts []wat.Inst) {
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

func (t *SPtr) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.underlying.EmitLoadFromAddr(addr, offset)
}

/**************************************
aSPtr:
**************************************/
type aSPtr struct {
	aStruct
	typ *SPtr
}

func newValue_SPtr(name string, kind ValueKind, typ *SPtr) *aSPtr {
	var v aSPtr
	v.typ = typ
	v.aStruct = *newValue_Struct(name, kind, typ.underlying)
	return &v
}

func (v *aSPtr) Type() ValueType { return v.typ }

func (v *aSPtr) raw() []wat.Value        { return v.aStruct.raw() }
func (v *aSPtr) EmitInit() []wat.Inst    { return v.aStruct.EmitInit() }
func (v *aSPtr) EmitPush() []wat.Inst    { return v.aStruct.EmitPush() }
func (v *aSPtr) EmitPop() []wat.Inst     { return v.aStruct.EmitPop() }
func (v *aSPtr) EmitRelease() []wat.Inst { return v.aStruct.EmitRelease() }

func (v *aSPtr) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	return v.aStruct.emitStoreToAddr(addr, offset)
}

func (v *aSPtr) emitGetValue() []wat.Inst {
	return v.typ.Base.EmitLoadFromAddr(v.aStruct.Extract("data"), 0)
}

func (v *aSPtr) emitSetValue(d Value) []wat.Inst {
	if !d.Type().Equal(v.typ.Base) && !v.typ.Base.Equal(v.typ._void) {
		logger.Fatal("Type not match")
		return nil
	}
	return d.emitStoreToAddr(v.aStruct.Extract("data"), 0)
}
