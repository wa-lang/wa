package wat2x64

import (
	"fmt"
	"io"
	"strings"
)

func (p *wat2X64Worker) gasComment(w io.Writer, text string) {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "#") {
		fmt.Fprintln(w, text)
	} else {
		fmt.Fprintln(w, "#", text)
	}
}

func (p *wat2X64Worker) gasIntelSyntax(w io.Writer) {
	fmt.Fprintln(w, ".intel_syntax noprefix")
}

func (p *wat2X64Worker) gasSectionDataStart(w io.Writer) {
	fmt.Fprintln(w, ".section .data")
}

func (p *wat2X64Worker) gasSectionTextStart(w io.Writer) {
	fmt.Fprintln(w, ".section .text")
}

func (p *wat2X64Worker) gasExtern(w io.Writer, name string) {
	fmt.Fprintln(w, ".extern", name)
}

func (p *wat2X64Worker) gasGlobal(w io.Writer, name string) {
	fmt.Fprintln(w, ".globl", name)
}

func (p *wat2X64Worker) gasDefI32(w io.Writer, name string, v int32) {
	fmt.Fprintf(w, "%s: .long %d\n", name, v)
}

func (p *wat2X64Worker) gasDefI64(w io.Writer, name string, v int64) {
	fmt.Fprintf(w, "%s: .quad %d\n", name, v)
}

func (p *wat2X64Worker) gasDefF32(w io.Writer, name string, v float32) {
	fmt.Fprintf(w, "%s: .float %f\n", name, v)
}

func (p *wat2X64Worker) gasDefF64(w io.Writer, name string, v float64) {
	fmt.Fprintf(w, "%s: .double %f\n", name, v)
}

func (p *wat2X64Worker) gasDefString(w io.Writer, name string, v string) {
	fmt.Fprintf(w, "%s: .asciz %q\n", name, v)
}

func (p *wat2X64Worker) gasDefConstInt(w io.Writer, name string, v int) {
	fmt.Fprintf(w, "%s = %d\n", name, v)
}

func (p *wat2X64Worker) gasFuncStart(w io.Writer, fnName string) {
	fmt.Fprintf(w, "%s:\n", fnName)
}

// 函数名字重定向
func (p *wat2X64Worker) gasSet(w io.Writer, src, dst string) {
	fmt.Fprintf(w, ".set %s, %s\n", src, dst)
}
