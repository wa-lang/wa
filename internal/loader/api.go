// 版权 @2021 凹语言 作者。保留所有权利。

package loader

import (
	"github.com/wa-lang/wa/internal/ast"
	"github.com/wa-lang/wa/internal/config"
	"github.com/wa-lang/wa/internal/ssa"
	"github.com/wa-lang/wa/internal/token"
	"github.com/wa-lang/wa/internal/types"
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
	Pkg   *types.Package // 类型检查后的包
	Info  *types.Info    // 包的类型检查信息
	Files []*ast.File    // AST语法树

	SSAPkg *ssa.Package
}

// 加载程序
// 入口 appPath 是包对应目录的路径
func LoadProgram(cfg *config.Config, appPath string) (*Program, error) {
	return newLoader(cfg).LoadProgram(appPath)
}
