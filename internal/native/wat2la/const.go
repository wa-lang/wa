// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"
	"sort"
)

// 注册常量
func (p *wat2laWorker) registerConst(x uint64) {
	p.constLitMap[x] = x
}

// 生成常量
func (p *wat2laWorker) buildConstList(w io.Writer) error {
	var xList = make([]uint64, 0, len(p.constLitMap))
	for x := range p.constLitMap {
		xList = append(xList, x)
	}
	sort.Slice(xList, func(i, j int) bool {
		return xList[i] < xList[j]
	})
	for _, x := range xList {
		fmt.Fprintf(w, "global %s: u64 = %d\n", p.getConstName(x), x)
	}
	return nil
}
