// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appobjdump

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/native/link/elf"
	"wa-lang.org/wa/internal/native/loong64"
	"wa-lang.org/wa/internal/native/riscv"
	"wa-lang.org/wa/internal/printer/tabwriter"
)

var CmdObjdump = &cli.Command{
	Hidden:    true,
	Name:      "objdump",
	Usage:     "dump elf text and data sections",
	ArgsUsage: "<file.elf>",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "prog",
			Usage: "print program data",
		},
		&cli.BoolFlag{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list elf section name",
		},
		&cli.StringFlag{
			Name:    "section",
			Aliases: []string{"s"},
			Usage:   "set section name",
		},
		&cli.BoolFlag{
			Name:  "zh",
			Usage: "set zh mode",
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
		isZhMode := c.Bool("zh")

		if sectionName != "" || listSection {
			cmdObjdump(infile, sectionName, listSection)
		} else {
			cmdProgdump(infile, isZhMode)
		}

		return nil
	},
}

// 打印 prog 数据
func cmdProgdump(filename string, isZhMode bool) {
	f, err := elf.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Printf("Class  : %v\n", f.Class)
	fmt.Printf("Version: %v\n", f.Version)
	fmt.Printf("OS/ABI : %v\n", f.OSABI)
	fmt.Printf("Machine: %v\n", f.Machine)
	fmt.Printf("Entry  : %x\n", f.Entry)
	fmt.Println()

	for _, p := range f.Progs {
		if p.Type != elf.PT_LOAD || p.Flags&elf.PF_R == 0 {
			continue // 跳过不可读部分
		}

		if p.Flags&elf.PF_X != 0 {
			textAddr := uint64(p.Vaddr)
			textData := make([]byte, p.Filesz)
			_, err := io.ReadFull(p.Open(), textData)
			if err != nil {
				fmt.Println("ERR:", err)
				continue
			}
			if f.Entry > textAddr {
				diff := f.Entry - textAddr
				textData = textData[diff:]
				textAddr = f.Entry
			}

			printProgText(f.Machine, textAddr, textData, isZhMode)
		} else {
			dataAddr := uint64(p.Vaddr)
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

func printProgText(machine elf.Machine, addr uint64, data []byte, isZhMode bool) {
	addrWidth := len(fmt.Sprintf("%x", addr))

	// 用于格式化指令对齐
	var buf bytes.Buffer
	var w = tabwriter.NewWriter(&buf, 1, 1, 1, ' ', 0)

	fmt.Fprintln(w, "[.text.]")
	for k := 0; k < len(data); k += 4 {
		fmt.Fprintf(w, "%0*X: ", addrWidth, addr+uint64(k))
		x := binary.LittleEndian.Uint32(data[k:][:4])
		instStr := []byte(decodeInst(machine, x, isZhMode))
		if idx := bytes.IndexByte(instStr, ' '); idx > 0 {
			instStr[idx] = '\t'
		}
		fmt.Fprintf(w, "%08X # %s\n", x, string(instStr))
	}
	fmt.Fprintln(w)

	w.Flush()
	fmt.Print(buf.String())
}

func decodeInst(machine elf.Machine, x uint32, isZhMode bool) string {
	switch machine {
	case elf.EM_LOONGARCH:
		as, arg, err := loong64.Decode(x)
		if err != nil {
			return "???"
		}
		if isZhMode {
			return loong64.AsmSyntaxEx(as, "", arg, loong64.ZhRegAliasString, loong64.ZhAsString)
		} else {
			return loong64.AsmSyntax(as, "", arg)
		}

	case elf.EM_RISCV:
		as, arg, err := riscv.Decode(x)
		if err != nil {
			return "???"
		}
		if isZhMode {
			return riscv.AsmSyntaxEx(as, "", arg, riscv.ZhRegAliasString, riscv.ZhAsString)
		} else {
			return riscv.AsmSyntax(as, "", arg)
		}

	default:
		panic(fmt.Sprintf("unsupport %v", machine))
	}
}

func printProgData(addr uint64, data []byte) {
	addrWidth := len(fmt.Sprintf("%x", addr))

	fmt.Printf("%-*s ", addrWidth, "[.data.]")
	for i := 0; i < 16; i++ {
		fmt.Printf(" %02X", i)
	}
	fmt.Println()

	for k := 0; k < len(data); k += 16 {
		fmt.Printf("%0*X: ", addrWidth, addr+uint64(k))
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
