// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strings"

	"wa-lang.org/wa/internal/native/abi"
	nativeast "wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/ast/astutil"
	"wa-lang.org/wa/internal/native/loong64"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

const (
	kFuncNamePrefix    = ".F."
	kFuncRetNamePrefix = ".F.ret."
)

const (
	kLabelPrefixName_brNext        = ".L.brNext."
	kLabelPrefixName_brCase        = ".L.brCase."
	kLabelPrefixName_brDefault     = ".L.brDefault."
	kLabelPrefixName_brFallthrough = ".L.brFallthrough."
	kLabelPrefixName_else          = ".L.else."
	kLabelPrefixName_end           = ".L.end."
	kLabelPrefixName_return        = ".L.return."
)

//
// 函数栈帧布局
// 参考 /docs/asm_abi_la64.md
//

func (p *wat2laWorker) buildFuncs(w io.Writer) error {
	if len(p.m.Funcs) == 0 {
		return nil
	}

	for _, f := range p.m.Funcs {
		p.localNames = nil
		p.localTypes = nil
		p.scopeTypes = nil
		p.scopeLabels = nil
		p.scopeLabelsSuffix = nil
		p.scopeStackBases = nil
		p.scopeResults = nil

		p.gasComment(w, "func "+f.Name+f.Type.String())
		p.gasSectionTextStart(w)
		if f.ExportName == f.Name {
			p.gasGlobal(w, kFuncNamePrefix+f.Name)
		}
		p.gasFuncStart(w, kFuncNamePrefix+f.Name)
		if err := p.buildFunc_body(w, f); err != nil {
			return err
		}
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

	assert(len(p.localNames) == 0)
	assert(len(p.localTypes) == 0)

	assert(len(p.scopeTypes) == 0)
	assert(len(p.scopeLabels) == 0)
	assert(len(p.scopeLabelsSuffix) == 0)
	assert(len(p.scopeStackBases) == 0)
	assert(len(p.scopeResults) == 0)

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
			Name: fmt.Sprintf("%s%d", kFuncRetNamePrefix, i),
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

	for _, x := range fn.Type.Params {
		p.localNames = append(p.localNames, x.Name)
		p.localTypes = append(p.localTypes, x.Type)
	}
	if len(fn.Body.Locals) > 0 {
		for _, x := range fn.Body.Locals {
			p.localNames = append(p.localNames, x.Name)
			p.localTypes = append(p.localTypes, x.Type)
			p.gasCommentInFunc(w, fmt.Sprintf("local %s: %s", x.Name, x.Type))
		}
		fmt.Fprintln(&bufHeader)
	}

	// 模拟构建函数栈帧
	// 后面需要根据参数是否走寄存器传递生成相关的入口代码和返回代码
	if err := astutil.BuildFuncFrame(p.cpuType, fnNative); err != nil {
		return err
	}

	// WASM 栈的开始位置
	// WASM 栈底第一个元素相对于 $fp 的偏移位置, 每个元素 8 字节
	// 栈指令依赖该偏移量
	p.fnWasmR0Base = 0 - fnNative.FrameSize
	p.fnMaxCallArgsSize = 0

	// 开始解析 wasm 指令
	var stk valueTypeStack
	var bufIns bytes.Buffer
	{
		stk.funcName = fn.Name

		assert(stk.Len() == 0)
		for _, ins := range fn.Body.Insts {
			if err := p.buildFunc_ins(&bufIns, fnNative, fn, &stk, ins); err != nil {
				return err
			}

			// 跳过后续的死代码分析
			if ins.Token().IsTerminal() {
				break
			}
		}
	}

	// 补充构建头部栈帧
	{
		frameSize := fnNative.FrameSize + stk.MaxDepth()*8 + p.fnMaxCallArgsSize

		// 对齐到 16 字节
		if x := frameSize; x%16 != 0 {
			x = ((x + 15) / 16) * 16
			frameSize = x
		}

		fmt.Fprintf(&bufHeader, "    addi.d  $sp, $sp, -16\n")
		fmt.Fprintf(&bufHeader, "    st.d    $ra, $sp, 8\n")
		fmt.Fprintf(&bufHeader, "    st.d    $fp, $sp, 0\n")
		fmt.Fprintf(&bufHeader, "    addi.d  $fp, $sp, 0\n")
		fmt.Fprintf(&bufHeader, "    addi.d  $sp, $sp, %d\n", 0-frameSize)
		fmt.Fprintln(&bufHeader)
	}

	// 走栈返回
	if len(fnNative.Type.Return) > 1 && fnNative.Type.Return[1].Reg == 0 {
		p.gasCommentInFunc(&bufHeader, "将返回地址备份到栈")
		fmt.Fprintf(&bufHeader, "    st.d $a0, $fp, %d # return address\n", 2*8)
		fmt.Fprintln(&bufHeader)
	}

	// 将寄存器参数备份到栈
	if len(fn.Type.Params) > 0 {
		p.gasCommentInFunc(&bufHeader, "将寄存器参数备份到栈")
		for i, arg := range fnNative.Type.Args {
			if arg.Reg == 0 {
				continue // 走栈的输入参数已经在栈中
			}

			// 将寄存器中的参数存储到对于的栈帧中
			switch fn.Type.Params[i].Type {
			case token.I32:
				fmt.Fprintf(&bufHeader, "    st.w %s, $fp, %d # save arg.%d\n",
					loong64.RegString(arg.Reg), arg.RBPOff,
					i,
				)
			case token.I64:
				fmt.Fprintf(&bufHeader, "    st.d %s, $fp, %d # save arg.%d\n",
					loong64.RegString(arg.Reg), arg.RBPOff,
					i,
				)
			case token.F32:
				fmt.Fprintf(&bufHeader, "    st.w %s, $fp, %d # save arg.%d\n",
					loong64.RegString(arg.Reg), arg.RBPOff,
					i,
				)
			case token.F64:
				fmt.Fprintf(&bufHeader, "    st.d %s, $fp, %d # save arg.%d\n",
					loong64.RegString(arg.Reg), arg.RBPOff,
					i,
				)
			default:
				panic("unreachable")
			}
		}
	} else {
		p.gasCommentInFunc(&bufHeader, "没有参数需要备份到栈")
		fmt.Fprintln(&bufHeader)
	}

	// 返回值初始化为 0
	if len(fnNative.Type.Return) > 0 {
		p.gasCommentInFunc(&bufHeader, "将返回值变量初始化为0")
		for i, ret := range fnNative.Type.Return {
			fmt.Fprintf(&bufHeader,
				"    st.d $zero, $fp, %d # ret.%d = 0\n",
				ret.RBPOff, i,
			)
		}
		fmt.Fprintln(&bufHeader)
	} else {
		p.gasCommentInFunc(&bufHeader, "没有返回值变量需要初始化为0")
		fmt.Fprintln(&bufHeader)
	}

	// 局部变量初始化为 0
	if len(fnNative.Body.Locals) > 0 {
		p.gasCommentInFunc(&bufHeader, "将局部变量初始化为0")
		for _, local := range fnNative.Body.Locals {
			fmt.Fprintf(&bufHeader, "    st.d   $zero, $fp, %d # local %s = 0\n",
				local.RBPOff, local.Name,
			)
		}
		fmt.Fprintln(&bufHeader)
	} else {
		p.gasCommentInFunc(&bufHeader, "没有局部变量需要初始化为0")
		fmt.Fprintln(&bufHeader)
	}

	// 将栈上的值保存到返回值变量
	assert(stk.Len() == len(fn.Type.Results))
	if n := len(fn.Type.Results); n > 0 {
		p.gasCommentInFunc(&bufIns, "return")
		p.gasCommentInFunc(&bufIns, "copy result from stack")
		for i := n - 1; i >= 0; i-- {
			reti := fnNative.Type.Return[i]
			sp := p.fnWasmR0Base - 8*stk.Pop(fn.Type.Results[i]) - 8
			switch fn.Type.Results[i] {
			case token.I32:
				fmt.Fprintf(&bufIns, "    ld.w $t0, $fp, %d\n", sp)
				fmt.Fprintf(&bufIns, "    st.w $t0, $fp, %d # ret.%d\n", reti.RBPOff, i)
			case token.I64:
				fmt.Fprintf(&bufIns, "    ld.d $t0, $fp, %d\n", sp)
				fmt.Fprintf(&bufIns, "    st.d $t0, $fp, %d # ret.%d\n", reti.RBPOff, i)
			case token.F32:
				fmt.Fprintf(&bufIns, "    ld.w $t0, $fp, %d\n", sp)
				fmt.Fprintf(&bufIns, "    st.w $t0, $fp, %d # ret.%d\n", reti.RBPOff, i)
			case token.F64:
				fmt.Fprintf(&bufIns, "    ld.d $t0, $fp, %d\n", sp)
				fmt.Fprintf(&bufIns, "    st.d $t0, $fp, %d # ret.%d\n", reti.RBPOff, i)
			default:
				panic("unreachable")
			}
		}
	}
	assert(stk.Len() == 0)

	// 根据ABI处理返回值
	// 需要将栈顶位置的结果转化为本地ABI规范的返回值
	{
		labelRetuenId := p.makeLabelId(kLabelPrefixName_return, fn.Name, "")

		// 返回代码位置
		fmt.Fprintln(&bufIns)
		p.gasCommentInFunc(&bufIns, "根据ABI处理返回值")
		p.gasFuncLabel(&bufIns, labelRetuenId)

		// 如果走内存, 返回地址
		if len(fn.Type.Results) > 1 && fnNative.Type.Return[1].Reg == 0 {
			fmt.Fprintf(&bufIns, "    ld.d $a0, $fp, %d # ret address\n",
				fnNative.Type.Return[0].RBPOff,
			)
		} else {
			// 如果走寄存器, 则复制寄存器
			for i, ret := range fnNative.Type.Return {
				if ret.Reg == 0 {
					continue
				}
				switch fn.Type.Results[i] {
				case token.I32:
					fmt.Fprintf(&bufIns, "    ld.w %v, $fp, %d # ret %s\n",
						loong64.RegString(ret.Reg),
						fnNative.Type.Return[i].RBPOff,
						fnNative.Type.Return[i].Name,
					)
				case token.I64:
					fmt.Fprintf(&bufIns, "    ld.d %v, $fp, %d # ret %s\n",
						loong64.RegString(ret.Reg),
						fnNative.Type.Return[i].RBPOff,
						fnNative.Type.Return[i].Name,
					)
				case token.F32:
					fmt.Fprintf(&bufIns, "    ld.w %v, $fp, %d # ret %s\n",
						loong64.RegString(ret.Reg),
						fnNative.Type.Return[i].RBPOff,
						fnNative.Type.Return[i].Name,
					)
				case token.F64:
					fmt.Fprintf(&bufIns, "    ld.d %v, $fp, %d # ret %s\n",
						loong64.RegString(ret.Reg),
						fnNative.Type.Return[i].RBPOff,
						fnNative.Type.Return[i].Name,
					)
				}
			}
		}
	}
	fmt.Fprintln(&bufIns)

	// 结束栈帧
	{
		p.gasCommentInFunc(&bufIns, "函数返回")
		fmt.Fprintln(&bufIns, "    addi.d  $sp, $fp, 0")
		fmt.Fprintln(&bufIns, "    ld.d    $ra, $sp, 8")
		fmt.Fprintln(&bufIns, "    ld.d    $fp, $sp, 0")
		fmt.Fprintln(&bufIns, "    addi.d  $sp, $sp, 16")
		fmt.Fprintln(&bufIns, "    jirl    $zero, $ra, 0")
		fmt.Fprintln(&bufIns)
	}

	// 头部赋值到 w
	io.Copy(w, &bufHeader)
	// 指令复制到 w
	io.Copy(w, &bufIns)

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
		p.gasCommentInFunc(w, "unreachable")
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kRuntimePanic)
		fmt.Fprintf(w, "    addi.d $t0, $t0, %%pc_lo12(%s)\n", kRuntimePanic)
		fmt.Fprintf(w, "    jirl $ra, $t0, 0\n")
		fmt.Fprintln(w)

	case token.INS_NOP:
		fmt.Fprintln(w, indent, "addi.w $zero, $zero, 0 # nop")
		fmt.Fprintln(w)

	case token.INS_BLOCK:
		i := i.(ast.Ins_Block)

		stkBase := stk.Len()
		defer func() { assert(stk.Len() == stkBase+len(i.Results)) }()

		labelName := i.Label
		labelSuffix := p.genNextId()
		labelBlockNextId := p.makeLabelId(kLabelPrefixName_brNext, labelName, labelSuffix)

		p.enterLabelScope(token.INS_BLOCK, stkBase, i.Label, labelSuffix, i.Results)
		defer p.leaveLabelScope()

		fmt.Fprintf(w, "    # block.begin(name=%s, suffix=%s)\n", labelName, labelSuffix)

		for _, ins := range i.List {
			if err := p.buildFunc_ins(w, fnNative, fn, stk, ins); err != nil {
				return err
			}

			// 跳过后续的死代码分析
			if ins.Token().IsTerminal() {
				break
			}
		}

		p.gasFuncLabel(w, labelBlockNextId)
		fmt.Fprintf(w, "    # block.end(name=%s, suffix=%s)\n", labelName, labelSuffix)
		fmt.Fprintln(w)

	case token.INS_LOOP:
		i := i.(ast.Ins_Loop)

		stkBase := stk.Len()
		defer func() { assert(stk.Len() == stkBase+len(i.Results)) }()

		labelName := i.Label
		labelSuffix := p.genNextId()
		labelLoopNextId := p.makeLabelId(kLabelPrefixName_brNext, labelName, labelSuffix)

		p.enterLabelScope(token.INS_LOOP, stkBase, i.Label, labelSuffix, i.Results)
		defer p.leaveLabelScope()

		fmt.Fprintf(w, "    # loop.begin(name=%s, suffix=%s)\n", labelName, labelSuffix)
		p.gasFuncLabel(w, labelLoopNextId)

		for _, ins := range i.List {
			if err := p.buildFunc_ins(w, fnNative, fn, stk, ins); err != nil {
				return err
			}

			// 跳过后续的死代码分析
			if ins.Token().IsTerminal() {
				break
			}
		}

		fmt.Fprintf(w, "    # loop.end(name=%s, suffix=%s)\n", labelName, labelSuffix)
		fmt.Fprintln(w)

	case token.INS_IF:
		i := i.(ast.Ins_If)

		labelName := i.Label
		labelSuffix := p.genNextId()
		labelIfElseId := p.makeLabelId(kLabelPrefixName_else, labelName, labelSuffix)
		labelNextId := p.makeLabelId(kLabelPrefixName_brNext, labelName, labelSuffix)

		// 龙芯没有 pop 指令，需要2个指令才能实现
		// 因此通过固定的偏移量，只需要一个指令

		sp0 := stk.Pop(token.I32)

		fmt.Fprintf(w, "    # if.begin(name=%s, suffix=%s)\n", labelName, labelSuffix)
		fmt.Fprintf(w, "    ld.w $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
		if len(i.Else) > 0 {
			fmt.Fprintf(w, "    beqz $t0, %s\n", labelIfElseId)
		} else {
			fmt.Fprintf(w, "    beqz $t0, %s\n", labelNextId)
		}

		stkBase := stk.Len()
		defer func() { assert(stk.Len() == stkBase+len(i.Results)) }()

		p.enterLabelScope(token.INS_IF, stkBase, i.Label, labelSuffix, i.Results)
		defer p.leaveLabelScope()

		fmt.Fprintf(w, "    # if.body(name=%s, suffix=%s)\n", labelName, labelSuffix)
		for _, ins := range i.Body {
			if err := p.buildFunc_ins(w, fnNative, fn, stk, ins); err != nil {
				return err
			}

			// 跳过后续的死代码分析
			if ins.Token().IsTerminal() {
				break
			}
		}
		fmt.Fprintf(w, "    b   %s\n", labelNextId)
		fmt.Fprintln(w)

		if len(i.Else) > 0 {
			p.Tracef("buildFunc_ins: %s%s begin: %v\n", indent, token.INS_ELSE, stk.String())
			defer func() { p.Tracef("buildFunc_ins: %s%s end: %v\n", indent, token.INS_ELSE, stk.String()) }()

			fmt.Fprintf(w, "    # if.else(name=%s, suffix=%s)\n", labelName, labelSuffix)
			fmt.Fprintf(w, "%s:\n", labelIfElseId)

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

				// 跳过后续的死代码分析
				if ins.Token().IsTerminal() {
					break
				}
			}
		}
		p.gasFuncLabel(w, labelNextId)
		fmt.Fprintf(w, "    # if.end(name=%s, suffix=%s)\n", labelName, labelSuffix)
		fmt.Fprintln(w)

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
		labelSuffix := p.findLabelSuffixId(i.X)
		labelNextId := p.makeLabelId(kLabelPrefixName_brNext, labelName, labelSuffix)
		labelScopeTyp := p.findLabelScopeType(i.X)

		destScopeStackBase := p.scopeStackBases[len(p.scopeLabels)-labelIdx-1]
		destScopeResults := p.scopeResults[len(p.scopeLabels)-labelIdx-1]

		currentScopeStackBase := p.scopeStackBases[len(p.scopeLabels)-1]

		fmt.Fprintf(w, "    # br %s\n", i.X)

		// br对应的block带返回值
		// 如果目标是 loop 则不需要处理结果, 因为还要继续循环
		if labelScopeTyp != token.INS_LOOP && len(destScopeResults) > 0 {
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
						fmt.Fprintf(w, "    ld.w $t0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
						fmt.Fprintf(w, "    st.w $t0, $fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8-8)
					case token.I64:
						fmt.Fprintf(w, "    ld.d $t0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
						fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8-8)
					case token.F32:
						fmt.Fprintf(w, "    fld.s $ft0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
						fmt.Fprintf(w, "    fst.s $ft0, $fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8-8)
					case token.F64:
						fmt.Fprintf(w, "    fld.d $ft0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
						fmt.Fprintf(w, "    fst.d $ft0, $fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8-8)
					default:
						unreachable()
					}
				}
			}
		}

		// 清除当前block的栈中除了目标返回值剩余的值
		// 这个操作只是为了退出当前block, 因为br已经是最后一个指令(br-if 和 br-table 需要不同的处理规则)
		// 外层的block需要一个清理之后的栈
		// 后续的死代码将被忽略
		for stk.Len() > currentScopeStackBase {
			stk.DropAny()
		}

		// 退出当前block时, stack已经被清理
		// 中间栈帧的数据会在外层block指令时处理
		assert(stk.Len() == currentScopeStackBase)

		fmt.Fprintf(w, "    b %s\n", labelNextId)
		fmt.Fprintln(w)

	case token.INS_BR_IF:
		i := i.(ast.Ins_BrIf)
		labelIdx := p.findLabelIndex(i.X)
		labelName := p.findLabelName(i.X)
		labelSuffix := p.findLabelSuffixId(i.X)
		labelBrNextId := p.makeLabelId(kLabelPrefixName_brNext, labelName, labelSuffix)
		labelBrFallthroughId := p.makeLabelId(kLabelPrefixName_brFallthrough, labelName, labelSuffix)
		labelScopeTyp := p.findLabelScopeType(i.X)

		destScopeStackBase := p.scopeStackBases[len(p.scopeLabels)-labelIdx-1]
		destScopeResults := p.scopeResults[len(p.scopeLabels)-labelIdx-1]

		currentScopeStackBase := p.scopeStackBases[len(p.scopeLabels)-1]

		// 弹出的是条件
		sp0 := stk.Pop(token.I32)

		fmt.Fprintf(w, "    # br_if %s\n", labelName)
		fmt.Fprintf(w, "    ld.w $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
		fmt.Fprintf(w, "    beqz $t0, %s\n", labelBrFallthroughId)

		// 对应的block带返回值
		// 如果目标是 loop 则不需要处理结果, 因为还要继续循环
		if labelScopeTyp != token.INS_LOOP && len(destScopeResults) > 0 {
			// 必须确保当前block的stk上有足够的返回值
			assert(currentScopeStackBase+len(destScopeResults) >= stk.Len())

			// 第一个返回值返回值的偏移地址
			firstResultOffset := stk.Len() - len(destScopeResults)

			// 如果返回值位置和目标block的base不一致则需要逐个复制
			if firstResultOffset > destScopeStackBase {
				// 返回值是逆序出栈
				for i := len(destScopeResults) - 1; i >= 0; i-- {
					xType := destScopeResults[i]
					reti := stk.Pop(xType)
					switch xType {
					case token.I32:
						fmt.Fprintf(w, "    ld.w $t0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
						fmt.Fprintf(w, "    st.w $t0, $fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8-8)
					case token.I64:
						fmt.Fprintf(w, "    ld.d $t0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
						fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8-8)
					case token.F32:
						fmt.Fprintf(w, "    ld.w $t0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
						fmt.Fprintf(w, "    st.w $t0, $fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8-8)
					case token.F64:
						fmt.Fprintf(w, "    ld.d $t0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
						fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8-8)
					default:
						unreachable()
					}
				}

				// br-if 因为涉及 else 分支, 需要确保 br-if 两个分支上的栈平衡
				// 因此需要将栈回退到当前 block 开始的状态
				for i := 0; i < len(destScopeResults); i++ {
					xType := destScopeResults[i]
					stk.Push(xType)
				}
			}
		}

		fmt.Fprintf(w, "    b %s\n", labelBrNextId)
		p.gasFuncLabel(w, labelBrFallthroughId)
		fmt.Fprintln(w)

	case token.INS_BR_TABLE:
		i := i.(ast.Ins_BrTable)
		assert(len(i.XList) > 1)

		// br-table的行为和br比较相似, 因此不涉及else部分不用担心栈平衡的问题.
		// 但是每个目标block的返回值必须完全一致

		labelSuffix := p.genNextId()

		// 默认分支
		defaultLabelIdx := p.findLabelIndex(i.XList[len(i.XList)-1])
		defaultLabelName := p.findLabelName(i.XList[len(i.XList)-1])
		defaultLabelSuffix := p.findLabelSuffixId(i.XList[len(i.XList)-1])
		defaultScopeResults := p.scopeResults[len(p.scopeLabels)-defaultLabelIdx-1]
		defaultLabelNextId := p.makeLabelId(kLabelPrefixName_brNext, defaultLabelName, defaultLabelSuffix)

		currentScopeStackBase := p.scopeStackBases[len(p.scopeLabels)-1]

		sp0 := stk.Pop(token.I32)

		fmt.Fprintf(w, "    # br_table\n")
		fmt.Fprintf(w, "    ld.w $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)

		// 生成跳转链
		for k := 0; k < len(i.XList); k++ {
			if k < len(i.XList)-1 {
				assert(k < 2048) // 小常量范围
				fmt.Fprintf(w, "    # br_table case %d\n", k)
				fmt.Fprintf(w, "    addi.d $t1, $zero, %d\n", k)
				fmt.Fprintf(w, "    sub.d  $t1, $t1, $t0\n")
				fmt.Fprintf(w, "    beqz   %s\n", p.makeLabelId(kLabelPrefixName_brCase, "", labelSuffix))
			} else {
				fmt.Fprintf(w, "    # br_table default\n")
				fmt.Fprintf(w, "    b   %s\n", p.makeLabelId(kLabelPrefixName_brDefault, "", labelSuffix))
			}
		}

		// 定义每个分支的跳转代码
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
				labelSuffix := p.findLabelSuffixId(i.XList[k])
				labelNextId := p.makeLabelId(kLabelPrefixName_brNext, labelName, labelSuffix)

				destScopeStackBase := p.scopeStackBases[len(p.scopeLabels)-labelIdx-1]
				destScopeResults := p.scopeResults[len(p.scopeLabels)-labelIdx-1]

				// 验证每个目标返回值必须和default一致
				assert(len(defaultScopeResults) == len(destScopeResults))
				for i := 0; i < len(defaultScopeResults); i++ {
					assert(defaultScopeResults[i] == destScopeResults[i])
				}

				if k < len(i.XList)-1 {
					p.gasFuncLabel(w, p.makeLabelId(kLabelPrefixName_brCase, "", labelSuffix))
				} else {
					p.gasFuncLabel(w, p.makeLabelId(kLabelPrefixName_brDefault, "", labelSuffix))
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
						for i := 0; i < len(destScopeResults); i++ {
							xType := destScopeResults[i]
							reti := retIdxList[i]
							switch xType {
							case token.I32:
								fmt.Fprintf(w, "    ld.w $t1, $fp, %d # copy result\n", p.fnWasmR0Base-reti*8-8)
								fmt.Fprintf(w, "    st.w $t1, $fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8-8)
							case token.I64:
								fmt.Fprintf(w, "    ld.d $t1, $fp, %d # copy result\n", p.fnWasmR0Base-reti*8-8)
								fmt.Fprintf(w, "    st.d $t1, $fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8-8)
							case token.F32:
								fmt.Fprintf(w, "    ld.w $t1, $fp, %d # copy result\n", p.fnWasmR0Base-reti*8-8)
								fmt.Fprintf(w, "    st.w $t1, $fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8-8)
							case token.F64:
								fmt.Fprintf(w, "    ld.d $t1, $fp, %d # copy result\n", p.fnWasmR0Base-reti*8-8)
								fmt.Fprintf(w, "    st.d $t1, $fp, %d\n", p.fnWasmR0Base-(destScopeStackBase+i)*8-8)
							default:
								unreachable()
							}
						}
					}
				}

				// 跳转到目标
				if k < len(i.XList)-1 {
					fmt.Fprintf(w, "    b %s\n", labelNextId)
				} else {
					fmt.Fprintf(w, "    b %s\n", defaultLabelNextId)
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
		fmt.Fprintln(w)

	case token.INS_RETURN:
		labelReturnId := p.makeLabelId(kLabelPrefixName_return, fn.Name, "")
		p.gasCommentInFunc(w, "return")
		switch len(fn.Type.Results) {
		case 0:
			fmt.Fprintf(w, "    b %s\n", labelReturnId)
		case 1:
			sp0 := stk.Pop(fn.Type.Results[0])
			switch fn.Type.Results[0] {
			case token.I32:
				fmt.Fprintf(w, "    ld.w $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
				fmt.Fprintf(w, "    st.w $t0, $fp, %d\n", fnNative.Type.Return[0].RBPOff)
				fmt.Fprintf(w, "    b    %s\n", labelReturnId)
			case token.I64:
				fmt.Fprintf(w, "    ld.d $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
				fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", fnNative.Type.Return[0].RBPOff)
				fmt.Fprintf(w, "    b    %s\n", labelReturnId)
			case token.F32:
				fmt.Fprintf(w, "    ld.w $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
				fmt.Fprintf(w, "    st.w $t0, $fp, %d\n", fnNative.Type.Return[0].RBPOff)
				fmt.Fprintf(w, "    b     %s\n", labelReturnId)
			case token.F64:
				fmt.Fprintf(w, "    ld.d $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
				fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", fnNative.Type.Return[0].RBPOff)
				fmt.Fprintf(w, "    b     %s\n", labelReturnId)
			default:
				unreachable()
			}
		default:
			for i := len(fn.Type.Results) - 1; i >= 0; i-- {
				xType := fn.Type.Results[i]
				spi := stk.Pop(xType)
				switch xType {
				case token.I32:
					fmt.Fprintf(w, "    $ld.w $t0, $fp, %d\n", p.fnWasmR0Base-spi*8-8)
					fmt.Fprintf(w, "    $st.w $t0, $fp, %d\n", fnNative.Type.Return[i].RBPOff)
				case token.I64:
					fmt.Fprintf(w, "    ld.d $t0, $fp, %d\n", p.fnWasmR0Base-spi*8-8)
					fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", fnNative.Type.Return[i].RBPOff)
				case token.F32:
					fmt.Fprintf(w, "    $ld.w $t0, $fp, %d\n", p.fnWasmR0Base-spi*8-8)
					fmt.Fprintf(w, "    $st.w $t0, $fp, %d\n", fnNative.Type.Return[i].RBPOff)
				case token.F64:
					fmt.Fprintf(w, "    ld.d $t0, $fp, %d\n", p.fnWasmR0Base-spi*8-8)
					fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", fnNative.Type.Return[i].RBPOff)
				default:
					unreachable()
				}
			}
			fmt.Fprintf(w, "    b    %s\n", labelReturnId)
		}
		assert(stk.Len() == 0)

	case token.INS_CALL:
		i := i.(ast.Ins_Call)

		fnCallType := p.findFuncType(i.X)
		fnCallIdx := p.findFuncIndex(i.X)

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
				Name: fmt.Sprintf("%s%d", kFuncRetNamePrefix, i),
				Type: wat2nativeType(typ),
				Cap:  1,
			}
		}

		// 模拟构建函数栈帧
		// 后面需要根据参数是否走寄存器传递生成相关的入口代码和返回代码
		if err := astutil.BuildFuncFrame(p.cpuType, fnCallNative); err != nil {
			return err
		}

		// 统计调用子函数需要的最大栈空间
		if fnCallNative.ArgsSize > p.fnMaxCallArgsSize {
			p.fnMaxCallArgsSize = fnCallNative.ArgsSize
		}

		// 参数列表
		// 出栈的顺序相反
		argList := make([]int, len(fnCallType.Params))
		for k := len(argList) - 1; k >= 0; k-- {
			x := fnCallType.Params[k]
			argList[k] = stk.Pop(x.Type)
		}

		p.gasCommentInFunc(w, fmt.Sprintf("call %s(...)", i.X))

		// 如果是走栈返回, 第一个是隐藏参数
		if len(fnCallNative.Type.Return) > 0 && fnCallNative.Type.Return[0].Reg == 0 {
			assert(p.cpuType == abi.LOONG64)
			fmt.Fprintf(w, "    addi.d $a0, $sp, %d # return address\n",
				fnCallNative.ArgsSize-len(fnCallNative.Type.Return)*8,
			)
		}

		// 准备调用参数
		for k, x := range fnCallType.Params {
			switch x.Type {
			case token.I32:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    ld.w %s, $fp, %d # arg %d\n",
						"$"+strings.ToLower(loong64.RegAliasString(arg.Reg)),
						p.fnWasmR0Base-argList[k]*8-8,
						k,
					)
				} else {
					fmt.Fprintf(w, "    ld.w $t0, $fp, %d\n",
						p.fnWasmR0Base-argList[k]*8-8,
					)
					fmt.Fprintf(w, "    st.w $t0, $sp, %d\n",
						arg.RSPOff,
					)
				}
			case token.I64:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    ld.d %s, $fp, %d # arg %d\n",
						"$"+strings.ToLower(loong64.RegAliasString(arg.Reg)),
						p.fnWasmR0Base-argList[k]*8-8,
						k,
					)
				} else {
					fmt.Fprintf(w, "    ld.d $t0, $fp, %d\n",
						p.fnWasmR0Base-argList[k]*8-8,
					)
					fmt.Fprintf(w, "    st.d $t0, $sp, %d\n",
						arg.RSPOff,
					)
				}
			case token.F32:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					if arg.Reg >= loong64.REG_FA0 && arg.Reg <= loong64.REG_A7 {
						fmt.Fprintf(w, "    fld.s %s, $fp, %d # arg %d\n",
							"$"+strings.ToLower(loong64.RegAliasString(arg.Reg)),
							p.fnWasmR0Base-argList[k]*8-8,
							k,
						)
					} else {
						fmt.Fprintf(w, "    ld.w %s, $fp, %d\n",
							"$"+strings.ToLower(loong64.RegAliasString(arg.Reg)),
							p.fnWasmR0Base+argList[k]*8,
						)
					}
				} else {
					fmt.Fprintf(w, "    fld.s $t0, $fp, %d\n",
						p.fnWasmR0Base-argList[k]*8-8,
					)
					fmt.Fprintf(w, "    fst.s $t0, $sp, %d\n",
						arg.RSPOff,
					)
				}
			case token.F64:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					if arg.Reg >= loong64.REG_FA0 && arg.Reg <= loong64.REG_A7 {
						fmt.Fprintf(w, "    fld.d %s, $fp, %d # arg %d\n",
							"$"+strings.ToLower(loong64.RegAliasString(arg.Reg)),
							p.fnWasmR0Base-argList[k]*8-8,
							k,
						)
					} else {
						fmt.Fprintf(w, "    ld.d %s, $fp, %d\n",
							"$"+strings.ToLower(loong64.RegAliasString(arg.Reg)),
							p.fnWasmR0Base-argList[k]*8-8,
						)
					}
				} else {
					fmt.Fprintf(w, "    fld.d $t0, $fp, %d\n",
						p.fnWasmR0Base-argList[k]*8-8,
					)
					fmt.Fprintf(w, "    fst.d $t0, $sp, %d\n",
						arg.RSPOff,
					)
				}
			default:
				unreachable()
			}
		}

		if fnCallIdx <= p.importFuncCount {
			fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kImportNamePrefix+i.X)
			fmt.Fprintf(w, "    addi.d $t0, $t0, %%pc_lo12(%s)\n", kImportNamePrefix+i.X)
			fmt.Fprintf(w, "    jirl $ra, $t0, 0\n")
		} else {
			fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kFuncNamePrefix+i.X)
			fmt.Fprintf(w, "    addi.d $t0, $t0, %%pc_lo12(%s)\n", kFuncNamePrefix+i.X)
			fmt.Fprintf(w, "    jirl $ra, $t0, 0\n")
		}

		// 提取返回值
		if len(fnCallNative.Type.Return) > 1 && fnCallNative.Type.Return[1].Reg == 0 {
			// 走栈返回结果, $a0 是地址
			// 根据 ABI 规范获取返回值, 再保存到 WASM 栈中
			for k, retType := range fnCallType.Results {
				reti := stk.Push(retType)
				switch retType {
				case token.I32:
					fmt.Fprintf(w, "    ld.w $t0, $a0, %d\n", k*8)
					fmt.Fprintf(w, "    st.w $t0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
				case token.I64:
					fmt.Fprintf(w, "    ld.d $t0, $a0, %d\n", k*8)
					fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
				case token.F32:
					fmt.Fprintf(w, "    ld.w $t0, $a0, %d\n", k*8)
					fmt.Fprintf(w, "    st.w $t0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
				case token.F64:
					fmt.Fprintf(w, "    ld.d $t0, $a0, %d\n", k*8)
					fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
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
					fmt.Fprintf(w, "    st.w %s, $fp, %d\n",
						"$"+strings.ToLower(loong64.RegAliasString(retNative.Reg)),
						p.fnWasmR0Base-reti*8-8,
					)
				case token.I64:
					fmt.Fprintf(w, "    st.d %s, $fp, %d\n",
						"$"+strings.ToLower(loong64.RegAliasString(retNative.Reg)),
						p.fnWasmR0Base-reti*8-8,
					)
				case token.F32:
					fmt.Fprintf(w, "    fst.s %s, $fp, %d\n",
						"$"+strings.ToLower(loong64.RegAliasString(retNative.Reg)),
						p.fnWasmR0Base-reti*8-8,
					)
				case token.F64:
					fmt.Fprintf(w, "    fst.d %s, $fp, %d\n",
						"$"+strings.ToLower(loong64.RegAliasString(retNative.Reg)),
						p.fnWasmR0Base-reti*8-8,
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
			Name: "", // 间接调用, 没有名字(可以尝试根据地址查询名字)
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
				Name: fmt.Sprintf("%s%d", kFuncRetNamePrefix, i),
				Type: wat2nativeType(typ),
				Cap:  1,
			}
		}

		// 模拟构建函数栈帧
		// 后面需要根据参数是否走寄存器传递生成相关的入口代码和返回代码
		if err := astutil.BuildFuncFrame(abi.LOONG64, fnCallNative); err != nil {
			return err
		}

		// 统计调用子函数需要的最大栈空间
		if fnCallNative.ArgsSize > p.fnMaxCallArgsSize {
			p.fnMaxCallArgsSize = fnCallNative.ArgsSize
		}

		// 获取函数地址
		p.gasCommentInFunc(w, "加载函数的地址")
		fmt.Fprintln(w)

		// 需要根据索引从函数列表查询
		sp0 := stk.Pop(token.I32)
		p.gasCommentInFunc(w, "t1 = table[?]")
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kTableAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kTableAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")
		fmt.Fprintf(w, "    ld.d      $t1, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
		fmt.Fprintf(w, "    slli.d    $t1, $t1, 3 # sizeof(i64) == 8\n")
		fmt.Fprintf(w, "    add.d     $t1, $t0, $t1\n")
		fmt.Fprintf(w, "    ld.d      $t1, $t1, 0\n")
		fmt.Fprintln(w)

		p.gasCommentInFunc(w, fmt.Sprintf("t2 = %s[t1]", kTableFuncIndexListName))
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kTableFuncIndexListName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kTableFuncIndexListName)
		fmt.Fprintf(w, "    slli.d    $t1, $t1, 3 # sizeof(i64) == 8\n")
		fmt.Fprintf(w, "    add.d     $t1, $t0, $t1\n")
		fmt.Fprintf(w, "    ld.d      $t2, $t1, 0\n")
		fmt.Fprintln(w)

		p.gasCommentInFunc(w, fmt.Sprintf("call_indirect %s(...)", "$t2"))
		p.gasCommentInFunc(w, fmt.Sprintf("type %s", fnCallType))

		// 如果是走栈返回, 第一个是隐藏参数
		if len(fnCallNative.Type.Return) > 0 && fnCallNative.Type.Return[0].Reg == 0 {
			assert(p.cpuType == abi.LOONG64)
			fmt.Fprintf(w, "    addi.d $a0, $sp, %d # return address\n",
				fnCallNative.ArgsSize-len(fnCallNative.Type.Return)*8,
			)
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
					fmt.Fprintf(w, "    ld.w %s, $fp, %d # arg %d\n",
						"$"+strings.ToLower(loong64.RegAliasString(arg.Reg)),
						p.fnWasmR0Base-argList[k]*8-8,
						k,
					)
				} else {
					fmt.Fprintf(w, "    ld.w $t1, $fp, %d\n",
						p.fnWasmR0Base-argList[k]*8-8,
					)
					fmt.Fprintf(w, "    st.w $t1, $fp, %d\n",
						arg.RBPOff,
					)
				}
			case token.I64:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    ld.d %s, $fp, %d # arg %d\n",
						"$"+strings.ToLower(loong64.RegAliasString(arg.Reg)),
						p.fnWasmR0Base-argList[k]*8-8,
						k,
					)
				} else {
					fmt.Fprintf(w, "    ld.d $t1, $fp, %d\n",
						p.fnWasmR0Base-argList[k]*8-8,
					)
					fmt.Fprintf(w, "    st.d $t1, $fp, %d\n",
						arg.RBPOff,
					)
				}
			case token.F32:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    fld.s %s, $fp, %d # arg %d\n",
						"$"+strings.ToLower(loong64.RegAliasString(arg.Reg)),
						p.fnWasmR0Base-argList[k]*8-8,
						k,
					)
				} else {
					fmt.Fprintf(w, "    fld.s $t1, $fp, %d\n",
						p.fnWasmR0Base-argList[k]*8-8,
					)
					fmt.Fprintf(w, "    fst.s $t1, $fp, %d\n",
						arg.RBPOff,
					)
				}
			case token.F64:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    fld.d %s, $fp, %d # arg %d\n",
						"$"+strings.ToLower(loong64.RegAliasString(arg.Reg)),
						p.fnWasmR0Base-argList[k]*8-8,
						k,
					)
				} else {
					fmt.Fprintf(w, "    fld.d $t1, $fp, %d\n",
						p.fnWasmR0Base-argList[k]*8-8,
					)
					fmt.Fprintf(w, "    fst.d $t1, $fp, %d\n",
						arg.RBPOff,
					)
				}
			default:
				unreachable()
			}
		}

		fmt.Fprintf(w, "    jirl $ra, $t2, 0\n")

		// 提取返回值
		// TODO: 类似代码有冗余
		if len(fnCallNative.Type.Return) > 1 && fnCallNative.Type.Return[1].Reg == 0 {
			// 走栈返回结果, $a0 是地址
			// 根据 ABI 规范获取返回值, 再保存到 WASM 栈中
			for k, retType := range fnCallType.Results {
				reti := stk.Push(retType)
				switch retType {
				case token.I32:
					fmt.Fprintf(w, "    ld.w $t0, $a0, %d\n", k*8)
					fmt.Fprintf(w, "    st.w $t0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
				case token.I64:
					fmt.Fprintf(w, "    ld.d $t0, $a0, %d\n", k*8)
					fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
				case token.F32:
					fmt.Fprintf(w, "    ld.w $t0, $a0, %d\n", k*8)
					fmt.Fprintf(w, "    st.w $t0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
				case token.F64:
					fmt.Fprintf(w, "    ld.d $t0, $a0, %d\n", k*8)
					fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", p.fnWasmR0Base-reti*8-8)
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
					fmt.Fprintf(w, "    st.w %s, $fp, %d\n",
						"$"+strings.ToLower(loong64.RegAliasString(retNative.Reg)),
						p.fnWasmR0Base-reti*8-8,
					)
				case token.I64:
					fmt.Fprintf(w, "    st.d %s, $fp, %d\n",
						"$"+strings.ToLower(loong64.RegAliasString(retNative.Reg)),
						p.fnWasmR0Base-reti*8-8,
					)
				case token.F32:
					fmt.Fprintf(w, "    fst.s %s, $fp, %d\n",
						"$"+strings.ToLower(loong64.RegAliasString(retNative.Reg)),
						p.fnWasmR0Base-reti*8-8,
					)
				case token.F64:
					fmt.Fprintf(w, "    fst.d %s, $fp, %d\n",
						"$"+strings.ToLower(loong64.RegAliasString(retNative.Reg)),
						p.fnWasmR0Base-reti*8-8,
					)
				default:
					unreachable()
				}
			}
		}

	case token.INS_DROP:
		sp0 := stk.DropAny()
		fmt.Fprintf(w, "    addi.w $zero, $zero, 0 # drop [fp%+d]\n", p.fnWasmR0Base-sp0*8-8)
		fmt.Fprintln(w)

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
		fmt.Fprintf(w, "    # select\n")
		switch valType {
		case token.I32:
			fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", p.fnWasmR0Base-spCondition*8-8)
			fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", p.fnWasmR0Base-spValueTrue*8-8)
			fmt.Fprintf(w, "    ld.w    $t2, $fp, %d\n", p.fnWasmR0Base-spValueFalse*8-8)
			fmt.Fprintf(w, "    masknez $t3, $t1, $t0\n")
			fmt.Fprintf(w, "    maskeqz $t4, $t2, $t0\n")
			fmt.Fprintf(w, "    or      $t0, $t3, $t4\n")
			fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", p.fnWasmR0Base-ret0*8-8)
		case token.I64:
			fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", p.fnWasmR0Base-spCondition*8-8)
			fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", p.fnWasmR0Base-spValueTrue*8-8)
			fmt.Fprintf(w, "    ld.d    $t2, $fp, %d\n", p.fnWasmR0Base-spValueFalse*8-8)
			fmt.Fprintf(w, "    masknez $t3, $t1, $t0\n")
			fmt.Fprintf(w, "    maskeqz $t4, $t2, $t0\n")
			fmt.Fprintf(w, "    or      $t0, $t3, $t4\n")
			fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", p.fnWasmR0Base-ret0*8-8)
		case token.F32:
			fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", p.fnWasmR0Base-spCondition*8-8)
			fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", p.fnWasmR0Base-spValueTrue*8-8)
			fmt.Fprintf(w, "    ld.w    $t2, $fp, %d\n", p.fnWasmR0Base-spValueFalse*8-8)
			fmt.Fprintf(w, "    masknez $t3, $t1, $t0\n")
			fmt.Fprintf(w, "    maskeqz $t4, $t2, $t0\n")
			fmt.Fprintf(w, "    or      $t0, $t3, $t4\n")
			fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", p.fnWasmR0Base-ret0*8-8)
		case token.F64:
			fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", p.fnWasmR0Base-spCondition*8-8)
			fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", p.fnWasmR0Base-spValueTrue*8-8)
			fmt.Fprintf(w, "    ld.d    $t2, $fp, %d\n", p.fnWasmR0Base-spValueFalse*8-8)
			fmt.Fprintf(w, "    masknez $t3, $t1, $t0\n")
			fmt.Fprintf(w, "    maskeqz $t4, $t2, $t0\n")
			fmt.Fprintf(w, "    or      $t0, $t3, $t4\n")
			fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", p.fnWasmR0Base-ret0*8-8)
		default:
			unreachable()
		}
		fmt.Fprintln(w)

	case token.INS_LOCAL_GET:
		i := i.(ast.Ins_LocalGet)
		xType := p.findLocalType(fn, i.X)
		xOff := p.findLocalOffset(fnNative, fn, i.X)
		ret0 := stk.Push(xType)
		fmt.Fprintf(w, "    # local.get %s\n", i.X)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "    ld.w $t0, $fp, %d\n", xOff)
			fmt.Fprintf(w, "    st.w $t0, $fp, %d\n", p.fnWasmR0Base-ret0*8-8)
			fmt.Fprintln(w)
		case token.I64:
			fmt.Fprintf(w, "    ld.d $t0, $fp, %d\n", xOff)
			fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", p.fnWasmR0Base-ret0*8-8)
			fmt.Fprintln(w)
		case token.F32:
			fmt.Fprintf(w, "    ld.w $t0, $fp, %d\n", xOff)
			fmt.Fprintf(w, "    st.w $t0, $fp, %d\n", p.fnWasmR0Base-ret0*8-8)
			fmt.Fprintln(w)
		case token.F64:
			fmt.Fprintf(w, "    ld.d $t0, $fp, %d\n", xOff)
			fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", p.fnWasmR0Base-ret0*8-8)
			fmt.Fprintln(w)
		default:
			unreachable()
		}
	case token.INS_LOCAL_SET:
		i := i.(ast.Ins_LocalSet)
		xType := p.findLocalType(fn, i.X)
		xOff := p.findLocalOffset(fnNative, fn, i.X)
		sp0 := stk.Pop(xType)
		fmt.Fprintf(w, "   # local.set %s\n", i.X)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "    ld.w $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    st.w $t0, $fp, %d\n", xOff)
		case token.I64:
			fmt.Fprintf(w, "    ld.d $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", xOff)
		case token.F32:
			fmt.Fprintf(w, "    ld.w $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    st.w $t0, $fp, %d\n", xOff)
		case token.F64:
			fmt.Fprintf(w, "    ld.d $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", xOff)
		default:
			unreachable()
		}
		fmt.Fprintln(w)

	case token.INS_LOCAL_TEE:
		i := i.(ast.Ins_LocalTee)
		xType := p.findLocalType(fn, i.X)
		xOff := p.findLocalOffset(fnNative, fn, i.X)
		sp0 := stk.Top(xType)
		fmt.Fprintf(w, "   # local.tee %s\n", i.X)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "    ld.w $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    st.w $t0, $fp, %d\n", xOff)
		case token.I64:
			fmt.Fprintf(w, "    ld.d $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", xOff)
		case token.F32:
			fmt.Fprintf(w, "    ld.w $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    st.w $t0, $fp, %d\n", xOff)
		case token.F64:
			fmt.Fprintf(w, "    ld.d $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    st.d $t0, $fp, %d\n", xOff)
		default:
			unreachable()
		}
		fmt.Fprintln(w)

	case token.INS_GLOBAL_GET:
		i := i.(ast.Ins_GlobalGet)
		xName := kGlobalNamePrefix + p.findGlobalName(i.X)
		xType := p.findGlobalType(i.X)
		ret0 := stk.Push(xType)
		fmt.Fprintf(w, "    # global.get %s\n", i.X)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "    pcalau12i $t1, %%pc_hi20(%s)\n", xName)
			fmt.Fprintf(w, "    addi.d    $t1, $t1, %%pc_lo12(%s)\n", xName)
			fmt.Fprintf(w, "    ld.w      $t0, $t1, 0\n")
			fmt.Fprintf(w, "    st.w      $t0, $fp, %d\n", p.fnWasmR0Base-ret0*8-8)
		case token.I64:
			fmt.Fprintf(w, "    pcalau12i $t1, %%pc_hi20(%s)\n", xName)
			fmt.Fprintf(w, "    addi.d    $t1, $t1, %%pc_lo12(%s)\n", xName)
			fmt.Fprintf(w, "    ld.d      $t0, $t1, 0\n")
			fmt.Fprintf(w, "    st.d      $t0, $fp, %d\n", p.fnWasmR0Base-ret0*8-8)
		case token.F32:
			fmt.Fprintf(w, "    pcalau12i $t1, %%pc_hi20(%s)\n", xName)
			fmt.Fprintf(w, "    addi.d    $t1, $t1, %%pc_lo12(%s)\n", xName)
			fmt.Fprintf(w, "    ld.w      $t0, $t1, 0\n")
			fmt.Fprintf(w, "    st.w      $t0, $fp, %d\n", p.fnWasmR0Base-ret0*8-8)
		case token.F64:
			fmt.Fprintf(w, "    pcalau12i $t1, %%pc_hi20(%s)\n", xName)
			fmt.Fprintf(w, "    addi.d    $t1, $t1, %%pc_lo12(%s)\n", xName)
			fmt.Fprintf(w, "    ld.d      $t0, $t1, 0\n")
			fmt.Fprintf(w, "    st.d      $t0, $fp, %d\n", p.fnWasmR0Base-ret0*8-8)
		default:
			unreachable()
		}
		fmt.Fprintln(w)

	case token.INS_GLOBAL_SET:
		i := i.(ast.Ins_GlobalSet)
		xName := kGlobalNamePrefix + p.findGlobalName(i.X)
		xType := p.findGlobalType(i.X)
		sp0 := stk.Pop(xType)
		fmt.Fprintf(w, "    # global.set %s\n", i.X)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "    pcalau12i $t1, %%pc_hi20(%s)\n", xName)
			fmt.Fprintf(w, "    addi.d    $t1, $t1, %%pc_lo12(%s)\n", xName)
			fmt.Fprintf(w, "    ld.w      $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    st.w      $t0, $t1, 0\n")
		case token.I64:
			fmt.Fprintf(w, "    pcalau12i $t1, %%pc_hi20(%s)\n", xName)
			fmt.Fprintf(w, "    addi.d    $t1, $t1, %%pc_lo12(%s)\n", xName)
			fmt.Fprintf(w, "    ld.d      $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    st.d      $t0, $t1, 0\n")
		case token.F32:
			fmt.Fprintf(w, "    pcalau12i $t1, %%pc_hi20(%s)\n", xName)
			fmt.Fprintf(w, "    addi.d    $t1, $t1, %%pc_lo12(%s)\n", xName)
			fmt.Fprintf(w, "    ld.w      $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    st.w      $t0, $t1, 0\n")
		case token.F64:
			fmt.Fprintf(w, "    pcalau12i $t1, %%pc_hi20(%s)\n", xName)
			fmt.Fprintf(w, "    addi.d    $t1, $t1, %%pc_lo12(%s)\n", xName)
			fmt.Fprintf(w, "    ld.d      $t0, $fp, %d\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    st.d      $t0, $t1, 0\n")
		default:
			unreachable()
		}
		fmt.Fprintln(w)

	case token.INS_TABLE_GET:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.FUNCREF) - 8 // funcref
		fmt.Fprintf(w, "    # table.get\n")
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kTableAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kTableAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")
		fmt.Fprintf(w, "    ld.w      $t1, $fp, %d # offset\n", sp0)
		fmt.Fprintf(w, "    alsl.d    $t1, $t1, $t0, 2 # $t1 = $t1<<(2+1) + $t0\n")
		fmt.Fprintf(w, "    ld.d      $t1, $t1, 0\n")
		fmt.Fprintf(w, "    st.w      $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_TABLE_SET:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.FUNCREF) - 8 // funcref
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		fmt.Fprintf(w, "    # table.set\n")
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kTableAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kTableAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")
		fmt.Fprintf(w, "    ld.w      $t1, $fp, %d # offset\n", sp1)
		fmt.Fprintf(w, "    alsl.d    $t1, $t1, $t0, 2 # $t1 = $t1<<(2+1) + $t0\n")
		fmt.Fprintf(w, "    ld.w      $t2, $fp, %d # funcref\n", sp0)
		fmt.Fprintf(w, "    ld.d      $t2, $t1, 0\n")
		fmt.Fprintln(w)

	case token.INS_I32_LOAD:
		i := i.(ast.Ins_I32Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.load\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 加载+入栈
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.w  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintf(w, "    st.w  $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LOAD:
		i := i.(ast.Ins_I64Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.load\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 加载+入栈
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.d  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintf(w, "    st.d  $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_LOAD:
		i := i.(ast.Ins_F32Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.load\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 加载+入栈
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.w  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintf(w, "    st.w  $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_LOAD:
		i := i.(ast.Ins_F64Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.load\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 加载+入栈
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.d  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintf(w, "    st.d  $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LOAD8_S:
		i := i.(ast.Ins_I32Load8S)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.load8_s\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 加载+入栈
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.w  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintf(w, "    st.b  $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LOAD8_U:
		i := i.(ast.Ins_I32Load8U)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.load8_u\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 加载+入栈
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.w  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintf(w, "    st.b  $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LOAD16_S:
		i := i.(ast.Ins_I32Load16S)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.load16_s\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 加载+入栈
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.w  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintf(w, "    st.h  $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LOAD16_U:
		i := i.(ast.Ins_I32Load16U)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.load16_u\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 加载+入栈
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.w  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintf(w, "    st.h  $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LOAD8_S:
		i := i.(ast.Ins_I64Load8S)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.load8_s\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 加载+入栈
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.d  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintf(w, "    st.b  $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LOAD8_U:
		i := i.(ast.Ins_I64Load8U)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.load8_u\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 加载+入栈
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.d  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintf(w, "    st.b  $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LOAD16_S:
		i := i.(ast.Ins_I64Load16S)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.load16_s\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 加载+入栈
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.d  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintf(w, "    st.h  $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LOAD16_U:
		i := i.(ast.Ins_I64Load16U)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.load16_u\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 加载+入栈
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.d  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintf(w, "    st.h  $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LOAD32_S:
		i := i.(ast.Ins_I64Load32S)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.load32_s\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 加载+入栈
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.d  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintf(w, "    st.w  $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LOAD32_U:
		i := i.(ast.Ins_I64Load32U)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.load32_u\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 加载+入栈
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.d  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintf(w, "    st.w  $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_STORE:
		i := i.(ast.Ins_I32Store)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # i32.store\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 出栈+保存
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    st.w  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_I64_STORE:
		i := i.(ast.Ins_I64Store)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # i64.store\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 出栈+保存
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.d  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    st.d  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_F32_STORE:
		i := i.(ast.Ins_F32Store)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # f32.store\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 出栈+保存
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    fld.s  $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fst.s  $ft1, $t0, %d\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_F64_STORE:
		i := i.(ast.Ins_F64Store)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # f64.store\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 出栈+保存
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    fld.d  $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fst.d  $ft1, $t0, %d\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_I32_STORE8:
		i := i.(ast.Ins_I32Store8)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # i32.store8\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 出栈+保存
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    st.b  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_I32_STORE16:
		i := i.(ast.Ins_I32Store16)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # i32.store16\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 出栈+保存
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    st.h  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_I64_STORE8:
		i := i.(ast.Ins_I64Store8)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # i64.store8\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 出栈+保存
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.d  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    st.b  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_I64_STORE16:
		i := i.(ast.Ins_I64Store16)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # i64.store16\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 出栈+保存
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.d  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    st.h  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_I64_STORE32:
		i := i.(ast.Ins_I64Store32)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # i64.store32\n")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加上目标地址
		fmt.Fprintf(w, "    ld.w  $t1, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    add.d $t0, $t0, $t1\n")

		// 出栈+保存
		assert(i.Offset < 2048)
		fmt.Fprintf(w, "    ld.d  $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    st.w  $t1, $t0, %d\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_MEMORY_SIZE:
		sp0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8
		fmt.Fprintf(w, "    # memory.size\n")
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryPagesName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryPagesName)
		fmt.Fprintf(w, "    ld.w      $t0, $t0, 0\n")
		fmt.Fprintf(w, "    st.w      $t0, $fp, %d\n", sp0)
		fmt.Fprintln(w)

	case token.INS_MEMORY_GROW:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		labelSuffix := p.genNextId()

		labelElse := p.makeLabelId(kLabelPrefixName_else, "", labelSuffix)
		labelEnd := p.makeLabelId(kLabelPrefixName_end, "", labelSuffix)

		fmt.Fprintf(w, "    # memory.grow\n")
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryPagesName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryPagesName)

		fmt.Fprintf(w, "    pcalau12i $t1, %%pc_hi20(%s)\n", kMemoryMaxPagesName)
		fmt.Fprintf(w, "    addi.d    $t1, $t1, %%pc_lo12(%s)\n", kMemoryMaxPagesName)

		fmt.Fprintf(w, "    ld.w      $t2, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    ld.w      $t3, $t0, 0 # pageNum\n")
		fmt.Fprintf(w, "    ld.w      $t4, $t1, 0 # maxPageNum\n")
		fmt.Fprintf(w, "    add.w     $t2, $t2, $t3\n")

		fmt.Fprintf(w, "    blt       $t4, $t2, %s # jmp if t2 > maxPageNum \n", labelElse)
		fmt.Fprintf(w, "    st.w      $t2, $t0, 0 # update pageNum\n")
		fmt.Fprintf(w, "    st.w      $t3, $fp, %d\n", ret0)
		fmt.Fprintf(w, "    b         %s\n", labelEnd)

		p.gasFuncLabel(w, labelElse)
		fmt.Fprintf(w, "    addi.w $t0, $zero, -1\n")
		fmt.Fprintf(w, "    st.w   $t0, $fp, %d\n", ret0)

		p.gasFuncLabel(w, labelEnd)
		fmt.Fprintln(w)

	case token.INS_MEMORY_INIT:
		i := i.(ast.Ins_MemoryInit)

		len := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		off := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		dst := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # memory.init")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 字符串是一个地址, 不需要再加载了
		fmt.Fprintf(w, "    pcalau12i $t1, %%pc_hi20(%s%d)\n", kMemoryDataPtrPrefix, i.DataIdx)
		fmt.Fprintf(w, "    addi.d    $t1, $t1, %%pc_lo12(%s%d)\n", kMemoryDataPtrPrefix, i.DataIdx)

		// 加载函数地址
		fmt.Fprintf(w, "    pcalau12i $t2, %%pc_hi20(%s)\n", kRuntimeMemcpy)
		fmt.Fprintf(w, "    addi.d    $t2, $t2, %%pc_lo12(%s)\n", kRuntimeMemcpy)

		// 准备调用参数
		fmt.Fprintf(w, "    ld.w  $a0, $fp, %d\n", dst)
		fmt.Fprintf(w, "    add.d $a0, $a0, $t0\n")
		fmt.Fprintf(w, "    ld.w  $a1, $fp, %d\n", off)
		fmt.Fprintf(w, "    add.d $a1, $a1, $t1\n")
		fmt.Fprintf(w, "    ld.w  $a2, $fp, %d\n", len)

		// 调用函数
		fmt.Fprintf(w, "    jirl  $ra, $t2, 0 # %s\n", kRuntimeMemcpy)
		fmt.Fprintln(w)

	case token.INS_MEMORY_COPY:
		len := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		src := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		dst := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # memory.copy")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加载函数地址
		fmt.Fprintf(w, "    pcalau12i $t1, %%pc_hi20(%s)\n", kRuntimeMemmove)
		fmt.Fprintf(w, "    addi.d    $t1, $t1, %%pc_lo12(%s)\n", kRuntimeMemmove)

		// 准备调用参数
		fmt.Fprintf(w, "    ld.w  $a0, $fp, %d\n", dst)
		fmt.Fprintf(w, "    add.d $a0, $a0, $t0\n")
		fmt.Fprintf(w, "    ld.w  $a1, $fp, %d\n", src)
		fmt.Fprintf(w, "    add.d $a1, $a1, $t0\n")
		fmt.Fprintf(w, "    ld.w  $a2, $fp, %d\n", len)

		// 调用函数
		fmt.Fprintf(w, "    jirl  $ra, $t1, 0 # %s\n", kRuntimeMemmove)
		fmt.Fprintln(w)

	case token.INS_MEMORY_FILL:
		len := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		val := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		dst := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # memory.fill")

		// 内存的开始地址
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

		// 加载函数地址
		fmt.Fprintf(w, "    pcalau12i $t1, %%pc_hi20(%s)\n", kRuntimeMemset)
		fmt.Fprintf(w, "    addi.d    $t1, $t1, %%pc_lo12(%s)\n", kRuntimeMemset)

		// 准备调用参数
		fmt.Fprintf(w, "    ld.w  $a0, $fp, %d\n", dst)
		fmt.Fprintf(w, "    add.d $a0, $a0, $t0\n")
		fmt.Fprintf(w, "    ld.w  $a1, $fp, %d\n", val)
		fmt.Fprintf(w, "    ld.w  $a2, $fp, %d\n", len)

		// 调用函数
		fmt.Fprintf(w, "    jirl  $ra, $t1, 0 # %s\n", kRuntimeMemset)
		fmt.Fprintln(w)

	case token.INS_I32_CONST:
		i := i.(ast.Ins_I32Const)

		sp0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.const %d\n", i.X)
		if i.X >= -2048 && i.X <= 2047 {
			fmt.Fprintf(w, "    addi.d $t0, $zero, %d\n", i.X)
			fmt.Fprintf(w, "    st.d   $t0, $fp, %d\n", sp0)
			fmt.Fprintln(w)
		} else {
			hi20 := uint32(i.X) >> 12
			lo12 := uint32(i.X) & 0xFFF
			fmt.Fprintf(w, "    lu12i.w $t0, 0x%X\n", hi20)
			fmt.Fprintf(w, "    ori     $t0, $t0, 0x%X\n", lo12)
			fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", sp0)
			fmt.Fprintln(w)
		}

	case token.INS_I64_CONST:
		i := i.(ast.Ins_I64Const)

		sp0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.const %d\n", i.X)
		if i.X >= -2048 && i.X <= 2047 {
			fmt.Fprintf(w, "    addi.d $t0, $zero, %d\n", i.X)
			fmt.Fprintf(w, "    st.d   $t0, $fp, %d\n", sp0)
			fmt.Fprintln(w)
		} else if int64(int32(i.X)) == i.X {
			hi20 := uint32(int32(i.X)) >> 12
			lo12 := uint32(int32(i.X)) & 0xFFF
			fmt.Fprintf(w, "    lu12i.w $t0, 0x%X\n", hi20)
			fmt.Fprintf(w, "    ori     $t0, $t0, 0x%X\n", lo12)
			fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", sp0)
			fmt.Fprintln(w)
		} else {
			val := uint64(i.X)
			hi20 := uint32(val) >> 12 & 0xFFFFF
			lo12 := uint32(val) & 0xFFF
			mid20 := (val >> 32) & 0xFFFFF
			top12 := (val >> 52) & 0xFFF

			fmt.Fprintf(w, "    lu12i.w $t0, 0x%X\n", hi20)
			fmt.Fprintf(w, "    ori     $t0, $t0, 0x%X\n", lo12)
			fmt.Fprintf(w, "    lu32i.d $t0, 0x%X\n", mid20)
			fmt.Fprintf(w, "    lu52i.d $t0, $t0, 0x%X\n", top12)
			fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", sp0)
			fmt.Fprintln(w)
		}

	case token.INS_F32_CONST:
		i := i.(ast.Ins_F32Const)
		bits := math.Float32bits(i.X)

		sp0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.const %f (raw: 0x%08X)\n", i.X, bits)
		if x := int32(bits); x >= -2048 && x <= 2047 {
			fmt.Fprintf(w, "    addi.d $t0, $zero, %d\n", x)
			fmt.Fprintf(w, "    st.d   $t0, $fp, %d\n", sp0)
			fmt.Fprintln(w)
		} else {
			hi20 := bits >> 12
			lo12 := bits & 0xFFF
			fmt.Fprintf(w, "    lu12i.w $t0, 0x%X\n", hi20)
			fmt.Fprintf(w, "    ori     $t0, $t0, 0x%X\n", lo12)
			fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", sp0)
			fmt.Fprintln(w)
		}

	case token.INS_F64_CONST:
		i := i.(ast.Ins_F64Const)
		bits := math.Float64bits(i.X)

		sp0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.const %f (raw: 0x%016X)\n", i.X, bits)
		if x := int64(bits); x >= -2048 && x <= 2047 {
			fmt.Fprintf(w, "    addi.d $t0, $zero, %d\n", x)
			fmt.Fprintf(w, "    st.d   $t0, $fp, %d\n", sp0)
			fmt.Fprintln(w)
		} else if int64(int32(x)) == x {
			hi20 := bits >> 12
			lo12 := bits & 0xFFF
			fmt.Fprintf(w, "    lu12i.w $t0, 0x%X\n", hi20)
			fmt.Fprintf(w, "    ori     $t0, $t0, 0x%X\n", lo12)
			fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", sp0)
			fmt.Fprintln(w)
		} else {
			hi20 := (uint32(bits) >> 12) & 0xFFFFF
			lo12 := uint32(bits) & 0xFFF
			mid20 := (bits >> 32) & 0xFFFFF
			top12 := (bits >> 52) & 0xFFF
			fmt.Fprintf(w, "    lu12i.w $t0, 0x%X\n", hi20)
			fmt.Fprintf(w, "    ori     $t0, $t0, 0x%X\n", lo12)
			fmt.Fprintf(w, "    lu32i.d $t0, 0x%X\n", mid20)
			fmt.Fprintf(w, "    lu52i.d $t0, $t0, 0x%X\n", top12)
			fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", sp0)
			fmt.Fprintln(w)
		}

	case token.INS_I32_EQZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.eqz\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sltui   $t1, $t0, 1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_EQ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.eq\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    xor     $t0, $t1, $t0\n")
		fmt.Fprintf(w, "    sltui   $t1, $t0, 1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_NE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.ne\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    xor     $t0, $t1, $t0\n")
		fmt.Fprintf(w, "    sltu    $t1, $zero, $t0\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LT_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.lt_s\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    slt     $t1, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LT_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.lt_u\n")
		fmt.Fprintf(w, "    ld.wu   $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.wu   $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sltu    $t1, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_GT_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.gt_s\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    slt     $t1, $t1, $t0\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_GT_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.gt_u\n")
		fmt.Fprintf(w, "    ld.wu   $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.wu   $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sltu    $t1, $t1, $t0\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LE_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.le_s\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    slt     $t1, $t1, $t0\n")
		fmt.Fprintf(w, "    xori    $t1, $t1, 1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LE_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.le_u\n")
		fmt.Fprintf(w, "    ld.wu   $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.wu   $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sltu    $t1, $t1, $t0\n")
		fmt.Fprintf(w, "    xori    $t1, $t1, 1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_GE_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.ge_s\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    slt     $t1, $t0, $t1\n")
		fmt.Fprintf(w, "    xori    $t1, $t1, 1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_GE_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.ge_u\n")
		fmt.Fprintf(w, "    ld.wu   $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.wu   $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sltu    $t1, $t0, $t1\n")
		fmt.Fprintf(w, "    xori    $t1, $t1, 1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_EQZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.eqz\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sltui   $t1, $t0, 1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_EQ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.eq\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    xor     $t0, $t1, $t0\n")
		fmt.Fprintf(w, "    sltui   $t1, $t0, 1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_NE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.ne\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    xor     $t0, $t1, $t0\n")
		fmt.Fprintf(w, "    sltu    $t1, $zero, $t0\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LT_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.lt_s\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    slt     $t1, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LT_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.lt_u\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sltu    $t1, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_GT_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.gt_s\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    slt     $t1, $t1, $t0\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_GT_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.gt_u\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sltu    $t1, $t1, $t0\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LE_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.le_s\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    slt     $t1, $t1, $t0\n")
		fmt.Fprintf(w, "    xori    $t1, $t1, 1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LE_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.le_u\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sltu    $t1, $t1, $t0\n")
		fmt.Fprintf(w, "    xori    $t1, $t1, 1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_GE_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.ge_s\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    slt     $t1, $t0, $t1\n")
		fmt.Fprintf(w, "    xori    $t1, $t1, 1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_GE_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.ge_u\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sltu    $t1, $t0, $t1\n")
		fmt.Fprintf(w, "    xori    $t1, $t1, 1\n")
		fmt.Fprintf(w, "    st.w    $t1, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_EQ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f32.eq\n")
		fmt.Fprintf(w, "    fld.s      $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.s      $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcmp.ceq.s $fcc0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    movcf2gr   $t0, $fcc0\n")
		fmt.Fprintf(w, "    st.w       $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_NE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f32.ne\n")
		fmt.Fprintf(w, "    fld.s       $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.s       $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcmp.cune.s $fcc0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    movcf2gr    $t0, $fcc0\n")
		fmt.Fprintf(w, "    st.w        $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_LT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f32.lt\n")
		fmt.Fprintf(w, "    fld.s      $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.s      $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcmp.clt.s $fcc0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    movcf2gr   $t0, $fcc0\n")
		fmt.Fprintf(w, "    st.w       $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_GT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f32.gt\n")
		fmt.Fprintf(w, "    fld.s      $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.s      $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcmp.clt.s $fcc0, $ft1, $ft0\n")
		fmt.Fprintf(w, "    movcf2gr   $t0, $fcc0\n")
		fmt.Fprintf(w, "    st.w       $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_LE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f32.le\n")
		fmt.Fprintf(w, "    fld.s      $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.s      $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcmp.cle.s $fcc0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    movcf2gr   $t0, $fcc0\n")
		fmt.Fprintf(w, "    st.w       $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_GE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f32.ge\n")
		fmt.Fprintf(w, "    fld.s      $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.s      $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcmp.cle.s $fcc0, $ft1, $ft0\n")
		fmt.Fprintf(w, "    movcf2gr   $t0, $fcc0\n")
		fmt.Fprintf(w, "    st.w       $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_EQ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f64.eq\n")
		fmt.Fprintf(w, "    fld.d      $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.d      $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcmp.ceq.d $fcc0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    movcf2gr   $t0, $fcc0\n")
		fmt.Fprintf(w, "    st.w       $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_NE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f64.ne\n")
		fmt.Fprintf(w, "    fld.d       $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.d       $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcmp.cune.d $fcc0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    movcf2gr    $t0, $fcc0\n")
		fmt.Fprintf(w, "    st.w        $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_LT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f64.lt\n")
		fmt.Fprintf(w, "    fld.d      $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.d      $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcmp.clt.d $fcc0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    movcf2gr   $t0, $fcc0\n")
		fmt.Fprintf(w, "    st.w       $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_GT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f64.gt\n")
		fmt.Fprintf(w, "    fld.d      $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.d      $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcmp.clt.d $fcc0, $ft1, $ft0\n")
		fmt.Fprintf(w, "    movcf2gr   $t0, $fcc0\n")
		fmt.Fprintf(w, "    st.w       $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_LE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f64.le\n")
		fmt.Fprintf(w, "    fld.d      $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.d      $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcmp.cle.d $fcc0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    movcf2gr   $t0, $fcc0\n")
		fmt.Fprintf(w, "    st.w       $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_GE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f64.ge\n")
		fmt.Fprintf(w, "    fld.d      $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.d      $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcmp.cle.d $fcc0, $ft1, $ft0\n")
		fmt.Fprintf(w, "    movcf2gr   $t0, $fcc0\n")
		fmt.Fprintf(w, "    st.w       $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_CLZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.clz\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    clz.w   $t0, $t0\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_CTZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.ctz\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    ctz.w   $t0, $t0\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_POPCNT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.popcnt\n")
		fmt.Fprintf(w, "    ld.wu   $t0, $fp, %d\n", sp0)

		// 分治法统计
		fmt.Fprintf(w, "    # t0 = t0 - ((t0 >> 1) & 0x55555555)\n")
		fmt.Fprintf(w, "    srli.w  $t1, $t0, 1\n")
		fmt.Fprintf(w, "    lu12i.w $t2, 0x55555\n")
		fmt.Fprintf(w, "    ori     $t2, $t2, 0x555\n")
		fmt.Fprintf(w, "    and     $t1, $t1, $t2\n")
		fmt.Fprintf(w, "    sub.w   $t0, $t0, $t1\n")

		fmt.Fprintf(w, "    # t0 = (t0 & 0x33333333) + ((t0 >> 2) & 0x33333333)\n")
		fmt.Fprintf(w, "    srli.w  $t1, $t0, 2\n")
		fmt.Fprintf(w, "    lu12i.w $t2, 0x33333\n")
		fmt.Fprintf(w, "    ori     $t2, $t2, 0x333\n")
		fmt.Fprintf(w, "    and     $t1, $t1, $t2\n")
		fmt.Fprintf(w, "    and     $t0, $t0, $t2\n")
		fmt.Fprintf(w, "    add.w   $t0, $t0, $t1\n")

		fmt.Fprintf(w, "    # t0 = (t0 + (t0 >> 4)) & 0x0F0F0F0F\n")
		fmt.Fprintf(w, "    srli.w  $t1, $t0, 4\n")
		fmt.Fprintf(w, "    add.w   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    lu12i.w $t2, 0x0F0F0\n")
		fmt.Fprintf(w, "    ori     $t2, $t2, 0xF0F\n")
		fmt.Fprintf(w, "    and     $t0, $t0, $t2\n")

		fmt.Fprintf(w, "    # t0 = (t0 * 0x01010101) >> 24\n")
		fmt.Fprintf(w, "    lu12i.w $t1, 0x01010\n")
		fmt.Fprintf(w, "    ori     $t1, $t1, 0x101\n")
		fmt.Fprintf(w, "    mul.w   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    srli.w  $t0, $t0, 24\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_ADD:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.add\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.w   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_SUB:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.sub\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sub.w   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_MUL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.mul\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    mul.w   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_DIV_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.div_s\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    div.w   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_DIV_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.div_u\n")
		fmt.Fprintf(w, "    ld.wu   $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.wu   $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    div.wu  $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_REM_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.rem_s\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    mod.w   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_REM_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.rem_u\n")
		fmt.Fprintf(w, "    ld.wu   $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.wu   $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    mod.wu  $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_AND:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.and\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    and     $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_OR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.or\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    or      $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_XOR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.xor\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    xor     $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_SHL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.shl\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sll.w   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_SHR_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.shr_s\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sra.w   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_SHR_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.shr_u\n")
		fmt.Fprintf(w, "    ld.wu   $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.wu   $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    srl.w   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_ROTL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.rotl\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sub.w   $t1, $zero, $t1\n")
		fmt.Fprintf(w, "    rotr.w  $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_ROTR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.rotr\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.w    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    rotr.w  $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_CLZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.clz\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    clz.d   $t0, $t0\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_CTZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.ctz\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    ctz.d   $t0, $t0\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_POPCNT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.popcnt\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp0)

		fmt.Fprintf(w, "    # x = x - ((x >> 1) & 0x5555555555555555)\n")
		fmt.Fprintf(w, "    srli.d  $t1, $t0, 1\n")
		fmt.Fprintf(w, "    lu12i.w $t2, 0x55555\n")
		fmt.Fprintf(w, "    ori     $t2, $t2, 0x555\n")
		fmt.Fprintf(w, "    lu32i.d $t2, 0x55555\n")
		fmt.Fprintf(w, "    lu52i.d $t2, $t2, 0x555\n")
		fmt.Fprintf(w, "    and     $t1, $t1, $t2\n")
		fmt.Fprintf(w, "    sub.d   $t0, $t0, $t1\n")

		fmt.Fprintf(w, "    # x = (x & 0x3333333333333333) + ((x >> 2) & 0x3333333333333333)\n")
		fmt.Fprintf(w, "    lu12i.w $t2, 0x33333\n")
		fmt.Fprintf(w, "    ori     $t2, $t2, 0x333\n")
		fmt.Fprintf(w, "    lu32i.d $t2, 0x33333\n")
		fmt.Fprintf(w, "    lu52i.d $t2, $t2, 0x333\n")
		fmt.Fprintf(w, "    srli.d  $t1, $t0, 2\n")
		fmt.Fprintf(w, "    and     $t0, $t0, $t2\n")
		fmt.Fprintf(w, "    and     $t1, $t1, $t2\n")
		fmt.Fprintf(w, "    add.d   $t0, $t0, $t1\n")

		fmt.Fprintf(w, "    # x = (x + (x >> 4)) & 0x0F0F0F0F0F0F0F0F\n")
		fmt.Fprintf(w, "    srli.d  $t1, $t0, 4\n")
		fmt.Fprintf(w, "    add.d   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    lu12i.w $t2, 0x0f0f0\n")
		fmt.Fprintf(w, "    ori     $t2, $t2, 0xf0f\n")
		fmt.Fprintf(w, "    lu32i.d $t2, 0x0f0f0\n")
		fmt.Fprintf(w, "    lu52i.d $t2, $t2, 0xf0f\n")
		fmt.Fprintf(w, "    and     $t0, $t0, $t2\n")

		fmt.Fprintf(w, "    # result = (x * 0x0101010101010101) >> 56\n")
		fmt.Fprintf(w, "    lu12i.w $t1, 0x01010\n")
		fmt.Fprintf(w, "    ori     $t1, $t1, 0x101\n")
		fmt.Fprintf(w, "    lu32i.d $t1, 0x01010\n")
		fmt.Fprintf(w, "    lu52i.d $t1, $t1, 0x101\n")
		fmt.Fprintf(w, "    mul.d   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    srli.d  $t0, $t0, 56\n")

		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_ADD:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.add\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    add.d   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_SUB:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.sub\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sub.d   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_MUL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.mul\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    mul.d   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_DIV_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.div_s\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    div.d   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_DIV_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.div_u\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    div.du  $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_REM_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.rem_s\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    mod.d   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_REM_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.rem_u\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    mod.wu  $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_AND:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.and\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    and     $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_OR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.or\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    or      $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_XOR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.xor\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    xor     $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_SHL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.shl\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sll.d   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_SHR_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.shr_s\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sra.d   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_SHR_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.shr_u\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    srl.d   $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_ROTL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.rotl\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    sub.d   $t1, $zero, $t1\n")
		fmt.Fprintf(w, "    rotr.d  $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_ROTR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.rotr\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    ld.d    $t1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    rotr.d  $t0, $t0, $t1\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_ABS:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.abs\n")
		fmt.Fprintf(w, "    fld.s   $ft0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fabs.s  $ft0, $ft0\n")
		fmt.Fprintf(w, "    fst.s   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_NEG:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.neg\n")
		fmt.Fprintf(w, "    fld.s   $ft0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fneg.s  $ft0, $ft0\n")
		fmt.Fprintf(w, "    fst.s   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_CEIL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.ceil\n")
		fmt.Fprintf(w, "    fld.s     $ft0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    frintrp.s $ft0, $ft0\n")
		fmt.Fprintf(w, "    fst.s     $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_FLOOR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.floor\n")
		fmt.Fprintf(w, "    fld.s     $ft0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    frintrm.s $ft0, $ft0\n")
		fmt.Fprintf(w, "    fst.s     $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_TRUNC:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.trunc\n")
		fmt.Fprintf(w, "    fld.s     $ft0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    frintrz.s $ft0, $ft0\n")
		fmt.Fprintf(w, "    fst.s     $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_NEAREST:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.nearest\n")
		fmt.Fprintf(w, "    fld.s      $ft0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    frintrne.s $ft0, $ft0\n")
		fmt.Fprintf(w, "    fst.s      $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_SQRT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.sqrt\n")
		fmt.Fprintf(w, "    fld.s   $ft0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fsqrt.s $ft0, $ft0\n")
		fmt.Fprintf(w, "    fst.s   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_ADD:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.add\n")
		fmt.Fprintf(w, "    fld.s   $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.s   $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fadd.s  $ft0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    fst.s   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_SUB:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.sub\n")
		fmt.Fprintf(w, "    fld.s   $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.s   $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fsub.s  $ft0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    fst.s   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_MUL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.mul\n")
		fmt.Fprintf(w, "    fld.s   $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.s   $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fmul.s  $ft0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    fst.s   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_DIV:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.div\n")
		fmt.Fprintf(w, "    fld.s   $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.s   $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fdiv.s  $ft0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    fst.s   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_MIN:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.min\n")
		fmt.Fprintf(w, "    fld.s   $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.s   $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fmin.s  $ft0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    fst.s   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_MAX:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.max\n")
		fmt.Fprintf(w, "    fld.s   $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.s   $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fmax.s  $ft0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    fst.s   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_COPYSIGN:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.copysign\n")
		fmt.Fprintf(w, "    fld.s       $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.s       $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcopysign.s $ft0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    fst.s       $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_ABS:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.abs\n")
		fmt.Fprintf(w, "    fld.d   $ft0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fabs.d  $ft0, $ft0\n")
		fmt.Fprintf(w, "    fst.d   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_NEG:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.neg\n")
		fmt.Fprintf(w, "    fld.d   $ft0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fneg.d  $ft0, $ft0\n")
		fmt.Fprintf(w, "    fst.d   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_CEIL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.ceil\n")
		fmt.Fprintf(w, "    fld.d     $ft0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    frintrp.d $ft0, $ft0\n")
		fmt.Fprintf(w, "    fst.d     $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_FLOOR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.floor\n")
		fmt.Fprintf(w, "    fld.d     $ft0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    frintrm.d $ft0, $ft0\n")
		fmt.Fprintf(w, "    fst.d     $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_TRUNC:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.trunc\n")
		fmt.Fprintf(w, "    fld.d     $ft0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    frintrz.d $ft0, $ft0\n")
		fmt.Fprintf(w, "    fst.d     $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_NEAREST:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.nearest\n")
		fmt.Fprintf(w, "    fld.d      $ft0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    frintrne.d $ft0, $ft0\n")
		fmt.Fprintf(w, "    fst.d      $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_SQRT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.sqrt\n")
		fmt.Fprintf(w, "    fld.d   $ft0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fsqrt.d $ft0, $ft0\n")
		fmt.Fprintf(w, "    fst.d   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_ADD:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.add\n")
		fmt.Fprintf(w, "    fld.d   $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.d   $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fadd.d  $ft0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    fst.d   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_SUB:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.sub\n")
		fmt.Fprintf(w, "    fld.d   $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.d   $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fsub.d  $ft0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    fst.d   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_MUL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.mul\n")
		fmt.Fprintf(w, "    fld.d   $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.d   $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fmul.d  $ft0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    fst.d   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_DIV:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.div\n")
		fmt.Fprintf(w, "    fld.d   $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.d   $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fdiv.d  $ft0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    fst.d   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_MIN:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.min\n")
		fmt.Fprintf(w, "    fld.d   $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.d   $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fmin.d  $ft0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    fst.d   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_MAX:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.max\n")
		fmt.Fprintf(w, "    fld.d   $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.d   $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fmax.d  $ft0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    fst.d   $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_COPYSIGN:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.copysign\n")
		fmt.Fprintf(w, "    fld.d       $ft0, $fp, %d\n", sp1)
		fmt.Fprintf(w, "    fld.d       $ft1, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcopysign.d $ft0, $ft0, $ft1\n")
		fmt.Fprintf(w, "    fst.d       $ft0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_WRAP_I64:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.wrap_i64\n")
		fmt.Fprintf(w, "    ld.d    $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    slli.w  $t0, $t0, 0\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_TRUNC_F32_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.trunc_f32_s\n")
		fmt.Fprintf(w, "    fld.s    $fa0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    ftintrz.w.s $fa0, $fa0\n") // 浮点转32位整型，存放在浮点寄存器
		fmt.Fprintf(w, "    movfr2gr.w $t0, $fa0\n")   // 移回通用寄存器
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_TRUNC_F32_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.trunc_f32_u\n")
		fmt.Fprintf(w, "    fld.s    $fa0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    ftintrz.l.s $fa0, $fa0\n")
		fmt.Fprintf(w, "    movfr2gr.d $t0, $fa0\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_TRUNC_F64_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.trunc_f64_s\n")
		fmt.Fprintf(w, "    fld.d    $fa0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    ftintrz.w.d $fa0, $fa0\n")
		fmt.Fprintf(w, "    movfr2gr.w $t0, $fa0\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_TRUNC_F64_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.trunc_f64_u\n")
		fmt.Fprintf(w, "    fld.d    $fa0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    ftintrz.l.d $fa0, $fa0\n")
		fmt.Fprintf(w, "    movfr2gr.d $t0, $fa0\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_EXTEND_I32_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.extend_i32_s\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp0) // ld.w 会自动符号扩展到64位
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_EXTEND_I32_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.extend_i32_u\n")
		fmt.Fprintf(w, "    ld.wu   $t0, $fp, %d\n", sp0) // ld.wu 零扩展到64位
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_TRUNC_F32_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.trunc_f32_s\n")
		fmt.Fprintf(w, "    fld.s    $fa0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    ftintrz.l.s $fa0, $fa0\n")
		fmt.Fprintf(w, "    movfr2gr.d $t0, $fa0\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_TRUNC_F32_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.trunc_f32_u\n")
		fmt.Fprintf(w, "    fld.s    $fa0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    ftintrz.l.s $fa0, $fa0\n")
		fmt.Fprintf(w, "    movfr2gr.d $t0, $fa0\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_TRUNC_F64_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.trunc_f64_s\n")
		fmt.Fprintf(w, "    fld.d    $fa0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    ftintrz.l.d $fa0, $fa0\n")
		fmt.Fprintf(w, "    movfr2gr.d $t0, $fa0\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_TRUNC_F64_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.trunc_f64_u\n")
		fmt.Fprintf(w, "    fld.d    $fa0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    ftintrz.l.d $fa0, $fa0\n")
		fmt.Fprintf(w, "    movfr2gr.d $t0, $fa0\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_CONVERT_I32_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.convert_i32_s\n")
		fmt.Fprintf(w, "    ld.w     $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    movgr2fr.w $fa0, $t0\n")
		fmt.Fprintf(w, "    ffint.s.w  $fa0, $fa0\n")
		fmt.Fprintf(w, "    fst.s    $fa0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_CONVERT_I32_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.convert_i32_u\n")
		fmt.Fprintf(w, "    ld.wu    $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    movgr2fr.d $fa0, $t0\n")
		fmt.Fprintf(w, "    ffint.s.l  $fa0, $fa0\n")
		fmt.Fprintf(w, "    fst.s    $fa0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_CONVERT_I64_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.convert_i64_s\n")
		fmt.Fprintf(w, "    ld.d     $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    movgr2fr.d $fa0, $t0\n")
		fmt.Fprintf(w, "    ffint.s.l  $fa0, $fa0\n") // 64位有符号整数 -> f32
		fmt.Fprintf(w, "    fst.s    $fa0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_CONVERT_I64_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.convert_i64_u\n")
		fmt.Fprintf(w, "    ld.d     $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    movgr2fr.d $fa0, $t0\n")
		fmt.Fprintf(w, "    ffint.s.l  $fa0, $fa0\n") // 临时替代方案，注意：这在处理大正数时会出错
		fmt.Fprintf(w, "    fst.s    $fa0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_DEMOTE_F64:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.demote_f64\n")
		fmt.Fprintf(w, "    fld.d    $fa0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcvt.s.d $fa0, $fa0\n")
		fmt.Fprintf(w, "    fst.s    $fa0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_CONVERT_I32_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.convert_i32_s\n")
		fmt.Fprintf(w, "    ld.w     $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    movgr2fr.w $fa0, $t0\n")
		fmt.Fprintf(w, "    ffint.d.w  $fa0, $fa0\n")
		fmt.Fprintf(w, "    fst.d    $fa0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_CONVERT_I32_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.convert_i32_u\n")
		fmt.Fprintf(w, "    ld.wu    $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    movgr2fr.d $fa0, $t0\n")
		fmt.Fprintf(w, "    ffint.d.l  $fa0, $fa0\n")
		fmt.Fprintf(w, "    fst.d    $fa0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_CONVERT_I64_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.convert_i64_s\n")
		fmt.Fprintf(w, "    ld.d     $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    movgr2fr.d $fa0, $t0\n")
		fmt.Fprintf(w, "    ffint.d.l  $fa0, $fa0\n")
		fmt.Fprintf(w, "    fst.d    $fa0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_CONVERT_I64_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.convert_i64_u\n")
		fmt.Fprintf(w, "    ld.d      $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    movgr2fr.d $fa0, $t0\n")
		fmt.Fprintf(w, "    ffint.d.l $fa0, $fa0\n")
		fmt.Fprintf(w, "    fst.d      $fa0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_PROMOTE_F32:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.promote_f32\n")
		fmt.Fprintf(w, "    fld.s    $fa0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    fcvt.d.s $fa0, $fa0\n")
		fmt.Fprintf(w, "    fst.d    $fa0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_REINTERPRET_F32:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.reinterpret_f32\n")
		fmt.Fprintf(w, "    fld.s    $fa0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    movfr2gr.s $t0, $fa0\n")
		fmt.Fprintf(w, "    st.w    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_REINTERPRET_F64:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.reinterpret_f64\n")
		fmt.Fprintf(w, "    fld.d    $fa0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    movfr2gr.d $t0, $fa0\n")
		fmt.Fprintf(w, "    st.d    $t0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_REINTERPRET_I32:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.reinterpret_i32\n")
		fmt.Fprintf(w, "    ld.w    $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    movgr2fr.w $fa0, $t0\n")
		fmt.Fprintf(w, "    fst.s    $fa0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_REINTERPRET_I64:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.reinterpret_i64\n")
		fmt.Fprintf(w, "    ld.d     $t0, $fp, %d\n", sp0)
		fmt.Fprintf(w, "    movgr2fr.d $fa0, $t0\n")
		fmt.Fprintf(w, "    fst.d    $fa0, $fp, %d\n", ret0)
		fmt.Fprintln(w)

	default:
		panic(fmt.Sprintf("unreachable: %T", i))
	}
	return nil
}
