// 版权 @2019 凹语言 作者。保留所有权利。

package app

import "github.com/wa-lang/wa/internal/config"

// 命令行选项
type Option struct {
	Debug        bool
	TargetArch   string
	TargetOS     string
	Clang        string
	Llc          string
	LD_StackSize int
	LD_MaxMemory int
}

func (opt *Option) Config() *config.Config {
	cfg := config.DefaultConfig()

	if opt.Debug {
		cfg.Debug = true
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
		cfg.WaSizes.MaxAlign = 4
		cfg.WaSizes.WordSize = 4
	case "amd64":
		cfg.WaSizes.MaxAlign = 8
		cfg.WaSizes.WordSize = 8
	case "arm64":
		cfg.WaSizes.MaxAlign = 8
		cfg.WaSizes.WordSize = 8
	default:
		panic("todo")
	}

	return cfg
}
