//go:build ignore

package main

func main() {
	println(fib(46))
}
func fib(n int) int {
	if n <= 2 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}
