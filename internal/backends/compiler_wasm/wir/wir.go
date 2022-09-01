package wir

import "github.com/wa-lang/wa/internal/backends/compiler_wasm/wir/wtypes"

type Module struct {
	Funcs []Function
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
