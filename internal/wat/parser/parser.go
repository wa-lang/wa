// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"fmt"
	"strconv"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/scanner"
	"wa-lang.org/wa/internal/wat/token"
)

type parser struct {
	filename string
	src      []byte

	scanner scanner.Scanner
	module  *ast.Module

	pos token.Pos
	tok token.Token
	lit string

	err error
}

type parserError struct {
	pos token.Position
	msg string
}

func (e *parserError) Error() string {
	return fmt.Sprintf("%v: %s", e.pos, e.msg)
}

func newParser(path string, src []byte) *parser {
	p := &parser{
		filename: path,
		src:      src,
	}

	return p
}

func (p *parser) ParseModule() (m *ast.Module, err error) {
	p.module = &ast.Module{
		File: token.NewFile(p.filename, len(p.src)),
	}

	defer func() {
		if r := recover(); r != nil {
			if errx, ok := r.(*parserError); ok {
				err = errx
			} else {
				panic(r)
			}
		}
	}()

	p.scanner.Init(p.module.File, p.src,
		func(pos token.Position, msg string) {
			if p.err != nil {
				p.err = &parserError{pos, msg}
			}
		},
		scanner.ScanComments,
	)

	p.parseFile()

	m = p.module
	err = p.err
	return
}

func (p *parser) errorf(pos token.Pos, format string, a ...interface{}) {
	p.err = &parserError{
		pos: p.module.File.Position(pos),
		msg: fmt.Sprintf(format, a...),
	}
	panic(p.err)
}

func (p *parser) next() {
	p.pos, p.tok, p.lit = p.scanner.Scan()
}

func (p *parser) acceptToken(expectToken token.Token) {
	if p.tok != expectToken {
		p.errorf(p.pos, "expect %v, got %v", expectToken, p.tok)
	}
	p.next()
}

func (p *parser) parseFile() {
	p.next()

	// 解析开头的注释
	for p.tok == token.COMMENT {
		p.parseComment()
	}

	// 解析 module 节点
	p.parseModule()

	// 解析结尾的注释
	for p.tok == token.COMMENT {
		p.parseComment()
	}

	// 结束解析
	p.acceptToken(token.EOF)
}

func (p *parser) parseComment() string {
	s := p.lit
	p.acceptToken(token.COMMENT)
	return s
}

func (p *parser) parseModule() {
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)

	// module 关键字
	p.acceptToken(token.MODULE)

	// 模块名字
	if p.tok == token.IDENT {
		p.module.Name = p.lit
		p.next()
	}

	// 解析模块主体
	for {
		// 解析注释
		if p.tok == token.COMMENT {
			p.parseComment()
			continue
		}

		// 解析section
		p.parseModuleSection()
	}
}

func (p *parser) parseModuleSection() {
	p.acceptToken(token.LPAREN)

	switch p.tok {
	case token.IMPORT:
		p.parseModuleSection_import()
		p.acceptToken(token.RPAREN)

	case token.MEMORY:
		p.parseModuleSection_memory()
		p.acceptToken(token.RPAREN)

	case token.EXPORT:
		p.parseModuleSection_export()
		p.acceptToken(token.RPAREN)

	case token.DATA:
		p.parseModuleSection_data()
		p.acceptToken(token.RPAREN)

	case token.FUNC:
		p.parseModuleSection_func()
		p.acceptToken(token.RPAREN)

	default:
		p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
	}
}

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

func (p *parser) parseModuleSection_memory() {
	p.acceptToken(token.MEMORY)

	if p.module.Memory == nil {
		p.module.Memory = &ast.Memory{}
	}

	if p.tok == token.IDENT {
		p.module.Memory.Name = p.lit
		p.acceptToken(token.IDENT)
	}

	p.module.Memory.Pages = p.parseIntLit()
}

func (p *parser) parseModuleSection_data() {
	p.acceptToken(token.DATA)

	// (data (i32.const 8) "hello world\n")

	p.acceptToken(token.LPAREN)
	p.acceptToken(token.INS_I32_CONST)
	p.parseIntLit()
	p.acceptToken(token.RPAREN)
	p.parseStringLit()
}

func (p *parser) parseModuleSection_func() {
	p.acceptToken(token.FUNC)

	// (func $main (export "_start")

	if p.tok == token.IDENT {
		p.acceptToken(token.IDENT)
	}

	if p.tok == token.LPAREN {
		p.acceptToken(token.LPAREN)

		switch p.tok {
		case token.EXPORT:
			p.acceptToken(token.EXPORT)
			p.parseStringLit()
			p.acceptToken(token.RPAREN)
		case token.PARAM:
			p.next()
			if p.tok == token.IDENT {
				p.parseIdent()
			}
			switch p.tok {
			case token.I32, token.I64, token.F32, token.F64:
				p.next()
			default:
			}
		case token.RPAREN:
			// todo
		default:
		}
	}

}

func (p *parser) parseModuleSection_export() {
	p.acceptToken(token.EXPORT)

	// 解析导出的名字, 字符串类型
	p.parseStringLit()

	// 导出的对象
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)

	// (memory 0)
	// (func $name)

	switch p.tok {
	case token.MEMORY:
		p.acceptToken(token.MEMORY)
		p.parseIntLit() // todo

	case token.FUNC:
		p.acceptToken(token.MEMORY)
		p.parseIdent() // todo

	default:
		p.errorf(p.pos, "expect int, got %q", p.lit)
	}
}

func (p *parser) parseStringLit() string {
	s := p.lit
	p.acceptToken(token.STRING)
	return s
}

func (p *parser) parseIntLit() int {
	pos, lit := p.pos, p.lit
	p.acceptToken(token.INT)

	n, err := strconv.Atoi(lit)
	if err != nil {
		p.errorf(pos, "expect int, got %q", lit)
	}
	return n
}

func (p *parser) parseIdent() string {
	s := p.lit
	p.acceptToken(token.IDENT)
	return s
}

func (p *parser) parseImportFuncType(spec *ast.ImportSpec) {
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)

	p.acceptToken(token.FUNC)

	spec.FuncName = p.parseIdent()
	spec.FuncType = &ast.FuncType{}

	if p.tok == token.LPAREN {
		p.acceptToken(token.LPAREN)

		switch p.tok {
		case token.PARAM:
			p.parseImportFuncType_param(spec)
			p.acceptToken(token.RPAREN)
		case token.RESULT:
			p.parseImportFuncType_result(spec)
			p.acceptToken(token.RPAREN)
			return
		default:
			p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
		}
	}

	if p.tok == token.LPAREN {
		p.acceptToken(token.LPAREN)

		switch p.tok {
		case token.RESULT:
			p.parseImportFuncType_result(spec)
			p.acceptToken(token.RPAREN)
		default:
			p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
		}
	}
}

func (p *parser) parseImportFuncType_param(spec *ast.ImportSpec) {
	p.acceptToken(token.PARAM)

	for {
		var field ast.Field
		if p.tok == token.IDENT {
			field.Name = p.lit
			p.next()
		}

		switch p.tok {
		case token.I32, token.I64, token.F32, token.F64:
			field.Type = p.lit
			spec.FuncType.Params = append(spec.FuncType.Params, field)
			p.next()
		case token.RPAREN:
			return
		default:
			p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
		}
	}
}

func (p *parser) parseImportFuncType_result(spec *ast.ImportSpec) {
	p.acceptToken(token.RESULT)

	for {
		switch p.tok {
		case token.I32, token.I64, token.F32, token.F64:
			spec.FuncType.ResultsType = append(spec.FuncType.ResultsType, p.lit)
			p.next()
		case token.RPAREN:
			return
		default:
			p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
		}
	}
}
