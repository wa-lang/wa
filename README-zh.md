<div align="center">
<h1>🇨🇳 凹语言™</h1>

[主页](https://wa-lang.org) | [Playground](https://wa-lang.org/playground) | [目标](https://wa-lang.org/goals.html) | [路线](https://wa-lang.org/smalltalk/st0002.html) | [社区](https://wa-lang.org/community) | [日志](https://wa-lang.org/changelog.html) | [论坛](https://github.com/wa-lang/wa/discussions)

凹语言™（凹读音“Wa”）是 针对 WASM 平台设计的通用编程语言，同时支持 Linux、macOS 和 Windows 等主流操作系统和 Chrome 等浏览器环境，同时也支持作为独立Shell脚本和被嵌入脚本模式执行。

![](docs/images/logo/logo-animate1.svg)

- 主页: [https://wa-lang.org](https://wa-lang.org)
- 仓库: [https://gitee.com/wa-lang/wa](https://gitee.com/wa-lang/wa)
- 开发工具: [Playground](https://wa-lang.org/playground), [VSCode 插件](https://marketplace.visualstudio.com/items?itemName=xxxDeveloper.vscode-wa), [Fleet 插件](https://github.com/wa-lang/fleet-wa), [Vim 插件](https://github.com/wa-lang/vim-wa)
- 开发组: [柴树杉(chai2010)](https://github.com/chai2010)、[丁尔男(Ending)](https://github.com/3dgen)、[史斌(Benshi)](https://github.com/benshi001)、[扈梦明(xxxDeveloper)](https://github.com/xxxDeveloper)、[刘云峰(leaftree)](https://github.com/leaftree)、[宋汝阳(ShiinaOrez)](https://github.com/ShiinaOrez)

> 说明: 凹语言的主仓库在 https://gitee.com/wa-lang/wa. 在 Github 同时有一个镜像仓库 https://github.com/wa-lang/wa. 凹语言代码除非特别声明，均以 AGPL-v3 开源协议授权, 具体可以参考 LICENSE 文件.

## Playground 在线预览

[https://wa-lang.org/playground](https://wa-lang.org/playground)

![[![](https://wa-lang.org/smalltalk/images/st0011-01.png)](https://wa-lang.org/playground)](https://wa-lang.org/st0011-03.png)


## 本地安装和测试:

1. `go install wa-lang.org/wa@latest`
2. `wa init -name=_examples/hi`
3. `wa run _examples/hi`

> 项目尚处于原型开源阶段，如果有共建和PR需求请参考 [如何贡献代码](https://wa-lang.org/community/contribute.html)。我们不再接受针对第三方依赖库修改的 PR。

## 例子: 凹语言

打印字符和调用函数：

```wa
import "fmt"

fn main {
	println("你好，凹语言！")
	println(add(40, 2))

	fmt.Println(1+1)
}

fn add(a: i32, b: i32) => i32 {
	return a+b
}
```

运行并输出结果:

```
$ go run main.go hello.wa 
你好，凹语言！
42
2
```

## 例子: 打印素数

打印 30 以内的素数:

```
# 版权 @2021 凹语言™ 作者。保留所有权利。

fn main {
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

凹语言本身也可以像 Lua 语言被嵌入 Go 宿主语言环境执行:

```
package main

import (
	"fmt"
	"wa-lang.org/wa/api"
)

func main() {
	output, err := api.RunCode(api.DefaultConfig(), "hello.wa", code)
	fmt.Print(string(output), err)
}
```

注：作为脚本执行目前只支持本地环境。
