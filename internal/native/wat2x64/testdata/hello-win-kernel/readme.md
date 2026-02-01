# Windows x64 裸机 API 调用原理

在不链接任何静态库（如 `msvcrt.lib`, `kernel32.lib`）的情况下，通过纯汇编实现对操作系统 API 的动态寻址与调用。

## 1. 函数寻址流程

程序要在没有任何“预置地址”的情况下运行，必须完成以下三个步骤：

### 第1步: 定位 DLL 基址 (通过 PEB)

Windows 每个进程在启动时，`GS` 寄存器会指向 **TEB (Thread Environment Block)**，而 `GS:[0x60]` 处存放着 **PEB (Process Environment Block)** 的指针。

1. **PEB -> Ldr (0x18)**: 获取加载器数据，包含了所有已加载模块的链表。
2. **Ldr -> InLoadOrderModuleList (0x10)**: 这是一个双向链表。
3. 遍历动态库依赖链表:
   - 第一个节点是 `EXE` 本身。
   - 第二个节点通常是 `ntdll.dll`。
   - 第三个节点通常是 `kernel32.dll`。
4. **DllBase (0x30)**: 获取该节点的基地址。

### 第2步: 解析导出表 (Export Table)

拿到 `kernel32.dll` 的基址后，我们需要像手动打开文件一样去解析它的 **PE 结构**。

1. **DOS Header**: 检查开头的 `MZ` 签名，获取 `e_lfanew` (偏移 `0x3c`)。
2. **NT Header**: 跳转到 PE 签名处。
3. **Data Directory (偏移 0x88)**: 导出表（Export Directory）的相对虚拟地址（RVA）。
4. **Export Directory 结构**:
   -  **AddressOfNames**: 函数名字符串数组。
   -  **AddressOfNameOrdinals**: 序号数组（作为名称和地址的索引桥梁）。
   -  **AddressOfFunctions**: 实际函数的地址数组。

### 第3步: 字符串比对与调用

由于没有链接器，需要手动执行 `strcmp`查找。

## 编译命令

编译时不应包含标准库：

```bash
# -nostdlib: 不连接标准库
# -Wl,-e_start: 指定入口点为 _start 标签
# -s: 去除符号表减小体积
gcc -nostdlib -s -o hello.exe hello.s -Wl,-e_start
```

