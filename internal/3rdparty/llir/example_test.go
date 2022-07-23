package llir_test

import (
	"fmt"

	"github.com/wa-lang/wa/internal/3rdparty/llir"
	constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

func Example() {
	// This example produces LLVM IR code equivalent to the following C code,
	// which implements a pseudo-random number generator.
	//
	//    int abs(int x);
	//
	//    int seed = 0;
	//
	//    // ref: https://en.wikipedia.org/wiki/Linear_congruential_generator
	//    //    a = 0x15A4E35
	//    //    c = 1
	//    int rand(void) {
	//       seed = seed*0x15A4E35 + 1;
	//       return abs(seed);
	//    }

	// Create convenience types and constants.
	i32 := types.I32
	zero := constant.NewInt(i32, 0)
	a := constant.NewInt(i32, 0x15A4E35) // multiplier of the PRNG.
	c := constant.NewInt(i32, 1)         // increment of the PRNG.

	// Create a new LLVM IR module.
	m := llir.NewModule()

	// Create an external function declaration and append it to the module.
	//
	//    int abs(int x);
	abs := m.NewFunc("abs", i32, llir.NewParam("x", i32))

	// Create a global variable definition and append it to the module.
	//
	//    int seed = 0;
	seed := m.NewGlobalDef("seed", zero)

	// Create a function definition and append it to the module.
	//
	//    int rand(void) { ... }
	rand := m.NewFunc("rand", i32)

	// Create an unnamed entry basic block and append it to the `rand` function.
	entry := rand.NewBlock("")

	// Create instructions and append them to the entry basic block.
	tmp1 := entry.NewLoad(types.I32, seed)
	tmp2 := entry.NewMul(tmp1, a)
	tmp3 := entry.NewAdd(tmp2, c)
	entry.NewStore(tmp3, seed)
	tmp4 := entry.NewCall(abs, tmp3)
	entry.NewRet(tmp4)

	// Print the LLVM IR assembly of the module.
	fmt.Println(m)

	// Output:
	//
	// @seed = global i32 0
	//
	// declare i32 @abs(i32 %x)
	//
	// define i32 @rand() {
	// 0:
	// 	%1 = load i32, i32* @seed
	// 	%2 = mul i32 %1, 22695477
	// 	%3 = add i32 %2, 1
	// 	store i32 %3, i32* @seed
	// 	%4 = call i32 @abs(i32 %3)
	// 	ret i32 %4
	// }
}
