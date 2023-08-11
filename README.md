<div align="center">
<h1>The Wa Programming Language</h1>

[简体中文](https://github.com/wa-lang/wa/blob/master/README-zh.md) | [English](https://github.com/wa-lang/wa/blob/master/README.md) 


</div>
<div align="center">

[![Build Status](https://github.com/wa-lang/wa/actions/workflows/wa.yml/badge.svg)](https://github.com/wa-lang/wa/actions/workflows/wa.yml)
[![Coverage Status](https://coveralls.io/repos/github/wa-lang/wa/badge.svg)](https://coveralls.io/github/wa-lang/wa)
[![GitHub release](https://img.shields.io/github/v/tag/wa-lang/wa.svg?label=release)](https://github.com/wa-lang/wa/releases)
[![license](https://img.shields.io/github/license/wa-lang/wa.svg)](https://github.com/wa-lang/wa/blob/master/LICENSE)

</div>

Wa is a general-purpose programming language designed for developing robustness and maintainability WebAssembly software.
Instead of requiring complex toolchains to set up, you can simply go install it - or run it in a browser.

![](docs/images/logo/logo-animate1.svg)

- Home: [https://wa-lang.org](https://wa-lang.org)
- Manual(Chinese): [https://wa-lang.org/man/](https://wa-lang.org/man/)
- Github: [https://github.com/wa-lang/wa](https://github.com/wa-lang/wa)
- Playground: [https://wa-lang.org/playground](https://wa-lang.org/playground)

> Note: Our canonical Git repository is located at https://gitee.com/wa-lang/wa. There is a mirror of the repository at https://github.com/wa-lang/wa. Unless otherwise noted, the Wa source files are distributed under the AGPL-v3 license found in the LICENSE file.

## Playground

[https://wa-lang.org/playground](https://wa-lang.org/playground)

![](https://wa-lang.org/playground-01.png)

## Snake Game

- Play: [https://wa-lang.org/wa/snake/](https://wa-lang.org/wa/snake/)
- Code: [waroot/examples/snake/README-en.md](waroot/examples/snake/README-en.md)

![](https://wa-lang.org/st0018-03.jpg)

## Example: Print Wa

Print rune and call function：

```wa
import "fmt"

global year: i32 = 2023

func main {
	println("hello, Wa!")
	println(add(40, 2), year)

	fmt.Println(1+1)
}

func add(a: i32, b: i32) => i32 {
	return a+b
}
```

Execute the program:

```
$ wa run hello.wa 
hello, Wa!
42 2023
2
```

## Example: Print Prime

Print prime numbers up to 30:

```wa
func main {
	for n := 2; n <= 30; n = n + 1 {
		isPrime: int = 1
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
$ cd waroot && wa run examples/prime
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

## Example: Print Prime with Chinese syntax

Print prime numbers up to 30:

```wz
引于 "书"

【启】：
  // 输出30以内的素数
  从n=2，到n>30，有n++：
    设素=1
    从i=2，到i*i>n，有i++：
      设x=n%i
      若x==0则：
        素=0
      。
    。
    若素!=0则：
      书·曰：n
    。
  。
。
```

Output is the same as the previous example.

More examples [waroot/examples](waroot/examples)

## Contributors

|Contributor|Contribution points|
| --- | --- |
|柴树杉| 35000|
|丁尔男| 42500|
|史斌  | 29000|
|扈梦明| 14000|
|赵普明| 17000|
|宋汝阳|  2000|
|刘云峰|  1000|
|王湘南|  1000|
|王泽龙|  1000|
|吴烜  |  3000|
|刘斌  |  2500|
|尹贻浩|  2000|
