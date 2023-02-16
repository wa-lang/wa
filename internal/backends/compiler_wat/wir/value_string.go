// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
)

/**************************************
String:
**************************************/
type String struct {
	*Struct
	_u8       ValueType
	_u32      ValueType
	_u8_block *Block
	_u8_ptr   ValueType
}

func (m *Module) GenValueType_String() *String {
	var str_t String
	t, ok := m.findValueType(str_t.Name())
	if ok {
		return t.(*String)
	}

	str_t._u8 = m.U8
	str_t._u32 = m.U32
	str_t._u8_block = m.GenValueType_Block(m.U8)
	str_t._u8_ptr = m.GenValueType_Ptr(m.U8)

	var members []Field
	members = append(members, NewField("block", str_t._u8_block))
	members = append(members, NewField("data", str_t._u8_ptr))
	members = append(members, NewField("len", str_t._u32))
	str_t.Struct = m.GenValueType_Struct(str_t.Name()+".underlying", members)
	m.regValueType(&str_t)
	return &str_t

}

func (t *String) Name() string           { return "string" }
func (t *String) size() int              { return t.Struct.size() }
func (t *String) align() int             { return t.Struct.align() }
func (t *String) onFree() int            { return t.Struct.onFree() }
func (t *String) Raw() []wat.ValueType   { return t.Struct.Raw() }
func (t *String) Equal(u ValueType) bool { _, ok := u.(*String); return ok }

func (t *String) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.Struct.EmitLoadFromAddr(addr, offset)
}

func (t *String) genAppendStrFunc() string {
	fn_name := "$" + t.Name() + ".appendstr"
	if currentModule.FindFunc(fn_name) != nil {
		return fn_name
	}

	var f Function
	f.InternalName = fn_name
	x := newValue_String("x", ValueKindLocal, t)
	y := newValue_String("y", ValueKindLocal, t)
	f.Params = append(f.Params, x)
	f.Params = append(f.Params, y)
	f.Results = append(f.Results, t)

	x_len := NewLocal("x_len", x.Extract("len").Type())
	f.Locals = append(f.Locals, x_len)
	f.Insts = append(f.Insts, x.Extract("len").EmitPush()...)
	f.Insts = append(f.Insts, x_len.EmitPop()...)

	y_len := NewLocal("y_len", y.Extract("len").Type())
	f.Locals = append(f.Locals, y_len)
	f.Insts = append(f.Insts, y.Extract("len").EmitPush()...)
	f.Insts = append(f.Insts, y_len.EmitPop()...)

	//gen new_len:
	new_len := NewLocal("new_len", x_len.Type())
	f.Locals = append(f.Locals, new_len)
	f.Insts = append(f.Insts, x_len.EmitPush()...)
	f.Insts = append(f.Insts, y_len.EmitPush()...)
	f.Insts = append(f.Insts, wat.NewInstAdd(wat.U32{}))
	f.Insts = append(f.Insts, new_len.EmitPop()...)

	item := NewLocal("item", t._u8)
	f.Locals = append(f.Locals, item)
	src := NewLocal("src", t._u8_ptr)
	f.Locals = append(f.Locals, src)
	dest := NewLocal("dest", t._u8_ptr)
	f.Locals = append(f.Locals, dest)
	item_size := NewConst(strconv.Itoa(t._u8.size()), t._u32)

	{ //if_false
		//gen new slice
		f.Insts = append(f.Insts, t._u8_block.emitHeapAlloc(new_len)...) //block

		f.Insts = append(f.Insts, wat.NewInstCall("$wa.RT.DupWatStack"))
		f.Insts = append(f.Insts, NewConst("16", t._u32).EmitPush()...)
		f.Insts = append(f.Insts, wat.NewInstAdd(wat.U32{})) //data
		f.Insts = append(f.Insts, wat.NewInstCall("$wa.RT.DupWatStack"))
		f.Insts = append(f.Insts, dest.EmitPop()...)     //dest
		f.Insts = append(f.Insts, new_len.EmitPush()...) //len

		//x->new
		{
			f.Insts = append(f.Insts, x.Extract("data").EmitPush()...)
			f.Insts = append(f.Insts, src.EmitPop()...)

			block := wat.NewInstBlock("block2")
			loop := wat.NewInstLoop("loop2")
			{
				loop.Insts = append(loop.Insts, x_len.EmitPush()...)
				loop.Insts = append(loop.Insts, wat.NewInstEqz(wat.U32{}))
				loop.Insts = append(loop.Insts, wat.NewInstIf([]wat.Inst{wat.NewInstBr("block2")}, nil, nil))

				//*dest = *src
				loop.Insts = append(loop.Insts, t._u8.EmitLoadFromAddr(src, 0)...)
				loop.Insts = append(loop.Insts, item.EmitPop()...)
				loop.Insts = append(loop.Insts, item.emitStoreToAddr(dest, 0)...)

				loop.Insts = append(loop.Insts, src.EmitPush()...)
				loop.Insts = append(loop.Insts, item_size.EmitPush()...)
				loop.Insts = append(loop.Insts, wat.NewInstAdd(wat.U32{}))
				loop.Insts = append(loop.Insts, src.EmitPop()...)

				loop.Insts = append(loop.Insts, dest.EmitPush()...)
				loop.Insts = append(loop.Insts, item_size.EmitPush()...)
				loop.Insts = append(loop.Insts, wat.NewInstAdd(wat.U32{}))
				loop.Insts = append(loop.Insts, dest.EmitPop()...)

				loop.Insts = append(loop.Insts, x_len.EmitPush()...)
				loop.Insts = append(loop.Insts, NewConst("1", t._u32).EmitPush()...)
				loop.Insts = append(loop.Insts, wat.NewInstSub(wat.U32{}))
				loop.Insts = append(loop.Insts, x_len.EmitPop()...)

				loop.Insts = append(loop.Insts, wat.NewInstBr("loop2"))
			}
			block.Insts = append(block.Insts, loop)
			f.Insts = append(f.Insts, block)
		}

		//y->new
		{
			f.Insts = append(f.Insts, y.Extract("data").EmitPush()...)
			f.Insts = append(f.Insts, src.EmitPop()...)

			block := wat.NewInstBlock("block3")
			loop := wat.NewInstLoop("loop3")
			{
				loop.Insts = append(loop.Insts, y_len.EmitPush()...)
				loop.Insts = append(loop.Insts, wat.NewInstEqz(wat.U32{}))
				loop.Insts = append(loop.Insts, wat.NewInstIf([]wat.Inst{wat.NewInstBr("block3")}, nil, nil))

				//*dest = *src
				loop.Insts = append(loop.Insts, t._u8.EmitLoadFromAddr(src, 0)...)
				loop.Insts = append(loop.Insts, item.EmitPop()...)
				loop.Insts = append(loop.Insts, item.emitStoreToAddr(dest, 0)...)

				loop.Insts = append(loop.Insts, src.EmitPush()...)
				loop.Insts = append(loop.Insts, item_size.EmitPush()...)
				loop.Insts = append(loop.Insts, wat.NewInstAdd(wat.U32{}))
				loop.Insts = append(loop.Insts, src.EmitPop()...)

				loop.Insts = append(loop.Insts, dest.EmitPush()...)
				loop.Insts = append(loop.Insts, item_size.EmitPush()...)
				loop.Insts = append(loop.Insts, wat.NewInstAdd(wat.U32{}))
				loop.Insts = append(loop.Insts, dest.EmitPop()...)

				loop.Insts = append(loop.Insts, y_len.EmitPush()...)
				loop.Insts = append(loop.Insts, NewConst("1", t._u32).EmitPush()...)
				loop.Insts = append(loop.Insts, wat.NewInstSub(wat.U32{}))
				loop.Insts = append(loop.Insts, y_len.EmitPop()...)

				loop.Insts = append(loop.Insts, wat.NewInstBr("loop3"))
			}
			block.Insts = append(block.Insts, loop)
			f.Insts = append(f.Insts, block)
		}
	}

	f.Insts = append(f.Insts, x.EmitRelease()...)
	f.Insts = append(f.Insts, y.EmitRelease()...)

	currentModule.AddFunc(&f)
	return fn_name
}

/**************************************
aString:
**************************************/
type aString struct {
	aStruct
	typ *String
}

func newValue_String(name string, kind ValueKind, typ *String) *aString {
	var v aString
	v.typ = typ
	v.aStruct = *newValue_Struct(name, kind, typ.Struct)
	if kind == ValueKindConst {
		v.aStruct.setFieldConstValue("block", NewConst("0", typ._u8_block))
		ptr := currentModule.AddDataSeg([]byte(name))
		v.aStruct.setFieldConstValue("data", NewConst(strconv.Itoa(ptr), typ._u8_ptr))
		v.aStruct.setFieldConstValue("len", NewConst(strconv.Itoa(len(name)), typ._u32))
	}

	return &v
}

func (v *aString) Type() ValueType { return v.typ }

func (v *aString) raw() []wat.Value        { return v.aStruct.raw() }
func (v *aString) EmitInit() []wat.Inst    { return v.aStruct.EmitInit() }
func (v *aString) EmitPush() []wat.Inst    { return v.aStruct.EmitPush() }
func (v *aString) EmitPop() []wat.Inst     { return v.aStruct.EmitPop() }
func (v *aString) EmitRelease() []wat.Inst { return v.aStruct.EmitRelease() }

func (v *aString) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	return v.aStruct.emitStoreToAddr(addr, offset)
}

func (v *aString) emitSub(low, high Value) (insts []wat.Inst) {
	//block
	insts = append(insts, v.Extract("block").EmitPush()...)

	//data:
	if low == nil {
		low = NewConst("0", v.typ._u32)
	}
	insts = append(insts, v.Extract("data").EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstAdd(wat.U32{}))

	//len:
	if high == nil {
		high = v.Extract("len")
	}
	insts = append(insts, high.EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstSub(wat.U32{}))

	return
}

func (v *aString) emitAt(index Value) (insts []wat.Inst) {
	insts = append(insts, v.Extract("data").EmitPush()...)
	insts = append(insts, index.EmitPush()...)
	insts = append(insts, wat.NewInstAdd(wat.I32{}))
	insts = append(insts, wat.NewInstLoad8u(0, 1))
	return
}

func (v *aString) emitAt_CommaOk(index Value) (insts []wat.Inst) {
	insts = append(insts, v.Extract("len").EmitPush()...)
	insts = append(insts, index.EmitPush()...)
	insts = append(insts, wat.NewInstGt(wat.U32{}))

	{
		var instsTrue, instsFalse []wat.Inst
		instsTrue = append(instsTrue, v.Extract("data").EmitPush()...)
		instsTrue = append(instsTrue, index.EmitPush()...)
		instsTrue = append(instsTrue, wat.NewInstAdd(wat.I32{}))
		instsTrue = append(instsTrue, wat.NewInstLoad8u(0, 1))
		instsTrue = append(instsTrue, wat.NewInstConst(wat.I32{}, "1"))

		instsFalse = append(instsFalse, wat.NewInstConst(wat.I32{}, "0"))
		instsFalse = append(instsFalse, wat.NewInstConst(wat.I32{}, "0"))
		inst_if := wat.NewInstIf(instsTrue, instsFalse, []wat.ValueType{wat.I32{}, wat.I32{}})
		insts = append(insts, inst_if)
	}

	return
}
