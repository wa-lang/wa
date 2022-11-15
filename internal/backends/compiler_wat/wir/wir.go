// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import "github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"

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

/**************************************
FuncSig:
**************************************/
type FuncSig struct {
	Params  []ValueType
	Results []ValueType
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
	emitLoadFromAddr(addr Value, offset int) []wat.Inst
}
