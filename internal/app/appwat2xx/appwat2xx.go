// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appwat2xx

import (
	"wa-lang.org/wa/internal/3rdparty/cli"
)

var CmdWat2XX = &cli.Command{
	Hidden: true,
	Name:   "wat2xx",
	Usage:  "convert a WebAssembly text file to xx",
	Subcommands: []*cli.Command{
		CmdWat2XX_wasm,
		CmdWat2XX_la64,
		CmdWat2XX_x64,
		CmdWat2XX_riscv,
		CmdWat2XX_arm64,
		CmdWat2XX_clang,
	},
}
