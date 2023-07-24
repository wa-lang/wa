package appyacc

import (
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
)

var CmdYacc = &cli.Command{
	Name:      "yacc",
	Usage:     "generates parsers for LALR(1) grammars",
	ArgsUsage: "<input>",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "l",
			Usage: "disable line directives",
		},
		&cli.StringFlag{
			Name:  "o",
			Usage: "set parser output file",
			Value: "y.wa",
		},
		&cli.StringFlag{
			Name:  "p",
			Usage: "name prefix to use in generated code",
			Value: "yy",
		},
		&cli.StringFlag{
			Name:  "v",
			Usage: "create parsing tables",
			Value: "y.output",
		},
		&cli.StringFlag{
			Name:  "c",
			Usage: "set copyright file",
			Value: "",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() != 1 {
			cli.ShowSubcommandHelpAndExit(c, 1)
		}
		InitFlags(Flags{
			Oflag:     c.String("o"),
			Vflag:     c.String("v"),
			Lflag:     c.Bool("l"),
			Prefix:    c.String("p"),
			Copyright: loadCopyright(c.String("c")),
		})
		Main(c.Args().First())
		return nil
	},
}

func loadCopyright(filename string) string {
	data, _ := os.ReadFile(filename)
	return string(data)
}
