# 凹语言和Go语言 Fib 递归版本基准对比

凹语言的优点可以通过一个简单的 Fibonacci 示例来说明。下面是凹语言和Go两种语言中实现的 fib 函数。

凹语言代码:

```wa
func init {
	println(fib(46))
}

func fib(n: int) => int {
	aux: func(n, acc1, acc2: int) => int
	aux = func(n, acc1, acc2: int) => int {
		switch n {
		case 0:
			return acc1
		case 1:
			return acc2
		default:
			return aux(n-1, acc2, acc1+acc2)
		}
	}
	return aux(n, 0, 1)
}
```

Go 代码:

```go
package main

func main() {
	println(fib(46))
}
func fib(n int) int {
	var aux func(n, acc1, acc2 int) int
	aux = func(n, acc1, acc2 int) int {
		switch n {
		case 0:
			return acc1
		case 1:
			return acc2
		default:
			return aux(n-1, acc2, acc1+acc2)
		}
	}
	return aux(n, 0, 1)
}
```

执行 `make` 输出结果如下：

```
$ make
wa -v
Wa version v0.17.0

go version
go version go1.21.0 darwin/amd64

wasmer -V
wasmer 4.3.7

wa build -optimize -target=wasi -output=fib_wa.wasm fib_wa.wa
GOOS=wasip1 GOARCH=wasm go build -o fib_go.wasm

du -sh fib_*.wasm
1.3M    fib_go.wasm
 12K    fib_wa.wasm

time wasmer fib_wa.wasm
1836311903
        0.04 real         0.03 user         0.01 sys
time wasmer fib_go.wasm
1836311903
        0.08 real         0.02 user         0.02 sys
```

简单总结:

- 凹语言是 v0.17.0 版本, Go 是 1.21.0 版本
- 凹语言输出的wasm体积为 12KB, Go 语言输出 1.3MB 大小的 wasm, 凹语言是Go的1/10大小
- 凹语言执行时间 0.12, Go 的执行时间是 0.26, 凹语言是Go的1/2执行时间

