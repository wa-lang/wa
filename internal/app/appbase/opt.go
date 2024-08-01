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
	WaBackend    string
	BuilgTags    []string
	TargetArch   string
	TargetOS     string
	LD_StackSize int
	LD_MaxMemory int
}

func (opt *Option) Config() *config.Config {
	cfg := config.DefaultConfig()

	if opt.Debug {
		cfg.Debug = true
	}
	if opt.WaBackend != "" {
		cfg.WaBackend = opt.WaBackend
	}
	if len(opt.BuilgTags) > 0 {
		cfg.BuilgTags = append(cfg.BuilgTags, opt.BuilgTags...)
	}
	if opt.TargetArch != "" {
		cfg.WaArch = opt.TargetArch
	}
	if opt.TargetOS != "" {
		cfg.WaOS = opt.TargetOS
	}
	if opt.LD_StackSize != 0 {
		cfg.LDFlags.StackSize = opt.LD_StackSize
	}
	if opt.LD_MaxMemory != 0 {
		cfg.LDFlags.MaxMemory = opt.LD_MaxMemory
	}

	switch cfg.WaArch {
	case "wasm":
		cfg.WaSizes.MaxAlign = 8
		cfg.WaSizes.WordSize = 4
	case "wasm64":
		cfg.WaSizes.MaxAlign = 8
		cfg.WaSizes.WordSize = 8
	default:
		panic("unknown WaArch: " + cfg.WaArch)
	}

	return cfg
}

// 构建命令行程序对象
func (opt *Option) Adjust() {
	if opt.TargetOS == "" {
		opt.TargetOS = config.WaOS_Default
	}
	if opt.TargetArch == "" {
		opt.TargetArch = config.WaArch_Default
	}
}

func BuildOptions(c *cli.Context, waBackend ...string) *Option {
	opt := &Option{
		Debug:        c.Bool("debug"),
		WaBackend:    config.WaBackend_Default,
		BuilgTags:    strings.Fields(c.String("tags")),
		LD_StackSize: c.Int("ld-stack-size"),
		LD_MaxMemory: c.Int("ld-max-memory"),
	}

	opt.TargetArch = "wasm"
	if len(waBackend) > 0 {
		opt.WaBackend = waBackend[0]
	}

	if target := c.String("target"); !config.CheckWaOS(target) {
		fmt.Printf("unknown target: %s\n", c.String("target"))
		os.Exit(1)
	}

	switch c.String("target") {
	case "", "wa", "walang":
		opt.TargetOS = config.WaOS_Default
	case config.WaOS_wasi:
		opt.TargetOS = config.WaOS_wasi
	case config.WaOS_unknown:
		opt.TargetOS = config.WaOS_unknown
	case config.WaOS_js:
		opt.TargetOS = config.WaOS_js
	default:
		fmt.Printf("unreachable: target: %s\n", c.String("target"))
		os.Exit(1)
	}

	opt.Adjust()
	return opt
}
