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
	Base ValueType
	_i32 ValueType
	_u32 ValueType
}

func (m *Module) GenValueType_Block(base ValueType) *Block {
	block_t := Block{Base: base}
	t, ok := m.findValueType(block_t.Name())
	if ok {
		return t.(*Block)
	}

	block_t._i32 = m.I32
	block_t._u32 = m.U32
	m.regValueType(&block_t)
	return &block_t
}

func (t *Block) Name() string         { return t.Base.Name() + ".$$block" }
func (t *Block) size() int            { return 4 }
func (t *Block) align() int           { return 4 }
func (t *Block) Raw() []wat.ValueType { return []wat.ValueType{wat.U32{}} }
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
	f.Insts = append(f.Insts, wat.NewInstCall("$wa.RT.Block.Release"))
	f.Insts = append(f.Insts, ptr.EmitPush()...)
	f.Insts = append(f.Insts, wat.NewInstConst(wat.U32{}, "0"))
	f.Insts = append(f.Insts, wat.NewInstStore(wat.U32{}, 0, 1))

	currentModule.AddFunc(&f)
	return currentModule.AddTableElem(f.InternalName)
}

func (t *Block) EmitLoadFromAddr(addr Value, offset int) (insts []wat.Inst) {
	insts = append(insts, addr.EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.U32{}, offset, 1))
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Retain"))
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
		insts = append(insts, NewConst(strconv.Itoa(t.Base.size()*c+16), t._u32).EmitPush()...)
		insts = append(insts, wat.NewInstCall("$waHeapAlloc"))

	default:
		if !item_count.Type().Equal(t._u32) && !item_count.Type().Equal(t._i32) {
			logger.Fatal("item_count should be u32|i32")
			return nil
		}

		insts = append(insts, item_count.EmitPush()...)
		insts = append(insts, NewConst(strconv.Itoa(t.Base.size()), t._u32).EmitPush()...)
		insts = append(insts, wat.NewInstMul(wat.U32{}))
		insts = append(insts, NewConst("16", t._u32).EmitPush()...)
		insts = append(insts, wat.NewInstAdd(wat.U32{}))
		insts = append(insts, wat.NewInstCall("$waHeapAlloc"))

	}

	insts = append(insts, item_count.EmitPush()...)                                      //item_count
	insts = append(insts, NewConst(strconv.Itoa(t.Base.onFree()), t._u32).EmitPush()...) //free_method
	insts = append(insts, NewConst(strconv.Itoa(t.Base.size()), t._u32).EmitPush()...)   //item_size
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Init"))

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
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Retain"))
	return
}

func (v *aBlock) EmitPop() (insts []wat.Inst) {
	insts = append(insts, v.EmitRelease()...)
	insts = append(insts, v.pop(v.name))
	return
}

func (v *aBlock) EmitRelease() (insts []wat.Inst) {
	insts = append(insts, v.push(v.name))
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Release"))
	return
}

func (v *aBlock) emitStoreToAddr(addr Value, offset int) (insts []wat.Inst) {
	insts = append(insts, addr.EmitPush()...)
	insts = append(insts, v.EmitPush()...)

	insts = append(insts, addr.EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.U32{}, offset, 1))
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Release"))

	insts = append(insts, wat.NewInstStore(toWatType(v.Type()), offset, 1))
	return
}
