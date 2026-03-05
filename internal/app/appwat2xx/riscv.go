// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appwat2xx

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
)

var CmdWat2XX_riscv = &cli.Command{
	Name:   "riscv",
	Usage:  "convert a WebAssembly text file to riscv",
	Action: CmdWat2XXAction_riscv,
}

func CmdWat2XXAction_riscv(c *cli.Context) error {
	fmt.Println("TODO")
	os.Exit(1)
	return nil
}
