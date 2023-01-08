package cir

import (
	"wa-lang.org/wa/internal/backends/compiler_c/cir/ctypes"
	"wa-lang.org/wa/internal/logger"
)

/**************************************
CIRStringer:
**************************************/
type CIRStringer interface {
	CIRString() string
}

/**************************************
Expr:
**************************************/
type Expr interface {
	CIRStringer
	Type() ctypes.Type
}

/**************************************
NotExpr:
**************************************/
type NotExpr struct {
	X Expr
}

func NewNotExpr(x Expr) *NotExpr {
	return &NotExpr{X: x}
}

func (e *NotExpr) CIRString() string {
	return "!" + e.X.CIRString()
}

func (e *NotExpr) Type() ctypes.Type {
	return ctypes.Bool
}

/**************************************
GetaddrExpr:
**************************************/
type GetaddrExpr struct {
	X Expr
}

func NewGetaddrExpr(x Expr) *GetaddrExpr {
	return &GetaddrExpr{X: x}
}

func (e *GetaddrExpr) CIRString() string {
	return "&" + e.X.CIRString()
}

func (e *GetaddrExpr) Type() ctypes.Type {
	return ctypes.NewPointerType(e.X.Type())
}

/**************************************
LoadExpr:
**************************************/
type LoadExpr struct {
	x Expr
}

func NewLoadExpr(expr Expr) *LoadExpr {
	switch expr.Type().(type) {
	case *ctypes.PointerType:
		return &LoadExpr{x: expr}

	case *ctypes.RefType:
		return &LoadExpr{x: expr}
	}

	logger.Fatal("expr is not a valid pointer or reference:", expr)
	return nil
}

func (e *LoadExpr) CIRString() string {
	switch e.x.Type().(type) {
	case *ctypes.PointerType:
		return "*" + e.x.CIRString()

	default:
		return e.x.CIRString() + ".GetValue()"
	}
}

func (e *LoadExpr) Type() ctypes.Type {
	switch t := e.x.Type().(type) {
	case *ctypes.PointerType:
		return t.Base

	case *ctypes.RefType:
		return t.Base
	}

	return nil
}

/**************************************
RawExpr:
**************************************/
type RawExpr struct {
	raw string
	typ ctypes.Type
}

func NewRawExpr(r string, t ctypes.Type) *RawExpr {
	return &RawExpr{raw: r, typ: t}
}

func (e *RawExpr) CIRString() string {
	return e.raw
}

func (e *RawExpr) Type() ctypes.Type {
	return e.typ
}

/**************************************
AddExpr:
**************************************/
type AddExpr struct {
	X Expr
	Y Expr
}

func NewAddExpr(x, y Expr) *AddExpr {
	return &AddExpr{X: x, Y: y}
}

func (e *AddExpr) CIRString() string {
	return "(" + e.X.CIRString() + " + " + e.Y.CIRString() + ")"
}

func (e *AddExpr) Type() ctypes.Type {
	return e.X.Type()
}

/**************************************
SubExpr:
**************************************/
type SubExpr struct {
	X Expr
	Y Expr
}

func NewSubExpr(x, y Expr) *SubExpr {
	return &SubExpr{X: x, Y: y}
}

func (e *SubExpr) CIRString() string {
	return "(" + e.X.CIRString() + " - " + e.Y.CIRString() + ")"
}

func (e *SubExpr) Type() ctypes.Type {
	return e.X.Type()
}

/**************************************
MulExpr:
**************************************/
type MulExpr struct {
	X Expr
	Y Expr
}

func NewMulExpr(x, y Expr) *MulExpr {
	return &MulExpr{X: x, Y: y}
}

func (e *MulExpr) CIRString() string {
	return "(" + e.X.CIRString() + " * " + e.Y.CIRString() + ")"
}

func (e *MulExpr) Type() ctypes.Type {
	return e.X.Type()
}

/**************************************
QuoExpr:
**************************************/
type QuoExpr struct {
	X Expr
	Y Expr
}

func NewQuoExpr(x, y Expr) *QuoExpr {
	return &QuoExpr{X: x, Y: y}
}

func (e *QuoExpr) CIRString() string {
	return "(" + e.X.CIRString() + " / " + e.Y.CIRString() + ")"
}

func (e *QuoExpr) Type() ctypes.Type {
	return e.X.Type()
}

/**************************************
EqlExpr:
**************************************/
type EqlExpr struct {
	X Expr
	Y Expr
}

func NewEqlExpr(x, y Expr) *EqlExpr {
	return &EqlExpr{X: x, Y: y}
}

func (e *EqlExpr) CIRString() string {
	return "(" + e.X.CIRString() + " == " + e.Y.CIRString() + ")"
}

func (e *EqlExpr) Type() ctypes.Type {
	return ctypes.Bool
}

/**************************************
CallExpr:
**************************************/
type CallExpr struct {
	Func Expr
	Args []Expr
}

func NewCallExpr(fn Expr, args []Expr) *CallExpr {
	if _, ok := fn.Type().(*ctypes.FuncType); !ok {
		logger.Fatal("func is not a valid Function")
		return nil
	}
	return &CallExpr{Func: fn, Args: args}
}

func (s *CallExpr) CIRString() string {
	str := s.Func.CIRString()
	str += "("
	for i, v := range s.Args {
		str += v.CIRString()
		if i < len(s.Args)-1 {
			str += ", "
		}
	}
	str += ")"
	return str
}

func (s *CallExpr) Type() ctypes.Type {
	return s.Func.Type().(*ctypes.FuncType).Ret
}

/**************************************
SelectExpr:
**************************************/
type SelectExpr struct {
	x, y Expr
}

func NewSelectExpr(x, y Expr) *SelectExpr {
	return &SelectExpr{x: x, y: y}
}

func (s *SelectExpr) CIRString() string {
	return s.x.CIRString() + "." + s.y.CIRString()
}

func (s *SelectExpr) Type() ctypes.Type {
	return s.y.Type()
}

/**************************************
IndexAddrExpr:
**************************************
type IndexAddrExpr struct {
	x, index Expr
	typ      ctypes.Type
}

func NewIndexAddrExpr(x, index Expr) *IndexAddrExpr {
	switch t := x.Type().(type) {
	case *ctypes.Array:
		return &IndexAddrExpr{x: x, index: index, typ: ctypes.NewPointerType(t.GetElem())}

	default:
		logger.Fatalf("Todo: NewIndexAddrExpr(), %T", t)
	}

	return nil
}

func (s *IndexAddrExpr) CIRString() string {
	return "&" + s.x.CIRString() + "[" + s.index.CIRString() + "]"
}

func (s *IndexAddrExpr) Type() ctypes.Type {
	return s.typ
}

*/
