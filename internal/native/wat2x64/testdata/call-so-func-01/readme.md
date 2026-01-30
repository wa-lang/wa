# 调用动态库

```
# 使用 nm 查看符号偏移
nm libhello.so | grep " T add"

# 或者使用 readelf（更准确）
readelf -s libhello.so | grep add
```

崩溃可以用 strace 跟着问题：

```
$ strace ./a.out
```