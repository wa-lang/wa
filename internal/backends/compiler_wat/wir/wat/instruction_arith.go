// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import (
	"strings"

	"wa-lang.org/wa/internal/logger"
)

type anInstruction struct {
}

func (i *anInstruction) isInstruction() {}

/**************************************
instAdd:
**************************************/
type instAdd struct {
	anInstruction
	typ ValueType
}

func NewInstAdd(t ValueType) *instAdd { return &instAdd{typ: t} }
func (i *instAdd) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString(i.typ.Name())
	sb.WriteString(".add")
}

/**************************************
instSub:
**************************************/
type instSub struct {
	anInstruction
	typ ValueType
}

func NewInstSub(t ValueType) *instSub { return &instSub{typ: t} }
func (i *instSub) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString(i.typ.Name())
	sb.WriteString(".sub")
}

/**************************************
instMul:
**************************************/
type instMul struct {
	anInstruction
	typ ValueType
}

func NewInstMul(t ValueType) *instMul { return &instMul{typ: t} }
func (i *instMul) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString(i.typ.Name())
	sb.WriteString(".mul")
}

/**************************************
instDiv:
**************************************/
type instDiv struct {
	anInstruction
	typ ValueType
}

func NewInstDiv(t ValueType) *instDiv { return &instDiv{typ: t} }
func (i *instDiv) Format(indent string, sb *strings.Builder) {
	switch i.typ.(type) {
	case I32:
		sb.WriteString(indent)
		sb.WriteString("i32.div_s")

	case U32:
		sb.WriteString(indent)
		sb.WriteString("i32.div_u")

	case I64:
		sb.WriteString(indent)
		sb.WriteString("i64.div_s")

	case U64:
		sb.WriteString(indent)
		sb.WriteString("i64.div_u")

	case F32:
		sb.WriteString(indent)
		sb.WriteString("f32.div")

	case F64:
		sb.WriteString(indent)
		sb.WriteString("f64.div")

	default:
		logger.Fatal("Todo")
	}
}

/**************************************
instRem:
**************************************/
type instRem struct {
	anInstruction
	typ ValueType
}

func NewInstRem(t ValueType) *instRem { return &instRem{typ: t} }
func (i *instRem) Format(indent string, sb *strings.Builder) {
	switch i.typ.(type) {
	case I32:
		sb.WriteString(indent)
		sb.WriteString("i32.rem_s")

	case U32:
		sb.WriteString(indent)
		sb.WriteString("i32.rem_u")

	case I64:
		sb.WriteString(indent)
		sb.WriteString("i64.rem_s")

	case U64:
		sb.WriteString(indent)
		sb.WriteString("i64.rem_u")

	default:
		logger.Fatal("Todo")
	}
}
