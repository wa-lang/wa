package cir

import "wa-lang.org/wa/internal/backends/compiler_c/cir/ctypes"

/**************************************
Var:
**************************************/
type Var struct {
	Name string
	typ  ctypes.Type

	AssociatedSSAObj interface{}
}

func (i *Var) CIRString() string {
	return i.Name
}

func NewVar(n string, t ctypes.Type) *Var {
	return &Var{Name: n, typ: t}
}

func (i *Var) Type() ctypes.Type {
	return i.typ
}
