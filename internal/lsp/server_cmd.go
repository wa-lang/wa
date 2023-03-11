// 版权 @2023 凹语言 作者。保留所有权利。

package lsp

import (
	"fmt"
)

func (p *LSPServer) executeCommand(params *ExecuteCommandParams) (interface{}, error) {
	if len(params.Arguments) != 1 {
		return nil, fmt.Errorf("invalid command")
	}
	return nil, nil
}

func (p *LSPServer) codeAction(uri string, params *CodeActionParams) ([]Command, error) {
	return nil, nil
}

func (p *LSPServer) completion(uri string, params *CompletionParams) ([]CompletionItem, error) {
	return nil, fmt.Errorf("TODO")
}

func (p *LSPServer) definition(uri string, params *TextDocumentPositionParams) ([]Location, error) {
	return nil, fmt.Errorf("TODO")
}

func (p *LSPServer) formatting(uri string, options FormattingOptions) ([]TextEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (p *LSPServer) hover(uri string, params *HoverParams) (*Hover, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *LSPServer) findRefs(uri string, params *TextDocumentPositionParams) ([]Location, error) {
	return nil, fmt.Errorf("TODO")
}

func (p *LSPServer) symbol(uri string) ([]DocumentSymbol, error) {
	return nil, fmt.Errorf("TODO")
}
