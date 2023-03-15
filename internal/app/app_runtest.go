// 版权 @2023 凹语言 作者。保留所有权利。

package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/tetratelabs/wazero/sys"
	"wa-lang.org/wa/internal/app/apputil"
	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/loader"
)

func (p *App) RunTest(filename string, appArgs ...string) error {
	cfg := p.opt.Config()
	cfg.UnitTest = true
	prog, err := loader.LoadProgram(cfg, filename)
	if err != nil {
		return err
	}

	// 凹中文的源码启动函数为【启】，对应的wat函数名应当是"$0xe590af"
	main := "main"
	if strings.HasSuffix(filename, ".wz") {
		main = "$0xe590af"
	}

	output, err := compiler_wat.New().Compile(prog, main)
	if err != nil {
		return err
	}

	if err = os.WriteFile("a.out.wat", []byte(output), 0666); err != nil {
		return err
	}

	stdoutStderr, err := apputil.RunWasm(cfg, "a.out.wat", appArgs...)
	if err != nil {
		if len(stdoutStderr) > 0 {
			fmt.Println(string(stdoutStderr))
		}
		if exitErr, ok := err.(*sys.ExitError); ok {
			os.Exit(int(exitErr.ExitCode()))
		}
		fmt.Println(err)
	}
	if len(stdoutStderr) > 0 {
		fmt.Println(string(stdoutStderr))
	}
	return nil
}
