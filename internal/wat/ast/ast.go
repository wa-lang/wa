// 版权 @2024 凹语言 作者。保留所有权利。

package ast

import "wa-lang.org/wa/internal/wat/token"

// Notes: AST 刻意忽略了注释信息!

// Module 表示一个 WASM 模块。
type Module struct {
	Name    string         // 模块的名字(可空)
	Types   []*TypeSection // 类型定义
	Imports []*ImportSpec  // 导入对象
	Exports []*ExportSpec  // 导出对象
	Memory  *Memory        // 内存对象, 包含 Data 信息(可空)
	Data    []*DataSection // 初始化数据, 可能重叠
	Table   *Table         // 表格对象, 包含 Elem 信息(可空)
	Elem    []*ElemSection // 初始化数据
	Globals []*Global      // 全局对象
	Funcs   []*Func        // 函数对象
	Start   string         // start 函数
}

// 类型信息
type TypeSection struct {
	Name string
	Type *FuncType
}

// 导入对象(仅支持函数)
type ImportSpec struct {
	ObjModule string      // 模块路径
	ObjName   string      // 实体名字
	ObjKind   token.Token // 导入类别: MEMORY, TABLE, FUNC, GLOBAL

	Memory     *Memory     // 导入内存
	TableName  string      // 导入Table编号
	GlobalName string      // 导入全局变量名字
	GlobalType token.Token // 导入全局变量类型: I32, I64, F32, F64
	FuncName   string      // 导入后函数名字
	FuncType   *FuncType   // 导入函数类型
}

// 导出对象
type ExportSpec struct {
	Name      string      // 导出名字
	Kind      token.Token // 导出类型
	FuncIdx   string      // 导出函数ID
	TableIdx  string      // 导出表格ID
	MemoryIdx string      // 导出内存ID
	GlobalIdx string      // 导出全局变量ID
}

// 内存信息
type Memory struct {
	Name     string // 内存对象的名字
	Pages    int    // 页数, 每页 64 KB
	MaxPages int    // 最大页数
}

// 初始化数据
type DataSection struct {
	Name   string // 名字
	Offset uint32 // 偏移量(0 表示自动计算)
	Value  []byte // 初始化值
}

// Table 信息
type Table struct {
	Name    string      // Table对象的名字
	Size    int         // 表格容量
	MaxSize int         // 最大容量
	Type    token.Token // 元素类型, 默认为 funcref, externref
}

// 表格元素数据
type ElemSection struct {
	Name   string   // 名字
	Offset uint32   // 偏移量(从 0 开始)
	Values []string // 初始化值, 引用的是其他 func 的名字
}

// 全局变量
type Global struct {
	Name       string      // 全局变量名
	ExportName string      // 导出名字
	Mutable    bool        // 是否可写
	Type       token.Token // 类型信息
	I32Value   int32       // 初始值(导入变量忽略)
	I64Value   int64       // 初始值(导入变量忽略)
	F32Value   float32     // 初始值(导入变量忽略)
	F64Value   float64     // 初始值(导入变量忽略)
}

// 函数定义
type Func struct {
	Name       string    // 函数名
	ExportName string    // 导出名字
	Type       *FuncType // 函数类型
	Body       *FuncBody // 函数体
}

// 函数类型
type FuncType struct {
	Params  []Field       // 参数名字
	Results []token.Token // 返回值类型: I32, I64, F32, F64
}

// 函数定义
type FuncBody struct {
	Locals []Field       // 局部变量
	Insts  []Instruction // 指令列表
}

// 参数和局部变量信息
type Field struct {
	Name string      // 变量名字
	Type token.Token // 变量类型: I32, I64, F32, F64
}
