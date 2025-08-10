// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// TODO: 删除该文件

package objabi

import "fmt"

// 机器码和对应的名字
type opSet struct {
	lo    As
	names []string
}

// 机器码区间集合
var aSpace []opSet

// 注册不同平台的指令区间
// 范围写死, 但是查表不要太大, 需要手动区分
// 能同时支持不同平台的混编吗
func RegisterOpcode(lo As, anames []string) {
	aSpace = append(aSpace, opSet{lo, anames})
}

// 指令机器码转为字符串名字
func (a As) String() string {
	if a < A_ARCHSPECIFIC {
		return Anames[a]
	}
	for i := range aSpace {
		as := &aSpace[i]
		if as.lo <= a && a < as.lo+As(len(as.names)) {
			return as.names[a-as.lo]
		}
	}
	return fmt.Sprintf("A???%d", a)
}
