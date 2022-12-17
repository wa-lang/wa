// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_wat

import (
	"strconv"
	"strings"

	"wa-lang.org/wa/internal/constant"
	"wa-lang.org/wa/internal/token"
	"wa-lang.org/wa/internal/types"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir"
	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/ssa"
)

type valueWrap struct {
	value          wir.Value
	force_register bool
}

type functionGenerator struct {
	module *wir.Module

	locals_map map[ssa.Value]valueWrap

	registers    []wir.Value
	cur_local_id int

	var_block_selector wir.Value
	var_current_block  wir.Value
	var_rets           []wir.Value
}

func newFunctionGenerator(module *wir.Module) *functionGenerator {
	return &functionGenerator{module: module, locals_map: make(map[ssa.Value]valueWrap)}
}

func (g *functionGenerator) getValue(i ssa.Value) valueWrap {
	if i == nil {
		return valueWrap{}
	}

	if v, ok := g.locals_map[i]; ok {
		return v
	}

	if v, ok := g.module.Globals_map[i]; ok {
		return valueWrap{value: v}
	}

	switch v := i.(type) {
	case *ssa.Const:
		//if v.Value == nil {
		//	return nil
		//}

		switch t := v.Type().(type) {
		case *types.Basic:
			switch t.Kind() {

			case types.Bool, types.UntypedBool:
				if constant.BoolVal(v.Value) {
					return valueWrap{value: wir.NewConst("1", wir.I32{})}
				} else {
					return valueWrap{value: wir.NewConst("0", wir.I32{})}
				}

			case types.Int:
				val, _ := constant.Int64Val(v.Value)
				return valueWrap{value: wir.NewConst(strconv.Itoa(int(val)), wir.I32{})}

			case types.Uint8:
				val, _ := constant.Uint64Val(v.Value)
				return valueWrap{value: wir.NewConst(strconv.Itoa(int(val)), wir.U8{})}

			case types.Int8:
				val, _ := constant.Uint64Val(v.Value)
				return valueWrap{value: wir.NewConst(strconv.Itoa(int(val)), wir.I8{})}

			case types.Uint16:
				val, _ := constant.Uint64Val(v.Value)
				return valueWrap{value: wir.NewConst(strconv.Itoa(int(val)), wir.U16{})}

			case types.Int16:
				val, _ := constant.Uint64Val(v.Value)
				return valueWrap{value: wir.NewConst(strconv.Itoa(int(val)), wir.I16{})}

			case types.Uint32:
				val, _ := constant.Uint64Val(v.Value)
				return valueWrap{value: wir.NewConst(strconv.Itoa(int(val)), wir.U32{})}

			case types.Int32:
				val, _ := constant.Int64Val(v.Value)
				if t.Name() == "rune" {
					return valueWrap{value: wir.NewConst(strconv.Itoa(int(val)), wir.RUNE{})}
				} else {
					return valueWrap{value: wir.NewConst(strconv.Itoa(int(val)), wir.I32{})}
				}

			case types.Float32:
				val, _ := constant.Float64Val(v.Value)
				return valueWrap{value: wir.NewConst(strconv.FormatFloat(val, 'f', -1, 32), wir.F32{})}

			case types.Float64:
				val, _ := constant.Float64Val(v.Value)
				return valueWrap{value: wir.NewConst(strconv.FormatFloat(val, 'f', -1, 64), wir.F64{})}

			case types.String, types.UntypedString:
				val := constant.StringVal(v.Value)
				return valueWrap{value: wir.NewConst(val, wir.NewString())}

			default:
				logger.Fatalf("Todo:%T %v", t, t.Kind())
			}

		case *types.Pointer:
			if v.Value == nil {
				return valueWrap{}
			} else {
				logger.Fatalf("Todo:%T", t)
			}

		case *types.Slice:
			if v.Value == nil {
				return valueWrap{value: wir.NewConst("0", wir.NewSlice(wir.ToWType(t.Elem())))}
			}
			logger.Fatalf("Todo:%T", t)

		default:
			logger.Fatalf("Todo:%T", t)
		}

	case ssa.Instruction:
		nv := valueWrap{value: g.addRegister(wir.ToWType(i.Type()))}
		g.locals_map[i] = nv
		return nv

	case *ssa.Function:
		fn_name, _ := GetFnMangleName(v)
		fn_sig := wir.NewFnSigFromSignature(v.Signature)

		return valueWrap{value: wir.GenConstFnValue(fn_name, fn_sig)}
	}

	logger.Fatal("Value not found:", i)
	return valueWrap{}
}

func (g *functionGenerator) genFunction(f *ssa.Function) *wir.Function {
	var wir_fn wir.Function
	{
		internal, external := GetFnMangleName(f)
		if len(f.LinkName()) > 0 {
			wir_fn.InternalName = f.LinkName()
		} else {
			wir_fn.InternalName = internal
		}
		if len(f.ExportName()) > 0 {
			wir_fn.ExternalName = f.ExportName()
		} else {
			wir_fn.ExternalName = external
		}
	}

	rets := f.Signature.Results()
	switch rets.Len() {
	case 0:
		break

	case 1:
		wir_fn.Results = append(wir_fn.Results, wir.ToWType(rets.At(0).Type()))

	default:
		typ := wir.ToWType(rets).(wir.Tuple)
		for _, f := range typ.Members {
			wir_fn.Results = append(wir_fn.Results, f.Type())
		}
	}

	for _, i := range f.FreeVars {
		fv := valueWrap{value: wir.NewLocal(wir.GenSymbolName(i.Name()), wir.ToWType(i.Type()))}
		wir_fn.Params = append(wir_fn.Params, fv.value)
		g.locals_map[i] = fv
	}
	for _, i := range f.Params {
		pv := valueWrap{value: wir.NewLocal(wir.GenSymbolName(i.Name()), wir.ToWType(i.Type()))}
		wir_fn.Params = append(wir_fn.Params, pv.value)
		g.locals_map[i] = pv
	}

	g.var_block_selector = wir.NewLocal("$block_selector", wir.I32{})
	g.registers = append(g.registers, g.var_block_selector)
	g.var_current_block = wir.NewLocal("$current_block", wir.I32{})
	g.registers = append(g.registers, g.var_current_block)
	for i, rt := range wir_fn.Results {
		rname := "$ret_" + strconv.Itoa(i)
		r := wir.NewLocal(rname, rt)
		g.var_rets = append(g.var_rets, r)
		g.registers = append(g.registers, r)
	}

	var block_temp wat.Inst
	//BlockSel:
	{
		inst := wat.NewInstBlock("$BlockSel")
		inst.Insts = append(inst.Insts, g.var_block_selector.EmitPush()...)
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

	for _, i := range g.registers {
		wir_fn.Insts = append(wir_fn.Insts, i.EmitInit()...)
	}

	wir_fn.Insts = append(wir_fn.Insts, block_temp)

	for _, r := range g.var_rets {
		wir_fn.Insts = append(wir_fn.Insts, r.EmitPush()...)
	}

	for _, i := range g.registers {
		wir_fn.Insts = append(wir_fn.Insts, i.EmitRelease()...)
	}

	for _, i := range wir_fn.Params {
		wir_fn.Insts = append(wir_fn.Insts, i.EmitRelease()...)
	}

	wir_fn.Locals = g.registers

	return &wir_fn
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
				b = append(b, wir.EmitAssginValue(g.var_current_block, wir.NewConst(strconv.Itoa(block.Index), wir.I32{}))...)
				b = append(b, wat.NewBlank())
				cur_block_assigned = true
			}
		}

		b = append(b, g.genInstruction(inst)...)
	}
	return b

}

func (g *functionGenerator) genInstruction(inst ssa.Instruction) (insts []wat.Inst) {
	insts = append(insts, wat.NewComment(inst.String()))

	switch inst := inst.(type) {

	case *ssa.If:
		insts = append(insts, g.genIf(inst)...)

	case *ssa.Store:
		insts = append(insts, g.genStore(inst)...)

	case *ssa.Jump:
		insts = append(insts, g.genJump(inst)...)

	case *ssa.Return:
		insts = append(insts, g.genReturn(inst)...)

	case ssa.Value:
		s, t := g.genValue(inst)
		if t != nil && !t.Equal(wir.VOID{}) {
			if v, ok := g.locals_map[inst]; ok {
				if !v.value.Type().Equal(t) {
					panic("Type not match")
				}
				s = append(s, v.value.EmitPop()...)
			} else {
				nv := g.addRegister(t)
				g.locals_map[inst] = valueWrap{value: nv}
				s = append(s, nv.EmitPop()...)
			}
		}
		insts = append(insts, s...)

	default:
		logger.Fatal("Todo:", inst.String())
	}
	insts = append(insts, wat.NewBlank())
	return
}

func (g *functionGenerator) genValue(v ssa.Value) ([]wat.Inst, wir.ValueType) {
	//if _, ok := g.locals_map[v]; ok {
	//	logger.Printf("Instruction already exist：%s\n", v)
	//}

	switch v := v.(type) {
	case *ssa.UnOp:
		return g.genUnOp(v)

	case *ssa.BinOp:
		return g.genBinOp(v)

	case *ssa.Call:
		return g.genCall(v)

	case *ssa.Phi:
		return g.genPhi(v)

	case *ssa.Alloc:
		return g.genAlloc(v)

	case *ssa.Extract:
		return g.genExtract(v)

	case *ssa.Field:
		return g.genFiled(v)

	case *ssa.FieldAddr:
		return g.genFieldAddr(v)

	case *ssa.IndexAddr:
		return g.genIndexAddr(v)

	case *ssa.Slice:
		return g.genSlice(v)

	case *ssa.MakeSlice:
		return g.genMakeSlice(v)

	case *ssa.Lookup:
		return g.genLookup(v)

	case *ssa.Convert:
		return g.genConvert(v)

	case *ssa.ChangeType:
		return g.genChangeType(v)

	case *ssa.MakeClosure:
		return g.genMakeClosre(v)
	}

	logger.Fatalf("Todo: %v, type: %T", v, v)
	return nil, nil
}

func (g *functionGenerator) genUnOp(inst *ssa.UnOp) (insts []wat.Inst, ret_type wir.ValueType) {
	switch inst.Op {
	case token.MUL: //*x
		return g.genLoad(inst.X)

	case token.SUB:
		x := g.getValue(inst.X)
		return wir.EmitUnOp(x.value, wat.OpCodeSub)

	default:
		logger.Fatal("Todo")
	}

	return
}

func (g *functionGenerator) genBinOp(inst *ssa.BinOp) ([]wat.Inst, wir.ValueType) {
	x := g.getValue(inst.X)
	y := g.getValue(inst.Y)

	switch inst.X.Type().Underlying().(type) {
	case *types.Basic:
		switch inst.Op {
		case token.ADD:
			return wir.EmitBinOp(x.value, y.value, wat.OpCodeAdd)

		case token.SUB:
			return wir.EmitBinOp(x.value, y.value, wat.OpCodeSub)

		case token.MUL:
			return wir.EmitBinOp(x.value, y.value, wat.OpCodeMul)

		case token.QUO:
			return wir.EmitBinOp(x.value, y.value, wat.OpCodeQuo)

		case token.REM:
			return wir.EmitBinOp(x.value, y.value, wat.OpCodeRem)

		case token.EQL:
			return wir.EmitBinOp(x.value, y.value, wat.OpCodeEql)

		case token.NEQ:
			return wir.EmitBinOp(x.value, y.value, wat.OpCodeNe)

		case token.LSS:
			return wir.EmitBinOp(x.value, y.value, wat.OpCodeLt)

		case token.GTR:
			return wir.EmitBinOp(x.value, y.value, wat.OpCodeGt)

		case token.LEQ:
			return wir.EmitBinOp(x.value, y.value, wat.OpCodeLe)

		case token.GEQ:
			return wir.EmitBinOp(x.value, y.value, wat.OpCodeGe)
		}

	default:
		logger.Fatalf("Todo: %v, type: %T, token:%v", inst, inst, inst.Op)
	}

	logger.Fatalf("Todo: %v, type: %T, token:%v", inst, inst, inst.Op)
	return nil, nil
}

func (g *functionGenerator) genCall(inst *ssa.Call) (insts []wat.Inst, ret_type wir.ValueType) {
	if inst.Call.IsInvoke() {
		logger.Fatal("Todo: genCall(), Invoke")
	}

	switch inst.Call.Value.(type) {
	case *ssa.Function:
		ret_type = wir.ToWType(inst.Call.Signature().Results())
		for _, v := range inst.Call.Args {
			insts = append(insts, g.getValue(v).value.EmitPush()...)
		}
		callee := inst.Call.StaticCallee()
		if callee.Parent() != nil {
			g.module.AddFunc(newFunctionGenerator(g.module).genFunction(callee))
		}

		if len(callee.LinkName()) > 0 {
			insts = append(insts, wat.NewInstCall(callee.LinkName()))
		} else {
			fn_internal_name, _ := GetFnMangleName(callee)
			insts = append(insts, wat.NewInstCall(fn_internal_name))
		}

	case *ssa.Builtin:
		return g.genBuiltin(inst.Common())

	case *ssa.MakeClosure:
		ret_type = wir.ToWType(inst.Type())
		var params []wir.Value
		for _, v := range inst.Call.Args {
			params = append(params, g.getValue(v).value)
		}
		closure := g.getValue(inst.Call.Value)
		insts = wir.EmitCallClosure(closure.value, params)

	default:
		ret_type = wir.ToWType(inst.Type())
		var params []wir.Value
		for _, v := range inst.Call.Args {
			params = append(params, g.getValue(v).value)
		}
		closure := g.getValue(inst.Call.Value)
		insts = wir.EmitCallClosure(closure.value, params)

	}

	return
}

func (g *functionGenerator) genBuiltin(call *ssa.CallCommon) (insts []wat.Inst, ret_type wir.ValueType) {
	switch call.Value.Name() {
	case "print", "println":
		for _, arg := range call.Args {
			arg := g.getValue(arg)
			switch arg.value.Type().(type) {
			case wir.I32:
				insts = append(insts, arg.value.EmitPush()...)
				insts = append(insts, wat.NewInstCall("$runtime.waPrintI32"))

			case wir.F32:
				insts = append(insts, arg.value.EmitPush()...)
				insts = append(insts, wat.NewInstCall("$runtime.waPrintF32"))

			case wir.F64:
				insts = append(insts, arg.value.EmitPush()...)
				insts = append(insts, wat.NewInstCall("$runtime.waPrintF64"))

			case wir.RUNE:
				insts = append(insts, arg.value.EmitPush()...)
				insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))

			case wir.String:
				insts = append(insts, wir.EmitPrintString(arg.value)...)

			default:
				logger.Fatalf("Todo: print(%T)", arg.value.Type())
			}
		}

		if call.Value.Name() == "println" {
			insts = append(insts, wir.NewConst(strconv.Itoa('\n'), wir.I32{}).EmitPush()...)
			insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))
		}
		ret_type = wir.VOID{}

	case "append":
		if len(call.Args) != 2 {
			panic("len(call.Args) != 2")
		}
		insts, ret_type = wir.EmitGenAppend(g.getValue(call.Args[0]).value, g.getValue(call.Args[1]).value)

	case "len":
		if len(call.Args) != 1 {
			panic("len(call.Args) != 1")
		}
		insts = wir.EmitGenLen(g.getValue(call.Args[0]).value)
		ret_type = wir.I32{}

	default:
		logger.Fatal("Todo:", call.Value)
	}
	return
}

func (g *functionGenerator) genPhiIter(preds []int, values []wir.Value) []wat.Inst {
	var insts []wat.Inst

	cond, _ := wir.EmitBinOp(g.var_current_block, wir.NewConst(strconv.Itoa(preds[0]), wir.I32{}), wat.OpCodeEql)
	insts = append(insts, cond...)

	trueInsts := values[0].EmitPush()
	var falseInsts []wat.Inst
	if len(preds) == 2 {
		falseInsts = values[1].EmitPush()
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
		values = append(values, g.getValue(v).value)
	}
	return g.genPhiIter(preds, values), wir.ToWType(inst.Type())
}

func (g *functionGenerator) genReturn(inst *ssa.Return) []wat.Inst {
	var insts []wat.Inst

	if len(inst.Results) != len(g.var_rets) {
		panic("len(inst.Results) != len(g.var_rets)")
	}

	for i := range inst.Results {
		insts = append(insts, wir.EmitAssginValue(g.var_rets[i], g.getValue(inst.Results[i]).value)...)
	}

	insts = append(insts, wat.NewInstBr("$BlockFnBody"))
	return insts
}

func (g *functionGenerator) genLoad(Addr ssa.Value) (insts []wat.Inst, ret_type wir.ValueType) {
	addr := g.getValue(Addr)

	if addr.force_register {
		insts = append(insts, addr.value.EmitPush()...)
		ret_type = addr.value.Type()
	} else {
		insts, ret_type = wir.EmitLoad(addr.value)
	}

	return
}

func (g *functionGenerator) genStore(inst *ssa.Store) []wat.Inst {
	addr := g.getValue(inst.Addr)
	val := g.getValue(inst.Val)

	if addr.force_register {
		return wir.EmitAssginValue(addr.value, val.value)
	} else {
		return wir.EmitStore(addr.value, val.value)
	}
}

func (g *functionGenerator) genIf(inst *ssa.If) []wat.Inst {
	cond := g.getValue(inst.Cond)
	if !cond.value.Type().Equal(wir.I32{}) {
		logger.Fatal("cond.type() != i32")
	}

	insts := cond.value.EmitPush()
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
		insts = wir.EmitAssginValue(g.var_block_selector, wir.NewConst(strconv.Itoa(dest), wir.I32{}))
		insts = append(insts, wat.NewInstBr("$BlockDisp"))
	} else {
		insts = append(insts, wat.NewInstBr("$Block_"+strconv.Itoa(dest-1)))
	}

	return insts
}

func (g *functionGenerator) genAlloc(inst *ssa.Alloc) (insts []wat.Inst, ret_type wir.ValueType) {
	typ := wir.ToWType(inst.Type().(*types.Pointer).Elem())
	if inst.Parent().ForceRegister() {
		nv := g.addRegister(typ)
		g.locals_map[inst] = valueWrap{value: nv, force_register: true}
		insts = append(insts, nv.EmitRelease()...)
		insts = append(insts, nv.EmitInit()...)
		ret_type = nil
	} else {
		if inst.Heap {
			insts, ret_type = wir.EmitHeapAlloc(typ)
		} else {
			insts, ret_type = wir.EmitStackAlloc(typ)
		}
	}

	return
}

func (g *functionGenerator) genExtract(inst *ssa.Extract) ([]wat.Inst, wir.ValueType) {
	v := g.getValue(inst.Tuple)
	return wir.EmitGenExtract(v.value, inst.Index)
}

func (g *functionGenerator) genFiled(inst *ssa.Field) ([]wat.Inst, wir.ValueType) {
	x := g.getValue(inst.X)
	field := inst.X.Type().Underlying().(*types.Struct).Field(inst.Field)
	fieldname := wir.GenSymbolName(field.Name())
	if field.Embedded() {
		if _, ok := field.Type().(*types.Named); ok {
			pkgname, _ := wir.GetPkgMangleName(field.Pkg().Path())
			fieldname = pkgname + "." + fieldname
		}
		fieldname = "$" + fieldname
	}

	return wir.EmitGenField(x.value, fieldname)
}

func (g *functionGenerator) genFieldAddr(inst *ssa.FieldAddr) ([]wat.Inst, wir.ValueType) {
	field := inst.X.Type().Underlying().(*types.Pointer).Elem().Underlying().(*types.Struct).Field(inst.Field)
	fieldname := wir.GenSymbolName(field.Name())
	if field.Embedded() {
		if _, ok := field.Type().(*types.Named); ok {
			pkgname, _ := wir.GetPkgMangleName(field.Pkg().Path())
			fieldname = pkgname + "." + fieldname
		}
		fieldname = "$" + fieldname
	}

	x := g.getValue(inst.X)
	if x.force_register {
		nv := wir.ExtractField(x.value, fieldname)
		g.locals_map[inst] = valueWrap{value: nv, force_register: true}
		return nil, nil
	} else {
		return wir.EmitGenFieldAddr(x.value, fieldname)
	}
}

func (g *functionGenerator) genIndexAddr(inst *ssa.IndexAddr) ([]wat.Inst, wir.ValueType) {
	if inst.Parent().ForceRegister() {
		logger.Fatal("ssa.IndexAddr is not available in ForceRegister-mode")
		return nil, nil
	}

	x := g.getValue(inst.X)
	id := g.getValue(inst.Index)

	return wir.EmitGenIndexAddr(x.value, id.value)
}

func (g *functionGenerator) genSlice(inst *ssa.Slice) ([]wat.Inst, wir.ValueType) {
	if inst.Parent().ForceRegister() {
		logger.Fatal("ssa.Slice is not available in ForceRegister-mode")
		return nil, nil
	}

	x := g.getValue(inst.X)
	var low, high wir.Value
	if inst.Low != nil {
		low = g.getValue(inst.Low).value
	}
	if inst.High != nil {
		high = g.getValue(inst.High).value
	}

	return wir.EmitGenSlice(x.value, low, high)
}

func (g *functionGenerator) genMakeSlice(inst *ssa.MakeSlice) ([]wat.Inst, wir.ValueType) {
	if inst.Parent().ForceRegister() {
		logger.Fatal("ssa.MakeSlice is not available in ForceRegister-mode")
		return nil, nil
	}

	elem := inst.Type().(*types.Slice).Elem()
	base_type := wir.ToWType(elem)
	Len := g.getValue(inst.Len)
	Cap := g.getValue(inst.Cap)
	return wir.EmitGenMakeSlice(base_type, Len.value, Cap.value)
}

func (g *functionGenerator) genLookup(inst *ssa.Lookup) ([]wat.Inst, wir.ValueType) {
	x := g.getValue(inst.X)
	index := g.getValue(inst.Index)

	return wir.EmitGenLookup(x.value, index.value, inst.CommaOk)
}

func (g *functionGenerator) genConvert(inst *ssa.Convert) (insts []wat.Inst, ret_type wir.ValueType) {
	x := g.getValue(inst.X)
	ret_type = wir.ToWType(inst.Type())
	insts = wir.EmitGenConvert(x.value, ret_type)
	return
}

func (g *functionGenerator) genChangeType(inst *ssa.ChangeType) (insts []wat.Inst, ret_type wir.ValueType) {
	ret_type = wir.ToWType(inst.Type())
	x := g.getValue(inst.X)
	insts = append(insts, x.value.EmitPush()...)
	return
}

func (g *functionGenerator) genMakeClosre(inst *ssa.MakeClosure) (insts []wat.Inst, ret_type wir.ValueType) {
	f := inst.Fn.(*ssa.Function)

	if f.Parent() != nil {
		return g.genMakeClosre_Anonymous(inst)
	}

	if len(f.FreeVars) == 1 && strings.HasSuffix(f.Name(), "$bound") {
		return g.genMakeClosre_Bound(inst)
	}

	panic("todo")
}

func (g *functionGenerator) genMakeClosre_Anonymous(inst *ssa.MakeClosure) (insts []wat.Inst, ret_type wir.ValueType) {
	f := inst.Fn.(*ssa.Function)

	g.module.AddFunc(newFunctionGenerator(g.module).genFunction(f))

	ret_type = wir.NewClosure(wir.NewFnSigFromSignature(f.Signature))
	if !ret_type.Equal(wir.ToWType(inst.Type())) {
		panic("ret_type != inst.Type()")
	}

	var st_free_data wir.Struct
	{
		var fields []wir.Field
		for _, freevar := range f.FreeVars {
			feild := wir.NewField(freevar.Name(), wir.ToWType(freevar.Type()))
			fields = append(fields, feild)
		}
		fn_internal_name, _ := GetFnMangleName(f)
		st_name := fn_internal_name + ".$warpdata"
		st_free_data = wir.NewStruct(st_name, fields)
	}

	var warp_fn_index int
	{
		var warp_fn wir.Function
		fn_name, _ := GetFnMangleName(f)
		warp_fn.InternalName = fn_name + ".$warpfn"
		for _, i := range f.Params {
			pa := valueWrap{value: wir.NewLocal(i.Name(), wir.ToWType(i.Type()))}
			warp_fn.Params = append(warp_fn.Params, pa.value)
		}
		warp_fn.Results = wir.NewFnSigFromSignature(f.Signature).Results

		dx := g.module.FindGlobal("$wa.RT.closure_data")
		data_ptr := wir.ExtractField(dx, "data")

		warp_fn.Insts = append(warp_fn.Insts, st_free_data.EmitLoadFromAddr(data_ptr, 0)...)
		warp_fn.Insts = append(warp_fn.Insts, dx.EmitRelease()...)
		warp_fn.Insts = append(warp_fn.Insts, dx.EmitInit()...)

		for _, i := range warp_fn.Params {
			warp_fn.Insts = append(warp_fn.Insts, i.EmitPush()...)
		}

		warp_fn.Insts = append(warp_fn.Insts, wat.NewInstCall(fn_name))

		g.module.AddFunc(&warp_fn)
		warp_fn_index = g.module.AddTableElem(warp_fn.InternalName)
	}

	closure := g.addRegister(wir.NewClosure(wir.NewFnSigFromSignature(f.Signature)))

	free_data := g.addRegister(st_free_data)
	{
		for i, freevar := range f.FreeVars {
			sv := g.getValue(inst.Bindings[i])
			insts = append(insts, sv.value.EmitPush()...)
			dv := wir.ExtractField(free_data, freevar.Name())
			insts = append(insts, dv.EmitPop()...)
		}
	}
	insts = append(insts, wir.NewConst(strconv.Itoa(warp_fn_index), wir.U32{}).EmitPush()...)
	insts = append(insts, wir.ExtractField(closure, "fn_index").EmitPop()...)
	{
		i, _ := wir.EmitHeapAlloc(st_free_data)
		insts = append(insts, i...)
		insts = append(insts, wir.ExtractField(closure, "data").EmitPop()...)
	}
	insts = append(insts, wir.EmitStore(wir.ExtractField(closure, "data"), free_data)...)
	insts = append(insts, free_data.EmitRelease()...)
	insts = append(insts, free_data.EmitInit()...)

	insts = append(insts, closure.EmitPush()...)
	return
}

func (g *functionGenerator) genMakeClosre_Bound(inst *ssa.MakeClosure) (insts []wat.Inst, ret_type wir.ValueType) {
	f := inst.Fn.(*ssa.Function)

	ret_type = wir.NewClosure(wir.NewFnSigFromSignature(f.Signature))
	if !ret_type.Equal(wir.ToWType(inst.Type())) {
		panic("ret_type != inst.Type()")
	}

	recv_type := wir.ToWType(f.FreeVars[0].Type())

	var warp_fn_index int
	{
		var warp_fn wir.Function
		fn_name, _ := GetFnMangleName(f.Object())
		warp_fn.InternalName = fn_name + ".$bound"
		for _, i := range f.Params {
			pa := valueWrap{value: wir.NewLocal(i.Name(), wir.ToWType(i.Type()))}
			warp_fn.Params = append(warp_fn.Params, pa.value)
		}
		warp_fn.Results = wir.NewFnSigFromSignature(f.Signature).Results

		dx := g.module.FindGlobal("$wa.RT.closure_data")
		data_ptr := wir.ExtractField(dx, "data")

		warp_fn.Insts = append(warp_fn.Insts, recv_type.EmitLoadFromAddr(data_ptr, 0)...)
		warp_fn.Insts = append(warp_fn.Insts, dx.EmitRelease()...)
		warp_fn.Insts = append(warp_fn.Insts, dx.EmitInit()...)

		for _, i := range warp_fn.Params {
			warp_fn.Insts = append(warp_fn.Insts, i.EmitPush()...)
		}

		warp_fn.Insts = append(warp_fn.Insts, wat.NewInstCall(fn_name))

		g.module.AddFunc(&warp_fn)
		warp_fn_index = g.module.AddTableElem(warp_fn.InternalName)
	}

	closure := g.addRegister(wir.NewClosure(wir.NewFnSigFromSignature(f.Signature)))

	insts = append(insts, wir.NewConst(strconv.Itoa(warp_fn_index), wir.U32{}).EmitPush()...)
	insts = append(insts, wir.ExtractField(closure, "fn_index").EmitPop()...)
	{
		i, _ := wir.EmitHeapAlloc(recv_type)
		insts = append(insts, i...)
		insts = append(insts, wir.ExtractField(closure, "data").EmitPop()...)
	}

	recv := g.getValue(inst.Bindings[0])
	insts = append(insts, wir.EmitStore(wir.ExtractField(closure, "data"), recv.value)...)

	insts = append(insts, closure.EmitPush()...)
	return
}

func (g *functionGenerator) addRegister(typ wir.ValueType) wir.Value {
	defer func() { g.cur_local_id++ }()
	name := "$T_" + strconv.Itoa(g.cur_local_id)
	v := wir.NewLocal(name, typ)
	g.registers = append(g.registers, v)
	return v
}

func (g *functionGenerator) genGetter(f *ssa.Function) *wir.Function {
	var wir_fn wir.Function
	if len(f.LinkName()) > 0 {
		wir_fn.InternalName = f.LinkName()
		wir_fn.ExternalName = f.LinkName()
	} else {
		wir_fn.InternalName, wir_fn.ExternalName = GetFnMangleName(f)
	}

	rets := f.Signature.Results()
	if rets.Len() > 1 {
		logger.Fatal("rets.Len() > 1")
		return nil
	}
	rtype := wir.ToWType(rets)
	wir_fn.Results = append(wir_fn.Results, rtype)

	if len(f.Params) != 1 {
		logger.Fatal("len(f.Params) != 1")
		return nil
	}
	if !wir.ToWType(f.Params[0].Type()).Equal(wir.U32{}) {
		logger.Fatal("addr_type != U32")
		return nil
	}
	addr := wir.NewLocal("addr", wir.NewPointer(rtype))
	wir_fn.Params = append(wir_fn.Params, addr)

	insts, _ := wir.EmitLoad(addr)
	wir_fn.Insts = append(wir_fn.Insts, insts...)

	return &wir_fn
}

func (g *functionGenerator) genSetter(f *ssa.Function) *wir.Function {
	var wir_fn wir.Function
	if len(f.LinkName()) > 0 {
		wir_fn.InternalName = f.LinkName()
		wir_fn.ExternalName = f.LinkName()
	} else {
		wir_fn.InternalName, wir_fn.ExternalName = GetFnMangleName(f)
	}

	rets := f.Signature.Results()
	if rets.Len() > 0 {
		logger.Fatal("rets.Len() > 0")
		return nil
	}

	if len(f.Params) != 2 {
		logger.Fatal("len(f.Params) != 2")
		return nil
	}
	if !wir.ToWType(f.Params[0].Type()).Equal(wir.U32{}) {
		logger.Fatal("addr_type != U32")
		return nil
	}

	value_type := wir.ToWType(f.Params[1].Type())

	addr := wir.NewLocal("addr", wir.NewPointer(value_type))
	wir_fn.Params = append(wir_fn.Params, addr)

	value := wir.NewLocal("data", value_type)
	wir_fn.Params = append(wir_fn.Params, value)

	insts := wir.EmitStore(addr, value)
	wir_fn.Insts = append(wir_fn.Insts, insts...)

	return &wir_fn
}
