// 版权 @2023 凹语言 作者。保留所有权利。

package appwasm2wat

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
)

var CmdWasm2wat = &cli.Command{
	Hidden:    true,
	Name:      "wasm2wat",
	Usage:     "convert wasm format to wasm text format",
	ArgsUsage: "<file.wasm>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "set output file",
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
			if n1, n2 := len(outfile), len(".wasm"); n1 > n2 {
				if s := outfile[n1-n2:]; strings.EqualFold(s, ".wasm") {
					outfile = outfile[:n1-n2]
				}
			}
			outfile += ".wat"
		}

		source, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_ = source

		return fmt.Errorf("wasm2wat: TODO")
	},
}
