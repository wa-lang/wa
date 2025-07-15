// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appp9asm

import (
	"flag"
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
	Usage:  "p9asm language assembly tool",
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

		if flags.PrintOut {
			ctxt.Debugasm = 1
		}
		ctxt.LineHist.TrimPathPrefix = flags.TrimPath
		ctxt.Flag_dynlink = flags.Dynlink
		if flags.Shared || flags.Dynlink {
			ctxt.Flag_shared = 1
		}

		fd, err := os.Create(flags.OutputFile)
		if err != nil {
			log.Fatal(err)
		}

		ctxt.Bso = obj.Binitw(os.Stdout)
		defer ctxt.Bso.Flush()

		ctxt.Diag = log.Fatalf
		output := obj.Binitw(fd)
		fmt.Fprintf(output, "wa object %s %s\n", obj.Getgoos(), obj.Getwaarch())
		fmt.Fprintf(output, "!\n")

		lexer, err := lex.NewLexer(name, ctxt, flags)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		parser := asm.NewParser(ctxt, arch, lexer, flags)

		pList := obj.Linknewplist(ctxt)
		var ok bool
		pList.Firstpc, ok = parser.Parse()
		if !ok {
			log.Printf("asm: assembly of %s failed", flag.Arg(0))
			os.Remove(flags.OutputFile)
			os.Exit(1)
		}
		obj.Writeobjdirect(ctxt, output)
		output.Flush()
		return nil
	},
}
