## Windows 10 环境测试

- 下载 `ubuntu-24.04.2-wsl-amd64.wsl`: https://ubuntu.com/desktop/wsl
- 安装 `wsl --import Ubuntu-24.04 C:\WSL\Ubuntu24 ubuntu-24.04.2-wsl-amd64.wsl`
- 启动 `wsl -d Ubuntu-24.04`, 进入 Linux 环境
- 安装 Go 语言
- 切换到仓库目录 `cd /mnt/.../wa` 
- 编译 `go build` 生成 `wa` 文件
- 执行 `./wa`

或者在 Windows 命令行环境执行:

- `wsl ./wa`
