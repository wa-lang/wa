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

func (p *wat2cWorker) buildFunc_body(w io.Writer, fn *ast.Func, cRetType string) error {
	p.Tracef("buildFunc_body: %s\n", fn.Name)

	var stk valueTypeStack
	var bufIns bytes.Buffer

	stk.funcName = fn.Name

	assert(len(p.scopeLabels) == 0)
	assert(len(p.scopeStackBases) == 0)

	if len(fn.Body.Locals) > 0 {
		for _, x := range fn.Body.Locals {
			localName := toCName(x.Name)
			p.localNames = append(p.localNames, localName)
			p.localTypes = append(p.localTypes, x.Type)
			switch x.Type {
			case token.I32:
				fmt.Fprintf(&bufIns, "  i32_t %s = 0;", localName)
			case token.I64:
				fmt.Fprintf(&bufIns, "  i64_t %s = 0;", localName)
			case token.F32:
				fmt.Fprintf(&bufIns, "  f32_t %s = 0;", localName)
			case token.F64:
				fmt.Fprintf(&bufIns, "  f64_t %s = 0;", localName)
			}
			if localName != x.Name {
				fmt.Fprintf(&bufIns, " // %s\n", x.Name)
			} else {
				fmt.Fprintf(&bufIns, "\n")
			}
		}
		fmt.Fprintln(&bufIns)
	}

	// 生成函数调用顺序
	fmt.Fprintf(&bufIns, "  WASM_TRACE();\n")
	fmt.Fprintln(&bufIns)

	assert(stk.Len() == 0)
	for _, ins := range fn.Body.Insts {
		if err := p.buildFunc_ins(&bufIns, fn, &stk, ins, 1); err != nil {
			return err
		}
	}

	// 返回值
	if len(fn.Type.Results) > 1 {
		fmt.Fprintf(w, "  %s $result;\n", cRetType)
	}

	// 固定类型的寄存器
	fmt.Fprintf(w, "  u32_t $R_u32;\n")
	fmt.Fprintf(w, "  u16_t $R_u16;\n")
	fmt.Fprintf(w, "  u8_t  $R_u8;\n")

	// 栈寄存器(union类型)
	fmt.Fprintf(w, "  val_t $R0")
	for i := 1; i < stk.MaxDepth(); i++ {
		fmt.Fprintf(w, ", $R%d", i)
	}
	fmt.Fprintf(w, ";\n")
	fmt.Fprintln(w)

	// 指令复制到 w
	io.Copy(w, &bufIns)

	// 有些函数最后的位置不是 return, 需要手动清理栈
	switch stk.LastInstruction().Token() {
	case token.INS_RETURN:
		// 已经处理过了
	case token.INS_UNREACHABLE:
		// 清空残余的栈, 不做类型校验
		for stk.Len() > 0 {
			stk.DropAny()
		}
	default:
		// 补充 return
		assert(stk.Len() == len(fn.Type.Results))

		const indent = "  "
		switch len(fn.Type.Results) {
		case 0:
			fmt.Fprintf(w, "%sreturn;\n", indent)
		case 1:
			sp0 := stk.Pop(fn.Type.Results[0])
			switch fn.Type.Results[0] {
			case token.I32:
				fmt.Fprintf(w, "%sreturn $R%d.i32;\n", indent, sp0)
			case token.I64:
				fmt.Fprintf(w, "%sreturn $R%d.i64;\n", indent, sp0)
			case token.F32:
				fmt.Fprintf(w, "%sreturn $R%d.f32;\n", indent, sp0)
			case token.F64:
				fmt.Fprintf(w, "%sreturn $R%d.f64;\n", indent, sp0)
			default:
				unreachable()
			}
		default:
			for i, xType := range fn.Type.Results {
				spi := stk.Pop(xType)
				switch xType {
				case token.I32:
					fmt.Fprintf(w, "%s$result.$R%d = $R%d.i32;\n", indent, i, spi)
				case token.I64:
					fmt.Fprintf(w, "%s$result.$R%d = $R%d.i64;\n", indent, i, spi)
				case token.F32:
					fmt.Fprintf(w, "%s$result.$R%d = $R%d.f32;\n", indent, i, spi)
				case token.F64:
					fmt.Fprintf(w, "%s$result.$R%d = $R%d.f64;\n", indent, i, spi)
				default:
					unreachable()
				}
			}
			fmt.Fprintf(w, "%sreturn $result;\n", indent)
		}
	}
	assert(stk.Len() == 0)

	// 补充生成 return
	if stk.LastInstruction().Token() != token.INS_RETURN {
		switch len(fn.Type.Results) {
		case 0:
			fmt.Fprintf(w, "  return;\n")
		case 1:
			fmt.Fprintf(w, "  return 0;\n")
		default:
			fmt.Fprintf(w, "  return $result;\n")
		}
	}

	return nil
}

func (p *wat2cWorker) buildFunc_ins(w io.Writer, fn *ast.Func, stk *valueTypeStack, i ast.Instruction, level int) error {
	stk.NextInstruction(i)

	indent := strings.Repeat("  ", level)

	p.Tracef("buildFunc_ins: %s%s begin: %v\n", indent, i.Token(), stk.String())
	defer func() { p.Tracef("buildFunc_ins: %s%s end: %v\n", indent, i.Token(), stk.String()) }()

	if p.ifUseMathX(i.Token()) {
		p.useMathX = true
	}

	switch tok := i.Token(); tok {
	case token.INS_UNREACHABLE:
		fmt.Fprintf(w, "%sabort(); // %s\n", indent, tok)
	case token.INS_NOP:
		fmt.Fprintf(w, "%s((void)0); // %s\n", indent, tok)

	case token.INS_BLOCK:
		i := i.(ast.Ins_Block)

		stkBase := stk.Len()
		defer func() { assert(stk.Len() == stkBase+len(i.Results)) }()

		p.enterLabelScope(stkBase, i.Label)
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
			fmt.Fprintf(w, "L_%s_next:;\n", toCName(i.Label))
		}

	case token.INS_LOOP:
		i := i.(ast.Ins_Loop)

		stkBase := stk.Len()
		defer func() { assert(stk.Len() == stkBase+len(i.Results)) }()

		p.enterLabelScope(stkBase, i.Label)
		defer p.leaveLabelScope()

		if i.Label != "" {
			fmt.Fprintf(w, "L_%s_next:;\n", toCName(i.Label))
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

	case token.INS_IF:
		i := i.(ast.Ins_If)

		sp0 := stk.Pop(token.I32)
		fmt.Fprintf(w, "%sif($R%d.i32) {\n", indent, sp0)

		stkBase := stk.Len()
		defer func() { assert(stk.Len() == stkBase+len(i.Results)) }()

		p.enterLabelScope(stkBase, i.Label)
		defer p.leaveLabelScope()

		for _, ins := range i.Body {
			if err := p.buildFunc_ins(w, fn, stk, ins, level+1); err != nil {
				return err
			}
		}

		if len(i.Else) > 0 {
			p.Tracef("buildFunc_ins: %s%s begin: %v\n", indent, token.INS_ELSE, stk.String())
			defer func() { p.Tracef("buildFunc_ins: %s%s end: %v\n", indent, token.INS_ELSE, stk.String()) }()

			// 这是静态分析, 需要消除 if 分支对栈分配的影响
			for _, retType := range i.Results {
				stk.Pop(retType)
			}

			// 重新开始
			assert(stk.Len() == stkBase)

			fmt.Fprintf(w, "%s} else {\n", indent)
			for _, ins := range i.Else {
				if err := p.buildFunc_ins(w, fn, stk, ins, level+1); err != nil {
					return err
				}
			}
		}
		fmt.Fprintf(w, "%s}\n", indent)
		if i.Label != "" {
			fmt.Fprintf(w, "L_%s_next:;\n", toCName(i.Label))
		}

	case token.INS_ELSE:
		unreachable()
	case token.INS_END:
		unreachable()

	case token.INS_BR:
		i := i.(ast.Ins_Br)

		labelIdx := p.findLabelIndex(i.X)
		labelName := p.findLabelName(i.X)

		stkBase := p.scopeStackBases[len(p.scopeLabels)-labelIdx-1]
		assert(stk.Len() == stkBase)

		fmt.Fprintf(w, "%sgoto L_%s_next;\n", indent, toCName(labelName))

	case token.INS_BR_IF:
		i := i.(ast.Ins_BrIf)
		labelIdx := p.findLabelIndex(i.X)
		labelName := p.findLabelName(i.X)

		stkBase := p.scopeStackBases[len(p.scopeLabels)-labelIdx-1]

		sp0 := stk.Pop(token.I32)
		assert(stk.Len() == stkBase)

		fmt.Fprintf(w, "%sif($R%d.i32) { goto L_%s_next; }\n",
			indent, sp0, toCName(labelName),
		)
	case token.INS_BR_TABLE:
		i := i.(ast.Ins_BrTable)
		assert(len(i.XList) > 1)
		sp0 := stk.Pop(token.I32)
		fmt.Fprintf(w, "%sswitch($R%d.i32) {\n", indent, sp0)
		for k := 0; k < len(i.XList)-1; k++ {
			fmt.Fprintf(w, "%scase %d: goto L_%s_next;\n",
				indent, k, toCName(p.findLabelName(i.XList[k])),
			)
		}
		fmt.Fprintf(w, "%sdefault: goto L_%s_next;\n",
			indent, toCName(p.findLabelName(i.XList[len(i.XList)-1])),
		)
		fmt.Fprintf(w, "%s}\n", indent)

	case token.INS_RETURN:
		switch len(fn.Type.Results) {
		case 0:
			fmt.Fprintf(w, "%sreturn;\n", indent)
		case 1:
			sp0 := stk.Pop(fn.Type.Results[0])
			switch fn.Type.Results[0] {
			case token.I32:
				fmt.Fprintf(w, "%sreturn $R%d.i32;\n", indent, sp0)
			case token.I64:
				fmt.Fprintf(w, "%sreturn $R%d.i64;\n", indent, sp0)
			case token.F32:
				fmt.Fprintf(w, "%sreturn $R%d.f32;\n", indent, sp0)
			case token.F64:
				fmt.Fprintf(w, "%sreturn $R%d.f64;\n", indent, sp0)
			default:
				unreachable()
			}
		default:
			for i, xType := range fn.Type.Results {
				spi := stk.Pop(xType)
				switch xType {
				case token.I32:
					fmt.Fprintf(w, "%s$result.$R%d = $R%d.i32;\n", indent, i, spi)
				case token.I64:
					fmt.Fprintf(w, "%s$result.$R%d = $R%d.i64;\n", indent, i, spi)
				case token.F32:
					fmt.Fprintf(w, "%s$result.$R%d = $R%d.f32;\n", indent, i, spi)
				case token.F64:
					fmt.Fprintf(w, "%s$result.$R%d = $R%d.f64;\n", indent, i, spi)
				default:
					unreachable()
				}
			}
			fmt.Fprintf(w, "%sreturn $result;\n", indent)
		}
		assert(stk.Len() == 0)

	case token.INS_CALL:
		i := i.(ast.Ins_Call)

		fnCallType := p.findFuncType(i.X)
		fnCallCRetType := p.getFuncCRetType(fnCallType, i.X)

		// 参数列表
		argList := make([]int, len(fnCallType.Params))
		for k, x := range fnCallType.Params {
			argList[k] = stk.Pop(x.Type)
		}

		// 需要定义临时变量保存返回值
		switch len(fnCallType.Results) {
		case 0:
			fmt.Fprintf(w, "%sfn_%s(", indent, toCName(i.X))
		case 1:
			ret0 := stk.Push(fnCallType.Results[0])
			switch fnCallType.Results[0] {
			case token.I32:
				fmt.Fprintf(w, "%s$R%d.i32 = fn_%s(", indent, ret0, toCName(i.X))
			case token.I64:
				fmt.Fprintf(w, "%s$R%d.i64 = fn_%s(", indent, ret0, toCName(i.X))
			case token.F32:
				fmt.Fprintf(w, "%s$R%d.f32 = fn_%s(", indent, ret0, toCName(i.X))
			case token.F64:
				fmt.Fprintf(w, "%s$R%d.f64 = fn_%s(", indent, ret0, toCName(i.X))
			default:
				unreachable()
			}
		}

		if len(fnCallType.Results) > 1 {
			fmt.Fprintf(w, "%s{\n", indent)
			fmt.Fprintf(w, "%s  %s $ret = fn_%s(", indent, fnCallCRetType, toCName(i.X))
		}

		for k, x := range fnCallType.Params {
			if k > 0 {
				fmt.Fprintf(w, ", ")
			}
			switch x.Type {
			case token.I32:
				fmt.Fprintf(w, "$R%d.i32", argList[k])
			case token.I64:
				fmt.Fprintf(w, "$R%d.i64", argList[k])
			case token.F32:
				fmt.Fprintf(w, "$R%d.f32", argList[k])
			case token.F64:
				fmt.Fprintf(w, "$R%d.f64", argList[k])
			default:
				unreachable()
			}
		}
		fmt.Fprintf(w, ");\n")

		// 复制到当前stk
		if len(fnCallType.Results) > 1 {
			for k, retType := range fnCallType.Results {
				reti := stk.Push(retType)
				switch retType {
				case token.I32:
					fmt.Fprintf(w, "%s  $R%d.i32 = $ret.$R%d;\n", indent, reti, k)
				case token.I64:
					fmt.Fprintf(w, "%s  $R%d.i64 = $ret.$R%d;\n", indent, reti, k)
				case token.F32:
					fmt.Fprintf(w, "%s  $R%d.f32 = $ret.$R%d;\n", indent, reti, k)
				case token.F64:
					fmt.Fprintf(w, "%s  $R%d.f64 = $ret.$R%d;\n", indent, reti, k)
				default:
					unreachable()
				}
			}
			fmt.Fprintf(w, "%s}\n", indent)
		}

	case token.INS_CALL_INDIRECT:
		i := i.(ast.Ins_CallIndirect)

		sp0 := stk.Pop(token.I32)
		fnCallType := p.findType(i.TypeIdx)
		fnCallCRetType := p.getFuncCRetType(fnCallType, "")

		// 参数列表
		argList := make([]int, len(fnCallType.Params))
		for k, x := range fnCallType.Params {
			argList[k] = stk.Pop(x.Type)
		}

		// 生成要调用函数的类型
		fmt.Fprintf(w, "%s{\n", indent)
		{
			if len(fnCallType.Results) > 1 {
				fmt.Fprintf(w, "%s  typedef struct {", indent)
				for i := 0; i < len(fnCallType.Results); i++ {
					if i > 0 {
						fmt.Fprintf(w, " ")
					}
					switch fnCallType.Results[i] {
					case token.I32:
						fmt.Fprintf(w, "i32_t $R%d;", i)
					case token.I64:
						fmt.Fprintf(w, "i64_t $R%d;", i)
					case token.F32:
						fmt.Fprintf(w, "f32_t $R%d;", i)
					case token.F64:
						fmt.Fprintf(w, "f64_t $R%d;", i)
					default:
						unreachable()
					}
				}
				fmt.Fprintf(w, "} %s;\n", fnCallCRetType)
			}

			fmt.Fprintf(w, "%s  typedef %s (*fn_t)(", indent, fnCallCRetType)
			if len(fnCallType.Params) > 0 {
				for i, x := range fnCallType.Params {
					var argName string
					if x.Name != "" {
						argName = toCName(x.Name)
					} else {
						argName = fmt.Sprintf("$arg%d", i)
					}
					if i > 0 {
						fmt.Fprintf(w, ", ")
					}
					switch x.Type {
					case token.I32:
						fmt.Fprintf(w, "i32_t %s", argName)
					case token.I64:
						fmt.Fprintf(w, "i64_t %s", argName)
					case token.F32:
						fmt.Fprintf(w, "f32_t %s", argName)
					case token.F64:
						fmt.Fprintf(w, "f64_t %s", argName)
					default:
						unreachable()
					}
				}
			}
			fmt.Fprintf(w, ");\n")

			// 调用函数
			switch len(fnCallType.Results) {
			case 0:
				fmt.Fprintf(w, "%s  ((fn_t)(WASM_TABLE_GET($R%d.i32)))(",
					indent, sp0,
				)
			case 1:
				ret0 := stk.Push(fnCallType.Results[0])
				switch fnCallType.Results[0] {
				case token.I32:
					fmt.Fprintf(w, "%s  $R%d.i32 = ((fn_t)(WASM_TABLE_GET($R%d.i32)))(",
						indent, ret0, sp0,
					)
				case token.I64:
					fmt.Fprintf(w, "%s  $R%d.i64 = ((fn_t)(WASM_TABLE_GET($R%d.i32)))(",
						indent, ret0, sp0,
					)
				case token.F32:
					fmt.Fprintf(w, "%s  $R%d.f32 = ((fn_t)(WASM_TABLE_GET($R%d.i32)))(",
						indent, ret0, sp0,
					)
				case token.F64:
					fmt.Fprintf(w, "%s  $R%d.f64 = ((fn_t)(WASM_TABLE_GET($R%d.i32)))(",
						indent, ret0, sp0,
					)
				default:
					unreachable()
				}

			default:
				fmt.Fprintf(w, "%s  %s $ret = ((fn_t)(WASM_TABLE_GET($R%d.i32)))(",
					indent, fnCallCRetType, sp0,
				)
			}
			for i, x := range fnCallType.Params {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				argi := argList[i]
				switch x.Type {
				case token.I32:
					fmt.Fprintf(w, "$R%d.i32", argi)
				case token.I64:
					fmt.Fprintf(w, "$R%d.i64", argi)
				case token.F32:
					fmt.Fprintf(w, "$R%d.f32", argi)
				case token.F64:
					fmt.Fprintf(w, "$R%d.f64", argi)
				default:
					unreachable()
				}
			}
			fmt.Fprintf(w, ");\n")

			// 保存返回值
			if len(fnCallType.Results) > 1 {
				for k, retType := range fnCallType.Results {
					reti := stk.Push(retType)
					switch retType {
					case token.I32:
						fmt.Fprintf(w, "%s  $R%d.i32 = $ret.$R%d;\n", indent, reti, k)
					case token.I64:
						fmt.Fprintf(w, "%s  $R%d.i64 = $ret.$R%d;\n", indent, reti, k)
					case token.F32:
						fmt.Fprintf(w, "%s  $R%d.f32 = $ret.$R%d;\n", indent, reti, k)
					case token.F64:
						fmt.Fprintf(w, "%s  $R%d.f64 = $ret.$R%d;\n", indent, reti, k)
					default:
						unreachable()
					}
				}
			}
		}
		fmt.Fprintf(w, "%s}\n", indent)

	case token.INS_DROP:
		sp0 := stk.DropAny()
		fmt.Fprintf(w, "%s// drop $R%d\n", indent, sp0)
	case token.INS_SELECT:
		i := i.(ast.Ins_Select)

		sp0 := stk.Pop(token.I32) // 判断条件

		// wasm 2.0 支持带类型
		valType := i.ResultTyp
		if valType == 0 {
			// 不带类型, 2个数值类型必须一样
			valType = stk.TopToken()
		}
		sp1 := stk.Pop(valType)
		sp2 := stk.Pop(valType)

		ret0 := stk.Push(token.I32)

		// 注意返回值的顺序
		// if sp0 != 0 { sp2 } else { sp1 }
		fmt.Fprintf(w, "%s$R%d.i32 = $R%d.i32? $R%d.i32: $R%d.i32;\n",
			indent, ret0, sp0, sp2, sp1,
		)
	case token.INS_LOCAL_GET:
		i := i.(ast.Ins_LocalGet)
		xType := p.findLocalType(fn, i.X)
		ret0 := stk.Push(xType)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "%s$R%d.i32 = %s;\n", indent, ret0, p.findLocalName(fn, i.X))
		case token.I64:
			fmt.Fprintf(w, "%s$R%d.i64 = %s;\n", indent, ret0, p.findLocalName(fn, i.X))
		case token.F32:
			fmt.Fprintf(w, "%s$R%d.f32 = %s;\n", indent, ret0, p.findLocalName(fn, i.X))
		case token.F64:
			fmt.Fprintf(w, "%s$R%d.f64 = %s;\n", indent, ret0, p.findLocalName(fn, i.X))
		default:
			unreachable()
		}
	case token.INS_LOCAL_SET:
		i := i.(ast.Ins_LocalSet)
		xType := p.findLocalType(fn, i.X)
		sp0 := stk.Pop(xType)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "%s%s = $R%d.i32;\n", indent, p.findLocalName(fn, i.X), sp0)
		case token.I64:
			fmt.Fprintf(w, "%s%s = $R%d.i64;\n", indent, p.findLocalName(fn, i.X), sp0)
		case token.F32:
			fmt.Fprintf(w, "%s%s = $R%d.f32;\n", indent, p.findLocalName(fn, i.X), sp0)
		case token.F64:
			fmt.Fprintf(w, "%s%s = $R%d.f64;\n", indent, p.findLocalName(fn, i.X), sp0)
		default:
			unreachable()
		}
	case token.INS_LOCAL_TEE:
		i := i.(ast.Ins_LocalTee)
		xType := p.findLocalType(fn, i.X)
		sp0 := stk.Top(xType)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "%s%s = $R%d.i32;\n", indent, p.findLocalName(fn, i.X), sp0)
		case token.I64:
			fmt.Fprintf(w, "%s%s = $R%d.i64;\n", indent, p.findLocalName(fn, i.X), sp0)
		case token.F32:
			fmt.Fprintf(w, "%s%s = $R%d.f32;\n", indent, p.findLocalName(fn, i.X), sp0)
		case token.F64:
			fmt.Fprintf(w, "%s%s = $R%d.f64;\n", indent, p.findLocalName(fn, i.X), sp0)
		default:
			unreachable()
		}
	case token.INS_GLOBAL_GET:
		i := i.(ast.Ins_GlobalGet)
		xType := p.findGlobalType(i.X)
		ret0 := stk.Push(xType)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "%s$R%d.i32 = var_%s;\n", indent, ret0, toCName(i.X))
		case token.I64:
			fmt.Fprintf(w, "%s$R%d.i64 = var_%s;\n", indent, ret0, toCName(i.X))
		case token.F32:
			fmt.Fprintf(w, "%s$R%d.f32 = var_%s;\n", indent, ret0, toCName(i.X))
		case token.F64:
			fmt.Fprintf(w, "%s$R%d.f64 = var_%s;\n", indent, ret0, toCName(i.X))
		default:
			unreachable()
		}
	case token.INS_GLOBAL_SET:
		i := i.(ast.Ins_GlobalSet)
		xType := p.findGlobalType(i.X)
		sp0 := stk.Pop(xType)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "%svar_%s = $R%d.i32;\n", indent, toCName(i.X), sp0)
		case token.I64:
			fmt.Fprintf(w, "%svar_%s = $R%d.i64;\n", indent, toCName(i.X), sp0)
		case token.F32:
			fmt.Fprintf(w, "%svar_%s = $R%d.f32;\n", indent, toCName(i.X), sp0)
		case token.F64:
			fmt.Fprintf(w, "%svar_%s = $R%d.f64;\n", indent, toCName(i.X), sp0)
		default:
			unreachable()
		}
	case token.INS_TABLE_GET:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.FUNCREF) // funcref
		fmt.Fprintf(w, "%s$R%d.ref = WASM_TABLE_GET($R%d.i32);\n", indent, sp0, ret0)
	case token.INS_TABLE_SET:
		sp0 := stk.Pop(token.FUNCREF) // funcref
		sp1 := stk.Pop(token.I32)
		fmt.Fprintf(w, "%sWASM_TABLE_GET($R%d.i32,$R%d.ref);\n", indent, sp1, sp0)
	case token.INS_I32_LOAD:
		i := i.(ast.Ins_I32Load)
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%smemcpy(&$R%d.i32, WASM_MEMORY_ADDR($R%d.i32+%d), 4);\n",
			indent, sp0, ret0, i.Offset,
		)
	case token.INS_I64_LOAD:
		i := i.(ast.Ins_I64Load)
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%smemcpy(&$R%d.i64, WASM_MEMORY_ADDR($R%d.i32+%d), 8);\n",
			indent, sp0, ret0, i.Offset,
		)
	case token.INS_F32_LOAD:
		i := i.(ast.Ins_F32Load)
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%smemcpy(&$R%d.f32, WASM_MEMORY_ADDR($R%d.i32+%d), 4);\n",
			indent, sp0, ret0, i.Offset,
		)
	case token.INS_F64_LOAD:
		i := i.(ast.Ins_I32Load)
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%smemcpy(&$R%d.f64, WASM_MEMORY_ADDR($R%d.i32+%d), 8);\n",
			indent, sp0, ret0, i.Offset,
		)
	case token.INS_I32_LOAD8_S:
		i := i.(ast.Ins_I32Load8S)
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%smemcpy(&$R_u8, WASM_MEMORY_ADDR($R%d.i32+%d), 1); $R%d.i32 = (i32_t)((i8_t)$R_u8);\n",
			indent, ret0, i.Offset, sp0,
		)
	case token.INS_I32_LOAD8_U:
		i := i.(ast.Ins_I32Load8U)
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%smemcpy(&$R_u8, WASM_MEMORY_ADDR($R%d.i32+%d), 1); $R%d.i32 = (i32_t)((u8_t)$R_u8);\n",
			indent, ret0, i.Offset, sp0,
		)
	case token.INS_I32_LOAD16_S:
		i := i.(ast.Ins_I32Load16S)
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%smemcpy(&$R_u16, WASM_MEMORY_ADDR($R%d.i32+%d), 2); $R%d.i32 = (i32_t)((i16_t)$R_u16);\n",
			indent, ret0, i.Offset, sp0,
		)
	case token.INS_I32_LOAD16_U:
		i := i.(ast.Ins_I32Load16U)
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%smemcpy(&$R_u16, WASM_MEMORY_ADDR($R%d.i32+%d), 2); $R%d.i32 = (i32_t)((u16_t)$R_u16);\n",
			indent, ret0, i.Offset, sp0,
		)
	case token.INS_I64_LOAD8_S:
		i := i.(ast.Ins_I64Load8S)
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%smemcpy(&$R_u8, WASM_MEMORY_ADDR($R%d.i32+%d), 1); $R%d.i64 = (i64_t)((i8_t)$R_u8);\n",
			indent, ret0, i.Offset, sp0,
		)
	case token.INS_I64_LOAD8_U:
		i := i.(ast.Ins_I64Load8U)
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%smemcpy(&$R_u8, WASM_MEMORY_ADDR($R%d.i32+%d), 1); $R%d.i64 = (i64_t)((u8_t)$R_u8);\n",
			indent, ret0, i.Offset, sp0,
		)
	case token.INS_I64_LOAD16_S:
		i := i.(ast.Ins_I64Load16S)
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%smemcpy(&$R_u16, WASM_MEMORY_ADDR($R%d.i32+%d), 2); $R%d.i64 = (i64_t)((i16_t)$R_u16);\n",
			indent, ret0, i.Offset, sp0,
		)
	case token.INS_I64_LOAD16_U:
		i := i.(ast.Ins_I64Load16U)
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%smemcpy(&$R_u16, WASM_MEMORY_ADDR($R%d.i32+%d), 2); $R%d.i64 = (i64_t)((u16_t)$R_u16);\n",
			indent, ret0, i.Offset, sp0,
		)
	case token.INS_I64_LOAD32_S:
		i := i.(ast.Ins_I64Load32S)
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%smemcpy(&$R_u32, WASM_MEMORY_ADDR($R%d.i32+%d), 4); $R%d.i64 = (i64_t)((i32_t)$R_u32);\n",
			indent, ret0, i.Offset, sp0,
		)
	case token.INS_I64_LOAD32_U:
		i := i.(ast.Ins_I64Load32U)
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%smemcpy(&$R_u32, WASM_MEMORY_ADDR($R%d.i32+%d), 4); $R%d.i64 = (i64_t)((u32_t)$R_u32);\n",
			indent, ret0, i.Offset, sp0,
		)
	case token.INS_I32_STORE:
		i := i.(ast.Ins_I32Store)
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		fmt.Fprintf(w, "%smemcpy(WASM_MEMORY_ADDR($R%d.i32+%d), &$R%d.i32, 4);\n",
			indent, sp1, i.Offset, sp0,
		)
	case token.INS_I64_STORE:
		i := i.(ast.Ins_I64Store)
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I32)
		fmt.Fprintf(w, "%smemcpy(WASM_MEMORY_ADDR($R%d.i32+%d), &$R%d.i64, 8);\n",
			indent, sp1, i.Offset, sp0,
		)
	case token.INS_F32_STORE:
		i := i.(ast.Ins_F32Store)
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.I32)
		fmt.Fprintf(w, "%smemcpy(WASM_MEMORY_ADDR($R%d.i32+%d), &$R%d.f32, 4);\n",
			indent, sp1, i.Offset, sp0,
		)
	case token.INS_F64_STORE:
		i := i.(ast.Ins_F64Store)
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.I32)
		fmt.Fprintf(w, "%smemcpy(WASM_MEMORY_ADDR($R%d.i32+%d), &$R%d.f64, 8);\n",
			indent, sp1, i.Offset, sp0,
		)
	case token.INS_I32_STORE8:
		i := i.(ast.Ins_I32Store8)
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		fmt.Fprintf(w, "%s$R_u8 = (u8_t)((i8_t)($R%d.i32)); memcpy(WASM_MEMORY_ADDR($R%d.i32+%d), &$R_u8, 1);\n",
			indent, sp0, sp1, i.Offset,
		)
	case token.INS_I32_STORE16:
		i := i.(ast.Ins_I32Store16)
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		fmt.Fprintf(w, "%s$R_u16 = (u16_t)((i16_t)($R%d.i32)); memcpy(WASM_MEMORY_ADDR($R%d.i32+%d), &$R_u16, 2);\n",
			indent, sp0, sp1, i.Offset,
		)
	case token.INS_I64_STORE8:
		i := i.(ast.Ins_I64Store8)
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I32)
		fmt.Fprintf(w, "%s$R_u8 = (u8_t)((i8_t)($R%d.i64)); memcpy(WASM_MEMORY_ADDR($R%d.i32+%d), &$R_u8, 1);\n",
			indent, sp0, sp1, i.Offset,
		)
	case token.INS_I64_STORE16:
		i := i.(ast.Ins_I64Store16)
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I32)
		fmt.Fprintf(w, "%s$R_u16 = (u16_t)((i16_t)($R%d.i64)); memcpy(WASM_MEMORY_ADDR($R%d.i32+%d), &$R_u16, 2);\n",
			indent, sp0, sp1, i.Offset,
		)
	case token.INS_I64_STORE32:
		i := i.(ast.Ins_I64Store32)
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I32)
		fmt.Fprintf(w, "%s$R_u32 = (u32_t)((i32_t)($R%d.i64)); memcpy(WASM_MEMORY_ADDR($R%d.i32+%d), &$R_u32, 4);\n",
			indent, sp0, sp1, i.Offset,
		)
	case token.INS_MEMORY_SIZE:
		sp0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = wasm_memory_size;\n", indent, sp0)
	case token.INS_MEMORY_GROW:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sif(wasm_memory_size+$R%d.i32 <= wasm_memory_init_max_pages) {\n",
			indent, sp0,
		)
		{
			fmt.Fprintf(w, "%si32_t $temp = wasm_memory_size;\n",
				indent+indent,
			)
			fmt.Fprintf(w, "%swasm_memory_size += $R%d.i32;\n",
				indent+indent, sp0,
			)
			fmt.Fprintf(w, "%s$R%d.i32 = $temp;\n",
				indent+indent, ret0,
			)
		}
		fmt.Fprintf(w, "%s} else {\n",
			indent,
		)
		{
			fmt.Fprintf(w, "%s$R%d.i32 = -1;\n",
				indent+indent, ret0,
			)
		}
		fmt.Fprintf(w, "%s}\n",
			indent,
		)

	case token.INS_MEMORY_INIT:
		i := i.(ast.Ins_MemoryInit)

		len := stk.Pop(token.I32)
		off := stk.Pop(token.I32)
		dst := stk.Pop(token.I32)

		var sb strings.Builder
		datai := p.m.Data[i.DataIdx].Value[off:][:len]
		for _, x := range datai {
			sb.WriteString(fmt.Sprintf("\\x%02x", x))
		}

		fmt.Fprintf(w, "%smemcpy((void*)WASM_MEMORY_ADDR($R%d.i32), (void*)(\"%s\"), %d);\n",
			indent, dst, sb.String(), len,
		)
	case token.INS_MEMORY_COPY:
		len := stk.Pop(token.I32)
		src := stk.Pop(token.I32)
		dst := stk.Pop(token.I32)
		fmt.Fprintf(w, "%smemcpy((void*)WASM_MEMORY_ADDR($R%d.i32), (void*)WASM_MEMORY_ADDR($R%d.i32), $R%d.i32);\n",
			indent, dst, src, len,
		)
	case token.INS_MEMORY_FILL:
		len := stk.Pop(token.I32)
		val := stk.Pop(token.I32)
		dst := stk.Pop(token.I32)
		fmt.Fprintf(w, "%smemset((void*)WASM_MEMORY_ADDR($R%d.i32), $R%d.i32, $R%d.i32);\n",
			indent, dst, val, len,
		)
	case token.INS_I32_CONST:
		i := i.(ast.Ins_I32Const)
		sp0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = %d;\n", indent, sp0, i.X)
	case token.INS_I64_CONST:
		i := i.(ast.Ins_I64Const)
		sp0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = %d;\n", indent, sp0, i.X)
	case token.INS_F32_CONST:
		i := i.(ast.Ins_F32Const)
		sp0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = %f;\n", indent, sp0, i.X)
	case token.INS_F64_CONST:
		i := i.(ast.Ins_F64Const)
		sp0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = %f;\n", indent, sp0, i.X)
	case token.INS_I32_EQZ:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.i32==0)? 1: 0;\n",
			indent, ret0, sp0,
		)
	case token.INS_I32_EQ:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.i32==$R%d.i32)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_NE:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.i32!=$R%d.i32)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_LT_S:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.i32<$R%d.i32)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_LT_U:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ((u32_t)($R%d.i32)<(u32_t)($R%d.i32))? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_GT_S:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.i32>$R%d.i32)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_GT_U:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ((u32_t)($R%d.i32)>(u32_t)($R%d.i32))? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_LE_S:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.i32<=$R%d.i32)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_LE_U:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ((u32_t)($R%d.i32)<=(u32_t)($R%d.i32))? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_GE_S:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.i32>=$R%d.i32)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_GE_U:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ((u32_t)($R%d.i32)>=(u32_t)($R%d.i32))? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_EQZ:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.i64==0)? 1: 0;\n",
			indent, ret0, sp0,
		)
	case token.INS_I64_EQ:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.i64==$R%d.i64)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_NE:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.i64!=$R%d.i64)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_LT_S:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.i64<$R%d.i64)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_LT_U:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ((u64_t)($R%d.i64)<(u64_t)($R%d.i64))? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_GT_S:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.i64>$R%d.i64)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_GT_U:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ((u64_t)($R%d.i64)>(u64_t)($R%d.i64))? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_LE_S:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.i64<=$R%d.i64)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_LE_U:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ((u64_t)($R%d.i64)<=(u64_t)($R%d.i64))? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_GE_S:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.i64>=$R%d.i64)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_GE_U:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ((u64_t)($R%d.i64)>=(u64_t)($R%d.i64))? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F32_EQ:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.f32==$R%d.f32)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F32_NE:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.f32!=$R%d.f32)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F32_LT:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.f32<$R%d.f32)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F32_GT:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.f32>$R%d.f32)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F32_LE:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.f32<=$R%d.f32)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F32_GE:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.f32>=$R%d.f32)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F64_EQ:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.f64==$R%d.f64)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F64_NE:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.f64!=$R%d.f64)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F64_LT:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.f64<$R%d.f64)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F64_GT:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.f64>$R%d.f64)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F64_LE:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.f64<=$R%d.f64)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F64_GE:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = ($R%d.f64>=$R%d.f64)? 1: 0;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_CLZ:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = I32_CLZ($R%d.i32);\n",
			indent, ret0, sp0,
		)
	case token.INS_I32_CTZ:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = I32_CTZ($R%d.i32);\n",
			indent, ret0, sp0,
		)
	case token.INS_I32_POPCNT:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = I32_POPCNT($R%d.i32);\n",
			indent, ret0, sp0,
		)
	case token.INS_I32_ADD:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = $R%d.i32 + $R%d.i32;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_SUB:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = $R%d.i32 - $R%d.i32;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_MUL:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = $R%d.i32 * $R%d.i32;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_DIV_S:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = $R%d.i32 / $R%d.i32;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_DIV_U:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = (i32_t)((u32_t)($R%d.i32)/(u32_t)($R%d.i32));\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_REM_S:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = $R%d.i32 %% $R%d.i32;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_REM_U:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = (i32_t)((u32_t)($R%d.i32)%%(u32_t)($R%d.i32));\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_AND:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = $R%d.i32 & $R%d.i32;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_OR:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = $R%d.i32 | $R%d.i32;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_XOR:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = $R%d.i32 ^ $R%d.i32;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_SHL:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = $R%d.i32 << ($R%d.i32&63);\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_SHR_S:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = $R%d.i32 >> ($R%d.i32&63);\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_SHR_U:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = (i32_t)((u32_t)($R%d.i32)>>(u32_t)($R%d.i32&63));\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_ROTL:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = I32_ROTL($R%d.i32, $R%d.i32);\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_ROTR:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = I32_ROTR($R%d.i32, $R%d.i32);\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_CLZ:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = I64_CLZ($R%d.i64);\n",
			indent, ret0, sp0,
		)
	case token.INS_I64_CTZ:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = I64_CTZ($R%d.i64);\n",
			indent, ret0, sp0,
		)
	case token.INS_I64_POPCNT:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = I64_POPCNT($R%d.i64);\n",
			indent, ret0, sp0,
		)
	case token.INS_I64_ADD:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = $R%d.i64 + $R%d.i64;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_SUB:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = $R%d.i64 - $R%d.i64;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_MUL:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = $R%d.i64 * $R%d.i64;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_DIV_S:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = $R%d.i64 / $R%d.i64;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_DIV_U:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = (i64_t)((u64_t)($R%d.i64)/(u64_t)($R%d.i64));\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_REM_S:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = (i64_t)((i64_t)($R%d.i64)%%(i64_t)($R%d.i64));\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_REM_U:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = (i64_t)((u64_t)($R%d.i64)%%(u64_t)($R%d.i64));\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_AND:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = $R%d.i64 & $R%d.i64;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_OR:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = $R%d.i64 | $R%d.i64;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_XOR:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = $R%d.i64 ^ $R%d.i64;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_SHL:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = $R%d.i64 << (((u64_t)$R%d.i64)&63);\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_SHR_S:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = $R%d.i64 >> (((u64_t)$R%d.i64)&63);\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_SHR_U:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = (i64_t)((u64_t)($R%d.i64)>>((u64_t)($R%d.i64)&63));\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_ROTL:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = I64_ROTL($R%d.i64, $R%d.i64);\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I64_ROTR:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = I64_ROTR($R%d.i64, $R%d.i64);\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F32_ABS:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = fabsf($R%d.f32);\n",
			indent, ret0, sp0,
		)
	case token.INS_F32_NEG:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = 0-$R%d.f32;\n",
			indent, ret0, sp0,
		)
	case token.INS_F32_CEIL:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = ceilf($R%d.f32);\n",
			indent, ret0, sp0,
		)
	case token.INS_F32_FLOOR:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = floorf($R%d.f32);\n",
			indent, ret0, sp0,
		)
	case token.INS_F32_TRUNC:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = truncf($R%d.f32);\n",
			indent, ret0, sp0,
		)
	case token.INS_F32_NEAREST:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = roundf($R%d.f32);\n",
			indent, ret0, sp0,
		)
	case token.INS_F32_SQRT:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = sqrtf($R%d.f32);\n",
			indent, ret0, sp0,
		)
	case token.INS_F32_ADD:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = $R%d.f32 + $R%d.f32;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F32_SUB:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = $R%d.f32 - $R%d.f32;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F32_MUL:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = $R%d.f32 * $R%d.f32;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F32_DIV:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = $R%d.f32 / $R%d.f32;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F32_MIN:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = fminf($R%d.f32, $R%d.f32);\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F32_MAX:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = fmaxf($R%d.f32, $R%d.f32);\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F32_COPYSIGN:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = copysignf($R%d.f32, $R%d.f32);\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F64_ABS:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = fabs($R%d.f64);\n",
			indent, ret0, sp0,
		)
	case token.INS_F64_NEG:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = 0-$R%d.f64;\n",
			indent, ret0, sp0,
		)
	case token.INS_F64_CEIL:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = ceil($R%d.f64);\n",
			indent, ret0, sp0,
		)
	case token.INS_F64_FLOOR:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = floor($R%d.f64);\n",
			indent, ret0, sp0,
		)
	case token.INS_F64_TRUNC:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = trunc($R%d.f64);\n",
			indent, ret0, sp0,
		)
	case token.INS_F64_NEAREST:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = round($R%d.f64);\n",
			indent, ret0, sp0,
		)
	case token.INS_F64_SQRT:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = sqrt($R%d.f64);\n",
			indent, ret0, sp0,
		)
	case token.INS_F64_ADD:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = $R%d.f64 + $R%d.f64;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F64_SUB:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = $R%d.f64 - $R%d.f64;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F64_MUL:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = $R%d.f64 * $R%d.f64;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F64_DIV:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = $R%d.f64 / $R%d.f64;\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F64_MIN:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = fmin($R%d.f64, $R%d.f64);\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F64_MAX:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = fmax($R%d.f64, $R%d.f64);\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_F64_COPYSIGN:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = copysign($R%d.f64, $R%d.f64);\n",
			indent, ret0, sp1, sp0,
		)
	case token.INS_I32_WRAP_I64:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = (i32_t)($R%d.i64);\n",
			indent, ret0, sp0,
		)
	case token.INS_I32_TRUNC_F32_S:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = (i32_t)(truncf($R%d.f32));\n",
			indent, ret0, sp0,
		)
	case token.INS_I32_TRUNC_F32_U:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = (i32_t)(u32_t)(truncf($R%d.f32));\n",
			indent, ret0, sp0,
		)
	case token.INS_I32_TRUNC_F64_S:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = (i32_t)(trunc($R%d.f64));\n",
			indent, ret0, sp0,
		)
	case token.INS_I32_TRUNC_F64_U:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%s$R%d.i32 = (i32_t)(u32_t)(trunc($R%d.f64));\n",
			indent, ret0, sp0,
		)
	case token.INS_I64_EXTEND_I32_S:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = (i64_t)($R%d.i32);\n",
			indent, ret0, sp0,
		)
	case token.INS_I64_EXTEND_I32_U:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = (i64_t)((u32_t)($R%d.i32));\n",
			indent, ret0, sp0,
		)
	case token.INS_I64_TRUNC_F32_S:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = (i64_t)(i32_t)(truncf($R%d.f32));\n",
			indent, ret0, sp0,
		)
	case token.INS_I64_TRUNC_F32_U:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = (i64_t)(u32_t)(truncf($R%d.f32));\n",
			indent, ret0, sp0,
		)
	case token.INS_I64_TRUNC_F64_S:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = (i64_t)(trunc($R%d.f64));\n",
			indent, ret0, sp0,
		)
	case token.INS_I64_TRUNC_F64_U:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%s$R%d.i64 = (i64_t)(u64_t)(trunc($R%d.f64));\n",
			indent, ret0, sp0,
		)
	case token.INS_F32_CONVERT_I32_S:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = (f32_t)($R%d.i32);\n",
			indent, ret0, sp0,
		)
	case token.INS_F32_CONVERT_I32_U:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = (f32_t)($R%d.u32);\n",
			indent, ret0, sp0,
		)
	case token.INS_F32_CONVERT_I64_S:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = (f32_t)($R%d.i64);\n",
			indent, ret0, sp0,
		)
	case token.INS_F32_CONVERT_I64_U:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = (f32_t)($R%d.u64);\n",
			indent, ret0, sp0,
		)
	case token.INS_F32_DEMOTE_F64:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d.f32 = (f32_t)($R%d.f64);\n",
			indent, ret0, sp0,
		)
	case token.INS_F64_CONVERT_I32_S:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = (f64_t)($R%d.i32);\n",
			indent, ret0, sp0,
		)
	case token.INS_F64_CONVERT_I32_U:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = (f64_t)($R%d.u32);\n",
			indent, ret0, sp0,
		)
	case token.INS_F64_CONVERT_I64_S:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = (f64_t)($R%d.i64);\n",
			indent, ret0, sp0,
		)
	case token.INS_F64_CONVERT_I64_U:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = (f64_t)($R%d.u64);\n",
			indent, ret0, sp0,
		)
	case token.INS_F64_PROMOTE_F32:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d.f64 = (f64_t)($R%d.f32);\n",
			indent, ret0, sp0,
		)
	case token.INS_I32_REINTERPRET_F32:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d = $R%d;\n",
			indent, ret0, sp0,
		)
	case token.INS_I64_REINTERPRET_F64:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d = $R%d;\n",
			indent, ret0, sp0,
		)
	case token.INS_F32_REINTERPRET_I32:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%s$R%d = $R%d;\n",
			indent, ret0, sp0,
		)
	case token.INS_F64_REINTERPRET_I64:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%s$R%d = $R%d;\n",
			indent, ret0, sp0,
		)
	default:
		panic(fmt.Sprintf("unreachable: %T", i))
	}
	return nil
}
