// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ld

import "encoding/binary"

var Linkarm = LinkArch{
	ByteOrder: binary.LittleEndian,
	Name:      "arm",
	Thechar:   '5',
	Minlc:     4,
	Ptrsize:   4,
	Regsize:   4,
}

var Linkarm64 = LinkArch{
	ByteOrder: binary.LittleEndian,
	Name:      "arm64",
	Thechar:   '7',
	Minlc:     4,
	Ptrsize:   8,
	Regsize:   8,
}

var Linkamd64 = LinkArch{
	ByteOrder: binary.LittleEndian,
	Name:      "amd64",
	Thechar:   '6',
	Minlc:     1,
	Ptrsize:   8,
	Regsize:   8,
}

var Link386 = LinkArch{
	ByteOrder: binary.LittleEndian,
	Name:      "386",
	Thechar:   '8',
	Minlc:     1,
	Ptrsize:   4,
	Regsize:   4,
}
