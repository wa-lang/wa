package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/scanner"
	"wa-lang.org/wa/internal/wat/token"
)

// The parser structure holds the parser's internal state.
type parser struct {
	m       *ast.Module
	errors  scanner.ErrorList
	scanner scanner.Scanner

	// Tracing/debugging
	trace  bool // == (mode & Trace != 0)
	indent int  // indentation used for tracing output

	// Comments
	//comments    []*ast.CommentGroup
	//leadComment *ast.CommentGroup // last lead comment
	//lineComment *ast.CommentGroup // last line comment

	// Next token
	pos int         // token position
	tok token.Token // one token look-ahead
	lit string      // token literal

	//unresolved []*ast.Ident      // unresolved identifiers
	//imports    []*ast.ImportSpec // list of imports
}

func (p *parser) init(path string, src []byte) {
	p.m = &ast.Module{
		File: token.NewFile(path, len(src)),
	}

	eh := func(pos token.Position, msg string) { p.errors.Add(pos, msg) }
	p.scanner.Init(p.m.File, src, eh, scanner.ScanComments)

	p.next()
}

func (p *parser) next() {
}

func (p *parser) parseFile() *ast.Module {
	return &ast.Module{}
}
