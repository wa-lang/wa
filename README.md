<div align="center">
<h1>The Wa Programming Language</h1>

[简体中文](https://github.com/wa-lang/wa/blob/master/README-zh.md) | [English](https://github.com/wa-lang/wa/blob/master/README.md) 


</div>
<div align="center">

[![](https://img.shields.io/badge/Wa-WASM-green?logo=data%3Aimage%2Fsvg%2Bxml%3Bcharset%3Dutf-8%3Bbase64%2CPHN2ZyB3aWR0aD0iMzAwIiBoZWlnaHQ9IjMwMCIgdmlld0JveD0iMCAwIDMwMCAzMDAiIGZpbGw9Im5vbmUiCiAgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8cGF0aCBmaWxsLXJ1bGU9ImV2ZW5vZGQiIGNsaXAtcnVsZT0iZXZlbm9kZCIgZD0iTTAgMjBDMCA4Ljk1NDMgOC45NTQzIDAgMjAgMEg4MEM5MS4wNDYgMCAxMDAgOC45NTQzIDEwMCAyMFYyNFY4MFYxMDBIMjAwVjgwVjI0VjIwQzIwMCA4Ljk1NDMgMjA4Ljk1NCAwIDIyMCAwSDI4MEMyOTEuMDQ2IDAgMzAwIDguOTU0MyAzMDAgMjBWNDRWODBWMjgwQzMwMCAyOTEuMDQ2IDI5MS4wNDYgMzAwIDI4MCAzMDBIMjBDOC45NTQzIDMwMCAwIDI5MS4wNDYgMCAyODBWODBWNDRWMjBaIiBmaWxsPSIjMDBCNUFCIi8%2BCiAgPHBhdGggZD0iTTUwIDU1QzUyLjc2MTQgNTUgNTUgNTIuNzYxNCA1NSA1MEM1NSA0Ny4yMzg2IDUyLjc2MTQgNDUgNTAgNDVDNDcuMjM4NiA0NSA0NSA0Ny4yMzg2IDQ1IDUwQzQ1IDUyLjc2MTQgNDcuMjM4NiA1NSA1MCA1NVoiIGZpbGw9IndoaXRlIi8%2BCiAgPHBhdGggZD0iTTI1MCA1NUMyNTIuNzYxIDU1IDI1NSA1Mi43NjE0IDI1NSA1MEMyNTUgNDcuMjM4NiAyNTIuNzYxIDQ1IDI1MCA0NUMyNDcuMjM5IDQ1IDI0NSA0Ny4yMzg2IDI0NSA1MEMyNDUgNTIuNzYxNCAyNDcuMjM5IDU1IDI1MCA1NVoiIGZpbGw9IndoaXRlIi8%2BCiAgPHBhdGggZD0iTTE1MCAxODBMMTg0IDIxNEwyMTggMTgwTTE1MCAxODBMMTE2IDIxNEw4MiAxODAiIGZpbGw9Im5vbmUiIHN0cm9rZT0id2hpdGUiIHN0cm9rZS13aWR0aD0iOCIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIi8%2BCjwvc3ZnPgo%3D)](https://wa-lang.org)
[![Build Status](https://github.com/wa-lang/wa/actions/workflows/wa.yml/badge.svg)](https://github.com/wa-lang/wa/actions/workflows/wa.yml)
[![Coverage Status](https://coveralls.io/repos/github/wa-lang/wa/badge.svg)](https://coveralls.io/github/wa-lang/wa)
[![GitHub release](https://img.shields.io/github/v/tag/wa-lang/wa.svg?label=release)](https://github.com/wa-lang/wa/releases)
[![license](https://img.shields.io/github/license/wa-lang/wa.svg)](https://github.com/wa-lang/wa/blob/master/LICENSE)
[![CNCF Landscape](https://img.shields.io/badge/CNCF%20Landscape-5699C6)](https://landscape.cncf.io/?item=wasm--languages--wa-lang)

</div>

Wa (Chinese name "凹", which pronounced "Wa") is a general-purpose programming language designed for for WebAssembly. The goal is to provide a simple, reliable, easy-to-use, statically typed language for high-performance web applications. The code generator and runtime are fully independently developed (not dependent on external projects such as LLVM). Currently, Wa is in the engineering trial stage.

![](docs/images/wa-chan/wa-chan-front-small-logo-animate1.svg)

- Home: [https://wa-lang.github.io/](https://wa-lang.github.io/)
- Manual: [https://wa-lang.github.io/man/en/](https://wa-lang.github.io/man/en/)
- Github: [https://github.com/wa-lang/wa](https://github.com/wa-lang/wa)
- Playground: [https://wa-lang.org/playground](https://wa-lang.org/playground)

> Note: The Wa compiler source files are distributed under the AGPL-v3 license. The Wa standard library is licensed under the MIT license. The programs you develop using the Wa language can be safely used commercially without open source. If you want to integrate the Wa compiler code into your own project, but do not want to be restricted by the infectiousness of AGPL-v3, you can contact us to customize the license agreement for you.

## How to contribute
The project is still in the prototype stage. If you have co-construction or PR needs, please refer to [How to contribute](./wca/readme_en.md). We no longer accept PRs for modifications to third-party libraries.

> Note: Our canonical Git repository is located at https://gitcode.com/wa-lang/wa. There is a mirror of the repository at https://github.com/wa-lang/wa.

> Note: Submitting a PR to this repository is deemed to be your recognition and acceptance of the [Wa Contributor Agreement](./wca/wca.md)(Aka WCA), but your PRs will not be reviewed or accepted until WCA is actually signed.

## Playground

[https://wa-lang.org/playground](https://wa-lang.org/playground)

![](docs/images/playground-01.png)

## Snake Game

- Play: [https://wa-lang.org/wa/snake/](https://wa-lang.org/wa/snake/)
- Code: [waroot/examples/snake/README-en.md](waroot/examples/snake/README-en.md)

![](docs/images/snake-01.jpg)

## WASM4 Game

- Wasm4/Snake: https://wa-lang.org/wa/w4-snake/
- Wasm4/2048: https://wa-lang.org/wa/w4-2048/

![](docs/images/wasm4-game-snake-2048.png)

- [Wasm4/Snake Code](waroot/examples/w4-snake)
- [Wasm4/2048 Code](waroot/examples/w4-2048)

## NES emulator

- Play: [https://wa-lang.org/nes/](https://wa-lang.org/nes/)
- Code: [https://github.com/wa-lang/nes](https://github.com/wa-lang/nes)

![](docs/images/nes-01.png)

## WebGPU Demo

- Play: [https://wa-lang.org/webgpu/](https://wa-lang.org/webgpu/)
- Code: [https://github.com/wa-lang/webgpu](https://github.com/wa-lang/webgpu)

![](docs/images/webgpu-01.png)

## P5 for creative coding

- https://wa-lang.org/smalltalk/st0037.html

![](docs/images/p5wa-01.png)

## Arduino Nano 33

- https://wa-lang.org/smalltalk/st0052.html

![](docs/images/arduino-nano33-01.png)


## Example: Print (Chinese)

Print string and call function：

```wz
注: 你好，世界！

引入 "书"

函数 主控:
	书·说("你好，凹语言中文版！")
完毕
```

Execute the program:

```
$ wa run hello.wz 
你好，凹语言中文版！
```

## Example: Print (English)

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
$ cd waroot && wa run -target=wasi examples/prime
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

## Example: Chrome Native AI

Chrome builtin Gemini Nano Demo:

```wa
import "ai"

func main {
	ai.RequestSession(func(session: ai.Session){
		session.PromptAsync("Who are you?", func(res: string) {
			println(res)
		})
	})
}
```

More examples [waroot/examples](waroot/examples)

## Contributors

|Contributor|Contribution points|
| --- | --- |
|柴树杉| 99650|
|丁尔男| 104150|
|史斌  | 10000|
|扈梦明| 60000|
|赵普明| 10000|
|宋汝阳|  2000|
|刘云峰|  1000|
|王潇南|  1000|
|王泽龙|  1000|
|吴烜  |  3000|
|刘斌  |  2500|
|尹贻浩|  2000|
|安博超 | 3000|
|yuqiaoyu| 600|
|qstesiro| 200|
|small_broken_gong|100|
|tk103331|100|
|蔡兴|3000|
|王任义|1000|
|imalasong|2000|
|杨刚|4000|
|崔爽|2000|
|李瑾|20000|
|王委委|100|
|雪碧|100|

## Contact us

Email: `dev{AT}wa-lang.org`
