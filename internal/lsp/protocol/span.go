// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// this file contains protocol<->span converters

package protocol

import (
	"fmt"
	"unicode/utf8"

	"wa-lang.org/wa/internal/lsp/span"
)

// CompareLocation defines a three-valued comparison over locations,
// lexicographically ordered by (URI, Range).
func CompareLocation(x, y Location) int {
	if x.URI != y.URI {
		if x.URI < y.URI {
			return -1
		} else {
			return +1
		}
	}
	return CompareRange(x.Range, y.Range)
}

// Format implements fmt.Formatter.
//
// See Range.Format for discussion of why the Formatter interface is
// implemented rather than Stringer.
func (p Position) Format(f fmt.State, _ rune) {
	fmt.Fprintf(f, "%v:%v", p.Line, p.Character)
}

// -- implementation helpers --

// UTF16Len returns the number of codes in the UTF-16 transcoding of s.
func UTF16Len(s []byte) int {
	var n int
	for len(s) > 0 {
		n++

		// Fast path for ASCII.
		if s[0] < 0x80 {
			s = s[1:]
			continue
		}

		r, size := utf8.DecodeRune(s)
		if r >= 0x10000 {
			n++ // surrogate pair
		}
		s = s[size:]
	}
	return n
}

type ColumnMapper struct {
	URI       span.URI
	Converter *span.TokenConverter
	Content   []byte
}

func URIFromSpanURI(uri span.URI) DocumentURI {
	return DocumentURI(uri)
}

func (u DocumentURI) SpanURI() span.URI {
	return span.URIFromURI(string(u))
}

func (m *ColumnMapper) Location(s span.Span) (Location, error) {
	rng, err := m.Range(s)
	if err != nil {
		return Location{}, err
	}
	return Location{URI: URIFromSpanURI(s.URI()), Range: rng}, nil
}

func (m *ColumnMapper) Range(s span.Span) (Range, error) {
	if span.CompareURI(m.URI, s.URI()) != 0 {
		return Range{}, fmt.Errorf("column mapper is for file %q instead of %q", m.URI, s.URI())
	}
	s, err := s.WithAll(m.Converter)
	if err != nil {
		return Range{}, err
	}
	start, err := m.Position(s.Start())
	if err != nil {
		return Range{}, err
	}
	end, err := m.Position(s.End())
	if err != nil {
		return Range{}, err
	}
	return Range{Start: start, End: end}, nil
}

func (m *ColumnMapper) Position(p span.Point) (Position, error) {
	chr, err := span.ToUTF16Column(p, m.Content)
	if err != nil {
		return Position{}, err
	}
	return Position{
		Line:      uint32(p.Line() - 1),
		Character: uint32(chr - 1),
	}, nil
}

func (m *ColumnMapper) Span(l Location) (span.Span, error) {
	return m.RangeSpan(l.Range)
}

func (m *ColumnMapper) RangeSpan(r Range) (span.Span, error) {
	start, err := m.Point(r.Start)
	if err != nil {
		return span.Span{}, err
	}
	end, err := m.Point(r.End)
	if err != nil {
		return span.Span{}, err
	}
	return span.New(m.URI, start, end).WithAll(m.Converter)
}

func (m *ColumnMapper) RangeToSpanRange(r Range) (span.Range, error) {
	spn, err := m.RangeSpan(r)
	if err != nil {
		return span.Range{}, err
	}
	return spn.Range(m.Converter)
}

func (m *ColumnMapper) PointSpan(p Position) (span.Span, error) {
	start, err := m.Point(p)
	if err != nil {
		return span.Span{}, err
	}
	return span.New(m.URI, start, start).WithAll(m.Converter)
}

func (m *ColumnMapper) Point(p Position) (span.Point, error) {
	line := int(p.Line) + 1
	offset, err := m.Converter.ToOffset(line, 1)
	if err != nil {
		return span.Point{}, err
	}
	lineStart := span.NewPoint(line, 1, offset)
	return span.FromUTF16Column(lineStart, int(p.Character)+1, m.Content)
}

func IsPoint(r Range) bool {
	return r.Start.Line == r.End.Line && r.Start.Character == r.End.Character
}

func CompareRange(a, b Range) int {
	if r := ComparePosition(a.Start, b.Start); r != 0 {
		return r
	}
	return ComparePosition(a.End, b.End)
}

func ComparePosition(a, b Position) int {
	if a.Line < b.Line {
		return -1
	}
	if a.Line > b.Line {
		return 1
	}
	if a.Character < b.Character {
		return -1
	}
	if a.Character > b.Character {
		return 1
	}
	return 0
}

func Intersect(a, b Range) bool {
	if a.Start.Line > b.End.Line || a.End.Line < b.Start.Line {
		return false
	}
	return !((a.Start.Line == b.End.Line) && a.Start.Character > b.End.Character ||
		(a.End.Line == b.Start.Line) && a.End.Character < b.Start.Character)
}

func (r Range) Format(f fmt.State, _ rune) {
	fmt.Fprintf(f, "%v:%v-%v:%v", r.Start.Line, r.Start.Character, r.End.Line, r.End.Character)
}
