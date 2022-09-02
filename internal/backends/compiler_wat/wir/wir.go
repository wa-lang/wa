package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wtypes"
)

type Module struct {
	BaseWat string
	Imports []Import
	Funcs   []*Function
}

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
	Name   string
	Result wtypes.ValueType
	Params []Value
	Locals []Value

	Insts []Instruction
}

/**************************************
FuncSig:
**************************************/
type FuncSig struct {
	Params  []wtypes.ValueType
	Results []wtypes.ValueType
}

/**************************************
Instruction:
**************************************/
type Instruction interface {
	Format(indent string) string
	isInstruction()
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
	Type() wtypes.ValueType
	Raw() []Value
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
