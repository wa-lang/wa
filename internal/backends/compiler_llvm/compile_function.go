// 版权 @2022 凹语言 作者。保留所有权利。

package compiler_llvm

import (
	"errors"
	"fmt"

	"github.com/wa-lang/wa/internal/ssa"
	"github.com/wa-lang/wa/internal/types"
)

func (p *Compiler) compileFunction(fn *ssa.Function) error {
	// Translate return type.
	var retType string
	rets := fn.Signature.Results()
	switch rets.Len() {
	case 0:
		retType = "void"
	case 1:
		retType = getTypeStr(rets.At(0).Type(), p.target)
		switch rets.At(0).Type().(type) {
		default:
			return errors.New("type '" + retType + "' can not be returned")
		case *types.Basic, *types.Pointer:
		}
	default:
		return errors.New("multiple return values are not supported")
	}
	p.output.WriteString("define " + retType + " @" + fn.Name() + "(")

	// Translate arguments.
	for i, v := range fn.Params {
		tyStr := getTypeStr(v.Type(), p.target)
		switch v.Type().(type) {
		default:
			return errors.New("type '" + tyStr + "' can not be used as argument")
		case *types.Basic, *types.Pointer, *types.Struct, *types.Array:
		}
		p.output.WriteString(tyStr)
		p.output.WriteString(" ")
		p.output.WriteString(getValueStr(v))
		if i < len(fn.Params)-1 {
			p.output.WriteString(", ")
		}
	}
	p.output.WriteString(") {\n")

	// Initialize all global variables at the beginning of the main function.
	if fn.Name() == "main" {
		p.output.WriteString("__basic_block_init:\n")
		p.output.WriteString("  call void @init()\n")
		p.output.WriteString("  br label %__basic_block_0\n\n")
	}

	// Translate Go SSA intermediate instructions.
	for i, b := range fn.Blocks {
		p.output.WriteString(fmt.Sprintf("__basic_block_%d:\n", i))
		for _, instr := range b.Instrs {
			if err := p.compileInstr(instr); err != nil {
				return err
			}
		}
		if i < len(fn.Blocks)-1 {
			p.output.WriteString("\n")
		}
	}
	p.output.WriteString("}\n\n")

	return nil
}

func (p *Compiler) compileInstr(instr ssa.Instruction) error {
	switch instr := instr.(type) {
	case *ssa.Return:
		switch len(instr.Results) {
		case 0:
			p.output.WriteString("  ret void\n")
		case 1: // ret %type %value
			p.output.WriteString("  ret ")
			p.output.WriteString(getTypeStr(instr.Results[0].Type(), p.target))
			p.output.WriteString(" ")
			p.output.WriteString(getValueStr(instr.Results[0]))
			p.output.WriteString("\n")
		default:
			return errors.New("multiple return values are not supported")
		}

	case ssa.Value:
		if err := p.compileValue(instr); err != nil {
			return err
		}

	case *ssa.Jump:
		p.output.WriteString(fmt.Sprintf("  br label %%__basic_block_%d\n", instr.Block().Succs[0].Index))

	case *ssa.If:
		p.output.WriteString(fmt.Sprintf("  br i1 %s, label %%__basic_block_%d, label %%__basic_block_%d\n", getValueStr(instr.Cond), instr.Block().Succs[0].Index, instr.Block().Succs[1].Index))

	case *ssa.Store:
		p.output.WriteString("  store ")
		p.output.WriteString(getTypeStr(instr.Val.Type(), p.target))
		p.output.WriteString(" ")
		p.output.WriteString(getValueStr(instr.Val))
		p.output.WriteString(", ")
		p.output.WriteString(getTypeStr(instr.Addr.Type(), p.target))
		p.output.WriteString(" ")
		p.output.WriteString(getValueStr(instr.Addr))
		p.output.WriteString("\n")

	default:
		p.output.WriteString("  ; " + instr.String() + "\n")
		// panic("unsupported IR '" + instr.String() + "'")
	}

	return nil
}
