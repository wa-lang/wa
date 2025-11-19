# 测试流程

1. Windows环境安装`esp-idf-tools-setup-espressif-ide-3.1.0-with-esp-idf-5.3.1.exe`
2. 进入ESP-IDF命令行环境, 执行`esptool.py version`命令查看版本
3. 复制`$IDF_PATH/examples/get-started/hello_world`例子
4. 进入工程目录执行命令`idf.py set-target esp32c3`
5. 将`build`子目录添加到`.gitignore`忽略文件
6. 配置工具`idf.py menuconfig`(可先跳过)
7. 执行构建命令`idf.py build`
  - 生成 bootloader.bin (位于 build/bootloader/)
  - 生成 partition-table.bin (位于 build/partition_table/)
  - 生成 hello_world.bin (位于 build/)
  - 以上几个是要烧录的文件
8. 执行烧录命令`idf.py flash`
9. 执行验证运行`idf.py monitor`
10. 修改打印的字符串, 重新build并验证

## 几个烧录文件

| 文件名 | 烧录地址 (默认) | 主要作用 |
| ------ | -------------- | -------- |
| bootloader.bin | 0x0 或 0x1000 | 一级引导：芯片上电后首先运行的代码，负责初始化 SPI Flash、找到分区表，并加载 二级引导器。
| partition-table.bin | 0x8000 | Flash 磁盘分区表：定义 Flash 存储器的逻辑布局，告诉引导器和应用程序固件、配置数据、OTA 槽位等都存放在哪里
| hello_world.bin | 0x10000 | 主应用程序：实际代码（hello_world），是芯片运行的最终目标

