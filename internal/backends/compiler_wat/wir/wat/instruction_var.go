// 版权 @2022 凹语言 作者。保留所有权利。

package wat

/**************************************
InstGetLocal:
**************************************/
type InstGetLocal struct {
	anInstruction
	name string
}

func NewInstGetLocal(name string) *InstGetLocal     { return &InstGetLocal{name: name} }
func (i *InstGetLocal) Format(indent string) string { return indent + "local.get $" + i.name }

/**************************************
InstSetLocal:
**************************************/
type InstSetLocal struct {
	anInstruction
	name string
}

func NewInstSetLocal(name string) *InstSetLocal     { return &InstSetLocal{name: name} }
func (i *InstSetLocal) Format(indent string) string { return indent + "local.set $" + i.name }

/**************************************
InstGetGlobal:
**************************************/
type InstGetGlobal struct {
	anInstruction
	name string
}

func NewInstGetGlobal(name string) *InstGetGlobal    { return &InstGetGlobal{name: name} }
func (i *InstGetGlobal) Format(indent string) string { return indent + "global.get $" + i.name }

/**************************************
InstSetGlobal:
**************************************/
type InstSetGlobal struct {
	anInstruction
	name string
}

func NewInstSetGlobal(name string) *InstSetGlobal    { return &InstSetGlobal{name: name} }
func (i *InstSetGlobal) Format(indent string) string { return indent + "global.set $" + i.name }
