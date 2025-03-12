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
	_base_ptr   *Ptr
	_void       ValueType
}

func (m *Module) GenValueType_Ref(base ValueType) *Ref {
	ref_t := Ref{Base: base}
	ref_t.name = base.Named() + ".$ref"
	t, ok := m.findValueType(ref_t.name)
	if ok {
		return t.(*Ref)
	}

	ref_t._base_block = m.GenValueType_Block(base)
	ref_t._void = m.VOID
	ref_t._base_ptr = m.GenValueType_Ptr(base)

	ref_t.underlying = m.genInternalStruct(ref_t.name + ".underlying")
	ref_t.underlying.AppendField(m.NewStructField("b", ref_t._base_block))
	ref_t.underlying.AppendField(m.NewStructField("d", ref_t._base_ptr))
	ref_t.underlying.Finish()

	m.addValueType(&ref_t)
	return &ref_t
}

func (t *Ref) Size() int            { return t.underlying.Size() }
func (t *Ref) align() int           { return t.underlying.align() }
func (t *Ref) Kind() TypeKind       { return kRef }
func (t *Ref) OnFree() int          { return t.underlying.OnFree() }
func (t *Ref) Raw() []wat.ValueType { return t.underlying.Raw() }
func (t *Ref) Equal(u ValueType) bool {
	if ut, ok := u.(*Ref); ok {
		// 在包含自身引用的结构体中，比较 Base 会引发无限嵌套
		//return t.Base.Equal(ut.Base)

		// 如果按照 wir.Module 的假设，任一种实类型均只存在一个实例，那么可用下列判断提前解除嵌套
		//if t == ut {
		//	return true
		//}

		// 存在风险，依赖于Name的唯一性
		return t.Base.Named() == ut.Base.Named()
	}
	return false
}

func (t *Ref) emitHeapAlloc() (insts []wat.Inst) {
	//insts = append(insts, wat.NewBlank())
	//insts = append(insts, wat.NewComment("Ref.emitHeapAlloc start"))

	insts = append(insts, t._base_block.emitHeapAlloc(NewConst("1", t.underlying._u32))...)

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
	insts = append(insts, wat.NewInstCall("runtime.stackAlloc"))

	//insts = append(insts, wat.NewComment("Ref.emitStackAlloc end"))
	//insts = append(insts, wat.NewBlank())
	return
}

func (t *Ref) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.underlying.EmitLoadFromAddr(addr, offset)
}

func (t *Ref) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	return t.underlying.EmitLoadFromAddrNoRetain(addr, offset)
}

func (t *Ref) newConstRef(ptr int) *aRef {
	v := newValue_Ref("", ValueKindConst, t)
	v.setFieldConstValue("b", NewConst("0", t._base_block))
	v.setFieldConstValue("d", NewConst(strconv.Itoa(ptr), t._base_ptr))
	return v
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

func (v *aRef) emitGetValue() []wat.Inst {
	return v.typ.Base.EmitLoadFromAddr(v.aStruct.ExtractByName("d"), 0)
}

func (v *aRef) emitSetValue(d Value) []wat.Inst {
	if !d.Type().Equal(v.typ.Base) && !v.typ.Base.Equal(v.typ._void) {
		logger.Fatal("Type not match")
		return nil
	}
	return d.emitStoreToAddr(v.aStruct.ExtractByName("d"), 0)
}

func (v *aRef) emitEq(r Value) (insts []wat.Inst, ok bool) {
	if !v.Type().Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}

	return v.ExtractByName("d").emitEq(r.(*aRef).ExtractByName("d"))
}

func (v *aRef) emitCompare(r Value) (insts []wat.Inst) {
	if !v.Type().Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}

	return v.ExtractByName("d").emitCompare(r.(*aRef).ExtractByName("d"))
}

func (v *aRef) getConstPtr() int {
	if v.kind != ValueKindConst {
		logger.Fatal("Must be a const")
	}

	i, _ := strconv.Atoi(v.ExtractByName("d").Name())
	return i
}

func (v *aRef) emitGenSetFinalizer(fn_id int) (insts []wat.Inst) {
	insts = append(insts, v.ExtractByName("b").EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(fn_id)))
	insts = append(insts, wat.NewInstCall("runtime.Block.SetFinalizer"))

	return
}
