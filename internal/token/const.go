// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

// 关键字定义为常量, 避免直接引用字符串错误

const (
	K_nil    = "nil"
	K_true   = "true"
	K_false  = "false"
	K_iota   = "iota"
	K_byte   = "byte"
	K_rune   = "rune"
	K_string = "string"

	K_init = "init"
	K_main = "main"
	K_this = "this"

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
	K_点 = "·" // 全角 .
	K_注 = "注" // 注: 注释

	K_引入 = "引入"

	K_常量 = "常量"
	K_全局 = "全局"
	K_类型 = "类型"
	K_函数 = "函数"

	K_结构 = "结构"
	K_字典 = "字典"
	K_接口 = "接口"

	K_设定 = "设定" // 局部变量

	K_如果 = "如果"
	K_或者 = "或者"
	K_否则 = "否则"

	K_找辙 = "找辙"
	K_有辙 = "有辙" // K_或者
	K_没辙 = "没辙" // K_否则

	K_循环 = "循环"
	K_迭代 = "迭代"
	K_继续 = "继续"
	K_跳出 = "跳出"

	K_押后 = "押后"
	K_返回 = "返回"

	K_区块 = "区块"
	K_完毕 = "完毕"
)

const (
	K_空 = "空" // nil
	K_真 = "真" // true
	K_假 = "假" // false

	K_布尔 = "布尔" // bool
	K_整型 = "整型" // int
	K_正整 = "正整" // uint
	K_单精 = "单精" // float32
	K_双精 = "双精" // float64
	K_单复 = "单复" // complex64
	K_双复 = "双复" // complex128
	K_针型 = "针型" // uintptr
	K_字串 = "字串" // string
	K_字节 = "字节" // byte
	K_符文 = "符文" // rune

	K_微整型 = "微整型" // i8
	K_短整型 = "短整型" // i16
	K_普整型 = "普整型" // i32
	K_长整型 = "长整型" // i64
	K_微正整 = "微正整" // u8
	K_短正整 = "短正整" // u16
	K_普正整 = "普正整" // u32
	K_长正整 = "长正整" // u64

	K_错误 = "错误" // error

	K_约塔 = "约塔" // iota
	K_准备 = "准备" // init
	K_主控 = "主控" // main
	K_我的 = "我的" // this

	K_追加 = "追加" // append
	K_容量 = "容量" // cap
	K_复数 = "复数" // complex
	K_拷贝 = "拷贝" // copy
	K_删除 = "删除" // delete
	K_虚部 = "虚部" // imag
	K_长度 = "长度" // len
	K_构建 = "构建" // make
	K_新建 = "新建" // new
	K_崩溃 = "崩溃" // panic
	K_输出 = "输出" // println
	K_打印 = "打印" // print
	K_实部 = "实部" // real

	K_断言 = "断言" // assert

	K_主包 = "主包" // main
	K_太初 = "太初" // builtin
	K_鸿蒙 = "鸿蒙" // unsafe
	K_周行 = "周行" // runtime
)
