// 版权 @2023 凹语言 作者。保留所有权利。

package apptest

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/wabt"
	"wa-lang.org/wa/internal/wazero"
	"wa-lang.org/wa/waroot"
)

func RunTest(cfg *config.Config, pkgpath string, appArgs ...string) {
	var pkgList = []string{pkgpath}
	if pkgpath == "std" {
		pkgList = waroot.GetStdPkgList()
	}
	for _, p := range pkgList {
		runTest(cfg, p, appArgs...)
	}
}

func runTest(cfg *config.Config, pkgpath string, appArgs ...string) {
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
	watOutput, err := compiler_wat.New().Compile(prog, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// wat 落盘(仅用于调试)
	os.WriteFile("a.out.wat", []byte(watOutput), 0666)

	// 编译为 wasm
	wasmBytes, err := wabt.Wat2Wasm([]byte(watOutput))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m, err := wazero.BuildModule(cfg, wasmName, wasmBytes, wasmArgs...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer m.Close()

	// 执行测试函数
	var firstError error
	for _, t := range mainPkg.TestInfo.Tests {
		_, stdout, stderr, err := m.RunFunc(mainPkg.Pkg.Path() + "." + t.Name)
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

		stdout = bytes.TrimSpace(stdout)
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
	for _, t := range mainPkg.TestInfo.Examples {
		_, stdout, stderr, err := m.RunFunc(mainPkg.Pkg.Path() + "." + t.Name)
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

		stdout = bytes.TrimSpace(stdout)
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
