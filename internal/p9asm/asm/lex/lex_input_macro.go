// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package lex

import "text/scanner"

// #define 宏定义
func (in *_Input) defineMacro() {
	tok := in.stk.Next()
	if tok != scanner.Ident {
		in.expectText("expected identifier after # directive")
	}
	name := in.stk.Text()

	prevCol := in.stk.Col()
	tok = in.stk.Next()

	var args []string
	var tokens []Token
	if tok != '\n' && tok != scanner.EOF {
		// 需要小心区分 `#define A(x)` 和 `#define A (x)`
		// 第一个是带参数的宏, 第二个没有参数, 区别在于第二个中间有空格
		if tok == '(' && in.stk.Col() == prevCol+1 {
			acceptArg := true // 宏是带参数的
			args = []string{} // 空参数和 nil 是有区别的

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
	}

	if in.macros[name] != nil {
		in.Error("redefinition of macro:", name)
	}
	in.macros[name] = &_MacroDefine{
		name:   name,
		args:   args,
		tokens: tokens,
	}
}

// #define 宏展开
func (in *_Input) invokeMacro(macro *_MacroDefine) {
	// len(macro.args) == 0 和 nil 是有区别的, 前者表示有参数的宏, 后者表示普通替换
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
