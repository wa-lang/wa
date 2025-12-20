// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/token"
)

// 目前宿主函数是固定的

func (p *wat2laWorker) buildImport(w io.Writer) error {
	if len(p.m.Imports) == 0 {
		return nil
	}

	// 同一个对象可能被导入多次
	var hostGlobalMap = make(map[string]bool)
	var hostFuncMap = make(map[string]bool)

	// 导入全局的只读变量
	for _, importSpec := range p.m.Imports {
		if importSpec.ObjKind != token.GLOBAL {
			continue
		}

		globalName := importSpec.ObjModule + "." + importSpec.ObjName
		globalType := importSpec.GlobalType

		// 已经处理过
		if hostGlobalMap[globalName] {
			continue
		}
		hostGlobalMap[globalName] = true

		fmt.Fprintf(w, "extern %s %s_%s;\n", p.getCType(globalType), p.opt.Prefix, toCName(globalName))
	}
	if len(hostGlobalMap) > 0 {
		fmt.Fprintln(w)
	}

	// 声明原始的宿主函数
	for _, importSpec := range p.m.Imports {
		if importSpec.ObjKind != token.FUNC {
			continue
		}

		fnName := importSpec.ObjModule + "." + importSpec.ObjName
		fnType := importSpec.FuncType

		// 已经处理过
		if hostFuncMap[fnName] {
			continue
		}
		hostFuncMap[fnName] = true

		// 返回值类型
		cRetType := p.getHostFuncCRetType(fnType)
		if len(fnType.Results) > 1 {
			panic("wat2c: host func donot support multi return value")
		}

		// 返回值通过栈传递
		fmt.Fprintf(w, "extern %s %s_%s(", cRetType, p.opt.Prefix, toCName(fnName))
		if len(fnType.Params) > 0 {
			for i, x := range fnType.Params {
				var argName string
				if x.Name != "" {
					argName = toCName(x.Name)
				} else {
					argName = fmt.Sprintf("arg%d", i)
				}
				if i > 0 {
					fmt.Fprint(w, ", ")
				}

				switch x.Type {
				case token.I32:
					fmt.Fprintf(w, "int32_t %v", argName)
				case token.I64:
					fmt.Fprintf(w, "int64_t %v", argName)
				case token.F32:
					fmt.Fprintf(w, "float %v", argName)
				case token.F64:
					fmt.Fprintf(w, "double %v", argName)
				default:
					unreachable()
				}
			}
		}
		fmt.Fprintf(w, ");\n")
	}
	if len(hostFuncMap) > 0 {
		fmt.Fprintln(w)
	}

	// 定义导入后的全局变量
	for _, importSpec := range p.m.Imports {
		if importSpec.ObjKind != token.GLOBAL {
			continue
		}

		fmt.Fprintf(w, "#define %s_%s %s_%s // import %s.%s\n",
			p.opt.Prefix, toCName(importSpec.GlobalName),
			p.opt.Prefix, toCName(importSpec.ObjModule+"."+importSpec.ObjName),
			importSpec.ObjModule, importSpec.ObjName,
		)
	}
	if len(hostGlobalMap) > 0 {
		fmt.Fprintln(w)
	}

	// 定义导入后的函数
	for _, importSpec := range p.m.Imports {
		if importSpec.ObjKind != token.FUNC {
			continue
		}

		fmt.Fprintf(w, "#define %s_%s %s_%s // import %s.%s\n",
			p.opt.Prefix, toCName(importSpec.FuncName),
			p.opt.Prefix, toCName(importSpec.ObjModule+"."+importSpec.ObjName),
			importSpec.ObjModule, importSpec.ObjName,
		)
	}
	if len(hostFuncMap) > 0 {
		fmt.Fprintln(w)
	}

	return nil
}
