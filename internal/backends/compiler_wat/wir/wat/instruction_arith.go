// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import (
	"github.com/wa-lang/wa/internal/logger"
)

type anInstruction struct {
}

func (i *anInstruction) isInstruction() {}

/**************************************
InstAdd:
**************************************/
type InstAdd struct {
	anInstruction
	typ ValueType
}

func NewInstAdd(t ValueType) *InstAdd          { return &InstAdd{typ: t} }
func (i *InstAdd) Format(indent string) string { return indent + i.typ.Name() + ".add" }

/**************************************
InstSub:
**************************************/
type InstSub struct {
	anInstruction
	typ ValueType
}

func NewInstSub(t ValueType) *InstSub          { return &InstSub{typ: t} }
func (i *InstSub) Format(indent string) string { return indent + i.typ.Name() + ".sub" }

/**************************************
InstMul:
**************************************/
type InstMul struct {
	anInstruction
	typ ValueType
}

func NewInstMul(t ValueType) *InstMul          { return &InstMul{typ: t} }
func (i *InstMul) Format(indent string) string { return indent + i.typ.Name() + ".mul" }

/**************************************
InstDiv:
**************************************/
type InstDiv struct {
	anInstruction
	typ ValueType
}

func NewInstDiv(t ValueType) *InstDiv { return &InstDiv{typ: t} }
func (i *InstDiv) Format(indent string) string {
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
InstRem:
**************************************/
type InstRem struct {
	anInstruction
	typ ValueType
}

func NewInstRem(t ValueType) *InstRem { return &InstRem{typ: t} }
func (i *InstRem) Format(indent string) string {
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
