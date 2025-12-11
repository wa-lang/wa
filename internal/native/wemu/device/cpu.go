// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package device

import "wa-lang.org/wa/internal/native/abi"

// CPU接口
type CPU interface {
	GetPC() uint64
	SetPC(v uint64)

	XRegNum() int
	GetXReg(i int) uint64
	SetXReg(i int, v uint64)

	FRegNum() int
	GetFReg(i int) float64
	SetFReg(i int, v float64)

	Reset(pc, sp uint64)
	StepRun(bus *Bus) error

	LookupRegister(regName string) (r abi.RegType, ok bool)
	RegAliasString(r abi.RegType) string

	InstString(x uint32) (string, error)
}
