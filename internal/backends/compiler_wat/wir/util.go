// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strings"

	"github.com/wa-lang/wa/internal/types"

	"github.com/wa-lang/wa/internal/logger"
)

func ToWType(from types.Type) ValueType {
	switch t := from.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool, types.UntypedBool, types.Int, types.UntypedInt:
			return I32{}

		case types.Int32:
			if t.Name() == "rune" {
				return RUNE{}
			} else {
				return I32{}
			}

		case types.Uint32:
			return U32{}

		case types.Int64:
			return I64{}

		case types.Uint64:
			return U64{}

		case types.Float32, types.UntypedFloat:
			return F32{}

		case types.Float64:
			return F64{}

		case types.Int8:
			return I8{}

		case types.Uint8:
			return U8{}

		case types.Int16:
			return I16{}

		case types.Uint16:
			return U16{}

		case types.String:
			return NewString()

		default:
			logger.Fatalf("Unknown type:%s", t)
			return nil
		}

	case *types.Tuple:
		switch t.Len() {
		case 0:
			return VOID{}

		case 1:
			return ToWType(t.At(0).Type())

		default:
			var feilds []ValueType
			for i := 0; i < t.Len(); i++ {
				feilds = append(feilds, ToWType(t.At(i).Type()))
			}
			return NewTuple(feilds)
		}

	case *types.Pointer:
		return NewRef(ToWType(t.Elem()))

	case *types.Named:
		switch ut := t.Underlying().(type) {
		case *types.Struct:
			var fs []Field
			for i := 0; i < ut.NumFields(); i++ {
				f := ut.Field(i)
				wtyp := ToWType(f.Type())
				if f.Embedded() {
					fs = append(fs, NewField("$"+wtyp.Name(), wtyp))
				} else {
					fs = append(fs, NewField(f.Name(), wtyp))
				}
			}
			return NewStruct(GetPkgMangleName(t.Obj().Pkg().Path())+t.Obj().Name(), fs)

		//case *types.Signature:
		//	ut.

		default:
			logger.Fatalf("Todo:%T", ut)
		}

	case *types.Array:
		return NewArray(ToWType(t.Elem()), int(t.Len()))

	case *types.Slice:
		return NewSlice(ToWType(t.Elem()))

	default:
		logger.Fatalf("Todo:%T", t)
	}

	return nil
}

func IsNumber(v Value) bool {
	switch v.Type().(type) {
	case I8, U8, I16, U16, I32, U32, I64, U64, F32, F64:
		return true
	}

	return false
}

func GetPkgMangleName(pkg_path string) string {
	return strings.ReplaceAll(pkg_path, "/", "$") + "."
}

func ExtractField(x Value, field_name string) Value {
	switch x := x.(type) {
	case *aStruct:
		return x.Extract(field_name)

	default:
		logger.Fatalf("Todo:%T", x)
	}

	return nil
}
