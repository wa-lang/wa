// 版权 @2021 凹语言 作者。保留所有权利。

package waroot

import (
	"embed"
	"io/fs"
)

//go:embed _example_template
var _exampleAppFS embed.FS

func GetExampleAppFS() fs.FS {
	fs, err := fs.Sub(_exampleAppFS, "_example_template")
	if err != nil {
		panic(err)
	}
	return fs
}
