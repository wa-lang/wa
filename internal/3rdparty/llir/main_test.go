package llir_test

import (
	"fmt"

	ir "github.com/wa-lang/wa/internal/3rdparty/llir"
	constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

func Example_main() {
	// This example produces LLVM IR code equivalent to the following C code:
	//
	//    int main() {
	//       int a = 32;
	//       int b = 16;
	//       return a + b;
	//    }
	//
	// Read: https://blog.felixangell.com/an-introduction-to-llvm-in-go for inspiration.

	// Create convenience types.
	i32 := types.I32

	// Create a new LLVM IR module.
	m := ir.NewModule()
	// int main() { ... }
	main := m.NewFunc("main", i32)

	// Create an unnamed entry basic block and append it to the `main` function.
	entry := main.NewBlock("")
	// Create instructions and append them to the entry basic block.

	// %a = alloca i32
	a := entry.NewAlloca(i32)
	a.SetName("a")

	// %b = alloca i32
	b := entry.NewAlloca(i32)
	b.SetName("b")

	// store i32 32, i32* %a
	entry.NewStore(constant.NewInt(i32, 32), a)

	// store i32 16, i32* %b
	entry.NewStore(constant.NewInt(i32, 16), b)

	// %1 = load i32, i32* %a
	tmpA := entry.NewLoad(types.I32, a)

	// %2 = load i32, i32* %b
	tmpB := entry.NewLoad(types.I32, b)

	// %3 = add nsw i32 %1, %2
	tmpC := entry.NewAdd(tmpA, tmpB)

	// ret i32 %3
	entry.NewRet(tmpC)

	// Print the LLVM IR assembly of the module.
	fmt.Println(m)

	// Output:
	//
	// define i32 @main() {
	// 0:
	// 	%a = alloca i32
	// 	%b = alloca i32
	// 	store i32 32, i32* %a
	// 	store i32 16, i32* %b
	// 	%1 = load i32, i32* %a
	// 	%2 = load i32, i32* %b
	// 	%3 = add i32 %1, %2
	// 	ret i32 %3
	// }
}
