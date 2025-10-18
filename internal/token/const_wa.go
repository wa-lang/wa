// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

// 关键字定义为常量, 避免直接引用字符串错误

const (
	K_dot = "."
)

const (
	K_break    = "break"
	K_case     = "case"
	K_const    = "const"
	K_continue = "continue"

	K_default = "default"
	K_defer   = "defer"
	K_else    = "else"
	K_for     = "for"

	K_func   = "func"
	K_global = "global"
	K_if     = "if"
	K_import = "import"

	K_interface = "interface"
	K_map       = "map"
	K_range     = "range"
	K_return    = "return"

	K_struct = "struct"
	K_switch = "switch"
	K_type   = "type"
)

const (
	K_nil   = "nil"
	K_true  = "true"
	K_false = "false"

	K_iota   = "iota"
	K_byte   = "byte"
	K_rune   = "rune"
	K_string = "string"
	K_any    = "any"

	K_bool       = "bool"       // bool
	K_int        = "int"        // int
	K_uint       = "uint"       // uint
	K_f32        = "f32"        // float32
	K_f64        = "f64"        // float64
	K_complex64  = "complex64"  // complex64
	K_complex128 = "complex128" // complex128
	K_uintptr    = "uintptr"    // uintptr

	K_i8  = "__wa_i8"  // i8
	K_i16 = "__wa_i16" // i16
	K_i32 = "i32"      // i32
	K_i64 = "i64"      // i64
	K_u8  = "u8"       // u8
	K_u16 = "u16"      // u16
	K_u32 = "u32"      // u32
	K_u64 = "u64"      // u64

	K_error = "error"
	K_Error = "Error"

	K_init = "init"
	K_main = "main"
	K_this = "this"

	K_append  = "append"  // append
	K_cap     = "cap"     // cap
	K_complex = "complex" // complex
	K_copy    = "copy"    // copy
	K_delete  = "delete"  // delete
	K_imag    = "imag"    // imag
	K_len     = "len"     // len
	K_make    = "make"    // make
	K_new     = "new"     // new
	K_panic   = "panic"   // panic
	K_println = "println" // println
	K_print   = "print"   // print
	K_real    = "real"    // real

	K_assert = "assert" // assert
	K_trace  = "trace"  // trace

	K_universe = "universe"
	K_unsafe   = "unsafe"
	K_runtime  = "runtime"

	// unsafe
	K_unsafe_Raw      = "Raw"
	K_unsafe_Alignof  = "Alignof"
	K_unsafe_Offsetof = "Offsetof"
	K_unsafe_Sizeof   = "Sizeof"

	// runtime
	K_runtime_SetFinalizer = "SetFinalizer"

	// 内置宏
	K__PACKAGE__ = "__PACKAGE__"
	K__FILE__    = "__FILE__"
	K__LINE__    = "__LINE__"
	K__COLUMN__  = "__COLUMN__"
	K__FUNC__    = "__FUNC__"
	K__POS__     = "__POS__"
)
