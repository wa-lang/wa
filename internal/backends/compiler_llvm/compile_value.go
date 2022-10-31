// 版权 @2022 凹语言 作者。保留所有权利。

package compiler_llvm

import (
	"fmt"

	"github.com/wa-lang/wa/internal/ssa"
	"github.com/wa-lang/wa/internal/token"
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

	default:
		p.output.WriteString("  ; " + val.Name() + " = " + val.String() + "\n")
		// panic("unsupported Value '" + val.Name() + " = " + val.String() + "'")
	}

	return nil
}

func (p *Compiler) compileUnOp(val *ssa.UnOp) error {
	switch val.Op {
	case token.SUB:
		p.output.WriteString("  ")
		p.output.WriteString("%" + val.Name())
		p.output.WriteString(" = ")
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
		p.output.WriteString("\n")

	default:
		p.output.WriteString("  ; " + val.Name() + " = " + val.String() + "\n")
		// panic("unsupported Value '" + val.Name() + " = " + val.String() + "'")
	}

	return nil
}

func (p *Compiler) compileBinOp(val *ssa.BinOp) error {
	sintOpMap := map[token.Token]string{
		token.ADD: "add",
		token.SUB: "sub",
		token.MUL: "mul",
		token.QUO: "sdiv",
		token.REM: "srem",
	}
	uintOpMap := map[token.Token]string{
		token.ADD: "add",
		token.SUB: "sub",
		token.MUL: "mul",
		token.QUO: "udiv",
		token.REM: "urem",
	}
	floatOpMap := map[token.Token]string{
		token.ADD: "fadd",
		token.SUB: "fsub",
		token.MUL: "fmul",
		token.QUO: "fdiv",
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

	if op, ok := opMap[val.Op]; ok {
		p.output.WriteString("  ")
		p.output.WriteString("%" + val.Name())
		p.output.WriteString(" = ")
		p.output.WriteString(op)
		p.output.WriteString(" ")
		p.output.WriteString(getTypeStr(val.X.Type(), p.target))
		p.output.WriteString(" ")
		p.output.WriteString(getValueStr(val.X))
		p.output.WriteString(", ")
		p.output.WriteString(getValueStr(val.Y))
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
		if !isVoidFunc(val) {
			p.output.WriteString("  %")
			p.output.WriteString(val.Name())
			p.output.WriteString(" = call ")
			p.output.WriteString(getTypeStr(val.Type(), p.target))
		} else {
			p.output.WriteString("  call void")
		}
		p.output.WriteString(" @")
		p.output.WriteString(val.Call.StaticCallee().Name())
		p.output.WriteString("(")
		// Emit parameters.
		for i, v := range val.Call.Args {
			p.output.WriteString(getTypeStr(v.Type(), p.target))
			p.output.WriteString(" ")
			p.output.WriteString(getValueStr(v))
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

func (p *Compiler) compilePrint(val *ssa.CallCommon, ln bool) error {
	index := len(p.fmts)
	size := int(0)
	format := ""

	// Formulate the format string.
	for _, arg := range val.Args {
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

	// Emit the call instruction and the first parameter.
	p.output.WriteString("  call i32 (i8*, ...) @printf(i8* getelementptr inbounds (")
	p.output.WriteString(fmt.Sprintf("[%d x i8], [%d x i8]* @printfmt%d, i16 0, i16 0)", size, size, index))

	// Emit other parameters and finish the call instruction.
	for _, arg := range val.Args {
		// Omit constant strings.
		if !isConstString(arg) {
			p.output.WriteString(", ")
			p.output.WriteString(getTypeStr(arg.Type(), p.target))
			p.output.WriteString(" ")
			p.output.WriteString(getValueStr(arg))
		}
	}
	p.output.WriteString(")\n")

	// Collect all format strings, and emit them as global variables later.
	p.fmts = append(p.fmts, FmtStr{format, size})

	return nil
}
