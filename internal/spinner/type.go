// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package spinner

import (
	"fmt"

	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/types"
	"wa-lang.org/wa/internal/wire"
)

/**************************************
本文用于将类型转为 wire.Type
**************************************/

//
func (b *Builder) BuildType(t types.Type) wire.Type {
	if t, ok := t.(*types.Basic); ok {
		switch t.Kind() {
		case types.Bool, types.UntypedBool:
			return b.module.Types.Bool

		case types.Uint8:
			return b.module.Types.U8

		case types.Uint16:
			return b.module.Types.U16

		case types.Uint32:
			return b.module.Types.U32

		case types.Uint64:
			return b.module.Types.U64

		case types.Int8:
			return b.module.Types.I8

		case types.Int16:
			return b.module.Types.I16

		case types.Int32:
			if t.Name() == "rune" {
				return b.module.Types.Rune
			} else {
				return b.module.Types.I32
			}

		case types.Int64:
			return b.module.Types.I64

		case types.Float32:
			return b.module.Types.F32

		case types.Float64, types.UntypedFloat:
			return b.module.Types.F64

		case types.Complex64:
			return b.module.Types.Complex64

		case types.Complex128:
			return b.module.Types.Complex128

		case types.Uint:
			return b.module.Types.Uint

		case types.Int, types.UntypedInt:
			return b.module.Types.Int

		case types.Uintptr:
			return b.module.Types.Uint

		case types.String:
			return b.module.Types.String

		default:
			logger.Fatalf("Unknown type:%s", t)
			return nil
		}
	}

	name := t.String()
	if v, ok := b.typeTable[name]; ok {
		return *v
	}

	var wtype wire.Type
	b.typeTable[name] = &wtype

	switch t := t.(type) {
	case *types.Tuple:
		switch t.Len() {
		case 0:
			wtype = b.module.Types.Void

		case 1:
			wtype = b.BuildType(t.At(0).Type())

		default:
			panic("Todo")
			//var feilds []wir.Type
			//for i := 0; i < t.Len(); i++ {
			//	feilds = append(feilds, tLib.compile(t.At(i).Type()))
			//}
			//newType = tLib.module.GenValueType_Tuple(feilds)
		}

	case *types.Pointer:
		delete(b.typeTable, name)
		wtype = b.module.Types.GenRef(b.BuildType(t.Elem()))
		b.typeTable[name] = &wtype

		//panic("Todo")
		//delete(tLib.typeTable, from_name)
		//newType = tLib.module.GenValueType_Ref(tLib.compile(t.Elem()))
		//tLib.typeTable[from_name] = &newType
		//uncommanFlag = true

	case *types.Array:
		panic("Todo")
		//newType = tLib.module.GenValueType_Array(tLib.compile(t.Elem()), int(t.Len()), "")

	case *types.Slice:
		panic("Todo")
		//newType = tLib.module.GenValueType_Slice(tLib.compile(t.Elem()), "")

	case *types.Map:
		panic("Todo")
		//ei, ok := tLib.typeTable["interface{}"]
		//if !ok {
		//	logger.Fatal("Can't find interface{} for map_imp.")
		//}
		//newType = tLib.module.GenValueType_Map(tLib.compile(t.Key()), tLib.compile(t.Elem()), "", *ei)

	case *types.Signature:
		panic("Todo")
		//newType = tLib.module.GenValueType_Closure(tLib.GenFnSig(t))

	case *types.Interface:
		panic("Todo")
		/*
			newType, _ = tLib.module.GenValueType_Interface("i`" + strconv.Itoa(tLib.anonInterfaceCount) + "`")
			tLib.anonInterfaceCount++

			for i := 0; i < t.NumMethods(); i++ {
				var method wir.Method
				method.Sig = tLib.GenFnSig(t.Method(i).Type().(*types.Signature))

				method.Name = t.Method(i).Name()

				var fnSig wir.FnSig
				fnSig.Params = append(fnSig.Params, tLib.module.GenValueType_Ref(tLib.module.VOID))
				fnSig.Params = append(fnSig.Params, method.Sig.Params...)
				fnSig.Results = method.Sig.Results

				method.FullFnName = tLib.module.AddFnSig(&fnSig)

				newType.AddMethod(method)
			} //*/

	case *types.Struct:
		panic("Todo")
		/*
			tStruct, found := tLib.module.GenValueType_Struct("s`" + strconv.Itoa(tLib.anonStructCount) + "`")
			newType = tStruct
			tLib.anonStructCount++

			if !found {
				for i := 0; i < t.NumFields(); i++ {
					sf := t.Field(i)
					dtyp := tLib.compile(sf.Type())
					var fname string
					if sf.Embedded() {
						//fname = "$" + dtyp.Name()
						fname = "$" + sf.Name()
					} else {
						fname = sf.Name()
					}
					df := tLib.module.NewStructField(fname, dtyp)
					tStruct.AppendField(df)

				}
				tStruct.Finish()
			} else {
				panic("???")
			} //*/

	case *types.Named:
		panic("Todo")

		/*
			pkg_name := ""
			if t.Obj().Pkg() != nil {
				pkg_name, _ = wir.GetPkgMangleName(t.Obj().Pkg().Path())
			}
			obj_name := wir.GenSymbolName(t.Obj().Name())
			type_name := pkg_name + "." + obj_name

			switch ut := t.Underlying().(type) {
			case *types.Basic:
				switch ut.Kind() {
				case types.Bool, types.UntypedBool:
					newType = tLib.module.GenValueType_bool(type_name)

				case types.Uint8:
					newType = tLib.module.GenValueType_u8(type_name)

				//case types.Int8:
				//	newType = tLib.module.GenValueType_i8(type_name)

				case types.Uint16:
					newType = tLib.module.GenValueType_u16(type_name)

				//case types.Int16:
				//	newType = tLib.module.GenValueType_i16(type_name)

				case types.Uint32:
					newType = tLib.module.GenValueType_u32(type_name)

				case types.Int32:
					if ut.Name() == "rune" {
						newType = tLib.module.GenValueType_rune(type_name)
					} else {
						newType = tLib.module.GenValueType_i32(type_name)
					}

				case types.Uint64:
					newType = tLib.module.GenValueType_u64(type_name)

				case types.Int64:
					newType = tLib.module.GenValueType_i64(type_name)

				case types.Float32:
					newType = tLib.module.GenValueType_f32(type_name)

				case types.Float64, types.UntypedFloat:
					newType = tLib.module.GenValueType_f64(type_name)

				case types.Complex64:
					newType = tLib.module.GenValueType_complex64(type_name)

				case types.Complex128:
					newType = tLib.module.GenValueType_complex128(type_name)

				case types.Uint:
					newType = tLib.module.GenValueType_uint(type_name)

				case types.Int, types.UntypedInt:
					newType = tLib.module.GenValueType_int(type_name)

				case types.String:
					newType = tLib.module.GenValueType_string(type_name)

				default:
					logger.Fatalf("Unknown type:%s", t)
					return nil
				}

			case *types.Array:
				newType = tLib.module.GenValueType_Array(tLib.compile(ut.Elem()), int(ut.Len()), type_name)

			case *types.Slice:
				newType = tLib.module.GenValueType_Slice(tLib.compile(ut.Elem()), type_name)

			case *types.Map:
				ei, ok := tLib.typeTable["interface{}"]
				if !ok {
					logger.Fatal("Can't find interface{} for map_imp.")
				}
				newType = tLib.module.GenValueType_Map(tLib.compile(ut.Key()), tLib.compile(ut.Elem()), type_name, *ei)

			case *types.Struct:
				tStruct, found := tLib.module.GenValueType_Struct(type_name)
				newType = tStruct
				if !found {
					for i := 0; i < ut.NumFields(); i++ {
						sf := ut.Field(i)
						dtyp := tLib.compile(sf.Type())
						if sf.Embedded() {
							//df := tLib.module.NewStructField("$"+dtyp.Name(), dtyp)
							df := tLib.module.NewStructField("$"+wir.GenSymbolName(sf.Name()), dtyp)
							tStruct.AppendField(df)
						} else {
							df := tLib.module.NewStructField(wir.GenSymbolName(sf.Name()), dtyp)
							tStruct.AppendField(df)
						}
					}
					tStruct.Finish()
				}

			case *types.Interface:
				itype, found := tLib.module.GenValueType_Interface(type_name)
				newType = itype
				if !found {
					for i := 0; i < ut.NumMethods(); i++ {
						var method wir.Method
						method.Sig = tLib.GenFnSig(ut.Method(i).Type().(*types.Signature))

						method.Name = ut.Method(i).Name()

						var fnSig wir.FnSig
						fnSig.Params = append(fnSig.Params, tLib.module.GenValueType_Ref(tLib.module.VOID))
						fnSig.Params = append(fnSig.Params, method.Sig.Params...)
						fnSig.Results = method.Sig.Results

						method.FullFnName = tLib.module.AddFnSig(&fnSig)

						itype.AddMethod(method)
					}
				}

			case *types.Signature:
				newType = tLib.module.GenValueType_Closure(tLib.GenFnSig(ut))

			default:
				logger.Fatalf("Todo:%T", ut)
			} //*/

	default:
		panic(fmt.Sprintf("Todo:%T", t))

	}

	return wtype

	// 需要处理 Signature
}
