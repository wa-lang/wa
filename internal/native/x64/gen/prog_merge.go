// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import "log"

// 合并决策子树副本
//
// 此函数通常在所有路径添加完成后调用.
// 为了优化结构, 它会引入跨节点链接, 将原本层级分明的树转变为有向无环图,
// 从而实现不同指令路径间公共逻辑的共享
func (dst *Prog) mergeCopy(src *Prog) {
	if dst.Action != src.Action {
		log.Printf("cannot merge %s|%s and %s|%s", dst.Path, dst.Action, src.Path, src.Action)
		return
	}

	for _, key := range src.keys() {
		if dst.Child[key] == nil {
			dst.Child[key] = src.Child[key]
		} else {
			dst.Child[key].mergeCopy(src.Child[key])
		}
	}
}
