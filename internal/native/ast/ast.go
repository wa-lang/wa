// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package ast

import "wa-lang.org/wa/internal/native/token"

// 汇编源文件
// 每个文件只是一个代码片段, 不能识别外部的符号类型, 只针对指令做简单的语义检查
// 只在链接阶段处理外部的符号依赖, 并做符号地址检查
type File struct {
	Pos     token.Pos // 位置
	Name    string    // 模块名
	Globals []*Global // 全局对象
	Funcs   []*Func   // 函数对象
}

// 全局对象
type Global struct {
	Pos  token.Pos // 位置
	Name string    // 全局变量名
	// TODO: 长度/地址对其/初始化值
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
	ArgsName  []string // 参数名字
	ArgsSize  []int    // 参数大小
	FrameSize int      // 函数帧大小
	RetSize   int      // 返回值大小
}

// 函数定义
type FuncBody struct {
	Insts []Instruction // 指令列表
}

// 机器指令
type Instruction struct {
	Pos token.Pos // 位置
	As  int       // 指令的值(平台相关)
	// TODO: 指令的类型和寻址方式
}

// 地址对象
type Addr struct{}
