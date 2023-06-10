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
	interface_t.underlying.AppendField(m.NewStructField("data", m.GenValueType_Ref(m.VOID)))
	interface_t.underlying.AppendField(m.NewStructField("itab", m.UPTR))
	interface_t.underlying.AppendField(m.NewStructField("eq", m.I32))
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

func (t *Interface) emitGenFromRef(x *aRef) (insts []wat.Inst) {
	insts = append(insts, x.EmitPush()...) //data

	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(x.Type().Hash())))
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(t.Hash())))
	insts = append(insts, wat.NewInstConst(wat.I32{}, "0"))
	insts = append(insts, wat.NewInstCall("$wa.runtime.getItab")) //itab

	insts = append(insts, wat.NewInstConst(wat.I32{}, "0")) //eq

	return
}

func (t *Interface) emitGenFromValue(x Value, xRefType *Ref, compID int) (insts []wat.Inst) {
	insts = append(insts, xRefType.emitHeapAlloc()...)
	insts = append(insts, x.emitStore(0)...) //data

	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(x.Type().Hash())))
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(t.Hash())))
	insts = append(insts, wat.NewInstConst(wat.I32{}, "0"))
	insts = append(insts, wat.NewInstCall("$wa.runtime.getItab")) //itab

	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(compID))) //eq

	return
}

func (t *Interface) emitGenFromInterface(x *aInterface) (insts []wat.Inst) {
	insts = append(insts, x.Extract("data").EmitPush()...) //data

	insts = append(insts, x.Extract("itab").EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.I32{}, 0, 4))
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(t.Hash())))
	insts = append(insts, wat.NewInstConst(wat.I32{}, "0"))
	insts = append(insts, wat.NewInstCall("$wa.runtime.getItab")) //itab

	insts = append(insts, x.Extract("eq").EmitPush()...) //eq

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
	if _, ok := destType.(*Ref); ok {
		ifBlock.True = v.Extract("data").EmitPush()
	} else {
		ifBlock.True = destType.EmitLoadFromAddr(v.Extract("data").(*aRef).Extract("data"), 0)
	}

	// false:
	ifBlock.False = NewConst("0", destType).EmitPush()

	if commaOk {
		ifBlock.Ret = append(ifBlock.Ret, wat.I32{})
		ifBlock.True = append(ifBlock.True, wat.NewInstConst(wat.I32{}, "1"))
		ifBlock.False = append(ifBlock.False, wat.NewInstConst(wat.I32{}, "0"))
	} else {
		ifBlock.False = append(ifBlock.False, wat.NewInstUnreachable())
	}

	insts = append(insts, ifBlock)
	return
}

func (v *aInterface) emitQueryInterface(destType ValueType, commaOk bool) (insts []wat.Inst) {
	insts = append(insts, v.Extract("data").EmitPush()...) //data, todo: nil check

	insts = append(insts, v.Extract("itab").EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.I32{}, 0, 4))
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(destType.Hash())))
	if commaOk {
		insts = append(insts, wat.NewInstConst(wat.I32{}, "1"))
	} else {
		insts = append(insts, wat.NewInstConst(wat.I32{}, "0"))
	}
	insts = append(insts, wat.NewInstCall("$wa.runtime.getItab")) //itab

	insts = append(insts, wat.NewInstCall("$wa.runtime.DupI32"))
	ifs := wat.NewInstIf(nil, nil, nil)
	insts = append(insts, ifs)
	if commaOk {
		ifs.Ret = []wat.ValueType{wat.I32{}, wat.I32{}}
		ifs.True = append(ifs.True, v.Extract("eq").EmitPush()...)    //eq
		ifs.True = append(ifs.True, wat.NewInstConst(wat.I32{}, "1")) //ok:true

		ifs.False = append(ifs.False, v.Extract("eq").EmitPush()...)    //eq
		ifs.False = append(ifs.False, wat.NewInstConst(wat.I32{}, "0")) //ok:false
	} else {
		ifs.False = append(ifs.False, wat.NewInstUnreachable())
		insts = append(insts, v.Extract("eq").EmitPush()...) //eq
	}

	return
}

func (v *aInterface) emitEq(r Value) (insts []wat.Inst, ok bool) {
	if !v.Type().Equal(r.Type()) {
		logger.Fatal("v.Type() != r.Type()")
	}

	d := r.(*aInterface)
	ins, _ := v.Extract("eq").emitEq(d.Extract("eq"))
	insts = append(insts, ins...)

	compEq := wat.NewInstIf(nil, nil, nil)
	compEq.Ret = append(compEq.Ret, wat.I32{})
	{
		compEq.True = append(compEq.True, v.Extract("eq").EmitPush()...)
		compEq.True = append(compEq.True, wat.NewInstConst(wat.I32{}, "-1"))
		compEq.True = append(compEq.True, wat.NewInstNe(wat.I32{}))

		compable := wat.NewInstIf(nil, nil, nil)
		compable.Ret = append(compable.Ret, wat.I32{})
		{
			compable.True = append(compable.True, v.Extract("eq").EmitPush()...)
			compable.True = append(compable.True, wat.NewInstEqz(wat.I32{}))

			isRef := wat.NewInstIf(nil, nil, nil)
			isRef.Ret = append(isRef.Ret, wat.I32{})

			ins, _ = v.Extract("data").emitEq(d.Extract("data"))
			isRef.True = ins

			isRef.False = append(isRef.False, v.Extract("data").(*aRef).Extract("data").EmitPush()...)
			isRef.False = append(isRef.False, d.Extract("data").(*aRef).Extract("data").EmitPush()...)
			isRef.False = append(isRef.False, v.Extract("eq").EmitPush()...)
			isRef.False = append(isRef.False, wat.NewInstCallIndirect("$wa.runtime.comp"))

			compable.True = append(compable.True, isRef)
		}
		compable.False = append(compable.False, wat.NewInstConst(wat.I32{}, "0"))
		compable.False = append(compable.False, wat.NewInstUnreachable())

		compEq.True = append(compEq.True, compable)
	}

	compEq.False = append(compEq.False, wat.NewInstConst(wat.I32{}, "0"))

	insts = append(insts, compEq)
	ok = true

	return
}
