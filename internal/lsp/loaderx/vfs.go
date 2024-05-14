// 版权 @2024 凹语言 作者。保留所有权利。

package loaderx

import (
	"io"
	"io/fs"
	"os"
	"time"

	"wa-lang.org/wa/internal/lsp/protocol"
)

var _ fs.FS = (*openedFS)(nil)

type openedFS struct {
	*Universe
}

func (p *openedFS) Open(name string) (fs.File, error) {
	if f, ok := p.Files[protocol.DocumentURI(name)]; ok {
		return &openedFile{
			path: name,
			openedFileInfo: openedFileInfo{
				name: name,
				f:    f,
			},
		}, nil
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

type openedFileInfo struct {
	name string
	f    *File
}

func (i *openedFileInfo) Name() string               { return i.name }
func (i *openedFileInfo) Size() int64                { return int64(len(i.f.Data)) }
func (i *openedFileInfo) Mode() fs.FileMode          { return i.f.Mode }
func (i *openedFileInfo) Type() fs.FileMode          { return i.f.Mode.Type() }
func (i *openedFileInfo) ModTime() time.Time         { return i.f.ModTime }
func (i *openedFileInfo) IsDir() bool                { return i.f.Mode&fs.ModeDir != 0 }
func (i *openedFileInfo) Sys() interface{}           { return i.f.Sys }
func (i *openedFileInfo) Info() (fs.FileInfo, error) { return i, nil }

// An openMapFile is a regular (non-directory) fs.File open for reading.
type openedFile struct {
	path string
	openedFileInfo
	offset int64
}

func (f *openedFile) Stat() (fs.FileInfo, error) { return &f.openedFileInfo, nil }

func (f *openedFile) Close() error { return nil }

func (f *openedFile) Read(b []byte) (int, error) {
	if f.offset >= int64(len(f.f.Data)) {
		return 0, io.EOF
	}
	if f.offset < 0 {
		return 0, &fs.PathError{Op: "read", Path: f.path, Err: fs.ErrInvalid}
	}
	n := copy(b, f.f.Data[f.offset:])
	f.offset += int64(n)
	return n, nil
}

func (f *openedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		// offset += 0
	case 1:
		offset += f.offset
	case 2:
		offset += int64(len(f.f.Data))
	}
	if offset < 0 || offset > int64(len(f.f.Data)) {
		return 0, &fs.PathError{Op: "seek", Path: f.path, Err: fs.ErrInvalid}
	}
	f.offset = offset
	return offset, nil
}

func (f *openedFile) ReadAt(b []byte, offset int64) (int, error) {
	if offset < 0 || offset > int64(len(f.f.Data)) {
		return 0, &fs.PathError{Op: "read", Path: f.path, Err: fs.ErrInvalid}
	}
	n := copy(b, f.f.Data[offset:])
	if n < len(b) {
		return n, io.EOF
	}
	return n, nil
}
