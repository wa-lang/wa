//go:build ignore

package main

import "fmt"

func main() {
	fmt.Printf("fib(%d) = %d\n", 46, fib(46))
}
func fib(n int) int {
	if n <= 2 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}
