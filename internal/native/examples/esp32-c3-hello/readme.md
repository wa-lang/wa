# ESP32-C3 例子

Windows 通过 Putty 查看串口数据.

烧写自定义镜像后将不能使用官方的监控功能.

目前的例子是替代的boot程序

```
$ esptool.py --chip esp32c3  --port COM3  --baud 460800 write_flash -z 0x0000 hello_riscv32_zh.ws.esp32.bin
```

在boot中还需要对串口进行初始化。

