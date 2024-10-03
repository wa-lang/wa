// 版权 @2024 凹语言 作者。保留所有权利。

package wat2c

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *wat2cWorker) buildFunc_body(w io.Writer, fn *ast.Func) error {
	var stk valueTypeStack
	var bufIns bytes.Buffer

	if len(fn.Body.Locals) > 0 {
		for _, x := range fn.Body.Locals {
			fmt.Fprintf(&bufIns, "\t// local %s: %s\n", x.Name, x.Type)
			fmt.Fprintf(&bufIns, "\t%s var_%s = 0;\n", toCType(x.Type), toCName(x.Name))
		}
		fmt.Fprintln(&bufIns)
	}

	for _, ins := range fn.Body.Insts {
		if err := p.buildFunc_ins(&bufIns, fn.Type, &stk, ins, 1); err != nil {
			return err
		}
	}

	// 未使用的可以省略
	if stkMaxDepth := stk.maxDepth(); stkMaxDepth > 0 {
		fmt.Fprintf(w, "\tint64_t   $stk[%d]; // todo: fix stk size\n", stkMaxDepth)
		fmt.Fprintf(w, "\tconst int $stk_max = %d;\n", stkMaxDepth)
		fmt.Fprintf(w, "\tint       $stk_top = %d;\n", 0)
		fmt.Fprintln(w)
	}

	// 指令复制到 w
	io.Copy(w, &bufIns)
	return nil
}

func (p *wat2cWorker) buildFunc_ins(w io.Writer, fnType *ast.FuncType, stk *valueTypeStack, i ast.Instruction, level int) error {
	indent := strings.Repeat("\t", level)
	switch tok := i.Token(); tok {
	case token.INS_UNREACHABLE:
		stk.nop()
		fmt.Fprintf(w, "%sabort(); // %s\n", indent, tok)
	case token.INS_NOP:
		stk.nop()
		fmt.Fprintf(w, "%s(0); // %s\n", indent, tok)
	case token.INS_BLOCK:
		i := i.(ast.Ins_Block)
		fmt.Fprintf(w, "L_%s_begin:\n", toCName(i.Label))
		fmt.Fprintf(w, "%s{\n", indent)
		for _, ins := range i.List {
			if err := p.buildFunc_ins(w, fnType, stk, ins, level+1); err != nil {
				return err
			}
		}
		fmt.Fprintf(w, "%s}\n", indent)
		fmt.Fprintf(w, "L_%s_end:\n", toCName(i.Label))

		// 记录返回结果导致的栈变化
		for _, result := range i.Results {
			stk.push(result)
		}

	case token.INS_LOOP:
		i := i.(ast.Ins_Loop)
		fmt.Fprintf(w, "L_%s_begin:\n", toCName(i.Label))
		fmt.Fprintf(w, "%sfor(;;) {\n", indent)
		for _, ins := range i.List {
			if err := p.buildFunc_ins(w, fnType, stk, ins, level+1); err != nil {
				return err
			}
		}
		fmt.Fprintf(w, "%s}\n", indent)
		fmt.Fprintf(w, "L_%s_end:\n", toCName(i.Label))

		// 记录返回结果导致的栈变化
		for _, result := range i.Results {
			stk.push(result)
		}

	case token.INS_IF:
		i := i.(ast.Ins_If)
		fmt.Fprintf(w, "L_%s_begin:\n", toCName(i.Label))
		fmt.Fprintf(w, "%sif($stk[$stk_top]) {\n", indent)
		for _, ins := range i.Body {
			if err := p.buildFunc_ins(w, fnType, stk, ins, level+1); err != nil {
				return err
			}
		}

		if len(i.Else) > 0 {
			fmt.Fprintf(w, "%s} else {\n", indent)
			fmt.Fprintf(w, "L_%s_else:\n", toCName(i.Label))
			for _, ins := range i.Else {
				if err := p.buildFunc_ins(w, fnType, stk, ins, level+1); err != nil {
					return err
				}
			}
		}
		fmt.Fprintf(w, "%s}\n", indent)
		fmt.Fprintf(w, "L_%s_end:\n", toCName(i.Label))

		// 记录返回结果导致的栈变化
		for _, result := range i.Results {
			stk.push(result)
		}

	case token.INS_ELSE:
		panic("unreachable")
	case token.INS_END:
		panic("unreachable")

	case token.INS_BR:
		stk.nop()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_BR_IF:
		// stk.pop()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_BR_TABLE:
		// stk.popN(2)
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_RETURN:
		for _, result := range fnType.Results {
			_ = result
			// todo: 等其他指令完整后再处理
			//if vt := stk.pop(); vt != result {
			//	panic(fmt.Sprintf("unreachable: %v != %v", vt, result))
			//}
		}
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_CALL:
		// todo: 要准备的函数参数个数
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_CALL_INDIRECT:
		stk.pop()
		// todo: 要准备的函数参数个数
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_DROP:
		stk.pop()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_SELECT:
		stk.popN(3)
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_TYPED_SELECT:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_LOCAL_GET:
		stk.pushN(1)
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_LOCAL_SET:
		//stk.pop()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_LOCAL_TEE:
		//stk.push(stk.pop())
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_GLOBAL_GET:
		stk.pushN(1)
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_GLOBAL_SET:
		//stk.pop()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_TABLE_GET:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_TABLE_SET:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_LOAD:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_LOAD:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_LOAD:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_LOAD:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_LOAD8_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)

	case token.INS_I32_LOAD8_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)

	case token.INS_I32_LOAD16_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)

	case token.INS_I32_LOAD16_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)

	case token.INS_I64_LOAD8_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_LOAD8_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_LOAD16_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_LOAD16_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_LOAD32_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_LOAD32_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_STORE:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_STORE:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_STORE:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_STORE:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_STORE8:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_STORE16:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_STORE8:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_STORE16:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_STORE32:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_MEMORY_SIZE:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_MEMORY_GROW:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_CONST:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_CONST:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_CONST:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_CONST:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_EQZ:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_EQ:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_NE:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_LT_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_LT_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_GT_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_GT_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_LE_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_LE_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_GE_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_GE_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_EQZ:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_EQ:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_NE:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_LT_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_LT_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_GT_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_GT_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_LE_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_LE_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_GE_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_GE_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_EQ:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_NE:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_LT:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_GT:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_LE:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_GE:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_EQ:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_NE:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_LT:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_GT:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_LE:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_GE:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_CLZ:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_CTZ:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_POPCNT:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_ADD:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_SUB:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_MUL:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_DIV_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_DIV_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_REM_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_REM_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_AND:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_OR:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_XOR:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_SHL:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_SHR_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_SHR_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_ROTL:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_ROTR:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_CLZ:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_CTZ:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_POPCNT:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_ADD:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_SUB:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_MUL:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_DIV_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_DIV_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_REM_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_REM_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_AND:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_OR:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_XOR:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_SHL:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_SHR_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_SHR_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_ROTL:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_ROTR:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_ABS:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_NEG:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_CEIL:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_FLOOR:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_TRUNC:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_NEAREST:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_SQRT:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_ADD:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_SUB:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_MUL:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_DIV:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_MIN:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_MAX:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_COPYSIGN:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_ABS:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_NEG:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_CEIL:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_FLOOR:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_TRUNC:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_NEAREST:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_SQRT:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_ADD:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_SUB:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_MUL:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_DIV:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_MIN:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_MAX:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_COPYSIGN:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_WRAP_I64:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_TRUNC_F32_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_TRUNC_F32_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_TRUNC_F64_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_TRUNC_F64_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_EXTEND_I32_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_EXTEND_I32_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_TRUNC_F32_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_TRUNC_F32_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_TRUNC_F64_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_TRUNC_F64_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_CONVERT_I32_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_CONVERT_I32_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_CONVERT_I64_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_CONVERT_I64_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_DEMOTE_F64:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_CONVERT_I32_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_CONVERT_I32_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_CONVERT_I64_S:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_CONVERT_I64_U:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_PROMOTE_F32:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I32_REINTERPRET_F32:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_I64_REINTERPRET_F64:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F32_REINTERPRET_I32:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
	case token.INS_F64_REINTERPRET_I64:
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)

	default:
		panic(fmt.Sprintf("unreachable: %T", i))
	}
	return nil
}
