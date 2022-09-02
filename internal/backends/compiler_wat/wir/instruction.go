package wir

import (
	"strconv"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wtypes"
	"github.com/wa-lang/wa/internal/logger"
)

type anInstruction struct {
}

func (i *anInstruction) isInstruction() {}

/**************************************
InstConst:
**************************************/
type InstConst struct {
	anInstruction
	value Value
}

func NewInstConst(v Value) *InstConst {
	if v.Kind() != ValueKindConst {
		logger.Fatal("newInstructionConst()只能接受常数")
	}
	return &InstConst{value: v}
}
func (i *InstConst) Format(indent string) string {
	switch i.value.Type().(type) {
	case wtypes.Int32:
		return indent + "i32.const " + i.value.Name()

	case wtypes.Int64:
		return indent + "i64.const " + i.value.Name()
	}

	logger.Fatalf("Todo %T", i.value.Type())
	return ""
}

/**************************************
InstGetLocal:
**************************************/
type InstGetLocal struct {
	anInstruction
	name string
}

func NewInstGetLocal(name string) *InstGetLocal     { return &InstGetLocal{name: name} }
func (i *InstGetLocal) Format(indent string) string { return indent + "local.get " + i.name }

/**************************************
instSetLocal:
**************************************/
type InstSetLocal struct {
	anInstruction
	name string
}

func NewInstSetLocal(name string) *InstSetLocal     { return &InstSetLocal{name: name} }
func (i *InstSetLocal) Format(indent string) string { return indent + "local.set " + i.name }

/**************************************
InstAdd:
**************************************/
type InstAdd struct {
	anInstruction
	typ wtypes.ValueType
}

func NewInstAdd(t wtypes.ValueType) *InstAdd   { return &InstAdd{typ: t} }
func (i *InstAdd) Format(indent string) string { return indent + i.typ.String() + ".add" }

/**************************************
InstSub:
**************************************/
type InstSub struct {
	anInstruction
	typ wtypes.ValueType
}

func NewInstSub(t wtypes.ValueType) *InstSub   { return &InstSub{typ: t} }
func (i *InstSub) Format(indent string) string { return indent + i.typ.String() + ".sub" }

/**************************************
InstMul:
**************************************/
type InstMul struct {
	anInstruction
	typ wtypes.ValueType
}

func NewInstMul(t wtypes.ValueType) *InstMul   { return &InstMul{typ: t} }
func (i *InstMul) Format(indent string) string { return indent + i.typ.String() + ".mul" }

/**************************************
InstDiv:
**************************************/
type InstDiv struct {
	anInstruction
	typ wtypes.ValueType
}

func NewInstDiv(t wtypes.ValueType) *InstDiv { return &InstDiv{typ: t} }
func (i *InstDiv) Format(indent string) string {
	switch i.typ.(type) {
	case wtypes.Int32:
		return indent + "i32.div_s"
	}
	logger.Fatal("Todo")
	return ""
}

/**************************************
InstRem:
**************************************/
type InstRem struct {
	anInstruction
	typ wtypes.ValueType
}

func NewInstRem(t wtypes.ValueType) *InstRem { return &InstRem{typ: t} }
func (i *InstRem) Format(indent string) string {
	switch i.typ.(type) {
	case wtypes.Int32:
		return indent + "i32.rem_s"
	}
	logger.Fatal("Todo")
	return ""
}

/**************************************
InstEq:
**************************************/
type InstEq struct {
	anInstruction
	typ wtypes.ValueType
}

func NewInstEq(t wtypes.ValueType) *InstEq    { return &InstEq{typ: t} }
func (i *InstEq) Format(indent string) string { return indent + i.typ.String() + ".eq" }

/**************************************
InstNe:
**************************************/
type InstNe struct {
	anInstruction
	typ wtypes.ValueType
}

func NewInstNe(t wtypes.ValueType) *InstNe    { return &InstNe{typ: t} }
func (i *InstNe) Format(indent string) string { return indent + i.typ.String() + ".ne" }

/**************************************
InstLt:
**************************************/
type InstLt struct {
	anInstruction
	typ wtypes.ValueType
}

func NewInstLt(t wtypes.ValueType) *InstLt { return &InstLt{typ: t} }
func (i *InstLt) Format(indent string) string {
	switch i.typ.(type) {
	case wtypes.Int32:
		return indent + "i32.lt_s"
	}
	logger.Fatal("Todo")
	return ""
}

/**************************************
InstGt:
**************************************/
type InstGt struct {
	anInstruction
	typ wtypes.ValueType
}

func NewInstGt(t wtypes.ValueType) *InstGt { return &InstGt{typ: t} }
func (i *InstGt) Format(indent string) string {
	switch i.typ.(type) {
	case wtypes.Int32:
		return indent + "i32.gt_s"
	}
	logger.Fatal("Todo")
	return ""
}

/**************************************
InstLe:
**************************************/
type InstLe struct {
	anInstruction
	typ wtypes.ValueType
}

func NewInstLe(t wtypes.ValueType) *InstLe { return &InstLe{typ: t} }
func (i *InstLe) Format(indent string) string {
	switch i.typ.(type) {
	case wtypes.Int32:
		return indent + "i32.le_s"
	}
	logger.Fatal("Todo")
	return ""
}

/**************************************
InstGe:
**************************************/
type InstGe struct {
	anInstruction
	typ wtypes.ValueType
}

func NewInstGe(t wtypes.ValueType) *InstGe { return &InstGe{typ: t} }
func (i *InstGe) Format(indent string) string {
	switch i.typ.(type) {
	case wtypes.Int32:
		return indent + "i32.ge_s"
	}
	logger.Fatal("Todo")
	return ""
}

/**************************************
InstCall:
**************************************/
type InstCall struct {
	anInstruction
	name string
}

func NewInstCall(name string) *InstCall         { return &InstCall{name: name} }
func (i *InstCall) Format(indent string) string { return indent + "call " + i.name }

/**************************************
InstBlock:
**************************************/
type InstBlock struct {
	anInstruction
	name  string
	Insts []Instruction
}

func NewInstBlock(name string) *InstBlock { return &InstBlock{name: name} }
func (i *InstBlock) Format(indent string) string {
	s := indent + "(block "
	s += i.name + "\n"
	for _, v := range i.Insts {
		s += v.Format(indent+"  ") + "\n"
	}
	s += indent + ") ;;" + i.name
	return s
}

/**************************************
InstLoop:
**************************************/
type InstLoop struct {
	anInstruction
	name  string
	Insts []Instruction
}

func NewInstLoop(name string) *InstLoop { return &InstLoop{name: name} }
func (i *InstLoop) Format(indent string) string {
	s := indent + "(loop "
	s += i.name + "\n"
	for _, v := range i.Insts {
		s += v.Format(indent+"  ") + "\n"
	}
	s += indent + ") ;;" + i.name
	return s
}

/**************************************
InstBr:
**************************************/
type InstBr struct {
	anInstruction
	Name string
}

func NewInstBr(name string) *InstBr           { return &InstBr{Name: name} }
func (i *InstBr) Format(indent string) string { return indent + "br " + i.Name }

/**************************************
InstBrTable:
**************************************/
type InstBrTable struct {
	anInstruction
	Table []int
}

func NewInstBrTable(t []int) *InstBrTable { return &InstBrTable{Table: t} }
func (i *InstBrTable) Format(indent string) string {
	s := indent + "br_table"
	for _, v := range i.Table {
		s += " " + strconv.Itoa(v)
	}
	return s
}

/**************************************
InstIf:
**************************************/
type InstIf struct {
	anInstruction
	True  []Instruction
	False []Instruction
	Ret   wtypes.ValueType
}

func NewInstIf(instsTrue, instsFalse []Instruction, ret wtypes.ValueType) *InstIf {
	return &InstIf{True: instsTrue, False: instsFalse, Ret: ret}
}
func (i *InstIf) Format(indent string) string {
	s := indent + "if"
	if !i.Ret.Equal(wtypes.Void{}) {
		s += " (result"
		rrs := i.Ret.Raw()
		for _, rr := range rrs {
			s += " " + rr.Name()
		}
		s += ")"
	}
	s += "\n"

	for _, v := range i.True {
		s += v.Format(indent+"  ") + "\n"
	}
	s += indent + "else\n"
	for _, v := range i.False {
		s += v.Format(indent+"  ") + "\n"
	}
	s += indent + "end"
	return s
}

/**************************************
InstReturn:
**************************************/
type InstReturn struct {
	anInstruction
}

func NewInstReturn() *InstReturn                  { return &InstReturn{} }
func (i *InstReturn) Format(indent string) string { return indent + "return" }
