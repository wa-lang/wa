// 版权 @2022 凹语言 作者。保留所有权利。

// 凹语言™ 功能 API 包。
package api

import (
	"io/fs"
	"strings"

	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/format"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/logger"
)

// 调试参数
var (
	FlagDebugMode = &config.DebugMode

	FlagEnableTrace_api      = &config.EnableTrace_api
	FlagEnableTrace_app      = &config.EnableTrace_app
	FlagEnableTrace_compiler = &config.EnableTrace_compiler
	FlagEnableTrace_loader   = &config.EnableTrace_loader
)

// 配置参数, 包含文件系统和 OS 等信息
type Config = config.Config

// 模块元信息, 主包路径
type Manifest = config.Manifest

// 模块元信息中的包信息
type Manifest_package = config.Manifest_package

// 程序对象, 包含全量的 AST 和 SSA 信息, 经过语义检查
type Program = loader.Program

// 包虚拟文件系统
type PkgVFS = config.PkgVFS

// 指针和整数大小
type StdSize = config.StdSizes

// 默认配置
func DefaultConfig() *Config {
	return config.DefaultConfig()
}

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

// 加载单文件程序
// 入口 appPath 是包对应目录的路径
func LoadProgramFile(cfg *config.Config, filename string, src interface{}) (*Program, error) {
	return loader.LoadProgramFile(cfg, filename, src)
}

// 基于 VFS 加载程序
// 入口 pkgPath 是包路径, 必须是 vfs.App 子包
func LoadProgramVFS(vfs *config.PkgVFS, cfg *config.Config, pkgPath string) (*Program, error) {
	return loader.LoadProgramVFS(vfs, cfg, pkgPath)
}

// 构建 wat 目标
func BuildFile(cfg *config.Config, filename string, src interface{}) (wat []byte, err error) {
	prog, err := LoadProgramFile(cfg, filename, src)
	if err != nil || prog == nil {
		logger.Tracef(&config.EnableTrace_api, "LoadProgramVFS failed, err = %v", err)
		return nil, err
	}

	// 凹中文的源码启动函数为【启】，对应的wat函数名应当是"$0xe590af"
	main := "main"
	if strings.HasSuffix(filename, ".wz") {
		main = "$0xe590af"
	}

	// 如果是运行整个package，则判断主包里是否有名为【启】的函数，如果有，则将其作为启动函数
	if filename == "." {
		for k := range prog.SSAMainPkg.Members {
			if k == "启" && prog.SSAMainPkg.Members[k].Type().Underlying().String() == "func()" {
				main = "$0xe590af"
			}
		}
	}

	watOut, err := compiler_wat.New().Compile(prog, main)
	return []byte(watOut), err
}

// 构建 wat 目标
func BuildVFS(cfg *config.Config, vfs *config.PkgVFS, appPkg string) (wat []byte, err error) {
	prog, err := LoadProgramVFS(vfs, cfg, appPkg)
	if err != nil || prog == nil {
		logger.Tracef(&config.EnableTrace_api, "LoadProgramVFS failed, err = %v", err)
		return nil, err
	}

	watOut, err := compiler_wat.New().Compile(prog, "main")
	return []byte(watOut), err
}

// 格式化代码
func FormatCode(filename, code string) (string, error) {
	data, _, err := format.File(nil, filename, code)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
