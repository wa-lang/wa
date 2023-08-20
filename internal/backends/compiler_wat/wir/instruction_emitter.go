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

	switch op {
	case wat.OpCodeSub:
		insts = append(insts, NewConst("0", ret_type).EmitPush()...)
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, wat.NewInstSub(toWatType(ret_type)))
		if ret_type.Equal(m.U8) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "255"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		} else if ret_type.Equal(m.U16) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "65535"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		}

	case wat.OpCodeXor:
		insts = append(insts, NewConst("-1", ret_type).EmitPush()...)
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, wat.NewInstXor(toWatType(ret_type)))
		if ret_type.Equal(m.U8) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "255"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		} else if ret_type.Equal(m.U16) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "65535"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		}

	case wat.OpCodeNot:
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, wat.NewInstEqz(toWatType(ret_type)))

	default:
		logger.Fatalf("Todo: %[1]v: %[1]T", op)

	}

	return
}

func (m *Module) EmitBinOp(x, y Value, op wat.OpCode) (insts []wat.Inst, ret_type ValueType) {
	for {
		if ut, ok := x.(*aDup); ok {
			x = ut.underlying
		} else {
			break
		}
	}

	for {
		if ut, ok := y.(*aDup); ok {
			y = ut.underlying
		} else {
			break
		}
	}

	switch op {
	case wat.OpCodeAdd:
		ret_type = x.Type()
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, y.EmitPush()...)
		if ret_type.Equal(m.STRING) {
			insts = append(insts, wat.NewInstCall(m.STRING.(*String).genFunc_Append()))
		} else {
			insts = append(insts, wat.NewInstAdd(toWatType(ret_type)))
		}

		if ret_type.Equal(m.U8) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "255"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		} else if ret_type.Equal(m.U16) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "65535"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		}

	case wat.OpCodeSub:
		ret_type = x.Type()
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, y.EmitPush()...)
		insts = append(insts, wat.NewInstSub(toWatType(ret_type)))

		if ret_type.Equal(m.U8) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "255"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		} else if ret_type.Equal(m.U16) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "65535"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		}

	case wat.OpCodeMul:
		ret_type = x.Type()
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, y.EmitPush()...)
		insts = append(insts, wat.NewInstMul(toWatType(ret_type)))

		if ret_type.Equal(m.U8) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "255"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		} else if ret_type.Equal(m.U16) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "65535"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		}

	case wat.OpCodeQuo:
		ret_type = x.Type()
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, y.EmitPush()...)
		insts = append(insts, wat.NewInstDiv(toWatType(ret_type)))

	case wat.OpCodeRem:
		ret_type = x.Type()
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, y.EmitPush()...)
		insts = append(insts, wat.NewInstRem(toWatType(ret_type)))

	case wat.OpCodeEql:
		ins, _ := x.emitEq(y)
		insts = append(insts, ins...)
		ret_type = m.BOOL

	case wat.OpCodeNe:
		ins, _ := x.emitEq(y)
		insts = append(insts, ins...)
		insts = append(insts, wat.NewInstEqz(wat.I32{}))
		ret_type = m.BOOL

	case wat.OpCodeLt:
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, y.EmitPush()...)
		insts = append(insts, wat.NewInstLt(toWatType(x.Type())))
		ret_type = m.BOOL

	case wat.OpCodeGt:
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, y.EmitPush()...)
		insts = append(insts, wat.NewInstGt(toWatType(x.Type())))
		ret_type = m.BOOL

	case wat.OpCodeLe:
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, y.EmitPush()...)
		insts = append(insts, wat.NewInstLe(toWatType(x.Type())))
		ret_type = m.BOOL

	case wat.OpCodeGe:
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, y.EmitPush()...)
		insts = append(insts, wat.NewInstGe(toWatType(x.Type())))
		ret_type = m.BOOL

	case wat.OpCodeAnd:
		ret_type = x.Type()
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, y.EmitPush()...)
		insts = append(insts, wat.NewInstAnd(toWatType(ret_type)))

	case wat.OpCodeOr:
		ret_type = x.Type()
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, y.EmitPush()...)
		insts = append(insts, wat.NewInstOr(toWatType(ret_type)))

	case wat.OpCodeXor:
		ret_type = x.Type()
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, y.EmitPush()...)
		insts = append(insts, wat.NewInstXor(toWatType(ret_type)))

	case wat.OpCodeShl:
		ret_type = x.Type()

		if x.Type().Size() <= 4 && y.Type().Size() == 8 {
			insts = append(insts, x.EmitPush()...)
			insts = append(insts, wat.NewInstConvert_i64_extend_i32_u())
			insts = append(insts, y.EmitPush()...)
			insts = append(insts, wat.NewInstShl(toWatType(y.Type())))
			insts = append(insts, wat.NewInstConvert_i32_wrap_i64())
		} else if x.Type().Size() == 8 && y.Type().Size() <= 4 {
			insts = append(insts, x.EmitPush()...)
			insts = append(insts, y.EmitPush()...)
			insts = append(insts, wat.NewInstConvert_i64_extend_i32_u())
			insts = append(insts, wat.NewInstShl(toWatType(ret_type)))
		} else if (x.Type().Size() <= 4 && y.Type().Size() <= 4) || (x.Type().Size() == 8 && y.Type().Size() == 8) {
			insts = append(insts, x.EmitPush()...)
			insts = append(insts, y.EmitPush()...)
			insts = append(insts, wat.NewInstShl(toWatType(ret_type)))
		} else {
			logger.Fatal("Unreachable")
		}

		if ret_type.Equal(m.U8) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "255"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		} else if ret_type.Equal(m.U16) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "65535"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		}

	case wat.OpCodeShr:
		ret_type = x.Type()

		if x.Type().Size() <= 4 && y.Type().Size() == 8 {
			insts = append(insts, x.EmitPush()...)
			insts = append(insts, wat.NewInstConvert_i64_extend_i32_u())
			insts = append(insts, y.EmitPush()...)
			insts = append(insts, wat.NewInstShr(toWatType(y.Type())))
			insts = append(insts, wat.NewInstConvert_i32_wrap_i64())
		} else if x.Type().Size() == 8 && y.Type().Size() <= 4 {
			insts = append(insts, x.EmitPush()...)
			insts = append(insts, y.EmitPush()...)
			insts = append(insts, wat.NewInstConvert_i64_extend_i32_u())
			insts = append(insts, wat.NewInstShr(toWatType(ret_type)))
		} else if (x.Type().Size() <= 4 && y.Type().Size() <= 4) || (x.Type().Size() == 8 && y.Type().Size() == 8) {
			insts = append(insts, x.EmitPush()...)
			insts = append(insts, y.EmitPush()...)
			insts = append(insts, wat.NewInstShr(toWatType(ret_type)))
		} else {
			logger.Fatal("Unreachable")
		}

	case wat.OpCodeAndNot:
		ret_type = x.Type()
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, y.EmitPush()...)
		insts = append(insts, NewConst("-1", y.Type()).EmitPush()...)
		insts = append(insts, wat.NewInstXor(toWatType(ret_type)))
		insts = append(insts, wat.NewInstAnd(toWatType(ret_type)))

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
	var field *StructField
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
			insts = append(insts, NewConst(strconv.Itoa(typ.Base.Size()), m.I32).EmitPush()...)
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
			insts = append(insts, NewConst(strconv.Itoa(typ.Base.Size()), m.I32).EmitPush()...)
			insts = append(insts, id.EmitPush()...)
			insts = append(insts, wat.NewInstMul(wat.I32{}))
			insts = append(insts, wat.NewInstAdd(wat.I32{}))
			ret_type = m.GenValueType_Ref(typ.Base)

		default:
			logger.Fatalf("Todo: %T", typ)
		}

	case *aSlice:
		base_type := x.Type().(*Slice).Base
		insts = append(insts, x.Extract("b").EmitPush()...)
		insts = append(insts, x.Extract("d").EmitPush()...)
		insts = append(insts, NewConst(strconv.Itoa(base_type.Size()), m.I32).EmitPush()...)
		insts = append(insts, id.EmitPush()...)
		insts = append(insts, wat.NewInstMul(wat.I32{}))
		insts = append(insts, wat.NewInstAdd(wat.I32{}))
		ret_type = m.GenValueType_Ref(base_type)

	default:
		logger.Fatalf("Todo: %T", x)
	}

	return
}

func (m *Module) EmitGenIndex(x, id Value) (insts []wat.Inst, ret_type ValueType) {
	if !id.Type().Equal(m.I32) {
		panic("index should be i32")
	}

	switch x := x.(type) {
	case *aArray:
		ret_type = x.Type().(*Array).Base
		insts = append(insts, x.emitIndexOf(id)...)

	default:
		logger.Fatalf("Todo: %T", x)
	}

	return
}

func (m *Module) EmitGenSlice(x, low, high, max Value) (insts []wat.Inst, ret_type ValueType) {
	switch x := x.(type) {
	case *aSlice:
		insts = x.emitSub(low, high, max)
		ret_type = x.Type()

	case *aString:
		insts = x.emitSub(low, high)
		ret_type = x.Type()

	case *aRef:
		switch btype := x.Type().(*Ref).Base.(type) {
		case *Slice:
			slt := m.GenValueType_Slice(btype.Base)
			insts = slt.emitGenFromRefOfSlice(x, low, high, max)
			ret_type = slt

		case *Array:
			slt := m.GenValueType_Slice(btype.Base)
			insts = slt.emitGenFromRefOfArray(x, low, high, max)
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
	slice_type := m.GenValueType_Slice(base_type)
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
	for {
		if u, ok := x.(*aDup); ok {
			x = u.underlying
		} else {
			break
		}
	}

	for {
		if ut, ok := typ.(*Dup); ok {
			typ = ut.Base
		} else {
			break
		}
	}

	if x.Type().Equal(typ) {
		insts = append(insts, x.EmitPush()...)
		return
	}

	xt := x.Type()

	switch {
	/*case typ.Equal(m.I8):
	insts = append(insts, x.EmitPush()...)
	switch {
	case xt.Equal(m.I8), xt.Equal(m.U8), xt.Equal(m.I16), xt.Equal(m.U16), xt.Equal(m.I32), xt.Equal(m.U32), xt.Equal(m.RUNE):
		break

	case xt.Equal(m.I64), xt.Equal(m.U64):
		insts = append(insts, wat.NewInstConvert_i32_wrap_i64())

	case xt.Equal(m.F32):
		insts = append(insts, wat.NewInstConvert_i32_trunc_f32_s())

	case xt.Equal(m.F64):
		insts = append(insts, wat.NewInstConvert_i32_trunc_f64_s())
	}
	return  */

	case typ.Equal(m.U8):
		insts = append(insts, x.EmitPush()...)
		switch {
		case xt.Equal(m.U8), xt.Equal(m.U16), xt.Equal(m.I32), xt.Equal(m.U32), xt.Equal(m.RUNE): //Todo: xt.Equal(m.I8), xt.Equal(m.I16)
			break

		case xt.Equal(m.I64), xt.Equal(m.U64):
			insts = append(insts, wat.NewInstConvert_i32_wrap_i64())

		case xt.Equal(m.F32):
			insts = append(insts, wat.NewInstConvert_i32_trunc_f32_s())

		case xt.Equal(m.F64):
			insts = append(insts, wat.NewInstConvert_i32_trunc_f64_s())
		}
		insts = append(insts, wat.NewInstConst(wat.I32{}, "255"))
		insts = append(insts, wat.NewInstAnd(wat.I32{}))
		return

	/*case typ.Equal(m.I16):
	insts = append(insts, x.EmitPush()...)
	switch {
	case xt.Equal(m.I8), xt.Equal(m.U8), xt.Equal(m.I16), xt.Equal(m.U16), xt.Equal(m.I32), xt.Equal(m.U32), xt.Equal(m.RUNE):
		break

	case xt.Equal(m.I64), xt.Equal(m.U64):
		insts = append(insts, wat.NewInstConvert_i32_wrap_i64())

	case xt.Equal(m.F32):
		insts = append(insts, wat.NewInstConvert_i32_trunc_f32_s())

	case xt.Equal(m.F64):
		insts = append(insts, wat.NewInstConvert_i32_trunc_f64_s())
	}
	return  */

	case typ.Equal(m.U16):
		insts = append(insts, x.EmitPush()...)
		switch {
		case xt.Equal(m.U8), xt.Equal(m.U16), xt.Equal(m.I32), xt.Equal(m.U32), xt.Equal(m.RUNE): //Todo:xt.Equal(m.I8), xt.Equal(m.I16)
			break

		case xt.Equal(m.I64), xt.Equal(m.U64):
			insts = append(insts, wat.NewInstConvert_i32_wrap_i64())

		case xt.Equal(m.F32):
			insts = append(insts, wat.NewInstConvert_i32_trunc_f32_s())

		case xt.Equal(m.F64):
			insts = append(insts, wat.NewInstConvert_i32_trunc_f64_s())
		}
		insts = append(insts, wat.NewInstConst(wat.I32{}, "65535"))
		insts = append(insts, wat.NewInstAnd(wat.I32{}))
		return

	case typ.Equal(m.I32), typ.Equal(m.U32), typ.Equal(m.RUNE):
		insts = append(insts, x.EmitPush()...)
		switch {
		case xt.Equal(m.U8), xt.Equal(m.U16), xt.Equal(m.I32), xt.Equal(m.U32), xt.Equal(m.RUNE): //Todo:xt.Equal(m.I8), xt.Equal(m.I16),
			break

		case xt.Equal(m.I64), xt.Equal(m.U64):
			insts = append(insts, wat.NewInstConvert_i32_wrap_i64())

		case xt.Equal(m.F32):
			insts = append(insts, wat.NewInstConvert_i32_trunc_f32_s())

		case xt.Equal(m.F64):
			insts = append(insts, wat.NewInstConvert_i32_trunc_f64_s())
		}
		return

	case typ.Equal(m.I64):
		insts = append(insts, x.EmitPush()...)
		switch {
		case xt.Equal(m.I32): //Todo: xt.Equal(m.I8), xt.Equal(m.I16)
			insts = append(insts, wat.NewInstConvert_i64_extend_i32_s())

		case xt.Equal(m.U8), xt.Equal(m.U16), xt.Equal(m.U32), xt.Equal(m.RUNE):
			insts = append(insts, wat.NewInstConvert_i64_extend_i32_u())

		case xt.Equal(m.I64), xt.Equal(m.U64):
			break

		case xt.Equal(m.F32):
			insts = append(insts, wat.NewInstConvert_i64_trunc_f32_s())

		case xt.Equal(m.F64):
			insts = append(insts, wat.NewInstConvert_i64_trunc_f64_s())
		}
		return

	case typ.Equal(m.U64):
		insts = append(insts, x.EmitPush()...)
		switch {
		case xt.Equal(m.I32): //Todo: xt.Equal(m.I8), xt.Equal(m.I16)
			insts = append(insts, wat.NewInstConvert_i64_extend_i32_u())

		case xt.Equal(m.U8), xt.Equal(m.U16), xt.Equal(m.U32), xt.Equal(m.RUNE):
			insts = append(insts, wat.NewInstConvert_i64_extend_i32_u())

		case xt.Equal(m.I64), xt.Equal(m.U64):
			break

		case xt.Equal(m.F32):
			insts = append(insts, wat.NewInstConvert_i64_trunc_f32_s())

		case xt.Equal(m.F64):
			insts = append(insts, wat.NewInstConvert_i64_trunc_f64_s())
		}
		return

	case typ.Equal(m.F32):
		insts = append(insts, x.EmitPush()...)
		switch {
		case xt.Equal(m.I32): //Todo: xt.Equal(m.I8), xt.Equal(m.I16)
			insts = append(insts, wat.NewInstConvert_f32_convert_i32_s())

		case xt.Equal(m.U8), xt.Equal(m.U16), xt.Equal(m.U32):
			insts = append(insts, wat.NewInstConvert_f32_convert_i32_u())

		case xt.Equal(m.I64):
			insts = append(insts, wat.NewInstConvert_f32_convert_i64_s())

		case xt.Equal(m.U64):
			insts = append(insts, wat.NewInstConvert_f32_convert_i64_u())

		case xt.Equal(m.F64):
			insts = append(insts, wat.NewInstConvert_f32_demote_f64())
		}
		return

	case typ.Equal(m.F64):
		insts = append(insts, x.EmitPush()...)
		switch {
		case xt.Equal(m.I32): //Todo: xt.Equal(m.I8), xt.Equal(m.I16)
			insts = append(insts, wat.NewInstConvert_f64_convert_i32_s())

		case xt.Equal(m.U8), xt.Equal(m.U16), xt.Equal(m.U32):
			insts = append(insts, wat.NewInstConvert_f64_convert_i32_u())

		case xt.Equal(m.I64):
			insts = append(insts, wat.NewInstConvert_f64_convert_i64_s())

		case xt.Equal(m.U64):
			insts = append(insts, wat.NewInstConvert_f64_convert_i64_u())

		case xt.Equal(m.F32):
			insts = append(insts, wat.NewInstConvert_f64_promote_f32())
		}
		return

	case typ.Equal(m.STRING):
		switch {
		case xt.Equal(m.BYTES):
			x := x.(*aSlice)
			insts = append(insts, NewConst("", m.STRING).EmitPush()...)
			insts = append(insts, x.Extract("b").EmitPush()...)
			insts = append(insts, x.Extract("d").EmitPush()...)
			insts = append(insts, x.Extract("l").EmitPush()...)
			insts = append(insts, wat.NewInstCall(m.STRING.(*String).genFunc_Append()))
		}
		return

	case typ.Equal(m.BYTES):
		switch {
		case xt.Equal(m.STRING):
			x := x.(*aString)
			insts = append(insts, NewConst("0", m.BYTES).EmitPush()...)
			insts = append(insts, x.Extract("b").EmitPush()...)
			insts = append(insts, x.Extract("d").EmitPush()...)
			insts = append(insts, x.Extract("l").EmitPush()...)
			insts = append(insts, x.Extract("l").EmitPush()...)
			insts = append(insts, wat.NewInstCall(m.BYTES.(*Slice).genAppendFunc()))
		}
		return
	}

	logger.Fatalf("Todo: %+v %+v", x, typ)
	return
}

func (m *Module) EmitGenAppend(x, y Value) (insts []wat.Inst, ret_type ValueType) {
	xtype := x.Type().(*Slice)
	insts = append(insts, x.EmitPush()...)
	insts = append(insts, y.EmitPush()...)

	if !x.Type().Equal(y.Type()) {
		if _, ok := y.Type().(*String); ok && xtype.Base.Equal(m.U8) {
			insts = append(insts, y.(*aString).Extract("l").EmitPush()...)
		} else {
			logger.Fatal("Type not match")
			return
		}
	}

	insts = append(insts, wat.NewInstCall(xtype.genAppendFunc()))
	ret_type = xtype

	return
}

func (m *Module) EmitGenLen(x Value) (insts []wat.Inst) {
	switch x := x.(type) {
	case *aArray:
		insts = NewConst(strconv.Itoa(x.Type().(*Array).Capacity), m.I32).EmitPush()

	case *aSlice:
		insts = x.Extract("l").EmitPush()

	case *aString:
		insts = x.Extract("l").EmitPush()

	default:
		logger.Fatalf("Todo: %T", x)
	}

	return
}

func (m *Module) EmitGenCap(x Value) (insts []wat.Inst) {
	switch x := x.(type) {
	case *aArray:
		insts = NewConst(strconv.Itoa(x.Type().(*Array).Capacity), m.I32).EmitPush()

	case *aSlice:
		insts = x.Extract("c").EmitPush()

	default:
		logger.Fatalf("Todo: %T", x)
	}

	return
}

func (m *Module) EmitGenCopy(x, y Value) (insts []wat.Inst) {
	xtype := x.Type().(*Slice)
	insts = append(insts, x.EmitPush()...)
	insts = append(insts, y.EmitPush()...)

	if !x.Type().Equal(y.Type()) {
		if _, ok := y.Type().(*String); ok && xtype.Base.Equal(m.U8) {
			insts = append(insts, y.(*aString).Extract("l").EmitPush()...)
		} else {
			logger.Fatal("Type not match")
			return
		}
	}

	insts = append(insts, wat.NewInstCall(xtype.genCopyFunc()))
	return
}

func (m *Module) EmitGenMakeInterface(x Value, itype ValueType) (insts []wat.Inst) {
	x_type := x.Type()
	m.markConcreteTypeUsed(x_type)
	m.markInterfaceUsed(itype)

	switch x := x.(type) {
	case *aRef:
		return itype.(*Interface).emitGenFromRef(x)

	default:
		compID := x_type.OnComp()
		if compID == 0 {
			var f Function
			f.InternalName = "$" + GenSymbolName(x_type.Name()) + ".$$compAddr"
			p0 := NewLocal("p0", m.GenValueType_Ptr(x_type))
			p1 := NewLocal("p1", m.GenValueType_Ptr(x_type))
			f.Params = append(f.Params, p0)
			f.Params = append(f.Params, p1)
			f.Results = append(f.Results, m.I32)

			v0 := NewLocal("v0", x_type)
			v1 := NewLocal("v1", x_type)
			f.Locals = append(f.Locals, v0)
			f.Locals = append(f.Locals, v1)

			//f.Insts = append(f.Insts, v0.EmitInit()...)
			//f.Insts = append(f.Insts, v1.EmitInit()...)

			f.Insts = append(f.Insts, x_type.EmitLoadFromAddr(p0, 0)...)
			f.Insts = append(f.Insts, v0.EmitPop()...)
			f.Insts = append(f.Insts, x_type.EmitLoadFromAddr(p1, 0)...)
			f.Insts = append(f.Insts, v1.EmitPop()...)

			if ins, ok := v0.emitEq(v1); ok {
				f.Insts = append(f.Insts, ins...)

				f.Insts = append(f.Insts, v0.EmitRelease()...)
				f.Insts = append(f.Insts, v1.EmitRelease()...)

				m.AddFunc(&f)
				compID = m.AddTableElem(f.InternalName)

			} else {
				compID = -1
			}

			x_type.setOnComp(compID)
		}

		ref_t := m.GenValueType_Ref(x_type)
		return itype.(*Interface).emitGenFromValue(x, ref_t, compID)
	}
}

func (m *Module) EmitGenChangeInterface(x Value, destType ValueType) (insts []wat.Inst) {
	m.markInterfaceUsed(x.Type())
	m.markInterfaceUsed(destType)
	return destType.(*Interface).emitGenFromInterface(x.(*aInterface))
}

func (m *Module) EmitGenTypeAssert(x Value, destType ValueType, commaOk bool) (insts []wat.Inst) {
	si := x.(*aInterface)
	m.markInterfaceUsed(x.Type())
	if di, ok := destType.(*Interface); ok {
		m.markInterfaceUsed(di)
		return si.emitQueryInterface(destType, commaOk)
	} else {
		m.markConcreteTypeUsed(destType)
		return si.emitGetData(destType, commaOk)
	}
}

func (m *Module) EmitInvoke(i Value, params []Value, mid int, typeName string) (insts []wat.Inst) {
	iface := i.(*aInterface)
	insts = append(insts, iface.Extract("d").EmitPush()...)

	for _, v := range params {
		insts = append(insts, v.EmitPush()...)
	}

	insts = append(insts, iface.Extract("itab").EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.I32{}, 8+mid*4, 4))

	insts = append(insts, wat.NewInstCallIndirect(typeName))
	return
}

func (m *Module) EmitPrintString(v Value) (insts []wat.Inst) {
	s := v.(*aString)

	insts = append(insts, s.Extract("d").EmitPush()...)
	insts = append(insts, s.Extract("l").EmitPush()...)
	insts = append(insts, wat.NewInstCall("$runtime.waPuts"))
	return
}

func (m *Module) EmitPrintInterface(v Value) (insts []wat.Inst) {
	i := v.(*aInterface)
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa('(')))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))

	insts = append(insts, i.Extract("d").(*aRef).Extract("d").EmitPush()...)
	insts = append(insts, wat.NewInstCall("$runtime.waPrintU32Ptr")) // data

	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(',')))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))

	insts = append(insts, i.Extract("itab").EmitPush()...)
	insts = append(insts, wat.NewInstCall("$runtime.waPrintU32Ptr")) // itab

	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(')')))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))
	return
}

func (m *Module) EmitStringValue(v Value) (insts []wat.Inst) {
	s := v.(*aString)
	insts = append(insts, s.Extract("d").EmitPush()...)
	insts = append(insts, s.Extract("l").EmitPush()...)
	return
}

func (m *Module) emitPrintValue(v Value) (insts []wat.Inst) {

	panic("Todo")
}
