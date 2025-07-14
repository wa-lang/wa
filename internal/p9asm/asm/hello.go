// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build ignore

package main

import (
	"bytes"
	"fmt"

	"wa-lang.org/wa/internal/p9asm/asm"
	"wa-lang.org/wa/internal/p9asm/asm/arch"
	"wa-lang.org/wa/internal/p9asm/asm/lex"
	"wa-lang.org/wa/internal/p9asm/obj"
)

func main() {
	arch := arch.Set(arch.AMD64)
	ctxt := obj.Linknew(arch.LinkArch)

	lexer := lex.NewLexer("hello.p9asm", ctxt)
	parser := asm.NewParser(ctxt, arch, lexer)

	prog, ok := parser.Parse()
	if !ok {
		panic("asm: assembly failed")
	}

	var buf bytes.Buffer
	output := obj.Binitw(&buf)
	obj.Writeobjdirect(ctxt, output)
	output.Flush()

	fmt.Printf("%+v\n", prog)
}

const code = `
TEXT foo(SB), 0, $0
	RET
`
