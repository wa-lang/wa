// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import "wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"

func SetCurrentModule(m *Module) { currentModule = m }

var currentModule *Module

/**************************************
Function:
**************************************/
type Function struct {
	Name    string
	Results []ValueType
	Params  []Value
	Locals  []Value

	Insts []wat.Inst
}

type ValueKind uint8

const (
	ValueKindLocal ValueKind = iota
	ValueKindGlobal_Value
	ValueKindGlobal_Pointer
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
	EmitPop() []wat.Inst
	EmitRelease() []wat.Inst
	emitStoreToAddr(addr Value, offset int) []wat.Inst
}

/**************************************
ValueType:
**************************************/
type ValueType interface {
	Name() string
	size() int
	align() int
	onFree() int
	Raw() []wat.ValueType
	Equal(ValueType) bool
	EmitLoadFromAddr(addr Value, offset int) []wat.Inst
}
