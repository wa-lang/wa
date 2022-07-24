// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_c

import (
	"strconv"

	"github.com/wa-lang/wa/internal/backends/compiler_c/cir"
	"github.com/wa-lang/wa/internal/backends/compiler_c/cir/cconstant"
	"github.com/wa-lang/wa/internal/backends/compiler_c/cir/ctypes"
	"github.com/wa-lang/wa/internal/constant"
	"github.com/wa-lang/wa/internal/logger"
	"github.com/wa-lang/wa/internal/ssa"
	"github.com/wa-lang/wa/internal/token"
	"github.com/wa-lang/wa/internal/types"
)

type functionGenerator struct {
	compiler *CompilerC

	locals  map[ssa.Value]cir.Var
	localID int

	params        []cir.VarDecl
	var_cur_block cir.Expr
}

func newFunctionGenerator(p *CompilerC) *functionGenerator {
	return &functionGenerator{compiler: p, locals: make(map[ssa.Value]cir.Var)}
}

func (g *functionGenerator) getValue(i interface{}) cir.Expr {
	for scope := g.compiler.curScope; scope != nil; scope = scope.Parent {
		for _, v := range scope.Vars {
			if v.AssociatedSSAObj == i {
				return &v.Var
			}
		}
	}

	for _, v := range g.params {
		if v.AssociatedSSAObj == i {
			return &v.Var
		}
	}

	switch v := i.(type) {
	case *ssa.Const:
		switch t := v.Type().(type) {
		case *types.Basic:
			switch t.Kind() {
			case types.Bool:
				return cconstant.NewBool(constant.BoolVal(v.Value))

			case types.Int:
				val, _ := constant.Int64Val(v.Value)
				return cconstant.NewInt(ctypes.Int64, val)

			case types.Float32:
				val, _ := constant.Float32Val(v.Value)
				return cconstant.NewFloat(val)

			case types.Float64:
				val, _ := constant.Float64Val(v.Value)
				return cconstant.NewDouble(val)

			case types.String:
				val := constant.StringVal(v.Value)
				return cconstant.NewString(val)

			default:
				logger.Fatal("Todo: getValue(), ", t)
			}

		case *types.Slice:
			logger.Fatal("Todo: getValue(), ", t)

		default:
			logger.Fatal("Todo: getValue(), ", t)

		}

		/*
			switch v.Value.Kind() {
			case constant.Kind(constant.Bool):
				return cconstant.NewBool(constant.BoolVal(v.Value))

			case constant.Kind(constant.Int):
				val, _ := constant.Int64Val(v.Value)
				return cconstant.NewInt(ctypes.Int64, val)

			case constant.Kind(constant.Float):
				val, _ := constant.Float64Val(v.Value)
				return cconstant.NewDouble(val)

			case constant.String:
				val := constant.StringVal(v.Value)
				return cconstant.NewString(val)

			default:
				logger.Fatal("Todo")
			}  //*/
	}

	logger.Fatal("Value not found:", i)
	return nil
}

func (g *functionGenerator) lookupStruct(n string) ctypes.Type {
	for scope := g.compiler.curScope; scope != nil; scope = scope.Parent {
		for _, s := range scope.Structs {
			if s.Struct.CIRString() == n {
				return &s.Struct
			}
		}
	}
	return nil
}

func (g *functionGenerator) genFunction(f *ssa.Function) {
	rets := f.Signature.Results()
	ret_type := cir.ToCType(rets)

	if rets.Len() > 1 {
		if g.lookupStruct(ret_type.CIRString()) == nil {
			switch s := ret_type.(type) {
			case *ctypes.Tuple:
				g.compiler.curScope.AddTupleDecl(*s)

			default:
				logger.Fatal("Type of multi-ret should be Tuple")
			}
		}
	}

	for _, i := range f.Params {
		pa := *cir.NewVarDecl(i.Name(), cir.ToCType(i.Type()))
		pa.AssociatedSSAObj = i
		g.params = append(g.params, pa)
	}
	fn := g.compiler.curScope.AddFuncDecl(f.Name(), ret_type, g.params)

	if len(f.Blocks) == 0 {
		return
	}
	fn.Body = cir.NewScope(g.compiler.curScope)
	g.compiler.curScope = fn.Body
	defer func() { g.compiler.curScope = fn.Body.Parent }()

	g.var_cur_block = &g.compiler.curScope.AddVarDecl("$T_Block_Current", ctypes.Uint32).Var

	for _, b := range f.Blocks {
		g.genBlock(b)
	}
}

func (g *functionGenerator) genBlock(block *ssa.BasicBlock) {
	if len(block.Instrs) == 0 {
		logger.Fatalf("Block:%s is empty", block)
	}

	g.compiler.curScope.AddBlackLine()
	g.compiler.curScope.AddLabel("$Block_" + strconv.Itoa(block.Index))

	cur_block_assigned := false
	for _, inst := range block.Instrs {
		if _, ok := inst.(*ssa.Phi); !ok {
			if !cur_block_assigned {
				g.compiler.curScope.AddAssignStmt(g.var_cur_block, cconstant.NewInt(ctypes.Uint32, int64(block.Index)))
				cur_block_assigned = true
			}
		}
		g.createInstruction(inst)
	}
}

func (g *functionGenerator) createInstruction(inst ssa.Instruction) {
	switch inst := inst.(type) {

	case *ssa.Phi:
		g.genPhi(inst)

	case *ssa.Alloc:
		g.genAlloc(inst)

	case *ssa.If:
		g.genIf(inst)

	case *ssa.Store:
		g.genStore(inst)

	case *ssa.Jump:
		g.genJump(inst)

	case *ssa.Return:
		g.genReturn(inst)

	case ssa.Value:
		v := g.genValue(inst)
		if v.Type().Equal(ctypes.Void) {
			g.compiler.curScope.AddExprStmt(v)
		} else {
			r := g.compiler.curScope.AddVarDecl(g.genRegister(), v.Type())
			r.AssociatedSSAObj = inst
			r.InitVal = v
		}

	default:
		logger.Fatal("Todo:", inst.String())
	}

}

func (g *functionGenerator) genValue(v ssa.Value) cir.Expr {
	if _, ok := g.locals[v]; ok {
		logger.Fatal("Instruction already exist：", v.String())
	}

	switch v := v.(type) {
	case *ssa.UnOp:
		return g.genUnOp(v)

	case *ssa.BinOp:
		return g.genBinOp(v)

	case *ssa.Call:
		return g.genCall(v)

	case *ssa.Extract:
		return g.genExtract(v)

	case *ssa.FieldAddr:
		return g.genFieldAddr(v)

	case *ssa.Field:
		return g.genField(v)

	case *ssa.IndexAddr:
		return g.genIndexAddr(v)

	case *ssa.Slice:

	}

	logger.Fatalf("Todo: %v, type: %T", v, v)
	return nil
}

func (g *functionGenerator) genUnOp(inst *ssa.UnOp) cir.Expr {
	switch inst.Op {
	case token.NOT:
		return cir.NewNotExpr(g.getValue(inst.X))

	case token.MUL: //*x
		return g.genLoad(inst.X)

	default:
		logger.Fatal("Todo")
	}

	logger.Fatal("Todo")
	return nil
}

func (g *functionGenerator) genLoad(addr ssa.Value) cir.Expr {
	return cir.NewLoadExpr(g.getValue(addr))
}

func (g *functionGenerator) genBinOp(inst *ssa.BinOp) cir.Expr {
	x := g.getValue(inst.X)
	y := g.getValue(inst.Y)

	switch inst.X.Type().Underlying().(type) {
	case *types.Basic:
		switch inst.Op {
		case token.ADD:
			return cir.NewAddExpr(x, y)

		case token.SUB:
			return cir.NewSubExpr(x, y)

		case token.MUL:
			return cir.NewMulExpr(x, y)

		case token.QUO:
			return cir.NewQuoExpr(x, y)

		case token.EQL:
			return cir.NewEqlExpr(x, y)
		}

	default:
		logger.Fatalf("Todo: %v, type: %T, token:%v", inst, inst, inst.Op)
	}

	logger.Fatalf("Todo: %v, type: %T, token:%v", inst, inst, inst.Op)
	return nil
}

func (g *functionGenerator) genCall(inst *ssa.Call) cir.Expr {
	if inst.Call.IsInvoke() {
		logger.Fatal("Todo: Invoke")
	}

	switch inst.Call.Value.(type) {
	case *ssa.Function:
		ret_type := cir.ToCType(inst.Call.Signature().Results())
		var params_type []ctypes.Type
		var params []cir.Expr
		for _, v := range inst.Call.Args {
			params_type = append(params_type, cir.ToCType(v.Type()))
			params = append(params, g.getValue(v))
		}
		fn_type := ctypes.NewFuncType(ret_type, params_type)
		fn := cir.NewVar(inst.Call.StaticCallee().Name(), fn_type)

		return cir.NewCallExpr(fn, params)

	case *ssa.Builtin:
		return g.genBuiltin(inst.Common())

	case *ssa.MakeClosure:
		logger.Fatal("Todo: MakeClosure")

	default:
		logger.Fatalf("Todo: type:%T", inst.Call.Value)
	}

	logger.Fatal("Todo")
	return nil
}

func (g *functionGenerator) genExtract(inst *ssa.Extract) cir.Expr {
	return cir.NewSelectExpr(g.getValue(inst.Tuple), cir.NewRawExpr("$m"+strconv.Itoa(inst.Index), cir.ToCType(inst.Type())))
}

func (g *functionGenerator) genAlloc(inst *ssa.Alloc) {
	if inst.Heap {
		ref_type := ctypes.NewRefType(cir.ToCType(inst.Type().(*types.Pointer).Elem()))
		v := cir.NewRawExpr(ref_type.CIRString()+"::New()", ref_type)
		r := g.compiler.curScope.AddVarDecl(g.genRegister(), v.Type())
		r.AssociatedSSAObj = inst
		r.InitVal = v
		return
	}

	c_type := cir.ToCType(inst.Type().(*types.Pointer).Elem())
	r := g.compiler.curScope.AddVarDecl(g.genRegister(), c_type)
	switch c_type := c_type.(type) {
	case *ctypes.Array:
		r.AssociatedSSAObj = inst

	case *ctypes.Slice:
		r.AssociatedSSAObj = inst

	default:
		rt := g.compiler.curScope.AddVarDecl(g.genRegister(), ctypes.NewPointerType(c_type))
		rt.InitVal = cir.NewGetaddrExpr(&r.Var)
		rt.AssociatedSSAObj = inst
	}
}

func (g *functionGenerator) genFieldAddr(inst *ssa.FieldAddr) cir.Expr {
	cx := g.getValue(inst.X)
	field := inst.X.Type().Underlying().(*types.Pointer).Elem().Underlying().(*types.Struct).Field(inst.Field)
	fieldname := field.Name()
	if field.Embedded() {
		fieldname = "$" + fieldname
	}
	switch cx.Type().(type) {
	case *ctypes.RefType:
		ct := ctypes.NewRefType(cir.ToCType(field.Type()))
		return cir.NewRawExpr("{ &"+cx.CIRString()+".GetRaw()->"+fieldname+", "+cx.CIRString()+".GetBlock() }", ct)

	case *ctypes.PointerType:
		ct := ctypes.NewPointerType(cir.ToCType(field.Type()))
		return cir.NewRawExpr("&"+cx.CIRString()+"->"+fieldname, ct)

	default:
		logger.Fatalf("Invalid type:%T", cx.Type())
	}

	return nil
}

func (g *functionGenerator) genField(inst *ssa.Field) cir.Expr {
	cx := g.getValue(inst.X)
	field := inst.X.Type().(*types.Struct).Field(inst.Field)
	fieldname := field.Name()
	if field.Embedded() {
		fieldname = "$" + fieldname
	}
	return cir.NewRawExpr(cx.CIRString()+"."+fieldname, cir.ToCType(field.Type()))
}

func (g *functionGenerator) genIndexAddr(inst *ssa.IndexAddr) cir.Expr {
	println(g.getValue(inst.Index))
	x := g.getValue(inst.X)
	switch t := x.Type().(type) {
	case *ctypes.Array:
		return cir.NewRawExpr("&"+x.CIRString()+".at("+g.getValue(inst.Index).CIRString()+")", ctypes.NewPointerType(t.GetElem()))

	case *ctypes.Slice:
		return cir.NewRawExpr("{ "+x.CIRString()+".At("+g.getValue(inst.Index).CIRString()+"), "+x.CIRString()+".GetBlock() }", ctypes.NewRefType(t.GetElem()))

	case *ctypes.PointerType:
		switch t := t.Base.(type) {
		case *ctypes.Array:
			return cir.NewRawExpr("&"+x.CIRString()+"->at("+g.getValue(inst.Index).CIRString()+")", ctypes.NewPointerType(t.GetElem()))

		case *ctypes.Slice:
			return cir.NewRawExpr("{ "+x.CIRString()+"->At("+g.getValue(inst.Index).CIRString()+"), "+x.CIRString()+".GetBlock() }", ctypes.NewRefType(t.GetElem()))

		default:
			logger.Fatalf("Todo: genIndexAddr(), %T %s", x.Type(), inst)
		}

	case *ctypes.RefType:
		switch t := t.Base.(type) {
		case *ctypes.Array:
			return cir.NewRawExpr("{ &"+x.CIRString()+".GetRaw()->at("+g.getValue(inst.Index).CIRString()+"), "+x.CIRString()+".GetBlock() }", ctypes.NewRefType(t.GetElem()))

		case *ctypes.Slice:
			return cir.NewRawExpr("{ "+x.CIRString()+".GetRaw()->At("+g.getValue(inst.Index).CIRString()+"), "+x.CIRString()+".GetBlock() }", ctypes.NewRefType(t.GetElem()))
		}

	default:
		logger.Fatalf("Todo: genIndexAddr(), %T %s", x.Type(), inst)

	}
	return nil
}

func (g *functionGenerator) genBuiltin(call *ssa.CallCommon) cir.Expr {
	switch call.Value.Name() {
	case "print", "println":
		var args []cir.Expr
		args = append(args, nil)
		var arg_type []ctypes.Type
		var f string
		for _, arg := range call.Args {
			if len(args) > 1 {
				f += " "
			}
			arg := g.getValue(arg)
			switch arg.Type().(type) {
			case *ctypes.StringType:
				arg_type = append(arg_type, &ctypes.StringType{})
				f += "%s"
				args = append(args, cir.NewRawExpr(arg.CIRString()+".c_str()", &ctypes.StringType{}))

			case *ctypes.BoolType:
			case *ctypes.IntType:
				arg_type = append(arg_type, ctypes.Int64)
				f += "%d"
				args = append(args, arg)

			case *ctypes.FloatType:
				arg_type = append(arg_type, ctypes.Double)
				f += "%lf"
				args = append(args, arg)

			default:
				logger.Fatalf("Todo: print(%s)", arg.Type().CIRString())
			}
		}
		if call.Value.Name() == "println" {
			f += "\\n"
		}

		args[0] = cconstant.NewString(f)
		fn_type := ctypes.NewFuncType(ctypes.Void, arg_type)
		fn := cir.NewVar("printf", fn_type)
		expr := cir.NewCallExpr(fn, args)
		return expr

	}
	logger.Fatal("Todo:", call.Value)
	return nil
}

func (g *functionGenerator) genPrint(v cir.Expr) cir.Expr {
	var args []cir.Expr
	var fn_type ctypes.Type
	switch v.Type().(type) {
	case *ctypes.StringType:
		args = append(args, cconstant.NewString("%s"))
		args = append(args, v)
		fn_type = ctypes.NewFuncType(ctypes.Void, []ctypes.Type{&ctypes.StringType{}, &ctypes.StringType{}})

	case *ctypes.BoolType:
	case *ctypes.IntType:
		args = append(args, cconstant.NewString("%d"))
		args = append(args, v)
		fn_type = ctypes.NewFuncType(ctypes.Void, []ctypes.Type{&ctypes.StringType{}, ctypes.Int64})

	case *ctypes.FloatType:
		args = append(args, cconstant.NewString("%lf"))
		args = append(args, v)
		fn_type = ctypes.NewFuncType(ctypes.Void, []ctypes.Type{&ctypes.StringType{}, ctypes.Double})

	default:
		logger.Fatalf("Todo: print(%s)", v.Type().Name())
		return nil
	}
	fn := cir.NewVar("printf", fn_type)
	expr := cir.NewCallExpr(fn, args)
	return expr
}

func (g *functionGenerator) genPhi(inst *ssa.Phi) {
	r := &g.compiler.curScope.AddVarDecl(g.genRegister(), cir.ToCType(inst.Type())).Var
	r.AssociatedSSAObj = inst

	var edges []cir.PhiEdge
	for i, v := range inst.Edges {
		var e cir.PhiEdge
		e.Incoming = inst.Block().Preds[i].Index
		e.Value = g.getValue(v)
		edges = append(edges, e)
	}

	g.compiler.curScope.AddPhiStmt(g.var_cur_block, r, edges)
}

func (g *functionGenerator) genRegister() string {
	defer func() { g.localID++ }()
	return "$T_" + strconv.Itoa(g.localID)
}

func (g *functionGenerator) genReturn(inst *ssa.Return) {
	var ret []cir.Expr
	for _, v := range inst.Results {
		ret = append(ret, g.getValue(v))
	}

	g.compiler.curScope.AddReturnStmt(ret)
}

func (g *functionGenerator) genIf(inst *ssa.If) {
	var succs [2]int
	succs[0] = inst.Block().Succs[0].Index
	succs[1] = inst.Block().Succs[1].Index

	g.compiler.curScope.AddIfStmt(g.getValue(inst.Cond), succs)
}

func (g *functionGenerator) genStore(inst *ssa.Store) {
	g.compiler.curScope.AddStoreStmt(g.getValue(inst.Addr), g.getValue(inst.Val))
}

func (g *functionGenerator) genJump(inst *ssa.Jump) {
	g.compiler.curScope.AddJumpStmt("$Block_" + strconv.Itoa(inst.Block().Succs[0].Index))
}
