// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w18n

// 本地化配置
type Local struct {
	Local string // 本地化的名字

	// 编译指令(Directive)
	D_wx                string
	D_wx_build          string
	D_wx_linkname       string
	D_wx_export         string
	D_wx_force_register string
	D_wx_runtime_getter string
	D_wx_runtime_setter string
	D_wx_runtime_sizer  string
	D_wx_generic        string
	D_wx_operator       string
	D_wx_embed          string
	D_wx_ignore         string

	// 关键字(Keyword)
	K_引入 string
	K_常量 string
	K_全局 string
	K_类型 string
	K_函数 string
	K_结构 string
	K_字典 string
	K_接口 string
	K_设定 string
	K_如果 string
	K_或者 string
	K_否则 string
	K_找辙 string
	K_有辙 string
	K_没辙 string
	K_循环 string
	K_迭代 string
	K_继续 string
	K_跳出 string
	K_押后 string
	K_返回 string
	K_区块 string
	K_完毕 string

	// 预定义: 内置宏
	Builtin__PACKAGE__ string
	Builtin__FILE__    string
	Builtin__LINE__    string
	Builtin__COLUMN__  string
	Builtin__FUNC__    string
	Builtin__POS__     string

	// 预定义: 常量
	Builtin_nil   string
	Builtin_true  string
	Builtin_false string
	Builtin_iota  string

	// 预定义: 类型
	Builtin_byte   string
	Builtin_rune   string
	Builtin_string string
	Builtin_any    string

	// 预定义: 类型
	Builtin_bool       string
	Builtin_int        string
	Builtin_uint       string
	Builtin_f32        string
	Builtin_f64        string
	Builtin_complex64  string
	Builtin_complex128 string
	Builtin_uintptr    string

	// 预定义: 类型
	Builtin___wa_i8  string
	Builtin___wa_i16 string
	Builtin_i32      string
	Builtin_i64      string
	Builtin_u8       string
	Builtin_u16      string
	Builtin_u32      string
	Builtin_u64      string

	// 预定义: 类型
	Builtin_error string
	Builtin_Error string

	// 预定义: 特殊函数名
	Builtin_init string
	Builtin_main string
	Builtin_this string

	// 预定义: 内置函数
	Builtin_append  string
	Builtin_cap     string
	Builtin_complex string
	Builtin_copy    string
	Builtin_delete  string
	Builtin_imag    string
	Builtin_len     string
	Builtin_make    string
	Builtin_new     string
	Builtin_panic   string
	Builtin_println string
	Builtin_print   string
	Builtin_real    string
	Builtin_assert  string
	Builtin_trace   string

	// unsafe 包
	Unsafe_Pointer    string
	Unsafe_Raw        string
	Unsafe_Alignof    string
	Unsafe_Offsetof   string
	Unsafe_Sizeof     string
	Unsafe_SliceData  string
	Unsafe_StringData string

	// 运行时包
	Runtime_SetFinalizer string
}
