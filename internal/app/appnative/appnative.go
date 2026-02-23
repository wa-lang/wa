// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appnative

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
)

var CmdNative = &cli.Command{
	Name:  "native",
	Usage: "toolchain for build native target",
	Subcommands: []*cli.Command{
		CmdNative_Build,
		CmdNative_Run,
		CmdNative_Test,
	},
}

var CmdNative_Build = &cli.Command{
	Name:  "build",
	Usage: "compile Wa source code in native mode",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "arch",
			Usage: "set app default arch(x64/loong64)",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "os",
			Usage: "set app default os(windows/linux)",
			Value: "",
		},
	},

	Action: func(c *cli.Context) error {
		fmt.Println("TODO")
		os.Exit(1)
		return nil
	},
}

var CmdNative_Run = &cli.Command{
	Name:  "run",
	Usage: "compile and run Wa program in native mode",

	Action: func(c *cli.Context) error {
		fmt.Println("TODO")
		os.Exit(1)
		return nil
	},
}

var CmdNative_Test = &cli.Command{
	Name:  "test",
	Usage: "test Wa packages in native mode",

	Action: func(c *cli.Context) error {
		fmt.Println("TODO")
		os.Exit(1)
		return nil
	},
}
