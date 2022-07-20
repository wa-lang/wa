# 凹语言™

凹语言™（凹读音“Wa”）是[柴树杉](https://github.com/chai2010)、[丁尔男](https://github.com/3dgen)和[史斌](https://github.com/benshi001)设计的实验性编程语言。

```
+---+    +---+
| o |    | o |
|   +----+   |
|            |
|     Wa     |
|            |
+------------+
```

## 设计哲学

- 披着 Go 和 Rust 语法外衣的 C++ 语言
- 凹语言™源码文件后缀为 `.wa`
- 支持 `*.wa` 和 `*.wa.go` 两种类型的源码，分别采用凹语言™语法和 uGo 语法，二者在 AST 层面统一
- 凹语言支持英文和中文关键字，二者在 AST 层面统一，格式化时可指定输出的类型
- 凹语言的 AST 支持 C/LLVM/WASM 等后端

## 例子: 打印素数

[./_examples/hello/hello.wa.go](./_examples/hello/hello.wa.go) 打印 30 以内的素数：

```
// 版权 @2021 凹语言™ 作者。保留所有权利。

fn main() {
	for n := 2; n <= 30; n = n + 1 {
		let isPrime int = 1
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

## 版权

版权 @2019 凹语言™ 作者。保留所有权利。
