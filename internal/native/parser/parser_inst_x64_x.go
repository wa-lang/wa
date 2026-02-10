// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
	"wa-lang.org/wa/internal/native/x64"
)

func (p *parser) parseInst_x64(fn *ast.Func) (inst *ast.Instruction) {
	assert(p.cpu == abi.X64Unix || p.cpu == abi.X64Windows)

	inst = new(ast.Instruction)
	inst.Arg = new(abi.AsArgument)

	inst.CPU = p.cpu
	inst.Doc = p.parseDocComment(&fn.Body.Comments, inst.Pos)
	if inst.Doc != nil {
		fn.Body.Objects = fn.Body.Objects[:len(fn.Body.Objects)-1]
	}

	defer func() {
		inst.Comment = p.parseTailComment(inst.Pos)
		p.consumeSemicolonList()
	}()

	if p.tok == token.IDENT {
		inst.Pos = p.pos
		inst.Label = p.parseIdent()
		p.acceptToken(token.COLON)

		// 后续如果不是指令则结束
		if !p.tok.IsAs() {
			return inst
		}
	}

	inst.Pos = p.pos
	inst.AsName = p.lit
	inst.As = p.parseAs()

	switch inst.As {
	default:
		p.errorf(p.pos, "%v is not x64 instruction", p.tok)
	case x64.AADD:
		panic("TODO")
	case x64.AADDSD:
		panic("TODO")
	case x64.AADDSS:
		panic("TODO")
	case x64.AAND:
		panic("TODO")
	case x64.ACALL:
		panic("TODO")
	case x64.ACDQ:
		panic("TODO")
	case x64.ACMP:
		panic("TODO")
	case x64.ACVTSI2SD:
		panic("TODO")
	case x64.ACVTTSD2SI:
		panic("TODO")
	case x64.ADIV:
		panic("TODO")
	case x64.ADIVSD:
		panic("TODO")
	case x64.ADIVSS:
		panic("TODO")
	case x64.AIDIV:
		panic("TODO")
	case x64.AIMUL:
		panic("TODO")
	case x64.AJA:
		panic("TODO")
	case x64.AJE:
		panic("TODO")
	case x64.AJMP:
		panic("TODO")
	case x64.ALEA:
		panic("TODO")
	case x64.AMOV:
		panic("TODO")
	case x64.AMOVABS:
		panic("TODO")
	case x64.AMOVQ:
		panic("TODO")
	case x64.AMOVSD:
		panic("TODO")
	case x64.AMOVSS:
		panic("TODO")
	case x64.AMOVZX:
		panic("TODO")
	case x64.AMULSD:
		panic("TODO")
	case x64.AMULSS:
		panic("TODO")
	case x64.ANOP:
		panic("TODO")
	case x64.AOR:
		panic("TODO")
	case x64.APOP:
		panic("TODO")
	case x64.APUSH:
		panic("TODO")
	case x64.ARET:
		panic("TODO")
	case x64.ASAR:
		panic("TODO")
	case x64.ASETA:
		panic("TODO")
	case x64.ASETAE:
		panic("TODO")
	case x64.ASETB:
		panic("TODO")
	case x64.ASETBE:
		panic("TODO")
	case x64.ASETE:
		panic("TODO")
	case x64.ASETG:
		panic("TODO")
	case x64.ASETGE:
		panic("TODO")
	case x64.ASETL:
		panic("TODO")
	case x64.ASETLE:
		panic("TODO")
	case x64.ASETNE:
		panic("TODO")
	case x64.ASETNP:
		panic("TODO")
	case x64.ASHL:
		panic("TODO")
	case x64.ASUB:
		panic("TODO")
	case x64.ASUBSD:
		panic("TODO")
	case x64.ASUBSS:
		panic("TODO")
	case x64.AUCOMISD:
		panic("TODO")
	case x64.AXOR:
		panic("TODO")
	}

	panic("unreachable")
}
