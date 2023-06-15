// 版权 @2023 凹语言 作者。保留所有权利。

package app

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/sys"
	"wa-lang.org/wa/internal/app/apputil"
	"wa-lang.org/wa/internal/app/waruntime"
	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/config"
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

	// 生成 wat 文件(main 函数为空)
	watOutput, err := compiler_wat.New().Compile(prog, "")
	if err != nil {
		return err
	}

	// 编译为 wasm
	if err = os.WriteFile("a.out.wat", []byte(watOutput), 0666); err != nil {
		return err
	}
	wasmBytes, stderr, err := apputil.RunWat2Wasm("a.out.wat")
	if err != nil {
		if s := sWithPrefix(string(stderr), "    "); s != "" {
			fmt.Println(s)
		}
		return err
	}
	// os.WriteFile("a.out.wasm", wasmBytes, 0666)

	// 构建 wasm 可执行实例
	// https://pkg.go.dev/github.com/tetratelabs/wazero@v1.0.0-pre.4#section-readme

	var r wazero.Runtime
	var ctx = context.Background()
	var stdoutBuffer = new(bytes.Buffer)
	var stderrBuffer = new(bytes.Buffer)
	{
		r = wazero.NewRuntime(ctx)
		defer r.Close(ctx)

		switch cfg.WaOS {
		case config.WaOS_arduino:
			if _, err = waruntime.ArduinoInstantiate(ctx, r); err != nil {
				if s := sWithPrefix(stderrBuffer.String(), "    "); s != "" {
					fmt.Println(s)
				}
				return err
			}
		case config.WaOS_chrome:
			if _, err = waruntime.ChromeInstantiate(ctx, r); err != nil {
				if s := sWithPrefix(stderrBuffer.String(), "    "); s != "" {
					fmt.Println(s)
				}
				return err
			}
		case config.WaOS_wasi:
			if _, err = waruntime.WasiInstantiate(ctx, r); err != nil {
				if s := sWithPrefix(stderrBuffer.String(), "    "); s != "" {
					fmt.Println(s)
				}
				return err
			}
		}
	}

	conf := wazero.NewModuleConfig().
		WithStdout(stdoutBuffer).
		WithStderr(stderrBuffer).
		WithStdin(os.Stdin).
		WithRandSource(rand.Reader).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithArgs("a.out.wasm").
		WithName("unittest")

	// 执行 init 函数
	compiled, err := r.CompileModule(ctx, wasmBytes)
	if err != nil {
		return err
	}

	wasmIns, err := r.InstantiateModule(ctx, compiled, conf)
	if err != nil {
		if s := sWithPrefix(stderrBuffer.String(), "    "); s != "" {
			fmt.Println(s)
		}
		return err
	}

	// 执行测试函数
	var firstError error
	for _, t := range mainPkg.TestInfo.Tests {
		stdoutBuffer.Reset()
		stderrBuffer.Reset()

		_, err := wasmIns.ExportedFunction(mainPkg.Pkg.Path() + "." + t.Name).Call(ctx)
		if err != nil {
			if s := sWithPrefix(stderrBuffer.String(), "    "); s != "" {
				fmt.Println(s)
			}
			return err
		}

		stdout := stdoutBuffer.Bytes()
		stderr := stderrBuffer.Bytes()

		stdout = bytes.TrimSpace(stdout)
		if t.Output != "" && t.Output == string(stdout) {
			continue
		}

		if err != nil {
			if firstError == nil {
				firstError = err
			}
			if _, ok := err.(*sys.ExitError); ok {
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
		output, err := compiler_wat.New().Compile(prog, t.Name)
		if err != nil {
			return err
		}

		if err = os.WriteFile("a.out.wat", []byte(output), 0666); err != nil {
			return err
		}

		stdout, stderr, err := apputil.RunWasmEx(cfg, "a.out.wat", appArgs...)

		stdout = bytes.TrimSpace(stdout)
		bOutputOK := t.Output == string(stdout)

		if err == nil && bOutputOK {
			continue
		}

		if err != nil {
			if firstError == nil {
				firstError = err
			}
			if _, ok := err.(*sys.ExitError); ok {
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

	return nil
}

func sWithPrefix(s, prefix string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, line := range lines {
		lines[i] = prefix + strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}
