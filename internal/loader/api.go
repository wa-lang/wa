// 版权 @2021 凹语言 作者。保留所有权利。

package loader

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/ssa"
	"wa-lang.org/wa/internal/token"
	"wa-lang.org/wa/internal/types"
)

// 程序对象
// 包含程序需要的全部信息
type Program struct {
	Cfg      *config.Config   // 配置信息
	Manifest *config.Manifest // 主包信息

	Fset *token.FileSet
	Pkgs map[string]*Package

	SSAProgram *ssa.Program
	SSAMainPkg *ssa.Package
}

// 单个包对象
type Package struct {
	Pkg          *types.Package // 类型检查后的包
	Info         *types.Info    // 包的类型检查信息
	Files        []*ast.File    // AST语法树
	WsFiles      []*WsFile      // 汇编代码
	WImportFiles []*WhostFile   // 宿主代码文件

	SSAPkg   *ssa.Package
	TestInfo *TestInfo
}

// 单元测试信息
type TestInfo struct {
	Files    []string // 测试文件
	Tests    []TestFuncInfo
	Benchs   []TestFuncInfo
	Examples []TestFuncInfo
}

// 测试函数信息
type TestFuncInfo struct {
	FuncPos     token.Pos // 函数位置
	Name        string    // 函数名, 不含包路径
	Output      string    // 期望输出, 为空表示不验证
	OutputPanic bool      // 异常信息
}

// 汇编代码文件
type WsFile struct {
	Name string // 文件名
	Code string // 汇编代码
}

// 宿主代码文件
type WhostFile struct {
	Name string // 文件名
	Code string // 宿主代码
}

// 加载程序
// 入口 appPath 是包对应目录的路径
func LoadProgram(cfg *config.Config, appPath string) (*Program, error) {
	return newLoader(cfg).LoadProgram(appPath)
}

// 加载单文件程序
// 入口 appPath 是包对应目录的路径
func LoadProgramFile(cfg *config.Config, filename string, src interface{}) (*Program, error) {
	return newLoader(cfg).LoadProgramFile(filename, src)
}

// 基于 VFS 加载程序
// 入口 pkgPath 是包路径, 必须是 vfs.App 子包
func LoadProgramVFS(vfs *config.PkgVFS, cfg *config.Config, pkgPath string) (*Program, error) {
	return newLoader(cfg).LoadProgramVFS(vfs, pkgPath)
}
