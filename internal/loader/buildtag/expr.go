// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildtag

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

// An Expr is a build tag constraint expression.
// The underlying concrete type is *AndExpr, *OrExpr, *NotExpr, or *TagExpr.
type Expr interface {
	// String returns the string form of the expression,
	// using the boolean syntax used in #wa:build lines.
	String() string

	// Eval reports whether the expression evaluates to true.
	// It calls ok(tag) as needed to find out whether a given build tag
	// is satisfied by the current build configuration.
	Eval(ok func(tag string) bool) bool

	// The presence of an isExpr method explicitly marks the type as an Expr.
	// Only implementations in this package should be used as Exprs.
	isExpr()
}

// A TagExpr is an Expr for the single tag Tag.
type TagExpr struct {
	Tag string // for example, “linux” or “cgo”
}

func (x *TagExpr) isExpr() {}

func (x *TagExpr) Eval(ok func(tag string) bool) bool {
	return ok(x.Tag)
}

func (x *TagExpr) String() string {
	return x.Tag
}

func tag(tag string) Expr { return &TagExpr{tag} }

// A NotExpr represents the expression !X (the negation of X).
type NotExpr struct {
	X Expr
}

func (x *NotExpr) isExpr() {}

func (x *NotExpr) Eval(ok func(tag string) bool) bool {
	return !x.X.Eval(ok)
}

func (x *NotExpr) String() string {
	s := x.X.String()
	switch x.X.(type) {
	case *AndExpr, *OrExpr:
		s = "(" + s + ")"
	}
	return "!" + s
}

func not(x Expr) Expr { return &NotExpr{x} }

// An AndExpr represents the expression X && Y.
type AndExpr struct {
	X, Y Expr
}

func (x *AndExpr) isExpr() {}

func (x *AndExpr) Eval(ok func(tag string) bool) bool {
	// Note: Eval both, to make sure ok func observes all tags.
	xok := x.X.Eval(ok)
	yok := x.Y.Eval(ok)
	return xok && yok
}

func (x *AndExpr) String() string {
	return andArg(x.X) + " && " + andArg(x.Y)
}

func andArg(x Expr) string {
	s := x.String()
	if _, ok := x.(*OrExpr); ok {
		s = "(" + s + ")"
	}
	return s
}

func and(x, y Expr) Expr {
	return &AndExpr{x, y}
}

// An OrExpr represents the expression X || Y.
type OrExpr struct {
	X, Y Expr
}

func (x *OrExpr) isExpr() {}

func (x *OrExpr) Eval(ok func(tag string) bool) bool {
	// Note: Eval both, to make sure ok func observes all tags.
	xok := x.X.Eval(ok)
	yok := x.Y.Eval(ok)
	return xok || yok
}

func (x *OrExpr) String() string {
	return orArg(x.X) + " || " + orArg(x.Y)
}

func orArg(x Expr) string {
	s := x.String()
	if _, ok := x.(*AndExpr); ok {
		s = "(" + s + ")"
	}
	return s
}

func or(x, y Expr) Expr {
	return &OrExpr{x, y}
}

// A SyntaxError reports a syntax error in a parsed build expression.
type SyntaxError struct {
	Offset int    // byte offset in input where error was detected
	Err    string // description of error
}

func (e *SyntaxError) Error() string {
	return e.Err
}

var errNotConstraint = errors.New("not a build constraint")

// Parse parses a single build constraint line of the form “#wa:build ...”
// and returns the corresponding boolean expression.
func Parse(line string) (Expr, error) {
	if text, ok := splitWaBuild(line); ok {
		return parseExpr(text)
	}
	return nil, errNotConstraint
}

// IsWaBuild reports whether the line of text is a “#wa:build” constraint.
// It only checks the prefix of the text, not that the expression itself parses.
func IsWaBuild(line string) bool {
	_, ok := splitWaBuild(line)
	return ok
}

// splitWaBuild splits apart the leading #wa:build prefix in line from the build expression itself.
// It returns "", false if the input is not a #wa:build line or if the input contains multiple lines.
func splitWaBuild(line string) (expr string, ok bool) {
	// A single trailing newline is OK; otherwise multiple lines are not.
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	if strings.Contains(line, "\n") {
		return "", false
	}

	if !strings.HasPrefix(line, "#wa:build") {
		return "", false
	}

	line = strings.TrimSpace(line)
	line = line[len("#wa:build"):]

	// If strings.TrimSpace finds more to trim after removing the #wa:build prefix,
	// it means that the prefix was followed by a space, making this a #wa:build line
	// (as opposed to a #wa:buildsomethingelse line).
	// If line is empty, we had "#wa:build" by itself, which also counts.
	trim := strings.TrimSpace(line)
	if len(line) == len(trim) && line != "" {
		return "", false
	}

	return trim, true
}

// An exprParser holds state for parsing a build expression.
type exprParser struct {
	s string // input string
	i int    // next read location in s

	tok   string // last token read
	isTag bool
	pos   int // position (start) of last token
}

// parseExpr parses a boolean build tag expression.
func parseExpr(text string) (x Expr, err error) {
	defer func() {
		if e := recover(); e != nil {
			if e, ok := e.(*SyntaxError); ok {
				err = e
				return
			}
			panic(e) // unreachable unless parser has a bug
		}
	}()

	p := &exprParser{s: text}
	x = p.or()
	if p.tok != "" {
		panic(&SyntaxError{Offset: p.pos, Err: "unexpected token " + p.tok})
	}
	return x, nil
}

// or parses a sequence of || expressions.
// On entry, the next input token has not yet been lexed.
// On exit, the next input token has been lexed and is in p.tok.
func (p *exprParser) or() Expr {
	x := p.and()
	for p.tok == "||" {
		x = or(x, p.and())
	}
	return x
}

// and parses a sequence of && expressions.
// On entry, the next input token has not yet been lexed.
// On exit, the next input token has been lexed and is in p.tok.
func (p *exprParser) and() Expr {
	x := p.not()
	for p.tok == "&&" {
		x = and(x, p.not())
	}
	return x
}

// not parses a ! expression.
// On entry, the next input token has not yet been lexed.
// On exit, the next input token has been lexed and is in p.tok.
func (p *exprParser) not() Expr {
	p.lex()
	if p.tok == "!" {
		p.lex()
		if p.tok == "!" {
			panic(&SyntaxError{Offset: p.pos, Err: "double negation not allowed"})
		}
		return not(p.atom())
	}
	return p.atom()
}

// atom parses a tag or a parenthesized expression.
// On entry, the next input token HAS been lexed.
// On exit, the next input token has been lexed and is in p.tok.
func (p *exprParser) atom() Expr {
	// first token already in p.tok
	if p.tok == "(" {
		pos := p.pos
		defer func() {
			if e := recover(); e != nil {
				if e, ok := e.(*SyntaxError); ok && e.Err == "unexpected end of expression" {
					e.Err = "missing close paren"
				}
				panic(e)
			}
		}()
		x := p.or()
		if p.tok != ")" {
			panic(&SyntaxError{Offset: pos, Err: "missing close paren"})
		}
		p.lex()
		return x
	}

	if !p.isTag {
		if p.tok == "" {
			panic(&SyntaxError{Offset: p.pos, Err: "unexpected end of expression"})
		}
		panic(&SyntaxError{Offset: p.pos, Err: "unexpected token " + p.tok})
	}
	tok := p.tok
	p.lex()
	return tag(tok)
}

// lex finds and consumes the next token in the input stream.
// On return, p.tok is set to the token text,
// p.isTag reports whether the token was a tag,
// and p.pos records the byte offset of the start of the token in the input stream.
// If lex reaches the end of the input, p.tok is set to the empty string.
// For any other syntax error, lex panics with a SyntaxError.
func (p *exprParser) lex() {
	p.isTag = false
	for p.i < len(p.s) && (p.s[p.i] == ' ' || p.s[p.i] == '\t') {
		p.i++
	}
	if p.i >= len(p.s) {
		p.tok = ""
		p.pos = p.i
		return
	}
	switch p.s[p.i] {
	case '(', ')', '!':
		p.pos = p.i
		p.i++
		p.tok = p.s[p.pos:p.i]
		return

	case '&', '|':
		if p.i+1 >= len(p.s) || p.s[p.i+1] != p.s[p.i] {
			panic(&SyntaxError{Offset: p.i, Err: "invalid syntax at " + string(rune(p.s[p.i]))})
		}
		p.pos = p.i
		p.i += 2
		p.tok = p.s[p.pos:p.i]
		return
	}

	tag := p.s[p.i:]
	for i, c := range tag {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' && c != '.' {
			tag = tag[:i]
			break
		}
	}
	if tag == "" {
		c, _ := utf8.DecodeRuneInString(p.s[p.i:])
		panic(&SyntaxError{Offset: p.i, Err: "invalid syntax at " + string(c)})
	}

	p.pos = p.i
	p.i += len(tag)
	p.tok = p.s[p.pos:p.i]
	p.isTag = true
	return
}

// isValidTag reports whether the word is a valid build tag.
// Tags must be letters, digits, underscores or dots.
// Unlike in Go identifiers, all digits are fine (e.g., "386").
func isValidTag(word string) bool {
	if word == "" {
		return false
	}
	for _, c := range word {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' && c != '.' {
			return false
		}
	}
	return true
}

// pushNot applies DeMorgan's law to push negations down the expression,
// so that only tags are negated in the result.
// (It applies the rewrites !(X && Y) => (!X || !Y) and !(X || Y) => (!X && !Y).)
func pushNot(x Expr, not bool) Expr {
	switch x := x.(type) {
	default:
		// unreachable
		return x
	case *NotExpr:
		if _, ok := x.X.(*TagExpr); ok && !not {
			return x
		}
		return pushNot(x.X, !not)
	case *TagExpr:
		if not {
			return &NotExpr{X: x}
		}
		return x
	case *AndExpr:
		x1 := pushNot(x.X, not)
		y1 := pushNot(x.Y, not)
		if not {
			return or(x1, y1)
		}
		if x1 == x.X && y1 == x.Y {
			return x
		}
		return and(x1, y1)
	case *OrExpr:
		x1 := pushNot(x.X, not)
		y1 := pushNot(x.Y, not)
		if not {
			return and(x1, y1)
		}
		if x1 == x.X && y1 == x.Y {
			return x
		}
		return or(x1, y1)
	}
}

// appendSplitAnd appends x to list while splitting apart any top-level && expressions.
// For example, appendSplitAnd({W}, X && Y && Z) = {W, X, Y, Z}.
func appendSplitAnd(list []Expr, x Expr) []Expr {
	if x, ok := x.(*AndExpr); ok {
		list = appendSplitAnd(list, x.X)
		list = appendSplitAnd(list, x.Y)
		return list
	}
	return append(list, x)
}

// appendSplitOr appends x to list while splitting apart any top-level || expressions.
// For example, appendSplitOr({W}, X || Y || Z) = {W, X, Y, Z}.
func appendSplitOr(list []Expr, x Expr) []Expr {
	if x, ok := x.(*OrExpr); ok {
		list = appendSplitOr(list, x.X)
		list = appendSplitOr(list, x.Y)
		return list
	}
	return append(list, x)
}
