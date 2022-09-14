// 版权 @2022 凹语言 作者。保留所有权利。

package wat

/**************************************
InstConst:
**************************************/
type InstConst struct {
	anInstruction
	typ     ValueType
	literal string
}

func NewInstConst(typ ValueType, literal string) *InstConst {
	return &InstConst{typ: typ, literal: literal}
}
func (i *InstConst) Format(indent string) string {
	return indent + i.typ.Name() + ".const " + i.literal
}
