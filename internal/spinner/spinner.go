// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package spinner

import (
	"fmt"
	"strconv"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/constant"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/token"
	"wa-lang.org/wa/internal/types"
	"wa-lang.org/wa/internal/wire"
)

/**************************************
本包用于将 AST 转为 wire
**************************************/

//-------------------------------------

/**************************************
Func: 包级函数
**************************************/
type Func struct {
	object   *types.Func
	sig      *types.Signature
	body     ast.Stmt
	pos, end int

	info *types.Info
	pkg  *types.Package

	wireFunc *wire.Function
}

/**************************************
Global: 包级全局变量
**************************************/
type Global struct {
	object *types.Var
}

/**************************************
Type: 包级具名类型
**************************************/
type Type struct {
	object *types.TypeName
}

/**************************************
Package: 包
**************************************/
type Package struct {
	pkg   *types.Package
	info  *types.Info
	files []*ast.File

	funcs   map[types.Object]*Func
	globals map[types.Object]*Global
	types   map[types.Object]*Type
}

func (sp *Spinner) createPackage(pkg *types.Package, files []*ast.File, info *types.Info) *Package {
	p := &Package{
		pkg:     pkg,
		info:    info,
		files:   files,
		funcs:   make(map[types.Object]*Func),
		globals: make(map[types.Object]*Global),
		types:   make(map[types.Object]*Type),
	}

	isMainPkg := pkg.Path() == sp.manifest.MainPkg

	for _, file := range files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				switch decl.Tok {
				case token.VAR, token.GLOBAL:
					for _, spec := range decl.Specs {
						for _, id := range spec.(*ast.ValueSpec).Names {
							if isBlankIdent(id) {
								continue
							}

							obj := info.Defs[id]
							v, ok := obj.(*types.Var)
							if !ok {
								panic("unexpected object: " + obj.String())
							}

							p.globals[v] = &Global{
								object: v,
							}
						}
					}

				case token.TYPE:
					for _, spec := range decl.Specs {
						id := spec.(*ast.TypeSpec).Name
						if isBlankIdent(id) {
							continue
						}

						obj := info.Defs[id]
						t, ok := obj.(*types.TypeName)
						if !ok {
							panic("unexpected object: " + obj.String())
						}

						p.types[t] = &Type{
							object: t,
						}
					}
				}

			case *ast.FuncDecl:
				id := decl.Name
				if isBlankIdent(id) {
					continue
				}

				obj := info.Defs[id]
				f, ok := obj.(*types.Func)
				if !ok {
					panic("unexpected object: " + obj.String())
				}

				fn := &Func{
					object: f,
					sig:    f.Type().(*types.Signature),
					body:   decl.Body,
					pos:    int(decl.Pos()),
					end:    int(decl.End()),
					info:   info,
					pkg:    pkg,
				}

				if isMainPkg {
					if ast.IsExported(decl.Name.Name) || sp.manifest.W2Mode && decl.Name.Name == token.K_主控 || !sp.manifest.W2Mode && decl.Name.Name == token.K_main {
						sp.funcsTodo[f] = true
					}
				}
				p.funcs[f] = fn
			}
		}
	}

	if isMainPkg {
		sp.mainPkg = p
	}
	sp.packages[pkg] = p
	return p
}

func (sp *Spinner) build() {
	sp.Module.Init()

	sp.typeTable = make(map[string]*wire.Type)
	sp.typeTable[sp.Module.Types.Void.Name()] = &sp.Module.Types.Void
	sp.typeTable[sp.Module.Types.Bool.Name()] = &sp.Module.Types.Bool
	sp.typeTable[sp.Module.Types.U8.Name()] = &sp.Module.Types.U8
	sp.typeTable[sp.Module.Types.U16.Name()] = &sp.Module.Types.U16
	sp.typeTable[sp.Module.Types.U32.Name()] = &sp.Module.Types.U32
	sp.typeTable[sp.Module.Types.U64.Name()] = &sp.Module.Types.U64
	sp.typeTable[sp.Module.Types.Uint.Name()] = &sp.Module.Types.Uint
	sp.typeTable[sp.Module.Types.I8.Name()] = &sp.Module.Types.I8
	sp.typeTable[sp.Module.Types.I16.Name()] = &sp.Module.Types.I16
	sp.typeTable[sp.Module.Types.I32.Name()] = &sp.Module.Types.I32
	sp.typeTable[sp.Module.Types.I64.Name()] = &sp.Module.Types.I64
	sp.typeTable[sp.Module.Types.Int.Name()] = &sp.Module.Types.Int
	sp.typeTable[sp.Module.Types.F32.Name()] = &sp.Module.Types.F32
	sp.typeTable[sp.Module.Types.F64.Name()] = &sp.Module.Types.F64
	//sp.typeTable[b.module.Types.Complex64.Name()] = &sp.Module.Types.Complex64
	//sp.typeTable[b.module.Types.Complex128.Name()] = &sp.Module.Types.Complex128
	sp.typeTable[sp.Module.Types.Rune.Name()] = &sp.Module.Types.Rune
	sp.typeTable[sp.Module.Types.String.Name()] = &sp.Module.Types.String

	for len(sp.funcsTodo) > 0 {
		for fnObj := range sp.funcsTodo {
			fb := CreateFuncBuilder(sp)
			pkg := sp.packages[fnObj.Pkg()]
			f := pkg.funcs[fnObj]
			fb.build(f)

			delete(sp.funcsTodo, fnObj)
			break
		}
	}
}

func (sp *Spinner) buildType(t types.Type) wire.Type {
	if t, ok := t.(*types.Basic); ok {
		switch t.Kind() {
		case types.Bool, types.UntypedBool:
			return sp.Module.Types.Bool

		case types.Uint8:
			return sp.Module.Types.U8

		case types.Uint16:
			return sp.Module.Types.U16

		case types.Uint32:
			return sp.Module.Types.U32

		case types.Uint64:
			return sp.Module.Types.U64

		case types.Int8:
			return sp.Module.Types.I8

		case types.Int16:
			return sp.Module.Types.I16

		case types.Int32:
			if t.Name() == "rune" {
				return sp.Module.Types.Rune
			} else {
				return sp.Module.Types.I32
			}

		case types.Int64:
			return sp.Module.Types.I64

		case types.Float32:
			return sp.Module.Types.F32

		case types.Float64, types.UntypedFloat:
			return sp.Module.Types.F64

		//case types.Complex64:
		//	return sp.Module.Types.Complex64
		//
		//case types.Complex128:
		//	return sp.Module.Types.Complex128

		case types.Uint:
			return sp.Module.Types.Uint

		case types.Int, types.UntypedInt:
			return sp.Module.Types.Int

		case types.Uintptr:
			return sp.Module.Types.Uint

		case types.String:
			return sp.Module.Types.String

		default:
			logger.Fatalf("Unknown type:%s", t)
			return nil
		}
	}

	name := t.String()
	if v, ok := sp.typeTable[name]; ok {
		return *v
	}

	var wtype wire.Type
	sp.typeTable[name] = &wtype

	switch t := t.(type) {
	case *types.Tuple:
		switch t.Len() {
		case 0:
			wtype = sp.Module.Types.Void

		case 1:
			wtype = sp.buildType(t.At(0).Type())

		default:
			var members []wire.Type
			for i := 0; i < t.Len(); i++ {
				members = append(members, sp.buildType(t.At(i).Type()))
			}
			wtype = sp.Module.Types.GenTuple(members)
		}

	case *types.Pointer:
		delete(sp.typeTable, name)
		wtype = sp.Module.Types.GenRef(sp.buildType(t.Elem()))
		sp.typeTable[name] = &wtype

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
		members := make([]wire.StructMember, t.NumFields())
		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			members[i].Type = sp.buildType(field.Type())

			if field.Embedded() {
				members[i].Name = "$" + field.Name()
			} else {
				members[i].Name = field.Name()
			}
		}
		wtype = sp.Module.Types.GenStruct(members)

		//panic("Todo")
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
		name := getMangleName(t.Obj().Pkg(), t.Obj().Name())
		named_type := sp.Module.Types.GenNamed(name)
		named_type.SetUnderlying(sp.buildType(t.Underlying()))
		wtype = named_type

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

func (sp *Spinner) constval(val constant.Value, typ types.Type, pos int) wire.Expr {
	t := sp.buildType(typ)

	if t.Equal(sp.Module.Types.Bool) {
		if constant.BoolVal(val) {
			return sp.Module.NewConst("1", t, pos)
		} else {
			return sp.Module.NewConst("0", t, pos)
		}
	} else if t.Equal(sp.Module.Types.U8) || t.Equal(sp.Module.Types.U16) || t.Equal(sp.Module.Types.U32) || t.Equal(sp.Module.Types.U64) || t.Equal(sp.Module.Types.Uint) {
		val, _ := constant.Uint64Val(val)
		return sp.Module.NewConst(strconv.Itoa(int(val)), t, pos)
	} else if t.Equal(sp.Module.Types.F32) {
		val, _ := constant.Float64Val(val)
		return sp.Module.NewConst(strconv.FormatFloat(val, 'f', -1, 32), t, pos)
	} else if t.Equal(sp.Module.Types.F64) {
		val, _ := constant.Float64Val(val)
		return sp.Module.NewConst(strconv.FormatFloat(val, 'f', -1, 64), t, pos)
		//} else if t.Equal(b.module.Types.Complex64) || t.Equal(b.module.Types.Complex128) {
		//	re, _ := constant.Float64Val(constant.Real(val))
		//	im, _ := constant.Float64Val(constant.Imag(val))
		//	res := strconv.FormatFloat(re, 'f', -1, 64)
		//	ims := strconv.FormatFloat(im, 'f', -1, 64)
		//	return sp.Module.NewConst(res+" "+ims, t, pos)
	} else if t.Equal(sp.Module.Types.I8) || t.Equal(sp.Module.Types.I16) || t.Equal(sp.Module.Types.I32) || t.Equal(sp.Module.Types.I64) || t.Equal(sp.Module.Types.Int) || t.Equal(sp.Module.Types.Rune) {
		val, _ := constant.Int64Val(val)
		return sp.Module.NewConst(strconv.Itoa(int(val)), t, pos)
	} else if t.Equal(sp.Module.Types.String) {
		val := constant.StringVal(val)
		return sp.Module.NewConst(val, t, pos)
	}

	panic("Todo")
}

func getMangleName(pkg *types.Package, name string) string {
	if pkg != nil {
		return pkg.Name() + "." + name
	}
	return name
}
