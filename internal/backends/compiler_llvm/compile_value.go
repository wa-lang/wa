// 版权 @2022 凹语言 作者。保留所有权利。

package compiler_llvm

import (
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
	p.output.WriteString("  ; " + val.Name() + " = " + val.String() + "\n")
	return nil
}
