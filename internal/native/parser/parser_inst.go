package parser

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
)

func (p *parser) parseInst(fn *ast.Func) (inst *ast.Instruction) {
	switch p.cpu {
	case abi.LOONG64:
		return p.parseInst_loong(fn)
	case abi.RISCV32, abi.RISCV64:
		return p.parseInst_riscv(fn)
	default:
		panic("unreachable")
	}
}
