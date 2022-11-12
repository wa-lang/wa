// 版权 @2022 凹语言 作者。保留所有权利。

package compiler_llvm

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/wa-lang/wa/internal/ssa"
	"github.com/wa-lang/wa/internal/token"
	"github.com/wa-lang/wa/internal/types"
)

func (p *Compiler) compileValue(val ssa.Value) error {
	switch val := val.(type) {
	case *ssa.UnOp:
		if err := p.compileUnOp(val); err != nil {
			return err
		}

	case *ssa.BinOp:
		if err := p.compileBinOp(val); err != nil {
			return err
		}

	case *ssa.Call:
		if err := p.compileCall(val); err != nil {
			return err
		}

	case *ssa.Phi:
		p.output.WriteString("  ")
		p.output.WriteString(getValueStr(val))
		p.output.WriteString(" = phi ")
		p.output.WriteString(getTypeStr(val.Type(), p.target))
		for i, v := range val.Edges {
			p.output.WriteString(fmt.Sprintf(" [ %s, %%__basic_block_%d ]", getValueStr(v), val.Block().Preds[i].Index))
			if i < len(val.Edges)-1 {
				p.output.WriteString(",")
			}
		}
		p.output.WriteString("\n")

	case *ssa.Alloc:
		p.output.WriteString("  ")
		p.output.WriteString(getValueStr(val))
		p.output.WriteString(" = alloca ")
		if pt, ok := val.Type().(*types.Pointer); ok {
			p.output.WriteString(getTypeStr(pt.Elem(), p.target))
		} else {
			panic("unsupported type '" + getTypeStr(val.Type(), p.target) + "' in IR Alloc")
		}
		p.output.WriteString("\n")

	case *ssa.IndexAddr:
		p.output.WriteString("  ")
		p.output.WriteString(getValueStr(val))
		p.output.WriteString(" = getelementptr inbounds ")
		if t, ok := val.X.Type().(*types.Pointer); ok {
			TyStr := getTypeStr(t, p.target)
			TyStr = TyStr[0 : len(TyStr)-1]
			p.output.WriteString(TyStr)
			p.output.WriteString(", ")
		} else {
			panic("a pointer type value is expected")
		}
		p.output.WriteString(getTypeStr(val.X.Type(), p.target))
		p.output.WriteString(" ")
		p.output.WriteString(getValueStr(val.X))
		p.output.WriteString(", ")
		p.output.WriteString(getTypeStr(val.Index.Type(), p.target))
		p.output.WriteString(" 0, ")
		p.output.WriteString(getTypeStr(val.Index.Type(), p.target))
		p.output.WriteString(" ")
		p.output.WriteString(getValueStr(val.Index))
		p.output.WriteString("\n")

	case *ssa.FieldAddr:
		p.output.WriteString("  ")
		p.output.WriteString(getValueStr(val))
		p.output.WriteString(" = getelementptr inbounds ")
		if t, ok := val.X.Type().(*types.Pointer); ok {
			TyStr := getTypeStr(t, p.target)
			TyStr = TyStr[0 : len(TyStr)-1]
			p.output.WriteString(TyStr)
			p.output.WriteString(", ")
		} else {
			panic("a pointer type value is expected")
		}
		p.output.WriteString(getTypeStr(val.X.Type(), p.target))
		p.output.WriteString(" ")
		p.output.WriteString(getValueStr(val.X))
		p.output.WriteString(", ")
		p.output.WriteString(getTypeStr(types.Typ[types.Int], p.target))
		p.output.WriteString(" 0, i32 ")
		p.output.WriteString(strconv.Itoa(val.Field))
		p.output.WriteString("\n")

	case *ssa.Convert:
		if err := p.compileConvert(val); err != nil {
			return err
		}

	default:
		p.output.WriteString("  ; " + val.Name() + " = " + val.String() + "\n")
		// panic("unsupported Value '" + val.Name() + " = " + val.String() + "'")
	}

	return nil
}

func (p *Compiler) compileConvert(val *ssa.Convert) error {
	frIsFloat, frIsSigned := checkType(val.X.Type())
	toIsFloat, toIsSigned := checkType(val.Type())
	frSize := getTypeSize(val.X.Type(), p.target)
	toSize := getTypeSize(val.Type(), p.target)

	p.output.WriteString("  ")
	p.output.WriteString(getValueStr(val))
	p.output.WriteString(" = ")

	switch {
	case !frIsFloat && !toIsFloat: // integral type -> integral type
		switch {
		case frSize == toSize:
			p.output.WriteString("add ")
			p.output.WriteString(getTypeStr(val.Type(), p.target))
			p.output.WriteString(" ")
			p.output.WriteString(getValueStr(val.X))
			p.output.WriteString(", 0\n")
			return nil

		case frSize > toSize:
			p.output.WriteString("trunc ")
		case frSize < toSize:
			if !toIsSigned || !frIsSigned {
				p.output.WriteString("zext ")
			} else {
				p.output.WriteString("sext ")
			}
		default:
			panic("unkown type convert: " + val.String())
		}

	case frIsFloat && toIsFloat: // float point type -> float point type
		switch {
		case frSize > toSize:
			p.output.WriteString("fptrunc ")
		case frSize < toSize:
			p.output.WriteString("fpext ")
		default:
			panic("unkown type convert: " + val.String())
		}

	case !frIsFloat && toIsFloat: // integeral type -> float point type
		if frIsSigned {
			p.output.WriteString("sitofp ")
		} else {
			p.output.WriteString("uitofp ")
		}

	case frIsFloat && !toIsFloat: //float point type -> integral type
		if toIsSigned {
			p.output.WriteString("fptosi ")
		} else {
			p.output.WriteString("fptoui ")
		}

	default:
		panic("unkown type convert: " + val.String())
	}

	p.output.WriteString(getTypeStr(val.X.Type(), p.target))
	p.output.WriteString(" ")
	p.output.WriteString(getValueStr(val.X))
	p.output.WriteString(" to ")
	p.output.WriteString(getTypeStr(val.Type(), p.target))
	p.output.WriteString("\n")
	return nil
}

func (p *Compiler) compileUnOp(val *ssa.UnOp) error {
	p.output.WriteString("  ")
	p.output.WriteString(getValueStr(val))
	p.output.WriteString(" = ")

	switch val.Op {
	case token.SUB:
		if isFloat, _ := checkType(val.X.Type()); isFloat {
			p.output.WriteString("fneg ")
			p.output.WriteString(getTypeStr(val.X.Type(), p.target))
			p.output.WriteString(" ")
		} else {
			p.output.WriteString("sub ")
			p.output.WriteString(getTypeStr(val.X.Type(), p.target))
			p.output.WriteString(" 0, ")
		}
		p.output.WriteString(getValueStr(val.X))

	case token.MUL:
		p.output.WriteString("load ")
		p.output.WriteString(getTypeStr(val.Type(), p.target))
		p.output.WriteString(", ")
		p.output.WriteString(getTypeStr(val.X.Type(), p.target))
		p.output.WriteString(" ")
		p.output.WriteString(getValueStr(val.X))

	default:
		p.output.WriteString("  ; " + val.Name() + " = " + val.String())
		// panic("unsupported Value '" + val.Name() + " = " + val.String() + "'")
	}

	p.output.WriteString("\n")
	return nil
}

func (p *Compiler) compileBinOp(val *ssa.BinOp) error {
	sintOpMap := map[token.Token]string{
		token.ADD: "add",
		token.SUB: "sub",
		token.MUL: "mul",
		token.QUO: "sdiv",
		token.REM: "srem",
		token.EQL: "icmp eq",
		token.NEQ: "icmp ne",
		token.LSS: "icmp slt",
		token.GTR: "icmp sgt",
		token.LEQ: "icmp sle",
		token.GEQ: "icmp sge",
	}
	uintOpMap := map[token.Token]string{
		token.ADD: "add",
		token.SUB: "sub",
		token.MUL: "mul",
		token.QUO: "udiv",
		token.REM: "urem",
		token.EQL: "icmp eq",
		token.NEQ: "icmp ne",
		token.LSS: "icmp ult",
		token.GTR: "icmp ugt",
		token.LEQ: "icmp ule",
		token.GEQ: "icmp uge",
	}
	floatOpMap := map[token.Token]string{
		token.ADD: "fadd",
		token.SUB: "fsub",
		token.MUL: "fmul",
		token.QUO: "fdiv",
		token.EQL: "fcmp oeq",
		token.NEQ: "fcmp une",
		token.LSS: "fcmp olt",
		token.GTR: "fcmp ogt",
		token.LEQ: "fcmp ole",
		token.GEQ: "fcmp oge",
	}

	// Type float, signed int, unsigned int each has its own LLVM-IR.
	var opMap map[token.Token]string
	isFloat, isSigned := checkType(val.X.Type())
	if isFloat {
		opMap = floatOpMap
	} else if isSigned {
		opMap = sintOpMap
	} else {
		opMap = uintOpMap
	}

	// Special process for float32 constants.
	xStr, yStr := "", ""
	if isConstFloat32(val.X) {
		xStr = p.wrapConstFloat32(val.X)
	} else {
		xStr = getValueStr(val.X)
	}
	if isConstFloat32(val.Y) {
		yStr = p.wrapConstFloat32(val.Y)
	} else {
		yStr = getValueStr(val.Y)
	}

	// Emit the binary-operator instruction.
	if op, ok := opMap[val.Op]; ok {
		p.output.WriteString("  ")
		p.output.WriteString("%" + val.Name())
		p.output.WriteString(" = ")
		p.output.WriteString(op)
		p.output.WriteString(" ")
		p.output.WriteString(getTypeStr(val.X.Type(), p.target))
		p.output.WriteString(" ")
		p.output.WriteString(xStr)
		p.output.WriteString(", ")
		p.output.WriteString(yStr)
		p.output.WriteString("\n")
		return nil
	}

	p.output.WriteString("  ; " + val.Name() + " = " + val.String() + "\n")
	return nil
	// panic("unsupported Value '" + val.Name() + " = " + val.String() + "'")
}

func (p *Compiler) compileCall(val *ssa.Call) error {
	switch val.Call.Value.(type) {
	case *ssa.Function:
		// Special process for float32 constants.
		paf32 := map[int]string{}
		for i, v := range val.Call.Args {
			if isConstFloat32(v) {
				paf32[i] = p.wrapConstFloat32(v)
			}
		}
		// Emit the call instruction.
		if !isVoidFunc(val) {
			ty := getRealType(val.Type())
			tyStr := getTypeStr(ty, p.target)
			switch ty.(type) {
			default:
				return errors.New("type '" + tyStr + "' can not be returned")
			case *types.Basic, *types.Pointer, *types.Array, *types.Struct:
				// Only allow scalar/pointer/array/struct types.
			}
			p.output.WriteString("  %")
			p.output.WriteString(val.Name())
			p.output.WriteString(" = call ")
			p.output.WriteString(tyStr)
		} else {
			p.output.WriteString("  call void")
		}
		// Emit the function name.
		p.output.WriteString(" @")
		callee := val.Call.StaticCallee()
		if len(callee.LinkName()) > 0 {
			p.output.WriteString(callee.LinkName())
		} else {
			p.output.WriteString(getNormalName(callee.Pkg.Pkg.Path() + "." + callee.Name()))
		}
		p.output.WriteString("(")
		// Emit parameters.
		for i, v := range val.Call.Args {
			ty := getRealType(v.Type())
			tyStr := getTypeStr(ty, p.target)
			switch ty.(type) {
			default:
				return errors.New("type '" + tyStr + "' can not be used as parameter")
			case *types.Basic, *types.Pointer, *types.Array, *types.Struct:
				// Only allow scalar/pointer/array/struct types.
			}
			p.output.WriteString(tyStr)
			p.output.WriteString(" ")
			if pv, ok := paf32[i]; ok {
				p.output.WriteString(pv)
			} else {
				p.output.WriteString(getValueStr(v))
			}
			if i < len(val.Call.Args)-1 {
				p.output.WriteString(", ")
			}
		}
		p.output.WriteString(")\n")

	case *ssa.Builtin:
		if err := p.compileBuiltin(val.Common()); err != nil {
			return err
		}

	default:
		p.output.WriteString("  ; " + val.Name() + " = " + val.String() + "\n")
		// panic("unsupported Value '" + val.Name() + " = " + val.String() + "'")
	}

	return nil
}

func (p *Compiler) compileBuiltin(val *ssa.CallCommon) error {
	switch val.Value.Name() {
	case "println":
		if err := p.compilePrint(val, true); err != nil {
			return err
		}

	case "print":
		if err := p.compilePrint(val, false); err != nil {
			return err
		}

	default:
		p.output.WriteString("  ; " + val.String() + "\n")
		// panic("unsupported builtin '" + val.String() + "'")
	}

	return nil
}

type printArg struct {
	name string
	typ  string
}

func (p *Compiler) compilePrint(val *ssa.CallCommon, ln bool) error {
	index := len(p.fmts)
	size := int(0)
	format := ""

	// Formulate the format string.
	for _, arg := range val.Args {
		ty := getRealType(arg.Type())
		if _, ok := ty.(*types.Basic); !ok {
			return errors.New("can only print scalar typed values")
		}
		f, s := getValueFmt(arg, p.target)
		format += f
		size += s
	}
	// Add a new line to the format string.
	if ln {
		format += "\\0A"
		size += 1
	}
	// End the format string.
	format += "\\00"
	size += 1

	// Emit some type promote instructions for f32/i16/i8/u16/u8 parameters.
	varMap := map[int]printArg{}
	for i, arg := range val.Args {
		argty := getRealType(arg.Type())
		if t, ok := argty.(*types.Basic); ok {
			instr, ty := "", ""
			switch t.Kind() {
			case types.Float32:
				instr, ty = "fpext", "double"
			case types.Uint16, types.Uint8:
				instr, ty = "zext", "i32"
			case types.Int16, types.Int8:
				instr, ty = "sext", "i32"
			default:
			}
			if len(instr) > 0 {
				name := fmt.Sprintf("%%ptmp%d", rand.Int())
				varMap[i] = printArg{name, ty}
				p.output.WriteString("  ")
				p.output.WriteString(name)
				p.output.WriteString(" = ")
				p.output.WriteString(instr)
				p.output.WriteString(" ")
				p.output.WriteString(getTypeStr(arg.Type(), p.target))
				p.output.WriteString(" ")
				p.output.WriteString(getValueStr(arg))
				p.output.WriteString(" to ")
				p.output.WriteString(ty)
				p.output.WriteString("\n")
			}
		}
	}

	// Emit the call instruction and the first parameter.
	p.output.WriteString("  call i32 (i8*, ...) @printf(i8* getelementptr inbounds (")
	p.output.WriteString(fmt.Sprintf("[%d x i8], [%d x i8]* @printfmt%d, i16 0, i16 0)", size, size, index))

	// Emit other parameters and finish the call instruction.
	for i, arg := range val.Args {
		// Omit constant strings.
		if !isConstString(arg) {
			p.output.WriteString(", ")
			if vn, ok := varMap[i]; ok {
				p.output.WriteString(vn.typ)
				p.output.WriteString(" ")
				p.output.WriteString(vn.name)
			} else {
				ty := getRealType(arg.Type())
				p.output.WriteString(getTypeStr(ty, p.target))
				p.output.WriteString(" ")
				p.output.WriteString(getValueStr(arg))
			}
		}
	}
	p.output.WriteString(")\n")

	// Collect all format strings, and emit them as global variables later.
	p.fmts = append(p.fmts, FmtStr{format, size})

	return nil
}

func (p *Compiler) wrapConstFloat32(val ssa.Value) string {
	if !isConstFloat32(val) {
		panic("a float32 constant is expected")
	}
	valStr := getValueStr(val)
	ret := fmt.Sprintf("%%tmp%d", rand.Int())
	p.output.WriteString("  ")
	p.output.WriteString(ret)
	p.output.WriteString(" = fptrunc double ")
	p.output.WriteString(valStr)
	p.output.WriteString(" to float\n")
	return ret
}
