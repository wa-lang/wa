package cir

import (
	"wa-lang.org/wa/internal/backends/compiler_c/cir/ctypes"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/types"
)

func ToCType(from types.Type) ctypes.Type {
	switch t := from.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			return ctypes.Bool

		case types.Int8, types.UntypedBool:
			return ctypes.Int8

		case types.Uint8:
			return ctypes.Uint8

		case types.Int16:
			return ctypes.Int16

		case types.Uint16:
			return ctypes.Uint16

		case types.Int, types.Int32, types.UntypedInt:
			return ctypes.Int32

		case types.Uint, types.Uint32:
			return ctypes.Uint32

		case types.Int64:
			return ctypes.Int64

		case types.Uint64:
			return ctypes.Uint64

		case types.Float32, types.UntypedFloat:
			return ctypes.Float

		case types.Float64:
			return ctypes.Double

		case types.String:
			return &ctypes.StringVar{}

		default:
			logger.Fatalf("Unknown type:%s", t)
			return nil
		}

	case *types.Tuple:
		switch t.Len() {
		case 0:
			return ctypes.Void

		case 1:
			return ToCType(t.At(0).Type())

		default:
			var vars []ctypes.Type
			for i := 0; i < t.Len(); i++ {
				vars = append(vars, ToCType(t.At(i).Type()))
			}
			return ctypes.NewTuple(vars)
		}

	case *types.Pointer:
		return ctypes.NewRefType(ToCType(t.Elem()))

	case *types.Named:
		switch ut := t.Underlying().(type) {
		case *types.Struct:
			var fs []ctypes.Field
			for i := 0; i < ut.NumFields(); i++ {
				f := ut.Field(i)
				var v ctypes.Field
				bt := ToCType(f.Type())
				if f.Embedded() {
					v = *ctypes.NewField("$"+bt.CIRString(), bt)
				} else {
					v = *ctypes.NewField(f.Name(), bt)
				}
				fs = append(fs, v)
			}
			return ctypes.NewStruct(t.Obj().Name(), fs)

		default:
			logger.Fatalf("ToCType Todo:%T", ut)
		}

	case *types.Array:
		return ctypes.NewArray(t.Len(), ToCType(t.Elem()))

	case *types.Slice:
		return ctypes.NewSlice(ToCType(t.Elem()))

	default:
		logger.Fatalf("ToCType Todo:%T", t)
	}

	return nil
}
