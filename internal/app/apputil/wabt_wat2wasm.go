// 版权 @2019 凹语言 作者。保留所有权利。

package apputil

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/sys"

	"wa-lang.org/wa/internal/app/waruntime"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wabt-go"
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
	stdout, stderr, err := runWasm(cfg, filename, wasmArgs...)
	stdoutStderr = append(stdout, stderr...)
	return
}

func RunWasmEx(cfg *config.Config, filename string, wasmArgs ...string) (stdout, stderr []byte, err error) {
	stdout, stderr, err = runWasm(cfg, filename, wasmArgs...)
	return
}

func RunWat2Wasm(filename string) (stdout, stderr []byte, err error) {
	watDir := getWatAbsDir(filename)
	stdout, stderr, err = xRunWat2Wasm_exe(watDir, filepath.Base(filename), "--output=-")
	return
}

func runWasm(cfg *config.Config, filename string, wasmArgs ...string) (stdout, stderr []byte, err error) {
	dst := filename
	if strings.HasSuffix(filename, ".wat") {
		dst = filename + ".wasm"
		defer func() {
			if err == nil {
				os.Remove(dst)
			}
		}()

		if stdout, stderr, err = xRunWat2Wasm_exe("", filename, "--output=-"); err != nil {
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

func xRunWat2Wasm_wasm(dir string, args ...string) (stdout, stderr []byte, err error) {
	if true {
		panic("disabled")
	}

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

func xInstallWat2wasm(path string) error {
	if path == "" {
		path, _ = os.Getwd()
	}

	var exePath string
	if isDir(path) {
		if runtime.GOOS == "windows" {
			exePath = filepath.Join(path, "wat2wasm.exe")
		} else {
			exePath = filepath.Join(path, "wat2wasm")
		}
	} else {
		exePath = path
	}

	if err := os.MkdirAll(filepath.Dir(exePath), 0777); err != nil {
		logger.Tracef(&config.EnableTrace_app, "install wat2wasm failed: %+v", err)
	}
	if err := os.WriteFile(exePath, wabt.LoadWat2Wasm(), 0777); err != nil {
		logger.Tracef(&config.EnableTrace_app, "install wat2wasm failed: %+v", err)
		return err
	}

	return nil
}

var muRunWat2Wasm sync.Mutex
var wat2wasmPath string

func xRunWat2Wasm_exe(_ string, args ...string) (stdout, stderr []byte, err error) {
	muRunWat2Wasm.Lock()
	defer muRunWat2Wasm.Unlock()

	if wat2wasmPath == "" {
		logger.Tracef(&config.EnableTrace_app, "wat2wasm not found")
		return nil, nil, errors.New("wat2wasm not found")
	}

	var bufStdout bytes.Buffer
	var bufStderr bytes.Buffer

	cmd := exec.Command(wat2wasmPath, args...)
	cmd.Stdout = &bufStdout
	cmd.Stderr = &bufStderr

	err = cmd.Run()
	stdout = bufStdout.Bytes()
	stderr = bufStderr.Bytes()
	return
}

func init() {
	const baseName = "wa.wat2wasm.exe"

	// 1. exe 同级目录存在 wat2wasm ?
	wat2wasmPath = filepath.Join(curExeDir(), baseName)
	if exeExists(wat2wasmPath) {
		return
	}

	// 2. 当前目录存在 wat2wasm ?
	cwd, _ := os.Getwd()
	wat2wasmPath = filepath.Join(cwd, baseName)
	if exeExists(wat2wasmPath) {
		return
	}

	// 3. 本地系统存在 wat2wasm ?
	if s, _ := exec.LookPath(baseName); s != "" {
		wat2wasmPath = s
		return
	}

	// 4. wat2wasm 安装到 exe 所在目录 ?
	wat2wasmPath = filepath.Join(curExeDir(), baseName)
	if err := os.WriteFile(wat2wasmPath, wabt.LoadWat2Wasm(), 0777); err != nil {
		logger.Tracef(&config.EnableTrace_app, "install wat2wasm failed: %+v", err)
		return
	}
}

// 是否为目录
func isDir(path string) bool {
	if fi, _ := os.Lstat(path); fi != nil && fi.IsDir() {
		return true
	}
	return false
}

// exe 文件存在
func exeExists(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	if !fi.Mode().IsRegular() {
		return false
	}
	return true
}

// 当前执行程序所在目录
func curExeDir() string {
	s, err := os.Executable()
	if err != nil {
		logger.Panicf("os.Executable() failed: %+v", err)
	}
	return filepath.Dir(s)
}
