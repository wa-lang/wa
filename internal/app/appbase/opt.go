// 版权 @2019 凹语言 作者。保留所有权利。

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
	Target       string
	LD_StackSize int
	LD_MaxMemory int
	Optimize     bool // 优化
}

func (opt *Option) Config() *config.Config {
	cfg := config.DefaultConfig()

	cfg.Target = opt.Target

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
		LD_StackSize: c.Int("ld-stack-size"),
		LD_MaxMemory: c.Int("ld-max-memory"),
		Optimize:     c.Bool("optimize"),
	}

	if target := c.String("target"); target != "" && !config.CheckWaOS(target) {
		fmt.Printf("unknown target: %s\n", c.String("target"))
		os.Exit(1)
	}

	switch c.String("target") {
	case "":
		// read from default or wa.mod
	case config.WaOS_js:
		opt.Target = config.WaOS_js
	case config.WaOS_wasi:
		opt.Target = config.WaOS_wasi
	case config.WaOS_wasm4:
		opt.Target = config.WaOS_wasm4
	case config.WaOS_arduino:
		opt.Target = config.WaOS_arduino
	case config.WaOS_unknown:
		opt.Target = config.WaOS_unknown
	default:
		fmt.Printf("unreachable: target: %s\n", c.String("target"))
		os.Exit(1)
	}

	return opt
}
