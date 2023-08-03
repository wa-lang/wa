// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm

import (
	"encoding/binary"
	"log"
	"os"
)

var zeros [512]byte

type OutBuf struct {
	off int64

	buf  []byte // backing store of mmap'd output file
	heap []byte // backing store for non-mmapped data

	name   string
	f      *os.File
	encbuf [8]byte // temp buffer used by WriteN methods
	isView bool    // true if created from View()
}

// Data returns the whole written OutBuf as a byte slice.
func (out *OutBuf) Data() []byte {
	return out.heap
}

// copyHeap copies the heap to the mmapped section of memory, returning true if
// a copy takes place.
func (out *OutBuf) copyHeap() bool {
	return false
}

// maxOutBufHeapLen limits the growth of the heap area.
const maxOutBufHeapLen = 10 << 20

// writeLoc determines the write location if a buffer is mmaped.
// We maintain two write buffers, an mmapped section, and a heap section for
// writing. When the mmapped section is full, we switch over the heap memory
// for writing.
func (out *OutBuf) writeLoc(lenToWrite int64) (int64, []byte) {
	// See if we have enough space in the mmaped area.
	bufLen := int64(len(out.buf))
	if out.off+lenToWrite <= bufLen {
		return out.off, out.buf
	}

	// Not enough space in the mmaped area, write to heap area instead.
	heapPos := out.off - bufLen
	heapLen := int64(len(out.heap))
	lenNeeded := heapPos + lenToWrite
	if lenNeeded > heapLen { // do we need to grow the heap storage?
		// The heap variables aren't protected by a mutex. For now, just bomb if you
		// try to use OutBuf in parallel. (Note this probably could be fixed.)
		if out.isView {
			panic("cannot write to heap in parallel")
		}
		// See if our heap would grow to be too large, and if so, copy it to the end
		// of the mmapped area.
		if heapLen > maxOutBufHeapLen && out.copyHeap() {
			heapPos -= heapLen
			lenNeeded = heapPos + lenToWrite
			heapLen = 0
		}
		out.heap = append(out.heap, make([]byte, lenNeeded-heapLen)...)
	}
	return heapPos, out.heap
}

func (out *OutBuf) SeekSet(p int64) {
	out.off = p
}

func (out *OutBuf) Offset() int64 {
	return out.off
}

// Write writes the contents of v to the buffer.
func (out *OutBuf) Write(v []byte) (int, error) {
	n := len(v)
	pos, buf := out.writeLoc(int64(n))
	copy(buf[pos:], v)
	out.off += int64(n)
	return n, nil
}

func (out *OutBuf) Write8(v uint8) {
	pos, buf := out.writeLoc(1)
	buf[pos] = v
	out.off++
}

// WriteByte is an alias for Write8 to fulfill the io.ByteWriter interface.
func (out *OutBuf) WriteByte(v byte) error {
	out.Write8(v)
	return nil
}

func (out *OutBuf) Write16(v uint16) {
	binary.LittleEndian.PutUint16(out.encbuf[:], v)
	out.Write(out.encbuf[:2])
}

func (out *OutBuf) Write32(v uint32) {
	binary.LittleEndian.PutUint32(out.encbuf[:], v)
	out.Write(out.encbuf[:4])
}

func (out *OutBuf) Write32b(v uint32) {
	binary.BigEndian.PutUint32(out.encbuf[:], v)
	out.Write(out.encbuf[:4])
}

func (out *OutBuf) Write64(v uint64) {
	binary.LittleEndian.PutUint64(out.encbuf[:], v)
	out.Write(out.encbuf[:8])
}

func (out *OutBuf) WriteString(s string) {
	pos, buf := out.writeLoc(int64(len(s)))
	n := copy(buf[pos:], s)
	if n != len(s) {
		log.Fatalf("WriteString truncated. buffer size: %d, offset: %d, len(s)=%d", len(out.buf), out.off, len(s))
	}
	out.off += int64(n)
}

// WriteStringN writes the first n bytes of s.
// If n is larger than len(s) then it is padded with zero bytes.
func (out *OutBuf) WriteStringN(s string, n int) {
	out.WriteStringPad(s, n, zeros[:])
}

// WriteStringPad writes the first n bytes of s.
// If n is larger than len(s) then it is padded with the bytes in pad (repeated as needed).
func (out *OutBuf) WriteStringPad(s string, n int, pad []byte) {
	if len(s) >= n {
		out.WriteString(s[:n])
	} else {
		out.WriteString(s)
		n -= len(s)
		for n > len(pad) {
			out.Write(pad)
			n -= len(pad)

		}
		out.Write(pad[:n])
	}
}
