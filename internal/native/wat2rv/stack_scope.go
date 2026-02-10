// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2rv

import (
	"strconv"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
	wattoken "wa-lang.org/wa/internal/wat/token"
)

// 块空间状态
type scopeContext struct {
	Type             wattoken.Token    // 区块的类型, 用于区别处理 br 时的返回值
	Label            string            // 嵌套的label查询, if/block/loop
	LabelSuffix      string            // 作用域唯一后缀(避免重名)
	StackBase        int               // if/block/loop, 开始的栈位置
	Result           []wattoken.Token  // 对应块的返回值数量和类型
	IgnoreStackCheck bool              // 是否忽略当前栈的精简检查(br和return指令需要忽略)
	InstList         []ast.Instruction // 用于 body 测试
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

func (p *scopeContextStack) GetReturnScopeContext() *scopeContext {
	return p.stack[0]
}

func (p *scopeContextStack) FindScopeContext(label string) *scopeContext {
	if label == "" {
		panic("wat2la: empty label name")
	}
	if idx, err := strconv.Atoi(label); err == nil {
		if idx >= 0 && idx < len(p.stack) {
			return p.stack[len(p.stack)-idx-1]
		} else {
			panic("wat2la: invalid label index")
		}
	} else {
		// 逆序从内向外部查找
		// 正常情况下 label 是唯一的, 不同顺序不影响结果
		for i := len(p.stack) - 1; i >= 0; i-- {
			if ctx := p.stack[i]; ctx.Label == label {
				return ctx
			}
		}
	}
	return nil
}
