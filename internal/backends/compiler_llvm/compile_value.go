// 版权 @2022 凹语言 作者。保留所有权利。

package compiler_llvm

import (
	"fmt"

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
		if !val.Heap {
			p.output.WriteString("  ")
			p.output.WriteString(getValueStr(val))
			p.output.WriteString(" = alloca ")
			if pt, ok := val.Type().(*types.Pointer); ok {
				p.output.WriteString(getTypeStr(pt.Elem(), p.target))
			} else {
				p.output.WriteString(getTypeStr(val.Type(), p.target))
			}
		} else {
			// TODO: Support heap allocation.
			p.output.WriteString("  ; " + val.Name() + " = " + val.String())
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

	default:
		p.output.WriteString("  ; " + val.Name() + " = " + val.String() + "\n")
		// panic("unsupported Value '" + val.Name() + " = " + val.String() + "'")
	}

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
