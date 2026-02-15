// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package link

import (
	"bytes"
	"encoding/binary"
	"time"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/link/pe"
)

// 生成 PE 格式可执行文件
func LinkEXE(prog *abi.LinkedProgram) ([]byte, error) {
	var buf bytes.Buffer

	// 加载镜像基础地址都是固定的
	assert(prog.X64ImageBase == pe.DefaultImageBase)

	// 1. DOS Header (64 bytes)
	dosHeader := make([]byte, 64)
	dosHeader[0], dosHeader[1] = 'M', 'Z'
	binary.LittleEndian.PutUint32(dosHeader[0x3C:], 0x40)
	buf.Write(dosHeader)

	// 2. PE Signature ("PE\0\0")
	buf.Write([]byte{'P', 'E', 0, 0})

	// 3. File Header
	fileHeader := pe.FileHeader{
		Machine:              pe.IMAGE_FILE_MACHINE_AMD64,
		NumberOfSections:     2, // .text 和 .data
		TimeDateStamp:        uint32(time.Now().Unix()),
		SizeOfOptionalHeader: uint16(binary.Size(pe.OptionalHeader64{})),
		Characteristics:      pe.IMAGE_FILE_EXECUTABLE_IMAGE | pe.IMAGE_FILE_LARGE_ADDRESS_AWARE,
	}
	binary.Write(&buf, binary.LittleEndian, fileHeader)

	// 4. 计算各段在文件中的偏移 (Header 之后)
	// Header 大小 = DOS(64) + PE(4) + FileHdr(20) + OptHdr(240) + 2*SecHdr(40) = 408
	// 将其对齐到 FileAlignment (512)
	headerSize := uint32(0x40 + 4 + 20 + binary.Size(pe.OptionalHeader64{}) + 2*40)
	fileSizeHeader := align(headerSize, pe.DefaultFileAlignment)

	textFileSize := align(uint32(len(prog.TextData)), pe.DefaultFileAlignment)
	dataFileSize := align(uint32(len(prog.DataData)), pe.DefaultFileAlignment)

	// 5. Optional Header
	optHeader := pe.OptionalHeader64{
		Magic:                 pe.PE64Magic,
		AddressOfEntryPoint:   uint32(prog.TextAddr - prog.X64ImageBase), // 假设.text在0x1000
		ImageBase:             pe.DefaultImageBase,
		SectionAlignment:      pe.DefaultSectionAlignment,
		FileAlignment:         pe.DefaultFileAlignment,
		MajorSubsystemVersion: 6,
		MinorSubsystemVersion: 0, // Windows Vista+
		SizeOfImage:           align(0x1000+uint32(len(prog.TextData)), pe.DefaultSectionAlignment) + align(uint32(len(prog.DataData)), pe.DefaultSectionAlignment),
		SizeOfHeaders:         fileSizeHeader,
		Subsystem:             pe.IMAGE_SUBSYSTEM_WINDOWS_CUI,
		DllCharacteristics:    0,
	}
	binary.Write(&buf, binary.LittleEndian, optHeader)

	// 6. Section Headers
	// .text 段
	textHeader := pe.SectionHeader64{
		Name:            [8]byte{'.', 't', 'e', 'x', 't'},
		VirtualSize:     uint32(len(prog.TextData)),
		VirtualAddress:  uint32(prog.TextAddr - prog.X64ImageBase),
		Size:            textFileSize,
		Offset:          fileSizeHeader,
		Characteristics: pe.IMAGE_SCN_CNT_CODE | pe.IMAGE_SCN_MEM_EXECUTE | pe.IMAGE_SCN_MEM_READ,
	}
	binary.Write(&buf, binary.LittleEndian, textHeader)

	// .data 段
	// TODO: 断言地址已经对齐
	dataHeader := pe.SectionHeader64{
		Name:            [8]byte{'.', 'd', 'a', 't', 'a'},
		VirtualSize:     uint32(len(prog.DataData)),
		VirtualAddress:  uint32(prog.DataAddr - prog.X64ImageBase),
		Size:            dataFileSize,
		Offset:          fileSizeHeader + textFileSize,
		Characteristics: pe.IMAGE_SCN_CNT_INITIALIZED_DATA | pe.IMAGE_SCN_MEM_READ | pe.IMAGE_SCN_MEM_WRITE,
	}
	binary.Write(&buf, binary.LittleEndian, dataHeader)

	// 7. 填充 Header 剩余空间到 FileAlignment
	buf.Write(make([]byte, int(fileSizeHeader)-buf.Len()))

	// 8. 写入段数据并对齐
	buf.Write(prog.TextData)
	buf.Write(make([]byte, int(textFileSize)-len(prog.TextData)))

	buf.Write(prog.DataData)
	buf.Write(make([]byte, int(dataFileSize)-len(prog.DataData)))

	return buf.Bytes(), nil
}

func align(n, a uint32) uint32 {
	if a > 0 && n%a != 0 {
		return (n/a + 1) * a
	}
	return n
}
