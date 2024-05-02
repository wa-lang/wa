// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/lsp/protocol"
	"wa-lang.org/wa/internal/token"
)

func (p *LSPServer) Definition(ctx context.Context, params *protocol.DefinitionParams) ([]protocol.Location, error) {
	p.logger.Println("Definition:", jsonMarshal(params))

	if !strings.HasSuffix(string(params.TextDocument.URI), ".wa") {
		return nil, nil
	}

	module, ok := p.loadModule(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}
	if err := module.LoadProgram(); err != nil {
		return nil, err
	}

	mainPkg := module.program.Pkgs[module.manifest.Pkg.Pkgpath] // todo(chai): fix
	p.logger.Println("manifest:", jsonMarshalIndent(module.manifest))

	if mainPkg == nil {
		return nil, fmt.Errorf("MainPkg is nil: %v", module.manifest.MainPkg)
	}

	baseName := filepath.Base(string(params.TextDocument.URI))

	path := params.TextDocument.URI.Path()
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	mapper := protocol.NewMapper(params.TextDocument.URI, content)

	var pos = token.NoPos
	var file *ast.File
	for _, f := range mainPkg.Files {
		// todo(chai): fix f.Pos() 位置
		p.logger.Println("file:", module.program.Fset.Position(f.Pos()), f.Name)
		tokFile := module.program.Fset.File(f.Pos())
		if tokFile.Name() == baseName {
			off, err := mapper.PositionOffset(params.Position)
			if err != nil {
				return nil, err
			}
			pos = token.Pos(tokFile.Base() + off)
			file = f
			break
		}
	}

	if pos == token.NoPos {
		return nil, fmt.Errorf("no pos")
	}

	var objIdent *ast.Ident
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		if pos < n.Pos() || pos > n.End() {
			return false
		}
		if x, ok := n.(*ast.Ident); ok {
			objIdent = x
			return false
		}
		return true
	})
	if objIdent == nil {
		return nil, nil
	}

	obj := mainPkg.Info.ObjectOf(objIdent)
	if obj == nil {
		return nil, nil
	}

	xpos := module.program.Fset.Position(obj.Pos())

	reply := []protocol.Location{
		{
			URI: params.TextDocument.URI,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(xpos.Line - 1),
					Character: uint32(xpos.Column - 1),
				},
				End: protocol.Position{
					Line:      uint32(xpos.Line - 1),
					Character: uint32(xpos.Column - 1),
				},
			},
		},
	}

	return reply, nil
}
