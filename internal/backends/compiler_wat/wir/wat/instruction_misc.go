package wat

import "strings"

/**************************************
instDrop:
**************************************/
type instDrop struct {
	anInstruction
}

func NewInstDrop() *instDrop { return &instDrop{} }
func (i *instDrop) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("drop")
}

/**************************************
comment:
**************************************/
type comment struct {
	anInstruction
	name string
}

func NewComment(name string) *comment { return &comment{name: name} }
func (i *comment) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString(";;")
	sb.WriteString(i.name)
}

/**************************************
blank:
**************************************/
type blank struct {
	anInstruction
}

func NewBlank() *blank                                     { return &blank{} }
func (i *blank) Format(indent string, sb *strings.Builder) {}

/**************************************
instUnreachable:
**************************************/
type instUnreachable struct {
	anInstruction
}

func NewInstUnreachable() *instUnreachable { return &instUnreachable{} }
func (i *instUnreachable) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("unreachable")
}
