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
	"wa-lang.org/wa/internal/p9asm/obj"
	"wa-lang.org/wa/internal/p9asm/objabi"
)

var _ TokenReader = (*_Input)(nil)

// _Input is the main input: a stack of readers and some macro definitions.
// It also handles #include processing (by pushing onto the input stack)
// and parses and instantiates macro definitions.
type _Input struct {
	ctxt            *obj.Link
	stk             _Stack
	includes        []string
	beginningOfLine bool
	macros          map[string]*_MacroDefine
	text            string // Text of last token returned by Next.
	peek            bool
	peekToken       ScanToken
	peekText        string
}

// #define 宏指令
type _MacroDefine struct {
	name   string   // 宏的名字
	args   []string // 宏函数的参数部分, len==0 和 nil 需要严格区分s
	tokens []Token  // 宏主体内容
}

func newInput(ctxt *obj.Link, name string, flags *arch.Flags) (*_Input, error) {
	if ctxt == nil {
		ctxt = new(obj.Link)
	}
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
		ctxt:            ctxt,
		includes:        append([]string{filepath.Dir(name)}, flags.IncludeDirs...),
		beginningOfLine: true,
		macros:          macros,
	}
	return p, nil
}

func (in *_Input) Push(r TokenReader) {
	if len(in.stk.tr) > 100 {
		in.Error("input recursion")
	}
	in.stk.Push(r)
}

func (in *_Input) Text() string { return in.text }
func (in *_Input) File() string { return in.stk.File() }
func (in *_Input) Line() int    { return in.stk.Line() }
func (in *_Input) Col() int     { return in.stk.Col() }

// TODO: 补充 pos 信息
func (in *_Input) Pos() objabi.Pos              { return objabi.NoPos }
func (in *_Input) SetPos(line int, file string) { in.stk.SetPos(line, file) }

func (in *_Input) Close() { in.stk.Close() }

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
			// 仅支持 #define 宏指令
			if !in.beginningOfLine {
				in.Error("'#' must be first item on line")
			}

			tok := in.stk.Next()
			if tok != scanner.Ident {
				in.expectText("expected identifier after '#'")
			}
			switch in.stk.Text() {
			case "define":
				in.defineMacro()
			default:
				in.Error("unexpected token after '#':", in.stk.Text())
			}

			in.beginningOfLine = true
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
			in.text = in.stk.Text()
			return tok
		}
	}
	in.Error("recursive macro invocation")
	return 0
}

func (in *_Input) Error(args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s:%d: %s", in.stk.File(), in.stk.Line(), fmt.Sprintln(args...))
	os.Exit(1)
}

// expectText is like Error but adds "got XXX" where XXX is a quoted representation of the most recent token.
func (in *_Input) expectText(args ...interface{}) {
	in.Error(append(args, "; got", strconv.Quote(in.stk.Text()))...)
}

func lookup(args []string, arg string) int {
	for i, a := range args {
		if a == arg {
			return i
		}
	}
	return -1
}
