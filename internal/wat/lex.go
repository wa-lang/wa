// 版权 @2023 凹语言 作者。保留所有权利。

package wat

import (
	"fmt"
)

func waLex(name, input string) (tokens, comments []Token) {
	l := newLexer(name, input)
	tokens = l.Tokens()
	comments = l.Comments()
	return
}

type watLexer struct {
	src      *watSourceStream
	tokens   []Token
	comments []Token
}

func newLexer(name, input string) *watLexer {
	return &watLexer{
		src: newSourceStream(name, input),
	}
}

func (p *watLexer) Tokens() []Token {
	if len(p.tokens) == 0 {
		p.run()
	}
	return p.tokens
}

func (p *watLexer) Comments() []Token {
	if len(p.tokens) == 0 {
		p.run()
	}
	return p.comments
}

func (p *watLexer) emit(kind TokenKind) {
	lit, pos := p.src.EmitToken()
	_ = pos

	if kind == EOF {
		pos = p.src.Pos()
		lit = ""
	}
	if kind == IDENT {
		kind = watLookup(lit)
	}
	p.tokens = append(p.tokens, Token{
		Kind:    kind,
		Literal: lit,
		//Pos:     Pos(pos + 1),
	})
}

func (p *watLexer) emitComment() {
	lit, pos := p.src.EmitToken()
	_ = pos

	p.comments = append(p.comments, Token{
		Kind:    COMMENT,
		Literal: lit,
		// Pos:     Pos(pos + 1),
	})
}

func (p *watLexer) errorf(format string, args ...interface{}) {
	tok := Token{
		//Kind:    ERROR,
		Literal: fmt.Sprintf(format, args...),
		//Pos:     Pos(p.src.Pos() + 1),
	}
	p.tokens = append(p.tokens, tok)
	panic(tok)
}

func (p *watLexer) run() (tokens []Token) {
	defer func() {
		tokens = p.tokens
		if r := recover(); r != nil {
			if _, ok := r.(Token); !ok {
				panic(r)
			}
		}
	}()

	for {
		r := p.src.Read()
		if r == rune(EOF) {
			p.emit(EOF)
			return
		}

		switch {

		// TODO: case xxx

		default:
			p.errorf("unrecognized character: %#U", r)
			return
		}
	}
}
