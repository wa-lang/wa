// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appp9link

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/p9asm/link/amd64"
	"wa-lang.org/wa/internal/p9asm/link/arm"
	"wa-lang.org/wa/internal/p9asm/link/arm64"
	"wa-lang.org/wa/internal/p9asm/link/x86"
	"wa-lang.org/wa/internal/p9asm/obj"
)

var CmdP9Link = &cli.Command{
	Hidden: true,
	Name:   "p9link",
	Usage:  "p9asm object link tool",
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
		switch obj.Getwaarch() {
		default:
			fmt.Fprintf(os.Stderr, "link: unknown architecture %q\n", obj.Getwaarch())
			os.Exit(2)
		case "386":
			x86.Main()
		case "amd64", "amd64p32":
			amd64.Main()
		case "arm":
			arm.Main()
		case "arm64":
			arm64.Main()
		}
		return nil
	},
}
