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
			fmt.Fprintf(&bufIns, "  // local %s: %s\n", x.Name, x.Type)
			fmt.Fprintf(&bufIns, "  %s var_%s = 0;\n", toCType(x.Type), toCName(x.Name))
		}
		fmt.Fprintln(&bufIns)
	}

	for _, ins := range fn.Body.Insts {
		if err := p.buildFunc_ins(&bufIns, fn.Type, &stk, ins, 1); err != nil {
			return err
		}
	}

	// 未使用的可以省略
	if stkMaxDepth := stk.MaxDepth(); stkMaxDepth > 0 {
		fmt.Fprintf(w, "  wasm_reg_t $reg[%d];\n", stkMaxDepth)
		fmt.Fprintln(w)
	}

	// 指令复制到 w
	io.Copy(w, &bufIns)
	return nil
}

func (p *wat2cWorker) buildFunc_ins(w io.Writer, fnType *ast.FuncType, stk *valueTypeStack, i ast.Instruction, level int) error {
	indent := strings.Repeat("  ", level)
	switch tok := i.Token(); tok {
	case token.INS_UNREACHABLE:
		stk.Nop()
		fmt.Fprintf(w, "%sabort(); // %s\n", indent, tok)
	case token.INS_NOP:
		stk.Nop()
		fmt.Fprintf(w, "%s(0); // %s\n", indent, tok)
	case token.INS_BLOCK:
		i := i.(ast.Ins_Block)
		fmt.Fprintf(w, "%s{ // block\n", indent)
		for _, ins := range i.List {
			if err := p.buildFunc_ins(w, fnType, stk, ins, level+1); err != nil {
				return err
			}
		}
		fmt.Fprintf(w, "%s}\n", indent)
		fmt.Fprintf(w, "L_%s_next:\n", toCName(i.Label))

	case token.INS_LOOP:
		stkLen := stk.Len()
		i := i.(ast.Ins_Loop)
		fmt.Fprintf(w, "L_%s_next:\n", toCName(i.Label))
		fmt.Fprintf(w, "%s{ // loop\n", indent)
		for _, ins := range i.List {
			if err := p.buildFunc_ins(w, fnType, stk, ins, level+1); err != nil {
				return err
			}
		}
		fmt.Fprintf(w, "%s}\n", indent)

		assert(stk.Len() == stkLen+len(i.Results))

	case token.INS_IF:
		stkLen := stk.Len()
		i := i.(ast.Ins_If)
		fmt.Fprintf(w, "%sif($reg[%d].i32) {\n", indent, stk.Len())
		for _, ins := range i.Body {
			if err := p.buildFunc_ins(w, fnType, stk, ins, level+1); err != nil {
				return err
			}
		}

		if len(i.Else) > 0 {
			fmt.Fprintf(w, "%s} else {\n", indent)
			for _, ins := range i.Else {
				if err := p.buildFunc_ins(w, fnType, stk, ins, level+1); err != nil {
					return err
				}
			}
		}
		fmt.Fprintf(w, "%s}\n", indent)
		fmt.Fprintf(w, "L_%s_next:\n", toCName(i.Label))

		assert(stk.Len() == stkLen+len(i.Results))

	case token.INS_ELSE:
		panic("unreachable")
	case token.INS_END:
		panic("unreachable")

	case token.INS_BR:
		stkLen := stk.Len()
		i := i.(ast.Ins_Br)
		stk.Nop()
		fmt.Fprintf(w, "%sgoto L_%s_next;\n", indent, toCName(i.X))
		assert(stk.Len() == stkLen)

	case token.INS_BR_IF:
		stkLen := stk.Len()
		i := i.(ast.Ins_BrIf)
		stk.Pop()
		fmt.Fprintf(w, "%sif($reg[%d].i32) { goto L_%s_next; }\n", indent, stk.Len(), toCName(i.X))
		assert(stk.Len() == stkLen-1)
	case token.INS_BR_TABLE:
		stkLen := stk.Len()

		// stk.popN(2)
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_RETURN:
		stkLen := stk.Len()

		for _, result := range fnType.Results {
			_ = result
			// todo: 等其他指令完整后再处理
			//if vt := stk.pop(); vt != result {
			//	panic(fmt.Sprintf("unreachable: %v != %v", vt, result))
			//}
		}
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_CALL:
		stkLen := stk.Len()

		// todo: 要准备的函数参数个数
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_CALL_INDIRECT:
		stkLen := stk.Len()

		stk.Pop()
		// todo: 要准备的函数参数个数
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_DROP:
		stkLen := stk.Len()

		stk.Pop()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_SELECT:
		stkLen := stk.Len()

		stk.PopN(3)
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_TYPED_SELECT:
		stkLen := stk.Len()

		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_LOCAL_GET:
		stkLen := stk.Len()

		stk.PushN(1)
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_LOCAL_SET:
		stkLen := stk.Len()

		//stk.pop()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_LOCAL_TEE:
		stkLen := stk.Len()

		//stk.push(stk.pop())
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_GLOBAL_GET:
		stkLen := stk.Len()

		stk.PushN(1)
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_GLOBAL_SET:
		stkLen := stk.Len()
		//stk.pop()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_TABLE_GET:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_TABLE_SET:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_LOAD:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LOAD:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_LOAD:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_LOAD:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_LOAD8_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)

	case token.INS_I32_LOAD8_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)

	case token.INS_I32_LOAD16_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)

	case token.INS_I32_LOAD16_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)

	case token.INS_I64_LOAD8_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LOAD8_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LOAD16_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LOAD16_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LOAD32_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LOAD32_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_STORE:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_STORE:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_STORE:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_STORE:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_STORE8:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_STORE16:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_STORE8:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_STORE16:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_STORE32:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_MEMORY_SIZE:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_MEMORY_GROW:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_CONST:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_CONST:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_CONST:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_CONST:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_EQZ:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_EQ:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_NE:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_LT_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_LT_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_GT_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_GT_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_LE_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_LE_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_GE_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_GE_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_EQZ:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_EQ:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_NE:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LT_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LT_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_GT_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_GT_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LE_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LE_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_GE_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_GE_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_EQ:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_NE:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_LT:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_GT:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_LE:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_GE:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_EQ:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_NE:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_LT:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_GT:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_LE:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_GE:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_CLZ:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_CTZ:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_POPCNT:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_ADD:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_SUB:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_MUL:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_DIV_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_DIV_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_REM_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_REM_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_AND:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_OR:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_XOR:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_SHL:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_SHR_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_SHR_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_ROTL:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_ROTR:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_CLZ:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_CTZ:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_POPCNT:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_ADD:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_SUB:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_MUL:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_DIV_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_DIV_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_REM_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_REM_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_AND:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_OR:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_XOR:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_SHL:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_SHR_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_SHR_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_ROTL:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_ROTR:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_ABS:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_NEG:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_CEIL:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_FLOOR:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_TRUNC:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_NEAREST:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_SQRT:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_ADD:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_SUB:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_MUL:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_DIV:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_MIN:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_MAX:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_COPYSIGN:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_ABS:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_NEG:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_CEIL:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_FLOOR:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_TRUNC:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_NEAREST:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_SQRT:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_ADD:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_SUB:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_MUL:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_DIV:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_MIN:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_MAX:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_COPYSIGN:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_WRAP_I64:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_TRUNC_F32_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_TRUNC_F32_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_TRUNC_F64_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_TRUNC_F64_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_EXTEND_I32_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_EXTEND_I32_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_TRUNC_F32_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_TRUNC_F32_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_TRUNC_F64_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_TRUNC_F64_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_CONVERT_I32_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_CONVERT_I32_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_CONVERT_I64_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_CONVERT_I64_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_DEMOTE_F64:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_CONVERT_I32_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_CONVERT_I32_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_CONVERT_I64_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_CONVERT_I64_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_PROMOTE_F32:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I32_REINTERPRET_F32:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_I64_REINTERPRET_F64:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F32_REINTERPRET_I32:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)
	case token.INS_F64_REINTERPRET_I64:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// todo: %T\n", indent, i)
		assert(stk.Len() == stkLen)

	default:
		panic(fmt.Sprintf("unreachable: %T", i))
	}
	return nil
}
