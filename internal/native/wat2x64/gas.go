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

func (p *wat2X64Worker) gasCommentInFunc(w io.Writer, text string) {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "#") {
		fmt.Fprintln(w, "    "+text)
	} else {
		fmt.Fprintln(w, "    #", text)
	}
}

func (p *wat2X64Worker) gasIntelSyntax(w io.Writer) {
	fmt.Fprintln(w, ".intel_syntax noprefix")
}

func (p *wat2X64Worker) gasSectionDataStart(w io.Writer) {
	fmt.Fprintln(w, ".section .data")
	fmt.Fprintln(w, ".align 8")
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

func (p *wat2X64Worker) gasDefArray(w io.Writer, name string, elemSize, len int) {
	fmt.Fprintf(w, "%s: .fill %d, %d, 0\n", name, len, elemSize)
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

func (p *wat2X64Worker) gasFuncLabel(w io.Writer, labelName string) {
	fmt.Fprintf(w, "%s:\n", labelName)
}

// 函数名字重定向
func (p *wat2X64Worker) gasSet(w io.Writer, src, dst string) {
	fmt.Fprintf(w, ".set %s, %s\n", src, dst)
}
