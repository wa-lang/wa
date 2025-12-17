// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package pcrel_test

import (
	"fmt"

	"wa-lang.org/wa/internal/native/pcrel"
)

func ExampleMakeLa64PCRel() {
	const (
		pc      = 0x12000010c
		message = 0x120010130
	)

	pc_hi20, pc_lo12 := pcrel.MakeLa64PCRel(message, pc)
	addr := pcrel.GetTargetAddressLa64(pc, pc_hi20, pc_lo12)

	fmt.Printf("pc_hi20: %d\n", pc_hi20)
	fmt.Printf("pc_lo12: %d\n", pc_lo12)
	fmt.Printf("message: 0x%X\n", addr)

	// Output:
	// pc_hi20: 16
	// pc_lo12: 304
	// message: 0x120010130
}

func ExampleMakeAbs() {
	const UART0 = 0x10000000
	pcrel_hi, pcrel_lo := pcrel.MakeAbs(UART0)

	fmt.Printf("pcrel_hi: %d\n", pcrel_hi)
	fmt.Printf("pcrel_lo: %d\n", pcrel_lo)
	addr := pcrel.GetTargetAddress(0, pcrel_hi, pcrel_lo)
	fmt.Printf("UART0: 0x%08X\n", addr)

	// Output:
	// pcrel_hi: 65536
	// pcrel_lo: 0
	// UART0: 0x10000000
}

func ExampleMakePCRel() {
	const message = 0x80000028
	const pc = 0x80000000

	pcrel_hi, pcrel_lo := pcrel.MakePCRel(message, pc)
	addr := pcrel.GetTargetAddress(pc, pcrel_hi, pcrel_lo)

	fmt.Printf("pcrel_hi: %d\n", pcrel_hi)
	fmt.Printf("pcrel_lo: %d\n", pcrel_lo)
	fmt.Printf("message: 0x%08X\n", addr)

	// Output:
	// pcrel_hi: 0
	// pcrel_lo: 40
	// message: 0x80000028
}

func ExampleMakePCRel_negative() {
	const message = 0x7FFFFFF0
	const pc = 0x80000000

	pcrel_hi, pcrel_lo := pcrel.MakePCRel(message, pc)
	addr := pcrel.GetTargetAddress(pc, pcrel_hi, pcrel_lo)

	fmt.Printf("message: 0x%08X\n", addr)

	// Output:
	// message: 0x7FFFFFF0
}

func ExampleMakePCRel_equal() {
	const message = 0x80000000
	const pc = 0x80000000

	pcrel_hi, pcrel_lo := pcrel.MakePCRel(message, pc)
	addr := pcrel.GetTargetAddress(pc, pcrel_hi, pcrel_lo)

	fmt.Printf("message: 0x%08X\n", addr)

	// Output:
	// message: 0x80000000
}
