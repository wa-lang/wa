// 版权 @2019 凹语言 作者。保留所有权利。

package app

import (
	"os/exec"
	"runtime"

	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/logger"
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
		p.opt.TargetOS = config.WaOS_Walang
	}
	if p.opt.TargetArch == "" {
		p.opt.TargetArch = config.WaArch_Wasm
	}

	return p
}

func (p *App) GetConfig() *config.Config {
	return p.opt.Config()
}
