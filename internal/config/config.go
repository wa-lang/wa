// 版权 @2019 凹语言 作者。保留所有权利。

package config

import (
	"io/fs"
	"os"
	"runtime"
)

// 字长和指针大小
type StdSizes struct {
	WordSize int64 // word size in bytes - must be >= 4 (32bits)
	MaxAlign int64 // maximum alignment in bytes - must be >= 1
}

// 包虚拟文件系统
type PkgVFS struct {
	App    fs.FS // 当前工程的 src 目录, 导入路径去掉前缀对应目录
	Std    fs.FS // 标准库, 导入路径对应目录
	Vendor fs.FS // 第三方库, 导入路径目录
}

// 通用配置信息
type Config struct {
	WaRoot   string   // 凹 程序根目录, src 目录下是包代码, 为空时用内置标准库实现
	WaArch   string   // 目标 CPU
	WaOS     string   // 目标 OS
	WaSizes  StdSizes // 指针大小
	Optimize bool     // 是否优化
	Debug    bool     // 调试模式
	LDFlags           // 链接参数
}

// 链接参数
type LDFlags struct {
	StackSize int // 栈大小
	MaxMemory int // 最大内存
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
