// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appobjdump

import (
	"debug/elf"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/native/loong64"
	"wa-lang.org/wa/internal/native/riscv"
)

var CmdObjdump = &cli.Command{
	Hidden:    true,
	Name:      "objdump",
	Usage:     "dump elf text and data sections",
	ArgsUsage: "<file.elf>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "section",
			Aliases: []string{"s"},
			Usage:   "set section name",
			Value:   ".text",
		},
		&cli.BoolFlag{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list elf section name",
		},
		&cli.BoolFlag{
			Name:  "prog",
			Usage: "print program data",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "no input file")
			os.Exit(1)
		}

		infile := c.Args().First()
		sectionName := c.String("section")
		listSection := c.Bool("list")
		prog := c.Bool("prog")

		if prog {
			cmdProgdump(infile)
		} else {
			cmdObjdump(infile, sectionName, listSection)
		}
		return nil
	},
}

// 打印 prog 数据
func cmdProgdump(filename string) {
	f, err := elf.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Printf("Machine: %v\n", f.Machine)
	fmt.Printf("Class  : %v\n", f.Class)
	fmt.Println()

	for _, p := range f.Progs {
		if p.Type != elf.PT_LOAD || p.Flags&elf.PF_R == 0 {
			continue // 跳过不可读部分
		}

		if p.Flags&elf.PF_X != 0 {
			textAddr := int64(p.Vaddr)
			textData := make([]byte, p.Filesz)
			_, err := io.ReadFull(p.Open(), textData)
			if err != nil {
				fmt.Println("ERR:", err)
				continue
			}
			printProgText(f.Machine, textAddr, textData)
		} else {
			dataAddr := int64(p.Vaddr)
			dataData := make([]byte, p.Filesz)
			_, err := io.ReadFull(p.Open(), dataData)
			if err != nil {
				fmt.Println("ERR:", err)
				continue
			}
			printProgData(dataAddr, dataData)
		}
	}
}

func printProgText(machine elf.Machine, addr int64, data []byte) {
	fmt.Printf("[.text.] ")
	for i := 0; i < 4; i++ {
		fmt.Printf("%02X ", i)
	}
	fmt.Println()

	for k := 0; k < len(data); k += 4 {
		fmt.Printf("%08X ", addr+int64(k))
		for i := 0; i < 4; i++ {
			if k+i < len(data) {
				fmt.Printf("%02X ", data[k+i])
			} else {
				fmt.Print("   ")
			}
		}
		x := binary.LittleEndian.Uint32(data[k:][:4])
		fmt.Println("#", decodeInst(machine, x))
	}
	fmt.Println()
}

func decodeInst(machine elf.Machine, x uint32) string {
	switch machine {
	case elf.EM_LOONGARCH:
		as, arg, err := loong64.Decode(x)
		if err != nil {
			return err.Error()
		}
		return loong64.AsmSyntax(as, "", arg)

	case elf.EM_RISCV:
		as, arg, err := riscv.Decode(x)
		if err != nil {
			return err.Error()
		}
		return riscv.AsmSyntax(as, "", arg)

	default:
		return fmt.Sprintf("unsupport %v", machine)
	}
}

func printProgData(addr int64, data []byte) {
	fmt.Printf("[.data.] ")
	for i := 0; i < 16; i++ {
		fmt.Printf("%02X ", i)
	}
	fmt.Println()

	for k := 0; k < len(data); k += 16 {
		fmt.Printf("%08X ", addr+int64(k))
		for i := 0; i < 16; i++ {
			if k+i < len(data) {
				fmt.Printf("%02X ", data[k+i])
			} else {
				fmt.Print("   ")
			}
		}
		for i := 0; i < 16; i++ {
			if k+i < len(data) {
				if unicode.IsPrint(rune(data[k+i])) {
					fmt.Printf("%c", data[k+i])
				} else {
					fmt.Print(".")
				}
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// 打印段数据
func cmdObjdump(filename string, sectionName string, listSection bool) {
	f, err := elf.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if listSection {
		for _, sec := range f.Sections {
			fmt.Println(sec.Name)
		}
		return
	}

	for _, sec := range f.Sections {
		if sec.Name == sectionName {
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
