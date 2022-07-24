package cir

/**************************************
Ident:
**************************************/
type Ident struct {
	Name string
}

func NewIdent(s string) *Ident {
	return &Ident{Name: s}
}

func (i *Ident) CIRString() string {
	return i.Name
}

func (i *Ident) IsExpr() {}
