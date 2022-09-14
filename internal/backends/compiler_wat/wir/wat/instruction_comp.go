// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import "github.com/wa-lang/wa/internal/logger"

/**************************************
InstEq:
**************************************/
type InstEq struct {
	anInstruction
	typ ValueType
}

func NewInstEq(t ValueType) *InstEq           { return &InstEq{typ: t} }
func (i *InstEq) Format(indent string) string { return indent + i.typ.Name() + ".eq" }

/**************************************
InstNe:
**************************************/
type InstNe struct {
	anInstruction
	typ ValueType
}

func NewInstNe(t ValueType) *InstNe           { return &InstNe{typ: t} }
func (i *InstNe) Format(indent string) string { return indent + i.typ.Name() + ".ne" }

/**************************************
InstLt:
**************************************/
type InstLt struct {
	anInstruction
	typ ValueType
}

func NewInstLt(t ValueType) *InstLt { return &InstLt{typ: t} }
func (i *InstLt) Format(indent string) string {
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
InstGt:
**************************************/
type InstGt struct {
	anInstruction
	typ ValueType
}

func NewInstGt(t ValueType) *InstGt { return &InstGt{typ: t} }
func (i *InstGt) Format(indent string) string {
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
InstLe:
**************************************/
type InstLe struct {
	anInstruction
	typ ValueType
}

func NewInstLe(t ValueType) *InstLe { return &InstLe{typ: t} }
func (i *InstLe) Format(indent string) string {
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
InstGe:
**************************************/
type InstGe struct {
	anInstruction
	typ ValueType
}

func NewInstGe(t ValueType) *InstGe { return &InstGe{typ: t} }
func (i *InstGe) Format(indent string) string {
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
