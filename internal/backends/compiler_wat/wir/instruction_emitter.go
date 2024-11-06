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
		if !lh.Type().Equal(rh.Type()) && !(lh.Type().Equal(m.I32) && rh.Type().Equal(m.RUNE) || lh.Type().Equal(m.RUNE) && rh.Type().Equal(m.I32)) {
			logger.Fatal("x.Type:", lh.Type().Named(), ", y.Type():", rh.Type().Named())
		}

		insts = append(insts, rh.EmitPush()...)
		insts = append(insts, lh.EmitPop()...)
	}

	return insts
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
	switch op {
	case wat.OpCodeAdd:
		ret_type = x.Type()
		insts = append(insts, x.EmitPushNoRetain()...)
		insts = append(insts, y.EmitPushNoRetain()...)

		if ret_type.Equal(m.STRING) {
			insts = append(insts, wat.NewInstCall(m.STRING.(*String).fnName_append))
		} else if ret_type.Equal(m.COMPLEX64) {
			insts = append(insts, m.COMPLEX64.(*Complex64).emitAdd()...)
		} else if ret_type.Equal(m.COMPLEX128) {
			insts = append(insts, m.COMPLEX128.(*Complex128).emitAdd()...)
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
		insts = append(insts, x.EmitPushNoRetain()...)
		insts = append(insts, y.EmitPushNoRetain()...)

		if ret_type.Equal(m.COMPLEX64) {
			insts = append(insts, m.COMPLEX64.(*Complex64).emitSub()...)
		} else if ret_type.Equal(m.COMPLEX128) {
			insts = append(insts, m.COMPLEX128.(*Complex128).emitSub()...)
		} else {
			insts = append(insts, wat.NewInstSub(toWatType(ret_type)))
		}

		if ret_type.Equal(m.U8) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "255"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		} else if ret_type.Equal(m.U16) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "65535"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		}

	case wat.OpCodeMul:
		ret_type = x.Type()
		insts = append(insts, x.EmitPushNoRetain()...)
		insts = append(insts, y.EmitPushNoRetain()...)

		if ret_type.Equal(m.COMPLEX64) {
			insts = append(insts, m.COMPLEX64.(*Complex64).emitMul()...)
		} else if ret_type.Equal(m.COMPLEX128) {
			insts = append(insts, m.COMPLEX128.(*Complex128).emitMul()...)
		} else {
			insts = append(insts, wat.NewInstMul(toWatType(ret_type)))
		}

		if ret_type.Equal(m.U8) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "255"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		} else if ret_type.Equal(m.U16) {
			insts = append(insts, wat.NewInstConst(wat.I32{}, "65535"))
			insts = append(insts, wat.NewInstAnd(wat.I32{}))
		}

	case wat.OpCodeQuo:
		ret_type = x.Type()
		insts = append(insts, x.EmitPushNoRetain()...)
		insts = append(insts, y.EmitPushNoRetain()...)

		if ret_type.Equal(m.COMPLEX64) {
			insts = append(insts, m.COMPLEX64.(*Complex64).emitDiv()...)
		} else if ret_type.Equal(m.COMPLEX128) {
			insts = append(insts, m.COMPLEX128.(*Complex128).emitDiv()...)
		} else {
			insts = append(insts, wat.NewInstDiv(toWatType(ret_type)))
		}

	case wat.OpCodeRem:
		ret_type = x.Type()
		insts = append(insts, x.EmitPushNoRetain()...)
		insts = append(insts, y.EmitPushNoRetain()...)
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
		insts = append(insts, x.EmitPushNoRetain()...)
		insts = append(insts, y.EmitPushNoRetain()...)

		if x.Type().Equal(m.STRING) {
			insts = append(insts, wat.NewInstCall("$wa.runtime.string_LSS"))
		} else {
			insts = append(insts, wat.NewInstLt(toWatType(x.Type())))
		}

		ret_type = m.BOOL

	case wat.OpCodeGt:
		insts = append(insts, x.EmitPushNoRetain()...)
		insts = append(insts, y.EmitPushNoRetain()...)

		if x.Type().Equal(m.STRING) {
			insts = append(insts, wat.NewInstCall("$wa.runtime.string_GTR"))
		} else {
			insts = append(insts, wat.NewInstGt(toWatType(x.Type())))
		}

		ret_type = m.BOOL

	case wat.OpCodeLe:
		insts = append(insts, x.EmitPushNoRetain()...)
		insts = append(insts, y.EmitPushNoRetain()...)

		if x.Type().Equal(m.STRING) {
			insts = append(insts, wat.NewInstCall("$wa.runtime.string_LEQ"))
		} else {
			insts = append(insts, wat.NewInstLe(toWatType(x.Type())))
		}
		ret_type = m.BOOL

	case wat.OpCodeGe:
		insts = append(insts, x.EmitPushNoRetain()...)
		insts = append(insts, y.EmitPushNoRetain()...)

		if x.Type().Equal(m.STRING) {
			insts = append(insts, wat.NewInstCall("$wa.runtime.string_GEQ"))
		} else {
			insts = append(insts, wat.NewInstGe(toWatType(x.Type())))
		}
		ret_type = m.BOOL

	case wat.OpCodeComp:
		if x.Type().Equal(m.STRING) {
			insts = append(insts, x.EmitPushNoRetain()...)
			insts = append(insts, y.EmitPushNoRetain()...)
			insts = append(insts, wat.NewInstCall("$wa.runtime.string_Comp"))
		} else {
			insts = append(insts, x.EmitPushNoRetain()...)
			insts = append(insts, y.EmitPushNoRetain()...)
			insts = append(insts, wat.NewInstLt(toWatType(x.Type())))

			inst_lt := wat.NewInstIf(nil, nil, nil)
			inst_lt.Ret = append(inst_lt.Ret, wat.I32{})
			inst_lt.True = append(inst_lt.True, wat.NewInstConst(wat.I32{}, "-1"))
			inst_lt.False = append(inst_lt.False, x.EmitPushNoRetain()...)
			inst_lt.False = append(inst_lt.False, y.EmitPushNoRetain()...)
			inst_lt.False = append(inst_lt.False, wat.NewInstGt(toWatType(x.Type())))

			insts = append(insts, inst_lt)
		}

		ret_type = m.I32

	case wat.OpCodeAnd:
		ret_type = x.Type()
		insts = append(insts, x.EmitPushNoRetain()...)
		insts = append(insts, y.EmitPushNoRetain()...)
		insts = append(insts, wat.NewInstAnd(toWatType(ret_type)))

	case wat.OpCodeOr:
		ret_type = x.Type()
		insts = append(insts, x.EmitPushNoRetain()...)
		insts = append(insts, y.EmitPushNoRetain()...)
		insts = append(insts, wat.NewInstOr(toWatType(ret_type)))

	case wat.OpCodeXor:
		ret_type = x.Type()
		insts = append(insts, x.EmitPushNoRetain()...)
		insts = append(insts, y.EmitPushNoRetain()...)
		insts = append(insts, wat.NewInstXor(toWatType(ret_type)))

	case wat.OpCodeShl:
		ret_type = x.Type()

		if x.Type().Size() <= 4 && y.Type().Size() == 8 {
			insts = append(insts, x.EmitPushNoRetain()...)
			insts = append(insts, y.EmitPushNoRetain()...)
			insts = append(insts, wat.NewInstConvert_i32_wrap_i64())
			insts = append(insts, wat.NewInstShl(toWatType(ret_type)))
		} else if x.Type().Size() == 8 && y.Type().Size() <= 4 {
			insts = append(insts, x.EmitPushNoRetain()...)
			insts = append(insts, y.EmitPushNoRetain()...)
			insts = append(insts, wat.NewInstConvert_i64_extend_i32_u())
			insts = append(insts, wat.NewInstShl(toWatType(ret_type)))
		} else if (x.Type().Size() <= 4 && y.Type().Size() <= 4) || (x.Type().Size() == 8 && y.Type().Size() == 8) {
			insts = append(insts, x.EmitPushNoRetain()...)
			insts = append(insts, y.EmitPushNoRetain()...)
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
			insts = append(insts, x.EmitPushNoRetain()...)
			insts = append(insts, y.EmitPushNoRetain()...)
			insts = append(insts, wat.NewInstConvert_i32_wrap_i64())
			insts = append(insts, wat.NewInstShr(toWatType(ret_type)))
		} else if x.Type().Size() == 8 && y.Type().Size() <= 4 {
			insts = append(insts, x.EmitPushNoRetain()...)
			insts = append(insts, y.EmitPushNoRetain()...)
			insts = append(insts, wat.NewInstConvert_i64_extend_i32_u())
			insts = append(insts, wat.NewInstShr(toWatType(ret_type)))
		} else if (x.Type().Size() <= 4 && y.Type().Size() <= 4) || (x.Type().Size() == 8 && y.Type().Size() == 8) {
			insts = append(insts, x.EmitPushNoRetain()...)
			insts = append(insts, y.EmitPushNoRetain()...)
			insts = append(insts, wat.NewInstShr(toWatType(ret_type)))
		} else {
			logger.Fatal("Unreachable")
		}

	case wat.OpCodeAndNot:
		ret_type = x.Type()
		insts = append(insts, x.EmitPushNoRetain()...)
		insts = append(insts, y.EmitPushNoRetain()...)
		insts = append(insts, NewConst("-1", y.Type()).EmitPushNoRetain()...)
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
	case ValueKindGlobal:
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

func (m *Module) EmitStore(addr, value Value, is_init bool) (insts []wat.Inst) {
	switch addr.Kind() {
	case ValueKindGlobal:
		if value == nil {
			value = NewConst("0", addr.Type())
		}
		if !addr.Type().Equal(value.Type()) {
			logger.Fatal("Type not match")
			return nil
		}

		insts = append(insts, value.EmitPush()...)
		insts = append(insts, addr.EmitPop()...)

	default:
		switch addr := addr.(type) {
		case *aRef:
			if value == nil {
				value = NewConst("0", addr.Type().(*Ref).Base)
			}
			if addr.Kind() != ValueKindConst || value.Kind() != ValueKindConst || !is_init {
				insts = append(insts, addr.emitSetValue(value)...)
			} else {
				m.DataSeg.Set(value.Bin(), addr.getConstPtr())
			}

		case *aPtr:
			if value == nil {
				value = NewConst("0", addr.Type().(*Ptr).Base)
			}
			insts = append(insts, addr.emitSetValue(value)...)

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
	f := x.(*aTuple).ExtractByID(id)
	insts = append(insts, f.EmitPush()...)
	ret_type = f.Type()
	return
}

func (m *Module) EmitGenField(x Value, field_id int) (insts []wat.Inst, ret_type ValueType) {
	switch x := x.(type) {
	case *aStruct:
		field := x.ExtractByID(field_id)
		insts = append(insts, field.EmitPush()...)
		ret_type = field.Type()

	default:
		logger.Fatalf("Todo:%T", x)
	}

	return
}

func (m *Module) EmitGenFieldAddr(x Value, field_id int) (insts []wat.Inst, ret_type ValueType, ret_val Value) {
	switch x := x.(type) {
	case *aRef:
		field := x.Type().(*Ref).Base.(*Struct).fields[field_id]
		r_type := m.GenValueType_Ref(field.Type())
		ret_type = r_type
		if x.Kind() != ValueKindConst || x.ExtractByName("b").Name() != "0" {
			insts = append(insts, x.EmitPush()...)
			insts = append(insts, NewConst(strconv.Itoa(field._start), m.I32).EmitPush()...)
			insts = append(insts, wat.NewInstAdd(wat.I32{}))
		} else {
			ret_val = r_type.newConstRef(x.getConstPtr() + field._start)
		}
	case *aPtr:
		field := x.Type().(*Ptr).Base.(*Struct).fields[field_id]
		ret_type = m.GenValueType_Ptr(field.Type())
		insts = append(insts, x.EmitPush()...)
		insts = append(insts, NewConst(strconv.Itoa(field._start), m.I32).EmitPush()...)
		insts = append(insts, wat.NewInstAdd(wat.I32{}))

	default:
		logger.Fatalf("Todo:%T", x.Type())
	}

	return
}

func (m *Module) EmitGenIndexAddr(x, id Value) (insts []wat.Inst, ret_type ValueType, ret_val Value) {
	if !id.Type().Equal(m.I32) && !id.Type().Equal(m.U32) {
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
			r_type := m.GenValueType_Ref(typ.Base)
			ret_type = r_type
			if x.Kind() != ValueKindConst || id.Kind() != ValueKindConst || x.ExtractByName("b").Name() != "0" {
				insts = append(insts, x.EmitPush()...)
				insts = append(insts, NewConst(strconv.Itoa(typ.Base.Size()), m.I32).EmitPush()...)
				insts = append(insts, id.EmitPush()...)
				insts = append(insts, wat.NewInstMul(wat.I32{}))
				insts = append(insts, wat.NewInstAdd(wat.I32{}))
			} else {
				ptr := x.getConstPtr()
				i, _ := strconv.Atoi(id.Name())
				ptr += r_type.Base.Size() * i
				ret_val = r_type.newConstRef(ptr)
			}

		default:
			logger.Fatalf("Todo: %T", typ)
		}

	case *aSlice:
		base_type := x.Type().(*Slice).Base
		insts = append(insts, x.ExtractByName("b").EmitPush()...)
		insts = append(insts, x.ExtractByName("d").EmitPush()...)
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
		insts = append(insts, x.emitIndexOf(m, id)...)

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
			insts = btype.emitGenFromRefOfSlice(x, low, high, max)
			ret_type = btype

		case *Array:
			slt := m.GenValueType_Slice(btype.Base, "")
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

func (m *Module) EmitGenMakeSlice(slice_type ValueType, Len, Cap Value) (insts []wat.Inst) {
	insts = slice_type.(*Slice).emitGenMake(Len, Cap)
	return
}

func (m *Module) EmitGenMakeMap(map_type ValueType) (insts []wat.Inst) {
	return map_type.(*Map).emitGenMake()
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

	case *aMap:
		if CommaOk {
			fileds := []ValueType{x.typ.Elem, m.BOOL}
			ret_type = m.GenValueType_Tuple(fileds)
		} else {
			ret_type = x.typ.Elem
		}

		insts = x.emitLookup(index, CommaOk)

	default:
		logger.Fatalf("Todo: %T", x)
	}

	return
}

func (m *Module) EmitGenConvert(x Value, typ ValueType) (insts []wat.Inst) {
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
			insts = append(insts, NewConst("", m.STRING).EmitPushNoRetain()...)
			insts = append(insts, x.ExtractByName("b").EmitPushNoRetain()...)
			insts = append(insts, x.ExtractByName("d").EmitPushNoRetain()...)
			insts = append(insts, x.ExtractByName("l").EmitPushNoRetain()...)
			insts = append(insts, wat.NewInstCall(m.STRING.(*String).fnName_append))
			return

		case xt.Equal(m.GenValueType_Slice(m.RUNE, "")):
			insts = append(insts, x.EmitPushNoRetain()...)
			insts = append(insts, wat.NewInstCall("runtime.stringFromRuneSlice"))
			return

		case xt.Equal(m.RUNE) || xt.Equal(m.I32):
			insts = append(insts, x.EmitPushNoRetain()...)
			insts = append(insts, wat.NewInstCall("runtime.stringFromRune"))
			return

		case xt.Equal(m.U8):
			insts = append(insts, x.EmitPushNoRetain()...)
			insts = append(insts, wat.NewInstCall("runtime.stringFromRune"))
			return
		}

	case typ.Equal(m.BYTES):
		switch {
		case xt.Equal(m.STRING):
			x := x.(*aString)
			insts = append(insts, NewConst("0", m.BYTES).EmitPushNoRetain()...)
			insts = append(insts, x.ExtractByName("b").EmitPushNoRetain()...)
			insts = append(insts, x.ExtractByName("d").EmitPushNoRetain()...)
			insts = append(insts, x.ExtractByName("l").EmitPushNoRetain()...)
			insts = append(insts, x.ExtractByName("l").EmitPushNoRetain()...)
			insts = append(insts, wat.NewInstCall(m.BYTES.(*Slice).genAppendFunc(m)))
			return
		}

	case typ.Equal(m.GenValueType_Slice(m.RUNE, "")):
		switch {
		case xt.Equal(m.STRING):
			insts = append(insts, x.EmitPushNoRetain()...)
			insts = append(insts, wat.NewInstCall("runtime.runeSliceFromString"))
			return
		}
	}

	logger.Fatalf("Todo: x.type: %s, dest_type: %s", x.Type().Named(), typ.Named())
	return
}

func (m *Module) EmitGenAppend(x, y Value) (insts []wat.Inst, ret_type ValueType) {
	xtype := x.Type().(*Slice)
	insts = append(insts, x.EmitPushNoRetain()...)
	insts = append(insts, y.EmitPushNoRetain()...)

	if !x.Type().Equal(y.Type()) {
		if _, ok := y.Type().(*String); ok && xtype.Base.Equal(m.U8) {
			insts = append(insts, y.(*aString).ExtractByName("l").EmitPushNoRetain()...)
		} else {
			logger.Fatal("Type not match")
			return
		}
	}

	insts = append(insts, wat.NewInstCall(xtype.genAppendFunc(m)))
	ret_type = xtype

	return
}

func (m *Module) EmitGenLen(x Value) (insts []wat.Inst) {
	switch x := x.(type) {
	case *aArray:
		insts = NewConst(strconv.Itoa(x.Type().(*Array).Capacity), m.I32).EmitPush()

	case *aSlice:
		insts = x.ExtractByName("l").EmitPush()

	case *aString:
		insts = x.ExtractByName("l").EmitPush()

	case *aMap:
		insts = x.emitLen()

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
		insts = x.ExtractByName("c").EmitPush()

	default:
		logger.Fatalf("Todo: %T", x)
	}

	return
}

func (m *Module) EmitGenCopy(x, y Value) (insts []wat.Inst) {
	xtype := x.Type().(*Slice)
	insts = append(insts, x.EmitPushNoRetain()...)
	insts = append(insts, y.EmitPushNoRetain()...)

	if !x.Type().Equal(y.Type()) {
		if _, ok := y.Type().(*String); ok && xtype.Base.Equal(m.U8) {
			insts = append(insts, y.(*aString).ExtractByName("l").EmitPushNoRetain()...)
		} else {
			logger.Fatal("Type not match")
			return
		}
	}

	insts = append(insts, wat.NewInstCall(xtype.genCopyFunc(m)))
	return
}

func (m *Module) EmitGenRaw(x Value) (insts []wat.Inst) {
	s := x.(*aSlice)
	return s.emitConvertToBytes()
}

func (m *Module) EmitGenSetFinalizer(x Value, fn_id int) (insts []wat.Inst) {
	return x.(*aRef).emitGenSetFinalizer(fn_id)
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
			f.InternalName = "$" + GenSymbolName(x_type.Named()) + ".$$compAddr"
			p0 := NewLocal("p0", m.GenValueType_Ptr(x_type))
			p1 := NewLocal("p1", m.GenValueType_Ptr(x_type))
			f.Params = append(f.Params, p0)
			f.Params = append(f.Params, p1)
			f.Results = append(f.Results, m.I32)

			v0 := NewLocal("v0", x_type)
			v1 := NewLocal("v1", x_type)
			f.Locals = append(f.Locals, v0)
			f.Locals = append(f.Locals, v1)

			f.Insts = append(f.Insts, p0.EmitPushNoRetain()...)
			{
				p0_valid := wat.NewInstIf(nil, nil, nil)
				p0_valid.True = append(p0_valid.True, x_type.EmitLoadFromAddr(p0, 0)...)
				p0_valid.True = append(p0_valid.True, v0.EmitPop()...)
				f.Insts = append(f.Insts, p0_valid)
			}

			f.Insts = append(f.Insts, p1.EmitPushNoRetain()...)
			{
				p1_valid := wat.NewInstIf(nil, nil, nil)
				p1_valid.True = append(p1_valid.True, x_type.EmitLoadFromAddr(p1, 0)...)
				p1_valid.True = append(p1_valid.True, v1.EmitPop()...)
				f.Insts = append(f.Insts, p1_valid)
			}

			f.Insts = append(f.Insts, v0.emitCompare(v1)...)

			f.Insts = append(f.Insts, v0.EmitRelease()...)
			f.Insts = append(f.Insts, v1.EmitRelease()...)

			m.AddFunc(&f)
			compID = m.AddTableElem(f.InternalName)

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

func (m *Module) EmitGenRange(x Value) (insts []wat.Inst, ret_type ValueType) {
	switch x := x.(type) {
	case *aString:
		ret_type, _ = m.findValueType("runtime.stringIter")
		insts = append(insts, x.ExtractByName("d").EmitPush()...)
		insts = append(insts, x.ExtractByName("l").EmitPush()...)
		insts = append(insts, wat.NewInstConst(wat.I32{}, "0"))

	case *aMap:
		ret_type, _ = m.findValueType("runtime.mapIter")
		insts = append(insts, x.aRef.EmitPush()...)
		insts = append(insts, wat.NewInstConst(wat.I32{}, "0"))

	default:
		logger.Fatalf("Todo:%T", x)
	}

	return
}

func (m *Module) EmitGenNext_String(iter Value) (insts []wat.Inst, ret_type ValueType) {
	fields := []ValueType{m.BOOL, m.INT, m.RUNE}
	ret_type = m.GenValueType_Tuple(fields)

	insts = append(insts, iter.EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstCall("runtime.next_rune"))
	insts = append(insts, iter.(*aStruct).ExtractByName("pos").EmitPop()...)
	return
}

func (m *Module) EmitGenNext_Map(iter Value, ktype ValueType, vtype ValueType) (insts []wat.Inst, ret_type ValueType) {
	fields := []ValueType{m.BOOL, ktype, vtype}
	ret_type = m.GenValueType_Tuple(fields)

	insts = append(insts, iter.EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstCall("runtime.map."+ktype.Named()+"."+vtype.Named()+".next"))
	insts = append(insts, iter.(*aStruct).ExtractByName("pos").EmitPop()...)
	return
}

func (m *Module) EmitInvoke(i Value, params []Value, mid int, typeName string) (insts []wat.Inst) {
	iface := i.(*aInterface)
	insts = append(insts, iface.ExtractByName("d").EmitPushNoRetain()...)

	for _, v := range params {
		insts = append(insts, v.EmitPushNoRetain()...)
	}

	insts = append(insts, iface.ExtractByName("itab").EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstLoad(wat.I32{}, 8+mid*4, 4))

	insts = append(insts, wat.NewInstCallIndirect(typeName))
	return
}

func (m *Module) EmitGenMapUpdate(ma, k, v Value) (insts []wat.Inst) {
	mi := ma.(*aMap)
	return mi.emitUpdate(k, v)
}

func (m *Module) EmitGenDelete(ma, k Value) (insts []wat.Inst) {
	mi := ma.(*aMap)
	return mi.emitDelete(k)
}

func (m *Module) EmitPrintString(v Value) (insts []wat.Inst) {
	s := v.(*aString)

	insts = append(insts, s.ExtractByName("d").EmitPushNoRetain()...)
	insts = append(insts, s.ExtractByName("l").EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstCall("$runtime.waPuts"))
	return
}

func (m *Module) EmitPrintInterface(v Value) (insts []wat.Inst) {
	i := v.(*aInterface)
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa('(')))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))

	insts = append(insts, i.ExtractByName("d").(*aRef).ExtractByName("d").EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstCall("$runtime.waPrintU32Ptr")) // data

	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(',')))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))

	insts = append(insts, i.ExtractByName("itab").EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstCall("$runtime.waPrintU32Ptr")) // itab

	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(')')))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))
	return
}

func (m *Module) EmitPrintRef(v Value) (insts []wat.Inst) {
	i := v.(*aRef)
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa('(')))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))

	insts = append(insts, i.ExtractByName("b").EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstCall("$runtime.waPrintU32Ptr")) // block
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(',')))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))

	insts = append(insts, i.ExtractByName("b").EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstLoad(wat.I32{}, 0, 1))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintI32")) // refcount
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(',')))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))

	insts = append(insts, i.ExtractByName("d").EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstCall("$runtime.waPrintU32Ptr")) // data

	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(')')))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))
	return
}

func (m *Module) EmitPrintClosure(v Value) (insts []wat.Inst) {
	i := v.(*aClosure)
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa('(')))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))

	insts = append(insts, i.ExtractByName("fn_index").EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstCall("$runtime.waPrintI32")) // fn_index
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(',')))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))

	insts = append(insts, i.ExtractByName("d").(*aRef).ExtractByName("b").EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstCall("$runtime.waPrintU32Ptr")) // data.block
	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(',')))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))

	insts = append(insts, i.ExtractByName("d").(*aRef).ExtractByName("b").EmitPushNoRetain()...)
	insts = append(insts, wat.NewInstLoad(wat.I32{}, 0, 1))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintI32")) // data.block.RefCount

	insts = append(insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(')')))
	insts = append(insts, wat.NewInstCall("$runtime.waPrintRune"))
	return
}

func (m *Module) EmitStringValue(v Value) (insts []wat.Inst) {
	s := v.(*aString)
	insts = append(insts, s.ExtractByName("d").EmitPush()...)
	insts = append(insts, s.ExtractByName("l").EmitPush()...)
	return
}

func (m *Module) emitPrintValue(v Value) (insts []wat.Inst) {

	panic("Todo")
}
