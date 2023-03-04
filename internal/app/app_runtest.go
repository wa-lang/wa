// 版权 @2023 凹语言 作者。保留所有权利。

package app

import (
	"os"
	"strings"
)

func (p *App) RunTest(path string) error {
	if path == "" {
		path, _ = os.Getwd()
	}

	if strings.HasSuffix(path, "...") {
		panic("TODO: test dir/...")
	}

	fi, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		panic("TODO: test path must be dir")
	}

	panic("todo")
}
