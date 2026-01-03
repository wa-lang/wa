// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"bytes"
	"fmt"

	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/wasm"
	watast "wa-lang.org/wa/internal/wat/ast"
	watparser "wa-lang.org/wa/internal/wat/parser"
	wattoken "wa-lang.org/wa/internal/wat/token"
)

const DebugMode = false

// Wat程序转译到 X64汇编
func Wat2X64(filename string, source []byte, osName string) (m *watast.Module, code []byte, err error) {
	return wat2x64(filename, source, osName)
}

func wat2x64(filename string, source []byte, osName string) (m *watast.Module, code []byte, err error) {
	m, err = watparser.ParseModule(filename, source)
	if err != nil {
		return m, nil, err
	}

	worker := newWat2X64Worker(m, osName)
	code, err = worker.BuildProgram()
	return
}

type wat2X64Worker struct {
	osName string

	m *watast.Module

	importGlobalCount int // 导入全局只读变量的数目
	importFuncCount   int // 导入函数的数目

	inlinedTypeIndices []*inlinedTypeIndex
	inlinedTypes       []*wasm.FunctionType

	localNames      []string           // 参数和局部变量名
	localTypes      []wattoken.Token   // 参数和局部变量类型
	scopeLabels     []string           // 嵌套的label查询, if/block/loop
	scopeStackBases []int              // if/block/loop, 开始的栈位置
	scopeResults    [][]wattoken.Token // 对应块的返回值数量和类型
	fnWasmR0Base    int                // 当前函数的WASM栈R0位置

	constLitMap map[uint64]uint64 // 常量列表

	dataSection []*ast.Global
	textSection []*ast.Func

	trace bool // 调试开关
}

type inlinedTypeIndex struct {
	section    wasm.SectionID
	idx        wasm.Index
	inlinedIdx wasm.Index
}

func newWat2X64Worker(mWat *watast.Module, osName string) *wat2X64Worker {
	p := &wat2X64Worker{m: mWat, osName: osName, trace: DebugMode}

	// 统计导入的global和func索引
	p.importGlobalCount = 0
	p.importFuncCount = 0
	for _, importSpec := range p.m.Imports {
		switch importSpec.ObjKind {
		case wattoken.GLOBAL:
			p.importGlobalCount++
		case wattoken.FUNC:
			p.importFuncCount++
		}
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
	p.dataSection = p.dataSection[:0]
	p.textSection = p.textSection[:0]

	p.constLitMap = map[uint64]uint64{}

	var out bytes.Buffer

	fmt.Fprintf(&out, "# 自动生成的代码, 不要手动修改!!!\n\n")

	if err := p.buildImport(&out); err != nil {
		return nil, err
	}

	if err := p.buildGlobal(&out); err != nil {
		return nil, err
	}

	if err := p.buildTable(&out); err != nil {
		return nil, err
	}

	if err := p.buildMemory(&out); err != nil {
		return nil, err
	}

	if err := p.buildFuncs(&out); err != nil {
		return nil, err
	}

	// 生成常量
	if err := p.buildConstList(&out); err != nil {
		return nil, err
	}

	// 内置函数
	switch p.osName {
	case "linux":
		if err := p.buildBuiltinLinux(&out); err != nil {
			return nil, err
		}
	case "windows":
		if err := p.buildBuiltinWindows(&out); err != nil {
			return nil, err
		}
	default:
		panic(fmt.Sprintf("unknown os %s", p.osName))
	}

	// 启动函数
	if err := p.buildStart(&out); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
