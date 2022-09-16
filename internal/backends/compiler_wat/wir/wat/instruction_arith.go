// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import (
	"github.com/wa-lang/wa/internal/logger"
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

func NewInstAdd(t ValueType) *instAdd          { return &instAdd{typ: t} }
func (i *instAdd) Format(indent string) string { return indent + i.typ.Name() + ".add" }

/**************************************
instSub:
**************************************/
type instSub struct {
	anInstruction
	typ ValueType
}

func NewInstSub(t ValueType) *instSub          { return &instSub{typ: t} }
func (i *instSub) Format(indent string) string { return indent + i.typ.Name() + ".sub" }

/**************************************
instMul:
**************************************/
type instMul struct {
	anInstruction
	typ ValueType
}

func NewInstMul(t ValueType) *instMul          { return &instMul{typ: t} }
func (i *instMul) Format(indent string) string { return indent + i.typ.Name() + ".mul" }

/**************************************
instDiv:
**************************************/
type instDiv struct {
	anInstruction
	typ ValueType
}

func NewInstDiv(t ValueType) *instDiv { return &instDiv{typ: t} }
func (i *instDiv) Format(indent string) string {
	switch i.typ.(type) {
	case I32:
		return indent + "i32.div_s"

	case U32:
		return indent + "i32.div_u"

	case I64:
		return indent + "i64.div_s"

	case U64:
		return indent + "i64.div_u"

	case F32:
		return indent + "f32.div"

	case F64:
		return indent + "f64.div"

	}
	logger.Fatal("Todo")
	return ""
}

/**************************************
instRem:
**************************************/
type instRem struct {
	anInstruction
	typ ValueType
}

func NewInstRem(t ValueType) *instRem { return &instRem{typ: t} }
func (i *instRem) Format(indent string) string {
	switch i.typ.(type) {
	case I32:
		return indent + "i32.rem_s"

	case U32:
		return indent + "i32.rem_u"

	case I64:
		return indent + "i64.rem_s"

	case U64:
		return indent + "i64.rem_u"
	}
	logger.Fatal("Todo")
	return ""
}
