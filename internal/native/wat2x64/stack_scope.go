// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"fmt"
	"strconv"

	"wa-lang.org/wa/internal/wat/token"
	wattoken "wa-lang.org/wa/internal/wat/token"
)

// 块空间状态
type scopeContext struct {
	Type             wattoken.Token   // 区块的类型, 用于区别处理 br 时的返回值
	Label            string           // 嵌套的label查询, if/block/loop
	LabelSuffix      string           // 作用域唯一后缀(避免重名)
	StackBase        int              // if/block/loop, 开始的栈位置
	Result           []wattoken.Token // 对应块的返回值数量和类型
	IgnoreStackCheck bool             // 是否忽略当前栈的精简检查(br和return指令需要忽略)
}

// 空间状态栈
type scopeContextStack struct {
	stack []*scopeContext
}

func (p *scopeContextStack) Len() int {
	return len(p.stack)
}
func (p *scopeContextStack) Top() *scopeContext {
	return p.stack[len(p.stack)-1]
}
func (p *scopeContextStack) Target(idx int) *scopeContext {
	return p.stack[idx]
}
func (p *scopeContextStack) TargetByLabelIndex(labelIdx int) *scopeContext {
	return p.stack[p.Len()-labelIdx-1]
}

func (p *scopeContextStack) EnterScope(typ token.Token, stkBase int, label, labelSuffix string, results []token.Token) {
	p.stack = append(p.stack, &scopeContext{
		Type:             typ,
		StackBase:        stkBase,
		Label:            label,
		LabelSuffix:      labelSuffix,
		Result:           results,
		IgnoreStackCheck: false,
	})
}
func (p *scopeContextStack) LeaveScope() {
	p.stack = p.stack[:len(p.stack)-1]
}

func (p *scopeContextStack) findLabelScopeType(label string) token.Token {
	if label == "" {
		panic("wat2x64: empty label name")
	}

	idx := p.findLabelIndex(label)
	if idx < len(p.stack) {
		return p.stack[len(p.stack)-idx-1].Type
	}
	panic(fmt.Sprintf("wat2x64: unknown label %q", label))
}

func (p *scopeContextStack) findLabelName(label string) string {
	if label == "" {
		panic("wat2x64: empty label name")
	}

	idx := p.findLabelIndex(label)
	if idx < len(p.stack) {
		return p.stack[len(p.stack)-idx-1].Label
	}
	panic(fmt.Sprintf("wat2x64: unknown label %q", label))
}

func (p *scopeContextStack) findLabelSuffixId(label string) string {
	if label == "" {
		panic("wat2x64: empty label suffix id")
	}

	idx := p.findLabelIndex(label)
	if idx < len(p.stack) {
		return p.stack[len(p.stack)-idx-1].LabelSuffix
	}
	panic(fmt.Sprintf("wat2x64: unknown label %q", label))
}

func (p *scopeContextStack) findLabelIndex(label string) int {
	if label == "" {
		panic("wat2x64: empty label index")
	}

	if idx, err := strconv.Atoi(label); err == nil {
		return idx
	}
	for i := 0; i < len(p.stack); i++ {
		if s := p.stack[len(p.stack)-i-1]; s.Label == label {
			return i
		}
	}
	panic(fmt.Sprintf("wat2x64: unknown label %q", label))
}
