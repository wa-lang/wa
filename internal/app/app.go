// 版权 @2019 凹语言 作者。保留所有权利。

package app

import (
	"fmt"
	"os/exec"
	"runtime"

	"wa-lang.org/wa/internal/app/appfmt"
	"wa-lang.org/wa/internal/app/appinit"
	"wa-lang.org/wa/internal/app/applex"
	"wa-lang.org/wa/internal/app/appplay"
	"wa-lang.org/wa/internal/app/apptest"
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/backends/compiler_c"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/token"
)

// 命令行程序对象
type App struct {
	opt  Option
	path string
	src  string
}

// 构建命令行程序对象
func NewApp(opt *Option) *App {
	logger.Tracef(&config.EnableTrace_app, "opt: %+v", opt)

	p := &App{}
	if opt != nil {
		p.opt = *opt
	}
	if p.opt.Clang == "" {
		if runtime.GOOS == "windows" {
			p.opt.Clang, _ = exec.LookPath("clang.exe")
		} else {
			p.opt.Clang, _ = exec.LookPath("clang")
		}
		if p.opt.Clang == "" {
			p.opt.Clang = "clang"
		}
	}
	if p.opt.Llc == "" {
		if runtime.GOOS == "windows" {
			p.opt.Llc, _ = exec.LookPath("llc.exe")
		} else {
			p.opt.Llc, _ = exec.LookPath("llc")
		}
		if p.opt.Llc == "" {
			p.opt.Llc = "llc"
		}
	}
	if p.opt.TargetOS == "" {
		p.opt.TargetOS = config.WaOS_Default
	}
	if p.opt.TargetArch == "" {
		p.opt.TargetArch = config.WaArch_Default
	}

	return p
}

func (p *App) GetConfig() *config.Config {
	return p.opt.Config()
}

func (p *App) RunTest(pkgpath string, appArgs ...string) {
	apptest.RunTest(p.opt.Config(), pkgpath, appArgs...)
}

func (p *App) AST(filename string) error {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(nil, fset, filename, nil, 0)
	if err != nil {
		return err
	}

	return ast.Print(fset, f)
}

func (p *App) CIR(filename string) error {
	cfg := p.opt.Config()
	prog, err := loader.LoadProgram(cfg, filename)
	if err != nil {
		return err
	}

	var c compiler_c.CompilerC
	c.CompilePackage(prog.SSAMainPkg)
	fmt.Println(c.String())

	return nil
}

func (p *App) Fmt(path string) error {
	return appfmt.Fmt(path)
}

func (p *App) Playground(addr string) error {
	return appplay.RunPlayground(addr)
}

func (p *App) InitApp(name, pkgpath string, update bool) error {
	return appinit.InitApp(name, pkgpath, update)
}

func (p *App) Lex(filename string) error {
	return applex.Lex(filename)
}
