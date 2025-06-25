// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import (
	"encoding/json"
)

type serializedFile struct {
	// fields correspond 1:1 to fields with same (lower-case) name in File
	Name  string
	Base  int
	Size  int
	Lines []int
	Infos []lineInfo
}

type serializedFileSet struct {
	Base  int
	Files []serializedFile
}

// 编码为JSON
func (s *FileSet) ToJson() []byte {
	var jsonBytes []byte
	var err error

	err = s.Write(func(data interface{}) error {
		jsonBytes, err = json.MarshalIndent(data, "", "\t")
		if err != nil {
			panic(err)
		}
		return nil
	})
	return jsonBytes
}

// 从JSON解码
func (s *FileSet) FromJson(jsonBytes []byte) error {
	return s.Read(func(x interface{}) error {
		return json.Unmarshal(jsonBytes, x)
	})
}

// Read calls decode to deserialize a file set into s; s must not be nil.
func (s *FileSet) Read(decode func(interface{}) error) error {
	var ss serializedFileSet
	if err := decode(&ss); err != nil {
		return err
	}

	s.mutex.Lock()
	s.base = ss.Base
	files := make([]*File, len(ss.Files))
	for i := 0; i < len(ss.Files); i++ {
		f := &ss.Files[i]
		files[i] = &File{
			set:   s,
			name:  f.Name,
			base:  f.Base,
			size:  f.Size,
			lines: f.Lines,
			infos: f.Infos,
		}
	}
	s.files = files
	s.last = nil
	s.mutex.Unlock()

	return nil
}

// Write calls encode to serialize the file set s.
func (s *FileSet) Write(encode func(interface{}) error) error {
	var ss serializedFileSet

	s.mutex.Lock()
	ss.Base = s.base
	files := make([]serializedFile, len(s.files))
	for i, f := range s.files {
		f.mutex.Lock()
		files[i] = serializedFile{
			Name:  f.name,
			Base:  f.base,
			Size:  f.size,
			Lines: append([]int(nil), f.lines...),
			Infos: append([]lineInfo(nil), f.infos...),
		}
		f.mutex.Unlock()
	}
	ss.Files = files
	s.mutex.Unlock()

	return encode(ss)
}
