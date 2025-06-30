// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appdap

import (
	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/dap"
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
		dap.Run()
		return nil
	},
}
