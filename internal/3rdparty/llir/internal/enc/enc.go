// Package enc implements encoding of identifiers for LLVM IR assembly.
package enc

import (
	"fmt"
	"strconv"
	"strings"
)

// GlobalName encodes a global name to its LLVM IR assembly representation.
//
// Examples:
//    "foo" -> "@foo"
//    "a b" -> `@"a b"`
//    "世" -> `@"\E4\B8\96"`
//    "2" -> `@"2"`
//
// References:
//    http://www.llvm.org/docs/LangRef.html#identifiers
func GlobalName(name string) string {
	// Positive numeric global names are quoted to distinguish global names from
	// global IDs; e.g.
	//
	//    @"2"
	if _, err := strconv.ParseUint(name, 10, 64); err == nil {
		return `@"` + name + `"`
	}
	return "@" + EscapeIdent(name)
}

// GlobalID encodes a global ID to its LLVM IR assembly representation.
//
// Examples:
//    "42" -> "@42"
//
// References:
//    http://www.llvm.org/docs/LangRef.html#identifiers
func GlobalID(id int64) string {
	if id < 0 {
		panic(fmt.Errorf("negative global ID (%d); should be represented as global name", id))
	}
	return "@" + strconv.FormatInt(id, 10)
}

// LocalName encodes a local name to its LLVM IR assembly representation.
//
// Examples:
//    "foo" -> "%foo"
//    "a b" -> `%"a b"`
//    "世" -> `%"\E4\B8\96"`
//    "2" -> `%"2"`
//
// References:
//    http://www.llvm.org/docs/LangRef.html#identifiers
func LocalName(name string) string {
	// Positive numeric local names are quoted to distinguish local names from
	// local IDs; e.g.
	//
	//    %"2"
	if _, err := strconv.ParseUint(name, 10, 64); err == nil {
		return `%"` + name + `"`
	}
	return "%" + EscapeIdent(name)
}

// LocalID encodes a local ID to its LLVM IR assembly representation.
//
// Examples:
//    "42" -> "%42"
//
// References:
//    http://www.llvm.org/docs/LangRef.html#identifiers
func LocalID(id int64) string {
	if id < 0 {
		panic(fmt.Errorf("negative local ID (%d); should be represented as local name", id))
	}
	return "%" + strconv.FormatInt(id, 10)
}

// LabelName encodes a label name to its LLVM IR assembly representation.
//
// Examples:
//    "foo" -> "foo:"
//    "a b" -> `"a b":`
//    "世" -> `"\E4\B8\96":`
//    "2" -> `"2":`
//
// References:
//    http://www.llvm.org/docs/LangRef.html#identifiers
func LabelName(name string) string {
	// Positive numeric label names are quoted to distinguish label names from
	// label IDs; e.g.
	//
	//    "2":
	if _, err := strconv.ParseUint(name, 10, 64); err == nil {
		return `"` + name + `":`
	}
	return EscapeIdent(name) + ":"
}

// LabelID encodes a label ID to its LLVM IR assembly representation.
//
// Examples:
//    "42" -> 42:
//
// References:
//    http://www.llvm.org/docs/LangRef.html#identifiers
func LabelID(id int64) string {
	if id < 0 {
		panic(fmt.Errorf("negative label ID (%d); should be represented as label name", id))
	}
	return strconv.FormatInt(id, 10) + ":"
}

// TypeName encodes a type name to its LLVM IR assembly representation.
//
// Examples:
//    "foo" -> "%foo"
//    "a b" -> `%"a b"`
//    "世" -> `%"\E4\B8\96"`
//    "2" -> `%2`
//
// References:
//    http://www.llvm.org/docs/LangRef.html#identifiers
func TypeName(name string) string {
	return "%" + EscapeIdent(name)
}

// AttrGroupID encodes a attribute group ID to its LLVM IR assembly
// representation.
//
// Examples:
//    "42" -> "#42"
//
// References:
//    http://www.llvm.org/docs/LangRef.html#identifiers
func AttrGroupID(id int64) string {
	return "#" + strconv.FormatInt(id, 10)
}

// ComdatName encodes a comdat name to its LLVM IR assembly representation.
//
// Examples:
//    "foo" -> $%foo"
//    "a b" -> `$"a b"`
//    "世" -> `$"\E4\B8\96"`
//
// References:
//    http://www.llvm.org/docs/LangRef.html#identifiers
func ComdatName(name string) string {
	return "$" + EscapeIdent(name)
}

// MetadataName encodes a metadata name to its LLVM IR assembly representation.
//
// Examples:
//    "foo" -> "!foo"
//    "a b" -> `!a\20b`
//    "世" -> `!\E4\B8\96`
//
// References:
//    http://www.llvm.org/docs/LangRef.html#identifiers
func MetadataName(name string) string {
	valid := func(b byte) bool {
		return strings.IndexByte(tail, b) != -1
	}
	if strings.ContainsRune(decimal, rune(name[0])) {
		// Escape first character if digit, to distinguish named from unnamed
		// metadata.
		return "!" + `\3` + name[:1] + string(Escape([]byte(name[1:]), valid))
	}
	return "!" + string(Escape([]byte(name), valid))
}

// MetadataID encodes a metadata ID to its LLVM IR assembly representation.
//
// Examples:
//    "42" -> "!42"
//
// References:
//    http://www.llvm.org/docs/LangRef.html#identifiers
func MetadataID(id int64) string {
	return "!" + strconv.FormatInt(id, 10)
}

const (
	// decimal specifies the decimal digit characters.
	decimal = "0123456789"
	// upper specifies the uppercase letters.
	upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// lower specifies the lowercase letters.
	lower = "abcdefghijklmnopqrstuvwxyz"
	// alpha specifies the alphabetic characters.
	alpha = upper + lower
	// head is the set of valid characters for the first character of an
	// identifier.
	head = alpha + "$-._"
	// tail is the set of valid characters for the remaining characters of an
	// identifier (i.e. all characters in the identifier except the first). All
	// characters of a label may be from the tail set, even the first character.
	tail = head + decimal
	// quotedIdent is the set of valid characters in quoted identifiers, which
	// excludes ASCII control characters, double quote, backslash and extended
	// ASCII characters.
	quotedIdent = " !#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"
)

// EscapeIdent replaces any characters which are not valid in identifiers with
// corresponding hexadecimal escape sequence (\XX).
func EscapeIdent(s string) string {
	replace := false
	extra := 0
	for i := 0; i < len(s); i++ {
		if strings.IndexByte(tail, s[i]) == -1 {
			// Check if a replacement is required.
			//
			// Note, there are characters which are not valid in an identifier
			// (e.g. '#') but are valid in a quoted identifier, and therefore
			// require a replacement (i.e. quoted identifier), but no extra
			// characters for the escape sequence.
			replace = true
		}
		if strings.IndexByte(quotedIdent, s[i]) == -1 {
			// Two extra bytes are required for each byte not valid in a quoted
			// identifier; e.g.
			//
			//    "\t" -> `\09`
			//    "世" -> `\E4\B8\96`
			extra += 2
		}
	}
	if !replace {
		return s
	}
	// Replace invalid characters.
	const hextable = "0123456789ABCDEF"
	buf := make([]byte, len(s)+extra)
	j := 0
	for i := 0; i < len(s); i++ {
		b := s[i]
		if strings.IndexByte(quotedIdent, b) != -1 {
			buf[j] = b
			j++
			continue
		}
		buf[j] = '\\'
		buf[j+1] = hextable[b>>4]
		buf[j+2] = hextable[b&0x0F]
		j += 3
	}
	// Add surrounding quotes.
	return `"` + string(buf) + `"`
}

// EscapeString replaces any characters in s categorized as invalid in string
// literals with corresponding hexadecimal escape sequence (\XX).
func EscapeString(s []byte) string {
	valid := func(b byte) bool {
		return ' ' <= b && b <= '~' && b != '"' && b != '\\'
	}
	return string(Escape(s, valid))
}

// Escape replaces any characters in s categorized as invalid by the valid
// function with corresponding hexadecimal escape sequence (\XX).
func Escape(s []byte, valid func(b byte) bool) string {
	// Check if a replacement is required.
	extra := 0
	for i := 0; i < len(s); i++ {
		if !valid(s[i]) {
			// Two extra bytes are required for each invalid byte; e.g.
			//    "#" -> `\23`
			//    "世" -> `\E4\B8\96`
			extra += 2
		}
	}
	if extra == 0 {
		return string(s)
	}
	// Replace invalid characters.
	const hextable = "0123456789ABCDEF"
	buf := make([]byte, len(s)+extra)
	j := 0
	for i := 0; i < len(s); i++ {
		b := s[i]
		if valid(b) {
			buf[j] = b
			j++
			continue
		}
		buf[j] = '\\'
		buf[j+1] = hextable[b>>4]
		buf[j+2] = hextable[b&0x0F]
		j += 3
	}
	return string(buf)
}

// Unescape replaces hexadecimal escape sequences (\xx) in s with their
// corresponding characters.
func Unescape(s string) []byte {
	if !strings.ContainsRune(s, '\\') {
		return []byte(s)
	}
	j := 0
	buf := []byte(s)
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b == '\\' {
			if len(s) > i+1 && s[i+1] == '\\' {
				b = '\\'
				i++
			} else if len(s) > i+2 {
				x1, ok := unhex(s[i+1])
				if ok {
					x2, ok := unhex(s[i+2])
					if ok {
						b = x1<<4 | x2
						i += 2
					}
				}
			}
		}
		if i != j {
			buf[j] = b
		}
		j++
	}
	return buf[:j]
}

// Quote returns s as a double-quoted string literal.
func Quote(s []byte) string {
	return `"` + string(EscapeString(s)) + `"`
}

// Unquote interprets s as a double-quoted string literal, returning the string
// value that s quotes.
func Unquote(s string) []byte {
	if len(s) < 2 {
		panic(fmt.Errorf("invalid length of quoted string; expected >= 2, got %d", len(s)))
	}
	if !strings.HasPrefix(s, `"`) {
		panic(fmt.Errorf("invalid quoted string `%s`; missing quote character prefix", s))
	}
	if !strings.HasSuffix(s, `"`) {
		panic(fmt.Errorf("invalid quoted string `%s`; missing quote character suffix", s))
	}
	// Skip double-quotes.
	s = s[1 : len(s)-1]
	return Unescape(s)
}

// unhex returns the numeric value represented by the hexadecimal digit b. It
// returns false if b is not a hexadecimal digit.
func unhex(b byte) (v byte, ok bool) {
	// This is an adapted copy of the unhex function from the strconv package,
	// which is governed by a BSD-style license.
	switch {
	case '0' <= b && b <= '9':
		return b - '0', true
	case 'a' <= b && b <= 'f':
		return b - 'a' + 10, true
	case 'A' <= b && b <= 'F':
		return b - 'A' + 10, true
	}
	return 0, false
}
