// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appp9asm

import (
	"fmt"
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
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "dump instructions as they are parsed",
		},
		&cli.StringFlag{
			Name:  "o",
			Usage: "output file; default foo.6 for /a/b/c/foo.s on amd64",
		},
		&cli.BoolFlag{
			Name:  "S",
			Usage: "print assembly and machine code",
		},
		&cli.StringFlag{
			Name:  "trimpath",
			Usage: "remove prefix from recorded source file paths",
		},
		&cli.BoolFlag{
			Name:  "shared",
			Usage: "generate code that can be linked into a shared library",
		},
		&cli.BoolFlag{
			Name:  "dynlink",
			Usage: "support references to Go symbols defined in other shared libraries",
		},
		&cli.StringSliceFlag{
			Name:  "D",
			Usage: "predefined symbol with optional simple value -D=identifer=value; can be set multiple times",
		},
		&cli.StringSliceFlag{
			Name:  "I",
			Usage: "include directory; can be set multiple times",
		},
	},
	Action: func(c *cli.Context) error {
		name := c.Args().First()

		flags := &arch.Flags{
			Debug:       c.Bool("debug"),
			OutputFile:  c.String("o"),
			PrintOut:    c.Bool("S"),
			TrimPath:    c.String("trimpath"),
			Shared:      c.Bool("shared"),
			Dynlink:     c.Bool("dynlink"),
			Defines:     c.StringSlice("D"),
			IncludeDirs: c.StringSlice("I"),
		}

		arch := arch.Set(arch.AMD64)
		ctxt := obj.Linknew(&x86.Linkamd64)

		lexer, err := lex.NewLexer(name, ctxt, flags)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		parser := asm.NewParser(ctxt, arch, lexer, flags)

		prog, ok := parser.Parse()
		if !ok {
			fmt.Fprintf(os.Stderr, "asm: assembly of %s failed", c.Args().First())
			os.Exit(1)
		}

		fmt.Printf("%+v\n", prog)
		return nil
	},
}
