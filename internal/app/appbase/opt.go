// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appbase

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/config"
)

// 命令行选项
type Option struct {
	Debug        bool
	RunFileMode  bool // 单文件运行模式
	WaBackend    string
	BuilgTags    []string
	TargetArch   string
	TargetOS     string
	Wat2CNative  bool
	Wat2CPrefix  string
	Wat2CExports map[string]string
	LD_StackSize int
	LD_MaxMemory int
	Optimize     bool // 优化
}

func (opt *Option) Config() *config.Config {
	cfg := config.DefaultConfig()

	cfg.TargetArch = opt.TargetArch
	cfg.TargetOS = opt.TargetOS

	if opt.Debug {
		cfg.Debug = true
	}
	if opt.WaBackend != "" {
		cfg.WaBackend = opt.WaBackend
	}
	if len(opt.BuilgTags) > 0 {
		cfg.BuilgTags = append(cfg.BuilgTags, opt.BuilgTags...)
	}
	if opt.LD_StackSize != 0 {
		cfg.LDFlags.StackSize = opt.LD_StackSize
	}
	if opt.LD_MaxMemory != 0 {
		cfg.LDFlags.MaxMemory = opt.LD_MaxMemory
	}

	cfg.WaSizes.MaxAlign = 8
	cfg.WaSizes.WordSize = 4

	return cfg
}

func BuildOptions(c *cli.Context, waBackend ...string) *Option {
	opt := &Option{
		Debug:        c.Bool("debug"),
		WaBackend:    config.WaBackend_Default,
		BuilgTags:    strings.Fields(c.String("tags")),
		Wat2CNative:  c.Bool("wat2c-native"),
		Wat2CPrefix:  c.String("wat2c-prefix"),
		Wat2CExports: c.StringSliceAsMap("wat2c-exports"),
		LD_StackSize: c.Int("ld-stack-size"),
		LD_MaxMemory: c.Int("ld-max-memory"),
		Optimize:     c.Bool("optimize"),
	}

	targetArch := c.String("arch")
	targetOS := c.String("target")

	if targetArch != "" && !config.CheckWaArch(targetArch) {
		fmt.Printf("unknown arch: %s\n", targetArch)
		os.Exit(1)
	}
	if targetOS != "" && !config.CheckWaOS(targetOS) {
		fmt.Printf("unknown target: %s\n", targetOS)
		os.Exit(1)
	}

	if targetArch == "" {
		targetArch = config.WaArch_Default
	}
	if targetOS == "" {
		targetOS = config.WaOS_Default
	}

	switch targetOS {
	case config.WaOS_js:
		assert(targetArch == config.WaArch_wasm)
		opt.TargetArch = config.WaArch_wasm
		opt.TargetOS = config.WaOS_js
	case config.WaOS_wasm4:
		assert(targetArch == config.WaArch_wasm)
		opt.TargetArch = config.WaArch_wasm
		opt.TargetOS = config.WaOS_wasm4
	case config.WaOS_arduino:
		opt.TargetOS = config.WaOS_arduino
	case config.WaOS_linux:
		switch targetArch {
		case config.WaArch_loong64:
			opt.TargetArch = config.WaArch_loong64
		case config.WaArch_riscv32:
			opt.TargetArch = config.WaArch_riscv32
		case config.WaArch_x64:
			opt.TargetArch = config.WaArch_x64
		default:
			assert(false)
		}
		opt.TargetOS = config.WaOS_linux
	case config.WaOS_windows:
		assert(targetArch == config.WaArch_x64)
		opt.TargetArch = config.WaArch_x64
		opt.TargetOS = config.WaOS_windows
	case config.WaOS_unknown:
		opt.TargetArch = targetArch
		opt.TargetOS = config.WaOS_unknown
	default:
		fmt.Printf("unreachable: target: %s/%s\n", targetOS, targetArch)
		os.Exit(1)
	}

	return opt
}

func assert(ok bool) {
	if !ok {
		panic("assert failed")
	}
}
