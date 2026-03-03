// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appesp32

import "wa-lang.org/wa/internal/3rdparty/cli"

var CmdESP32 = &cli.Command{
	Hidden: true,
	Name:   "esp32",
	Usage:  "toolchain for esp32",
	Subcommands: []*cli.Command{
		CmdESP32Build,
		CmdESP32Flash,
		CmdESP32Dump,
	},
}
