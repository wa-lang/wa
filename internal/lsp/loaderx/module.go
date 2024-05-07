// 版权 @2024 凹语言 作者。保留所有权利。

package loaderx

import "wa-lang.org/wa/internal/lsp/protocol"

// 文档的路径
type URI = protocol.DocumentURI

// 模块类型
type ModType int

const (
	ModType_Normal ModType = iota // 普通模块
	ModType_File                  // 单个文件, 映射到 __main__ 包
	ModType_Std                   // 标准库中的文件, 属于对应的包
)

// 模块
// 每个Wa代码只能属于一个模块
// 该模块只能依赖 Std, 不会再对其他代码产生依赖
type Module struct {
	ModType                     // 模块类型
	SrcRoot string              // 模块根目录
	PkgRoot string              // 对应的 import 根路径
	Files   map[URI]*File       // 工作区文件信息(可能包含被修改的 Std 文件)
	Pkgs    map[string]*Package // 全部包信息(包含 Std)
}

// 包信息
type Package struct {
	UriPath   URI      // 目录的路径
	PkgPath   string   // 属于的包路径
	FileNames []string // 文件名列表
	// 其他类型信息
}

// 文件
// 可包含当前未保存的编辑内容, 有版本号信息
// 文件发生变化时重建对应的包和被依赖的包
// 需要定期检查文件是否被外部修改
type File struct {
	Uri     string  // 资源路径
	Version uint32  // 文件编辑的版本号
	Content *string // 当前的最新内容, nil 表示未加载
	PkgPath string  // 属于的包路径
}
