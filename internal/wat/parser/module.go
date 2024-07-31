// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// https://webassembly.github.io/spec/core/text/modules.html

// module ::= (module id? (m:modulefield))
//
// modulefield ::= type
//              |  import
//              |  func
//              |  table
//              |  memory
//              |  global
//              |  export
//              |  start
//              |  elem
//              |  data

// 解析 module
// (module $name ...)
func (p *parser) parseModule() {
	// 头部注释
	p.consumeComments()

	p.acceptToken(token.LPAREN)
	p.acceptToken(token.MODULE)

	// 模块名字
	if p.tok == token.IDENT {
		p.module.Name = p.lit
		p.next()
	}

	// 解析模块Section
	for {
		p.consumeComments()
		if p.tok == token.LPAREN {
			p.parseModuleSection()
		} else {
			break
		}
	}

	// 补充导出全局变量/函数
	for _, g := range p.module.Globals {
		if g.ExportName != "" {
			p.module.Exports = append(p.module.Exports, &ast.ExportSpec{
				Name:      g.ExportName,
				Kind:      token.GLOBAL,
				GlobalIdx: g.Name,
			})
		}
	}
	for _, fn := range p.module.Funcs {
		if fn.ExportName != "" {
			p.module.Exports = append(p.module.Exports, &ast.ExportSpec{
				Name:    fn.ExportName,
				Kind:    token.FUNC,
				FuncIdx: fn.Name,
			})
		}
	}

	p.consumeComments()
	p.acceptToken(token.RPAREN)
}

// 解析模块段
func (p *parser) parseModuleSection() {
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)
	defer p.consumeComments()

	switch p.tok {
	default:
		p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)

	case token.TYPE:
		typ := p.parseModuleSection_type()
		p.module.Types = append(p.module.Types, typ)

	case token.IMPORT:
		spec := p.parseModuleSection_import()
		p.module.Imports = append(p.module.Imports, spec)
	case token.EXPORT:
		spec := p.parseModuleSection_export()
		p.module.Exports = append(p.module.Exports, spec)

	case token.MEMORY:
		mem := p.parseModuleSection_memory()
		p.module.Memory = mem
	case token.DATA:
		dataSection := p.parseModuleSection_data()
		p.module.Data = append(p.module.Data, dataSection)

	case token.TABLE:
		tab := p.parseModuleSection_table()
		p.module.Table = tab

	case token.ELEM:
		elemSection := p.parseModuleSection_elem()
		p.module.Elem = append(p.module.Elem, elemSection)

	case token.GLOBAL:
		g := p.parseModuleSection_global()
		p.module.Globals = append(p.module.Globals, g)

	case token.FUNC:
		fn := p.parseModuleSection_func()
		p.module.Funcs = append(p.module.Funcs, fn)

	case token.START:
		p.module.Start = p.parseModuleSection_start()
	}
}
