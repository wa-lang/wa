<div align="center">
<h1>The Wa Programming Language</h1>

[简体中文](https://github.com/wa-lang/wa/blob/master/README-zh.md) | [English](https://github.com/wa-lang/wa/blob/master/README.md) 


</div>
<div align="center">

[![Build Status](https://github.com/wa-lang/wa/actions/workflows/wa.yml/badge.svg)](https://github.com/wa-lang/wa/actions/workflows/wa.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/wa-lang/wa)](https://goreportcard.com/report/github.com/wa-lang/wa)
[![Coverage Status](https://coveralls.io/repos/github/wa-lang/wa/badge.svg)](https://coveralls.io/github/wa-lang/wa)
[![GitHub release](https://img.shields.io/github/v/tag/wa-lang/wa.svg?label=release)](https://github.com/wa-lang/wa/releases)
[![Go Reference](https://pkg.go.dev/badge/wa-lang.org/wa.svg)](https://pkg.go.dev/wa-lang.org/wa)
[![license](https://img.shields.io/github/license/wa-lang/wa.svg)](https://github.com/wa-lang/wa/blob/master/LICENSE)

</div>

Wa is a general-purpose programming language designed for developing robustness and maintainability WebAssembly software.
Instead of requiring complex toolchains to set up, you can simply go install it - or run it in a browser.

![](docs/images/logo/logo-animate1.svg)

- Home: [https://wa-lang.org](https://wa-lang.org)
- Github: [https://github.com/wa-lang/wa](https://github.com/wa-lang/wa)
- Playground: [https://wa-lang.org/playground](https://wa-lang.org/playground)

> Note: Our canonical Git repository is located at https://gitee.com/wa-lang/wa. There is a mirror of the repository at https://github.com/wa-lang/wa. Unless otherwise noted, the Wa source files are distributed under the AGPL-v3 license found in the LICENSE file.

## Playground

[https://wa-lang.org/playground](https://wa-lang.org/playground)

![](https://wa-lang.org/playground-01.png)

## Snake Game

- Play: [https://wa-lang.org/wa/snake/](https://wa-lang.org/wa/snake/)
- Code: [_examples/snake/README-en.md](_examples/snake/README-en.md)

![](https://wa-lang.org/st0018-03.jpg)

## Install and Run:

Go >= 1.17

1. `go install wa-lang.org/wa@latest`
2. `wa init -name=_examples/hi`
3. `wa run _examples/hi`

> The Wa project is still in very early stage. If you want to submit PR, please read the [Contribution Guide(Chinese)](https://wa-lang.org/community/contribute.html). We do not accept PR only about 3rdparty changes.

## Example: Print Wa

Print rune and call function：

```wa
# Copyright @2019-2022 The Wa author. All rights reserved.

import "fmt"

func main {
	println("hello, Wa!")
	println(add(40, 2))

	fmt.Println(1+1)
}

func add(a: i32, b: i32) => i32 {
	return a+b
}
```

Execute the program:

```
$ go run main.go hello.wa 
hello, Wa!
42
2
```

## Example: Print Prime

Print prime numbers up to 30:

```
func main {
	for n := 2; n <= 30; n = n + 1 {
		var isPrime int = 1
		for i := 2; i*i <= n; i = i + 1 {
			if x := n % i; x == 0 {
				isPrime = 0
			}
		}
		if isPrime != 0 {
			println(n)
		}
	}
}
```

Execute the program:

```
$ go run main.go run _examples/prime
2
3
5
7
11
13
17
19
23
29
```

More examples [_examples](_examples)

