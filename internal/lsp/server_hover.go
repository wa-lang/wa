// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"wa-lang.org/wa/internal/lsp/protocol"
	"wa-lang.org/wa/internal/token"
)

func (p *LSPServer) Hover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	p.logger.Println("Hover:", jsonMarshal(params))

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
	for _, f := range mainPkg.Files {
		// todo(chai): fix f.Pos() 位置
		p.logger.Println("file:", module.program.Fset.Position(f.Pos()), f.Name)
		if module.program.Fset.File(f.Pos()).Name() == baseName {
			off, err := mapper.PositionOffset(params.Position)
			if err != nil {
				return nil, err
			}
			pos = f.Pos() + token.Pos(off)
			break
		}
	}

	if pos == token.NoPos {
		return nil, fmt.Errorf("no pos")
	}

	scope := mainPkg.Pkg.Scope().Innermost(pos)
	for i, name := range scope.Names() {
		obj := scope.Lookup(name)
		p.logger.Println(i, name, obj.Pos())
	}

	obj := scope.LookupByPos(pos)
	if obj == nil {
		return nil, nil
	}

	rng, err := mapper.OffsetRange(int(obj.Pos()), int(obj.Node().End()))
	if err != nil {
		return nil, err
	}

	reply := &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  "wa",
			Value: obj.String(),
		},
		Range: rng,
	}

	return reply, nil
}
