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

func (check *Checker) tryUnaryOperatorCall(
	x *operand, e *ast.UnaryExpr,
	op token.Token,
) bool {
	return false // todo(chai)
}

func (check *Checker) tryBinaryOperatorCall(
	x *operand, e *ast.BinaryExpr,
	lhs, rhs ast.Expr,
	op token.Token,
) bool {
	return false // todo(chai)
}
