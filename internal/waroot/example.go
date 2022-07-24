// 版权 @2021 凹语言 作者。保留所有权利。

package waroot

import (
	"embed"
	"io/fs"
)

//go:embed _example_app
var _exampleAppFS embed.FS

//go:embed _example_vendor
var _exampleVendorFS embed.FS

func GetExampleAppFS() fs.FS {
	fs, err := fs.Sub(_exampleAppFS, "_example_app")
	if err != nil {
		panic(err)
	}
	return fs
}

func GetExampleVendorFS() fs.FS {
	fs, err := fs.Sub(_exampleVendorFS, "_example_vendor")
	if err != nil {
		panic(err)
	}
	return fs
}
