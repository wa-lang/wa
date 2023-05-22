package wat

/**************************************
instDrop:
**************************************/
type instDrop struct {
	anInstruction
}

func NewInstDrop() *instDrop                    { return &instDrop{} }
func (i *instDrop) Format(indent string) string { return indent + "drop" }

/**************************************
comment:
**************************************/
type comment struct {
	anInstruction
	name string
}

func NewComment(name string) *comment          { return &comment{name: name} }
func (i *comment) Format(indent string) string { return indent + ";;" + i.name }

/**************************************
blank:
**************************************/
type blank struct {
	anInstruction
}

func NewBlank() *blank                       { return &blank{} }
func (i *blank) Format(indent string) string { return "" }

/**************************************
instUnreachable:
**************************************/
type instUnreachable struct {
	anInstruction
}

func NewInstinstUnreachable() *instUnreachable         { return &instUnreachable{} }
func (i *instUnreachable) Format(indent string) string { return indent + "unreachable" }
