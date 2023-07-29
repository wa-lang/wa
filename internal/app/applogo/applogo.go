// 版权 @2023 凹语言 作者。保留所有权利。

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
