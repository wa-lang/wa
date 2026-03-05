// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appwat2xx

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
)

var CmdWat2XX_arm64 = &cli.Command{
	Name:   "arm64",
	Usage:  "convert a WebAssembly text file to arm64",
	Action: CmdWat2XXAction_arm64,
}

func CmdWat2XXAction_arm64(c *cli.Context) error {
	fmt.Println("TODO")
	os.Exit(1)
	return nil
}
