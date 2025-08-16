// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv64

import "io"

// 内存读写接口
type Memory interface {
	io.WriterAt
	io.WriterAt
}

// 默认的内存实现
func NewMemory() Memory {
	return nil
}
