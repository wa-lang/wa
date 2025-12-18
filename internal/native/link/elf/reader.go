// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package elf

import (
	"errors"
	"io"
	"os"
)

// errorReader returns error from all operations.
type errorReader struct {
	error
}

func (r errorReader) Read(p []byte) (n int, err error) {
	return 0, r.error
}

func (r errorReader) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, r.error
}

func (r errorReader) Seek(offset int64, whence int) (int64, error) {
	return 0, r.error
}

func (r errorReader) Close() error {
	return r.error
}

// readSeekerFromReader converts an io.Reader into an io.ReadSeeker.
// In general Seek may not be efficient, but it is optimized for
// common cases such as seeking to the end to find the length of the
// data.
type readSeekerFromReader struct {
	reset  func() (io.Reader, error)
	r      io.Reader
	size   int64
	offset int64
}

func (r *readSeekerFromReader) start() {
	x, err := r.reset()
	if err != nil {
		r.r = errorReader{err}
	} else {
		r.r = x
	}
	r.offset = 0
}

func (r *readSeekerFromReader) Read(p []byte) (n int, err error) {
	if r.r == nil {
		r.start()
	}
	n, err = r.r.Read(p)
	r.offset += int64(n)
	return n, err
}

func (r *readSeekerFromReader) Seek(offset int64, whence int) (int64, error) {
	var newOffset int64
	switch whence {
	case io.SeekStart:
		newOffset = offset
	case io.SeekCurrent:
		newOffset = r.offset + offset
	case io.SeekEnd:
		newOffset = r.size + offset
	default:
		return 0, os.ErrInvalid
	}

	switch {
	case newOffset == r.offset:
		return newOffset, nil

	case newOffset < 0, newOffset > r.size:
		return 0, os.ErrInvalid

	case newOffset == 0:
		r.r = nil

	case newOffset == r.size:
		r.r = errorReader{io.EOF}

	default:
		if newOffset < r.offset {
			// Restart at the beginning.
			r.start()
		}
		// Read until we reach offset.
		var buf [512]byte
		for r.offset < newOffset {
			b := buf[:]
			if newOffset-r.offset < int64(len(buf)) {
				b = buf[:newOffset-r.offset]
			}
			if _, err := r.Read(b); err != nil {
				return 0, err
			}
		}
	}
	r.offset = newOffset
	return r.offset, nil
}

type nobitsSectionReader struct{}

func (*nobitsSectionReader) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, errors.New("unexpected read from SHT_NOBITS section")
}

// ReadDataAt reads n bytes from the input stream at off, but avoids
// allocating all n bytes if n is large. This avoids crashing the program
// by allocating all n bytes in cases where n is incorrect.
func saferio_ReadDataAt(r io.ReaderAt, n uint64, off int64) ([]byte, error) {
	const chunk = 10 << 20 // 10M

	if int64(n) < 0 || n != uint64(int(n)) {
		// n is too large to fit in int, so we can't allocate
		// a buffer large enough. Treat this as a read failure.
		return nil, io.ErrUnexpectedEOF
	}

	if n < chunk {
		buf := make([]byte, n)
		_, err := r.ReadAt(buf, off)
		if err != nil {
			// io.SectionReader can return EOF for n == 0,
			// but for our purposes that is a success.
			if err != io.EOF || n > 0 {
				return nil, err
			}
		}
		return buf, nil
	}

	var buf []byte
	buf1 := make([]byte, chunk)
	for n > 0 {
		next := n
		if next > chunk {
			next = chunk
		}
		_, err := r.ReadAt(buf1[:next], off)
		if err != nil {
			return nil, err
		}
		buf = append(buf, buf1[:next]...)
		n -= next
		off += int64(next)
	}
	return buf, nil
}
