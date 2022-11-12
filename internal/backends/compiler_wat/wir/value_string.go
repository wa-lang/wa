// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
)

/**************************************
String:
**************************************/
type String struct {
	Struct
}

func NewString() String {
	var v String
	var m []Field
	m = append(m, NewField("block", NewBlock(U8{})))
	m = append(m, NewField("data", NewPointer(U8{})))
	m = append(m, NewField("len", U32{}))
	v.Struct = NewStruct("string", m)
	return v
}
func (t String) Name() string           { return "string" }
func (t String) size() int              { return 12 }
func (t String) align() int             { return 4 }
func (t String) onFree() int            { return t.Struct.onFree() }
func (t String) Raw() []wat.ValueType   { return t.Struct.Raw() }
func (t String) Equal(u ValueType) bool { _, ok := u.(String); return ok }

func (t String) emitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.Struct.emitLoadFromAddr(addr, offset)
}

func (t String) genAppendStrFunc() string {
	fn_name := "$" + t.Name() + ".appendstr"
	if currentModule.findFunc(fn_name) != nil {
		return fn_name
	}

	var f Function
	f.Name = fn_name
	x := newValueString("x", ValueKindLocal)
	y := newValueString("y", ValueKindLocal)
	f.Params = append(f.Params, x)
	f.Params = append(f.Params, y)
	f.Result = t

	x_len := NewLocal("x_len", x.underlying.Extract("len").Type())
	f.Locals = append(f.Locals, x_len)
	f.Insts = append(f.Insts, x.underlying.Extract("len").EmitPush()...)
	f.Insts = append(f.Insts, x_len.EmitPop()...)

	y_len := NewLocal("y_len", y.underlying.Extract("len").Type())
	f.Locals = append(f.Locals, y_len)
	f.Insts = append(f.Insts, y.underlying.Extract("len").EmitPush()...)
	f.Insts = append(f.Insts, y_len.EmitPop()...)

	//gen new_len:
	new_len := NewLocal("new_len", x_len.Type())
	f.Locals = append(f.Locals, new_len)
	f.Insts = append(f.Insts, x_len.EmitPush()...)
	f.Insts = append(f.Insts, y_len.EmitPush()...)
	f.Insts = append(f.Insts, wat.NewInstAdd(wat.U32{}))
	f.Insts = append(f.Insts, new_len.EmitPop()...)

	item := NewLocal("item", U8{})
	f.Locals = append(f.Locals, item)
	src := NewLocal("src", NewPointer(U8{}))
	f.Locals = append(f.Locals, src)
	dest := NewLocal("dest", NewPointer(U8{}))
	f.Locals = append(f.Locals, dest)
	item_size := NewConst(strconv.Itoa(U8{}.size()), U32{})

	{ //if_false
		//gen new slice
		f.Insts = append(f.Insts, NewBlock(U8{}).emitHeapAlloc(new_len)...) //block

		f.Insts = append(f.Insts, wat.NewInstCall("$wa.RT.DupWatStack"))
		f.Insts = append(f.Insts, NewConst("16", U32{}).EmitPush()...)
		f.Insts = append(f.Insts, wat.NewInstAdd(wat.U32{})) //data
		f.Insts = append(f.Insts, wat.NewInstCall("$wa.RT.DupWatStack"))
		f.Insts = append(f.Insts, dest.EmitPop()...)     //dest
		f.Insts = append(f.Insts, new_len.EmitPush()...) //len

		//x->new
		{
			f.Insts = append(f.Insts, x.underlying.Extract("data").EmitPush()...)
			f.Insts = append(f.Insts, src.EmitPop()...)

			block := wat.NewInstBlock("block2")
			loop := wat.NewInstLoop("loop2")
			{
				loop.Insts = append(loop.Insts, x_len.EmitPush()...)
				loop.Insts = append(loop.Insts, wat.NewInstEqz(wat.U32{}))
				loop.Insts = append(loop.Insts, wat.NewInstIf([]wat.Inst{wat.NewInstBr("block2")}, nil, nil))

				//*dest = *src
				loop.Insts = append(loop.Insts, U8{}.emitLoadFromAddr(src, 0)...)
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
			f.Insts = append(f.Insts, block)
		}

		//y->new
		{
			f.Insts = append(f.Insts, y.underlying.Extract("data").EmitPush()...)
			f.Insts = append(f.Insts, src.EmitPop()...)

			block := wat.NewInstBlock("block3")
			loop := wat.NewInstLoop("loop3")
			{
				loop.Insts = append(loop.Insts, y_len.EmitPush()...)
				loop.Insts = append(loop.Insts, wat.NewInstEqz(wat.U32{}))
				loop.Insts = append(loop.Insts, wat.NewInstIf([]wat.Inst{wat.NewInstBr("block3")}, nil, nil))

				//*dest = *src
				loop.Insts = append(loop.Insts, U8{}.emitLoadFromAddr(src, 0)...)
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
	aValue
	underlying aStruct
}

func newValueString(name string, kind ValueKind) *aString {
	var v aString
	string_type := NewString()
	if kind == ValueKindConst {
		string_type.findFieldByName("block").const_val = NewConst("0", NewBlock(U8{}))
		ptr := currentModule.AddDataSeg([]byte(name))
		string_type.findFieldByName("data").const_val = NewConst(strconv.Itoa(ptr), NewPointer(U8{}))
		string_type.findFieldByName("len").const_val = NewConst(strconv.Itoa(len(name)), U32{})
	}
	v.aValue = aValue{name: name, kind: kind, typ: string_type}
	v.underlying = *newValueStruct(name, kind, string_type.Struct)
	return &v
}

func (v *aString) raw() []wat.Value        { return v.underlying.raw() }
func (v *aString) EmitInit() []wat.Inst    { return v.underlying.EmitInit() }
func (v *aString) EmitPush() []wat.Inst    { return v.underlying.EmitPush() }
func (v *aString) EmitPop() []wat.Inst     { return v.underlying.EmitPop() }
func (v *aString) EmitRelease() []wat.Inst { return v.underlying.EmitRelease() }

func (v *aString) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	return v.underlying.emitStoreToAddr(addr, offset)
}

func (v *aString) emitSub(low, high Value) (insts []wat.Inst) {
	//block
	insts = append(insts, v.underlying.Extract("block").EmitPush()...)

	//data:
	if low == nil {
		low = NewConst("0", U32{})
	}
	insts = append(insts, v.underlying.Extract("data").EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstAdd(wat.U32{}))

	//len:
	if high == nil {
		high = v.underlying.Extract("len")
	}
	insts = append(insts, high.EmitPush()...)
	insts = append(insts, low.EmitPush()...)
	insts = append(insts, wat.NewInstSub(wat.U32{}))

	return
}
