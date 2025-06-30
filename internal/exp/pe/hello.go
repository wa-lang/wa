// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build ignore

package main

import (
	"debug/pe"
	"fmt"
	"os"
)

func main() {
	filePath := "testdata/pts-tinype/hh2.golden.exe"

	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	peFile, err := pe.NewFile(f)
	if err != nil {
		panic(err) // panic: size of data directories(24) is inconsistent with number of data directories(2)
	}
	defer peFile.Close()

	fmt.Printf("=== FileHeader ===\n")
	fmt.Printf("Machine: 0x%x\n", peFile.FileHeader.Machine)
	fmt.Printf("Number of Sections: %d\n", peFile.FileHeader.NumberOfSections)
	fmt.Printf("TimeDateStamp: 0x%x\n", peFile.FileHeader.TimeDateStamp)

	fmt.Printf("\n=== OptionalHeader ===\n")
	switch oh := peFile.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		fmt.Printf("32-bit, EntryPoint: 0x%x, ImageBase: 0x%x\n", oh.AddressOfEntryPoint, oh.ImageBase)
	case *pe.OptionalHeader64:
		fmt.Printf("64-bit, EntryPoint: 0x%x, ImageBase: 0x%x\n", oh.AddressOfEntryPoint, oh.ImageBase)
	default:
		fmt.Printf("Unknown OptionalHeader format\n")
	}

	fmt.Printf("\n=== Sections ===\n")
	for _, sec := range peFile.Sections {
		fmt.Printf("Name: %s | VirtualAddr: 0x%x | Size: 0x%x\n",
			sec.Name, sec.VirtualAddress, sec.Size,
		)
	}
}
