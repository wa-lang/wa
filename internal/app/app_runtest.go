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

func (p *App) RunTest(pkgpath string, appArgs ...string) error {
	cfg := p.opt.Config()
	cfg.UnitTest = true
	prog, err := loader.LoadProgram(cfg, pkgpath)
	if err != nil {
		return err
	}

	startTime := time.Now()
	mainPkg := prog.Pkgs[prog.Manifest.MainPkg]

	if len(mainPkg.TestInfo.Files) == 0 {
		fmt.Printf("?    %s [no test files]\n", prog.Manifest.MainPkg)
		return nil
	}
	if len(mainPkg.TestInfo.Funcs) == 0 {
		fmt.Printf("ok   %s %v\n", prog.Manifest.MainPkg, time.Now().Sub(startTime))
		return nil
	}

	var lastError error
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

		if err != nil {
			lastError = err
			if _, ok := err.(*sys.ExitError); ok {
				fmt.Printf("---- %s.%s\n", prog.Manifest.MainPkg, main)
				if s := sWithPrefix(string(stdoutStderr), "    "); s != "" {
					fmt.Println(s)
				}
			} else {
				fmt.Println(err)
			}
		}
	}

	if lastError != nil {
		fmt.Printf("FAIL %s %v\n", prog.Manifest.MainPkg, time.Now().Sub(startTime).Round(time.Microsecond))
		os.Exit(1)
	}
	fmt.Printf("ok   %s %v\n", prog.Manifest.MainPkg, time.Now().Sub(startTime).Round(time.Microsecond))

	return nil
}

func sWithPrefix(s, prefix string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, line := range lines {
		lines[i] = prefix + strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}
