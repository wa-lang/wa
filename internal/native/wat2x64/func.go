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
	kLabelPrefixName_brNext        = ".L.brNext."
	kLabelPrefixName_brCase        = ".L.brCase."
	kLabelPrefixName_brDefault     = ".L.brDefault."
	kLabelPrefixName_brFallthrough = ".L.brFallthrough."
	kLabelPrefixName_else          = ".L.else."
	kLabelPrefixName_end           = ".L.end."
)

//
// 函数栈帧布局
// 参考 /docs/asm_abi_x64.md
//

func (p *wat2X64Worker) buildFuncs(w io.Writer) error {
	if len(p.m.Funcs) == 0 {
		return nil
	}

	for _, f := range p.m.Funcs {
		p.gasComment(w, "func "+f.Name+f.Type.String())
		p.gasSectionTextStart(w)
		if f.ExportName == f.Name {
			p.gasGlobal(w, kFuncNamePrefix+fixName(f.Name))
		}
		p.gasFuncStart(w, kFuncNamePrefix+fixName(f.Name))
		if err := p.buildFunc_body(w, f); err != nil {
			return err
		}
	}

	return nil
}

func (p *wat2X64Worker) buildFunc_body(w io.Writer, fn *ast.Func) error {
	p.Tracef("buildFunc_body: %s\n", fn.Name)

	assert(len(fn.Type.Results) == len(fn.Body.Results))

	// 转化为汇编的结构, 准备构建函数栈帧
	fnNative := wat2nativeFunc(fn.Name, fn.Type, fn.Locals)

	// 模拟构建函数栈帧
	// 后面需要根据参数是否走寄存器传递生成相关的入口代码和返回代码
	if err := astutil.BuildFuncFrame(p.cpuType, fnNative); err != nil {
		return err
	}

	// WASM 栈的开始位置
	// WASM 栈底第一个元素相对于 rbp 的偏移位置, 每个元素 8 字节
	// 栈指令依赖该偏移量
	p.fnWasmR0Base = 0 - fnNative.FrameSize
	p.fnMaxCallArgsSize = 0

	// 开始解析 wasm 指令
	var stk valueTypeStack
	var scopeStack scopeContextStack
	var bufHeader bytes.Buffer
	var bufBody bytes.Buffer
	var bufReturn bytes.Buffer
	{
		stk.funcName = fn.Name
		assert(stk.Len() == 0)

		// 定义局部变量
		if len(fn.Locals) > 0 {
			for _, x := range fn.Locals {
				p.gasCommentInFunc(w, fmt.Sprintf("local %s: %s", x.Name, x.Type))
			}
			fmt.Fprintln(&bufHeader)
		}

		// 便于函数主体指令
		if err := p.buildFunc_ins(&bufBody, fnNative, fn, &stk, &scopeStack, *fn.Body); err != nil {
			return err
		}
	}

	// 补充构建头部栈帧
	// 需要先统计 WASM 栈的最大深度和调用子函数的参数和返回值空间
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

	// 如果走栈返回
	// 先将调用者传入的返回值栈地址寄存器参数备份
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
				fmt.Fprintf(&bufHeader, "    mov [rbp%+d], %v # save arg.%d\n",
					arg.RBPOff, x64.Reg32String(arg.Reg),
					i,
				)
			case token.I64:
				fmt.Fprintf(&bufHeader, "    mov [rbp%+d], %v # save arg.%d\n",
					arg.RBPOff, x64.RegString(arg.Reg),
					i,
				)
			case token.F32:
				fmt.Fprintf(&bufHeader, "    movss dword ptr [rbp%+d], %v # save arg.%d\n",
					arg.RBPOff, x64.RegString(arg.Reg),
					i,
				)
			case token.F64:
				fmt.Fprintf(&bufHeader, "    movsd qword ptr [rbp%+d], %v # save arg.%d\n",
					arg.RBPOff, x64.RegString(arg.Reg),
					i,
				)
			default:
				panic("unreachable")
			}
		}
		fmt.Fprintln(&bufHeader)
	} else {
		p.gasCommentInFunc(&bufHeader, "没有参数需要备份到栈")
		fmt.Fprintln(&bufHeader)
	}

	// 返回值初始化为 0
	if len(fnNative.Type.Return) > 0 {
		p.gasCommentInFunc(&bufHeader, "将返回值变量初始化为0")
		for i, ret := range fnNative.Type.Return {
			fmt.Fprintf(&bufHeader,
				"    mov dword ptr [rbp%+d], 0 # ret.%d = 0\n",
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
			fmt.Fprintf(&bufHeader, "    mov dword ptr [rbp%+d], 0 # local %s = 0\n",
				local.RBPOff, local.Name,
			)
		}
		fmt.Fprintln(&bufHeader)
	} else {
		p.gasCommentInFunc(&bufHeader, "没有局部变量需要初始化为0")
		fmt.Fprintln(&bufHeader)
	}

	// 准备物理函数的返回值处理
	p.gasCommentInFunc(&bufReturn, "根据ABI处理返回值")
	fmt.Fprintln(&bufReturn)

	// 将栈上的值保存到返回值变量
	// 这是函数Body作为一个Block的返回值
	// 栈上多余的数据以及在 return 和 br 指令时刻搬运完成
	assert(stk.Len() == len(fn.Type.Results))
	if n := len(fn.Type.Results); n > 0 {
		p.gasCommentInFunc(&bufReturn, "将栈上数据复制到返回值变量")
		for i := n - 1; i >= 0; i-- {
			reti := fnNative.Type.Return[i]
			sp := p.fnWasmR0Base - 8*stk.Pop(fn.Type.Results[i]) - 8
			switch fn.Type.Results[i] {
			case token.I32:
				fmt.Fprintf(&bufReturn, "    mov eax,  dword ptr [rbp%+d]\n", sp)
				fmt.Fprintf(&bufReturn, "    mov dword ptr [rbp%+d], eax # ret.%d\n", reti.RBPOff, i)
			case token.I64:
				fmt.Fprintf(&bufReturn, "    mov rax, qword ptr [rbp%+d]\n", sp)
				fmt.Fprintf(&bufReturn, "    mov qword ptr [rbp%+d], rax # ret.%d\n", reti.RBPOff, i)
			case token.F32:
				fmt.Fprintf(&bufReturn, "    mov eax, dword ptr [rbp%+d]\n", sp)
				fmt.Fprintf(&bufReturn, "    mov dword ptr [rbp%+d], eax # ret.%d\n", reti.RBPOff, i)
			case token.F64:
				fmt.Fprintf(&bufReturn, "    mov rax, qword ptr [rbp%+d]\n", sp)
				fmt.Fprintf(&bufReturn, "    mov qword ptr [rbp%+d], rax # ret.%d\n", reti.RBPOff, i)
			default:
				panic("unreachable")
			}
		}
		fmt.Fprintln(&bufReturn)
	}
	assert(stk.Len() == 0)

	// 根据ABI处理返回值
	// 需要将栈顶位置的结果转化为本地ABI规范的返回值
	{
		// 如果走内存, 返回地址
		if len(fn.Type.Results) > 1 && fnNative.Type.Return[1].Reg == 0 {
			p.gasCommentInFunc(&bufReturn, "将返回地址复制到寄存器")
			fmt.Fprintf(&bufReturn, "    mov rax, [rbp%+d] # ret address\n",
				fnNative.Type.Return[0].RBPOff,
			)
		} else {
			p.gasCommentInFunc(&bufReturn, "将返回值变量复制到寄存器")

			// 如果走寄存器, 则复制寄存器
			for i, ret := range fnNative.Type.Return {
				if ret.Reg == 0 {
					continue
				}
				switch fn.Type.Results[i] {
				case token.I32:
					fmt.Fprintf(&bufReturn, "    mov %v, [rbp%+d] # ret.%d\n",
						x64.Reg32String(ret.Reg),
						fnNative.Type.Return[i].RBPOff,
						i,
					)
				case token.I64:
					fmt.Fprintf(&bufReturn, "    mov %v, [rbp%+d] # ret.%d\n",
						x64.RegString(ret.Reg),
						fnNative.Type.Return[i].RBPOff,
						i,
					)
				case token.F32:
					fmt.Fprintf(&bufReturn, "    movss %v, qword ptr [rbp%+d] # ret.%d\n",
						x64.RegString(ret.Reg),
						fnNative.Type.Return[i].RBPOff,
						i,
					)
				case token.F64:
					fmt.Fprintf(&bufReturn, "    movsd %v, qword ptr [rbp%+d] # ret.%d\n",
						x64.RegString(ret.Reg),
						fnNative.Type.Return[i].RBPOff,
						i,
					)
				}
			}
		}
	}
	fmt.Fprintln(&bufReturn)

	// 结束栈帧
	{
		p.gasCommentInFunc(&bufReturn, "函数返回")
		fmt.Fprintln(&bufReturn, "    mov rsp, rbp")
		fmt.Fprintln(&bufReturn, "    pop rbp")
		fmt.Fprintln(&bufReturn, "    ret")
		fmt.Fprintln(&bufReturn)
	}

	// 头部复制到 w
	io.Copy(w, &bufHeader)
	// 指令复制到 w
	io.Copy(w, &bufBody)
	// 尾部复制到 w
	io.Copy(w, &bufReturn)

	return nil
}

func (p *wat2X64Worker) buildFunc_ins(
	w io.Writer,
	fnNative *nativeast.Func, fn *ast.Func,
	stk *valueTypeStack, scopeStack *scopeContextStack,
	i ast.Instruction,
) error {
	stk.NextInstruction(i)

	p.Tracef("buildFunc_ins: %s begin: %v\n", i.Token(), stk.String())
	defer func() { p.Tracef("buildFunc_ins: %s end: %v\n", i.Token(), stk.String()) }()

	defer func() {
		if i.Token().IsTerminal() {
			assert(scopeStack.Top().IgnoreStackCheck)
		}
	}()

	switch tok := i.Token(); tok {
	case token.INS_UNREACHABLE:
		currentScopeContext := scopeStack.Top()
		currentScopeContext.IgnoreStackCheck = true

		assert(stk.Len() >= currentScopeContext.StackBase)

		// 将静态分析栈平衡
		for stk.Len() > currentScopeContext.StackBase {
			stk.DropAny()
		}
		for _, retType := range currentScopeContext.Result {
			stk.Push(retType)
		}

		// 运行时调用异常函数退出
		fmt.Fprintf(w, "    call %s # unreachable\n", kRuntimePanic)
		fmt.Fprintln(w)

	case token.INS_NOP:
		fmt.Fprintf(w, "    nop # nop\n")
		fmt.Fprintln(w)

	case token.INS_BLOCK:
		i := i.(ast.Ins_Block)

		labelName := fixName(i.Label)
		labelSuffix := p.genNextId()
		labelBlockNextId := p.makeLabelId(kLabelPrefixName_brNext, labelName, labelSuffix)

		if isSameInstList(fn.Body.List, i.List) {
			fmt.Fprintf(w, "    # fn.body.begin(name=%s, suffix=%s)\n", labelName, labelSuffix)
			fmt.Fprintln(w)
		} else {
			fmt.Fprintf(w, "    # block.begin(name=%s, suffix=%s)\n", labelName, labelSuffix)
			fmt.Fprintln(w)
		}

		// 编译块内的指令
		scopeStack.EnterScope(token.INS_BLOCK, stk.Len(), labelName, labelSuffix, i.Results)
		{
			scopeCtx := scopeStack.Top()
			scopeCtx.InstList = i.List

			for _, ins := range i.List {
				if err := p.buildFunc_ins(w, fnNative, fn, stk, scopeStack, ins); err != nil {
					return err
				}

				// 跳过后续的死代码分析
				if ins.Token().IsTerminal() {
					break
				}
			}

			if scopeCtx.IgnoreStackCheck {
				// 验证和搬运工作已经在 return 和 br 指令处完成
				// 可能跨越了多个块, 只能在调转指令处定位到目标块进行检查
				n := scopeCtx.StackBase + len(i.Results)
				assert(stk.Len() >= n)

				// 这里只需要重置栈
				for stk.Len() > scopeCtx.StackBase {
					stk.DropAny()
				}

				// 并且填入合适的块返回值, 以后后续检查继续进行
				for _, reti := range i.Results {
					stk.Push(reti)
				}

			} else {
				// 从 end 正常结束需要精确验证栈匹配
				assert(stk.Len() == scopeCtx.StackBase+len(i.Results))
				for k, expect := range i.Results {
					if got := stk.TokenAt(scopeCtx.StackBase + k); got != expect {
						panic(fmt.Sprintf("expect = %v, got = %v", expect, got))
					}
				}
			}

			// 用于 br/br-if/br-table 指令跳转
			// block 的跳转地址在结尾
			p.gasFuncLabel(w, labelBlockNextId)
		}
		scopeStack.LeaveScope()

		if isSameInstList(fn.Body.List, i.List) {
			fmt.Fprintf(w, "    # fn.body.end(name=%s, suffix=%s)\n", labelName, labelSuffix)
			fmt.Fprintln(w)
		} else {
			fmt.Fprintf(w, "    # block.end(name=%s, suffix=%s)\n", labelName, labelSuffix)
			fmt.Fprintln(w)
		}

	case token.INS_LOOP:
		i := i.(ast.Ins_Loop)

		labelName := fixName(i.Label)
		labelSuffix := p.genNextId()
		labelLoopNextId := p.makeLabelId(kLabelPrefixName_brNext, labelName, labelSuffix)

		fmt.Fprintf(w, "    # loop.begin(name=%s, suffix=%s)\n", labelName, labelSuffix)

		// 编译块内的指令
		scopeStack.EnterScope(token.INS_LOOP, stk.Len(), labelName, labelSuffix, i.Results)
		{
			scopeCtx := scopeStack.Top()

			// 用于 br/br-if/br-table 指令跳转
			// loop 的跳转地址在开头
			p.gasFuncLabel(w, labelLoopNextId)

			for _, ins := range i.List {
				if err := p.buildFunc_ins(w, fnNative, fn, stk, scopeStack, ins); err != nil {
					return err
				}

				// 跳过后续的死代码分析
				if ins.Token().IsTerminal() {
					break
				}
			}

			if scopeCtx.IgnoreStackCheck {
				// 验证和搬运工作已经在 return 和 br 指令处完成
				// 可能跨越了多个块, 只能在调转指令处定位到目标块进行检查
				n := scopeCtx.StackBase + len(i.Results)
				assert(stk.Len() >= n)

				// 这里只需要重置栈
				for stk.Len() > scopeCtx.StackBase {
					stk.DropAny()
				}

				// 并且填入合适的块返回值, 以后后续检查继续进行
				for _, reti := range i.Results {
					stk.Push(reti)
				}

			} else {
				// 从 end 正常结束需要精确验证栈匹配
				assert(stk.Len() == scopeCtx.StackBase+len(i.Results))
				for k, expect := range i.Results {
					if got := stk.TokenAt(scopeCtx.StackBase + k); got != expect {
						panic(fmt.Sprintf("expect = %v, got = %v", expect, got))
					}
				}
			}
		}
		scopeStack.LeaveScope()

		fmt.Fprintf(w, "    # loop.end(name=%s, suffix=%s)\n", labelName, labelSuffix)
		fmt.Fprintln(w)

	case token.INS_IF:
		i := i.(ast.Ins_If)

		labelName := fixName(i.Label)
		labelSuffix := p.genNextId()
		labelIfElseId := p.makeLabelId(kLabelPrefixName_else, labelName, labelSuffix)
		labelNextId := p.makeLabelId(kLabelPrefixName_brNext, labelName, labelSuffix)

		sp0 := stk.Pop(token.I32)

		fmt.Fprintf(w, "    # if.begin(name=%s, suffix=%s)\n", labelName, labelSuffix)
		fmt.Fprintf(w, "    mov eax, [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
		fmt.Fprintf(w, "    cmp eax, 0\n")
		if len(i.Else) > 0 {
			fmt.Fprintf(w, "    je  %s\n", labelIfElseId)
			fmt.Fprintln(w)
		} else {
			fmt.Fprintf(w, "    je  %s\n", labelNextId)
			fmt.Fprintln(w)
		}

		fmt.Fprintf(w, "    # if.body(name=%s, suffix=%s)\n", labelName, labelSuffix)

		// 编译 if 块内的指令
		scopeStack.EnterScope(token.INS_IF, stk.Len(), labelName, labelSuffix, i.Results)
		{
			scopeCtx := scopeStack.Top()

			for _, ins := range i.Body {
				if err := p.buildFunc_ins(w, fnNative, fn, stk, scopeStack, ins); err != nil {
					return err
				}

				// 跳过后续的死代码分析
				if ins.Token().IsTerminal() {
					break
				}
			}

			if scopeCtx.IgnoreStackCheck {
				// 验证和搬运工作已经在 return 和 br 指令处完成
				// 可能跨越了多个块, 只能在调转指令处定位到目标块进行检查
				n := scopeCtx.StackBase + len(i.Results)
				assert(stk.Len() >= n)

				// 这里只需要重置栈
				for stk.Len() > scopeCtx.StackBase {
					stk.DropAny()
				}

				// 并且填入合适的块返回值, 以后后续检查继续进行
				for _, reti := range i.Results {
					stk.Push(reti)
				}

			} else {
				// 从 end 正常结束需要精确验证栈匹配
				assert(stk.Len() == scopeCtx.StackBase+len(i.Results))
				for k, expect := range i.Results {
					if got := stk.TokenAt(scopeCtx.StackBase + k); got != expect {
						panic(fmt.Sprintf("expect = %v, got = %v", expect, got))
					}
				}
			}

			// 如果有 else 分支则要将栈上的返回值也清除
			// else 分析时会再次生成, 以确保栈的平衡
			if len(i.Else) > 0 {
				assert(stk.Len() == scopeCtx.StackBase+len(i.Results))
				for k := len(i.Results) - 1; k >= 0; k-- {
					retType := i.Results[k]
					stk.Pop(retType)
				}
			}

			fmt.Fprintf(w, "    jmp %s\n", labelNextId)
			fmt.Fprintln(w)
		}
		scopeStack.LeaveScope()

		// 编译 else 块内的指令
		scopeStack.EnterScope(token.INS_ELSE, stk.Len(), labelName, labelSuffix, i.Results)
		if len(i.Else) > 0 {
			scopeCtx := scopeStack.Top()

			// 生成 else 入口标签点
			p.gasFuncLabel(w, labelIfElseId)

			for _, ins := range i.Else {
				if err := p.buildFunc_ins(w, fnNative, fn, stk, scopeStack, ins); err != nil {
					return err
				}

				// 跳过后续的死代码分析
				if ins.Token().IsTerminal() {
					break
				}
			}

			if scopeCtx.IgnoreStackCheck {
				// 验证和搬运工作已经在 return 和 br 指令处完成
				// 可能跨越了多个块, 只能在调转指令处定位到目标块进行检查
				n := scopeCtx.StackBase + len(i.Results)
				assert(stk.Len() >= n)

				// 这里只需要重置栈
				for stk.Len() > scopeCtx.StackBase {
					stk.DropAny()
				}

				// 并且填入合适的块返回值, 以后后续检查继续进行
				for _, reti := range i.Results {
					stk.Push(reti)
				}

			} else {
				// 从 end 正常结束需要精确验证栈匹配
				assert(stk.Len() == scopeCtx.StackBase+len(i.Results))
				for k, expect := range i.Results {
					if got := stk.TokenAt(scopeCtx.StackBase + k); got != expect {
						panic(fmt.Sprintf("expect = %v, got = %v", expect, got))
					}
				}
			}

			assert(stk.Len() == scopeCtx.StackBase+len(i.Results))
		}
		scopeStack.LeaveScope()

		// 用于 br/br-if/br-table 指令跳转
		// if/else 块的跳转地址在结尾
		p.gasFuncLabel(w, labelNextId)
		fmt.Fprintf(w, "    # if.end(name=%s, suffix=%s)\n", labelName, labelSuffix)
		fmt.Fprintln(w)

	case token.INS_ELSE:
		unreachable()
	case token.INS_END:
		unreachable()

	case token.INS_BR:
		i := i.(ast.Ins_Br)

		// 设置当前 block 为非正常的 end 结束
		currentScopeContext := scopeStack.Top()
		currentScopeContext.IgnoreStackCheck = true

		// br会根据目标block的返回值个数, 从当前block产生的栈中去返回值,
		// 至于中间被跳过的block栈数据全部被丢弃.

		destScopeContex := scopeStack.FindScopeContext(i.X)
		labelNextId := p.makeLabelId(kLabelPrefixName_brNext, fixName(destScopeContex.Label), destScopeContex.LabelSuffix)

		fmt.Fprintf(w, "    # br %s\n", i.X)

		// br对应的block带返回值
		// 如果目标是 loop 则不需要处理结果, 因为还要继续循环
		if destScopeContex.Type != token.INS_LOOP && len(destScopeContex.Result) > 0 {
			// 必须确保当前block的stk上有足够的返回值
			assert(stk.Len() >= destScopeContex.StackBase+len(destScopeContex.Result))

			// 第一个返回值返回值的偏移地址
			firstResultBase := stk.Len() - len(destScopeContex.Result)

			// 验证栈上的值和返回值类型匹配
			for k, expect := range destScopeContex.Result {
				if got := stk.TokenAt(firstResultBase + k); got != expect {
					panic(fmt.Sprintf("expect = %v, got = %v", expect, got))
				}
			}

			// 如果返回值位置和目标block的base不一致则需要逐个复制
			if firstResultBase > destScopeContex.StackBase {
				// 返回值是逆序出栈
				// 注意: 这里只是搬运, 不能改变栈的状态
				for i := len(destScopeContex.Result) - 1; i >= 0; i-- {
					switch xType := destScopeContex.Result[i]; xType {
					case token.I32:
						fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
						fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
					case token.I64:
						fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
						fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
					case token.F32:
						fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
						fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
					case token.F64:
						fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
						fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
					default:
						unreachable()
					}
				}
			}
		}

		// br 指令不对栈做清理工作

		fmt.Fprintf(w, "    jmp %s\n", labelNextId)
		fmt.Fprintln(w)

	case token.INS_BR_IF:
		i := i.(ast.Ins_BrIf)

		// 设置当前 block 为非正常的 end 结束
		currentScopeContext := scopeStack.Top()
		currentScopeContext.IgnoreStackCheck = true

		destScopeContex := scopeStack.FindScopeContext(i.X)
		labelBrNextId := p.makeLabelId(kLabelPrefixName_brNext, destScopeContex.Label, destScopeContex.LabelSuffix)
		labelBrFallthroughId := p.makeLabelId(kLabelPrefixName_brFallthrough, destScopeContex.Label, destScopeContex.LabelSuffix)

		// 弹出的是条件
		sp0 := stk.Pop(token.I32)

		fmt.Fprintf(w, "    # br_if %s.%s\n", destScopeContex.Label, destScopeContex.LabelSuffix)
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
		fmt.Fprintf(w, "    cmp eax, 0\n")
		fmt.Fprintf(w, "    je  %s\n", labelBrFallthroughId)

		// 对应的block带返回值
		// 如果目标是 loop 则不需要处理结果, 因为还要继续循环
		if destScopeContex.Type != token.INS_LOOP && len(destScopeContex.Result) > 0 {
			// 必须确保当前block的stk上有足够的返回值
			assert(currentScopeContext.StackBase+len(destScopeContex.Result) >= stk.Len())

			// 第一个返回值返回值的偏移地址
			firstResultBase := stk.Len() - len(destScopeContex.Result)

			// 验证栈上的值和返回值类型匹配
			for k, expect := range destScopeContex.Result {
				if got := stk.TokenAt(firstResultBase + k); got != expect {
					panic(fmt.Sprintf("expect = %v, got = %v", expect, got))
				}
			}

			// 如果返回值位置和目标block的base不一致则需要逐个复制
			if firstResultBase > destScopeContex.StackBase {
				// 返回值是逆序出栈
				// 注意: 这里只是搬运, 不能改变栈的状态
				for i := len(destScopeContex.Result) - 1; i >= 0; i-- {
					switch xType := destScopeContex.Result[i]; xType {
					case token.I32:
						fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
						fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
					case token.I64:
						fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
						fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
					case token.F32:
						fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
						fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
					case token.F64:
						fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
						fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
					default:
						unreachable()
					}
				}
			}
		}

		// br 指令不对栈做清理工作

		fmt.Fprintf(w, "    jmp %s\n", labelBrNextId)
		p.gasFuncLabel(w, labelBrFallthroughId)
		fmt.Fprintln(w)

	case token.INS_BR_TABLE:
		i := i.(ast.Ins_BrTable)
		assert(len(i.XList) > 1)

		// 设置当前 block 为非正常的 end 结束
		currentScopeContext := scopeStack.Top()
		currentScopeContext.IgnoreStackCheck = true

		// br-table的行为和br比较相似, 因此不涉及else部分不用担心栈平衡的问题.
		// 但是每个目标block的返回值必须完全一致

		labelSuffix := p.genNextId()

		// 默认分支
		defaultScopeContex := scopeStack.FindScopeContext(i.XList[len(i.XList)-1])
		defaultLabelNextId := p.makeLabelId(kLabelPrefixName_brNext, fixName(defaultScopeContex.Label), defaultScopeContex.LabelSuffix)

		sp0 := stk.Pop(token.I32)

		fmt.Fprintf(w, "    # br_table\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)

		// 生成跳转链
		for k := 0; k < len(i.XList); k++ {
			if k < len(i.XList)-1 {
				fmt.Fprintf(w, "    # br_table case %d\n", k)
				fmt.Fprintf(w, "    cmp eax, %d\n", k)
				fmt.Fprintf(w, "    je  %s\n", p.makeLabelId(kLabelPrefixName_brCase, i.XList[k], labelSuffix))
			} else {
				fmt.Fprintf(w, "    # br_table default\n")
				fmt.Fprintf(w, "    jmp %s\n", p.makeLabelId(kLabelPrefixName_brDefault, "", labelSuffix))
			}
		}

		// 定义每个分支的跳转代码
		{
			// 当前block的返回值位置是相同的, 只能统一取一次
			var retIdxList = make([]int, len(defaultScopeContex.Result))
			for k := len(defaultScopeContex.Result) - 1; k >= 0; k-- {
				xTyp := defaultScopeContex.Result[k]
				retIdxList[k] = stk.Pop(xTyp)
			}

			for k := 0; k < len(i.XList); k++ {
				destScopeContex := scopeStack.FindScopeContext(i.XList[k])
				labelNextId := p.makeLabelId(kLabelPrefixName_brNext, fixName(destScopeContex.Label), labelSuffix)

				// 验证每个目标返回值必须和default一致
				assert(len(defaultScopeContex.Result) == len(destScopeContex.Result))
				for i := 0; i < len(defaultScopeContex.Result); i++ {
					assert(defaultScopeContex.Result[i] == destScopeContex.Result[i])
				}

				// 生成跳转的标签
				if k < len(i.XList)-1 {
					p.gasFuncLabel(w, p.makeLabelId(kLabelPrefixName_brCase, i.XList[k], labelSuffix))
				} else {
					p.gasFuncLabel(w, p.makeLabelId(kLabelPrefixName_brDefault, "", labelSuffix))
				}

				// 带返回值的情况
				if len(destScopeContex.Result) > 0 {
					// 必须确保当前block的stk上有足够的返回值
					assert(stk.Len() >= destScopeContex.StackBase+len(destScopeContex.Result))

					// 第一个返回值返回值的偏移地址
					firstResultBase := stk.Len() - len(destScopeContex.Result)

					// 验证栈上的值和返回值类型匹配
					for k, expect := range destScopeContex.Result {
						if got := stk.TokenAt(firstResultBase + k); got != expect {
							panic(fmt.Sprintf("expect = %v, got = %v", expect, got))
						}
					}

					// 如果返回值位置和目标block的base不一致则需要逐个复制
					if firstResultBase > destScopeContex.StackBase {
						// 返回值是逆序出栈
						for i := 0; i < len(destScopeContex.Result); i++ {
							switch xType := destScopeContex.Result[i]; xType {
							case token.I32:
								fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
								fmt.Fprintf(w, "    mov dword ptr [rbp%+d], r10d\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
							case token.I64:
								fmt.Fprintf(w, "    mov r10, qword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
								fmt.Fprintf(w, "    mov qword ptr [rbp%+d], r10\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
							case token.F32:
								fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
								fmt.Fprintf(w, "    mov dword ptr [rbp%+d], r10d\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
							case token.F64:
								fmt.Fprintf(w, "    mov r10, qword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
								fmt.Fprintf(w, "    mov qword ptr [rbp%+d], r10\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
							default:
								unreachable()
							}
						}
					}
				}

				// 跳转到目标
				if k < len(i.XList)-1 {
					fmt.Fprintf(w, "    jmp %s\n", labelNextId)
				} else {
					fmt.Fprintf(w, "    jmp %s\n", defaultLabelNextId)
				}
			}
		}

		// br 指令不对栈做清理工作

		fmt.Fprintf(w, "    # br_table end\n")
		fmt.Fprintln(w)

	case token.INS_RETURN:
		// 设置当前 block 为非正常的 end 结束
		currentScopeContext := scopeStack.Top()
		currentScopeContext.IgnoreStackCheck = true

		// return 是 br 指令的语法糖
		// br会根据目标block的返回值个数, 从当前block产生的栈中去返回值,
		// 至于中间被跳过的block栈数据全部被丢弃.

		destScopeContex := scopeStack.GetReturnScopeContext()
		assert(isSameInstList(fn.Body.List, destScopeContex.InstList))

		labelNextId := p.makeLabelId(kLabelPrefixName_brNext, fixName(destScopeContex.Label), destScopeContex.LabelSuffix)

		p.gasCommentInFunc(w, "return")

		// br对应的block带返回值
		// 如果目标是 loop 则不需要处理结果, 因为还要继续循环
		if destScopeContex.Type != token.INS_LOOP && len(destScopeContex.Result) > 0 {
			// 必须确保当前block的stk上有足够的返回值
			assert(stk.Len() >= destScopeContex.StackBase+len(destScopeContex.Result))

			// 第一个返回值返回值的偏移地址
			firstResultBase := stk.Len() - len(destScopeContex.Result)

			// 验证栈上的值和返回值类型匹配
			for k, expect := range destScopeContex.Result {
				if got := stk.TokenAt(firstResultBase + k); got != expect {
					panic(fmt.Sprintf("expect = %v, got = %v", expect, got))
				}
			}

			// 如果返回值位置和目标block的base不一致则需要逐个复制
			if firstResultBase > destScopeContex.StackBase {
				// 返回值是逆序出栈
				// 注意: 这里只是搬运, 不能改变栈的状态
				for i := len(destScopeContex.Result) - 1; i >= 0; i-- {
					switch xType := destScopeContex.Result[i]; xType {
					case token.I32:
						fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
						fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
					case token.I64:
						fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
						fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
					case token.F32:
						fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
						fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
					case token.F64:
						fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d] # copy result\n", p.fnWasmR0Base-(firstResultBase+i)*8-8)
						fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", p.fnWasmR0Base-(destScopeContex.StackBase+i)*8-8)
					default:
						unreachable()
					}
				}
			}
		}

		// return 只是跳转到 body 块的结尾, 清理和返回值在后面处理
		fmt.Fprintf(w, "    jmp %s\n", labelNextId)
		fmt.Fprintln(w)

	case token.INS_CALL:
		i := i.(ast.Ins_Call)

		fnCallType := p.findFuncType(i.X)
		fnCallIdx := p.findFuncIndex(i.X)

		// 构建被调用函数的栈帧信息
		fnCallNative := wat2nativeFunc(i.X, fnCallType, nil)

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
		if len(fnCallNative.Type.Return) > 1 && fnCallNative.Type.Return[1].Reg == 0 {
			if p.cpuType == abi.X64Unix {
				fmt.Fprintf(w, "    lea rdi, [rsp%+d] # return address\n",
					fnCallNative.ArgsSize-len(fnCallNative.Type.Return)*8,
				)
			} else {
				fmt.Fprintf(w, "    lea rcx, [rsp%+d] # return address\n",
					fnCallNative.ArgsSize-len(fnCallNative.Type.Return)*8,
				)
			}
		}

		// 准备调用参数
		for k, x := range fnCallType.Params {
			switch x.Type {
			case token.I32:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    mov %s, dword ptr [rbp%+d] # arg %d\n",
						x64.Reg32String(arg.Reg),
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
					fmt.Fprintf(w, "    movss %s, dword ptr [rbp%+d] # arg %d\n",
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

		if fnCallIdx < len(p.m.Imports) {
			fmt.Fprintf(w, "    call %s\n", kImportNamePrefix+fixName(i.X))
		} else {
			fmt.Fprintf(w, "    call %s\n", kFuncNamePrefix+fixName(i.X))
		}

		// 提取返回值
		if len(fnCallNative.Type.Return) > 1 && fnCallNative.Type.Return[1].Reg == 0 {
			// 走栈返回结果, 地址通过 rcx 传入, 通过 rax 返回
			// 根据 ABI 规范获取返回值, 再保存到 WASM 栈中
			for k, retType := range fnCallType.Results {
				reti := stk.Push(retType)
				switch retType {
				case token.I32:
					fmt.Fprintf(w, "    mov r10d, dword ptr [rax%+d]\n", k*8)
					fmt.Fprintf(w, "    mov dword ptr [rbp%+d], r10d\n", p.fnWasmR0Base-reti*8-8)
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
						x64.Reg32String(retNative.Reg),
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
		fmt.Fprintln(w)

	case token.INS_CALL_INDIRECT:
		i := i.(ast.Ins_CallIndirect)

		fnCallType := p.findType(i.TypeIdx)

		// 构建被调用函数的栈帧信息
		// 间接调用, 没有名字(可以尝试根据地址查询名字)
		fnCallNative := wat2nativeFunc("", fnCallType, nil)

		// 模拟构建函数栈帧
		// 后面需要根据参数是否走寄存器传递生成相关的入口代码和返回代码
		if err := astutil.BuildFuncFrame(p.cpuType, fnCallNative); err != nil {
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
			if p.cpuType == abi.X64Unix {
				fmt.Fprintf(w, "    lea rdi, [rsp%+d] # return address\n",
					fnCallNative.ArgsSize-len(fnCallNative.Type.Return)*8,
				)
			} else {
				fmt.Fprintf(w, "    lea rcx, [rsp%+d] # return address\n",
					fnCallNative.ArgsSize-len(fnCallNative.Type.Return)*8,
				)
			}
		}

		// 准备调用参数
		for k, x := range fnCallType.Params {
			switch x.Type {
			case token.I32:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    mov %s, dword ptr [rbp%+d] # arg %d\n",
						x64.Reg32String(arg.Reg),
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
						p.fnWasmR0Base+argList[k]*8-8,
					)
					fmt.Fprintf(w, "    mov qword ptr [rsp%+d], rax\n",
						arg.RSPOff,
					)
				}
			case token.F32:
				if arg := fnCallNative.Type.Args[k]; arg.Reg != 0 {
					fmt.Fprintf(w, "    movss %s, dword ptr [rbp%+d] # arg %d\n",
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
		fmt.Fprintf(w, "    call %s\n", r11)

		// 提取返回值
		if len(fnCallNative.Type.Return) > 1 && fnCallNative.Type.Return[1].Reg == 0 {
			// 走栈返回结果, 地址通过 rcx 传入, 通过 rax 返回
			// 根据 ABI 规范获取返回值, 再保存到 WASM 栈中
			for k, retType := range fnCallType.Results {
				reti := stk.Push(retType)
				switch retType {
				case token.I32:
					fmt.Fprintf(w, "    mov r10d, dword ptr [rax%+d]\n", k*8)
					fmt.Fprintf(w, "    mov dword ptr [rbp%+d], r10d\n", p.fnWasmR0Base-reti*8-8)
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
						x64.Reg32String(retNative.Reg),
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
		fmt.Fprintln(w)

	case token.INS_DROP:
		sp0 := stk.DropAny()
		fmt.Fprintf(w, "    nop # drop [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
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

		p.gasCommentInFunc(w, "select")
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-spCondition*8-8)
		fmt.Fprintf(w, "    test eax, eax\n")

		switch valType {
		case token.I32:
			fmt.Fprintf(w, "    mov    r10d, dword ptr [rbp%+d]\n", p.fnWasmR0Base-spValueFalse*8-8)
			fmt.Fprintf(w, "    mov    r11d, dword ptr [rbp%+d]\n", p.fnWasmR0Base-spValueTrue*8-8)
			fmt.Fprintf(w, "    cmovne r10d, r11d\n")
			fmt.Fprintf(w, "    mov    dword ptr [rbp%+d], r10d", p.fnWasmR0Base-ret0*8-8)
		case token.I64:
			fmt.Fprintf(w, "    mov    r10, qword ptr [rbp%+d]\n", p.fnWasmR0Base-spValueFalse*8-8)
			fmt.Fprintf(w, "    mov    r11, qword ptr [rbp%+d]\n", p.fnWasmR0Base-spValueTrue*8-8)
			fmt.Fprintf(w, "    cmovne r10, r11\n")
			fmt.Fprintf(w, "    mov    qword ptr [rbp%+d], r10", p.fnWasmR0Base-ret0*8-8)
		case token.F32:
			fmt.Fprintf(w, "    mov    r10d, dword ptr [rbp%+d]\n", p.fnWasmR0Base-spValueFalse*8-8)
			fmt.Fprintf(w, "    mov    r11d, dword ptr [rbp%+d]\n", p.fnWasmR0Base-spValueTrue*8-8)
			fmt.Fprintf(w, "    cmovne r10d, r11d\n")
			fmt.Fprintf(w, "    mov    dword ptr [rbp%+d], r10d", p.fnWasmR0Base-ret0*8-8)
		case token.F64:
			fmt.Fprintf(w, "    mov    r10, qword ptr [rbp%+d]\n", p.fnWasmR0Base-spValueFalse*8-8)
			fmt.Fprintf(w, "    mov    r11, qword ptr [rbp%+d]\n", p.fnWasmR0Base-spValueTrue*8-8)
			fmt.Fprintf(w, "    cmovne r10, r11\n")
			fmt.Fprintf(w, "    mov    qword ptr [rbp%+d], r10", p.fnWasmR0Base-ret0*8-8)
		default:
			unreachable()
		}
		fmt.Fprintln(w)

	case token.INS_LOCAL_GET:
		i := i.(ast.Ins_LocalGet)
		xType := p.findLocalType(fn, i.X)
		xOff := p.findLocalOffset(fnNative, fn, i.X)
		ret0 := stk.Push(xType)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "    # local.get %s i32\n", i.X)
			fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", xOff)
			fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", p.fnWasmR0Base-ret0*8-8)
			fmt.Fprintln(w)
		case token.I64:
			fmt.Fprintf(w, "    # local.get %s i64\n", i.X)
			fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", xOff)
			fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", p.fnWasmR0Base-ret0*8-8)
			fmt.Fprintln(w)
		case token.F32:
			fmt.Fprintf(w, "    # local.get %s f32\n", i.X)
			fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", xOff)
			fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", p.fnWasmR0Base-ret0*8-8)
			fmt.Fprintln(w)
		case token.F64:
			fmt.Fprintf(w, "    # local.get %s f64\n", i.X)
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
			fmt.Fprintf(w, "    # local.set %s i32\n", i.X)
			fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", xOff)
			fmt.Fprintln(w)
		case token.I64:
			fmt.Fprintf(w, "    # local.set %s i64\n", i.X)
			fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", xOff)
			fmt.Fprintln(w)
		case token.F32:
			fmt.Fprintf(w, "    # local.set %s f32\n", i.X)
			fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", xOff)
			fmt.Fprintln(w)
		case token.F64:
			fmt.Fprintf(w, "    # local.set %s f64\n", i.X)
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
			fmt.Fprintf(w, "    # local.tee %s i32\n", i.X)
			fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", xOff)
		case token.I64:
			fmt.Fprintf(w, "    # local.tee %s i64\n", i.X)
			fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", xOff)
		case token.F32:
			fmt.Fprintf(w, "    # local.tee %s f32\n", i.X)
			fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", xOff)
		case token.F64:
			fmt.Fprintf(w, "    # local.tee %s f64\n", i.X)
			fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", xOff)
		default:
			unreachable()
		}
		fmt.Fprintln(w)

	case token.INS_GLOBAL_GET:
		i := i.(ast.Ins_GlobalGet)
		xType := p.findGlobalType(i.X)
		ret0 := stk.Push(xType)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "    # global.get %s i32\n", i.X)
			fmt.Fprintf(w, "    mov eax, dword ptr [rip+%s]\n", kGlobalNamePrefix+p.findGlobalName(i.X))
			fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", p.fnWasmR0Base-ret0*8-8)
		case token.I64:
			fmt.Fprintf(w, "    # global.get %s i64\n", i.X)
			fmt.Fprintf(w, "    mov rax, qword ptr [rip+%s]\n", kGlobalNamePrefix+p.findGlobalName(i.X))
			fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", p.fnWasmR0Base-ret0*8-8)
		case token.F32:
			fmt.Fprintf(w, "    # global.get %s f32\n", i.X)
			fmt.Fprintf(w, "    mov eax, dword ptr [rip+%s]\n", kGlobalNamePrefix+p.findGlobalName(i.X))
			fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", p.fnWasmR0Base-ret0*8-8)
		case token.F64:
			fmt.Fprintf(w, "    # global.get %s f64\n", i.X)
			fmt.Fprintf(w, "    mov rax, qword ptr [rip+%s]\n", kGlobalNamePrefix+p.findGlobalName(i.X))
			fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", p.fnWasmR0Base-ret0*8-8)
		default:
			unreachable()
		}
		fmt.Fprintln(w)

	case token.INS_GLOBAL_SET:
		i := i.(ast.Ins_GlobalSet)
		xType := p.findGlobalType(i.X)
		sp0 := stk.Pop(xType)
		switch xType {
		case token.I32:
			fmt.Fprintf(w, "    # global.set %s i32\n", i.X)
			fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    mov dword ptr [rip+%s], eax\n", kGlobalNamePrefix+p.findGlobalName(i.X))
		case token.I64:
			fmt.Fprintf(w, "    # global.set %s i64\n", i.X)
			fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    mov qword ptr [rip+%s], rax\n", kGlobalNamePrefix+p.findGlobalName(i.X))
		case token.F32:
			fmt.Fprintf(w, "    # global.set %s f32\n", i.X)
			fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    mov dword ptr [rip+%s], eax\n", kGlobalNamePrefix+p.findGlobalName(i.X))
		case token.F64:
			fmt.Fprintf(w, "    # global.set %s f64\n", i.X)
			fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", p.fnWasmR0Base-sp0*8-8)
			fmt.Fprintf(w, "    mov qword ptr [rip+%s], rax\n", kGlobalNamePrefix+p.findGlobalName(i.X))
		default:
			unreachable()
		}
		fmt.Fprintln(w)

	case token.INS_TABLE_GET:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.FUNCREF) - 8 // funcref
		fmt.Fprintf(w, "    # table.get\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rip + %s]\n", kTableAddrName)
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov rax, dword ptr [rax+r10*8]\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_TABLE_SET:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.FUNCREF) - 8 // funcref
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		fmt.Fprintf(w, "    # table.set\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rip + %s]\n", kTableAddrName)
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov r11d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov dword ptr [rax+r11*8], r10\n")
		fmt.Fprintln(w)

	case token.INS_I32_LOAD:
		i := i.(ast.Ins_I32Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.load\n")
		fmt.Fprintf(w, "    mov rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [r10%+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LOAD:
		i := i.(ast.Ins_I64Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.load\n")
		fmt.Fprintf(w, "    mov rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [r10%+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_LOAD:
		i := i.(ast.Ins_F32Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.load\n")
		fmt.Fprintf(w, "    mov   rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add   r10, rax\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [r10%+d]\n", i.Offset)
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_LOAD:
		i := i.(ast.Ins_F64Load)

		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.load\n")
		fmt.Fprintf(w, "    mov   rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add   r10, rax\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [r10%+d]\n", i.Offset)
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LOAD8_S:
		i := i.(ast.Ins_I32Load8S)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.load8_s\n")
		fmt.Fprintf(w, "    mov   rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add   r10, rax\n")
		fmt.Fprintf(w, "    movsx eax, byte ptr [r10%+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LOAD8_U:
		i := i.(ast.Ins_I32Load8U)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.load8_u\n")
		fmt.Fprintf(w, "    mov   rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add   r10, rax\n")
		fmt.Fprintf(w, "    movzx eax, byte ptr [r10%+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LOAD16_S:
		i := i.(ast.Ins_I32Load16S)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.load16_s\n")
		fmt.Fprintf(w, "    mov   rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add   r10, rax\n")
		fmt.Fprintf(w, "    movsx eax, word ptr [r10%+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LOAD16_U:
		i := i.(ast.Ins_I32Load16U)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.load16_u\n")
		fmt.Fprintf(w, "    mov   rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add   r10, rax\n")
		fmt.Fprintf(w, "    movzx eax, word ptr [r10%+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LOAD8_S:
		i := i.(ast.Ins_I64Load8S)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.load8_s\n")
		fmt.Fprintf(w, "    mov   rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add   r10, rax\n")
		fmt.Fprintf(w, "    movsx rax, byte ptr [r10%+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov   qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LOAD8_U:
		i := i.(ast.Ins_I64Load8U)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.load8_u\n")
		fmt.Fprintf(w, "    mov   rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add   r10, rax\n")
		fmt.Fprintf(w, "    movzx rax, byte ptr [r10%+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov   qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LOAD16_S:
		i := i.(ast.Ins_I64Load16S)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.load16_s\n")
		fmt.Fprintf(w, "    mov   rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add   r10, rax\n")
		fmt.Fprintf(w, "    movsx rax, word ptr [r10%+d]\n", i.Offset) //
		fmt.Fprintf(w, "    mov   qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LOAD16_U:
		i := i.(ast.Ins_I64Load16U)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.load16_u\n")
		fmt.Fprintf(w, "    mov   rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add   r10, rax\n")
		fmt.Fprintf(w, "    movzx rax, word ptr [r10%+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov   qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LOAD32_S:
		i := i.(ast.Ins_I64Load32S)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.load32_s\n")
		fmt.Fprintf(w, "    mov    rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov    r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add    r10, rax\n")
		fmt.Fprintf(w, "    movsxd rax, dword ptr [r10%+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov    qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LOAD32_U:
		i := i.(ast.Ins_I64Load32U)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.load32_u\n")
		fmt.Fprintf(w, "    mov rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [r10%+d]\n", i.Offset)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_STORE:
		i := i.(ast.Ins_I32Store)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # i32.store\n")
		fmt.Fprintf(w, "    mov rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov dword ptr [r10%+d], eax\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_I64_STORE:
		i := i.(ast.Ins_I64Store)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # i64.store\n")
		fmt.Fprintf(w, "    mov rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov qword ptr [r10%+d], rax\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_F32_STORE:
		i := i.(ast.Ins_F32Store)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # f32.store\n")
		fmt.Fprintf(w, "    mov   rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add   r10, rax\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss dword ptr [r10%+d], xmm4\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_F64_STORE:
		i := i.(ast.Ins_F64Store)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # f64.store\n")
		fmt.Fprintf(w, "    mov   rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add   r10, rax\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd qword ptr [r10%+d], xmm4\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_I32_STORE8:
		i := i.(ast.Ins_I32Store8)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # i32.store8\n")
		fmt.Fprintf(w, "    mov rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov byte ptr [r10%+d], al\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_I32_STORE16:
		i := i.(ast.Ins_I32Store16)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # i64.store16\n")
		fmt.Fprintf(w, "    mov rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov word ptr [r10%+d], ax\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_I64_STORE8:
		i := i.(ast.Ins_I64Store8)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # i64.store8\n")
		fmt.Fprintf(w, "    mov rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov byte ptr [r10%+d], al\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_I64_STORE16:
		i := i.(ast.Ins_I64Store16)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # i64.store16\n")
		fmt.Fprintf(w, "    mov rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov word ptr [r10%+d], ax\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_I64_STORE32:
		i := i.(ast.Ins_I64Store32)
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		fmt.Fprintf(w, "    # i64.store32\n")
		fmt.Fprintf(w, "    mov rax, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add r10, rax\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov dword ptr [r10%+d], eax\n", i.Offset)
		fmt.Fprintln(w)

	case token.INS_MEMORY_SIZE:
		sp0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8
		fmt.Fprintf(w, "    # memory.size\n")
		fmt.Fprintf(w, "    mov rax, [rip+%s]\n", kMemoryPagesName)
		fmt.Fprintf(w, "    mov [rbp%+d], rax\n", sp0)
		fmt.Fprintln(w)

	case token.INS_MEMORY_GROW:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		labelSuffix := p.genNextId()

		labelElse := p.makeLabelId(kLabelPrefixName_else, "", labelSuffix)
		labelEnd := p.makeLabelId(kLabelPrefixName_end, "", labelSuffix)

		fmt.Fprintf(w, "    # memory.grow\n")
		fmt.Fprintf(w, "    mov r10, [rip+%s]\n", kMemoryPagesName)
		fmt.Fprintf(w, "    mov r11, [rip+%s]\n", kMemoryMaxPagesName)
		fmt.Fprintf(w, "    mov rax, [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    add rax, r10\n")
		fmt.Fprintf(w, "    cmp rax, r11\n")
		fmt.Fprintf(w, "    ja  %s\n", labelElse)

		fmt.Fprintf(w, "    mov [rip+%s], rax\n", kMemoryPagesName)
		fmt.Fprintf(w, "    mov [rbp%+d], r10\n", ret0)
		fmt.Fprintf(w, "    jmp %s\n", labelEnd)

		p.gasFuncLabel(w, labelElse)
		fmt.Fprintf(w, "    mov rax, -1\n")
		fmt.Fprintf(w, "    mov [rbp%+d], rax\n", ret0)

		p.gasFuncLabel(w, labelEnd)
		fmt.Fprintln(w)

	case token.INS_MEMORY_INIT:
		i := i.(ast.Ins_MemoryInit)

		len := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		off := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		dst := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		// 参数寄存器
		regArg0 := "rcx"
		regArg1 := "rdx"
		regArg2 := "r8"

		if p.cpuType == abi.X64Unix {
			regArg0 = "rdi"
			regArg1 = "rsi"
			regArg2 = "rdx"
		}

		fmt.Fprintf(w, "    # memory.init")
		fmt.Fprintf(w, "    mov  %s, [rip + %s]\n", regArg0, kMemoryAddrName)
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", dst)
		fmt.Fprintf(w, "    add  %s, rax\n", regArg0)
		fmt.Fprintf(w, "    lea  %s, [rip + %s%d]\n", regArg1, kMemoryDataPtrPrefix, i.DataIdx)
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", off)
		fmt.Fprintf(w, "    add  %s, rax\n", regArg1)
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", len)
		fmt.Fprintf(w, "    mov  %s, rax\n", regArg2)
		fmt.Fprintf(w, "    call %s\n", kRuntimeMemcpy)
		fmt.Fprintln(w)

	case token.INS_MEMORY_COPY:
		len := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		src := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		dst := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		// 参数寄存器
		regArg0 := "rcx"
		regArg1 := "rdx"
		regArg2 := "r8"

		if p.cpuType == abi.X64Unix {
			regArg0 = "rdi"
			regArg1 = "rsi"
			regArg2 = "rdx"
		}

		fmt.Fprintf(w, "    # memory.copy\n")
		fmt.Fprintf(w, "    mov  %s, [rip + %s]\n", regArg0, kMemoryAddrName)
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", dst)
		fmt.Fprintf(w, "    add  %s, rax\n", regArg0)
		fmt.Fprintf(w, "    mov  %s, [rip + %s]\n", regArg1, kMemoryAddrName)
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", src)
		fmt.Fprintf(w, "    add  %s, rax\n", regArg1)
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", len)
		fmt.Fprintf(w, "    mov  %s, rax\n", regArg2)
		fmt.Fprintf(w, "    call %s\n", kRuntimeMemmove)
		fmt.Fprintln(w)

	case token.INS_MEMORY_FILL:
		len := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		val := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		dst := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8

		// 参数寄存器
		regArg0 := "rcx"
		regArg1 := "rdx"
		regArg2 := "r8"

		if p.cpuType == abi.X64Unix {
			regArg0 = "rdi"
			regArg1 = "rsi"
			regArg2 = "rdx"
		}

		fmt.Fprintf(w, "    # memory.fill\n")
		fmt.Fprintf(w, "    mov  %s, [rip + %s]\n", regArg0, kMemoryAddrName)
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", dst)
		fmt.Fprintf(w, "    add  %s, rax\n", regArg0)
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", val)
		fmt.Fprintf(w, "    mov  %s, rax\n", regArg1)
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", len)
		fmt.Fprintf(w, "    mov  %s, rax\n", regArg2)
		fmt.Fprintf(w, "    call %s\n", kRuntimeMemset)
		fmt.Fprintln(w)

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
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.eqz\n")
		fmt.Fprintf(w, "    mov   eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   eax, 0  # (eax==0)?\n")
		fmt.Fprintf(w, "    sete  al      # al = (eax==0)? 1: 0\n")
		fmt.Fprintf(w, "    movzx eax, al # eax = al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_EQ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.eq\n")
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10d, r11d\n")
		fmt.Fprintf(w, "    sete  al      # al = (r10d==r11d)? 1: 0\n")
		fmt.Fprintf(w, "    movzx eax, al # eax = al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_NE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.ne\n")
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10d, r11d\n")
		fmt.Fprintf(w, "    setne al      # al = (r10d==r11d)? 1: 0\n")
		fmt.Fprintf(w, "    movzx eax, al # eax = al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LT_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.lt_s\n")
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10d, r11d\n")
		fmt.Fprintf(w, "    setl  al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LT_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.lt_u\n")
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10d, r11d\n")
		fmt.Fprintf(w, "    setb  al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_GT_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.gt_s\n")
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10d, r11d\n")
		fmt.Fprintf(w, "    setg  al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_GT_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.gt_u\n")
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10d, r11d\n")
		fmt.Fprintf(w, "    seta  al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LE_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.le_s\n")
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10d, r11d\n")
		fmt.Fprintf(w, "    setle al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_LE_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.le_u\n")
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10d, r11d\n")
		fmt.Fprintf(w, "    setbe al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_GE_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.ge_s\n")
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10d, r11d\n")
		fmt.Fprintf(w, "    setge al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_GE_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.ge_u\n")
		fmt.Fprintf(w, "    mov   r10d, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10d, r11d\n")
		fmt.Fprintf(w, "    setae al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_EQZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.eqz\n")
		fmt.Fprintf(w, "    mov   rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   rax, 0  # (rax==0)?\n")
		fmt.Fprintf(w, "    sete  al      # al = (rax==0)? 1: 0\n")
		fmt.Fprintf(w, "    movzx eax, al # eax = al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_EQ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.eq\n")
		fmt.Fprintf(w, "    mov   r10, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10, r11\n")
		fmt.Fprintf(w, "    sete  al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_NE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.ne\n")
		fmt.Fprintf(w, "    mov   r10, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10, r11\n")
		fmt.Fprintf(w, "    setne al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LT_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.lt_s\n")
		fmt.Fprintf(w, "    mov   r10, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10, r11\n")
		fmt.Fprintf(w, "    setl  al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LT_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.lt_u\n")
		fmt.Fprintf(w, "    mov   r10, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10, r11\n")
		fmt.Fprintf(w, "    setb  al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_GT_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.gt_s\n")
		fmt.Fprintf(w, "    mov   r10, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10, r11\n")
		fmt.Fprintf(w, "    setg  al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_GT_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.gt_u\n")
		fmt.Fprintf(w, "    mov   r10, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10, r11\n")
		fmt.Fprintf(w, "    seta  al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LE_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.le_s\n")
		fmt.Fprintf(w, "    mov   r10, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10, r11\n")
		fmt.Fprintf(w, "    setle al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_LE_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.le_u\n")
		fmt.Fprintf(w, "    mov   r10, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10, r11\n")
		fmt.Fprintf(w, "    setbe al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_GE_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.ge_s\n")
		fmt.Fprintf(w, "    mov   r10, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10, r11\n")
		fmt.Fprintf(w, "    setge al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_GE_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i64.ge_u\n")
		fmt.Fprintf(w, "    mov   r10, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov   r11, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cmp   r10, r11\n")
		fmt.Fprintf(w, "    setae al\n")
		fmt.Fprintf(w, "    movzx eax, al\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_EQ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f32.eq\n")
		fmt.Fprintf(w, "    movss   xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movss   xmm5, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    sete    al\n")
		fmt.Fprintf(w, "    setnp   cl # set if not NaN\n")
		fmt.Fprintf(w, "    and     al, cl\n")
		fmt.Fprintf(w, "    movzx   eax, al\n")
		fmt.Fprintf(w, "    mov     dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_NE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f32.ne\n")
		fmt.Fprintf(w, "    movss   xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movss   xmm5, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    setne   al\n")
		fmt.Fprintf(w, "    setp    cl # set if NaN\n")
		fmt.Fprintf(w, "    or      al, cl\n")
		fmt.Fprintf(w, "    movzx   eax, al\n")
		fmt.Fprintf(w, "    mov     dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_LT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f32.lt\n")
		fmt.Fprintf(w, "    movss   xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movss   xmm5, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    setb    al\n")
		fmt.Fprintf(w, "    setnp   cl # set if not NaN\n")
		fmt.Fprintf(w, "    and al, cl\n")
		fmt.Fprintf(w, "    movzx   eax, al\n")
		fmt.Fprintf(w, "    mov     dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_GT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f32.gt\n")
		fmt.Fprintf(w, "    movss   xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movss   xmm5, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    seta    al\n")
		fmt.Fprintf(w, "    setnp   cl # set if not NaN\n")
		fmt.Fprintf(w, "    and     al, cl\n")
		fmt.Fprintf(w, "    movzx   eax, al\n")
		fmt.Fprintf(w, "    mov     dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_LE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f32.le\n")
		fmt.Fprintf(w, "    movss   xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movss   xmm5, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    setbe   al\n")
		fmt.Fprintf(w, "    setnp   cl # set if not NaN\n")
		fmt.Fprintf(w, "    and     al, cl\n")
		fmt.Fprintf(w, "    movzx   eax, al\n")
		fmt.Fprintf(w, "    mov     dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_GE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f32.ge\n")
		fmt.Fprintf(w, "    movss   xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movss   xmm5, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    ucomiss xmm4, xmm5\n")
		fmt.Fprintf(w, "    setae   al\n")
		fmt.Fprintf(w, "    setnp   cl # set if not NaN\n")
		fmt.Fprintf(w, "    and     al, cl\n")
		fmt.Fprintf(w, "    movzx   eax, al\n")
		fmt.Fprintf(w, "    mov     dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_EQ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f64.eq\n")
		fmt.Fprintf(w, "    movsd   xmm4, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movsd   xmm5, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    ucomisd xmm4, xmm5\n")
		fmt.Fprintf(w, "    sete    al\n")
		fmt.Fprintf(w, "    setnp   cl # set if not NaN\n")
		fmt.Fprintf(w, "    and al, cl\n")
		fmt.Fprintf(w, "    movzx   eax, al\n")
		fmt.Fprintf(w, "    mov     dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_NE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f64.ne\n")
		fmt.Fprintf(w, "    movsd   xmm4, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movsd   xmm5, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    ucomisd xmm4, xmm5\n")
		fmt.Fprintf(w, "    setne   al\n")
		fmt.Fprintf(w, "    setp    cl # set if NaN\n")
		fmt.Fprintf(w, "    or      al, cl\n")
		fmt.Fprintf(w, "    movzx   eax, al\n")
		fmt.Fprintf(w, "    mov     dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_LT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f64.lt\n")
		fmt.Fprintf(w, "    movsd   xmm4, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movsd   xmm5, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    ucomisd xmm4, xmm5\n")
		fmt.Fprintf(w, "    setb    al\n")
		fmt.Fprintf(w, "    setnp   cl # set if not NaN\n")
		fmt.Fprintf(w, "    and     al, cl\n")
		fmt.Fprintf(w, "    movzx   eax, al\n")
		fmt.Fprintf(w, "    mov     dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_GT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f64.gt\n")
		fmt.Fprintf(w, "    movsd   xmm4, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movsd   xmm5, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    ucomisd xmm4, xmm5\n")
		fmt.Fprintf(w, "    seta    al\n")
		fmt.Fprintf(w, "    setnp   cl # set if not NaN\n")
		fmt.Fprintf(w, "    and     al, cl\n")
		fmt.Fprintf(w, "    movzx   eax, al\n")
		fmt.Fprintf(w, "    mov     dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_LE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f64.le\n")
		fmt.Fprintf(w, "    movsd   xmm4, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movsd   xmm5, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    ucomisd xmm4, xmm5\n")
		fmt.Fprintf(w, "    setbe   al\n")
		fmt.Fprintf(w, "    setnp   cl # set if not NaN\n")
		fmt.Fprintf(w, "    and     al, cl\n")
		fmt.Fprintf(w, "    movzx   eax, al\n")
		fmt.Fprintf(w, "    mov     dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_GE:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # f64.ge\n")
		fmt.Fprintf(w, "    movsd   xmm4, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    movsd   xmm5, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    ucomisd xmm4, xmm5\n")
		fmt.Fprintf(w, "    setae   al\n")
		fmt.Fprintf(w, "    setnp   cl # set if not NaN\n")
		fmt.Fprintf(w, "    and     al, cl\n")
		fmt.Fprintf(w, "    movzx   eax, al\n")
		fmt.Fprintf(w, "    mov     dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_CLZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.clz\n")
		fmt.Fprintf(w, "    mov   eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    lzcnt eax, eax\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_CTZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.ctz\n")
		fmt.Fprintf(w, "    mov   eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    tzcnt eax, eax\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_POPCNT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.popcnt\n")
		fmt.Fprintf(w, "    mov    eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    popcnt eax, eax\n")
		fmt.Fprintf(w, "    mov    dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_ADD:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.add\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_SUB:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.sub\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    sub eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_MUL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.mul\n")
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    imul eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov  dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_DIV_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.div_s\n")
		fmt.Fprintf(w, "    push rdx\n")
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cdq  # edx = copysign(eax)\n")
		fmt.Fprintf(w, "    idiv dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov  dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintf(w, "    pop  rdx\n")
		fmt.Fprintln(w)

	case token.INS_I32_DIV_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.div_u\n")
		fmt.Fprintf(w, "    push rdx\n")
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    xor  edx, edx # 无符号高位清零\n")
		fmt.Fprintf(w, "    div  dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov  dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintf(w, "    pop  rdx\n")
		fmt.Fprintln(w)

	case token.INS_I32_REM_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.rem_s\n")
		fmt.Fprintf(w, "    push rdx\n")
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cdq  # edx = copysign(eax)\n")
		fmt.Fprintf(w, "    idiv dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov  dword ptr [rbp%+d], edx\n", ret0)
		fmt.Fprintf(w, "    pop  rdx\n")
		fmt.Fprintln(w)

	case token.INS_I32_REM_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.rem_u\n")
		fmt.Fprintf(w, "    push rdx\n")
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    xor  edx, edx # 无符号高位清零\n")
		fmt.Fprintf(w, "    div  dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov  dword ptr [rbp%+d], edx\n", ret0)
		fmt.Fprintf(w, "    pop  rdx\n")
		fmt.Fprintln(w)

	case token.INS_I32_AND:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.and\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    and eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_OR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.or\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    or  eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_XOR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.xor\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    xor eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_SHL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.shl\n")
		fmt.Fprintf(w, "    push rcx\n")
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov  ecx, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    shl  eax, cl # cl 是 ecx 低8位\n")
		fmt.Fprintf(w, "    mov  dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintf(w, "    pop  rcx\n")
		fmt.Fprintln(w)

	case token.INS_I32_SHR_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.shr_s\n")
		fmt.Fprintf(w, "    push rcx\n")
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov  ecx, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    sar  eax, cl # cl 是 ecx 低8位\n")
		fmt.Fprintf(w, "    mov  dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintf(w, "    pop  rcx\n")
		fmt.Fprintln(w)

	case token.INS_I32_SHR_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.shr_u\n")
		fmt.Fprintf(w, "    push rcx\n")
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov  ecx, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    shr  eax, cl # cl 是 ecx 低8位\n")
		fmt.Fprintf(w, "    mov  dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintf(w, "    pop  rcx\n")
		fmt.Fprintln(w)

	case token.INS_I32_ROTL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.rotl\n")
		fmt.Fprintf(w, "    push rcx\n")
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov  ecx, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    rol  eax, cl # cl 是 ecx 低8位\n")
		fmt.Fprintf(w, "    mov  dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintf(w, "    pop  rcx\n")
		fmt.Fprintln(w)

	case token.INS_I32_ROTR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.rotr\n")
		fmt.Fprintf(w, "    push rcx\n")
		fmt.Fprintf(w, "    mov  eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov  ecx, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    ror  eax, cl # cl 是 ecx 低8位\n")
		fmt.Fprintf(w, "    mov  dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintf(w, "    pop  rcx\n")
		fmt.Fprintln(w)

	case token.INS_I64_CLZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.clz\n")
		fmt.Fprintf(w, "    mov   rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    lzcnt rax, rax\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_CTZ:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.ctz\n")
		fmt.Fprintf(w, "    mov   rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    tzcnt rax, rax\n")
		fmt.Fprintf(w, "    mov   dword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_POPCNT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.popcnt\n")
		fmt.Fprintf(w, "    mov    rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    popcnt rax, rax\n")
		fmt.Fprintf(w, "    mov    dword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_ADD:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.add\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    add rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_SUB:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.sub\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    sub rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_MUL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.mul\n")
		fmt.Fprintf(w, "    mov  rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    imul rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov  qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_DIV_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.div_s\n")
		fmt.Fprintf(w, "    push rdx\n")
		fmt.Fprintf(w, "    mov  rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cqo  # rdx = copysign(rax)\n")
		fmt.Fprintf(w, "    idiv qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov  qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintf(w, "    pop rdx\n")
		fmt.Fprintln(w)

	case token.INS_I64_DIV_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.div_u\n")
		fmt.Fprintf(w, "    push rdx\n")
		fmt.Fprintf(w, "    mov  rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    xor  rdx, rdx # 无符号高位清零\n")
		fmt.Fprintf(w, "    div  qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov  qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintf(w, "    pop  rdx\n")
		fmt.Fprintln(w)

	case token.INS_I64_REM_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.rem_s\n")
		fmt.Fprintf(w, "    push rdx\n")
		fmt.Fprintf(w, "    mov  rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    cqo  # rdx = copysign(rax)\n")
		fmt.Fprintf(w, "    idiv qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov  qword ptr [rbp%+d], rdx\n", ret0)
		fmt.Fprintf(w, "    pop  rdx\n")
		fmt.Fprintln(w)

	case token.INS_I64_REM_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.rem_u\n")
		fmt.Fprintf(w, "    push rdx\n")
		fmt.Fprintf(w, "    mov  rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    xor  rdx, rdx # 无符号高位清零\n")
		fmt.Fprintf(w, "    div  qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov  qword ptr [rbp%+d], rdx\n", ret0)
		fmt.Fprintf(w, "    pop  rdx\n")
		fmt.Fprintln(w)

	case token.INS_I64_AND:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.and\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    and rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_OR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.or\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    or  rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_XOR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.xor\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    xor rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_SHL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.shl\n")
		fmt.Fprintf(w, "    push rcx\n")
		fmt.Fprintf(w, "    mov  rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov  rcx, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    shl  rax, cl # cl 是 rcx 低8位\n")
		fmt.Fprintf(w, "    mov  qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintf(w, "    pop  rcx\n")
		fmt.Fprintln(w)

	case token.INS_I64_SHR_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.shr_s\n")
		fmt.Fprintf(w, "    push rcx\n")
		fmt.Fprintf(w, "    mov  rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov  rcx, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    sar  rax, cl # cl 是 rcx 低8位\n")
		fmt.Fprintf(w, "    mov  qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintf(w, "    pop  rcx\n")
		fmt.Fprintln(w)

	case token.INS_I64_SHR_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.shr_u\n")
		fmt.Fprintf(w, "    push rcx\n")
		fmt.Fprintf(w, "    mov  rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov  rcx, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    shr  rax, cl # cl 是 rcx 低8位\n")
		fmt.Fprintf(w, "    mov  qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintf(w, "    pop  rcx\n")
		fmt.Fprintln(w)

	case token.INS_I64_ROTL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.rotl\n")
		fmt.Fprintf(w, "    push rcx\n")
		fmt.Fprintf(w, "    mov  rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov  rcx, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    rol  rax, cl # cl 是 rcx 低8位\n")
		fmt.Fprintf(w, "    mov  qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintf(w, "    pop  rcx\n")
		fmt.Fprintln(w)

	case token.INS_I64_ROTR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.rotr\n")
		fmt.Fprintf(w, "    push rcx\n")
		fmt.Fprintf(w, "    mov  rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov  rcx, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    ror  rax, cl # cl 是 rcx 低8位\n")
		fmt.Fprintf(w, "    mov  qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintf(w, "    pop  rcx\n")
		fmt.Fprintln(w)

	case token.INS_F32_ABS:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.abs\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    and eax, 0x7FFFFFFF\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_NEG:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.neg\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    xor rax, 0x80000000\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_CEIL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.ceil\n")
		fmt.Fprintf(w, "    movss   xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundss xmm4, xmm4, 2\n")
		fmt.Fprintf(w, "    movss   dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_FLOOR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.floor\n")
		fmt.Fprintf(w, "    movss   xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundss xmm4, xmm4, 1\n")
		fmt.Fprintf(w, "    movss   dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_TRUNC:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.trunc\n")
		fmt.Fprintf(w, "    movss   xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundss xmm4, xmm4, 3\n")
		fmt.Fprintf(w, "    movss   dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_NEAREST:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.nearest\n")
		fmt.Fprintf(w, "    movss   xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundss xmm4, xmm4, 0\n")
		fmt.Fprintf(w, "    movss   dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_SQRT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.sqrt\n")
		fmt.Fprintf(w, "    movss  xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    sqrtss xmm4, xmm4\n")
		fmt.Fprintf(w, "    movss  dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_ADD:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.add\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    addss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_SUB:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.sub\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    subss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_MUL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.mul\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mulss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_DIV:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.div\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    divss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_MIN:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.min\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    minss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_MAX:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.max\n")
		fmt.Fprintf(w, "    movss xmm4, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    maxss xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movss dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_COPYSIGN:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.copysign\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov r10d, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    and eax, 0x7FFFFFFF\n")
		fmt.Fprintf(w, "    and r10d, 0x80000000\n")
		fmt.Fprintf(w, "    or  eax, r10d\n")
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_ABS:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.abs\n")
		fmt.Fprintf(w, "    mov    rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movabs r11, 0x7FFFFFFFFFFFFFFF\n")
		fmt.Fprintf(w, "    and    rax, r11\n")
		fmt.Fprintf(w, "    mov    qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_NEG:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.neg\n")
		fmt.Fprintf(w, "    mov    rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movabs r11, 0x8000000000000000\n")
		fmt.Fprintf(w, "    xor    rax, r11\n")
		fmt.Fprintf(w, "    mov    qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_CEIL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.ceil\n")
		fmt.Fprintf(w, "    movsd   xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundsd xmm4, xmm4, 2\n")
		fmt.Fprintf(w, "    movsd   qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_FLOOR:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.floor\n")
		fmt.Fprintf(w, "    movsd   xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundsd xmm4, xmm4, 1\n")
		fmt.Fprintf(w, "    movsd   qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_TRUNC:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.trunc\n")
		fmt.Fprintf(w, "    movsd   xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundsd xmm4, xmm4, 3\n")
		fmt.Fprintf(w, "    movsd   qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_NEAREST:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.nearest\n")
		fmt.Fprintf(w, "    movsd   xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    roundsd xmm4, xmm4, 0\n")
		fmt.Fprintf(w, "    movsd   qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_SQRT:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.sqrt\n")
		fmt.Fprintf(w, "    movsd  xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    sqrtsd xmm4, xmm4\n")
		fmt.Fprintf(w, "    movsd  qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_ADD:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.add\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    addsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_SUB:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.sub\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    subsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_MUL:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.mul\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mulsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_DIV:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.div\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    divsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_MIN:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.min\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    minsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_MAX:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.max\n")
		fmt.Fprintf(w, "    movsd xmm4, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    maxsd xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movsd qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_COPYSIGN:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		sp1 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.copysign\n")
		fmt.Fprintf(w, "    mov    rax, qword ptr [rbp%+d]\n", sp1)
		fmt.Fprintf(w, "    mov    r10, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    movabs r11, 0x7FFFFFFFFFFFFFFF\n")
		fmt.Fprintf(w, "    and    rax, r11\n")
		fmt.Fprintf(w, "    movabs r11, 0x8000000000000000\n")
		fmt.Fprintf(w, "    and    r10, r11\n")
		fmt.Fprintf(w, "    or     rax, r10\n")
		fmt.Fprintf(w, "    mov    qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_WRAP_I64:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.wrap_i64\n")
		fmt.Fprintf(w, "    mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_TRUNC_F32_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.trunc_f32_s\n")
		fmt.Fprintf(w, "    movss     xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttss2si eax, xmm4\n")
		fmt.Fprintf(w, "    mov       dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_TRUNC_F32_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.trunc_f32_u\n")
		fmt.Fprintf(w, "    movss     xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttss2si rax, xmm4\n")
		fmt.Fprintf(w, "    mov       dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_TRUNC_F64_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.trunc_f64_s\n")
		fmt.Fprintf(w, "    movsd     xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttsd2si eax, xmm4\n")
		fmt.Fprintf(w, "    mov       dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_TRUNC_F64_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		fmt.Fprintf(w, "    # i32.trunc_f64_u\n")
		fmt.Fprintf(w, "    movsd     xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttsd2si rax, xmm4\n")
		fmt.Fprintf(w, "    mov       dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_EXTEND_I32_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.extend_i32_s\n")
		fmt.Fprintf(w, "    movsxd rax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov    qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_EXTEND_I32_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.extend_i32_u\n")
		fmt.Fprintf(w, "    mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_TRUNC_F32_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.trunc_f32_s\n")
		fmt.Fprintf(w, "    movss     xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttss2si rax, xmm4\n")
		fmt.Fprintf(w, "    mov       qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_TRUNC_F32_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.trunc_f32_u\n")
		fmt.Fprintf(w, "    movss     xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttss2si rax, xmm4\n")
		fmt.Fprintf(w, "    mov       qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_TRUNC_F64_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.trunc_f64_s\n")
		fmt.Fprintf(w, "    movsd     xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttsd2si rax, xmm4\n")
		fmt.Fprintf(w, "    mov       qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_TRUNC_F64_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		fmt.Fprintf(w, "    # i64.trunc_f64_u\n")
		fmt.Fprintf(w, "    movsd     xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvttsd2si rax, xmm4\n")
		fmt.Fprintf(w, "    mov       qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_CONVERT_I32_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.convert_i32_s\n")
		fmt.Fprintf(w, "    mov      eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2ss xmm4, eax\n")
		fmt.Fprintf(w, "    movss    dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_CONVERT_I32_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.convert_i32_u\n")
		fmt.Fprintf(w, "    mov      eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2ss xmm4, rax\n")
		fmt.Fprintf(w, "    movss    dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_CONVERT_I64_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.convert_i64_s\n")
		fmt.Fprintf(w, "    mov      rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2ss xmm4, rax\n")
		fmt.Fprintf(w, "    movss    dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_CONVERT_I64_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.convert_i64_u\n")
		fmt.Fprintf(w, "    mov      rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2ss xmm4, rax\n")
		fmt.Fprintf(w, "    movss    dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_DEMOTE_F64:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		fmt.Fprintf(w, "    # f32.demote_f64\n")
		fmt.Fprintf(w, "    movsd    xmm4, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsd2ss xmm4, xmm4\n")
		fmt.Fprintf(w, "    movss    dword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_CONVERT_I32_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.convert_i32_s\n")
		fmt.Fprintf(w, "    mov      eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2sd xmm4, eax\n")
		fmt.Fprintf(w, "    movsd    qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_CONVERT_I32_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.convert_i32_u\n")
		fmt.Fprintf(w, "    mov      eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2sd xmm4, rax\n")
		fmt.Fprintf(w, "    movsd    qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_CONVERT_I64_S:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.convert_i64_s\n")
		fmt.Fprintf(w, "    mov      rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2sd xmm4, rax\n")
		fmt.Fprintf(w, "    movsd    qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_CONVERT_I64_U:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.convert_i64_u\n")
		fmt.Fprintf(w, "    mov      rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtsi2sd xmm4, rax\n")
		fmt.Fprintf(w, "    movsd    qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_PROMOTE_F32:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		fmt.Fprintf(w, "    # f64.promote_f32\n")
		fmt.Fprintf(w, "    movss    xmm4, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    cvtss2sd xmm4, xmm4\n")
		fmt.Fprintf(w, "    movsd    qword ptr [rbp%+d], xmm4\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I32_REINTERPRET_F32:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I32) - 8

		assert(sp0 == ret0)
		fmt.Fprintf(w, "    # i32.reinterpret_f32\n")
		fmt.Fprintf(w, "    # mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    # mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_I64_REINTERPRET_F64:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.F64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.I64) - 8

		assert(sp0 == ret0)
		fmt.Fprintf(w, "    # i64.reinterpret_f64\n")
		fmt.Fprintf(w, "    # mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    # mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F32_REINTERPRET_I32:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I32) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F32) - 8

		assert(sp0 == ret0)
		fmt.Fprintf(w, "    # f32.reinterpret_i32\n")
		fmt.Fprintf(w, "    # mov eax, dword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    # mov dword ptr [rbp%+d], eax\n", ret0)
		fmt.Fprintln(w)

	case token.INS_F64_REINTERPRET_I64:
		sp0 := p.fnWasmR0Base - 8*stk.Pop(token.I64) - 8
		ret0 := p.fnWasmR0Base - 8*stk.Push(token.F64) - 8

		assert(sp0 == ret0)
		fmt.Fprintf(w, "    # f64.reinterpret_i64\n")
		fmt.Fprintf(w, "    # mov rax, qword ptr [rbp%+d]\n", sp0)
		fmt.Fprintf(w, "    # mov qword ptr [rbp%+d], rax\n", ret0)
		fmt.Fprintln(w)

	default:
		panic(fmt.Sprintf("unreachable: %T", i))
	}
	return nil
}
