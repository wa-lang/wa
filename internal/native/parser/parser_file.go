// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"math"
	"os"

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
			ext := &ast.Extern{Pos: p.pos, Tok: token.GAS_EXTERN}

			p.acceptToken(token.GAS_EXTERN)
			ext.Name = p.parseIdent()
			p.prog.Externs = append(p.prog.Externs, ext)
			p.consumeSemicolonList()

			p.prog.Externs = append(p.prog.Externs, ext)
			p.prog.Objects = append(p.prog.Objects, ext)

		case token.EXTERN_zh:
			ext := &ast.Extern{Pos: p.pos, Tok: token.EXTERN_zh}

			p.acceptToken(token.EXTERN_zh)
			ext.Name = p.parseIdent()
			p.prog.Externs = append(p.prog.Externs, ext)
			p.consumeSemicolonList()

			p.prog.Externs = append(p.prog.Externs, ext)
			p.prog.Objects = append(p.prog.Objects, ext)

		case token.GAS_GLOBL:
			p.acceptToken(token.GAS_GLOBL)

			// 声明为导出符号, 后续解析到真实定义的时候再合并信息
			globalName := p.parseIdent()
			p.gasGlobl[globalName] = true
			p.consumeSemicolonList()

		case token.GAS_SECTION:
			beginTok := p.pos
			p.acceptToken(token.GAS_SECTION)

			p.gasSectionName = p.parseIdent()
			if _, ok := p.gasSectionAlign[p.gasSectionName]; !ok {
				switch p.gasSectionName {
				case ".text", ".init", ".fini":
					if p.cpu == abi.X64Unix || p.cpu == abi.X64Windows {
						p.gasSectionAlign[p.gasSectionName] = 4
					} else {
						p.gasSectionAlign[p.gasSectionName] = 2
					}
				default:
					if p.cpu == abi.X64Unix || p.cpu == abi.X64Windows {
						p.gasSectionAlign[p.gasSectionName] = 1
					} else {
						p.gasSectionAlign[p.gasSectionName] = 0
					}
				}
			}
			p.gasAlign = p.gasSectionAlign[p.gasSectionName]
			p.consumeSemicolonList()

			switch p.gasSectionName {
			case ".text", ".init", ".fini":
				funcObj := &ast.Func{
					Pos:  beginTok,
					Type: new(ast.FuncType),
					Body: new(ast.FuncBody),

					Section: p.gasSectionName,
				}

				// 函数段强制4字节对齐
				if p.cpu == abi.X64Unix || p.cpu == abi.X64Windows {
					assert(p.gasAlign == 4)
				} else {
					assert(p.gasAlign == 2)
				}

				// 关联注释需要马上解析
				funcObj.Doc = p.parseDocComment(&p.prog.Comments, funcObj.Pos)
				if funcObj.Doc != nil {
					p.prog.Objects = p.prog.Objects[:len(p.prog.Objects)-1]
				}

				// 函数定义时候经先常导出符号
				if p.tok == token.GAS_GLOBL {
					p.acceptToken(token.GAS_GLOBL)

					// 声明为导出符号, 后续解析到真实定义的时候再合并信息
					globlName := p.parseIdent()
					p.gasGlobl[globlName] = true
					p.consumeSemicolonList()
				}

				// 开始解析函数定义的标签
				funcObj.Name = p.parseIdent()
				p.acceptToken(token.COLON)

				for {
					if p.err != nil {
						break
					}
					if p.tok == token.EOF {
						break
					}

					// 注释
					if p.tok == token.COMMENT {
						// 插入空行
						if len(funcObj.Body.Objects) > 0 {
							prevObj := funcObj.Body.Objects[len(funcObj.Body.Objects)-1]
							prevLine := p.posLine(prevObj.BeginPos())
							curLine := p.posLine(p.pos)
							if curLine-prevLine > 1 {
								funcObj.Body.Objects = append(funcObj.Body.Objects, &ast.BlankLine{
									Pos: p.pos - 1,
								})
							}
						}

						commentObj := p.parseCommentGroup(true)
						funcObj.Body.Comments = append(funcObj.Body.Comments, commentObj)
						funcObj.Body.Objects = append(funcObj.Body.Objects, commentObj)
						continue
					}

					// 解析指令
					if p.tok == token.IDENT || p.tok.IsAs() {
						// 插入空行
						if len(funcObj.Body.Objects) > 0 {
							prevObj := funcObj.Body.Objects[len(funcObj.Body.Objects)-1]
							prevLine := p.posLine(prevObj.BeginPos())
							curLine := p.posLine(p.pos)
							if curLine-prevLine > 1 {
								funcObj.Body.Objects = append(funcObj.Body.Objects, &ast.BlankLine{
									Pos: p.pos - 1,
								})
							}
						}

						// 解析指令的关联注释
						instDoc := p.parseDocComment(&funcObj.Body.Comments, p.pos)
						if instDoc != nil {
							funcObj.Body.Objects = funcObj.Body.Objects[:len(funcObj.Body.Objects)-1]
						}

						inst := p.parseInst(funcObj)
						inst.Doc = instDoc

						funcObj.Body.Insts = append(funcObj.Body.Insts, inst)
						funcObj.Body.Objects = append(funcObj.Body.Objects, inst)
						continue
					}

					// 未知 token
					break
				}

				p.prog.Funcs = append(p.prog.Funcs, funcObj)
				p.prog.Objects = append(p.prog.Objects, funcObj)

				// 最后一个指令如果是注释, 则提到全局对象
				if len(funcObj.Body.Objects) > 0 {
					lastInstObj := funcObj.Body.Objects[len(funcObj.Body.Objects)-1]
					if commentObj, ok := lastInstObj.(*ast.CommentGroup); ok {
						funcObj.Body.Comments = funcObj.Body.Comments[:len(funcObj.Body.Comments)-1]
						funcObj.Body.Objects = funcObj.Body.Objects[:len(funcObj.Body.Objects)-1]

						p.prog.Comments = append(p.prog.Comments, commentObj)
						p.prog.Objects = append(p.prog.Objects, commentObj)
					}
				}

				// 删除结尾的空行
				for len(funcObj.Body.Objects) > 0 {
					lastInstObj := funcObj.Body.Objects[len(funcObj.Body.Objects)-1]
					if _, ok := lastInstObj.(*ast.BlankLine); ok {
						funcObj.Body.Objects = funcObj.Body.Objects[:len(funcObj.Body.Objects)-1]
					} else {
						break
					}
				}
			}

		case token.GAS_ALIGN:
			p.acceptToken(token.GAS_ALIGN)
			assert(p.gasSectionName != "")
			p.gasAlign = p.parseIntLit()
			p.gasSectionAlign[p.gasSectionName] = p.gasAlign

			switch p.gasSectionName {
			case ".text", ".init", ".fini":
				// 函数必须是4字节对齐
				if p.cpu == abi.X64Unix || p.cpu == abi.X64Windows {
					assert(p.gasAlign == 4)
				} else {
					assert(p.gasAlign == 2)
				}

			default:
				// 数据段只有: 1/2/4/8/16字节
				if p.cpu == abi.X64Unix || p.cpu == abi.X64Windows {
					assert(p.gasAlign == 1 || p.gasAlign == 2 || p.gasAlign == 4 || p.gasAlign == 8 || p.gasAlign == 16)
				} else {
					assert(p.gasAlign >= 0 && p.gasAlign <= 4)
				}
			}
			p.consumeSemicolonList()

		case token.IDENT:
			switch p.gasSectionName {
			case "":
				p.errorf(p.pos, "%s missing section name", p.lit)

			case ".data", ".radata", ".bss":
				beginPos := p.pos
				name := p.parseIdent()
				p.acceptToken(token.COLON)

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
					token.GAS_ASCII,
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

				// 验证初始化值的合法性, 填充类型和Size
				switch globalObj.TypeTok {
				case token.GAS_BYTE:
					globalObj.Type = token.I8
					globalObj.Size = 1
					v := int(globalObj.Init.Lit.ConstV.(int64))
					assert(v >= math.MinInt8 && v < math.MaxUint8)
				case token.GAS_SHORT:
					globalObj.Type = token.I16
					globalObj.Size = 2
					v := int(globalObj.Init.Lit.ConstV.(int64))
					assert(v >= math.MinInt16 && v < math.MaxUint16)
				case token.GAS_LONG:
					if p.cpu == abi.RISCV32 && globalObj.Init.Symbal != "" {
						globalObj.Type = token.I32
						globalObj.Size = 4
					} else {
						globalObj.Type = token.I32
						globalObj.Size = 4
						v := globalObj.Init.Lit.ConstV.(int64)
						assert(v >= math.MinInt32 && v < math.MaxUint32)
					}
				case token.GAS_QUAD:
					if globalObj.Init.Symbal != "" {
						globalObj.Type = token.I64
						globalObj.Size = 8
					} else {
						globalObj.Type = token.I64
						globalObj.Size = 8
						_ = globalObj.Init.Lit.ConstV.(int64)
					}
				case token.GAS_FLOAT:
					globalObj.Type = token.F32
					globalObj.Size = 4
					_ = globalObj.Init.Lit.ConstV.(float64)
				case token.GAS_DOUBLE:
					globalObj.Type = token.F64
					globalObj.Size = 8
					_ = globalObj.Init.Lit.ConstV.(float64)
				case token.GAS_ASCII:
					globalObj.Type = token.Bin
					globalObj.Size = len(globalObj.Init.Lit.ConstV.(string))
				case token.GAS_SKIP:
					globalObj.Type = token.Bin
					globalObj.Size = int(globalObj.Init.Lit.ConstV.(int64))
				case token.GAS_INCBIN:
					globalObj.Type = token.Bin
					filename := globalObj.Init.Lit.ConstV.(string)
					if fi, err := os.Lstat(filename); err != nil {
						p.errorf(p.pos, "%v %v failed: %v", token.GAS_INCBIN, filename, err)
					} else {
						const maxSize = 2 << 20
						if fi.Size() > maxSize {
							p.errorf(p.pos, "%v %v file size large than 2MB", token.GAS_INCBIN, filename)
						}
						globalObj.Size = int(fi.Size())
					}
				default:
					panic("unreachable")
				}

				p.prog.Globals = append(p.prog.Globals, globalObj)
				p.prog.Objects = append(p.prog.Objects, globalObj)

			case ".text", ".init", ".fini":
				// 函数已经提前解析
				p.errorf(p.pos, "unreachable")

			default:
				p.errorf(p.pos, "invalid section name: %s", p.gasSectionName)
			}

		case token.GLOBAL_zh, token.READONLY_zh:
			p.gasSectionName = ""
			p.gasAlign = 0

			globalObj := p.parseGlobal(p.tok)
			p.prog.Globals = append(p.prog.Globals, globalObj)
			p.prog.Objects = append(p.prog.Objects, globalObj)

		case token.FUNC_zh:
			p.gasSectionName = ""
			p.gasAlign = 0

			funcObj := p.parseFunc()
			p.prog.Funcs = append(p.prog.Funcs, funcObj)
			p.prog.Objects = append(p.prog.Objects, funcObj)

		default:
			p.errorf(p.pos, "unkonw token: %v", p.tok)
		}
	}

	// 收集信息导出符号信息
	for _, g := range p.prog.Globals {
		if p.gasGlobl[g.Name] {
			g.ExportName = g.Name
		}
	}
	for _, fn := range p.prog.Funcs {
		if p.gasGlobl[fn.Name] {
			fn.ExportName = fn.Name
		}
	}
}
