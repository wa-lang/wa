// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import "github.com/wa-lang/wa/internal/logger"

/**************************************
instEqz:
**************************************/
type instEqz struct {
	anInstruction
	typ ValueType
}

func NewInstEqz(t ValueType) *instEqz          { return &instEqz{typ: t} }
func (i *instEqz) Format(indent string) string { return indent + i.typ.Name() + ".eqz" }

/**************************************
instEq:
**************************************/
type instEq struct {
	anInstruction
	typ ValueType
}

func NewInstEq(t ValueType) *instEq           { return &instEq{typ: t} }
func (i *instEq) Format(indent string) string { return indent + i.typ.Name() + ".eq" }

/**************************************
instNe:
**************************************/
type instNe struct {
	anInstruction
	typ ValueType
}

func NewInstNe(t ValueType) *instNe           { return &instNe{typ: t} }
func (i *instNe) Format(indent string) string { return indent + i.typ.Name() + ".ne" }

/**************************************
instLt:
**************************************/
type instLt struct {
	anInstruction
	typ ValueType
}

func NewInstLt(t ValueType) *instLt { return &instLt{typ: t} }
func (i *instLt) Format(indent string) string {
	switch i.typ.(type) {
	case I32:
		return indent + "i32.lt_s"

	case U32:
		return indent + "i32.lt_u"

	case I64:
		return indent + "i64.lt_s"

	case U64:
		return indent + "i64.lt_u"

	case F32:
		return indent + "f32.lt"

	case F64:
		return indent + "f64.lt"
	}
	logger.Fatal("Todo")
	return ""
}

/**************************************
instGt:
**************************************/
type instGt struct {
	anInstruction
	typ ValueType
}

func NewInstGt(t ValueType) *instGt { return &instGt{typ: t} }
func (i *instGt) Format(indent string) string {
	switch i.typ.(type) {
	case I32:
		return indent + "i32.gt_s"

	case U32:
		return indent + "i32.gt_u"

	case I64:
		return indent + "i64.gt_s"

	case U64:
		return indent + "i64.gt_u"

	case F32:
		return indent + "f32.gt"

	case F64:
		return indent + "f64.gt"
	}
	logger.Fatal("Todo")
	return ""
}

/**************************************
instLe:
**************************************/
type instLe struct {
	anInstruction
	typ ValueType
}

func NewInstLe(t ValueType) *instLe { return &instLe{typ: t} }
func (i *instLe) Format(indent string) string {
	switch i.typ.(type) {
	case I32:
		return indent + "i32.le_s"

	case U32:
		return indent + "i32.le_u"

	case I64:
		return indent + "i64.le_s"

	case U64:
		return indent + "i64.le_u"

	case F32:
		return indent + "f32.le"

	case F64:
		return indent + "f64.le"
	}
	logger.Fatal("Todo")
	return ""
}

/**************************************
instGe:
**************************************/
type instGe struct {
	anInstruction
	typ ValueType
}

func NewInstGe(t ValueType) *instGe { return &instGe{typ: t} }
func (i *instGe) Format(indent string) string {
	switch i.typ.(type) {
	case I32:
		return indent + "i32.ge_s"

	case U32:
		return indent + "i32.ge_u"

	case I64:
		return indent + "i64.ge_s"

	case U64:
		return indent + "i64.ge_u"

	case F32:
		return indent + "f32.ge"

	case F64:
		return indent + "f64.ge"
	}
	logger.Fatal("Todo")
	return ""
}
