// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
)

/**************************************
Slice:
**************************************/
type Slice struct {
	Base ValueType
	Struct
}

func NewSlice(base ValueType) Slice {
	var v Slice
	v.Base = base
	var m []Field
	m = append(m, NewField("block", NewBlock(base)))
	m = append(m, NewField("data", NewPointer(base)))
	m = append(m, NewField("len", U32{}))
	m = append(m, NewField("cap", U32{}))
	v.Struct = NewStruct(base.Name()+".$$slice", m)
	return v
}
func (t Slice) Name() string         { return t.Base.Name() + ".$$slice" }
func (t Slice) size() int            { return t.Struct.size() }
func (t Slice) align() int           { return t.Struct.align() }
func (t Slice) onFree() int          { return t.Struct.onFree() }
func (t Slice) Raw() []wat.ValueType { return t.Struct.Raw() }
func (t Slice) Equal(u ValueType) bool {
	if ut, ok := u.(Slice); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}

func (t Slice) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.Struct.EmitLoadFromAddr(addr, offset)
}

/*这个函数极其不优雅*/
func (t Slice) emitGenFromRefOfSlice(x *aRef, low, high Value) (insts []wat.Inst) {
	//block
	insts = append(insts, x.Extract("data").EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.U32{}, 0, 1))
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Retain"))

	//data
	if low == nil {
		low = NewConst("0", U32{})
	}
	insts = append(insts, x.Extract("data").EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.U32{}, 4, 1))
	insts = append(insts, NewConst(strconv.Itoa(x.Type().(Ref).Base.(Array).Base.size()), U32{}).EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstMul(wat.U32{}))
	insts = append(insts, wat.NewInstAdd(wat.U32{}))

	//len:
	if high == nil {
		insts = append(insts, x.Extract("data").EmitPush()...)
		insts = append(insts, wat.NewInstLoad(wat.U32{}, 12, 1))
	} else {
		insts = append(insts, high.EmitPush()...)
	}
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstSub(wat.U32{}))

	//cap:
	insts = append(insts, x.Extract("data").EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.U32{}, 12, 1))
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstSub(wat.U32{}))

	return
}

func (t Slice) emitGenFromRefOfArray(x *aRef, low, high Value) (insts []wat.Inst) {
	//block
	insts = append(insts, x.Extract("block").EmitPush()...)

	//data
	if low == nil {
		low = NewConst("0", U32{})
	}
	insts = append(insts, x.Extract("data").EmitPush()...)
	insts = append(insts, NewConst(strconv.Itoa(x.Type().(Ref).Base.(Array).Base.size()), U32{}).EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstMul(wat.U32{}))
	insts = append(insts, wat.NewInstAdd(wat.U32{}))

	array_len := NewConst(strconv.Itoa(x.Type().(Ref).Base.(Array).Capacity), U32{})

	//len:
	if high == nil {
		high = array_len
	}
	insts = append(insts, high.EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstSub(wat.U32{}))

	//cap:
	insts = append(insts, array_len.EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstSub(wat.U32{}))

	return
}

func (t Slice) genAppendFunc() string {
	fn_name := "$" + t.Name() + ".append"
	if currentModule.findFunc(fn_name) != nil {
		return fn_name
	}

	var f Function
	f.Name = fn_name
	x := newValueSlice("x", ValueKindLocal, t.Base)
	y := newValueSlice("y", ValueKindLocal, t.Base)
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

	//if_new_len_le_cap
	f.Insts = append(f.Insts, new_len.EmitPush()...)
	f.Insts = append(f.Insts, x.Extract("cap").EmitPush()...)
	f.Insts = append(f.Insts, wat.NewInstLe(wat.U32{}))

	item := NewLocal("item", t.Base)
	f.Locals = append(f.Locals, item)
	src := NewLocal("src", NewPointer(t.Base))
	f.Locals = append(f.Locals, src)
	dest := NewLocal("dest", NewPointer(t.Base))
	f.Locals = append(f.Locals, dest)
	item_size := NewConst(strconv.Itoa(t.Base.size()), U32{})

	inst_if := wat.NewInstIf(nil, nil, t.Raw())
	{ //if_true
		var if_true []wat.Inst

		if_true = append(if_true, x.Extract("block").EmitPush()...)
		if_true = append(if_true, x.Extract("data").EmitPush()...)
		if_true = append(if_true, new_len.EmitPush()...)
		if_true = append(if_true, x.Extract("cap").EmitPush()...)

		//get src
		if_true = append(if_true, y.Extract("data").EmitPush()...)
		if_true = append(if_true, src.EmitPop()...)

		//get dest
		if_true = append(if_true, x.Extract("data").EmitPush()...)
		if_true = append(if_true, item_size.EmitPush()...)
		if_true = append(if_true, x_len.EmitPush()...)
		if_true = append(if_true, wat.NewInstMul(wat.U32{}))
		if_true = append(if_true, wat.NewInstAdd(wat.U32{}))
		if_true = append(if_true, dest.EmitPop()...)

		block := wat.NewInstBlock("block1")
		loop := wat.NewInstLoop("loop1")
		{
			loop.Insts = append(loop.Insts, y_len.EmitPush()...)
			loop.Insts = append(loop.Insts, wat.NewInstEqz(wat.U32{}))
			loop.Insts = append(loop.Insts, wat.NewInstIf([]wat.Inst{wat.NewInstBr("block1")}, nil, nil))

			//*dest = *src
			loop.Insts = append(loop.Insts, t.Base.EmitLoadFromAddr(src, 0)...)
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
			loop.Insts = append(loop.Insts, NewConst("1", U32{}).EmitPush()...)
			loop.Insts = append(loop.Insts, wat.NewInstSub(wat.U32{}))
			loop.Insts = append(loop.Insts, y_len.EmitPop()...)

			loop.Insts = append(loop.Insts, wat.NewInstBr("loop1"))
		}
		block.Insts = append(block.Insts, loop)
		if_true = append(if_true, block)

		inst_if.True = if_true
	}

	{ //if_false
		var if_false []wat.Inst

		new_cap := NewLocal("new_cap", U32{})
		f.Locals = append(f.Locals, new_cap)
		//gen new slice
		if_false = append(if_false, new_len.EmitPush()...)
		if_false = append(if_false, NewConst("2", U32{}).EmitPush()...)
		if_false = append(if_false, wat.NewInstMul(wat.U32{}))
		if_false = append(if_false, new_cap.EmitPop()...)
		if_false = append(if_false, NewBlock(t.Base).emitHeapAlloc(new_cap)...) //block

		if_false = append(if_false, wat.NewInstCall("$wa.RT.DupWatStack"))
		if_false = append(if_false, NewConst("16", U32{}).EmitPush()...)
		if_false = append(if_false, wat.NewInstAdd(wat.U32{})) //data
		if_false = append(if_false, wat.NewInstCall("$wa.RT.DupWatStack"))
		if_false = append(if_false, dest.EmitPop()...)     //dest
		if_false = append(if_false, new_len.EmitPush()...) //len
		if_false = append(if_false, new_cap.EmitPush()...) //cap

		//x->new
		{
			if_false = append(if_false, x.Extract("data").EmitPush()...)
			if_false = append(if_false, src.EmitPop()...)

			block := wat.NewInstBlock("block2")
			loop := wat.NewInstLoop("loop2")
			{
				loop.Insts = append(loop.Insts, x_len.EmitPush()...)
				loop.Insts = append(loop.Insts, wat.NewInstEqz(wat.U32{}))
				loop.Insts = append(loop.Insts, wat.NewInstIf([]wat.Inst{wat.NewInstBr("block2")}, nil, nil))

				//*dest = *src
				loop.Insts = append(loop.Insts, t.Base.EmitLoadFromAddr(src, 0)...)
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
				loop.Insts = append(loop.Insts, NewConst("1", U32{}).EmitPush()...)
				loop.Insts = append(loop.Insts, wat.NewInstSub(wat.U32{}))
				loop.Insts = append(loop.Insts, x_len.EmitPop()...)

				loop.Insts = append(loop.Insts, wat.NewInstBr("loop2"))
			}
			block.Insts = append(block.Insts, loop)
			if_false = append(if_false, block)
		}

		//y->new
		{
			if_false = append(if_false, y.Extract("data").EmitPush()...)
			if_false = append(if_false, src.EmitPop()...)

			block := wat.NewInstBlock("block3")
			loop := wat.NewInstLoop("loop3")
			{
				loop.Insts = append(loop.Insts, y_len.EmitPush()...)
				loop.Insts = append(loop.Insts, wat.NewInstEqz(wat.U32{}))
				loop.Insts = append(loop.Insts, wat.NewInstIf([]wat.Inst{wat.NewInstBr("block3")}, nil, nil))

				//*dest = *src
				loop.Insts = append(loop.Insts, t.Base.EmitLoadFromAddr(src, 0)...)
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
				loop.Insts = append(loop.Insts, NewConst("1", U32{}).EmitPush()...)
				loop.Insts = append(loop.Insts, wat.NewInstSub(wat.U32{}))
				loop.Insts = append(loop.Insts, y_len.EmitPop()...)

				loop.Insts = append(loop.Insts, wat.NewInstBr("loop3"))
			}
			block.Insts = append(block.Insts, loop)
			if_false = append(if_false, block)
		}

		inst_if.False = if_false
	}

	f.Insts = append(f.Insts, inst_if)
	f.Insts = append(f.Insts, x.EmitRelease()...)
	f.Insts = append(f.Insts, y.EmitRelease()...)

	currentModule.AddFunc(&f)
	return fn_name
}

/**************************************
aSlice:
**************************************/
type aSlice struct {
	aStruct
	typ Slice
}

func newValueSlice(name string, kind ValueKind, base_type ValueType) *aSlice {
	var v aSlice
	v.typ = NewSlice(base_type)
	v.aStruct = *newValueStruct(name, kind, v.typ.Struct)
	return &v
}

func (v *aSlice) raw() []wat.Value        { return v.aStruct.raw() }
func (v *aSlice) EmitInit() []wat.Inst    { return v.aStruct.EmitInit() }
func (v *aSlice) EmitPush() []wat.Inst    { return v.aStruct.EmitPush() }
func (v *aSlice) EmitPop() []wat.Inst     { return v.aStruct.EmitPop() }
func (v *aSlice) EmitRelease() []wat.Inst { return v.aStruct.EmitRelease() }

func (v *aSlice) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	return v.aStruct.emitStoreToAddr(addr, offset)
}

func (v *aSlice) emitSub(low, high Value) (insts []wat.Inst) {
	//block
	insts = append(insts, v.Extract("block").EmitPush()...)

	//data:
	if low == nil {
		low = NewConst("0", U32{})
	}
	insts = append(insts, v.Extract("data").EmitPush()...)
	insts = append(insts, NewConst(strconv.Itoa(v.Type().(Slice).Base.size()), U32{}).EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstMul(wat.U32{}))
	insts = append(insts, wat.NewInstAdd(wat.U32{}))

	//len:
	if high == nil {
		high = v.Extract("len")
	}
	insts = append(insts, high.EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstSub(wat.U32{}))

	//cap:
	insts = append(insts, v.Extract("cap").EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstSub(wat.U32{}))

	return
}
