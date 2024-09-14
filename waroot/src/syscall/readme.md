# syscall 包规范

syscall 包对应不同的目标平台, 每个子包负责导入目标平台最小的 API.同时可以提供基于汇编胶水代码封装后更为友好一点的API(不宜过度), 子包不得依赖其他包.

## 目录结构

比如 `syscall/js` 对应 js 目标平台, 目录结构如下:

```
$ tree js
js
├── readme.md
├── api.wa
├── api_example_test.wa
├── z_abi.wa
└── z_abi.wat.ws
```

- readme.md: 包文档
- api.wa: 对外提供凹语言风格的 API 包装
- api_example_test.wa: 用法示例
- z_abi.wa: 导入的宿主函数或汇编实现的函数
- z_abi.wat.ws: 汇编语言代码


## 导入目标平台函数

只能导入目标平台对应的宿主函数, 例如:

```
#wa:import syscall_js proc_exit
func __import__proc_exit(code: i32)
```

导入名字以 `__import__` 为前缀.

## 汇编语言实现内部函数

只能在当前包对应的名字空间定义汇编函数, 例如:

```
#wa:linkname $syscall/js.__linkname__string_to_ptr
func __linkname__string_to_ptr(s: string) => i32
```

汇编函数名字以 `__linkname__` 为前缀.
