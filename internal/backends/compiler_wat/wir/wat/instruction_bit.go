// 版权 @2022 凹语言 作者。保留所有权利。

package wat

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
