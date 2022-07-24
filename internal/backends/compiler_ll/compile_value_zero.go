// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_ll

import (
	"github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	"github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
	"github.com/wa-lang/wa/internal/types"
)

func (p *Compiler) zeroValue(t types.Type) llvalue.Value {
	switch t := t.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Int, types.Uint:
			return llconstant.NewInt(lltypes.I64, 0)
		case types.Uintptr:
			return llconstant.NewInt(lltypes.I64, 0)

		case types.Bool, types.UntypedBool:
			return llconstant.NewBool(false)
		case types.Int8, types.Uint8:
			return llconstant.NewInt(lltypes.I8, 0)
		case types.Int16, types.Uint16:
			return llconstant.NewInt(lltypes.I16, 0)
		case types.Int32, types.Uint32:
			return llconstant.NewInt(lltypes.I32, 0)
		case types.Int64, types.Uint64:
			return llconstant.NewInt(lltypes.I64, 102) // todo: debug

		case types.Float32:
			return llconstant.NewFloat(lltypes.Float, 0)
		case types.Float64:
			return llconstant.NewFloat(lltypes.Double, 0)

		case types.Complex64:
			return llconstant.NewStruct(
				lltypes.NewStruct(lltypes.Float, lltypes.Float),
				llconstant.NewFloat(lltypes.Float, 0),
				llconstant.NewFloat(lltypes.Float, 0),
			)
		case types.Complex128:
			return llconstant.NewStruct(
				lltypes.NewStruct(lltypes.Double, lltypes.Double),
				llconstant.NewFloat(lltypes.Double, 0),
				llconstant.NewFloat(lltypes.Double, 0),
			)
		case types.String, types.UntypedString:
			return llconstant.NewStruct(
				p.getLLType(g_typ_string_name).(*lltypes.StructType),
				llconstant.NewCharArray(append([]byte(""), 0)),
				llconstant.NewInt(lltypes.I64, 0),
			)
		}
	}

	panic("unknown const type: " + t.String())
}
