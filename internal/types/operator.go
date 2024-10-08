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
									if len(fnDecl.Type.Params.List[0].Names) == 1 {
										typName.ops.Unary_ADD = funcs[0]
										continue
									}
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
									if len(fnDecl.Type.Params.List[0].Names) == 1 {
										typName.ops.Unary_SUB = funcs[0]
										continue
									}
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

func (check *Checker) tryFixOperatorCall(expr ast.Expr) ast.Expr {
	defer func(ctxt context, indent int) {
		check.context = ctxt
		check.indent = indent
	}(check.context, check.indent)

	check.context.ignoreFuncLitBody = true

	switch expr := expr.(type) {
	case *ast.BinaryExpr:
		var x, y operand
		check.rawExpr(&x, expr.X, nil)
		check.rawExpr(&y, expr.Y, nil)
		if check.tryBinaryOperatorCall(&x, &y, expr.X, expr.Y, expr.Op) {
			return x.expr
		}
	case *ast.UnaryExpr:
		var x operand
		check.rawExpr(&x, expr.X, nil)
		if check.tryUnaryOperatorCall(&x, expr) {
			return x.expr
		}
	}
	return nil
}

func (check *Checker) tryUnaryOperatorCall(x *operand, e *ast.UnaryExpr) bool {
	if x.typ == nil {
		return false
	}

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
	}
	if fn == nil {
		return false
	}

	if err := check.tryUnaryOpFunc(fn, x, e); err != nil {
		return false
	}

	x.mode = value
	x.typ = fn.typ.(*Signature).results.vars[0].typ

	if fn.pkg == check.pkg {
		x.expr = &ast.CallExpr{
			Fun:  &ast.Ident{Name: "#{func}:" + fn.name},
			Args: []ast.Expr{e.X},
		}
	} else {
		check.ensureOperatorCallPkgImported(fn.pkg)
		x.expr = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "#{pkg}:" + fn.pkg.path},
				Sel: &ast.Ident{Name: fn.name},
			},
			Args: []ast.Expr{e.X},
		}
	}

	check.hasCallOrRecv = true
	return true
}

func (check *Checker) tryBinaryOperatorCall(
	x, y *operand, lhs, rhs ast.Expr, op token.Token,
) bool {
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
	xFuncs := check.getBinOpFuncs(xNamed, op)
	yFuncs := check.getBinOpFuncs(yNamed, op)
	if len(xFuncs) == 0 && len(yFuncs) == 0 {
		return false
	}

	var fnMatched *Func
	for _, fn := range xFuncs {
		if err := check.tryBinOpFunc(fn, x, y, lhs, rhs); err == nil {
			fnMatched = fn
			break
		}
	}
	if fnMatched == nil {
		for _, fn := range yFuncs {
			if err := check.tryBinOpFunc(fn, x, y, lhs, rhs); err == nil {
				fnMatched = fn
				break
			}
		}
	}
	if fnMatched == nil {
		return false
	}

	x.mode = value
	x.typ = fnMatched.typ.(*Signature).results.vars[0].typ

	if fnMatched.pkg == check.pkg {
		x.expr = &ast.CallExpr{
			Fun:  &ast.Ident{Name: "#{func}:" + fnMatched.name},
			Args: []ast.Expr{lhs, rhs},
		}
	} else {
		check.ensureOperatorCallPkgImported(fnMatched.pkg)
		x.expr = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "#{pkg}:" + fnMatched.pkg.path},
				Sel: &ast.Ident{Name: fnMatched.name},
			},
			Args: []ast.Expr{lhs, rhs},
		}
	}

	switch op {
	case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
		x.expr = &ast.BinaryExpr{
			X: x.expr, Op: op, Y: &ast.BasicLit{Kind: token.INT, Value: "0"},
		}
		assert(x.typ == Typ[Int])
		x.typ = Typ[Bool]
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

func (check *Checker) tryBinOpFunc(fn *Func, x, y *operand, lhs, rhs ast.Expr) (err error) {
	firstErrBak := check.firstErr
	check.firstErr = nil

	defer func() { check.firstErr = firstErrBak }()

	defer check.handleBailout(&err)

	check.rawExpr(x, lhs, nil)
	check.rawExpr(y, rhs, nil)

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
			return typ.ops.SUB
		case token.MUL:
			return typ.ops.MUL
		case token.QUO:
			return typ.ops.QUO
		case token.REM:
			return typ.ops.REM
		}
	}
	return nil
}

// 确保运算符重载的函数对应的包被导入
func (check *Checker) ensureOperatorCallPkgImported(pkg *Package) {
	assert(pkg != check.pkg)

	pkgname := "#{pkg}:" + pkg.path
	if check.pkg.scope.Lookup(pkgname) != nil {
		return
	}

	obj := NewPkgName(token.NoPos, check.pkg, pkgname, pkg)
	obj.setNode(&ast.ImportSpec{
		Name: &ast.Ident{Name: pkgname},
		Path: &ast.BasicLit{Kind: token.STRING, Value: pkg.path},
		Comment: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "internal: only for operator overloading call"},
			},
		},
	})

	check.declare(check.pkg.scope, nil, obj, token.NoPos)
	check.pkg.imports = append(check.pkg.imports, pkg)
}
