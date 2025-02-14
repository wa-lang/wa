# syscall/wasm4

WASM4 游戏平台: https://wasm4.org/

Wasm4 是一个使用 WebAssembly 构建复古风格游戏的框架。它提供了一个 160 x 160 像素、64K 内存的游戏主机。通过使用 WebAssembly 技术使得开发的游戏能够在所有网页浏览器和一些低端设备上运行。随着凹语言支持Wasm4平台，现在开发者也能使用凹语言轻松开发Wasm4游戏。

## 快速入门

使用 `wa` 命令生存一个Wasm4示例工程:

```
$ wa init -wasm4
$ tree hello
hello
├── README.md
├── src
│   └── main.wa
└── wa.mod
```

在 hello 目录生成一个 Wasm4 版本的你好世界例子。其中 main.wa 代码如下：

```wa
import "syscall/wasm4"
```

首先是导入`syscall/wasm4`包，然后定义Update函数：

```wa
global smiley = [8]byte{...}

#wa:export update
func Update {
	wasm4.SetDrawColors(2, 0, 0, 0)
	wasm4.Text("Hello from Wa!", 10, 10)

	gamepad := wasm4.GetGamePad1()
	if gamepad&wasm4.BUTTON_1 != 0 {
		wasm4.SetDrawColors(4, 0, 0, 0)
	}

	wasm4.Blit(smiley[:], 76, 76, 8, 8, wasm4.BLIT_1BPP)
	wasm4.Text("Press X to blink", 16, 90)
}
```

首先是调用`wasm4.SetDrawColors`设置绘制颜色，然后调用`wasm4.Text`在屏幕的指定坐标绘制文字。然后根据`wasm4.GetGamePad1()`获得游戏按键状态，并有条件调整绘制颜色。最后`wasm4.Blit()`调用绘制一个笑脸精灵。

进入hello目录编译和执行：

```
$ wa build -target=wasm4
$ w4 run output/hello.wasm
```
