// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appp9asm

import (
	"fmt"
	"log"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/p9asm/asm"
	"wa-lang.org/wa/internal/p9asm/asm/arch"
	"wa-lang.org/wa/internal/p9asm/asm/lex"
	"wa-lang.org/wa/internal/p9asm/obj"
	"wa-lang.org/wa/internal/p9asm/obj/x86"
)

var CmdP9Asm = &cli.Command{
	Hidden: true,
	Name:   "p9asm",
	Usage:  "plan9 assembly language tool",
	Flags:  []cli.Flag{},
	Action: func(c *cli.Context) error {
		name := c.Args().First()
		data, _ := os.ReadFile(name)

		arch := arch.Set(arch.AMD64)
		ctxt := obj.Linknew(&x86.Linkamd64)

		lexer := lex.NewLexer(name, data)
		parser := asm.NewParser(ctxt, arch, lexer)

		prog, ok := parser.Parse()
		if !ok {
			log.Printf("asm: assembly of %s failed", c.Args().First())
			os.Exit(1)
		}

		fmt.Printf("%+v\n", prog)
		return nil
	},
}
