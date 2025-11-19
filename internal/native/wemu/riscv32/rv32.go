// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv32

import "wa-lang.org/wa/internal/native/wemu/device"

// 寄存器整数
type RVUInt = uint32

var _ device.CPU = (*CPU)(nil)
