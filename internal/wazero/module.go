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

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"wa-lang.org/wa/internal/app/waruntime"
	"wa-lang.org/wa/internal/config"
)

// wasm 模块, 可多次执行
type Module struct {
	cfg *config.Config

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
func BuildModule(
	cfg *config.Config, wasmName string, wasmBytes []byte, wasmArgs ...string,
) (*Module, error) {
	m := &Module{
		cfg:       cfg,
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
func (p *Module) RunInitOnce() (stdout, stderr []byte, err error) {
	p.wazeroOnce.Do(func() { p.runInitFunc() })
	stdout = p.stdoutBuffer.Bytes()
	stderr = p.stderrBuffer.Bytes()
	err = p.wazeroInitErr
	return
}

// 执行指定函数(init会被强制执行一次)
func (p *Module) RunFunc(name string, args ...uint64) (result []uint64, stdout, stderr []byte, err error) {
	stdout, stderr, err = p.RunInitOnce()
	if err != nil {
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
		WithArgs(append([]string{p.wasmName}, p.wasmArgs...)...)

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

	switch p.cfg.WaOS {
	case config.WaOS_arduino:
		if _, err = waruntime.ArduinoInstantiate(p.wazeroCtx, p.wazeroRuntime); err != nil {
			p.wazeroInitErr = err
			return err
		}
	case config.WaOS_chrome:
		if _, err = waruntime.ChromeInstantiate(p.wazeroCtx, p.wazeroRuntime); err != nil {
			p.wazeroInitErr = err
			return err
		}
	case config.WaOS_wasi:
		if _, err = waruntime.WasiInstantiate(p.wazeroCtx, p.wazeroRuntime); err != nil {
			p.wazeroInitErr = err
			return err
		}
	}

	return nil
}

func (p *Module) runInitFunc() {
	if err := p.wazeroInitErr; err != nil {
		return
	}
	if p.wazeroModule != nil {
		return
	}
	p.wazeroModule, p.wazeroInitErr = p.wazeroRuntime.InstantiateModule(
		p.wazeroCtx, p.wazeroCompileModule, p.wazeroConf,
	)
}
