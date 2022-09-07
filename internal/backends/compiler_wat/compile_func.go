// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_wat

import (
	"strconv"

	"github.com/wa-lang/wa/internal/constant"
	"github.com/wa-lang/wa/internal/token"
	"github.com/wa-lang/wa/internal/types"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir"
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
	"github.com/wa-lang/wa/internal/ssa"
)

type functionGenerator struct {
	module *wat.Module

	locals_map map[ssa.Value]wir.Value

	registers    []wir.Value
	cur_local_id int

	var_block_selector wir.Value
	var_current_block  wir.Value
	var_ret            wir.Value
}

func newFunctionGenerator(p *Compiler) *functionGenerator {
	return &functionGenerator{module: &p.module, locals_map: make(map[ssa.Value]wir.Value)}
}

func (g *functionGenerator) getValue(i ssa.Value) wir.Value {
	if i == nil {
		return nil
	}

	if v, ok := g.locals_map[i]; ok {
		return v
	}

	switch v := i.(type) {
	case *ssa.Const:
		switch t := v.Type().(type) {
		case *types.Basic:
			switch t.Kind() {
			case types.Bool:
				logger.Fatalf("Todo:%T", t)

			case types.Int:
				val, _ := constant.Int64Val(v.Value)
				return wir.NewConstI32(int32(val))

			case types.Int32:
				val, _ := constant.Int64Val(v.Value)
				return wir.NewConstI32(int32(val))

			case types.Float32:
				logger.Fatalf("Todo:%T", t)

			case types.Float64:
				logger.Fatalf("Todo:%T", t)

			case types.String, types.UntypedString:
				logger.Fatalf("Todo:%T", t)

			default:
				logger.Fatalf("Todo:%T", t)
			}

		case *types.Slice:
			logger.Fatalf("Todo:%T", t)

		default:
			logger.Fatalf("Todo:%T", t)
		}

	case ssa.Instruction:
		nv := g.addRegister(wir.ToWType(i.Type()))
		g.locals_map[i] = nv
		return nv
	}

	logger.Fatal("Value not found:", i)
	return nil
}

func (g *functionGenerator) genFunction(f *ssa.Function) *wat.Function {
	var wir_fn wir.Function
	wir_fn.Name = f.Name()

	rets := f.Signature.Results()
	wir_fn.Result = wir.ToWType(rets)
	if rets.Len() > 1 {
		logger.Fatal("Todo")
	}

	for _, i := range f.Params {
		pa := wir.NewVar(i.Name(), wir.ValueKindLocal, wir.ToWType(i.Type()))
		wir_fn.Params = append(wir_fn.Params, pa)
		g.locals_map[i] = pa
	}

	g.var_block_selector = wir.NewVar("$block_selector", wir.ValueKindLocal, wir.Int32{})
	g.registers = append(g.registers, g.var_block_selector)
	wir_fn.Insts = append(wir_fn.Insts, g.var_block_selector.EmitInit()...)

	g.var_current_block = wir.NewVar("$current_block", wir.ValueKindLocal, wir.Int32{})
	g.registers = append(g.registers, g.var_current_block)
	wir_fn.Insts = append(wir_fn.Insts, g.var_current_block.EmitInit()...)

	if !wir_fn.Result.Equal(wir.Void{}) {
		g.var_ret = wir.NewVar("$ret", wir.ValueKindLocal, wir_fn.Result)
		g.registers = append(g.registers, g.var_ret)
		wir_fn.Insts = append(wir_fn.Insts, g.var_ret.EmitInit()...)
	}

	var block_temp wat.Inst
	//BlockSel:
	{
		inst := wat.NewInstBlock("$BlockSel")
		inst.Insts = append(inst.Insts, g.var_block_selector.EmitGet()...)
		t := make([]int, len(f.Blocks)+1)
		for i := range f.Blocks {
			t[i] = i
		}
		t[len(f.Blocks)] = 0
		inst.Insts = append(inst.Insts, wat.NewInstBrTable(t))
		block_temp = inst
	}

	for i, b := range f.Blocks {
		block := wat.NewInstBlock("$Block_" + strconv.Itoa(i))
		block.Insts = append(block.Insts, block_temp)
		block.Insts = append(block.Insts, g.genBlock(b)...)
		block_temp = block
	}

	//BlockDisp
	{
		inst := wat.NewInstLoop("$BlockDisp")
		inst.Insts = append(inst.Insts, block_temp)
		block_temp = inst
	}

	//BlockFnBody
	{
		inst := wat.NewInstBlock("$BlockFnBody")
		inst.Insts = append(inst.Insts, block_temp)
		block_temp = inst
	}

	wir_fn.Insts = append(wir_fn.Insts, block_temp)

	if !wir_fn.Result.Equal(wir.Void{}) {
		wir_fn.Insts = append(wir_fn.Insts, g.var_ret.EmitGet()...)
	}

	wir_fn.Locals = g.registers

	return wir_fn.ToWatFunc()
}

func (g *functionGenerator) genBlock(block *ssa.BasicBlock) []wat.Inst {
	if len(block.Instrs) == 0 {
		logger.Fatalf("Block:%s is empty", block)
	}

	cur_block_assigned := false
	var b []wat.Inst
	for _, inst := range block.Instrs {
		if _, ok := inst.(*ssa.Phi); !ok {
			if !cur_block_assigned {
				b = append(b, wir.EmitAssginValue(g.var_current_block, wir.NewConstI32(int32(block.Index)))...)
				cur_block_assigned = true
			}
		}

		b = append(b, g.genInstruction(inst)...)
	}
	return b

}

func (g *functionGenerator) genInstruction(inst ssa.Instruction) []wat.Inst {
	switch inst := inst.(type) {

	case *ssa.Alloc:
		logger.Fatalf("Todo:%T", inst)

	case *ssa.If:
		return g.genIf(inst)

	case *ssa.Store:
		logger.Fatalf("Todo:%T", inst)

	case *ssa.Jump:
		return g.genJump(inst)

	case *ssa.Return:
		return g.genReturn(inst)

	case *ssa.Extract:
		logger.Fatalf("Todo:%T", inst)

	case *ssa.Field:
		logger.Fatalf("Todo:%T", inst)

	case ssa.Value:
		s, t := g.genValue(inst)
		if !t.Equal(wir.Void{}) {
			v := g.getValue(inst)
			s = append(s, v.EmitSet()...)
			g.locals_map[inst] = v
		}
		return s

	default:
		logger.Fatal("Todo:", inst.String())
	}
	return nil
}

func (g *functionGenerator) genValue(v ssa.Value) ([]wat.Inst, wir.ValueType) {
	//if _, ok := g.locals_map[v]; ok {
	//	logger.Printf("Instruction already exist：%s\n", v)
	//}

	switch v := v.(type) {
	case *ssa.UnOp:
		logger.Fatalf("Todo: %v, type: %T", v, v)

	case *ssa.BinOp:
		return g.genBinOp(v)

	case *ssa.Call:
		return g.genCall(v)

	case *ssa.Phi:
		return g.genPhi(v)

	case *ssa.FieldAddr:
		logger.Fatalf("Todo: %v, type: %T", v, v)

	case *ssa.IndexAddr:
		logger.Fatalf("Todo: %v, type: %T", v, v)

	case *ssa.Slice:
		logger.Fatalf("Todo: %v, type: %T", v, v)
	}

	logger.Fatalf("Todo: %v, type: %T", v, v)
	return nil, nil
}

func (g *functionGenerator) genBinOp(inst *ssa.BinOp) ([]wat.Inst, wir.ValueType) {
	x := g.getValue(inst.X)
	y := g.getValue(inst.Y)

	switch inst.X.Type().Underlying().(type) {
	case *types.Basic:
		switch inst.Op {
		case token.ADD:
			return wir.EmitBinOp(x, y, wat.OpCodeAdd)

		case token.SUB:
			return wir.EmitBinOp(x, y, wat.OpCodeSub)

		case token.MUL:
			return wir.EmitBinOp(x, y, wat.OpCodeMul)

		case token.QUO:
			return wir.EmitBinOp(x, y, wat.OpCodeQuo)

		case token.REM:
			return wir.EmitBinOp(x, y, wat.OpCodeRem)

		case token.EQL:
			return wir.EmitBinOp(x, y, wat.OpCodeEql)

		case token.NEQ:
			return wir.EmitBinOp(x, y, wat.OpCodeNe)

		case token.LSS:
			return wir.EmitBinOp(x, y, wat.OpCodeLt)

		case token.GTR:
			return wir.EmitBinOp(x, y, wat.OpCodeGt)

		case token.LEQ:
			return wir.EmitBinOp(x, y, wat.OpCodeLe)

		case token.GEQ:
			return wir.EmitBinOp(x, y, wat.OpCodeGe)
		}

	default:
		logger.Fatalf("Todo: %v, type: %T, token:%v", inst, inst, inst.Op)
	}

	logger.Fatalf("Todo: %v, type: %T, token:%v", inst, inst, inst.Op)
	return nil, nil
}

func (g *functionGenerator) genCall(inst *ssa.Call) ([]wat.Inst, wir.ValueType) {
	if inst.Call.IsInvoke() {
		logger.Fatal("Todo: genCall(), Invoke")
	}

	switch inst.Call.Value.(type) {
	case *ssa.Function:
		ret_type := wir.ToWType(inst.Call.Signature().Results())
		var insts []wat.Inst
		for _, v := range inst.Call.Args {
			insts = append(insts, g.getValue(v).EmitGet()...)
		}
		insts = append(insts, wat.NewInstCall(inst.Call.StaticCallee().Name()))
		return insts, ret_type

	case *ssa.Builtin:
		return g.genBuiltin(inst.Common())

	case *ssa.MakeClosure:
		logger.Fatal("Todo: genCall(), MakeClosure")

	default:
		logger.Fatalf("Todo: type:%T", inst.Call.Value)
	}

	logger.Fatal("Todo")

	return nil, nil
}

func (g *functionGenerator) genBuiltin(call *ssa.CallCommon) ([]wat.Inst, wir.ValueType) {
	switch call.Value.Name() {
	case "print", "println":
		var insts []wat.Inst
		for _, arg := range call.Args {
			arg := g.getValue(arg)
			switch arg.Type().(type) {
			case wir.Int32:
				insts = append(insts, arg.EmitGet()...)
				insts = append(insts, wat.NewInstCall("$print_i32"))

			default:
				logger.Fatalf("Todo: print(%T)", arg.Type())
			}
		}

		if call.Value.Name() == "println" {
			insts = append(insts, wir.NewConstI32('\n').EmitGet()...)
			insts = append(insts, wat.NewInstCall("$print_rune"))
		}

		return insts, wir.Void{}
	}
	logger.Fatal("Todo:", call.Value)
	return nil, nil
}

func (g *functionGenerator) genPhiIter(preds []int, values []wir.Value) []wat.Inst {
	var insts []wat.Inst

	cond, _ := wir.EmitBinOp(g.var_current_block, wir.NewConstI32(int32(preds[0])), wat.OpCodeEql)
	insts = append(insts, cond...)

	trueInsts := values[0].EmitGet()
	var falseInsts []wat.Inst
	if len(preds) == 2 {
		falseInsts = values[1].EmitGet()
	} else {
		falseInsts = g.genPhiIter(preds[1:], values[1:])
	}
	insts = append(insts, wat.NewInstIf(trueInsts, falseInsts, values[0].Type().Raw()))

	return insts
}
func (g *functionGenerator) genPhi(inst *ssa.Phi) ([]wat.Inst, wir.ValueType) {
	var preds []int
	var values []wir.Value
	for i, v := range inst.Edges {
		preds = append(preds, inst.Block().Preds[i].Index)
		values = append(values, g.getValue(v))
	}
	return g.genPhiIter(preds, values), wir.ToWType(inst.Type())
}

func (g *functionGenerator) genReturn(inst *ssa.Return) []wat.Inst {
	var insts []wat.Inst

	switch len(inst.Results) {
	case 0:
		break

	case 1:
		insts = append(insts, wir.EmitAssginValue(g.var_ret, g.getValue(inst.Results[0]))...)

	default:
		logger.Fatal("Todo")
	}

	insts = append(insts, wat.NewInstBr("$BlockFnBody"))
	return insts
}

func (g *functionGenerator) genIf(inst *ssa.If) []wat.Inst {
	cond := g.getValue(inst.Cond)
	if !cond.Type().Equal(wir.Int32{}) {
		logger.Fatal("cond.type() != i32")
	}

	insts := cond.EmitGet()
	instsTrue := g.genJumpID(inst.Block().Index, inst.Block().Succs[0].Index)
	instsFalse := g.genJumpID(inst.Block().Index, inst.Block().Succs[1].Index)
	insts = append(insts, wat.NewInstIf(instsTrue, instsFalse, nil))

	return insts
}

func (g *functionGenerator) genJump(inst *ssa.Jump) []wat.Inst {
	return g.genJumpID(inst.Block().Index, inst.Block().Succs[0].Index)
}

func (g *functionGenerator) genJumpID(cur, dest int) []wat.Inst {
	var insts []wat.Inst

	if cur >= dest {
		insts = wir.EmitAssginValue(g.var_block_selector, wir.NewConstI32(int32(dest)))
		insts = append(insts, wat.NewInstBr("$BlockDisp"))
	} else {
		insts = append(insts, wat.NewInstBr("$Block_"+strconv.Itoa(dest-1)))
	}

	return insts
}

func (g *functionGenerator) addRegister(typ wir.ValueType) wir.Value {
	defer func() { g.cur_local_id++ }()
	name := "$T_" + strconv.Itoa(g.cur_local_id)
	v := wir.NewVar(name, wir.ValueKindLocal, typ)
	g.registers = append(g.registers, v)
	return v
}
