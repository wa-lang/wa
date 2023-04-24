// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_wat

import (
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
	if v, ok := tLib.typeTable[from.String()]; ok {
		return *v
	}

	var newType wir.ValueType
	tLib.typeTable[from.String()] = &newType
	uncommanFlag := false

	switch t := from.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool, types.UntypedBool, types.Int, types.UntypedInt:
			newType = tLib.module.I32

		case types.Int32:
			if t.Name() == "rune" {
				newType = tLib.module.RUNE
			} else {
				newType = tLib.module.I32
			}

		case types.Uint32, types.Uintptr:
			newType = tLib.module.U32

		case types.Int64:
			newType = tLib.module.I64

		case types.Uint64:
			newType = tLib.module.U64

		case types.Float32, types.UntypedFloat:
			newType = tLib.module.F32

		case types.Float64:
			newType = tLib.module.F64

		case types.Int8:
			newType = tLib.module.I8

		case types.Uint8:
			newType = tLib.module.U8

		case types.Int16:
			newType = tLib.module.I16

		case types.Uint16:
			newType = tLib.module.U16

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
		newType = tLib.module.GenValueType_Ref(tLib.compile(t.Elem()))
		uncommanFlag = true

	case *types.Named:
		switch ut := t.Underlying().(type) {
		case *types.Struct:
			pkg_name, _ := wir.GetPkgMangleName(t.Obj().Pkg().Path())
			obj_name := wir.GenSymbolName(t.Obj().Name())
			tStruct, found := tLib.module.GenValueType_Struct(pkg_name + "." + obj_name)
			if !found {
				for i := 0; i < ut.NumFields(); i++ {
					sf := ut.Field(i)
					dtyp := tLib.compile(sf.Type())
					if sf.Embedded() {
						df := tLib.module.NewStructField("$"+dtyp.Name(), dtyp)
						tStruct.AppendField(df)
					} else {
						df := tLib.module.NewStructField(wir.GenSymbolName(sf.Name()), dtyp)
						tStruct.AppendField(df)
					}
				}
				tStruct.Finish()
			}

			newType = tStruct
			uncommanFlag = true

		case *types.Interface:
			if ut.NumMethods() == 0 {
				return tLib.compile(ut)
			}
			pkg_name, _ := wir.GetPkgMangleName(t.Obj().Pkg().Path())
			obj_name := wir.GenSymbolName(t.Obj().Name())

			newType = tLib.module.GenValueType_Interface(pkg_name + "." + obj_name)

			for i := 0; i < ut.NumMethods(); i++ {
				var method wir.Method
				method.Sig = tLib.GenFnSig(ut.Method(i).Type().(*types.Signature))
				method.Name = ut.Method(i).Name()
				method.FullFnName = pkg_name + "." + obj_name + "." + method.Name

				newType.AddMethod(method)

				var fntype wir.FnType
				fntype.FnSig.Params = append(fntype.FnSig.Params, tLib.module.GenValueType_Ref(tLib.module.VOID))
				fntype.FnSig.Params = append(fntype.FnSig.Params, method.Sig.Params...)
				fntype.FnSig.Results = method.Sig.Results
				fntype.Name = method.FullFnName
				tLib.module.AddFnType(&fntype)
			}

		case *types.Signature:
			newType = tLib.module.GenValueType_Closure(tLib.GenFnSig(ut))

		default:
			logger.Fatalf("Todo:%T", ut)
		}

	case *types.Array:
		newType = tLib.module.GenValueType_Array(tLib.compile(t.Elem()), int(t.Len()))

	case *types.Slice:
		newType = tLib.module.GenValueType_Slice(tLib.compile(t.Elem()))

	case *types.Signature:
		newType = tLib.module.GenValueType_Closure(tLib.GenFnSig(t))

	case *types.Interface:
		if t.NumMethods() != 0 {
			panic("NumMethods of interface{} != 0")
		}
		newType = tLib.module.GenValueType_Interface("interface")

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
			method.FullFnName, _ = wir.GetFnMangleName(mfn)

			newType.AddMethod(method)
		}
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
