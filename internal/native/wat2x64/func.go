// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"wa-lang.org/wa/internal/native/abi"
	nativeast "wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/ast/astutil"
	"wa-lang.org/wa/internal/native/x64"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

const (
	kFuncNamePrefix    = ".F."
	kFuncRetNamePrefix = ".F.ret."
)

const (
	kLabelName_return = ".L.return"

	kLabelNamePreifx_blockBegin = ".L.block.begin."
	kLabelNamePreifx_blockEnd   = ".L.block.end."

	kLabelNamePreifx_ifBegin = ".L.if.begin."
	kLabelNamePreifx_ifBody  = ".L.if.body."
	kLabelNamePreifx_ifElse  = ".L.if.else."
	kLabelNamePreifx_ifEnd   = ".L.if.end."

	kLabelNamePreifx_loopBegin = ".L.loop.begin."
	kLabelNamePreifx_loopEnd   = ".L.loop.end."

	kLabelNamePreifx_next = ".L.next."
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

func (p *wat2X64Worker) buildFuncs(w io.Writer) error {
	if len(p.m.Funcs) == 0 {
		return nil
	}

	for _, f := range p.m.Funcs {
		p.localNames = nil
		p.localTypes = nil
		p.scopeLabels = nil
		p.scopeStackBases = nil
		p.scopeResults = nil

		p.gasComment(w, "func "+f.Name)
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

func (p *wat2X64Worker) buildFunc_body(w io.Writer, fn *ast.Func) error {
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

	// 至少要有一个指令
	if len(fn.Body.Insts) == 0 {
		fn.Body.Insts = []ast.Instruction{
			ast.Ins_Return{OpToken: ast.OpToken(token.INS_RETURN)},
		}
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
	if err := astutil.BuildFuncFrame(abi.X64, fnNative); err != nil {
		return err
	}

	// WASM 栈的开始位置
	// WASM 栈底第一个元素相对于 rbp 的偏移位置, 每个元素 8 字节
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

		fmt.Fprintf(&bufHeader, "    push rbp\n")
		fmt.Fprintf(&bufHeader, "    mov  rbp, rsp\n")
		fmt.Fprintf(&bufHeader, "    sub  rsp, %d\n", frameSize)
		fmt.Fprintln(&bufHeader)
	}

	// 走栈返回
	if len(fnNative.Type.Return) > 1 && fnNative.Type.Return[1].Reg == 0 {
		p.gasCommentInFunc(&bufHeader, "将返回地址备份到栈")
		fmt.Fprintf(&bufHeader, "    mov [rbp%+d], rcx # return address\n", 2*8)
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
				fmt.Fprintf(&bufHeader, "    mov [rbp%+d], %v # save arg %s\n",
					arg.RBPOff, x64.RegString(arg.Reg),
					arg.Name,
				)
			case token.I64:
				fmt.Fprintf(&bufHeader, "    mov [rbp%+d], %v # save arg %s\n",
					arg.RBPOff, x64.RegString(arg.Reg),
					arg.Name,
				)
			case token.F32:
				fmt.Fprintf(&bufHeader, "    movss qword ptr [rbp%+d], %v # save arg %s\n",
					arg.RBPOff, x64.RegString(arg.Reg),
					arg.Name,
				)
			case token.F64:
				fmt.Fprintf(&bufHeader, "    movsd qword ptr [rbp%+d], %v # save arg %s\n",
					arg.RBPOff, x64.RegString(arg.Reg),
					arg.Name,
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
		for _, ret := range fnNative.Type.Return {
			fmt.Fprintf(&bufHeader,
				"    mov dword ptr [rbp%+d], 0 # ret %s = 0\n",
				ret.RBPOff, ret.Name,
			)
		}
	} else {
		p.gasCommentInFunc(&bufHeader, "没有返回值变量需要初始化为0")
		fmt.Fprintln(&bufHeader)
	}

	// 局部变量初始化为 0
	if len(fnNative.Body.Locals) > 0 {
		p.gasCommentInFunc(&bufHeader, "将局部变量初始化为0")
		for _, local := range fnNative.Body.Locals {
			fmt.Fprintf(&bufHeader, "    mov dword ptr [rbp%+d], 0 # local %s = 0\n",
				local.RBPOff, local.Name,
			)
		}
		fmt.Fprintln(&bufHeader)
	} else {
		p.gasCommentInFunc(&bufHeader, "没有局部变量需要初始化为0")
		fmt.Fprintln(&bufHeader)
	}

	// 根据ABI处理返回值
	// 需要将栈顶位置的结果转化为本地ABI规范的返回值
	{
		// 返回代码位置
		fmt.Fprintln(&bufIns)
		p.gasCommentInFunc(&bufIns, "根据ABI处理返回值")
		p.gasFuncLabel(&bufIns, kLabelName_return)

		// 如果走内存, 返回地址
		if len(fn.Type.Results) > 1 && fnNative.Type.Return[1].Reg == 0 {
			fmt.Fprintf(&bufIns, "    mov rax, [rbp%+d] # ret address\n",
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
					fmt.Fprintf(&bufIns, "    mov %v, [rbp%+d] # ret %s\n",
						x64.RegString(ret.Reg),
						fnNative.Type.Return[i].RBPOff,
						fnNative.Type.Return[i].Name,
					)
				case token.I64:
					fmt.Fprintf(&bufIns, "    mov %v, [rbp%+d] # ret %s\n",
						x64.RegString(ret.Reg),
						fnNative.Type.Return[i].RBPOff,
						fnNative.Type.Return[i].Name,
					)
				case token.F32:
					fmt.Fprintf(&bufIns, "    movss %v, qword ptr [rbp%+d] # ret %s\n",
						x64.RegString(ret.Reg),
						fnNative.Type.Return[i].RBPOff,
						fnNative.Type.Return[i].Name,
					)
				case token.F64:
					fmt.Fprintf(&bufIns, "    movsd %v, qword ptr [rbp%+d] # ret.%s\n",
						x64.RegString(ret.Reg),
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
		fmt.Fprintln(&bufIns, "    mov rsp, rbp")
		fmt.Fprintln(&bufIns, "    pop rbp")
		fmt.Fprintln(&bufIns, "    ret")
		fmt.Fprintln(&bufIns)
	}

	// 头部赋值到 w
	io.Copy(w, &bufHeader)
	// 指令复制到 w
	io.Copy(w, &bufIns)

	// 有些函数最后的位置不是 return, 需要手动清理栈(验证栈正确性)
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

func (p *wat2X64Worker) buildFunc_ins(
	w io.Writer, fnNative *nativeast.Func,
	fn *ast.Func, stk *valueTypeStack, i ast.Instruction,
) error {
	stk.NextInstruction(i)

	p.Tracef("buildFunc_ins: %s begin: %v\n", i.Token(), stk.String())
	defer func() { p.Tracef("buildFunc_ins: %s end: %v\n", i.Token(), stk.String()) }()

	switch tok := i.Token(); tok {
	case token.INS_UNREACHABLE:
		fmt.Fprintf(w, "    call %s # unreachable\n", kRuntimePanic)
	case token.INS_NOP:
		fmt.Fprintf(w, "    nop\n")

	case token.INS_BLOCK:
		i := i.(ast.Ins_Block)

		stkBase := stk.Len()
		defer func() { assert(stk.Len() == stkBase+len(i.Results)) }()

		p.enterLabelScope(stkBase, i.Label, i.Results)
		defer p.leaveLabelScope()

		if i.Label != "" {
			p.gasFuncLabel(w, kLabelNamePreifx_blockBegin+i.Label)
		}

		for _, ins := range i.List {
			if err := p.buildFunc_ins(w, fnNative, fn, stk, ins); err != nil {
				return err
			}
		}

		if i.Label != "" {
			p.gasFuncLabel(w, kLabelNamePreifx_next+i.Label)
			p.gasFuncLabel(w, kLabelNamePreifx_blockEnd+i.Label)
		}

	case token.INS_LOOP:
		i := i.(ast.Ins_Loop)

		stkBase := stk.Len()
		defer func() { assert(stk.Len() == stkBase+len(i.Results)) }()

		p.enterLabelScope(stkBase, i.Label, i.Results)
		defer p.leaveLabelScope()

		if i.Label != "" {
			p.gasFuncLabel(w, kLabelNamePreifx_loopBegin+i.Label)
			p.gasFuncLabel(w, kLabelNamePreifx_next+i.Label)
		}
		for _, ins := range i.List {
			if err := p.buildFunc_ins(w, fnNative, fn, stk, ins); err != nil {
				return err
			}
		}
		if i.Label != "" {
			p.gasFuncLabel(w, kLabelNamePreifx_loopEnd+i.Label)
		}

	case token.INS_IF:
		i := i.(ast.Ins_If)

		if i.Label != "" {
			p.gasFuncLabel(w, kLabelNamePreifx_ifBegin+i.Label)
		}

		// 龙芯没有 pop 指令，需要2个指令才能实现
		// 因此通过固定的偏移量，只需要一个指令

		sp0 := stk.Pop(token.I32)

		fmt.Fprintf(w, "    mov eax, [rbp%+d]\n", p.fnWasmR0Base-sp0*8)
		fmt.Fprintf(w, "    test eax, eax\n")
		if len(i.Else) > 0 {
			fmt.Fprintf(w, "    jne %s%s # if eax != 0 { jmp body }\n", kLabelNamePreifx_ifBody, i.Label)
		} else {
			fmt.Fprintf(w, "    je %s%s # if eax != 0 { jmp end }\n", kLabelNamePreifx_ifEnd, i.Label)
		}

		stkBase := stk.Len()
		defer func() { assert(stk.Len() == stkBase+len(i.Results)) }()

		p.enterLabelScope(stkBase, i.Label, i.Results)
		defer p.leaveLabelScope()

		if i.Label != "" {
			p.gasFuncLabel(w, kLabelNamePreifx_ifBody+i.Label)
		}
		for _, ins := range i.Body {
			if err := p.buildFunc_ins(w, fnNative, fn, stk, ins); err != nil {
				return err
			}
		}

		if len(i.Else) > 0 {
			p.Tracef("buildFunc_ins: %s begin: %v\n", token.INS_ELSE, stk.String())
			defer func() { p.Tracef("buildFunc_ins: %s end: %v\n", token.INS_ELSE, stk.String()) }()

			if i.Label != "" {
				p.gasFuncLabel(w, kLabelNamePreifx_ifElse+i.Label)
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
			p.gasFuncLabel(w, kLabelNamePreifx_next+i.Label)
			p.gasFuncLabel(w, kLabelNamePreifx_ifEnd+i.Label)
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
				fmt.Fprintf(w, "    # copy jmp %s result\n", labelName)
				for i := len(destScopeResults) - 1; i >= 0; i-- {
					xType := destScopeResults[i]
					reti := stk.Pop(xType)
					switch xType {
					case token.I32:
						fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-reti*8)
						fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
					case token.I64:
						fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", p.fnWasmR0Base-reti*8)
						fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
					case token.F32:
						fmt.Fprintf(w, "    movss xmm0, dword ptr [rbp%+d]\n", p.fnWasmR0Base-reti*8)
						fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm0\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
					case token.F64:
						fmt.Fprintf(w, "    movsd xmm0, qword ptr [rbp%+d]\n", p.fnWasmR0Base-reti*8)
						fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm0\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
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

		fmt.Fprintf(w, "    jmp %s%s\n", kLabelNamePreifx_next, labelName)

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
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8)
		fmt.Fprintf(w, "    test eax, eax\n")
		fmt.Fprintf(w, "    jne  %s\n", kLabelNamePreifx_next+labelName)

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
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8)
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
						fmt.Fprintf(w, "    # copy br_table %s result\n", labelName)
						for i := 0; i < len(destScopeResults); i++ {
							xType := destScopeResults[i]
							reti := retIdxList[i]
							switch xType {
							case token.I32:
								fmt.Fprintf(w, "    mov edx, dword ptr [rbp%+d]\n", p.fnWasmR0Base-reti*8)
								fmt.Fprintf(w, "    mov dword ptr [rbp%+d], edx\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
							case token.I64:
								fmt.Fprintf(w, "    mov rdx, qword ptr [rbp%+d]\n", p.fnWasmR0Base-reti*8)
								fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rdx\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
							case token.F32:
								fmt.Fprintf(w, "    movss xmm0, dword ptr [rbp%+d]\n", p.fnWasmR0Base-reti*8)
								fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm0\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
							case token.F64:
								fmt.Fprintf(w, "    movsd xmm0, qword ptr [rbp%+d]\n", p.fnWasmR0Base-reti*8)
								fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm0\n", p.fnWasmR0Base-(destScopeStackBase+i)*8)
							default:
								unreachable()
							}
						}
					}
				}

				if k < len(i.XList)-1 {
					fmt.Fprintf(w, "    # br_table case %d\n", k)
					fmt.Fprintf(w, "    cmp eax, %d\n", k)
					fmt.Fprintf(w, "    jne %s\n", kLabelNamePreifx_next+labelName)
				} else {
					assert(labelName == defaultLabelName)
					fmt.Fprintf(w, "    # br_table default\n")
					fmt.Fprintf(w, "    jmp %s\n", kLabelNamePreifx_next+defaultLabelName)
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
			fmt.Fprintf(w, "    jmp %s\n", kLabelName_return)
		case 1:
			sp0 := stk.Pop(fn.Type.Results[0])
			switch fn.Type.Results[0] {
			case token.I32:
				fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8)
				fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", fnNative.Type.Return[0].RBPOff)
				fmt.Fprintf(w, "    jmp %s\n", kLabelName_return)
			case token.I64:
				fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8)
				fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", fnNative.Type.Return[0].RBPOff)
				fmt.Fprintf(w, "    jmp %s\n", kLabelName_return)
			case token.F32:
				fmt.Fprintf(w, "    movss xmm0, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8)
				fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm0\n", fnNative.Type.Return[0].RBPOff)
				fmt.Fprintf(w, "    jmp   %s\n", kLabelName_return)
			case token.F64:
				fmt.Fprintf(w, "    movsd xmm0, qword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8)
				fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm0\n", fnNative.Type.Return[0].RBPOff)
				fmt.Fprintf(w, "    jmp   %s\n", kLabelName_return)
			default:
				unreachable()
			}
		default:
			for i := len(fn.Type.Results) - 1; i >= 0; i-- {
				xType := fn.Type.Results[i]
				spi := stk.Pop(xType)
				switch xType {
				case token.I32:
					fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-spi*8)
					fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", fnNative.Type.Return[i].RBPOff)
				case token.I64:
					fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", p.fnWasmR0Base-spi*8)
					fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", fnNative.Type.Return[i].RBPOff)
				case token.F32:
					fmt.Fprintf(w, "    movss xmm0, dword ptr [rbp%+d]\n", p.fnWasmR0Base-spi*8)
					fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm0\n", fnNative.Type.Return[i].RBPOff)
				case token.F64:
					fmt.Fprintf(w, "    movsd xmm0, qword ptr [rbp%+d]\n", p.fnWasmR0Base-spi*8)
					fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm0\n", fnNative.Type.Return[i].RBPOff)
				default:
					unreachable()
				}
			}
			fmt.Fprintf(w, "    mov    rax, [rbp%+d] # return address\n", 2*8)
			fmt.Fprintf(w, "    jmp    %s\n", kLabelName_return)
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
		if err := astutil.BuildFuncFrame(abi.X64, fnCallNative); err != nil {
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
		if len(fnCallNative.Type.Return) > 1 && fnCallNative.Type.Return[1].Reg == 0 {
			assert(p.osName == abi.WINDOWS.String())
			fmt.Fprintf(w, "    lea rcx, [rsp%+d] # return address\n", len(fnCallType.Params)*8)
		}

		// 准备调用参数
		for k, x := range fnCallType.Params {
			switch x.Type {
			case token.I32:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    mov %s, dword ptr [rbp%+d] # arg %d\n",
						x64.RegString(arg.Reg),
						p.fnWasmR0Base-argList[k]*8-8,
						k,
					)
				} else {
					fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n",
						p.fnWasmR0Base-argList[k]*8-8,
					)
					fmt.Fprintf(w, "    mov dword ptr [rsp%+d], eax\n",
						arg.RSPOff,
					)
				}
			case token.I64:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    mov %s, qword ptr [rbp%+d] # arg %d\n",
						x64.RegString(arg.Reg),
						p.fnWasmR0Base-argList[k]*8-8,
						k,
					)
				} else {
					fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n",
						p.fnWasmR0Base-argList[k]*8-8,
					)
					fmt.Fprintf(w, "    mov qword ptr [rsp%+d], rax\n",
						arg.RSPOff,
					)
				}
			case token.F32:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    movss %s, qword ptr [rbp%+d] # arg %d\n",
						x64.RegString(arg.Reg),
						p.fnWasmR0Base-argList[k]*8-8,
						k,
					)
				} else {
					fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n",
						p.fnWasmR0Base-argList[k]*8-8,
					)
					fmt.Fprintf(w, "    movss dword ptr [rsp%+d], xmm4\n",
						arg.RSPOff,
					)
				}
			case token.F64:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    movsd %s, qword ptr [rbp%+d] # arg %d\n",
						x64.RegString(arg.Reg),
						p.fnWasmR0Base-argList[k]*8-8,
						k,
					)
				} else {
					fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n",
						p.fnWasmR0Base-argList[k]*8-8,
					)
					fmt.Fprintf(w, "    movsd qword ptr [rsp%+d], xmm4\n",
						arg.RSPOff,
					)
				}
			default:
				unreachable()
			}
		}

		if fnCallIdx <= p.importFuncCount {
			fmt.Fprintf(w, "    call %s\n", kImportNamePrefix+i.X)
		} else {
			fmt.Fprintf(w, "    call %s\n", kFuncNamePrefix+i.X)
		}

		// 提取返回值
		if len(fnCallNative.Type.Return) > 1 && fnCallNative.Type.Return[1].Reg == 0 {
			// 走栈返回结果, 地址通过 rcx 传入, 通过 rax 返回
			// 根据 ABI 规范获取返回值, 再保存到 WASM 栈中
			for k, retType := range fnCallType.Results {
				reti := stk.Push(retType)
				switch retType {
				case token.I32:
					fmt.Fprintf(w, "    mov r10, dword ptr [rax%+d]\n", k*8)
					fmt.Fprintf(w, "    mov dword ptr [rbp%+d], r10\n", p.fnWasmR0Base-reti*8-8)
				case token.I64:
					fmt.Fprintf(w, "    mov r10, qword ptr [rax%+d]\n", k*8)
					fmt.Fprintf(w, "    mov qword ptr [rbp%+d], r10\n", p.fnWasmR0Base-reti*8-8)
				case token.F32:
					fmt.Fprintf(w, "    movss xmm4, dword ptr [rax%+d]\n", k*8)
					fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", p.fnWasmR0Base-reti*8-8)
				case token.F64:
					fmt.Fprintf(w, "    movsd xmm4, qword ptr [rax%+d]\n", k*8)
					fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", p.fnWasmR0Base-reti*8-8)
				default:
					unreachable()
				}
			}
		} else {
			// 走寄存器返回
			// 需要转换为从WASM栈返回
			for k, retType := range fnCallType.Results {
				retNative := fnCallNative.Type.Return[k]
				assert(retNative.Reg != 0)

				reti := stk.Push(retType)
				switch retType {
				case token.I32:
					fmt.Fprintf(w, "    mov dword ptr [rbp%+d], %s\n",
						p.fnWasmR0Base-reti*8-8,
						x64.RegString(retNative.Reg),
					)
				case token.I64:
					fmt.Fprintf(w, "    mov qword ptr [rbp%+d], %s\n",
						p.fnWasmR0Base-reti*8-8,
						x64.RegString(retNative.Reg),
					)
				case token.F32:
					fmt.Fprintf(w, "    movss dword ptr [rbp%+d], %s\n",
						p.fnWasmR0Base-reti*8-8,
						x64.RegString(retNative.Reg),
					)
				case token.F64:
					fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], %s\n",
						p.fnWasmR0Base-reti*8-8,
						x64.RegString(retNative.Reg),
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
		if err := astutil.BuildFuncFrame(abi.X64, fnCallNative); err != nil {
			return err
		}

		// 统计调用子函数需要的最大栈空间
		if fnCallNative.ArgsSize > p.fnMaxCallArgsSize {
			p.fnMaxCallArgsSize = fnCallNative.ArgsSize
		}

		// 获取函数地址
		sp0 := p.fnWasmR0Base - stk.Pop(token.I32)*8 - 8
		p.gasCommentInFunc(w, "加载函数的地址")
		fmt.Fprintln(w)

		const r10 = "r10"
		p.gasCommentInFunc(w, fmt.Sprintf("%s = table[?]", r10))
		fmt.Fprintf(w, "    mov  rax, [rip+%s]\n", kTableAddrName)
		fmt.Fprintf(w, "    mov  %s, [rbp%+d]\n", r10, sp0)
		fmt.Fprintf(w, "    mov  %s, [rax+%s*8]\n", r10, r10)
		fmt.Fprintln(w)

		const r11 = "r11"
		p.gasCommentInFunc(w, fmt.Sprintf("%s = %s[%s]", r11, kTableFuncIndexListName, r10))
		fmt.Fprintf(w, "    lea  rax, [rip+%s]\n", kTableFuncIndexListName)
		fmt.Fprintf(w, "    mov  %s, [rax+%s*8]\n", r11, r10)
		fmt.Fprintln(w)

		p.gasCommentInFunc(w, fmt.Sprintf("call_indirect %s(...)", r11))
		p.gasCommentInFunc(w, fmt.Sprintf("type %s", fnCallType))

		// 参数列表
		// 出栈的顺序相反
		argList := make([]int, len(fnCallType.Params))
		for k := len(argList) - 1; k >= 0; k-- {
			x := fnCallType.Params[k]
			argList[k] = stk.Pop(x.Type)
		}

		// 如果是走栈返回, 第一个是隐藏参数
		if len(fnCallNative.Type.Return) > 1 && fnCallNative.Type.Return[1].Reg == 0 {
			assert(p.osName == abi.WINDOWS.String())
			fmt.Fprintf(w, "    lea rcx, [rsp%+d] # return address\n", len(fnCallType.Params)*8)
		}

		// 准备调用参数
		for k, x := range fnCallType.Params {
			switch x.Type {
			case token.I32:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    mov %s, dword ptr [rbp%+d] # arg %d\n",
						x64.RegString(arg.Reg),
						p.fnWasmR0Base-argList[k]*8-8,
						k,
					)
				} else {
					fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n",
						p.fnWasmR0Base-argList[k]*8-8,
					)
					fmt.Fprintf(w, "    mov dword ptr [rsp%+d], eax\n",
						arg.RSPOff,
					)
				}
			case token.I64:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    mov %s, qword ptr [rbp%+d] # arg %d\n",
						x64.RegString(arg.Reg),
						p.fnWasmR0Base-argList[k]*8-8,
						k,
					)
				} else {
					fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n",
						p.fnWasmR0Base+argList[k]*8,
					)
					fmt.Fprintf(w, "    mov qword ptr [rsp%+d], rax\n",
						arg.RSPOff,
					)
				}
			case token.F32:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    movss %s, qword ptr [rbp%+d] # arg %d\n",
						x64.RegString(arg.Reg),
						p.fnWasmR0Base-argList[k]*8-8,
						k,
					)
				} else {
					fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n",
						p.fnWasmR0Base+argList[k]*8,
					)
					fmt.Fprintf(w, "    movss dword ptr [rsp%+d], xmm4\n",
						arg.RSPOff,
					)
				}
			case token.F64:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    movsd %s, qword ptr [rbp%+d] # arg %d\n",
						x64.RegString(arg.Reg),
						p.fnWasmR0Base-argList[k]*8-8,
						k,
					)
				} else {
					fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n",
						p.fnWasmR0Base+argList[k]*8,
					)
					fmt.Fprintf(w, "    movsd qword ptr [rsp%+d], xmm4\n",
						arg.RSPOff,
					)
				}
			default:
				unreachable()
			}
		}
		fmt.Fprintf(w, "    call %s\n", r11)

		// 提取返回值
		if len(fnCallNative.Type.Return) > 1 && fnCallNative.Type.Return[1].Reg == 0 {
			// 走栈返回结果, 地址通过 rcx 传入, 通过 rax 返回
			// 根据 ABI 规范获取返回值, 再保存到 WASM 栈中
			for k, retType := range fnCallType.Results {
				reti := stk.Push(retType)
				switch retType {
				case token.I32:
					fmt.Fprintf(w, "    mov r10, dword ptr [rax%+d]\n", k*8)
					fmt.Fprintf(w, "    mov dword ptr [rbp%+d], r10\n", p.fnWasmR0Base-reti*8-8)
				case token.I64:
					fmt.Fprintf(w, "    mov r10, qword ptr [rax%+d]\n", k*8)
					fmt.Fprintf(w, "    mov qword ptr [rbp%+d], r10\n", p.fnWasmR0Base-reti*8-8)
				case token.F32:
					fmt.Fprintf(w, "    movss xmm4, dword ptr [rax%+d]\n", k*8)
					fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", p.fnWasmR0Base-reti*8-8)
				case token.F64:
					fmt.Fprintf(w, "    movsd xmm4, qword ptr [rax%+d]\n", k*8)
					fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", p.fnWasmR0Base-reti*8-8)
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
					fmt.Fprintf(w, "    mov dword ptr [rbp%+d], %s\n",
						p.fnWasmR0Base-reti*8-8,
						x64.RegString(retNative.Reg),
					)
				case token.I64:
					fmt.Fprintf(w, "    mov qword ptr [rbp%+d], %s\n",
						p.fnWasmR0Base-reti*8-8,
						x64.RegString(retNative.Reg),
					)
				case token.F32:
					fmt.Fprintf(w, "    movss dword ptr [rbp%+d], %s\n",
						p.fnWasmR0Base-reti*8-8,
						x64.RegString(retNative.Reg),
					)
				case token.F64:
					fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], %s\n",
						p.fnWasmR0Base-reti*8-8,
						x64.RegString(retNative.Reg),
					)
				default:
					unreachable()
				}
			}
		}

	case token.INS_DROP:
		sp0 := stk.DropAny()
		fmt.Fprintf(w, "    nop # drop [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
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

		p.gasCommentInFunc(w, "select")
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-spCondition*8)
		fmt.Fprintf(w, "    test eax, eax\n")

		switch valType {
		case token.I32:
			fmt.Fprintf(w, "    mov    r10, dword ptr [rbp - %d]\n", p.fnWasmR0Base-spValueFalse*8)
			fmt.Fprintf(w, "    mov    r11, dword ptr [rbp - %d]\n", p.fnWasmR0Base-spValueTrue*8)
			fmt.Fprintf(w, "    cmovne r10, r11")
			fmt.Fprintf(w, "    mov    dword ptr [rbp - %d], r10", p.fnWasmR0Base-ret0*8)
		case token.I64:
			fmt.Fprintf(w, "    mov    r10, qword ptr [rbp - %d]\n", p.fnWasmR0Base-spValueFalse*8)
			fmt.Fprintf(w, "    mov    r11, qword ptr [rbp - %d]\n", p.fnWasmR0Base-spValueTrue*8)
			fmt.Fprintf(w, "    cmovne r10, r11")
			fmt.Fprintf(w, "    mov    qword ptr [rbp - %d], r10", p.fnWasmR0Base-ret0*8)
		case token.F32:
			fmt.Fprintf(w, "    movss  xmm4, dword ptr [rbp - %d]\n", p.fnWasmR0Base-spValueFalse*8)
			fmt.Fprintf(w, "    movss  xmm5, dword ptr [rbp - %d]\n", p.fnWasmR0Base-spValueTrue*8)
			fmt.Fprintf(w, "    cmovne xmm4, xmm5")
			fmt.Fprintf(w, "    movss  dword ptr [rbp - %d], xmm4", p.fnWasmR0Base-ret0*8)
		case token.F64:
			fmt.Fprintf(w, "    movsd  xmm4, qword ptr [rbp - %d]\n", p.fnWasmR0Base-spValueFalse*8)
			fmt.Fprintf(w, "    movsd  xmm5, qword ptr [rbp - %d]\n", p.fnWasmR0Base-spValueTrue*8)
			fmt.Fprintf(w, "    cmovne xmm4, xmm5")
			fmt.Fprintf(w, "    movsd  qword ptr [rbp - %d], xmm4", p.fnWasmR0Base-ret0*8)
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
			fmt.Fprintf(w, "    # local.get %s i32\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", xOff)
			fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", p.fnWasmR0Base-ret0*8-8)
			fmt.Fprintln(w)
		case token.I64:
			fmt.Fprintf(w, "    # local.get %s i64\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", xOff)
			fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", p.fnWasmR0Base-ret0*8-8)
			fmt.Fprintln(w)
		case token.F32:
			fmt.Fprintf(w, "    # local.get %s f32\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", xOff)
			fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", p.fnWasmR0Base-ret0*8-8)
			fmt.Fprintln(w)
		case token.F64:
			fmt.Fprintf(w, "    # local.get %s f64\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", xOff)
			fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", p.fnWasmR0Base-ret0*8-8)
			fmt.Fprintln(w)
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
			fmt.Fprintf(w, "    # local.set %s i32\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", xOff)
			fmt.Fprintln(w)
		case token.I64:
			fmt.Fprintf(w, "    # local.set %s i64\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", xOff)
			fmt.Fprintln(w)
		case token.F32:
			fmt.Fprintf(w, "    # local.set %s f32\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", xOff)
			fmt.Fprintln(w)
		case token.F64:
			fmt.Fprintf(w, "    # local.set %s f64\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", xOff)
			fmt.Fprintln(w)
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
			fmt.Fprintf(w, "    # local.tee %s i32\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8)
			fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", xOff)
		case token.I64:
			fmt.Fprintf(w, "    # local.tee %s i64\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8)
			fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", xOff)
		case token.F32:
			fmt.Fprintf(w, "    # local.tee %s f32\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8)
			fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", xOff)
		case token.F64:
			fmt.Fprintf(w, "    # local.tee %s f64\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8)
			fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", xOff)
		default:
			unreachable()
		}
	case token.INS_GLOBAL_GET:
		i := i.(ast.Ins_GlobalGet)
		xType := p.findGlobalType(i.X)
		ret0 := stk.Push(xType)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "    # global.get %s i32\n", i.X)
			fmt.Fprintf(w, "    mov eax, dword ptr [rip+%s]\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", p.fnWasmR0Base-ret0*8)
		case token.I64:
			fmt.Fprintf(w, "    # global.get %s i64\n", i.X)
			fmt.Fprintf(w, "    mov rax, qword ptr [rip+%s]\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", p.fnWasmR0Base-ret0*8)
		case token.F32:
			fmt.Fprintf(w, "    # global.get %s f32\n", i.X)
			fmt.Fprintf(w, "    mov eax, dword ptr [rip+%s]\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", p.fnWasmR0Base-ret0*8)
		case token.F64:
			fmt.Fprintf(w, "    # global.get %s f64\n", i.X)
			fmt.Fprintf(w, "    mov rax, qword ptr [rip+%s]\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", p.fnWasmR0Base-ret0*8)
		default:
			unreachable()
		}
	case token.INS_GLOBAL_SET:
		i := i.(ast.Ins_GlobalSet)
		xType := p.findGlobalType(i.X)
		sp0 := stk.Pop(xType)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "    # global.set %s i32\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8)
			fmt.Fprintf(w, "    mov dword ptr [rip+%s], eax\n", p.findLocalName(fn, i.X))
		case token.I64:
			fmt.Fprintf(w, "    # global.set %s i64\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8)
			fmt.Fprintf(w, "    mov qword ptr [rip+%s], rax\n", p.findLocalName(fn, i.X))
		case token.F32:
			fmt.Fprintf(w, "    # global.set %s f32\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8)
			fmt.Fprintf(w, "    mov dword ptr [rip+%s], eax\n", p.findLocalName(fn, i.X))
		case token.F64:
			fmt.Fprintf(w, "    # global.set %s f64\n", p.findLocalName(fn, i.X))
			fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8)
			fmt.Fprintf(w, "    mov qword ptr [rip+%s], rax\n", p.findLocalName(fn, i.X))
		default:
			unreachable()
		}

	case token.INS_TABLE_GET:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.FUNCREF) // funcref
		fmt.Fprintf(w, "    # table.get\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kTableAddrName)
		fmt.Fprintf(w, "    movsxd r10, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov rax, dword ptr [rax+r10]\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], rax\n", ret0)
	case token.INS_TABLE_SET:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.FUNCREF) // funcref
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		fmt.Fprintf(w, "    # table.set\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kTableAddrName)
		fmt.Fprintf(w, "    movsxd r10, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsxd r11, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov rax, qword ptr [rax+r11]\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+r11], r10\n")

	case token.INS_I32_LOAD:
		i := i.(ast.Ins_I32Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.load\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov eax, dword [r10 %+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov dword [rbp%+d], eax\n", ret0)

	case token.INS_I64_LOAD:
		i := i.(ast.Ins_I64Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.load\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov rax, qword [r10 %+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov qword [rbp%+d], rax\n", ret0)

	case token.INS_F32_LOAD:
		i := i.(ast.Ins_F32Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.load\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    movss xmm4, dword [r10 %+d]\n", i.Offset)
		fmt.Fprintf(w, "    movss dword [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_LOAD:
		i := i.(ast.Ins_I32Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.load\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    movsd xmm4, dword [r10 %+d]\n", i.Offset)
		fmt.Fprintf(w, "    movsd qword [rbp%+d], xmm4\n", ret0)

	case token.INS_I32_LOAD8_S:
		i := i.(ast.Ins_I32Load8S)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.load8_s\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    movsx eax, byte [r10 %+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov dword [rbp%+d], eax\n", ret0)

	case token.INS_I32_LOAD8_U:
		i := i.(ast.Ins_I32Load8U)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.load8_u\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    movzx eax, byte [r10 %+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov dword [rbp%+d], eax\n", ret0)

	case token.INS_I32_LOAD16_S:
		i := i.(ast.Ins_I32Load16S)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.load16_s\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    movsx eax, word [r10 %+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov dword [rbp%+d], eax\n", ret0)

	case token.INS_I32_LOAD16_U:
		i := i.(ast.Ins_I32Load16U)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.load16_u\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    movzx eax, word [r10 %+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov dword [rbp%+d], eax\n", ret0)

	case token.INS_I64_LOAD8_S:
		i := i.(ast.Ins_I64Load8S)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.load8_s\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    movsx rax, byte [r10 %+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov qword [rbp%+d], rax\n", ret0)

	case token.INS_I64_LOAD8_U:
		i := i.(ast.Ins_I64Load8U)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.load8_u\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    movzx rax, byte [r10 %+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov qword [rbp%+d], rax\n", ret0)

	case token.INS_I64_LOAD16_S:
		i := i.(ast.Ins_I64Load16S)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.load16_s\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    movsx eax, word [r10 %+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov qword [rbp%+d], eax\n", ret0)

	case token.INS_I64_LOAD16_U:
		i := i.(ast.Ins_I64Load16U)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.load16_u\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    movzx eax, word [r10 %+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov qword [rbp%+d], eax\n", ret0)

	case token.INS_I64_LOAD32_S:
		i := i.(ast.Ins_I64Load32S)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.load32_s\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    movsx eax, dword [r10 %+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov qword [rbp%+d], eax\n", ret0)

	case token.INS_I64_LOAD32_U:
		i := i.(ast.Ins_I64Load32U)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.load32_u\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    movzx eax, dword [r10 %+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov qword [rbp%+d], eax\n", ret0)

	case token.INS_I32_STORE:
		i := i.(ast.Ins_I32Store)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)

		fmt.Fprintf(w, "    # i32.store\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov eax, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov dword [r10 %+d], eax\n", i.Offset)

	case token.INS_I64_STORE:
		i := i.(ast.Ins_I64Store)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)

		fmt.Fprintf(w, "    # i64.store\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov rax, qword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov qword [r10 %+d], rax\n", i.Offset)

	case token.INS_F32_STORE:
		i := i.(ast.Ins_F32Store)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)

		fmt.Fprintf(w, "    # f32.store\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov xmm4, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov dword [r10 %+d], xmm4\n", i.Offset)

	case token.INS_F64_STORE:
		i := i.(ast.Ins_F64Store)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)

		fmt.Fprintf(w, "    # f64.store\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov xmm4, qword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov qword [r10 %+d], xmm4\n", i.Offset)

	case token.INS_I32_STORE8:
		i := i.(ast.Ins_I32Store8)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)

		fmt.Fprintf(w, "    # i32.store8\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov eax, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov byte [r10 %+d], eax\n", i.Offset)

	case token.INS_I32_STORE16:
		i := i.(ast.Ins_I32Store16)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)

		fmt.Fprintf(w, "    # i64.store16\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov eax, dword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov word [r10 %+d], eax\n", i.Offset)

	case token.INS_I64_STORE8:
		i := i.(ast.Ins_I64Store8)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)

		fmt.Fprintf(w, "    # i64.store8\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov eax, qword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov byte [r10 %+d], eax\n", i.Offset)

	case token.INS_I64_STORE16:
		i := i.(ast.Ins_I64Store16)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)

		fmt.Fprintf(w, "    # i64.store16\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov eax, qword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov word [r10 %+d], eax\n", i.Offset)

	case token.INS_I64_STORE32:
		i := i.(ast.Ins_I64Store32)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)

		fmt.Fprintf(w, "    # i64.store32\n")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10, dword [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov eax, qword [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov dword [r10 %+d], eax\n", i.Offset)

	case token.INS_MEMORY_SIZE:
		sp0 := p.fnWasmR0Base - 8*stk.Push(token.I32)
		fmt.Fprintf(w, "    # memory.size\n")
		fmt.Fprintf(w, "    lea rax, [rip+%s]\n", kMemoryPagesName)
		fmt.Fprintf(w, "    mov [rbp%+d], rax\n", sp0)

	case token.INS_MEMORY_GROW:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		labelBody := p.gasGenNextId(kLabelNamePreifx_ifBody)
		labelElse := p.gasGenNextId(kLabelNamePreifx_ifElse)
		labelEnd := p.gasGenNextId(kLabelNamePreifx_ifEnd)

		fmt.Fprintf(w, "    # memory.grow\n")
		fmt.Fprintf(w, "    lea r10, [rip+%s]\n", kMemoryPagesName)
		fmt.Fprintf(w, "    lea r11, [rip+%s]\n", kMemoryMaxPagesName)
		fmt.Fprintf(w, "    mov rax, [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add rax, r10\n")
		fmt.Fprintf(w, "    cmp rax, r11\n")
		fmt.Fprintf(w, "    jg  %s\n", labelElse)

		p.gasFuncLabel(w, labelBody)
		fmt.Fprintf(w, "    mov [rbp%+d], rax\n", ret0)
		fmt.Fprintf(w, "    jmp %s\n", labelEnd)

		p.gasFuncLabel(w, labelElse)
		fmt.Fprintf(w, "    mov rax, -1\n")
		fmt.Fprintf(w, "    mov [rbp%+d], rax\n", ret0)

		p.gasFuncLabel(w, labelEnd)

	case token.INS_MEMORY_INIT:
		panic("unsupport memory.init")

	case token.INS_MEMORY_COPY:
		panic("unsupport memory.copy")

	case token.INS_MEMORY_FILL:
		panic("unsupport memory.fill")

	case token.INS_I32_CONST:
		i := i.(ast.Ins_I32Const)
		sp0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8
		fmt.Fprintf(w, "    # i32.const %d\n", i.X)
		fmt.Fprintf(w, "    mov eax, %d\n", i.X)
		fmt.Fprintf(w, "    mov [rbp%+d], eax\n", sp0)
		fmt.Fprintln(w)

	case token.INS_I64_CONST:
		i := i.(ast.Ins_I64Const)
		aa := stk.Push(token.I64)
		sp0 := p.fnWasmR0Base - 8*aa - 8
		fmt.Fprintf(w, "    # i64.const %d\n", i.X)
		fmt.Fprintf(w, "    movabs rax, %d\n", int64(i.X))
		fmt.Fprintf(w, "    mov    [rbp%+d], rax\n", sp0)
		fmt.Fprintln(w)

	case token.INS_F32_CONST:
		i := i.(ast.Ins_F32Const)
		sp0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8
		fmt.Fprintf(w, "    # f32.const %f\n", i.X)
		fmt.Fprintf(w, "    mov  eax, 0x%X\n", math.Float32bits(i.X))
		fmt.Fprintf(w, "    movd [rbp%+d], eax\n", sp0)
		fmt.Fprintln(w)

	case token.INS_F64_CONST:
		i := i.(ast.Ins_F64Const)
		sp0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8
		fmt.Fprintf(w, "    # f64.const %f\n", i.X)
		fmt.Fprintf(w, "    movabs rax, 0x%X\n", math.Float64bits(i.X))
		fmt.Fprintf(w, "    movq   [rbp%+d], xmm4\n", sp0)
		fmt.Fprintln(w)

	case token.INS_I32_EQZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.eqz\n")
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    test eax, eax\n")
		fmt.Fprintf(w, "    setz al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_EQ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.eq\n")
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10d, r11d\n")
		fmt.Fprintf(w, "    sete al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_NE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.ne\n")
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10d, r11d\n")
		fmt.Fprintf(w, "    setne al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_LT_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.lt_s\n")
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10d, r11d\n")
		fmt.Fprintf(w, "    setl al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_LT_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.lt_u\n")
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10d, r11d\n")
		fmt.Fprintf(w, "    setb al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_GT_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.gt_s\n")
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10d, r11d\n")
		fmt.Fprintf(w, "    setg al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_GT_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.gt_u\n")
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10d, r11d\n")
		fmt.Fprintf(w, "    seta al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_LE_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.le_s\n")
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10d, r11d\n")
		fmt.Fprintf(w, "    setle al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_LE_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.le_u\n")
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10d, r11d\n")
		fmt.Fprintf(w, "    setbe al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_GE_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.ge_s\n")
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10d, r11d\n")
		fmt.Fprintf(w, "    setge al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_GE_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.ge_u\n")
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10d, r11d\n")
		fmt.Fprintf(w, "    setae al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I64_EQZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i64.eqz\n")
		fmt.Fprintf(w, "    mov  rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    test rax, rax\n")
		fmt.Fprintf(w, "    setz al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I64_EQ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i64.eq\n")
		fmt.Fprintf(w, "    mov r10, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10, r11\n")
		fmt.Fprintf(w, "    sete al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I64_NE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i64.ne\n")
		fmt.Fprintf(w, "    mov r10, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10, r11\n")
		fmt.Fprintf(w, "    setne al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I64_LT_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i64.lt_s\n")
		fmt.Fprintf(w, "    mov r10, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10, r11\n")
		fmt.Fprintf(w, "    setl al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I64_LT_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i64.lt_u\n")
		fmt.Fprintf(w, "    mov r10, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10, r11\n")
		fmt.Fprintf(w, "    setb al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I64_GT_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i64.gt_s\n")
		fmt.Fprintf(w, "    mov r10, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10, r11\n")
		fmt.Fprintf(w, "    setg al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I64_GT_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i64.gt_u\n")
		fmt.Fprintf(w, "    mov r10, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10, r11\n")
		fmt.Fprintf(w, "    seta al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I64_LE_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i64.le_s\n")
		fmt.Fprintf(w, "    mov r10, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10, r11\n")
		fmt.Fprintf(w, "    setle al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I64_LE_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i64.le_u\n")
		fmt.Fprintf(w, "    mov r10, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10, r11\n")
		fmt.Fprintf(w, "    setbe al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I64_GE_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i64.ge_s\n")
		fmt.Fprintf(w, "    mov r10, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10, r11\n")
		fmt.Fprintf(w, "    setge al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I64_GE_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i64.ge_u\n")
		fmt.Fprintf(w, "    mov r10, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cmp r10, r11\n")
		fmt.Fprintf(w, "    setae al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_F32_EQ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # f32.eq\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss xmm5, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    setz al\n")
		fmt.Fprintf(w, "    setnp cl\n")
		fmt.Fprintf(w, "    and al, cl\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_F32_NE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # f32.ne\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss xmm5, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    setne al\n")
		fmt.Fprintf(w, "    setp cl\n")
		fmt.Fprintf(w, "    or al, cl\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_F32_LT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # f32.lt\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss xmm5, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    setb al\n")
		fmt.Fprintf(w, "    setnp cl\n")
		fmt.Fprintf(w, "    and al, cl\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_F32_GT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # f32.gt\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss xmm5, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    seta al\n")
		fmt.Fprintf(w, "    setnp cl\n")
		fmt.Fprintf(w, "    and al, cl\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_F32_LE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # f32.le\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss xmm5, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    setbe al\n")
		fmt.Fprintf(w, "    setnp cl\n")
		fmt.Fprintf(w, "    and al, cl\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_F32_GE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # f32.ge\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss xmm5, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    setae al\n")
		fmt.Fprintf(w, "    setnp cl\n")
		fmt.Fprintf(w, "    and al, cl\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_F64_EQ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # f64.eq\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd xmm5, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    setz al\n")
		fmt.Fprintf(w, "    setnp cl\n")
		fmt.Fprintf(w, "    and al, cl\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_F64_NE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # f64.ne\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd xmm5, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    setne al\n")
		fmt.Fprintf(w, "    setp cl\n")
		fmt.Fprintf(w, "    or al, cl\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_F64_LT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # f64.lt\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd xmm5, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    setb al\n")
		fmt.Fprintf(w, "    setnp cl\n")
		fmt.Fprintf(w, "    and al, cl\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_F64_GT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # f64.gt\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd xmm5, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    seta al\n")
		fmt.Fprintf(w, "    setnp cl\n")
		fmt.Fprintf(w, "    and al, cl\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_F64_LE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # f64.le\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd xmm5, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    setbe al\n")
		fmt.Fprintf(w, "    setnp cl\n")
		fmt.Fprintf(w, "    and al, cl\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_F64_GE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # f64.ge\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd xmm5, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    setae al\n")
		fmt.Fprintf(w, "    setnp cl\n")
		fmt.Fprintf(w, "    and al, cl\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_CLZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.clz\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    lzcnt eax, eax\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_CTZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.ctz\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    tzcnt eax, eax\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_POPCNT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.popcnt\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    popcnt eax, eax\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_ADD:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.add\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_SUB:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.sub\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    sub eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_MUL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.mul\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mul eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_DIV_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.div_s\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+8], edx\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cdq\n")
		fmt.Fprintf(w, "    idiv dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintf(w, "    mov edx, qword ptr [rbp+8]\n")

	case token.INS_I32_DIV_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.div_u\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+8], edx\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    xor edx, edx\n")
		fmt.Fprintf(w, "    div dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintf(w, "    mov edx, qword ptr [rbp+8]\n")

	case token.INS_I32_REM_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.rem_s\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+8], edx\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cdq\n")
		fmt.Fprintf(w, "    idiv dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], edx\n", ret0)
		fmt.Fprintf(w, "    mov edx, qword ptr [rbp+8]\n")

	case token.INS_I32_REM_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.rem_u\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+8], edx\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    xor edx, edx\n")
		fmt.Fprintf(w, "    div dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], edx\n", ret0)
		fmt.Fprintf(w, "    mov edx, qword ptr [rbp+8]\n")

	case token.INS_I32_AND:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.and\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    and eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_OR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.or\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    or eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_XOR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.xor\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    xor eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_SHL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.shl\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+0], rcx\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov ecx, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    shl eax, cl # cl 是 ecx 低8位\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintf(w, "    mov rcx, qword ptr [rbp+0]\n")

	case token.INS_I32_SHR_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.shr_s\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+0], rcx\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov ecx, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    sar eax, cl # cl 是 ecx 低8位\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintf(w, "    mov rcx, qword ptr [rbp+0]\n")

	case token.INS_I32_SHR_U:
		sp0 := stk.Pop(token.I32)
		sp1 := stk.Pop(token.I32)
		ret0 := stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.shr_u\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+0], rcx\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov ecx, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    shr eax, cl # cl 是 ecx 低8位\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintf(w, "    mov rcx, qword ptr [rbp+0]\n")

	case token.INS_I32_ROTL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.rotl\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+0], rcx\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov ecx, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    rol eax, cl # cl 是 ecx 低8位\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintf(w, "    mov rcx, qword ptr [rbp+0]\n")

	case token.INS_I32_ROTR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.rotr\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+0], rcx\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov ecx, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    ror eax, cl # cl 是 ecx 低8位\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintf(w, "    mov rcx, qword ptr [rbp+0]\n")

	case token.INS_I64_CLZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i64.clz\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    lzcnt rax, rax\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], rax\n", ret0)

	case token.INS_I64_CTZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i64.ctz\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    tzcnt rax, rax\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], rax\n", ret0)

	case token.INS_I64_POPCNT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i64.popcnt\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    popcnt rax, rax\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], rax\n", ret0)

	case token.INS_I64_ADD:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.add\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)

	case token.INS_I64_SUB:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.sub\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    sub rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)

	case token.INS_I64_MUL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.mul\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mul rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)

	case token.INS_I64_DIV_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.div_s\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+8], rdx\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cdq\n")
		fmt.Fprintf(w, "    idiv qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintf(w, "    mov rdx, qword ptr [rbp+8]\n")

	case token.INS_I64_DIV_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.div_u\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+8], rdx\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    xor rdx, rdx\n")
		fmt.Fprintf(w, "    div qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintf(w, "    mov rdx, qword ptr [rbp+8]\n")

	case token.INS_I64_REM_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.rem_s\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+8], rdx\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cdq\n")
		fmt.Fprintf(w, "    idiv qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rdx\n", ret0)
		fmt.Fprintf(w, "    mov rdx, qword ptr [rbp+8]\n")

	case token.INS_I64_REM_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.rem_u\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+8], rdx\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    xor rdx, rdx\n")
		fmt.Fprintf(w, "    div qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rdx\n", ret0)
		fmt.Fprintf(w, "    mov rdx, qword ptr [rbp+8]\n")

	case token.INS_I64_AND:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.and\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    and rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)

	case token.INS_I64_OR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.or\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    or  rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)

	case token.INS_I64_XOR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.xor\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    xor rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)

	case token.INS_I64_SHL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.shl\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+0], rcx\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov rcx, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    shl rax, cl # cl 是 rcx 低8位\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintf(w, "    mov rcx, qword ptr [rbp+0]\n")

	case token.INS_I64_SHR_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.shr_s\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+0], rcx\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov rcx, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    sar rax, cl # cl 是 rcx 低8位\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintf(w, "    mov rcx, qword ptr [rbp+0]\n")

	case token.INS_I64_SHR_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.shr_u\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+0], rcx\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov rcx, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    shr rax, cl # cl 是 rcx 低8位\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintf(w, "    mov rcx, qword ptr [rbp+0]\n")

	case token.INS_I64_ROTL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.rotl\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+0], rcx\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov rcx, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    rol rax, cl # cl 是 rcx 低8位\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintf(w, "    mov rcx, qword ptr [rbp+0]\n")

	case token.INS_I64_ROTR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.rotr\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp+0], rcx\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov rcx, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    ror rax, cl # cl 是 rcx 低8位\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintf(w, "    mov rcx, qword ptr [rbp+0]\n")

	case token.INS_F32_ABS:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.abs\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov   eax, 0x7FFFFFFF\n")
		fmt.Fprintf(w, "    movd  xmm5, eax\n")
		fmt.Fprintf(w, "    andps xmm4, xmm5\n")
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_NEG:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.neg\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov eax, 0x80000000\n")
		fmt.Fprintf(w, "    movd xmm5, eax\n")
		fmt.Fprintf(w, "    xorps xmm4, xmm5\n")
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_CEIL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.ceil\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundss xmm4, xmm4, 0x02\n")
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_FLOOR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.floor\n")
		fmt.Fprintf(w, "    movss   xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundss xmm4, xmm4, 0x01\n")
		fmt.Fprintf(w, "    movss   dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_TRUNC:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.trunc\n")
		fmt.Fprintf(w, "    movss   xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundss xmm4, xmm4, 0x03\n")
		fmt.Fprintf(w, "    movss   dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_NEAREST:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.nearest\n")
		fmt.Fprintf(w, "    movss   xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundss xmm4, xmm4, 0x00\n")
		fmt.Fprintf(w, "    movss   dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_SQRT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.sqrt\n")
		fmt.Fprintf(w, "    movss  xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    sqrtss xmm4, xmm4\n")
		fmt.Fprintf(w, "    movss  dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_ADD:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.add\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    addss xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_SUB:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.sub\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    subss xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_MUL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.mul\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mulss xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_DIV:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.div\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    divss xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_MIN:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.min\n")
		fmt.Fprintf(w, "    movss   xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss   xmm5, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    minss   xmm4, xmm5\n")
		fmt.Fprintf(w, "    movss   dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_MAX:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		labelNotEqual := p.gasGenNextId(kLabelNamePreifx_ifEnd)

		fmt.Fprintf(w, "    # f32.max\n")
		fmt.Fprintf(w, "    movss   xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss   xmm5, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    maxss   xmm4, xmm5\n")
		fmt.Fprintf(w, "    maxss   xmm5, xmm4\n")
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    jne     %s\n", labelNotEqual)
		fmt.Fprintf(w, "    andps   xmm4, xmm5\n")
		fmt.Fprintf(w, "%s:\n", labelNotEqual)
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_COPYSIGN:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.copysign\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss xmm5, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov eax, 0x7FFFFFFF\n")
		fmt.Fprintf(w, "    movd xmm6, eax\n")
		fmt.Fprintf(w, "    movd r10d, xmm4\n")
		fmt.Fprintf(w, "    and r10d, eax\n")
		fmt.Fprintf(w, "    not eax\n")
		fmt.Fprintf(w, "    movd r11d, xmm5\n")
		fmt.Fprintf(w, "    and r11d, eax\n")
		fmt.Fprintf(w, "    or r10d, r11d\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], r10d\n", ret0)

	case token.INS_F64_ABS:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.abs\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov rax, 0x7FFFFFFFFFFFFFFF\n")
		fmt.Fprintf(w, "    movq xmm5, rax\n")
		fmt.Fprintf(w, "    andpd xmm4, xmm5\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_NEG:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.neg\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov rax, 0x8000000000000000\n")
		fmt.Fprintf(w, "    movq xmm5, rax\n")
		fmt.Fprintf(w, "    xorpd xmm4, xmm5\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_CEIL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.ceil\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundsd xmm4, xmm4, 0x02\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_FLOOR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.floor\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundsd xmm4, xmm4, 0x01\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_TRUNC:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.trunc\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundsd xmm4, xmm4, 0x03\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_NEAREST:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.nearest\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundsd xmm4, xmm4, 0x00\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_SQRT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.sqrt\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    sqrtsd xmm4, xmm4\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_ADD:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.add\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd xmm5, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    addsd xmm4, xmm5\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_SUB:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.sub\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd xmm5, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    subsd xmm4, xmm5\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_MUL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.mul\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd xmm5, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mulsd xmm4, xmm5\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_DIV:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.div\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd xmm5, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    divsd xmm4, xmm5\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_MIN:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		labelIfDone := p.gasGenNextId(kLabelNamePreifx_ifEnd)

		fmt.Fprintf(w, "    # f64.min\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd xmm5, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    minsd xmm4, xmm5\n")
		fmt.Fprintf(w, "    minsd xmm5, xmm4\n")
		fmt.Fprintf(w, "    ucomisd xmm4, xmm5\n")
		fmt.Fprintf(w, "    jne %s\n", labelIfDone)
		fmt.Fprintf(w, "    orpd xmm4, xmm5\n")
		fmt.Fprintf(w, "%s:\n", labelIfDone)
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_MAX:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		labelNotEqual := p.gasGenNextId(kLabelNamePreifx_ifEnd)

		fmt.Fprintf(w, "    # f64.max\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd xmm5, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    maxsd xmm4, xmm5\n")
		fmt.Fprintf(w, "    maxsd xmm5, xmm4\n")
		fmt.Fprintf(w, "    ucomisd xmm4, xmm5\n")
		fmt.Fprintf(w, "    jne %s\n", labelNotEqual)
		fmt.Fprintf(w, "    andpd xmm4, xmm5\n")
		fmt.Fprintf(w, "%s:\n", labelNotEqual)
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_COPYSIGN:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.copysign\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r10, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov r11, 0x7FFFFFFFFFFFFFFF\n")
		fmt.Fprintf(w, "    and rax, r11\n")
		fmt.Fprintf(w, "    not r11\n")
		fmt.Fprintf(w, "    and r10, r11\n")
		fmt.Fprintf(w, "    or  rax, r10\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)

	case token.INS_I32_WRAP_I64:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.wrap_i64\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_TRUNC_F32_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.trunc_f32_s\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttss2si eax, xmm4\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_TRUNC_F32_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.trunc_f32_u\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttss2si rax, xmm4\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_TRUNC_F64_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.trunc_f64_s\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttsd2si eax, xmm4\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I32_TRUNC_F64_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32)

		fmt.Fprintf(w, "    # i32.trunc_f64_u\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttsd2si rax, xmm4\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I64_EXTEND_I32_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.extend_i32_s\n")
		fmt.Fprintf(w, "    movsxd rax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)

	case token.INS_I64_EXTEND_I32_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.extend_i32_u\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)

	case token.INS_I64_TRUNC_F32_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.trunc_f32_s\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttss2si rax, xmm4\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)

	case token.INS_I64_TRUNC_F32_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.trunc_f32_u\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttss2si rax, xmm4\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)

	case token.INS_I64_TRUNC_F64_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.trunc_f64_s\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttsd2si rax, xmm4\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)

	case token.INS_I64_TRUNC_F64_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64)

		fmt.Fprintf(w, "    # i64.trunc_f64_u\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttsd2si rax, xmm4\n")
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)

	case token.INS_F32_CONVERT_I32_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.convert_i32_s\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2ss xmm4, eax\n")
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_CONVERT_I32_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.convert_i32_u\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2ss xmm4, rax\n")
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_CONVERT_I64_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.convert_i64_s\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2ss xmm4, rax\n")
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_CONVERT_I64_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.convert_i64_u\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2ss xmm4, rax\n")
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F32_DEMOTE_F64:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.demote_f64\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsd2ss xmm4, xmm4\n")
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_CONVERT_I32_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.convert_i32_s\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2sd xmm4, eax\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_CONVERT_I32_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.convert_i32_u\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2sd xmm4, rax\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_CONVERT_I64_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.convert_i64_s\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2sd xmm4, rax\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_CONVERT_I64_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.convert_i64_u\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2sd xmm4, rax\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_F64_PROMOTE_F32:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.promote_f32\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtss2sd xmm4, xmm4\n")
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)

	case token.INS_I32_REINTERPRET_F32:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # i32.reinterpret_f32\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_I64_REINTERPRET_F64:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # i64.reinterpret_f64\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)

	case token.INS_F32_REINTERPRET_I32:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32)

		fmt.Fprintf(w, "    # f32.reinterpret_i32\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)

	case token.INS_F64_REINTERPRET_I64:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64)
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64)

		fmt.Fprintf(w, "    # f64.reinterpret_i64\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)

	default:
		panic(fmt.Sprintf("unreachable: %T", i))
	}
	return nil
}
