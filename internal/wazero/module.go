// 版权 @2023 凹语言 作者。保留所有权利。

package wazero

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"wa-lang.org/wa/internal/3rdparty/wazero"
	"wa-lang.org/wa/internal/3rdparty/wazero/api"
	"wa-lang.org/wa/internal/config"
)

// wasm 模块, 可多次执行
type Module struct {
	wasmName  string
	wasmBytes []byte
	wasmArgs  []string

	stdoutBuffer bytes.Buffer
	stderrBuffer bytes.Buffer

	wazeroOnce          sync.Once
	wazeroCtx           context.Context
	wazeroConf          wazero.ModuleConfig
	wazeroRuntime       wazero.Runtime
	wazeroCompileModule wazero.CompiledModule
	wazeroModule        api.Module
	wazeroInitErr       error
}

// 构建模块(会执行编译)
func BuildModule(wasmName string, wasmBytes []byte, wasmArgs ...string) (*Module, error) {
	m := &Module{
		wasmName:  wasmName,
		wasmBytes: wasmBytes,
		wasmArgs:  wasmArgs,
	}
	if err := m.buildModule(); err != nil {
		return nil, err
	}
	return m, nil
}

// 执行初始化, 仅执行一次
func (p *Module) RunMain(mainFunc string) (stdout, stderr []byte, err error) {
	p.wazeroModule, p.wazeroInitErr = p.wazeroRuntime.InstantiateModule(
		p.wazeroCtx, p.wazeroCompileModule, p.wazeroConf,
	)

	if mainFunc != "" {
		fn := p.wazeroModule.ExportedFunction(mainFunc)
		if fn == nil && mainFunc != "_main" {
			err = fmt.Errorf("wazero: func %q not found", mainFunc)
			return
		}

		_, err = fn.Call(p.wazeroCtx)
		if err != nil {
			stdout = p.stdoutBuffer.Bytes()
			stderr = p.stderrBuffer.Bytes()
			return
		}
	}

	stdout = p.stdoutBuffer.Bytes()
	stderr = p.stderrBuffer.Bytes()
	err = p.wazeroInitErr
	return
}

// 执行指定函数(init会被强制执行一次)
func (p *Module) RunFunc(name string, args ...uint64) (result []uint64, stdout, stderr []byte, err error) {
	if p.wazeroModule == nil {
		p.wazeroModule, p.wazeroInitErr = p.wazeroRuntime.InstantiateModule(
			p.wazeroCtx, p.wazeroCompileModule, p.wazeroConf,
		)
	}
	if p.wazeroInitErr != nil {
		stdout = p.stdoutBuffer.Bytes()
		stderr = p.stderrBuffer.Bytes()
		err = p.wazeroInitErr
		return
	}

	p.stdoutBuffer.Reset()
	p.stderrBuffer.Reset()
	fn := p.wazeroModule.ExportedFunction(name)
	if fn == nil {
		err = fmt.Errorf("wazero: func %q not found", name)
		return
	}

	result, err = fn.Call(p.wazeroCtx, args...)
	stdout = p.stdoutBuffer.Bytes()
	stderr = p.stderrBuffer.Bytes()
	return
}

// 关闭模块
func (p *Module) Close() error {
	var err error
	if p.wazeroRuntime != nil {
		err = p.wazeroRuntime.Close(p.wazeroCtx)
		p.wazeroRuntime = nil
	}
	return err
}

// 判断目标类型
func ReadImportModuleName(wasmBytes []byte) (string, error) {
	wazeroCtx := context.Background()
	rt := wazero.NewRuntime(wazeroCtx)

	var err error
	wazeroCompileModule, err := rt.CompileModule(wazeroCtx, wasmBytes)
	if err != nil {
		return "", err
	}

	for _, importedFunc := range wazeroCompileModule.ImportedFunctions() {
		moduleName, funcName, isImport := importedFunc.Import()
		if !isImport {
			continue
		}

		switch moduleName {
		case "syscall_js":
			return config.WaOS_js, nil
		case "wasi_snapshot_preview1":
			return config.WaOS_wasi, nil
		case "arduino":
			return config.WaOS_arduino, nil
		case "env":
			if funcName == "blitSub" {
				return config.WaOS_wasm4, nil
			}
		}
	}
	return "", nil
}

// 是否包含用户自定义的宿主函数
func HasUnknownImportFunc(wasmBytes []byte) bool {
	wazeroCtx := context.Background()
	rt := wazero.NewRuntime(wazeroCtx)

	var err error
	wazeroCompileModule, err := rt.CompileModule(wazeroCtx, wasmBytes)
	if err != nil {
		return false
	}

	for _, importedFunc := range wazeroCompileModule.ImportedFunctions() {
		moduleName, _, isImport := importedFunc.Import()
		if !isImport {
			continue
		}

		switch moduleName {
		case "syscall_js", "wasi_snapshot_preview1", "arduino":
		default: // wasm4, unknown, ...
			return true
		}
	}
	return false
}

func (p *Module) buildModule() error {
	p.wazeroCtx = context.Background()

	p.wazeroConf = wazero.NewModuleConfig().
		WithStdout(&p.stdoutBuffer).
		WithStderr(&p.stderrBuffer).
		WithStdin(os.Stdin).
		WithRandSource(rand.Reader).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithArgs(append([]string{p.wasmName}, p.wasmArgs...)...).
		WithName(p.wasmName)

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
			p.wazeroConf = p.wazeroConf.WithEnv(key, value)
		}
	}

	p.wazeroRuntime = wazero.NewRuntime(p.wazeroCtx)

	var err error
	p.wazeroCompileModule, err = p.wazeroRuntime.CompileModule(p.wazeroCtx, p.wasmBytes)
	if err != nil {
		p.wazeroInitErr = err
		return err
	}

	// 根据导入的函数识别宿主类型
	var waOS = config.WaOS_unknown
	for _, importedFunc := range p.wazeroCompileModule.ImportedFunctions() {
		moduleName, funcName, isImport := importedFunc.Import()
		if !isImport {
			continue
		}

		if moduleName == "syscall_js" && funcName == "print_str" {
			waOS = config.WaOS_js
			break
		}

		if moduleName == "wasi_snapshot_preview1" {
			waOS = config.WaOS_wasi
			break
		}

		if moduleName == "arduino" {
			waOS = config.WaOS_arduino
			break
		}

		if moduleName == "env" && funcName == "blitSub" {
			waOS = config.WaOS_wasm4
			break
		}
	}

	switch waOS {
	case config.WaOS_unknown:
		if _, err = UnknownInstantiate(p.wazeroCtx, p.wazeroRuntime); err != nil {
			p.wazeroInitErr = err
			return err
		}
	case config.WaOS_wasi:
		if _, err = WasiInstantiate(p.wazeroCtx, p.wazeroRuntime); err != nil {
			p.wazeroInitErr = err
			return err
		}
	case config.WaOS_wasm4:
		panic("wasm4: TODO") // 浏览器执行
	case config.WaOS_js:
		if _, err = JsInstantiate(p.wazeroCtx, p.wazeroRuntime); err != nil {
			p.wazeroInitErr = err
			return err
		}
	case config.WaOS_arduino:
		if _, err = ArduinoInstantiate(p.wazeroCtx, p.wazeroRuntime); err != nil {
			p.wazeroInitErr = err
			return err
		}

	default:
		return fmt.Errorf("unknown waos: %q", waOS)
	}

	return nil
}
