// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file sets up the universe scope and the unsafe package.

package types

import "wa-lang.org/wa/internal/token"

// Typ contains the predeclared *Basic types indexed by their
// corresponding BasicKind.
//
// The *Basic type for Typ[Byte] will have the name "uint8".
// Use Universe.Lookup("byte").Type() to obtain the specific
// alias basic type named "byte" (and analogous for "rune").
var Typ = []*Basic{
	Invalid: {Invalid, 0, "invalid type"},

	Bool:          {Bool, IsBoolean, "bool"},
	Int:           {Int, IsInteger, "int"},
	Int8:          {Int8, IsInteger, "__wa_int8"},
	Int16:         {Int16, IsInteger, "__wa_int16"},
	Int32:         {Int32, IsInteger, "int32"},
	Int64:         {Int64, IsInteger, "int64"},
	Uint:          {Uint, IsInteger | IsUnsigned, "uint"},
	Uint8:         {Uint8, IsInteger | IsUnsigned, "uint8"},
	Uint16:        {Uint16, IsInteger | IsUnsigned, "uint16"},
	Uint32:        {Uint32, IsInteger | IsUnsigned, "uint32"},
	Uint64:        {Uint64, IsInteger | IsUnsigned, "uint64"},
	Uintptr:       {Uintptr, IsInteger | IsUnsigned, "uintptr"},
	Float32:       {Float32, IsFloat, "float32"},
	Float64:       {Float64, IsFloat, "float64"},
	Complex64:     {Complex64, IsComplex, "complex64"},
	Complex128:    {Complex128, IsComplex, "complex128"},
	String:        {String, IsString, "string"},
	UnsafePointer: {UnsafePointer, 0, token.K_unsafe_Pointer},

	UntypedBool:    {UntypedBool, IsBoolean | IsUntyped, "untyped bool"},
	UntypedInt:     {UntypedInt, IsInteger | IsUntyped, "untyped int"},
	UntypedRune:    {UntypedRune, IsInteger | IsUntyped, "untyped rune"},
	UntypedFloat:   {UntypedFloat, IsFloat | IsUntyped, "untyped float"},
	UntypedComplex: {UntypedComplex, IsComplex | IsUntyped, "untyped complex"},
	UntypedString:  {UntypedString, IsString | IsUntyped, "untyped string"},
	UntypedNil:     {UntypedNil, IsUntyped, "untyped nil"},
}

// A builtinId is the id of a builtin function.
type builtinId int

const (
	// universe scope
	_Append builtinId = iota
	_Cap
	_Complex
	_Copy
	_Delete
	_Imag
	_Len
	_Make
	_New
	_Panic
	_Print
	_Println
	_Real

	// package unsafe
	_unsafe_Alignof
	_unsafe_Offsetof
	_unsafe_Sizeof
	_unsafe_Raw
	_unsafe_SliceData

	// package runtime
	_runtime_SetFinalizer

	// testing support
	_Assert
	_Trace
)

// 之前的全局对象变成函数
// 根据当前包的语言模式选择对应的版本

func (p *Checker) XUniverse() *Scope {
	if p.pkg.W2Mode {
		return WzUniverse
	} else {
		return WaUniverse
	}
}
func (p *Checker) XUnsafe() *Package {
	if p.pkg.W2Mode {
		return WzUnsafe
	} else {
		return WaUnsafe
	}
}

func (p *Checker) _universeIota() *Const {
	if p.pkg.W2Mode {
		return wzUniverseIota
	} else {
		return waUniverseIota
	}
}
func (p *Checker) _universeByte() *Basic {
	if p.pkg.W2Mode {
		return wzUniverseByte
	} else {
		return waUniverseByte
	}
}
func (p *Checker) _universeRune() *Basic {
	if p.pkg.W2Mode {
		return wzUniverseRune
	} else {
		return waUniverseRune
	}
}

func (p *Checker) _universe__PACKAGE__() *Const {
	if p.pkg.W2Mode {
		return wzUniverse__PACKAGE__
	} else {
		return waUniverse__PACKAGE__
	}
}
func (p *Checker) _universe__FILE__() *Const {
	if p.pkg.W2Mode {
		return wzUniverse__FILE__
	} else {
		return waUniverse__FILE__
	}
}
func (p *Checker) _universe__LINE__() *Const {
	if p.pkg.W2Mode {
		return wzUniverse__LINE__
	} else {
		return waUniverse__LINE__
	}
}
func (p *Checker) _universe__COLUMN__() *Const {
	if p.pkg.W2Mode {
		return wzUniverse__COLUMN__
	} else {
		return waUniverse__COLUMN__
	}
}
func (p *Checker) _universe__FUNC__() *Const {
	if p.pkg.W2Mode {
		return wzUniverse__FUNC__
	} else {
		return waUniverse__FUNC__
	}
}
func (p *Checker) _universe__POS__() *Const {
	if p.pkg.W2Mode {
		return wzUniverse__POS__
	} else {
		return waUniverse__POS__
	}
}
