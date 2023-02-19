// 版权 @2022 凹语言 作者。保留所有权利。

package wat

/**************************************
instXor:
**************************************/
type instXor struct {
	anInstruction
	typ ValueType
}

func NewInstXor(t ValueType) *instXor          { return &instXor{typ: t} }
func (i *instXor) Format(indent string) string { return indent + i.typ.Name() + ".xor" }
