// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import "wa-lang.org/wa/internal/logger"

/**************************************
instAnd:
**************************************/
type instAnd struct {
	anInstruction
	typ ValueType
}

func NewInstAnd(t ValueType) *instAnd          { return &instAnd{typ: t} }
func (i *instAnd) Format(indent string) string { return indent + i.typ.Name() + ".and" }

/**************************************
instOr:
**************************************/
type instOr struct {
	anInstruction
	typ ValueType
}

func NewInstOr(t ValueType) *instOr           { return &instOr{typ: t} }
func (i *instOr) Format(indent string) string { return indent + i.typ.Name() + ".or" }

/**************************************
instXor:
**************************************/
type instXor struct {
	anInstruction
	typ ValueType
}

func NewInstXor(t ValueType) *instXor          { return &instXor{typ: t} }
func (i *instXor) Format(indent string) string { return indent + i.typ.Name() + ".xor" }

/**************************************
instShl:
**************************************/
type instShl struct {
	anInstruction
	typ ValueType
}

func NewInstShl(t ValueType) *instShl          { return &instShl{typ: t} }
func (i *instShl) Format(indent string) string { return indent + i.typ.Name() + ".shl" }

/**************************************
instShr:
**************************************/
type instShr struct {
	anInstruction
	typ ValueType
}

func NewInstShr(t ValueType) *instShr { return &instShr{typ: t} }
func (i *instShr) Format(indent string) string {
	switch i.typ.(type) {
	case I32, I64:
		return indent + i.typ.Name() + ".shr_s"

	case U32, U64:
		return indent + i.typ.Name() + ".shr_u"
	}
	logger.Fatal("Todo")
	return ""
}
