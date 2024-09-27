# syscall/js 包

该包提供了 js 目标平台定义的 API 函数, 一般面向浏览器环境.

## 例子

新建 `hello.wa` 文件:

```wa
import "syscall/js"

func main {
	js.PrintString("hello js\n")
}
```

- `wa run -target=js hello.wa`: 执行
- `wa build -target=js hello.wa`: 编译

