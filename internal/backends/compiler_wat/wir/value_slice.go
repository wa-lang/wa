// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
Slice:
**************************************/
type Slice struct {
	tCommon
	Base        ValueType
	underlying  *Struct
	_i32        ValueType
	_u32        ValueType
	_base_block *Block
	_base_ptr   *Ptr
}

func (m *Module) GenValueType_Slice(base ValueType, name string) *Slice {
	slice_t := Slice{Base: base}
	if len(name) > 0 {
		slice_t.name = name
	} else {
		slice_t.name = base.Named() + ".$slice"
	}
	t, ok := m.findValueType(slice_t.name)
	if ok {
		return t.(*Slice)
	}

	slice_t._i32 = m.I32
	slice_t._u32 = m.U32
	slice_t._base_block = m.GenValueType_Block(base)
	slice_t._base_ptr = m.GenValueType_Ptr(base)

	slice_t.underlying = m.genInternalStruct(slice_t.name + ".underlying")
	slice_t.underlying.AppendField(m.NewStructField("b", slice_t._base_block))
	slice_t.underlying.AppendField(m.NewStructField("d", slice_t._base_ptr))
	slice_t.underlying.AppendField(m.NewStructField("l", slice_t._u32))
	slice_t.underlying.AppendField(m.NewStructField("c", slice_t._u32))
	slice_t.underlying.Finish()

	m.addValueType(&slice_t)
	return &slice_t
}

func (t *Slice) Size() int            { return t.underlying.Size() }
func (t *Slice) align() int           { return t.underlying.align() }
func (t *Slice) Kind() TypeKind       { return kSlice }
func (t *Slice) OnFree() int          { return t.underlying.OnFree() }
func (t *Slice) Raw() []wat.ValueType { return t.underlying.Raw() }
func (t *Slice) Equal(u ValueType) bool {
	if ut, ok := u.(*Slice); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}

func (t *Slice) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.underlying.EmitLoadFromAddr(addr, offset)
}

func (t *Slice) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	return t.underlying.EmitLoadFromAddrNoRetain(addr, offset)
}

/*这个函数极其不优雅*/
func (t *Slice) emitGenFromRefOfSlice(x *aRef, low, high, max Value) (insts []wat.Inst) {
	//block
	insts = append(insts, x.ExtractByName("d").EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.U32{}, 0, 1))
	insts = append(insts, wat.NewInstCall("runtime.Block.Retain"))

	//data
	if low == nil {
		low = NewConst("0", t._u32)
	}
	insts = append(insts, x.ExtractByName("d").EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.U32{}, 4, 1))

	//insts = append(insts, NewConst(strconv.Itoa(x.Type().(*Ref).Base.(*Slice).Base.Size()), t._u32).EmitPush()...)
	insts = append(insts, NewConst(strconv.Itoa(t.Base.Size()), t._u32).EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstMul(wat.U32{}))
	insts = append(insts, wat.NewInstAdd(wat.U32{}))

	//len:
	if high == nil {
		insts = append(insts, x.ExtractByName("d").EmitPush()...)
		insts = append(insts, wat.NewInstLoad(wat.U32{}, 12, 1))
	} else {
		insts = append(insts, high.EmitPush()...)
	}
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstSub(wat.U32{}))

	//cap:
	if max == nil {
		insts = append(insts, x.ExtractByName("d").EmitPush()...)
		insts = append(insts, wat.NewInstLoad(wat.U32{}, 12, 1))
	} else {
		insts = append(insts, max.EmitPush()...)
	}
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstSub(wat.U32{}))

	return
}

func (t *Slice) emitGenFromRefOfArray(x *aRef, low, high, max Value) (insts []wat.Inst) {
	//block
	insts = append(insts, x.ExtractByName("b").EmitPush()...)

	//data
	if low == nil {
		low = NewConst("0", t._u32)
	}
	insts = append(insts, x.ExtractByName("d").EmitPush()...)
	//insts = append(insts, NewConst(strconv.Itoa(x.Type().(*Ref).Base.(*Array).Base.Size()), t._u32).EmitPush()...)
	insts = append(insts, NewConst(strconv.Itoa(t.Base.Size()), t._u32).EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstMul(wat.U32{}))
	insts = append(insts, wat.NewInstAdd(wat.U32{}))

	array_len := NewConst(strconv.Itoa(x.Type().(*Ref).Base.(*Array).Capacity), t._u32)

	//len:
	if high == nil {
		high = array_len
	}
	insts = append(insts, high.EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstSub(wat.U32{}))

	//cap:
	if max == nil {
		insts = append(insts, array_len.EmitPush()...)
	} else {
		insts = append(insts, max.EmitPush()...)
	}
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstSub(wat.U32{}))

	return
}

func (t *Slice) emitGenMake(Len, Cap Value) (insts []wat.Inst) {
	//block, data:
	insts = append(insts, t._base_block.emitHeapAlloc(Cap)...)

	//len:
	if !Len.Type().Equal(t._u32) && !Len.Type().Equal(t._i32) {
		logger.Fatal("Len should be u32")
		return nil
	}
	insts = append(insts, Len.EmitPush()...)

	//cap:
	insts = append(insts, Cap.EmitPush()...)

	return
}

func (t *Slice) genAppendFunc(m *Module) string {
	fn_name := "$" + GenSymbolName(t.Named()) + ".append"
	if m.FindFunc(fn_name) != nil {
		return fn_name
	}

	var f Function
	f.InternalName = fn_name
	x := newValue_Slice("x", ValueKindLocal, t)
	y := newValue_Slice("y", ValueKindLocal, t)
	f.Params = append(f.Params, x)
	f.Params = append(f.Params, y)
	f.Results = append(f.Results, t)

	item := NewLocal("item", t.Base)
	f.Locals = append(f.Locals, item)
	//f.Insts = append(f.Insts, item.EmitInit()...)

	x_len := NewLocal("x_len", x.ExtractByName("l").Type())
	f.Locals = append(f.Locals, x_len)
	f.Insts = append(f.Insts, x.ExtractByName("l").EmitPush()...)
	f.Insts = append(f.Insts, x_len.EmitPop()...)

	y_len := NewLocal("y_len", y.ExtractByName("l").Type())
	f.Locals = append(f.Locals, y_len)
	f.Insts = append(f.Insts, y.ExtractByName("l").EmitPush()...)
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
	f.Insts = append(f.Insts, x.ExtractByName("c").EmitPush()...)
	f.Insts = append(f.Insts, wat.NewInstLe(wat.U32{}))

	src := NewLocal("src", t._base_ptr)
	f.Locals = append(f.Locals, src)
	dest := NewLocal("dest", t._base_ptr)
	f.Locals = append(f.Locals, dest)
	item_size := NewConst(strconv.Itoa(t.Base.Size()), t._u32)

	inst_if := wat.NewInstIf(nil, nil, t.Raw())
	{ //if_true
		var if_true []wat.Inst

		if_true = append(if_true, x.ExtractByName("b").EmitPush()...)
		if_true = append(if_true, x.ExtractByName("d").EmitPush()...)
		if_true = append(if_true, new_len.EmitPush()...)
		if_true = append(if_true, x.ExtractByName("c").EmitPush()...)

		//get src
		if_true = append(if_true, y.ExtractByName("d").EmitPush()...)
		if_true = append(if_true, src.EmitPop()...)

		//get dest
		if_true = append(if_true, x.ExtractByName("d").EmitPush()...)
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
			loop.Insts = append(loop.Insts, NewConst("1", t._u32).EmitPush()...)
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

		new_cap := NewLocal("new_cap", t._u32)
		f.Locals = append(f.Locals, new_cap)
		//gen new slice
		if_false = append(if_false, new_len.EmitPush()...)
		if_false = append(if_false, NewConst("2", t._u32).EmitPush()...)
		if_false = append(if_false, wat.NewInstMul(wat.U32{}))
		if_false = append(if_false, new_cap.EmitPop()...)
		if_false = append(if_false, t._base_block.emitHeapAlloc(new_cap)...) //block, data

		if_false = append(if_false, wat.NewInstCall("runtime.DupI32"))
		if_false = append(if_false, dest.EmitPop()...)     //dest
		if_false = append(if_false, new_len.EmitPush()...) //len
		if_false = append(if_false, new_cap.EmitPush()...) //cap

		//x->new
		{
			if_false = append(if_false, x.ExtractByName("d").EmitPush()...)
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
				loop.Insts = append(loop.Insts, NewConst("1", t._u32).EmitPush()...)
				loop.Insts = append(loop.Insts, wat.NewInstSub(wat.U32{}))
				loop.Insts = append(loop.Insts, x_len.EmitPop()...)

				loop.Insts = append(loop.Insts, wat.NewInstBr("loop2"))
			}
			block.Insts = append(block.Insts, loop)
			if_false = append(if_false, block)
		}

		//y->new
		{
			if_false = append(if_false, y.ExtractByName("d").EmitPush()...)
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
				loop.Insts = append(loop.Insts, NewConst("1", t._u32).EmitPush()...)
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
	f.Insts = append(f.Insts, item.EmitRelease()...)

	m.AddFunc(&f)
	return fn_name
}

func (t *Slice) genCopyFunc(m *Module) string {
	var f Function
	f.InternalName = "$" + GenSymbolName(t.Base.Named()+".$slice") + ".copy"
	if m.FindFunc(f.InternalName) != nil {
		return f.InternalName
	}

	d := newValue_Slice("d", ValueKindLocal, t)
	s := newValue_Slice("s", ValueKindLocal, t)
	f.Params = append(f.Params, d)
	f.Params = append(f.Params, s)
	f.Results = append(f.Results, t._u32)

	item := NewLocal("item", t.Base)
	f.Locals = append(f.Locals, item)
	//f.Insts = append(f.Insts, item.EmitInit()...)

	count := NewLocal("count", d.ExtractByName("l").Type())
	f.Locals = append(f.Locals, count)
	{
		f.Insts = append(f.Insts, d.ExtractByName("l").EmitPush()...)
		f.Insts = append(f.Insts, s.ExtractByName("l").EmitPush()...)
		f.Insts = append(f.Insts, wat.NewInstGt(toWatType(count.Type())))

		ifs := wat.NewInstIf(nil, nil, nil)
		f.Insts = append(f.Insts, ifs)

		ifs.True = append(ifs.True, s.ExtractByName("l").EmitPush()...)
		ifs.True = append(ifs.True, count.EmitPop()...)
		ifs.False = append(ifs.False, d.ExtractByName("l").EmitPush()...)
		ifs.False = append(ifs.False, count.EmitPop()...)
	}
	f.Insts = append(f.Insts, count.EmitPush()...) //ret size

	dp := NewLocal("dp", d.ExtractByName("d").Type())
	f.Locals = append(f.Locals, dp)
	sp := NewLocal("sp", s.ExtractByName("d").Type())
	f.Locals = append(f.Locals, sp)
	item_size := NewLocal("item_size", d.ExtractByName("l").Type())
	f.Locals = append(f.Locals, item_size)
	{
		f.Insts = append(f.Insts, d.ExtractByName("d").EmitPush()...)
		f.Insts = append(f.Insts, s.ExtractByName("d").EmitPush()...)
		f.Insts = append(f.Insts, wat.NewInstLt(toWatType(d.ExtractByName("d").Type())))

		ifs := wat.NewInstIf(nil, nil, nil)
		f.Insts = append(f.Insts, ifs)
		// dp<sp
		ifs.True = append(ifs.True, d.ExtractByName("d").EmitPush()...)
		ifs.True = append(ifs.True, dp.EmitPop()...)
		ifs.True = append(ifs.True, s.ExtractByName("d").EmitPush()...)
		ifs.True = append(ifs.True, sp.EmitPop()...)
		ifs.True = append(ifs.True, wat.NewInstConst(wat.I32{}, strconv.Itoa(t.Base.Size())))
		ifs.True = append(ifs.True, item_size.EmitPop()...)

		// dp>sp
		ifs.False = append(ifs.False, count.EmitPush()...)
		ifs.False = append(ifs.False, wat.NewInstConst(wat.I32{}, "1"))
		ifs.False = append(ifs.False, wat.NewInstSub(wat.I32{}))
		ifs.False = append(ifs.False, wat.NewInstConst(wat.I32{}, strconv.Itoa(t.Base.Size())))
		ifs.False = append(ifs.False, wat.NewInstMul(wat.I32{}))
		ifs.False = append(ifs.False, item_size.EmitPop()...)

		ifs.False = append(ifs.False, d.ExtractByName("d").EmitPush()...)
		ifs.False = append(ifs.False, item_size.EmitPush()...)
		ifs.False = append(ifs.False, wat.NewInstAdd(wat.I32{}))
		ifs.False = append(ifs.False, dp.EmitPop()...)

		ifs.False = append(ifs.False, s.ExtractByName("d").EmitPush()...)
		ifs.False = append(ifs.False, item_size.EmitPush()...)
		ifs.False = append(ifs.False, wat.NewInstAdd(wat.I32{}))
		ifs.False = append(ifs.False, sp.EmitPop()...)

		ifs.False = append(ifs.False, wat.NewInstConst(wat.I32{}, "0"))
		ifs.False = append(ifs.False, wat.NewInstConst(wat.I32{}, strconv.Itoa(t.Base.Size())))
		ifs.False = append(ifs.False, wat.NewInstSub(wat.I32{}))
		ifs.False = append(ifs.False, item_size.EmitPop()...)
	}

	b0 := wat.NewInstBlock("b0")
	f.Insts = append(f.Insts, b0)

	l0 := wat.NewInstLoop("l0")
	b0.Insts = append(b0.Insts, l0)

	l0.Insts = append(l0.Insts, count.EmitPush()...)
	l0.Insts = append(l0.Insts, wat.NewInstEqz(wat.I32{}))
	{
		ifs := wat.NewInstIf(nil, nil, nil)
		l0.Insts = append(l0.Insts, ifs)

		ifs.True = append(ifs.True, wat.NewInstBr("b0"))

		ifs.False = append(ifs.False, t.Base.EmitLoadFromAddr(sp, 0)...)
		ifs.False = append(ifs.False, item.EmitPop()...)
		ifs.False = append(ifs.False, item.emitStoreToAddr(dp, 0)...)

		ifs.False = append(ifs.False, sp.EmitPush()...)
		ifs.False = append(ifs.False, item_size.EmitPush()...)
		ifs.False = append(ifs.False, wat.NewInstAdd(wat.I32{}))
		ifs.False = append(ifs.False, sp.EmitPop()...)

		ifs.False = append(ifs.False, dp.EmitPush()...)
		ifs.False = append(ifs.False, item_size.EmitPush()...)
		ifs.False = append(ifs.False, wat.NewInstAdd(wat.I32{}))
		ifs.False = append(ifs.False, dp.EmitPop()...)

		ifs.False = append(ifs.False, count.EmitPush()...)
		ifs.False = append(ifs.False, wat.NewInstConst(wat.I32{}, "1"))
		ifs.False = append(ifs.False, wat.NewInstSub(wat.I32{}))
		ifs.False = append(ifs.False, count.EmitPop()...)

		ifs.False = append(ifs.False, wat.NewInstBr("l0"))
	}

	f.Insts = append(f.Insts, item.EmitRelease()...)

	m.AddFunc(&f)
	return f.InternalName
}

/**************************************
aSlice:
**************************************/
type aSlice struct {
	aStruct
	typ *Slice
}

func newValue_Slice(name string, kind ValueKind, typ *Slice) *aSlice {
	var v aSlice
	v.typ = typ
	v.aStruct = *newValue_Struct(name, kind, typ.underlying)
	return &v
}

func (v *aSlice) Type() ValueType { return v.typ }

func (v *aSlice) emitSub(low, high, max Value) (insts []wat.Inst) {
	//block
	insts = append(insts, v.ExtractByName("b").EmitPush()...)

	//data:
	if low == nil {
		low = NewConst("0", v.typ._u32)
	}
	insts = append(insts, v.ExtractByName("d").EmitPush()...)
	insts = append(insts, NewConst(strconv.Itoa(v.typ.Base.Size()), v.typ._u32).EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstMul(wat.U32{}))
	insts = append(insts, wat.NewInstAdd(wat.U32{}))

	//len:
	if high == nil {
		high = v.ExtractByName("l")
	}
	insts = append(insts, high.EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstSub(wat.U32{}))

	//cap:
	if max == nil {
		insts = append(insts, v.ExtractByName("c").EmitPush()...)
	} else {
		insts = append(insts, max.EmitPush()...)
	}
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstSub(wat.U32{}))

	return
}

func (v *aSlice) emitEq(r Value) ([]wat.Inst, bool) {
	if r.Kind() != ValueKindConst {
		return nil, false
	}

	r_s, ok := r.(*aSlice)
	if !ok {
		logger.Fatal("r is not a Slice")
	}

	if r_s.Name() != "0" {
		logger.Fatal("r should be nil")
	}

	return v.ExtractByName("d").emitEq(r_s.ExtractByName("d"))
}

func (v *aSlice) emitCompare(r Value) (insts []wat.Inst) {
	if !r.Type().Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}

	t := r.(*aSlice)

	insts = append(insts, v.ExtractByName("d").EmitPushNoRetain()...)
	insts = append(insts, t.ExtractByName("d").EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstLt(toWatType(v.ExtractByName("d").Type())))

	instDLe := wat.NewInstIf(nil, nil, []wat.ValueType{wat.I32{}})

	instDLe.True = append(instDLe.True, wat.NewInstConst(wat.I32{}, "-1"))

	instDLe.False = append(instDLe.False, v.ExtractByName("d").EmitPushNoRetain()...)
	instDLe.False = append(instDLe.False, t.ExtractByName("d").EmitPushNoRetain()...)
	instDLe.False = append(instDLe.False, wat.NewInstGt(toWatType(v.ExtractByName("d").Type())))

	instLLe := wat.NewInstIf(nil, nil, []wat.ValueType{wat.I32{}})
	instLLe.True = append(instLLe.True, wat.NewInstConst(wat.I32{}, "1"))
	instLLe.False = v.ExtractByName("l").emitCompare(t.ExtractByName("l"))

	instDLe.False = append(instDLe.False, instLLe)

	insts = append(insts, instDLe)

	return
}

func (v *aSlice) emitConvertToBytes() (insts []wat.Inst) {
	// block:
	insts = append(insts, v.ExtractByName("b").EmitPush()...)

	// data:
	insts = append(insts, v.ExtractByName("d").EmitPush()...)

	// len:
	insts = append(insts, v.ExtractByName("l").EmitPush()...)
	insts = append(insts, wat.NewInstConst(wat.U32{}, strconv.Itoa(v.typ.Base.Size())))
	insts = append(insts, wat.NewInstMul(wat.U32{}))

	// cap:
	insts = append(insts, wat.NewInstCall("runtime.DupI32"))

	return
}
