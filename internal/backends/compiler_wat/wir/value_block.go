// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

/**************************************
Block:
**************************************/
type Block struct {
	Base ValueType
}

func NewBlock(base ValueType) Block  { return Block{Base: base} }
func (t Block) Name() string         { return t.Base.Name() + ".$$block" }
func (t Block) size() int            { return 4 }
func (t Block) align() int           { return 4 }
func (t Block) Raw() []wat.ValueType { return []wat.ValueType{wat.I32{}} }
func (t Block) Equal(u ValueType) bool {
	if ut, ok := u.(Block); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}
func (t Block) onFree(m *Module) int {
	var f Function
	f.Name = "$" + t.Name() + ".$$onFree"
	if i := m.findTableElem(f.Name); i != 0 {
		return i
	}

	f.Result = VOID{}
	ptr := NewLocal("$ptr", I32{})
	f.Params = append(f.Params, ptr)

	f.Insts = append(f.Insts, ptr.EmitPush()...)
	f.Insts = append(f.Insts, wat.NewInstLoad(wat.I32{}, 0, 1))
	f.Insts = append(f.Insts, wat.NewInstCall("$wa.RT.Block.Release"))
	f.Insts = append(f.Insts, ptr.EmitPush()...)
	f.Insts = append(f.Insts, wat.NewInstConst(wat.I32{}, "0"))
	f.Insts = append(f.Insts, wat.NewInstStore(wat.I32{}, 0, 1))

	return m.addTableFunc(&f)
}

func (t Block) emitLoadFromAddr(addr Value, offset int) (insts []wat.Inst) {
	insts = append(insts, addr.EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.I32{}, offset, 1))
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Retain"))
	return
}

func (t Block) emitHeapAlloc(item_count Value, module *Module) (insts []wat.Inst) {
	switch item_count.Kind() {
	case ValueKindConst:
		c, err := strconv.Atoi(item_count.Name())
		if err != nil {
			logger.Fatalf("%v\n", err)
			return nil
		}
		insts = append(insts, NewConst(strconv.Itoa(t.Base.size()*c+16), I32{}).EmitPush()...)
		insts = append(insts, wat.NewInstCall("$waHeapAlloc"))

	default:
		if !item_count.Type().Equal(I32{}) {
			logger.Fatal("item_count should be i32")
			return nil
		}

		insts = append(insts, item_count.EmitPush()...)
		insts = append(insts, NewConst(strconv.Itoa(t.Base.size()), I32{}).EmitPush()...)
		insts = append(insts, wat.NewInstMul(wat.I32{}))
		insts = append(insts, NewConst("16", I32{}).EmitPush()...)
		insts = append(insts, wat.NewInstAdd(wat.I32{}))
		insts = append(insts, wat.NewInstCall("$waHeapAlloc"))

	}

	insts = append(insts, item_count.EmitPush()...)                                           //item_count
	insts = append(insts, NewConst(strconv.Itoa(t.Base.onFree(module)), I32{}).EmitPush()...) //free_method
	insts = append(insts, NewConst(strconv.Itoa(t.Base.size()), I32{}).EmitPush()...)         //item_size
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Init"))

	return
}

/**************************************
aBlock:
**************************************/
type aBlock struct {
	aValue
}

func newValueBlock(name string, kind ValueKind, base_type ValueType) *aBlock {
	return &aBlock{aValue: aValue{name: name, kind: kind, typ: NewBlock(base_type)}}
}

func (v *aBlock) raw() []wat.Value {
	return []wat.Value{wat.NewVarI32(v.name)}
}

func (v *aBlock) EmitInit() (insts []wat.Inst) {
	insts = append(insts, wat.NewInstConst(wat.I32{}, "0"))
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
	insts = append(insts, wat.NewInstLoad(wat.I32{}, offset, 1))
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Release"))

	insts = append(insts, wat.NewInstStore(toWatType(v.Type()), offset, 1))
	return
}
