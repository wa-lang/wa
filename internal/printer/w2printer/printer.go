// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package printer implements printing of AST nodes.
package w2printer

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

const (
	maxNewlines = 2     // max. number of newlines between source text
	debug       = false // enable for debugging
	infinity    = 1 << 30
)

type whiteSpace byte

const (
	ignore   = whiteSpace(0)
	blank    = whiteSpace(' ')
	vtab     = whiteSpace('\v')
	newline  = whiteSpace('\n')
	formfeed = whiteSpace('\f')
	indent   = whiteSpace('>')
	unindent = whiteSpace('<')
)

// A pmode value represents the current printer mode.
type pmode int

const (
	noExtraBlank     pmode = 1 << iota // disables extra blank after /*-style comment
	noExtraLinebreak                   // disables extra line break after /*-style comment
)

type commentInfo struct {
	cindex         int               // current comment index
	comment        *ast.CommentGroup // = printer.comments[cindex]; or nil
	commentOffset  int               // = printer.posFor(printer.comments[cindex].List[0].Pos()).Offset; or infinity
	commentNewline bool              // true if the comment group contains newlines
}

type printer struct {
	// Configuration (does not change after initialization)
	Config
	fset *token.FileSet

	// Current state
	output       []byte       // raw printer result
	indent       int          // current indentation
	level        int          // level == 0: outside composite literal; level > 0: inside composite literal
	mode         pmode        // current printer mode
	endAlignment bool         // if set, terminate alignment immediately
	impliedSemi  bool         // if set, a linebreak implies a semicolon
	lastTok      token.Token  // last token printed (token.ILLEGAL if it's whitespace)
	prevOpen     token.Token  // previous non-brace "open" token (, [, or token.ILLEGAL
	wsbuf        []whiteSpace // delayed white space

	// Positions
	// The out position differs from the pos position when the result
	// formatting differs from the source formatting (in the amount of
	// white space). If there's a difference and SourcePos is set in
	// ConfigMode, //line directives are used in the output to restore
	// original source positions for a reader.
	pos     token.Position // current position in AST (source) space
	out     token.Position // current position in output space
	last    token.Position // value of pos after calling writeString
	linePtr *int           // if set, record out.Line for the next token in *linePtr

	// The list of all source comments, in order of appearance.
	comments        []*ast.CommentGroup // may be nil
	useNodeComments bool                // if not set, ignore lead and line comments of nodes

	// Information about p.comments[p.cindex]; set up by nextComment.
	commentInfo

	// Cache of already computed node sizes.
	nodeSizes map[ast.Node]int

	// Cache of most recently computed line position.
	cachedPos  token.Pos
	cachedLine int // line corresponding to cachedPos
}

func (p *printer) init(cfg *Config, fset *token.FileSet, nodeSizes map[ast.Node]int) {
	p.Config = *cfg
	p.fset = fset
	p.pos = token.Position{Line: 1, Column: 1}
	p.out = token.Position{Line: 1, Column: 1}
	p.wsbuf = make([]whiteSpace, 0, 16) // whitespace sequences are short
	p.nodeSizes = nodeSizes
	p.cachedPos = -1
}

// fprint implements Fprint and takes a nodesSizes map for setting up the printer state.
func (cfg *Config) fprint(output io.Writer, fset *token.FileSet, node interface{}, nodeSizes map[ast.Node]int) (err error) {
	// print node
	var p printer
	p.init(cfg, fset, nodeSizes)
	if err = p.printNode(node); err != nil {
		return
	}
	// print outstanding comments
	p.impliedSemi = false // EOF acts like a newline
	p.flush(token.Position{Offset: infinity, Line: infinity}, token.EOF)

	// redirect output through a trimmer to eliminate trailing whitespace
	// (Input to a tabwriter must be untrimmed since trailing tabs provide
	// formatting information. The tabwriter could provide trimming
	// functionality but no tabwriter is used when RawFormat is set.)
	output = &trimmer{output: output}

	// redirect output through a tabwriter if necessary
	if cfg.Mode&RawFormat == 0 {
		minwidth := cfg.Tabwidth

		padchar := byte('\t')
		if cfg.Mode&UseSpaces != 0 {
			padchar = ' '
		}

		twmode := tabwriter.DiscardEmptyColumns
		if cfg.Mode&TabIndent != 0 {
			minwidth = 0
			twmode |= tabwriter.TabIndent
		}

		output = tabwriter.NewWriter(output, minwidth, cfg.Tabwidth, 1, padchar, twmode)
	}

	// write printer result via tabwriter/trimmer to output
	if _, err = output.Write(p.output); err != nil {
		return
	}

	// flush tabwriter, if any
	if tw, _ := output.(*tabwriter.Writer); tw != nil {
		err = tw.Flush()
	}

	return
}

func (p *printer) flush(next token.Position, tok token.Token) (wroteNewline, droppedFF bool) {
	if p.commentBefore(next) {
		// if there are comments before the next item, intersperse them
		wroteNewline, droppedFF = p.intersperseComments(next, tok)
	} else {
		// otherwise, write any leftover whitespace
		p.writeWhitespace(len(p.wsbuf))
	}
	return
}

// print prints a list of "items" (roughly corresponding to syntactic
// tokens, but also including whitespace and formatting information).
// It is the only print function that should be called directly from
// any of the AST printing functions in nodes.go.
//
// Whitespace is accumulated until a non-whitespace token appears. Any
// comments that need to appear before that token are printed first,
// taking into account the amount and structure of any pending white-
// space for best comment placement. Then, any leftover whitespace is
// printed, followed by the actual token.
func (p *printer) print(args ...interface{}) {
	for _, arg := range args {
		// information about the current arg
		var data string
		var isLit bool
		var impliedSemi bool // value for p.impliedSemi after this arg

		// record previous opening token, if any
		switch p.lastTok {
		case token.ILLEGAL:
			// ignore (white space)
		case token.LPAREN, token.LBRACK:
			p.prevOpen = p.lastTok
		default:
			// other tokens followed any opening token
			p.prevOpen = token.ILLEGAL
		}

		switch x := arg.(type) {
		case pmode:
			// toggle printer mode
			p.mode ^= x
			continue

		case whiteSpace:
			if x == ignore {
				// don't add ignore's to the buffer; they
				// may screw up "correcting" unindents (see
				// LabeledStmt)
				continue
			}
			i := len(p.wsbuf)
			if i == cap(p.wsbuf) {
				// Whitespace sequences are very short so this should
				// never happen. Handle gracefully (but possibly with
				// bad comment placement) if it does happen.
				p.writeWhitespace(i)
				i = 0
			}
			p.wsbuf = p.wsbuf[0 : i+1]
			p.wsbuf[i] = x
			if x == newline || x == formfeed {
				// newlines affect the current state (p.impliedSemi)
				// and not the state after printing arg (impliedSemi)
				// because comments can be interspersed before the arg
				// in this case
				p.impliedSemi = false
			}
			p.lastTok = token.ILLEGAL
			continue

		case *ast.Ident:
			data = x.Name
			impliedSemi = true
			p.lastTok = token.IDENT

		case *ast.BasicLit:
			data = x.Value
			isLit = true
			impliedSemi = true
			p.lastTok = x.Kind

		case token.Token:
			s := x.String()
			if mayCombine(p.lastTok, s[0]) {
				// the previous and the current token must be
				// separated by a blank otherwise they combine
				// into a different incorrect token sequence
				// (except for token.INT followed by a '.' this
				// should never happen because it is taken care
				// of via binary expression formatting)
				if len(p.wsbuf) != 0 {
					p.internalError("whitespace buffer not empty")
				}
				p.wsbuf = p.wsbuf[0:1]
				p.wsbuf[0] = ' '
			}
			data = s
			// some keywords followed by a newline imply a semicolon
			switch x {
			case token.BREAK, token.CONTINUE, token.RETURN,
				token.INC, token.DEC, token.RPAREN, token.RBRACK, token.RBRACE:
				impliedSemi = true

			case token.Zh_跳出, token.Zh_继续, token.Zh_返回:
				impliedSemi = true
			}
			p.lastTok = x

		case token.Pos:
			if x.IsValid() {
				p.pos = p.posFor(x) // accurate position of next item
			}
			continue

		case string:
			// incorrect AST - print error message
			data = x
			isLit = true
			impliedSemi = true
			p.lastTok = token.STRING

		default:
			fmt.Fprintf(os.Stderr, "print: unsupported argument %v (%T)\n", arg, arg)
			panic("wa-lang.org/wa/internal/printer type")
		}
		// data != ""

		next := p.pos // estimated/accurate position of next item
		wroteNewline, droppedFF := p.flush(next, p.lastTok)

		// intersperse extra newlines if present in the source and
		// if they don't cause extra semicolons (don't do this in
		// flush as it will cause extra newlines at the end of a file)
		if !p.impliedSemi {
			n := nlimit(next.Line - p.pos.Line)
			// don't exceed maxNewlines if we already wrote one
			if wroteNewline && n == maxNewlines {
				n = maxNewlines - 1
			}
			if n > 0 {
				ch := byte('\n')
				if droppedFF {
					ch = '\f' // use formfeed since we dropped one before
				}
				p.writeByte(ch, n)
				impliedSemi = false
			}
		}

		// the next token starts now - record its line number if requested
		if p.linePtr != nil {
			*p.linePtr = p.out.Line
			p.linePtr = nil
		}

		p.writeString(next, data, isLit)
		p.impliedSemi = impliedSemi
	}
}

func (p *printer) internalError(msg ...interface{}) {
	if debug {
		fmt.Print(p.pos.String() + ": ")
		fmt.Println(msg...)
		panic("wa-lang.org/wa/internal/printer")
	}
}

// recordLine records the output line number for the next non-whitespace
// token in *linePtr. It is used to compute an accurate line number for a
// formatted construct, independent of pending (not yet emitted) whitespace
// or comments.
func (p *printer) recordLine(linePtr *int) {
	p.linePtr = linePtr
}

// linesFrom returns the number of output lines between the current
// output line and the line argument, ignoring any pending (not yet
// emitted) whitespace or comments. It is used to compute an accurate
// size (in number of lines) for a formatted construct.
func (p *printer) linesFrom(line int) int {
	return p.out.Line - line
}

func (p *printer) posFor(pos token.Pos) token.Position {
	// not used frequently enough to cache entire token.Position
	return p.fset.PositionFor(pos, false /* absolute position */)
}

func (p *printer) lineFor(pos token.Pos) int {
	if pos != p.cachedPos {
		p.cachedPos = pos
		p.cachedLine = p.fset.PositionFor(pos, false /* absolute position */).Line
	}
	return p.cachedLine
}

// distanceFrom returns the column difference between from and p.pos (the current
// estimated position) if both are on the same line; if they are on different lines
// (or unknown) the result is infinity.
func (p *printer) distanceFrom(from token.Pos) int {
	if from.IsValid() && p.pos.IsValid() {
		if f := p.posFor(from); f.Line == p.pos.Line {
			return p.pos.Column - f.Column
		}
	}
	return infinity
}

// Print as many newlines as necessary (but at least min newlines) to get to
// the current line. ws is printed before the first line break. If newSection
// is set, the first line break is printed as formfeed. Returns 0 if no line
// breaks were printed, returns 1 if there was exactly one newline printed,
// and returns a value > 1 if there was a formfeed or more than one newline
// printed.
//
// TODO(gri): linebreak may add too many lines if the next statement at "line"
//
//	is preceded by comments because the computation of n assumes
//	the current position before the comment and the target position
//	after the comment. Thus, after interspersing such comments, the
//	space taken up by them is not considered to reduce the number of
//	linebreaks. At the moment there is no easy way to know about
//	future (not yet interspersed) comments in this function.
func (p *printer) linebreak(line, min int, ws whiteSpace, newSection bool) (nbreaks int) {
	n := nlimit(line - p.pos.Line)
	if n < min {
		n = min
	}
	if n > 0 {
		p.print(ws)
		if newSection {
			p.print(formfeed)
			n--
			nbreaks = 2
		}
		nbreaks += n
		for ; n > 0; n-- {
			p.print(newline)
		}
	}
	return
}
