// 版权 @2024 凹语言 作者。保留所有权利。

package printer

import (
	"fmt"
	"strings"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *watPrinter) printFuncs() error {
	if p.isFuncEmpty() {
		return nil
	}
	for _, fn := range p.m.Funcs {
		fmt.Fprintf(p.w, "%s(func %s", p.indent, p.identOrIndex(fn.Name))

		if fn.ExportName != "" {
			fmt.Fprintf(p.w, " (export %q)", fn.ExportName)
		}

		if len(fn.Type.Params) > 0 {
			for _, x := range fn.Type.Params {
				if x.Name != "" {
					fmt.Fprintf(p.w, " (param %s %v)", p.identOrIndex(x.Name), x.Type)
				} else {
					fmt.Fprintf(p.w, " (param %v)", x.Type)
				}
			}
		}
		if len(fn.Type.Results) > 0 {
			fmt.Fprintf(p.w, " (result")
			for _, x := range fn.Type.Results {
				fmt.Fprintf(p.w, " %v", x)
			}
			fmt.Fprint(p.w, ")")
		}

		if len(fn.Body.Locals) != 0 || len(fn.Body.Insts) != 0 {
			fmt.Fprintln(p.w)
			p.printFuncs_body(fn)
			fmt.Fprintf(p.w, "%s)\n", p.indent)
		} else {
			fmt.Fprintf(p.w, ")\n")
		}
	}

	return nil
}

func (p *watPrinter) printFuncs_body(fn *ast.Func) {
	for _, local := range fn.Body.Locals {
		fmt.Fprint(p.w, p.indent+p.indent)
		fmt.Fprint(p.w, "(local")
		if local.Name != "" {
			fmt.Fprintf(p.w, " %s", p.identOrIndex(local.Name))
		}
		fmt.Fprintf(p.w, " %v)\n", local.Type)
	}

	for _, ins := range fn.Body.Insts {
		p.printFuncs_body_ins(fn, ins, 0)
	}

	return
}

func (p *watPrinter) printFuncs_indent(blkLevel int) {
	fmt.Fprint(p.w, p.indent+p.indent)
	for i := 0; i < blkLevel; i++ {
		fmt.Fprint(p.w, p.indent)
	}
}

func (p *watPrinter) printFuncs_body_ins(fn *ast.Func, ins ast.Instruction, blkLevel int) {
	p.printFuncs_indent(blkLevel)

	switch tok := ins.Token(); tok {
	default:
		panic("unreachable")
	case token.INS_UNREACHABLE:
		fmt.Fprintln(p.w, tok)
	case token.INS_NOP:
		fmt.Fprintln(p.w, tok)

	case token.INS_BLOCK:
		insBlock := ins.(ast.Ins_Block)

		fmt.Fprint(p.w, tok)
		if s := insBlock.Label; s != "" {
			fmt.Fprintf(p.w, " %s", p.identOrIndex(s))
		}
		if len(insBlock.Results) > 0 {
			fmt.Fprint(p.w, "(result")
			for _, x := range insBlock.Results {
				fmt.Fprint(p.w, " ", x)
			}
			fmt.Fprint(p.w, ")")
		}
		fmt.Fprintln(p.w)

		for _, x := range insBlock.List {
			p.printFuncs_body_ins(fn, x, blkLevel+1)
		}

		p.printFuncs_indent(blkLevel)
		fmt.Fprintln(p.w, token.INS_END)

	case token.INS_LOOP:
		insLoop := ins.(ast.Ins_Loop)

		fmt.Fprint(p.w, tok)
		if s := insLoop.Label; s != "" {
			fmt.Fprintf(p.w, " %s", p.identOrIndex(s))
		}
		if len(insLoop.Results) > 0 {
			fmt.Fprint(p.w, "(result")
			for _, x := range insLoop.Results {
				fmt.Fprint(p.w, " ", x)
			}
			fmt.Fprint(p.w, ")")
		}
		fmt.Fprintln(p.w)

		for _, x := range insLoop.List {
			p.printFuncs_body_ins(fn, x, blkLevel+1)
		}

		p.printFuncs_indent(blkLevel)
		fmt.Fprintln(p.w, token.INS_END)

	case token.INS_IF:
		insIf := ins.(ast.Ins_If)

		fmt.Fprint(p.w, tok)
		if s := insIf.Label; s != "" {
			fmt.Fprintf(p.w, " %s", p.identOrIndex(s))
		}
		if len(insIf.Results) > 0 {
			fmt.Fprint(p.w, "(result")
			for _, x := range insIf.Results {
				fmt.Fprint(p.w, " ", x)
			}
			fmt.Fprint(p.w, ")")
		}
		fmt.Fprintln(p.w)

		for _, x := range insIf.Body {
			p.printFuncs_body_ins(fn, x, blkLevel+1)
		}

		if len(insIf.Else) > 0 {
			p.printFuncs_indent(blkLevel)
			fmt.Fprintln(p.w, token.INS_ELSE)

			for _, x := range insIf.Else {
				p.printFuncs_body_ins(fn, x, blkLevel+1)
			}
		}

		p.printFuncs_indent(blkLevel)
		fmt.Fprintln(p.w, token.INS_END)

	case token.INS_ELSE:
		panic("unreachable")
	case token.INS_END:
		panic("unreachable")

	case token.INS_BR:
		fmt.Fprintln(p.w, tok, p.identOrIndex(ins.(ast.Ins_Br).X))
	case token.INS_BR_IF:
		fmt.Fprintln(p.w, tok, p.identOrIndex(ins.(ast.Ins_BrIf).X))
	case token.INS_BR_TABLE:
		var sb strings.Builder
		for _, s := range ins.(ast.Ins_BrTable).XList {
			sb.WriteRune(' ')
			sb.WriteString(p.identOrIndex(s))
		}
		fmt.Fprintln(p.w, tok, sb.String())
	case token.INS_RETURN:
		fmt.Fprintln(p.w, tok)
	case token.INS_CALL:
		fmt.Fprintln(p.w, tok, p.identOrIndex(ins.(ast.Ins_Call).X))
	case token.INS_CALL_INDIRECT:
		insCallIndirect := ins.(ast.Ins_CallIndirect)
		fmt.Fprint(p.w, tok)
		if s := insCallIndirect.TableIdx; s != "" {
			fmt.Fprintf(p.w, " %s", p.identOrIndex(s))
		}
		fmt.Fprintf(p.w, " (type %s)", p.identOrIndex(insCallIndirect.TypeIdx))
		fmt.Fprintln(p.w)

	case token.INS_DROP:
		fmt.Fprintln(p.w, tok)
	case token.INS_SELECT:
		insSelect := ins.(ast.Ins_Select)
		if insSelect.ResultTyp != 0 {
			fmt.Fprintln(p.w, tok, "(result "+insSelect.ResultTyp.String()+")")
		} else {
			fmt.Fprintln(p.w, tok)
		}
	case token.INS_LOCAL_GET:
		fmt.Fprintln(p.w, tok, p.identOrIndex(ins.(ast.Ins_LocalGet).X))
	case token.INS_LOCAL_SET:
		fmt.Fprintln(p.w, tok, p.identOrIndex(ins.(ast.Ins_LocalSet).X))
	case token.INS_LOCAL_TEE:
		fmt.Fprintln(p.w, tok, p.identOrIndex(ins.(ast.Ins_LocalTee).X))
	case token.INS_GLOBAL_GET:
		fmt.Fprintln(p.w, tok, p.identOrIndex(ins.(ast.Ins_GlobalGet).X))
	case token.INS_GLOBAL_SET:
		fmt.Fprintln(p.w, tok, p.identOrIndex(ins.(ast.Ins_GlobalSet).X))
	case token.INS_TABLE_GET:
		fmt.Fprintln(p.w, tok, p.identOrIndex(ins.(ast.Ins_TableGet).TableIdx))
	case token.INS_TABLE_SET:
		fmt.Fprintln(p.w, tok, p.identOrIndex(ins.(ast.Ins_TableSet).TableIdx))
	case token.INS_I32_LOAD:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I32Load)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 4 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I64_LOAD:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I64Load)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 8 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_F32_LOAD:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_F32Load)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 4 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_F64_LOAD:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_F64Load)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 8 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I32_LOAD8_S:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I32Load8S)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 4 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I32_LOAD8_U:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I32Load8U)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 4 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I32_LOAD16_S:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I32Load16S)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 2 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I32_LOAD16_U:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I32Load16U)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 2 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I64_LOAD8_S:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I64Load8S)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 8 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I64_LOAD8_U:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I64Load8U)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 8 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I64_LOAD16_S:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I64Load16S)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 8 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I64_LOAD16_U:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I64Load16U)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 2 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I64_LOAD32_S:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I64Load32S)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 4 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I64_LOAD32_U:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I64Load32U)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 8 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I32_STORE:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I32Store)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 4 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I64_STORE:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I64Store)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 2 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_F32_STORE:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_F32Store)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 4 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_F64_STORE:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_F64Store)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 8 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I32_STORE8:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I32Store8)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 4 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I32_STORE16:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I32Store16)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 2 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I64_STORE8:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I64Store8)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 8 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I64_STORE16:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I64Store16)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 2 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_I64_STORE32:
		fmt.Fprint(p.w, tok)
		insLoad := ins.(ast.Ins_I64Store32)
		if x := insLoad.Offset; x != 0 {
			fmt.Fprintf(p.w, " offset=%d", x)
		}
		if x := insLoad.Align; x != 8 {
			fmt.Fprintf(p.w, " align=%d", x)
		}
		fmt.Fprintln(p.w)
	case token.INS_MEMORY_SIZE:
		fmt.Fprintln(p.w, tok)
	case token.INS_MEMORY_GROW:
		fmt.Fprintln(p.w, tok)
	case token.INS_MEMORY_INIT:
		fmt.Fprintln(p.w, tok, ins.(ast.Ins_MemoryInit).DataIdx)
	case token.INS_MEMORY_COPY:
		fmt.Fprintln(p.w, tok)
	case token.INS_MEMORY_FILL:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_CONST:
		fmt.Fprintln(p.w, tok, ins.(ast.Ins_I32Const).X)
	case token.INS_I64_CONST:
		fmt.Fprintln(p.w, tok, ins.(ast.Ins_I64Const).X)
	case token.INS_F32_CONST:
		fmt.Fprintln(p.w, tok, ins.(ast.Ins_F32Const).X)
	case token.INS_F64_CONST:
		fmt.Fprintln(p.w, tok, ins.(ast.Ins_F64Const).X)
	case token.INS_I32_EQZ:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_EQ:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_NE:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_LT_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_LT_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_GT_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_GT_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_LE_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_LE_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_GE_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_GE_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_EQZ:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_EQ:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_NE:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_LT_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_LT_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_GT_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_GT_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_LE_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_LE_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_GE_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_GE_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_EQ:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_NE:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_LT:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_GT:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_LE:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_GE:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_EQ:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_NE:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_LT:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_GT:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_LE:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_GE:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_CLZ:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_CTZ:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_POPCNT:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_ADD:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_SUB:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_MUL:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_DIV_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_DIV_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_REM_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_REM_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_AND:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_OR:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_XOR:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_SHL:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_SHR_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_SHR_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_ROTL:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_ROTR:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_CLZ:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_CTZ:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_POPCNT:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_ADD:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_SUB:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_MUL:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_DIV_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_DIV_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_REM_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_REM_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_AND:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_OR:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_XOR:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_SHL:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_SHR_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_SHR_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_ROTL:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_ROTR:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_ABS:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_NEG:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_CEIL:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_FLOOR:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_TRUNC:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_NEAREST:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_SQRT:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_ADD:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_SUB:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_MUL:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_DIV:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_MIN:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_MAX:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_COPYSIGN:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_ABS:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_NEG:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_CEIL:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_FLOOR:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_TRUNC:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_NEAREST:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_SQRT:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_ADD:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_SUB:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_MUL:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_DIV:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_MIN:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_MAX:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_COPYSIGN:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_WRAP_I64:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_TRUNC_F32_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_TRUNC_F32_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_TRUNC_F64_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_TRUNC_F64_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_EXTEND_I32_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_EXTEND_I32_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_TRUNC_F32_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_TRUNC_F32_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_TRUNC_F64_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_TRUNC_F64_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_CONVERT_I32_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_CONVERT_I32_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_CONVERT_I64_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_CONVERT_I64_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_DEMOTE_F64:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_CONVERT_I32_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_CONVERT_I32_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_CONVERT_I64_S:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_CONVERT_I64_U:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_PROMOTE_F32:
		fmt.Fprintln(p.w, tok)
	case token.INS_I32_REINTERPRET_F32:
		fmt.Fprintln(p.w, tok)
	case token.INS_I64_REINTERPRET_F64:
		fmt.Fprintln(p.w, tok)
	case token.INS_F32_REINTERPRET_I32:
		fmt.Fprintln(p.w, tok)
	case token.INS_F64_REINTERPRET_I64:
		fmt.Fprintln(p.w, tok)
	}
}

func (p *watPrinter) identOrIndex(idOrIdx string) string {
	if ch := idOrIdx[0]; ch >= '0' && ch <= '9' {
		return idOrIdx
	} else {
		return "$" + idOrIdx
	}
}
