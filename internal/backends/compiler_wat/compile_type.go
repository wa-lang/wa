// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_wat

import (
	"strconv"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/ssa"
	"wa-lang.org/wa/internal/types"
)

type typeLib struct {
	module    *wir.Module
	ssaProg   *ssa.Program
	typeTable map[string]*wrapType

	usedConcreteTypes []*wrapType
	usedInterfaces    []*wrapType
}

type wrapMethod struct {
	sig        wir.FnSig
	name       string
	fullFnName string
}

type wrapType struct {
	wirType wir.ValueType
	methods []wrapMethod
}

func newTypeLib(m *wir.Module, prog *ssa.Program) *typeLib {
	te := typeLib{module: m, ssaProg: prog}
	te.typeTable = make(map[string]*wrapType)
	return &te
}

func (tLib *typeLib) find(v types.Type) *wrapType {
	if v, ok := tLib.typeTable[v.String()]; ok {
		return v
	}

	return nil
}

func (tLib *typeLib) compile(from types.Type) wir.ValueType {
	if v, ok := tLib.typeTable[from.String()]; ok {
		return v.wirType
	}

	//s := from.String()
	//println(s)

	var newType wrapType
	tLib.typeTable[from.String()] = &newType
	uncommanFlag := false

	switch t := from.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool, types.UntypedBool, types.Int, types.UntypedInt:
			newType.wirType = tLib.module.I32

		case types.Int32:
			if t.Name() == "rune" {
				newType.wirType = tLib.module.RUNE
			} else {
				newType.wirType = tLib.module.I32
			}

		case types.Uint32, types.Uintptr:
			newType.wirType = tLib.module.U32

		case types.Int64:
			newType.wirType = tLib.module.I64

		case types.Uint64:
			newType.wirType = tLib.module.U64

		case types.Float32, types.UntypedFloat:
			newType.wirType = tLib.module.F32

		case types.Float64:
			newType.wirType = tLib.module.F64

		case types.Int8:
			newType.wirType = tLib.module.I8

		case types.Uint8:
			newType.wirType = tLib.module.U8

		case types.Int16:
			newType.wirType = tLib.module.I16

		case types.Uint16:
			newType.wirType = tLib.module.U16

		case types.String:
			newType.wirType = tLib.module.STRING

		default:
			logger.Fatalf("Unknown type:%s", t)
			return nil
		}

	case *types.Tuple:
		switch t.Len() {
		case 0:
			newType.wirType = tLib.module.VOID

		case 1:
			newType.wirType = tLib.compile(t.At(0).Type())

		default:
			var feilds []wir.ValueType
			for i := 0; i < t.Len(); i++ {
				feilds = append(feilds, tLib.compile(t.At(i).Type()))
			}
			newType.wirType = tLib.module.GenValueType_Tuple(feilds)
		}

	case *types.Pointer:
		newType.wirType = tLib.module.GenValueType_Ref(tLib.compile(t.Elem()))
		uncommanFlag = true

	case *types.Named:
		switch ut := t.Underlying().(type) {
		case *types.Struct:
			//Todo: 待解决类型嵌套包含的问题
			var fs []wir.Field
			for i := 0; i < ut.NumFields(); i++ {
				f := ut.Field(i)
				wtyp := tLib.compile(f.Type())
				if f.Embedded() {
					fs = append(fs, wir.NewField("$"+wtyp.Name(), wtyp))
				} else {
					fs = append(fs, wir.NewField(wir.GenSymbolName(f.Name()), wtyp))
				}
			}
			pkg_name, _ := wir.GetPkgMangleName(t.Obj().Pkg().Path())
			obj_name := wir.GenSymbolName(t.Obj().Name())
			tStruct := tLib.module.GenValueType_Struct(pkg_name+"."+obj_name, fs)

			newType.wirType = tStruct
			uncommanFlag = true

		case *types.Interface:
			if ut.NumMethods() == 0 {
				return tLib.compile(ut)
			}
			pkg_name, _ := wir.GetPkgMangleName(t.Obj().Pkg().Path())
			obj_name := wir.GenSymbolName(t.Obj().Name())

			for i := 0; i < ut.NumMethods(); i++ {
				var method wrapMethod
				method.sig = tLib.GenFnSig(ut.Method(i).Type().(*types.Signature))
				method.name = ut.Method(i).Name()
				method.fullFnName = pkg_name + "." + obj_name + "." + method.name
				newType.methods = append(newType.methods, method)

				var fntype wir.FnType
				fntype.FnSig.Params = append(fntype.FnSig.Params, tLib.module.GenValueType_Ref(tLib.module.VOID))
				fntype.FnSig.Params = append(fntype.FnSig.Params, method.sig.Params...)
				fntype.FnSig.Results = method.sig.Results
				fntype.Name = method.fullFnName
				tLib.module.AddFnType(&fntype)
			}
			newType.wirType = tLib.module.GenValueType_Interface(pkg_name + "." + obj_name)

		case *types.Signature:
			newType.wirType = tLib.module.GenValueType_Closure(tLib.GenFnSig(ut))

		default:
			logger.Fatalf("Todo:%T", ut)
		}

	case *types.Array:
		newType.wirType = tLib.module.GenValueType_Array(tLib.compile(t.Elem()), int(t.Len()))

	case *types.Slice:
		newType.wirType = tLib.module.GenValueType_Slice(tLib.compile(t.Elem()))

	case *types.Signature:
		newType.wirType = tLib.module.GenValueType_Closure(tLib.GenFnSig(t))

	case *types.Interface:
		if t.NumMethods() != 0 {
			panic("NumMethods of interface{} != 0")
		}
		newType.wirType = tLib.module.GenValueType_Interface("interface")

	default:
		logger.Fatalf("Todo:%T", t)
	}

	if uncommanFlag {
		methodset := tLib.ssaProg.MethodSets.MethodSet(from)
		for i := 0; i < methodset.Len(); i++ {
			sel := methodset.At(i)
			mfn := tLib.ssaProg.MethodValue(sel)

			CompileFunc(mfn, tLib, tLib.module)

			var method wrapMethod
			method.name = mfn.Name()
			method.sig = tLib.GenFnSig(mfn.Signature)
			method.fullFnName, _ = wir.GetFnMangleName(mfn)

			newType.methods = append(newType.methods, method)
		}
	}

	return newType.wirType
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

func (tLib *typeLib) markConcreteTypeUsed(t types.Type) {
	v, ok := tLib.typeTable[t.String()]
	if !ok {
		logger.Fatal("Type not found.")
	}

	if v.wirType.Hash() != 0 {
		return
	}

	tLib.usedConcreteTypes = append(tLib.usedConcreteTypes, v)
	v.wirType.SetHash(len(tLib.usedConcreteTypes))
}

func (tLib *typeLib) markInterfaceUsed(t types.Type) {
	v, ok := tLib.typeTable[t.String()]
	if !ok {
		logger.Fatal("Type not found.")
	}

	if v.wirType.Hash() != 0 {
		return
	}

	tLib.usedInterfaces = append(tLib.usedInterfaces, v)
	v.wirType.SetHash(-len(tLib.usedInterfaces))
}

func (tLib *typeLib) buildItab() {
	var itabs []byte
	t_itab := tLib.typeTable["runtime._itab"].wirType

	for _, conrete := range tLib.usedConcreteTypes {
		for _, iface := range tLib.usedInterfaces {
			fits := true

			vtable := make([]int, len(iface.methods))

			for mid, method := range iface.methods {
				found := false
				for _, d := range conrete.methods {
					if d.name == method.name && d.sig.Equal(&method.sig) {
						found = true
						vtable[mid] = tLib.module.AddTableElem(d.fullFnName)
						break
					}
				}

				if !found {
					fits = false
					break
				}
			}

			var addr int
			if fits {
				var itab []byte
				header := wir.NewConst("0", t_itab)
				itab = append(itab, header.Bin()...)
				for _, v := range vtable {
					fnid := wir.NewConst(strconv.Itoa(v), tLib.module.U32)
					itab = append(itab, fnid.Bin()...)
				}

				addr = tLib.module.DataSeg.Append(itab, 8)
			}

			itabs = append(itabs, wir.NewConst(strconv.Itoa(addr), tLib.module.U32).Bin()...)
		}
	}

	itabs_ptr := tLib.module.DataSeg.Append(itabs, 8)
	tLib.module.SetGlobalInitValue("$wa.RT._itabsPtr", strconv.Itoa(itabs_ptr))
	tLib.module.SetGlobalInitValue("$wa.RT._interfaceCount", strconv.Itoa(len(tLib.usedInterfaces)))
	tLib.module.SetGlobalInitValue("$wa.RT._concretTypeCount", strconv.Itoa(len(tLib.usedConcreteTypes)))
}

func (p *Compiler) compileType(t *ssa.Type) {
	p.tLib.compile(t.Type())
}
