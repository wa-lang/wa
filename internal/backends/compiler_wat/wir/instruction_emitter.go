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

func EmitUnOp(x Value, op wat.OpCode) (insts []wat.Inst, ret_type ValueType) {
	if !IsNumber(x) {
		logger.Fatal("Todo")
	}
	ret_type = x.Type()
	insts = append(insts, NewConst("0", ret_type).EmitPush()...)
	insts = append(insts, x.EmitPush()...)
	insts = append(insts, wat.NewInstSub(toWatType(ret_type)))
	return
}

func EmitBinOp(x, y Value, op wat.OpCode) (insts []wat.Inst, ret_type ValueType) {
	rtype := binOpMatchType(x.Type(), y.Type())

	insts = append(insts, x.EmitPush()...)
	insts = append(insts, y.EmitPush()...)

	switch op {
	case wat.OpCodeAdd:
		string_type := NewString()
		if rtype.Equal(string_type) {
			insts = append(insts, wat.NewInstCall(string_type.genAppendStrFunc()))
		} else {
			insts = append(insts, wat.NewInstAdd(toWatType(rtype)))
		}
		ret_type = rtype

	case wat.OpCodeSub:
		insts = append(insts, wat.NewInstSub(toWatType(rtype)))
		ret_type = rtype

	case wat.OpCodeMul:
		insts = append(insts, wat.NewInstMul(toWatType(rtype)))
		ret_type = rtype

	case wat.OpCodeQuo:
		insts = append(insts, wat.NewInstDiv(toWatType(rtype)))
		ret_type = rtype

	case wat.OpCodeRem:
		insts = append(insts, wat.NewInstRem(toWatType(rtype)))
		ret_type = rtype

	case wat.OpCodeEql:
		insts = append(insts, wat.NewInstEq(toWatType(rtype)))
		ret_type = I32{}

	case wat.OpCodeNe:
		insts = append(insts, wat.NewInstNe(toWatType(rtype)))
		ret_type = I32{}

	case wat.OpCodeLt:
		insts = append(insts, wat.NewInstLt(toWatType(rtype)))
		ret_type = I32{}

	case wat.OpCodeGt:
		insts = append(insts, wat.NewInstGt(toWatType(rtype)))
		ret_type = I32{}

	case wat.OpCodeLe:
		insts = append(insts, wat.NewInstLe(toWatType(rtype)))
		ret_type = I32{}

	case wat.OpCodeGe:
		insts = append(insts, wat.NewInstGe(toWatType(rtype)))
		ret_type = I32{}

	default:
		logger.Fatal("Todo")
	}

	return
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

func EmitHeapAlloc(typ ValueType) (insts []wat.Inst, ret_type ValueType) {
	ref_typ := NewRef(typ)
	ret_type = ref_typ
	insts = ref_typ.emitHeapAlloc()
	return
}

func EmitStackAlloc(typ ValueType) (insts []wat.Inst, ret_type ValueType) {
	return EmitHeapAlloc(typ)
	//ref_typ := NewRef(typ)
	//ret_type = ref_typ
	//insts = ref_typ.emitStackAlloc(module)
	//return
}

func EmitGenExtract(x Value, id int) (insts []wat.Inst, ret_type ValueType) {
	f := x.(*aTuple).Extract(id)
	insts = append(insts, f.EmitPush()...)
	ret_type = f.Type()
	return
}

func EmitGenField(x Value, field_name string) (insts []wat.Inst, ret_type ValueType) {
	switch x := x.(type) {
	case *aStruct:
		field := x.Extract(field_name)
		insts = append(insts, field.EmitPush()...)
		ret_type = field.Type()

	default:
		logger.Fatalf("Todo:%T", x)
	}

	return
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
		ret_type = NewPointer(field.Type())

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

	case *aSlice:
		base_type := x.Type().(Slice).Base
		insts = append(insts, x.Extract("block").EmitPush()...)
		insts = append(insts, x.Extract("data").EmitPush()...)
		insts = append(insts, NewConst(strconv.Itoa(base_type.size()), I32{}).EmitPush()...)
		insts = append(insts, id.EmitPush()...)
		insts = append(insts, wat.NewInstMul(wat.I32{}))
		insts = append(insts, wat.NewInstAdd(wat.I32{}))
		ret_type = NewRef(base_type)

	default:
		logger.Fatalf("Todo: %T", x)
	}

	return
}

func EmitGenSlice(x, low, high Value) (insts []wat.Inst, ret_type ValueType) {
	switch x := x.(type) {
	case *aSlice:
		insts = x.emitSub(low, high)
		ret_type = x.Type()

	case *aString:
		insts = x.emitSub(low, high)
		ret_type = x.Type()

	case *aRef:
		switch btype := x.Type().(Ref).Base.(type) {
		case Slice:
			slt := NewSlice(btype.Base)
			insts = slt.emitGenFromRefOfSlice(x, low, high)
			ret_type = slt

		case Array:
			slt := NewSlice(btype.Base)
			insts = slt.emitGenFromRefOfArray(x, low, high)
			ret_type = slt

		default:
			logger.Fatalf("Todo: %T", btype)
		}

	default:
		logger.Fatalf("Todo: %T", x)
	}

	return
}

func EmitGenLookup(x, index Value, CommaOk bool) (insts []wat.Inst, ret_type ValueType) {
	switch x := x.(type) {
	case *aString:
		if CommaOk {
			insts = x.emitAt_CommaOk(index)
			fileds := []ValueType{U8{}, I32{}}
			ret_type = NewTuple(fileds)
		} else {
			insts = x.emitAt(index)
			ret_type = U8{}
		}

	default:
		logger.Fatalf("Todo: %T", x)
	}

	return
}

func EmitGenConvert(x Value, typ ValueType) (insts []wat.Inst) {
	src_raw_type := x.Type().Raw()
	dest_raw_type := typ.Raw()
	if len(src_raw_type) != len(dest_raw_type) {
		logger.Fatalf("Todo: %T %T", x, typ)
		panic("Todo")
	}
	for i := range src_raw_type {
		if src_raw_type[i].Name() != dest_raw_type[i].Name() {
			logger.Fatalf("Todo: %T %T", x, typ)
			panic("Todo")
		}
	}

	insts = append(insts, x.EmitPush()...)
	return
}

func EmitGenAppend(x, y Value) (insts []wat.Inst, ret_type ValueType) {
	if !x.Type().Equal(y.Type()) {
		logger.Fatal("Type not match")
		return
	}

	stype := x.Type().(Slice)
	insts = append(insts, x.EmitPush()...)
	insts = append(insts, y.EmitPush()...)
	insts = append(insts, wat.NewInstCall(stype.genAppendFunc()))
	ret_type = stype

	return
}

func EmitGenLen(x Value) (insts []wat.Inst) {
	switch x := x.(type) {
	case *aArray:
		insts = NewConst(strconv.Itoa(x.Type().(Array).Capacity), I32{}).EmitPush()

	case *aSlice:
		insts = x.Extract("len").EmitPush()

	case *aString:
		insts = x.Extract("len").EmitPush()

	default:
		logger.Fatalf("Todo: %T", x)
	}

	return
}

func EmitPrintString(v Value) (insts []wat.Inst) {
	s := v.(*aString)

	insts = append(insts, s.Extract("data").EmitPush()...)
	insts = append(insts, s.Extract("len").EmitPush()...)
	insts = append(insts, wat.NewInstCall("$runtime.waPuts"))
	return
}

func emitPrintValue(v Value) (insts []wat.Inst) {

	panic("Todo")
}
