// 版权 @2023 凹语言 作者。保留所有权利。

package apptest

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/wat/watutil"
	"wa-lang.org/wa/internal/wazero"
	"wa-lang.org/wa/waroot/src"
)

var CmdTest = &cli.Command{
	Name:  "test",
	Usage: "test Wa packages",
	Flags: []cli.Flag{
		appbase.MakeFlag_target(),
		appbase.MakeFlag_tags(),
		&cli.StringFlag{
			Name:  "run",
			Usage: "set run func file name pattern",
			Value: "",
		},
	},
	Action: func(c *cli.Context) error {
		opt := appbase.BuildOptions(c)

		var pkgpath = "."
		var runPattern = c.String("run")
		var appArgs []string

		if c.Args().Len() > 0 {
			pkgpath = c.Args().First()
		}
		if c.Args().Len() > 1 {
			appArgs = c.Args().Slice()[1:]
		}

		RunTest(opt.Config(), pkgpath, runPattern, appArgs...)
		return nil
	},
}

func RunTest(cfg *config.Config, pkgpath, runPattern string, appArgs ...string) {
	var pkgList = []string{pkgpath}
	if pkgpath == "std" {
		pkgList = src.GetStdTestPkgList()
	}
	for _, p := range pkgList {
		runTest(cfg, p, runPattern, appArgs...)
	}
}

func runTest(cfg *config.Config, pkgpath, runPattern string, appArgs ...string) {
	cfg.UnitTest = true
	prog, err := loader.LoadProgram(cfg, pkgpath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	startTime := time.Now()
	mainPkg := prog.Pkgs[prog.Manifest.MainPkg]
	wasmName := "unittest://" + pkgpath
	wasmArgs := []string{}

	if len(mainPkg.TestInfo.Files) == 0 {
		fmt.Printf("?    %s [no test files]\n", prog.Manifest.MainPkg)
		return
	}

	// 生成 wat 文件(main 函数为空)
	watOutput, err := compiler_wat.New().Compile(prog)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// wat 落盘(仅用于调试)
	os.WriteFile("a.out.wat", []byte(watOutput), 0666)

	// 编译为 wasm
	wasmBytes, err := watutil.Wat2Wasm("a.out.wat", []byte(watOutput))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m, err := wazero.BuildModule(wasmName, wasmBytes, wasmArgs...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer m.Close()

	// 执行测试函数
	var firstError error
	for i := 0; i < len(mainPkg.TestInfo.Tests); i++ {
		t := mainPkg.TestInfo.Tests[i]

		if runPattern != "" {
			matched, err := filepath.Match(runPattern, t.Name)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if !matched {
				continue
			}
		}

		tFuncName := mainPkg.Pkg.Path() + "." + t.Name
		tFuncName = strings.ReplaceAll(tFuncName, "/", "$")
		_, stdout, stderr, err := m.RunFunc(tFuncName)
		if t.OutputPanic {
			stdout = fmtGotOutput(stdout)
			expect, got := t.Output, string(stdout)

			if exitCode, _ := wazero.AsExitError(err); exitCode == 0 {
				fmt.Printf("---- %s.%s\n", prog.Manifest.MainPkg, t.Name)
				fmt.Printf("    expect panic, got = nil\n")

				if firstError == nil {
					firstError = fmt.Errorf("expect(panic) = %q, got = %q", expect, "nil")
				}
				continue
			}

			if !strings.HasPrefix(got, "panic: "+expect) { // panic: ${expect} (pos)
				fmt.Printf("---- %s.%s\n", prog.Manifest.MainPkg, t.Name)
				fmt.Printf("    expect(panic) = %q, got = %q\n", expect, got)

				if firstError == nil {
					firstError = fmt.Errorf("expect(panic) = %q, got = %q", expect, got)
				}
			}

			// 重新加载
			{
				m, err = wazero.BuildModule(wasmName, wasmBytes, wasmArgs...)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				// 临时方案: defer 太多
				defer m.Close()
			}

			continue
		}
		if err != nil {
			if len(stdout) > 0 {
				if s := sWithPrefix(string(stdout), "    "); s != "" {
					fmt.Println(s)
				}
			}
			if len(stderr) > 0 {
				if s := sWithPrefix(string(stderr), "    "); s != "" {
					fmt.Println(s)
				}
			} else {
				if s := sWithPrefix(err.Error(), "    "); s != "" {
					fmt.Println(s)
				}
			}

			os.Exit(1)
		}

		stdout = fmtGotOutput(stdout)
		if t.Output != "" && t.Output == string(stdout) {
			continue
		}

		if err != nil {
			if firstError == nil {
				firstError = err
			}
			if _, ok := wazero.AsExitError(err); ok {
				fmt.Printf("---- %s.%s\n", prog.Manifest.MainPkg, t.Name)
				if s := sWithPrefix(string(stdout), "    "); s != "" {
					fmt.Println(s)
				}
				if s := sWithPrefix(string(stderr), "    "); s != "" {
					fmt.Println(s)
				}
			} else {
				fmt.Println(err)
			}
		}

		if t.Output != "" {
			if expect, got := t.Output, string(stdout); expect != got {
				if firstError == nil {
					firstError = fmt.Errorf("expect = %q, got = %q", expect, got)
				}
				fmt.Printf("---- %s.%s\n", prog.Manifest.MainPkg, t.Name)
				fmt.Printf("    expect = %q, got = %q\n", expect, got)
			}
		}
	}

	for i := 0; i < len(mainPkg.TestInfo.Examples); i++ {
		t := mainPkg.TestInfo.Examples[i]

		if runPattern != "" {
			matched, err := filepath.Match(runPattern, t.Name)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if !matched {
				continue
			}
		}

		tFuncName := mainPkg.Pkg.Path() + "." + t.Name
		tFuncName = strings.ReplaceAll(tFuncName, "/", "$")
		_, stdout, stderr, err := m.RunFunc(tFuncName)
		if t.OutputPanic {
			stdout = fmtGotOutput(stdout)
			expect, got := t.Output, string(stdout)

			if exitCode, _ := wazero.AsExitError(err); exitCode == 0 {
				fmt.Printf("---- %s.%s\n", prog.Manifest.MainPkg, t.Name)
				fmt.Printf("    expect panic, got = nil\n")

				if firstError == nil {
					firstError = fmt.Errorf("expect panic, got = nil")
				}
				continue
			}

			if !strings.HasPrefix(got, "panic: "+expect) { // panic: ${expect} (pos)
				fmt.Printf("---- %s.%s\n", prog.Manifest.MainPkg, t.Name)
				fmt.Printf("    expect(panic) = %q, got = %q\n", expect, got)

				if firstError == nil {
					firstError = fmt.Errorf("expect(panic) = %q, got = %q", expect, got)
				}
			}

			// 重新加载
			{
				m, err = wazero.BuildModule(wasmName, wasmBytes, wasmArgs...)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				// 临时方案: defer 太多
				defer m.Close()
			}

			continue
		}
		if err != nil {
			if len(stdout) > 0 {
				if s := sWithPrefix(string(stdout), "    "); s != "" {
					fmt.Println(s)
				}
			}
			if len(stderr) > 0 {
				if s := sWithPrefix(string(stderr), "    "); s != "" {
					fmt.Println(s)
				}
			}

			os.Exit(1)
		}

		stdout = fmtGotOutput(stdout)
		if t.Output != "" && t.Output == string(stdout) {
			continue
		}

		if err != nil {
			if firstError == nil {
				firstError = err
			}
			if _, ok := wazero.AsExitError(err); ok {
				fmt.Printf("---- %s.%s\n", prog.Manifest.MainPkg, t.Name)
				if s := sWithPrefix(string(stdout), "    "); s != "" {
					fmt.Println(s)
				}
				if s := sWithPrefix(string(stderr), "    "); s != "" {
					fmt.Println(s)
				}
			} else {
				fmt.Println(err)
			}
		}

		if t.Output != "" {
			if expect, got := t.Output, string(stdout); expect != got {
				if firstError == nil {
					firstError = fmt.Errorf("expect = %q, got = %q", expect, got)
				}
				fmt.Printf("---- %s.%s\n", prog.Manifest.MainPkg, t.Name)
				fmt.Printf("    expect = %q, got = %q\n", expect, got)
			}
		}
	}
	if firstError != nil {
		fmt.Printf("FAIL %s %v\n", prog.Manifest.MainPkg, time.Now().Sub(startTime).Round(time.Millisecond))
		os.Exit(1)
	}

	fmt.Printf("ok   %s %v\n", prog.Manifest.MainPkg, time.Now().Sub(startTime).Round(time.Millisecond))

	return
}

func sWithPrefix(s, prefix string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, line := range lines {
		lines[i] = prefix + strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}

func fmtGotOutput(stdout []byte) []byte {
	stdout = bytes.TrimSpace(stdout)
	lines := bytes.Split(stdout, []byte("\n"))
	for i, s := range lines {
		lines[i] = bytes.TrimSpace(s)
	}
	return bytes.Join(lines, []byte("\n"))
}
