// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package ast

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/token"
)

// 汇编源文件
// 每个文件只是一个代码片段, 不能识别外部的符号类型, 只针对指令做简单的语义检查
// 只在链接阶段处理外部的符号依赖, 并做符号地址检查
type File struct {
	Pos      token.Pos  // 位置
	Consts   []*Const   // 全局常量
	Globals  []*Global  // 全局对象
	Funcs    []*Func    // 函数对象
	Start    string     // start 函数
	Comments []*Comment // 注释列表, 可根据 pos 定位
}

// 注释
type Comment struct {
	Pos  token.Pos // 位置
	Text string    // 注释文本
}

// 全局常量
type Const struct {
	Pos   token.Pos // 位置
	Name  string    // 常量名
	Value LitValue  // 常量面值
}

// 全局对象
type Global struct {
	Pos  token.Pos   // 位置
	Name string      // 全局变量名
	Type token.Token // 类型, 缺少时用 .Size 表示, 语法糖
	Size int         // 大小, 对应结构体
	Init []InitValue // 初始数据
}

// 初始化的面值
type InitValue struct {
	Offset     int         // 相对偏移
	Type       token.Token // 目标类型(可以i64转型为i32)
	LitValue   *LitValue   // 常量面值
	ConstValue *Const      // 全局常量
	GlobalAddr *Global     // 全局变量地址
}

// 常量面值
type LitValue struct {
	Type       token.Token // 初始化值的类型
	IntValue   int64       // I32/I64 值
	FloatValue float64     // F32/F64 值
	StrValue   *string     // 字符串, 指针非空时有效
}

// 函数对象
type Func struct {
	Pos  token.Pos // 位置
	Name string    // 函数名
	Type *FuncType // 函数类型
	Body *FuncBody // 函数体
}

// 函数类型
type FuncType struct {
	Args   []Field     // 参数列表
	Return token.Token // 返回值类型
}

// 函数定义
type FuncBody struct {
	Locals []Field       // 局部变量
	Insts  []Instruction // 指令列表
}

// 参数/局部变量
type Field struct {
	Pos  token.Pos   // 位置
	Type token.Token // 类型
	Name string      // 名字
}

// 机器指令
type Instruction struct {
	Pos   token.Pos       // 位置
	Label string          // 指令对应的 Label
	As    abi.As          // 汇编指令
	Arg   *abi.AsArgument // 指令参数
}
