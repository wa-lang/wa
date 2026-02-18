// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build ignore

package main

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/native/link/pe"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: peinfo file.exe")
		return
	}

	f, err := pe.OpenPE64(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Println("=== File Header ===")
	fmt.Printf("Machine: 0x%x\n", f.FileHeader.Machine)
	fmt.Printf("Sections: %d\n", f.FileHeader.NumberOfSections)

	if opt := f.OptionalHeader; true {
		fmt.Println("\n=== Optional Header ===")
		fmt.Printf("ImageBase: 0x%x\n", opt.ImageBase)
		fmt.Printf("EntryPoint RVA: 0x%x\n", opt.AddressOfEntryPoint)
		fmt.Printf("Subsystem: %d\n", opt.Subsystem)
	}

	fmt.Println("\n=== Sections ===")
	for _, s := range f.Sections {
		fmt.Printf(
			"%s RVA=0x%x Raw=0x%x Size=0x%x\n",
			s.Name,
			s.VirtualAddress,
			s.Offset,
			s.Size,
		)
	}

	fmt.Println("\n=== Imports ===")

	imports, err := f.ImportedSymbols()
	if err == nil {
		for _, imp := range imports {
			fmt.Println(imp)
		}
	}
}
