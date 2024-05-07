// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/ast/astutil"
	"wa-lang.org/wa/internal/lsp/protocol"
	"wa-lang.org/wa/internal/token"
	"wa-lang.org/wa/internal/types"
)

func (p *LSPServer) Completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	p.logger.Println("Completion:", jsonMarshal(params))

	// todo: 需要获取到最新的还未保存的文件内容

	if !strings.HasSuffix(string(params.TextDocument.URI), ".wa") {
		return nil, nil
	}

	module, ok := p.loadModule(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}
	if err := module.LoadProgram(); err != nil {
		p.logger.Println("LoadProgram:", err)
		return nil, err
	}

	mainPkg := module.program.Pkgs[module.manifest.Pkg.Pkgpath] // todo(chai): fix
	p.logger.Println("manifest:", jsonMarshalIndent(module.manifest))

	if mainPkg == nil {
		return nil, fmt.Errorf("MainPkg is nil: %v", module.manifest.MainPkg)
	}

	baseName := filepath.Base(string(params.TextDocument.URI))

	content, err := os.ReadFile(params.TextDocument.URI.Path())
	if err != nil {
		p.logger.Println("ReadFile:", params.TextDocument.URI.Path(), err)
		return nil, err
	}

	if s, ok := p.fileMap[params.TextDocument.URI.Path()]; ok {
		content = []byte(s)
	}

	mapper := protocol.NewMapper(params.TextDocument.URI, content)

	var pos = token.NoPos
	var file *ast.File
	for _, f := range mainPkg.Files {
		// todo(chai): fix f.Pos() 位置
		p.logger.Println("file:", module.program.Fset.Position(f.Pos()), f.Name)
		tokFile := module.program.Fset.File(f.Pos())
		if tokFile.Name() == baseName {
			params_Position := params.Position
			params_Position.Character--
			off, err := mapper.PositionOffset(params_Position)
			if err != nil {
				p.logger.Println("mapper.PositionOffset(params.Position):", err)
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

	path, _ := astutil.PathEnclosingInterval(file, pos-1, pos-1)
	if path == nil {
		return nil, fmt.Errorf("cannot find node enclosing position")
	}

	objIdent, _ := path[0].(*ast.Ident)
	if objIdent == nil {
		return nil, nil
	}

	obj := mainPkg.Info.ObjectOf(objIdent)
	if obj == nil {
		return nil, nil
	}

	var completionItems []protocol.CompletionItem
	switch obj := obj.(type) {
	case *types.PkgName:
		for _, name := range obj.Imported().Scope().Names() {
			completionItems = append(completionItems, protocol.CompletionItem{
				Label:      name,
				InsertText: name,
			})
		}
	default:
		// todo
	}

	reply := &protocol.CompletionList{
		Items: completionItems,
	}
	return reply, nil
}
