// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package pe

import (
	"encoding/binary"
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

type ImportDirectory struct {
	OriginalFirstThunk uint32
	TimeDateStamp      uint32
	ForwarderChain     uint32
	Name               uint32
	FirstThunk         uint32

	dll string
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

	var dosheader [96]byte
	if _, err := r.ReadAt(dosheader[0:], 0); err != nil {
		return nil, err
	}
	var base int64
	if dosheader[0] == 'M' && dosheader[1] == 'Z' {
		signoff := int64(binary.LittleEndian.Uint32(dosheader[0x3c:]))
		var sign [4]byte
		r.ReadAt(sign[:], signoff)
		if !(sign[0] == 'P' && sign[1] == 'E' && sign[2] == 0 && sign[3] == 0) {
			return nil, fmt.Errorf("invalid PE file signature: % x", sign)
		}
		base = signoff + 4
	} else {
		base = int64(0)
	}
	sr.Seek(base, io.SeekStart)
	if err := binary.Read(sr, binary.LittleEndian, &f.FileHeader); err != nil {
		return nil, err
	}
	switch f.FileHeader.Machine {
	case IMAGE_FILE_MACHINE_AMD64,
		IMAGE_FILE_MACHINE_I386,
		IMAGE_FILE_MACHINE_UNKNOWN:
		// ok
	default:
		return nil, fmt.Errorf("unrecognized PE machine: %#x", f.FileHeader.Machine)
	}

	var err error

	// Seek past file header.
	_, err = sr.Seek(base+int64(binary.Size(f.FileHeader)), io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Read optional header.
	f.OptionalHeader, err = readOptionalHeader(sr, f.FileHeader.SizeOfOptionalHeader)
	if err != nil {
		return nil, err
	}

	// Process sections.
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
	// If optional header size is 0, return empty optional header.
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

	if !read(&ohMagic) {
		return nil, fmt.Errorf("failure to read optional header magic: %v", err)

	}

	switch ohMagic {
	case PE64Magic: // PE32+
		var (
			oh64 OptionalHeader64
			// There can be 0 or more data directories. So the minimum size of optional
			// header is calculated by subtracting oh64.DataDirectory size from oh64 size.
			oh64MinSz = binary.Size(oh64) - binary.Size(oh64.DataDirectory)
		)

		if sz < uint16(oh64MinSz) {
			return nil, fmt.Errorf("optional header size(%d) is less minimum size (%d) for PE32+ optional header", sz, oh64MinSz)
		}

		// Init oh64 fields
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

		dd, err := readDataDirectories(r, sz-uint16(oh64MinSz), oh64.NumberOfRvaAndSizes)
		if err != nil {
			return nil, err
		}

		copy(oh64.DataDirectory[:], dd)

		return &oh64, nil
	default:
		return nil, fmt.Errorf("optional header has unexpected Magic of 0x%x", ohMagic)
	}
}

func (f *File) ImportedSymbols() ([]string, error) {
	if f.OptionalHeader == nil {
		return nil, nil
	}

	// grab the number of data directory entries
	var dd_length uint32
	dd_length = f.OptionalHeader.NumberOfRvaAndSizes

	// check that the length of data directory entries is large
	// enough to include the imports directory.
	if dd_length < IMAGE_DIRECTORY_ENTRY_IMPORT+1 {
		return nil, nil
	}

	// grab the import data directory entry
	var idd DataDirectory
	idd = f.OptionalHeader.DataDirectory[IMAGE_DIRECTORY_ENTRY_IMPORT]

	// figure out which section contains the import directory table
	var ds *Section
	ds = nil
	for _, s := range f.Sections {
		if s.Offset == 0 {
			continue
		}
		// We are using distance between s.VirtualAddress and idd.VirtualAddress
		// to avoid potential overflow of uint32 caused by addition of s.VirtualSize
		// to s.VirtualAddress.
		if s.VirtualAddress <= idd.VirtualAddress && idd.VirtualAddress-s.VirtualAddress < s.VirtualSize {
			ds = s
			break
		}
	}

	// didn't find a section, so no import libraries were found
	if ds == nil {
		return nil, nil
	}

	d, err := ds.Data()
	if err != nil {
		return nil, err
	}

	// seek to the virtual address specified in the import data directory
	d = d[idd.VirtualAddress-ds.VirtualAddress:]

	// start decoding the import directory
	var ida []ImportDirectory
	for len(d) >= 20 {
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
			if va&0x8000000000000000 > 0 { // is Ordinal
				// TODO: add dynimport ordinal support.
			} else {
				fn, _ := getString(names, int(uint32(va)-ds.VirtualAddress+2))
				all = append(all, fn+":"+dt.dll)
			}
		}
	}

	return all, nil
}

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
