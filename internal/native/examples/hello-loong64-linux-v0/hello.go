// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build ingore

package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	flagFilename    = flag.String("f", "", "set file name")
	flagSection     = flag.String("s", ".text", "set elf section name")
	flagListSection = flag.Bool("l", false, "list elf section name")
)

func main() {
	flag.Parse()

	if *flagFilename == "" || *flagSection == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := elf.Open(*flagFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if *flagListSection {
		for _, sec := range f.Sections {
			fmt.Println(sec.Name)
		}
		return
	}

	for _, sec := range f.Sections {
		if sec.Name == *flagSection {
			fmt.Printf("Name     : %v\n", sec.Name)
			fmt.Printf("Addr     : 0x%08x\n", sec.Addr)
			fmt.Printf("Addralign: %v\n", sec.Addralign)
			fmt.Printf("Offset   : 0x%08x\n", sec.Offset)
			fmt.Printf("Size     : %v\n", sec.Size)
			fmt.Printf("FileSize : %v\n", sec.FileSize)

			data, err := sec.Data()
			if err != nil {
				panic(err)
			}

			if len(data) > 0 {
				fmt.Println(strings.Repeat("-", 8+16*3))
				fmt.Print(strings.Repeat(" ", 8+1))
				for i := 0; i < 16; i++ {
					fmt.Printf("%02X ", i)
				}
				fmt.Println()
			}

			for i, d := range data {
				if i%16 == 0 {
					if i > 0 {
						fmt.Println()
					}
					fmt.Printf("%08x ", sec.Offset+uint64(i))
				}
				fmt.Printf("%02X ", d)
			}
		}
	}
}
