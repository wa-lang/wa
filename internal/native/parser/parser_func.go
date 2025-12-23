// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"fmt"
	"sort"
	"strings"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

// 函数签名只用于文档注释, 不做语义检查

// func $add[prop1=val1,prop2=val2](%a:i32, %b:i32, %c:i32) => f64 {
//     local %d: i32 # 局部变量必须先声明, i32 大小的空间
//     # 指令
// Loop:
// }

// 函数 $add(%a:普整, %b:普整, %c:普整) => 双精:
//     local %d: i32 # 局部变量必须先声明, i32 大小的空间
//     # 指令
// Loop:
// 完毕

func (p *parser) parseFunc(tok token.Token) *ast.Func {
	p.fnArgRet.Reset()

	fn := &ast.Func{
		Pos:  p.pos,
		Tok:  tok,
		Body: new(ast.FuncBody),
	}

	fn.Doc = p.parseDocComment(&p.prog.Comments, fn.Pos)
	if fn.Doc != nil {
		p.prog.Objects = p.prog.Objects[:len(p.prog.Objects)-1]
	}

	fn.Tok = p.acceptTokenAorB(token.FUNC, token.FUNC_zh)
	fn.Name = p.parseIdent()

	if p.tok == token.LBRACK {
		p.parseFunc_prop(fn)
	}

	if p.tok == token.LPAREN {
		p.parseFunc_args(fn)
	}
	if p.tok == token.ARROW {
		p.parseFunc_return(fn)
	}

	p.parseFunc_body(fn)
	p.consumeSemicolonList()

	return fn
}

func (p *parser) parseFunc_prop(fn *ast.Func) {
	p.acceptToken(token.LBRACK)
	defer p.acceptToken(token.RBRACK)

	for p.tok == token.IDENT {
		propKey := p.parseIdent()
		propVal := ""

		// 检查属性是否重复定义
		for _, s := range fn.Prop {
			if s == propKey || strings.HasPrefix(s, propKey+"=") {
				p.errorf(p.pos, "prop %s exists", propKey)
			}
		}

		if p.tok == token.ASSIGN {
			p.acceptToken(token.ASSIGN)
			propVal = p.lit
			p.next()
		}
		if propVal != "" {
			fn.Prop = append(fn.Prop, fmt.Sprintf("%s=%s", propKey, propVal))
		} else {
			fn.Prop = append(fn.Prop, propKey)
		}
	}
}

func (p *parser) parseFunc_args(fn *ast.Func) {
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)

	for p.tok == token.IDENT {
		argPos := p.pos
		argName := p.parseIdent()

		p.acceptToken(token.COLON)

		arg := &ast.Local{
			Pos:  argPos,
			Name: argName,
		}

		switch p.tok {
		case token.I32, token.I64, token.F32, token.F64:
			arg.Type = p.tok
			p.acceptToken(p.tok)
		case token.I32_zh, token.I64_zh, token.F32_zh, token.F64_zh:
			arg.Type = p.tok
			p.acceptToken(p.tok)
		default:
			p.errorf(p.pos, "expect argument type(i32/i64/f32/f64/ptr), got %v", p.tok)
		}

		arg.Reg, arg.Off = p.fnArgRet.AllocArg(arg.Type)
		fn.Type.Args = append(fn.Type.Args, arg)

		if p.tok == token.COMMA {
			p.acceptToken(token.COMMA)
			continue
		} else {
			break
		}
	}
}

func (p *parser) parseFunc_return(fn *ast.Func) {
	p.acceptToken(token.ARROW)

	// 一组返回值
	if p.tok == token.LPAREN {
		p.acceptToken(token.LPAREN)
		defer p.acceptToken(token.RPAREN)

		for p.tok == token.IDENT {
			retPos := p.pos
			retName := p.parseIdent()

			p.acceptToken(token.COLON)

			ret := &ast.Local{
				Pos:  retPos,
				Name: retName,
			}

			switch p.tok {
			case token.I32, token.I64, token.F32, token.F64:
				ret.Type = p.tok
				p.acceptToken(p.tok)
			case token.I32_zh, token.I64_zh, token.F32_zh, token.F64_zh:
				ret.Type = p.tok
				p.acceptToken(p.tok)
			default:
				p.errorf(p.pos, "expect return type(i32/i64/f32/f64/ptr), got %v", p.tok)
			}

			ret.Reg, ret.Off = p.fnArgRet.AllocArg(ret.Type)
			fn.Type.Return = append(fn.Type.Args, ret)

			if p.tok == token.COMMA {
				p.acceptToken(token.COMMA)
				continue
			} else {
				break
			}
		}
		return
	}

	switch p.tok {
	case token.I32, token.I64, token.F32, token.F64:
		retReg, retOff := p.fnArgRet.AllocArg(p.tok)
		fn.Type.Return = append(fn.Type.Args, &ast.Local{
			Pos:  p.pos,
			Type: p.tok,
			Reg:  retReg,
			Off:  retOff,
		})
		p.acceptToken(p.tok)
	case token.I32_zh, token.I64_zh, token.F32_zh, token.F64_zh:
		retReg, retOff := p.fnArgRet.AllocArg(p.tok)
		fn.Type.Return = append(fn.Type.Args, &ast.Local{
			Pos:  p.pos,
			Type: p.tok,
			Reg:  retReg,
			Off:  retOff,
		})
		p.acceptToken(p.tok)
	default:
		p.errorf(p.pos, "expect return type(i32/i64/f32/f64/ptr), got %v", p.tok)
	}
}

func (p *parser) parseFunc_body(fn *ast.Func) {
	assert(p.cpu == abi.RISCV64 || p.cpu == abi.RISCV32 || p.cpu == abi.LOONG64)

	fn.Body.Pos = p.pos

	if fn.Tok == token.FUNC_zh {
		p.acceptToken(token.COLON)
		defer p.acceptToken(token.END_zh)
	} else {
		p.acceptToken(token.LBRACE)
		defer p.acceptToken(token.RBRACE)
	}

Loop:
	for {
		switch p.tok {
		case token.EOF, token.ILLEGAL:
			break Loop
		case token.RBRACE, token.END_zh:
			break Loop
		case token.COMMENT:
			commentObj := p.parseCommentGroup(false)
			fn.Body.Comments = append(fn.Body.Comments, commentObj)
			fn.Body.Objects = append(fn.Body.Objects, commentObj)
		case token.LOCAL, token.LOCAL_zh:
			if len(fn.Body.Insts) > 0 {
				p.errorf(p.pos, "local must before the instruction list")
			}
			localObj := p.parseFunc_body_local(fn)
			fn.Body.Locals = append(fn.Body.Locals, localObj)
			fn.Body.Objects = append(fn.Body.Objects, localObj)

		default:

			// 给局部变量在栈上预留空间预留
			// 为了减少栈内存碎片, 局部变量会重新排序
			{
				locals := append([]*ast.Local{}, fn.Body.Locals...)
				sort.Slice(locals, func(i, j int) bool {
					return locals[i].Type.NumberTypeSize() < locals[j].Type.NumberTypeSize()
				})

				// 局部变量分配
				for _, local := range locals {
					local.Off = p.fnArgRet.AllocLocal(local.Type, local.Cap)
				}
			}

			if p.tok == token.IDENT || p.tok.IsAs() {
				inst := p.parseInst(fn)
				fn.Body.Insts = append(fn.Body.Insts, inst)
				fn.Body.Objects = append(fn.Body.Objects, inst)
			} else {
				p.errorf(p.pos, "unknow as %v", p.tok)
			}
		}
	}
}

func (p *parser) parseFunc_body_local(fn *ast.Func) *ast.Local {
	local := &ast.Local{Pos: p.pos}

	local.Doc = p.parseDocComment(&fn.Body.Comments, local.Pos)
	if local.Doc != nil {
		fn.Body.Objects = fn.Body.Objects[:len(fn.Body.Objects)-1]
	}

	local.Tok = p.acceptTokenAorB(token.LOCAL, token.LOCAL_zh)
	local.Name = p.parseIdent()
	p.acceptToken(token.COLON)

	// 解析局部变量的容量
	if p.tok == token.LBRACK {
		p.acceptToken(token.LBRACK)
		local.Cap = p.parseIntLit()
		p.acceptToken(token.RBRACK)
	} else {
		local.Cap = 1
	}

	switch p.tok {
	case token.I32, token.I64,
		token.U32, token.U64,
		token.F32, token.F64,
		token.I32_zh, token.I64_zh,
		token.U32_zh, token.U64_zh,
		token.F32_zh, token.F64_zh:
		local.Type = p.tok
		p.acceptToken(p.tok)
	default:
		p.errorf(p.pos, "expect local type(i32/i64/f32/f64/ptr), got %v", p.tok)
	}

	p.consumeSemicolonList()
	return local
}
