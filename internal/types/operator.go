// 版权 @2019 凹语言 作者。保留所有权利。

package types

import (
	"errors"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/ast/astutil"
	"wa-lang.org/wa/internal/token"
)

// 类型基础运算符
type typeOperator struct {
	ADD []*Func // +
	SUB []*Func // -
	MUL []*Func // *
	QUO []*Func // /
	REM []*Func // %

	Unary_ADD *Func // +x
	Unary_SUB *Func // -x
	Unary_XOR *Func // ^x
	Unary_NOT *Func // !x
}

func (check *Checker) lookupOperatorFuncs(pkg *Package, names ...string) (funcs []*Func, err error) {
	for _, s := range names {
		xObj := pkg.Scope().Lookup(s)
		if xObj == nil {
			return nil, errors.New("not found")
		}
		xFn, ok := xObj.(*Func)
		if !ok {
			return nil, errors.New("not function")
		}
		funcs = append(funcs, xFn)
	}
	return
}

// 预处运算符重载
func (check *Checker) processTypeOperators() {
	if check.conf.DisableGeneric {
		return
	}

	for obj := range check.objMap {
		// 仅自定义的新类型支持重载
		typName, _ := obj.(*TypeName)
		if typName == nil {
			continue
		}

		assert(typName.ops == nil)
		typName.ops = &typeOperator{}

		if info := astutil.ParseCommentInfo(typName.NodeDoc()); len(info.Operator) != 0 {
			for _, ops := range info.Operator {
				assert(len(ops) > 1)

				// 这里只是查询到重载的全局函数, 并未做合法性验证
				funcs, err := check.lookupOperatorFuncs(typName.Pkg(), ops[1:]...)
				if err != nil {
					check.errorf(obj.Pos(), "%s operator %v", obj.Name(), err)
					return
				}
				if len(funcs) == 0 {
					continue
				}

				switch ops[0] {
				case "^":
					assert(len(funcs) == 1)
					typName.ops.Unary_XOR = funcs[0]
				case "!":
					assert(len(funcs) == 1)
					typName.ops.Unary_NOT = funcs[0]

				case "+":
					if len(funcs) == 1 {
						if typ := funcs[0].typ; typ != nil {
							if sig := typ.(*Signature); sig.Params().Len() == 1 {
								typName.ops.Unary_ADD = funcs[0]
								continue
							}
						}
						if node := funcs[0].node; node != nil {
							if fnDecl, ok := node.(*ast.FuncDecl); ok {
								assert(fnDecl.Type != nil)
								assert(fnDecl.Type.Params != nil)
								if len(fnDecl.Type.Params.List) == 1 {
									typName.ops.Unary_ADD = funcs[0]
									continue
								}
							}
						}
					}
					typName.ops.ADD = funcs
				case "-":
					if len(funcs) == 1 {
						if typ := funcs[0].typ; typ != nil {
							if sig := typ.(*Signature); sig.Params().Len() == 1 {
								typName.ops.Unary_SUB = funcs[0]
								continue
							}
						}
						if node := funcs[0].node; node != nil {
							if fnDecl, ok := node.(*ast.FuncDecl); ok {
								if len(fnDecl.Type.Params.List) == 1 {
									typName.ops.Unary_SUB = funcs[0]
									continue
								}
							}
						}
					}
					typName.ops.SUB = funcs

				case "*":
					typName.ops.MUL = funcs
				case "/":
					typName.ops.QUO = funcs
				case "%":
					typName.ops.REM = funcs

				default:
					check.errorf(obj.Pos(), "%s operator %s invalid", obj.Name(), ops[0])
					return
				}
			}
		}
	}
}

func (check *Checker) tryUnaryOperatorCall(x *operand, e *ast.UnaryExpr) bool {
	if true {
		return false // todo(chai): debug
	}

	assert(x.typ != nil)

	var xNamed *Named
	if v, ok := x.typ.(*Named); ok {
		xNamed = v
	}
	if xNamed == nil || xNamed.obj == nil || xNamed.obj.ops == nil {
		return false
	}

	var fn *Func
	switch e.Op {
	case token.ADD:
		fn = xNamed.obj.ops.Unary_ADD
	case token.SUB:
		fn = xNamed.obj.ops.Unary_SUB
	case token.XOR:
		fn = xNamed.obj.ops.Unary_XOR
	case token.NOT:
		fn = xNamed.obj.ops.Unary_NOT
	}
	if fn == nil {
		return false
	}

	err := check.tryUnaryOpFunc(fn, x, e)
	if err == nil {
		return false
	}

	x.mode = value
	x.typ = fn.typ.(*Signature).results.vars[0].typ
	x.expr = &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: fn.pkg.path},
			Sel: &ast.Ident{Name: fn.pkg.name},
		},
		Args: []ast.Expr{e.X},
	}

	check.hasCallOrRecv = true
	return true
}

func (check *Checker) tryBinaryOperatorCall(x, y *operand, e *ast.BinaryExpr) bool {
	if true {
		return false // todo(chai): debug
	}

	var xNamed, yNamed *Named
	if v, ok := x.typ.(*Named); ok {
		xNamed = v
	}
	if v, ok := y.typ.(*Named); ok {
		yNamed = v
	}

	// 至少有1个是自定义类型
	if xNamed == nil && yNamed == nil {
		return false
	}

	// 根据左右顺序匹配
	xFuncs := check.getBinOpFuncs(xNamed, e.Op)
	yFuncs := check.getBinOpFuncs(yNamed, e.Op)

	var fnMached *Func
	if fnMached == nil {
		for _, fn := range xFuncs {
			if err := check.tryBinOpFunc(fn, x, y, e); err == nil {
				fnMached = fn
				break
			}
		}
	}
	if fnMached == nil {
		for _, fn := range yFuncs {
			if err := check.tryBinOpFunc(fn, x, y, e); err == nil {
				fnMached = fn
				break
			}
		}
	}
	if fnMached == nil {
		return false
	}

	x.mode = value
	x.typ = fnMached.typ.(*Signature).results.vars[0].typ
	x.expr = &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: fnMached.pkg.path},
			Sel: &ast.Ident{Name: fnMached.pkg.name},
		},
		Args: []ast.Expr{e.X},
	}

	check.hasCallOrRecv = true
	return true
}

func (check *Checker) tryUnaryOpFunc(fn *Func, x *operand, e *ast.UnaryExpr) (err error) {
	assert(fn != nil)

	firstErrBak := check.firstErr
	check.firstErr = nil

	defer func() { check.firstErr = firstErrBak }()

	defer check.handleBailout(&err)

	check.rawExpr(x, e.X, nil)

	assert(fn.typ != nil)
	sig, _ := fn.typ.(*Signature)
	assert(sig != nil)
	assert(sig.recv == nil)
	assert(sig.params.Len() == 1)
	assert(sig.results.Len() == 1)

	// 检查参数是否匹配
	check.argument(sig, 0, x, token.NoPos, "")
	return
}

func (check *Checker) tryBinOpFunc(fn *Func, x, y *operand, e *ast.BinaryExpr) (err error) {
	firstErrBak := check.firstErr
	check.firstErr = nil

	defer func() { check.firstErr = firstErrBak }()

	defer check.handleBailout(&err)

	check.rawExpr(x, e.X, nil)
	check.rawExpr(y, e.Y, nil)

	assert(fn.typ != nil)
	sig, _ := fn.typ.(*Signature)
	assert(sig != nil)
	assert(sig.recv == nil)
	assert(sig.params.Len() == 2)
	assert(sig.results.Len() == 1)

	// 检查参数是否匹配
	check.argument(sig, 0, x, token.NoPos, "")
	check.argument(sig, 1, y, token.NoPos, "")
	return
}

func (check *Checker) getBinOpFuncs(x *Named, op token.Token) []*Func {
	if x == nil {
		return nil
	}
	if typ := x.obj; typ != nil && typ.ops != nil {
		switch op {
		case token.ADD:
			return typ.ops.ADD
		case token.SUB:
			return typ.ops.ADD
		case token.MUL:
			return typ.ops.ADD
		case token.QUO:
			return typ.ops.ADD
		case token.REM:
			return typ.ops.ADD
		}
	}
	return nil
}
