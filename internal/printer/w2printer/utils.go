// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"fmt"
	"strings"
	"unicode"

	"wa-lang.org/wa/internal/token"
)

func assert(cond bool, a ...string) {
	if !cond {
		if len(a) > 0 {
			msg := fmt.Sprint(a)
			panic("wa-lang.org/wa/internal/parser/w2printer internal error: " + msg)
		} else {
			panic("wa-lang.org/wa/internal/parser/w2printer internal error: " + "assert failed")
		}
	}
}

// nlines limits n to maxNewlines.
func nlimit(n int) int {
	if n > maxNewlines {
		n = maxNewlines
	}
	return n
}

func mayCombine(prev token.Token, next byte) (b bool) {
	switch prev {
	case token.INT:
		b = next == '.' // 1.
	case token.ADD:
		b = next == '+' // ++
	case token.SUB:
		b = next == '-' // --
	case token.QUO:
		b = next == '*' // /*
	case token.LSS:
		b = next == '-' || next == '<' // <- or <<
	case token.AND:
		b = next == '&' || next == '^' // && or &^
	}
	return
}

// Returns true if s contains only white space
// (only tabs and blanks can appear in the printer's context).
func isBlank(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > ' ' {
			return false
		}
	}
	return true
}

// commonPrefix returns the common prefix of a and b.
func commonPrefix(a, b string) string {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] && (a[i] <= ' ' || a[i] == '*') {
		i++
	}
	return a[0:i]
}

// trimRight returns s with trailing whitespace removed.
func trimRight(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// stripCommonPrefix removes a common prefix from /*-style comment lines (unless no
// comment line is indented, all but the first line have some form of space prefix).
// The prefix is computed using heuristics such that is likely that the comment
// contents are nicely laid out after re-printing each line using the printer's
// current indentation.
func stripCommonPrefix(lines []string) {
	if len(lines) <= 1 {
		return // at most one line - nothing to do
	}
	// len(lines) > 1

	// The heuristic in this function tries to handle a few
	// common patterns of /*-style comments: Comments where
	// the opening /* and closing */ are aligned and the
	// rest of the comment text is aligned and indented with
	// blanks or tabs, cases with a vertical "line of stars"
	// on the left, and cases where the closing */ is on the
	// same line as the last comment text.

	// Compute maximum common white prefix of all but the first,
	// last, and blank lines, and replace blank lines with empty
	// lines (the first line starts with /* and has no prefix).
	// In cases where only the first and last lines are not blank,
	// such as two-line comments, or comments where all inner lines
	// are blank, consider the last line for the prefix computation
	// since otherwise the prefix would be empty.
	//
	// Note that the first and last line are never empty (they
	// contain the opening /* and closing */ respectively) and
	// thus they can be ignored by the blank line check.
	prefix := ""
	prefixSet := false
	if len(lines) > 2 {
		for i, line := range lines[1 : len(lines)-1] {
			if isBlank(line) {
				lines[1+i] = "" // range starts with lines[1]
			} else {
				if !prefixSet {
					prefix = line
					prefixSet = true
				}
				prefix = commonPrefix(prefix, line)
			}

		}
	}
	// If we don't have a prefix yet, consider the last line.
	if !prefixSet {
		line := lines[len(lines)-1]
		prefix = commonPrefix(line, line)
	}

	/*
	 * Check for vertical "line of stars" and correct prefix accordingly.
	 */
	lineOfStars := false
	if i := strings.Index(prefix, "*"); i >= 0 {
		// Line of stars present.
		if i > 0 && prefix[i-1] == ' ' {
			i-- // remove trailing blank from prefix so stars remain aligned
		}
		prefix = prefix[0:i]
		lineOfStars = true
	} else {
		// No line of stars present.
		// Determine the white space on the first line after the /*
		// and before the beginning of the comment text, assume two
		// blanks instead of the /* unless the first character after
		// the /* is a tab. If the first comment line is empty but
		// for the opening /*, assume up to 3 blanks or a tab. This
		// whitespace may be found as suffix in the common prefix.
		first := lines[0]
		if isBlank(first[2:]) {
			// no comment text on the first line:
			// reduce prefix by up to 3 blanks or a tab
			// if present - this keeps comment text indented
			// relative to the /* and */'s if it was indented
			// in the first place
			i := len(prefix)
			for n := 0; n < 3 && i > 0 && prefix[i-1] == ' '; n++ {
				i--
			}
			if i == len(prefix) && i > 0 && prefix[i-1] == '\t' {
				i--
			}
			prefix = prefix[0:i]
		} else {
			// comment text on the first line
			suffix := make([]byte, len(first))
			n := 2 // start after opening /*
			for n < len(first) && first[n] <= ' ' {
				suffix[n] = first[n]
				n++
			}
			if n > 2 && suffix[2] == '\t' {
				// assume the '\t' compensates for the /*
				suffix = suffix[2:n]
			} else {
				// otherwise assume two blanks
				suffix[0], suffix[1] = ' ', ' '
				suffix = suffix[0:n]
			}
			// Shorten the computed common prefix by the length of
			// suffix, if it is found as suffix of the prefix.
			prefix = strings.TrimSuffix(prefix, string(suffix))
		}
	}

	// Handle last line: If it only contains a closing */, align it
	// with the opening /*, otherwise align the text with the other
	// lines.
	last := lines[len(lines)-1]
	closing := "*/"
	i := strings.Index(last, closing) // i >= 0 (closing is always present)
	if isBlank(last[0:i]) {
		// last line only contains closing */
		if lineOfStars {
			closing = " */" // add blank to align final star
		}
		lines[len(lines)-1] = prefix + closing
	} else {
		// last line contains more comment text - assume
		// it is aligned like the other lines and include
		// in prefix computation
		prefix = commonPrefix(prefix, last)
	}

	// Remove the common prefix from all but the first and empty lines.
	for i, line := range lines {
		if i > 0 && line != "" {
			lines[i] = line[len(prefix):]
		}
	}
}
