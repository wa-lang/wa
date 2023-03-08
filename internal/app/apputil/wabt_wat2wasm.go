// 版权 @2019 凹语言 作者。保留所有权利。

package apputil

import (
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

	"wa-lang.org/wa/internal/app/waruntime"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wabt-go"
)

var wat2wasmPath string

func RunWasm(cfg *config.Config, filename string, wasmArgs ...string) (stdoutStderr []byte, err error) {
	dst := filename
	if strings.HasSuffix(filename, ".wat") {
		dst += ".wasm"
		if stdoutStderr, err = RunWat2Wasm(filename, "-o", dst); err != nil {
			return stdoutStderr, err
		}
		defer os.Remove(dst)
	}

	wasmBytes, err := os.ReadFile(dst)
	if err != nil {
		return nil, err
	}

	wasmExe := filepath.Base(filename)

	conf := wazero.NewModuleConfig().
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithStdin(os.Stdin).
		WithRandSource(rand.Reader).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithArgs(append([]string{wasmExe}, wasmArgs...)...)

	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	code, err := r.CompileModule(ctx, wasmBytes)
	if err != nil {
		return nil, err
	}

	switch cfg.WaOS {
	case config.WaOS_Arduino:
		if _, err = waruntime.ArduinoInstantiate(ctx, r); err != nil {
			return nil, err
		}
	case config.WaOS_Chrome:
		if _, err = waruntime.ChromeInstantiate(ctx, r); err != nil {
			return nil, err
		}
	case config.WaOS_Wasi:
		if _, err = waruntime.WasiInstantiate(ctx, r); err != nil {
			return nil, err
		}
	}

	_, err = r.InstantiateModule(ctx, code, conf)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func InstallWat2wasm(path string) error {
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

func RunWat2Wasm(args ...string) (stdoutStderr []byte, err error) {
	muRunWat2Wasm.Lock()
	defer muRunWat2Wasm.Unlock()

	if wat2wasmPath == "" {
		logger.Tracef(&config.EnableTrace_app, "wat2wasm not found")
		return nil, errors.New("wat2wasm not found")
	}
	cmd := exec.Command(wat2wasmPath, args...)
	stdoutStderr, err = cmd.CombinedOutput()
	return
}

func init() {
	var baseName = "wat2wasm"
	if runtime.GOOS == "windows" {
		baseName += ".exe"
	}

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

	// 4. wat2wasm 安装到同级目录存 ?
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
