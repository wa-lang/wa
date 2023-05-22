// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
Interface:
**************************************/
type Interface struct {
	tCommon
	name       string
	underlying *Struct
}

func (m *Module) GenValueType_Interface(name string) *Interface {
	t, ok := m.findValueType(name)
	if ok {
		return t.(*Interface)
	}

	var interface_t Interface
	interface_t.name = name

	interface_t.underlying = m.genInternalStruct(interface_t.Name() + ".underlying")
	interface_t.underlying.AppendField(m.NewStructField("data", m.GenValueType_SPtr(m.VOID)))
	interface_t.underlying.AppendField(m.NewStructField("itab", m.UPTR))
	interface_t.underlying.Finish()

	m.addValueType(&interface_t)
	return &interface_t
}

func (t *Interface) Name() string         { return t.name }
func (t *Interface) Size() int            { return t.underlying.Size() }
func (t *Interface) align() int           { return t.underlying.align() }
func (t *Interface) Kind() TypeKind       { return kInterface }
func (t *Interface) onFree() int          { return t.underlying.onFree() }
func (t *Interface) Raw() []wat.ValueType { return t.underlying.Raw() }
func (t *Interface) Equal(u ValueType) bool {
	if ut, ok := u.(*Interface); ok {
		return t.name == ut.name
	}
	return false
}

func (t *Interface) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.underlying.EmitLoadFromAddr(addr, offset)
}

func (t *Interface) emitGenFromSPtr(x *aSPtr) (insts []wat.Inst) {
	insts = append(insts, x.EmitPush()...) //data

	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(x.Type().Hash())))
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(t.Hash())))
	insts = append(insts, wat.NewInstConst(wat.I32{}, "0"))
	insts = append(insts, wat.NewInstCall("$wa.RT.getItab")) //itab

	return
}

func (t *Interface) emitGenFromInterface(x *aInterface) (insts []wat.Inst) {
	insts = append(insts, x.Extract("data").EmitPush()...) //data

	insts = append(insts, x.Extract("itab").EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.I32{}, 0, 4))
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(t.Hash())))
	insts = append(insts, wat.NewInstConst(wat.I32{}, "0"))
	insts = append(insts, wat.NewInstCall("$wa.RT.getItab")) //itab

	return
}

/**************************************
aInterface:
**************************************/
type aInterface struct {
	aStruct
	typ *Interface
}

func newValue_Interface(name string, kind ValueKind, typ *Interface) *aInterface {
	var v aInterface
	v.typ = typ
	v.aStruct = *newValue_Struct(name, kind, typ.underlying)
	return &v
}

func (v *aInterface) Type() ValueType { return v.typ }

func (v *aInterface) raw() []wat.Value        { return v.aStruct.raw() }
func (v *aInterface) EmitInit() []wat.Inst    { return v.aStruct.EmitInit() }
func (v *aInterface) EmitPush() []wat.Inst    { return v.aStruct.EmitPush() }
func (v *aInterface) EmitPop() []wat.Inst     { return v.aStruct.EmitPop() }
func (v *aInterface) EmitRelease() []wat.Inst { return v.aStruct.EmitRelease() }

func (v *aInterface) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	return v.aStruct.emitStoreToAddr(addr, offset)
}

func (v *aInterface) emitGetData(destType ValueType, commaOk bool) (insts []wat.Inst) {
	insts = append(insts, v.Extract("itab").EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.I32{}, 0, 4))

	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(destType.Hash())))
	insts = append(insts, wat.NewInstEq(wat.I32{}))

	ifBlock := wat.NewInstIf(nil, nil, nil)
	ifBlock.Ret = destType.Raw()

	// true:
	if _, ok := destType.(*SPtr); ok {
		ifBlock.True = v.Extract("data").EmitPush()
	} else {
		ifBlock.True = destType.EmitLoadFromAddr(v.Extract("data"), 0)
	}

	// false:
	ifBlock.False = NewConst("0", destType).EmitPush()

	if commaOk {
		ifBlock.Ret = append(ifBlock.Ret, wat.I32{})
		ifBlock.True = append(ifBlock.True, wat.NewInstConst(wat.I32{}, "1"))
		ifBlock.False = append(ifBlock.False, wat.NewInstConst(wat.I32{}, "0"))
	} else {
		ifBlock.False = append(ifBlock.False, wat.NewInstinstUnreachable())
	}

	insts = append(insts, ifBlock)
	return
}

func (v *aInterface) emitQueryInterface(destType ValueType, commaOk bool) (insts []wat.Inst) {
	logger.Fatal("Todo")
	return
}
