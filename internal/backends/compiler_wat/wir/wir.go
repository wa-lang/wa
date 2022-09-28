// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import "github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"

/**************************************
Function:
**************************************/
type Function struct {
	Name   string
	Result ValueType
	Params []Value
	Locals []Value

	Insts []wat.Inst
}

/**************************************
FuncSig:
**************************************/
//type FuncSig struct {
//	Params  []ValueType
//	Results []ValueType
//}

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
	EmitPop() []wat.Inst
	EmitRelease() []wat.Inst
	emitLoad(addr Value) []wat.Inst
	emitStore(addr Value) []wat.Inst
}

/**************************************
ValueType:
**************************************/
type ValueType interface {
	byteSize() int
	Raw() []wat.ValueType
	Equal(ValueType) bool
}
