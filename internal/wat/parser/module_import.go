// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *parser) parseModuleSection_import() *ast.ImportSpec {
	p.acceptToken(token.IMPORT)

	spec := &ast.ImportSpec{}

	// 宿主模块和名字
	spec.ModulePath = p.parseStringLit()
	spec.FuncPath = p.parseStringLit()

	// 导入函数的类型
	p.parseImportFuncType(spec)

	return spec
}
