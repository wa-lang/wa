// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_wat

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/ssa"
	"wa-lang.org/wa/internal/types"
)

type typeLib struct {
	prog      *loader.Program
	module    *wir.Module
	typeTable map[string]*wir.ValueType

	pendingMethods []*ssa.Function

	anonStructCount    int
	anonInterfaceCount int

	MaxTypeSize int
}

func newTypeLib(m *wir.Module, prog *loader.Program) *typeLib {
	te := typeLib{module: m, prog: prog}
	te.typeTable = make(map[string]*wir.ValueType)
	return &te
}

func (tLib *typeLib) find(v types.Type) wir.ValueType {
	if v, ok := tLib.typeTable[v.String()]; ok {
		return *v
	}

	return nil
}

func (tLib *typeLib) compile(from types.Type) wir.ValueType {
	from_name := from.String()
	if v, ok := tLib.typeTable[from_name]; ok {
		return *v
	}

	var newType wir.ValueType
	tLib.typeTable[from_name] = &newType
	uncommanFlag := false

	switch t := from.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool, types.UntypedBool:
			newType = tLib.module.BOOL

		case types.Uint8:
			newType = tLib.module.U8

		//case types.Int8:
		//	newType = tLib.module.I8

		case types.Uint16:
			newType = tLib.module.U16

		//case types.Int16:
		//	newType = tLib.module.I16

		case types.Uint32:
			newType = tLib.module.U32

		case types.Int32:
			if t.Name() == "rune" {
				newType = tLib.module.RUNE
			} else {
				newType = tLib.module.I32
			}

		case types.Uint64:
			newType = tLib.module.U64

		case types.Int64:
			newType = tLib.module.I64

		case types.Float32:
			newType = tLib.module.F32

		case types.Float64, types.UntypedFloat:
			newType = tLib.module.F64

		case types.Complex64:
			newType = tLib.module.COMPLEX64

		case types.Complex128:
			newType = tLib.module.COMPLEX128

		case types.Uint:
			newType = tLib.module.UINT

		case types.Int, types.UntypedInt:
			newType = tLib.module.INT

		case types.Uintptr:
			newType = tLib.module.UPTR

		case types.String:
			newType = tLib.module.STRING

		default:
			logger.Fatalf("Unknown type:%s", t)
			return nil
		}

	case *types.Tuple:
		switch t.Len() {
		case 0:
			newType = tLib.module.VOID

		case 1:
			newType = tLib.compile(t.At(0).Type())

		default:
			var feilds []wir.ValueType
			for i := 0; i < t.Len(); i++ {
				feilds = append(feilds, tLib.compile(t.At(i).Type()))
			}
			newType = tLib.module.GenValueType_Tuple(feilds)
		}

	case *types.Pointer:
		delete(tLib.typeTable, from_name)
		newType = tLib.module.GenValueType_Ref(tLib.compile(t.Elem()))
		tLib.typeTable[from_name] = &newType
		uncommanFlag = true

	case *types.Array:
		newType = tLib.module.GenValueType_Array(tLib.compile(t.Elem()), int(t.Len()), "")

	case *types.Slice:
		newType = tLib.module.GenValueType_Slice(tLib.compile(t.Elem()), "")

	case *types.Map:
		ei, ok := tLib.typeTable["interface{}"]
		if !ok {
			logger.Fatal("Can't find interface{} for map_imp.")
		}
		newType = tLib.module.GenValueType_Map(tLib.compile(t.Key()), tLib.compile(t.Elem()), "", *ei)

	case *types.Signature:
		newType = tLib.module.GenValueType_Closure(tLib.GenFnSig(t))

	case *types.Interface:
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
		}

	case *types.Struct:
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
		}
	case *types.Named:
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
		}

	default:
		logger.Fatalf("Todo:%T", t)
	}

	if uncommanFlag {
		methodset := tLib.prog.SSAProgram.MethodSets.MethodSet(from)
		for i := 0; i < methodset.Len(); i++ {
			sel := methodset.At(i)
			mfn := tLib.prog.SSAProgram.MethodValue(sel)

			tLib.pendingMethods = append(tLib.pendingMethods, mfn)

			var method wir.Method
			method.Sig = tLib.GenFnSig(mfn.Signature)
			method.Name = mfn.Name()
			method.FullFnName, _ = wir.GetFnMangleName(mfn, tLib.prog.Manifest.MainPkg)

			newType.AddMethod(method)
		}
	}

	tsize := newType.Size()
	if tsize > tLib.MaxTypeSize {
		tLib.MaxTypeSize = tsize
	}

	return newType
}

func (tl *typeLib) GenFnSig(s *types.Signature) wir.FnSig {
	var sig wir.FnSig
	for i := 0; i < s.Params().Len(); i++ {
		typ := tl.compile(s.Params().At(i).Type())
		sig.Params = append(sig.Params, typ)
	}
	for i := 0; i < s.Results().Len(); i++ {
		typ := tl.compile(s.Results().At(i).Type())
		sig.Results = append(sig.Results, typ)
	}
	return sig
}

func (tLib *typeLib) finish() {
	for len(tLib.pendingMethods) > 0 {
		mfn := tLib.pendingMethods[0]
		tLib.pendingMethods = tLib.pendingMethods[1:]
		CompileFunc(mfn, tLib.prog, tLib, tLib.module)
	}
}

func (p *Compiler) compileType(t *ssa.Type) {
	p.tLib.compile(t.Type())
}
