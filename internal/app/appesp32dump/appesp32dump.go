// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appesp32dump

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/native/link/esp32"
)

var CmdESP32Dump = &cli.Command{
	Hidden: true,
	Name:   "esp32dump",
	Usage:  "dump esp32 image file",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "hello",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "no input file")
			os.Exit(1)
		}

		f, err := esp32.Load(c.Args().First())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(f.Dump())
		return nil
	},
}
