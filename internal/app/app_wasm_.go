// 版权 @2019 凹语言 作者。保留所有权利。

package app

import (
	"strings"

	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/loader"
)

func (p *App) WASM(filename string) ([]byte, error) {
	cfg := p.opt.Config()
	prog, err := loader.LoadProgram(cfg, filename)
	if err != nil {
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

	output, err := compiler_wat.New().Compile(prog, main)

	if err != nil {
		return nil, err
	}

	return []byte(output), nil
}
