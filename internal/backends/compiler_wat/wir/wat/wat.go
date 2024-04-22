// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import "strings"

/**************************************
Import:
**************************************/
type Import interface {
	Format(indent string) string
	ModuleName() string
	ObjName() string
	Type() ObjType
}

/**************************************
Function:
**************************************/
type Function struct {
	InternalName string
	ExternalName string
	Results      []ValueType
	Params       []Value
	Locals       []Value

	Insts []Inst
}

/**************************************
Table:
**************************************/
type Table struct {
	Elems []string
}

/**************************************
FuncType:
**************************************/
type FuncType struct {
	Name string
	FuncSig
}

/**************************************
FuncSig:
**************************************/
type FuncSig struct {
	Params  []ValueType
	Results []ValueType
}

/**************************************
Inst:
**************************************/
type Inst interface {
	Format(indent string, sb *strings.Builder)
	isInstruction()
}

/**************************************
Value:
**************************************/
type Value interface {
	Name() string
	Type() ValueType
}

/**************************************
ValueType:
**************************************/
type ValueType interface {
	Name() string
	isValueType()
}

/**************************************
OpCode:
**************************************/
type OpCode int32

const (
	OpCodeAdd OpCode = iota
	OpCodeSub
	OpCodeMul
	OpCodeQuo
	OpCodeRem
	OpCodeEql
	OpCodeNe
	OpCodeLt
	OpCodeGt
	OpCodeLe
	OpCodeGe
	OpCodeComp
	OpCodeAnd
	OpCodeOr
	OpCodeXor
	OpCodeNot
	OpCodeShl
	OpCodeShr
	OpCodeAndNot
)

/**************************************
ObjType:
**************************************/
type ObjType int32

const (
	ObjTypeFunc ObjType = iota
	ObjTypeMem
	ObjTypeTable
	ObjTypeGlobal
)

/**************************************
Global:
**************************************/
type Global struct {
	V         Value
	IsMut     bool
	InitValue string
	NameExp   string
}
