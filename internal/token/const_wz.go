// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

// 关键字定义为常量, 避免直接引用字符串错误

const (
	K_点 = "·" // 全角 .
	K_注 = "注" // 注: 注释
)

const (
	K_引入 = "引入" // import

	K_常量 = "常量" // const
	K_全局 = "全局" // global
	K_类型 = "类型" // type
	K_函数 = "函数" // func

	K_结构 = "结构" // struct
	K_字典 = "字典" // map
	K_接口 = "接口" // interface

	K_设定 = "设定" // local

	K_如果 = "如果" // if
	K_或者 = "或者" // else if
	K_否则 = "否则" // else

	K_找辙 = "找辙" // switch
	K_有辙 = "有辙" // case
	K_没辙 = "没辙" // default

	K_循环 = "循环" // for
	K_迭代 = "迭代" // range
	K_继续 = "继续" // continue
	K_跳出 = "跳出" // break

	K_押后 = "押后" // defer
	K_返回 = "返回" // return

	K_区块 = "区块" // {
	K_完毕 = "完毕" // }
)

const (
	K_空 = "空" // nil
	K_真 = "真" // true
	K_假 = "假" // false

	K_约塔 = "约塔" // iota
	K_字节 = "字节" // byte
	K_符文 = "符文" // rune
	K_字串 = "字串" // string
	K_皮囊 = "皮囊" // any

	K_布尔 = "布尔" // bool
	K_整型 = "整型" // int
	K_正整 = "正整" // uint
	K_单精 = "单精" // float32
	K_双精 = "双精" // float64
	K_单复 = "单复" // complex64
	K_双复 = "双复" // complex128
	K_针型 = "针型" // uintptr

	K_微整型 = "微整型" // i8
	K_短整型 = "短整型" // i16
	K_普整型 = "普整型" // i32
	K_长整型 = "长整型" // i64
	K_微正整 = "微正整" // u8
	K_短正整 = "短正整" // u16
	K_普正整 = "普正整" // u32
	K_长正整 = "长正整" // u64

	K_错误 = "错误" // error
	K_报错 = "报错" // err.Error

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
	K_跟踪 = "跟踪" // trace

	K_主包 = "主包" // main
	K_太初 = "太初" // builtin/universe
	K_洪荒 = "洪荒" // unsafe
	K_丹田 = "丹田" // runtime

	// unsafe
	K_unsafe_原生    = "原生"    // unsafe.Raw
	K_unsafe_对齐倍数  = "对齐倍数"  // unsafe.Alignof
	K_unsafe_字节偏移量 = "字节偏移量" // unsafe.Offsetof
	K_unsafe_字节大小  = "字节大小"  // unsafe.Sizeof

	// runtime
	K_runtime_设置终结函数 = "设置终结函数" // runtime.SetFinalizer

	// 内置宏
	K__包__  = "__包__"
	K__文件__ = "__文件__"
	K__行号__ = "__行号__"
	K__列号__ = "__列号__"
	K__函数__ = "__函数__"
	K__位置__ = "__位置__"
)

// 待候选使用的名字
// "鸿蒙" => unsafe
// "照见" => reflect
// "三省" => reflect
// "气海" => runtime ?
// "檀中" => runtime ?
