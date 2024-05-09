// 版权 @2023 凹语言 作者。保留所有权利。

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
