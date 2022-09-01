package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wasm/wir/wtypes"
	"github.com/wa-lang/wa/internal/logger"
)

func EmitPushValue(x Value) []Instruction {
	var insts []Instruction
	vs := x.Raw()
	if len(vs) > 1 {
		for _, v := range vs {
			insts = append(insts, EmitPushValue(v)...)
		}
	}

	switch x.Kind() {
	case ValueKindConst:
		insts = append(insts, NewInstConst(x))

	case ValueKindLocal:
		insts = append(insts, NewInstGetLocal(x.Name()))

	case ValueKindGlobal:
		logger.Fatal("Todo")
	}

	return insts
}

func EmitPopValue(x Value) []Instruction {
	var insts []Instruction
	vs := x.Raw()
	if len(vs) > 1 {
		for _, v := range vs {
			insts = append(insts, EmitPopValue(v)...)
		}
	}

	switch x.Kind() {
	case ValueKindConst:
		logger.Fatal("不可Pop至常数")

	case ValueKindLocal:
		insts = append(insts, NewInstSetLocal(x.Name()))

	case ValueKindGlobal:
		logger.Fatal("Todo")
	}

	return insts
}

func EmitAssginValue(lh, rh Value) []Instruction {
	var insts []Instruction

	if rh == nil {
		ls := lh.Raw()
		for _, v := range ls {
			c := NewConst(v.Type(), nil)
			insts = append(insts, EmitPushValue(c)...)
			insts = append(insts, EmitPopValue(v)...)
		}
	} else {
		if !lh.Type().Equal(rh.Type()) {
			logger.Fatal("x.Type() != y.Type()")
		}

		ls := lh.Raw()
		rs := rh.Raw()

		for i := range ls {
			insts = append(insts, EmitPushValue(rs[i])...)
			insts = append(insts, EmitPopValue(ls[i])...)
		}
	}

	return insts
}

func EmitConvertValueType(from, to wtypes.ValueType) {
	logger.Fatal("Todo")
}

func EmitBinOp(x, y Value, op OpCode) ([]Instruction, wtypes.ValueType) {
	var insts []Instruction
	rtype := binOpMatchType(x.Type(), y.Type())

	insts = append(insts, EmitPushValue(x)...)
	insts = append(insts, EmitPushValue(y)...)

	switch op {
	case OpCodeAdd:
		insts = append(insts, NewInstAdd(rtype))

	case OpCodeSub:
		insts = append(insts, NewInstSub(rtype))

	case OpCodeMul:
		insts = append(insts, NewInstMul(rtype))

	case OpCodeQuo:
		insts = append(insts, NewInstDiv(rtype))

	case OpCodeRem:
		insts = append(insts, NewInstRem(rtype))

	case OpCodeEql:
		insts = append(insts, NewInstEq(rtype))

	case OpCodeNe:
		insts = append(insts, NewInstNe(rtype))

	case OpCodeLt:
		insts = append(insts, NewInstLt(rtype))

	case OpCodeGt:
		insts = append(insts, NewInstGt(rtype))

	case OpCodeLe:
		insts = append(insts, NewInstLe(rtype))

	case OpCodeGe:
		insts = append(insts, NewInstGe(rtype))

	default:
		logger.Fatal("Todo")
	}

	return insts, rtype
}

func binOpMatchType(x, y wtypes.ValueType) wtypes.ValueType {
	if x.Equal(y) {
		return x
	}

	logger.Fatalf("Todo %T %T", x, y)
	return nil
}
