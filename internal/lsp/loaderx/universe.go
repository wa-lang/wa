// 版权 @2024 凹语言 作者。保留所有权利。

package loaderx

import (
	"io/fs"
	"log"
	"time"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/lsp/protocol"
	"wa-lang.org/wa/internal/token"
	"wa-lang.org/wa/internal/types"
)

type (
	URI         = protocol.URI
	DocumentURI = protocol.DocumentURI
	PkgpathURI  string // import "pkgpath"
)

type Config struct {
	WaOS   string
	WaRoot string

	logger *log.Logger
}

// 模块
// 每个Wa代码只能属于一个模块
// 该模块只能依赖 Std, 不会再对其他代码产生依赖
type Universe struct {
	Version         int32  // 版本号
	WaOS            string // 系统
	WaRoot          string // 模块根目录
	WorkspaceFolder []URI  // 工作区目录

	Fset  *token.FileSet          //
	Files map[DocumentURI]*File   // 工作区文件信息(可能包含被修改的 Std 文件)
	Pkgs  map[PkgpathURI]*Package // 全部包信息(包含 Std)

	logger *log.Logger
}

// 包信息
type Package struct {
	Version   int32    // 版本号
	UriPath   URI      // 目录的路径
	PkgPath   string   // 属于的包路径
	FileNames []string // 文件名列表

	Pkg   *types.Package // 类型检查后的包
	Info  *types.Info    // 具名对象信息
	Files []*ast.File    // 语法树
}

// 文件
// 可包含当前未保存的编辑内容, 有版本号信息
// 文件发生变化时重建对应的包和被依赖的包
// 需要定期检查文件是否被外部修改
type File struct {
	Version int32       // 版本号
	FileUri DocumentURI // 资源路径
	PkgPath string      // 属于的包路径

	Data    []byte      // file content
	Mode    fs.FileMode // FileInfo.Mode
	ModTime time.Time   // FileInfo.ModTime
	Sys     interface{} // FileInfo.Sys

	Parsed       bool
	ParsedResult *ast.Node
	ParsedErr    error

	Saved bool
}
