package applsp

import (
	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/lsp"
)

var CmdLsp = &cli.Command{
	Hidden: true,
	Name:   "lsp",
	Usage:  "run Wa langugage server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "log-file",
			Usage: "set log file",
			Value: "",
		},
	},
	Action: func(c *cli.Context) error {
		opt := &lsp.Option{LogFile: c.String("log-file")}
		return lsp.NewLSPServer(opt).Run()
	},
}
