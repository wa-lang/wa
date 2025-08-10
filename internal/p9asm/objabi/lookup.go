// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package objabi

import (
	"fmt"
)

// 根据名字查找汇编指令, 失败返回 objabi.AXXX
func LookupAs(asName string) As {
	if asName == "XXX" {
		return AXXX
	}
	for i, s := range Anames {
		if s == asName {
			return ABase + As(i)
		}
	}
	return AXXX
}

// 汇编指令转字符串格式
func AsString(as As) string {
	if ABase <= as && as < A_ARCHSPECIFIC {
		return Anames[as-ABase]
	}
	return fmt.Sprintf("objabi.As(%d)", as-ABase)
}
