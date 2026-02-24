// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appnative

import (
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

var CmdNative_Test = &cli.Command{
	Hidden: true,
	Name:   "test",
	Usage:  "test Wa packages in native mode",
	Action: CmdTestAction,
}
