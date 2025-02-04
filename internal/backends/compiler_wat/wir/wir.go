// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import "wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"

func SetCurrentModule(m *Module) { currentModule = m }

var currentModule *Module

/**************************************
Function:
**************************************/
type Function struct {
	InternalName     string
	ExternalName     string
	ExportName       string // add by chai, 临时添加, 用于 fix export name 的 js 名字问题
	ExplicitExported bool
	Results          []ValueType
	Params           []Value
	Locals           []Value

	Insts []wat.Inst
}

/**************************************
Global:
**************************************/
type Global struct {
	Name     string
	Name_exp string
	val      Value
	Type     ValueType
	init_val Value
}

type ValueKind uint8

const (
	ValueKindLocal ValueKind = iota
	ValueKindGlobal
	ValueKindConst
)

/**************************************
Value:
**************************************/
type Value interface {
	Name() string
	Kind() ValueKind
	Type() ValueType
	raw() []wat.Value
	EmitInit() []wat.Inst
	EmitPush() []wat.Inst
	EmitPushNoRetain() []wat.Inst
	EmitPop() []wat.Inst
	EmitRelease() []wat.Inst
	emitStoreToAddr(addr Value, offset int) []wat.Inst
	emitStore(offset int) []wat.Inst
	Bin() []byte

	emitEq(r Value) ([]wat.Inst, bool)
	emitCompare(r Value) []wat.Inst
}

/**************************************
ValueType:
**************************************/
type ValueType interface {
	Named() string
	Size() int
	align() int
	Kind() TypeKind
	OnFree() int
	Raw() []wat.ValueType
	Equal(ValueType) bool
	EmitLoadFromAddr(addr Value, offset int) []wat.Inst
	EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst

	Hash() int
	SetHash(h int)

	AddMethod(m Method)
	NumMethods() int
	Method(i int) Method

	typeInfoAddr() int

	OnComp() int
	setOnComp(c int)
}

/**************************************
Method:
**************************************/
type Method struct {
	Sig        FnSig
	Name       string
	FullFnName string
}
