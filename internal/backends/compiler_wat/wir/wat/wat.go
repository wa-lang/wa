package wat

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
	Name    string
	Results []ValueType
	Params  []Value
	Locals  []Value

	Insts []Inst
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
	Format(indent string) string
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
