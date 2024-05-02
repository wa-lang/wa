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

func (p *LSPServer) Hover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	p.logger.Println("Hover:", jsonMarshal(params))

	if !strings.HasSuffix(string(params.TextDocument.URI), ".wa") {
		return nil, nil
	}

	p.logger.Println("Hover.text:", p.fileMap[params.TextDocument.URI.Path()])

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
	var base int

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
			base = tokFile.Base()
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

	rng, err := mapper.OffsetRange(int(objIdent.Pos())-base, int(objIdent.End())-base)
	if err != nil {
		return nil, err
	}

	reply := &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  "wa",
			Value: obj.Type().String() + " @ " + module.program.Fset.Position(obj.Pos()).String(),
		},
		Range: rng,
	}

	return reply, nil
}
