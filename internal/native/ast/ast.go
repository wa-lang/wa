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
	Pos         token.Pos               // 位置
	CPU         abi.CPUType             // CPU类型
	Doc         *CommentGroup           // 关联文档
	IntelSyntax *GasIntelSyntaxNoprefix // X64必须采用 Intel 语法
	Externs     []*GasExtern            // 外部的符号
	Consts      []*Const                // 全局常量
	Globals     []*Global               // 全局对象
	Funcs       []*Func                 // 函数对象
	Comments    []*CommentGroup         // 孤立的注释
	Objects     []Object                // 保序的列表
}

// 单行注释
type Comment struct {
	Pos  token.Pos // # 位置
	Text string    // 注释文本
}

// 一组相邻的注释
type CommentGroup struct {
	TopLevel bool       // 是否为全局类型的注释
	List     []*Comment // len(List) > 0
}

// X64 汇编特有的头部
type GasIntelSyntaxNoprefix struct {
	Pos token.Pos // 位置
}

// 外部符号
type GasExtern struct {
	Pos  token.Pos // 位置
	Name string    // 符号名字
}

// 基本的面值
type BasicLit struct {
	Pos       token.Pos   // 位置
	LitKind   token.Token // INT/FLOAT/CHAR/STRING, 默认类型 INT=>I64, FLOAT=>F64, CHAR => I32
	LitString string      // 原始的字符串: 42, 0x7f, 3.14, 1e-9, 'a', '\x7f', "foo" or `\m\n\o`
	ConstV    interface{} // 解析后的常量值, 对应的类型: int64, float64, string
}

// 全局常量
type Const struct {
	Pos     token.Pos     // 位置
	Tok     token.Token   // 关键字(可能有多语言)
	Doc     *CommentGroup // 关联文档
	Name    string        // 常量名
	Value   *BasicLit     // 常量面值
	Comment *Comment      // 尾部单行注释
}

// 全局对象
type Global struct {
	Pos        token.Pos         // 位置
	Tok        token.Token       // 关键字(可能有多语言)
	Doc        *CommentGroup     // 关联文档
	Name       string            // 全局变量名
	TypeTok    token.Token       // 原始的类型
	Type       token.Type        // 全局变量的类型(结合初始值类型推导)
	Size       int               // 内存大小(没有类型信息)
	Init       *InitValue        // 初始数据
	Objects    []Object          // 保序的对象
	LinkInfo   *abi.LinkedSymbol // 链接信息
	Section    string            // 段名
	Align      int               // 对齐
	ExportName string            // 导出名
}

// 初始化的面值
type InitValue struct {
	Pos     token.Pos // 位置
	Lit     *BasicLit // 或字面值
	Symbal  string    // 或标识符(只能是常量或地址)
	Comment *Comment  // 尾部单行注释
}

// 函数对象
type Func struct {
	Pos        token.Pos         // 位置
	Tok        token.Token       // 关键字(可能有多语言)
	Doc        *CommentGroup     // 关联文档
	Name       string            // 函数名
	Type       *FuncType         // 函数类型
	ArgsSize   int               // 调用该函数需要的栈大小(参数/返回值)
	FrameSize  int               // 函数内栈帧大小(局部变量/临时空间), 不包含头部(rip/rbp)
	BodySize   int               // 指令大小(没有类型信息)
	Body       *FuncBody         // 函数体
	LinkInfo   *abi.LinkedSymbol // 链接信息
	Section    string            // 段名
	ExportName string            // 导出的名字
}

// 函数类型
type FuncType struct {
	Pos    token.Pos // 位置
	Args   []*Local  // 参数列表
	Return []*Local  // 返回值类型
}

// 函数定义
type FuncBody struct {
	Pos      token.Pos       // 位置
	Locals   []*Local        // 局部变量
	Insts    []*Instruction  // 指令列表
	Comments []*CommentGroup // 孤立的注释
	Objects  []Object        // 保序对象列表(包含空行)
}

// 参数/返回值/局部变量
// 对应寄存器或栈帧的偏移量都是固定的
type Local struct {
	Pos     token.Pos     // 位置
	Tok     token.Token   // 关键字(可能有多语言)
	Doc     *CommentGroup // 关联文档
	Name    string        // 名字
	Type    token.Type    // 类型
	Comment *Comment      // 尾部单行注释
	Reg     abi.RegType   // 对应的寄存器, 在生成机器码阶段计算
	RBPOff  int           // 相对于Rbp的偏移地址, 在生成机器码阶段计算
	RSPOff  int           // 相对于Rsp的偏移地址, 用于调用方根据 rsp 定位
	Cap     int           // 容量, 元素个数, 默认为1
}

// 机器指令
// 函数体内的独立注释作为空指令记录
type Instruction struct {
	CPU     abi.CPUType     // CPU类型, 用于格式化
	Pos     token.Pos       // 位置
	Doc     *CommentGroup   // 关联文档
	Label   string          // 指令对应的 Label, 可以只是 Label
	As      abi.As          // 汇编指令
	AsName  string          // 指令的名字(支持多语言)
	Arg     *abi.AsArgument // 指令参数
	Comment *Comment        // 尾部单行注释
}

// 空行, 用于指令间分隔打印
type BlankLine struct {
	Pos token.Pos
}
