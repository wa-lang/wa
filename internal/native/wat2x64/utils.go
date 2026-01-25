// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"fmt"

	nativeast "wa-lang.org/wa/internal/native/ast"
	nativetok "wa-lang.org/wa/internal/native/token"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func assert(condition bool, args ...interface{}) {
	if !condition {
		if msg := fmt.Sprint(args...); msg != "" {
			panic(fmt.Sprintf("assert failed, %s", msg))
		} else {
			panic("assert failed")
		}
	}
}

func unreachable() {
	panic("unreachable")
}

func wat2nativeType(typ token.Token) nativetok.Token {
	switch typ {
	case token.I32:
		return nativetok.I32
	case token.I64:
		return nativetok.I64
	case token.F32:
		return nativetok.F32
	case token.F64:
		return nativetok.F64
	default:
		panic("unreachable")
	}
}

// 转化为本地函数结构(不含指令)
func wat2nativeFunc(fnName string, fnType *ast.FuncType, fnLocals []ast.Field) *nativeast.Func {
	fnNative := &nativeast.Func{
		Name: fnName,
		Type: &nativeast.FuncType{
			Args:   make([]*nativeast.Local, len(fnType.Params)),
			Return: make([]*nativeast.Local, len(fnType.Results)),
		},
		Body: &nativeast.FuncBody{
			Locals: make([]*nativeast.Local, len(fnLocals)),
		},
	}
	for i, arg := range fnType.Params {
		fnNative.Type.Args[i] = &nativeast.Local{
			Name: arg.Name,
			Type: wat2nativeType(arg.Type),
			Cap:  1,
		}
	}
	for i, typ := range fnType.Results {
		fnNative.Type.Return[i] = &nativeast.Local{
			Name: fmt.Sprintf("%s%d", kFuncRetNamePrefix, i),
			Type: wat2nativeType(typ),
			Cap:  1,
		}
	}
	for i, local := range fnLocals {
		fnNative.Body.Locals[i] = &nativeast.Local{
			Name: local.Name,
			Type: wat2nativeType(local.Type),
			Cap:  1,
		}
	}
	return fnNative
}
