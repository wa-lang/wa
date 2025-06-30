// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package applogo

import "wa-lang.org/wa/internal/3rdparty/cli"

var CmdLogo = &cli.Command{
	Name:  "logo",
	Usage: "print Wa text format logo",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "more",
			Aliases: []string{"m"},
			Usage:   "print more text logos",
		},
		&cli.BoolFlag{
			Name:  "svg",
			Usage: "print svg logo",
		},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("svg") {
			PrintLogoSvg()
			return nil
		}
		PrintLogo(c.Bool("more"))
		return nil
	},
}
