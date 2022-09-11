// 版权 @2019 凹语言 作者。保留所有权利。

package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"

	"github.com/wa-lang/wa/internal/config"
	"github.com/wa-lang/wa/internal/logger"
	"github.com/wa-lang/wabt-go"
)

var wat2wasmPath string

func RunWasm(filename string) (stdoutStderr []byte, err error) {
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

	ctx := context.Background()
	r := wazero.NewRuntimeWithConfig(ctx,
		wazero.NewRuntimeConfig().WithWasmCore2(),
	)
	defer r.Close(ctx)

	_, err = r.NewModuleBuilder("wasi_snapshot_preview1").
		ExportFunction("fd_write", func(m api.Module, fd, iov, iov_len, p_nwritten uint32) uint32 {
			pos, _ := m.Memory().ReadUint32Le(ctx, iov)
			n, _ := m.Memory().ReadUint32Le(ctx, iov+4)
			bytes, _ := m.Memory().Read(ctx, pos, n)
			fmt.Print(string(bytes))
			return 0
		}).
		ExportFunction("waPuts", func(m api.Module, pos, len uint32) {
			bytes, _ := m.Memory().Read(ctx, pos, len)
			fmt.Print(string(bytes))
			return
		}).
		ExportFunction("waPrintI32", func(m api.Module, v uint32) {
			fmt.Print(v)
			return
		}).
		ExportFunction("waPrintRune", func(m api.Module, ch uint32) {
			fmt.Printf("%c", rune(ch))
			return
		}).
		Instantiate(ctx, r)
	if err != nil {
		return nil, err
	}

	_, err = r.InstantiateModuleFromBinary(ctx, wasmBytes)
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

func RunWat2Wasm(args ...string) (stdoutStderr []byte, err error) {
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
