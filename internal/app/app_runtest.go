// 版权 @2023 凹语言 作者。保留所有权利。

package app

import (
	"fmt"
	"os"
	"strings"
	"time"

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

	startTime := time.Now()
	mainPkg := prog.Pkgs[prog.Manifest.MainPkg]

	// TODO: 测试指定路径的包

	if len(mainPkg.TestInfo.Files) == 0 {
		fmt.Printf("?    %s [no test files]\n", prog.Manifest.MainPkg)
		return nil
	}
	if len(mainPkg.TestInfo.Funcs) == 0 {
		fmt.Printf("ok   %s %v\n", prog.Manifest.MainPkg, time.Now().Sub(startTime))
		return nil
	}

	for _, main := range mainPkg.TestInfo.Funcs {
		output, err := compiler_wat.New().Compile(prog, main)
		if err != nil {
			return err
		}

		if err = os.WriteFile("a.out.wat", []byte(output), 0666); err != nil {
			return err
		}

		stdoutStderr, err := apputil.RunWasm(cfg, "a.out.wat", appArgs...)
		if err == nil {
			continue
		}

		if exitErr, ok := err.(*sys.ExitError); ok {
			fmt.Printf("---- %s.%s [%v]\n", prog.Manifest.MainPkg, main, time.Now().Sub(startTime))
			if s := sWithPrefix(string(stdoutStderr), "    "); s != "" {
				fmt.Println(s)
			}
			os.Exit(int(exitErr.ExitCode()))
		}
		fmt.Println(err)
	}

	fmt.Printf("ok   %s %v\n", prog.Manifest.MainPkg, time.Now().Sub(startTime))
	return nil
}

func sWithPrefix(s, prefix string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, line := range lines {
		lines[i] = prefix + strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}
