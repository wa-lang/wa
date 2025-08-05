// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"sort"
	"strconv"
	"strings"
)

// 位置指针
type Pos int

// 包含文件名和行列号的完整位置信息
// 行列号都是从1开始
type Position struct {
	Filename string // 文件名, 可缺省
	Offset   int    // 偏移量, 从0开始
	Line     int    // 行位置, 从1开始
	Column   int    // 列位置, 以字节计算, 从1开始
}

// 有效的位置必须有一个有效的行位置
func (pos *Position) IsValid() bool { return pos.Line > 0 }

// 位置信息的字符串表示有以下几种格式:
//
//	file:line:column    完整的文件名和行列号位置
//	file:line           文件名和行位置
//	line:column         行列位置
//	line                只有行位置
//	file                只有文件名
//	-                   彻底无效
func (pos Position) String() string {
	var sb strings.Builder
	sb.WriteString(pos.Filename)
	if pos.IsValid() {
		if sb.Len() > 0 {
			sb.WriteByte(':')
		}
		sb.WriteString(strconv.Itoa(pos.Line))
		if pos.Column != 0 {
			sb.WriteByte(':')
			sb.WriteString(strconv.Itoa(pos.Column))
		}
	}
	if sb.Len() == 0 {
		sb.WriteByte('-')
	}
	return sb.String()
}

// 表示文件的文件名/内容大小/行位置表格索引
type File struct {
	name string // 文件名
	base Pos    // 基础偏移量
	size int    // 文件大小

	lines []int // 行的位置列表, 第一必须从0开始
}

func NewFile(name string, base Pos, size int) *File {
	return &File{name: name, base: base, size: size, lines: []int{0}}
}

func (f *File) Name() string   { return f.name }
func (f *File) Base() Pos      { return f.base }
func (f *File) Size() int      { return f.size }
func (f *File) LineCount() int { return len(f.lines) }

// 根据文件的内容解析每行位置
func (f *File) SetLinesForContent(content []byte) {
	var lines []int
	var linePos = 0 // 第一行的偏移位置
	for offset, b := range content {
		if linePos >= 0 {
			lines = append(lines, linePos)
		}
		linePos = -1
		if b == '\n' {
			// 遇到换行, 产生一个新的行偏移位置
			linePos = offset + 1
		}
	}
	// 覆盖之前的行位置信息
	f.lines = lines
}

// 查询第几行对于的偏移量, 不含base部分
func (f *File) LineStart(line int) int {
	if line < 1 {
		panic("illegal line number (line numbering starts at 1)")
	}
	if line > len(f.lines) {
		panic("illegal line number")
	}
	return f.lines[line-1]
}

// 根据 pos 查询行号, 包含 base 处理
func (f *File) Line(p Pos) int {
	return f.Position(p).Line
}

// 根据 pos 查询完整的位置信息, 包含 base 处理
func (f *File) Position(p Pos) (pos Position) {
	if int(p) < 0 || int(p) > f.size {
		panic("illegal Pos value")
	}
	pos.Filename = f.name
	pos.Offset = int(p - f.base)
	if i := sort.SearchInts(f.lines, pos.Offset); i > 0 {
		pos.Line, pos.Column = i, pos.Offset-f.lines[i]
	}
	return
}
