package llmetadata

import (
	"strings"

	"github.com/wa-lang/wa/internal/3rdparty/llir/internal/enc"
	enum "github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
)

// diFlagsString returns the string representation of the given debug
// information flags.
func diFlagsString(flags enum.DIFlag) string {
	if flags == enum.DIFlagZero {
		return flags.String()
	}
	var ss []string
	if flag := flags & 0x3; flag != 0 {
		ss = append(ss, flag.String())
	}
	for mask := enum.DIFlagFirst; mask <= enum.DIFlagLast; mask <<= 1 {
		if flags&mask != 0 {
			ss = append(ss, mask.String())
		}
	}
	return strings.Join(ss, " | ")
}

// dispFlagsString returns the string representation of the given subprogram
// specific flags.
func dispFlagsString(flags enum.DISPFlag) string {
	if flags == enum.DISPFlagZero {
		return flags.String()
	}
	var ss []string
	for mask := enum.DISPFlagFirst; mask <= enum.DISPFlagLast; mask <<= 1 {
		if flags&mask != 0 {
			ss = append(ss, mask.String())
		}
	}
	return strings.Join(ss, " | ")
}

// TODO: fix string representation for all enums which are defined in the
// grammar as `FooEnum | FooInt`, in the same way as dwarfTagString.

// dwarfTagString returns the string representation of the given Dwarf tag.
func dwarfTagString(tag enum.DwarfTag) string {
	s := tag.String()
	if strings.HasPrefix(s, "DwarfTag(") && strings.HasSuffix(s, ")") {
		return s[len("DwarfTag(") : len(s)-len(")")]
	}
	return s
}

// quote returns s as a double-quoted string literal.
func quote(s string) string {
	return enc.Quote([]byte(s))
}
