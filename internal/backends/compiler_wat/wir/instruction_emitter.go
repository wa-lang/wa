// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

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
		r = I32{}

	case wat.OpCodeNe:
		insts = append(insts, wat.NewInstNe(rtype))
		r = I32{}

	case wat.OpCodeLt:
		insts = append(insts, wat.NewInstLt(rtype))
		r = I32{}

	case wat.OpCodeGt:
		insts = append(insts, wat.NewInstGt(rtype))
		r = I32{}

	case wat.OpCodeLe:
		insts = append(insts, wat.NewInstLe(rtype))
		r = I32{}

	case wat.OpCodeGe:
		insts = append(insts, wat.NewInstGe(rtype))
		r = I32{}

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
	switch addr.Kind() {
	case ValueKindGlobal_Value:
		insts = append(insts, addr.EmitPush()...)
		ret_type = addr.Type()

	default:
		switch addr := addr.(type) {
		case *aRef:
			insts = append(insts, addr.emitGetValue()...)
			ret_type = addr.Type().(Ref).Base

		case *aPointer:
			insts = append(insts, addr.emitGetValue()...)
			ret_type = addr.Type().(Pointer).Base

		default:
			logger.Fatalf("Todo %v", addr)
		}
	}

	return
}

func EmitStore(addr, value Value) (insts []wat.Inst) {
	switch addr.Kind() {
	case ValueKindGlobal_Value:
		if value == nil {
			zero_value := NewConst("0", addr.Type())
			insts = append(insts, zero_value.EmitPush()...)
			insts = append(insts, addr.EmitPop()...)
		} else {
			if !addr.Type().Equal(value.Type()) {
				logger.Fatal("Type not match")
				return nil
			}
			insts = append(insts, value.EmitPush()...)
			insts = append(insts, addr.EmitPop()...)
		}

	default:
		switch addr := addr.(type) {
		case *aRef:
			if value == nil {
				zero_value := NewConst("0", addr.Type().(Ref).Base)
				insts = append(insts, addr.emitSetValue(zero_value)...)
			} else {
				insts = append(insts, addr.emitSetValue(value)...)
			}

		case *aPointer:
			if value == nil {
				zero_value := NewConst("0", addr.Type().(Pointer).Base)
				insts = append(insts, addr.emitSetValue(zero_value)...)
			} else {
				insts = append(insts, addr.emitSetValue(value)...)
			}

		default:
			logger.Fatalf("Todo %v", addr)
		}
	}

	return
}

func EmitHeapAlloc(typ ValueType, module *Module) (insts []wat.Inst, ret_type ValueType) {
	ref_typ := NewRef(typ)
	ret_type = ref_typ
	insts = ref_typ.emitHeapAlloc(module)
	return
}

func EmitStackAlloc(typ ValueType, module *Module) (insts []wat.Inst, ret_type ValueType) {
	return EmitHeapAlloc(typ, module)
	//ref_typ := NewRef(typ)
	//ret_type = ref_typ
	//insts = ref_typ.emitStackAlloc(module)
	//return
}

func EmitGenFieldAddr(x Value, field_name string) (insts []wat.Inst, ret_type ValueType) {
	insts = append(insts, x.EmitPush()...)
	var field *Field
	switch addr := x.(type) {
	case *aRef:
		field = addr.Type().(Ref).Base.(Struct).findFieldByName(field_name)

		ret_type = NewRef(field.Type())
	case *aPointer:
		field = addr.Type().(Pointer).Base.(Struct).findFieldByName(field_name)

	default:
		logger.Fatalf("Todo:%T", x.Type())
	}

	insts = append(insts, NewConst(strconv.Itoa(field._start), I32{}).EmitPush()...)
	insts = append(insts, wat.NewInstAdd(wat.I32{}))
	return
}

func EmitGenIndexAddr(x, id Value) (insts []wat.Inst, ret_type ValueType) {
	if !id.Type().Equal(I32{}) {
		panic("index should be i32")
	}

	switch x := x.(type) {
	case *aPointer:
		switch typ := x.Type().(Pointer).Base.(type) {
		case Array:
			insts = append(insts, x.EmitPush()...)
			insts = append(insts, NewConst(strconv.Itoa(typ.Base.size()), I32{}).EmitPush()...)
			insts = append(insts, id.EmitPush()...)
			insts = append(insts, wat.NewInstMul(wat.I32{}))
			insts = append(insts, wat.NewInstAdd(wat.I32{}))
			ret_type = NewPointer(typ.Base)

		default:
			logger.Fatalf("Todo: %T", typ)
		}

	case *aRef:
		switch typ := x.Type().(Ref).Base.(type) {
		case Array:
			insts = append(insts, x.EmitPush()...)
			insts = append(insts, NewConst(strconv.Itoa(typ.Base.size()), I32{}).EmitPush()...)
			insts = append(insts, id.EmitPush()...)
			insts = append(insts, wat.NewInstMul(wat.I32{}))
			insts = append(insts, wat.NewInstAdd(wat.I32{}))
			ret_type = NewRef(typ.Base)

		default:
			logger.Fatalf("Todo: %T", typ)
		}

	default:
		logger.Fatalf("Todo: %T", x)
	}

	return
}
