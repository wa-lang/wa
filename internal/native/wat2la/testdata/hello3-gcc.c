// 1. 参数在栈上, 如何对齐
// 2. 浮点数参数超出时, 如何利用整数寄存器

#include <stdint.h>

struct ret_t {
    int32_t v0;
    int64_t v1;
    int32_t v2;
    int64_t v3;
};

struct ret_t add10(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int32_t a8, int64_t a9) {
    struct ret_t v;
    v.v3 = a8+a9;
    return v;
}

int main() {
    struct ret_t x = add10(0, 1, 2, 3, 4, 5, 6, 7, 8, 9);
    return 0;
}
