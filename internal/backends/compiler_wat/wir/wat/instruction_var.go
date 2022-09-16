// 版权 @2022 凹语言 作者。保留所有权利。

package wat

/**************************************
instGetLocal:
**************************************/
type instGetLocal struct {
	anInstruction
	name string
}

func NewInstGetLocal(name string) *instGetLocal     { return &instGetLocal{name: name} }
func (i *instGetLocal) Format(indent string) string { return indent + "local.get $" + i.name }

/**************************************
instSetLocal:
**************************************/
type instSetLocal struct {
	anInstruction
	name string
}

func NewInstSetLocal(name string) *instSetLocal     { return &instSetLocal{name: name} }
func (i *instSetLocal) Format(indent string) string { return indent + "local.set $" + i.name }

/**************************************
instGetGlobal:
**************************************/
type instGetGlobal struct {
	anInstruction
	name string
}

func NewInstGetGlobal(name string) *instGetGlobal    { return &instGetGlobal{name: name} }
func (i *instGetGlobal) Format(indent string) string { return indent + "global.get $" + i.name }

/**************************************
instSetGlobal:
**************************************/
type instSetGlobal struct {
	anInstruction
	name string
}

func NewInstSetGlobal(name string) *instSetGlobal    { return &instSetGlobal{name: name} }
func (i *instSetGlobal) Format(indent string) string { return indent + "global.set $" + i.name }
