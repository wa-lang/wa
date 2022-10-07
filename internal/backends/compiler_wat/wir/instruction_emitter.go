// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

func EmitAssginValue(lh, rh Value) []wat.Inst {
	var insts []wat.Inst

	if rh == nil {
		insts = append(insts, lh.EmitRelease()...)
		insts = append(insts, lh.EmitInit()...)
	} else {
		if !lh.Type().Equal(rh.Type()) {
			logger.Fatal("x.Type() != y.Type()")
		}

		insts = append(insts, rh.EmitPush()...)
		insts = append(insts, lh.EmitPop()...)
	}

	return insts
}

func EmitConvertValueType(from, to ValueType) {
	logger.Fatal("Todo")
}

func EmitBinOp(x, y Value, op wat.OpCode) ([]wat.Inst, ValueType) {
	var insts []wat.Inst
	r := binOpMatchType(x.Type(), y.Type())
	if len(r.Raw()) != 1 {
		logger.Fatalf("Todo %T", r)
		return nil, nil
	}
	rtype := r.Raw()[0]

	insts = append(insts, x.EmitPush()...)
	insts = append(insts, y.EmitPush()...)

	switch op {
	case wat.OpCodeAdd:
		insts = append(insts, wat.NewInstAdd(rtype))

	case wat.OpCodeSub:
		insts = append(insts, wat.NewInstSub(rtype))

	case wat.OpCodeMul:
		insts = append(insts, wat.NewInstMul(rtype))

	case wat.OpCodeQuo:
		insts = append(insts, wat.NewInstDiv(rtype))

	case wat.OpCodeRem:
		insts = append(insts, wat.NewInstRem(rtype))

	case wat.OpCodeEql:
		insts = append(insts, wat.NewInstEq(rtype))

	case wat.OpCodeNe:
		insts = append(insts, wat.NewInstNe(rtype))

	case wat.OpCodeLt:
		insts = append(insts, wat.NewInstLt(rtype))

	case wat.OpCodeGt:
		insts = append(insts, wat.NewInstGt(rtype))

	case wat.OpCodeLe:
		insts = append(insts, wat.NewInstLe(rtype))

	case wat.OpCodeGe:
		insts = append(insts, wat.NewInstGe(rtype))

	default:
		logger.Fatal("Todo")
	}

	return insts, r
}

func binOpMatchType(x, y ValueType) ValueType {
	if x.Equal(y) {
		return x
	}

	logger.Fatalf("Todo %T %T", x, y)
	return nil
}

func EmitLoad(addr Value) (insts []wat.Inst, ret_type ValueType) {
	switch addr := addr.(type) {
	case *VarRef:
		ret_type = addr.Type().(Ref).Base
		insts = addr.EmitLoad()

	default:
		logger.Fatalf("Todo %v", addr)
	}
	return
}

func EmitStore(addr, value Value) (insts []wat.Inst) {
	switch addr := addr.(type) {
	case *VarRef:
		insts = addr.EmitStore(value)

	default:
		logger.Fatalf("Todo %v", addr)
	}

	return
}

func EmitHeapAlloc(typ ValueType, module *Module) (insts []wat.Inst, ret_type ValueType) {
	ret_type = NewRef(typ)
	insts = NewVarRef("", ValueKindLocal, typ).emitHeapAlloc(module)
	return
}

func EmitStackAlloc(typ ValueType, module *Module) (insts []wat.Inst, ret_type ValueType) {
	ret_type = NewRef(typ)
	insts = NewVarRef("", ValueKindLocal, typ).emitStackAlloc(module)
	return
}
