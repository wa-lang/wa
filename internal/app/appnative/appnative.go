// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appnative

import (
	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
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
		appbase.MakeFlag_output(),
		appbase.MakeFlag_arch(),
		appbase.MakeFlag_target(),
		appbase.MakeFlag_tags(),
	},

	Action: CmdBuildAction,
}

var CmdNative_Run = &cli.Command{
	Name:  "run",
	Usage: "compile and run Wa program in native mode",
	Flags: []cli.Flag{
		appbase.MakeFlag_target(),
		appbase.MakeFlag_tags(),
		appbase.MakeFlag_optimize(),
	},

	Action: CmdRunAction,
}

var CmdNative_Test = &cli.Command{
	Hidden: true,
	Name:   "test",
	Usage:  "test Wa packages in native mode",
	Action: CmdTestAction,
}
