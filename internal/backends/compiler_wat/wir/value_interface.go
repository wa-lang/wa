// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
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
	interface_t.underlying.AppendField(m.NewStructField("data", m.GenValueType_Ref(m.GenValueType_Ref(m.VOID))))
	interface_t.underlying.AppendField(m.NewStructField("itab", m.U32))
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

func (t *Interface) emitGenMake(x Value, x_ref *Ref) (insts []wat.Inst) {
	insts = append(insts, x_ref.emitHeapAlloc()...)
	insts = append(insts, x.emitStore(0)...)

	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(x.Type().Hash())))
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(t.Hash())))
	insts = append(insts, wat.NewInstConst(wat.I32{}, "0"))
	insts = append(insts, wat.NewInstCall("$wa.RT.getItab"))

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
