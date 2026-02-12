// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"fmt"
	"log"
	"sort"
)

func (p *Prog) keys() []string {
	var keys []string
	for key := range p.Child {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// 根据指定的 action 和 key 在决策树中前进或开辟路径
//
// 1. 如果当前节点 p 尚未定义 Action，则赋予其指定的动作. 否则新动作必须与原动作一致(禁止多种判定逻辑)
// 2. 如果存在对应 key 的子分支, 则直接返回该子树节点
// 3. 如果不存在，则新建一个 Prog 节点作为子树并返回
func (p *Prog) walk(action, key, text, opcode string) *Prog {
	if p.Action == "" {
		p.Action = action
	} else if p.Action != action {
		log.Printf("%s; %s: conflicting paths %s and %s|%s %s\n", text, opcode, p.findChildLeaf(), p.Path, action, key)
		return new(Prog)
	}
	q := p.Child[key]
	if q == nil {
		if p.Child == nil {
			p.Child = make(map[string]*Prog)
		}
		q = new(Prog)
		q.Path = fmt.Sprintf("%s|%s %s", p.Path, action, key)
		p.Child[key] = q
	}
	return q
}

// 寻找一个叶子节点并返回完整路径
//
// 主要用于错误信息中, 展示某个特定的子树分支最终会导向哪条指令
func (p *Prog) findChildLeaf() string {
	for {
		if len(p.Child) == 0 {
			return p.Path
		}
		p = p.Child[p.keys()[0]]
	}
}
