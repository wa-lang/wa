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
	inst.ArgX64 = new(abi.X64Argument)

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

	switch x64.AsOpFormatType(inst.As) {
	case x64.OpFormatType_NoArgs:
		return

	case x64.OpFormatType_Imm:
		inst.ArgX64.Dst = p.parseX64Operand()
		assert(inst.ArgX64.Dst.Kind == abi.X64Operand_Imm)
		return
	case x64.OpFormatType_Reg:
		inst.ArgX64.Dst = p.parseX64Operand()
		assert(inst.ArgX64.Dst.Kind == abi.X64Operand_Reg)
		return
	case x64.OpFormatType_Mem:
		inst.ArgX64.Dst = p.parseX64Operand()
		assert(inst.ArgX64.Dst.Kind == abi.X64Operand_Mem)
		return
	case x64.OpFormatType_Any:
		inst.ArgX64.Dst = p.parseX64Operand()
		return

	case x64.OpFormatType_Imm2Reg:
		inst.ArgX64.Dst = p.parseX64Operand()
		p.acceptToken(token.COMMA)
		inst.ArgX64.Src = p.parseX64Operand()
		assert(inst.ArgX64.Src.Kind == abi.X64Operand_Imm)
		assert(inst.ArgX64.Dst.Kind == abi.X64Operand_Reg)
		return
	case x64.OpFormatType_Imm2Mem:
		inst.ArgX64.Dst = p.parseX64Operand()
		p.acceptToken(token.COMMA)
		inst.ArgX64.Src = p.parseX64Operand()
		assert(inst.ArgX64.Src.Kind == abi.X64Operand_Imm)
		assert(inst.ArgX64.Dst.Kind == abi.X64Operand_Mem)
		return
	case x64.OpFormatType_Reg2Reg:
		inst.ArgX64.Dst = p.parseX64Operand()
		p.acceptToken(token.COMMA)
		inst.ArgX64.Src = p.parseX64Operand()
		assert(inst.ArgX64.Src.Kind == abi.X64Operand_Reg)
		assert(inst.ArgX64.Dst.Kind == abi.X64Operand_Reg)
		return
	case x64.OpFormatType_Mem2Reg:
		inst.ArgX64.Dst = p.parseX64Operand()
		p.acceptToken(token.COMMA)
		inst.ArgX64.Src = p.parseX64Operand()
		assert(inst.ArgX64.Src.Kind == abi.X64Operand_Mem)
		assert(inst.ArgX64.Dst.Kind == abi.X64Operand_Reg)
		return
	case x64.OpFormatType_Reg2Mem:
		inst.ArgX64.Dst = p.parseX64Operand()
		p.acceptToken(token.COMMA)
		inst.ArgX64.Src = p.parseX64Operand()
		assert(inst.ArgX64.Src.Kind == abi.X64Operand_Reg)
		assert(inst.ArgX64.Dst.Kind == abi.X64Operand_Mem)
		return

	case x64.OpFormatType_Any2Any:
		inst.ArgX64.Dst = p.parseX64Operand()
		p.acceptToken(token.COMMA)
		inst.ArgX64.Src = p.parseX64Operand()
		switch inst.ArgX64.Src.Kind {
		case abi.X64Operand_Imm:
			dstKind := inst.ArgX64.Dst.Kind
			assert(dstKind == abi.X64Operand_Reg || dstKind == abi.X64Operand_Mem)
		case abi.X64Operand_Reg:
			dstKind := inst.ArgX64.Dst.Kind
			assert(dstKind == abi.X64Operand_Reg || dstKind == abi.X64Operand_Mem)
		case abi.X64Operand_Mem:
			dstKind := inst.ArgX64.Dst.Kind
			assert(dstKind == abi.X64Operand_Reg)
		}
		return

	default:
		p.errorf(p.pos, "%v is not x64 instruction", p.tok)
	}

	panic("unreachable")
}

// mov eax, 1
// mov rbp, rsp
// mov [rip + .Wa.Memory.addr], rax
// mov r8b, [rsi]
// mov [rdi], r8b
// mov byte ptr [rcx], '-'
// mov rdi, qword ptr [rbp-8] # arg 0
// mov qword ptr [rbp-8], rax

func (p *parser) parseX64Operand() *abi.X64Operand {
	op := &abi.X64Operand{}

	// 检查是否有内存大小前缀 (byte ptr, dword ptr, etc.)
	if ptrType := p.tryParseX64PtrType(p.lit); ptrType != 0 {
		op.PtrTyp = ptrType
		p.next()
		p.acceptIdentToken(token.IDENT, "ptr")
	}

	switch p.tok {
	case token.LBRACK:
		if op.PtrTyp == 0 {
			p.errorf(p.pos, "pointer type missing")
		}

		// 处理内存寻址 [rip + symbol] 或 [reg + offset]
		op.Kind = abi.X64Operand_Mem
		p.next()

		// 处理首个组件(寄存器/符号/数字)
		if p.tok.IsRegister() {
			op.Reg = p.parseRegister()
		} else if p.tok == token.IDENT {
			op.Symbol = p.parseIdent()
		} else if p.tok == token.INT {
			op.Offset = int64(p.parseIntLit())
		}

		// 解析偏移量或符号: + symbol / - 8
		for p.tok == token.ADD || p.tok == token.SUB {
			sign := int64(1)
			if p.tok == token.SUB {
				sign = -1
			}
			p.next()

			switch p.tok {
			case token.INT:
				if op.Offset != 0 {
					p.errorf(p.pos, "offset(%d) exists", op.Offset)
				}
				op.Offset = int64(p.parseIntLit()) * sign
			case token.IDENT:
				if op.Symbol != "" {
					p.errorf(p.pos, "symbol(%s) exists", op.Symbol)
				}
				op.Symbol = p.parseIdent()
			default:
				panic("unreachable")
			}
		}
		p.acceptToken(token.RBRACK)
		return op

	case token.INT, token.CHAR:
		op.Kind = abi.X64Operand_Imm
		op.Imm = int64(p.parseIntLit())
		return op

	case token.IDENT:
		op.Kind = abi.X64Operand_Imm
		op.Symbol = p.parseIdent()
		return op

	default:
		if p.tok.IsRegister() {
			op.Kind = abi.X64Operand_Reg
			op.Reg = p.parseRegister()
			return op
		}
		p.errorf(p.pos, "unexpected x64 operand: %v", p.tok)
		return nil
	}
}

func (p *parser) tryParseX64PtrType(lit string) abi.X64PtrType {
	switch lit {
	case "byte":
		return abi.X64BytePtr
	case "word":
		return abi.X64WordPtr
	case "dword":
		return abi.X64DWordPtr
	case "qword":
		return abi.X64QWordPtr
	}
	return 0
}
