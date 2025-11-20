// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package esp32

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

type File struct {
	Head         ImageHeader
	Segments     []SegmentHeader
	SegmentsData [][]byte
	Checksum     uint8
	SHA256       []byte
}

func Load(path string) (*File, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	f := new(File)
	if err := binary.Read(r, binary.LittleEndian, &f.Head); err != nil {
		return nil, err
	}

	for i := 0; i < int(f.Head.SegmentCount); i++ {
		var segHdr SegmentHeader
		if err := binary.Read(r, binary.LittleEndian, &segHdr); err != nil {
			return f, err
		}

		data := make([]byte, segHdr.DataLength)

		f.Segments = append(f.Segments, segHdr)
		f.SegmentsData = append(f.SegmentsData, data)

		if _, err := r.Read(data); err != nil {
			return f, err
		}
	}

	// 读取 Checksum
	{
		// 计算文件大小
		totalSize := int64(ImageHeaderSize)
		for _, data := range f.SegmentsData {
			totalSize += SegmentHeaderSize + int64(len(data))
		}

		// 计算偏移量
		aligned := (totalSize + 15) &^ 15
		checksumOffset := aligned - 1

		if _, err := r.Seek(checksumOffset, io.SeekStart); err != nil {
			return nil, err
		}

		var cs [1]byte
		if _, err := io.ReadFull(r, cs[:]); err != nil {
			return nil, err
		}
		f.Checksum = cs[0]
	}

	// 读取 sha256
	if f.Head.HashAppended != 0 {
		f.SHA256 = make([]byte, 32)
		if _, err := io.ReadFull(r, f.SHA256); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (f *File) Dump() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "ESP32 Image Info\n")
	fmt.Fprintf(&sb, "Magic:        0x%02X\n", f.Head.Magic)
	fmt.Fprintf(&sb, "SegmentCount: %d\n", f.Head.SegmentCount)
	fmt.Fprintf(&sb, "EntryAddr:    0x%08X\n", f.Head.EntryAddr)
	fmt.Fprintf(&sb, "WP Pin:       0x%02X\n", f.Head.WPPin)
	fmt.Fprintf(&sb, "ClkQio:       0x%02X\n", f.Head.SpiMode)
	fmt.Fprintf(&sb, "\n")

	fmt.Fprintf(&sb, "Segments:\n")
	for i, seg := range f.Segments {
		fmt.Fprintf(&sb, "  #%d LoadAddr=0x%08X Length=%d (0x%X)\n",
			i, seg.LoadAddr, seg.DataLength, seg.DataLength)
	}
	fmt.Fprintf(&sb, "\n")

	// checksum
	calc := f.ComputeChecksum()
	ok := ""
	if calc == f.Checksum {
		ok = " (OK)"
	} else {
		ok = " (MISMATCH)"
	}

	fmt.Fprintf(&sb, "Checksum:\n")
	fmt.Fprintf(&sb, "  File:     0x%02X\n", f.Checksum)
	fmt.Fprintf(&sb, "  Computed: 0x%02X%s\n", calc, ok)
	fmt.Fprintf(&sb, "\n")

	// sha256 hash
	if len(f.SHA256) > 0 {
		fmt.Fprintf(&sb, "SHA256 (image):\n")
		hexStr := strings.ToUpper(hex.EncodeToString(f.SHA256))

		// 分组显示：每 32 hex = 16 bytes = 一行
		for i := 0; i < len(hexStr); i += 32 {
			end := i + 32
			if end > len(hexStr) {
				end = len(hexStr)
			}
			fmt.Fprintf(&sb, "  %s\n", hexStr[i:end])
		}
	}

	return sb.String()
}

// 计算 XOR 校验和(不包含段头部数据)
func (f *File) ComputeChecksum() uint8 {
	cs := uint8(0)
	for _, data := range f.SegmentsData {
		for _, b := range data {
			cs ^= b
		}
	}
	cs ^= 0xEF
	return cs
}
