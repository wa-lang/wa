// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import "log"

// 检查决策树是否格式正确, 并进行规范化处理
//
// 核心逻辑是将 any 节点合并到具体的决策分支中.
// 最终确保每个节点要么只有一个 any 子节点(作为无操作跳转),
// 要么完全没有 any 子节点, 从而消除解码时的多路径歧义
func (p *Prog) check() {
	if p.Child["any"] != nil && len(p.Child) > 1 {
		for _, key := range p.keys() {
			if key != "any" {
				p.Child[key].mergeCopy(p.Child["any"])
			}
		}
		if allKeys[p.Action] == nil {
			log.Printf("%s: unknown key space for %s=any", p.Path, p.Action)
		}
		for _, key := range allKeys[p.Action] {
			if p.Child[key] == nil {
				p.Child[key] = p.Child["any"]
			}
		}
		delete(p.Child, "any")
	}

	for _, q := range p.Child {
		q.check()
	}

	switch p.Action {
	case "op", "read", "arg":
		if len(p.Child) > 1 {
			log.Printf("%s: multiple children for action=%s: %v", p.Path, p.Action, p.keys())
		}
	}
}
