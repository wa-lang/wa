// 版权 @2019 凹语言 作者。保留所有权利。

package apputil

import (
	"bytes"
	"context"
	"crypto/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/sys"

	"wa-lang.org/wa/internal/app/waruntime"
	"wa-lang.org/wa/internal/config"
	wabt_wasm "wa-lang.org/wabt-go/wabt-wasm"
)

func getWatAbsDir(filename string) string {
	dir := filepath.Dir(filename)
	abs, err := filepath.Abs(dir)
	if err != nil {
		return filepath.ToSlash(dir)
	}
	return filepath.ToSlash(abs)
}

func RunWasm(cfg *config.Config, filename string, wasmArgs ...string) (stdoutStderr []byte, err error) {
	stdout, strerr, err := runWasm(cfg, filename, wasmArgs...)
	stdoutStderr = append(stdout, strerr...)
	return
}

func RunWat2Wasm(filename string) (stdout, stderr []byte, err error) {
	watDir := getWatAbsDir(filename)
	stdout, stderr, err = runWat2Wasm(watDir, "/"+filepath.Base(filename), "--output=-")
	return
}

func runWasm(cfg *config.Config, filename string, wasmArgs ...string) (stdout, stderr []byte, err error) {
	dst := filename
	if strings.HasSuffix(filename, ".wat") {
		dst = filename[:len(filename)-len(".wat")] + ".wasm"
		defer func() {
			if err == nil {
				os.Remove(dst)
			}
		}()

		// 将目录 wat 文件目录映射到根目录
		watDir := getWatAbsDir(filename)
		if stdout, stderr, err = runWat2Wasm(watDir, "/"+filepath.Base(filename), "--output=-"); err != nil {
			return stdout, stderr, err
		}

		os.WriteFile(dst, stdout, 0666)
	}

	wasmBytes, err := os.ReadFile(dst)
	if err != nil {
		return nil, nil, err
	}

	wasmExe := filepath.Base(filename)

	stdoutBuffer := new(bytes.Buffer)
	stderrBuffer := new(bytes.Buffer)

	conf := wazero.NewModuleConfig().
		WithStdout(stdoutBuffer).
		WithStderr(stderrBuffer).
		WithStdin(os.Stdin).
		WithRandSource(rand.Reader).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithArgs(append([]string{wasmExe}, wasmArgs...)...)

	// TODO: Windows 可能导致异常, 临时屏蔽
	if runtime.GOOS != "windows" {
		for _, s := range os.Environ() {
			var key, value string
			if kv := strings.Split(s, "="); len(kv) >= 2 {
				key = kv[0]
				value = kv[1]
			} else if len(kv) >= 1 {
				key = kv[0]
			}
			conf = conf.WithEnv(key, value)
		}
	}

	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	code, err := r.CompileModule(ctx, wasmBytes)
	if err != nil {
		return stdoutBuffer.Bytes(), stderrBuffer.Bytes(), err
	}

	switch cfg.WaOS {
	case config.WaOS_arduino:
		if _, err = waruntime.ArduinoInstantiate(ctx, r); err != nil {
			return stdoutBuffer.Bytes(), stderrBuffer.Bytes(), err
		}
	case config.WaOS_chrome:
		if _, err = waruntime.ChromeInstantiate(ctx, r); err != nil {
			return stdoutBuffer.Bytes(), stderrBuffer.Bytes(), err
		}
	case config.WaOS_wasi:
		if _, err = waruntime.WasiInstantiate(ctx, r); err != nil {
			return stdoutBuffer.Bytes(), stderrBuffer.Bytes(), err
		}
	}

	_, err = r.InstantiateModule(ctx, code, conf)
	if err != nil {
		if exitErr, ok := err.(*sys.ExitError); ok {
			if exitErr.ExitCode() == 0 {
				return stdoutBuffer.Bytes(), stderrBuffer.Bytes(), nil
			}
		}
		return stdoutBuffer.Bytes(), stderrBuffer.Bytes(), err
	}

	return stdoutBuffer.Bytes(), stderrBuffer.Bytes(), nil
}

func runWat2Wasm(dir string, args ...string) (stdout, stderr []byte, err error) {
	stdoutBuffer := new(bytes.Buffer)
	stderrBuffer := new(bytes.Buffer)

	conf := wazero.NewModuleConfig().
		WithFS(os.DirFS(dir)).
		WithStdout(stdoutBuffer).
		WithStderr(stderrBuffer).
		WithStdin(os.Stdin).
		WithRandSource(rand.Reader).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithArgs(append([]string{"wat2wasm.wasm"}, args...)...)

	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	code, err := r.CompileModule(ctx, []byte(wabt_wasm.LoadWat2Wasm()))
	if err != nil {
		return stdoutBuffer.Bytes(), stderrBuffer.Bytes(), err
	}

	if _, err = waruntime.WasiInstantiate(ctx, r); err != nil {
		return stdoutBuffer.Bytes(), stderrBuffer.Bytes(), err
	}

	_, err = r.InstantiateModule(ctx, code, conf)
	if err != nil {
		if exitErr, ok := err.(*sys.ExitError); ok {
			if exitErr.ExitCode() == 0 {
				return stdoutBuffer.Bytes(), stderrBuffer.Bytes(), nil
			}
		}
		return stdoutBuffer.Bytes(), stderrBuffer.Bytes(), err
	}

	return stdoutBuffer.Bytes(), stderrBuffer.Bytes(), nil
}
