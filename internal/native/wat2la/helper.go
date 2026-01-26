// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"strconv"

	nativeast "wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *wat2laWorker) Tracef(foramt string, a ...interface{}) {
	if p.trace {
		fmt.Printf(foramt, a...)
	}
}

func (p *wat2laWorker) findGlobalName(ident string) string {
	if ident == "" {
		panic("wat2la: empty local name")
	}

	if idx, err := strconv.Atoi(ident); err == nil {
		if idx < 0 || idx >= len(p.m.Globals) {
			panic(fmt.Sprintf("wat2la: unknown global %q", ident))
		}
		return p.m.Globals[idx].Name
	}
	for _, g := range p.m.Globals {
		if g.Name == ident {
			return g.Name
		}
	}
	panic("unreachable")
}

func (p *wat2laWorker) findGlobalType(ident string) token.Token {
	if ident == "" {
		panic("wat2la: empty global name")
	}

	// 全局变量是索引类型
	if idx, err := strconv.Atoi(ident); err == nil {
		if idx < 0 || idx >= len(p.m.Globals) {
			panic(fmt.Sprintf("wat2x64: unknown global %q", ident))
		}

		// 模块内部定义的全局变量
		return p.m.Globals[idx].Type
	}

	// 查找本地定义的全局对象
	for _, g := range p.m.Globals {
		if g.Name == ident {
			return g.Type
		}
	}
	panic("unreachable")
}

func (p *wat2laWorker) findLocalOffset(fnNative *nativeast.Func, fn *ast.Func, ident string) int {
	if ident == "" {
		panic("wat2la: empty local name")
	}

	if idx, err := strconv.Atoi(ident); err == nil {
		if idx < 0 || idx >= len(fn.Type.Params)+len(fn.Locals) {
			panic(fmt.Sprintf("wat2la: unknown local %q", ident))
		}
		if idx < len(fn.Type.Params) {
			return fnNative.Type.Args[idx].RBPOff
		}
		n := idx - len(fn.Type.Params)
		return fnNative.Body.Locals[n].RBPOff
	}
	for idx, arg := range fn.Type.Params {
		if arg.Name == ident {
			return fnNative.Type.Args[idx].RBPOff
		}
	}
	for idx, arg := range fn.Locals {
		if arg.Name == ident {
			return fnNative.Body.Locals[idx].RBPOff
		}
	}
	panic("unreachable")
}

func (p *wat2laWorker) findLocalType(fn *ast.Func, ident string) token.Token {
	if ident == "" {
		panic("wat2la: empty local name")
	}

	if idx, err := strconv.Atoi(ident); err == nil {
		if idx < 0 || idx >= len(fn.Type.Params)+len(fn.Type.Results)+len(fn.Locals) {
			panic(fmt.Sprintf("wat2la: unknown local %q", ident))
		}
		if idx < len(fn.Type.Params) {
			return fn.Type.Params[idx].Type
		}
		idx = idx - len(fn.Type.Params)
		if idx < len(fn.Type.Results) {
			return fn.Type.Results[idx]
		}
		idx = idx - len(fn.Type.Results)
		return fn.Locals[idx].Type
	}
	for _, arg := range fn.Type.Params {
		if arg.Name == ident {
			return arg.Type
		}
	}
	for _, local := range fn.Locals {
		if local.Name == ident {
			return local.Type
		}
	}
	panic("unreachable")
}

func (p *wat2laWorker) findType(ident string) *ast.FuncType {
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

func (p *wat2laWorker) findFuncType(ident string) *ast.FuncType {
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

func (p *wat2laWorker) findFuncIndex(ident string) int {
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
			return len(p.m.Imports) + i
		}
	}

	panic(fmt.Sprintf("wat2la: unknown func %q", ident))
}

func (p *wat2laWorker) makeLabelId(prefix, name, suffix string) string {
	if name != "" && suffix != "" && suffix[0] != '.' {
		suffix = "." + suffix
	}
	return prefix + name + suffix
}
