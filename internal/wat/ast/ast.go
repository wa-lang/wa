// 版权 @2024 凹语言 作者。保留所有权利。

package ast

import "wa-lang.org/wa/internal/wat/token"

// Module 表示一个 WASM 模块。
type Module struct {
	File *token.File // 文件位置信息

	Doc  string // 关联注释
	Name string // 模块的名字(可空)

	Imports []*ImportSpec // 导入对象
	Memory  *Memory       // 内存对象, 包含 Data 信息(可空)
	Table   *Table        // 表格对象, 包含 Elem 信息(可空)
	Globals []Global      // 全局对象
	Funcs   []Func        // 函数对象
}

// 节点信息
type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
}

// 指令对应接口
type Instruction interface {
	Node
	isInstruction()
}

// 导入对象(仅支持函数)
type ImportSpec struct {
	Doc        string    // 关联的注释
	ModulePath string    // 模块路径
	FuncPath   string    // 函数路径
	FuncName   string    // 导入后的名字
	FuncType   *FuncType // 函数类型
	EndPos     token.Pos // 结束位置
}

// 内存信息
type Memory struct {
	Doc        string // 注释
	Name       string // 内存对象的名字
	ExportName string // 导出名字

	Pages int           // 页数, 每页 64 KB
	Data  []DataSection // 初始化数据, 可能重叠
}

// 初始化数据
type DataSection struct {
	Doc    string // 注释
	Name   string // 名字
	Offset uint32 // 偏移量(0 表示自动计算)
	Value  []byte // 初始化值
}

// Table 信息
type Table struct {
	Doc        string // 注释
	Name       string // Table对象的名字
	ExportName string // 导出名字

	Size int           // 表格容量
	Type string        // 元素类型, 默认为 anyfunc
	Elem []ElemSection // 初始化数据
}

// 表格元素数据
type ElemSection struct {
	Doc    string // 注释
	Offset uint32 // 偏移量(从 0 开始)
	Value  string // 初始化值, 引用的是其他 func 的名字
}

// 全局变量
type Global struct {
	Doc        string // 注释
	Name       string // 全局变量名
	ExportName string // 导出名字

	Type    string // 类型信息
	Value   string // 初始值(导入变量忽略)
	Mutable bool   // 是否可写
}

// 函数定义
type Func struct {
	Doc        string // 注释
	Name       string // 函数名
	ExportName string // 导出名字
	IsStart    bool   // start 函数

	Type *FuncType // 函数类型
	Body *FuncBody // 函数体
}

// 函数类型
type FuncType struct {
	Params      []Field  // 参数名字
	ResultsType []string // 返回值类型
}

// 函数定义
type FuncBody struct {
	Locals []Field       // 局部变量
	Insts  []Instruction // 指令列表
}

// 参数和局部变量信息
type Field struct {
	Doc  string // 注释
	Name string // 变量名字
	Type string // 变量类型
}
