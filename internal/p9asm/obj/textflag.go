// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package obj

// 函数不再包含标志位, 默认都是固定栈

const (
	DUPOK  = 1 // 可以出现多个重名符号, 取第一个
	RODATA = 2 // 只读数据段
	NOPTR  = 4 // 不包含指针的数据
)
