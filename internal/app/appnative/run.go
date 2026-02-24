// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appnative

import (
	"fmt"
	"os"
	"os/exec"

	"wa-lang.org/wa/internal/3rdparty/cli"
)

var CmdNative_Run = &cli.Command{
	Name:  "run",
	Usage: "compile and run Wa program in native mode",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "tags",
			Usage: "set build tags",
		},
		&cli.BoolFlag{
			Name:  "optimize",
			Usage: "enable optimize flag",
		},
	},

	Action: CmdRunAction,
}

func CmdRunAction(c *cli.Context) error {
	exePath, err := doCmdBuildAction(c)
	if err != nil {
		fmt.Println("appnative.BuildApp:", err)
		os.Exit(1)
		return nil
	}

	cmd := exec.Command(exePath, c.Args().Slice()...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		fmt.Println(err)
		os.Exit(1)
		return nil
	}

	fmt.Print(string(output))
	return nil
}
