// 版权 @2024 凹语言 作者。保留所有权利。

package printer

import (
	"fmt"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// (import "syscall_js" "print_bool" (func $syscall$js.__import__print_bool (param i32)))

func (p *watPrinter) printImport() error {
	if len(p.m.Imports) > 0 {
		fmt.Fprintln(p.w)
	}
	for _, importSpec := range p.m.Imports {
		switch importSpec.ObjKind {
		case token.GLOBAL:
			panic("TODO")
		case token.FUNC:
			fmt.Fprint(p.w, p.indent)
			fmt.Fprintf(p.w, "(import %q %q", importSpec.ObjModule, importSpec.ObjName)
			p.printImport_func(importSpec)
			fmt.Fprint(p.w, ")\n")

		case token.MEMORY:
		case token.TABLE:
		default:
			panic("unreachable")
		}
	}
	return nil
}

func (p *watPrinter) printImport_func(importSpec *ast.ImportSpec) {
	fmt.Fprintf(p.w, " (func $%s", importSpec.FuncName)

	fnType := importSpec.FuncType
	if len(fnType.Params) > 0 {
		for _, x := range fnType.Params {
			fmt.Fprintf(p.w, " (param %v)", x.Type)
		}
	}
	if len(fnType.Results) > 0 {
		fmt.Fprintf(p.w, " (result")
		for _, x := range fnType.Results {
			fmt.Fprintf(p.w, " %v)", x)
		}
		fmt.Fprint(p.w, ")")
	}

	fmt.Fprint(p.w, ")")
	return
}

/*

type ImportSpec struct {
	ObjModule string      // 模块路径
	ObjName   string      // 实体名字
	ObjKind   token.Token // 导入类别: MEMORY, TABLE, FUNC, GLOBAL

	Memory     *Memory     // 导入内存
	TableName  string      // 导入Table编号
	GlobalName string      // 导入全局变量名字
	GlobalType token.Token // 导入全局变量类型: I32, I64, F32, F64
	FuncName   string      // 导入后函数名字
	FuncType   *FuncType   // 导入函数类型
}
*/
