// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
Block:
**************************************/
type Block struct {
	tCommon
	Base ValueType
	_i32 ValueType
	_u32 ValueType
}

func (m *Module) GenValueType_Block(base ValueType) *Block {
	block_t := Block{Base: base}
	//t, ok := m.findValueType(block_t.Name())
	//if ok {
	//	return t.(*Block)
	//}
	//m.addValueType(&block_t)

	block_t._i32 = m.I32
	block_t._u32 = m.U32
	return &block_t
}

func (t *Block) Name() string         { return t.Base.Name() + ".$$block" }
func (t *Block) Size() int            { return 4 }
func (t *Block) align() int           { return 4 }
func (t *Block) Kind() TypeKind       { return kBlock }
func (t *Block) Raw() []wat.ValueType { return []wat.ValueType{wat.U32{}} }

func (t *Block) typeInfoAddr() int {
	logger.Fatalf("Internal type: %s shouldn't have typeInfo.", t.Name())
	return 0
}

func (t *Block) Equal(u ValueType) bool {
	if ut, ok := u.(*Block); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}

func (t *Block) onFree() int {
	var f Function
	f.InternalName = "$" + GenSymbolName(t.Name()) + ".$$onFree"
	if i := currentModule.findTableElem(f.InternalName); i != 0 {
		return i
	}

	ptr := NewLocal("ptr", t._u32)
	f.Params = append(f.Params, ptr)

	f.Insts = append(f.Insts, ptr.EmitPush()...)
	f.Insts = append(f.Insts, wat.NewInstLoad(wat.U32{}, 0, 1))
	f.Insts = append(f.Insts, wat.NewInstCall("$Release"))
	f.Insts = append(f.Insts, ptr.EmitPush()...)
	f.Insts = append(f.Insts, wat.NewInstConst(wat.U32{}, "0"))
	f.Insts = append(f.Insts, wat.NewInstStore(wat.U32{}, 0, 1))

	currentModule.AddFunc(&f)
	return currentModule.AddTableElem(f.InternalName)
}

func (t *Block) EmitLoadFromAddr(addr Value, offset int) (insts []wat.Inst) {
	insts = append(insts, addr.EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.U32{}, offset, 1))
	insts = append(insts, wat.NewInstCall("$Retain"))
	return
}

func (t *Block) emitHeapAlloc(item_count Value) (insts []wat.Inst) {
	switch item_count.Kind() {
	case ValueKindConst:
		c, err := strconv.Atoi(item_count.Name())
		if err != nil {
			logger.Fatalf("%v\n", err)
			return nil
		}
		insts = append(insts, NewConst(strconv.Itoa(t.Base.Size()*c+16), t._u32).EmitPush()...)
		insts = append(insts, wat.NewInstCall("$waHeapAlloc"))

	default:
		if !item_count.Type().Equal(t._u32) && !item_count.Type().Equal(t._i32) {
			logger.Fatal("item_count should be u32|i32")
			return nil
		}

		insts = append(insts, item_count.EmitPush()...)
		insts = append(insts, NewConst(strconv.Itoa(t.Base.Size()), t._u32).EmitPush()...)
		insts = append(insts, wat.NewInstMul(wat.U32{}))
		insts = append(insts, NewConst("16", t._u32).EmitPush()...)
		insts = append(insts, wat.NewInstAdd(wat.U32{}))
		insts = append(insts, wat.NewInstCall("$waHeapAlloc"))

	}

	insts = append(insts, item_count.EmitPush()...)                                      //item_count
	insts = append(insts, NewConst(strconv.Itoa(t.Base.onFree()), t._u32).EmitPush()...) //free_method
	insts = append(insts, NewConst(strconv.Itoa(t.Base.Size()), t._u32).EmitPush()...)   //item_size
	insts = append(insts, wat.NewInstCall("$wa.runtime.Block.Init"))

	return
}

/**************************************
aBlock:
**************************************/
type aBlock struct {
	aValue
}

func newValue_Block(name string, kind ValueKind, typ *Block) *aBlock {
	return &aBlock{aValue: aValue{name: name, kind: kind, typ: typ}}
}

func (v *aBlock) raw() []wat.Value {
	return []wat.Value{wat.NewVarU32(v.name)}
}

func (v *aBlock) EmitInit() (insts []wat.Inst) {
	insts = append(insts, wat.NewInstConst(wat.U32{}, "0"))
	insts = append(insts, v.pop(v.name))
	return
}

func (v *aBlock) EmitPush() (insts []wat.Inst) {
	insts = append(insts, v.push(v.name))
	if v.Kind() == ValueKindConst && v.Name() == "0" {
		return
	}

	insts = append(insts, wat.NewInstCall("$Retain"))
	return
}

func (v *aBlock) EmitPop() (insts []wat.Inst) {
	insts = append(insts, v.EmitRelease()...)
	insts = append(insts, v.pop(v.name))
	return
}

func (v *aBlock) EmitRelease() (insts []wat.Inst) {
	if v.Kind() == ValueKindConst && v.Name() == "0" {
		return
	}
	insts = append(insts, v.push(v.name))
	insts = append(insts, wat.NewInstCall("$Release"))
	return
}

func (v *aBlock) emitStoreToAddr(addr Value, offset int) (insts []wat.Inst) {
	//if addr.Kind() == ValueKindGlobal_Pointer {
	//	if v.Kind() == ValueKindConst && v.Name() == "0" {
	//		return
	//	}
	//}
	insts = append(insts, addr.EmitPush()...)                    // a
	insts = append(insts, v.EmitPush()...)                       // a v
	insts = append(insts, addr.EmitPush()...)                    // a v a
	insts = append(insts, wat.NewInstLoad(wat.U32{}, offset, 1)) // a v o
	insts = append(insts, wat.NewInstCall("$Release"))           // a v

	insts = append(insts, wat.NewInstStore(toWatType(v.Type()), offset, 1))
	return
}

func (v *aBlock) emitStore(offset int) (insts []wat.Inst) {
	insts = append(insts, wat.NewInstCall("$wa.runtime.DupI32"))  // a
	insts = append(insts, wat.NewInstCall("$wa.runtime.DupI32"))  // a a
	insts = append(insts, v.EmitPush()...)                        // a a v
	insts = append(insts, wat.NewInstCall("$wa.runtime.SwapI32")) // a v a
	insts = append(insts, wat.NewInstLoad(wat.U32{}, offset, 1))  // a v o
	insts = append(insts, wat.NewInstCall("$Release"))            // a v

	insts = append(insts, wat.NewInstStore(toWatType(v.Type()), offset, 1))
	return
}

func (v *aBlock) Bin() (b []byte) {
	if v.Kind() != ValueKindConst {
		panic("Value.bin(): const only!")
	}

	b = make([]byte, 4)
	i, _ := strconv.Atoi(v.Name())
	b[0] = byte(i & 0xFF)
	b[1] = byte((i >> 8) & 0xFF)
	b[2] = byte((i >> 16) & 0xFF)
	b[3] = byte((i >> 24) & 0xFF)

	return
}

func (v *aBlock) emitEq(r Value) ([]wat.Inst, bool) {
	//logger.Fatal("aBlock shouldn't be compared.")
	return nil, false
}
