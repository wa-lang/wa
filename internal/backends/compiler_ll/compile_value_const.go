// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_ll

import (
	"github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	"github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
	"github.com/wa-lang/wa/internal/constant"
	"github.com/wa-lang/wa/internal/ssa"
	"github.com/wa-lang/wa/internal/types"
)

// 暂时全部用 I64 表示 int 和 指针

func (p *Compiler) constValue(expr *ssa.Const) llvalue.Value {
	t := expr.Type().Underlying()

	switch t := t.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Int, types.Uint:
			n, _ := constant.Int64Val(expr.Value)
			return llconstant.NewInt(lltypes.I64, n)
		case types.Uintptr:
			n, _ := constant.Int64Val(expr.Value)
			return llconstant.NewInt(lltypes.I64, n)

		case types.Bool, types.UntypedBool:
			return llconstant.NewBool(constant.BoolVal(expr.Value))
		case types.Int8, types.Uint8:
			n, _ := constant.Int64Val(expr.Value)
			return llconstant.NewInt(lltypes.I8, n)
		case types.Int16, types.Uint16:
			n, _ := constant.Int64Val(expr.Value)
			return llconstant.NewInt(lltypes.I16, n)
		case types.Int32, types.Uint32:
			n, _ := constant.Int64Val(expr.Value)
			return llconstant.NewInt(lltypes.I32, n)
		case types.Int64, types.Uint64:
			n, _ := constant.Int64Val(expr.Value)
			return llconstant.NewInt(lltypes.I64, n)

		case types.Float32:
			v, _ := constant.Float64Val(expr.Value)
			return llconstant.NewFloat(lltypes.Float, v)
		case types.Float64:
			v, _ := constant.Float64Val(expr.Value)
			return llconstant.NewFloat(lltypes.Double, v)

		case types.Complex64:
			real_v, _ := constant.Float64Val(constant.Real(expr.Value))
			imag_v, _ := constant.Float64Val(constant.Imag(expr.Value))
			return llconstant.NewStruct(
				p.getLLType(g_typ_complex64_name).(*lltypes.StructType),
				llconstant.NewFloat(lltypes.Float, real_v),
				llconstant.NewFloat(lltypes.Float, imag_v),
			)
		case types.Complex128:
			real_v, _ := constant.Float64Val(constant.Real(expr.Value))
			imag_v, _ := constant.Float64Val(constant.Imag(expr.Value))
			return llconstant.NewStruct(
				p.getLLType(g_typ_complex128_name).(*lltypes.StructType),
				llconstant.NewFloat(lltypes.Double, real_v),
				llconstant.NewFloat(lltypes.Double, imag_v),
			)
		case types.String, types.UntypedString:
			s := constant.StringVal(expr.Value)
			return llconstant.NewStruct(
				p.getLLType(g_typ_string_name).(*lltypes.StructType),
				llconstant.NewCharArray(append([]byte(s), 0)),
				llconstant.NewInt(lltypes.I64, int64(len(s)+1)),
			)
		}
	}

	panic("unknown const type: " + t.String())
}
