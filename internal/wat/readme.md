# Wat 文本到 wasm 转化

目标不是解析完整的 wat 语法, 而是能满足凹语言输出的 wat 格式.

## Wat 格式的子集

- 函数指令不支持折叠
- 只支持行注释, 不支持多行块注释
- 每个指令一行, 单指令之间不会出现行注释
- 对象前出现的是关联注释, 其他注释全部丢弃
- 转义字符串扩展: '\n', '\r', '\t', '\\', '\"'