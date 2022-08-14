// 版权 @2022 凹语言 作者。保留所有权利。

// 凹语言™ 功能 API 包。
package api

import (
	"io/fs"
	"os"
	"os/exec"
	"runtime"

	"github.com/wa-lang/wa/internal/backends/compiler_ll"
	"github.com/wa-lang/wa/internal/backends/compiler_ll/builtin"
	"github.com/wa-lang/wa/internal/config"
	"github.com/wa-lang/wa/internal/loader"
	"github.com/wa-lang/wa/internal/logger"
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

// 基于 VFS 加载程序
// 入口 pkgPath 是包路径, 必须是 vfs.App 子包
func LoadProgramVFS(vfs *config.PkgVFS, cfg *config.Config, pkgPath string) (*Program, error) {
	return loader.LoadProgramVFS(vfs, cfg, pkgPath)
}

// 执行 vfs 中的程序
func RunVFS(vfs *config.PkgVFS, appPkg string, arg ...string) (output []byte, err error) {
	cfg := config.DefaultConfig()
	prog, err := LoadProgramVFS(vfs, cfg, appPkg)
	if err != nil || prog == nil {
		logger.Tracef(&config.EnableTrace_api, "LoadProgramVFS failed, err = %v", err)
		return nil, err
	}

	a_out := "a.out"
	if runtime.GOOS == "windows" {
		a_out = "a.exe"
	}
	output, err = buildProgram(prog, a_out)
	if err != nil {
		logger.Tracef(&config.EnableTrace_api, "buildProgram failed, err = %v", err)
		logger.Tracef(&config.EnableTrace_api, "buildProgram failed, output = %v", string(output))
		return nil, err
	}

	output, err = exec.Command(a_out, arg...).CombinedOutput()
	if err != nil {
		return output, err
	}

	return output, nil
}

func buildProgram(prog *Program, outfile string) (output []byte, err error) {
	llOutput, err := compiler_ll.New().Compile(prog)
	if err != nil {
		return nil, err
	}

	const (
		_a_out_ll     = "_a.out.ll"
		_a_builtin_ll = "_a_builtin.out.ll"
	)

	defer func() {
		if err == nil {
			os.Remove(_a_out_ll)
			os.Remove(_a_builtin_ll)
		}
	}()

	llButiltin := builtin.GetBuiltinLL(prog.Cfg.WaOS, prog.Cfg.WaArch)
	if err := os.WriteFile(_a_builtin_ll, []byte(llButiltin), 0666); err != nil {
		return nil, err
	}

	if err := os.WriteFile(_a_out_ll, []byte(llOutput), 0666); err != nil {
		return nil, err
	}

	args := []string{
		"-Wno-override-module", "-o", outfile,
		_a_out_ll, _a_builtin_ll,
	}

	clangExe := "clang"
	if runtime.GOOS == "windows" {
		clangExe += ".exe"
	}
	clangPath, err := exec.LookPath(clangExe)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(clangPath, args...)
	data, err := cmd.CombinedOutput()
	return data, err
}

// TODO: 解析 ast/语义/SSA 分阶段解析, 放到 Program 中
// TODO: Program 编译到不同后端的函数
