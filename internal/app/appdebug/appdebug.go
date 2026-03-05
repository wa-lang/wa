// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appdebug

import "wa-lang.org/wa/internal/3rdparty/cli"

var CmdDebug = &cli.Command{
	Hidden: true,
	Name:   "native",
	Usage:  "debug box (lsp|dap|...)",
	Subcommands: []*cli.Command{
		CmdLsp,
		CmdDap,
	},
}
