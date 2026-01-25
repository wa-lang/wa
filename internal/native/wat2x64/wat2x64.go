// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"bytes"
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	watast "wa-lang.org/wa/internal/wat/ast"
	watparser "wa-lang.org/wa/internal/wat/parser"
)

const DebugMode = false

// Wat程序转译到 X64汇编
func Wat2X64(filename string, source []byte, cpuType abi.CPUType) (m *watast.Module, code []byte, err error) {
	m, err = watparser.ParseModule(filename, source)
	if err != nil {
		return m, nil, err
	}

	worker := newWat2X64Worker(filename, m, cpuType)
	code, err = worker.BuildProgram()
	return
}

type wat2X64Worker struct {
	m *watast.Module

	filename string
	cpuType  abi.CPUType

	fnWasmR0Base      int // 当前函数的WASM栈R0位置
	fnMaxCallArgsSize int // 调用子函数需要的最大空间

	nextId int64 // 用于生成唯一ID

	trace bool // 调试开关
}

func newWat2X64Worker(filename string, mWat *watast.Module, cpuType abi.CPUType) *wat2X64Worker {
	assert(cpuType == abi.X64Windows || cpuType == abi.X64Unix)

	p := &wat2X64Worker{
		m:        mWat,
		cpuType:  cpuType,
		filename: filename,
		trace:    DebugMode,
	}

	// 如果 start 字段为空, 则尝试用 _start 导出函数替代
	if p.m.Start == "" {
		for _, fn := range p.m.Funcs {
			if fn.ExportName == "_start" {
				p.m.Start = fn.Name
			}
		}
	}

	return p
}

func (p *wat2X64Worker) BuildProgram() (code []byte, err error) {
	var out bytes.Buffer

	fmt.Fprintf(&out, "# 源文件: %s, ABI: %s\n", p.filename, p.cpuType)
	fmt.Fprintf(&out, "# 自动生成的代码, 不要手动修改!!!\n\n")

	p.gasIntelSyntax(&out)
	fmt.Fprintln(&out)

	// 生成运行时函数
	p.buildRuntimeHead(&out)

	// 导入函数
	if err := p.buildImport(&out); err != nil {
		return nil, err
	}

	// 构建内存
	if err := p.buildMemory(&out); err != nil {
		return nil, err
	}

	// 构建表格
	if err := p.buildTable(&out); err != nil {
		return nil, err
	}

	// 构建全局变量
	if err := p.buildGlobal(&out); err != nil {
		return nil, err
	}

	// 启动函数
	if err := p.buildStart(&out); err != nil {
		return nil, err
	}

	// 生成运行时函数
	if err := p.buildRuntimeImpl(&out); err != nil {
		return nil, err
	}

	// 构建函数
	if err := p.buildFuncs(&out); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func (p *wat2X64Worker) genNextId() string {
	nextId := p.nextId
	p.nextId++
	return fmt.Sprintf("%08X", uint64(nextId))
}
