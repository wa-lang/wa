// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appbase

import (
	"fmt"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/config"
)

// 输出路径
func MakeFlag_output() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "set output file",
		Value:   "",
	}
}

// 构建的指令集
func MakeFlag_arch() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "arch",
		Usage: fmt.Sprintf("set arch type (%s)", strings.Join(config.WaArch_List, "|")),
		Value: "",
	}
}

// 构建的目标
func MakeFlag_target() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "target",
		Usage: fmt.Sprintf("set target type (%s)", strings.Join(config.WaOS_List, "|")),
		Value: "",
	}
}

// 通过 wat2c 同时生成native的C代码
func MakeFlag_wat2c_native() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  "wat2c-native",
		Usage: "use wat2c generated native c code",
	}
}

// wat2c 生成名字的前缀
func MakeFlag_wat2c_prefix() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "wat2c-prefix",
		Usage: "name prefix to use in wat2c generated code",
		Value: "app",
	}
}

// 构建的 Tags
func MakeFlag_tags() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "tags",
		Usage: "set build tags",
	}
}

func MakeFlag_ld_stack_size() *cli.IntFlag {
	return &cli.IntFlag{
		Name:  "ld-stack-size",
		Usage: "set stack size",
	}
}

func MakeFlag_ld_max_memory() *cli.IntFlag {
	return &cli.IntFlag{
		Name:  "ld-max-memory",
		Usage: "set max memory size",
	}
}

func MakeFlag_optimize() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  "optimize",
		Usage: "enable optimize flag",
	}
}
