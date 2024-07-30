// 版权 @2024 凹语言 作者。保留所有权利。

package scanner

import (
	"strings"
	"testing"

	"wa-lang.org/wa/internal/wat/token"
)

const /* class */ (
	special = iota
	literal
	operator
	keyword
	instruction
)

func tokenclass(tok token.Token) int {
	switch {
	case tok.IsLiteral():
		return literal
	case tok.IsOperator():
		return operator
	case tok.IsKeyword():
		if tok.IsIsntruction() {
			return instruction
		}
		return keyword
	case tok.IsIsntruction():
		return instruction
	}
	return special
}

type elt struct {
	tok   token.Token
	lit   string
	class int
}

var tokens = [...]elt{
	// 分隔符
	{token.LPAREN, "(", operator},
	{token.RPAREN, ")", operator},
	{token.ASSIGN, "=", operator},

	// 单行注释(不支持多行注释)
	{token.COMMENT, ";; a comment\n", special},

	// 标识符
	{token.IDENT, "$foobar", literal},
	{token.IDENT, "$$foobar", literal},
	{token.IDENT, "$foobar.abc.123", literal},
	{token.IDENT, "$bar9876", literal},

	// 字面值
	{token.INT, "0", literal},
	{token.INT, "1", literal},
	{token.INT, "123456789012345678890", literal},
	{token.INT, "01234567", literal},
	{token.INT, "0xcafebabe", literal},
	{token.FLOAT, "0.", literal},
	{token.FLOAT, ".0", literal},
	{token.FLOAT, "3.14159265", literal},
	{token.FLOAT, "1e0", literal},
	{token.FLOAT, "1e+100", literal},
	{token.FLOAT, "1e-100", literal},
	{token.FLOAT, "2.71828e-1000", literal},
	{token.CHAR, "'a'", literal},
	{token.STRING, `"wasi_snapshot_preview1"`, literal},
	//{token.STRING, `"abc\12\ab\CD\n\t\r"`, literal},

	// Keywords
	{token.I32, "i32", keyword},
	{token.I64, "i64", keyword},
	{token.F32, "f32", keyword},
	{token.F64, "f64", keyword},

	{token.MUT, "mut", keyword},
	{token.FUNCREF, "funcref", keyword},
	{token.OFFSET, "offset", keyword},

	{token.MODULE, "module", keyword},
	{token.IMPORT, "import", keyword},
	{token.EXPORT, "export", keyword},
	{token.MEMORY, "memory", keyword},
	{token.DATA, "data", keyword},
	{token.TABLE, "table", keyword},
	{token.ELEM, "elem", keyword},
	{token.TYPE, "type", keyword},
	{token.GLOBAL, "global", keyword},
	{token.FUNC, "func", keyword},
	{token.PARAM, "param", keyword},
	{token.RESULT, "result", keyword},
	{token.LOCAL, "local", keyword},
	{token.START, "start", keyword},

	// TODO: 指令
	{token.INS_GLOBAL_GET, "global.get", instruction},
	{token.INS_GLOBAL_SET, "global.set", instruction},
	{token.INS_LOCAL_GET, "local.get", instruction},
	{token.INS_LOCAL_SET, "local.set", instruction},

	{token.INS_I32_CONST, "i32.const", instruction},
	{token.INS_I64_CONST, "i64.const", instruction},
	{token.INS_F32_CONST, "f32.const", instruction},
	{token.INS_F64_CONST, "f64.const", instruction},

	{token.INS_I32_ADD, "i32.add", instruction},
	{token.INS_I64_ADD, "i64.add", instruction},
	{token.INS_F32_ADD, "f32.add", instruction},
	{token.INS_F64_ADD, "f64.add", instruction},

	{token.INS_I32_SUB, "i32.sub", instruction},
	{token.INS_I64_SUB, "i64.sub", instruction},
	{token.INS_F32_SUB, "f32.sub", instruction},
	{token.INS_F64_SUB, "f64.sub", instruction},

	{token.INS_I32_LE_S, "i32.le_s", instruction},
	{token.INS_I64_LE_S, "i64.le_s", instruction},

	{token.INS_I32_EQ, "i32.eq", instruction},
	{token.INS_I64_EQ, "i64.eq", instruction},
	{token.INS_F32_EQ, "f32.eq", instruction},
	{token.INS_F64_EQ, "f64.eq", instruction},

	{token.INS_I32_STORE, "i32.store", instruction},

	{token.INS_CALL, "call", instruction},
	{token.INS_CALL_INDIRECT, "call_indirect", instruction},

	{token.INS_NOP, "nop", instruction},
	{token.INS_UNREACHABLE, "unreachable", instruction},
	{token.INS_BLOCK, "block", instruction},
	{token.INS_END, "end", instruction},
	{token.INS_IF, "if", instruction},
	{token.INS_ELSE, "else", instruction},
	{token.INS_LOOP, "loop", instruction},
	{token.INS_BR, "br", instruction},
	{token.INS_BR_IF, "br_if", instruction},
	{token.INS_BR_TABLE, "br_table", instruction},
	{token.INS_DROP, "drop", instruction},
}

const whitespace = "  \t  \n\n\n" // to separate tokens

var source = func() []byte {
	var src []byte
	for _, t := range tokens {
		src = append(src, t.lit...)
		src = append(src, whitespace...)
	}
	return src
}()

func newlineCount(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			n++
		}
	}
	return n
}

// Verify that calling Scan() provides the correct results.
func TestScan(t *testing.T) {
	whitespace_linecount := newlineCount(whitespace)

	// error handler
	eh := func(pos token.Position, msg string) {
		t.Errorf("%v: error handler called (msg = %s)", pos, msg)
	}

	// verify scan
	var s Scanner
	s.Init(token.NewFile("", len(source)), source, eh, ScanComments|dontInsertSemis)

	// set up expected position
	epos := token.Position{Line: 1, Column: 1}

	index := 0
	for {
		_, tok, lit := s.Scan()

		// check position
		if tok == token.EOF {
			// correction for EOF
			epos.Line = newlineCount(string(source))
			epos.Column = 2
		}
		// checkPos(t, lit, pos, epos)

		// check token
		e := elt{token.EOF, "", special}
		if index < len(tokens) {
			e = tokens[index]
			index++
		}
		if tok != e.tok {
			t.Fatalf("%d: bad token for %q: got %s, expected %s", index, lit, tok, e.tok)
		}

		// check token class
		if tokenclass(tok) != e.class {
			t.Fatalf("%d: bad class for %q: got %d, expected %d", index, lit, tokenclass(tok), e.class)
		}

		// check literal
		elit := ""
		switch e.tok {
		case token.COMMENT:
			// no CRs in comments
			elit = string(stripCR([]byte(e.lit), e.lit[0] == '('))
			//-style comment literal doesn't contain newline
			if elit[0] == ';' {
				elit = elit[0 : len(elit)-1]
			}
		case token.IDENT:
			elit = e.lit[1:]
		case token.STRING:
			elit = e.lit[1 : len(e.lit)-1]
		default:
			if e.tok.IsLiteral() {
				// no CRs in raw string literals
				elit = e.lit
			} else if e.tok.IsKeyword() {
				elit = e.lit
			} else if e.tok.IsIsntruction() {
				elit = e.lit
			}
		}
		if lit != elit {
			t.Fatalf("bad literal for %q: got %q, expected %q", lit, lit, elit)
		}

		if tok == token.EOF {
			break
		}

		// update position
		epos.Offset += len(e.lit) + len(whitespace)
		epos.Line += newlineCount(e.lit) + whitespace_linecount

	}

	if s.ErrorCount != 0 {
		t.Errorf("found %d errors", s.ErrorCount)
	}
}

func TestStripCR(t *testing.T) {
	for _, test := range []struct{ have, want string }{
		{"//\n", "//\n"},
		{"//\r\n", "//\n"},
		{"//\r\r\r\n", "//\n"},
		{"//\r*\r/\r\n", "//*/\n"},
		{"/**/", "/**/"},
		{"/*\r/*/", "/*/*/"},
		{"/*\r*/", "/**/"},
		{"/**\r/*/", "/**\r/*/"},
		{"/*\r/\r*\r/*/", "/*/*\r/*/"},
		{"/*\r\r\r\r*/", "/**/"},
	} {
		got := string(stripCR([]byte(test.have), len(test.have) >= 2 && test.have[1] == '*'))
		if got != test.want {
			t.Errorf("stripCR(%q) = %q; want %q", test.have, got, test.want)
		}
	}
}

func TestNumbers(t *testing.T) {
	for k, test := range []struct {
		tok              token.Token
		src, tokens, err string
	}{
		// binaries
		{token.INT, "0b0", "0b0", ""},
		{token.INT, "0b1010", "0b1010", ""},
		{token.INT, "0B1110", "0B1110", ""},

		{token.INT, "0b", "0b", "binary literal has no digits"},
		{token.INT, "0b0190", "0b0190", "invalid digit '9' in binary literal"},
		{token.INT, "0b01a0", "0b01 a0", ""}, // only accept 0-9

		{token.FLOAT, "0b.", "0b.", "invalid radix point in binary literal"},
		{token.FLOAT, "0b.1", "0b.1", "invalid radix point in binary literal"},
		{token.FLOAT, "0b1.0", "0b1.0", "invalid radix point in binary literal"},
		{token.FLOAT, "0b1e10", "0b1e10", "'e' exponent requires decimal mantissa"},
		{token.FLOAT, "0b1P-1", "0b1P-1", "'P' exponent requires hexadecimal mantissa"},

		// octals
		{token.INT, "0o0", "0o0", ""},
		{token.INT, "0o1234", "0o1234", ""},
		{token.INT, "0O1234", "0O1234", ""},

		{token.INT, "0o", "0o", "octal literal has no digits"},
		{token.INT, "0o8123", "0o8123", "invalid digit '8' in octal literal"},
		{token.INT, "0o1293", "0o1293", "invalid digit '9' in octal literal"},
		{token.INT, "0o12a3", "0o12 a3", ""}, // only accept 0-9

		{token.FLOAT, "0o.", "0o.", "invalid radix point in octal literal"},
		{token.FLOAT, "0o.2", "0o.2", "invalid radix point in octal literal"},
		{token.FLOAT, "0o1.2", "0o1.2", "invalid radix point in octal literal"},
		{token.FLOAT, "0o1E+2", "0o1E+2", "'E' exponent requires decimal mantissa"},
		{token.FLOAT, "0o1p10", "0o1p10", "'p' exponent requires hexadecimal mantissa"},

		// 0-octals
		{token.INT, "0", "0", ""},
		{token.INT, "0123", "0123", ""},

		{token.INT, "08123", "08123", "invalid digit '8' in octal literal"},
		{token.INT, "01293", "01293", "invalid digit '9' in octal literal"},
		//{token.INT, "0F.", "0 F .", ""}, // only accept 0-9
		//{token.INT, "0123F.", "0123 F .", ""},
		//{token.INT, "0123456x", "0123456 x", ""},

		// decimals
		{token.INT, "1", "1", ""},
		{token.INT, "1234", "1234", ""},

		//{token.INT, "1f", "1 f", ""}, // only accept 0-9

		// decimal floats
		{token.FLOAT, "0.", "0.", ""},
		{token.FLOAT, "123.", "123.", ""},
		{token.FLOAT, "0123.", "0123.", ""},

		{token.FLOAT, ".0", ".0", ""},
		{token.FLOAT, ".123", ".123", ""},
		{token.FLOAT, ".0123", ".0123", ""},

		{token.FLOAT, "0.0", "0.0", ""},
		{token.FLOAT, "123.123", "123.123", ""},
		{token.FLOAT, "0123.0123", "0123.0123", ""},

		{token.FLOAT, "0e0", "0e0", ""},
		{token.FLOAT, "123e+0", "123e+0", ""},
		{token.FLOAT, "0123E-1", "0123E-1", ""},

		{token.FLOAT, "0.e+1", "0.e+1", ""},
		{token.FLOAT, "123.E-10", "123.E-10", ""},
		{token.FLOAT, "0123.e123", "0123.e123", ""},

		{token.FLOAT, ".0e-1", ".0e-1", ""},
		{token.FLOAT, ".123E+10", ".123E+10", ""},
		{token.FLOAT, ".0123E123", ".0123E123", ""},

		{token.FLOAT, "0.0e1", "0.0e1", ""},
		{token.FLOAT, "123.123E-10", "123.123E-10", ""},
		{token.FLOAT, "0123.0123e+456", "0123.0123e+456", ""},

		{token.FLOAT, "0e", "0e", "exponent has no digits"},
		{token.FLOAT, "0E+", "0E+", "exponent has no digits"},
		{token.FLOAT, "1e+f", "1e+ f", "exponent has no digits"},
		{token.FLOAT, "0p0", "0p0", "'p' exponent requires hexadecimal mantissa"},
		{token.FLOAT, "1.0P-1", "1.0P-1", "'P' exponent requires hexadecimal mantissa"},

		// hexadecimals
		{token.INT, "0x0", "0x0", ""},
		{token.INT, "0x1234", "0x1234", ""},
		{token.INT, "0xcafef00d", "0xcafef00d", ""},
		{token.INT, "0XCAFEF00D", "0XCAFEF00D", ""},

		{token.INT, "0x", "0x", "hexadecimal literal has no digits"},
		{token.INT, "0x1g", "0x1 g", ""},

		// hexadecimal floats
		{token.FLOAT, "0x0p0", "0x0p0", ""},
		{token.FLOAT, "0x12efp-123", "0x12efp-123", ""},
		{token.FLOAT, "0xABCD.p+0", "0xABCD.p+0", ""},
		{token.FLOAT, "0x.0189P-0", "0x.0189P-0", ""},
		{token.FLOAT, "0x1.ffffp+1023", "0x1.ffffp+1023", ""},

		{token.FLOAT, "0x.", "0x.", "hexadecimal literal has no digits"},
		{token.FLOAT, "0x0.", "0x0.", "hexadecimal mantissa requires a 'p' exponent"},
		{token.FLOAT, "0x.0", "0x.0", "hexadecimal mantissa requires a 'p' exponent"},
		{token.FLOAT, "0x1.1", "0x1.1", "hexadecimal mantissa requires a 'p' exponent"},
		{token.FLOAT, "0x1.1e0", "0x1.1e0", "hexadecimal mantissa requires a 'p' exponent"},
		{token.FLOAT, "0x1.2gp1a", "0x1.2 gp1a", "hexadecimal mantissa requires a 'p' exponent"},
		{token.FLOAT, "0x0p", "0x0p", "exponent has no digits"},
		{token.FLOAT, "0xeP-", "0xeP-", "exponent has no digits"},
		{token.FLOAT, "0x1234PAB", "0x1234P AB", "exponent has no digits"},
		{token.FLOAT, "0x1.2p1a", "0x1.2p1 a", ""},

		// separators
		{token.INT, "0b_1000_0001", "0b_1000_0001", ""},
		{token.INT, "0o_600", "0o_600", ""},
		{token.INT, "0_466", "0_466", ""},
		{token.INT, "1_000", "1_000", ""},
		{token.FLOAT, "1_000.000_1", "1_000.000_1", ""},
		{token.INT, "0x_f00d", "0x_f00d", ""},
		{token.FLOAT, "0x_f00d.0p1_2", "0x_f00d.0p1_2", ""},

		{token.INT, "0b__1000", "0b__1000", "'_' must separate successive digits"},
		{token.INT, "0o60___0", "0o60___0", "'_' must separate successive digits"},
		{token.INT, "0466_", "0466_", "'_' must separate successive digits"},
		{token.FLOAT, "1_.", "1_.", "'_' must separate successive digits"},
		{token.FLOAT, "0._1", "0._1", "'_' must separate successive digits"},
		{token.FLOAT, "2.7_e0", "2.7_e0", "'_' must separate successive digits"},
		{token.INT, "0x___0", "0x___0", "'_' must separate successive digits"},
		{token.FLOAT, "0x1.0_p0", "0x1.0_p0", "'_' must separate successive digits"},
	} {
		var s Scanner
		var err string
		s.Init(token.NewFile("", len(test.src)), []byte(test.src), func(_ token.Position, msg string) {
			if err == "" {
				err = msg
			}
		}, 0)
		for i, want := range strings.Split(test.tokens, " ") {
			err = ""
			_, tok, lit := s.Scan()

			if i == 0 {
				if tok != test.tok {
					t.Fatalf("%d: %q: got token %s; want %s", k, test.src, tok, test.tok)
				}
				if err != test.err {
					t.Fatalf("%d:%q: got error %q; want %q", k, test.src, err, test.err)
				}
			}

			if lit != want {
				t.Fatalf("%d: %q: got literal %q (%s); want %s", k, test.src, lit, tok, want)
			}
		}

		// make sure we read all
		_, tok, _ := s.Scan()
		if tok != token.EOF {
			t.Fatalf("%q: got %s; want EOF", test.src, tok)
		}
	}
}
