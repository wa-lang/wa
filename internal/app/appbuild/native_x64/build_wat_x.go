// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package native_x64

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/wat2x64"
)

func BuildApp_wat(opt *appbase.Option, input, outfile string) (mainFunc string, wasmBytes, fsetBytes []byte, err error) {
	// wa build -arch=x64 -target=linux input.wat
	// wa build -arch=x64 -target=windows input.wat

	if s := opt.TargetOS; s != "" && s != config.WaOS_linux && s != config.WaOS_windows {
		panic(fmt.Sprintf("x64 donot support %s", s))
	}

	cpuType := abi.X64Windows
	if strings.EqualFold(opt.TargetOS, config.WaOS_linux) {
		cpuType = abi.X64Unix
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
	_, nasmBytes, err := wat2x64.Wat2X64(input, watData, cpuType)
	if err != nil {
		fmt.Printf("wat2x64 %s failed: %v\n", input, err)
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
