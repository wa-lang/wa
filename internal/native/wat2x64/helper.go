// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"fmt"
	"strconv"

	nativeast "wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *wat2X64Worker) Tracef(foramt string, a ...interface{}) {
	if p.trace {
		fmt.Printf(foramt, a...)
	}
}

func (p *wat2X64Worker) findGlobalType(ident string) token.Token {
	if ident == "" {
		panic("wat2la: empty global name")
	}

	// 全局变量是索引类型
	if idx, err := strconv.Atoi(ident); err == nil {
		if idx < 0 {
			panic(fmt.Sprintf("wat2la: unknown global %q", ident))
		}

		// 是导入的全局变量
		if idx < p.importGlobalCount {
			var nextIndex int
			for _, importSpec := range p.m.Imports {
				if importSpec.ObjKind != token.GLOBAL {
					continue
				}
				// 找到导入对象
				if nextIndex == idx {
					return importSpec.GlobalType
				}
				// 更新索引
				nextIndex++
			}
		}

		// 模块内部定义的全局变量
		return p.m.Globals[idx-p.importGlobalCount].Type
	}

	// 从导入对象开始查找
	for _, importSpec := range p.m.Imports {
		if importSpec.ObjKind != token.GLOBAL {
			continue
		}
		// 找到导入对象
		if importSpec.GlobalName == ident {
			return importSpec.GlobalType
		}
	}

	// 查找本地定义的全局对象
	for _, g := range p.m.Globals {
		if g.Name == ident {
			return g.Type
		}
	}
	panic("unreachable")
}

func (p *wat2X64Worker) findLocalOffset(fnNative *nativeast.Func, fn *ast.Func, ident string) int {
	if ident == "" {
		panic("wat2la: empty local name")
	}

	if idx, err := strconv.Atoi(ident); err == nil {
		if idx < 0 || idx >= len(fn.Type.Params)+len(fn.Body.Locals) {
			panic(fmt.Sprintf("wat2la: unknown local %q", ident))
		}
		if idx < len(fn.Type.Params) {
			return fnNative.Type.Args[idx].Off
		}
		n := idx - len(fn.Type.Params)
		return fnNative.Body.Locals[n].Off
	}
	for idx, arg := range fn.Type.Params {
		if arg.Name == ident {
			return fnNative.Type.Args[idx].Off
		}
	}
	for idx, arg := range fn.Body.Locals {
		if arg.Name == ident {
			return fnNative.Body.Locals[idx].Off
		}
	}
	panic("unreachable")
}

func (p *wat2X64Worker) findLocalType(fn *ast.Func, ident string) token.Token {
	if ident == "" {
		panic("wat2la: empty local name")
	}

	if idx, err := strconv.Atoi(ident); err == nil {
		if idx < 0 || idx >= len(fn.Type.Params)+len(fn.Body.Locals) {
			panic(fmt.Sprintf("wat2la: unknown local %q", ident))
		}
		return p.localTypes[idx]
	}
	for idx, arg := range fn.Type.Params {
		if arg.Name == ident {
			return p.localTypes[idx]
		}
	}
	for idx, arg := range fn.Body.Locals {
		if arg.Name == ident {
			return p.localTypes[len(fn.Type.Params)+idx]
		}
	}
	panic("unreachable")
}

func (p *wat2X64Worker) findLocalName(fn *ast.Func, ident string) string {
	if ident == "" {
		panic("wat2la: empty local name")
	}

	if idx, err := strconv.Atoi(ident); err == nil {
		if idx < 0 || idx >= len(fn.Type.Params)+len(fn.Body.Locals) {
			panic(fmt.Sprintf("wat2la: unknown local %q", ident))
		}
		return p.localNames[idx]
	}
	for idx, arg := range fn.Type.Params {
		if arg.Name == ident {
			return p.localNames[idx]
		}
	}
	for idx, arg := range fn.Body.Locals {
		if arg.Name == ident {
			return p.localNames[len(fn.Type.Params)+idx]
		}
	}
	panic("unreachable")
}

func (p *wat2X64Worker) findType(ident string) *ast.FuncType {
	if ident == "" {
		panic("wat2la: empty ident")
	}

	if idx, err := strconv.Atoi(ident); err == nil {
		if idx < 0 || idx >= len(p.m.Types) {
			panic(fmt.Sprintf("wat2la: unknown type %q", ident))
		}
		return p.m.Types[idx].Type
	}
	for _, x := range p.m.Types {
		if x.Name == ident {
			return x.Type
		}
	}
	panic(fmt.Sprintf("wat2la: unknown type %q", ident))
}

func (p *wat2X64Worker) findFuncType(ident string) *ast.FuncType {
	if ident == "" {
		panic("wat2la: empty ident")
	}

	// 查找导入的函数类型
	for _, importSpec := range p.m.Imports {
		if importSpec.ObjKind != token.FUNC {
			continue
		}
		if importSpec.FuncName == ident {
			return importSpec.FuncType
		}
	}

	// 查找本地定义的函数
	for _, fn := range p.m.Funcs {
		if fn.Name == ident {
			return fn.Type
		}
	}

	panic(fmt.Sprintf("wat2la: unknown func %q", ident))
}

func (p *wat2X64Worker) findFuncIndex(ident string) int {
	if ident == "" {
		panic("wat2la: empty ident")
	}

	if idx, err := strconv.Atoi(ident); err == nil {
		return idx
	}

	// 导入函数
	{
		var nextIndex int
		for _, importSpec := range p.m.Imports {
			if importSpec.ObjKind != token.FUNC {
				continue
			}
			if ident == importSpec.FuncName {
				return nextIndex
			}
			nextIndex++
		}
	}

	// 查找本地定义的函数
	for i, fn := range p.m.Funcs {
		if fn.Name == ident {
			return p.importFuncCount + i
		}
	}

	panic(fmt.Sprintf("wat2la: unknown func %q", ident))
}

func (p *wat2X64Worker) findLabelName(label string) string {
	if label == "" {
		panic("wat2la: empty label")
	}

	idx := p.findLabelIndex(label)
	if idx < len(p.scopeLabels) {
		return p.scopeLabels[len(p.scopeLabels)-idx-1]
	}
	panic(fmt.Sprintf("wat2la: unknown label %q", label))
}

func (p *wat2X64Worker) findLabelIndex(label string) int {
	if label == "" {
		panic("wat2la: empty label")
	}

	if idx, err := strconv.Atoi(label); err == nil {
		return idx
	}
	for i := 0; i < len(p.scopeLabels); i++ {
		if s := p.scopeLabels[len(p.scopeLabels)-i-1]; s == label {
			return i
		}
	}
	panic(fmt.Sprintf("wat2la: unknown label %q", label))
}

func (p *wat2X64Worker) enterLabelScope(stkBase int, label string, results []token.Token) {
	p.scopeLabels = append(p.scopeLabels, label)
	p.scopeStackBases = append(p.scopeStackBases, stkBase)
	p.scopeResults = append(p.scopeResults, results)
}
func (p *wat2X64Worker) leaveLabelScope() {
	p.scopeLabels = p.scopeLabels[:len(p.scopeLabels)-1]
	p.scopeStackBases = p.scopeStackBases[:len(p.scopeStackBases)-1]
	p.scopeResults = p.scopeResults[:len(p.scopeResults)-1]
}
