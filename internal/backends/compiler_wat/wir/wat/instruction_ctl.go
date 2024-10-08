// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import (
	"fmt"
	"strconv"
	"strings"

	"wa-lang.org/wa/internal/logger"
)

/**************************************
instCall:
**************************************/
type instCall struct {
	anInstruction
	name string
}

func NewInstCall(name string) *instCall { return &instCall{name: name} }
func (i *instCall) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("call $")
	sb.WriteString(i.name)
}

/**************************************
instCallIndirect:
**************************************/
type instCallIndirect struct {
	anInstruction
	func_type string
}

func NewInstCallIndirect(func_type string) *instCallIndirect {
	return &instCallIndirect{func_type: func_type}
}
func (i *instCallIndirect) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("call_indirect (type $")
	sb.WriteString(i.func_type)
	sb.WriteString(")")
}

/**************************************
instBlock:
**************************************/
type instBlock struct {
	anInstruction
	name  string
	Insts []Inst
	Ret   []ValueType
}

func NewInstBlock(name string) *instBlock { return &instBlock{name: name} }
func (i *instBlock) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("block")

	if len(i.Ret) > 0 {
		sb.WriteString(" (result")
		for _, r := range i.Ret {
			sb.WriteString(" ")
			sb.WriteString(r.Name())
		}
		sb.WriteString(")")
	}

	if len(i.name) > 0 {
		sb.WriteString(" $")
		sb.WriteString(i.name)
	}

	sb.WriteString("\n")

	indent_t := indent + "  "
	for _, v := range i.Insts {
		v.Format(indent_t, sb)
		sb.WriteString("\n")
	}

	sb.WriteString(indent)
	sb.WriteString("end ;;")
	sb.WriteString(i.name)
}

/**************************************
instLoop:
**************************************/
type instLoop struct {
	anInstruction
	name  string
	Insts []Inst
}

func NewInstLoop(name string) *instLoop { return &instLoop{name: name} }
func (i *instLoop) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("loop $")

	sb.WriteString(i.name)
	sb.WriteString("\n")

	indent_t := indent + "  "
	for _, v := range i.Insts {
		v.Format(indent_t, sb)
		sb.WriteString("\n")
	}

	sb.WriteString(indent)
	sb.WriteString("end ;;")
	sb.WriteString(i.name)
}

/**************************************
instBr:
**************************************/
type instBr struct {
	anInstruction
	label interface{}
}

func NewInstBr(label interface{}) *instBr { return &instBr{label: label} }
func (i *instBr) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	switch l := i.label.(type) {
	case string:
		sb.WriteString("br $")
		sb.WriteString(l)

	case int:
		sb.WriteString("br ")
		sb.WriteString(fmt.Sprint(l))

	default:
		logger.Fatalf("Invalid br type:%T", l)
	}
}

/**************************************
instBrIf:
**************************************/
type instBrIf struct {
	anInstruction
	label interface{}
}

func NewInstBrIf(label interface{}) *instBrIf { return &instBrIf{label: label} }
func (i *instBrIf) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	switch l := i.label.(type) {
	case string:
		sb.WriteString("br_if $")
		sb.WriteString(l)

	case int:
		sb.WriteString("br_if ")
		sb.WriteString(fmt.Sprint(l))

	default:
		logger.Fatalf("Invalid br_if type:%T", l)
	}
}

/**************************************
instBrTable:
**************************************/
type instBrTable struct {
	anInstruction
	Table []int
}

func NewInstBrTable(t []int) *instBrTable { return &instBrTable{Table: t} }
func (i *instBrTable) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("br_table")
	for _, v := range i.Table {
		sb.WriteString(" ")
		sb.WriteString(strconv.Itoa(v))
	}
}

/**************************************
instIf:
**************************************/
type instIf struct {
	anInstruction
	True  []Inst
	False []Inst
	Ret   []ValueType
}

func NewInstIf(instsTrue, instsFalse []Inst, ret []ValueType) *instIf {
	return &instIf{True: instsTrue, False: instsFalse, Ret: ret}
}
func (i *instIf) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("if")

	if len(i.Ret) > 0 {
		sb.WriteString(" (result")
		for _, r := range i.Ret {
			sb.WriteString(" ")
			sb.WriteString(r.Name())
		}
		sb.WriteString(")")
	}
	sb.WriteString("\n")

	indent_t := indent + "  "
	for _, v := range i.True {
		v.Format(indent_t, sb)
		sb.WriteString("\n")
	}
	sb.WriteString(indent)
	sb.WriteString("else\n")
	for _, v := range i.False {
		v.Format(indent_t, sb)
		sb.WriteString("\n")
	}

	sb.WriteString(indent)
	sb.WriteString("end")
}

/**************************************
instReturn:
**************************************/
type instReturn struct {
	anInstruction
}

func NewInstReturn() *instReturn { return &instReturn{} }
func (i *instReturn) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("return")
}
