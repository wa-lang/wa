// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

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

func isSameInstList(s1, s2 []ast.Instruction) bool {
	h1 := (*reflect.SliceHeader)(unsafe.Pointer(&s1))
	h2 := (*reflect.SliceHeader)(unsafe.Pointer(&s2))
	return h1.Data == h2.Data
}

func fixName(s string) string {
	if strings.ContainsAny(s, "/`") {
		s = strings.ReplaceAll(s, "/", ".")
		s = strings.ReplaceAll(s, "`", ".")
	}
	return s
}

// C 语言中 \x 转义序列是“贪婪”的, 会一直读取尽可能多的十六进制字符作为转义值的一部分，不会自动终止。
// 因此 \x000a 会被解析为 \x0a, 而不是 \x00 和 "0a" 字符串
// 解决的办法是通过字符串强制切割: \x00""0a

// 对于 gas 汇编中, 采用3个数字的八进制转义
func toGasString(dataValue []byte) string {
	var sb strings.Builder
	for _, b := range dataValue {
		// 仅保留安全的 ASCII 可打印字符, 避开双引号 " 和反斜杠 \
		if b >= 32 && b <= 126 && b != '"' && b != '\\' {
			sb.WriteByte(b)
		} else {
			sb.WriteString(fmt.Sprintf("\\%03o", b))
		}
	}
	return sb.String()
}
