// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package bio

import (
	"bufio"
	"io"
	"log"
	"os"
)

const Beof = -1

type Biobuf struct {
	f       *os.File
	r       *bufio.Reader
	w       *bufio.Writer
	linelen int
}

func Bopenw(name string) (*Biobuf, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return &Biobuf{f: f, w: bufio.NewWriter(f)}, nil
}

func Bopenr(name string) (*Biobuf, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return &Biobuf{f: f, r: bufio.NewReader(f)}, nil
}

func Binitw(w io.Writer) *Biobuf {
	return &Biobuf{w: bufio.NewWriter(w)}
}

func (b *Biobuf) Write(p []byte) (int, error) {
	return b.w.Write(p)
}

func (b *Biobuf) WriteString(s string) (int, error) {
	return b.w.WriteString(s)
}

func (b *Biobuf) Bseek(offset int64, whence int) int64 {
	if b.w != nil {
		if err := b.w.Flush(); err != nil {
			log.Fatalf("writing output: %v", err)
		}
	} else if b.r != nil {
		if whence == 1 {
			offset -= int64(b.r.Buffered())
		}
	}
	off, err := b.f.Seek(offset, whence)
	if err != nil {
		log.Fatalf("seeking in output: %v", err)
	}
	if b.r != nil {
		b.r.Reset(b.f)
	}
	return off
}

func (b *Biobuf) Boffset() int64 {
	if b.w != nil {
		if err := b.w.Flush(); err != nil {
			log.Fatalf("writing output: %v", err)
		}
	}
	off, err := b.f.Seek(0, 1)
	if err != nil {
		log.Fatalf("seeking in output [0, 1]: %v", err)
	}
	if b.r != nil {
		off -= int64(b.r.Buffered())
	}
	return off
}

func (b *Biobuf) Flush() error {
	return b.w.Flush()
}

func (b *Biobuf) Bputc(c byte) {
	b.w.WriteByte(c)
}

func (b *Biobuf) Bread(p []byte) int {
	n, err := io.ReadFull(b.r, p)
	if n == 0 {
		if err != nil && err != io.EOF {
			n = -1
		}
	}
	return n
}

func (b *Biobuf) Bgetc() int {
	c, err := b.r.ReadByte()
	if err != nil {
		return -1
	}
	return int(c)
}

func (b *Biobuf) Bgetrune() int {
	r, _, err := b.r.ReadRune()
	if err != nil {
		return -1
	}
	return int(r)
}

func (b *Biobuf) Bungetrune() {
	b.r.UnreadRune()
}

func (b *Biobuf) Read(p []byte) (int, error) {
	return b.r.Read(p)
}

func (b *Biobuf) Peek(n int) ([]byte, error) {
	return b.r.Peek(n)
}

func (b *Biobuf) Brdline(delim int) string {
	s, err := b.r.ReadBytes(byte(delim))
	if err != nil {
		log.Fatalf("reading input: %v", err)
	}
	b.linelen = len(s)
	return string(s)
}

func (b *Biobuf) Brdstr(delim int, cut int) string {
	s, err := b.r.ReadString(byte(delim))
	if err != nil {
		log.Fatalf("reading input: %v", err)
	}
	if len(s) > 0 && cut > 0 {
		s = s[:len(s)-1]
	}
	return s
}

func (b *Biobuf) Blinelen() int {
	return b.linelen
}

func (b *Biobuf) Bterm() error {
	var err error
	if b.w != nil {
		err = b.w.Flush()
	}
	err1 := b.f.Close()
	if err == nil {
		err = err1
	}
	return err
}
