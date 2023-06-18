# Docker + WebAssembly 3 分钟入门: 凹语言版

2022 年底，Docker 发布了对 WebAssembly 支持的预览版本，通过集成 WasmEdge 运行时支持WASM模块。Docker 运行时已经准备就绪，那么如何构建出 WASM 模块呢？目前支持 Wasm 的语言有很多，这里我们选择凹语言来构建 Wasm 镜像。Docker 官方博文：https://www.docker.com/blog/docker-wasm-technical-preview/

## 1. 凹语言到WASM模块

凹语言是针对 WebAssembly 设计的语言，也是国内第一个实现纯浏览器内编译、执行全链路的自研静态类型的编译型通用编程语言。这里我们再尝试通过凹语言来构造 Docker 的 Wasm 镜像。

凹语言是 Go 语言开发的编译器，因此需要本地先安装 Go1.17+ 版本环境。然后基于 dev-wasi 分支最新代码构造出 wa 语言编译器命令。或者通过以下命令安装：`go install wa-lang.org/wa@dev-wasi`，安装的命令默认在 `$HOME/go/bin` 目录。确保本地 `wa` 命令行可以使用，可以通过 `wa -v` 查看版本。

凹语言环境配置好之后，创建 hello.wa 文件:

```
// 版权 @2019 凹语言 作者。保留所有权利。

import "fmt"
import "runtime"

func main {
	println("你好，凹语言！", runtime.WAOS)
	println(add(40, 2))

	fmt.Println(1+1)
}

func add(a: i32, b: i32) => i32 {
	return a+b
}
```

然后在命令行运行程序：

```
$ wa run hello.wa
你好，凹语言！wasi
42
2
```

一切正常！

## 2. 构建 wasm 模块

通过 `wa build hello.wa` 命令生成 a.out.wasm 模块，大小约 3.6KB。然后通过 wasmer 执行：

```
$ wa build hello.wa
$ wasmer a.out.wasm
你好，凹语言！wasi
42
2
```

也可以通过 wabt 等辅助工具测试。输出结果说明一切正常。

## 3. 配置 Docker wasm 环境

完整的 Docker wasm 官方文档可以看这里：https://docs.docker.com/desktop/wasm/ 。按照好最新的 Docker 之后，按照 https://docs.docker.com/desktop/containerd/#enabling-the-containerd-image-store-feature 的提示（Settings页面的Experimental菜单），打开“Use containerd for pulling and storing images”特性。

## 4. 构建 Docker wasm 镜像

在当前目录创建 Dockerfile，内容如下：

```dockerfile
FROM scratch
ADD a.out.wasm /hello.wasm
ENTRYPOINT ["hello.wasm"]
```

执行以下命令创建 Docker wasm 镜像：

```
$ docker buildx build --platform wasi/wasm32 -t wa-lang/hello-world .
[+] Building 5.2s (8/8) FINISHED
 => [internal] load build definition from Dockerfile                   0.1s
 => => transferring dockerfile: 191B                                   0.0s
 => [internal] load .dockerignore                                      0.2s
 => => transferring context: 2B                                        0.0s
 => resolve image config for docker.io/docker/dockerfile:1             3.9s
 => [auth] docker/dockerfile:pull token for registry-1.docker.io       0.0s
 => CACHED docker-image://docker.io/docker/dockerfile:1@sha256:d2...   0.2s
 => => resolve docker.io/docker/dockerfile:1@sha256:d2d74ff22a0e4...   0.2s
 => [internal] load build context                                      0.1s
 => => transferring context: 3.71kB                                    0.0s
 => [1/1] ADD a.out.wasm /hello.wasm                                   0.1s
 => exporting to image                                                 0.3s
 => => exporting layers                                                0.1s
 => => exporting manifest sha256:47686ab02b26ea0ed51261e17a32bc07...   0.0s
 => => exporting config sha256:95cf2c26b4bcf8ae39c99060a0aa108198...   0.0s
 => => naming to docker.io/wa-lang/hello-world:latest                  0.0s
 => => unpacking to docker.io/wa-lang/hello-world:latest               0.1s
$
```

需要注意的是这里使用的是 `buildx` 子命令，并且输出的是 `wasi/wasm32` 目标平台。完成后可以通过 `docker image list` 查看新生成的镜像。

## 5. 执行 Docker wasm 镜像

通过以下命令执行 Docker wasm 镜像：

```
$ docker run --rm \
	--name=wasm-hello \
	--runtime=io.containerd.wasmedge.v1 \
	--platform=wasi/wasm32 \
	docker.io/wa-lang/hello-world:latest
你好，凹语言！wasi
42
2
```

首先要选择 `io.containerd.wasmedge.v1` 运行时，同时也u有指定 `wasi/wasm32` 平台类型。如果一切顺利就可以看到输出结果了。

## 6. 总结展望

Docker 的创始人曾经说过，如果 Wasm 技术早点出现那么就不会有 Docker 这个技术。而目前 Docker 对 wasm 的支持也说明了其本身的应用场景。其实 Wasm 虽然诞生于 Web 领域，但是在 Web 之外也有相当广泛的应用场景。不过目前主流编程语言都不是原生为 Wasm 设计的，我们希望通过新的凹语言为 Wasm 提供更好的支持和体验。同时希望未来凹语言可以通过 Wasm 服务更多的场景。
