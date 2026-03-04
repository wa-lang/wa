// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appasm2xx

import "wa-lang.org/wa/internal/3rdparty/cli"

var CmdAsm2xx = &cli.Command{
	Hidden: true,
	Name:   "asm2xx",
	Usage:  "convert native asm to executable file",
	Subcommands: []*cli.Command{
		CmdAsm2xx_wasm,
		CmdAsm2xx_elf,
		CmdAsm2xx_pe,
	},
}
