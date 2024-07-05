// 版权 @2023 凹语言 作者。保留所有权利。

package wat

// Module 表示一个 WASM 模块。
type Module struct {
	Doc  string // 注释
	Name string // 模块的名字(可空)

	InitFn string // 初始化函数
	MainFn string // 开始函数名字(可空)

	Memory *Memory // 内存对象, 包含 Data 信息(可空)
	Table  *Table  // 表格对象, 包含 Elem 信息(可空)

	Globals []Global // 全局对象
	Funcs   []Func   // 函数对象
}

// 内存信息
type Memory struct {
	Doc        string   // 注释
	Name       string   // 内存对象的名字
	ImportPath []string // 导入路径(仅针对导入对象)
	ExportName string   // 导出名字

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
	Doc        string   // 注释
	Name       string   // Table对象的名字
	ImportPath []string // 导入路径(仅针对导入对象)
	ExportName string   // 导出名字

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
	Doc        string   // 注释
	Name       string   // 全局变量名
	ImportPath []string // 导入路径(仅针对导入对象)
	ExportName string   // 导出名字

	Type    string // 类型信息
	Value   string // 初始值(导入变量忽略)
	Mutable bool   // 是否可写

}

// 函数定义
type Func struct {
	Doc        string   // 注释
	Name       string   // 函数名
	ImportPath []string // 导入路径(仅针对导入对象)
	ExportName string   // 导出名字

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

// 指令对应接口
type Instruction interface {
	DocString() string
	WATString() string // 支持缩进

	isInstruction()
}
