// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appnative

import (
	"fmt"
	"os"
	"os/exec"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
)

func CmdRunAction(c *cli.Context) error {
	input := c.Args().First()

	if input == "" || input == "." {
		input, _ = os.Getwd()
	}

	var opt = appbase.BuildOptions(c)
	if appbase.HasExt(input, ".wa", ".wz") {
		// 执行单个 wa 脚本, 避免写磁盘
		opt.RunFileMode = true
	}

	exePath, err := BuildApp(opt, input, "")
	if err != nil {
		fmt.Println("appbuild.BuildApp:", err)
		os.Exit(1)
		return nil
	}

	var appArgs []string
	if c.NArg() > 1 {
		appArgs = c.Args().Slice()
	}

	cmd := exec.Command(exePath, appArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(output)
		fmt.Println(err)
		os.Exit(1)
		return nil
	}

	fmt.Println(output)
	return nil
}
