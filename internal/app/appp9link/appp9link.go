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
)

var CmdP9Link = &cli.Command{
	Hidden: true,
	Name:   "p9link",
	Usage:  "p9asm object link tool",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "arch",
			Usage: "set target architecture (386|amd64|arm|arm64)",
		},
	},
	Action: func(c *cli.Context) error {
		// obj 文件隐含了 arch 信息
		waarch := c.String("arch")
		//_ = obj.Getwaarch()
		switch waarch {
		default:
			fmt.Fprintf(os.Stderr, "link: unknown architecture %q\n", waarch)
			os.Exit(2)
		case "386":
			x86.Main()
		case "amd64":
			amd64.Main()
		case "arm":
			arm.Main()
		case "arm64":
			arm64.Main()
		}
		return nil
	},
}
