// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"wa-lang.org/wa/internal/native/abi"
	nativeast "wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/ast/astutil"
	"wa-lang.org/wa/internal/native/loong64"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

//
// 函数栈帧布局
// 参考 /docs/asm_abi_la64.md
//
// 补充说明:
// - GP 以 global 变量保存
// - 预留出足够的 WasmStack 空间, 避免和子函数返回的数据空间冲突
// - WasmStack 上的内存也通过 FP 定位
//

func (p *wat2laWorker) buildFuncs(w io.Writer) error {
	if len(p.m.Funcs) == 0 {
		return nil
	}

	for _, f := range p.m.Funcs {
		fmt.Fprintf(w, "func $%s", f.Name)
		if len(f.Type.Params) > 0 {
			fmt.Fprint(w, "(")
			for i, x := range f.Type.Params {
				if i > 0 {
					fmt.Fprint(w, ", ")
				}
				if x.Name != "" {
					fmt.Fprintf(w, "$%s: %v", x.Name, x.Type)
				} else {
					fmt.Fprintf(w, "$arg.%d, %v", i, x.Type)
				}
			}
			fmt.Fprint(w, ")")
		}
		if len(f.Type.Results) > 1 {
			fmt.Fprintf(w, " => (")
			for i, x := range f.Type.Results {
				if i > 0 {
					fmt.Fprint(w, ", ")
				}
				fmt.Fprintf(w, "$ret.%d: %v", i, x)
			}
			fmt.Fprint(w, ")")
		} else if len(f.Type.Results) == 1 {
			fmt.Fprintf(w, "%v", f.Type.Results[0])
		}
		fmt.Fprintln(w, " {")

		// 翻译函数实现
		{
			p.localNames = nil
			p.localTypes = nil
			p.scopeLabels = nil
			p.scopeStackBases = nil
			p.scopeResults = nil

			if err := p.buildFunc_body(w, f); err != nil {
				return err
			}
		}

		fmt.Fprintln(w, "}")
	}

	return nil
}

func (p *wat2laWorker) buildFunc_body(w io.Writer, fn *ast.Func) error {
	p.Tracef("buildFunc_body: %s\n", fn.Name)

	var bufHeader bytes.Buffer

	if p.m.Memory != nil {
		addrType := p.m.Memory.AddrType
		assert(addrType == token.I32 || addrType == token.I64)
	}

	assert(p.m.Memory.AddrType == token.I32)

	assert(len(p.scopeLabels) == 0)
	assert(len(p.scopeStackBases) == 0)
	assert(len(p.scopeResults) == 0)

	if len(fn.Body.Locals) > 0 {
		for _, x := range fn.Body.Locals {
			p.localNames = append(p.localNames, x.Name)
			p.localTypes = append(p.localTypes, x.Type)
			fmt.Fprintf(&bufHeader, "    local %s: %v\n", x.Name, x.Type)
		}
		fmt.Fprintln(&bufHeader)
	}

	// 转化为汇编的结构, 准备构建函数栈帧
	fnNative := &nativeast.Func{
		Name: fn.Name,
		Type: &nativeast.FuncType{
			Args:   make([]*nativeast.Local, len(fn.Type.Params)),
			Return: make([]*nativeast.Local, len(fn.Type.Results)),
		},
		Body: &nativeast.FuncBody{
			Locals: make([]*nativeast.Local, len(fn.Body.Locals)),
		},
	}
	for i, arg := range fn.Type.Params {
		fnNative.Type.Args[i] = &nativeast.Local{
			Name: arg.Name,
			Type: wat2nativeType(arg.Type),
			Cap:  1,
		}
	}
	for i, typ := range fn.Type.Results {
		fnNative.Type.Return[i] = &nativeast.Local{
			Name: fmt.Sprintf("$ret.%d", i),
			Type: wat2nativeType(typ),
			Cap:  1,
		}
	}
	for i, local := range fn.Body.Locals {
		fnNative.Body.Locals[i] = &nativeast.Local{
			Name: local.Name,
			Type: wat2nativeType(local.Type),
			Cap:  1,
		}
	}

	// 模拟构建函数栈帧
	// 后面需要根据参数是否走寄存器传递生成相关的入口代码和返回代码
	if err := astutil.BuildFuncFrame(abi.LOONG64, fnNative); err != nil {
		return err
	}

	// 第一步构建 ra 和 fp
	{
		fmt.Fprintf(&bufHeader, "     # 栈帧开始\n")
		fmt.Fprintln(&bufHeader, "    addi.d $sp, $sp, -16 # $sp = $sp - 16")
		fmt.Fprintln(&bufHeader, "    st.d   $ra, $sp, 8   # save $ra")
		fmt.Fprintln(&bufHeader, "    st.d   $fp, $sp, 0   # save $fp")

		fmt.Fprintf(&bufHeader, "    # memory address\n")
		fmt.Fprintf(&bufHeader, "    pcalau12i %s, %%pc_hi20(%s)\n", kMemoryReg, kMemoryName)
		fmt.Fprintf(&bufHeader, "    addi.d    %s, %s, %%pc_lo12(%s)\n", kMemoryReg, kMemoryReg, kMemoryName)

		fmt.Fprintf(&bufHeader, "    # table address\n")
		fmt.Fprintf(&bufHeader, "    pcalau12i %s, %%pc_hi20(%s)\n", kTableReg, kTableName)
		fmt.Fprintf(&bufHeader, "    addi.d    %s, %s, %%pc_lo12(%s)\n", kTableReg, kTableReg, kTableName)
	}

	// 将寄存器参数备份到栈
	if len(fn.Type.Params) > 0 {
		fmt.Fprintf(&bufHeader, "    # 将寄存器参数备份到栈\n")
		for i, arg := range fnNative.Type.Args {
			if arg.Reg == 0 {
				continue // 走栈的输入参数不需要
			}

			// 将寄存器中的参数存储到对于的栈帧中
			switch fn.Type.Params[i].Type {
			case token.I32:
				fmt.Fprintf(&bufHeader, "    st.w %v, $fp, %d # save %s\n",
					loong64.RegString(arg.Reg),
					arg.Off,
					arg.Name,
				)
			case token.I64:
				fmt.Fprintf(&bufHeader, "    st.d %v, $fp, %d # save %s\n",
					loong64.RegString(arg.Reg),
					arg.Off,
					arg.Name,
				)
			case token.F32:
				fmt.Fprintf(&bufHeader, "    fst.s %v, $fp, %d # save %s\n",
					loong64.RegString(arg.Reg),
					arg.Off,
					arg.Name,
				)
			case token.F64:
				fmt.Fprintf(&bufHeader, "    fst.d %v, $fp, %d # save %s\n",
					loong64.RegString(arg.Reg),
					arg.Off,
					arg.Name,
				)
			default:
				panic("unreachable")
			}
		}
	}

	// 返回值初始化为 0
	for i, ret := range fnNative.Type.Return {
		switch fn.Type.Results[i] {
		case token.I32:
			fmt.Fprintf(&bufHeader, "    st.w zero, $fp, %d # save %s\n",
				ret.Off,
				ret.Name,
			)
		case token.I64:
			fmt.Fprintf(&bufHeader, "    st.d zero, $fp, %d # save %s\n",
				ret.Off,
				ret.Name,
			)
		case token.F32:
			fmt.Fprintf(&bufHeader, "    fst.s zero, $fp, %d # save %s\n",
				ret.Off,
				ret.Name,
			)
		case token.F64:
			fmt.Fprintf(&bufHeader, "    fst.d zero, $fp, %d # save %s\n",
				ret.Off,
				ret.Name,
			)
		default:
			panic("unreachable")
		}
	}

	// 局部遍历初始化为 0
	for i, ret := range fnNative.Type.Return {
		switch fn.Body.Locals[i].Type {
		case token.I32:
			fmt.Fprintf(&bufHeader, "    st.w zero, $fp, %d # save %s\n",
				ret.Off,
				ret.Name,
			)
		case token.I64:
			fmt.Fprintf(&bufHeader, "    st.d zero, $fp, %d # save %s\n",
				ret.Off,
				ret.Name,
			)
		case token.F32:
			fmt.Fprintf(&bufHeader, "    fst.s zero, $fp, %d # save %s\n",
				ret.Off,
				ret.Name,
			)
		case token.F64:
			fmt.Fprintf(&bufHeader, "    fst.d zero, $fp, %d # save %s\n",
				ret.Off,
				ret.Name,
			)
		default:
			panic("unreachable")
		}
	}

	// WASM 栈的开始位置
	{
		// WASM 栈底第一个元素相对于 FP 的偏移位置, 每个元素 8 字节
		n := len(fnNative.Type.Args) + len(fnNative.Type.Return) + len(fnNative.Body.Locals)
		p.fnWasmR0Base = 0 - (n+2+1)*8
	}

	// 至少要有一个指令
	if len(fn.Body.Insts) == 0 {
		fn.Body.Insts = []ast.Instruction{
			ast.Ins_Return{OpToken: ast.OpToken(token.INS_RETURN)},
		}
	}

	// 开始解析 wasm 指令
	var stk valueTypeStack
	var bufIns bytes.Buffer

	stk.funcName = fn.Name

	assert(stk.Len() == 0)
	for _, ins := range fn.Body.Insts {
		if err := p.buildFunc_ins(&bufIns, fnNative, fn, &stk, ins); err != nil {
			return err
		}
	}

	// 头部代码: 更新栈底精确位置
	{
		// 输入参数/返回值/局部遍历
		n := len(fnNative.Type.Args) + len(fnNative.Type.Return) + len(fnNative.Body.Locals)

		// WASM 栈的最大深度
		n += stk.MaxDepth()

		// SP 必须对齐到 16 字节
		if n%16 != 0 {
			n = ((n + 15) / 16) * 16
		}

		fmt.Fprintf(&bufHeader, "    addi.d $sp, $sp, -%d  # $sp = $sp - len(args)*8\n", n*8)
	}

	// 根据ABI处理返回值
	{
		// 返回代码位置
		fmt.Fprintf(&bufIns, "$L.%s.return:\n", fn.Name)

		// 如果走内存, 返回地址
		if len(fn.Type.Results) > 1 && fnNative.Type.Return[1].Reg == 0 {
			fmt.Fprintf(&bufIns, "    addi.d a0, $fp, %d # ret.%s\n",
				fnNative.Type.Return[0].Off,
				fnNative.Type.Return[0].Name,
			)
		} else {
			for i, ret := range fnNative.Type.Return {
				if ret.Reg == 0 {
					continue
				}
				fmt.Fprintf(&bufIns, "    ld.d %v, $fp, %d # ret.%s\n",
					loong64.RegString(ret.Reg),
					fnNative.Type.Return[i].Off,
					fnNative.Type.Return[i].Name,
				)
			}
		}
	}

	// 结束栈帧
	{
		fmt.Fprintf(&bufIns, "     # 栈帧结束\n")
		fmt.Fprintln(&bufIns, "    addi.d   $sp, $fp, 0 # $sp = $fp")
		fmt.Fprintln(&bufIns, "    ld.d $ra, $fp, -8    # restore $ra")
		fmt.Fprintln(&bufIns, "    ld.d $fp, $fp, -16   # restore $fp")
		fmt.Fprintln(&bufIns, "    jr $ra               # return")
	}

	// 头部赋值到 w
	io.Copy(w, &bufHeader)
	// 指令复制到 w
	io.Copy(w, &bufIns)

	// 有些函数最后的位置不是 return, 需要手动清理栈
	switch tok := stk.LastInstruction().Token(); tok {
	case token.INS_RETURN:
		// 已经处理过了
	case token.INS_UNREACHABLE:
		// 清空残余的栈, 不做类型校验
		if tok == token.INS_UNREACHABLE {
			for stk.Len() > 0 {
				stk.DropAny()
			}
		}

	default:
		// 补充 return
		assert(stk.Len() == len(fn.Type.Results))

		const indent = "  "
		switch len(fn.Type.Results) {
		case 0:
		case 1:
			stk.Pop(fn.Type.Results[0])
		default:
			for _, xType := range fn.Type.Results {
				stk.Pop(xType)
			}
		}
	}
	assert(stk.Len() == 0)

	return nil
}

func (p *wat2laWorker) buildFunc_ins(
	w io.Writer, fnNative *nativeast.Func,
	fn *ast.Func, stk *valueTypeStack, i ast.Instruction,
) error {
	stk.NextInstruction(i)

	const indent = "    "
	p.Tracef("buildFunc_ins: %s%s begin: %v\n", indent, i.Token(), stk.String())
	defer func() { p.Tracef("buildFunc_ins: %s%s end: %v\n", indent, i.Token(), stk.String()) }()

	switch tok := i.Token(); tok {
	case token.INS_UNREACHABLE:
		fmt.Fprintln(w, indent, "bl $builtin.panic # unreachable")
	case token.INS_NOP:
		fmt.Fprintln(w, indent, "addi zero, zero, 0 # nop")

	case token.INS_BLOCK:
		i := i.(ast.Ins_Block)

		stkBase := stk.Len()
		defer func() { assert(stk.Len() == stkBase+len(i.Results)) }()

		p.enterLabelScope(stkBase, i.Label, i.Results)
		defer p.leaveLabelScope()

		if i.Label != "" {
			fmt.Fprintf(w, "L.block.%s.begin:\n", i.Label)
		}

		for _, ins := range i.List {
			if err := p.buildFunc_ins(w, fnNative, fn, stk, ins); err != nil {
				return err
			}
		}

		if i.Label != "" {
			fmt.Fprintf(w, "L.block.%s.end:\n\n", i.Label)
		}

	case token.INS_LOOP:
		i := i.(ast.Ins_Loop)

		stkBase := stk.Len()
		defer func() { assert(stk.Len() == stkBase+len(i.Results)) }()

		p.enterLabelScope(stkBase, i.Label, i.Results)
		defer p.leaveLabelScope()

		if i.Label != "" {
			fmt.Fprintf(w, "L.loop.%s.begin:\n", i.Label)
		}
		for _, ins := range i.List {
			if err := p.buildFunc_ins(w, fnNative, fn, stk, ins); err != nil {
				return err
			}
		}
		if i.Label != "" {
			fmt.Fprintf(w, "L.loop.%s.end:\n", i.Label)
		}

	case token.INS_IF:
		i := i.(ast.Ins_If)

		if i.Label != "" {
			fmt.Fprintf(w, "L.if.%s.begin:\n", i.Label)
		}

		// 龙芯没有 pop 指令，需要2个指令才能实现
		// 因此通过固定的偏移量，只需要一个指令

		sp0 := stk.Pop(token.I32)

		fmt.Fprintf(w, "    ld.w t0, fp, %d\n", p.fnWasmR0Base-sp0*8)
		fmt.Fprintf(w, "    bne  t0, zero, L.if.%s.body\n", i.Label)
		if len(i.Else) > 0 {
			fmt.Fprintf(w, "    bne  t0, zero, L.if.%s.else\n", i.Label)
		} else {
			fmt.Fprintf(w, "    bne  t0, zero, L.if.%s.end\n", i.Label)
		}

		stkBase := stk.Len()
		defer func() { assert(stk.Len() == stkBase+len(i.Results)) }()

		p.enterLabelScope(stkBase, i.Label, i.Results)
		defer p.leaveLabelScope()

		if i.Label != "" {
			fmt.Fprintf(w, "L.if.%s.body:\n", i.Label)
		}
		for _, ins := range i.Body {
			if err := p.buildFunc_ins(w, fnNative, fn, stk, ins); err != nil {
				return err
			}
		}

		if len(i.Else) > 0 {
			p.Tracef("buildFunc_ins: %s%s begin: %v\n", indent, token.INS_ELSE, stk.String())
			defer func() { p.Tracef("buildFunc_ins: %s%s end: %v\n", indent, token.INS_ELSE, stk.String()) }()

			if i.Label != "" {
				fmt.Fprintf(w, "L.if.%s.else:\n", i.Label)
			}

			// 这是静态分析, 需要消除 if 分支对栈分配的影响
			for _, retType := range i.Results {
				stk.Pop(retType)
			}

			// 重新开始
			assert(stk.Len() == stkBase)

			for _, ins := range i.Else {
				if err := p.buildFunc_ins(w, fnNative, fn, stk, ins); err != nil {
					return err
				}
			}
		}
		if i.Label != "" {
			fmt.Fprintf(w, "L.if.%s.end:\n", i.Label)
		}

	case token.INS_ELSE:
		unreachable()
	case token.INS_END:
		unreachable()

	case token.INS_BR:
		i := i.(ast.Ins_Br)

		// br会根据目标block的返回值个数, 从当前block产生的栈中去返回值,
		// 至于中间被跳过的block栈数据全部被丢弃.

		labelIdx := p.findLabelIndex(i.X)
		labelName := p.findLabelName(i.X)

		destScopeStackBase := p.scopeStackBases[len(p.scopeLabels)-labelIdx-1]
		destScopeResults := p.scopeResults[len(p.scopeLabels)-labelIdx-1]

		currentScopeStackBase := p.scopeStackBases[len(p.scopeLabels)-1]

		// br对应的block带返回值
		if len(destScopeResults) > 0 {
			// 必须确保当前block的stk上有足够的返回值
			assert(currentScopeStackBase+len(destScopeResults) >= stk.Len())

			// 第一个返回值返回值的偏移地址
			firstResultOffset := stk.Len() - len(destScopeResults)

			// 如果返回值位置和目标block的base不一致则需要逐个复制
			if firstResultOffset > destScopeStackBase {
				// 返回值是逆序出栈
				fmt.Fprintf(w, "%s# copy br %s result\n", indent, labelName)
				for i := len(destScopeResults) - 1; i >= 0; i-- {
					xType := destScopeResults[i]
					reti := stk.Pop(xType)
					switch xType {
					case token.I32:
						fmt.Fprintf(w, "    ld.w t0, fp, %d\n", p.fnWasmR0Base-reti*8)
						fmt.Fprintf(w, "    st.w t0, fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
					case token.I64:
						fmt.Fprintf(w, "    ld.d t0, fp, %d\n", p.fnWasmR0Base-reti*8)
						fmt.Fprintf(w, "    st.d t0, fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
					case token.F32:
						fmt.Fprintf(w, "    fld.s ft0, fp, %d\n", p.fnWasmR0Base-reti*8)
						fmt.Fprintf(w, "    fst.s ft0, fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
					case token.F64:
						fmt.Fprintf(w, "    fld.d ft0, fp, %d\n", p.fnWasmR0Base-reti*8)
						fmt.Fprintf(w, "    fst.d ft0, fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
					default:
						unreachable()
					}
				}
			}
		}

		// 清除当前block的栈中除了目标返回值剩余的值
		// 这个操作只是为了退出当前block, 因为br已经是最后一个指令
		// 外层的block需要一个清理之后的栈
		for stk.Len() > currentScopeStackBase {
			stk.DropAny()
		}

		// 退出当前block时, stack已经被清理
		// 中间栈帧的数据会在外层block指令时处理
		assert(stk.Len() == currentScopeStackBase)

		fmt.Fprintf(w, "    b L.%s.next\n", labelName)

	case token.INS_BR_IF:
		i := i.(ast.Ins_BrIf)
		labelIdx := p.findLabelIndex(i.X)
		labelName := p.findLabelName(i.X)

		scopeStackBase := p.scopeStackBases[len(p.scopeLabels)-labelIdx-1]
		scopeResults := p.scopeResults[len(p.scopeLabels)-labelIdx-1]

		// 而br-if因为涉及else分支(需要维持栈平衡), 当前block后续的指令假设br-if只消耗了一个i32用于条件,
		// 因此这种条件br的目标block不能带返回值.
		assert(len(scopeResults) == 0)

		// 如果是跨越多个Block, 只需要丢弃中间block的栈数据即可
		if scopeStackBase > stk.Len() {
			// 这里是生成运行时代码, 因此不需要涉及栈丢弃的逻辑
			// 跨越多个block对应的stk的变量位置会在依次处理外层block时对应上
		}

		sp0 := stk.Pop(token.I32)
		assert(stk.Len() == scopeStackBase)

		fmt.Fprintf(w, "    # br_if %s\n", labelName)
		fmt.Fprintf(w, "    ld.w t0, fp, %d\n", p.fnWasmR0Base-sp0*8)
		fmt.Fprintf(w, "    bne t0, zero L.%s.next\n", labelName)

	case token.INS_BR_TABLE:
		i := i.(ast.Ins_BrTable)
		assert(len(i.XList) > 1)

		// br-table的行为和br比较相似, 因此不涉及else部分不用担心栈平衡的问题.
		// 但是每个目标block的返回值必须完全一致

		// 默认分支
		defaultLabelIdx := p.findLabelIndex(i.XList[len(i.XList)-1])
		defaultLabelName := p.findLabelName(i.XList[len(i.XList)-1])
		defaultScopeResults := p.scopeResults[len(p.scopeLabels)-defaultLabelIdx-1]

		currentScopeStackBase := p.scopeStackBases[len(p.scopeLabels)-1]

		sp0 := stk.Pop(token.I32)

		fmt.Fprintf(w, "    # br_table\n")
		fmt.Fprintf(w, "    ld.w t0, fp, %d\n", p.fnWasmR0Base-sp0*8)
		{
			// 当前block的返回值位置是相同的, 只能统一取一次
			var retIdxList = make([]int, len(defaultScopeResults))
			for k := len(defaultScopeResults) - 1; k >= 0; k-- {
				xTyp := defaultScopeResults[k]
				retIdxList[k] = stk.Pop(xTyp)
			}

			for k := 0; k < len(i.XList); k++ {
				labelIdx := p.findLabelIndex(i.XList[k])
				labelName := p.findLabelName(i.XList[k])

				destScopeStackBase := p.scopeStackBases[len(p.scopeLabels)-labelIdx-1]
				destScopeResults := p.scopeResults[len(p.scopeLabels)-labelIdx-1]

				// 验证每个目标返回值必须和default一致
				assert(len(defaultScopeResults) == len(destScopeResults))
				for i := 0; i < len(defaultScopeResults); i++ {
					assert(defaultScopeResults[i] == destScopeResults[i])
				}

				// 带返回值的情况
				if len(destScopeResults) > 0 {
					// 必须确保当前block的stk上有足够的返回值
					assert(currentScopeStackBase+len(destScopeResults) >= stk.Len())

					// 第一个返回值返回值的偏移地址
					firstResultOffset := stk.Len() - len(destScopeResults)

					// 如果返回值位置和目标block的base不一致则需要逐个复制
					if firstResultOffset > destScopeStackBase {
						// 返回值是逆序出栈
						fmt.Fprintf(w, "%s# copy br_table %s result\n", indent, labelName)
						for i := 0; i < len(destScopeResults); i++ {
							xType := destScopeResults[i]
							reti := retIdxList[i]
							switch xType {
							case token.I32:
								fmt.Fprintf(w, "    ld.w t1, fp, %d;\n", p.fnWasmR0Base-reti*8)
								fmt.Fprintf(w, "    st.w t1, fp, %d;\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
							case token.I64:
								fmt.Fprintf(w, "    ld.d t1, fp, %d;\n", p.fnWasmR0Base-reti*8)
								fmt.Fprintf(w, "    st.d t1, fp, %d;\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
							case token.F32:
								fmt.Fprintf(w, "    fld.w t1, fp, %d;\n", p.fnWasmR0Base-reti*8)
								fmt.Fprintf(w, "    fst.w t1, fp, %d;\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
							case token.F64:
								fmt.Fprintf(w, "    fld.d t1, fp, %d;\n", p.fnWasmR0Base-reti*8)
								fmt.Fprintf(w, "    fst.d t1, fp, %d;\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
							default:
								unreachable()
							}
						}
					}
				}

				if k < len(i.XList)-1 {
					fmt.Fprintf(w, "    # br_table case %d\n", k)
					fmt.Fprintf(w, "    ld.w t1, zero, %d;\n", k)
					fmt.Fprintf(w, "    bne  t0, t1 L.%s.next\n", labelName)
				} else {
					assert(labelName == defaultLabelName)
					fmt.Fprintf(w, "    # br_table default\n")
					fmt.Fprintf(w, "    b L.%s.next;\n", defaultLabelName)
				}
			}
		}

		// 清除当前block的栈中除了目标返回值剩余的值
		// 这个操作只是为了退出当前block, 因为br已经是最后一个指令
		// 外层的block需要一个清理之后的栈
		for stk.Len() > currentScopeStackBase {
			stk.DropAny()
		}

		// 退出当前block时, stack已经被清理
		// 中间栈帧的数据会在外层block指令时处理
		assert(stk.Len() == currentScopeStackBase)

		fmt.Fprintf(w, "    # br_table end\n")

	case token.INS_RETURN:
		switch len(fn.Type.Results) {
		case 0:
			fmt.Fprintf(w, "    b L.%s.return\n", fn.Name)
		case 1:
			sp0 := stk.Pop(fn.Type.Results[0])
			switch fn.Type.Results[0] {
			case token.I32:
				fmt.Fprintf(w, "    ld.w t0, fp, %d\n", p.fnWasmR0Base-sp0*8)
				fmt.Fprintf(w, "    st.w t0, fp, %d\n", fnNative.Type.Return[0].Off)
				fmt.Fprintf(w, "    b    L.%s.return\n", fn.Name)
			case token.I64:
				fmt.Fprintf(w, "    ld.d t0, fp, %d\n", p.fnWasmR0Base-sp0*8)
				fmt.Fprintf(w, "    st.d t0, fp, %d\n", fnNative.Type.Return[0].Off)
				fmt.Fprintf(w, "    b    L.%s.return\n", fn.Name)
			case token.F32:
				fmt.Fprintf(w, "    fld.s t0, fp, %d\n", p.fnWasmR0Base-sp0*8)
				fmt.Fprintf(w, "    fst.s t0, fp, %d\n", fnNative.Type.Return[0].Off)
				fmt.Fprintf(w, "    b    L.%s.return\n", fn.Name)
			case token.F64:
				fmt.Fprintf(w, "    fld.d t0, fp, %d\n", p.fnWasmR0Base-sp0*8)
				fmt.Fprintf(w, "    fst.d t0, fp, %d\n", fnNative.Type.Return[0].Off)
				fmt.Fprintf(w, "    b    L.%s.return\n", fn.Name)
			default:
				unreachable()
			}
		default:
			for i := len(fn.Type.Results) - 1; i >= 0; i-- {
				xType := fn.Type.Results[i]
				spi := stk.Pop(xType)
				switch xType {
				case token.I32:
					fmt.Fprintf(w, "    ld.w t0, fp, %d\n", p.fnWasmR0Base-spi*8)
					fmt.Fprintf(w, "    st.w t0, fp, %d\n", fnNative.Type.Return[i].Off)
				case token.I64:
					fmt.Fprintf(w, "    ld.d t0, fp, %d\n", p.fnWasmR0Base-spi*8)
					fmt.Fprintf(w, "    st.d t0, fp, %d\n", fnNative.Type.Return[i].Off)
				case token.F32:
					fmt.Fprintf(w, "    fld.s t0, fp, %d\n", p.fnWasmR0Base-spi*8)
					fmt.Fprintf(w, "    fst.s t0, fp, %d\n", fnNative.Type.Return[i].Off)
				case token.F64:
					fmt.Fprintf(w, "    fld.d t0, fp, %d\n", p.fnWasmR0Base-spi*8)
					fmt.Fprintf(w, "    fst.d t0, fp, %d\n", fnNative.Type.Return[i].Off)
				default:
					unreachable()
				}
			}
			fmt.Fprintf(w, "    b    L.%s.return\n", fn.Name)
		}
		assert(stk.Len() == 0)

	case token.INS_CALL:
		i := i.(ast.Ins_Call)

		fnCallType := p.findFuncType(i.X)

		// 构建被调用函数的栈帧信息
		fnCallNative := &nativeast.Func{
			Name: i.X,
			Type: &nativeast.FuncType{
				Args:   make([]*nativeast.Local, len(fnCallType.Params)),
				Return: make([]*nativeast.Local, len(fnCallType.Results)),
			},
			Body: &nativeast.FuncBody{
				// 不需要局部变量信息
			},
		}
		for i, arg := range fnCallType.Params {
			fnCallNative.Type.Args[i] = &nativeast.Local{
				Name: arg.Name,
				Type: wat2nativeType(arg.Type),
				Cap:  1,
			}
		}
		for i, typ := range fnCallType.Results {
			fnCallNative.Type.Return[i] = &nativeast.Local{
				Name: fmt.Sprintf("$ret.%d", i),
				Type: wat2nativeType(typ),
				Cap:  1,
			}
		}

		// 模拟构建函数栈帧
		// 后面需要根据参数是否走寄存器传递生成相关的入口代码和返回代码
		if err := astutil.BuildFuncFrame(abi.LOONG64, fnCallNative); err != nil {
			return err
		}

		// 参数列表
		// 出栈的顺序相反
		argList := make([]int, len(fnCallType.Params))
		for k := len(argList) - 1; k >= 0; k-- {
			x := fnCallType.Params[k]
			argList[k] = stk.Pop(x.Type)
		}

		// 准备调用参数
		for k, x := range fnCallType.Params {
			switch x.Type {
			case token.I32:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    ld.w %s, fp, %d",
						loong64.RegString(arg.Reg),
						p.fnWasmR0Base+argList[k]*8,
					)
				} else {
					fmt.Fprintf(w, "    ld.w t1, fp, %d",
						p.fnWasmR0Base+argList[k]*8,
					)
					fmt.Fprintf(w, "    st.w t1, fp, %d",
						arg.Off,
					)
				}
			case token.I64:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    ld.d %s, fp, %d",
						loong64.RegString(arg.Reg),
						p.fnWasmR0Base+argList[k]*8,
					)
				} else {
					fmt.Fprintf(w, "    ld.d t1, fp, %d",
						p.fnWasmR0Base+argList[k]*8,
					)
					fmt.Fprintf(w, "    st.d t1, fp, %d",
						arg.Off,
					)
				}
			case token.F32:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    fld.s %s, fp, %d",
						loong64.RegString(arg.Reg),
						p.fnWasmR0Base+argList[k]*8,
					)
				} else {
					fmt.Fprintf(w, "    fld.s t1, fp, %d",
						p.fnWasmR0Base+argList[k]*8,
					)
					fmt.Fprintf(w, "    fst.s t1, fp, %d",
						arg.Off,
					)
				}
			case token.F64:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    fld.d %s, fp, %d",
						loong64.RegString(arg.Reg),
						p.fnWasmR0Base+argList[k]*8,
					)
				} else {
					fmt.Fprintf(w, "    fld.d t1, fp, %d",
						p.fnWasmR0Base+argList[k]*8,
					)
					fmt.Fprintf(w, "    fst.d t1, fp, %d",
						arg.Off,
					)
				}
			default:
				unreachable()
			}
		}
		fmt.Fprintf(w, "    bl %s", i.X)

		// 提取返回值
		if len(fnCallNative.Type.Return) > 1 && fnCallNative.Type.Return[1].Reg == 0 {
			// 走栈返回结果, a0 是地址
			// 根据 ABI 规范获取返回值, 再保存到 WASM 栈中
			for k, retType := range fnCallType.Results {
				reti := stk.Push(retType)
				switch retType {
				case token.I32:
					fmt.Fprintf(w, "    ld.w t0, a0, %d\n", k*8)
					fmt.Fprintf(w, "    st.w t0, fp, %d\n", p.fnWasmR0Base-reti*8)
				case token.I64:
					fmt.Fprintf(w, "    ld.d t0, a0, %d\n", k*8)
					fmt.Fprintf(w, "    st.d t0, fp, %d\n", p.fnWasmR0Base-reti*8)
				case token.F32:
					fmt.Fprintf(w, "    fld.s t0, a0, %d\n", k*8)
					fmt.Fprintf(w, "    fst.s t0, fp, %d\n", p.fnWasmR0Base-reti*8)
				case token.F64:
					fmt.Fprintf(w, "    fld.d t0, a0, %d\n", k*8)
					fmt.Fprintf(w, "    fst.d t0, fp, %d\n", p.fnWasmR0Base-reti*8)
				default:
					unreachable()
				}
			}
		} else {
			// 走寄存器返回
			for k, retType := range fnCallType.Results {
				retNative := fnCallNative.Type.Return[k]
				assert(retNative.Reg != 0)

				reti := stk.Push(retType)
				switch retType {
				case token.I32:
					fmt.Fprintf(w, "    st.w %s, fp, %d\n",
						loong64.RegString(retNative.Reg),
						p.fnWasmR0Base-reti*8,
					)
				case token.I64:
					fmt.Fprintf(w, "    st.d %s, fp, %d\n",
						loong64.RegString(retNative.Reg),
						p.fnWasmR0Base-reti*8,
					)
				case token.F32:
					fmt.Fprintf(w, "    fst.s %s, fp, %d\n",
						loong64.RegString(retNative.Reg),
						p.fnWasmR0Base-reti*8,
					)
				case token.F64:
					fmt.Fprintf(w, "    fst.d %s, fp, %d\n",
						loong64.RegString(retNative.Reg),
						p.fnWasmR0Base-reti*8,
					)
				default:
					unreachable()
				}
			}
		}

	case token.INS_CALL_INDIRECT:
		i := i.(ast.Ins_CallIndirect)

		fnCallType := p.findType(i.TypeIdx)

		// 构建被调用函数的栈帧信息
		fnCallNative := &nativeast.Func{
			Name: "", // 间接调用, 没有名字(TODO: 可以尝试根据地址查询名字)
			Type: &nativeast.FuncType{
				Args:   make([]*nativeast.Local, len(fnCallType.Params)),
				Return: make([]*nativeast.Local, len(fnCallType.Results)),
			},
			Body: &nativeast.FuncBody{
				// 不需要局部变量信息
			},
		}
		for i, arg := range fnCallType.Params {
			fnCallNative.Type.Args[i] = &nativeast.Local{
				Name: arg.Name,
				Type: wat2nativeType(arg.Type),
				Cap:  1,
			}
		}
		for i, typ := range fnCallType.Results {
			fnCallNative.Type.Return[i] = &nativeast.Local{
				Name: fmt.Sprintf("$ret.%d", i),
				Type: wat2nativeType(typ),
				Cap:  1,
			}
		}

		// 模拟构建函数栈帧
		// 后面需要根据参数是否走寄存器传递生成相关的入口代码和返回代码
		if err := astutil.BuildFuncFrame(abi.LOONG64, fnCallNative); err != nil {
			return err
		}

		// 获取函数地址
		sp0 := stk.Pop(token.I32)
		fnRegName := "t2"
		fmt.Fprintf(w, "    ld.d %s, fp, %d\n", fnRegName, p.fnWasmR0Base-sp0*8)

		// 参数列表
		// 出栈的顺序相反
		argList := make([]int, len(fnCallType.Params))
		for k := len(argList) - 1; k >= 0; k-- {
			x := fnCallType.Params[k]
			argList[k] = stk.Pop(x.Type)
		}

		// 准备调用参数
		for k, x := range fnCallType.Params {
			switch x.Type {
			case token.I32:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    ld.w %s, fp, %d",
						loong64.RegString(arg.Reg),
						p.fnWasmR0Base+argList[k]*8,
					)
				} else {
					fmt.Fprintf(w, "    ld.w t1, fp, %d",
						p.fnWasmR0Base+argList[k]*8,
					)
					fmt.Fprintf(w, "    st.w t1, fp, %d",
						arg.Off,
					)
				}
			case token.I64:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    ld.d %s, fp, %d",
						loong64.RegString(arg.Reg),
						p.fnWasmR0Base+argList[k]*8,
					)
				} else {
					fmt.Fprintf(w, "    ld.d t1, fp, %d",
						p.fnWasmR0Base+argList[k]*8,
					)
					fmt.Fprintf(w, "    st.d t1, fp, %d",
						arg.Off,
					)
				}
			case token.F32:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    fld.s %s, fp, %d",
						loong64.RegString(arg.Reg),
						p.fnWasmR0Base+argList[k]*8,
					)
				} else {
					fmt.Fprintf(w, "    fld.s t1, fp, %d",
						p.fnWasmR0Base+argList[k]*8,
					)
					fmt.Fprintf(w, "    fst.s t1, fp, %d",
						arg.Off,
					)
				}
			case token.F64:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    fld.d %s, fp, %d",
						loong64.RegString(arg.Reg),
						p.fnWasmR0Base+argList[k]*8,
					)
				} else {
					fmt.Fprintf(w, "    fld.d t1, fp, %d",
						p.fnWasmR0Base+argList[k]*8,
					)
					fmt.Fprintf(w, "    fst.d t1, fp, %d",
						arg.Off,
					)
				}
			default:
				unreachable()
			}
		}
		fmt.Fprintf(w, "    bl %s\n", fnRegName)

		// 提取返回值
		if len(fnCallNative.Type.Return) > 1 && fnCallNative.Type.Return[1].Reg == 0 {
			// 走栈返回结果, a0 是地址
			// 根据 ABI 规范获取返回值, 再保存到 WASM 栈中
			for k, retType := range fnCallType.Results {
				reti := stk.Push(retType)
				switch retType {
				case token.I32:
					fmt.Fprintf(w, "    ld.w t0, a0, %d\n", k*8)
					fmt.Fprintf(w, "    st.w t0, fp, %d\n", p.fnWasmR0Base-reti*8)
				case token.I64:
					fmt.Fprintf(w, "    ld.d t0, a0, %d\n", k*8)
					fmt.Fprintf(w, "    st.d t0, fp, %d\n", p.fnWasmR0Base-reti*8)
				case token.F32:
					fmt.Fprintf(w, "    fld.s t0, a0, %d\n", k*8)
					fmt.Fprintf(w, "    fst.s t0, fp, %d\n", p.fnWasmR0Base-reti*8)
				case token.F64:
					fmt.Fprintf(w, "    fld.d t0, a0, %d\n", k*8)
					fmt.Fprintf(w, "    fst.d t0, fp, %d\n", p.fnWasmR0Base-reti*8)
				default:
					unreachable()
				}
			}
		} else {
			// 走寄存器返回
			for k, retType := range fnCallType.Results {
				retNative := fnCallNative.Type.Return[k]
				assert(retNative.Reg != 0)

				reti := stk.Push(retType)
				switch retType {
				case token.I32:
					fmt.Fprintf(w, "    st.w %s, fp, %d\n",
						loong64.RegString(retNative.Reg),
						p.fnWasmR0Base-reti*8,
					)
				case token.I64:
					fmt.Fprintf(w, "    st.d %s, fp, %d\n",
						loong64.RegString(retNative.Reg),
						p.fnWasmR0Base-reti*8,
					)
				case token.F32:
					fmt.Fprintf(w, "    fst.s %s, fp, %d\n",
						loong64.RegString(retNative.Reg),
						p.fnWasmR0Base-reti*8,
					)
				case token.F64:
					fmt.Fprintf(w, "    fst.d %s, fp, %d\n",
						loong64.RegString(retNative.Reg),
						p.fnWasmR0Base-reti*8,
					)
				default:
					unreachable()
				}
			}
		}

	case token.INS_DROP:
		sp0 := stk.DropAny()
		fmt.Fprintf(w, "    addi zero, zero, 0 # drop R%d\n", sp0)
	case token.INS_SELECT:
		i := i.(ast.Ins_Select)

		spCondition := stk.Pop(token.I32) // 判断条件

		// wasm 2.0 支持带类型
		valType := i.ResultTyp
		if valType == 0 {
			// 不带类型, 2个数值类型必须一样
			valType = stk.TopToken()
		}
		spValueFalse := stk.Pop(valType)
		spValueTrue := stk.Pop(valType)

		ret0 := stk.Push(valType)

		// 注意返回值的顺序
		// if sp0 != 0 { sp2 } else { sp1 }
		switch valType {
		case token.I32:
			fmt.Fprintf(w, "    # select\n")
			fmt.Fprintf(w, "    ld.w t1, fp, %d\n", p.fnWasmR0Base-spCondition*8)
			fmt.Fprintf(w, "    ld.w t2, fp, %d\n", p.fnWasmR0Base-spValueTrue*8)
			fmt.Fprintf(w, "    ld.w t3, fp, %d\n", p.fnWasmR0Base-spValueFalse*8)
			fmt.Fprintf(w, "    masknez t0, t2, t1")
			fmt.Fprintf(w, "    maskeqz t0, t3, t1")
			fmt.Fprintf(w, "    st.w t0, fp, %d\n", p.fnWasmR0Base-ret0*8)
		case token.I64:
			fmt.Fprintf(w, "    # select\n")
			fmt.Fprintf(w, "    ld.w t1, fp, %d\n", p.fnWasmR0Base-spCondition*8)
			fmt.Fprintf(w, "    ld.d t2, fp, %d\n", p.fnWasmR0Base-spValueTrue*8)
			fmt.Fprintf(w, "    ld.d t3, fp, %d\n", p.fnWasmR0Base-spValueFalse*8)
			fmt.Fprintf(w, "    masknez t0, t2, t1")
			fmt.Fprintf(w, "    maskeqz t0, t3, t1")
			fmt.Fprintf(w, "    st.d t0, fp, %d\n", p.fnWasmR0Base-ret0*8)
		case token.F32:
			// TODO: 浮点数选择需要验证
			fmt.Fprintf(w, "    # select\n")
			fmt.Fprintf(w, "    ld.w t1, fp, %d\n", p.fnWasmR0Base-spCondition*8)
			fmt.Fprintf(w, "    fld.s ft2, fp, %d\n", p.fnWasmR0Base-spValueTrue*8)
			fmt.Fprintf(w, "    fld.s ft3, fp, %d\n", p.fnWasmR0Base-spValueFalse*8)
			fmt.Fprintf(w, "    movgr2cf fcc0, t1")
			fmt.Fprintf(w, "    fsel ft0, ft2, ft3, ")
			fmt.Fprintf(w, "    fst.s t1, fp, %d\n", p.fnWasmR0Base-ret0*8)
		case token.F64:
			// TODO: 浮点数选择需要验证
			fmt.Fprintf(w, "    # select\n")
			fmt.Fprintf(w, "    ld.w t1, fp, %d\n", p.fnWasmR0Base-spCondition*8)
			fmt.Fprintf(w, "    fld.d ft2, fp, %d\n", p.fnWasmR0Base-spValueTrue*8)
			fmt.Fprintf(w, "    fld.d ft3, fp, %d\n", p.fnWasmR0Base-spValueFalse*8)
			fmt.Fprintf(w, "    movgr2cf fcc0, t1")
			fmt.Fprintf(w, "    fsel ft0, ft2, ft3, ")
			fmt.Fprintf(w, "    fst.d t1, fp, %d\n", p.fnWasmR0Base-ret0*8)
		default:
			unreachable()
		}
	case token.INS_LOCAL_GET:
		i := i.(ast.Ins_LocalGet)
		xType := p.findLocalType(fn, i.X)
		xOff := p.findLocalOffset(fnNative, fn, i.X)
		ret0 := stk.Push(xType)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "    ld.w t0, fp, %d # %s\n", xOff, p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    st.w t0, fp, %d # %s\n", p.fnWasmR0Base-ret0*8, p.findLocalName(fn, i.X))
		case token.I64:
			fmt.Fprintf(w, "    ld.d t0, fp, %d # %s\n", xOff, p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    st.d t0, fp, %d # %s\n", p.fnWasmR0Base-ret0*8, p.findLocalName(fn, i.X))
		case token.F32:
			fmt.Fprintf(w, "    fld.s ft0, fp, %d # %s\n", xOff, p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    fst.s ft0, fp, %d # %s\n", p.fnWasmR0Base-ret0*8, p.findLocalName(fn, i.X))
		case token.F64:
			fmt.Fprintf(w, "    fld.d ft0, fp, %d # %s\n", xOff, p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    fst.d ft0, fp, %d # %s\n", p.fnWasmR0Base-ret0*8, p.findLocalName(fn, i.X))
		default:
			unreachable()
		}
	case token.INS_LOCAL_SET:
		i := i.(ast.Ins_LocalSet)
		xType := p.findLocalType(fn, i.X)
		xOff := p.findLocalOffset(fnNative, fn, i.X)
		sp0 := stk.Pop(xType)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "    ld.w t0, fp, %d # %s\n", p.fnWasmR0Base-sp0*8, p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    st.w t0, fp, %d # %s\n", xOff, p.findLocalName(fn, i.X))
		case token.I64:
			fmt.Fprintf(w, "    ld.d t0, fp, %d # %s\n", p.fnWasmR0Base-sp0*8, p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    st.d t0, fp, %d # %s\n", xOff, p.findLocalName(fn, i.X))
		case token.F32:
			fmt.Fprintf(w, "    fld.w ft0, fp, %d # %s\n", p.fnWasmR0Base-sp0*8, p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    fst.w ft0, fp, %d # %s\n", xOff, p.findLocalName(fn, i.X))
		case token.F64:
			fmt.Fprintf(w, "    fld.w ft0, fp, %d # %s\n", p.fnWasmR0Base-sp0*8, p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    fst.w ft0, fp, %d # %s\n", xOff, p.findLocalName(fn, i.X))
		default:
			unreachable()
		}
	case token.INS_LOCAL_TEE:
		i := i.(ast.Ins_LocalTee)
		xType := p.findLocalType(fn, i.X)
		xOff := p.findLocalOffset(fnNative, fn, i.X)
		sp0 := stk.Top(xType)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "    ld.w t0, fp, %d # %s\n", p.fnWasmR0Base-sp0*8, p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    st.w t0, fp, %d # %s\n", xOff, p.findLocalName(fn, i.X))
		case token.I64:
			fmt.Fprintf(w, "    ld.d t0, fp, %d # %s\n", p.fnWasmR0Base-sp0*8, p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    st.d t0, fp, %d # %s\n", xOff, p.findLocalName(fn, i.X))
		case token.F32:
			fmt.Fprintf(w, "    fld.w ft0, fp, %d # %s\n", p.fnWasmR0Base-sp0*8, p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    fst.w ft0, fp, %d # %s\n", xOff, p.findLocalName(fn, i.X))
		case token.F64:
			fmt.Fprintf(w, "    fld.w ft0, fp, %d # %s\n", p.fnWasmR0Base-sp0*8, p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    fst.w ft0, fp, %d # %s\n", xOff, p.findLocalName(fn, i.X))
		default:
			unreachable()
		}
	case token.INS_GLOBAL_GET:
		i := i.(ast.Ins_GlobalGet)
		xType := p.findGlobalType(i.X)
		ret0 := stk.Push(xType)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "    # global.get %s\n", i.X)
			fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20(%s)\n", i.X)
			fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12(%s)\n", i.X)
			fmt.Fprintf(w, "    ld.w      t0, t1, 0\n")
			fmt.Fprintf(w, "    st.w      t0, fp, %d # %s\n", p.fnWasmR0Base-ret0*8, p.findLocalName(fn, i.X))
		case token.I64:
			fmt.Fprintf(w, "    # global.get %s\n", i.X)
			fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20(%s)\n", i.X)
			fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12(%s)\n", i.X)
			fmt.Fprintf(w, "    ld.d      t0, t1, 0\n")
			fmt.Fprintf(w, "    st.d      t0, fp, %d # %s\n", p.fnWasmR0Base-ret0*8, p.findLocalName(fn, i.X))
		case token.F32:
			fmt.Fprintf(w, "    # global.get %s\n", i.X)
			fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20(%s)\n", i.X)
			fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12(%s)\n", i.X)
			fmt.Fprintf(w, "    fld.s     ft0, t1, 0\n")
			fmt.Fprintf(w, "    fst.s     ft0, fp, %d # %s\n", p.fnWasmR0Base-ret0*8, p.findLocalName(fn, i.X))
		case token.F64:
			fmt.Fprintf(w, "    # global.get %s\n", i.X)
			fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20(%s)\n", i.X)
			fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12(%s)\n", i.X)
			fmt.Fprintf(w, "    fld.d     ft0, t1, 0\n")
			fmt.Fprintf(w, "    fst.d     ft0, fp, %d # %s\n", p.fnWasmR0Base-ret0*8, p.findLocalName(fn, i.X))
		default:
			unreachable()
		}
	case token.INS_GLOBAL_SET:
		i := i.(ast.Ins_GlobalSet)
		xType := p.findGlobalType(i.X)
		sp0 := stk.Pop(xType)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "    # global.set %s\n", i.X)
			fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20(%s)\n", i.X)
			fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12(%s)\n", i.X)
			fmt.Fprintf(w, "    ld.w      t0, fp, %d\n", p.fnWasmR0Base-sp0*8)
			fmt.Fprintf(w, "    st.w      t0, t1, 0\n")
		case token.I64:
			fmt.Fprintf(w, "    # global.set %s\n", i.X)
			fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20(%s)\n", i.X)
			fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12(%s)\n", i.X)
			fmt.Fprintf(w, "    ld.d      t0, fp, %d\n", p.fnWasmR0Base-sp0*8)
			fmt.Fprintf(w, "    st.d      t0, t1, 0\n")
		case token.F32:
			fmt.Fprintf(w, "    # global.set %s\n", i.X)
			fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20(%s)\n", i.X)
			fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12(%s)\n", i.X)
			fmt.Fprintf(w, "    fld.s     ft0, fp, %d\n", p.fnWasmR0Base-sp0*8)
			fmt.Fprintf(w, "    fst.s     ft0, t1, 0\n")
		case token.F64:
			fmt.Fprintf(w, "    # global.set %s\n", i.X)
			fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20(%s)\n", i.X)
			fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12(%s)\n", i.X)
			fmt.Fprintf(w, "    fld.d     ft0, fp, %d\n", p.fnWasmR0Base-sp0*8)
			fmt.Fprintf(w, "    fst.d     ft0, t1, 0\n")
		default:
			unreachable()
		}
	case token.INS_TABLE_GET:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.FUNCREF) // funcref
		fmt.Fprintf(w, "    # table.get\n")
		fmt.Fprintf(w, "    ld.w  t0, fp, %d # t0 = [pop]\n", sp0)
		fmt.Fprintf(w, "    add.d t0, t0, %s # t0 = table + t0\n", kTableReg)
		fmt.Fprintf(w, "    ld.w  t0, t0, 0  # t0 = [t0]\n")
		fmt.Fprintf(w, "    st.w  t0, fp, %d # push t0\n", ret0)
	case token.INS_TABLE_SET:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.FUNCREF) // funcref
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		fmt.Fprintf(w, "    # table.get\n")
		fmt.Fprintf(w, "    ld.w  t2, fp, %d # t2 = pop\n", sp0)
		fmt.Fprintf(w, "    ld.w  t0, fp, %d # t0 = $pop\n", sp1)
		fmt.Fprintf(w, "    add.d t0, t0, %s # t0 = table + t0\n", kTableReg)
		fmt.Fprintf(w, "    st.w  t2, t0, 0\n")
	case token.INS_I32_LOAD:
		i := i.(ast.Ins_I32Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.load\n")
		fmt.Fprintf(w, "    ld.w  t0, fp, %d # t0 = [pop]\n", sp0)
		fmt.Fprintf(w, "    add.d t0, t0, %s # t0 = mempry + t0\n", kMemoryReg)
		fmt.Fprintf(w, "    ld.w  t0, t0, %d # t0 = [t0+off]\n", i.Offset)
		fmt.Fprintf(w, "    st.w  t0, fp, %d # push t0\n", ret0)

	case token.INS_I64_LOAD:
		i := i.(ast.Ins_I64Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.load\n")
		fmt.Fprintf(w, "    ld.w  t0, fp, %d # t0 = [pop]\n", sp0)
		fmt.Fprintf(w, "    add.d t0, t0, %s # t0 = mempry + t0\n", kMemoryReg)
		fmt.Fprintf(w, "    ld.d  t0, t0, %d # t0 = [t0+off]\n", i.Offset)
		fmt.Fprintf(w, "    st.d  t0, fp, %d # push t0\n", ret0)

	case token.INS_F32_LOAD:
		i := i.(ast.Ins_F32Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # i64.load\n")
		fmt.Fprintf(w, "    ld.w  t0, fp, %d # t0 = [pop]\n", sp0)
		fmt.Fprintf(w, "    add.d t0, t0, %s # t0 = mempry + t0\n", kMemoryReg)
		fmt.Fprintf(w, "    fld.s  t0, t0, %d # t0 = [t0+off]\n", i.Offset)
		fmt.Fprintf(w, "    fst.s  t0, fp, %d # push t0\n", ret0)

	case token.INS_F64_LOAD:
		i := i.(ast.Ins_I32Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # i64.load\n")
		fmt.Fprintf(w, "    ld.w  t0, fp, %d # t0 = [pop]\n", sp0)
		fmt.Fprintf(w, "    add.d t0, t0, %s # t0 = mempry + t0\n", kMemoryReg)
		fmt.Fprintf(w, "    fld.d  t0, t0, %d # t0 = [t0+off]\n", i.Offset)
		fmt.Fprintf(w, "    fst.d  t0, fp, %d # push t0\n", ret0)

	case token.INS_I32_LOAD8_S:
		// i := i.(ast.Ins_I32Load8S)
		// sp0 := stk.Pop(token.I32)
		// ret0 := stk.Push(token.I32)
		panic("TODO")

	case token.INS_I32_LOAD8_U:
		// i := i.(ast.Ins_I32Load8U)
		// sp0 := stk.Pop(token.I32)
		// ret0 := stk.Push(token.I32)
		panic("TODO")

	case token.INS_I32_LOAD16_S:
		// i := i.(ast.Ins_I32Load16S)
		// sp0 := stk.Pop(token.I32)
		// ret0 := stk.Push(token.I32)
		panic("TODO")

	case token.INS_I32_LOAD16_U:
		// i := i.(ast.Ins_I32Load16U)
		// sp0 := stk.Pop(token.I32)
		// ret0 := stk.Push(token.I32)
		panic("TODO")

	case token.INS_I64_LOAD8_S:
		// i := i.(ast.Ins_I64Load8S)
		// sp0 := stk.Pop(token.I32)
		// ret0 := stk.Push(token.I64)
		panic("TODO")

	case token.INS_I64_LOAD8_U:
		// i := i.(ast.Ins_I64Load8U)
		// sp0 := stk.Pop(token.I32)
		// ret0 := stk.Push(token.I64)
		panic("TODO")

	case token.INS_I64_LOAD16_S:
		// i := i.(ast.Ins_I64Load16S)
		// sp0 := stk.Pop(token.I32)
		// ret0 := stk.Push(token.I64)
		panic("TODO")

	case token.INS_I64_LOAD16_U:
		// i := i.(ast.Ins_I64Load16U)
		// sp0 := stk.Pop(token.I32)
		// ret0 := stk.Push(token.I64)
		panic("TODO")

	case token.INS_I64_LOAD32_S:
		// i := i.(ast.Ins_I64Load32S)
		// sp0 := stk.Pop(token.I32)
		// ret0 := stk.Push(token.I64)
		panic("TODO")

	case token.INS_I64_LOAD32_U:
		// i := i.(ast.Ins_I64Load32U)
		// sp0 := stk.Pop(token.I32)
		// ret0 := stk.Push(token.I64)
		panic("TODO")

	case token.INS_I32_STORE:
		// i := i.(ast.Ins_I32Store)
		// sp0 := stk.Pop(token.I32)
		// sp1 := stk.Pop(token.I32)
		panic("TODO")

	case token.INS_I64_STORE:
		// i := i.(ast.Ins_I64Store)
		// sp0 := stk.Pop(token.I64)
		// sp1 := stk.Pop(token.I32)
		panic("TODO")

	case token.INS_F32_STORE:
		// i := i.(ast.Ins_F32Store)
		// sp0 := stk.Pop(token.F32)
		// sp1 := stk.Pop(token.I32)
		panic("TODO")

	case token.INS_F64_STORE:
		// i := i.(ast.Ins_F64Store)
		// sp0 := stk.Pop(token.F64)
		// sp1 := stk.Pop(token.I32)
		panic("TODO")

	case token.INS_I32_STORE8:
		// i := i.(ast.Ins_I32Store8)
		// sp0 := stk.Pop(token.I32)
		// sp1 := stk.Pop(token.I32)
		panic("TODO")

	case token.INS_I32_STORE16:
		// i := i.(ast.Ins_I32Store16)
		// sp0 := stk.Pop(token.I32)
		// sp1 := stk.Pop(token.I32)
		panic("TODO")

	case token.INS_I64_STORE8:
		// i := i.(ast.Ins_I64Store8)
		// sp0 := stk.Pop(token.I64)
		// sp1 := stk.Pop(token.I32)
		panic("TODO")

	case token.INS_I64_STORE16:
		// i := i.(ast.Ins_I64Store16)
		// sp0 := stk.Pop(token.I64)
		// sp1 := stk.Pop(token.I32)
		panic("TODO")

	case token.INS_I64_STORE32:
		// i := i.(ast.Ins_I64Store32)
		// sp0 := stk.Pop(token.I64)
		// sp1 := stk.Pop(token.I32)
		panic("TODO")

	case token.INS_MEMORY_SIZE:
		sp0 := p.fnWasmR0Base - 8*stk.Push(token.I32)
		fmt.Fprintf(w, "    # memory.size\n")
		fmt.Fprintf(w, "    pcalau12i t0, %%pc_hi20(%s)\n", kMemorySizeName)
		fmt.Fprintf(w, "    addi.d    t0, t0, %%pc_lo12(%s)\n", kMemorySizeName)
		fmt.Fprintf(w, "    st.w      t0, fp, %d\n", sp0)
	case token.INS_MEMORY_GROW:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		// 最大内存在启动时就分配好, 这里只是调整全局变量

		fmt.Fprintf(w, "    # memory.grow\n")
		fmt.Fprintf(w, "    pcalau12i t0, %%pc_hi20(%s)\n", kMemorySizeName)
		fmt.Fprintf(w, "    addi.d    t0, t0, %%pc_lo12(%s)\n", kMemorySizeName)

		fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20(%s)\n", kMemoryMaxSizeName)
		fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12(%s)\n", kMemoryMaxSizeName)

		fmt.Fprintf(w, "    ld.w      t3, fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d     t3, t3, t1\n")

		fmt.Fprintf(w, "    blt       t1, t3, L.xxx.else\n")
		fmt.Fprintf(w, "    st.w      t3, %s, 0\n", kMemorySizeName)
		fmt.Fprintf(w, "    st.w      t0, fp, %d\n", ret0)
		fmt.Fprintf(w, "    b         L.xxx.else\n")
		fmt.Fprintf(w, "L.xxx.else:\n")
		fmt.Fprintf(w, "    lui.w     t4, -1\n")
		fmt.Fprintf(w, "    st.w      t4, fp, %d\n", ret0)
		fmt.Fprintf(w, "L.xxx.end:\n")

	case token.INS_MEMORY_INIT:
		i := i.(ast.Ins_MemoryInit)

		len := stk.Pop(token.I32)
		off := stk.Pop(token.I32)
		dst := stk.Pop(token.I32)

		// data 字符串绑定到 text 段, 对于一个地址
		// 然后通过地址来初始化

		// 需要实现 memcpy/memset 函数

		var sb strings.Builder
		datai := p.m.Data[i.DataIdx].Value[off:][:len]
		for _, x := range datai {
			sb.WriteString(fmt.Sprintf("\\x%02x", x))
		}

		fmt.Fprintf(w, "%smemcpy(&memory[R%d.i32], (void*)(\"%s\"), %d); // %s\n",
			indent, dst, sb.String(), len,
			insString(i),
		)
	case token.INS_MEMORY_COPY:
		len := stk.Pop(token.I32)
		src := stk.Pop(token.I32)
		dst := stk.Pop(token.I32)
		fmt.Fprintf(w, "%smemcpy(&memory[R%d.i32], &memory[R%d.i32], R%d.i32); // %s\n",
			indent, dst, src, len,
			insString(i),
		)
	case token.INS_MEMORY_FILL:
		len := stk.Pop(token.I32)
		val := stk.Pop(token.I32)
		dst := stk.Pop(token.I32)
		fmt.Fprintf(w, "%smemset(&memory[R%d.i32], R%d.i32, R%d.i32); // %s\n",
			indent, dst, val, len,
			insString(i),
		)
	case token.INS_I32_CONST:
		i := i.(ast.Ins_I32Const)
		sp0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = %d; // %s\n", indent, sp0, i.X, insString(i))
	case token.INS_I64_CONST:
		i := i.(ast.Ins_I64Const)
		sp0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = %d; // %s\n", indent, sp0, i.X, insString(i))
	case token.INS_F32_CONST:
		i := i.(ast.Ins_F32Const)
		sp0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = %f; // %s\n", indent, sp0, i.X, insString(i))
	case token.INS_F64_CONST:
		i := i.(ast.Ins_F64Const)
		sp0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = %f; // %s\n", indent, sp0, i.X, insString(i))
	case token.INS_I32_EQZ:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.i32==0)? 1: 0; // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I32_EQ:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.i32==R%d.i32)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_NE:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.i32!=R%d.i32)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_LT_S:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.i32<R%d.i32)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_LT_U:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = ((uint32_t)(R%d.i32)<(uint32_t)(R%d.i32))? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_GT_S:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.i32>R%d.i32)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_GT_U:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = ((uint32_t)(R%d.i32)>(uint32_t)(R%d.i32))? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_LE_S:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.i32<=R%d.i32)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_LE_U:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = ((uint32_t)(R%d.i32)<=(uint32_t)(R%d.i32))? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_GE_S:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.i32>=R%d.i32)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_GE_U:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = ((uint32_t)(R%d.i32)>=(uint32_t)(R%d.i32))? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_EQZ:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.i64==0)? 1: 0; // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I64_EQ:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.i64==R%d.i64)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_NE:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.i64!=R%d.i64)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_LT_S:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.i64<R%d.i64)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_LT_U:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = ((uint64_t)(R%d.i64)<(uint64_t)(R%d.i64))? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_GT_S:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.i64>R%d.i64)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_GT_U:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = ((uint64_t)(R%d.i64)>(uint64_t)(R%d.i64))? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_LE_S:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.i64<=R%d.i64)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_LE_U:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = ((uint64_t)(R%d.i64)<=(uint64_t)(R%d.i64))? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_GE_S:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.i64>=R%d.i64)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_GE_U:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = ((uint64_t)(R%d.i64)>=(uint64_t)(R%d.i64))? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F32_EQ:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.f32==R%d.f32)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F32_NE:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.f32!=R%d.f32)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F32_LT:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.f32<R%d.f32)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F32_GT:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.f32>R%d.f32)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F32_LE:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.f32<=R%d.f32)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F32_GE:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.f32>=R%d.f32)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F64_EQ:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.f64==R%d.f64)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F64_NE:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.f64!=R%d.f64)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F64_LT:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.f64<R%d.f64)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F64_GT:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.f64>R%d.f64)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F64_LE:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.f64<=R%d.f64)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F64_GE:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (R%d.f64>=R%d.f64)? 1: 0; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_CLZ:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = I32_CLZ(R%d.i32); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I32_CTZ:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = I32_CTZ(R%d.i32); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I32_POPCNT:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = I32_POPCNT(R%d.i32); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I32_ADD:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = R%d.i32 + R%d.i32; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_SUB:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = R%d.i32 - R%d.i32; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_MUL:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = R%d.i32 * R%d.i32; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_DIV_S:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = R%d.i32 / R%d.i32; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_DIV_U:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (int32_t)((uint32_t)(R%d.i32)/(uint32_t)(R%d.i32)); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_REM_S:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = R%d.i32 %% R%d.i32; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_REM_U:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (int32_t)((uint32_t)(R%d.i32)%%(uint32_t)(R%d.i32)); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_AND:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = R%d.i32 & R%d.i32; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_OR:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = R%d.i32 | R%d.i32; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_XOR:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = R%d.i32 ^ R%d.i32; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_SHL:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = R%d.i32 << (R%d.i32&63); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_SHR_S:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = R%d.i32 >> (R%d.i32&63); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_SHR_U:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (int32_t)((uint32_t)(R%d.i32)>>(uint32_t)(R%d.i32&63)); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_ROTL:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = I32_ROTL(R%d.i32, R%d.i32); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_ROTR:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = I32_ROTR(R%d.i32, R%d.i32); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_CLZ:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = I64_CLZ(R%d.i64); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I64_CTZ:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = I64_CTZ(R%d.i64); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I64_POPCNT:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = I64_POPCNT(R%d.i64); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I64_ADD:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = R%d.i64 + R%d.i64; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_SUB:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = R%d.i64 - R%d.i64; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_MUL:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = R%d.i64 * R%d.i64; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_DIV_S:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = R%d.i64 / R%d.i64; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_DIV_U:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = (int64_t)((uint64_t)(R%d.i64)/(uint64_t)(R%d.i64)); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_REM_S:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = (int64_t)((int64_t)(R%d.i64)%%(int64_t)(R%d.i64)); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_REM_U:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = (int64_t)((uint64_t)(R%d.i64)%%(uint64_t)(R%d.i64)); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_AND:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = R%d.i64 & R%d.i64; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_OR:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = R%d.i64 | R%d.i64; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_XOR:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = R%d.i64 ^ R%d.i64; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_SHL:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = R%d.i64 << (((uint64_t)R%d.i64)&63); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_SHR_S:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = R%d.i64 >> (((uint64_t)R%d.i64)&63); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_SHR_U:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = (int64_t)((uint64_t)(R%d.i64)>>((uint64_t)(R%d.i64)&63)); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_ROTL:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = I64_ROTL(R%d.i64, R%d.i64); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I64_ROTR:
		sp0 := stk.Pop(token.I64)
		sp1 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = I64_ROTR(R%d.i64, R%d.i64); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F32_ABS:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = fabsf(R%d.f32); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F32_NEG:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = 0-R%d.f32; // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F32_CEIL:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = ceilf(R%d.f32); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F32_FLOOR:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = floorf(R%d.f32); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F32_TRUNC:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = truncf(R%d.f32); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F32_NEAREST:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = roundf(R%d.f32); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F32_SQRT:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = sqrtf(R%d.f32); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F32_ADD:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = R%d.f32 + R%d.f32; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F32_SUB:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = R%d.f32 - R%d.f32; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F32_MUL:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = R%d.f32 * R%d.f32; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F32_DIV:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = R%d.f32 / R%d.f32; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F32_MIN:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = fminf(R%d.f32, R%d.f32); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F32_MAX:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = fmaxf(R%d.f32, R%d.f32); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F32_COPYSIGN:
		sp0 := stk.Pop(token.F32)
		sp1 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = copysignf(R%d.f32, R%d.f32); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F64_ABS:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = fabs(R%d.f64); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F64_NEG:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = 0-R%d.f64; // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F64_CEIL:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = ceil(R%d.f64); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F64_FLOOR:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = floor(R%d.f64); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F64_TRUNC:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = trunc(R%d.f64); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F64_NEAREST:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = round(R%d.f64); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F64_SQRT:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = sqrt(R%d.f64); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F64_ADD:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = R%d.f64 + R%d.f64; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F64_SUB:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = R%d.f64 - R%d.f64; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F64_MUL:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = R%d.f64 * R%d.f64; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F64_DIV:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = R%d.f64 / R%d.f64; // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F64_MIN:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = fmin(R%d.f64, R%d.f64); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F64_MAX:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = fmax(R%d.f64, R%d.f64); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_F64_COPYSIGN:
		sp0 := stk.Pop(token.F64)
		sp1 := stk.Pop(token.F64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = copysign(R%d.f64, R%d.f64); // %s\n",
			indent, ret0, sp1, sp0,
			insString(i),
		)
	case token.INS_I32_WRAP_I64:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (int32_t)(R%d.i64); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I32_TRUNC_F32_S:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (int32_t)(truncf(R%d.f32)); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I32_TRUNC_F32_U:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (int32_t)(uint32_t)(truncf(R%d.f32)); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I32_TRUNC_F64_S:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (int32_t)(trunc(R%d.f64)); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I32_TRUNC_F64_U:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I32)
		fmt.Fprintf(w, "%sR%d.i32 = (int32_t)(uint32_t)(trunc(R%d.f64)); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I64_EXTEND_I32_S:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = (int64_t)(R%d.i32); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I64_EXTEND_I32_U:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = (int64_t)((uint32_t)(R%d.i32)); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I64_TRUNC_F32_S:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = (int64_t)(int32_t)(truncf(R%d.f32)); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I64_TRUNC_F32_U:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = (int64_t)(uint32_t)(truncf(R%d.f32)); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I64_TRUNC_F64_S:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = (int64_t)(trunc(R%d.f64)); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I64_TRUNC_F64_U:
		sp0 := stk.Pop(token.F64)
		ret0 := stk.Push(token.I64)
		fmt.Fprintf(w, "%sR%d.i64 = (int64_t)(uint64_t)(trunc(R%d.f64)); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F32_CONVERT_I32_S:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = (float)(R%d.i32); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F32_CONVERT_I32_U:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = (float)(R%d.u32); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F32_CONVERT_I64_S:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = (float)(R%d.i64); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F32_CONVERT_I64_U:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = (float)(R%d.u64); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F32_DEMOTE_F64:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d.f32 = (float)(R%d.f64); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F64_CONVERT_I32_S:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = (double)(R%d.i32); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F64_CONVERT_I32_U:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = (double)(R%d.u32); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F64_CONVERT_I64_S:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = (double)(R%d.i64); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F64_CONVERT_I64_U:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = (double)(R%d.u64); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F64_PROMOTE_F32:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d.f64 = (double)(R%d.f32); // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I32_REINTERPRET_F32:
		sp0 := stk.Pop(token.F32)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d = R%d; // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_I64_REINTERPRET_F64:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d = R%d; // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F32_REINTERPRET_I32:
		sp0 := stk.Pop(token.I32)
		ret0 := stk.Push(token.F32)
		fmt.Fprintf(w, "%sR%d = R%d; // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	case token.INS_F64_REINTERPRET_I64:
		sp0 := stk.Pop(token.I64)
		ret0 := stk.Push(token.F64)
		fmt.Fprintf(w, "%sR%d = R%d; // %s\n",
			indent, ret0, sp0,
			insString(i),
		)
	default:
		panic(fmt.Sprintf("unreachable: %T", i))
	}
	return nil
}
