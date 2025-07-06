// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appp9asm

import (
	"flag"
	"log"
	"os"
	"runtime"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/p9asm/arch"
	"wa-lang.org/wa/internal/p9asm/asm"
	"wa-lang.org/wa/internal/p9asm/flags"
	"wa-lang.org/wa/internal/p9asm/lex"
	"wa-lang.org/wa/internal/p9asm/obj"
)

var CmdP9Asm = &cli.Command{
	Hidden: true,
	Name:   "p9asm",
	Usage:  "plan9 assembly language tool",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "arch",
			Value: runtime.GOARCH,
		},
	},
	Action: func(c *cli.Context) error {
		log.SetFlags(0)
		log.SetPrefix("asm: ")

		architecture := arch.Set(c.String("arch"))
		if architecture == nil {
			log.Fatalf("asm: unrecognized architecture %s", c.String("arch"))
		}

		ctxt := obj.Linknew(architecture.LinkArch)

		ctxt.LineHist.TrimPathPrefix = *flags.TrimPath

		lexer := lex.NewLexer(flag.Arg(0), ctxt)
		parser := asm.NewParser(ctxt, architecture, lexer)

		prog, ok := parser.Parse()
		if !ok {
			log.Printf("asm: assembly of %s failed", c.Args().First())
			os.Exit(1)
		}

		_ = prog // TODO
		return nil
	},
}
