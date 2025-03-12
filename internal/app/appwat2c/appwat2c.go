// 版权 @2023 凹语言 作者。保留所有权利。

package appwat2c

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/wat/watutil"
)

var CmdWat2c = &cli.Command{
	Hidden:    true,
	Name:      "wat2c",
	Usage:     "convert a WebAssembly text file to a C source and header",
	ArgsUsage: "<file.wat>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "set code output file",
			Value:   "a.out.c",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "no input file")
			os.Exit(1)
		}

		infile := c.Args().First()
		outfile := c.String("output")

		if outfile == "" {
			outfile = infile
			if n1, n2 := len(outfile), len(".wat"); n1 > n2 {
				if s := outfile[n1-n2:]; strings.EqualFold(s, ".wat") {
					outfile = outfile[:n1-n2]
				}
			}
			outfile += ".c"
		}
		if !strings.HasSuffix(outfile, ".c") {
			outfile += ".c"
		}

		hdrfile := outfile[:len(outfile)-2] + ".h"

		source, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		code, header, err := watutil.Wat2C(infile, source)
		if err != nil {
			os.WriteFile(outfile, code, 0666)
			fmt.Println(err)
			os.Exit(1)
		}

		err = os.WriteFile(hdrfile, header, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = os.WriteFile(outfile, code, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return nil
	},
}
