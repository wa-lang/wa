<div align="center">
<h1>🇨🇳 凹语言™ The Wa Programming Language</h1>

[主页](https://wa-lang.org) | [Playground](https://wa-lang.org/playground) | [目标](https://wa-lang.org/goals.html) | [路线](https://wa-lang.org/smalltalk/st0002.html) | [社区](https://wa-lang.org/community) | [日志](https://wa-lang.org/changelog.html) | [论坛](https://github.com/wa-lang/wa/discussions)

</div>
<div align="center">

[![Build Status](https://github.com/wa-lang/wa/actions/workflows/wa.yml/badge.svg)](https://github.com/wa-lang/wa/actions/workflows/wa.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/wa-lang/wa)](https://goreportcard.com/report/github.com/wa-lang/wa)
[![Coverage Status](https://coveralls.io/repos/github/wa-lang/wa/badge.svg)](https://coveralls.io/github/wa-lang/wa)
[![GitHub release](https://img.shields.io/github/v/tag/wa-lang/wa.svg?label=release)](https://github.com/wa-lang/wa/releases)

</div>

凹语言™（凹读音“Wa”）是 [柴树杉](https://github.com/chai2010)、[丁尔男](https://github.com/3dgen)、[史斌](https://github.com/benshi001) 针对 WASM 平台设计的的通用编程语言，支持 Linux、macOS 和 Windows 等主流操作系统和 Chrome 等浏览器环境，同时也支持作为独立Shell脚本和被嵌入脚本模式执行。

![](docs/images/logo/logo-animate1-blue.svg)

- 主页：[https://wa-lang.org](https://wa-lang.org)
- 扩展工具：[Playground](https://wa-lang.org/playground), [VSCode 插件](https://marketplace.visualstudio.com/items?itemName=xxxDeveloper.vscode-wa)
- 开发组成员：[柴树杉](https://github.com/chai2010)、[丁尔男](https://github.com/3dgen)、[史斌](https://github.com/benshi001)、[扈梦明](https://github.com/xxxDeveloper)、[刘云峰](https://github.com/leaftree)

安装和测试:

1. `go install github.com/wa-lang/wa@latest`
2. `wa init -name=_examples/hi`
3. `wa run _examples/hi`

> 项目尚处于原型开源阶段，如果有共建和PR需求请 [入群交流](https://wa-lang.org/community/index.html)。

> [VS Code 插件支持](https://marketplace.visualstudio.com/items?itemName=xxxDeveloper.vscode-wa)

## 设计目标

- 披着 Go 语法外衣的 C 语言；
- 凹语言™源码文件后缀为 `.wa`；
- 凹语言™编译器兼容 WaGo 语法。WaGo 是 Go 真子集。使用 WaGo 语法的源码文件后缀为 `.wa.go`。凹语法与 WaGo 语法在 AST 层面一致；
- 凹语言™支持中文/英文双语关键字，即任一关键字均有中文及英文版，二者在语法层面等价。

更多细节请参考 [凹语言™项目目标](docs/goals.md)

## 例子: 你好, 凹语言

打印字符和调用函数：

```
fn main() {
	print('凹')
	print('语')
	print('言')
	print('\n')

	println(add(40, 2))
}

fn add(a: i32, b: i32) => i32 {
	return a+b
}
```

运行并输出结果:

```
$ go run main.go hello.wa 
凹语言
42
```

## 例子: 打印素数

打印 30 以内的素数：

```
# 版权 @2021 凹语言™ 作者。保留所有权利。

fn main() {
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

运行并输出结果:

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

更多例子 [_examples](_examples)

## 作为脚本执行

凹语言本身也可以像 Lua 语言被嵌入 Go 宿主语言环境执行：

```go
package main

import (
	"fmt"
	"github.com/wa-lang/wa/api"
)

func main() {
	output, err := api.RunCode("hello.wa", "fn main() { println(40+2) }")
	fmt.Print(string(output), err)
}
```

注：作为脚本执行目前只支持本地环境。

## 版权

版权 @2019-2022 凹语言™ 作者。保留所有权利。

