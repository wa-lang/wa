package llutil

import (
	"fmt"

	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// Layout specifies how data is to be laid out in memory.
type Layout interface {
	// SizeOf returns the size of the given type in number of bits.
	SizeOf(typ types.Type) int
}

// DefaultLayout provides a default implementation of data layout, which specifies the size of types and how data is to be laid out in memory.
//
// Users may embed this struct and implement the Sizeof method to provide custom handling of specific types and their size in number of bits.
type DefaultLayout struct{}

// SizeOf returns the size of the given type in number of bits.
func (l DefaultLayout) SizeOf(typ types.Type) int {
	switch typ := typ.(type) {
	case *types.IntType:
		return int(typ.BitSize)
	case *types.FloatType:
		switch typ.Kind {
		case types.FloatKindHalf:
			return 16
		case types.FloatKindFloat:
			return 32
		case types.FloatKindDouble:
			return 64
		case types.FloatKindFP128:
			return 128
		case types.FloatKindX86_FP80:
			return 80
		case types.FloatKindPPC_FP128:
			return 128
		default:
			panic(fmt.Errorf("support for size of on floating-point type of kind %v not yet implemented", typ.Kind))
		}
	}
	panic(fmt.Errorf("support for size of on type %T not yet implemented", typ))
}
