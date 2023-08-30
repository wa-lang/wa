package wat

import "bytes"

func makeAlign(i, align int) int {
	if align == 1 || align == 0 {
		return i
	}
	return (i + align - 1) / align * align
}

/**************************************
DataSeg:
**************************************/
type DataSeg struct {
	start int
	data  []byte
}

func NewDataSeg(start int) *DataSeg {
	return &DataSeg{start: start}
}

func (s *DataSeg) Append(data []byte, align int) (ptr int) {
	ptr = bytes.Index(s.data, data)
	if ptr != -1 {
		ptr += s.start
		return
	}

	ptr = s.Alloc(len(data), align)
	s.Set(data, ptr)
	return
}

func (s *DataSeg) Alloc(size, align int) (ptr int) {
	p := s.start + len(s.data)
	ptr = makeAlign(p, align)
	d := ptr + size - p
	s.data = append(s.data, make([]byte, d)...)
	return
}

func (s *DataSeg) Set(data []byte, ptr int) {
	if copy(s.data[ptr-s.start:], data) != len(data) {
		panic("len(dst) < len(src)")
	}
}

func (s *DataSeg) Size() int {
	return s.start + len(s.data)
}
