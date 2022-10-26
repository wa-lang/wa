// 版权 @2022 凹语言 作者。保留所有权利。

package compiler_llvm

import (
	"errors"

	"github.com/wa-lang/wa/internal/ssa"
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
	default:
		return errors.New("multiple return values are not supported")
	}
	p.output.WriteString("define " + retType + " @" + fn.Name() + "(")

	// Translate arguments.
	for i, v := range fn.Params {
		p.output.WriteString(getTypeStr(v.Type(), p.target) + " %" + v.Name())
		if i < len(fn.Params)-1 {
			p.output.WriteString(", ")
		}
	}
	p.output.WriteString(") {\n")

	// Emit fake function body.
	p.output.WriteString("  ret " + retType)
	if retType != "void" {
		p.output.WriteString(" 0")
	}

	p.output.WriteString("\n}\n\n")
	return nil
}
