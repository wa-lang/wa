// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

func (m *Module) EmitAssginValue(lh, rh Value) []wat.Inst {
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

func (m *Module) EmitConvertValueType(from, to ValueType) {
	logger.Fatal("Todo")
}

func (m *Module) EmitUnOp(x Value, op wat.OpCode) (insts []wat.Inst, ret_type ValueType) {
	if !IsNumber(x) {
		logger.Fatal("Todo")
	}
	ret_type = x.Type()
	insts = append(insts, NewConst("0", ret_type).EmitPush()...)
	insts = append(insts, x.EmitPush()...)
	insts = append(insts, wat.NewInstSub(toWatType(ret_type)))
	return
}

func (m *Module) EmitBinOp(x, y Value, op wat.OpCode) (insts []wat.Inst, ret_type ValueType) {
	rtype := m.binOpMatchType(x.Type(), y.Type())

	insts = append(insts, x.EmitPush()...)
	insts = append(insts, y.EmitPush()...)

	switch op {
	case wat.OpCodeAdd:
		if rtype.Equal(m.STRING) {
			insts = append(insts, wat.NewInstCall(m.STRING.(*String).genAppendStrFunc()))
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
		ret_type = m.I32

	case wat.OpCodeNe:
		insts = append(insts, wat.NewInstNe(toWatType(rtype)))
		ret_type = m.I32

	case wat.OpCodeLt:
		insts = append(insts, wat.NewInstLt(toWatType(rtype)))
		ret_type = m.I32

	case wat.OpCodeGt:
		insts = append(insts, wat.NewInstGt(toWatType(rtype)))
		ret_type = m.I32

	case wat.OpCodeLe:
		insts = append(insts, wat.NewInstLe(toWatType(rtype)))
		ret_type = m.I32

	case wat.OpCodeGe:
		insts = append(insts, wat.NewInstGe(toWatType(rtype)))
		ret_type = m.I32

	default:
		logger.Fatal("Todo")
	}

	return
}

func (m *Module) binOpMatchType(x, y ValueType) ValueType {
	if x.Equal(y) {
		return x
	}

	logger.Fatalf("Todo %T %T", x, y)
	return nil
}

func (m *Module) EmitLoad(addr Value) (insts []wat.Inst, ret_type ValueType) {
	switch addr.Kind() {
	case ValueKindGlobal_Value:
		insts = append(insts, addr.EmitPush()...)
		ret_type = addr.Type()

	default:
		switch addr := addr.(type) {
		case *aRef:
			insts = append(insts, addr.emitGetValue()...)
			ret_type = addr.Type().(*Ref).Base

		case *aPtr:
			insts = append(insts, addr.emitGetValue()...)
			ret_type = addr.Type().(*Ptr).Base

		default:
			logger.Fatalf("Todo %v", addr)
		}
	}

	return
}

func (m *Module) EmitStore(addr, value Value) (insts []wat.Inst) {
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
				zero_value := NewConst("0", addr.Type().(*Ref).Base)
				insts = append(insts, addr.emitSetValue(zero_value)...)
			} else {
				insts = append(insts, addr.emitSetValue(value)...)
			}

		case *aPtr:
			if value == nil {
				zero_value := NewConst("0", addr.Type().(*Ptr).Base)
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

func (m *Module) EmitHeapAlloc(typ ValueType) (insts []wat.Inst, ret_type ValueType) {
	ref_typ := m.GenValueType_Ref(typ)
	ret_type = ref_typ
	insts = ref_typ.emitHeapAlloc()
	return
}

func (m *Module) EmitStackAlloc(typ ValueType) (insts []wat.Inst, ret_type ValueType) {
	return m.EmitHeapAlloc(typ)
	//ref_typ := NewRef(typ)
	//ret_type = ref_typ
	//insts = ref_typ.emitStackAlloc(module)
	//return
}

func (m *Module) EmitGenExtract(x Value, id int) (insts []wat.Inst, ret_type ValueType) {
	f := x.(*aTuple).Extract(id)
	insts = append(insts, f.EmitPush()...)
	ret_type = f.Type()
	return
}

func (m *Module) EmitGenField(x Value, field_name string) (insts []wat.Inst, ret_type ValueType) {
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

func (m *Module) EmitGenFieldAddr(x Value, field_name string) (insts []wat.Inst, ret_type ValueType) {
	insts = append(insts, x.EmitPush()...)
	var field *Field
	switch addr := x.(type) {
	case *aRef:
		field = addr.Type().(*Ref).Base.(*Struct).findFieldByName(field_name)
		ret_type = m.GenValueType_Ref(field.Type())
	case *aPtr:
		field = addr.Type().(*Ptr).Base.(*Struct).findFieldByName(field_name)
		ret_type = m.GenValueType_Ptr(field.Type())

	default:
		logger.Fatalf("Todo:%T", x.Type())
	}

	insts = append(insts, NewConst(strconv.Itoa(field._start), m.I32).EmitPush()...)
	insts = append(insts, wat.NewInstAdd(wat.I32{}))
	return
}

func (m *Module) EmitGenIndexAddr(x, id Value) (insts []wat.Inst, ret_type ValueType) {
	if !id.Type().Equal(m.I32) {
		panic("index should be i32")
	}

	switch x := x.(type) {
	case *aPtr:
		switch typ := x.Type().(*Ptr).Base.(type) {
		case *Array:
			insts = append(insts, x.EmitPush()...)
			insts = append(insts, NewConst(strconv.Itoa(typ.Base.size()), m.I32).EmitPush()...)
			insts = append(insts, id.EmitPush()...)
			insts = append(insts, wat.NewInstMul(wat.I32{}))
			insts = append(insts, wat.NewInstAdd(wat.I32{}))
			ret_type = m.GenValueType_Ptr(typ.Base)

		default:
			logger.Fatalf("Todo: %T", typ)
		}

	case *aRef:
		switch typ := x.Type().(*Ref).Base.(type) {
		case *Array:
			insts = append(insts, x.EmitPush()...)
			insts = append(insts, NewConst(strconv.Itoa(typ.Base.size()), m.I32).EmitPush()...)
			insts = append(insts, id.EmitPush()...)
			insts = append(insts, wat.NewInstMul(wat.I32{}))
			insts = append(insts, wat.NewInstAdd(wat.I32{}))
			ret_type = m.GenValueType_Ref(typ.Base)

		default:
			logger.Fatalf("Todo: %T", typ)
		}

	case *aSlice:
		base_type := x.Type().(*Slice).Base
		insts = append(insts, x.Extract("block").EmitPush()...)
		insts = append(insts, x.Extract("data").EmitPush()...)
		insts = append(insts, NewConst(strconv.Itoa(base_type.size()), m.I32).EmitPush()...)
		insts = append(insts, id.EmitPush()...)
		insts = append(insts, wat.NewInstMul(wat.I32{}))
		insts = append(insts, wat.NewInstAdd(wat.I32{}))
		ret_type = m.GenValueType_Ref(base_type)

	default:
		logger.Fatalf("Todo: %T", x)
	}

	return
}

func (m *Module) EmitGenSlice(x, low, high Value) (insts []wat.Inst, ret_type ValueType) {
	switch x := x.(type) {
	case *aSlice:
		insts = x.emitSub(low, high)
		ret_type = x.Type()

	case *aString:
		insts = x.emitSub(low, high)
		ret_type = x.Type()

	case *aRef:
		switch btype := x.Type().(*Ref).Base.(type) {
		case *Slice:
			slt := m.genValueType_Slice(btype.Base)
			insts = slt.emitGenFromRefOfSlice(x, low, high)
			ret_type = slt

		case *Array:
			slt := m.genValueType_Slice(btype.Base)
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

func (m *Module) EmitGenMakeSlice(base_type ValueType, Len, Cap Value) (insts []wat.Inst, ret_type ValueType) {
	slice_type := m.genValueType_Slice(base_type)
	insts = slice_type.emitGenMake(Len, Cap)
	ret_type = slice_type
	return
}

func (m *Module) EmitGenLookup(x, index Value, CommaOk bool) (insts []wat.Inst, ret_type ValueType) {
	switch x := x.(type) {
	case *aString:
		if CommaOk {
			insts = x.emitAt_CommaOk(index)
			fileds := []ValueType{m.U8, m.I32}
			ret_type = m.GenValueType_Tuple(fileds)
		} else {
			insts = x.emitAt(index)
			ret_type = m.U8
		}

	default:
		logger.Fatalf("Todo: %T", x)
	}

	return
}

func (m *Module) EmitGenConvert(x Value, typ ValueType) (insts []wat.Inst) {
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

func (m *Module) EmitGenAppend(x, y Value) (insts []wat.Inst, ret_type ValueType) {
	if !x.Type().Equal(y.Type()) {
		logger.Fatal("Type not match")
		return
	}

	stype := x.Type().(*Slice)
	insts = append(insts, x.EmitPush()...)
	insts = append(insts, y.EmitPush()...)
	insts = append(insts, wat.NewInstCall(stype.genAppendFunc()))
	ret_type = stype

	return
}

func (m *Module) EmitGenLen(x Value) (insts []wat.Inst) {
	switch x := x.(type) {
	case *aArray:
		insts = NewConst(strconv.Itoa(x.Type().(*Array).Capacity), m.I32).EmitPush()

	case *aSlice:
		insts = x.Extract("len").EmitPush()

	case *aString:
		insts = x.Extract("len").EmitPush()

	default:
		logger.Fatalf("Todo: %T", x)
	}

	return
}

func (m *Module) EmitPrintString(v Value) (insts []wat.Inst) {
	s := v.(*aString)

	insts = append(insts, s.Extract("data").EmitPush()...)
	insts = append(insts, s.Extract("len").EmitPush()...)
	insts = append(insts, wat.NewInstCall("$runtime.waPuts"))
	return
}

func (m *Module) emitPrintValue(v Value) (insts []wat.Inst) {

	panic("Todo")
}
