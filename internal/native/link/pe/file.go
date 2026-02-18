// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package pe

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

type File struct {
	FileHeader
	OptionalHeader *OptionalHeader64
	Sections       []*Section
	closer         io.Closer
}

type Section struct {
	SectionHeader64
	sr *io.SectionReader
}

func (s *Section) Data() ([]byte, error) {
	if s.sr == nil || s.Size == 0 {
		return nil, nil
	}
	buf := make([]byte, s.Size)
	n, err := s.sr.ReadAt(buf, 0)
	if err == io.EOF && n == int(s.Size) {
		err = nil
	}
	return buf[:n], err
}

func OpenPE64(name string) (*File, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	ff, err := NewFilePE64(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	ff.closer = f
	return ff, nil
}

func (f *File) Close() error {
	var err error
	if f.closer != nil {
		err = f.closer.Close()
		f.closer = nil
	}
	return err
}

func NewFilePE64(r io.ReaderAt) (*File, error) {
	f := new(File)
	sr := io.NewSectionReader(r, 0, 1<<63-1)

	// 读取 DOS 头部
	// MS-DOS Stub 是可选的, 读取时忽略
	var dosheader [KDosHeaderSize]byte
	if n, err := r.ReadAt(dosheader[0:], 0); err != nil {
		return nil, err
	} else if n < KDosHeaderSize {
		return nil, errors.New("invalid PE file")
	}
	if dosheader[0] != 'M' || dosheader[1] != 'Z' {
		return nil, errors.New("invalid PE file")
	}

	// 从 DOS 头部获取 NtHeader64 的偏移地址
	_NTHeaderOffset := int64(binary.LittleEndian.Uint32(dosheader[0x3c:]))

	// 验证 PE 头部
	var sign [4]byte
	if n, err := r.ReadAt(sign[:], _NTHeaderOffset); err != nil {
		return nil, fmt.Errorf("invalid PE file: %s", err)
	} else if n < 4 {
		return nil, fmt.Errorf("invalid PE file: read failed")
	}
	if sign[0] != 'P' && sign[1] != 'E' && sign[2] != 0 && sign[3] != 0 {
		return nil, fmt.Errorf("invalid PE file signature: % x", sign)
	}

	// 读取 NtHeader64.FileHeader
	sr.Seek(_NTHeaderOffset+4, io.SeekStart)
	if err := binary.Read(sr, binary.LittleEndian, &f.FileHeader); err != nil {
		return nil, err
	}

	// 判断 PE 程序对于的 CPU 类型
	// 仅支持 amd64 和 386 两个类型
	if m := f.FileHeader.Machine; m != IMAGE_FILE_MACHINE_AMD64 && m != IMAGE_FILE_MACHINE_I386 {
		return nil, fmt.Errorf("unrecognized PE machine: %#x", f.FileHeader.Machine)
	}

	// 读取 OptionalHeader64
	_, err := sr.Seek(_NTHeaderOffset+4+int64(binary.Size(f.FileHeader)), io.SeekStart)
	if err != nil {
		return nil, err
	}
	f.OptionalHeader, err = readOptionalHeader(sr, f.FileHeader.SizeOfOptionalHeader)
	if err != nil {
		return nil, err
	}

	// 读取全部的段头表
	// 这里是真正的数据
	f.Sections = make([]*Section, f.FileHeader.NumberOfSections)
	for i := 0; i < int(f.FileHeader.NumberOfSections); i++ {
		sh := new(SectionHeader64)
		if err := binary.Read(sr, binary.LittleEndian, sh); err != nil {
			return nil, err
		}
		s := new(Section)
		s.SectionHeader64 = SectionHeader64{
			VirtualSize:          sh.VirtualSize,
			VirtualAddress:       sh.VirtualAddress,
			Size:                 sh.Size,
			Offset:               sh.Offset,
			PointerToRelocations: sh.PointerToRelocations,
			PointerToLineNumbers: sh.PointerToLineNumbers,
			NumberOfRelocations:  sh.NumberOfRelocations,
			NumberOfLineNumbers:  sh.NumberOfLineNumbers,
			Characteristics:      sh.Characteristics,
		}
		if sh.Offset != 0 {
			s.sr = io.NewSectionReader(r, int64(s.SectionHeader64.Offset), int64(s.SectionHeader64.Size))
		} else {
			s.sr = nil
		}
		f.Sections[i] = s
	}

	return f, nil
}

func readOptionalHeader(r io.ReadSeeker, sz uint16) (*OptionalHeader64, error) {
	if sz == 0 {
		return nil, nil
	}

	var (
		// First couple of bytes in option header state its type.
		// We need to read them first to determine the type and
		// validity of optional header.
		ohMagic   uint16
		ohMagicSz = binary.Size(ohMagic)
	)

	// If optional header size is greater than 0 but less than its magic size, return error.
	if sz < uint16(ohMagicSz) {
		return nil, fmt.Errorf("optional header size is less than optional header magic size")
	}

	// read reads from io.ReadSeeke, r, into data.
	var err error
	read := func(data interface{}) bool {
		err = binary.Read(r, binary.LittleEndian, data)
		return err == nil
	}

	// PE32+
	if !read(&ohMagic) {
		return nil, fmt.Errorf("failure to read optional header magic: %v", err)
	}
	if ohMagic != PE64Magic {
		return nil, fmt.Errorf("optional header has unexpected Magic of 0x%x", ohMagic)
	}

	var (
		oh64 OptionalHeader64
		// There can be 0 or more data directories. So the minimum size of optional
		// header is calculated by subtracting oh64.DataDirectory size from oh64 size.
		oh64MinSz = binary.Size(oh64) - binary.Size(oh64.DataDirectory)
	)

	// 判断有问题, ohMagic 已经被读取过了
	if sz < uint16(oh64MinSz) {
		return nil, fmt.Errorf("optional header size(%d) is less minimum size (%d) for PE32+ optional header", sz, oh64MinSz)
	}

	// 读取 OptionalHeader64
	// 不包含 Magic 和 DataDirectory 部分
	oh64.Magic = ohMagic
	if !read(&oh64.MajorLinkerVersion) ||
		!read(&oh64.MinorLinkerVersion) ||
		!read(&oh64.SizeOfCode) ||
		!read(&oh64.SizeOfInitializedData) ||
		!read(&oh64.SizeOfUninitializedData) ||
		!read(&oh64.AddressOfEntryPoint) ||
		!read(&oh64.BaseOfCode) ||
		!read(&oh64.ImageBase) ||
		!read(&oh64.SectionAlignment) ||
		!read(&oh64.FileAlignment) ||
		!read(&oh64.MajorOperatingSystemVersion) ||
		!read(&oh64.MinorOperatingSystemVersion) ||
		!read(&oh64.MajorImageVersion) ||
		!read(&oh64.MinorImageVersion) ||
		!read(&oh64.MajorSubsystemVersion) ||
		!read(&oh64.MinorSubsystemVersion) ||
		!read(&oh64.Win32VersionValue) ||
		!read(&oh64.SizeOfImage) ||
		!read(&oh64.SizeOfHeaders) ||
		!read(&oh64.CheckSum) ||
		!read(&oh64.Subsystem) ||
		!read(&oh64.DllCharacteristics) ||
		!read(&oh64.SizeOfStackReserve) ||
		!read(&oh64.SizeOfStackCommit) ||
		!read(&oh64.SizeOfHeapReserve) ||
		!read(&oh64.SizeOfHeapCommit) ||
		!read(&oh64.LoaderFlags) ||
		!read(&oh64.NumberOfRvaAndSizes) {
		return nil, fmt.Errorf("failure to read PE32+ optional header: %v", err)
	}

	// 读取 DataDirectory
	dd, err := readDataDirectories(r, sz-uint16(oh64MinSz), oh64.NumberOfRvaAndSizes)
	if err != nil {
		return nil, err
	}
	copy(oh64.DataDirectory[:], dd)

	return &oh64, nil
}

// 读取 DataDirectory
func readDataDirectories(r io.ReadSeeker, sz uint16, n uint32) ([]DataDirectory, error) {
	ddSz := uint64(binary.Size(DataDirectory{}))
	if uint64(sz) != uint64(n)*ddSz {
		return nil, fmt.Errorf("size of data directories(%d) is inconsistent with number of data directories(%d)", sz, n)
	}

	dd := make([]DataDirectory, n)
	if err := binary.Read(r, binary.LittleEndian, dd); err != nil {
		return nil, fmt.Errorf("failure to read data directories: %v", err)
	}

	return dd, nil
}

// 从 DataDirectory 读取导入符号信息
func (f *File) ImportedSymbols() ([]string, error) {
	if f.OptionalHeader == nil {
		return nil, nil
	}

	// 检查 DataDirectory 数量
	// 最后可能有一个空元素作为结束标志, 可能要多加一个
	if f.OptionalHeader.NumberOfRvaAndSizes < IMAGE_DIRECTORY_ENTRY_IMPORT+1 {
		return nil, nil
	}

	// 获取导入表对于的槽位
	idd := f.OptionalHeader.DataDirectory[IMAGE_DIRECTORY_ENTRY_IMPORT]

	// 查找导入表数据对应哪个段表格区间
	var ds *Section
	for _, s := range f.Sections {
		if s.Offset == 0 {
			continue
		}
		// 以下的计算是为了避免 uint32 溢出, 用于判断 idd 是否在此段数据区间
		// 问题: 并没有检查 idd.Size
		if s.VirtualAddress <= idd.VirtualAddress && idd.VirtualAddress-s.VirtualAddress < s.VirtualSize {
			ds = s
			break
		}
	}
	if ds == nil {
		return nil, nil
	}

	// 读取段数据
	d, err := ds.Data()
	if err != nil {
		return nil, err
	}

	// 跳过段数据开头空白区域
	d = d[idd.VirtualAddress-ds.VirtualAddress:]

	// 每个 ImportDirectory 头部位 20字节
	const kImportDirectoryMinSize = 20

	// 从段数据解析导入符号信息
	var ida []ImportDirectory
	for len(d) >= kImportDirectoryMinSize {
		var dt ImportDirectory
		dt.OriginalFirstThunk = binary.LittleEndian.Uint32(d[0:4])
		dt.TimeDateStamp = binary.LittleEndian.Uint32(d[4:8])
		dt.ForwarderChain = binary.LittleEndian.Uint32(d[8:12])
		dt.Name = binary.LittleEndian.Uint32(d[12:16])
		dt.FirstThunk = binary.LittleEndian.Uint32(d[16:20])
		d = d[20:]
		if dt.OriginalFirstThunk == 0 {
			break
		}
		ida = append(ida, dt)
	}

	names, _ := ds.Data()
	var all []string
	for _, dt := range ida {
		// 读取动态库的名字
		dt.dll, _ = getString(names, int(dt.Name-ds.VirtualAddress))
		d, _ = ds.Data()
		// seek to OriginalFirstThunk
		d = d[dt.OriginalFirstThunk-ds.VirtualAddress:]
		for len(d) > 0 {
			va := binary.LittleEndian.Uint64(d[0:8])
			d = d[8:]
			if va == 0 {
				break
			}
			if va&0x8000000000000000 > 0 {
				// TODO: 按序号导入函数
			} else {
				// 按照函数名导入
				// 格式 ExitProcess:KERNEL32.dll
				fn, _ := getString(names, int(uint32(va)-ds.VirtualAddress+2))
				all = append(all, fn+":"+dt.dll)
			}
		}
	}

	return all, nil
}

// 读取C风格的字符串, 以 \0 结尾
func getString(section []byte, start int) (string, bool) {
	if start < 0 || start >= len(section) {
		return "", false
	}

	for end := start; end < len(section); end++ {
		if section[end] == 0 {
			return string(section[start:end]), true
		}
	}
	return "", false
}
