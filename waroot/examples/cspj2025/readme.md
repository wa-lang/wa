# CSP-J 环境说明

因为涉及读写stdin或命令行参数，需要依赖以下的宿主函数：

```
#wa:import syscall_js get_stdin_size
func get_stdin_size => i32

#wa:import syscall_js get_stdin_data
func get_stdin_data(ptr: i32)

#wa:linkname $wa.runtime.slice_to_ptr
func byte_slice_to_ptr(t: []byte) => i32
```

在内置的`syscall_js`环境已经包含了以上的函数。如果是自定义的wasm环境需要手动配置。

Windows系统标准输入有两种方式：

- `echo 290es1q0 | prog.exe` 将 "290es1q0" 字符串作为 stdin 数据
- `prog.exe < number.in` 将 number.in 文件内容作为 stdin 数据

如果没有指定标准输入，同时程序里有执行了从 stdin 读取的操作，可能会发生阻塞。可以通过输入按 Ctrl + Z（Linux下按 Ctrl + D），然后按 Enter 方式模拟 stdin 文件结束。

## BUG

中文版底层尚有问题，正在修复中。
