// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package objabi

import "fmt"

// 寄存器区间
type regSet struct {
	lo    RBaseType              // 开始
	hi    RBaseType              // 结束(开区间)
	Rconv func(RBaseType) string // 用于打印
}

// 用于注册不同架构下的寄存器
// 不同架构的寄存器编号处于不同空间不会冲突
var regSpace []regSet

// 注册不同平台的寄存器区间
func RegisterRegister(lo, hi RBaseType, Rconv func(RBaseType) string) {
	regSpace = append(regSpace, regSet{lo, hi, Rconv})
}

// 将寄存器编号转为字符串名字
func (reg RBaseType) String() string {
	if reg == REG_NONE {
		return "NONE"
	}
	for i := range regSpace {
		rs := &regSpace[i]
		if rs.lo <= reg && reg < rs.hi {
			return rs.Rconv(reg)
		}
	}
	return fmt.Sprintf("R???%d", reg)
}
