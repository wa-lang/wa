// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

func (p *parser) parseFile() {
	// 解析开头的关联文档
	p.prog.Doc = p.parseCommentGroup(true)

	// 解析代码主体
	for {
		if p.err != nil {
			return
		}
		if p.tok == token.EOF {
			break
		}

		switch p.tok {
		case token.COMMENT:
			commentObj := p.parseCommentGroup(true)
			p.prog.Comments = append(p.prog.Comments, commentObj)
			p.prog.Objects = append(p.prog.Objects, commentObj)

		case token.GAS_X64_INTEL_SYNTAX:
			if p.cpu == abi.X64Unix || p.cpu == abi.X64Windows {
				p.prog.IntelSyntax = &ast.GasIntelSyntaxNoprefix{
					Pos: p.pos,
				}
				p.acceptToken(token.GAS_X64_INTEL_SYNTAX)
				p.acceptToken(token.GAS_X64_NOPREFIX)
			} else {
				p.errorf(p.pos, "%v only enabled on X64 CPU", p.tok)
			}

		case token.GAS_EXTERN:
			ext := &ast.GasExtern{Pos: p.pos}

			p.acceptToken(token.GAS_EXTERN)
			ext.Name = p.parseIdent()
			p.prog.Externs = append(p.prog.Externs, ext)
			p.consumeSemicolonList()

			p.prog.Externs = append(p.prog.Externs, ext)
			p.prog.Objects = append(p.prog.Objects, ext)

		case token.GAS_GLOBL:
			p.acceptToken(token.GAS_GLOBL)

			// 声明为导出符号, 后续解析到真实定义的时候再合并信息
			globalName := p.parseIdent()
			p.gasGlobal[globalName] = true
			p.consumeSemicolonList()

		case token.GAS_SECTION:
			p.acceptToken(token.GAS_SECTION)

			p.gasSectionName = p.parseIdent()
			p.gasAlign = 0

			p.consumeSemicolonList()

		case token.GAS_ALIGN:
			p.acceptToken(token.GAS_ALIGN)
			p.gasAlign = p.parseIntLit()
			p.consumeSemicolonList()

		case token.IDENT:
			switch p.gasSectionName {
			case "":
				p.errorf(p.pos, "%s missing section name", p.lit)

			case ".data", ".radata", ".bss":
				beginPos := p.pos
				name := p.parseIdent()
				switch p.tok {
				case token.ASSIGN:
					p.next()

					constObj := &ast.Const{
						Pos:  beginPos,
						Name: name,
					}

					constObj.Doc = p.parseDocComment(&p.prog.Comments, constObj.Pos)
					if constObj.Doc != nil {
						p.prog.Objects = p.prog.Objects[:len(p.prog.Objects)-1]
					}

					// TODO: 重新基础常量解析, 包含标识符
					constObj.Value = p.parseBasicLit()
					if constObj.Value.LitKind == token.STRING {
						p.errorf(constObj.Value.Pos, "const(%s) donot support string type", constObj.Value.LitString)
					}
					constObj.Comment = p.parseTailComment(constObj.Pos)
					p.consumeSemicolonList()

					p.prog.Consts = append(p.prog.Consts, constObj)
					p.prog.Objects = append(p.prog.Objects, constObj)

				case token.COLON:
					p.next()

					globalObj := &ast.Global{
						Pos:  beginPos,
						Name: name,
						Init: &ast.InitValue{},

						Section: p.gasSectionName,
						Align:   p.gasAlign,
					}

					globalObj.Doc = p.parseDocComment(&p.prog.Comments, globalObj.Pos)
					if globalObj.Doc != nil {
						p.prog.Objects = p.prog.Objects[:len(p.prog.Objects)-1]
					}

					// 解析类型Token
					globalObj.TypeTok = p.acceptToken(
						token.GAS_BYTE,
						token.GAS_SHORT,
						token.GAS_LONG,
						token.GAS_QUAD,
						token.GAS_ASSCII,
						token.GAS_ASSCIZ,
						token.GAS_SKIP,
						token.GAS_INCBIN,
					)

					// 解析初始化值
					globalObj.Init.Pos = p.pos
					if p.tok == token.IDENT {
						globalObj.Init.Symbal = p.parseIdent()
						globalObj.Init.Comment = p.parseTailComment(globalObj.Init.Pos)
						p.consumeSemicolonList()
					} else {
						globalObj.Init.Lit = p.parseBasicLit()
						globalObj.Init.Comment = p.parseTailComment(globalObj.Init.Pos)
						p.consumeSemicolonList()
					}

					// TODO: 验证初始化值的合法性, 填充类型和Size

					p.prog.Globals = append(p.prog.Globals, globalObj)
					p.prog.Objects = append(p.prog.Objects, globalObj)

				default:
					p.errorf(p.pos, "expect %v or %v, got %v", token.ASSIGN, token.COLON, p.tok)
				}

			case ".text", ".init", ".fini":
				funcObj := &ast.Func{
					Pos:  p.pos,
					Tok:  p.tok,
					Type: new(ast.FuncType),
					Body: new(ast.FuncBody),

					Section: p.gasSectionName,
				}

				assert(p.cpu == abi.RISCV64 || p.cpu == abi.RISCV32 || p.cpu == abi.LOONG64)

				funcObj.Body.Pos = p.pos

				for {
					if p.err != nil {
						return
					}
					if p.tok == token.EOF {
						return
					}

					// 注释
					if p.tok == token.COMMENT {
						commentObj := p.parseCommentGroup(false)
						funcObj.Body.Comments = append(funcObj.Body.Comments, commentObj)
						funcObj.Body.Objects = append(funcObj.Body.Objects, commentObj)
						continue
					}

					// 解析指令
					if p.tok == token.IDENT || p.tok.IsAs() {
						inst := p.parseInst(funcObj)
						funcObj.Body.Insts = append(funcObj.Body.Insts, inst)
						funcObj.Body.Objects = append(funcObj.Body.Objects, inst)
						continue
					}

					// 未知 token
					break
				}

				p.prog.Funcs = append(p.prog.Funcs, funcObj)
				p.prog.Objects = append(p.prog.Objects, funcObj)

			default:
				p.errorf(p.pos, "invalid section name: %s", p.gasSectionName)
			}

		case token.CONST_zh:
			p.gasSectionName = ""
			p.gasAlign = 0

			constObj := p.parseConst(p.tok)
			p.prog.Consts = append(p.prog.Consts, constObj)
			p.prog.Objects = append(p.prog.Objects, constObj)

		case token.GLOBAL_zh:
			p.gasSectionName = ""
			p.gasAlign = 0

			globalObj := p.parseGlobal(p.tok)
			p.prog.Globals = append(p.prog.Globals, globalObj)
			p.prog.Objects = append(p.prog.Objects, globalObj)

		case token.FUNC_zh:
			p.gasSectionName = ""
			p.gasAlign = 0

			funcObj := p.parseFunc(p.tok)
			p.prog.Funcs = append(p.prog.Funcs, funcObj)
			p.prog.Objects = append(p.prog.Objects, funcObj)

		default:
			p.errorf(p.pos, "unkonw token: %v", p.tok)
		}
	}

	// 收集信息导出符号信息
	for _, g := range p.prog.Globals {
		if p.gasGlobal[g.Name] {
			g.Exported = true
		}
	}
	for _, fn := range p.prog.Funcs {
		if p.gasGlobal[fn.Name] {
			fn.ExportName = fn.Name
		}
	}
}
