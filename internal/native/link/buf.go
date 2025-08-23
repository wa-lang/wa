// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package link

import (
	"errors"
	"fmt"
	"io"
)

var (
	_ io.WriterAt    = (*_MemBuffer)(nil)
	_ io.WriteSeeker = (*_MemBuffer)(nil)
)

type _MemBuffer struct {
	buf []byte
	off int64
}

func (m *_MemBuffer) Seek(offset int64, whence int) (int64, error) {
	var newOff int64
	switch whence {
	case io.SeekStart:
		newOff = offset
	case io.SeekCurrent:
		newOff = m.off + offset
	case io.SeekEnd:
		newOff = int64(len(m.buf)) + offset
	}
	if newOff < 0 {
		return 0, fmt.Errorf("negative position")
	}
	m.off = newOff
	return newOff, nil
}

func (m *_MemBuffer) Write(p []byte) (int, error) {
	if n := int(m.off) + len(p); n > cap(m.buf) {
		m.grow(n)
	}
	copy(m.buf[m.off:], p)
	m.off += int64(len(p))
	return len(p), nil
}

func (m *_MemBuffer) WriteAt(p []byte, off int64) (int, error) {
	if off < 0 {
		return 0, errors.New("negative offset")
	}
	if n := int(off) + len(p); n > cap(m.buf) {
		m.grow(n)
	}
	copy(m.buf[off:], p)
	return len(p), nil
}

func (m *_MemBuffer) Bytes() []byte {
	return m.buf
}

func (m *_MemBuffer) grow(n int) {
	newBuf := make([]byte, n)
	copy(newBuf, m.buf)
	m.buf = newBuf
}
