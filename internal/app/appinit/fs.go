// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appinit

import (
	"embed"
	"io/fs"
)

//go:embed misc/_example_app
var _exampleAppFS embed.FS

func waroot_GetExampleAppFS() fs.FS {
	fs, err := fs.Sub(_exampleAppFS, "misc/_example_app")
	if err != nil {
		panic(err)
	}
	return fs
}
