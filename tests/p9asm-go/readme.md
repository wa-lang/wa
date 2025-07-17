
## Windows 10 环境测试

准备 WSL 测试环境：

- 下载 `ubuntu-24.04.2-wsl-amd64.wsl`: https://ubuntu.com/desktop/wsl
- 安装 `wsl --import Ubuntu-24.04 C:\WSL\Ubuntu24 ubuntu-24.04.2-wsl-amd64.wsl`
- 启动 `wsl -d Ubuntu-24.04`
- 执行 `apt-get update`
- 执行 `apt-get install golang`

编译汇编程序并验证结果：

- 执行 `go tool asm -p main -o main_linux_amd64.o main_linux_amd64.s`, Go1.11+ 需要 `-p main` 指定包路径
- 执行 `go tool link -H linux -o a.out main_linux_amd64.o`
- 执行 `./a.out`
- 验证返回值 `echo $?`
