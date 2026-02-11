// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import "wa-lang.org/wa/internal/native/ast"

// mov eax, 1
// mov rbp, rsp
// mov [rip + .Wa.Memory.addr], rax
// mov r8b, [rsi]
// mov [rdi], r8b
// mov byte ptr [rcx], '-'
// mov rdi, qword ptr [rbp-8] # arg 0
// mov qword ptr [rbp-8], rax

func (p *parser) parseInst_x64_mov(fn *ast.Func, inst *ast.Instruction) {
	panic("TODO")
}
