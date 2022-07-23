package llutil

import (
	"log"

	constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
)

// Simplify returns an equivalent (and potentially simplified) constant to
// the constant expression.
func Simplify(c constant.Constant) constant.Constant {
	switch c := c.(type) {
	case *constant.ExprAdd:
		x, ok := Simplify(c.X).(*constant.Int)
		y, ok2 := Simplify(c.Y).(*constant.Int)
		if ok && ok2 {
			z := constant.NewInt(x.Typ, 0)
			z.X = z.X.Add(x.X, y.X)
			return z
		}
		return c
	case *constant.ExprSub:
		x, ok := Simplify(c.X).(*constant.Int)
		y, ok2 := Simplify(c.Y).(*constant.Int)
		if ok && ok2 {
			z := constant.NewInt(x.Typ, 0)
			z.X = z.X.Sub(x.X, y.X)
			return z
		}
		return c
	case *constant.ExprMul:
		x, ok := Simplify(c.X).(*constant.Int)
		y, ok2 := Simplify(c.Y).(*constant.Int)
		if ok && ok2 {
			z := constant.NewInt(x.Typ, 0)
			z.X = z.X.Mul(x.X, y.X)
			return z
		}
		return c
	case *constant.ExprSDiv:
		// TODO: if we need to handle signed division differently from unsigned division.
		x, ok := Simplify(c.X).(*constant.Int)
		y, ok2 := Simplify(c.Y).(*constant.Int)
		if ok && ok2 {
			z := constant.NewInt(x.Typ, 0)
			z.X = z.X.Div(x.X, y.X)
			return z
		}
		return c
	case *constant.ExprUDiv:
		x, ok := Simplify(c.X).(*constant.Int)
		y, ok2 := Simplify(c.Y).(*constant.Int)
		if ok && ok2 {
			z := constant.NewInt(x.Typ, 0)
			z.X = z.X.Div(x.X, y.X)
			return z
		}
		return c
	case *constant.ExprFAdd:
		x, ok := Simplify(c.X).(*constant.Float)
		y, ok2 := Simplify(c.Y).(*constant.Float)
		if ok && ok2 {
			z := constant.NewFloat(x.Typ, 0)
			z.X = z.X.Add(x.X, y.X)
			return z
		}
		return c
	case *constant.ExprFSub:
		x, ok := Simplify(c.X).(*constant.Float)
		y, ok2 := Simplify(c.Y).(*constant.Float)
		if ok && ok2 {
			z := constant.NewFloat(x.Typ, 0)
			z.X = z.X.Sub(x.X, y.X)
			return z
		}
		return c
	case *constant.ExprFMul:
		x, ok := Simplify(c.X).(*constant.Float)
		y, ok2 := Simplify(c.Y).(*constant.Float)
		if ok && ok2 {
			z := constant.NewFloat(x.Typ, 0)
			z.X = z.X.Mul(x.X, y.X)
			return z
		}
		return c
	default:
		log.Printf("support for simplifying constant expression %T not yet implemented; returning original constant expression", c)
		return c
	}
}
