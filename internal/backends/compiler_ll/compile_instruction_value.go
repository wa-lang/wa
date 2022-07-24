// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_ll

import (
	"fmt"

	"github.com/wa-lang/wa/internal/3rdparty/llir"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
	"github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
	"github.com/wa-lang/wa/internal/ssa"
	"github.com/wa-lang/wa/internal/token"
)

func (p *Compiler) getValue(expr ssa.Value) (v llvalue.Value) {
	if x, ok := p.curLocals[expr]; ok {
		return x
	}
	defer func() {
		p.curLocals[expr] = v
	}()

	switch expr := expr.(type) {
	case *ssa.Const:
		return p.compileValue_Const(expr)
	case *ssa.Function:
		return p.compileValue_Function(expr)
	case *ssa.Global:
		return p.compileValue_Global(expr)
	case *ssa.Convert:
		return p.compileValue_Convert(expr)
	case *ssa.UnOp:
		return p.compileValue_UnOp(expr)
	case *ssa.Parameter:
		return p.compileValue_Parameter(expr)
	case *ssa.BinOp:
		return p.compileValue_BinOp(expr)
	case *ssa.Phi:
		return p.compileValue_Phi(expr)
	case *ssa.Call:
		return p.compileValue_Call(expr)
	}

	fmt.Printf("todo: Compiler.getValue(expr.type=%T, expr.val=%+v)\n", expr, expr)
	return nil
}

func (p *Compiler) compileValue_Const(expr *ssa.Const) (v llvalue.Value) {
	return p.constValue(expr)
}

func (p *Compiler) compileValue_Function(expr *ssa.Function) (v llvalue.Value) {
	panic("TODO")
}

func (p *Compiler) compileValue_Global(expr *ssa.Global) (v llvalue.Value) {
	return p.getGlobal(expr)
}

func (p *Compiler) compileValue_Convert(expr *ssa.Convert) (v llvalue.Value) {
	llirBlock := p.curBlockEntries[expr.Block()]
	return llirBlock.NewTrunc(p.getValue(expr.X), p.toLLType(expr.Type()))
}

func (p *Compiler) compileValue_UnOp(expr *ssa.UnOp) (v llvalue.Value) {
	llirBlock := p.curBlockEntries[expr.Block()]
	switch expr.Op {
	case token.MUL:
		return llirBlock.NewLoad(p.toLLType(expr.Type()), p.getValue(expr.X))
	default:
		panic("TODO")
	}
}

func (p *Compiler) compileValue_Parameter(expr *ssa.Parameter) (v llvalue.Value) {
	for i, x := range expr.Parent().Params {
		if x == expr {
			return p.curFunc.Params[i]
		}
	}
	return nil
}
func (p *Compiler) compileValue_BinOp(expr *ssa.BinOp) (v llvalue.Value) {
	llirBlock := p.curBlockEntries[expr.Block()]
	switch expr.Op {
	case token.ADD:
		x := p.getValue(expr.X)
		y := p.getValue(expr.Y)
		if p.isStrType(expr.X.Type()) {
			llirBlock := p.curBlockEntries[expr.Block()]
			return llirBlock.NewCall(p.getLLFunc("ugo_cstring_join"), x, y)
		}
		if p.isFloatType(expr.X.Type()) {
			return llirBlock.NewFAdd(x, y)
		}
		return llirBlock.NewAdd(x, y)
	case token.SUB:
		x := p.getValue(expr.X)
		y := p.getValue(expr.Y)
		if p.isFloatType(expr.X.Type()) {
			return llirBlock.NewFSub(x, y)
		}
		return llirBlock.NewSub(x, y)
	case token.MUL:
		x := p.getValue(expr.X)
		y := p.getValue(expr.Y)
		if p.isFloatType(expr.X.Type()) {
			return llirBlock.NewFMul(x, y)
		}
		return llirBlock.NewMul(x, y)
	case token.QUO:
		x := p.getValue(expr.X)
		y := p.getValue(expr.Y)
		if p.isFloatType(expr.X.Type()) {
			return llirBlock.NewFDiv(x, y)
		}
		return llirBlock.NewSDiv(x, y)
	case token.REM: // %
		x := p.getValue(expr.X)
		y := p.getValue(expr.Y)
		if p.isFloatType(expr.X.Type()) {
			return llirBlock.NewFRem(x, y)
		}
		return llirBlock.NewSRem(x, y)

	case token.EQL: // ==
		x := p.getValue(expr.X)
		y := p.getValue(expr.Y)
		if p.isStrType(expr.X.Type()) {
			llirBlock := p.curBlockEntries[expr.Block()]
			cmpRet := llirBlock.NewCall(p.getLLFunc("ugo_cstring_cmp"), x, y)
			cond := llirBlock.NewICmp(llenum.IPredNE, cmpRet, llconstant.NewInt(lltypes.I32, 0))
			return cond
		}
		if p.isFloatType(expr.X.Type()) {
			return llirBlock.NewFCmp(llenum.FPredOEQ, x, y)
		}
		return llirBlock.NewICmp(llenum.IPredEQ, x, y)

	case token.NEQ: // !=
		x := p.getValue(expr.X)
		y := p.getValue(expr.Y)
		if p.isStrType(expr.X.Type()) {
			llirBlock := p.curBlockEntries[expr.Block()]
			cmpRet := llirBlock.NewCall(p.getLLFunc("ugo_cstring_cmp"), x, y)
			cond := llirBlock.NewICmp(llenum.IPredEQ, cmpRet, llconstant.NewInt(lltypes.I32, 0))
			return cond
		}
		if p.isFloatType(expr.X.Type()) {
			return llirBlock.NewFCmp(llenum.FPredONE, x, y)
		}
		return llirBlock.NewICmp(llenum.IPredNE, x, y)

	case token.LSS: // <
		x := p.getValue(expr.X)
		y := p.getValue(expr.Y)
		if p.isFloatType(expr.X.Type()) {
			return llirBlock.NewFCmp(llenum.FPredOLT, x, y)
		}
		return llirBlock.NewICmp(llenum.IPredSLT, x, y)
	case token.GTR: // >
		x := p.getValue(expr.X)
		y := p.getValue(expr.Y)
		if p.isFloatType(expr.X.Type()) {
			return llirBlock.NewFCmp(llenum.FPredOGT, x, y)
		}
		return llirBlock.NewICmp(llenum.IPredSGT, x, y)

	case token.LEQ: // <=
		x := p.getValue(expr.X)
		y := p.getValue(expr.Y)
		if p.isFloatType(expr.X.Type()) {
			return llirBlock.NewFCmp(llenum.FPredOLE, x, y)
		}
		return llirBlock.NewICmp(llenum.IPredSLE, x, y)

	case token.GEQ: // <=
		x := p.getValue(expr.X)
		y := p.getValue(expr.Y)
		if p.isFloatType(expr.X.Type()) {
			return llirBlock.NewFCmp(llenum.FPredOGE, x, y)
		}
		return llirBlock.NewICmp(llenum.IPredSGE, x, y)
	}
	panic(fmt.Sprintf("TODO: op: %+v", expr.Op))
}

func (p *Compiler) compileValue_Phi(expr *ssa.Phi) (v llvalue.Value) {
	llirBlock := p.curBlockEntries[expr.Block()]

	var incomings []*llir.Incoming
	for i, edge := range expr.Edges {
		x := p.zeroValue(edge.Type())
		pred := p.curBlockEntries[expr.Block().Preds[i]]
		incomings = append(incomings, llir.NewIncoming(x, pred))
	}

	phi := llirBlock.NewPhi(incomings...)
	p.phis[expr] = phi
	return phi
}

func (p *Compiler) compileValue_Call(expr *ssa.Call) (v llvalue.Value) {
	if builtin, ok := expr.Call.Value.(*ssa.Builtin); ok {
		return p.callBuiltin(expr, builtin)
	}

	if expr.Call.Method == nil {
		if fn, ok := expr.Call.Value.(*ssa.Function); ok {
			var mangledName string
			if fn.Pkg.Pkg.Name() == "main" {
				mangledName = fmt.Sprintf("%s.%s", fn.Pkg.Pkg.Name(), fn.Name())
			} else {
				mangledName = fmt.Sprintf("%s.%s", fn.Pkg.Pkg.Path(), fn.Name())
			}
			fn := p.getFunc(mangledName)
			var args []llvalue.Value
			for _, arg := range expr.Common().Args {
				args = append(args, p.getValue(arg))
			}

			llirBlock := p.curBlockEntries[expr.Block()]
			return llirBlock.NewCall(fn, args...)
		}
	}

	panic("TODO")
}
