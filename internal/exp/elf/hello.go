// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build ignore

package main

import (
	"debug/elf"
	"fmt"
	"log"
	"os"
)

func main() {
	filePath := "testdata/tiny/tiny.elf.bin"

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	ef, err := elf.NewFile(f)
	if err != nil {
		log.Fatal(err)
	}
	defer ef.Close()

	fmt.Printf("=== ELF Header ===\n")
	fmt.Printf("Class: %v\n", ef.Class) // ELFCLASS32 / ELFCLASS64
	fmt.Printf("Data: %v\n", ef.Data)   // LSB/MSB
	fmt.Printf("OSABI: %v\n", ef.OSABI)
	fmt.Printf("Type: %v\n", ef.Type)       // ET_EXEC, ET_DYN, ET_REL
	fmt.Printf("Machine: %v\n", ef.Machine) // EM_RISCV, EM_X86_64 等
	fmt.Printf("Entry: 0x%x\n", ef.Entry)   // 程序入口地址

	fmt.Printf("\n=== Sections ===\n")
	for _, sec := range ef.Sections {
		fmt.Printf("Name: %-16s Type: %-12v Addr: 0x%x Size: 0x%x\n",
			sec.Name, sec.Type, sec.Addr, sec.Size,
		)
	}

	fmt.Printf("\n=== Segments (Program Headers) ===\n")
	for i, prog := range ef.Progs {
		fmt.Printf("[%d] Type: %-8v VirtAddr: 0x%x FileSize: 0x%x MemSize: 0x%x Flags: %v\n",
			i, prog.Type, prog.Vaddr, prog.Filesz, prog.Memsz, prog.Flags,
		)

		fmt.Printf("Off: %x, Filesz: %d\n", prog.Off, prog.Filesz)
		data := make([]byte, prog.Filesz)
		if _, err := prog.ReadAt(data, 0); err != nil {
			panic(err)
		}
		fmt.Printf("data: %x\n", data)
	}
}
