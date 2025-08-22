// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"io"
	"log"
	"os"

	"wa-lang.org/wa/internal/native/riscv"
)

const (
	// ELF 常量
	ELFCLASS64  = 2
	ELFDATA2LSB = 1
	EV_CURRENT  = 1

	ET_EXEC  = 2
	EM_RISCV = 243

	PT_LOAD = 1

	PF_X = 1
	PF_W = 2
	PF_R = 4

	EI_NIDENT = 16

	ELF64_EHDR_SIZE = 64
	ELF64_PHDR_SIZE = 56
)

type Elf64_Ehdr struct {
	Ident     [EI_NIDENT]byte
	Type      uint16
	Machine   uint16
	Version   uint32
	Entry     uint64
	Phoff     uint64
	Shoff     uint64
	Flags     uint32
	Ehsize    uint16
	Phentsize uint16
	Phnum     uint16
	Shentsize uint16
	Shnum     uint16
	Shstrndx  uint16
}

type Elf64_Phdr struct {
	Type   uint32
	Flags  uint32
	Offset uint64
	Vaddr  uint64
	Paddr  uint64
	Filesz uint64
	Memsz  uint64
	Align  uint64
}

func writeAt(w *os.File, off int64, b []byte) {
	if _, err := w.Seek(off, io.SeekStart); err != nil {
		log.Fatalf("seek %d: %v", off, err)
	}
	if _, err := w.Write(b); err != nil {
		log.Fatalf("write at %d: %v", off, err)
	}
}

var (
	textAddr = 0x80000000
	textData = genTextData()

	dataAddr = 0x8000003c
	dataData = []byte("Hello native/risc-v!\n\x00")

	flagOutput = flag.String("o", "a.out.exe", "set output file")
)

func main() {
	flag.Parse()

	// 文件布局：
	// [ELF64_Ehdr (64)] [2*Elf64_Phdr (2*56=112)] [text] [data]
	ehOff := int64(0)
	phOff := int64(ELF64_EHDR_SIZE)
	textOff := int64(ELF64_EHDR_SIZE + 2*ELF64_PHDR_SIZE)
	dataOff := textOff + int64(len(textData))

	// 构造 ELF 头
	var eh Elf64_Ehdr
	eh.Ident[0] = 0x7f
	eh.Ident[1] = 'E'
	eh.Ident[2] = 'L'
	eh.Ident[3] = 'F'
	eh.Ident[4] = ELFCLASS64
	eh.Ident[5] = ELFDATA2LSB
	eh.Ident[6] = EV_CURRENT // 版本
	// 其余 Ident 字节保持 0（System V / 无特别 OSABI）

	eh.Type = ET_EXEC
	eh.Machine = EM_RISCV
	eh.Version = EV_CURRENT
	eh.Entry = uint64(textAddr)
	eh.Phoff = uint64(phOff)
	eh.Shoff = 0 // 不写节区表
	eh.Flags = 0
	eh.Ehsize = ELF64_EHDR_SIZE
	eh.Phentsize = ELF64_PHDR_SIZE
	eh.Phnum = 1
	if len(dataData) > 0 {
		eh.Phnum = 2
	}
	eh.Shentsize = 0
	eh.Shnum = 0
	eh.Shstrndx = 0

	// Program Header: .text (RX)
	textPh := Elf64_Phdr{
		Type:   PT_LOAD,
		Flags:  PF_R | PF_X,
		Offset: uint64(textOff),
		Vaddr:  uint64(textAddr),
		Paddr:  uint64(textAddr), // 裸机：物理=虚拟
		Filesz: uint64(len(textData)),
		Memsz:  uint64(len(textData)),
		Align:  1, // 设为 1，避免 vaddr/offset 对齐约束
	}

	// Program Header: .data (RW)
	dataPh := Elf64_Phdr{
		Type:   PT_LOAD,
		Flags:  PF_R | PF_W,
		Offset: uint64(dataOff),
		Vaddr:  uint64(dataAddr),
		Paddr:  uint64(dataAddr),
		Filesz: uint64(len(dataData)),
		Memsz:  uint64(len(dataData)),
		Align:  1,
	}

	// 开始写文件
	f, err := os.Create(*flagOutput)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// 写 ELF 头
	if _, err := f.Seek(ehOff, io.SeekStart); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(f, binary.LittleEndian, &eh); err != nil {
		log.Fatalf("write ehdr: %v", err)
	}

	// 写 Program Headers
	if _, err := f.Seek(phOff, io.SeekStart); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(f, binary.LittleEndian, &textPh); err != nil {
		log.Fatalf("write phdr(text): %v", err)
	}
	if err := binary.Write(f, binary.LittleEndian, &dataPh); err != nil {
		log.Fatalf("write phdr(data): %v", err)
	}

	// 写段内容
	writeAt(f, textOff, textData)
	writeAt(f, dataOff, dataData)
}

func genTextData() []byte {
	const start_pc = 0x80000000
	instList := []struct {
		pc  int64
		as  riscv.As
		arg *riscv.AsArgument
	}{
		{0, riscv.AAUIPC, &riscv.AsArgument{Rd: riscv.REG_A0, Imm: 0}},
		{0, riscv.AADDI, &riscv.AsArgument{Rd: riscv.REG_A0, Rs1: riscv.REG_A0, Imm: 60}},
		{0, riscv.ALBU, &riscv.AsArgument{Rd: riscv.REG_A1, Rs1: riscv.REG_A0, Imm: 0}},
		{0, riscv.ABEQ, &riscv.AsArgument{Rs1: riscv.REG_A1, Rs2: riscv.REG_ZERO, Imm: 24}},
		{0, riscv.ALUI, &riscv.AsArgument{Rd: riscv.REG_T0, Imm: 0x10000}},
		{0, riscv.AADDI, &riscv.AsArgument{Rd: riscv.REG_T0, Rs1: riscv.REG_T0, Imm: 0}},
		{0, riscv.ASB, &riscv.AsArgument{Rs1: riscv.REG_T0, Rs2: riscv.REG_A1, Imm: 0}},
		{0, riscv.AADDI, &riscv.AsArgument{Rd: riscv.REG_A0, Rs1: riscv.REG_A0, Imm: 1}},
		{0, riscv.AJAL, &riscv.AsArgument{Rd: riscv.REG_ZERO, Imm: -24}},
		{0, riscv.ALUI, &riscv.AsArgument{Rd: riscv.REG_T0, Imm: 0x100}},
		{0, riscv.AADDI, &riscv.AsArgument{Rd: riscv.REG_T0, Rs1: riscv.REG_T0, Imm: 0}},
		{0, riscv.ALUI, &riscv.AsArgument{Rd: riscv.REG_T1, Imm: 0x5}},
		{0, riscv.AADDI, &riscv.AsArgument{Rd: riscv.REG_T1, Rs1: riscv.REG_T1, Imm: 0x555}},
		{0, riscv.ASW, &riscv.AsArgument{Rs1: riscv.REG_T0, Rs2: riscv.REG_T1, Imm: 0}},
		{0, riscv.AJAL, &riscv.AsArgument{Rd: riscv.REG_ZERO, Imm: 0}},
	}

	var buf bytes.Buffer
	for i, inst := range instList {
		inst.pc = start_pc + int64(i)*4
		x, err := riscv.EncodeRV64(inst.as, inst.arg)
		if err != nil {
			log.Fatal(i, err)
		}
		err = binary.Write(&buf, binary.LittleEndian, x)
		if err != nil {
			log.Fatal(i, err)
		}
	}

	return buf.Bytes()
}
