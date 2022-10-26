// 版权 @2022 凹语言 作者。保留所有权利。

package compiler_llvm

import (
	"strings"

	"github.com/wa-lang/wa/internal/types"
)

func getArch(arch string) string {
	pos := strings.Index(arch, "-")
	if pos < 0 {
		return arch
	}
	return arch[0:pos]
}

func getTypeStr(from types.Type, target string) string {
	// feasible types on different targets
	defInt := map[string]string{
		"avr":     "i16",
		"thumb":   "i16",
		"arm":     "i32",
		"aarch64": "i64",
		"riscv32": "i32",
		"riscv64": "i64",
		"x86":     "i32",
		"x86_64":  "i64",
	}

	// fixed types
	expTy := map[types.BasicKind]string{
		types.Bool:         "i1",
		types.UntypedBool:  "i1",
		types.Int8:         "i8",
		types.Uint8:        "i8",
		types.Int16:        "i16",
		types.Uint16:       "i16",
		types.Int32:        "i32",
		types.Uint32:       "i32",
		types.Int64:        "i64",
		types.Uint64:       "i64",
		types.Float32:      "f32",
		types.Float64:      "f64",
		types.UntypedFloat: "f32",
	}

	switch t := from.(type) {
	case *types.Basic:
		// return fixed types
		if ty, ok := expTy[t.Kind()]; ok {
			return ty
		}
		// return feasible types
		switch t.Kind() {
		case types.Int, types.Uint, types.UntypedInt:
			if ty, ok := defInt[getArch(target)]; ok {
				return ty
			}
			return "i64"
		// should never reach here
		default:
			panic("unknown basic type")
		}
	default:
		// TODO: support pointer, array and struct types
		panic("unknown type")
	}
}
