// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
Array:
**************************************/
type Array struct {
	tCommon
	Base       ValueType
	underlying *Struct
	Capacity   int
}

func (m *Module) GenValueType_Array(base ValueType, capacity int, name string) *Array {

	arr_t := Array{Base: base, Capacity: capacity}
	if len(name) > 0 {
		arr_t.name = name
	} else {
		arr_t.name = base.Named() + ".$array" + strconv.Itoa(capacity)
	}

	t, ok := m.findValueType(arr_t.name)
	if ok {
		return t.(*Array)
	}

	arr_t.underlying = m.genInternalStruct(arr_t.name + ".underlying")
	for i := 0; i < capacity; i++ {
		field := m.NewStructField("m"+strconv.Itoa(i), base)
		arr_t.underlying.AppendField(field)
	}
	arr_t.underlying.Finish()

	m.addValueType(&arr_t)
	return &arr_t
}

func (t *Array) Size() int            { return t.underlying.Size() }
func (t *Array) align() int           { return t.underlying.align() }
func (t *Array) Kind() TypeKind       { return kArray }
func (t *Array) Raw() []wat.ValueType { return t.underlying.Raw() }
func (t *Array) OnFree() int          { return t.underlying.OnFree() }

func (t *Array) Equal(u ValueType) bool {
	if ut, ok := u.(*Array); ok {
		return t.Base.Equal(ut.Base) && t.Capacity == ut.Capacity
	}
	return false
}

func (t *Array) EmitLoadFromAddr(addr Value, offset int) (insts []wat.Inst) {
	return t.underlying.EmitLoadFromAddr(addr, offset)
}

func (t *Array) EmitLoadFromAddrNoRetain(addr Value, offset int) (insts []wat.Inst) {
	return t.underlying.EmitLoadFromAddrNoRetain(addr, offset)
}

func (t *Array) genFunc_IndexOf(m *Module) string {
	if t.Capacity == 0 {
		return ""
	}

	fn_name := "$" + t.Named() + ".$IndexOf"
	if m.FindFunc(fn_name) != nil {
		return fn_name
	}

	var f Function
	f.InternalName = fn_name
	x := newValue_Array("x", ValueKindLocal, t)
	id := newValue_Basic("id", ValueKindLocal, t.underlying._u32)
	f.Params = append(f.Params, x)
	f.Params = append(f.Params, id)
	f.Results = append(f.Results, t.Base)

	ret := NewLocal("ret", t.Base)
	f.Locals = append(f.Locals, ret)

	var block_pre wat.Inst
	{
		table := make([]int, t.Capacity+1)
		for i := 0; i < t.Capacity; i++ {
			table[i] = i
		}
		table[t.Capacity] = t.Capacity - 1
		block_sel := wat.NewInstBlock("block_sel")
		block_sel.Insts = append(block_sel.Insts, id.EmitPush()...)
		block_sel.Insts = append(block_sel.Insts, wat.NewInstBrTable(table))
		block_pre = block_sel
	}
	for i := 0; i < t.Capacity; i++ {
		block := wat.NewInstBlock("block" + strconv.Itoa(i))
		block.Insts = append(block.Insts, block_pre)

		block.Insts = append(block.Insts, x.ExtractByName("m"+strconv.Itoa(i)).EmitPush()...)
		block.Insts = append(block.Insts, ret.EmitPop()...)
		block.Insts = append(block.Insts, wat.NewInstBr("block"+strconv.Itoa(t.Capacity-1)))

		block_pre = block
	}

	f.Insts = append(f.Insts, block_pre)
	f.Insts = append(f.Insts, ret.EmitPush()...)
	m.AddFunc(&f)
	return fn_name
}

/**************************************
aArray:
**************************************/
type aArray struct {
	aStruct
	typ *Array
}

func newValue_Array(name string, kind ValueKind, typ *Array) *aArray {
	var v aArray
	v.typ = typ
	v.aStruct = *newValue_Struct(name, kind, typ.underlying)
	return &v
}

func (v *aArray) Type() ValueType { return v.typ }

func (v *aArray) emitStoreToAddr(addr Value, offset int) (insts []wat.Inst) {
	if !addr.Type().(*Ptr).Base.Equal(v.Type()) {
		logger.Fatal("Type not match")
		return nil
	}

	return v.aStruct.emitStoreToAddr(addr, offset)
}

func (v *aArray) emitIndexOf(m *Module, id Value) (insts []wat.Inst) {
	fn_name := v.typ.genFunc_IndexOf(m)
	if len(fn_name) == 0 {
		zero_value := NewConst("0", v.typ.Base)
		insts = append(insts, zero_value.EmitPush()...)
		return
	}

	insts = append(insts, v.EmitPushNoRetain()...)
	insts = append(insts, id.EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstCall(fn_name))

	return
}
