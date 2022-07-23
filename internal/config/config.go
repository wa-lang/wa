// 版权 @2019 凹语言 作者。保留所有权利。

package config

import (
	"os"
	"runtime"
)

// 字长和指针大小
type StdSizes struct {
	WordSize int64 // word size in bytes - must be >= 4 (32bits)
	MaxAlign int64 // maximum alignment in bytes - must be >= 1
}

// 通用配置信息
type Config struct {
	WaRoot   string   // 凹 程序根目录, src 目录下是包代码
	WaArch   string   // 目标 CPU
	WaOS     string   // 目标 OS
	WaSizes  StdSizes // 指针大小
	Optimize bool     // 是否优化
	Debug    bool     // 调试模式
}

func (p *Config) Clone() *Config {
	var q = *p
	return &q
}

func DefaultConfig() *Config {
	p := &Config{}

	if p.WaArch == "" {
		if s := os.Getenv("WAARCH"); s != "" {
			p.WaArch = s
		} else {
			p.WaArch = runtime.GOARCH
		}
	}
	if p.WaOS == "" {
		if s := os.Getenv("WAOS"); s != "" {
			p.WaOS = s
		} else {
			p.WaOS = runtime.GOOS
		}
	}
	if p.WaRoot == "" {
		if s := os.Getenv("WAROOT"); s != "" {
			p.WaRoot = s
		}
	}

	return p
}
