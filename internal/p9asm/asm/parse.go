// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package asm implements the parser and instruction generator for the assembler.
// TODO: Split apart?
package asm

import (
	"text/scanner"

	"wa-lang.org/wa/internal/p9asm/asm/arch"
	"wa-lang.org/wa/internal/p9asm/asm/lex"
	"wa-lang.org/wa/internal/p9asm/obj"
)

type Parser struct {
	Flags *arch.Flags

	lex           lex.TokenReader
	lineNum       int   // Line number in source file.
	histLineNum   int32 // Cumulative line number across source files.
	errorLine     int32 // (Cumulative) line number of last error.
	errorCount    int   // Number of errors.
	pc            int64 // virtual PC; count of Progs; doesn't advance for GLOBL or DATA.
	input         []lex.Token
	inputPos      int
	pendingLabels []string // Labels to attach to next instruction.
	labels        map[string]*obj.Prog
	toPatch       []Patch
	addr          []obj.Addr
	arch          *arch.Arch
	ctxt          *obj.Link
	firstProg     *obj.Prog
	lastProg      *obj.Prog
	dataAddr      map[string]int64 // Most recent address for DATA for this symbol.

	panicOnError bool // 将错误以异常方式抛出, 用于单元测试
}

type Patch struct {
	prog  *obj.Prog
	label string
}

func NewParser(ctxt *obj.Link, ar *arch.Arch, lexer lex.TokenReader, flags *arch.Flags) *Parser {
	if flags == nil {
		flags = new(arch.Flags)
	}
	return &Parser{
		Flags:    flags,
		ctxt:     ctxt,
		arch:     ar,
		lex:      lexer,
		labels:   make(map[string]*obj.Prog),
		dataAddr: make(map[string]int64),
	}
}

// TODO(chai2010): 返回错误
func (p *Parser) Parse() (*obj.Prog, bool) {
	// 每次解析一行
	for p.parseLine() {
	}
	if p.errorCount > 0 {
		return nil, false
	}

	// 回填 label 跳转的地址
	for _, patch := range p.toPatch {
		targetProg := p.labels[patch.label]
		if targetProg == nil {
			p.errorf("undefined label %s", patch.label)
		} else {
			patch.prog.To = obj.Addr{
				Type: obj.TYPE_BRANCH,
				Val:  targetProg,
			}
		}
	}

	// 清空 toPatch 列表, 已经完成
	p.toPatch = p.toPatch[:0]

	return p.firstProg, true
}

// WORD [ arg {, arg} ] (';' | '\n')
func (p *Parser) parseLine() bool {
	// Skip newlines.
	var tok lex.ScanToken
	for {
		tok = p.lex.Next()
		// We save the line number here so error messages from this instruction
		// are labeled with this line. Otherwise we complain after we've absorbed
		// the terminating newline and the line numbers are off by one in errors.
		p.lineNum = p.lex.Line()
		p.histLineNum = lex.HistLine()
		switch tok {
		case '\n', ';':
			continue
		case scanner.EOF:
			return false
		}
		break
	}
	// First item must be an identifier.
	if tok != scanner.Ident {
		p.errorf("expected identifier, found %q", p.lex.Text())
		return false // Might as well stop now.
	}
	word := p.lex.Text()
	var cond string
	operands := make([][]lex.Token, 0, 3)
	// Zero or more comma-separated operands, one per loop.
	nesting := 0
	colon := -1
	for tok != '\n' && tok != ';' {
		// Process one operand.
		items := make([]lex.Token, 0, 3)
		for {
			tok = p.lex.Next()
			if len(operands) == 0 && len(items) == 0 {
				if (p.arch.LinkArch.Thechar == '5' || p.arch.LinkArch.Thechar == '7') && tok == '.' {
					// ARM conditionals.
					tok = p.lex.Next()
					str := p.lex.Text()
					if tok != scanner.Ident {
						p.errorf("ARM condition expected identifier, found %s", str)
					}
					cond = cond + "." + str
					continue
				}
				if tok == ':' {
					// Labels.
					p.pendingLabels = append(p.pendingLabels, word)
					return true
				}
			}
			if tok == scanner.EOF {
				p.errorf("unexpected EOF")
				return false
			}
			// Split operands on comma. Also, the old syntax on x86 for a "register pair"
			// was AX:DX, for which the new syntax is DX, AX. Note the reordering.
			if tok == '\n' || tok == ';' || (nesting == 0 && (tok == ',' || tok == ':')) {
				if tok == ':' {
					// Remember this location so we can swap the operands below.
					if colon >= 0 {
						p.errorf("invalid ':' in operand")
					}
					colon = len(operands)
				}
				break
			}
			if tok == '(' || tok == '[' {
				nesting++
			}
			if tok == ')' || tok == ']' {
				nesting--
			}
			items = append(items, lex.Make(tok, p.lex.Text()))
		}
		if len(items) > 0 {
			operands = append(operands, items)
			if colon >= 0 && len(operands) == colon+2 {
				// AX:DX becomes DX, AX.
				operands[colon], operands[colon+1] = operands[colon+1], operands[colon]
				colon = -1
			}
		} else if len(operands) > 0 || tok == ',' || colon >= 0 {
			// Had a separator with nothing after.
			p.errorf("missing operand")
		}
	}

	switch word {
	case "DATA":
		p.asmData(operands)
		return true
	case "FUNCDATA":
		p.asmFuncData(operands)
		return true
	case "GLOBL":
		p.asmGlobl(operands)
		return true
	case "PCDATA":
		p.asmPCData(operands)
		return true
	case "TEXT":
		p.asmText(operands)
		return true
	}

	if i, ok := p.arch.Instructions[word]; ok {
		p.instruction(i, word, cond, operands)
		return true
	}
	p.errorf("unrecognized instruction %q", word)
	return true
}

// address parses the operand into a link address structure.
func (p *Parser) address(operand []lex.Token) obj.Addr {
	p.input = operand
	p.inputPos = 0

	addr := obj.Addr{}
	p.operand(&addr)
	return addr
}
