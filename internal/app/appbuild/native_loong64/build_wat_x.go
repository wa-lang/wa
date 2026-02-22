// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package native_loong64

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/native/wat2la"
)

func BuildApp_wat(opt *appbase.Option, input, outfile string) (mainFunc string, wasmBytes, fsetBytes []byte, err error) {
	// wa build -arch=loong64 -target=linux input.wat

	if s := opt.TargetOS; s != "" && s != config.WaOS_linux {
		panic(fmt.Sprintf("loong64 donot support %s", s))
	}

	// 设置默认输出目标
	if outfile == "" {
		outfile = appbase.ReplaceExt(input, ".wat", ".wa.s")
	}

	watData, err := os.ReadFile(input)
	if err != nil {
		fmt.Printf("read %s failed: %v\n", input, err)
		os.Exit(1)
	}
	_, nasmBytes, err := wat2la.Wat2LA64(input, watData, false)
	if err != nil {
		fmt.Printf("wat2la %s failed: %v\n", input, err)
		os.Exit(1)
	}

	// 保存汇编到文件
	err = os.WriteFile(outfile, nasmBytes, 0666)
	if err != nil {
		fmt.Printf("write %s failed: %v\n", outfile, err)
		os.Exit(1)
	}

	// OK
	return "", nasmBytes, nil, nil
}
