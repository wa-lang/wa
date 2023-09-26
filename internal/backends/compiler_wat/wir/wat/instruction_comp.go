// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import (
	"strings"

	"wa-lang.org/wa/internal/logger"
)

/**************************************
instEqz:
**************************************/
type instEqz struct {
	anInstruction
	typ ValueType
}

func NewInstEqz(t ValueType) *instEqz { return &instEqz{typ: t} }
func (i *instEqz) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString(i.typ.Name())
	sb.WriteString(".eqz")
}

/**************************************
instEq:
**************************************/
type instEq struct {
	anInstruction
	typ ValueType
}

func NewInstEq(t ValueType) *instEq { return &instEq{typ: t} }
func (i *instEq) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString(i.typ.Name())
	sb.WriteString(".eq")
}

/**************************************
instNe:
**************************************/
type instNe struct {
	anInstruction
	typ ValueType
}

func NewInstNe(t ValueType) *instNe { return &instNe{typ: t} }
func (i *instNe) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString(i.typ.Name())
	sb.WriteString(".ne")
}

/**************************************
instLt:
**************************************/
type instLt struct {
	anInstruction
	typ ValueType
}

func NewInstLt(t ValueType) *instLt { return &instLt{typ: t} }
func (i *instLt) Format(indent string, sb *strings.Builder) {
	switch i.typ.(type) {
	case I32:
		sb.WriteString(indent)
		sb.WriteString("i32.lt_s")

	case U32:
		sb.WriteString(indent)
		sb.WriteString("i32.lt_u")

	case I64:
		sb.WriteString(indent)
		sb.WriteString("i64.lt_s")

	case U64:
		sb.WriteString(indent)
		sb.WriteString("i64.lt_u")

	case F32:
		sb.WriteString(indent)
		sb.WriteString("f32.lt")

	case F64:
		sb.WriteString(indent)
		sb.WriteString("f64.lt")

	default:
		logger.Fatal("Todo")
	}
}

/**************************************
instGt:
**************************************/
type instGt struct {
	anInstruction
	typ ValueType
}

func NewInstGt(t ValueType) *instGt { return &instGt{typ: t} }
func (i *instGt) Format(indent string, sb *strings.Builder) {
	switch i.typ.(type) {
	case I32:
		sb.WriteString(indent)
		sb.WriteString("i32.gt_s")

	case U32:
		sb.WriteString(indent)
		sb.WriteString("i32.gt_u")

	case I64:
		sb.WriteString(indent)
		sb.WriteString("i64.gt_s")

	case U64:
		sb.WriteString(indent)
		sb.WriteString("i64.gt_u")

	case F32:
		sb.WriteString(indent)
		sb.WriteString("f32.gt")

	case F64:
		sb.WriteString(indent)
		sb.WriteString("f64.gt")

	default:
		logger.Fatal("Todo")
	}
}

/**************************************
instLe:
**************************************/
type instLe struct {
	anInstruction
	typ ValueType
}

func NewInstLe(t ValueType) *instLe { return &instLe{typ: t} }
func (i *instLe) Format(indent string, sb *strings.Builder) {
	switch i.typ.(type) {
	case I32:
		sb.WriteString(indent)
		sb.WriteString("i32.le_s")

	case U32:
		sb.WriteString(indent)
		sb.WriteString("i32.le_u")

	case I64:
		sb.WriteString(indent)
		sb.WriteString("i64.le_s")

	case U64:
		sb.WriteString(indent)
		sb.WriteString("i64.le_u")

	case F32:
		sb.WriteString(indent)
		sb.WriteString("f32.le")

	case F64:
		sb.WriteString(indent)
		sb.WriteString("f64.le")

	default:
		logger.Fatal("Todo")
	}
}

/**************************************
instGe:
**************************************/
type instGe struct {
	anInstruction
	typ ValueType
}

func NewInstGe(t ValueType) *instGe { return &instGe{typ: t} }
func (i *instGe) Format(indent string, sb *strings.Builder) {
	switch i.typ.(type) {
	case I32, I64:
		sb.WriteString(indent)
		sb.WriteString(i.typ.Name())
		sb.WriteString(".ge_s")

	case U32, U64:
		sb.WriteString(indent)
		sb.WriteString(i.typ.Name())
		sb.WriteString(".ge_u")

	case F32, F64:
		sb.WriteString(indent)
		sb.WriteString(i.typ.Name())
		sb.WriteString(".ge")

	default:
		logger.Fatal("Todo")
	}
}
