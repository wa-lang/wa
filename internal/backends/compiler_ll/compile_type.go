// 版权 @2021 凹语言 作者。保留所有权利。

package compiler

import (
	"github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	"github.com/wa-lang/wa/internal/ssa"
	"github.com/wa-lang/wa/internal/types"
)

func (p *Compiler) compileType(pkg *ssa.Package, typMember *ssa.Type) error {
	if typ, ok := typMember.Type().(*types.Named); ok {
		if ty, ok := typ.Underlying().(*types.Struct); ok {
			p.module.NewTypeDef(typ.Obj().Id(), p.toLLType(ty))
		}
	}

	return nil
}

func (p *Compiler) isBoolType(typ types.Type) bool {
	if t, ok := typ.Underlying().(*types.Basic); ok {
		return t.Info()&types.IsBoolean != 0
	}
	return false
}
func (p *Compiler) isIntType(typ types.Type) bool {
	if t, ok := typ.Underlying().(*types.Basic); ok {
		return t.Info()&types.IsInteger != 0
	}
	return false
}
func (p *Compiler) isUIntType(typ types.Type) bool {
	if t, ok := typ.Underlying().(*types.Basic); ok {
		return t.Info()&types.IsUnsigned != 0
	}
	return false
}
func (p *Compiler) isFloatType(typ types.Type) bool {
	if t, ok := typ.Underlying().(*types.Basic); ok {
		return t.Info()&types.IsFloat != 0
	}
	return false
}
func (p *Compiler) isComplexType(typ types.Type) bool {
	if t, ok := typ.Underlying().(*types.Basic); ok {
		return t.Info()&types.IsComplex != 0
	}
	return false
}
func (p *Compiler) isStrType(typ types.Type) bool {
	if t, ok := typ.Underlying().(*types.Basic); ok {
		return t.Info()&types.IsString != 0
	}
	return false
}

func (p *Compiler) llIsStrType(typ lltypes.Type) bool {
	for _, t := range p.module.TypeDefs {
		if t == typ {
			return t.Name() == g_typ_string_name
		}
	}
	return false
}

func (p *Compiler) waIsStrType(typ types.Type) bool {
	if x, ok := typ.(*types.Basic); ok {
		if x.Kind() == types.String || x.Kind() == types.UntypedString {
			return true
		}
	}
	return false
}

func (p *Compiler) waIsComplexType(typ types.Type) bool {
	if x, ok := typ.(*types.Basic); ok {
		if x.Kind() == types.Complex64 || x.Kind() == types.Complex128 {
			return true
		}
	}
	return false
}

func (p *Compiler) isPointerType(typ types.Type) bool {
	if _, ok := typ.(*types.Pointer); ok {
		return true
	} else if typ, ok := typ.(*types.Basic); ok && typ.Kind() == types.UnsafePointer {
		return true
	} else {
		return false
	}
}

func (p *Compiler) toLLFuncType(t types.Type) *lltypes.FuncType {
	return p.toLLType(t).(*lltypes.FuncType)
}

func (p *Compiler) toLLType(t types.Type) lltypes.Type {
	switch t := t.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Int, types.Uint:
			return lltypes.I64
		case types.Uintptr:
			return lltypes.I64
		case types.UnsafePointer:
			return lltypes.NewPointer(nil)

		case types.Bool, types.UntypedBool:
			return lltypes.I1
		case types.Int8, types.Uint8:
			return lltypes.I8
		case types.Int16, types.Uint16:
			return lltypes.I16
		case types.Int32, types.Uint32:
			return lltypes.I32
		case types.Int64, types.Uint64:
			return lltypes.I64

		case types.Float32:
			return lltypes.Float
		case types.Float64:
			return lltypes.Double

		case types.Complex64:
			structType := lltypes.NewStruct(lltypes.Float, lltypes.Float)
			return structType
		case types.Complex128:
			return lltypes.NewStruct(lltypes.Double, lltypes.Double)

		case types.String, types.UntypedString:
			return p.getLLType(g_typ_string_name)

		default:
			panic("unreachable")
		}

	case *types.Pointer:
		return lltypes.NewPointer(p.toLLType(t.Elem()))

	case *types.Array:
		return lltypes.NewArray(uint64(t.Len()), p.toLLType(t.Elem()))

	case *types.Slice:
		return lltypes.NewStruct(
			p.toLLType(t.Elem()),
			lltypes.I64, // len
			lltypes.I64, // cap
		)

	case *types.Struct:
		fields := make([]lltypes.Type, t.NumFields())
		for i := 0; i < t.NumFields(); i++ {
			fields[i] = p.toLLType(t.Field(i).Type())
		}
		return lltypes.NewStruct(fields...)

	case *types.Tuple:
		fields := make([]lltypes.Type, t.Len())
		for i := 0; i < t.Len(); i++ {
			fields[i] = p.toLLType(t.At(i).Type())
		}
		return lltypes.NewStruct(fields...)

	case *types.Map:
		// todo

	case *types.Interface:
		// todo

	case *types.Named:
		if st, ok := t.Underlying().(*types.Struct); ok {
			fields := make([]lltypes.Type, st.NumFields())
			for i := 0; i < st.NumFields(); i++ {
				fields[i] = p.toLLType(st.Field(i).Type())
			}
			t2 := lltypes.NewStruct(fields...)
			// t2.TypeName = "TODO"
			return t2
		}
		return p.toLLType(t.Underlying())

	case *types.Signature:
		var returnType lltypes.Type
		switch t.Results().Len() {
		case 0:
			returnType = lltypes.Void
		case 1:
			returnType = p.toLLType(t.Results().At(0).Type())
		default:
			fields := make([]lltypes.Type, t.Results().Len())
			for i := 0; i < t.Results().Len(); i++ {
				fields[i] = p.toLLType(t.Results().At(i).Type())
			}
			returnType = lltypes.NewStruct(fields...)
		}

		var paramTypes []lltypes.Type
		if t.Recv() != nil {
			recv := p.toLLType(t.Recv().Type())

			paramTypes = append(paramTypes, recv)
		}

		for i := 0; i < t.Params().Len(); i++ {
			subType := p.toLLType(t.Params().At(i).Type())
			paramTypes = append(paramTypes, subType)
		}

		return lltypes.NewFunc(returnType, paramTypes...)

	default:
		panic("unknown type: " + t.String())
	}
	return nil
}
