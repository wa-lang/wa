// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package obj

// A LineStack is an entry in the recorded line history.
// Although the history at any given line number is a stack,
// the record for all line processed forms a tree, with common
// stack prefixes acting as parents.
type LineStack struct {
	Parent    *LineStack // parent in inclusion stack
	Lineno    int        // virtual line number where this entry takes effect
	File      string     // file name used to open source file, for error messages
	AbsFile   string     // absolute file name, for pcln tables
	FileLine  int        // line number in file at Lineno
	Directive bool
}

func (stk *LineStack) fileLineAt(lineno int) int {
	return stk.FileLine + lineno - stk.Lineno
}

// The span of valid linenos in the recorded line history can be broken
// into a set of ranges, each with a particular stack.
// A LineRange records one such range.
type LineRange struct {
	Start int        // starting lineno
	Stack *LineStack // top of stack for this range
}

// Does s have t as a path prefix?
// That is, does s == t or does s begin with t followed by a slash?
// For portability, we allow ASCII case folding, so that hasPathPrefix("a/b/c", "A/B") is true.
// Similarly, we allow slash folding, so that hasPathPrefix("a/b/c", "a\\b") is true.
// We do not allow full Unicode case folding, for fear of causing more confusion
// or harm than good. (For an example of the kinds of things that can go wrong,
// see http://article.gmane.org/gmane.linux.kernel/1853266.)
func hasPathPrefix(s string, t string) bool {
	if len(t) > len(s) {
		return false
	}
	var i int
	for i = 0; i < len(t); i++ {
		cs := int(s[i])
		ct := int(t[i])
		if 'A' <= cs && cs <= 'Z' {
			cs += 'a' - 'A'
		}
		if 'A' <= ct && ct <= 'Z' {
			ct += 'a' - 'A'
		}
		if cs == '\\' {
			cs = '/'
		}
		if ct == '\\' {
			ct = '/'
		}
		if cs != ct {
			return false
		}
	}
	return i >= len(s) || s[i] == '/' || s[i] == '\\'
}
