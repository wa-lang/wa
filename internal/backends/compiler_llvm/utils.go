// 版权 @2022 凹语言 作者。保留所有权利。

package compiler_llvm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wa-lang/wa/internal/constant"
	"github.com/wa-lang/wa/internal/ssa"
	"github.com/wa-lang/wa/internal/types"
)

func getArch(arch string) string {
	pos := strings.Index(arch, "-")
	if pos < 0 {
		return arch
	}
	return arch[0:pos]
}

func isVoidFunc(val *ssa.Call) bool {
	return val.Type().String() == "()"
}

func isConstString(val ssa.Value) bool {
	if _, ok := val.(*ssa.Const); !ok {
		return false
	}
	if t, ok := val.Type().(*types.Basic); ok {
		if t.Kind() == types.String {
			return true
		}
	}
	return false
}

func isConstFloat32(val ssa.Value) bool {
	if _, ok := val.(*ssa.Const); !ok {
		return false
	}
	if t, ok := val.Type().(*types.Basic); ok {
		return t.Kind() == types.Float32
	}
	return false
}

func checkType(ty types.Type) (isFloat bool, isSigned bool) {
	switch t := ty.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Float32, types.Float64, types.UntypedFloat:
			return true, true
		case types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.Uint:
			return false, false
		case types.Bool, types.UntypedBool, types.Int8, types.Int16,
			types.Int32, types.Int64, types.Int, types.UntypedInt:
			return false, true
		case types.String:
			return false, false
		default:
			panic("unknown basic type")
		}
	default:
		panic("basic type is expected")
	}
}

type fmtInfo struct {
	str string
	sz  int
}

func getValueFmt(val ssa.Value, target string) (string, int) {
	defIntFmt := map[string]fmtInfo{
		"avr":     {"%d", 2},
		"thumb":   {"%d", 2},
		"arm":     {"%d", 2},
		"aarch64": {"%ld", 3},
		"riscv32": {"%d", 2},
		"riscv64": {"%ld", 3},
		"x86":     {"%d", 2},
		"x86_64":  {"%ld", 3},
	}

	// Directly combine constant strings to the format string.
	if isConstString(val) {
		str := getValueStr(val)
		siz := len(str)
		str = strings.ReplaceAll(str, "\"", "\\22")
		return str, siz
	}

	switch t := val.Type().(type) {
	case *types.Basic:
		switch t.Kind() {
		// Handle fixed types.
		case types.Int8, types.Int16, types.Int32,
			types.Bool, types.UntypedBool:
			return "%d", 2
		case types.Uint8, types.Uint16, types.Uint32:
			return "%u", 2
		case types.Int64:
			return "%ld", 3
		case types.Uint64:
			return "%lu", 3
		case types.Float32, types.Float64, types.UntypedFloat:
			return "%lf", 3
		// Handle feasible types.
		case types.Int, types.Uint, types.UntypedInt:
			if fmt, ok := defIntFmt[getArch(target)]; ok {
				return fmt.str, fmt.sz
			}
			return "%ld", 3
		case types.String:
			return "%s", 2
		// should never reach here
		default:
			panic("unknown basic type")
		}
	default:
		// TODO: How to print complex tpyes, such arrays and structs?
		panic("unknown type")
	}
}

func getTypeStr(ty types.Type, target string) string {
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
		types.Float32:      "float",
		types.Float64:      "double",
		types.UntypedFloat: "float",
	}

	switch t := ty.(type) {
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

	case *types.Array:
		return fmt.Sprintf("[%d x %s]", t.Len(), getTypeStr(t.Elem(), target))

	case *types.Pointer:
		return getTypeStr(t.Elem(), target) + "*"

	case *types.Struct:
		var retTy strings.Builder
		retTy.WriteString("{")
		for i := 0; i < t.NumFields(); i++ {
			retTy.WriteString(getTypeStr(t.Field(i).Type(), target))
			if i != t.NumFields()-1 {
				retTy.WriteString(", ")
			}
		}
		retTy.WriteString("}")
		return retTy.String()

	case *types.Named:
		return getTypeStr(t.Underlying(), target)

	default:
		panic("unknown type")
	}
}

func getTypeSize(ty types.Type, target string) int {
	// feasible types on different targets
	defInt := map[string]int{
		"avr":     2,
		"thumb":   2,
		"arm":     4,
		"aarch64": 8,
		"riscv32": 4,
		"riscv64": 8,
		"x86":     4,
		"x86_64":  8,
	}
	// fixed types
	expTy := map[types.BasicKind]int{
		types.Bool:         1,
		types.UntypedBool:  1,
		types.Int8:         1,
		types.Uint8:        1,
		types.Int16:        2,
		types.Uint16:       2,
		types.Int32:        4,
		types.Uint32:       4,
		types.Int64:        8,
		types.Uint64:       8,
		types.Float32:      4,
		types.Float64:      8,
		types.UntypedFloat: 4,
	}

	switch t := ty.(type) {
	case *types.Basic:
		// return size of a fixed type
		if sz, ok := expTy[t.Kind()]; ok {
			return sz
		}
		// return size of a feasible type
		switch t.Kind() {
		case types.Int, types.Uint, types.UntypedInt:
			if sz, ok := defInt[getArch(target)]; ok {
				return sz
			}
			return 8
		// should never reach here
		default:
			panic("unknown basic type")
		}

	default:
		panic("unknown type")
	}
}

func getValueStr(val ssa.Value) string {
	switch c := val.(type) {
	case *ssa.Const:
		// Get the full content of the constant string.
		if isConstString(val) {
			return constant.StringVal(c.Value)
		}
		// Drop the type information for non-string constants.
		valStr := val.String()
		if pos := strings.Index(valStr, ":"); pos > 0 {
			valStr = valStr[0:pos]
		}
		// Special form for float32/float64 constants as LLVM-IR requested.
		if isFloat, _ := checkType(val.Type()); isFloat {
			if f, err := strconv.ParseFloat(valStr, 64); err == nil {
				valStr = fmt.Sprintf("%e", f)
			}
		}
		return valStr

	case *ssa.Parameter:
		return "%" + val.Name()

	case *ssa.Global:
		return "@" + getNormalName(val.Name())

	default:
		return "%" + val.Name()
	}
}

func getNormalName(name string) string {
	name = strings.ReplaceAll(name, ".", "__")
	name = strings.ReplaceAll(name, "$", "___")
	return name
}
