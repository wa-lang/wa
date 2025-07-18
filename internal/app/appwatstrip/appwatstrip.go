// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appwatstrip

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/wat/watutil/watstrip"
)

var CmdWatStrip = &cli.Command{
	Hidden:    true,
	Name:      "wat-strip",
	Usage:     "remove unused func and global in WebAssembly text file",
	ArgsUsage: "<file.wat>",
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "no input file")
			os.Exit(1)
		}

		infile := c.Args().First()

		source, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		watBytes, err := watstrip.WatStrip(infile, source)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Stdout.Write(watBytes)
		return nil
	},
}
