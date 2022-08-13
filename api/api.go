// 版权 @2022 凹语言 作者。保留所有权利。

// 凹语言™ 功能 API 包。
package api

import (
	"io/fs"

	"github.com/wa-lang/wa/internal/config"
	"github.com/wa-lang/wa/internal/loader"
)

// 配置参数, 包含文件系统和 OS 等信息
type Config = config.Config

// 模块元信息, 主包路径
type Manifest = config.Manifest

// 模块元信息中的包信息
type Manifest_package = config.Manifest_package

// 程序对象, 包含全量的 AST 和 SSA 信息, 经过语义检查
type Program = loader.Program

// 指针和整数大小
type StdSize = config.StdSizes

// 加载 WaModFile 文件
// 如果 vfs 为空则从本地文件系统读取
func LoadManifest(vfs fs.FS, appPath string) (p *Manifest, err error) {
	return config.LoadManifest(vfs, appPath)
}

// 加载程序
// 入口 appPath 是包对应目录的路径
func LoadProgram(cfg *config.Config, appPath string) (*Program, error) {
	return loader.LoadProgram(cfg, appPath)
}

// TODO: Program 编译到不同后端的函数
