// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appp9objdump

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/p9asm/objfile"
)

var CmdP9Objdump = &cli.Command{
	Hidden: true,
	Name:   "p9objdump",
	Usage:  "disassembles executable file start end",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "s",
			Usage: "only dump symbols matching this regexp",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "no input file")
			os.Exit(1)
		}

		var symregexp = c.String("s")
		var symRE *regexp.Regexp

		if symregexp != "" {
			re, err := regexp.Compile(symregexp)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid -s regexp: %v", err)
				os.Exit(1)
			}
			symRE = re
		}

		f, err := objfile.Open(c.Args().First())
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		dis, err := f.Disasm()
		if err != nil {
			fmt.Fprintf(os.Stderr, "disassemble %s: %v", c.Args().First(), err)
			os.Exit(1)
		}

		var start, end uint64 = 0, ^uint64(0)
		if c.NArg() == 3 {
			// disassembly of PC range
			start, err = strconv.ParseUint(strings.TrimPrefix(flag.Arg(1), "0x"), 16, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid start PC: %v", err)
				os.Exit(1)

			}
			end, err = strconv.ParseUint(strings.TrimPrefix(flag.Arg(2), "0x"), 16, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid end PC: %v", err)
				os.Exit(1)
			}
		}

		dis.Print(os.Stdout, symRE, start, end)
		return nil
	},
}
