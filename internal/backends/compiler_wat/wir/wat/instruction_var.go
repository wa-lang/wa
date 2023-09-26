// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import "strings"

/**************************************
instGetLocal:
**************************************/
type instGetLocal struct {
	anInstruction
	name string
}

func NewInstGetLocal(name string) *instGetLocal { return &instGetLocal{name: name} }
func (i *instGetLocal) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("local.get $")
	sb.WriteString(i.name)
}

/**************************************
instSetLocal:
**************************************/
type instSetLocal struct {
	anInstruction
	name string
}

func NewInstSetLocal(name string) *instSetLocal { return &instSetLocal{name: name} }
func (i *instSetLocal) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("local.set $")
	sb.WriteString(i.name)
}

/**************************************
instGetGlobal:
**************************************/
type instGetGlobal struct {
	anInstruction
	name string
}

func NewInstGetGlobal(name string) *instGetGlobal { return &instGetGlobal{name: name} }
func (i *instGetGlobal) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("global.get $")
	sb.WriteString(i.name)
}

/**************************************
instSetGlobal:
**************************************/
type instSetGlobal struct {
	anInstruction
	name string
}

func NewInstSetGlobal(name string) *instSetGlobal { return &instSetGlobal{name: name} }
func (i *instSetGlobal) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("global.set $")
	sb.WriteString(i.name)
}
