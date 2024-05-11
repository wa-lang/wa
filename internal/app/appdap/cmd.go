package appdap

import (
	"fmt"

	"wa-lang.org/wa/internal/3rdparty/cli"
)

var CmdDap = &cli.Command{
	Hidden: true,
	Name:   "dap",
	Usage:  "run DAP server",
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
	},
	Action: func(c *cli.Context) error {
		fmt.Println("TODO: DAP")
		return nil
	},
}
