// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/scanner"

	"wa-lang.org/wa/internal/p9asm/asm/arch"
)

var _ TokenReader = (*_Input)(nil)

// _Input is the main input: a stack of readers and some macro definitions.
// It also handles #include processing (by pushing onto the input stack)
// and parses and instantiates macro definitions.
type _Input struct {
	stk             _Stack
	includes        []string
	beginningOfLine bool
	ifdefStack      []bool
	macros          map[string]*_MacroDefine
	text            string // Text of last token returned by Next.
	peek            bool
	peekToken       ScanToken
	peekText        string
}

// #define 宏指令
type _MacroDefine struct {
	name   string   // 宏的名字
	args   []string // 宏函数的参数部分
	tokens []Token  // 宏主体内容
}

func newInput(name string, flags *arch.Flags) (*_Input, error) {
	if flags == nil {
		flags = new(arch.Flags)
	}

	// 解析命令行预定义的宏
	// 比如 -D -VERSION=123
	macros := make(map[string]*_MacroDefine)
	for _, name := range flags.Defines {
		value := "1"
		i := strings.IndexRune(name, '=')
		if i > 0 {
			name, value = name[:i], name[i+1:]
		}
		tokens := LexString(name)
		if len(tokens) != 1 || tokens[0].ScanToken != scanner.Ident {
			return nil, fmt.Errorf("asm: parsing -D: %q is not a valid identifier name", tokens[0])
		}
		macros[name] = &_MacroDefine{
			name:   name,
			tokens: LexString(value),
		}
	}

	p := &_Input{
		// include directories: look in source dir, then -I directories.
		includes:        append([]string{filepath.Dir(name)}, flags.IncludeDirs...),
		beginningOfLine: true,
		macros:          macros,
	}
	return p, nil
}

func (in *_Input) Error(args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s:%d: %s", in.stk.File(), in.stk.Line(), fmt.Sprintln(args...))
	os.Exit(1)
}

// expectText is like Error but adds "got XXX" where XXX is a quoted representation of the most recent token.
func (in *_Input) expectText(args ...interface{}) {
	in.Error(append(args, "; got", strconv.Quote(in.stk.Text()))...)
}

// enabled reports whether the input is enabled by an ifdef, or is at the top level.
func (in *_Input) enabled() bool {
	return len(in.ifdefStack) == 0 || in.ifdefStack[len(in.ifdefStack)-1]
}

func (in *_Input) expectNewline(directive string) {
	tok := in.stk.Next()
	if tok != '\n' {
		in.expectText("expected newline after", directive)
	}
}

func (in *_Input) Next() ScanToken {
	if in.peek {
		in.peek = false
		tok := in.peekToken
		in.text = in.peekText
		return tok
	}
	// If we cannot generate a token after 100 macro invocations, we're in trouble.
	// The usual case is caught by Push, below, but be safe.
	for nesting := 0; nesting < 100; {
		tok := in.stk.Next()
		switch tok {
		case '#':
			if !in.beginningOfLine {
				in.Error("'#' must be first item on line")
			}
			in.beginningOfLine = in.hash()
		case scanner.Ident:
			// Is it a macro name?
			name := in.stk.Text()
			macro := in.macros[name]
			if macro != nil {
				nesting++
				in.invokeMacro(macro)
				continue
			}
			fallthrough
		default:
			in.beginningOfLine = tok == '\n'
			if in.enabled() {
				in.text = in.stk.Text()
				return tok
			}
		}
	}
	in.Error("recursive macro invocation")
	return 0
}

func (in *_Input) Text() string {
	return in.text
}

func (in *_Input) File() string {
	return in.stk.File()
}

func (in *_Input) Line() int {
	return in.stk.Line()
}

func (in *_Input) Col() int {
	return in.stk.Col()
}

func (in *_Input) SetPos(line int, file string) {
	in.stk.SetPos(line, file)
}

// hash processes a # preprocessor directive. It returns true iff it completes.
func (in *_Input) hash() bool {
	// We have a '#'; it must be followed by a known word (define, include, etc.).
	tok := in.stk.Next()
	if tok != scanner.Ident {
		in.expectText("expected identifier after '#'")
	}
	if !in.enabled() {
		// Can only start including again if we are at #else or #endif.
		// We let #line through because it might affect errors.
		switch in.stk.Text() {
		case "else", "endif", "line":
			// Press on.
		default:
			return false
		}
	}
	switch in.stk.Text() {
	case "define":
		in.define()
	case "else":
		in.else_()
	case "endif":
		in.endif()
	case "ifdef":
		in.ifdef(true)
	case "ifndef":
		in.ifdef(false)
	case "include":
		in.include()
	case "line":
		in.line()
	case "undef":
		in.undef()
	default:
		in.Error("unexpected token after '#':", in.stk.Text())
	}
	return true
}

// macroName returns the name for the macro being referenced.
func (in *_Input) macroName() string {
	// We use the Stack's input method; no macro processing at this stage.
	tok := in.stk.Next()
	if tok != scanner.Ident {
		in.expectText("expected identifier after # directive")
	}
	// Name is alphanumeric by definition.
	return in.stk.Text()
}

// #define processing.
func (in *_Input) define() {
	name := in.macroName()
	args, tokens := in.macroDefinition(name)
	in.defineMacro(name, args, tokens)
}

// defineMacro stores the macro definition in the _Input.
func (in *_Input) defineMacro(name string, args []string, tokens []Token) {
	if in.macros[name] != nil {
		in.Error("redefinition of macro:", name)
	}
	in.macros[name] = &_MacroDefine{
		name:   name,
		args:   args,
		tokens: tokens,
	}
}

// macroDefinition returns the list of formals and the tokens of the definition.
// The argument list is nil for no parens on the definition; otherwise a list of
// formal argument names.
func (in *_Input) macroDefinition(name string) ([]string, []Token) {
	prevCol := in.stk.Col()
	tok := in.stk.Next()
	if tok == '\n' || tok == scanner.EOF {
		return nil, nil // No definition for macro
	}
	var args []string
	// The C preprocessor treats
	//	#define A(x)
	// and
	//	#define A (x)
	// distinctly: the first is a macro with arguments, the second without.
	// Distinguish these cases using the column number, since we don't
	// see the space itself. Note that text/scanner reports the position at the
	// end of the token. It's where you are now, and you just read this token.
	if tok == '(' && in.stk.Col() == prevCol+1 {
		// Macro has arguments. Scan list of formals.
		acceptArg := true
		args = []string{} // Zero length but not nil.
	Loop:
		for {
			tok = in.stk.Next()
			switch tok {
			case ')':
				tok = in.stk.Next() // First token of macro definition.
				break Loop
			case ',':
				if acceptArg {
					in.Error("bad syntax in definition for macro:", name)
				}
				acceptArg = true
			case scanner.Ident:
				if !acceptArg {
					in.Error("bad syntax in definition for macro:", name)
				}
				arg := in.stk.Text()
				if i := lookup(args, arg); i >= 0 {
					in.Error("duplicate argument", arg, "in definition for macro:", name)
				}
				args = append(args, arg)
				acceptArg = false
			default:
				in.Error("bad definition for macro:", name)
			}
		}
	}
	var tokens []Token
	// Scan to newline. Backslashes escape newlines.
	for tok != '\n' {
		if tok == '\\' {
			tok = in.stk.Next()
			if tok != '\n' && tok != '\\' {
				in.Error(`can only escape \ or \n in definition for macro:`, name)
			}
		}
		tokens = append(tokens, Token{tok, in.stk.Text()})
		tok = in.stk.Next()
	}
	return args, tokens
}

func lookup(args []string, arg string) int {
	for i, a := range args {
		if a == arg {
			return i
		}
	}
	return -1
}

// invokeMacro pushes onto the input Stack a Slice that holds the macro definition with the actual
// parameters substituted for the formals.
// Invoking a macro does not touch the PC/line history.
func (in *_Input) invokeMacro(macro *_MacroDefine) {
	// If the macro has no arguments, just substitute the text.
	if macro.args == nil {
		in.Push(newSlice(in.stk.File(), in.stk.Line(), macro.tokens))
		return
	}
	tok := in.stk.Next()
	if tok != '(' {
		// If the macro has arguments but is invoked without them, all we push is the macro name.
		// First, put back the token.
		in.peekToken = tok
		in.peekText = in.text
		in.peek = true
		in.Push(newSlice(in.stk.File(), in.stk.Line(), []Token{{macroName, macro.name}}))
		return
	}
	actuals := in.argsFor(macro)
	var tokens []Token
	for _, tok := range macro.tokens {
		if tok.ScanToken != scanner.Ident {
			tokens = append(tokens, tok)
			continue
		}
		substitution := actuals[tok.Text]
		if substitution == nil {
			tokens = append(tokens, tok)
			continue
		}
		tokens = append(tokens, substitution...)
	}
	in.Push(newSlice(in.stk.File(), in.stk.Line(), tokens))
}

// argsFor returns a map from formal name to actual value for this argumented macro invocation.
// The opening parenthesis has been absorbed.
func (in *_Input) argsFor(macro *_MacroDefine) map[string][]Token {
	var args [][]Token
	// One macro argument per iteration. Collect them all and check counts afterwards.
	for argNum := 0; ; argNum++ {
		tokens, tok := in.collectArgument(macro)
		args = append(args, tokens)
		if tok == ')' {
			break
		}
	}
	// Zero-argument macros are tricky.
	if len(macro.args) == 0 && len(args) == 1 && args[0] == nil {
		args = nil
	} else if len(args) != len(macro.args) {
		in.Error("wrong arg count for macro", macro.name)
	}
	argMap := make(map[string][]Token)
	for i, arg := range args {
		argMap[macro.args[i]] = arg
	}
	return argMap
}

// collectArgument returns the actual tokens for a single argument of a macro.
// It also returns the token that terminated the argument, which will always
// be either ',' or ')'. The starting '(' has been scanned.
func (in *_Input) collectArgument(macro *_MacroDefine) ([]Token, ScanToken) {
	nesting := 0
	var tokens []Token
	for {
		tok := in.stk.Next()
		if tok == scanner.EOF || tok == '\n' {
			in.Error("unterminated arg list invoking macro:", macro.name)
		}
		if nesting == 0 && (tok == ')' || tok == ',') {
			return tokens, tok
		}
		if tok == '(' {
			nesting++
		}
		if tok == ')' {
			nesting--
		}
		tokens = append(tokens, Token{tok, in.stk.Text()})
	}
}

// #ifdef and #ifndef processing.
func (in *_Input) ifdef(truth bool) {
	name := in.macroName()
	in.expectNewline("#if[n]def")
	if _, defined := in.macros[name]; !defined {
		truth = !truth
	}
	in.ifdefStack = append(in.ifdefStack, truth)
}

// #else processing
func (in *_Input) else_() {
	in.expectNewline("#else")
	if len(in.ifdefStack) == 0 {
		in.Error("unmatched #else")
	}
	in.ifdefStack[len(in.ifdefStack)-1] = !in.ifdefStack[len(in.ifdefStack)-1]
}

// #endif processing.
func (in *_Input) endif() {
	in.expectNewline("#endif")
	if len(in.ifdefStack) == 0 {
		in.Error("unmatched #endif")
	}
	in.ifdefStack = in.ifdefStack[:len(in.ifdefStack)-1]
}

// #include processing.
func (in *_Input) include() {
	// Find and parse string.
	tok := in.stk.Next()
	if tok != scanner.String {
		in.expectText("expected string after #include")
	}
	name, err := strconv.Unquote(in.stk.Text())
	if err != nil {
		in.Error("unquoting include file name: ", err)
	}
	in.expectNewline("#include")
	// Push tokenizer for file onto stack.
	fd, err := os.Open(name)
	if err != nil {
		for _, dir := range in.includes {
			fd, err = os.Open(filepath.Join(dir, name))
			if err == nil {
				break
			}
		}
		if err != nil {
			in.Error("#include:", err)
		}
	}
	in.Push(newTokenizer(name, fd, fd))
}

// #line processing.
func (in *_Input) line() {
	// Only need to handle Plan 9 format: #line 337 "filename"
	tok := in.stk.Next()
	if tok != scanner.Int {
		in.expectText("expected line number after #line")
	}
	line, err := strconv.Atoi(in.stk.Text())
	if err != nil {
		in.Error("error parsing #line (cannot happen):", err)
	}
	tok = in.stk.Next()
	if tok != scanner.String {
		in.expectText("expected file name in #line")
	}
	file, err := strconv.Unquote(in.stk.Text())
	if err != nil {
		in.Error("unquoting #line file name: ", err)
	}
	tok = in.stk.Next()
	if tok != '\n' {
		in.Error("unexpected token at end of #line: ", tok)
	}
	in.stk.SetPos(line, file)
}

// #undef processing
func (in *_Input) undef() {
	name := in.macroName()
	if in.macros[name] == nil {
		in.Error("#undef for undefined macro:", name)
	}
	// Newline must be next.
	tok := in.stk.Next()
	if tok != '\n' {
		in.Error("syntax error in #undef for macro:", name)
	}
	delete(in.macros, name)
}

func (in *_Input) Push(r TokenReader) {
	if len(in.stk.tr) > 100 {
		in.Error("input recursion")
	}
	in.stk.Push(r)
}

func (in *_Input) Close() {
	in.stk.Close()
}
