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

	p.labelScope = nil

	if len(fn.Body.Locals) > 0 {
		for _, x := range fn.Body.Locals {
			localName := toCName(x.Name)
			p.localNames = append(p.localNames, localName)
			if localName != x.Name {
				fmt.Fprintf(&bufIns, "  val_t %s = {0}; // %s\n", localName, x.Name)
			} else {
				fmt.Fprintf(&bufIns, "  val_t %s = {0};\n", localName)
			}
		}
		fmt.Fprintln(&bufIns)
	}

	for i, ins := range fn.Body.Insts {
		if err := p.buildFunc_ins(&bufIns, fn, &stk, ins, 1); err != nil {
			return err
		}
		// 手动补充最后一个 return
		if i == len(fn.Body.Insts)-1 && ins.Token() != token.INS_RETURN {
			insReturn := ast.Ins_Return{OpToken: ast.OpToken(token.INS_RETURN)}
			if err := p.buildFunc_ins(&bufIns, fn, &stk, insReturn, 1); err != nil {
				return err
			}
		}
	}

	// 未使用的可以省略
	if stkMaxDepth := stk.MaxDepth(); stkMaxDepth > 0 {
		fmt.Fprintf(w, "  fn_%s_ret_t $result;\n", toCName(fn.Name))
		fmt.Fprintf(w, "  val_t $R[%d];\n", stkMaxDepth)
		fmt.Fprintf(w, "  u32_t $R_u32;\n")
		fmt.Fprintf(w, "  u16_t $R_u16;\n")
		fmt.Fprintf(w, "  u8_t  $R_u8;\n")
		fmt.Fprintln(w)
	} else {
		fmt.Fprintf(w, "  fn_%s_ret_t $result;\n", toCName(fn.Name))
		fmt.Fprintf(w, "  val_t $R[%d];\n", 1)
		fmt.Fprintf(w, "  u32_t $R_u32;\n")
		fmt.Fprintf(w, "  u16_t $R_u16;\n")
		fmt.Fprintf(w, "  u8_t  $R_u8;\n")
		fmt.Fprintln(w)
	}

	// 指令复制到 w
	io.Copy(w, &bufIns)
	return nil
}

func (p *wat2cWorker) buildFunc_ins(w io.Writer, fn *ast.Func, stk *valueTypeStack, i ast.Instruction, level int) error {
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
		p.enterLabelScope(i.Label)
		defer p.leaveLabelScope()

		if i.Label != "" {
			fmt.Fprintf(w, "%s{ // block $%s\n", indent, i.Label)
		} else {
			fmt.Fprintf(w, "%s{ // block\n", indent)
		}
		for _, ins := range i.List {
			if err := p.buildFunc_ins(w, fn, stk, ins, level+1); err != nil {
				return err
			}
		}
		fmt.Fprintf(w, "%s}\n", indent)
		if i.Label != "" {
			fmt.Fprintf(w, "L_%s_next:\n", toCName(i.Label))
		}
	case token.INS_LOOP:
		stkLen := stk.Len()
		i := i.(ast.Ins_Loop)
		p.enterLabelScope(i.Label)
		defer p.leaveLabelScope()

		if i.Label != "" {
			fmt.Fprintf(w, "L_%s_next:\n", toCName(i.Label))
		}
		if i.Label != "" {
			fmt.Fprintf(w, "%s{ // loop $%s\n", indent, i.Label)
		} else {
			fmt.Fprintf(w, "%s{ // loop\n", indent)
		}
		for _, ins := range i.List {
			if err := p.buildFunc_ins(w, fn, stk, ins, level+1); err != nil {
				return err
			}
		}
		fmt.Fprintf(w, "%s}\n", indent)

		assert(stk.Len() == stkLen+len(i.Results))

	case token.INS_IF:
		stkLen := stk.Len()
		i := i.(ast.Ins_If)
		p.enterLabelScope(i.Label)
		defer p.leaveLabelScope()

		fmt.Fprintf(w, "%sif($R[%d].i32) {\n", indent, stk.Len())
		for _, ins := range i.Body {
			if err := p.buildFunc_ins(w, fn, stk, ins, level+1); err != nil {
				return err
			}
		}

		if len(i.Else) > 0 {
			fmt.Fprintf(w, "%s} else {\n", indent)
			for _, ins := range i.Else {
				if err := p.buildFunc_ins(w, fn, stk, ins, level+1); err != nil {
					return err
				}
			}
		}
		fmt.Fprintf(w, "%s}\n", indent)
		if i.Label != "" {
			fmt.Fprintf(w, "L_%s_next:\n", toCName(i.Label))
		}
		assert(stk.Len() == stkLen+len(i.Results))

	case token.INS_ELSE:
		panic("unreachable")
	case token.INS_END:
		panic("unreachable")

	case token.INS_BR:
		stkLen := stk.Len()
		i := i.(ast.Ins_Br)
		assert(i.X != "")
		stk.Nop()
		fmt.Fprintf(w, "%sgoto L_%s_next;\n", indent, toCName(p.findLabelName(i.X)))
		assert(stk.Len() == stkLen)

	case token.INS_BR_IF:
		stkLen := stk.Len()
		i := i.(ast.Ins_BrIf)
		assert(i.X != "")
		stk.Pop()
		fmt.Fprintf(w, "%sif($R[%d].i32) { goto L_%s_next; }\n",
			indent, stk.Len(), toCName(p.findLabelName(i.X)),
		)
		assert(stk.Len() == stkLen-1)
	case token.INS_BR_TABLE:
		stkLen := stk.Len()
		i := i.(ast.Ins_BrTable)
		assert(len(i.XList) > 1)
		fmt.Fprintf(w, "%sswitch($R[%d].i32) {\n", indent, stk.TopIdx())
		stk.Pop()
		for k := 0; k < len(i.XList)-1; k++ {
			assert(i.XList[k] != "")
			fmt.Fprintf(w, "%scase %d: goto L_%s_next;\n",
				indent, k, toCName(p.findLabelName(i.XList[k])),
			)
		}
		fmt.Fprintf(w, "%sdefault: goto L_%s_next;\n",
			indent, toCName(p.findLabelName(i.XList[len(i.XList)-1])),
		)
		fmt.Fprintf(w, "%s}\n", indent)
		assert(stk.Len() == stkLen-1)
	case token.INS_RETURN:
		stkLen := stk.Len()
		assert(stkLen == len(fn.Type.Results))
		for i := 0; i < len(fn.Type.Results); i++ {
			fmt.Fprintf(w, "%s$result.$R[%d] = $R[%d];\n", indent, i, i)
		}
		fmt.Fprintf(w, "%sreturn $result;\n", indent)
		stk.PopN(len(fn.Type.Results))
		assert(stk.Len() == 0)
	case token.INS_CALL:
		stkLen := stk.Len()
		i := i.(ast.Ins_Call)

		fnType := p.findFuncType(i.X)
		assert(stkLen >= len(fnType.Params))

		stk.PopN(len(fnType.Params))
		resultOff := stk.Len()
		stk.PushN(len(fnType.Results))

		// 需要定义临时变量保存返回值
		if len(fnType.Results) > 0 {
			fmt.Fprintf(w, "%s{ fn_%s_ret_t $ret = fn_%s(", indent, toCName(i.X), toCName(i.X))
		} else {
			fmt.Fprintf(w, "%sfn_%s(", indent, toCName(i.X))
		}
		for i := 0; i < len(fnType.Params); i++ {
			if i > 0 {
				fmt.Fprintf(w, ", ")
			}
			fmt.Fprintf(w, "$R[%d]", resultOff+i)
		}
		fmt.Fprintf(w, ");")

		// 复制到当前stk
		if len(fnType.Results) > 0 {
			fmt.Fprintf(w, " memcpy(&$R[%d], &$ret, sizeof($ret)); }\n", resultOff)
		} else {
			fmt.Fprintf(w, "\n")
		}

	case token.INS_CALL_INDIRECT:
		stkLen := stk.Len()
		i := i.(ast.Ins_CallIndirect)

		fnType := p.findType(i.TypeIdx)
		assert(stkLen >= len(fnType.Params))

		tableAddr := stk.TopIdx()
		stk.PopN(len(fnType.Params) + 1)
		resultOff := stk.Len()
		stk.PushN(len(fnType.Results))

		// 生成要定义函数的类型
		fmt.Fprintf(w, "%s{\n", indent)
		{
			fmt.Fprintf(w, "%stypedef struct { val_t $R[%d]; } fn_ret_t;\n",
				indent, len(fnType.Results),
			)
			fmt.Fprintf(w, "%stypedef fn_ret_t (*fn_t)(", indent)
			if len(fnType.Params) > 0 {
				for i, x := range fnType.Params {
					if i > 0 {
						fmt.Fprintf(w, ", ")
					}
					if x.Name != "" {
						fmt.Fprintf(w, "val_t %v", toCName(x.Name))
					} else {
						fmt.Fprintf(w, "val_t $arg%d", i)
					}
				}
			}
			fmt.Fprintf(w, ");\n")

			// 调用函数
			if len(fnType.Results) > 0 {
				fmt.Fprintf(w, "%sfn_ret_t $ret = ((fn_t)(wasm_table[$R[%d].i32]))(&$R[%d]",
					indent, tableAddr, resultOff,
				)
			} else {
				fmt.Fprintf(w, "%s((fn_t)(wasm_table[$R[%d].i32]))(&$R[%d]",
					indent, tableAddr, resultOff,
				)
			}
			for i := 0; i < len(fnType.Params); i++ {
				fmt.Fprintf(w, ", $R[%d]", resultOff+i)
			}
			fmt.Fprintf(w, ");\n")

			// 保存返回值
			if len(fnType.Results) > 0 {
				fmt.Fprintf(w, "%smemcpy(&$R[%d], &$ret, sizeof($ret));\n", indent, resultOff)
			}
		}

		fmt.Fprintf(w, "%s}\n", indent)

	case token.INS_DROP:
		stkLen := stk.Len()
		sp := stk.TopIdx()
		stk.Pop()
		fmt.Fprintf(w, "%s$R[%d]; // drop\n", indent, sp)
		assert(stk.Len() == stkLen)
	case token.INS_SELECT:
		stkLen := stk.Len()
		assert(stkLen >= 3)
		sp := stk.TopIdx()
		stk.PopN(3)
		fmt.Fprintf(w, "%s$R[%d-2] = $R[%d].i32? $R[%d-1]: $R[%d-2];\n",
			indent, sp, sp, sp, sp,
		)
		stk.PushN(1)
		assert(stk.Len() == stkLen-3+1)
	case token.INS_LOCAL_GET:
		stkLen := stk.Len()
		i := i.(ast.Ins_LocalGet)
		stk.PushN(1)
		fmt.Fprintf(w, "%s$R[%d] = %s;\n", indent, stk.TopIdx(), p.findLocalName(fn, i.X))
		assert(stk.Len() == stkLen+1)
	case token.INS_LOCAL_SET:
		stkLen := stk.Len()
		i := i.(ast.Ins_LocalSet)
		fmt.Fprintf(w, "%s%s = $R[%d];\n", indent, p.findLocalName(fn, i.X), stk.TopIdx())
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_LOCAL_TEE:
		stkLen := stk.Len()
		i := i.(ast.Ins_LocalTee)
		fmt.Fprintf(w, "%s%s = $R[%d];\n", indent, p.findLocalName(fn, i.X), stk.TopIdx())
		assert(stk.Len() == stkLen)
	case token.INS_GLOBAL_GET:
		stkLen := stk.Len()
		i := i.(ast.Ins_GlobalGet)
		fmt.Fprintf(w, "%s$R[%d] = var_%s;\n", indent, stk.TopIdx(), toCName(i.X))
		assert(stk.Len() == stkLen+1)
	case token.INS_GLOBAL_SET:
		stkLen := stk.Len()
		i := i.(ast.Ins_GlobalSet)
		fmt.Fprintf(w, "%svar_%s = $R[%d];\n", indent, toCName(i.X), stk.TopIdx())
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_TABLE_GET:
		stkLen := stk.Len()
		sp := stk.TopIdx()
		fmt.Fprintf(w, "%s$R[%d].i64 = (i64_t)(wasm_table[$R[%d].i32]);\n", indent, sp, sp)
		assert(stk.Len() == stkLen)
	case token.INS_TABLE_SET:
		stkLen := stk.Len()
		sp := stk.TopIdx()
		fmt.Fprintf(w, "%swasm_table[$R[%d].i32] = (uintptr_t)($R[%d].i64);\n", indent, sp, sp)
		assert(stk.Len() == stkLen)
	case token.INS_I32_LOAD:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		i := i.(ast.Ins_I32Load)
		fmt.Fprintf(w, "%smemcpy(&$R[%d].i32, &wasm_memoy[$R[%d].i32+%d], 4);\n",
			indent, stk.TopIdx(), stk.TopIdx(), i.Offset,
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LOAD:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		i := i.(ast.Ins_I64Load)
		fmt.Fprintf(w, "%smemcpy(&$R[%d].i64, &wasm_memoy[$R[%d].i32+%d], 8);\n",
			indent, stk.TopIdx(), stk.TopIdx(), i.Offset,
		)
		assert(stk.Len() == stkLen)
	case token.INS_F32_LOAD:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		i := i.(ast.Ins_F32Load)
		fmt.Fprintf(w, "%smemcpy(&$R[%d].f32, &wasm_memoy[$R[%d].i32+%d], 4);\n",
			indent, stk.TopIdx(), stk.TopIdx(), i.Offset,
		)
		assert(stk.Len() == stkLen)
	case token.INS_F64_LOAD:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		i := i.(ast.Ins_I32Load)
		fmt.Fprintf(w, "%smemcpy(&$R[%d].f64, &wasm_memoy[$R[%d].i32+%d], 8);\n",
			indent, stk.TopIdx(), stk.TopIdx(), i.Offset,
		)
		assert(stk.Len() == stkLen)
	case token.INS_I32_LOAD8_S:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		i := i.(ast.Ins_I32Load8S)
		fmt.Fprintf(w, "%smemcpy(&$R_u8, &wasm_memoy[$R[%d].i32+%d], 1); $R[%d].i32 = (i32_t)((i8_t)$R_u8);\n",
			indent, stk.TopIdx(), i.Offset, stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I32_LOAD8_U:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		i := i.(ast.Ins_I32Load8U)
		fmt.Fprintf(w, "%smemcpy(&$R_u8, &wasm_memoy[$R[%d].i32+%d], 1); $R[%d].i32 = (i32_t)((ui8_t)$R_u8);\n",
			indent, stk.TopIdx(), i.Offset, stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I32_LOAD16_S:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		i := i.(ast.Ins_I32Load16S)
		fmt.Fprintf(w, "%smemcpy(&$R_u16, &wasm_memoy[$R[%d].i32+%d], 2); $R[%d].i32 = (i32_t)((i16_t)$R_u16);\n",
			indent, stk.TopIdx(), i.Offset, stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I32_LOAD16_U:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		i := i.(ast.Ins_I32Load16U)
		fmt.Fprintf(w, "%smemcpy(&$R_u16, &wasm_memoy[$R[%d].i32+%d], 2); $R[%d].i32 = (i32_t)((ui16_t)$R_u16);\n",
			indent, stk.TopIdx(), i.Offset, stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LOAD8_S:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		i := i.(ast.Ins_I64Load8S)
		fmt.Fprintf(w, "%smemcpy(&$R_u8, &wasm_memoy[$R[%d].i32+%d], 1); $R[%d].i64 = (i64_t)((i8_t)$R_u8);\n",
			indent, stk.TopIdx(), i.Offset, stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LOAD8_U:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		i := i.(ast.Ins_I64Load8U)
		fmt.Fprintf(w, "%smemcpy(&$R_u8, &wasm_memoy[$R[%d].i32+%d], 1); $R[%d].i64 = (i64_t)((ui8_t)$R_u8);\n",
			indent, stk.TopIdx(), i.Offset, stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LOAD16_S:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		i := i.(ast.Ins_I64Load16S)
		fmt.Fprintf(w, "%smemcpy(&$R_u16, &wasm_memoy[$R[%d].i32+%d], 2); $R[%d].i64 = (i64_t)((i16_t)$R_u16);\n",
			indent, stk.TopIdx(), i.Offset, stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LOAD16_U:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		i := i.(ast.Ins_I64Load16U)
		fmt.Fprintf(w, "%smemcpy(&$R_u16, &wasm_memoy[$R[%d].i32+%d], 2); $R[%d].i64 = (i64_t)((ui16_t)$R_u16);\n",
			indent, stk.TopIdx(), i.Offset, stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LOAD32_S:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		i := i.(ast.Ins_I64Load32S)
		fmt.Fprintf(w, "%smemcpy(&$R_u32, &wasm_memoy[$R[%d].i32+%d], 4); $R[%d].i64 = (i64_t)((i32_t)$R_u32);\n",
			indent, stk.TopIdx(), i.Offset, stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_LOAD32_U:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		i := i.(ast.Ins_I64Load32U)
		fmt.Fprintf(w, "%smemcpy(&$R_u32, &wasm_memoy[$R[%d].i32+%d], 4); $R[%d].i64 = (i64_t)((ui32_t)$R_u32);\n",
			indent, stk.TopIdx(), i.Offset, stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I32_STORE:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		i := i.(ast.Ins_I32Store)
		fmt.Fprintf(w, "%smemcpy(&wasm_memoy[$R[%d].i32+%d], &$R[%d].i32, 4);\n",
			indent, stk.TopIdx()-1, i.Offset, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-2)
	case token.INS_I64_STORE:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		i := i.(ast.Ins_I64Store)
		fmt.Fprintf(w, "%smemcpy(&wasm_memoy[$R[%d].i32+%d], &$R[%d].i64, 8);\n",
			indent, stk.TopIdx()-1, i.Offset, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-2)
	case token.INS_F32_STORE:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		i := i.(ast.Ins_F32Store)
		fmt.Fprintf(w, "%smemcpy(&wasm_memoy[$R[%d].i32+%d], &$R[%d].f32, 4);\n",
			indent, stk.TopIdx()-1, i.Offset, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-2)
	case token.INS_F64_STORE:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		i := i.(ast.Ins_F64Store)
		fmt.Fprintf(w, "%smemcpy(&wasm_memoy[$R[%d].i32+%d], &$R[%d].f64, 8);\n",
			indent, stk.TopIdx()-1, i.Offset, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-2)
	case token.INS_I32_STORE8:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		i := i.(ast.Ins_I32Store8)
		fmt.Fprintf(w, "%s$R_u8 = (ui8_t)((i8_t)($R[%d].i32)); memcpy(&wasm_memoy[$R[%d].i32+%d], &$R_u8, 1);\n",
			indent, stk.TopIdx(), stk.TopIdx()-1, i.Offset,
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-2)
	case token.INS_I32_STORE16:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		i := i.(ast.Ins_I32Store16)
		fmt.Fprintf(w, "%s$R_u16 = (ui16_t)((i16_t)($R[%d].i32)); memcpy(&wasm_memoy[$R[%d].i32+%d], &$R_u16, 2);\n",
			indent, stk.TopIdx(), stk.TopIdx()-1, i.Offset,
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-2)
	case token.INS_I64_STORE8:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		i := i.(ast.Ins_I64Store8)
		fmt.Fprintf(w, "%s$R_u8 = (ui8_t)((i8_t)($R[%d].i64)); memcpy(&wasm_memoy[$R[%d].i32+%d], &$R_u8, 1);\n",
			indent, stk.TopIdx(), stk.TopIdx()-1, i.Offset,
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-2)
	case token.INS_I64_STORE16:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		i := i.(ast.Ins_I64Store16)
		fmt.Fprintf(w, "%s$R_u16 = (ui16_t)((i16_t)($R[%d].i64)); memcpy(&wasm_memoy[$R[%d].i32+%d], &$R_u16, 2);\n",
			indent, stk.TopIdx(), stk.TopIdx()-1, i.Offset,
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-2)
	case token.INS_I64_STORE32:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		i := i.(ast.Ins_I64Store32)
		fmt.Fprintf(w, "%s$R_u32 = (ui32_t)((i32_t)($R[%d].i64)); memcpy(&wasm_memoy[$R[%d].i32+%d], &$R_u32, 4);\n",
			indent, stk.TopIdx(), stk.TopIdx()-1, i.Offset,
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-2)
	case token.INS_MEMORY_SIZE:
		stkLen := stk.Len()
		stk.PushN(1)
		fmt.Fprintf(w, "%s$R[%d].i32 = wasm_memoy_size;\n", indent, stk.TopIdx())
		assert(stk.Len() == stkLen+1)
	case token.INS_MEMORY_GROW:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%swasm_memoy_size += $R[%d].i32; $R[%d].i32 = wasm_memoy_size;\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I32_CONST:
		stkLen := stk.Len()
		i := i.(ast.Ins_I32Const)
		stk.PushN(1)
		fmt.Fprintf(w, "%s$R[%d].i32 = %d;\n", indent, stk.TopIdx(), i.X)
		assert(stk.Len() == stkLen+1)
	case token.INS_I64_CONST:
		stkLen := stk.Len()
		i := i.(ast.Ins_I64Const)
		stk.PushN(1)
		fmt.Fprintf(w, "%s$R[%d].i64 = %d;\n", indent, stk.TopIdx(), i.X)
		assert(stk.Len() == stkLen+1)
	case token.INS_F32_CONST:
		stkLen := stk.Len()
		i := i.(ast.Ins_F32Const)
		stk.PushN(1)
		fmt.Fprintf(w, "%s$R[%d].f32 = %f;\n", indent, stk.TopIdx(), i.X)
		assert(stk.Len() == stkLen+1)
	case token.INS_F64_CONST:
		stkLen := stk.Len()
		i := i.(ast.Ins_F64Const)
		stk.PushN(1)
		fmt.Fprintf(w, "%s$R[%d].f64 = %f;\n", indent, stk.TopIdx(), i.X)
		assert(stk.Len() == stkLen+1)
	case token.INS_I32_EQZ:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].i32==0)? 1: 0;\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)

	case token.INS_I32_EQ:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].i32==$R[%d].i32)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_NE:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].i32!=$R[%d].i32)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_LT_S:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].i32<$R[%d].i32)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_LT_U:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ((ui32_t)($R[%d].i32)<(ui32_t)($R[%d].i32))? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_GT_S:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].i32>$R[%d].i32)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_GT_U:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ((ui32_t)($R[%d].i32)>(ui32_t)($R[%d].i32))? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_LE_S:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].i32<=$R[%d].i32)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_LE_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i32 = ((ui32_t)($R[%d].i32)<=(ui32_t)($R[%d].i32))? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_GE_S:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].i32>=$R[%d].i32)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_GE_U:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ((ui32_t)($R[%d].i32)>=(ui32_t)($R[%d].i32))? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_EQZ:
		stkLen := stk.Len()
		assert(stkLen >= 1)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].i64==0)? 1: 0;\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_EQ:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].i64==$R[%d].i64)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_NE:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].i64!=$R[%d].i64)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_LT_S:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].i64<$R[%d].i64)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_LT_U:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ((u64_t)($R[%d].i64)<(u64_t)($R[%d].i64))? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_GT_S:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].i64>$R[%d].i64)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_GT_U:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ((u64_t)($R[%d].i64)>(u64_t)($R[%d].i64))? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_LE_S:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].i64<=$R[%d].i64)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_LE_U:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ((u64_t)($R[%d].i64)<=(u64_t)($R[%d].i64))? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_GE_S:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].i64>=$R[%d].i64)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_GE_U:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ((u64_t)($R[%d].i64)>=(u64_t)($R[%d].i64))? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_F32_EQ:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].f32==$R[%d].f32)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_F32_NE:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].f32!=$R[%d].f32)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_F32_LT:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].f32<$R[%d].f32)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_F32_GT:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].f32>$R[%d].f32)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_F32_LE:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].f32<=$R[%d].f32)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_F32_GE:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].f32>=$R[%d].f32)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_F64_EQ:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].f64==$R[%d].f64)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_F64_NE:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].f64!=$R[%d].f64)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_F64_LT:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].f64<$R[%d].f64)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_F64_GT:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].f64>$R[%d].f64)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_F64_LE:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].f64<=$R[%d].f64)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_F64_GE:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = ($R[%d].f64>=$R[%d].f64)? 1: 0;\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.PopN(2)
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_CLZ:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i32 = I32_CLZ($R[%d].i32);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I32_CTZ:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i32 = I32_CTZ($R[%d].i32);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I32_POPCNT:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i32 = I32_POPCNT($R[%d].u32);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I32_ADD:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 += $R[%d].i32;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_SUB:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 -= $R[%d].i32;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_MUL:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 *= $R[%d].i32;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_DIV_S:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 /= $R[%d].i32;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_DIV_U:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = (i32_t)((ui32_t)($R[%d].i32)/(ui32_t)($R[%d].i32));\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_REM_S:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 %%= $R[%d].i32;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_REM_U:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = (i32_t)((ui32_t)($R[%d].i32)%%(ui32_t)($R[%d].i32));\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_AND:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 &= $R[%d].i32;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_OR:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 |= $R[%d].i32;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_XOR:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 ^= $R[%d].i32;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_SHL:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 <<= $R[%d].i32;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_SHR_S:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 >>= $R[%d].i32;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_SHR_U:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i32 = (i32_t)((ui32_t)($R[%d].i32)>>(ui32_t)($R[%d].i32));\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_ROTL:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i32 = I32_ROTL($R[%d].i32, $R[%d].i32);\n",
			indent, stk.TopIdx(), stk.TopIdx(), stk.TopIdx()-1,
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I32_ROTR:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i32 = I32_ROTR($R[%d].i32, $R[%d].i32);\n",
			indent, stk.TopIdx(), stk.TopIdx(), stk.TopIdx()-1,
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_CLZ:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i32 = I64_CLZ($R[%d].i64);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_CTZ:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i32 = I64_CTZ($R[%d].i64);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_POPCNT:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i32 = I64_POPCNT($R[%d].i64);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_ADD:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i64 += $R[%d].i64;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_SUB:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i64 -= $R[%d].i64;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_MUL:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i64 *= $R[%d].i64;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_DIV_S:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i64 /= $R[%d].i64;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_DIV_U:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i64 = (i64_t)((u64_t)($R[%d].i64)/(u64_t)($R[%d].i64));\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_REM_S:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i64 = (i64_t)((i64_t)($R[%d].i64)%%(i64_t)($R[%d].i64));\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_REM_U:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i64 = (i64_t)((u64_t)($R[%d].i64)/(u64_t)($R[%d].i64));\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_AND:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i64 &= $R[%d].i64;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_OR:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i64 |= $R[%d].i64;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_XOR:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i64 ^= $R[%d].i64;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_SHL:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i64 <<= $R[%d].i64;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_SHR_S:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i64 >>= $R[%d].i64;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_SHR_U:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].i64 = (i64_t)((u64_t)($R[%d].i64)>>(u64_t)($R[%d].i64));\n",
			indent, stk.TopIdx()-1, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_ROTL:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i64 = I64_ROTL($R[%d].i64, $R[%d].i32);\n",
			indent, stk.TopIdx(), stk.TopIdx(), stk.TopIdx()-1,
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_I64_ROTR:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i64 = I64_ROTR($R[%d].i64, $R[%d].i32);\n",
			indent, stk.TopIdx(), stk.TopIdx(), stk.TopIdx()-1,
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_F32_ABS:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f32 = fabsf($R[%d].f32);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F32_NEG:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f32 = 0-$R[%d].f32;\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F32_CEIL:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f32 = ceilf($R[%d].f32);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F32_FLOOR:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f32 = floorf($R[%d].f32);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F32_TRUNC:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f32 = truncf($R[%d].f32);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F32_NEAREST:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f32 = roundf($R[%d].f32);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F32_SQRT:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f32 = sqrtf($R[%d].f32);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F32_ADD:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].f32 += $R[%d].f32;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_F32_SUB:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].f32 -= $R[%d].f32;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_F32_MUL:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].f32 *= $R[%d].f32;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_F32_DIV:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].f32 /= $R[%d].f32;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_F32_MIN:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].f32 = fminf($R[%d].f32, $R[%d].f32);\n",
			indent, stk.TopIdx()-1, stk.TopIdx(), stk.TopIdx()-1,
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_F32_MAX:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].f32 = fmaxf($R[%d].f32, $R[%d].f32);\n",
			indent, stk.TopIdx()-1, stk.TopIdx(), stk.TopIdx()-1,
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_F32_COPYSIGN:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f32 = copysignf($R[%d].f32, $R[%d].f32);\n",
			indent, stk.TopIdx()-1, stk.TopIdx(), stk.TopIdx()-1,
		)
		assert(stk.Len() == stkLen)
	case token.INS_F64_ABS:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f64 = fabs($R[%d].f64);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F64_NEG:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f64 = 0-$R[%d].f64;\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F64_CEIL:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f64 = ceil($R[%d].f64);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F64_FLOOR:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f64 = floor($R[%d].f64);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F64_TRUNC:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f64 = trunc($R[%d].f64);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F64_NEAREST:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f64 = round($R[%d].f64);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F64_SQRT:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f64 = sqrt($R[%d].f64);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F64_ADD:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].f64 += $R[%d].f64;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_F64_SUB:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].f64 -= $R[%d].f64;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_F64_MUL:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].f64 *= $R[%d].f64;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_F64_DIV:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].f64 /= $R[%d].f64;\n",
			indent, stk.TopIdx()-1, stk.TopIdx(),
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_F64_MIN:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].f64 = fmin($R[%d].f64, $R[%d].f64);\n",
			indent, stk.TopIdx()-1, stk.TopIdx(), stk.TopIdx()-1,
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_F64_MAX:
		stkLen := stk.Len()
		assert(stkLen >= 2)
		fmt.Fprintf(w, "%s$R[%d].f64 = fmax($R[%d].f64, $R[%d].f64);\n",
			indent, stk.TopIdx()-1, stk.TopIdx(), stk.TopIdx()-1,
		)
		stk.Pop()
		assert(stk.Len() == stkLen-1)
	case token.INS_F64_COPYSIGN:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f64 = copysign($R[%d].f64, $R[%d].f64);\n",
			indent, stk.TopIdx()-1, stk.TopIdx(), stk.TopIdx()-1,
		)
		assert(stk.Len() == stkLen)
	case token.INS_I32_WRAP_I64:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i32 = (i32_t)($R[%d].i63&0xFFFFFFFF);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I32_TRUNC_F32_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i32 = (i32_t)(truncf($R[%d].f32));\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I32_TRUNC_F32_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i32 = (u32_t)(truncf($R[%d].f32));\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I32_TRUNC_F64_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i32 = (i32_t)(trunc($R[%d].f64));\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I32_TRUNC_F64_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i32 = (i32_t)(trunc($R[%d].f64));\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_EXTEND_I32_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i64 = (i64_t)($R[%d].i32);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_EXTEND_I32_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i64 = (i64_t)((u32_t)($R[%d].i32));\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_TRUNC_F32_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i64 = (i64_t)(truncf($R[%d].f32));\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_TRUNC_F32_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i64 = (i64_t)(truncf($R[%d].f32));\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_TRUNC_F64_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i64 = (i64_t)(trunc($R[%d].f64));\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_TRUNC_F64_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].i64 = (i64_t)(trunc($R[%d].f64));\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F32_CONVERT_I32_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f32 = (f32_t)($R[%d].i32);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F32_CONVERT_I32_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f32 = (f32_t)($R[%d].u32);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F32_CONVERT_I64_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f32 = (f32_t)($R[%d].i64);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F32_CONVERT_I64_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f32 = (f32_t)($R[%d].u64);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F32_DEMOTE_F64:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f32 = (f32_t)($R[%d].f64);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F64_CONVERT_I32_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f64 = (f64_t)($R[%d].i32);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F64_CONVERT_I32_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f64 = (f64_t)($R[%d].u32);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F64_CONVERT_I64_S:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f64 = (f64_t)($R[%d].i64);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F64_CONVERT_I64_U:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f64 = (f64_t)($R[%d].u64);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F64_PROMOTE_F32:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s$R[%d].f64 = (f64_t)($R[%d].f32);\n",
			indent, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I32_REINTERPRET_F32:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// %T: $R[%d].i32 <-- $R[%d].i32;\n",
			indent, i, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_I64_REINTERPRET_F64:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// %T: $R[%d].i64 <-- $R[%d].f64;\n",
			indent, i, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F32_REINTERPRET_I32:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// %T: $R[%d].f32 <-- $R[%d].i32;\n",
			indent, i, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)
	case token.INS_F64_REINTERPRET_I64:
		stkLen := stk.Len()
		fmt.Fprintf(w, "%s// %T: $R[%d].f64 <-- $R[%d].i64;\n",
			indent, i, stk.TopIdx(), stk.TopIdx(),
		)
		assert(stk.Len() == stkLen)

	default:
		panic(fmt.Sprintf("unreachable: %T", i))
	}
	return nil
}
