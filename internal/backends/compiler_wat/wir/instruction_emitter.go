package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

//func EmitPushValue(x Value) []Instruction {
//	var insts []Instruction
//	vs := x.Raw()
//	if len(vs) > 1 {
//		for _, v := range vs {
//			insts = append(insts, EmitPushValue(v)...)
//		}
//	}
//
//	switch x.Kind() {
//	case ValueKindConst:
//		insts = append(insts, NewInstConst(x))
//
//	case ValueKindLocal:
//		insts = append(insts, NewInstGetLocal(x.Name()))
//
//	case ValueKindGlobal:
//		logger.Fatal("Todo")
//	}
//
//	return insts
//}

//func EmitPopValue(x Value) []Instruction {
//	var insts []Instruction
//	vs := x.Raw()
//	if len(vs) > 1 {
//		for _, v := range vs {
//			insts = append(insts, EmitPopValue(v)...)
//		}
//	}
//
//	switch x.Kind() {
//	case ValueKindConst:
//		logger.Fatal("不可Pop至常数")
//
//	case ValueKindLocal:
//		insts = append(insts, NewInstSetLocal(x.Name()))
//
//	case ValueKindGlobal:
//		logger.Fatal("Todo")
//	}
//
//	return insts
//}

func EmitAssginValue(lh, rh Value) []wat.Inst {
	var insts []wat.Inst

	if rh == nil {
		insts = append(insts, lh.EmitRelease()...)
		insts = append(insts, lh.EmitInit()...)
	} else {
		if !lh.Type().Equal(rh.Type()) {
			logger.Fatal("x.Type() != y.Type()")
		}

		insts = append(insts, rh.EmitGet()...)
		insts = append(insts, lh.EmitSet()...)
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

	insts = append(insts, x.EmitGet()...)
	insts = append(insts, y.EmitGet()...)

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
