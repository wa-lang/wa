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
	Pos     token.Pos // 位置
	Name    string    // 模块名
	Globals []*Global // 全局对象
	Funcs   []*Func   // 函数对象
	Start   string    // start 函数
}

// 全局对象
type Global struct {
	Pos  token.Pos // 位置
	Name string    // 全局变量名
	Addr int64     // 内存地址
	Size int       // 大小

	InitAddr int64  // 初始数据开始地址
	InitData []byte // 初始数据
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
	Pos token.Pos       // 位置
	As  abi.As          // 汇编指令
	Arg *abi.AsArgument // 指令参数
}
