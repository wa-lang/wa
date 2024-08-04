package applsp

import (
	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/lsp"
)

var CmdLsp = &cli.Command{
	Hidden: false,
	Name:   "lsp",
	Usage:  "run Wa langugage server (dev)",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "wa-root",
			Usage: "set wa root",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "wa-os",
			Usage: "set wa os",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "log-file",
			Usage: "set log file",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "sync-file-dir",
			Usage: "set sync file dir",
			Value: "",
		},
	},
	Action: func(c *cli.Context) error {
		opt := &lsp.Option{
			LogFile:     c.String("log-file"),
			SyncFileDir: c.String("sync-file-dir"),
			WaOS:        c.String("wa-os"),
			WaRoot:      c.String("wa-root"),
		}
		return lsp.NewLSPServer(opt).Run()
	},
}
