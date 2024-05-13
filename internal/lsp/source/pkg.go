// 版权 @2024 凹语言 作者。保留所有权利。

package source

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/lsp/span"
	"wa-lang.org/wa/internal/types"
)

// Package represents a Go package that has been type-checked. It maintains
// only the relevant fields of a *go/packages.Package.
type Package interface {
	ID() string
	Name() string
	PkgPath() string
	CompiledWaFiles() []*ParsedWaFile
	File(uri span.URI) (*ParsedWaFile, error)
	GetSyntax() []*ast.File
	GetTypes() *types.Package
	GetTypesInfo() *types.Info
	GetTypesSizes() types.Sizes
	IsIllTyped() bool
	ForTest() string
	GetImport(pkgPath string) (Package, error)
	MissingDependencies() []string
	Imports() []Package
	HasListOrParseErrors() bool
	HasTypeErrors() bool
}
