<div align="center">
<h1>凹语言</h1>

[主页](https://wa-lang.org) | [Playground](https://wa-lang.org/playground) | [路线](https://wa-lang.org/smalltalk/st0001.html) | [社区](https://wa-lang.org/community) | [日志](https://wa-lang.org/guide/changelog.html)

</div>

凹语言（凹读音“Wā”）是针对 WebAssembly 设计的编程语言，目标是为高性能网页应用提供一门简洁、可靠、易用、强类型的编译型通用语言。凹语言的代码生成器及运行时为全自主研发（不依赖于LLVM等外部项目），实现了全链路自主可控。目前凹语言处于工程试用阶段。

![](docs/images/wa-chan/wa-chan-front-small-logo-animate1.svg)

- 主页: [https://wa-lang.org](https://wa-lang.org)
- 参考手册: [https://wa-lang.org/man/](https://wa-lang.org/man/)
- 仓库(GitCode): [https://gitcode.com/wa-lang/wa](https://gitcode.com/wa-lang/wa)
- 仓库(Gitee): [https://gitee.com/wa-lang/wa](https://gitee.com/wa-lang/wa)
- 仓库(Github): [https://github.com/wa-lang/wa](https://github.com/wa-lang/wa)
- Playground: [https://wa-lang.org/playground](https://wa-lang.org/playground)

> 说明: 凹语言编译器代码以 AGPL-v3 开源协议授权, 标准库以 MIT 协议授权，这意味着您使用凹语言开发的程序可以安全商用无需开源。若您希望在自己的项目中整合凹语言编译器的代码，而又不希望受 AGPL-v3 的传染性限制，您可以联系我们单独为您定制授权协议。

## 如何参与开发

项目尚处于原型开源阶段，如果有共建和PR需求请参考 [如何贡献代码](./wca/readme.md)。我们不再接受针对第三方依赖库修改的 PR。

向本仓库提交PR视同您认可并接受[凹语言贡献者协议](./wca/wca.md)，但在实际签署之前，您的PR不会被评审或接受。

> 特别注意：与 Issue 相比，发起 PR 更容易获得贡献点（贡献点可用于参加回馈活动，如：[首次凹语言贡献者回馈活动](https://wa-lang.org/smalltalk/st0078.html)）。当您在项目中找到问题发起 Issue后，不妨与我们联系，我们会帮助您将 Issue 转为 PR。

## Playground 在线预览

[https://wa-lang.org/playground](https://wa-lang.org/playground)

![](docs/images/playground-01.png)

## 贪吃蛇游戏

- [https://wa-lang.org/wa/snake/](https://wa-lang.org/wa/snake/)
- [https://wa-lang.org/smalltalk/st0018.html](https://wa-lang.org/smalltalk/st0018.html)

![](docs/images/snake-01.jpg)


## WASM4游戏

- Wasm4/Snake: https://wa-lang.org/wa/w4-snake/
  - 贪吃蛇源码 (英文): https://gitcode.com/wa-lang/wa/tree/master/waroot/examples/w4-snake/
  - 贪吃蛇源码 (中文): https://gitcode.com/wa-lang/wa/tree/master/waroot/examples/w4-snake-wz/
- Wasm4/2048: https://wa-lang.org/wa/w4-2048/

![](docs/images/wasm4-game-snake-2048.png)

- [Wasm4/Snake Code](waroot/examples/w4-snake)
- [Wasm4/2048 Code](waroot/examples/w4-2048)

## NES小霸王游戏机模拟器

- Play: [https://wa-lang.org/nes/](https://wa-lang.org/nes/)
- Code: [https://gitee.com/wa-lang/nes-wa](https://gitee.com/wa-lang/nes-wa)

![](docs/images/nes-01.png)

## WebGPU 模拟土星和小行星

- Play: [https://wa-lang.org/webgpu/](https://wa-lang.org/webgpu/)
- Code: [https://gitee.com/wa-lang/webgpu](https://gitee.com/wa-lang/webgpu)

![](docs/images/webgpu-01.png)

## P5 儿童编程

- https://wa-lang.org/smalltalk/st0037.html

![](docs/images/p5wa-01.png)

## Arduino Nano 33 开发板

- https://wa-lang.org/smalltalk/st0052.html

![](docs/images/arduino-nano33-01.png)

## 例子: 凹语言(中文版)

打印字符串：

```wz
注: 你好，世界！

引入 "书"

函数·主控:
	书·说("你好，凹语言中文版！")
完毕
```

运行并输出结果:

```
$ wa run hello.wz 
你好，凹语言中文版！
```

## 例子: 凹语言(英文版)

打印字符和调用函数：

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

<!--
```w2
引入 "格式化"

全局 年 = 2023

定义 主控:
	输出("你好, 凹语言!")
	输出(加法(40, 2), 年)

	格式化.输出(1+1)
完毕

定义 加法(甲, 乙: 整数) => 整数:
	返回 甲+乙
完毕
```
-->

运行并输出结果:

```
$ wa run hello.wa 
你好，凹语言！
42 2023
2
```

## 例子: 打印素数

打印 30 以内的素数:

```wa
// 版权 @2021 凹语言™ 作者。保留所有权利。

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

<!--
```w2
注: 版权 @2021 凹语言™ 作者。保留所有权利。

定义 主控:
	循环 甲 := 2; 甲 <= 30; 甲 = 甲 + 1:
		是素数 := 1
		循环 乙 := 2; 乙*乙 <= 甲; 乙 = 乙 + 1:
			如果 丙 := 甲 % 乙; 丙 == 0:
				是素数 = 0
			完毕
		完毕
		如果 是素数 != 0:
			输出(n)
		完毕
	完毕
完毕
```
-->

运行并输出结果:

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

## 例子：Chrome本地AI

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


更多例子 [waroot/examples](waroot/examples)

## 贡献者名单

|贡献者|贡献点|
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

贡献点变更记录见 [waroot/cplog](waroot/cplog) 目录。


## 联系我们
电子邮箱：dev@wa-lang.org

微信号：walang_dev
![](https://wa-lang.org/wechatgroup.jpg)
