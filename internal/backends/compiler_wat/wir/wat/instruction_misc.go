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
