// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package link

import (
	"encoding/binary"
	"io"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/link/elf"
)

func _LinkELF_X64_linux(prog *abi.LinkedProgram) ([]byte, error) {
	var (
		ehOff = int64(0)
		phOff = int64(elf.ELF64HDRSIZE)
	)

	var eh elf.ElfHeader64
	copy(eh.Ident[:], elf.ELFMAG)
	eh.Ident[elf.EI_CLASS] = byte(elf.ELFCLASS64)   // 64位
	eh.Ident[elf.EI_DATA] = byte(elf.ELFDATA2LSB)   // 小端序
	eh.Ident[elf.EI_VERSION] = byte(elf.EV_CURRENT) // 文件版本
	eh.Type = uint16(elf.ET_EXEC)                   // 可执行程序
	eh.Machine = uint16(elf.EM_X86_64)              // CPU类型: X64
	eh.Version = uint32(elf.EV_CURRENT)             // ELF版本
	eh.Entry = uint64(prog.Entry)                   // 程序开始地址
	eh.Phoff = uint64(phOff)                        // 程序头位置
	eh.Shoff = 0                                    // 不写节区表
	eh.Flags = 0                                    // 处理器特殊标志
	eh.Ehsize = elf.ELF64HDRSIZE                    // 文件头大小
	eh.Phentsize = elf.ELF64PHDRSIZE                // 程序头大小
	eh.Phnum = 2                                    // 程序头表中表项的数量(text 和 data)
	eh.Shentsize = 0                                // 节头表中每一个表项的大小
	eh.Shnum = 0                                    // 节头表中表项的数量
	eh.Shstrndx = 0                                 // 节头表中与节名字表相对应的表项的索引

	// 最大的页大小
	const maxPageSize = 64 << 10

	// 程序头: .text (RX)
	textPh := elf.ElfProgHeader64{
		Type:   elf.PT_LOAD,
		Flags:  elf.PF_R | elf.PF_X,        // 可读+执行, 不可修改
		Off:    uint64(0),                  // 指令在 prog.Entry-prog.TextAddr 开始位置
		Vaddr:  uint64(prog.TextAddr),      // 虚拟内存地址
		Paddr:  uint64(prog.TextAddr),      // 物理内存地址
		Filesz: uint64(len(prog.TextData)), // 数据段文件大小
		Memsz:  uint64(len(prog.TextData)), // 内存大小
		Align:  maxPageSize,                // 设置为 64KB 页对齐
	}

	// 程序头: .data
	dataPh := elf.ElfProgHeader64{
		Type:   elf.PT_LOAD,
		Flags:  elf.PF_R | elf.PF_W,        // 可读写
		Off:    uint64(len(prog.TextData)), // 数据段 offset
		Vaddr:  uint64(prog.DataAddr),      // 虚拟内存地址
		Paddr:  uint64(prog.DataAddr),      // 物理内存地址
		Filesz: uint64(len(prog.DataData)), // 数据段文件大小
		Memsz:  uint64(len(prog.DataData)), // 内存大小
		Align:  maxPageSize,                // 设置为 64KB 页对齐
	}

	// 构造内存缓存
	var buf _MemBuffer

	// 写 ELF 头
	buf.Seek(ehOff, io.SeekStart)
	if err := binary.Write(&buf, binary.LittleEndian, &eh); err != nil {
		return nil, err
	}

	// 写程序头
	buf.Seek(int64(eh.Phoff), io.SeekStart)
	if err := binary.Write(&buf, binary.LittleEndian, &textPh); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, &dataPh); err != nil {
		return nil, err
	}

	// 填充指令段开头的数据
	// Entry 可能不是第一个指令
	const fileHeaderSpaceSize = elf.ELF64HDRSIZE + elf.ELF64PHDRSIZE*2
	assert((prog.Entry - prog.TextAddr) >= fileHeaderSpaceSize)
	assert(len(prog.TextData) > int(fileHeaderSpaceSize))

	// 写段内容
	buf.WriteAt(prog.TextData[fileHeaderSpaceSize:], fileHeaderSpaceSize)
	buf.WriteAt(prog.DataData, int64(len(prog.TextData)))

	// OK
	return buf.Bytes(), nil
}
