// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/loong64"
	"wa-lang.org/wa/internal/native/token"
)

func (p *parser) parseInst_loong(fn *ast.Func) (inst *ast.Instruction) {
	assert(p.cpu == abi.LOONG64)

	inst = new(ast.Instruction)
	inst.Arg = new(abi.AsArgument)

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

	// 查询指令的参数格式
	if !loong64.AsValid(inst.As) {
		p.errorf(p.pos, "%v is not loong instruction", inst.As)
	}

	// 查询指令的参数元信息
	for _, argTyp := range loong64.AsArgs(inst.As) {
		_ = argTyp // TODO
	}

	panic("native/parser.parser.parseInst_loong: TODO")
}
