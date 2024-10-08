# wat2c 翻译后C代码性能对比

注意，这个fib时递归实现，执行性能比较低。

`fib(46)` 递归实现, `-O1`优化：

```
$ make wat2c
go build -o fib_go_native.exe fib_go_native.go
clang -O1 -o fib_c_native.exe _fib_c_native.c
wa wat2c -o fib_wat2c_native.c fib_wat.txt && clang -O1 -o fib_wat2c_native.exe fib_wat2c_main.c
time ./fib_go_native.exe
fib(46) = 1836311903
       10.28 real         8.63 user         0.15 sys
time ./fib_c_native.exe
fib(46) = 1836311903
        5.39 real         5.03 user         0.05 sys
time ./fib_wat2c_native.exe
fib(46) = 1836311903
        6.38 real         5.23 user         0.06 sys
```

wat转译到C代码的执行性能接近本地C版本, 比Go快20-30%.
