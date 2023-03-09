// 版权 @2023 凹语言 作者。保留所有权利。

package app

import (
	"os"
	"strings"

	"wa-lang.org/wa/internal/format"
)

func (p *App) Fmt(path string) error {
	if path == "" {
		path, _ = os.Getwd()
	}

	if strings.HasSuffix(path, "...") {
		panic("TODO: fmt dir/...")
	}

	fi, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		panic("TODO: fmt dir")
	}

	code, err := format.File(nil, path, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(path, code, 0666)
}
