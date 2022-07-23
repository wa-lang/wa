// This example program analyses an LLVM IR module to produce a callgraph in
// Graphviz DOT format.
package llir_test

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/asm"
	"github.com/llir/llvm/ir"
)

func Example_callgraph() {
	// Parse LLVM IR assembly file.
	m, err := asm.ParseFile("testdata/eval.ll")
	if err != nil {
		panic(err)
	}
	// Produce callgraph of module.
	callgraph := genCallgraph(m)
	// Output callgraph in Graphviz DOT format.
	fmt.Println(callgraph)

	// Output:
	//
	// digraph {
	// 	"@add"
	// 	"@sub"
	// 	"@f"
	// 	"@f" -> "@printf"
	// 	"@main"
	// 	"@main" -> "@add"
	// 	"@main" -> "@sub"
	// 	"@main" -> "@f"
	// 	"@printf"
	// }
}

// genCallgraph returns the callgraph in Graphviz DOT format of the given LLVM
// IR module.
func genCallgraph(m *ir.Module) string {
	buf := &strings.Builder{}
	buf.WriteString("digraph {\n")
	// For each function of the module.
	for _, f := range m.Funcs {
		// Add caller node.
		caller := f.Ident()
		fmt.Fprintf(buf, "\t%q\n", caller)
		// For each basic block of the function.
		for _, block := range f.Blocks {
			// For each non-branching instruction of the basic block.
			for _, inst := range block.Insts {
				// Type switch on instruction to find call instructions.
				switch inst := inst.(type) {
				case *ir.InstCall:
					callee := inst.Callee.Ident()
					// Add edges from caller to callee.
					fmt.Fprintf(buf, "\t%q -> %q\n", caller, callee)
				}
			}
			// Terminator of basic block.
			switch term := block.Term.(type) {
			case *ir.TermRet:
				// do something.
				_ = term
			}
		}
	}
	buf.WriteString("}")
	return buf.String()
}
