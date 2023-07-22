// 版权 @2023 凹语言 作者。保留所有权利。

package appbuild

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/loader"
)

func Build(opt *appbase.Option, filename string) ([]byte, error) {
	if _, err := os.Lstat(filename); err != nil {
		return nil, fmt.Errorf("%q not found", filename)
	}
	cfg := opt.Config()
	prog, err := loader.LoadProgram(cfg, filename)
	if err != nil {
		return nil, err
	}

	output, err := compiler_wat.New().Compile(prog, "main")

	if err != nil {
		return nil, err
	}

	return []byte(output), nil
}
