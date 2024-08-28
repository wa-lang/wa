// 版权 @2019 凹语言 作者。保留所有权利。

package config

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"

	"wa-lang.org/wa/internal/version"
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
	Target    string   // 目标平台
	WatOutput string   // 输出的 wat 文件路径
	WaBackend string   // 编译器后端
	WaSizes   StdSizes // 指针大小
	BuilgTags []string // 条件编译的标志
	UnitTest  bool     // 单元测试模式
	Optimize  bool     // 是否优化
	Debug     bool     // 调试模式
	LDFlags            // 链接参数
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

	if p.WaBackend == "" {
		p.WaBackend = WaBackend_Default
	}

	return p
}

// Waroot 是否有效
// 需要保持和 waroot 处理一致
func isWarootValid() (warootDir string, ok bool) {
	if s, _ := os.UserHomeDir(); s != "" {
		warootDir = filepath.Join(s, "wa")
	}

	d, err := os.ReadFile(filepath.Join(warootDir, "VERSION"))
	if err != nil {
		return "", false
	}

	ver := string(bytes.TrimSpace(d))
	if ver != version.Version {
		return "", false
	}

	return warootDir, true
}
