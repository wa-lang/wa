package wir

import (
	"github.com/wa-lang/wa/internal/types"

	"github.com/wa-lang/wa/internal/backends/compiler_wasm/wir/wtypes"
	"github.com/wa-lang/wa/internal/logger"
)

func ToWType(from types.Type) wtypes.ValueType {
	switch t := from.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			return wtypes.Int32{}

		case types.Int8, types.UntypedBool:
			logger.Fatalf("ToWType Todo:%T", t)

		case types.Uint8:
			logger.Fatalf("ToWType Todo:%T", t)

		case types.Int16:
			logger.Fatalf("ToWType Todo:%T", t)

		case types.Uint16:
			logger.Fatalf("ToWType Todo:%T", t)

		case types.Int, types.Int32, types.UntypedInt:
			return wtypes.Int32{}

		case types.Uint, types.Uint32:
			logger.Fatalf("ToWType Todo:%T", t)

		case types.Int64:
			logger.Fatalf("ToWType Todo:%T", t)

		case types.Uint64:
			logger.Fatalf("ToWType Todo:%T", t)

		case types.Float32, types.UntypedFloat:
			logger.Fatalf("ToWType Todo:%T", t)

		case types.Float64:
			logger.Fatalf("ToWType Todo:%T", t)

		case types.String:
			logger.Fatalf("ToWType Todo:%T", t)

		default:
			logger.Fatalf("Unknown type:%s", t)
			return nil
		}

	case *types.Tuple:
		switch t.Len() {
		case 0:
			return wtypes.Void{}

		case 1:
			return ToWType(t.At(0).Type())

		default:
			logger.Fatalf("Todo type:%s", t)
		}

	default:
		logger.Fatalf("Todo:%T", t)
	}

	return nil
}
