package llutil

import (
	"strings"

	"github.com/wa-lang/wa/internal/3rdparty/llir"
)

// Comment is an LLVM IR comment represented as a pseudo-instruction. Comment
// implements llir.Instruction.
type Comment struct {
	// Comment text; may contain multiple lines.
	Text string

	// embed llir.Instruction to satisfy the llir.Instruction interface.
	llir.Instruction
}

// LLString returns the LLVM syntax representation of the value.
func (c *Comment) LLString() string {
	// handle multi-line comments.
	text := strings.ReplaceAll(c.Text, "\n", "; ")
	return "; " + text
}

// NewComment returns a new LLVM IR comment represented as a pseudo-instruction.
// Text may contain multiple lines.
func NewComment(text string) *Comment {
	return &Comment{
		Text: text,
	}
}
